package stylegen

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// A components-layer rule that overrides compiled presentation with a PLAIN
// DECLARATION — `display: flex`, not `@apply flex` — cannot be judged by
// comparing class names: there is no class. It has to be judged on the CSS
// property it sets, against the properties the component's own utilities set.
//
// Nothing in Go knows what properties `size-4` or `has-[>svg]:px-2` expand to;
// only Tailwind does, and only against this repo's theme. So the utilities are
// compiled — once per check, one probe class per utility — and the property
// names are read straight out of Tailwind's output. Deriving them from a
// hand-written property table instead would be exactly the silent-false-
// negative machine this gate exists to replace.
//
// The same compile answers a second question the gate used to ignore: WHEN does
// the utility apply. `hover:bg-accent` and an unconditional `background-color`
// rule name the same property but never compete, so a property-only oracle
// reports a contest that cannot happen — a false positive, which during a
// migration wave is the failure mode that gets the gate disabled. The state is
// read the same way the property is: out of Tailwind's own output, as the
// selector Tailwind emits around the probe class plus whatever at-rules it
// wraps it in. Hand-maintaining a variant->selector table here would reintroduce
// exactly the drift the property table was avoided to prevent.
//
// The compile is LAZY. It happens only when a rule that actually targets a
// composed marker carries plain declarations, so a tree whose foundation rules
// all name unmigrated components never shells out at all.

// probeClassPrefix names the throwaway class each probed utility is applied to.
const probeClassPrefix = "gsxui-layercheck-probe-"

var probeClassPattern = regexp.MustCompile(`\.` + regexp.QuoteMeta(probeClassPrefix) + `(\d+)\b`)

// stateProperty is one CSS property together with the state the declaring rule
// applies in. Contest is decided on the pair: an unconditional
// `background-color` rule contests `bg-x` but not `hover:bg-x`, and a `:hover`
// rule contests `hover:bg-x` but not `bg-x`.
type stateProperty struct {
	state    string
	property string
}

// utilityPropertyResolver answers "in which states does this set of utilities
// declare which CSS properties?", compiling on first use and caching for the
// rest of the check.
type utilityPropertyResolver struct {
	root string
	sets map[string]utilitySets

	byUtility map[string]map[stateProperty]struct{}
	err       error
	done      bool
}

// statePropertiesOf returns the union of the (state, property) pairs declared by
// every utility in the set.
func (r *utilityPropertyResolver) statePropertiesOf(utilities []string) (map[stateProperty]struct{}, error) {
	if err := r.resolve(); err != nil {
		return nil, err
	}
	union := map[stateProperty]struct{}{}
	for _, utility := range utilities {
		for pair := range r.byUtility[utility] {
			union[pair] = struct{}{}
		}
	}
	return union, nil
}

func (r *utilityPropertyResolver) resolve() error {
	if r.done {
		return r.err
	}
	r.done = true
	seen := map[string]struct{}{}
	var all []string
	for _, set := range r.sets {
		for _, utility := range append(append([]string{}, set.own...), set.descendant...) {
			if _, dup := seen[utility]; dup {
				continue
			}
			seen[utility] = struct{}{}
			all = append(all, utility)
		}
	}
	sort.Strings(all)
	r.byUtility, r.err = compileUtilityProperties(r.root, all)
	return r.err
}

// compileUtilityProperties runs Tailwind over one probe class per utility and
// reads back the state and properties each one declares.
func compileUtilityProperties(root string, utilities []string) (map[string]map[stateProperty]struct{}, error) {
	byUtility := make(map[string]map[stateProperty]struct{}, len(utilities))
	for _, utility := range utilities {
		byUtility[utility] = map[stateProperty]struct{}{}
	}
	if len(utilities) == 0 {
		return byUtility, nil
	}

	binary := filepath.Join(root, "node_modules", ".bin", "tailwindcss")
	if _, err := os.Stat(binary); err != nil {
		return nil, fmt.Errorf(
			"the layer gate must resolve compiled utilities to CSS properties, which needs the "+
				"Tailwind CLI at %s — run `npm install`: %w", binary, err)
	}

	// The probe lives inside root so that `@import \"tailwindcss\"` resolves
	// through root/node_modules, and it imports foundation.css because that is
	// where this repo's @theme lives — without it `bg-primary` is unknown.
	dir, err := os.MkdirTemp(root, ".layercheck-probe-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	var probe strings.Builder
	probe.WriteString("@import \"tailwindcss\" source(none);\n")
	probe.WriteString("@import \"tw-animate-css\";\n")
	fmt.Fprintf(&probe, "@import %q;\n", filepath.ToSlash(filepath.Join(root, "assets", "css", "foundation.css")))
	for index, utility := range utilities {
		fmt.Fprintf(&probe, ".%s%d { @apply %s; }\n", probeClassPrefix, index, utility)
	}

	input := filepath.Join(dir, "probe.css")
	output := filepath.Join(dir, "probe.out.css")
	if err := os.WriteFile(input, []byte(probe.String()), 0o644); err != nil {
		return nil, err
	}
	command := exec.Command(binary, "-i", input, "-o", output)
	command.Dir = root
	if combined, err := command.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("compile component utilities with Tailwind: %w\n%s", err, combined)
	}
	compiled, err := os.ReadFile(output)
	if err != nil {
		return nil, err
	}

	found, err := probedProperties(output, compiled)
	if err != nil {
		return nil, err
	}
	for index, pairs := range found {
		if index < 0 || index >= len(utilities) {
			continue
		}
		for pair := range pairs {
			byUtility[utilities[index]][pair] = struct{}{}
		}
	}
	for _, utility := range utilities {
		if len(byUtility[utility]) == 0 {
			return nil, fmt.Errorf(
				"utility %q compiled to no declarations — the layer gate cannot tell which "+
					"properties it contests, so it would silently pass every plain-declaration "+
					"override of it", utility)
		}
	}
	return byUtility, nil
}

// probedProperties maps each probe index to the (state, property) pairs its
// compiled rules declare. Tailwind's output nests @supports INSIDE rulesets and
// wraps rulesets in @media, so attribution walks the enclosing block headers: a
// header naming a probe class contributes the state, and every enclosing at-rule
// header contributes to it too. A CSS grammar parser is the wrong tool here —
// tdewolff's rejects a nested at-rule inside a qualified rule — and only block
// structure, selectors and property names are needed, so this scans for exactly
// those.
//
// One utility legitimately emits several rules with different states (a `dark:`
// utility that also needs an unconditional custom-property fallback, say); each
// contributes its own pairs.
func probedProperties(filename string, src []byte) (map[int]map[stateProperty]struct{}, error) {
	out := map[int]map[stateProperty]struct{}{}
	type frame struct {
		header  string
		indices []int
	}
	var stack []frame

	record := func(segment []byte) {
		declaration := strings.TrimSpace(string(segment))
		if declaration == "" || strings.HasPrefix(declaration, "@") {
			return
		}
		colon := strings.IndexByte(declaration, ':')
		if colon < 0 {
			return
		}
		property := strings.ToLower(strings.TrimSpace(declaration[:colon]))
		if property == "" || strings.ContainsAny(property, " \t\n(){}") {
			return
		}
		var atRules []string
		for _, f := range stack {
			if strings.HasPrefix(strings.TrimSpace(f.header), "@") {
				atRules = append(atRules, f.header)
			}
		}
		for _, f := range stack {
			for _, index := range f.indices {
				state := probeSelectorState(atRules, f.header, index)
				if out[index] == nil {
					out[index] = map[stateProperty]struct{}{}
				}
				out[index][stateProperty{state: state, property: property}] = struct{}{}
			}
		}
	}

	segment := 0
	for i := 0; i < len(src); i++ {
		switch c := src[i]; c {
		case '\'', '"':
			for i++; i < len(src); i++ {
				if src[i] == '\\' {
					i++
					continue
				}
				if src[i] == c {
					break
				}
			}
		case '/':
			if i+1 < len(src) && src[i+1] == '*' {
				if end := bytes.Index(src[i+2:], []byte("*/")); end >= 0 {
					i += 2 + end + 1
				} else {
					i = len(src)
				}
			}
		case '{':
			header := string(src[segment:i])
			var indices []int
			for _, match := range probeClassPattern.FindAllSubmatch(src[segment:i], -1) {
				index, err := strconv.Atoi(string(match[1]))
				if err != nil {
					continue
				}
				indices = append(indices, index)
			}
			stack = append(stack, frame{header: header, indices: indices})
			segment = i + 1
		case '}':
			record(src[segment:i])
			if len(stack) == 0 {
				return nil, fmt.Errorf("%s: unbalanced block at offset %d", filename, i)
			}
			stack = stack[:len(stack)-1]
			segment = i + 1
		case ';':
			record(src[segment:i])
			segment = i + 1
		}
	}
	if len(stack) != 0 {
		return nil, fmt.Errorf("%s: %d unclosed block(s)", filename, len(stack))
	}
	return out, nil
}

// probeSelectorState reduces one compiled probe ruleset to its comparable state.
// Tailwind may emit a selector suffix (`.probe0:hover`), an ancestor form
// (`.dark .probe0`, or its `:is(.dark *)` encoding), an at-rule wrapper, or any
// combination; all three are captured so that like is compared with like.
func probeSelectorState(atRules []string, selector string, index int) string {
	token := fmt.Sprintf(".%s%d", probeClassPrefix, index)
	// A selector list on one probe rule is unusual, but if it happens each
	// complex selector could carry a different state. Only the ones naming this
	// probe count, and the narrowest reading is to keep them all.
	var states []string
	for _, complex := range splitSelectorList(selector) {
		ancestors, suffix, ok := subjectConditions(complex, func(compound string) (string, bool) {
			return removeToken(compound, token)
		})
		if !ok {
			continue
		}
		states = append(states, canonicalState(atRules, ancestors, suffix))
	}
	if len(states) == 0 {
		// The probe class appears only in a non-subject position — the utility
		// styles something other than the element it is on. Nothing an authored
		// rule on the marker can be compared against, so record an unmatchable
		// state rather than pretending it is unconditional.
		return "\x00unmatched"
	}
	sort.Strings(states)
	return strings.Join(states, ";")
}

// canonicalState renders the conditions under which a rule applies to its
// subject as one comparable string: the at-rules it is wrapped in, the ancestor
// conditions above it, and the subject compound's own trailing conditions.
//
// Two states contest only when these strings are equal. That is deliberately
// strict: a state this function cannot reduce to the same text on both sides —
// a `@media (width >= 40rem)` wrapper on the authored side against Tailwind's
// own `md:` output, say — compares unequal and is therefore NOT reported as a
// contest. Missing a real contest costs a silent dead rule; inventing one blocks
// legitimate work and gets the gate switched off, so the error is taken in that
// direction on purpose. TestCanonicalStateDoesNotContestAcrossUncomparableMedia
// pins the choice.
func canonicalState(atRules, ancestors []string, suffix string) string {
	var wrappers []string
	for _, at := range atRules {
		normalized := normalizeSpace(at)
		// Tailwind implements `hover:` as a `:hover` rule inside
		// `@media (hover: hover)`; an authored `:hover` rule carries no such
		// wrapper. Keeping it would make EVERY hover-versus-hover contest
		// invisible, which is the one variant the interaction rules in
		// assets/css/styles/default.css use most. It is dropped so that the
		// `:hover` suffix — which both sides do carry — decides.
		if strings.HasPrefix(normalized, "@media") && strings.Contains(normalized, "(hover: hover)") {
			continue
		}
		wrappers = append(wrappers, normalized)
	}
	sort.Strings(wrappers)
	scopes := make([]string, 0, len(ancestors))
	for _, ancestor := range ancestors {
		scopes = append(scopes, stripSpace(ancestor))
	}
	sort.Strings(scopes)
	return strings.Join(wrappers, "&&") + "|" + strings.Join(scopes, "&&") + "|" + stripSpace(suffix)
}

// subjectConditions splits a complex selector into the ancestor conditions above
// its subject and the conditions the subject compound carries itself, once
// `remove` has taken the identifying token — the probe class, or the slot marker
// attribute — out of that compound.
//
// It reports false when the identifying token is not in the subject compound,
// which means the selector styles a DESCENDANT of the identified element rather
// than the element itself. That case is handled by the caller, not here.
func subjectConditions(complex string, remove func(string) (string, bool)) (ancestors []string, suffix string, ok bool) {
	compounds := splitCompounds(complex)
	if len(compounds) == 0 {
		return nil, "", false
	}
	rest, ok := remove(compounds[len(compounds)-1])
	if !ok {
		return nil, "", false
	}
	ancestors = scopeAncestors(compounds[:len(compounds)-1])
	// Tailwind encodes an ancestor condition as a suffix: `dark:` compiles to
	// `.probe0:is(.dark *)`, not `.dark .probe0`. Lifting it back out is what
	// lets it compare equal to an authored `.dark [marker]` rule.
	rest, lifted := liftAncestorPseudos(rest)
	ancestors = append(ancestors, scopeAncestors(lifted)...)
	return ancestors, rest, true
}

// scopeAncestors keeps the ancestor compounds that can actually change whether a
// rule applies to a given element in a given state — those naming a class or an
// id, such as `.dark`. A pure scoping ancestor like
// `:where(body:not([data-theme-button-preview]))` narrows WHICH documents the
// rule reaches, not which state of the element it targets, and both sides of the
// comparison would otherwise never line up.
func scopeAncestors(compounds []string) []string {
	var out []string
	for _, compound := range compounds {
		if strings.ContainsAny(compound, ".#") {
			out = append(out, compound)
		}
	}
	return out
}

// liftAncestorPseudos pulls `:is(X *)` / `:where(X *)` suffixes — Tailwind's
// encoding of an ancestor condition — out of a compound, returning what is left
// and the ancestor selectors it found.
func liftAncestorPseudos(compound string) (string, []string) {
	var lifted []string
	out := compound
	for _, name := range []string{":is(", ":where("} {
		for {
			index := strings.Index(strings.ToLower(out), name)
			if index < 0 {
				break
			}
			open := index + len(name) - 1
			end := skipBalanced(out, open, '(', ')')
			args := out[open+1 : max(end-1, open+1)]
			trimmed := strings.TrimSpace(args)
			if !strings.HasSuffix(trimmed, "*") {
				break
			}
			ancestor := strings.TrimSpace(strings.TrimSuffix(trimmed, "*"))
			if ancestor == "" {
				break
			}
			lifted = append(lifted, ancestor)
			out = out[:index] + out[end:]
		}
	}
	return out, lifted
}

// removeToken deletes one whole selector token — `.gsxui-layercheck-probe-3` —
// from a compound, reporting whether it was there. The match is
// boundary-checked so that probe 3 is never found inside probe 31.
func removeToken(compound, token string) (string, bool) {
	for offset := 0; ; {
		index := strings.Index(compound[offset:], token)
		if index < 0 {
			return compound, false
		}
		index += offset
		after := index + len(token)
		if after >= len(compound) || !isAttributeNameByte(compound[after]) {
			return compound[:index] + compound[after:], true
		}
		offset = index + 1
	}
}

// removeMarkerAttribute deletes the `[data-gsxui-slot-…]` attribute selector
// naming the given marker from a compound, reporting whether it was there, then
// drops any `:where()`/`:is()` wrapper it emptied out. `:where([marker]):hover`
// has to reduce to `:hover`, not to `:where():hover`.
func removeMarkerAttribute(compound, marker string) (string, bool) {
	for i := 0; i < len(compound); i++ {
		if compound[i] != '[' {
			continue
		}
		end := skipBalanced(compound, i, '[', ']')
		if !containsAttributeName(compound[i:end], marker) {
			i = end - 1
			continue
		}
		// A marker attribute carrying a value — [data-gsxui-slot-x="y"] — is a
		// narrower target than the bare marker, but the marker is still what the
		// rule selects; the value is part of the state, so only the name is
		// removed and the rest stays in the suffix.
		inner := strings.TrimSpace(compound[i+1 : max(end-1, i+1)])
		remainder := strings.TrimSpace(strings.TrimPrefix(inner, marker))
		out := compound[:i]
		if remainder != "" {
			out += "[" + marker + remainder + "]"
		}
		out += compound[end:]
		return dropEmptyLogicalPseudos(out), true
	}
	return compound, false
}

// dropEmptyLogicalPseudos removes `:where()` and `:is()` left behind with no
// arguments once the marker attribute was taken out of them.
func dropEmptyLogicalPseudos(compound string) string {
	for _, empty := range []string{":where()", ":is()"} {
		compound = strings.ReplaceAll(compound, empty, "")
	}
	return compound
}

func normalizeSpace(s string) string { return strings.Join(strings.Fields(s), " ") }

func stripSpace(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case ' ', '\t', '\n', '\r':
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

package stylegen

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
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

// ownedProperties maps each CSS property a component's compiled utilities set to
// EVERY state in which they set it. The list is a DISJUNCTION: the property is
// governed by the utilities layer whenever any one of those states holds, so an
// authored rule is shadowed as soon as it implies any single entry.
type ownedProperties map[string][]statePredicate

// shadows reports whether the compiled utilities take over property in a state
// the authored predicate guarantees. This is the whole directional relation: the
// authored rule is dead when, WHENEVER IT APPLIES, a compiled utility applies
// too — so it is authored that must imply generated, never the other way round.
func (owned ownedProperties) shadows(property string, authored statePredicate) bool {
	return slices.ContainsFunc(owned[property], authored.implies)
}

func (owned ownedProperties) add(property string, state statePredicate) {
	for _, existing := range owned[property] {
		if existing.box == state.box && slices.Equal(existing.conditions, state.conditions) {
			return
		}
	}
	owned[property] = append(owned[property], state)
}

// utilityPropertyResolver answers "in which states does this set of utilities
// declare which CSS properties?", compiling on first use and caching for the
// rest of the check.
type utilityPropertyResolver struct {
	root string
	sets map[string]utilitySets

	byUtility map[string]ownedProperties
	err       error
	done      bool
}

// statePropertiesOf returns the union, over every utility in the set, of the
// states in which each CSS property is declared.
func (r *utilityPropertyResolver) statePropertiesOf(utilities []string) (ownedProperties, error) {
	if err := r.resolve(); err != nil {
		return nil, err
	}
	union := ownedProperties{}
	for _, utility := range utilities {
		for property, states := range r.byUtility[utility] {
			for _, state := range states {
				union.add(property, state)
			}
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
func compileUtilityProperties(root string, utilities []string) (map[string]ownedProperties, error) {
	byUtility := make(map[string]ownedProperties, len(utilities))
	for _, utility := range utilities {
		byUtility[utility] = ownedProperties{}
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
	for index, owned := range found {
		if index < 0 || index >= len(utilities) {
			continue
		}
		for property, states := range owned {
			for _, state := range states {
				byUtility[utilities[index]].add(property, state)
			}
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
func probedProperties(filename string, src []byte) (map[int]ownedProperties, error) {
	out := map[int]ownedProperties{}
	type frame struct {
		header  string
		indices []int
	}
	var stack []frame
	var recordErr error

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
				states, err := probeSelectorStates(atRules, f.header, index)
				if err != nil {
					if recordErr == nil {
						recordErr = fmt.Errorf("%s: compiled selector %q: %w", filename, normalizeSpace(f.header), err)
					}
					continue
				}
				if out[index] == nil {
					out[index] = ownedProperties{}
				}
				for _, state := range states {
					out[index].add(property, state)
				}
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
	if recordErr != nil {
		return nil, recordErr
	}
	return out, nil
}

// probeSelectorStates reduces one compiled probe ruleset to the states it
// applies in. Tailwind may emit a selector suffix (`.probe0:hover`), an ancestor
// form (`.dark .probe0`, or its `:is(.dark *)` encoding), an at-rule wrapper, or
// any combination; all three are captured so that like is compared with like.
//
// A selector LIST on one probe rule is a disjunction — `.p:hover, .p:focus`
// governs the property in either state — so each branch becomes its own
// predicate and the caller reports a contest when the authored rule implies ANY
// of them. Folding them into one comparable string, which is what this used to
// do, invented a compound state that holds in neither case.
func probeSelectorStates(atRules []string, selector string, index int) ([]statePredicate, error) {
	token := fmt.Sprintf(".%s%d", probeClassPrefix, index)
	var states []statePredicate
	for _, complex := range splitSelectorList(selector) {
		ancestors, suffix, ok := subjectConditions(complex, func(compound string) (string, bool) {
			return removeToken(compound, token)
		})
		if !ok {
			continue
		}
		state, err := buildPredicate(atRules, ancestors, suffix)
		if err != nil {
			return nil, err
		}
		states = append(states, state)
	}
	if len(states) == 0 {
		// The probe class appears only in a non-subject position — the utility
		// styles something other than the element it is on. Nothing an authored
		// rule on the marker can be compared against, so record an unmatchable
		// state rather than pretending it is unconditional.
		return []statePredicate{unmatchablePredicate()}, nil
	}
	return states, nil
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

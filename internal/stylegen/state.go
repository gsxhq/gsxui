package stylegen

import (
	"fmt"
	"sort"
	"strings"
)

// SHADOWING IS DIRECTIONAL.
//
// A components-layer rule is dead when the compiled utility beats it EVERY TIME
// the authored rule applies. That is implication, not equality:
//
//	authored :hover  vs generated unconditional -> DEAD (the utility also
//	                                               applies during hover)
//	authored unconditional vs generated :hover  -> ALIVE (the authored rule
//	                                               still governs the rest)
//
// Modelling the relation as string equality — which is what this gate did — got
// the first row wrong and passed a rule that can never take effect. So a state
// is modelled here as an explicit PREDICATE: the conjunction of every condition
// that must hold for the rule to apply to its subject, plus which box it
// targets. Implication is then containment of conjuncts, because a conjunction
// implies each of its conjuncts and nothing else.
//
// The model FAILS LOUDLY on any selector form it cannot decompose. The spec
// claims this hazard is a build failure; treating an un-modelled condition as
// "not contesting" would make that claim false, and quietly so.

// statePredicate is the condition under which a rule applies to its subject.
//
// conditions is a CONJUNCTION: every one must hold. box names the pseudo-element
// the rule paints, empty for the element itself. box is compared by EQUALITY,
// never by implication — `[marker]::after` and a utility on `[marker]` style two
// different boxes and can never shadow one another, in either direction.
type statePredicate struct {
	box        string
	conditions []string
}

// unmatchableBox marks a predicate nothing can imply. It is used for a compiled
// utility whose probe class appears only in a non-subject position: the utility
// styles something other than the element it is on, so no authored rule on the
// marker can be lined up against it.
const unmatchableBox = "\x00unmatched"

func (p statePredicate) String() string {
	if len(p.conditions) == 0 && p.box == "" {
		return "unconditional"
	}
	return strings.Join(p.conditions, "&&") + p.box
}

// implies reports whether this predicate guarantees other — whether, whenever a
// rule with this state applies, a rule with state other applies too.
//
// A conjunction implies each of its conjuncts, so p implies other exactly when
// every one of other's conditions is also one of p's. The unconditional
// predicate has no conditions, is therefore implied by everything, and implies
// only itself: it is the weakest condition, which is what makes an unconditional
// utility shadow a :hover rule while the reverse does not hold.
func (p statePredicate) implies(other statePredicate) bool {
	if p.box != other.box {
		return false
	}
	if p.box == unmatchableBox {
		return false
	}
	held := make(map[string]struct{}, len(p.conditions))
	for _, condition := range p.conditions {
		held[condition] = struct{}{}
	}
	for _, condition := range other.conditions {
		if _, ok := held[condition]; !ok {
			return false
		}
	}
	return true
}

// unmatchablePredicate is the state of a utility that cannot be lined up against
// any authored rule on the marker.
func unmatchablePredicate() statePredicate { return statePredicate{box: unmatchableBox} }

// buildPredicate assembles the state of one rule from the three independent
// places its conditions come from: the at-rules it is nested in, the ancestor
// conditions above its subject, and the subject compound's own conditions. All
// three are conjunctive — every one must hold — so they land in one condition
// set and implication sees through all of them alike. A `.dark`-scoped rule
// therefore implies the same rule unscoped, exactly as `:hover` implies
// unconditional.
func buildPredicate(atRules, ancestors []string, suffix string) (statePredicate, error) {
	var conditions []string
	for _, at := range atRules {
		normalized := normalizeSpace(at)
		// Tailwind implements `hover:` as a `:hover` rule inside
		// `@media (hover: hover)`; an authored `:hover` rule carries no such
		// wrapper. Keeping it would make EVERY hover-versus-hover contest
		// invisible, which is the one variant the interaction rules in
		// assets/css/styles/default.css use most. It is dropped so that the
		// `:hover` condition — which both sides do carry — decides.
		if strings.HasPrefix(normalized, "@media") && strings.Contains(normalized, "(hover: hover)") {
			continue
		}
		conditions = append(conditions, "@"+stripSpace(normalized))
	}
	for _, ancestor := range ancestors {
		conditions = append(conditions, "ancestor:"+stripSpace(ancestor))
	}
	own, box, err := compoundConditions(suffix)
	if err != nil {
		return statePredicate{}, err
	}
	conditions = append(conditions, own...)
	sort.Strings(conditions)
	conditions = dedupe(conditions)
	return statePredicate{box: box, conditions: conditions}, nil
}

func dedupe(sorted []string) []string {
	out := sorted[:0]
	for i, s := range sorted {
		if i == 0 || sorted[i-1] != s {
			out = append(out, s)
		}
	}
	return out
}

// compoundConditions decomposes one compound selector — what is left of a
// subject compound once the identifying marker or probe class has been taken out
// of it — into its conjuncts plus the pseudo-element box it targets.
//
// Every form it accepts is one it can compare; every form it does not accept is
// an error, named, rather than a silently empty condition set. An empty set
// means UNCONDITIONAL, the weakest predicate, which is implied by everything —
// so mis-reading an unsupported form as empty would report shadowing that does
// not happen. Erring is the only safe way to not know.
func compoundConditions(compound string) (conditions []string, box string, err error) {
	setBox := func(candidate string) error {
		if box != "" && box != candidate {
			return fmt.Errorf("two pseudo-elements (%s and %s) in one compound", box, candidate)
		}
		box = candidate
		return nil
	}
	for i := 0; i < len(compound); {
		switch c := compound[i]; {
		case c == ' ' || c == '\t' || c == '\n' || c == '\r':
			i++
		case c == '*':
			i++
		case c == '.' || c == '#':
			end := skipIdent(compound, i+1)
			if end == i+1 {
				return nil, "", fmt.Errorf("empty class or id name at offset %d", i)
			}
			conditions = append(conditions, compound[i:end])
			i = end
		case c == '[':
			end := skipBalanced(compound, i, '[', ']')
			if end > len(compound) || compound[end-1] != ']' {
				return nil, "", fmt.Errorf("unterminated attribute selector at offset %d", i)
			}
			condition, err := normalizeAttributeCondition(compound[i:end])
			if err != nil {
				return nil, "", err
			}
			conditions = append(conditions, condition)
			i = end
		case c == ':':
			nested, nestedBox, end, err := pseudoCondition(compound, i)
			if err != nil {
				return nil, "", err
			}
			conditions = append(conditions, nested...)
			if nestedBox != "" {
				if err := setBox(nestedBox); err != nil {
					return nil, "", err
				}
			}
			i = end
		case isIdentStartByte(c):
			end := skipIdent(compound, i)
			conditions = append(conditions, "type:"+strings.ToLower(compound[i:end]))
			i = end
		default:
			return nil, "", fmt.Errorf(
				"unsupported selector character %q at offset %d — the layer gate cannot decide "+
					"whether this condition is implied by a compiled utility's own", string(c), i)
		}
	}
	return conditions, box, nil
}

// legacyPseudoElements are the four pseudo-elements CSS2 spelled with a single
// colon. `:before` and `::before` are the same box, so they must reduce to the
// same text or an authored `:before` rule would compare against a compiled
// `::before` utility as a different target.
var legacyPseudoElements = map[string]bool{
	"before": true, "after": true, "first-line": true, "first-letter": true,
}

// statePseudoClasses is every argument-less pseudo-class the model understands
// as a CONDITION on the subject. It is a closed list on purpose: a pseudo-class
// that is not here is one whose meaning the gate has not been taught, and the
// gate says so instead of guessing.
//
// The list covers every form the authored stylesheets actually use — see
// TestSupportedStateFormsCoverTheAuthoredStylesheets, which re-derives the forms
// in assets/css and web/ and fails when one falls outside the model.
var statePseudoClasses = map[string]bool{
	"active": true, "any-link": true, "autofill": true, "checked": true,
	"default": true, "defined": true, "disabled": true, "empty": true,
	"enabled": true, "first-child": true, "first-of-type": true,
	"focus": true, "focus-visible": true, "focus-within": true,
	"fullscreen": true, "hover": true, "in-range": true, "indeterminate": true,
	"invalid": true, "last-child": true, "last-of-type": true, "link": true,
	"modal": true, "only-child": true, "only-of-type": true, "open": true,
	"optional": true, "out-of-range": true, "placeholder-shown": true,
	"popover-open": true, "read-only": true, "read-write": true,
	"required": true, "root": true, "target": true, "user-invalid": true,
	"user-valid": true, "valid": true, "visited": true,
}

// functionalStatePseudoClasses is every pseudo-class taking arguments that the
// model treats as ONE opaque condition. Opaque is sound: two of them compare
// equal only when they are textually the same condition, and implication over a
// conjunction never needs to look inside one. `:is()`/`:where()` are handled
// separately because they are the two that can be vacuous or can be flattened.
var functionalStatePseudoClasses = map[string]bool{
	"dir": true, "has": true, "lang": true, "not": true, "nth-child": true,
	"nth-last-child": true, "nth-last-of-type": true, "nth-of-type": true,
	"state": true,
}

// pseudoCondition consumes one pseudo-class or pseudo-element starting at the
// ':' at index, returning the conditions it contributes, the pseudo-element box
// it names if it is one, and the offset just past it.
func pseudoCondition(compound string, index int) (conditions []string, box string, end int, err error) {
	i := index + 1
	doubled := false
	if i < len(compound) && compound[i] == ':' {
		doubled = true
		i++
	}
	nameStart := i
	i = skipIdent(compound, i)
	name := strings.ToLower(compound[nameStart:i])
	if name == "" {
		return nil, "", 0, fmt.Errorf("empty pseudo name at offset %d", index)
	}
	args := ""
	hasArgs := false
	if i < len(compound) && compound[i] == '(' {
		close := skipBalanced(compound, i, '(', ')')
		if close > len(compound) || compound[close-1] != ')' {
			return nil, "", 0, fmt.Errorf("unterminated %q argument list at offset %d", name, i)
		}
		args = compound[i+1 : close-1]
		hasArgs = true
		i = close
	}

	if doubled || legacyPseudoElements[name] {
		// A pseudo-element is a different BOX, not a state of the element, so it
		// is normalized to its modern double-colon spelling and compared by
		// equality. Any name is accepted: an unknown pseudo-element still names a
		// box that is not the element, which is a decision the model can make.
		element := "::" + name
		if hasArgs {
			element += "(" + stripSpace(args) + ")"
		}
		return nil, element, i, nil
	}

	switch name {
	case "where", "is":
		if !hasArgs {
			return nil, "", 0, fmt.Errorf(":%s used without an argument list at offset %d", name, index)
		}
		// An EMPTY branch is what is left when the marker attribute was lifted out
		// of a matches-any list — `:where([slot-a],[slot-b])` on marker a becomes
		// `:where(,[slot-b])`. That branch matches the subject unconditionally, so
		// the whole disjunction is vacuous. Splitting WITHOUT dropping empties is
		// what makes that visible; splitSelectorList throws the evidence away.
		var branches []string
		for _, branch := range splitSelectorListKeepingEmpty(args) {
			if strings.TrimSpace(branch) == "" {
				return nil, "", i, nil
			}
			branches = append(branches, branch)
		}
		if len(branches) == 0 {
			// `:where()` matches nothing extra — it is vacuous.
			return nil, "", i, nil
		}
		if len(branches) == 1 {
			// A single branch is a conjunction, not a disjunction, so it can be
			// flattened — which is what lets an authored `:where([aria-invalid])`
			// compare equal to a compiled `[aria-invalid]`.
			if inner := splitCompounds(branches[0]); len(inner) == 1 {
				nested, nestedBox, err := compoundConditions(inner[0])
				if err != nil {
					return nil, "", 0, err
				}
				return nested, nestedBox, i, nil
			}
		}
		// A real disjunction is one opaque condition. Branch order carries no
		// meaning, so it is sorted; `:is` and `:where` differ only in specificity,
		// which this model does not judge, so they reduce to one spelling.
		normalized := make([]string, 0, len(branches))
		for _, branch := range branches {
			normalized = append(normalized, stripSpace(branch))
		}
		sort.Strings(normalized)
		return []string{":is(" + strings.Join(normalized, ",") + ")"}, "", i, nil
	}

	if !hasArgs {
		if !statePseudoClasses[name] {
			return nil, "", 0, fmt.Errorf(
				"unsupported pseudo-class %q — the layer gate does not know when it holds, so it "+
					"cannot tell whether a compiled utility shadows this rule", ":"+name)
		}
		return []string{":" + name}, "", i, nil
	}
	if !functionalStatePseudoClasses[name] {
		return nil, "", 0, fmt.Errorf(
			"unsupported functional pseudo-class %q — the layer gate does not know when it holds, "+
				"so it cannot tell whether a compiled utility shadows this rule", ":"+name+"()")
	}
	return []string{":" + name + "(" + stripSpace(args) + ")"}, "", i, nil
}

// normalizeAttributeCondition renders one attribute selector in a single
// spelling, so that an authored `[data-state="open"]` and Tailwind's own
// `[data-state=open]` are recognised as the same condition rather than as two
// conditions neither of which implies the other.
func normalizeAttributeCondition(selector string) (string, error) {
	inner := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(selector, "["), "]"))
	if inner == "" {
		return "", fmt.Errorf("empty attribute selector %q", selector)
	}
	// The operator is the first of =, ~=, |=, ^=, $=, *= at the top level.
	operator := -1
	for i := 0; i < len(inner); i++ {
		if inner[i] == '=' {
			operator = i
			break
		}
	}
	if operator < 0 {
		return "[" + strings.ToLower(stripSpace(inner)) + "]", nil
	}
	prefix := ""
	nameEnd := operator
	if operator > 0 && strings.ContainsRune("~|^$*", rune(inner[operator-1])) {
		prefix = string(inner[operator-1])
		nameEnd = operator - 1
	}
	name := strings.ToLower(strings.TrimSpace(inner[:nameEnd]))
	if name == "" {
		return "", fmt.Errorf("attribute selector %q names no attribute", selector)
	}
	value := strings.TrimSpace(inner[operator+1:])
	// A trailing i/s case-sensitivity flag is part of the condition, kept as-is.
	flag := ""
	if len(value) > 2 && (strings.HasSuffix(value, " i") || strings.HasSuffix(value, " s") ||
		strings.HasSuffix(value, " I") || strings.HasSuffix(value, " S")) {
		flag = strings.ToLower(value[len(value)-1:])
		value = strings.TrimSpace(value[:len(value)-2])
	}
	if len(value) >= 2 && (value[0] == '"' || value[0] == '\'') && value[len(value)-1] == value[0] {
		value = value[1 : len(value)-1]
	}
	out := "[" + name + prefix + "=\"" + value + "\""
	if flag != "" {
		out += " " + flag
	}
	return out + "]", nil
}

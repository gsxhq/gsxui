package stylegen

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"
)

// predicate is a test-side constructor: the conditions in sorted, deduped form,
// which is the shape buildPredicate produces.
func predicate(box string, conditions ...string) statePredicate {
	sort.Strings(conditions)
	return statePredicate{box: box, conditions: dedupe(conditions)}
}

// TestStatePredicateImplicationIsDirectional is the relation itself, stated as a
// table. Every row is a claim about the browser: "whenever the left state holds,
// does the right state hold too?". Equality gets rows 2 and 5 wrong in opposite
// directions, which is why the model is implication.
func TestStatePredicateImplicationIsDirectional(t *testing.T) {
	t.Parallel()

	unconditional := predicate("")
	hover := predicate("", ":hover")
	focus := predicate("", ":focus")
	hoverDark := predicate("", ":hover", "ancestor:.dark")
	open := predicate("", `[data-state="open"]`)
	after := predicate("::after")

	tests := []struct {
		name              string
		authored          statePredicate
		generated         statePredicate
		wantAuthoredIsDed bool
	}{
		{"unconditional under unconditional", unconditional, unconditional, true},
		{"hover under unconditional", hover, unconditional, true},
		{"hover under hover", hover, hover, true},
		{"unconditional under hover", unconditional, hover, false},
		{"hover under focus", hover, focus, false},

		{"hover in dark under hover", hoverDark, hover, true},
		{"hover in dark under dark", hoverDark, predicate("", "ancestor:.dark"), true},
		{"hover under hover in dark", hover, hoverDark, false},
		{"attribute condition under unconditional", open, unconditional, true},
		{"unconditional under attribute condition", unconditional, open, false},
		{"pseudo-element under unconditional", after, unconditional, false},
		{"unconditional under pseudo-element", unconditional, after, false},
		{"pseudo-element under itself", after, after, true},
		{"anything under an unmatchable utility", hover, unmatchablePredicate(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.authored.implies(tt.generated); got != tt.wantAuthoredIsDed {
				t.Errorf("%v.implies(%v) = %v, want %v",
					tt.authored, tt.generated, got, tt.wantAuthoredIsDed)
			}
		})
	}
}

// TestBuildPredicateNormalizesEquivalentSpellings pins the reductions that let
// the authored side and Tailwind's output line up at all. Each pair is two
// spellings of ONE condition; if they did not reduce to the same text the model
// would quietly conclude that neither implies the other.
func TestBuildPredicateNormalizesEquivalentSpellings(t *testing.T) {
	t.Parallel()

	tests := []struct{ name, left, right string }{
		{"quoted and unquoted attribute values", `[data-state="open"]`, `[data-state=open]`},
		{"single-branch :where flattens", `:where([aria-invalid="true"])`, `[aria-invalid="true"]`},
		{"single-branch :is flattens", `:is([data-size="lg"])`, `[data-size="lg"]`},
		{"legacy pseudo-element spelling", `:after`, `::after`},
		{"a matches-any list with an emptied branch is vacuous",
			`:where(,[data-gsxui-slot-menubar-item])`, ``},
		{"condition order does not matter", `:hover[data-state="open"]`, `[data-state="open"]:hover`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			left, err := buildPredicate(nil, nil, tt.left)
			if err != nil {
				t.Fatalf("buildPredicate(%q) error = %v", tt.left, err)
			}
			right, err := buildPredicate(nil, nil, tt.right)
			if err != nil {
				t.Fatalf("buildPredicate(%q) error = %v", tt.right, err)
			}
			if left.box != right.box || !slices.Equal(left.conditions, right.conditions) {
				t.Errorf("buildPredicate(%q) = %v, buildPredicate(%q) = %v — want the same predicate",
					tt.left, left, tt.right, right)
			}
		})
	}
}

// TestBuildPredicateRefusesToGuess is the fail-closed half. An empty condition
// set means UNCONDITIONAL, the weakest predicate, which is implied by
// everything — so reading a form the model does not understand as "no
// conditions" would report shadowing that does not happen. Not knowing must be
// an error, and the error must name what it could not read.
func TestBuildPredicateRefusesToGuess(t *testing.T) {
	t.Parallel()

	tests := []struct{ suffix, want string }{
		{":no-such-pseudo-class", ":no-such-pseudo-class"},
		{":invented-thing(a)", ":invented-thing()"},
		{"&", "unsupported selector character"},
		{"::before::after", "two pseudo-elements"},
	}
	for _, tt := range tests {
		t.Run(tt.suffix, func(t *testing.T) {
			t.Parallel()
			got, err := buildPredicate(nil, nil, tt.suffix)
			if err == nil {
				t.Fatalf("buildPredicate(%q) = %v, want an error", tt.suffix, got)
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("buildPredicate(%q) error = %v, want it to name %q", tt.suffix, err, tt.want)
			}
		})
	}
}

// authoredStateForms re-derives, from the stylesheets the gate actually reads,
// every distinct state a selector-level rule against a slot marker puts its
// subject in.
func authoredStateForms(t *testing.T) map[string][]string {
	t.Helper()
	root := repoRoot(t)
	forms := map[string][]string{}
	for _, relative := range layerCheckedStylesheets {
		path := filepath.Join(root, filepath.FromSlash(relative))
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		rules, err := layeredRules(path, src)
		if err != nil {
			t.Fatal(err)
		}
		for _, rule := range rules {
			for _, complex := range splitSelectorList(rule.selector) {
				marker := ""
				for _, compound := range splitCompounds(complex) {
					index := strings.Index(compound, slotMarkerPrefix)
					if index < 0 {
						continue
					}
					end := index
					for end < len(compound) && isAttributeNameByte(compound[end]) {
						end++
					}
					marker = compound[index:end]
				}
				if marker == "" {
					continue
				}
				_, suffix, ok := subjectConditions(complex, func(compound string) (string, bool) {
					return removeMarkerAttribute(compound, marker)
				})
				if !ok {
					continue // targets a descendant; that path drops the suffix by design
				}
				where := fmt.Sprintf("%s: %s", relative, normalizeSelectorKey(complex))
				forms[stripSpace(suffix)] = append(forms[stripSpace(suffix)], where)
			}
		}
	}
	return forms
}

// TestSupportedStateFormsCoverTheAuthoredStylesheets is the usability half of
// failing closed. Refusing to guess is only safe if the model is broad enough
// that real work does not drown in "unsupported selector" failures — 52
// components are about to migrate through these files, and a gate that rejects
// the CSS the repo is already written in gets switched off.
//
// So this re-derives every state form the authored stylesheets really use and
// demands the model read all of them. A new form arriving with a migration is a
// deliberate line in statePseudoClasses, not a surprise on the day it lands.
func TestSupportedStateFormsCoverTheAuthoredStylesheets(t *testing.T) {
	t.Parallel()

	forms := authoredStateForms(t)
	if len(forms) < 100 {
		t.Fatalf("only %d state forms found — the derivation stopped working", len(forms))
	}
	for suffix, where := range forms {
		if _, err := buildPredicate(nil, nil, suffix); err != nil {
			t.Errorf("state form %q is outside the model: %v\n  used by %d rule(s), e.g. %s",
				suffix, err, len(where), where[0])
		}
	}
}

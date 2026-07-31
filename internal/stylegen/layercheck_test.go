package stylegen

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"
)

// TestComponentComposedMarkersFindsMarkersOnMigratedComponents pins the
// scanner: a data-gsxui-slot-* marker counts only when it sits on an element
// that renders through a component whose presentation is now compiled.
func TestComponentComposedMarkersFindsMarkersOnMigratedComponents(t *testing.T) {
	t.Parallel()

	markers, err := ComponentComposedMarkers(repoRoot(t))
	if err != nil {
		t.Fatalf("ComponentComposedMarkers() error = %v", err)
	}
	for _, want := range []string{
		"data-gsxui-slot-carousel-previous",
		"data-gsxui-slot-carousel-next",
		"data-gsxui-slot-input-group-button",
	} {
		if !slices.Contains(markers, want) {
			t.Errorf("ComponentComposedMarkers() = %q, missing %q", markers, want)
		}
	}
	// carousel-control-label sits on a <span> inside the Button, not on the
	// Button itself, so its presentation does not come from the utilities layer.
	if slices.Contains(markers, "data-gsxui-slot-carousel-control-label") {
		t.Errorf("ComponentComposedMarkers() includes a marker that is not on a migrated component")
	}
	if !slices.IsSorted(markers) {
		t.Errorf("ComponentComposedMarkers() = %q, want sorted", markers)
	}
}

func TestCheckLayerPrecedenceAcceptsTheCurrentTree(t *testing.T) {
	t.Parallel()

	if err := CheckLayerPrecedence(repoRoot(t)); err != nil {
		t.Fatalf("CheckLayerPrecedence() = %v, want nil", err)
	}
}

// TestLayerCheckedStylesheetsCoverEveryAuthoredStylesheet keeps the gate's file
// list honest. The list is spelled out rather than globbed so that no
// stylesheet enters scope unvetted — and this is the other half of that trade:
// a new stylesheet under assets/css or web/ must be added to the list (or to
// the exemptions, with a reason) rather than silently going unchecked.
//
// layerCheckedSheets() is the list PLUS the globbed contents of
// layerCheckedStylesheetDirs, so the one directory that is in scope by
// construction — the default style's per-component sheets — does not need a line
// per component. Everywhere else still fails here until someone decides.
func TestLayerCheckedStylesheetsCoverEveryAuthoredStylesheet(t *testing.T) {
	t.Parallel()

	root := repoRoot(t)
	var found []string
	for _, dir := range []string{filepath.Join("assets", "css"), "web"} {
		err := filepath.WalkDir(filepath.Join(root, dir), func(path string, entry fs.DirEntry, err error) error {
			if err != nil || entry.IsDir() || filepath.Ext(path) != ".css" {
				return err
			}
			relative, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			found = append(found, filepath.ToSlash(relative))
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	sort.Strings(found)

	listed, err := layerCheckedSheets(root)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(listed)
	if !slices.Equal(found, listed) {
		t.Errorf("layerCheckedStylesheets = %q, but the authored stylesheets are %q", listed, found)
	}
	for _, exemption := range repoExemptions(t) {
		if !slices.Contains(listed, exemption.key.file) {
			t.Errorf("exemption %q names a stylesheet the gate does not list", exemption.key.file)
		}
		if strings.TrimSpace(exemption.reason) == "" {
			t.Errorf("exemption %v carries no reason", exemption.key)
		}
	}
}

// repoExemptions builds the committed tree's exemption set, derived groups
// included.
func repoExemptions(t *testing.T) []layerCheckExemption {
	t.Helper()
	exemptions, err := layerCheckExemptions(repoRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	return exemptions
}

// handWrittenExemption returns the first exemption someone actually wrote down.
// The per-entry staleness tests below are about hand-written entries: a derived
// group's entries are arithmetic over a recipe, judged as a group instead (see
// TestDerivedGroupStalenessIsJudgedWholeNotPerEntry).
func handWrittenExemption(t *testing.T) layerCheckExemption {
	t.Helper()
	for _, exemption := range repoExemptions(t) {
		if exemption.derived == "" {
			return exemption
		}
	}
	t.Fatal("no hand-written exemption left to exercise per-entry staleness with")
	return layerCheckExemption{}
}

// repoViolations runs the gate's own collection pass over the committed tree,
// every style, exemptions NOT applied, so a test can ask exactly what the gate
// saw.
func repoViolations(t *testing.T) []violation {
	t.Helper()
	found, err := allLayerViolations(repoRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	return found
}

// TestLayerCheckExemptionsAreAllStillLoadBearing is what makes an exemption
// SELF-INVALIDATING. The old guard only demanded that web/site-button.css hold
// at least one violation somewhere, so deleting the deliberate fallback and
// replacing it with a different, genuinely wrong rule kept the guard green while
// the exemption silently moved to cover the wrong rule. This asserts each
// exemption key still names a violation that really occurs.
func TestLayerCheckExemptionsAreAllStillLoadBearing(t *testing.T) {
	t.Parallel()

	if stale := staleExemptions(repoViolations(t), repoExemptions(t)); len(stale) != 0 {
		t.Errorf("stale exemption(s):\n%s", strings.Join(stale, "\n\n"))
	}
}

// TestExemptionsAreTheOnlyThingKeepingTheTreeGreen pins the other half: the
// exemptions are load-bearing, not decorative. Without them the committed tree
// fails the gate.
func TestExemptionsAreTheOnlyThingKeepingTheTreeGreen(t *testing.T) {
	t.Parallel()

	found := repoViolations(t)
	if len(found) == 0 {
		t.Fatal("no violations at all — web/site-button.css no longer contests compiled Button presentation")
	}
	exempt := exemptionIndex(repoExemptions(t))
	for _, v := range found {
		if _, ok := exempt[v.key]; !ok {
			t.Errorf("unexempted violation on the committed tree:\n%s", v.message)
		}
	}
}

// TestStaleExemptionFailsTheBuild is the self-invalidation, exercised. An
// exemption whose violation no longer exists is reported — and the reason it
// carried is printed with it, so whoever has to re-decide can see what was
// decided before.
func TestStaleExemptionFailsTheBuild(t *testing.T) {
	t.Parallel()

	found := repoViolations(t)

	// Deleting the exempted rule: the selector no longer occurs at all.
	deleted := handWrittenExemption(t)
	deleted.key.selector = ":where([data-gsxui-slot-button]):this-rule-was-deleted"
	stale := staleExemptions(found, []layerCheckExemption{deleted})
	if len(stale) != 1 {
		t.Fatalf("staleExemptions() = %v, want the deleted rule reported", stale)
	}
	if !strings.Contains(stale[0], deleted.reason) {
		t.Errorf("stale diagnostic must carry the exemption's reason, got %q", stale[0])
	}

	// Replacing the exempted rule with a different one under the SAME selector:
	// a file-keyed exemption would still have called this live, which is exactly
	// the hole this key closes.
	replaced := handWrittenExemption(t)
	replaced.key.contested = "bg-some-utility-nobody-wrote"
	if stale := staleExemptions(found, []layerCheckExemption{replaced}); len(stale) != 1 {
		t.Fatalf("staleExemptions() = %v, want the replaced rule reported", stale)
	}

	// And the real list is not stale, so the two above are not vacuous.
	if stale := staleExemptions(found, repoExemptions(t)); len(stale) != 0 {
		t.Fatalf("the committed exemptions are stale: %v", stale)
	}
}

// TestExemptionThatStillHoldsPassesAndSurfacesItsReason is the positive case:
// a live exemption is not reported, and its reason is the one the gate would
// print if it ever went stale.
func TestExemptionThatStillHoldsPassesAndSurfacesItsReason(t *testing.T) {
	t.Parallel()

	live := handWrittenExemption(t)
	if stale := staleExemptions(repoViolations(t), []layerCheckExemption{live}); len(stale) != 0 {
		t.Fatalf("a live exemption was reported stale: %v", stale)
	}
	reason, ok := exemptionIndex(repoExemptions(t))[live.key]
	if !ok || reason == "" {
		t.Fatal("exemption index lost the reason")
	}
	if !strings.Contains(reason, ".spec.ts") {
		t.Errorf("reason must point at what pins the behaviour, got %q", reason)
	}
}

// injectCSS copies the repo into a temp root and splices extra source into the
// named stylesheet after the named anchor, returning the temp root.
func injectCSS(t *testing.T, relative, anchor, extra string) string {
	t.Helper()
	root := t.TempDir()
	copyRepoFixture(t, root)
	path := filepath.Join(root, filepath.FromSlash(relative))
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	injected := bytes.Replace(src, []byte(anchor), []byte(anchor+"\n"+extra), 1)
	if bytes.Equal(injected, src) {
		t.Fatalf("fixture did not change — the %q anchor moved", anchor)
	}
	if err := os.WriteFile(path, injected, 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// injectDefaultCSS splices into the default style's components layer, the gate's
// original and still most common subject.
//
// The target is styles/default/root.css rather than styles/default.css: the
// default style is now one file per component and the aggregator holds only
// @import lines plus the @layer utilities block, so it has no components layer to
// splice into. root.css carries the style's own :root declaration inside `@layer
// components {`, and it is the one file in that directory no component migration
// deletes — every other file there is a component's rules and will be removed
// when that component migrates, which would break these fixtures.
func injectDefaultCSS(t *testing.T, anchor, extra string) string {
	t.Helper()
	return injectCSS(t, "assets/css/styles/default/root.css", anchor, extra)
}

// writeGSXFixture adds a new ui/*.gsx file to an already-copied fixture root
// (see injectCSS/copyRepoFixture), for tests that need a composition shape no
// existing ui/*.gsx has rather than splicing into a stylesheet.
func writeGSXFixture(t *testing.T, root, name, src string) {
	t.Helper()
	path := filepath.Join(root, "ui", name)
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
}

// sharedMarkerFixtureGSX stamps one new marker, data-gsxui-slot-shared-control,
// onto both a <Card> and a <Badge> — two already-migrated components — the
// same way ui/input-group.gsx stamps data-gsxui-slot-input-group-control onto
// both an <Input> and a <Textarea>. Card and Badge both apply a "rounded-*"
// utility on their root slot (rounded-xl and rounded-4xl respectively), so a
// components-layer rule targeting the shared marker with a conflicting
// rounded-md contests both.
const sharedMarkerFixtureGSX = `package ui

import "github.com/gsxhq/gsx"

component SharedMarkerFixtureCard(attrs gsx.Attrs) {
	<Card { attrs... } data-gsxui-slot-shared-control></Card>
}

component SharedMarkerFixtureBadge(attrs gsx.Attrs) {
	<Badge variant="default" { attrs... } data-gsxui-slot-shared-control></Badge>
}
`

// TestCheckLayerPrecedenceRejectsARawDeclarationOverride is finding 1's core
// case. foundation.css is written in plain CSS declarations, not @apply, and
// the gate used to read only assets/css/styles/default.css — so a mechanics
// rule that goes dead the day its component migrates was invisible twice over.
func TestCheckLayerPrecedenceRejectsARawDeclarationOverride(t *testing.T) {
	t.Parallel()

	root := injectCSS(t, "assets/css/foundation.css", "@layer components {",
		`  :where([data-gsxui-slot-carousel-previous]) { display: block; }`)
	err := CheckLayerPrecedence(root)
	if err == nil {
		t.Fatal("CheckLayerPrecedence() = nil, want an error for the raw-declaration override")
	}
	if !strings.Contains(err.Error(), "carousel-previous") {
		t.Errorf("error must name the offending marker, got %q", err)
	}
	if !strings.Contains(err.Error(), "display") {
		t.Errorf("error must name the offending property, got %q", err)
	}
	if !strings.Contains(err.Error(), "foundation.css") {
		t.Errorf("error must name the offending stylesheet, got %q", err)
	}
}

// TestCheckLayerPrecedenceAllowsNonCompetingRawDeclarations is the other side:
// resolving utilities to real property names is what keeps the widened scope
// from flagging every mechanics rule in foundation.css. Neither Button nor
// Carousel (the two components whose utilities reach this marker) sets
// z-index, so this rule still applies in the browser. It used to use
// position/inset; Carousel's own migration made those competing properties,
// which is the gate working, not the fixture aging badly.
func TestCheckLayerPrecedenceAllowsNonCompetingRawDeclarations(t *testing.T) {
	t.Parallel()

	root := injectCSS(t, "assets/css/foundation.css", "@layer components {",
		`  :where([data-gsxui-slot-carousel-previous]) { z-index: 3; }`)
	if err := CheckLayerPrecedence(root); err != nil {
		t.Fatalf("CheckLayerPrecedence() = %v, want nil for a non-competing raw declaration", err)
	}
}

// TestCheckLayerPrecedenceAllowsAnUnconditionalRuleAgainstAHoverOnlyUtility is
// the false positive the migration wave would have hit 153 times over. Button
// declares text-decoration-line only under `hover:underline`; an unconditional
// rule setting it never competes with that, so flagging it would block correct
// work for no cascade reason at all.
func TestCheckLayerPrecedenceAllowsAnUnconditionalRuleAgainstAHoverOnlyUtility(t *testing.T) {
	t.Parallel()

	root := injectCSS(t, "assets/css/foundation.css", "@layer components {",
		`  [data-gsxui-slot-carousel-previous] { text-decoration-line: underline; }`)
	if err := CheckLayerPrecedence(root); err != nil {
		t.Fatalf("CheckLayerPrecedence() = %v, want nil — an unconditional rule cannot contest hover:underline", err)
	}
}

// TestCheckLayerPrecedenceRejectsAHoverRuleAgainstAHoverUtility is the other
// side of the same coin: when the states DO match, the contest is real and the
// rule really is dead.
func TestCheckLayerPrecedenceRejectsAHoverRuleAgainstAHoverUtility(t *testing.T) {
	t.Parallel()

	root := injectCSS(t, "assets/css/foundation.css", "@layer components {",
		`  [data-gsxui-slot-carousel-previous]:hover { text-decoration-line: underline; }`)
	err := CheckLayerPrecedence(root)
	if err == nil {
		t.Fatal("CheckLayerPrecedence() = nil, want an error — :hover contests hover:underline")
	}
	if !strings.Contains(err.Error(), "text-decoration-line") {
		t.Errorf("error must name the contested property, got %q", err)
	}
	if !strings.Contains(err.Error(), "carousel-previous") {
		t.Errorf("error must name the offending marker, got %q", err)
	}
}

// TestCheckLayerPrecedenceStillRejectsUnconditionalAgainstUnconditional guards
// against over-correcting: making the oracle variant-aware must not make it
// blind to the plain case it already caught.
func TestCheckLayerPrecedenceStillRejectsUnconditionalAgainstUnconditional(t *testing.T) {
	t.Parallel()

	root := injectCSS(t, "assets/css/foundation.css", "@layer components {",
		`  [data-gsxui-slot-carousel-previous] { display: block; }`)
	err := CheckLayerPrecedence(root)
	if err == nil {
		t.Fatal("CheckLayerPrecedence() = nil, want an error — display contests inline-flex")
	}
	if !strings.Contains(err.Error(), "display") {
		t.Errorf("error must name the contested property, got %q", err)
	}
}

// TestCheckLayerPrecedenceRejectsARuleNarrowedOnlyByAMediaQuery is the same
// directional reading applied to an at-rule. A `@media` wrapper NARROWS when the
// authored rule applies; it does not narrow when the unconditional utility
// applies. So at every width where the authored rule is live the utility is live
// too, and the rule is dead — even though the compiled Button carries no @media
// wrapper anywhere for the two states to compare equal.
//
// The oracle used to pass this, on the grounds that it could not line the two
// states up. Implication does not need to line them up: a conjunction implies
// each of its conjuncts, and dropping the wrapper is exactly that.
func TestCheckLayerPrecedenceRejectsARuleNarrowedOnlyByAMediaQuery(t *testing.T) {
	t.Parallel()

	root := injectCSS(t, "assets/css/foundation.css", "@layer components {",
		strings.Join([]string{
			`  @media (width >= 40rem) {`,
			`    [data-gsxui-slot-carousel-previous] { display: block; }`,
			`  }`,
		}, "\n"))
	err := CheckLayerPrecedence(root)
	if err == nil {
		t.Fatal("CheckLayerPrecedence() = nil, want an error — the utility applies at every width too")
	}
	if !strings.Contains(err.Error(), "display") {
		t.Errorf("error must name the contested property, got %q", err)
	}
}

// TestCheckLayerPrecedenceAllowsAMediaQueryNarrowingAHoverOnlyUtility keeps the
// row above from being an accident of "any @media contests anything". The
// authored rule is unconditional apart from the width, and the utility applies
// only on hover, so outside hover the authored rule still governs.
func TestCheckLayerPrecedenceAllowsAMediaQueryNarrowingAHoverOnlyUtility(t *testing.T) {
	t.Parallel()

	root := injectCSS(t, "assets/css/foundation.css", "@layer components {",
		strings.Join([]string{
			`  @media (width >= 40rem) {`,
			`    [data-gsxui-slot-carousel-previous] { text-decoration-line: underline; }`,
			`  }`,
		}, "\n"))
	if err := CheckLayerPrecedence(root); err != nil {
		t.Fatalf("CheckLayerPrecedence() = %v, want nil — hover:underline does not apply outside hover", err)
	}
}

// SHADOWING IS DIRECTIONAL — the whole of P1c, as the five rows that define the
// relation. Each row is a claim about the browser, and the oracle must agree:
//
//	authored state | generated state | dead? | why
//	---------------|-----------------|-------|-----------------------------------
//	unconditional  | unconditional   | YES   | same conditions
//	:hover         | unconditional   | YES   | the utility applies during hover too
//	:hover         | :hover          | YES   | same condition
//	unconditional  | :hover          | NO    | the rule still governs the rest
//	:hover         | :focus          | NO    | neither implies the other
//
// The properties are chosen so the generated side is fixed by the recipe, not by
// the test: Button sets border-radius unconditionally (`rounded-lg`/`rounded-4xl`)
// and text-decoration-line only under `hover:underline`.
//
// Row two is the hole. An equality oracle passed it — the authored `:hover` and
// the generated unconditional state are not the same string — and the rule it
// passed can never take effect.
func TestCheckLayerPrecedenceShadowingIsDirectionalForRawDeclarations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		rule     string
		wantDead bool
		property string
	}{
		{
			name:     "unconditional rule under an unconditional utility",
			rule:     `  [data-gsxui-slot-carousel-previous] { border-radius: 0px; }`,
			wantDead: true,
			property: "border-radius",
		},
		{
			name:     "hover rule under an unconditional utility",
			rule:     `  [data-gsxui-slot-carousel-previous]:hover { border-radius: 0px; }`,
			wantDead: true,
			property: "border-radius",
		},
		{
			name:     "hover rule under a hover utility",
			rule:     `  [data-gsxui-slot-carousel-previous]:hover { text-decoration-line: underline; }`,
			wantDead: true,
			property: "text-decoration-line",
		},
		{
			name:     "unconditional rule under a hover utility",
			rule:     `  [data-gsxui-slot-carousel-previous] { text-decoration-line: underline; }`,
			wantDead: false,
			property: "text-decoration-line",
		},
		{
			name:     "focus rule under a hover utility",
			rule:     `  [data-gsxui-slot-carousel-previous]:focus { text-decoration-line: underline; }`,
			wantDead: false,
			property: "text-decoration-line",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root := injectCSS(t, "assets/css/foundation.css", "@layer components {", tt.rule)
			err := CheckLayerPrecedence(root)
			if !tt.wantDead {
				if err != nil {
					t.Fatalf("CheckLayerPrecedence() = %v, want nil — the rule still applies somewhere", err)
				}
				return
			}
			if err == nil {
				t.Fatal("CheckLayerPrecedence() = nil, want an error — the rule can never take effect")
			}
			if !strings.Contains(err.Error(), tt.property) {
				t.Errorf("error must name the contested property %q, got %q", tt.property, err)
			}
			if !strings.Contains(err.Error(), "carousel-previous") {
				t.Errorf("error must name the offending marker, got %q", err)
			}
		})
	}
}

// TestCheckLayerPrecedenceShadowingIsDirectionalForApply is the same five rows
// through the OTHER path, to prove the implication model did not regress it.
//
// This path compares CLASS NAMES through the repo's Tailwind merger rather than
// resolving them to properties, so it reads state from two places at once: the
// utility's own variant prefix, and — for the unconditional-utility rows — the
// enclosing selector. That is why it already caught row two, and it still does.
//
// It has a blind spot the property path does not: with the state in the SELECTOR
// and an unprefixed utility, `:hover { @apply underline }` against a compiled
// `hover:underline` is two different class names to the merger, so row three
// goes unreported here. Pinned below as a known limitation rather than left to
// be discovered — the property path, which is the one the migration's raw
// declarations take, gets it right.
func TestCheckLayerPrecedenceShadowingIsDirectionalForApply(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		rule     string
		wantDead bool
	}{
		{
			name:     "unconditional rule under an unconditional utility",
			rule:     `  [data-gsxui-slot-carousel-previous] { @apply rounded-none; }`,
			wantDead: true,
		},
		{
			name:     "hover rule under an unconditional utility",
			rule:     `  [data-gsxui-slot-carousel-previous]:hover { @apply rounded-none; }`,
			wantDead: true,
		},
		{
			name:     "hover utility under a hover utility",
			rule:     `  [data-gsxui-slot-carousel-previous] { @apply hover:bg-transparent; }`,
			wantDead: true,
		},
		{
			name:     "unconditional rule under a hover utility",
			rule:     `  [data-gsxui-slot-carousel-previous] { @apply underline; }`,
			wantDead: false,
		},
		{
			name:     "focus utility under a hover utility",
			rule:     `  [data-gsxui-slot-carousel-previous] { @apply focus:underline; }`,
			wantDead: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root := injectDefaultCSS(t, "@layer components {", tt.rule)
			err := CheckLayerPrecedence(root)
			if tt.wantDead && err == nil {
				t.Fatal("CheckLayerPrecedence() = nil, want an error — the rule can never take effect")
			}
			if !tt.wantDead && err != nil {
				t.Fatalf("CheckLayerPrecedence() = %v, want nil — the rule still applies somewhere", err)
			}
		})
	}
}

// TestApplyPathStillMissesSelectorStateAgainstAVariantUtility pins the blind
// spot named above, so that it is a recorded limitation rather than a belief. If
// someone teaches the @apply path to read selector state, this test fails and
// they delete it — which is the point.
func TestApplyPathStillMissesSelectorStateAgainstAVariantUtility(t *testing.T) {
	t.Parallel()

	root := injectDefaultCSS(t, "@layer components {",
		`  [data-gsxui-slot-carousel-previous]:hover { @apply underline; }`)
	if err := CheckLayerPrecedence(root); err != nil {
		t.Fatalf("CheckLayerPrecedence() = %v — the @apply path learned selector state; "+
			"delete this test and fold the row into the table above", err)
	}
}

// TestCheckLayerPrecedenceFailsLoudlyOnAnUnsupportedSelectorForm is what makes
// the spec's claim true. The gate says it converts this hazard into a build
// failure; a state it cannot model, read as "no conditions", would instead be
// the weakest possible predicate — implied by everything, contesting nothing,
// silently. It must refuse, and it must name the selector so the refusal is
// actionable.
func TestCheckLayerPrecedenceFailsLoudlyOnAnUnsupportedSelectorForm(t *testing.T) {
	t.Parallel()

	root := injectCSS(t, "assets/css/foundation.css", "@layer components {",
		`  [data-gsxui-slot-carousel-previous]:no-such-pseudo-class { display: block; }`)
	err := CheckLayerPrecedence(root)
	if err == nil {
		t.Fatal("CheckLayerPrecedence() = nil, want an error — the gate cannot model this state")
	}
	for _, want := range []string{
		":no-such-pseudo-class",               // what it could not read
		"[data-gsxui-slot-carousel-previous]", // where
		"assets/css/foundation.css",           // which file
		"internal/stylegen/state.go",          // what to do about it
	} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error must mention %q, got:\n%s", want, err)
		}
	}
}

// TestCheckLayerPrecedenceChecksEveryStyleNotJustTheDefault is P1b. Maia is
// installable and its Button utilities differ substantially from Nova's, so a
// rule can be dead under Maia while a default-only gate stays green. `select-none`
// is in Maia's Button recipe and not in Nova's, which makes user-select a
// property exactly one style takes over.
func TestCheckLayerPrecedenceChecksEveryStyleNotJustTheDefault(t *testing.T) {
	t.Parallel()

	if DefaultStyle == "maia" {
		t.Fatal("this test relies on maia NOT being the default style")
	}
	root := injectCSS(t, "assets/css/foundation.css", "@layer components {",
		`  [data-gsxui-slot-carousel-previous] { user-select: none; }`)
	err := CheckLayerPrecedence(root)
	if err == nil {
		t.Fatal("CheckLayerPrecedence() = nil, want an error — the rule is dead under maia")
	}
	if !strings.Contains(err.Error(), "maia") {
		t.Errorf("the diagnostic must name the style that creates the conflict, got:\n%s", err)
	}
	if strings.Contains(err.Error(), "under style nova") {
		t.Errorf("nova does not set user-select, so it must not be blamed, got:\n%s", err)
	}
	if !strings.Contains(err.Error(), "user-select") {
		t.Errorf("error must name the contested property, got:\n%s", err)
	}
}

// TestEveryStyleIsChecked is the structural half: a style authored under
// registry/styles that the gate never looks at is a silent hole, so the pass
// must produce violations keyed to each one.
func TestEveryStyleIsChecked(t *testing.T) {
	t.Parallel()

	root := repoRoot(t)
	styles, err := styleNames(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(styles) < 2 {
		t.Fatalf("styleNames() = %q — this test needs more than one style to mean anything", styles)
	}
	checked := map[string]struct{}{}
	for _, v := range repoViolations(t) {
		checked[v.key.style] = struct{}{}
	}
	for _, style := range styles {
		if _, ok := checked[style]; !ok {
			t.Errorf("style %q produced no violations at all — it is not being checked", style)
		}
	}
}

// TestExemptionForOneStyleDoesNotExemptAnother is the consequence of putting the
// style in the key, exercised. web/site-button.css contests the same utilities
// under both styles today; an exemption written for Nova alone must leave the
// Maia violation reported, or the exemption list quietly reintroduces exactly
// the blind spot P1b describes.
func TestExemptionForOneStyleDoesNotExemptAnother(t *testing.T) {
	t.Parallel()

	found := repoViolations(t)
	novaOnly := handWrittenExemption(t)
	novaOnly.key.style = "nova"
	if !slices.ContainsFunc(found, func(v violation) bool { return v.key == novaOnly.key }) {
		t.Fatalf("fixture moved: %v is no longer a violation under nova", novaOnly.key)
	}
	maiaKey := novaOnly.key
	maiaKey.style = "maia"
	if !slices.ContainsFunc(found, func(v violation) bool { return v.key == maiaKey }) {
		t.Fatalf("fixture moved: %v is no longer a violation under maia", maiaKey)
	}

	exempt := exemptionIndex([]layerCheckExemption{novaOnly})
	if _, ok := exempt[novaOnly.key]; !ok {
		t.Error("the nova-only exemption does not cover its own violation")
	}
	if _, ok := exempt[maiaKey]; ok {
		t.Error("a nova-only exemption silently forgave the maia violation of the same rule")
	}
}

// TestStaleExemptionIsPerStyle is the other direction: an exemption naming a
// style under which the rule is NOT contested has no subject, and saying so is
// what stops a style list from drifting into a blanket waiver.
func TestStaleExemptionIsPerStyle(t *testing.T) {
	t.Parallel()

	found := repoViolations(t)
	for _, style := range []string{"no-such-style", ""} {
		invented := handWrittenExemption(t)
		invented.key.style = style
		stale := staleExemptions(found, []layerCheckExemption{invented})
		if len(stale) != 1 {
			t.Fatalf("staleExemptions() for style %q = %v, want the exemption reported — a violation "+
				"is always attributed to a named style, so an exemption naming any other one "+
				"(including none at all) has no subject", style, stale)
		}
		if !strings.Contains(stale[0], "under style "+style) {
			t.Errorf("the stale diagnostic must name the style, got %q", stale[0])
		}
	}
}

// TestComposedTargetIsDeterministic pins the tie-break. Marker ownership lives
// in a map, and one compound selector can name markers from two components; the
// answer must not depend on map iteration order.
func TestComposedTargetIsDeterministic(t *testing.T) {
	t.Parallel()

	markers := map[string][]string{
		"data-gsxui-slot-card":          {"card"},
		"data-gsxui-slot-sidebar-inset": {"sidebar"},
	}
	const selector = "[data-gsxui-slot-card][data-gsxui-slot-sidebar-inset]"
	marker, components, _, ok := composedTarget(selector, markers)
	if !ok {
		t.Fatal("composedTarget() found no marker")
	}
	for range 200 {
		gotMarker, gotComponents, _, gotOK := composedTarget(selector, markers)
		if !gotOK || gotMarker != marker || !slices.Equal(gotComponents, components) {
			t.Fatalf("composedTarget() = %q/%q, want a stable %q/%q", gotMarker, gotComponents, marker, components)
		}
	}
	if marker != "data-gsxui-slot-sidebar-inset" {
		t.Errorf("composedTarget() = %q, want the longest matching marker", marker)
	}
}

// TestComposedMarkersAllowsSharedOwnership pins the ruling: a marker stamped
// by two migrated components (see sharedMarkerFixtureGSX) is not an error —
// composedMarkers must return every owner, not just the first, or the second.
func TestComposedMarkersAllowsSharedOwnership(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	copyRepoFixture(t, root)
	writeGSXFixture(t, root, "shared_marker_fixture.gsx", sharedMarkerFixtureGSX)

	byMarker, err := composedMarkers(root)
	if err != nil {
		t.Fatalf("composedMarkers() error = %v, want a shared marker to be accepted, not rejected", err)
	}
	owners := byMarker["data-gsxui-slot-shared-control"]
	want := []string{"badge", "card"}
	if !slices.Equal(owners, want) {
		t.Errorf("composedMarkers()[shared-control] = %q, want %q (sorted, both owners)", owners, want)
	}
}

// TestCheckLayerPrecedenceNamesEveryOwnerOfASharedMarker is the end-to-end
// half of the same ruling: the gate must still fire when a components-layer
// rule targets a marker two migrated components share, and its diagnostic
// must name every contesting owner — not just one — since that name is the
// only thing telling a reader whose utilities are winning.
func TestCheckLayerPrecedenceNamesEveryOwnerOfASharedMarker(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	copyRepoFixture(t, root)
	writeGSXFixture(t, root, "shared_marker_fixture.gsx", sharedMarkerFixtureGSX)
	path := filepath.Join(root, "assets", "css", "styles", "default", "root.css")
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	const anchor = "@layer components {"
	const extra = `  [data-gsxui-slot-shared-control] { @apply rounded-md; }`
	injected := bytes.Replace(src, []byte(anchor), []byte(anchor+"\n"+extra), 1)
	if bytes.Equal(injected, src) {
		t.Fatal("fixture did not change — the anchor moved")
	}
	if err := os.WriteFile(path, injected, 0o644); err != nil {
		t.Fatal(err)
	}

	err = CheckLayerPrecedence(root)
	if err == nil {
		t.Fatal("CheckLayerPrecedence() = nil, want an error: rounded-md contests both Card's and Badge's own rounded-* utility")
	}
	if !strings.Contains(err.Error(), "badge, card") {
		t.Errorf("error must name every contesting owner sorted, got %q", err)
	}
}

func TestCheckLayerPrecedenceRejectsAComponentsLayerOverride(t *testing.T) {
	t.Parallel()

	root := injectDefaultCSS(t, "@layer components {",
		`  [data-gsxui-slot-carousel-previous] { @apply rounded-full; }`)
	err := CheckLayerPrecedence(root)
	if err == nil {
		t.Fatal("CheckLayerPrecedence() = nil, want an error for the components-layer override")
	}
	if !strings.Contains(err.Error(), "carousel-previous") {
		t.Errorf("error must name the offending marker, got %q", err)
	}
	if !strings.Contains(err.Error(), "@layer utilities") {
		t.Errorf("error must say what to do, got %q", err)
	}
	if !strings.Contains(err.Error(), "rounded-full") {
		t.Errorf("error must name the offending utility, got %q", err)
	}
}

// TestCheckLayerPrecedenceRejectsAZeroSpecificityUtilitiesOverride pins the
// second half of the invariant. Moving to @layer utilities is not enough:
// within one layer the cascade falls back to specificity, and :where() is
// deliberately zero-specificity, so the rule still loses to a plain utility.
func TestCheckLayerPrecedenceRejectsAZeroSpecificityUtilitiesOverride(t *testing.T) {
	t.Parallel()

	// The utilities escape hatch stays in the aggregator, not in the
	// per-component files, so this one injects there rather than via
	// injectDefaultCSS.
	root := injectCSS(t, "assets/css/styles/default.css", "@layer utilities {",
		`  :where([data-gsxui-slot-carousel-previous]) { @apply rounded-none; }`)
	err := CheckLayerPrecedence(root)
	if err == nil {
		t.Fatal("CheckLayerPrecedence() = nil, want an error for the zero-specificity override")
	}
	if !strings.Contains(err.Error(), "carousel-previous") {
		t.Errorf("error must name the offending marker, got %q", err)
	}
	if !strings.Contains(err.Error(), "specificity") {
		t.Errorf("error must explain the specificity half of the invariant, got %q", err)
	}
	if strings.Contains(err.Error(), "Move this rule to\n  @layer utilities") {
		t.Errorf("error must not tell an already-correct layer to move, got %q", err)
	}
}

// TestCheckLayerPrecedenceReportsEveryViolation — a sweep is more useful than
// a bisect when a migration lands.
func TestCheckLayerPrecedenceReportsEveryViolation(t *testing.T) {
	t.Parallel()

	root := injectDefaultCSS(t, "@layer components {", strings.Join([]string{
		`  [data-gsxui-slot-carousel-previous] { @apply rounded-full; }`,
		`  [data-gsxui-slot-carousel-next] { @apply h-10; }`,
	}, "\n"))
	err := CheckLayerPrecedence(root)
	if err == nil {
		t.Fatal("CheckLayerPrecedence() = nil, want errors")
	}
	if !strings.Contains(err.Error(), "carousel-previous") || !strings.Contains(err.Error(), "carousel-next") {
		t.Errorf("CheckLayerPrecedence() must report every violation, got %q", err)
	}
}

// TestCheckLayerPrecedenceAllowsNonCompetingComponentsLayerRules pins the
// other side: a property that no compiled utility contests stays legal in
// @layer components. Without this the gate would push unrelated CSS into the
// utilities layer for no reason. z-index, not the positioning geometry this
// used to assert on — Carousel's migration made position/inset competing on
// this marker.
func TestCheckLayerPrecedenceAllowsNonCompetingComponentsLayerRules(t *testing.T) {
	t.Parallel()

	root := injectDefaultCSS(t, "@layer components {",
		`  [data-gsxui-slot-carousel-previous] { @apply z-10; }`)
	if err := CheckLayerPrecedence(root); err != nil {
		t.Fatalf("CheckLayerPrecedence() = %v, want nil for a non-competing rule", err)
	}
}

// TestCheckLayerPrecedenceIgnoresMarkersOnUnmigratedComponents keeps the gate
// scoped to presentation that is actually compiled: a components-layer rule on
// a marker no migrated component renders still wins by layer order.
func TestCheckLayerPrecedenceIgnoresMarkersOnUnmigratedComponents(t *testing.T) {
	t.Parallel()

	root := injectDefaultCSS(t, "@layer components {",
		`  [data-gsxui-slot-carousel-control-label] { @apply rounded-full h-10 text-xs; }`)
	if err := CheckLayerPrecedence(root); err != nil {
		t.Fatalf("CheckLayerPrecedence() = %v, want nil for an unmigrated marker", err)
	}
}

func TestSelectorSpecificity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		selector string
		want     specificity
	}{
		{":where([data-gsxui-slot-input-group-button])", specificity{}},
		{"[data-gsxui-slot-carousel-previous]", specificity{0, 1, 0}},
		{":where([a])[data-size=\"xs\"]", specificity{0, 1, 0}},
		{":where([a])[data-size=\"xs\"] > svg:not([class*=\"size-\"])", specificity{0, 2, 1}},
		{"div p", specificity{0, 0, 2}},
		{"#id .cls span::before", specificity{1, 1, 2}},
		{":is(.a, #b) span", specificity{1, 0, 1}},
		{"a:hover", specificity{0, 1, 1}},
		{"*", specificity{}},
	}
	for _, tt := range tests {
		if got := selectorSpecificity(tt.selector); got != tt.want {
			t.Errorf("selectorSpecificity(%q) = %v, want %v", tt.selector, got, tt.want)
		}
	}
}

// TestSiteButtonExemptionsAreDerivedFromButtonsRecipe pins the property that
// replaced 112 hand-copied entries: every violation the generated fallback
// commits is forgiven, and the forgiveness comes from Button's own recipe rather
// than from a list anyone maintains. Change Button's recipe and this stays true
// with no edit here — which is the whole point.
func TestSiteButtonExemptionsAreDerivedFromButtonsRecipe(t *testing.T) {
	t.Parallel()

	exempt := exemptionIndex(repoExemptions(t))
	var covered int
	for _, v := range repoViolations(t) {
		if v.key.file != siteButtonPath {
			continue
		}
		covered++
		reason, ok := exempt[v.key]
		if !ok {
			t.Errorf("the generated fallback commits a violation nothing derived forgives:\n%s", v.message)
			continue
		}
		if !strings.Contains(reason, "GENERATED") {
			t.Errorf("a %s violation is forgiven by a hand-written reason: %q", siteButtonPath, reason)
		}
	}
	if covered == 0 {
		t.Fatalf("no %s violations at all — the fallback no longer contests compiled Button presentation", siteButtonPath)
	}
}

// TestSiteButtonExemptionsRefuseToDescribeADriftedFile is what keeps deriving
// the exemptions from being a blanket file waiver. The derivation describes the
// generator's output; if the committed file is not that output, it describes
// something that is not on disk, and forgiving it would be exactly the
// "everything in this file is fine" hole a file-keyed exemption opens.
func TestSiteButtonExemptionsRefuseToDescribeADriftedFile(t *testing.T) {
	t.Parallel()

	root := injectCSS(t, siteButtonPath, "@layer components {",
		"  :where(body) :where([data-gsxui-slot-button]) {\n    @apply rounded-full;\n  }")
	if err := CheckLayerPrecedence(root); err == nil {
		t.Fatal("a hand-edited web/site-button.css was accepted")
	} else if !strings.Contains(err.Error(), "drifted") {
		t.Errorf("the diagnostic must name the drift, got %v", err)
	}
}

// TestDerivedGroupStalenessIsJudgedWholeNotPerEntry pins the one way derived
// entries differ from written ones. A derived entry that forgives nothing is
// arithmetic, not a decision that lost its subject; a derived GROUP that
// forgives nothing is a mechanism that should be deleted, and is reported.
func TestDerivedGroupStalenessIsJudgedWholeNotPerEntry(t *testing.T) {
	t.Parallel()

	found := repoViolations(t)
	occurring := map[exemptionKey]struct{}{}
	for _, v := range found {
		occurring[v.key] = struct{}{}
	}
	var derived layerCheckExemption
	for _, exemption := range repoExemptions(t) {
		if _, ok := occurring[exemption.key]; ok && exemption.derived != "" {
			derived = exemption
			break
		}
	}
	if derived.derived == "" {
		t.Fatal("fixture moved: no derived exemption forgives anything")
	}

	dead := derived
	dead.key.contested = "bg-some-utility-nobody-wrote"
	if stale := staleExemptions(found, []layerCheckExemption{derived, dead}); len(stale) != 0 {
		t.Errorf("a derived entry was reported stale on its own: %v", stale)
	}

	orphan := derived
	orphan.key.selector = ":where([data-gsxui-slot-button]):this-rule-was-deleted"
	stale := staleExemptions(found, []layerCheckExemption{orphan})
	if len(stale) != 1 {
		t.Fatalf("staleExemptions() = %v, want the group with nothing left to forgive reported", stale)
	}
	if !strings.Contains(stale[0], siteButtonDerivedGroup) {
		t.Errorf("the diagnostic must name the group, got %q", stale[0])
	}
}

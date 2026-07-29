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

	listed := slices.Clone(layerCheckedStylesheets)
	sort.Strings(listed)
	if !slices.Equal(found, listed) {
		t.Errorf("layerCheckedStylesheets = %q, but the authored stylesheets are %q", listed, found)
	}
	for _, exemption := range layerCheckExemptions {
		if !slices.Contains(listed, exemption.key.file) {
			t.Errorf("exemption %q names a stylesheet the gate does not list", exemption.key.file)
		}
		if strings.TrimSpace(exemption.reason) == "" {
			t.Errorf("exemption %v carries no reason", exemption.key)
		}
	}
}

// repoViolations runs the gate's own collection pass over the committed tree,
// exemptions NOT applied, so a test can ask exactly what the gate saw.
func repoViolations(t *testing.T) []violation {
	t.Helper()
	root := repoRoot(t)
	markers, err := composedMarkers(root)
	if err != nil {
		t.Fatal(err)
	}
	sets, err := componentUtilities(root, markers)
	if err != nil {
		t.Fatal(err)
	}
	found, err := layerViolations(root, markers, sets, &utilityPropertyResolver{root: root, sets: sets})
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

	if stale := staleExemptions(repoViolations(t), layerCheckExemptions); len(stale) != 0 {
		t.Errorf("stale exemption(s):\n%s", strings.Join(stale, "\n\n"))
	}
}

// TestSiteButtonExemptionsAreTheOnlyThingKeepingTheTreeGreen pins the other
// half: the exemptions are load-bearing, not decorative. Without them the
// committed tree fails the gate.
func TestSiteButtonExemptionsAreTheOnlyThingKeepingTheTreeGreen(t *testing.T) {
	t.Parallel()

	found := repoViolations(t)
	if len(found) == 0 {
		t.Fatal("no violations at all — web/site-button.css no longer contests compiled Button presentation")
	}
	exempt := exemptionIndex(layerCheckExemptions)
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
	deleted := layerCheckExemptions[0]
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
	replaced := layerCheckExemptions[0]
	replaced.key.contested = "bg-some-utility-nobody-wrote"
	if stale := staleExemptions(found, []layerCheckExemption{replaced}); len(stale) != 1 {
		t.Fatalf("staleExemptions() = %v, want the replaced rule reported", stale)
	}

	// And the real list is not stale, so the two above are not vacuous.
	if stale := staleExemptions(found, layerCheckExemptions); len(stale) != 0 {
		t.Fatalf("the committed exemptions are stale: %v", stale)
	}
}

// TestExemptionThatStillHoldsPassesAndSurfacesItsReason is the positive case:
// a live exemption is not reported, and its reason is the one the gate would
// print if it ever went stale.
func TestExemptionThatStillHoldsPassesAndSurfacesItsReason(t *testing.T) {
	t.Parallel()

	live := layerCheckExemptions[0]
	if stale := staleExemptions(repoViolations(t), []layerCheckExemption{live}); len(stale) != 0 {
		t.Fatalf("a live exemption was reported stale: %v", stale)
	}
	reason, ok := exemptionIndex(layerCheckExemptions)[live.key]
	if !ok || reason == "" {
		t.Fatal("exemption index lost the reason")
	}
	if !strings.Contains(reason, "basic-demo-presentation.spec.ts") {
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

// injectDefaultCSS splices into the default style sheet, the gate's original
// and still most common subject.
func injectDefaultCSS(t *testing.T, anchor, extra string) string {
	t.Helper()
	return injectCSS(t, "assets/css/styles/default.css", anchor, extra)
}

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
// from flagging every mechanics rule in foundation.css. Button sets neither
// position nor inset, so this rule still applies in the browser.
func TestCheckLayerPrecedenceAllowsNonCompetingRawDeclarations(t *testing.T) {
	t.Parallel()

	root := injectCSS(t, "assets/css/foundation.css", "@layer components {",
		`  :where([data-gsxui-slot-carousel-previous]) { position: absolute; left: -3rem; }`)
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

// TestCanonicalStateDoesNotContestAcrossUncomparableMedia pins the deliberate
// false negative documented on canonicalState. An authored @media wrapper has no
// counterpart in the compiled utilities of this component, and rather than guess
// at equivalence the oracle declines to call it a contest.
func TestCanonicalStateDoesNotContestAcrossUncomparableMedia(t *testing.T) {
	t.Parallel()

	root := injectCSS(t, "assets/css/foundation.css", "@layer components {",
		strings.Join([]string{
			`  @media (width >= 40rem) {`,
			`    [data-gsxui-slot-carousel-previous] { display: block; }`,
			`  }`,
		}, "\n"))
	if err := CheckLayerPrecedence(root); err != nil {
		t.Fatalf("CheckLayerPrecedence() = %v, want nil for a state the oracle cannot compare", err)
	}
}

// TestComposedTargetIsDeterministic pins the tie-break. Marker ownership lives
// in a map, and one compound selector can name markers from two components; the
// answer must not depend on map iteration order.
func TestComposedTargetIsDeterministic(t *testing.T) {
	t.Parallel()

	markers := map[string]string{
		"data-gsxui-slot-card":          "card",
		"data-gsxui-slot-sidebar-inset": "sidebar",
	}
	const selector = "[data-gsxui-slot-card][data-gsxui-slot-sidebar-inset]"
	marker, component, _, ok := composedTarget(selector, markers)
	if !ok {
		t.Fatal("composedTarget() found no marker")
	}
	for range 200 {
		gotMarker, gotComponent, _, gotOK := composedTarget(selector, markers)
		if !gotOK || gotMarker != marker || gotComponent != component {
			t.Fatalf("composedTarget() = %q/%q, want a stable %q/%q", gotMarker, gotComponent, marker, component)
		}
	}
	if marker != "data-gsxui-slot-sidebar-inset" {
		t.Errorf("composedTarget() = %q, want the longest matching marker", marker)
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

	root := injectDefaultCSS(t, "@layer utilities {",
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
// other side: positioning geometry that no compiled utility contests stays
// legal in @layer components. Without this the gate would push unrelated CSS
// into the utilities layer for no reason.
func TestCheckLayerPrecedenceAllowsNonCompetingComponentsLayerRules(t *testing.T) {
	t.Parallel()

	root := injectDefaultCSS(t, "@layer components {",
		`  [data-gsxui-slot-carousel-previous] { @apply absolute -left-12; }`)
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

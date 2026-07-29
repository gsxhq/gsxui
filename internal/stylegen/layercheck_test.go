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
	for relative := range layerCheckExemptions {
		if !slices.Contains(listed, relative) {
			t.Errorf("exemption %q names a stylesheet the gate does not list", relative)
		}
	}
}

// TestSiteButtonFallbackWouldViolateWithoutItsExemption proves the exemption is
// load-bearing rather than decorative: web/site-button.css really does hold
// zero-specificity components-layer rules against the migrated Button's marker,
// and the gate really does read that file now. If the fallback is ever deleted
// or rewritten to win the cascade, this test says so and the exemption can go.
func TestSiteButtonFallbackWouldViolateWithoutItsExemption(t *testing.T) {
	t.Parallel()

	const relative = "web/site-button.css"
	if _, exempt := layerCheckExemptions[relative]; !exempt {
		t.Fatalf("%s is expected to be exempt", relative)
	}
	root := repoRoot(t)
	markers, err := composedMarkers(root)
	if err != nil {
		t.Fatal(err)
	}
	sets, err := componentUtilities(root, markers)
	if err != nil {
		t.Fatal(err)
	}
	src, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relative)))
	if err != nil {
		t.Fatal(err)
	}
	rules, err := layeredRules(relative, src)
	if err != nil {
		t.Fatal(err)
	}
	properties := &utilityPropertyResolver{root: root, sets: sets}
	var violations []string
	for _, rule := range rules {
		found, err := rule.violations(relative, markers, sets, properties)
		if err != nil {
			t.Fatal(err)
		}
		violations = append(violations, found...)
	}
	if len(violations) == 0 {
		t.Fatalf("%s no longer contests compiled Button presentation — drop its exemption", relative)
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

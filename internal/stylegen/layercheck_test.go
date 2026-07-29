package stylegen

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
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

// injectDefaultCSS copies the repo into a temp root and splices extra source in
// after the named anchor, returning the temp root.
func injectDefaultCSS(t *testing.T, anchor, extra string) string {
	t.Helper()
	root := t.TempDir()
	copyRepoFixture(t, root)
	path := filepath.Join(root, "assets", "css", "styles", "default.css")
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

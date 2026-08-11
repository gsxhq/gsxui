package ui_test

import (
	"html"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/internal/stylegen"
	"github.com/gsxhq/gsxui/merge"
	"github.com/gsxhq/gsxui/ui"
)

// novaNativeSelectRecipe loads the default style's NativeSelect recipe once.
// ui.NativeSelect is generated output now, so the classes it renders are the
// default style's concrete utilities — same pattern as novaKbdRecipe.
var novaNativeSelectRecipe = sync.OnceValue(func() recipe.Style {
	path := filepath.Join("..", "registry", "styles", stylegen.DefaultStyle, "native-select.css")
	src, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	style, err := recipe.ParseStyle(path, src)
	if err != nil {
		panic(err)
	}
	return style
})

func nativeSelectRecipeClass(class string, caller ...string) string {
	rule, ok := novaNativeSelectRecipe().Lookup(class)
	if !ok {
		panic("default style declares no recipe " + class)
	}
	classes := append(append([]string(nil), rule.Utilities...), caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

// The wrapper carries a literal "group/native-select" marker class
// (registry/canonical/native-select.gsx) that nativeSelectRecipeClass
// cannot see, for the same reason documented on field_test.go's
// canonicalFieldClass — see the style-porter report's "missing
// group/<name> marker" entry. Renders BEFORE the recipe's own utilities
// (matching NativeSelect's own class={"group/native-select",
// nativeSelect.Wrapper(), attrs.Class()} call), so it is built directly
// rather than through nativeSelectRecipeClass's fixed [recipe..., caller...]
// ordering.
func canonicalNativeSelectWrapperClass(caller ...string) string {
	rule, ok := novaNativeSelectRecipe().Lookup("gsxui-recipe-native-select-wrapper")
	if !ok {
		panic("default style declares no recipe gsxui-recipe-native-select-wrapper")
	}
	classes := append([]string{"group/native-select"}, rule.Utilities...)
	classes = append(classes, caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

func canonicalNativeSelectClass(caller ...string) string {
	return nativeSelectRecipeClass("gsxui-recipe-native-select", caller...)
}

func TestNativeSelectDefault(t *testing.T) {
	got := render(t, ui.NativeSelect(gsx.Raw("<option>x</option>"), nil))
	for _, want := range []string{
		`<div ` + canonicalNativeSelectWrapperClass() + ` data-gsxui-slot-native-select-wrapper>`,
		`<select ` + canonicalNativeSelectClass() + ` data-gsxui-slot-native-select`,
		"<option>x</option></select>",
		"</div>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestNativeSelectIconDependency(t *testing.T) {
	// The wrapper renders icon.ChevronDown as the trigger's chevron — the
	// native-select -> icon import internal/registry derives
	// Deps("native-select") from. This is a render-level proof the icon is
	// actually wired in, not just imported.
	got := render(t, ui.NativeSelect(nil, nil))
	if !strings.Contains(got, "<svg") {
		t.Errorf("expected chevron svg in render\nin: %s", got)
	}
	for _, want := range []string{
		`data-gsxui-slot-icon`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// Caller class merges onto the WRAPPER, not the <select> — width intent
// (w-full, w-[180px]) must size the positioned box the chevron anchors to,
// or the chevron detaches from the select's border. The select always
// fills the wrapper with w-full.
func TestNativeSelectCallerClassMerges(t *testing.T) {
	got := render(t, ui.NativeSelect(nil, gsx.Attrs{{Key: "class", Value: "w-full"}}))
	// The wrapper's own w-fit and the caller's w-full are both plain
	// utilities on the same element now, so merge.Merge drops w-fit rather
	// than the old zero-specificity :where() rule losing the cascade.
	wantWrapper := `<div ` + canonicalNativeSelectWrapperClass("w-full") + ` data-gsxui-slot-native-select-wrapper>`
	if !strings.Contains(got, wantWrapper) {
		t.Errorf("caller w-full should land on the wrapper\nwant: %s\nin: %s", wantWrapper, got)
	}
	if strings.Contains(got, "w-fit") {
		t.Errorf("caller w-full should have displaced the wrapper's own w-fit\nin: %s", got)
	}
	if strings.Count(got, "w-full") != 2 {
		// once on the wrapper (the caller's), once on the <select> (its own).
		t.Errorf("caller class must merge exactly once\nin: %s", got)
	}
}

func TestNativeSelectAttrsFallThrough(t *testing.T) {
	got := render(t, ui.NativeSelect(nil, gsx.Attrs{{Key: "id", Value: "s1"}, {Key: "name", Value: "country"}, {Key: "aria-label", Value: "Country"}}))
	for _, want := range []string{`id="s1"`, `name="country"`, `aria-label="Country"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestNativeSelectForwardsCallerPresenceMarkerToSelectWithoutPrefixFiltering(t *testing.T) {
	got := render(t, ui.NativeSelect(nil, gsx.Attrs{
		{Key: "data-gsxui-slot-calendar-dropdown", Value: true},
		{Key: "id", Value: "month"},
		{Key: "aria-label", Value: "Month"},
	}))
	requirePresenceAttributesOnSameTag(t, got, "<div", "data-gsxui-slot-native-select-wrapper")
	requirePresenceAttributesOnSameTag(t, got, `id="month"`,
		"data-gsxui-slot-calendar-dropdown",
		"data-gsxui-slot-native-select",
	)
	wrapper := openingTagContaining(t, got, "<div")
	if strings.Contains(wrapper, "data-gsxui-slot-calendar-dropdown") {
		t.Errorf("caller select marker leaked onto wrapper\ntag: %s", wrapper)
	}
}

func TestNativeSelectOptionDefault(t *testing.T) {
	got := render(t, ui.NativeSelectOption("us", false, false, gsx.Raw("United States"), nil))
	want := `<option value="us">United States</option>`
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestNativeSelectOptionSelectedAttr(t *testing.T) {
	got := render(t, ui.NativeSelectOption("us", true, false, gsx.Raw("United States"), nil))
	if !strings.Contains(got, " selected") || strings.Contains(got, `selected="`) {
		t.Errorf("selected attr should render bare, not stringified\nin: %s", got)
	}

	got = render(t, ui.NativeSelectOption("us", false, false, gsx.Raw("United States"), nil))
	if strings.Contains(got, "selected") {
		t.Errorf("selected=false should omit the attribute entirely\nin: %s", got)
	}
}

func TestNativeSelectOptionDisabledAttr(t *testing.T) {
	got := render(t, ui.NativeSelectOption("us", false, true, gsx.Raw("United States"), nil))
	if !strings.Contains(got, " disabled") || strings.Contains(got, `disabled="`) {
		t.Errorf("disabled attr should render bare\nin: %s", got)
	}

	got = render(t, ui.NativeSelectOption("us", false, false, gsx.Raw("United States"), nil))
	if strings.Contains(got, "disabled") {
		t.Errorf("disabled=false should omit the attribute entirely\nin: %s", got)
	}
}

func TestNativeSelectOptionAttrsFallThrough(t *testing.T) {
	got := render(t, ui.NativeSelectOption("us", false, false, gsx.Raw("US"), gsx.Attrs{{Key: "data-testid", Value: "opt-us"}}))
	if !strings.Contains(got, `data-testid="opt-us"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestNativeSelectGroupDefault(t *testing.T) {
	got := render(t, ui.NativeSelectGroup("Countries", gsx.Raw("<option>US</option>"), nil))
	want := `<optgroup label="Countries"><option>US</option></optgroup>`
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestNativeSelectGroupAttrsFallThrough(t *testing.T) {
	got := render(t, ui.NativeSelectGroup("Countries", gsx.Raw("x"), gsx.Attrs{{Key: "id", Value: "g1"}}))
	if !strings.Contains(got, `id="g1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestNativeSelectPinned(t *testing.T) {
	// Presentation lives in the stylesheet; the render pin covers structure.
	got := render(t, ui.NativeSelect(gsx.Raw(`<option value="us">United States</option>`), nil))
	want := `<div ` + canonicalNativeSelectWrapperClass() + ` data-gsxui-slot-native-select-wrapper><select ` + canonicalNativeSelectClass() + ` data-gsxui-slot-native-select><option value="us">United States</option></select><svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" data-gsxui-slot-icon><path d="m6 9 6 6 6-6"/></svg></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

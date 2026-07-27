package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestNativeSelectDefault(t *testing.T) {
	got := render(t, ui.NativeSelect(gsx.Raw("<option>x</option>"), nil))
	for _, want := range []string{
		`<div data-gsxui-slot-native-select-wrapper>`,
		`<select data-gsxui-slot-native-select`,
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
	if !strings.Contains(got, `<div class="w-full" data-gsxui-slot-native-select-wrapper>`) {
		t.Errorf("caller w-full should land on the wrapper\nin: %s", got)
	}
	if strings.Count(got, `class="w-full"`) != 1 {
		t.Errorf("caller class must appear exactly once\nin: %s", got)
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
		{Key: "data-gsxui-slot-calendar-dropdown", Value: gsx.Toggle(true)},
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
	want := `<div data-gsxui-slot-native-select-wrapper><select data-gsxui-slot-native-select><option value="us">United States</option></select><svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" data-gsxui-slot-icon><path d="m6 9 6 6 6-6"/></svg></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

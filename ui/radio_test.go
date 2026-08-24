package ui_test

import (
	"html"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestRadioDefault(t *testing.T) {
	got := render(t, ui.Radio(nil))
	for _, want := range []string{
		`<input type="radio"`,
		`data-gsxui-slot-radio`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestRadioCallerClassMerges(t *testing.T) {
	// Not a bare class="size-6" match: Radio now carries its own recipe
	// class, so the caller's size-6 merges in (replacing the recipe's own
	// size-4, same-property override) alongside it.
	got := render(t, ui.Radio(gsx.Attrs{{Key: "class", Value: "size-6"}}))
	if strings.Count(got, "size-6") != 1 {
		t.Errorf("caller class must merge in exactly once\nin: %s", got)
	}
	if strings.Count(got, "class=") != 1 {
		t.Errorf("expected exactly one class= attribute\nin: %s", got)
	}
}

func TestRadioAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Radio(gsx.Attrs{{Key: "id", Value: "r1"}, {Key: "name", Value: "plan"}, {Key: "value", Value: "pro"}, {Key: "aria-label", Value: "Pro"}}))
	for _, want := range []string{`id="r1"`, `name="plan"`, `value="pro"`, `aria-label="Pro"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestRadioCheckedAttr(t *testing.T) {
	got := render(t, ui.Radio(gsx.Attrs{{Key: "checked", Value: true}}))
	if !strings.Contains(got, " checked") || strings.Contains(got, `checked="`) {
		t.Errorf("checked attr should render bare, not stringified\nin: %s", got)
	}

	got = render(t, ui.Radio(gsx.Attrs{{Key: "checked", Value: false}}))
	if strings.Contains(got, `" checked`) || strings.Contains(got, `checked="false"`) {
		t.Errorf("checked=false should omit the attribute entirely\nin: %s", got)
	}
}

func TestRadioDisabledAttr(t *testing.T) {
	got := render(t, ui.Radio(gsx.Attrs{{Key: "disabled", Value: true}}))
	if !strings.Contains(got, " disabled") || strings.Contains(got, `disabled="`) {
		t.Errorf("disabled attr should render bare\nin: %s", got)
	}
}

func TestRadioPinned(t *testing.T) {
	// The class expectation derives from the default style's recipe CSS
	// (registry/styles/nova/radio.css) — see TestCheckboxPinned.
	got := render(t, ui.Radio(nil))
	want := `<input type="radio" class="` +
		html.EscapeString(strings.Join(styleRecipeUtilities("radio", "gsxui-recipe-radio"), " ")) +
		`" data-gsxui-slot-radio>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
	// data-checked:/aria-invalid:aria-checked: were dead selectors — our
	// <input type="radio"> is native. See the style-porter report's
	// "Radio/Checkbox data-checked: -> native :checked:" entry.
	if strings.Contains(got, "data-checked") || strings.Contains(got, "data-unchecked") {
		t.Errorf("dead data-checked/data-unchecked vocabulary in render\nin: %s", got)
	}
}

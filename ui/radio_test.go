package ui_test

import (
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
	// Presentation lives in the stylesheet; the render pin covers structure.
	got := render(t, ui.Radio(nil))
	want := `<input type="radio" class="border-input dark:bg-input/30 data-checked:bg-primary data-checked:text-primary-foreground dark:data-checked:bg-primary data-checked:border-primary aria-invalid:aria-checked:border-primary aria-invalid:border-destructive focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 dark:aria-invalid:border-destructive/50 flex size-4 rounded-full focus-visible:ring-3 aria-invalid:ring-3" data-gsxui-slot-radio>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

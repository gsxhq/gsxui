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
	want := `<input type="radio" class="aspect-square size-4 shrink-0 appearance-none rounded-full border border-input transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:bg-input/30 dark:aria-invalid:ring-destructive/40 checked:border-primary checked:bg-primary checked:text-primary-foreground dark:checked:bg-primary" data-gsxui-slot-radio>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

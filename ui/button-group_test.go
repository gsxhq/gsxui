package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

// TestButtonGroupDefaultPinned pins the zero-value (horizontal) orientation:
// data-orientation defaults via the house |> default pattern, and the class
// list picks the "horizontal" cva block (rounded-l-none/border-l-0/
// rounded-r-none on non-edge children) via a switch value-form, the same
// idiom badge.gsx uses for its own variant map.
func TestButtonGroupDefaultPinned(t *testing.T) {
	got := render(t, ui.ButtonGroup("", gsx.Raw("x"), nil))
	want := `<div role="group" data-orientation="horizontal" data-gsxui-slot-button-group>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestButtonGroupVerticalPinned(t *testing.T) {
	got := render(t, ui.ButtonGroup("vertical", gsx.Raw("x"), nil))
	want := `<div role="group" data-orientation="vertical" data-gsxui-slot-button-group>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestButtonGroupAttrsFallThrough(t *testing.T) {
	got := render(t, ui.ButtonGroup("", nil, gsx.Attrs{{Key: "id", Value: "bg1"}}))
	if !strings.Contains(got, `id="bg1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestButtonGroupCallerClassIsForwardedOnce(t *testing.T) {
	got := render(t, ui.ButtonGroup("", nil, gsx.Attrs{{Key: "class", Value: "gap-2"}}))
	if strings.Count(got, `class="gap-2"`) != 1 {
		t.Errorf("caller class must be the only class and render once\nin: %s", got)
	}
}

func TestButtonGroupTextPinned(t *testing.T) {
	got := render(t, ui.ButtonGroupText(gsx.Raw("x"), nil))
	want := `<div data-gsxui-slot-button-group-text>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestButtonGroupTextAttrsFallThrough(t *testing.T) {
	got := render(t, ui.ButtonGroupText(nil, gsx.Attrs{{Key: "id", Value: "t1"}}))
	if !strings.Contains(got, `id="t1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// TestButtonGroupSeparatorDefaultPinned proves the "vertical" default (the
// opposite of Separator's own "horizontal" default). Separator is migrated
// onto the slot axis now, so its own utilities render as a real class here;
// bg-input and the vertical h-auto still win over Separator's own bg-border
// and h-full, but via the unwrapped default.css rule in @layer utilities
// (assets/css/styles/default.css's ButtonGroupSeparator comment), not
// tailwind-merge — Separator's class and this element's other marker-keyed
// rule both render, layer/specificity settles the conflict.
func TestButtonGroupSeparatorDefaultPinned(t *testing.T) {
	got := render(t, ui.ButtonGroupSeparator("", nil))
	want := `<div role="none" data-orientation="vertical" ` + canonicalSeparatorClass("vertical") + ` data-gsxui-slot-button-group-separator data-gsxui-slot-separator></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestButtonGroupSeparatorOrientationOverride(t *testing.T) {
	got := render(t, ui.ButtonGroupSeparator("horizontal", nil))
	if !strings.Contains(got, `data-orientation="horizontal"`) {
		t.Errorf("missing data-orientation=horizontal override\nin: %s", got)
	}
}

func TestButtonGroupSeparatorAttrsFallThrough(t *testing.T) {
	got := render(t, ui.ButtonGroupSeparator("", gsx.Attrs{{Key: "id", Value: "sep1"}}))
	if !strings.Contains(got, `id="sep1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// TestButtonGroupSeparatorComposesSeparator proves ButtonGroupSeparator
// actually renders through ui.Separator (role="none" + the base
// data-[orientation=...] selectors both come from Separator itself), the
// button-group -> separator dependency internal/registry derives.
func TestButtonGroupSeparatorComposesSeparator(t *testing.T) {
	got := render(t, ui.ButtonGroupSeparator("", nil))
	for _, want := range []string{`role="none"`, `data-gsxui-slot-button-group-separator data-gsxui-slot-separator`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q (expected from ui.Separator)\nin: %s", want, got)
		}
	}
}

// Realistic composition: two buttons split by a separator, the
// button-group-separator demo shape.
func TestButtonGroupWithSeparator(t *testing.T) {
	got := render(t, ui.ButtonGroup("",
		gsx.Fragment(
			ui.Button("secondary", "sm", "", false, gsx.Raw("Copy"), nil),
			ui.ButtonGroupSeparator("", nil),
			ui.Button("secondary", "sm", "", false, gsx.Raw("Paste"), nil),
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-slot-button-group`,
		`>Copy</button>`,
		`data-gsxui-slot-button-group-separator data-gsxui-slot-separator`,
		`>Paste</button>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

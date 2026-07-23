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
	want := `<div role="group" data-slot="button-group" data-orientation="horizontal" class="flex w-fit items-stretch has-[&gt;[data-slot=button-group]]:gap-2 [&amp;&gt;*]:focus-visible:relative [&amp;&gt;*]:focus-visible:z-10 has-[select[aria-hidden=true]:last-child]:[&amp;&gt;[data-slot=select-trigger]:last-of-type]:rounded-r-md [&amp;&gt;[data-slot=select-trigger]:not([class*=&#39;w-&#39;])]:w-fit [&amp;&gt;input]:flex-1 [&amp;&gt;*:not(:first-child)]:rounded-l-none [&amp;&gt;*:not(:first-child)]:border-l-0 [&amp;&gt;*:not(:last-child)]:rounded-r-none">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestButtonGroupVerticalPinned(t *testing.T) {
	got := render(t, ui.ButtonGroup("vertical", gsx.Raw("x"), nil))
	want := `<div role="group" data-slot="button-group" data-orientation="vertical" class="flex w-fit items-stretch has-[&gt;[data-slot=button-group]]:gap-2 [&amp;&gt;*]:focus-visible:relative [&amp;&gt;*]:focus-visible:z-10 has-[select[aria-hidden=true]:last-child]:[&amp;&gt;[data-slot=select-trigger]:last-of-type]:rounded-r-md [&amp;&gt;[data-slot=select-trigger]:not([class*=&#39;w-&#39;])]:w-fit [&amp;&gt;input]:flex-1 flex-col [&amp;&gt;*:not(:first-child)]:rounded-t-none [&amp;&gt;*:not(:first-child)]:border-t-0 [&amp;&gt;*:not(:last-child)]:rounded-b-none">x</div>`
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

func TestButtonGroupCallerClassMerges(t *testing.T) {
	got := render(t, ui.ButtonGroup("", nil, gsx.Attrs{{Key: "class", Value: "gap-2"}}))
	if !strings.Contains(got, "gap-2") {
		t.Errorf("missing caller class gap-2\nin: %s", got)
	}
}

// ButtonGroupText carries no data-slot in shadcn's own source (unlike every
// other button-group part) — ported as-is, see docs/jsx-parity.md.
func TestButtonGroupTextPinned(t *testing.T) {
	got := render(t, ui.ButtonGroupText(gsx.Raw("x"), nil))
	want := `<div class="flex items-center gap-2 rounded-md border bg-muted px-4 text-sm font-medium shadow-xs [&amp;_svg]:pointer-events-none [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
	if strings.Contains(got, "data-slot") {
		t.Errorf("ButtonGroupText must not carry data-slot, matching shadcn's own source\nin: %s", got)
	}
}

func TestButtonGroupTextAttrsFallThrough(t *testing.T) {
	got := render(t, ui.ButtonGroupText(nil, gsx.Attrs{{Key: "id", Value: "t1"}}))
	if !strings.Contains(got, `id="t1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// TestButtonGroupSeparatorDefaultPinned proves the "vertical" default (the
// opposite of Separator's own "horizontal" default) and that
// data-[orientation=vertical]:h-auto / bg-input win their respective
// tailwind-merge conflicts against Separator's own base classes
// (data-[orientation=vertical]:h-full, bg-border — both dropped).
func TestButtonGroupSeparatorDefaultPinned(t *testing.T) {
	got := render(t, ui.ButtonGroupSeparator("", nil))
	want := `<div role="none" data-orientation="vertical" class="shrink-0 data-[orientation=horizontal]:h-px data-[orientation=horizontal]:w-full data-[orientation=vertical]:w-px relative m-0! self-stretch bg-input data-[orientation=vertical]:h-auto" data-slot="button-group-separator"></div>`
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
	for _, want := range []string{`role="none"`, "data-[orientation=horizontal]:h-px"} {
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
		`data-slot="button-group"`,
		`>Copy</button>`,
		`data-slot="button-group-separator"`,
		`>Paste</button>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

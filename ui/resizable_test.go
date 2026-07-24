package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestResizablePanelGroupPinned(t *testing.T) {
	got := render(t, ui.ResizablePanelGroup("", gsx.Raw("x"), nil))
	want := `<div data-slot="resizable-panel-group" data-gsxui-resizable aria-orientation="horizontal" class="flex h-full w-full aria-[orientation=vertical]:flex-col">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestResizablePanelGroupVertical(t *testing.T) {
	got := render(t, ui.ResizablePanelGroup("vertical", gsx.Raw("x"), nil))
	if !strings.Contains(got, `aria-orientation="vertical"`) {
		t.Errorf("want vertical group orientation\nin: %s", got)
	}
}

func TestResizablePanelDefaultSizeIsInlineFlexGrow(t *testing.T) {
	// Review round 2, CRITICAL: a percentage flex-basis can't resolve
	// against an indefinite main-axis container (e.g. a plain min-h-[...]
	// vertical group), which made non-max-w-pinned examples simply not
	// resizable at all. Fixed to a proportional flex-grow weight against a
	// zero LENGTH basis (flex: <n> 1 0px), which always resolves — the
	// split is correct on first paint with JS disabled regardless of
	// whether the group's own size is definite.
	got := render(t, ui.ResizablePanel("50%", "", "", gsx.Raw("x"), nil))
	if !strings.Contains(got, `style="flex: 50 1 0px"`) {
		t.Errorf("want inline flex-grow weight\nin: %s", got)
	}
}

func TestResizablePanelSizedPanelPinned(t *testing.T) {
	// Full-render pin: defaultSize="20%" becomes flex: 20 1 0px verbatim —
	// the numeric part of defaultSize is the grow weight, 1 is shrink, 0px
	// is the (always-resolvable) basis.
	got := render(t, ui.ResizablePanel("20%", "", "", gsx.Raw("x"), nil))
	want := `<div data-slot="resizable-panel" style="flex: 20 1 0px" class="min-w-0 min-h-0 overflow-hidden">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestResizablePanelUnsizedPanelPinned(t *testing.T) {
	// The complementary branch: an unsized panel renders flex: 1 1 0px, an
	// equal-weight share (no class-based flex-1/grow-0 split as of round 2
	// — see the FIX entry in the package doc comment).
	got := render(t, ui.ResizablePanel("", "", "", gsx.Raw("x"), nil))
	want := `<div data-slot="resizable-panel" style="flex: 1 1 0px" class="min-w-0 min-h-0 overflow-hidden">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestResizableHandleOrientationIsInverted(t *testing.T) {
	// The handle's aria-orientation names the RULE, not the group: a
	// horizontal rule (h-px w-full) divides a vertical (flex-col) group.
	// Callers pass the group's orientation; the component inverts.
	h := render(t, ui.ResizableHandle("horizontal", false, nil))
	if !strings.Contains(h, `aria-orientation="vertical"`) {
		t.Errorf("horizontal group wants a vertical rule\nin: %s", h)
	}
	v := render(t, ui.ResizableHandle("vertical", false, nil))
	if !strings.Contains(v, `aria-orientation="horizontal"`) {
		t.Errorf("vertical group wants a horizontal rule\nin: %s", v)
	}
}

func TestResizableHandleWithHandleRendersNovaPill(t *testing.T) {
	// nova's grip is an empty pill, not new-york-v4's icon-in-a-box.
	got := render(t, ui.ResizableHandle("horizontal", true, nil))
	for _, want := range []string{
		`role="separator"`,
		`tabindex="0"`,
		"h-6 w-1 rounded-lg bg-border",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
	if strings.Contains(got, "svg") {
		t.Errorf("nova's grip carries no icon glyph\nin: %s", got)
	}
}

func TestResizableHandleWithoutHandleHasNoGrip(t *testing.T) {
	got := render(t, ui.ResizableHandle("horizontal", false, nil))
	if strings.Contains(got, "h-6 w-1 rounded-lg bg-border") {
		t.Errorf("withHandle=false must not render the grip\nin: %s", got)
	}
}

func TestResizableHandlePinnedHorizontal(t *testing.T) {
	// Full-render pin (review round 1, item 6): only a Contains check on
	// the grip was pinned before, so a dropped token anywhere else in the
	// handle's own class string (verbatim from the source map) was
	// invisible to CI — the group's own string already gets this
	// treatment, matching house style (ui/toggle-group_test.go).
	got := render(t, ui.ResizableHandle("horizontal", false, nil))
	want := `<div data-slot="resizable-handle" role="separator" aria-orientation="vertical" tabindex="0" class="relative flex w-px items-center justify-center bg-border after:absolute after:inset-y-0 after:left-1/2 after:w-1 after:-translate-x-1/2 focus-visible:ring-1 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:outline-hidden aria-[orientation=horizontal]:h-px aria-[orientation=horizontal]:w-full aria-[orientation=horizontal]:after:left-0 aria-[orientation=horizontal]:after:h-1 aria-[orientation=horizontal]:after:w-full aria-[orientation=horizontal]:after:translate-x-0 aria-[orientation=horizontal]:after:-translate-y-1/2 [&amp;[aria-orientation=horizontal]&gt;div]:rotate-90 shrink-0 cursor-col-resize"></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestResizableHandlePinnedVertical(t *testing.T) {
	got := render(t, ui.ResizableHandle("vertical", false, nil))
	want := `<div data-slot="resizable-handle" role="separator" aria-orientation="horizontal" tabindex="0" class="relative flex w-px items-center justify-center bg-border after:absolute after:inset-y-0 after:left-1/2 after:w-1 after:-translate-x-1/2 focus-visible:ring-1 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:outline-hidden aria-[orientation=horizontal]:h-px aria-[orientation=horizontal]:w-full aria-[orientation=horizontal]:after:left-0 aria-[orientation=horizontal]:after:h-1 aria-[orientation=horizontal]:after:w-full aria-[orientation=horizontal]:after:translate-x-0 aria-[orientation=horizontal]:after:-translate-y-1/2 [&amp;[aria-orientation=horizontal]&gt;div]:rotate-90 shrink-0 cursor-row-resize"></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

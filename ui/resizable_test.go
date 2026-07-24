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

func TestResizablePanelDefaultSizeIsInlineFlexBasis(t *testing.T) {
	// Server-rendered geometry: the split is correct on first paint with
	// JS disabled, so defaultSize must land as a real inline style.
	got := render(t, ui.ResizablePanel("50%", "", "", gsx.Raw("x"), nil))
	if !strings.Contains(got, `style="flex-basis: 50%"`) {
		t.Errorf("want inline flex-basis\nin: %s", got)
	}
}

func TestResizablePanelUnsizedHasNoFlexBasis(t *testing.T) {
	got := render(t, ui.ResizablePanel("", "", "", gsx.Raw("x"), nil))
	if strings.Contains(got, "flex-basis") {
		t.Errorf("unsized panel must not stamp flex-basis\nin: %s", got)
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

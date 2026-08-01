package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestSheetPinnedParts(t *testing.T) {
	if got, want := render(t, ui.Sheet(gsx.Raw("x"), nil)), `<div class="contents" data-gsxui-slot-sheet data-gsxui-slot-dialog>x</div>`; got != want {
		t.Errorf("root mismatch\n got: %s\nwant: %s", got, want)
	}
	if got, want := render(t, ui.SheetTrigger(gsx.Raw("Open"), nil)), `<button type="button" aria-haspopup="dialog" aria-expanded="false" data-gsxui-slot-sheet-trigger>Open</button>`; got != want {
		t.Errorf("trigger mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestSheetContentSideAxis(t *testing.T) {
	for _, side := range []string{"right", "left", "top", "bottom"} {
		input := side
		if side == "right" {
			input = ""
		}
		got := render(t, ui.SheetContent(input, true, gsx.Raw("x"), nil))
		for _, want := range []string{
			`data-state="closed"`,
			`data-side="` + side + `"`,
			`data-gsxui-slot-sheet-content data-gsxui-slot-dialog-content`,
		} {
			if !strings.Contains(got, want) {
				t.Errorf("%s missing %q\nin: %s", side, want, got)
			}
		}
		// Sheet migrated to the slot axis: content now legitimately renders
		// its own resolved chrome plus exactly one side class. Pin the whole
		// attribute value rather than merely asserting a class exists — the
		// side axis is the only place a wrong switch arm would be invisible.
		if want := ` class="` + sheetContentClass(side) + `"`; !strings.Contains(got, want) {
			t.Errorf("%s content class mismatch\nwant: %s\n  in: %s", side, want, got)
		}
	}
}

// sheetContentClass is SheetContent's exact resolved class attribute value
// for one side (no caller classes). Sheet's content deliberately does not
// reuse Dialog's, so this is its own chrome, not a narrowing of
// dialogContentClass.
func sheetContentClass(side string) string {
	const base = "m-0 flex-col gap-4 shadow-lg transition ease-in-out max-h-none bg-background text-foreground fixed z-50 text-sm duration-200 outline-none data-[state=closed]:animate-out data-[state=open]:animate-in backdrop:bg-[var(--overlay)] backdrop:backdrop-blur-xs backdrop:duration-200 data-[state=open]:backdrop:animate-in data-[state=open]:backdrop:fade-in-0 data-[state=closed]:backdrop:animate-out data-[state=closed]:backdrop:fade-out-0 open:flex"
	switch side {
	case "right":
		return base + " inset-y-0 end-0 start-auto h-full w-3/4 border-s sm:max-w-sm data-[state=closed]:slide-out-to-end data-[state=open]:slide-in-from-end"
	case "left":
		return base + " inset-y-0 start-0 end-auto h-full w-3/4 border-e sm:max-w-sm data-[state=closed]:slide-out-to-start data-[state=open]:slide-in-from-start"
	case "top":
		return base + " inset-x-0 top-0 bottom-auto h-auto w-full max-w-none border-b data-[state=closed]:slide-out-to-top data-[state=open]:slide-in-from-top"
	case "bottom":
		return base + " inset-x-0 bottom-0 top-auto h-auto w-full max-w-none border-t data-[state=closed]:slide-out-to-bottom data-[state=open]:slide-in-from-bottom"
	}
	panic("unknown sheet side " + side)
}

// sheetCloseButtonClass and sheetCloseIconClass are the resolved classes of
// the close parts SheetContent injects.
func sheetCloseButtonClass() string {
	return "absolute top-3 end-3 rounded-xs opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:ring-2 focus:ring-ring focus:ring-offset-2 focus:outline-hidden disabled:pointer-events-none"
}

func sheetCloseIconClass() string { return "size-4" }

func TestSheetInjectedCloseParts(t *testing.T) {
	got := render(t, ui.SheetContent("", false, gsx.Raw("x"), nil))
	for _, want := range []string{
		`data-gsxui-slot-sheet-close-button data-gsxui-slot-sheet-close`,
		`data-gsxui-slot-sheet-close-icon`,
		`data-gsxui-slot-sheet-close-label`,
		`data-gsxui-dialog-close`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	hidden := render(t, ui.SheetContent("", true, gsx.Raw("x"), nil))
	if strings.Contains(hidden, "sheet-close") {
		t.Errorf("hideCloseButton must omit injected close parts\nin: %s", hidden)
	}
}

func TestSheetSemanticParts(t *testing.T) {
	tests := []struct {
		got  string
		slot string
		hook string
	}{
		{render(t, ui.SheetHeader(gsx.Raw("x"), nil)), "sheet-header", ""},
		{render(t, ui.SheetFooter(gsx.Raw("x"), nil)), "sheet-footer", ""},
		{render(t, ui.SheetTitle(gsx.Raw("x"), nil)), "sheet-title", "data-gsxui-slot-dialog-title"},
		{render(t, ui.SheetDescription(gsx.Raw("x"), nil)), "sheet-description", "data-gsxui-slot-dialog-description"},
		{render(t, ui.SheetClose(gsx.Raw("x"), nil)), "sheet-close", "data-gsxui-dialog-close"},
	}
	for _, tt := range tests {
		if !strings.Contains(tt.got, `data-gsxui-slot-`+tt.slot) ||
			(tt.hook != "" && !strings.Contains(tt.got, tt.hook)) {
			t.Errorf("%s semantic part mismatch\nin: %s", tt.slot, tt.got)
		}
	}
}

func TestSheetContentComposesPresenceMarkersOnDialog(t *testing.T) {
	got := render(t, ui.SheetContent("", true, gsx.Raw("x"), nil))
	requirePresenceAttributesOnSameTag(t, got, "data-gsxui-slot-sheet-content",
		"data-gsxui-slot-dialog-content",
	)
}

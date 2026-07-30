package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestSheetPinnedParts(t *testing.T) {
	if got, want := render(t, ui.Sheet(gsx.Raw("x"), nil)), `<div class="contents" data-gsxui-dialog data-gsxui-slot-sheet data-gsxui-slot-dialog>x</div>`; got != want {
		t.Errorf("root mismatch\n got: %s\nwant: %s", got, want)
	}
	if got, want := render(t, ui.SheetTrigger(gsx.Raw("Open"), nil)), `<button data-gsxui-dialog-trigger type="button" aria-haspopup="dialog" aria-expanded="false" data-gsxui-slot-sheet-trigger>Open</button>`; got != want {
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
		if strings.Contains(got, ` class=`) {
			t.Errorf("%s default content must not render a class\nin: %s", side, got)
		}
	}
}

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
		{render(t, ui.SheetTitle(gsx.Raw("x"), nil)), "sheet-title", "data-gsxui-dialog-title"},
		{render(t, ui.SheetDescription(gsx.Raw("x"), nil)), "sheet-description", "data-gsxui-dialog-description"},
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
	requirePresenceAttributesOnSameTag(t, got, "data-gsxui-dialog-content",
		"data-gsxui-slot-sheet-content",
		"data-gsxui-slot-dialog-content",
	)
}

package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestDrawerPinnedParts(t *testing.T) {
	if got, want := render(t, ui.Drawer(gsx.Raw("x"), nil)), `<div data-gsxui-dialog data-gsxui-slot-drawer data-gsxui-slot-dialog>x</div>`; got != want {
		t.Errorf("root mismatch\n got: %s\nwant: %s", got, want)
	}
	if got, want := render(t, ui.DrawerTrigger(gsx.Raw("Open"), nil)), `<button data-gsxui-dialog-trigger type="button" aria-haspopup="dialog" aria-expanded="false" data-gsxui-slot-drawer-trigger>Open</button>`; got != want {
		t.Errorf("trigger mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestDrawerContentSideAxis(t *testing.T) {
	for _, side := range []string{"bottom", "top", "left", "right"} {
		input := side
		if side == "bottom" {
			input = ""
		}
		got := render(t, ui.DrawerContent(input, gsx.Raw("x"), nil))
		for _, want := range []string{
			`data-state="closed"`,
			`data-side="` + side + `"`,
			`data-gsxui-slot-drawer-content data-gsxui-slot-dialog-content`,
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

func TestDrawerHandleBottomOnly(t *testing.T) {
	for _, side := range []string{"", "bottom"} {
		got := render(t, ui.DrawerContent(side, gsx.Raw("x"), nil))
		if strings.Count(got, `data-gsxui-slot-drawer-handle`) != 1 {
			t.Errorf("%q drawer must have one handle\nin: %s", side, got)
		}
	}
	for _, side := range []string{"top", "left", "right"} {
		got := render(t, ui.DrawerContent(side, gsx.Raw("x"), nil))
		if strings.Contains(got, "drawer-handle") {
			t.Errorf("%s drawer must omit handle\nin: %s", side, got)
		}
	}
}

func TestDrawerSemanticParts(t *testing.T) {
	tests := []struct {
		got  string
		slot string
		hook string
	}{
		{render(t, ui.DrawerHeader(gsx.Raw("x"), nil)), "drawer-header", ""},
		{render(t, ui.DrawerFooter(gsx.Raw("x"), nil)), "drawer-footer", ""},
		{render(t, ui.DrawerTitle(gsx.Raw("x"), nil)), "drawer-title", "data-gsxui-dialog-title"},
		{render(t, ui.DrawerDescription(gsx.Raw("x"), nil)), "drawer-description", "data-gsxui-dialog-description"},
		{render(t, ui.DrawerClose(gsx.Raw("x"), nil)), "drawer-close", "data-gsxui-dialog-close"},
	}
	for _, tt := range tests {
		if !strings.Contains(tt.got, `data-gsxui-slot-`+tt.slot) ||
			(tt.hook != "" && !strings.Contains(tt.got, tt.hook)) {
			t.Errorf("%s semantic part mismatch\nin: %s", tt.slot, tt.got)
		}
	}
}

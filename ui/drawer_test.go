package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestDrawerPinnedParts(t *testing.T) {
	if got, want := render(t, ui.Drawer(gsx.Raw("x"), nil)), `<div class="contents" data-gsxui-slot-drawer data-gsxui-slot-dialog>x</div>`; got != want {
		t.Errorf("root mismatch\n got: %s\nwant: %s", got, want)
	}
	if got, want := render(t, ui.DrawerTrigger(gsx.Raw("Open"), nil)), `<button type="button" aria-haspopup="dialog" aria-expanded="false" data-gsxui-slot-drawer-trigger>Open</button>`; got != want {
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
		// Drawer migrated to the slot axis: content now legitimately renders
		// its own resolved chrome plus exactly one side class. Pin the whole
		// attribute value rather than merely asserting a class exists — the
		// side axis is the only place a wrong switch arm would be invisible.
		if want := ` class="` + drawerContentClass(side) + `"`; !strings.Contains(got, want) {
			t.Errorf("%s content class mismatch\nwant: %s\n  in: %s", side, want, got)
		}
	}
}

// drawerContentClass is DrawerContent's exact resolved class attribute value
// for one side (no caller classes), HTML-escaped as it renders — the "&" of
// the header text-align arbitrary variant comes out as "&amp;". Drawer's
// content deliberately does not reuse Dialog's, so this is its own chrome,
// not a narrowing of dialogContentClass.
func drawerContentClass(side string) string {
	const base = "m-0 flex-col gap-4 shadow-lg transition ease-in-out bg-popover text-popover-foreground fixed z-50 text-sm duration-200 outline-none data-[state=closed]:animate-out data-[state=open]:animate-in backdrop:bg-[var(--overlay)] backdrop:backdrop-blur-xs backdrop:duration-200 data-[state=open]:backdrop:animate-in data-[state=open]:backdrop:fade-in-0 data-[state=closed]:backdrop:animate-out data-[state=closed]:backdrop:fade-out-0 open:flex"
	switch side {
	case "bottom":
		return base + " inset-x-0 bottom-0 top-auto mt-24 h-auto max-h-[80vh] w-full max-w-none rounded-t-xl border-t data-[state=closed]:slide-out-to-bottom data-[state=open]:slide-in-from-bottom [&amp;_[data-gsxui-slot-drawer-header]]:text-center"
	case "top":
		return base + " inset-x-0 top-0 bottom-auto mb-24 h-auto max-h-[80vh] w-full max-w-none rounded-b-xl border-b data-[state=closed]:slide-out-to-top data-[state=open]:slide-in-from-top [&amp;_[data-gsxui-slot-drawer-header]]:text-center"
	case "left":
		return base + " inset-y-0 start-0 end-auto h-full max-h-none w-3/4 rounded-e-xl border-e sm:max-w-sm data-[state=closed]:slide-out-to-start data-[state=open]:slide-in-from-start md:[&amp;_[data-gsxui-slot-drawer-header]]:text-start"
	case "right":
		return base + " inset-y-0 end-0 start-auto h-full max-h-none w-3/4 rounded-s-xl border-s sm:max-w-sm data-[state=closed]:slide-out-to-end data-[state=open]:slide-in-from-end md:[&amp;_[data-gsxui-slot-drawer-header]]:text-start"
	}
	panic("unknown drawer side " + side)
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
		{render(t, ui.DrawerTitle(gsx.Raw("x"), nil)), "drawer-title", "data-gsxui-slot-dialog-title"},
		{render(t, ui.DrawerDescription(gsx.Raw("x"), nil)), "drawer-description", "data-gsxui-slot-dialog-description"},
		{render(t, ui.DrawerClose(gsx.Raw("x"), nil)), "drawer-close", "data-gsxui-dialog-close"},
	}
	for _, tt := range tests {
		if !strings.Contains(tt.got, `data-gsxui-slot-`+tt.slot) ||
			(tt.hook != "" && !strings.Contains(tt.got, tt.hook)) {
			t.Errorf("%s semantic part mismatch\nin: %s", tt.slot, tt.got)
		}
	}
}

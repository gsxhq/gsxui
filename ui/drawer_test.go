package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

// TestDrawerStructure is an integration-shaped smoke test composing every
// part together, mirroring TestSheetStructure.
func TestDrawerStructure(t *testing.T) {
	got := render(t, ui.Drawer(
		gsx.Fragment(
			ui.DrawerTrigger(gsx.Raw("Open"), nil),
			ui.DrawerContent("", gsx.Fragment(
				ui.DrawerHeader(gsx.Fragment(
					ui.DrawerTitle(gsx.Raw("Move Goal"), nil),
					ui.DrawerDescription(gsx.Raw("Set your daily activity goal."), nil),
				), nil),
				ui.DrawerFooter(ui.DrawerClose(gsx.Raw("Cancel"), nil), nil),
			), nil),
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-dialog`,         // root hook (inherited from Dialog)
		`data-slot="drawer"`,        // root slot override
		`data-gsxui-dialog-trigger`, // trigger hook
		`data-slot="drawer-trigger"`,
		`aria-haspopup="dialog"`,
		`aria-expanded="false"`,
		"<dialog",
		`data-gsxui-dialog-content`, // dialog.js still selects on this
		`data-state="closed"`,
		`data-side="bottom"`, // default direction stamp
		`data-slot="drawer-content"`,
		`data-slot="drawer-handle"`, // handle bar visible for bottom (the default)
		`data-slot="drawer-header"`,
		`data-slot="drawer-title"`, ">Move Goal<",
		`data-slot="drawer-description"`,
		`data-slot="drawer-footer"`,
		`data-slot="drawer-close"`, ">Cancel<",
		`data-gsxui-dialog-close`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestDrawerRootPinned(t *testing.T) {
	got := render(t, ui.Drawer(gsx.Raw("x"), nil))
	want := `<div data-gsxui-dialog class="contents" data-slot="drawer">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestDrawerTriggerPinned(t *testing.T) {
	got := render(t, ui.DrawerTrigger(gsx.Raw("Open"), nil))
	want := `<button data-slot="drawer-trigger" data-gsxui-dialog-trigger type="button" aria-haspopup="dialog" aria-expanded="false">Open</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestDrawerContentPinnedBottom pins the exact full render for
// DrawerContent("", ...) — the zero-value direction (bottom, vaul's own
// default, distinct from Sheet's "right" default — see ui/drawer.gsx's own
// doc comment). Every class token verified against
// docs/superpowers/plans/2026-07-24-tier3-source-map-wrapped.md `## drawer`
// §3's recommended base + bottom class strings (nova bg-popover/
// text-popover-foreground, h-1 handle, rounded-t-xl), transcribed-not-
// verified per the map's own caveat — the controller's browser pass
// verifies these render correctly, mirroring sheet's own six ADAPTs found
// only in-browser.
func TestDrawerContentPinnedBottom(t *testing.T) {
	got := render(t, ui.DrawerContent("", gsx.Raw("x"), nil))
	want := `<dialog data-slot="drawer-content" data-gsxui-dialog-content data-state="closed" data-side="bottom" class="fixed z-50 m-0 open:flex flex-col gap-4 bg-popover text-popover-foreground text-sm shadow-lg transition ease-in-out duration-200 data-[state=closed]:animate-out data-[state=open]:animate-in backdrop:bg-black/10 backdrop:duration-200 supports-backdrop-filter:backdrop:backdrop-blur-xs data-[state=open]:backdrop:animate-in data-[state=open]:backdrop:fade-in-0 data-[state=closed]:backdrop:animate-out data-[state=closed]:backdrop:fade-out-0 inset-x-0 bottom-0 top-auto w-full max-w-none h-auto mt-24 max-h-[80vh] rounded-t-xl border-t data-[state=closed]:slide-out-to-bottom data-[state=open]:slide-in-from-bottom"><div data-slot="drawer-handle" class="mx-auto mt-4 h-1 w-[100px] shrink-0 rounded-full bg-muted"></div>x</dialog>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestDrawerContentPinnedTop(t *testing.T) {
	got := render(t, ui.DrawerContent("top", gsx.Raw("x"), nil))
	want := `<dialog data-slot="drawer-content" data-gsxui-dialog-content data-state="closed" data-side="top" class="fixed z-50 m-0 open:flex flex-col gap-4 bg-popover text-popover-foreground text-sm shadow-lg transition ease-in-out duration-200 data-[state=closed]:animate-out data-[state=open]:animate-in backdrop:bg-black/10 backdrop:duration-200 supports-backdrop-filter:backdrop:backdrop-blur-xs data-[state=open]:backdrop:animate-in data-[state=open]:backdrop:fade-in-0 data-[state=closed]:backdrop:animate-out data-[state=closed]:backdrop:fade-out-0 inset-x-0 top-0 bottom-auto w-full max-w-none h-auto mb-24 max-h-[80vh] rounded-b-xl border-b data-[state=closed]:slide-out-to-top data-[state=open]:slide-in-from-top">x</dialog>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestDrawerContentPinnedLeft(t *testing.T) {
	got := render(t, ui.DrawerContent("left", gsx.Raw("x"), nil))
	want := `<dialog data-slot="drawer-content" data-gsxui-dialog-content data-state="closed" data-side="left" class="fixed z-50 m-0 open:flex flex-col gap-4 bg-popover text-popover-foreground text-sm shadow-lg transition ease-in-out duration-200 data-[state=closed]:animate-out data-[state=open]:animate-in backdrop:bg-black/10 backdrop:duration-200 supports-backdrop-filter:backdrop:backdrop-blur-xs data-[state=open]:backdrop:animate-in data-[state=open]:backdrop:fade-in-0 data-[state=closed]:backdrop:animate-out data-[state=closed]:backdrop:fade-out-0 inset-y-0 left-0 right-auto h-full max-h-none w-3/4 rounded-r-xl border-r sm:max-w-sm data-[state=closed]:slide-out-to-left data-[state=open]:slide-in-from-left">x</dialog>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestDrawerContentPinnedRight(t *testing.T) {
	got := render(t, ui.DrawerContent("right", gsx.Raw("x"), nil))
	want := `<dialog data-slot="drawer-content" data-gsxui-dialog-content data-state="closed" data-side="right" class="fixed z-50 m-0 open:flex flex-col gap-4 bg-popover text-popover-foreground text-sm shadow-lg transition ease-in-out duration-200 data-[state=closed]:animate-out data-[state=open]:animate-in backdrop:bg-black/10 backdrop:duration-200 supports-backdrop-filter:backdrop:backdrop-blur-xs data-[state=open]:backdrop:animate-in data-[state=open]:backdrop:fade-in-0 data-[state=closed]:backdrop:animate-out data-[state=closed]:backdrop:fade-out-0 inset-y-0 right-0 left-auto h-full max-h-none w-3/4 rounded-l-xl border-l sm:max-w-sm data-[state=closed]:slide-out-to-right data-[state=open]:slide-in-from-right">x</dialog>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestDrawerHandleBarBottomOnly proves the handle-bar visibility rule: the
// mobile-sheet drag-handle affordance renders only for the bottom
// direction (vaul's own default and convention), and is absent — not
// merely hidden via CSS — for top/left/right, since direction is
// server-known and gsxui has no vaul underneath to key a group-data
// selector off of (see ui/drawer.gsx's own doc comment).
func TestDrawerHandleBarBottomOnly(t *testing.T) {
	for _, dir := range []string{"top", "left", "right"} {
		got := render(t, ui.DrawerContent(dir, gsx.Raw("x"), nil))
		if strings.Contains(got, "drawer-handle") {
			t.Errorf("direction %q must not render the handle bar\nin: %s", dir, got)
		}
	}
	got := render(t, ui.DrawerContent("bottom", gsx.Raw("x"), nil))
	if !strings.Contains(got, `data-slot="drawer-handle"`) {
		t.Errorf("bottom direction must render the handle bar\nin: %s", got)
	}
}

func TestDrawerContentAttrsFallThrough(t *testing.T) {
	got := render(t, ui.DrawerContent("", gsx.Raw("x"), gsx.Attrs{{Key: "id", Value: "d1"}, {Key: "class", Value: "max-w-md"}}))
	for _, want := range []string{`id="d1"`, "max-w-md"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// TestDrawerContentDisplayGated mirrors TestSheetContentDisplayGated: the
// base class gates display behind the open: variant rather than shipping
// an unscoped flex, which would defeat the UA's closed-dialog display:none.
func TestDrawerContentDisplayGated(t *testing.T) {
	got := render(t, ui.DrawerContent("", gsx.Raw("x"), nil))
	if !strings.Contains(got, "open:flex") {
		t.Errorf("want open:flex (display-gated), got: %s", got)
	}
	for _, tok := range strings.Fields(got) {
		if tok == "flex" {
			t.Errorf("unscoped flex token present (defeats closed-state display:none)\nin: %s", got)
		}
	}
}

func TestDrawerHeaderPinned(t *testing.T) {
	got := render(t, ui.DrawerHeader(gsx.Raw("x"), nil))
	want := `<div data-slot="drawer-header" class="flex flex-col gap-0.5 p-4 text-center">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestDrawerFooterPinned(t *testing.T) {
	got := render(t, ui.DrawerFooter(gsx.Raw("x"), nil))
	want := `<div data-slot="drawer-footer" class="mt-auto flex flex-col gap-2 p-4">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestDrawerTitlePinned(t *testing.T) {
	got := render(t, ui.DrawerTitle(gsx.Raw("Move Goal"), nil))
	want := `<h2 data-slot="drawer-title" class="font-medium text-foreground">Move Goal</h2>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestDrawerDescriptionPinned(t *testing.T) {
	got := render(t, ui.DrawerDescription(gsx.Raw("Set your daily activity goal."), nil))
	want := `<p data-slot="drawer-description" class="text-sm text-muted-foreground">Set your daily activity goal.</p>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestDrawerClosePinned(t *testing.T) {
	got := render(t, ui.DrawerClose(gsx.Raw("Cancel"), nil))
	want := `<button data-slot="drawer-close" data-gsxui-dialog-close type="button">Cancel</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

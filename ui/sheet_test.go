package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

// TestSheetStructure is an integration-shaped smoke test composing every
// part together, mirroring TestDialogStructure/TestAlertDialogStructure.
func TestSheetStructure(t *testing.T) {
	got := render(t, ui.Sheet(
		gsx.Fragment(
			ui.SheetTrigger(gsx.Raw("Open"), nil),
			ui.SheetContent("", false, gsx.Fragment(
				ui.SheetHeader(gsx.Fragment(
					ui.SheetTitle(gsx.Raw("Edit profile"), nil),
					ui.SheetDescription(gsx.Raw("Make changes to your profile here."), nil),
				), nil),
				ui.SheetFooter(ui.SheetClose(gsx.Raw("Save changes"), nil), nil),
			), nil),
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-dialog`,         // root hook (inherited from Dialog)
		`data-slot="sheet"`,         // root slot override
		`data-gsxui-dialog-trigger`, // trigger hook
		`data-slot="sheet-trigger"`,
		`aria-haspopup="dialog"`,
		`aria-expanded="false"`,
		"<dialog",
		`data-gsxui-dialog-content`, // dialog.js still selects on this
		`data-state="closed"`,
		`data-side="right"`, // default side stamp
		`data-slot="sheet-content"`,
		`data-slot="sheet-header"`,
		`data-slot="sheet-title"`, ">Edit profile<",
		`data-slot="sheet-description"`,
		`data-slot="sheet-footer"`,
		`data-slot="sheet-close"`, ">Save changes<",
		`data-gsxui-dialog-close`,
		`aria-label="Close"`, // the injected X close button
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestSheetRootPinned(t *testing.T) {
	got := render(t, ui.Sheet(gsx.Raw("x"), nil))
	want := `<div data-gsxui-dialog class="contents" data-slot="sheet">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestSheetTriggerPinned(t *testing.T) {
	got := render(t, ui.SheetTrigger(gsx.Raw("Open"), nil))
	want := `<button data-slot="sheet-trigger" data-gsxui-dialog-trigger type="button" aria-haspopup="dialog" aria-expanded="false">Open</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestSheetContentPinnedDefaultSide pins the exact full render for
// SheetContent("", false, ...) — the zero-value side (right, shadcn's own
// default) and hideCloseButton false (shadcn's showCloseButton default of
// true), so the injected X renders. Every class token verified against
// registry/new-york-v4/ui/sheet.tsx's own base + side="right" blocks, plus
// the three ADAPTs documented in ui/sheet.gsx's own SheetContent doc
// comment (open:flex, text-foreground, m-0) and the folded ::backdrop
// (dialog's own ADAPT, reused here).
func TestSheetContentPinnedDefaultSide(t *testing.T) {
	got := render(t, ui.SheetContent("", false, gsx.Raw("x"), nil))
	want := `<dialog data-slot="sheet-content" data-gsxui-dialog-content data-state="closed" data-side="right" class="fixed z-50 m-0 max-h-none open:flex flex-col gap-4 bg-background text-foreground shadow-lg transition ease-in-out data-[state=closed]:animate-out data-[state=closed]:duration-300 data-[state=open]:animate-in data-[state=open]:duration-500 backdrop:bg-black/50 inset-y-0 right-0 left-auto h-full w-3/4 border-l data-[state=closed]:slide-out-to-right data-[state=open]:slide-in-from-right sm:max-w-sm">x<button type="button" data-slot="sheet-close" data-gsxui-dialog-close aria-label="Close" class="absolute top-4 right-4 rounded-xs opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:ring-2 focus:ring-ring focus:ring-offset-2 focus:outline-hidden disabled:pointer-events-none data-[state=open]:bg-secondary"><svg class="size-4" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg></button></dialog>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestSheetContentSideLeft(t *testing.T) {
	got := render(t, ui.SheetContent("left", true, gsx.Raw("x"), nil))
	for _, want := range []string{
		`data-side="left"`,
		"inset-y-0 left-0 right-auto h-full w-3/4 border-r",
		"data-[state=closed]:slide-out-to-left",
		"data-[state=open]:slide-in-from-left",
		"sm:max-w-sm",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if strings.Contains(got, "right-0") {
		t.Errorf("left side must not carry right-anchoring tokens\nin: %s", got)
	}
}

func TestSheetContentSideTop(t *testing.T) {
	got := render(t, ui.SheetContent("top", true, gsx.Raw("x"), nil))
	for _, want := range []string{
		`data-side="top"`,
		"inset-x-0 top-0 bottom-auto w-full max-w-none h-auto border-b",
		"data-[state=closed]:slide-out-to-top",
		"data-[state=open]:slide-in-from-top",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if strings.Contains(got, "sm:max-w-sm") {
		t.Errorf("top side must not carry the left/right sm:max-w-sm token\nin: %s", got)
	}
}

func TestSheetContentSideBottom(t *testing.T) {
	got := render(t, ui.SheetContent("bottom", true, gsx.Raw("x"), nil))
	for _, want := range []string{
		`data-side="bottom"`,
		"inset-x-0 bottom-0 top-auto w-full max-w-none h-auto border-t",
		"data-[state=closed]:slide-out-to-bottom",
		"data-[state=open]:slide-in-from-bottom",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// TestSheetContentHideCloseButton proves hideCloseButton omits the injected
// X, same contract as DialogContent's own hideCloseButton.
func TestSheetContentHideCloseButton(t *testing.T) {
	got := render(t, ui.SheetContent("right", true, gsx.Raw("x"), nil))
	if strings.Contains(got, `aria-label="Close"`) {
		t.Errorf("hideCloseButton must omit the X button\nin: %s", got)
	}
}

// TestSheetContentDisplayGated proves the base class gates display behind
// the open: variant (this port's ADAPT, see ui/sheet.gsx's own doc
// comment) rather than shipping shadcn's unscoped `flex`, which would
// defeat the UA's closed-dialog display:none.
func TestSheetContentDisplayGated(t *testing.T) {
	got := render(t, ui.SheetContent("", true, gsx.Raw("x"), nil))
	if !strings.Contains(got, "open:flex") {
		t.Errorf("want open:flex (display-gated), got: %s", got)
	}
	// The bare, unscoped "flex" token must not appear on its own (flex-col
	// contains "flex" as a substring, so check for the exact space-bounded
	// token instead of a raw Contains).
	for _, tok := range strings.Fields(got) {
		if tok == "flex" {
			t.Errorf("unscoped flex token present (defeats closed-state display:none)\nin: %s", got)
		}
	}
}

func TestSheetContentAttrsFallThrough(t *testing.T) {
	got := render(t, ui.SheetContent("", true, gsx.Raw("x"), gsx.Attrs{{Key: "id", Value: "s1"}, {Key: "class", Value: "max-w-md"}}))
	for _, want := range []string{`id="s1"`, "max-w-md"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestSheetHeaderPinned(t *testing.T) {
	got := render(t, ui.SheetHeader(gsx.Raw("x"), nil))
	want := `<div data-slot="sheet-header" class="flex flex-col gap-1.5 p-4">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestSheetFooterPinned(t *testing.T) {
	got := render(t, ui.SheetFooter(gsx.Raw("x"), nil))
	want := `<div data-slot="sheet-footer" class="mt-auto flex flex-col gap-2 p-4">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestSheetTitlePinned(t *testing.T) {
	got := render(t, ui.SheetTitle(gsx.Raw("Edit profile"), nil))
	want := `<h2 data-slot="sheet-title" class="font-semibold text-foreground">Edit profile</h2>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestSheetTitleCallerClassMerges proves the fallthrough class merge path
// (representative pin per the task's test requirements): a caller class
// conflicting on the `font-semibold` utility must win.
func TestSheetTitleCallerClassMerges(t *testing.T) {
	got := render(t, ui.SheetTitle(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "font-bold"}}))
	if strings.Contains(got, "font-semibold") {
		t.Errorf("base font-semibold should be dropped by caller font-bold\nin: %s", got)
	}
	for _, want := range []string{"font-bold", "text-foreground"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestSheetDescriptionPinned(t *testing.T) {
	got := render(t, ui.SheetDescription(gsx.Raw("Make changes to your profile here."), nil))
	want := `<p data-slot="sheet-description" class="text-sm text-muted-foreground">Make changes to your profile here.</p>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestSheetClosePinned(t *testing.T) {
	got := render(t, ui.SheetClose(gsx.Raw("Save changes"), nil))
	want := `<button data-slot="sheet-close" data-gsxui-dialog-close type="button">Save changes</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

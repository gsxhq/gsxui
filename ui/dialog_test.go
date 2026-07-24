package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestDialogStructure(t *testing.T) {
	got := render(t, ui.Dialog(
		gsx.Fragment(
			ui.DialogTrigger(gsx.Raw("Open"), nil),
			ui.DialogContent(false, gsx.Fragment(
				ui.DialogHeader(gsx.Fragment(
					ui.DialogTitle(gsx.Raw("Are you sure?"), nil),
					ui.DialogDescription(gsx.Raw("This cannot be undone."), nil),
				), nil),
				ui.DialogFooter(false, ui.DialogClose(gsx.Raw("Cancel"), nil), nil),
			), nil),
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-dialog`,         // root hook
		`class="contents"`,          // root is layout-neutral
		`data-gsxui-dialog-trigger`, // trigger hook
		`aria-haspopup="dialog"`,    // trigger a11y: server-rendered
		`aria-expanded="false"`,     // trigger a11y: initial state
		"<dialog",                   // native element
		`data-state="closed"`,       // server-rendered initial state
		`data-slot="dialog-content"`,
		`data-slot="dialog-title"`, ">Are you sure?<",
		`data-slot="dialog-description"`,
		`data-gsxui-dialog-close`, // both DialogClose and the X button
		`aria-label="Close"`,      // the injected X close button
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestDialogPinned(t *testing.T) {
	// Exact full-render pin for DialogContent(false, ...) with a title child,
	// verified token-by-token against shadcn's DialogContent + close-button
	// classes (registry/new-york-v4/ui/dialog.tsx) and docs/jsx-parity.md's
	// ADAPT lines: text-foreground added, close button drops the
	// data-[state=open]: pair (dialog stamps state, not the close button),
	// and the Overlay's bg-black/50 is folded into backdrop:bg-black/30 backdrop:duration-200 supports-backdrop-filter:backdrop:backdrop-blur-sm data-[state=open]:backdrop:animate-in data-[state=open]:backdrop:fade-in-0 data-[state=closed]:backdrop:animate-out data-[state=closed]:backdrop:fade-out-0 on
	// the native ::backdrop (Portal/Overlay replaced by the top layer).
	got := render(t, ui.DialogContent(false, ui.DialogTitle(gsx.Raw("Title"), nil), nil))
	want := `<dialog data-slot="dialog-content" data-gsxui-dialog-content data-state="closed" class="fixed top-[50%] left-[50%] z-50 open:grid w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] gap-4 rounded-xl border bg-background p-4 text-sm text-foreground duration-200 outline-none sm:max-w-sm data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95 backdrop:bg-black/10 backdrop:duration-200 supports-backdrop-filter:backdrop:backdrop-blur-xs data-[state=open]:backdrop:animate-in data-[state=open]:backdrop:fade-in-0 data-[state=closed]:backdrop:animate-out data-[state=closed]:backdrop:fade-out-0"><h2 data-slot="dialog-title" class="text-base leading-none font-medium">Title</h2><button type="button" data-slot="dialog-close" data-gsxui-dialog-close aria-label="Close" class="absolute top-2 right-2 rounded-xs opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:ring-2 focus:ring-ring focus:ring-offset-2 focus:outline-hidden disabled:pointer-events-none [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 6 6 18"></path><path d="m6 6 12 12"></path></svg></button></dialog>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestDialogHideCloseButton(t *testing.T) {
	got := render(t, ui.DialogContent(true, gsx.Raw("x"), nil))
	if strings.Contains(got, `aria-label="Close"`) {
		t.Errorf("hideCloseButton must omit the X button\nin: %s", got)
	}
}

func TestDialogFooterShowCloseButton(t *testing.T) {
	got := render(t, ui.DialogFooter(true, gsx.Raw("x"), nil))
	for _, want := range []string{
		`data-slot="dialog-footer"`,
		`data-gsxui-dialog-close`,
		`data-variant="outline"`,
		">Close<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

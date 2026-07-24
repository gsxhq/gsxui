package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

// TestAlertDialogStructure is an integration-shaped smoke test composing
// every part together, mirroring TestDialogStructure — it also proves the
// distinguishing contract: data-gsxui-dialog-static + role="alertdialog"
// present, and NO injected close-X button (aria-label="Close" absent).
func TestAlertDialogStructure(t *testing.T) {
	got := render(t, ui.AlertDialog(
		gsx.Fragment(
			ui.AlertDialogTrigger(gsx.Raw("Delete"), nil),
			ui.AlertDialogContent(gsx.Fragment(
				ui.AlertDialogHeader(gsx.Fragment(
					ui.AlertDialogTitle(gsx.Raw("Are you absolutely sure?"), nil),
					ui.AlertDialogDescription(gsx.Raw("This action cannot be undone."), nil),
				), nil),
				ui.AlertDialogFooter(gsx.Fragment(
					ui.AlertDialogCancel(gsx.Raw("Cancel"), nil),
					ui.AlertDialogAction(gsx.Raw("Continue"), nil),
				), nil),
			), nil),
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-dialog`,         // root hook (inherited from Dialog)
		`data-slot="alert-dialog"`,  // root slot override
		`data-gsxui-dialog-trigger`, // trigger hook
		`data-slot="alert-dialog-trigger"`,
		`aria-haspopup="dialog"`,
		`aria-expanded="false"`,
		"<dialog",
		`data-gsxui-dialog-content`, // dialog.js still selects on this
		`data-gsxui-dialog-static`,  // the opt-out this task adds
		`role="alertdialog"`,
		`data-slot="alert-dialog-content"`,
		`data-slot="alert-dialog-header"`,
		`data-slot="alert-dialog-title"`, ">Are you absolutely sure?<",
		`data-slot="alert-dialog-description"`,
		`data-slot="alert-dialog-footer"`,
		`data-slot="alert-dialog-cancel"`, `data-variant="outline"`, ">Cancel<",
		`data-slot="alert-dialog-action"`, `data-variant="default"`, ">Continue<",
		`data-gsxui-dialog-close`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if strings.Contains(got, `aria-label="Close"`) {
		t.Errorf("AlertDialogContent must never render the injected X close button\nin: %s", got)
	}
}

func TestAlertDialogRootPinned(t *testing.T) {
	got := render(t, ui.AlertDialog(gsx.Raw("x"), nil))
	want := `<div data-gsxui-dialog class="contents" data-slot="alert-dialog">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestAlertDialogTriggerPinned(t *testing.T) {
	got := render(t, ui.AlertDialogTrigger(gsx.Raw("Delete"), nil))
	want := `<button data-slot="alert-dialog-trigger" data-gsxui-dialog-trigger type="button" aria-haspopup="dialog" aria-expanded="false">Delete</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestAlertDialogContentPinned is the exact full-render pin for
// AlertDialogContent — see the component's own doc comment in
// ui/alert-dialog.gsx for why the class string ends up byte-identical to
// DialogContent's own default (every alert-dialog.tsx content token is
// already present there; the one that isn't, a bare `grid`, is dropped as
// an ADAPT). hideCloseButton is always true here (no injected X button,
// unlike DialogContent's own TestDialogPinned) and data-gsxui-dialog-static
// + role="alertdialog" are the two attrs that distinguish this from a plain
// DialogContent render.
func TestAlertDialogContentPinned(t *testing.T) {
	got := render(t, ui.AlertDialogContent(gsx.Raw("x"), nil))
	want := `<dialog data-gsxui-dialog-content data-state="closed" class="fixed top-[50%] left-[50%] z-50 open:grid w-full translate-x-[-50%] translate-y-[-50%] gap-4 rounded-xl border bg-background p-4 text-sm text-foreground duration-200 outline-none data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95 backdrop:bg-black/10 backdrop:duration-200 supports-backdrop-filter:backdrop:backdrop-blur-xs data-[state=open]:backdrop:animate-in data-[state=open]:backdrop:fade-in-0 data-[state=closed]:backdrop:animate-out data-[state=closed]:backdrop:fade-out-0 max-w-xs sm:max-w-sm" data-slot="alert-dialog-content" role="alertdialog" data-gsxui-dialog-static="true">x</dialog>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestAlertDialogContentNoCloseButton proves hideCloseButton is hardcoded
// true — there is no param to turn the injected X back on, unlike
// DialogContent's own hideCloseButton bool.
func TestAlertDialogContentNoCloseButton(t *testing.T) {
	got := render(t, ui.AlertDialogContent(gsx.Raw("x"), nil))
	if strings.Contains(got, `aria-label="Close"`) {
		t.Errorf("AlertDialogContent must never render the injected X button\nin: %s", got)
	}
}

func TestAlertDialogContentAttrsFallThrough(t *testing.T) {
	got := render(t, ui.AlertDialogContent(gsx.Raw("x"), gsx.Attrs{{Key: "id", Value: "ad1"}, {Key: "class", Value: "max-w-md"}}))
	for _, want := range []string{`id="ad1"`, "max-w-md"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// TestAlertDialogHeaderPinned pins the current-upstream unconditional grid
// recipe (registry/new-york-v4/ui/alert-dialog.tsx's AlertDialogHeader
// base, HEAD f31ed8198) — always centered, two grid rows — not the
// pre-refactor flex/sm:text-left shape; see ui/alert-dialog.gsx's own doc
// comment on AlertDialogHeader for why this base is not scoped out along
// with the size/Media-conditional tokens.
func TestAlertDialogHeaderPinned(t *testing.T) {
	got := render(t, ui.AlertDialogHeader(gsx.Raw("x"), nil))
	want := `<div data-slot="alert-dialog-header" class="grid grid-rows-[auto_1fr] place-items-center gap-1.5 text-center">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestAlertDialogFooterPinned(t *testing.T) {
	got := render(t, ui.AlertDialogFooter(gsx.Raw("x"), nil))
	want := `<div data-slot="alert-dialog-footer" class="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end bg-muted/50 -mx-4 -mb-4 rounded-b-xl border-t p-4">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestAlertDialogTitlePinned(t *testing.T) {
	got := render(t, ui.AlertDialogTitle(gsx.Raw("Are you absolutely sure?"), nil))
	want := `<h2 data-slot="alert-dialog-title" class="text-base font-medium">Are you absolutely sure?</h2>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestAlertDialogTitleCallerClassMerges proves the fallthrough class merge
// path (representative pin per the task's test requirements): a caller
// class conflicting on the `text-base` font-size utility must win.
func TestAlertDialogTitleCallerClassMerges(t *testing.T) {
	got := render(t, ui.AlertDialogTitle(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "text-xl"}}))
	if strings.Contains(got, "text-base") {
		t.Errorf("base text-base should be dropped by caller text-xl\nin: %s", got)
	}
	for _, want := range []string{"text-xl", "font-medium"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// class token order (text-sm text-muted-foreground) matches current
// upstream exactly — see ui/alert-dialog.gsx's own comment on
// AlertDialogDescription.
func TestAlertDialogDescriptionPinned(t *testing.T) {
	got := render(t, ui.AlertDialogDescription(gsx.Raw("This action cannot be undone."), nil))
	want := `<p data-slot="alert-dialog-description" class="text-sm text-muted-foreground">This action cannot be undone.</p>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestAlertDialogActionPinned proves AlertDialogAction composes ui.Button
// at its default variant/size (shadcn's own buttonVariants() default) plus
// the close-wiring data attribute.
func TestAlertDialogActionPinned(t *testing.T) {
	got := render(t, ui.AlertDialogAction(gsx.Raw("Continue"), nil))
	want := `<button data-variant="default" data-size="default" type="button" class="inline-flex shrink-0 items-center justify-center rounded-lg border border-transparent bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 bg-primary text-primary-foreground hover:bg-primary/90 h-8 gap-1.5 px-2.5 has-[&gt;svg]:px-2" data-slot="alert-dialog-action" data-gsxui-dialog-close="true">Continue</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestAlertDialogCancelPinned proves AlertDialogCancel composes ui.Button
// with variant="outline" (shadcn's own buttonVariants({variant:"outline"}))
// plus the close-wiring data attribute — the button-variant/close-wiring
// half of this task's "Action/Cancel button variants + close wiring" pin
// requirement (Action's own pin above covers the default-variant half).
func TestAlertDialogCancelPinned(t *testing.T) {
	got := render(t, ui.AlertDialogCancel(gsx.Raw("Cancel"), nil))
	for _, want := range []string{
		`data-variant="outline"`,
		"border bg-background", // Button's outline variant classes
		`data-slot="alert-dialog-cancel"`,
		`data-gsxui-dialog-close="true"`,
		">Cancel<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

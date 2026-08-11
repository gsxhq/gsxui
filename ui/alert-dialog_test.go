package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/merge"
	"github.com/gsxhq/gsxui/ui"
)

// alertDialogContentClass is AlertDialogContent's exact resolved class
// attribute: AlertDialog's own content utilities (copied verbatim from
// ui/alert-dialog.gsx's AlertDialogContent, which now narrows Dialog's box
// via data-[size=...] variants rather than a plain max-w-* override) merged
// into <ui.DialogContent>'s, by the same merge.Merge the renderer uses.
// Derived rather than pasted so the assertion keeps proving the mechanism —
// the narrowing works because merge REPLACES Dialog's max-w-[calc(...)] with
// AlertDialog's own gap-4/rounded-xl/p-4/etc. and size-scoped max-w-xs
// utilities rather than appending to it.
func alertDialogContentClass() string {
	return merge.Merge([]string{dialogContentClass(), "data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 bg-popover text-popover-foreground ring-foreground/10 gap-4 rounded-xl p-4 ring-1 duration-100 data-[size=default]:max-w-xs data-[size=sm]:max-w-xs data-[size=default]:sm:max-w-sm"})
}

func TestAlertDialogPinnedParts(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"root", render(t, ui.AlertDialog(gsx.Raw("x"), nil)), `<div class="contents" data-gsxui-slot-alert-dialog data-gsxui-slot-dialog>x</div>`},
		{"trigger", render(t, ui.AlertDialogTrigger(gsx.Raw("Delete"), nil)), `<button type="button" aria-haspopup="dialog" aria-expanded="false" data-gsxui-slot-alert-dialog-trigger>Delete</button>`},
		{"content", render(t, ui.AlertDialogContent(gsx.Raw("x"), nil)), `<dialog class="` + alertDialogContentClass() + `" data-state="closed" role="alertdialog" data-gsxui-dialog-static data-gsxui-slot-alert-dialog-content data-gsxui-slot-dialog-content>x</dialog>`},
		{"header", render(t, ui.AlertDialogHeader(gsx.Raw("x"), nil)), `<div class="grid grid-rows-[auto_1fr] place-items-center gap-1.5 text-center sm:group-data-[size=default]/alert-dialog-content:place-items-start sm:group-data-[size=default]/alert-dialog-content:text-left" data-gsxui-slot-alert-dialog-header>x</div>`},
		{"footer", render(t, ui.AlertDialogFooter(gsx.Raw("x"), nil)), `<div class="bg-muted/50 -mx-4 -mb-4 rounded-b-xl border-t p-4 flex flex-col-reverse" data-gsxui-slot-alert-dialog-footer>x</div>`},
		{"title", render(t, ui.AlertDialogTitle(gsx.Raw("x"), nil)), `<h2 class="text-base font-medium" data-gsxui-slot-alert-dialog-title data-gsxui-slot-dialog-title>x</h2>`},
		{"description", render(t, ui.AlertDialogDescription(gsx.Raw("x"), nil)), `<p class="text-muted-foreground *:[a]:hover:text-foreground text-sm text-balance md:text-pretty *:[a]:underline *:[a]:underline-offset-3" data-gsxui-slot-alert-dialog-description data-gsxui-slot-dialog-description>x</p>`},
		{"action", render(t, ui.AlertDialogAction(gsx.Raw("x"), nil)), `<button data-variant="default" data-size="default" type="button" ` + canonicalButtonClass("default", "default") + ` data-gsxui-dialog-close data-gsxui-slot-alert-dialog-action data-gsxui-slot-button>x</button>`},
		{"cancel", render(t, ui.AlertDialogCancel(gsx.Raw("x"), nil)), `<button data-variant="outline" data-size="default" type="button" ` + canonicalButtonClass("outline", "default") + ` data-gsxui-dialog-close data-gsxui-slot-alert-dialog-cancel data-gsxui-slot-button>x</button>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("pinned render mismatch\n got: %s\nwant: %s", tt.got, tt.want)
			}
		})
	}
}

func TestAlertDialogContentNoCloseButton(t *testing.T) {
	got := render(t, ui.AlertDialogContent(gsx.Raw("x"), nil))
	if strings.Contains(got, "dialog-close") {
		t.Errorf("alert dialog must not inject a close button\nin: %s", got)
	}
}

func TestAlertDialogCallerAttrsFallThroughOnce(t *testing.T) {
	got := render(t, ui.AlertDialogContent(gsx.Raw("x"), gsx.Attrs{
		{Key: "id", Value: "confirm"},
		{Key: "class", Value: "caller"},
	}))
	if strings.Count(got, `id="confirm"`) != 1 {
		t.Errorf("id must render once\nin: %s", got)
	}
	if !strings.Contains(got, "caller") || strings.Count(got, "class=") != 1 {
		t.Errorf("caller class must be merged into the one class attribute\nin: %s", got)
	}
}

func TestAlertDialogActionComposesPresenceMarkersOnButton(t *testing.T) {
	got := render(t, ui.AlertDialogAction(gsx.Raw("Delete"), nil))
	requirePresenceAttributesOnSameTag(t, got, "data-gsxui-dialog-close",
		"data-gsxui-slot-alert-dialog-action",
		"data-gsxui-slot-button",
	)
}

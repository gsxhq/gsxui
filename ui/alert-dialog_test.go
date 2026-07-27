package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestAlertDialogPinnedParts(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"root", render(t, ui.AlertDialog(gsx.Raw("x"), nil)), `<div data-gsxui-dialog data-gsxui-slot-alert-dialog data-gsxui-slot-dialog>x</div>`},
		{"trigger", render(t, ui.AlertDialogTrigger(gsx.Raw("Delete"), nil)), `<button data-gsxui-dialog-trigger type="button" aria-haspopup="dialog" aria-expanded="false" data-gsxui-slot-alert-dialog-trigger>Delete</button>`},
		{"content", render(t, ui.AlertDialogContent(gsx.Raw("x"), nil)), `<dialog data-gsxui-dialog-content data-state="closed" role="alertdialog" data-gsxui-dialog-static data-gsxui-slot-alert-dialog-content data-gsxui-slot-dialog-content>x</dialog>`},
		{"header", render(t, ui.AlertDialogHeader(gsx.Raw("x"), nil)), `<div data-gsxui-slot-alert-dialog-header>x</div>`},
		{"footer", render(t, ui.AlertDialogFooter(gsx.Raw("x"), nil)), `<div data-gsxui-slot-alert-dialog-footer>x</div>`},
		{"title", render(t, ui.AlertDialogTitle(gsx.Raw("x"), nil)), `<h2 data-gsxui-dialog-title data-gsxui-slot-alert-dialog-title>x</h2>`},
		{"description", render(t, ui.AlertDialogDescription(gsx.Raw("x"), nil)), `<p data-gsxui-dialog-description data-gsxui-slot-alert-dialog-description>x</p>`},
		{"action", render(t, ui.AlertDialogAction(gsx.Raw("x"), nil)), `<button data-variant="default" data-size="default" type="button" data-gsxui-dialog-close data-gsxui-slot-alert-dialog-action data-gsxui-slot-button>x</button>`},
		{"cancel", render(t, ui.AlertDialogCancel(gsx.Raw("x"), nil)), `<button data-variant="outline" data-size="default" type="button" data-gsxui-dialog-close data-gsxui-slot-alert-dialog-cancel data-gsxui-slot-button>x</button>`},
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
	for _, want := range []string{`id="confirm"`, `class="caller"`} {
		if strings.Count(got, want) != 1 {
			t.Errorf("%q must render once\nin: %s", want, got)
		}
	}
}

func TestAlertDialogActionComposesPresenceMarkersOnButton(t *testing.T) {
	got := render(t, ui.AlertDialogAction(gsx.Raw("Delete"), nil))
	requirePresenceAttributesOnSameTag(t, got, "data-gsxui-dialog-close",
		"data-gsxui-slot-alert-dialog-action",
		"data-gsxui-slot-button",
	)
}

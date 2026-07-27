package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestDialogStructure(t *testing.T) {
	got := render(t, ui.Dialog(gsx.Fragment(
		ui.DialogTrigger(gsx.Raw("Open"), nil),
		ui.DialogContent(false, gsx.Fragment(
			ui.DialogTitle(gsx.Raw("Title"), nil),
			ui.DialogDescription(gsx.Raw("Description"), nil),
		), nil),
	), nil))
	for _, want := range []string{
		`data-gsxui-slot="dialog"`,
		`data-gsxui-slot="dialog-trigger"`,
		`aria-haspopup="dialog"`,
		`aria-expanded="false"`,
		`data-gsxui-slot="dialog-content"`,
		`data-state="closed"`,
		`data-gsxui-dialog-title`,
		`data-gsxui-dialog-description`,
		`data-gsxui-slot="dialog-close dialog-close-button"`,
		`data-gsxui-slot="dialog-close-icon"`,
		`data-gsxui-slot="dialog-close-label"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if strings.Contains(got, ` class=`) || strings.Contains(got, `data-slot=`) {
		t.Errorf("default dialog markup must contain no presentation classes or legacy slots\nin: %s", got)
	}
}

func TestDialogHideCloseButton(t *testing.T) {
	got := render(t, ui.DialogContent(true, gsx.Raw("x"), nil))
	if strings.Contains(got, "dialog-close") {
		t.Errorf("hideCloseButton must omit injected close parts\nin: %s", got)
	}
}

func TestDialogFooterShowCloseButton(t *testing.T) {
	got := render(t, ui.DialogFooter(true, gsx.Raw("x"), nil))
	for _, want := range []string{
		`data-gsxui-slot="dialog-footer"`,
		`data-gsxui-slot="button dialog-footer-close"`,
		`data-variant="outline"`,
		`data-gsxui-dialog-close`,
		">Close<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestDialogCallerClassOnly(t *testing.T) {
	got := render(t, ui.DialogContent(true, gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "caller"}}))
	if strings.Count(got, `class="caller"`) != 1 || strings.Count(got, `class=`) != 1 {
		t.Errorf("caller class must be the only class\nin: %s", got)
	}
}

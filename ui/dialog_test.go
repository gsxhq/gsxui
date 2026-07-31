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
		`data-gsxui-slot-dialog`,
		`data-gsxui-slot-dialog-trigger`,
		`aria-haspopup="dialog"`,
		`aria-expanded="false"`,
		`data-gsxui-slot-dialog-content`,
		`data-state="closed"`,
		`data-gsxui-dialog-title`,
		`data-gsxui-dialog-description`,
		`data-gsxui-slot-dialog-close-button data-gsxui-slot-dialog-close`,
		`data-gsxui-slot-dialog-close-icon`,
		`data-gsxui-slot-dialog-close-label`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	// Dialog is migrated to the slot axis, so its parts now render their own
	// resolved recipe classes (presentation lives in the stylesheet; this
	// structural pin only needs to confirm no LEGACY slot attribute leaked
	// through).
	if strings.Contains(got, `data-slot=`) {
		t.Errorf("default dialog markup must contain no legacy slots\nin: %s", got)
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
		`data-gsxui-slot-dialog-footer`,
		`data-gsxui-slot-dialog-footer-close data-gsxui-slot-button`,
		`data-variant="outline"`,
		`data-gsxui-dialog-close`,
		">Close<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// dialogContentClass is DialogContent's exact resolved class attribute value
// (no caller classes), for tests composing <ui.DialogContent> (AlertDialog,
// Command) that need to assert the composed markup verbatim.
func dialogContentClass() string {
	return "fixed z-50 text-sm duration-200 outline-none data-[state=closed]:animate-out data-[state=open]:animate-in backdrop:bg-[var(--overlay)] backdrop:backdrop-blur-xs backdrop:duration-200 data-[state=open]:backdrop:animate-in data-[state=open]:backdrop:fade-in-0 data-[state=closed]:backdrop:animate-out data-[state=closed]:backdrop:fade-out-0 top-1/2 left-1/2 w-full max-w-[calc(100%-2rem)] -translate-x-1/2 -translate-y-1/2 gap-4 rounded-xl border bg-background p-4 text-foreground open:grid data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95 sm:max-w-sm"
}

func TestDialogCallerClassOnly(t *testing.T) {
	got := render(t, ui.DialogContent(true, gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "caller"}}))
	if !strings.Contains(got, "caller") || strings.Count(got, `class=`) != 1 {
		t.Errorf("caller class must be merged into the one class attribute\nin: %s", got)
	}
}

package dialog_test

import (
	"context"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/dialog"
)

func render(t *testing.T, n gsx.Node) string {
	t.Helper()
	var sb strings.Builder
	if err := n.Render(context.Background(), &sb); err != nil {
		t.Fatal(err)
	}
	return sb.String()
}

func TestDialogStructure(t *testing.T) {
	got := render(t, dialog.Dialog(
		gsx.Fragment(
			dialog.DialogTrigger(gsx.Raw("Open"), nil),
			dialog.DialogContent(false, gsx.Fragment(
				dialog.DialogHeader(gsx.Fragment(
					dialog.DialogTitle(gsx.Raw("Are you sure?"), nil),
					dialog.DialogDescription(gsx.Raw("This cannot be undone."), nil),
				), nil),
				dialog.DialogFooter(dialog.DialogClose(gsx.Raw("Cancel"), nil), nil),
			), nil),
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-dialog`,         // root hook
		`class="contents"`,          // root is layout-neutral
		`data-gsxui-dialog-trigger`, // trigger hook
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

func TestDialogHideCloseButton(t *testing.T) {
	got := render(t, dialog.DialogContent(true, gsx.Raw("x"), nil))
	if strings.Contains(got, `aria-label="Close"`) {
		t.Errorf("hideCloseButton must omit the X button\nin: %s", got)
	}
}

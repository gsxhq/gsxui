package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

// TestToasterPinned pins the ONLY server-rendered surface of the sonner
// port: the always-present bottom-right region every toast <li> is
// client-constructed into (see ui/sonner.js). The toast markup itself has
// no Go render path — it is built by JS — so there is nothing else to pin
// here (command.js's precedent: Go-pin only the server surface).
func TestToasterPinned(t *testing.T) {
	got := render(t, ui.Toaster(nil))
	want := `<section aria-label="Notifications" tabindex="-1"><ol data-slot="toaster" data-gsxui-toaster class="pointer-events-none fixed z-100 flex flex-col gap-2 p-6 bottom-0 right-0"></ol></section>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// Attrs merge onto the <ol> (a caller can retarget position or add ids),
// same passthrough contract every other component's root carries.
func TestToasterAttrsMerge(t *testing.T) {
	got := render(t, ui.Toaster(gsx.Attrs{{Key: "id", Value: "my-toaster"}}))
	if !strings.Contains(got, `id="my-toaster"`) {
		t.Errorf("caller id missing\nin: %s", got)
	}
	if !strings.Contains(got, `data-gsxui-toaster`) {
		t.Errorf("data-gsxui-toaster hook missing\nin: %s", got)
	}
}

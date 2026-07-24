package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestCommandPinned(t *testing.T) {
	got := render(t, ui.Command(gsx.Raw("x"), nil))
	want := `<div data-slot="command" data-gsxui-command class="flex h-full w-full flex-col overflow-hidden rounded-xl bg-popover p-1 text-popover-foreground">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCommandInputPinned(t *testing.T) {
	got := render(t, ui.CommandInput("Search...", nil))
	for _, want := range []string{
		`<div data-slot="command-input-wrapper" class="flex h-9 items-center gap-2 border-b px-3">`,
		`data-slot="command-input"`,
		`data-gsxui-command-input`,
		`role="combobox"`,
		`aria-autocomplete="list"`,
		`placeholder="Search..."`,
		"flex h-10 w-full rounded-md bg-transparent py-3 text-sm outline-hidden placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestCommandItemPinned(t *testing.T) {
	got := render(t, ui.CommandItem("calendar", gsx.Raw("Calendar"), nil))
	for _, want := range []string{
		`data-slot="command-item"`,
		`data-gsxui-command-item`,
		`data-value="calendar"`,
		`role="option"`,
		`aria-selected="false"`,
		"data-[selected=true]:bg-accent data-[selected=true]:text-accent-foreground",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// CommandGroup renders its heading child only when non-empty.
func TestCommandGroupHeading(t *testing.T) {
	with := render(t, ui.CommandGroup("Settings", gsx.Raw("x"), nil))
	if !strings.Contains(with, `<div data-slot="command-group-heading" class="px-2 py-1.5 text-xs font-medium text-muted-foreground">Settings</div>`) {
		t.Errorf("missing heading in: %s", with)
	}
	without := render(t, ui.CommandGroup("", gsx.Raw("x"), nil))
	if strings.Contains(without, "command-group-heading") {
		t.Errorf("empty heading should render no heading element: %s", without)
	}
}

// CommandEmpty is server-rendered hidden — command.js reveals it when a
// query matches nothing.
func TestCommandEmptyHidden(t *testing.T) {
	got := render(t, ui.CommandEmpty(gsx.Raw("No results."), nil))
	if !strings.Contains(got, " hidden") {
		t.Errorf("CommandEmpty must render hidden: %s", got)
	}
}

// CommandDialog composes DialogContent (command → dialog dep) with the p-0
// override winning the merge, the sr-only a11y header inside the dialog for
// wireA11y, and the hotkey hook attribute.
func TestCommandDialogComposition(t *testing.T) {
	got := render(t, ui.CommandDialog("", "", gsx.Raw("x"), nil))
	for _, want := range []string{
		`data-gsxui-command-dialog`,
		`data-gsxui-dialog-content`,
		`<div data-slot="dialog-header" class="flex flex-col gap-2 text-center sm:text-left sr-only">`,
		">Command Palette</h2>",
		">Search for a command to run...</p>",
		`data-slot="command"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if strings.Contains(got, " p-4 ") {
		t.Errorf("dialog base p-4 must lose to CommandDialog's p-0: %s", got)
	}
	if !strings.Contains(got, " p-0") {
		t.Errorf("CommandDialog's p-0 override missing: %s", got)
	}
}

func TestCommandCallerClassMerges(t *testing.T) {
	got := render(t, ui.Command(nil, gsx.Attrs{{Key: "class", Value: "rounded-lg"}}))
	if strings.Contains(got, "rounded-xl") {
		t.Errorf("base rounded-xl should be dropped by caller rounded-lg\nin: %s", got)
	}
	if !strings.Contains(got, "rounded-lg") {
		t.Errorf("caller rounded-lg missing\nin: %s", got)
	}
}

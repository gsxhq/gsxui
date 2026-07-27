package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestCommandPinned(t *testing.T) {
	got := render(t, ui.Command(gsx.Raw("x"), nil))
	want := `<div data-gsxui-command data-gsxui-slot-command>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCommandInputPinned(t *testing.T) {
	got := render(t, ui.CommandInput("Search...", nil))
	for _, want := range []string{
		`<div data-gsxui-command-input-wrapper data-gsxui-slot-command-input-wrapper>`,
		`data-gsxui-slot-command-input`,
		`data-gsxui-command-input`,
		`role="combobox"`,
		`aria-autocomplete="list"`,
		`placeholder="Search..."`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestCommandItemPinned(t *testing.T) {
	got := render(t, ui.CommandItem("calendar", gsx.Raw("Calendar"), nil))
	for _, want := range []string{
		`data-gsxui-slot-command-item`,
		`data-gsxui-command-item`,
		`data-value="calendar"`,
		`role="option"`,
		`aria-selected="false"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// CommandGroup renders its heading child only when non-empty.
func TestCommandGroupHeading(t *testing.T) {
	with := render(t, ui.CommandGroup("Settings", gsx.Raw("x"), nil))
	if !strings.Contains(with, `<div data-gsxui-command-group-heading data-gsxui-slot-command-group-heading>Settings</div>`) {
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

// CommandDialog composes DialogContent and Command through ordered styling
// tokens, while its a11y header remains visually hidden by its semantic token.
func TestCommandDialogComposition(t *testing.T) {
	got := render(t, ui.CommandDialog("", "", ui.CommandInput("Search", nil), nil))
	for _, want := range []string{
		`data-gsxui-command-dialog`,
		`data-gsxui-dialog-content`,
		`data-gsxui-slot-command-dialog data-gsxui-slot-dialog`,
		`data-gsxui-slot-command-dialog-content data-gsxui-slot-dialog-content`,
		`data-gsxui-slot-command-dialog-header data-gsxui-slot-dialog-header`,
		">Command Palette</h2>",
		">Search for a command to run...</p>",
		`data-gsxui-slot-command-dialog-command data-gsxui-slot-command`,
		`data-gsxui-slot-command-input-wrapper`,
		`data-gsxui-slot-command-input`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if strings.Contains(got, ` class=`) {
		t.Errorf("CommandDialog must not pass an internal presentation class\nin: %s", got)
	}
}

func TestCommandCallerClassMerges(t *testing.T) {
	got := render(t, ui.Command(nil, gsx.Attrs{{Key: "class", Value: "rounded-lg"}}))
	if strings.Count(got, `class="rounded-lg"`) != 1 || strings.Count(got, `class=`) != 1 {
		t.Errorf("caller rounded-lg must be the only class and render once\nin: %s", got)
	}
}

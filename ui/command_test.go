package ui_test

import (
	"html"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/internal/stylegen"
	"github.com/gsxhq/gsxui/merge"
	"github.com/gsxhq/gsxui/ui"
)

// novaCommandRecipe loads the default style's Command recipe once. ui.Command
// and its parts are generated output now, so the classes they render are the
// default style's concrete utilities — same pattern as novaKbdRecipe in
// kbd_test.go.
var novaCommandRecipe = sync.OnceValue(func() recipe.Style {
	path := filepath.Join("..", "registry", "styles", stylegen.DefaultStyle, "command.css")
	src, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	style, err := recipe.ParseStyle(path, src)
	if err != nil {
		panic(err)
	}
	return style
})

// canonicalCommandClass renders the class attribute a Command slot emits,
// with any caller classes merged in the same order the generated .gsx does.
func canonicalCommandClass(slot string, caller ...string) string {
	class := "gsxui-recipe-command"
	if slot != "" {
		class += "-" + slot
	}
	rule, ok := novaCommandRecipe().Lookup(class)
	if !ok {
		panic("default style declares no recipe " + class)
	}
	classes := append([]string(nil), rule.Utilities...)
	classes = append(classes, caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

func TestCommandPinned(t *testing.T) {
	got := render(t, ui.Command(gsx.Raw("x"), nil))
	want := `<div data-gsxui-command ` + canonicalCommandClass("") + ` data-gsxui-slot-command>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCommandInputPinned(t *testing.T) {
	got := render(t, ui.CommandInput("Search...", nil))
	for _, want := range []string{
		`<div data-gsxui-command-input-wrapper ` + canonicalCommandClass("input-wrapper") + ` data-gsxui-slot-command-input-wrapper>`,
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
	if !strings.Contains(with, `<div data-gsxui-command-group-heading `+canonicalCommandClass("group-heading")+` data-gsxui-slot-command-group-heading>Settings</div>`) {
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
	got := render(t, ui.CommandDialog(
		"",
		"",
		ui.DialogTrigger(gsx.Text("Open"), nil),
		ui.CommandInput("Search", nil),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-command-dialog`,
		`data-gsxui-dialog-content`,
		`data-gsxui-dialog-trigger`,
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
	// Dialog/DialogContent/DialogHeader are migrated to the slot axis, so
	// CommandDialog's composed markup now legitimately carries their
	// resolved recipe classes (Command itself still adds none of its own).
}

// Command is migrated to the slot axis, so a caller's class no longer IS the
// whole attribute — it merges into the recipe's own resolved utilities, and
// tailwind-merge drops the recipe's rounded-xl in favour of the caller's
// rounded-lg. The pin asserts that exact merged result rather than the
// pre-migration exact-attribute match (migration playbook Step 9).
func TestCommandCallerClassMerges(t *testing.T) {
	got := render(t, ui.Command(nil, gsx.Attrs{{Key: "class", Value: "rounded-lg"}}))
	want := canonicalCommandClass("", "rounded-lg")
	if strings.Count(got, want) != 1 || strings.Count(got, `class=`) != 1 {
		t.Errorf("caller rounded-lg must merge into the recipe class exactly once\nwant: %s\nin: %s", want, got)
	}
	if !strings.Contains(got, "rounded-lg") || strings.Contains(got, "rounded-xl") {
		t.Errorf("caller radius must replace the recipe's own\nin: %s", got)
	}
}

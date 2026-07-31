package ui_test

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/internal/stylegen"
	"github.com/gsxhq/gsxui/ui"
)

func assertMenuCSSOnlyMarkup(t *testing.T, got string, slots ...string) {
	t.Helper()
	if strings.Contains(got, `data-slot=`) {
		t.Errorf("legacy data-slot must not render\nin: %s", got)
	}
	if strings.Contains(got, ` class="`) {
		t.Errorf("built-in presentation class must not render\nin: %s", got)
	}
	for _, slot := range slots {
		for name := range strings.FieldsSeq(slot) {
			if !strings.Contains(got, `data-gsxui-slot-`+name) {
				t.Errorf("missing slot marker %q\nin: %s", name, got)
			}
		}
	}
}

// novaComboboxRecipe loads the default style's Combobox recipe, the same way
// novaButtonRecipe loads Button's: Combobox is migrated to the slot axis, so
// its trigger and clear render Button's canonical roles WITH Combobox's own
// slot utilities merged on top as caller classes.
var novaComboboxRecipe = sync.OnceValue(func() recipe.Style {
	path := filepath.Join("..", "registry", "styles", stylegen.DefaultStyle, "combobox.css")
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

func comboboxRecipeUtilities(class string) []string {
	rule, ok := novaComboboxRecipe().Lookup(class)
	if !ok {
		panic("default style declares no recipe " + class)
	}
	return rule.Utilities
}

func assertComboboxCanonicalButtonRoles(t *testing.T, got string) {
	t.Helper()

	// Both buttons compose InputGroupButton -> Button and carry the canonical
	// ghost/icon-xs role. Combobox's own migration adds its trigger and clear
	// slot utilities on top, as ordinary caller classes merged by merge.Merge
	// — which is why the trigger's own
	// [&_svg:not([class*='size-'])]:size-4 REPLACES Button's icon-xs size-3
	// rather than following it, and why the expectation runs the merger
	// instead of concatenating strings.
	for _, want := range []string{
		canonicalButtonClass("ghost", "icon-xs", comboboxRecipeUtilities("gsxui-recipe-combobox-trigger")...),
		canonicalButtonClass("ghost", "icon-xs", comboboxRecipeUtilities("gsxui-recipe-combobox-clear")...),
	} {
		if strings.Count(got, want) != 1 {
			t.Errorf("Combobox trigger and clear must each render exact canonical Button roles\nwant: %s\nin: %s", want, got)
		}
	}
}

func TestComboboxCSSOnlyContract(t *testing.T) {
	got := render(t, ui.Combobox(
		"framework",
		"next",
		gsx.Fragment(
			ui.ComboboxInput("Search", true, true, false, nil, nil),
			ui.ComboboxContent(
				ui.ComboboxList(
					gsx.Fragment(
						ui.ComboboxGroup(
							gsx.Fragment(
								ui.ComboboxLabel(gsx.Raw("Framework"), nil),
								ui.ComboboxItem("next", true, gsx.Raw("Next"), nil),
							),
							nil,
						),
						ui.ComboboxEmpty(gsx.Raw("Empty"), nil),
						ui.ComboboxSeparator(nil),
					),
					nil,
				),
				nil,
			),
		),
		nil,
	))
	assertMenuMarkupSlots(t, got,
		"combobox",
		"combobox-bridge",
		"input-group combobox-input-group",
		"input input-group-control combobox-input",
		"button input-group-button combobox-trigger",
		"icon combobox-trigger-icon",
		"button input-group-button combobox-clear",
		"combobox-content",
		"combobox-list",
		"combobox-group",
		"combobox-label",
		"combobox-item",
		"combobox-item-indicator",
		"combobox-empty",
		"combobox-separator",
	)
	assertComboboxCanonicalButtonRoles(t, got)
	for _, hook := range []string{
		`data-gsxui-combobox-group`,
		`data-gsxui-combobox-label`,
		`data-gsxui-combobox-item`,
		`data-gsxui-combobox-empty`,
		`data-gsxui-combobox-separator`,
	} {
		if !strings.Contains(got, hook) {
			t.Errorf("missing Combobox behavior hook %q\nin: %s", hook, got)
		}
	}
}

func assertMenuMarkupSlots(t *testing.T, got string, slots ...string) {
	t.Helper()
	if strings.Contains(got, `data-slot=`) {
		t.Errorf("legacy data-slot must not render\nin: %s", got)
	}
	for _, slot := range slots {
		for name := range strings.FieldsSeq(slot) {
			if !strings.Contains(got, `data-gsxui-slot-`+name) {
				t.Errorf("missing slot marker %q\nin: %s", name, got)
			}
		}
	}
}

// TestDropdownMarkupContract used to be TestDropdownCSSOnlyContract:
// DropdownMenu migrated to the slot axis (registry/canonical/dropdown-menu.gsx),
// so its markup now legitimately carries a class= attribute (the recipe's
// compiled utilities) — assertMenuCSSOnlyMarkup's "no class=" assertion no
// longer holds, and this switches to assertMenuMarkupSlots (marker presence
// only), the same downgrade Card/Badge/Alert's own migrations made to their
// prior CSS-only pins.
func TestDropdownMarkupContract(t *testing.T) {
	got := render(t, ui.DropdownMenu(
		gsx.Fragment(
			ui.DropdownMenuTrigger(gsx.Raw("Open"), nil),
			ui.DropdownMenuContent(
				gsx.Fragment(
					ui.DropdownMenuGroup(ui.DropdownMenuItem("", gsx.Raw("Item"), nil), nil),
					ui.DropdownMenuCheckboxItem(true, "checked", gsx.Raw("Checked"), nil),
					ui.DropdownMenuRadioGroup("one", ui.DropdownMenuRadioItem(true, "one", gsx.Raw("One"), nil), nil),
					ui.DropdownMenuLabel(gsx.Raw("Label"), nil),
					ui.DropdownMenuSeparator(nil),
					ui.DropdownMenuShortcut(gsx.Raw("⌘K"), nil),
					ui.DropdownMenuSub(
						gsx.Fragment(
							ui.DropdownMenuSubTrigger(gsx.Raw("More"), nil),
							ui.DropdownMenuSubContent(gsx.Raw("Sub"), nil),
						),
						nil,
					),
				),
				nil,
			),
		),
		nil,
	))
	assertMenuMarkupSlots(t, got,
		"dropdown-menu",
		"dropdown-menu-trigger",
		"dropdown-menu-content",
		"dropdown-menu-group",
		"dropdown-menu-item",
		"dropdown-menu-checkbox-item",
		"dropdown-menu-checkbox-item-indicator",
		"dropdown-menu-radio-group",
		"dropdown-menu-radio-item",
		"dropdown-menu-radio-item-indicator",
		"dropdown-menu-label",
		"dropdown-menu-separator",
		"dropdown-menu-shortcut",
		"dropdown-menu-sub",
		"dropdown-menu-sub-trigger",
		"dropdown-menu-sub-content",
	)
}

// TestContextMenuMarkupContract used to be TestContextMenuCSSOnlyContract —
// see TestDropdownMarkupContract's own comment for why.
func TestContextMenuMarkupContract(t *testing.T) {
	got := render(t, ui.ContextMenu(
		gsx.Fragment(
			ui.ContextMenuTrigger(gsx.Raw("Open"), nil),
			ui.ContextMenuContent(
				gsx.Fragment(
					ui.ContextMenuGroup(ui.ContextMenuItem("", gsx.Raw("Item"), nil), nil),
					ui.ContextMenuCheckboxItem(true, "checked", gsx.Raw("Checked"), nil),
					ui.ContextMenuRadioGroup("one", ui.ContextMenuRadioItem(true, "one", gsx.Raw("One"), nil), nil),
					ui.ContextMenuLabel(gsx.Raw("Label"), nil),
					ui.ContextMenuSeparator(nil),
					ui.ContextMenuShortcut(gsx.Raw("⌘K"), nil),
					ui.ContextMenuSub(
						gsx.Fragment(
							ui.ContextMenuSubTrigger(gsx.Raw("More"), nil),
							ui.ContextMenuSubContent(gsx.Raw("Sub"), nil),
						),
						nil,
					),
				),
				nil,
			),
		),
		nil,
	))
	assertMenuMarkupSlots(t, got,
		"context-menu",
		"context-menu-trigger",
		"context-menu-content",
		"context-menu-group",
		"context-menu-item",
		"context-menu-checkbox-item",
		"context-menu-checkbox-item-indicator",
		"context-menu-radio-group",
		"context-menu-radio-item",
		"context-menu-radio-item-indicator",
		"context-menu-label",
		"context-menu-separator",
		"context-menu-shortcut",
		"context-menu-sub",
		"context-menu-sub-trigger",
		"context-menu-sub-content",
	)
}

// TestMenubarMarkupContract used to be TestMenubarCSSOnlyContract — see
// TestDropdownMarkupContract's own comment for why.
func TestMenubarMarkupContract(t *testing.T) {
	got := render(t, ui.Menubar(
		ui.MenubarMenu(
			gsx.Fragment(
				ui.MenubarTrigger(gsx.Raw("File"), nil),
				ui.MenubarContent(
					gsx.Fragment(
						ui.MenubarGroup(ui.MenubarItem("", gsx.Raw("Item"), nil), nil),
						ui.MenubarCheckboxItem(true, "checked", gsx.Raw("Checked"), nil),
						ui.MenubarRadioGroup("one", ui.MenubarRadioItem(true, "one", gsx.Raw("One"), nil), nil),
						ui.MenubarLabel(gsx.Raw("Label"), nil),
						ui.MenubarSeparator(nil),
						ui.MenubarShortcut(gsx.Raw("⌘K"), nil),
						ui.MenubarSub(
							gsx.Fragment(
								ui.MenubarSubTrigger(gsx.Raw("More"), nil),
								ui.MenubarSubContent(gsx.Raw("Sub"), nil),
							),
							nil,
						),
					),
					nil,
				),
			),
			nil,
		),
		nil,
	))
	assertMenuMarkupSlots(t, got,
		"menubar",
		"menubar-menu",
		"menubar-trigger",
		"menubar-content",
		"menubar-group",
		"menubar-item",
		"menubar-checkbox-item",
		"menubar-checkbox-item-indicator",
		"menubar-radio-group",
		"menubar-radio-item",
		"menubar-radio-item-indicator",
		"menubar-label",
		"menubar-separator",
		"menubar-shortcut",
		"menubar-sub",
		"menubar-sub-trigger",
		"menubar-sub-content",
	)
}

// TestNavigationMenuMarkupContract used to be TestNavigationMenuCSSOnlyContract
// — see TestDropdownMarkupContract's own comment for why.
func TestNavigationMenuMarkupContract(t *testing.T) {
	got := render(t, ui.NavigationMenu(
		ui.NavigationMenuList(
			gsx.Fragment(
				ui.NavigationMenuItem(
					gsx.Fragment(
						ui.NavigationMenuTrigger(gsx.Raw("Products"), nil),
						ui.NavigationMenuContent(
							ui.NavigationMenuLink(false, "", gsx.Raw("Docs"), nil),
							nil,
						),
					),
					nil,
				),
				ui.NavigationMenuItem(
					ui.NavigationMenuLink(false, "trigger", gsx.Raw("Home"), nil),
					nil,
				),
				ui.NavigationMenuIndicator(nil),
			),
			nil,
		),
		nil,
	))
	assertMenuMarkupSlots(t, got,
		"navigation-menu",
		"navigation-menu-list",
		"navigation-menu-item",
		"navigation-menu-trigger",
		"icon navigation-menu-trigger-icon",
		"navigation-menu-content",
		"navigation-menu-link",
		"navigation-menu-link navigation-menu-trigger",
		"navigation-menu-indicator",
		"navigation-menu-indicator-arrow",
	)
	for _, want := range []string{
		`data-viewport="false"`,
		`data-state="closed"`,
		`data-side="bottom"`,
		`data-variant="default"`,
		`data-variant="trigger"`,
		`data-active="false"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing NavigationMenu axis %q\nin: %s", want, got)
		}
	}
}

// TestCommandComboboxDropdownContextMenuMenubarNavigationMenuCallerClassesRemainCallerOnly
// used to assert every one of these six components renders class="caller-only"
// verbatim (nothing else in the attribute). That held when all six were
// CSS-only, but Dropdown/ContextMenu/Menubar/NavigationMenu migrated to the
// slot axis: their canonical .gsx now renders class={ <recipe accessors> }
// BEFORE { attrs... }, so the caller's class merges in alongside the
// recipe's own compiled utilities into ONE class attribute, rather than
// arriving as the entire attribute — the same class="X" -> bare Contains(X)
// downgrade the migration playbook's Step 9 documents for every migrated
// component's caller-class pin. Command and then Combobox followed, so all
// six now only assert the caller's token survives, merged in exactly once —
// the `exact` arm is kept (with no case selecting it) because it is the
// assertion an UNMIGRATED component must still satisfy, and the next
// component added here may well be one.
func TestCommandComboboxDropdownContextMenuMenubarNavigationMenuCallerClassesRemainCallerOnly(t *testing.T) {
	tests := []struct {
		name  string
		node  gsx.Node
		exact bool // true: class="caller-only" is the WHOLE attribute (unmigrated)
	}{
		{"Command", ui.Command(nil, gsx.Attrs{{Key: "class", Value: "caller-only"}}), false},
		{"Combobox", ui.Combobox("", "", nil, gsx.Attrs{{Key: "class", Value: "caller-only"}}), false},
		{"Dropdown", ui.DropdownMenuContent(nil, gsx.Attrs{{Key: "class", Value: "caller-only"}}), false},
		{"ContextMenu", ui.ContextMenuContent(nil, gsx.Attrs{{Key: "class", Value: "caller-only"}}), false},
		{"Menubar", ui.Menubar(nil, gsx.Attrs{{Key: "class", Value: "caller-only"}}), false},
		{"NavigationMenu", ui.NavigationMenu(nil, gsx.Attrs{{Key: "class", Value: "caller-only"}}), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := render(t, tt.node)
			if tt.exact {
				if strings.Count(got, `class="caller-only"`) != 1 || strings.Count(got, `class=`) != 1 {
					t.Errorf("caller class must be the only class and render once\nin: %s", got)
				}
				return
			}
			if strings.Count(got, "caller-only") != 1 || strings.Count(got, `class=`) != 1 {
				t.Errorf("caller class must merge in exactly once, in exactly one class attribute\nin: %s", got)
			}
		})
	}
}

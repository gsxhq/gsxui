package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
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

func assertComboboxCanonicalButtonRoles(t *testing.T, got string) {
	t.Helper()

	want := canonicalButtonClass("ghost", "icon-xs")
	if strings.Count(got, want) != 2 || strings.Count(got, `class=`) != 2 {
		t.Errorf("Combobox trigger and clear must each render exact canonical Button roles\nwant: %s\nin: %s", want, got)
	}
}

func TestCommandCSSOnlyContract(t *testing.T) {
	got := render(t, ui.Command(
		gsx.Fragment(
			ui.CommandInput("Search", nil),
			ui.CommandList(
				gsx.Fragment(
					ui.CommandEmpty(gsx.Raw("Empty"), nil),
					ui.CommandGroup("Group", ui.CommandItem("item", gsx.Raw("Item"), nil), nil),
					ui.CommandSeparator(nil),
				),
				nil,
			),
			ui.CommandShortcut(gsx.Raw("⌘K"), nil),
		),
		nil,
	))
	assertMenuCSSOnlyMarkup(t, got,
		"command",
		"command-input-wrapper",
		"command-input",
		"command-list",
		"command-empty",
		"command-group",
		"command-group-heading",
		"command-item",
		"command-separator",
		"command-shortcut",
	)
	for _, hook := range []string{
		`data-gsxui-command`,
		`data-gsxui-command-input`,
		`data-gsxui-command-list`,
		`data-gsxui-command-empty`,
		`data-gsxui-command-group`,
		`data-gsxui-command-group-heading`,
		`data-gsxui-command-item`,
		`data-gsxui-command-separator`,
	} {
		if !strings.Contains(got, hook) {
			t.Errorf("missing Command behavior hook %q\nin: %s", hook, got)
		}
	}
}

func TestCommandDialogCSSOnlyComposition(t *testing.T) {
	got := render(t, ui.CommandDialog("", "", nil, ui.CommandInput("Search", nil), nil))
	assertMenuCSSOnlyMarkup(t, got,
		"dialog command-dialog",
		"dialog-content command-dialog-content",
		"dialog-header command-dialog-header",
		"dialog-title",
		"dialog-description",
		"command command-dialog-command",
		"command-input-wrapper",
		"command-input",
		"dialog-close dialog-close-button",
	)
	for _, want := range []string{
		`data-gsxui-dialog-content`,
		`data-gsxui-command-dialog`,
		`data-gsxui-command`,
		`data-gsxui-command-input`,
		`>Command Palette</h2>`,
		`>Search for a command to run...</p>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing CommandDialog contract %q\nin: %s", want, got)
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

func TestDropdownCSSOnlyContract(t *testing.T) {
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
	assertMenuCSSOnlyMarkup(t, got,
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

func TestContextMenuCSSOnlyContract(t *testing.T) {
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
	assertMenuCSSOnlyMarkup(t, got,
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

func TestMenubarCSSOnlyContract(t *testing.T) {
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
	assertMenuCSSOnlyMarkup(t, got,
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

func TestNavigationMenuCSSOnlyContract(t *testing.T) {
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
	assertMenuCSSOnlyMarkup(t, got,
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

func TestCommandComboboxDropdownContextMenuMenubarNavigationMenuCallerClassesRemainCallerOnly(t *testing.T) {
	tests := []struct {
		name string
		node gsx.Node
	}{
		{"Command", ui.Command(nil, gsx.Attrs{{Key: "class", Value: "caller-only"}})},
		{"Combobox", ui.Combobox("", "", nil, gsx.Attrs{{Key: "class", Value: "caller-only"}})},
		{"Dropdown", ui.DropdownMenuContent(nil, gsx.Attrs{{Key: "class", Value: "caller-only"}})},
		{"ContextMenu", ui.ContextMenuContent(nil, gsx.Attrs{{Key: "class", Value: "caller-only"}})},
		{"Menubar", ui.Menubar(nil, gsx.Attrs{{Key: "class", Value: "caller-only"}})},
		{"NavigationMenu", ui.NavigationMenu(nil, gsx.Attrs{{Key: "class", Value: "caller-only"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := render(t, tt.node)
			if strings.Count(got, `class="caller-only"`) != 1 || strings.Count(got, `class=`) != 1 {
				t.Errorf("caller class must be the only class and render once\nin: %s", got)
			}
		})
	}
}

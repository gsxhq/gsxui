package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestContextMenuStructure(t *testing.T) {
	got := render(t, ui.ContextMenu(
		gsx.Fragment(
			ui.ContextMenuTrigger(gsx.Raw("Right click"), nil),
			ui.ContextMenuContent(
				gsx.Fragment(
					ui.ContextMenuLabel(gsx.Raw("Actions"), nil),
					ui.ContextMenuSeparator(nil),
					ui.ContextMenuItem("", gsx.Fragment(
						gsx.Raw("Back"),
						ui.ContextMenuShortcut(gsx.Raw("⌘["), nil),
					), nil),
					ui.ContextMenuItem("destructive", gsx.Raw("Delete"), nil),
				),
				nil,
			),
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-slot-context-menu`, `data-gsxui-contextmenu`,
		`data-gsxui-slot-context-menu-trigger`, `data-gsxui-contextmenu-trigger`,
		`data-gsxui-slot-context-menu-content`, `data-gsxui-contextmenu-content`,
		`popover="auto"`, `role="menu"`, `data-state="closed"`,
		`data-gsxui-slot-context-menu-label`, ">Actions<",
		`data-gsxui-slot-context-menu-separator`, `role="separator"`,
		`data-gsxui-slot-context-menu-item`, `data-gsxui-contextmenu-item`,
		`role="menuitem"`, `tabindex="-1"`,
		`data-variant="default"`, ">Back<",
		`data-gsxui-slot-context-menu-shortcut`, ">⌘[<",
		`data-variant="destructive"`, ">Delete<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	assertMenuCSSOnlyMarkup(t, got,
		"context-menu",
		"context-menu-trigger",
		"context-menu-content",
		"context-menu-label",
		"context-menu-separator",
		"context-menu-item",
		"context-menu-shortcut",
	)
}

func TestContextMenuItemVariants(t *testing.T) {
	for input, want := range map[string]string{"": "default", "destructive": "destructive"} {
		got := render(t, ui.ContextMenuItem(input, gsx.Raw("x"), nil))
		if !strings.Contains(got, `data-variant="`+want+`"`) {
			t.Errorf("variant %q must stamp %q\nin: %s", input, want, got)
		}
		if strings.Contains(got, `class=`) {
			t.Errorf("variant presentation must live in CSS\nin: %s", got)
		}
	}
}

func TestContextMenuContentCallerClassMerges(t *testing.T) {
	got := render(t, ui.ContextMenuContent(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "z-10"}}))
	if strings.Count(got, `class="z-10"`) != 1 || strings.Count(got, `class=`) != 1 {
		t.Errorf("caller z-10 must be the only class and render once\nin: %s", got)
	}
}

func TestContextMenuAttrsFallThrough(t *testing.T) {
	got := render(t, ui.ContextMenuContent(gsx.Raw("x"), gsx.Attrs{
		{Key: "id", Value: "my-menu"},
		{Key: "aria-label", Value: "Actions"},
	}))
	for _, want := range []string{`id="my-menu"`, `aria-label="Actions"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestContextMenuPopoverAndSideAttrsOverridable(t *testing.T) {
	got := render(t, ui.ContextMenuSubContent(gsx.Raw("x"), gsx.Attrs{
		{Key: "popover", Value: "manual"},
		{Key: "data-side", Value: "top"},
	}))
	for _, want := range []string{`popover="manual"`, `data-side="top"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing caller override %q\nin: %s", want, got)
		}
	}
}

func TestContextMenuTriggerPinned(t *testing.T) {
	got := render(t, ui.ContextMenuTrigger(gsx.Raw("Right click here"), nil))
	want := `<div data-gsxui-contextmenu-trigger data-gsxui-slot-context-menu-trigger>Right click here</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestContextMenuContentPinned(t *testing.T) {
	got := render(t, ui.ContextMenuContent(gsx.Raw("x"), nil))
	want := `<div data-gsxui-contextmenu-content popover="auto" role="menu" tabindex="-1" data-state="closed" data-gsxui-slot-context-menu-content>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
	if strings.Contains(got, `data-side=`) {
		t.Errorf("cursor-positioned top-level context content must not invent a side\nin: %s", got)
	}
}

func TestContextMenuCheckedItemsStampIndicatorState(t *testing.T) {
	checkbox := render(t, ui.ContextMenuCheckboxItem(true, "status", gsx.Raw("Status"), nil))
	for _, want := range []string{
		`role="menuitemcheckbox"`, `aria-checked="true"`, `data-state="checked"`,
		`data-gsxui-contextmenu-checkbox-item`,
		`data-gsxui-slot-context-menu-checkbox-item-indicator`,
	} {
		if !strings.Contains(checkbox, want) {
			t.Errorf("checkbox missing %q\nin: %s", want, checkbox)
		}
	}
	radio := render(t, ui.ContextMenuRadioItem(true, "top", gsx.Raw("Top"), nil))
	for _, want := range []string{
		`role="menuitemradio"`, `aria-checked="true"`, `data-state="checked"`,
		`data-gsxui-contextmenu-radio-item`,
		`data-gsxui-slot-context-menu-radio-item-indicator`,
	} {
		if !strings.Contains(radio, want) {
			t.Errorf("radio missing %q\nin: %s", want, radio)
		}
	}
}

func TestContextMenuSubNestsContentInsideParentContentPinned(t *testing.T) {
	got := render(t, ui.ContextMenuContent(
		ui.ContextMenuSub(
			gsx.Fragment(
				ui.ContextMenuSubTrigger(gsx.Raw("More"), nil),
				ui.ContextMenuSubContent(gsx.Raw("INNER"), nil),
			),
			nil,
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-slot-context-menu-content`,
		`data-gsxui-slot-context-menu-sub`,
		`data-gsxui-slot-context-menu-sub-trigger`,
		`aria-haspopup="menu"`,
		`aria-expanded="false"`,
		`data-state="closed"`,
		`data-gsxui-slot-icon`,
		`data-gsxui-slot-context-menu-sub-content`,
		`data-side="right"`,
		">INNER<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing nested submenu contract %q\nin: %s", want, got)
		}
	}
	assertMenuCSSOnlyMarkup(t, got,
		"context-menu-content",
		"context-menu-sub",
		"context-menu-sub-trigger",
		"context-menu-sub-content",
	)
}

func TestContextMenuPinned(t *testing.T) {
	got := render(t, ui.ContextMenuItem("", gsx.Raw("Back"), nil))
	want := `<div data-gsxui-contextmenu-item data-variant="default" role="menuitem" tabindex="-1" data-gsxui-slot-context-menu-item>Back</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

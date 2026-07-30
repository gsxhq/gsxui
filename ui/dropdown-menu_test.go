package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestDropdownStructure(t *testing.T) {
	got := render(t, ui.DropdownMenu(
		gsx.Fragment(
			ui.DropdownMenuTrigger(gsx.Raw("Open"), nil),
			ui.DropdownMenuContent(gsx.Fragment(
				ui.DropdownMenuLabel(gsx.Raw("Actions"), nil),
				ui.DropdownMenuSeparator(nil),
				ui.DropdownMenuItem("", gsx.Fragment(
					gsx.Raw("Edit"),
					ui.DropdownMenuShortcut(gsx.Raw("⌘E"), nil),
				), nil),
				ui.DropdownMenuItem("destructive", gsx.Raw("Delete"), nil),
			), nil),
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-slot-dropdown-menu`, `data-gsxui-dropdown`,
		`data-gsxui-slot-dropdown-menu-trigger`, `data-gsxui-dropdown-trigger`,
		`aria-haspopup="menu"`, `aria-expanded="false"`,
		`data-gsxui-slot-dropdown-menu-content`, `data-gsxui-dropdown-content`,
		`popover="auto"`, `role="menu"`, `data-state="closed"`, `data-side="bottom"`,
		`data-gsxui-slot-dropdown-menu-label`, ">Actions<",
		`data-gsxui-slot-dropdown-menu-separator`, `role="separator"`,
		`data-gsxui-slot-dropdown-menu-item`, `data-gsxui-dropdown-item`,
		`role="menuitem"`, `tabindex="-1"`,
		`data-variant="default"`, ">Edit<",
		`data-gsxui-slot-dropdown-menu-shortcut`, ">⌘E<",
		`data-variant="destructive"`, ">Delete<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	assertMenuMarkupSlots(t, got,
		"dropdown-menu",
		"dropdown-menu-trigger",
		"dropdown-menu-content",
		"dropdown-menu-label",
		"dropdown-menu-separator",
		"dropdown-menu-item",
		"dropdown-menu-shortcut",
	)
}

func TestDropdownItemVariants(t *testing.T) {
	// Migrated to the slot axis: variant presentation now lives in the
	// recipe's compiled class (dropdownMenu.ItemVariant), so a class
	// attribute legitimately renders here now.
	for input, want := range map[string]string{"": "default", "destructive": "destructive"} {
		got := render(t, ui.DropdownMenuItem(input, gsx.Raw("x"), nil))
		if !strings.Contains(got, `data-variant="`+want+`"`) {
			t.Errorf("variant %q must stamp %q\nin: %s", input, want, got)
		}
		if strings.Count(got, `class=`) != 1 {
			t.Errorf("variant presentation must render exactly one class attribute\nin: %s", got)
		}
	}
}

func TestDropdownContentCallerClassMerges(t *testing.T) {
	got := render(t, ui.DropdownMenuContent(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "z-10"}}))
	if strings.Count(got, "z-10") != 1 || strings.Count(got, `class=`) != 1 {
		t.Errorf("caller z-10 must merge in exactly once, in exactly one class attribute\nin: %s", got)
	}
}

func TestDropdownAttrsFallThrough(t *testing.T) {
	got := render(t, ui.DropdownMenuContent(gsx.Raw("x"), gsx.Attrs{
		{Key: "id", Value: "my-menu"},
		{Key: "aria-label", Value: "Actions"},
	}))
	for _, want := range []string{`id="my-menu"`, `aria-label="Actions"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestDropdownPopoverAndSideAttrsOverridable(t *testing.T) {
	got := render(t, ui.DropdownMenuSubContent(gsx.Raw("x"), gsx.Attrs{
		{Key: "popover", Value: "manual"},
		{Key: "data-side", Value: "left"},
	}))
	for _, want := range []string{`popover="manual"`, `data-side="left"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing caller override %q\nin: %s", want, got)
		}
	}
	for _, replaced := range []string{`popover="auto"`, `data-side="right"`} {
		if strings.Contains(got, replaced) {
			t.Errorf("caller override must replace %q\nin: %s", replaced, got)
		}
	}
}

func TestDropdownMenuCheckedItemsStampIndicatorState(t *testing.T) {
	checkbox := render(t, ui.DropdownMenuCheckboxItem(true, "show-toolbar", gsx.Raw("Toolbar"), nil))
	for _, want := range []string{
		`role="menuitemcheckbox"`,
		`aria-checked="true"`,
		`data-state="checked"`,
		`data-gsxui-dropdown-checkbox-item`,
		`data-value="show-toolbar"`,
		`data-gsxui-slot-dropdown-menu-checkbox-item-indicator`,
	} {
		if !strings.Contains(checkbox, want) {
			t.Errorf("checkbox missing %q\nin: %s", want, checkbox)
		}
	}
	radio := render(t, ui.DropdownMenuRadioItem(true, "top", gsx.Raw("Top"), nil))
	for _, want := range []string{
		`role="menuitemradio"`,
		`aria-checked="true"`,
		`data-state="checked"`,
		`data-value="top"`,
		`data-gsxui-slot-dropdown-menu-radio-item-indicator`,
	} {
		if !strings.Contains(radio, want) {
			t.Errorf("radio missing %q\nin: %s", want, radio)
		}
	}
	unchecked := render(t, ui.DropdownMenuCheckboxItem(false, "show-toolbar", gsx.Raw("Toolbar"), nil))
	for _, want := range []string{`aria-checked="false"`, `data-state="unchecked"`} {
		if !strings.Contains(unchecked, want) {
			t.Errorf("unchecked checkbox missing %q\nin: %s", want, unchecked)
		}
	}
}

func TestDropdownMenuSubNestsContentInsideParentContentPinned(t *testing.T) {
	got := render(t, ui.DropdownMenuContent(
		ui.DropdownMenuSub(
			gsx.Fragment(
				ui.DropdownMenuSubTrigger(gsx.Raw("More"), nil),
				ui.DropdownMenuSubContent(gsx.Raw("INNER"), nil),
			),
			nil,
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-slot-dropdown-menu-content`,
		`data-gsxui-slot-dropdown-menu-sub`,
		`data-gsxui-slot-dropdown-menu-sub-trigger`,
		`aria-haspopup="menu"`,
		`aria-expanded="false"`,
		`data-state="closed"`,
		`data-gsxui-slot-icon`,
		`data-gsxui-slot-dropdown-menu-sub-content`,
		`data-side="right"`,
		">INNER<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing nested submenu contract %q\nin: %s", want, got)
		}
	}
	if strings.Index(got, `data-gsxui-slot-dropdown-menu-sub-content`) >
		strings.LastIndex(got, "</div>") {
		t.Errorf("submenu content must remain DOM-nested\nin: %s", got)
	}
	assertMenuMarkupSlots(t, got,
		"dropdown-menu-content",
		"dropdown-menu-sub",
		"dropdown-menu-sub-trigger",
		"dropdown-menu-sub-content",
	)
}

func TestDropdownPinned(t *testing.T) {
	got := render(t, ui.DropdownMenuItem("", gsx.Raw("Edit"), nil))
	want := `<div class="relative flex cursor-default items-center gap-1.5 rounded-md px-1.5 py-1 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-disabled:pointer-events-none [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 [&amp;_svg:not([class*=&#39;text-&#39;])]:text-muted-foreground text-foreground" data-gsxui-dropdown-item data-variant="default" role="menuitem" tabindex="-1" data-gsxui-slot-dropdown-menu-item>Edit</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestDropdownContentPinned(t *testing.T) {
	got := render(t, ui.DropdownMenuContent(gsx.Raw("x"), nil))
	want := `<div class="z-50 max-h-96 min-w-[8rem] origin-top-left overflow-x-hidden overflow-y-auto rounded-lg border bg-popover p-1 text-popover-foreground shadow-md opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 [&amp;:popover-open]:opacity-100 [&amp;:popover-open]:scale-100 starting:[&amp;:popover-open]:opacity-0 starting:[&amp;:popover-open]:scale-95 data-[side=bottom]:starting:[&amp;:popover-open]:-translate-y-2 data-[side=left]:starting:[&amp;:popover-open]:translate-x-2 data-[side=right]:starting:[&amp;:popover-open]:-translate-x-2 data-[side=top]:starting:[&amp;:popover-open]:translate-y-2" data-gsxui-dropdown-content popover="auto" role="menu" tabindex="-1" data-state="closed" data-side="bottom" data-gsxui-slot-dropdown-menu-content>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

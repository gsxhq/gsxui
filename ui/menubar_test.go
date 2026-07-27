package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestMenubarTriggerRovingAndAria(t *testing.T) {
	got := render(t, ui.MenubarTrigger(gsx.Raw("File"), nil))
	for _, want := range []string{
		`data-gsxui-menubar-trigger`,
		`type="button"`,
		`aria-haspopup="menu"`,
		`aria-expanded="false"`,
		`data-state="closed"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if strings.Contains(got, `tabindex=`) {
		t.Errorf("server markup leaves graceful tab stops for JS normalization\nin: %s", got)
	}
}

func TestMenubarTriggerDisabledRendersNativeAttr(t *testing.T) {
	got := render(t, ui.MenubarTrigger(gsx.Raw("File"), gsx.Attrs{{Key: "disabled", Value: true}}))
	if !strings.Contains(got, ` disabled`) {
		t.Errorf("disabled trigger must render the native attribute\nin: %s", got)
	}
}

func TestMenubarStructure(t *testing.T) {
	got := render(t, ui.Menubar(
		ui.MenubarMenu(
			gsx.Fragment(
				ui.MenubarTrigger(gsx.Raw("File"), nil),
				ui.MenubarContent(gsx.Fragment(
					ui.MenubarItem("", gsx.Fragment(
						gsx.Raw("New Tab"),
						ui.MenubarShortcut(gsx.Raw("⌘T"), nil),
					), nil),
					ui.MenubarSeparator(nil),
					ui.MenubarItem("destructive", gsx.Raw("Delete"), nil),
				), nil),
			),
			nil,
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-slot="menubar"`, `role="menubar"`, `data-gsxui-menubar`,
		`data-gsxui-slot="menubar-menu"`, `data-gsxui-menubar-menu`,
		`data-gsxui-slot="menubar-trigger"`, `data-gsxui-menubar-trigger`,
		`aria-haspopup="menu"`, `aria-expanded="false"`,
		`data-gsxui-slot="menubar-content"`, `data-gsxui-menubar-content`,
		`popover="auto"`, `role="menu"`, `data-state="closed"`, `data-side="bottom"`,
		`data-gsxui-slot="menubar-item"`, `data-gsxui-menubar-item`,
		`role="menuitem"`, `tabindex="-1"`, `data-variant="default"`,
		">New Tab<", `data-gsxui-slot="menubar-shortcut"`, ">⌘T<",
		`data-gsxui-slot="menubar-separator"`, `role="separator"`,
		`data-variant="destructive"`, ">Delete<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	assertMenuCSSOnlyMarkup(t, got,
		"menubar",
		"menubar-menu",
		"menubar-trigger",
		"menubar-content",
		"menubar-item",
		"menubar-shortcut",
		"menubar-separator",
	)
}

func TestMenubarItemVariants(t *testing.T) {
	for input, want := range map[string]string{"": "default", "destructive": "destructive"} {
		got := render(t, ui.MenubarItem(input, gsx.Raw("x"), nil))
		if !strings.Contains(got, `data-variant="`+want+`"`) {
			t.Errorf("variant %q must stamp %q\nin: %s", input, want, got)
		}
		if strings.Contains(got, `class=`) {
			t.Errorf("variant presentation must live in CSS\nin: %s", got)
		}
	}
}

func TestMenubarContentCallerClassMerges(t *testing.T) {
	got := render(t, ui.MenubarContent(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "z-10"}}))
	if strings.Count(got, `class="z-10"`) != 1 || strings.Count(got, `class=`) != 1 {
		t.Errorf("caller z-10 must be the only class and render once\nin: %s", got)
	}
}

func TestMenubarAttrsFallThrough(t *testing.T) {
	got := render(t, ui.MenubarContent(gsx.Raw("x"), gsx.Attrs{
		{Key: "id", Value: "my-menu"},
		{Key: "aria-label", Value: "File"},
	}))
	for _, want := range []string{`id="my-menu"`, `aria-label="File"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestMenubarPopoverAndSideAttrsOverridable(t *testing.T) {
	got := render(t, ui.MenubarSubContent(gsx.Raw("x"), gsx.Attrs{
		{Key: "popover", Value: "manual"},
		{Key: "data-side", Value: "bottom"},
	}))
	for _, want := range []string{`popover="manual"`, `data-side="bottom"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing caller override %q\nin: %s", want, got)
		}
	}
}

func TestMenubarCheckedItemsStampIndicatorState(t *testing.T) {
	checkbox := render(t, ui.MenubarCheckboxItem(true, "toolbar", gsx.Raw("Toolbar"), nil))
	for _, want := range []string{
		`role="menuitemcheckbox"`, `aria-checked="true"`, `data-state="checked"`,
		`data-gsxui-menubar-checkbox-item`,
		`data-gsxui-slot="menubar-checkbox-item-indicator"`,
	} {
		if !strings.Contains(checkbox, want) {
			t.Errorf("checkbox missing %q\nin: %s", want, checkbox)
		}
	}
	radio := render(t, ui.MenubarRadioItem(true, "top", gsx.Raw("Top"), nil))
	for _, want := range []string{
		`role="menuitemradio"`, `aria-checked="true"`, `data-state="checked"`,
		`data-gsxui-menubar-radio-item`,
		`data-gsxui-slot="menubar-radio-item-indicator"`,
	} {
		if !strings.Contains(radio, want) {
			t.Errorf("radio missing %q\nin: %s", want, radio)
		}
	}
}

func TestMenubarSubNestsContentInsideParentContentPinned(t *testing.T) {
	got := render(t, ui.MenubarContent(
		ui.MenubarSub(
			gsx.Fragment(
				ui.MenubarSubTrigger(gsx.Raw("More"), nil),
				ui.MenubarSubContent(gsx.Raw("INNER"), nil),
			),
			nil,
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-slot="menubar-content"`,
		`data-gsxui-slot="menubar-sub"`,
		`data-gsxui-slot="menubar-sub-trigger"`,
		`aria-haspopup="menu"`,
		`aria-expanded="false"`,
		`data-state="closed"`,
		`data-gsxui-slot="icon"`,
		`data-gsxui-slot="menubar-sub-content"`,
		`data-side="right"`,
		">INNER<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing nested submenu contract %q\nin: %s", want, got)
		}
	}
	assertMenuCSSOnlyMarkup(t, got,
		"menubar-content",
		"menubar-sub",
		"menubar-sub-trigger",
		"menubar-sub-content",
	)
}

func TestMenubarPinnedParts(t *testing.T) {
	tests := []struct {
		name string
		node gsx.Node
		want string
	}{
		{
			name: "root",
			node: ui.Menubar(gsx.Raw("x"), nil),
			want: `<div data-gsxui-menubar role="menubar" data-gsxui-slot="menubar">x</div>`,
		},
		{
			name: "menu",
			node: ui.MenubarMenu(gsx.Raw("x"), nil),
			want: `<div data-gsxui-menubar-menu data-gsxui-slot="menubar-menu">x</div>`,
		},
		{
			name: "item",
			node: ui.MenubarItem("", gsx.Raw("Print"), nil),
			want: `<div data-gsxui-menubar-item data-variant="default" role="menuitem" tabindex="-1" data-gsxui-slot="menubar-item">Print</div>`,
		},
		{
			name: "content",
			node: ui.MenubarContent(gsx.Raw("x"), nil),
			want: `<div data-gsxui-menubar-content popover="auto" role="menu" tabindex="-1" data-state="closed" data-side="bottom" data-gsxui-slot="menubar-content">x</div>`,
		},
		{
			name: "label",
			node: ui.MenubarLabel(gsx.Raw("People"), nil),
			want: `<div data-gsxui-slot="menubar-label">People</div>`,
		},
		{
			name: "separator",
			node: ui.MenubarSeparator(nil),
			want: `<div role="separator" data-gsxui-slot="menubar-separator"></div>`,
		},
		{
			name: "shortcut",
			node: ui.MenubarShortcut(gsx.Raw("⌘T"), nil),
			want: `<span data-gsxui-slot="menubar-shortcut">⌘T</span>`,
		},
		{
			name: "group",
			node: ui.MenubarGroup(gsx.Raw("x"), nil),
			want: `<div role="group" data-gsxui-slot="menubar-group">x</div>`,
		},
		{
			name: "radio-group",
			node: ui.MenubarRadioGroup("benoit", gsx.Raw("x"), nil),
			want: `<div data-gsxui-menubar-radio-group role="group" data-value="benoit" data-gsxui-slot="menubar-radio-group">x</div>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := render(t, tt.node); got != tt.want {
				t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, tt.want)
			}
		})
	}
}

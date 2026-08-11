package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestNavigationMenuTriggerAria(t *testing.T) {
	got := render(t, ui.NavigationMenuTrigger(gsx.Raw("Products"), nil))
	for _, want := range []string{
		`aria-expanded="false"`,
		`data-gsxui-slot-navigation-menu-trigger`,
		`data-state="closed"`,
		`data-gsxui-slot-navigation-menu-trigger`,
		`data-gsxui-slot-navigation-menu-trigger-icon data-gsxui-slot-icon`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestNavigationMenuListIsAList(t *testing.T) {
	got := render(t, ui.NavigationMenuList(gsx.Raw("x"), nil))
	if !strings.Contains(got, `data-gsxui-slot-navigation-menu-list`) {
		t.Errorf("want the list slot\nin: %s", got)
	}
	if !strings.HasPrefix(got, "<ul") {
		t.Errorf("want an actual <ul>, not a div\nin: %s", got)
	}
}

func TestNavigationMenuRootPinned(t *testing.T) {
	// group/navigation-menu was added (registry/canonical/navigation-menu.gsx)
	// so NavigationMenuContent's own group-data-[viewport=false]/
	// navigation-menu: selectors have a real ancestor to match — see the
	// style-porter report's "missing group/<name> marker" entry.
	got := render(t, ui.NavigationMenu(gsx.Raw("x"), nil))
	want := `<nav class="group/navigation-menu max-w-max flex" data-viewport="false" data-gsxui-slot-navigation-menu>x</nav>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestNavigationMenuViewportPartRemoved(t *testing.T) {
	got := render(t, ui.NavigationMenuItem(gsx.Fragment(
		ui.NavigationMenuTrigger(gsx.Raw("Products"), nil),
		ui.NavigationMenuContent(gsx.Raw("links"), nil),
	), nil))
	if strings.Contains(got, "navigation-menu-viewport") {
		t.Errorf("navigation-menu-viewport must not appear anywhere\nin: %s", got)
	}
}

func TestNavigationMenuPinnedStructuralParts(t *testing.T) {
	tests := []struct {
		name string
		node gsx.Node
		want string
	}{
		{
			name: "list",
			node: ui.NavigationMenuList(gsx.Raw("x"), nil),
			want: `<ul class="gap-0 flex" data-gsxui-slot-navigation-menu-list>x</ul>`,
		},
		{
			name: "item",
			node: ui.NavigationMenuItem(gsx.Raw("x"), nil),
			want: `<li class="relative" data-gsxui-slot-navigation-menu-item>x</li>`,
		},
		{
			// transition-none suppresses an implicit transition:all that
			// duration-300 alone leaves active — see the style-porter
			// report's "duration-N alone" entry.
			name: "content",
			node: ui.NavigationMenuContent(gsx.Raw("x"), nil),
			want: `<div class="transition-none data-[motion^=from-]:animate-in data-[motion^=to-]:animate-out data-[motion^=from-]:fade-in data-[motion^=to-]:fade-out data-[motion=from-end]:slide-in-from-right-52 data-[motion=from-start]:slide-in-from-left-52 data-[motion=to-end]:slide-out-to-right-52 data-[motion=to-start]:slide-out-to-left-52 group-data-[viewport=false]/navigation-menu:bg-popover group-data-[viewport=false]/navigation-menu:text-popover-foreground group-data-[viewport=false]/navigation-menu:data-[state=open]:animate-in group-data-[viewport=false]/navigation-menu:data-[state=closed]:animate-out group-data-[viewport=false]/navigation-menu:data-[state=closed]:zoom-out-95 group-data-[viewport=false]/navigation-menu:data-[state=open]:zoom-in-95 group-data-[viewport=false]/navigation-menu:data-[state=open]:fade-in-0 group-data-[viewport=false]/navigation-menu:data-[state=closed]:fade-out-0 group-data-[viewport=false]/navigation-menu:ring-foreground/10 p-1 ease-[cubic-bezier(0.22,1,0.36,1)] group-data-[viewport=false]/navigation-menu:rounded-lg group-data-[viewport=false]/navigation-menu:shadow group-data-[viewport=false]/navigation-menu:ring-1 group-data-[viewport=false]/navigation-menu:duration-300" popover="manual" data-state="closed" data-side="bottom" data-gsxui-slot-navigation-menu-content>x</div>`,
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

func TestNavigationMenuTriggerPinned(t *testing.T) {
	got := render(t, ui.NavigationMenuTrigger(gsx.Raw("Products"), nil))
	for _, want := range []string{
		`<button`,
		`data-gsxui-slot-navigation-menu-trigger`,
		`type="button"`,
		`aria-expanded="false"`,
		`data-state="closed"`,
		`data-gsxui-slot-navigation-menu-trigger`,
		`>Products <svg`,
		`data-gsxui-slot-navigation-menu-trigger-icon data-gsxui-slot-icon`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing trigger contract %q\nin: %s", want, got)
		}
	}
	// Two class attributes legitimately render here now: the trigger button's
	// own recipe class, and its child chevron icon's own recipe class.
	if strings.Count(got, `class=`) != 2 {
		t.Errorf("trigger and icon presentation must render exactly two class attributes\nin: %s", got)
	}
}

func TestNavigationMenuLinkVariantReplacesClassHelper(t *testing.T) {
	for input, wantVariant := range map[string]string{"": "default", "trigger": "trigger"} {
		got := render(t, ui.NavigationMenuLink(false, input, gsx.Raw("Docs"), nil))
		for _, want := range []string{
			`data-gsxui-slot-navigation-menu-link`,
			`data-variant="` + wantVariant + `"`,
			`data-active="false"`,
		} {
			if !strings.Contains(got, want) {
				t.Errorf("variant %q missing %q\nin: %s", input, want, got)
			}
		}
		slot := `data-gsxui-slot-navigation-menu-link`
		if wantVariant == "trigger" {
			slot = `data-gsxui-slot-navigation-menu-trigger data-gsxui-slot-navigation-menu-link`
		}
		if !strings.Contains(got, slot) {
			t.Errorf("variant %q missing slot order %q\nin: %s", input, slot, got)
		}
		if strings.Count(got, `class=`) != 1 {
			t.Errorf("variant presentation must render exactly one class attribute\nin: %s", got)
		}
	}
}

func TestNavigationMenuContentHasNoFixedWidthUtility(t *testing.T) {
	got := render(t, ui.NavigationMenuContent(gsx.Raw("x"), nil))
	for _, dead := range []string{"w-full", "md:w-auto"} {
		if strings.Contains(got, dead) {
			t.Errorf("Content must not carry %q\nin: %s", dead, got)
		}
	}
}

func TestNavigationMenuOnlyOnePopoverPerItem(t *testing.T) {
	got := render(t, ui.NavigationMenu(
		ui.NavigationMenuList(
			ui.NavigationMenuItem(gsx.Fragment(
				ui.NavigationMenuTrigger(gsx.Raw("Products"), nil),
				ui.NavigationMenuContent(gsx.Raw(`<a>Dialog</a>`), nil),
			), nil),
			nil,
		),
		nil,
	))
	if n := strings.Count(got, `popover="`); n != 1 {
		t.Errorf("want exactly 1 popover-backed element, got %d\nin: %s", n, got)
	}
	if !strings.Contains(got, "Dialog") {
		t.Errorf("Content's own children must render, unoccluded\nin: %s", got)
	}
}

func TestNavigationMenuLinkActive(t *testing.T) {
	got := render(t, ui.NavigationMenuLink(true, "", gsx.Raw("Docs"), nil))
	if !strings.Contains(got, `data-active="true"`) {
		t.Errorf("active link must reflect data-active=true\nin: %s", got)
	}
}

func TestNavigationMenuIndicatorPinned(t *testing.T) {
	got := render(t, ui.NavigationMenuIndicator(nil))
	for _, want := range []string{
		`data-gsxui-slot-navigation-menu-indicator`,
		`data-state="hidden"`,
		`data-gsxui-slot-navigation-menu-indicator`,
		`data-gsxui-slot-navigation-menu-indicator-arrow`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("indicator missing %q\nin: %s", want, got)
		}
	}
	// Two class attributes legitimately render here now: the indicator's own
	// recipe class, and its child arrow div's own recipe class.
	if strings.Count(got, `class=`) != 2 {
		t.Errorf("indicator and arrow presentation must render exactly two class attributes\nin: %s", got)
	}
}

func TestNavigationMenuCallerClassMerges(t *testing.T) {
	got := render(t, ui.NavigationMenuContent(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "p-8"}}))
	if strings.Count(got, "p-8") != 1 || strings.Count(got, `class=`) != 1 {
		t.Errorf("caller p-8 must merge in exactly once, in exactly one class attribute\nin: %s", got)
	}
}

func TestNavigationMenuAttrsFallThrough(t *testing.T) {
	got := render(t, ui.NavigationMenuTrigger(gsx.Raw("x"), gsx.Attrs{{Key: "aria-label", Value: "Open products"}}))
	if !strings.Contains(got, `aria-label="Open products"`) {
		t.Errorf("missing aria-label\nin: %s", got)
	}
}

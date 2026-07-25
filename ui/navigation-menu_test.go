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
		`data-gsxui-navigation-menu-trigger`,
		`data-state="closed"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestNavigationMenuListIsAList(t *testing.T) {
	got := render(t, ui.NavigationMenuList(gsx.Raw("x"), nil))
	if !strings.Contains(got, `data-slot="navigation-menu-list"`) {
		t.Errorf("want the list slot\nin: %s", got)
	}
	if !strings.HasPrefix(got, "<ul") {
		t.Errorf("want an actual <ul>, not a div\nin: %s", got)
	}
}

func TestNavigationMenuRootPinned(t *testing.T) {
	// Exact full-render pin for the default (viewport=true) root: data-viewport
	// stamped "true", the shared NavigationMenuViewport auto-rendered as the
	// last child after the caller's own children.
	got := render(t, ui.NavigationMenu("", gsx.Raw("x"), nil))
	if !strings.HasPrefix(got, `<nav data-slot="navigation-menu" data-gsxui-navigation-menu data-viewport="true" class="group/navigation-menu relative flex max-w-max flex-1 items-center justify-center">x`) {
		t.Errorf("root prefix mismatch\ngot: %s", got)
	}
	if !strings.Contains(got, `data-slot="navigation-menu-viewport"`) {
		t.Errorf("want the auto-rendered viewport as a child\nin: %s", got)
	}
	if !strings.HasSuffix(got, "</nav>") {
		t.Errorf("want a closing </nav>\nin: %s", got)
	}
}

func TestNavigationMenuRootViewportFalseSkipsViewport(t *testing.T) {
	got := render(t, ui.NavigationMenu("false", gsx.Raw("x"), nil))
	if !strings.Contains(got, `data-viewport="false"`) {
		t.Errorf("want data-viewport=false\nin: %s", got)
	}
	if strings.Contains(got, `data-slot="navigation-menu-viewport"`) {
		t.Errorf("viewport=false must not auto-render NavigationMenuViewport\nin: %s", got)
	}
	want := `<nav data-slot="navigation-menu" data-gsxui-navigation-menu data-viewport="false" class="group/navigation-menu relative flex max-w-max flex-1 items-center justify-center">x</nav>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestNavigationMenuListPinned(t *testing.T) {
	got := render(t, ui.NavigationMenuList(gsx.Raw("x"), nil))
	want := `<ul data-slot="navigation-menu-list" data-gsxui-navigation-menu-list class="group flex flex-1 list-none items-center justify-center gap-1">x</ul>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestNavigationMenuItemPinned(t *testing.T) {
	got := render(t, ui.NavigationMenuItem(gsx.Raw("x"), nil))
	want := `<li data-slot="navigation-menu-item" data-gsxui-navigation-menu-item class="relative">x</li>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestNavigationMenuTriggerPinned(t *testing.T) {
	got := render(t, ui.NavigationMenuTrigger(gsx.Raw("Products"), nil))
	want := `<button data-slot="navigation-menu-trigger" data-gsxui-navigation-menu-trigger type="button" aria-expanded="false" data-state="closed" class="group inline-flex h-9 w-max items-center justify-center rounded-lg px-4 py-2 text-sm font-medium transition-[color,box-shadow] outline-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1 disabled:pointer-events-none disabled:opacity-50 data-[state=open]:bg-accent/50 data-[state=open]:text-accent-foreground data-[state=open]:hover:bg-accent data-[state=open]:focus:bg-accent group">Products<svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="relative top-[1px] ml-1 size-3 transition duration-300 group-data-[state=open]:rotate-180"><path d="m6 9 6 6 6-6"/></svg></button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestNavigationMenuTriggerStyleHelper(t *testing.T) {
	got := ui.NavigationMenuTriggerStyle()
	want := "group inline-flex h-9 w-max items-center justify-center rounded-lg px-4 py-2 text-sm font-medium transition-[color,box-shadow] outline-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1 disabled:pointer-events-none disabled:opacity-50 data-[state=open]:bg-accent/50 data-[state=open]:text-accent-foreground data-[state=open]:hover:bg-accent data-[state=open]:focus:bg-accent"
	if got != want {
		t.Errorf("got %q\nwant %q", got, want)
	}
	if !strings.HasPrefix(got, "group ") {
		t.Errorf("base string must itself start with the group marker (Trigger's own class adds a second, redundant one)\ngot: %s", got)
	}
}

func TestNavigationMenuContentPinned(t *testing.T) {
	got := render(t, ui.NavigationMenuContent(gsx.Raw("x"), nil))
	want := `<div data-slot="navigation-menu-content" data-gsxui-navigation-menu-content popover="manual" data-state="closed" data-side="bottom" class="top-0 left-0 w-full p-2 pr-2.5 md:absolute md:w-auto group-data-[viewport=false]/navigation-menu:top-full group-data-[viewport=false]/navigation-menu:mt-1.5 group-data-[viewport=false]/navigation-menu:overflow-hidden group-data-[viewport=false]/navigation-menu:rounded-lg group-data-[viewport=false]/navigation-menu:border group-data-[viewport=false]/navigation-menu:bg-popover group-data-[viewport=false]/navigation-menu:text-popover-foreground group-data-[viewport=false]/navigation-menu:shadow **:data-[slot=navigation-menu-link]:focus:ring-0 **:data-[slot=navigation-menu-link]:focus:outline-none opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:opacity-0 starting:open:scale-95 data-[side=bottom]:starting:open:-translate-y-2 data-[side=left]:starting:open:translate-x-2 data-[side=right]:starting:open:-translate-x-2 data-[side=top]:starting:open:translate-y-2">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestNavigationMenuViewportPinned(t *testing.T) {
	got := render(t, ui.NavigationMenuViewport(nil))
	want := `<div class="absolute top-full left-0 isolate z-50 flex justify-center"><div data-slot="navigation-menu-viewport" data-gsxui-navigation-menu-viewport popover="manual" data-state="closed" data-side="bottom" class="origin-top relative mt-1.5 w-full overflow-hidden rounded-lg border bg-popover text-popover-foreground shadow md:w-auto opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:opacity-0 starting:open:scale-95 data-[side=bottom]:starting:open:-translate-y-2 data-[side=left]:starting:open:translate-x-2 data-[side=right]:starting:open:-translate-x-2 data-[side=top]:starting:open:translate-y-2"></div></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestNavigationMenuLinkPinned(t *testing.T) {
	got := render(t, ui.NavigationMenuLink(false, gsx.Raw("Docs"), nil))
	want := `<a data-slot="navigation-menu-link" data-gsxui-navigation-menu-link data-active="false" class="flex flex-col gap-1 rounded-lg p-2 text-sm transition-all outline-none in-data-[slot=navigation-menu-content]:rounded-md hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1 data-[active=true]:bg-accent/50 data-[active=true]:text-accent-foreground data-[active=true]:hover:bg-accent data-[active=true]:focus:bg-accent [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 [&amp;_svg:not([class*=&#39;text-&#39;])]:text-muted-foreground">Docs</a>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestNavigationMenuLinkActive(t *testing.T) {
	got := render(t, ui.NavigationMenuLink(true, gsx.Raw("Docs"), nil))
	for _, want := range []string{`data-active="true"`} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestNavigationMenuIndicatorPinned(t *testing.T) {
	got := render(t, ui.NavigationMenuIndicator(nil))
	want := `<div data-slot="navigation-menu-indicator" data-gsxui-navigation-menu-indicator data-state="hidden" class="top-full z-[1] flex h-1.5 items-end justify-center overflow-hidden opacity-0 transition-opacity duration-200 data-[state=visible]:opacity-100"><div class="relative top-[60%] h-2 w-2 rotate-45 rounded-tl-sm bg-border shadow-md"></div></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestNavigationMenuCallerClassMerges(t *testing.T) {
	got := render(t, ui.NavigationMenuContent(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "w-96"}}))
	if strings.Contains(got, "w-full") {
		t.Errorf("caller w-96 must drop default w-full\nin: %s", got)
	}
	if !strings.Contains(got, "w-96") || !strings.Contains(got, "p-2") {
		t.Errorf("want w-96 plus surviving structural classes\nin: %s", got)
	}
}

func TestNavigationMenuAttrsFallThrough(t *testing.T) {
	got := render(t, ui.NavigationMenuTrigger(gsx.Raw("x"), gsx.Attrs{{Key: "aria-label", Value: "Open products"}}))
	if !strings.Contains(got, `aria-label="Open products"`) {
		t.Errorf("missing aria-label\nin: %s", got)
	}
}

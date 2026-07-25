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
	// FIX ROUND 1: the shared-viewport mode (and its `viewport` param) is
	// gone — see ui/navigation-menu.gsx's own file header GAP paragraph. The
	// root no longer takes a viewport param, stamps no data-viewport
	// attribute, and never auto-renders a NavigationMenuViewport child.
	got := render(t, ui.NavigationMenu(gsx.Raw("x"), nil))
	want := `<nav data-slot="navigation-menu" data-gsxui-navigation-menu class="group/navigation-menu relative flex max-w-max flex-1 items-center justify-center">x</nav>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestNavigationMenuViewportPartRemoved(t *testing.T) {
	// FIX ROUND 1: a v1 shipped NavigationMenuViewport as a second,
	// coincident popover meant to sit behind the active Content — measured
	// live, it painted an OPAQUE surface directly OVER the content instead
	// (both are promoted to the top layer; whichever is shown second wins
	// the paint order), making every link inside unreachable. The part is
	// removed outright, not kept as a no-op — the viewport={false}
	// configuration this port now ships never renders it at all in
	// upstream shadcn either. This test asserts the removal stays
	// removed: no navigation-menu-viewport slot exists anywhere a fully
	// composed menu renders (see TestNavigationMenuOnlyOnePopoverPerItem
	// below for the "nothing occludes Content" half of this guarantee).
	got := render(t, ui.NavigationMenuItem(gsx.Fragment(
		ui.NavigationMenuTrigger(gsx.Raw("Products"), nil),
		ui.NavigationMenuContent(gsx.Raw("links"), nil),
	), nil))
	if strings.Contains(got, "navigation-menu-viewport") {
		t.Errorf("navigation-menu-viewport must not appear anywhere — the part is removed\nin: %s", got)
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
	// FIX ROUND 1: the border/bg/shadow/rounded chrome is now UNCONDITIONAL
	// — no group-data-[viewport=false]/navigation-menu: gate, since there is
	// only one mode. This IS the "each Content carries the panel chrome"
	// pin the coordinator asked for: bg-popover/border/shadow/rounded-lg all
	// appear in the base class, live on every render, not behind a selector
	// that depends on an ancestor state this port never stamps.
	got := render(t, ui.NavigationMenuContent(gsx.Raw("x"), nil))
	// tailwind-merge drops the earlier same-group "top-0" in favor of the
	// later "top-full" (both are `top-*` utilities) — a real conflict
	// resolution, not a bug: the class string's own authoring order (block
	// 1's top-0 written before block 2's top-full) already decided the
	// winner, matching what the ORIGINAL two-block shadcn source already
	// did whenever its own viewport=false ancestor state was active.
	want := `<div data-slot="navigation-menu-content" data-gsxui-navigation-menu-content popover="manual" data-state="closed" data-side="bottom" class="left-0 w-full p-2 pr-2.5 md:absolute md:w-auto top-full mt-1.5 overflow-hidden rounded-lg border bg-popover text-popover-foreground shadow **:data-[slot=navigation-menu-link]:focus:ring-0 **:data-[slot=navigation-menu-link]:focus:outline-none opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:opacity-0 starting:open:scale-95 data-[side=bottom]:starting:open:-translate-y-2 data-[side=left]:starting:open:translate-x-2 data-[side=right]:starting:open:-translate-x-2 data-[side=top]:starting:open:translate-y-2">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
	for _, chrome := range []string{"bg-popover", "border", "shadow", "rounded-lg", "overflow-hidden"} {
		if !strings.Contains(got, chrome) {
			t.Errorf("Content must carry its own panel chrome unconditionally — missing %q\nin: %s", chrome, got)
		}
	}
}

// TestNavigationMenuOnlyOnePopoverPerItem is the "nothing occludes Content"
// half of FIX ROUND 1's verification: a fully composed Trigger+Content pair
// (the only two parts that can ever be popover-backed in this port now that
// NavigationMenuViewport is gone) renders exactly one popover="..." — no
// second top-layer element exists that could paint an opaque surface over
// Content the way the removed Viewport did.
func TestNavigationMenuOnlyOnePopoverPerItem(t *testing.T) {
	got := render(t, ui.NavigationMenuItem(gsx.Fragment(
		ui.NavigationMenuTrigger(gsx.Raw("Products"), nil),
		ui.NavigationMenuContent(gsx.Raw(`<a>Dialog</a>`), nil),
	), nil))
	if n := strings.Count(got, `popover="`); n != 1 {
		t.Errorf("want exactly 1 popover-backed element (Content itself), got %d\nin: %s", n, got)
	}
	if !strings.Contains(got, "Dialog") {
		t.Errorf("Content's own children must render, unoccluded\nin: %s", got)
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

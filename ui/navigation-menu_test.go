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
	// FIX ROUND 2: data-viewport="false" is stamped by default (a FIX ROUND 1
	// version stamped nothing at all) — this keeps NavigationMenuContent's
	// own group-data-[viewport=false]/navigation-menu: chrome GATE live
	// (verbatim against upstream's own selector prefix; the tokens behind it
	// are their own ledgered ADAPTs, not byte-identical — see
	// ui/navigation-menu.gsx's own doc comments). Like every server-authored
	// attribute here, a caller-supplied data-viewport overrides this default
	// (codegen's standard attrs.Has hatch), which would silently strip
	// Content's chrome since it's gated on this exact attribute.
	got := render(t, ui.NavigationMenu(gsx.Raw("x"), nil))
	want := `<nav data-slot="navigation-menu" data-gsxui-navigation-menu data-viewport="false" class="group/navigation-menu relative flex max-w-max flex-1 items-center justify-center">x</nav>`
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
	// upstream shadcn either.
	got := render(t, ui.NavigationMenuItem(gsx.Fragment(
		ui.NavigationMenuTrigger(gsx.Raw("Products"), nil),
		ui.NavigationMenuContent(gsx.Raw("links"), nil),
	), nil))
	if strings.Contains(got, "navigation-menu-viewport") {
		t.Errorf("navigation-menu-viewport must not appear anywhere — the part is removed\nin: %s", got)
	}
}

func TestNavigationMenuListPinned(t *testing.T) {
	// gap-0 is nova's own metric (.cn-navigation-menu-list), FIX ROUND 2 —
	// previously omitted (new-york-v4's own gap-1 was carried instead).
	got := render(t, ui.NavigationMenuList(gsx.Raw("x"), nil))
	want := `<ul data-slot="navigation-menu-list" data-gsxui-navigation-menu-list class="group flex flex-1 list-none items-center justify-center gap-0">x</ul>`
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
	// FIX ROUND 2: rounded-lg/px-2.5 py-1.5/transition-all are nova's own
	// metrics (previously the port under-adopted them, carrying
	// new-york-v4's own px-4 py-2/transition-[color,box-shadow] instead).
	// A literal space text node now separates {children} from the chevron,
	// matching new-york-v4's own {children}{" "} (navigation-menu.tsx:76).
	got := render(t, ui.NavigationMenuTrigger(gsx.Raw("Products"), nil))
	want := `<button data-slot="navigation-menu-trigger" data-gsxui-navigation-menu-trigger type="button" aria-expanded="false" data-state="closed" class="group inline-flex h-9 w-max items-center justify-center rounded-lg px-2.5 py-1.5 text-sm font-medium transition-all outline-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1 disabled:pointer-events-none disabled:opacity-50 data-[state=open]:bg-accent/50 data-[state=open]:text-accent-foreground data-[state=open]:hover:bg-accent data-[state=open]:focus:bg-accent group">Products <svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="relative top-[1px] ml-1 size-3 transition duration-300 group-data-[state=open]:rotate-180" data-gsxui-slot="icon"><path d="m6 9 6 6 6-6"/></svg></button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestNavigationMenuTriggerStyleHelper(t *testing.T) {
	got := ui.NavigationMenuTriggerStyle()
	want := "group inline-flex h-9 w-max items-center justify-center rounded-lg px-2.5 py-1.5 text-sm font-medium transition-all outline-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1 disabled:pointer-events-none disabled:opacity-50 data-[state=open]:bg-accent/50 data-[state=open]:text-accent-foreground data-[state=open]:hover:bg-accent data-[state=open]:focus:bg-accent"
	if got != want {
		t.Errorf("got %q\nwant %q", got, want)
	}
	if !strings.HasPrefix(got, "group ") {
		t.Errorf("base string must itself start with the group marker (Trigger's own class adds a second, redundant one)\ngot: %s", got)
	}
}

func TestNavigationMenuContentPinned(t *testing.T) {
	// FIX ROUND 2: the chrome block's group-data-[viewport=false]/
	// navigation-menu: prefix is back (NavigationMenu now stamps
	// data-viewport="false" by default, so the selector is live absent a
	// caller override — see TestNavigationMenuRootPinned), and p-1 replaces
	// new-york-v4's own p-2 pr-2.5 (nova metric, previously omitted).
	// w-full/md:w-auto are GONE — the CRITICAL mobile-overflow fix: a
	// fixed-positioned element's containing block is the viewport, not its
	// DOM ancestor, so w-full below md resolved to ~100vw. See
	// TestNavigationMenuContentHasNoFixedWidthUtility below for the
	// regression guard on that fix specifically.
	got := render(t, ui.NavigationMenuContent(gsx.Raw("x"), nil))
	want := `<div data-slot="navigation-menu-content" data-gsxui-navigation-menu-content popover="manual" data-state="closed" data-side="bottom" class="top-0 left-0 p-1 md:absolute group-data-[viewport=false]/navigation-menu:top-full group-data-[viewport=false]/navigation-menu:mt-1.5 group-data-[viewport=false]/navigation-menu:overflow-hidden group-data-[viewport=false]/navigation-menu:rounded-lg group-data-[viewport=false]/navigation-menu:border group-data-[viewport=false]/navigation-menu:bg-popover group-data-[viewport=false]/navigation-menu:text-popover-foreground group-data-[viewport=false]/navigation-menu:shadow **:data-[slot=navigation-menu-link]:focus:ring-0 **:data-[slot=navigation-menu-link]:focus:outline-none opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:opacity-0 starting:open:scale-95 data-[side=bottom]:starting:open:-translate-y-2 data-[side=left]:starting:open:translate-x-2 data-[side=right]:starting:open:-translate-x-2 data-[side=top]:starting:open:translate-y-2">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
	for _, chrome := range []string{"bg-popover", "border", "shadow", "rounded-lg", "overflow-hidden"} {
		if !strings.Contains(got, chrome) {
			t.Errorf("Content must carry its own panel chrome — missing %q\nin: %s", chrome, got)
		}
	}
}

// TestNavigationMenuContentHasNoFixedWidthUtility is the CRITICAL FIX ROUND 2
// regression guard: NavigationMenuContent must never carry a `w-full` (or
// any bare `w-*`) utility. ui/navigation-menu.js sets position:fixed on
// Content at every breakpoint — for a fixed element the containing block is
// the viewport itself, not the NavigationMenuItem <li> — so a `w-full`
// authored to mean "100% of the li" below new-york-v4's own md breakpoint
// instead resolves to ~100% of the VIEWPORT, overflowing the right edge and
// forcing horizontal page scroll (measured live on the built site before
// this fix). Go can't execute CSS layout, so this test asserts the
// render-contract half of the fix — the utility token itself is gone — not
// the computed layout; the file header's own doc comment on
// NavigationMenuContent records the live measurement.
func TestNavigationMenuContentHasNoFixedWidthUtility(t *testing.T) {
	got := render(t, ui.NavigationMenuContent(gsx.Raw("x"), nil))
	for _, dead := range []string{"w-full", "md:w-auto"} {
		if strings.Contains(got, dead) {
			t.Errorf("Content must not carry %q — fixed positioning makes it resolve against the viewport, not the trigger's own <li>\nin: %s", dead, got)
		}
	}
}

// TestNavigationMenuOnlyOnePopoverPerItem is the "nothing occludes Content"
// half of FIX ROUND 1's verification, corrected in FIX ROUND 2 to actually
// guard the bug it names: an earlier version of this test composed only
// NavigationMenuItem(Trigger, Content) directly, which never exercises
// NavigationMenu's own root — but the v1 bug this guards against was the
// ROOT auto-rendering a NavigationMenuViewport child
// ({viewport && <NavigationMenuViewport/>}), a path that composition never
// reaches. Reintroducing that exact regression at the root would have left
// the old test green. This version renders the full
// NavigationMenu>List>Item>Trigger+Content tree and counts popover="..."
// across the WHOLE tree, so a reintroduced root-level Viewport (or any
// other stray popover-backed element) would push the count to 2 and fail.
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
		t.Errorf("want exactly 1 popover-backed element (Content itself) across the whole tree, got %d\nin: %s", n, got)
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
	// FIX ROUND 3: pointer-events-none added — the indicator stays
	// permanently in the layout (only opacity toggles), so once JS
	// absolutizes it under the list it would otherwise be a real,
	// hit-testable strip even while visually hidden at opacity-0.
	got := render(t, ui.NavigationMenuIndicator(nil))
	want := `<div data-slot="navigation-menu-indicator" data-gsxui-navigation-menu-indicator data-state="hidden" class="pointer-events-none top-full z-[1] flex h-1.5 items-end justify-center overflow-hidden opacity-0 transition-opacity duration-200 data-[state=visible]:opacity-100"><div class="relative top-[60%] h-2 w-2 rotate-45 rounded-tl-sm bg-border shadow-md"></div></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestNavigationMenuCallerClassMerges(t *testing.T) {
	// Content carries no width utility any more (the CRITICAL fix above), so
	// the merge-drops-the-default check uses padding instead: p-1 is
	// Content's own nova-metric default.
	got := render(t, ui.NavigationMenuContent(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "p-8"}}))
	if strings.Contains(got, "p-1 ") || strings.HasSuffix(got, "p-1") {
		t.Errorf("caller p-8 must drop default p-1\nin: %s", got)
	}
	if !strings.Contains(got, "p-8") || !strings.Contains(got, "md:absolute") {
		t.Errorf("want p-8 plus surviving structural classes\nin: %s", got)
	}
}

func TestNavigationMenuAttrsFallThrough(t *testing.T) {
	got := render(t, ui.NavigationMenuTrigger(gsx.Raw("x"), gsx.Attrs{{Key: "aria-label", Value: "Open products"}}))
	if !strings.Contains(got, `aria-label="Open products"`) {
		t.Errorf("missing aria-label\nin: %s", got)
	}
}

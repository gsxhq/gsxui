package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestSidebarProviderStampsServerState(t *testing.T) {
	// The whole point of the design: state arrives as a parameter and is
	// rendered, so there is no flash and no hydration step.
	open := render(t, ui.SidebarProvider(true, gsx.Raw("x"), nil))
	if !strings.Contains(open, `data-state="expanded"`) {
		t.Errorf("want expanded\nin: %s", open)
	}
	closed := render(t, ui.SidebarProvider(false, gsx.Raw("x"), nil))
	if !strings.Contains(closed, `data-state="collapsed"`) {
		t.Errorf("want collapsed\nin: %s", closed)
	}
}

func TestSidebarProviderCarriesWidthVars(t *testing.T) {
	got := render(t, ui.SidebarProvider(true, gsx.Raw("x"), nil))
	for _, want := range []string{"--sidebar-width", "--sidebar-width-icon", `data-slot="sidebar-wrapper"`} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestSidebarRendersBothTrees(t *testing.T) {
	// Mobile is CSS-gated, not JS-swapped: the Sheet tree and the desktop
	// tree both exist in the DOM, gated md:hidden (on the Sheet root, see
	// TestSidebarMobileGateOnSheetRootNotContent below) / hidden md:block.
	got := render(t, ui.Sidebar(true, "", "", "", gsx.Raw("CONTENT"), nil))
	if strings.Count(got, "CONTENT") != 2 {
		t.Errorf("want children rendered in both trees, got %d\nin: %s", strings.Count(got, "CONTENT"), got)
	}
	for _, want := range []string{`data-mobile="true"`, "md:hidden", "md:block"} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

// TestSidebarMobileGateOnSheetRootNotContent is the regression test for
// review round 1's IMPORTANT 3: md:hidden must live on the Sheet root
// (class="contents md:hidden"), never on SheetContent's own <dialog> class,
// where it lost a specificity fight against the dialog's own open:flex
// rule once the sheet was open — see ui/sidebar.gsx's own package doc
// comment FIX entry for the full trace.
func TestSidebarMobileGateOnSheetRootNotContent(t *testing.T) {
	got := render(t, ui.Sidebar(true, "", "", "", gsx.Raw("x"), nil))
	if !strings.Contains(got, `class="contents md:hidden"`) {
		t.Errorf("want md:hidden on the Sheet root's contents class\nin: %s", got)
	}
	dialogStart := strings.Index(got, "<dialog")
	dialogClassEnd := strings.Index(got[dialogStart:], ">")
	dialogOpenTag := got[dialogStart : dialogStart+dialogClassEnd]
	if strings.Contains(dialogOpenTag, "md:hidden") {
		t.Errorf("md:hidden must NOT be on the <dialog> element itself (loses to open:flex)\nin: %s", dialogOpenTag)
	}
}

func TestSidebarCollapsibleNoneIsFlat(t *testing.T) {
	// The reference short-circuits to one plain div; no gap, no container,
	// no Sheet, children rendered exactly once.
	got := render(t, ui.Sidebar(true, "", "", "none", gsx.Raw("CONTENT"), nil))
	if strings.Count(got, "CONTENT") != 1 {
		t.Errorf("collapsible=none renders children once, got %d\nin: %s", strings.Count(got, "CONTENT"), got)
	}
	if strings.Contains(got, `data-mobile="true"`) {
		t.Errorf("collapsible=none must not render the Sheet tree\nin: %s", got)
	}
}

func TestSidebarStampsVariantSideCollapsible(t *testing.T) {
	got := render(t, ui.Sidebar(true, "right", "floating", "icon", gsx.Raw("x"), nil))
	for _, want := range []string{
		`data-side="right"`,
		`data-variant="floating"`,
		// The one gsxui-invented attribute ui/sidebar.js actually depends
		// on (review round 1, IMPORTANT 4): the constant configured
		// collapsible MODE, stamped separately from data-collapsible
		// (which is state-gated — see the next test) so a collapse ->
		// expand -> collapse cycle has something to restore it from.
		`data-gsxui-sidebar-collapsible="icon"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

// TestSidebarOpenDrivesDesktopState is the regression test for review
// round 1's CRITICAL 1: Sidebar's own desktop root must render
// data-state/data-collapsible FROM `open`, exactly mirroring the
// reference's `state === "collapsed" ? collapsible : ""` ternary — before
// this fix, Sidebar had no `open` parameter at all and the desktop root
// was permanently stuck at data-state="expanded"/data-collapsible="",
// regardless of what SidebarProvider was told.
func TestSidebarOpenDrivesDesktopState(t *testing.T) {
	open := render(t, ui.Sidebar(true, "", "", "icon", gsx.Raw("x"), nil))
	for _, want := range []string{`data-state="expanded"`, `data-collapsible=""`} {
		if !strings.Contains(open, want) {
			t.Errorf("open=true: want %q\nin: %s", want, open)
		}
	}
	closed := render(t, ui.Sidebar(false, "", "", "icon", gsx.Raw("x"), nil))
	for _, want := range []string{`data-state="collapsed"`, `data-collapsible="icon"`} {
		if !strings.Contains(closed, want) {
			t.Errorf("open=false: want %q\nin: %s", want, closed)
		}
	}
	// data-state="expanded" must not appear at all when closed (would mean
	// the ternary regressed to always-expanded).
	if strings.Contains(closed, `data-state="expanded"`) {
		t.Errorf("open=false must not render data-state=expanded anywhere\nin: %s", closed)
	}
}

func TestSidebarTriggerCarriesGsxuiHook(t *testing.T) {
	// data-gsxui-sidebar-trigger is the documented public JS-attachment
	// hook (docs/jsx-parity.md `## dialog` MECHANISM) — sidebar.js matches
	// it alongside data-slot="sidebar-trigger" so a caller's own styled
	// trigger wires up without stamping data-slot.
	got := render(t, ui.SidebarTrigger(nil))
	for _, want := range []string{`data-slot="sidebar-trigger"`, `data-gsxui-sidebar-trigger`} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestSidebarRailPinned(t *testing.T) {
	got := render(t, ui.SidebarRail(nil))
	want := `<button type="button" data-sidebar="rail" data-slot="sidebar-rail" data-gsxui-sidebar-rail aria-label="Toggle Sidebar" tabindex="-1" title="Toggle Sidebar" class="absolute inset-y-0 z-20 hidden w-4 -translate-x-1/2 transition-all ease-linear group-data-[side=left]:-right-4 group-data-[side=right]:left-0 after:absolute after:inset-y-0 after:left-1/2 after:w-[2px] hover:after:bg-sidebar-border sm:flex in-data-[side=left]:cursor-w-resize in-data-[side=right]:cursor-e-resize [[data-side=left][data-state=collapsed]_&amp;]:cursor-e-resize [[data-side=right][data-state=collapsed]_&amp;]:cursor-w-resize group-data-[collapsible=offcanvas]:translate-x-0 group-data-[collapsible=offcanvas]:after:left-full hover:group-data-[collapsible=offcanvas]:bg-sidebar [[data-side=left][data-collapsible=offcanvas]_&amp;]:-right-2 [[data-side=right][data-collapsible=offcanvas]_&amp;]:-left-2 [[data-mobile]_&amp;]:hidden"></button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestSidebarInsetPinned(t *testing.T) {
	got := render(t, ui.SidebarInset(gsx.Raw("x"), nil))
	want := `<main data-slot="sidebar-inset" class="relative flex w-full flex-1 flex-col bg-background md:peer-data-[variant=inset]:m-2 md:peer-data-[variant=inset]:ml-0 md:peer-data-[variant=inset]:rounded-xl md:peer-data-[variant=inset]:shadow-sm md:peer-data-[variant=inset]:peer-data-[state=collapsed]:ml-2">x</main>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestSidebarMenuButtonVariantSizeMatrix exact-match pins every
// variant x size combination (2 x 3 = 6), the full sidebarMenuButtonVariants
// cva surface, no tooltip.
func TestSidebarMenuButtonVariantSizeMatrix(t *testing.T) {
	cases := []struct {
		variant, size string
		want          string
	}{
		{"", "", `<button type="button" data-slot="sidebar-menu-button" data-sidebar="menu-button" data-size="default" data-active="false" class="peer/menu-button flex w-full items-center gap-2 overflow-hidden rounded-md p-2 text-left ring-sidebar-ring outline-hidden transition-[width,height,padding] group-has-data-[sidebar=menu-action]/menu-item:pr-8 group-data-[collapsible=icon]:size-8! group-data-[collapsible=icon]:p-2! focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 data-[active=true]:bg-sidebar-accent data-[active=true]:font-medium data-[active=true]:text-sidebar-accent-foreground data-[state=open]:hover:bg-sidebar-accent data-[state=open]:hover:text-sidebar-accent-foreground [&amp;&gt;span:last-child]:truncate [&amp;&gt;svg]:size-4 [&amp;&gt;svg]:shrink-0 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground h-8 text-sm">x</button>`},
		{"", "sm", `<button type="button" data-slot="sidebar-menu-button" data-sidebar="menu-button" data-size="sm" data-active="false" class="peer/menu-button flex w-full items-center gap-2 overflow-hidden rounded-md p-2 text-left ring-sidebar-ring outline-hidden transition-[width,height,padding] group-has-data-[sidebar=menu-action]/menu-item:pr-8 group-data-[collapsible=icon]:size-8! group-data-[collapsible=icon]:p-2! focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 data-[active=true]:bg-sidebar-accent data-[active=true]:font-medium data-[active=true]:text-sidebar-accent-foreground data-[state=open]:hover:bg-sidebar-accent data-[state=open]:hover:text-sidebar-accent-foreground [&amp;&gt;span:last-child]:truncate [&amp;&gt;svg]:size-4 [&amp;&gt;svg]:shrink-0 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground h-7 text-xs">x</button>`},
		{"", "lg", `<button type="button" data-slot="sidebar-menu-button" data-sidebar="menu-button" data-size="lg" data-active="false" class="peer/menu-button flex w-full items-center gap-2 overflow-hidden rounded-md p-2 text-left ring-sidebar-ring outline-hidden transition-[width,height,padding] group-has-data-[sidebar=menu-action]/menu-item:pr-8 group-data-[collapsible=icon]:size-8! focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 data-[active=true]:bg-sidebar-accent data-[active=true]:font-medium data-[active=true]:text-sidebar-accent-foreground data-[state=open]:hover:bg-sidebar-accent data-[state=open]:hover:text-sidebar-accent-foreground [&amp;&gt;span:last-child]:truncate [&amp;&gt;svg]:size-4 [&amp;&gt;svg]:shrink-0 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground h-12 text-sm group-data-[collapsible=icon]:p-0!">x</button>`},
		{"outline", "", `<button type="button" data-slot="sidebar-menu-button" data-sidebar="menu-button" data-size="default" data-active="false" class="peer/menu-button flex w-full items-center gap-2 overflow-hidden rounded-md p-2 text-left ring-sidebar-ring outline-hidden transition-[width,height,padding] group-has-data-[sidebar=menu-action]/menu-item:pr-8 group-data-[collapsible=icon]:size-8! group-data-[collapsible=icon]:p-2! focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 data-[active=true]:bg-sidebar-accent data-[active=true]:font-medium data-[active=true]:text-sidebar-accent-foreground data-[state=open]:hover:bg-sidebar-accent data-[state=open]:hover:text-sidebar-accent-foreground [&amp;&gt;span:last-child]:truncate [&amp;&gt;svg]:size-4 [&amp;&gt;svg]:shrink-0 bg-background shadow-[0_0_0_1px_var(--sidebar-border)] hover:bg-sidebar-accent hover:text-sidebar-accent-foreground hover:shadow-[0_0_0_1px_var(--sidebar-accent)] h-8 text-sm">x</button>`},
		{"outline", "sm", `<button type="button" data-slot="sidebar-menu-button" data-sidebar="menu-button" data-size="sm" data-active="false" class="peer/menu-button flex w-full items-center gap-2 overflow-hidden rounded-md p-2 text-left ring-sidebar-ring outline-hidden transition-[width,height,padding] group-has-data-[sidebar=menu-action]/menu-item:pr-8 group-data-[collapsible=icon]:size-8! group-data-[collapsible=icon]:p-2! focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 data-[active=true]:bg-sidebar-accent data-[active=true]:font-medium data-[active=true]:text-sidebar-accent-foreground data-[state=open]:hover:bg-sidebar-accent data-[state=open]:hover:text-sidebar-accent-foreground [&amp;&gt;span:last-child]:truncate [&amp;&gt;svg]:size-4 [&amp;&gt;svg]:shrink-0 bg-background shadow-[0_0_0_1px_var(--sidebar-border)] hover:bg-sidebar-accent hover:text-sidebar-accent-foreground hover:shadow-[0_0_0_1px_var(--sidebar-accent)] h-7 text-xs">x</button>`},
		{"outline", "lg", `<button type="button" data-slot="sidebar-menu-button" data-sidebar="menu-button" data-size="lg" data-active="false" class="peer/menu-button flex w-full items-center gap-2 overflow-hidden rounded-md p-2 text-left ring-sidebar-ring outline-hidden transition-[width,height,padding] group-has-data-[sidebar=menu-action]/menu-item:pr-8 group-data-[collapsible=icon]:size-8! focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 data-[active=true]:bg-sidebar-accent data-[active=true]:font-medium data-[active=true]:text-sidebar-accent-foreground data-[state=open]:hover:bg-sidebar-accent data-[state=open]:hover:text-sidebar-accent-foreground [&amp;&gt;span:last-child]:truncate [&amp;&gt;svg]:size-4 [&amp;&gt;svg]:shrink-0 bg-background shadow-[0_0_0_1px_var(--sidebar-border)] hover:bg-sidebar-accent hover:text-sidebar-accent-foreground hover:shadow-[0_0_0_1px_var(--sidebar-accent)] h-12 text-sm group-data-[collapsible=icon]:p-0!">x</button>`},
	}
	for _, c := range cases {
		got := render(t, ui.SidebarMenuButton(false, c.variant, c.size, "", gsx.Raw("x"), nil))
		if got != c.want {
			t.Errorf("variant=%q size=%q pinned render mismatch\n got: %s\nwant: %s", c.variant, c.size, got, c.want)
		}
	}
}

func TestSidebarMenuActionPinned(t *testing.T) {
	got := render(t, ui.SidebarMenuAction(false, gsx.Raw("x"), nil))
	want := `<button type="button" data-slot="sidebar-menu-action" data-sidebar="menu-action" class="absolute top-1.5 right-1 flex aspect-square w-5 items-center justify-center rounded-md p-0 text-sidebar-foreground ring-sidebar-ring outline-hidden transition-transform peer-hover/menu-button:text-sidebar-accent-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 [&amp;&gt;svg]:size-4 [&amp;&gt;svg]:shrink-0 after:absolute after:-inset-2 md:after:hidden peer-data-[size=sm]/menu-button:top-1 peer-data-[size=default]/menu-button:top-1.5 peer-data-[size=lg]/menu-button:top-2.5 group-data-[collapsible=icon]:hidden">x</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestSidebarMenuActionShowOnHoverPinned(t *testing.T) {
	got := render(t, ui.SidebarMenuAction(true, gsx.Raw("x"), nil))
	want := `<button type="button" data-slot="sidebar-menu-action" data-sidebar="menu-action" class="absolute top-1.5 right-1 flex aspect-square w-5 items-center justify-center rounded-md p-0 text-sidebar-foreground ring-sidebar-ring outline-hidden transition-transform peer-hover/menu-button:text-sidebar-accent-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 [&amp;&gt;svg]:size-4 [&amp;&gt;svg]:shrink-0 after:absolute after:-inset-2 md:after:hidden peer-data-[size=sm]/menu-button:top-1 peer-data-[size=default]/menu-button:top-1.5 peer-data-[size=lg]/menu-button:top-2.5 group-data-[collapsible=icon]:hidden group-focus-within/menu-item:opacity-100 group-hover/menu-item:opacity-100 peer-data-[active=true]/menu-button:text-sidebar-accent-foreground data-[state=open]:opacity-100 md:opacity-0">x</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
	if !strings.Contains(got, "md:opacity-0") {
		t.Errorf("showOnHover=true must add the md:opacity-0 block\nin: %s", got)
	}
}

func TestSidebarMenuSubButtonPinned(t *testing.T) {
	sm := render(t, ui.SidebarMenuSubButton("sm", false, gsx.Raw("x"), nil))
	wantSm := `<a data-slot="sidebar-menu-sub-button" data-sidebar="menu-sub-button" data-size="sm" data-active="false" class="flex h-7 min-w-0 -translate-x-px items-center gap-2 overflow-hidden rounded-md px-2 text-sidebar-foreground ring-sidebar-ring outline-hidden hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 [&amp;&gt;span:last-child]:truncate [&amp;&gt;svg]:size-4 [&amp;&gt;svg]:shrink-0 [&amp;&gt;svg]:text-sidebar-accent-foreground data-[active=true]:bg-sidebar-accent data-[active=true]:text-sidebar-accent-foreground text-xs group-data-[collapsible=icon]:hidden">x</a>`
	if sm != wantSm {
		t.Errorf("size=sm pinned render mismatch\n got: %s\nwant: %s", sm, wantSm)
	}

	md := render(t, ui.SidebarMenuSubButton("md", true, gsx.Raw("x"), nil))
	wantMd := `<a data-slot="sidebar-menu-sub-button" data-sidebar="menu-sub-button" data-size="md" data-active="true" class="flex h-7 min-w-0 -translate-x-px items-center gap-2 overflow-hidden rounded-md px-2 text-sidebar-foreground ring-sidebar-ring outline-hidden hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 [&amp;&gt;span:last-child]:truncate [&amp;&gt;svg]:size-4 [&amp;&gt;svg]:shrink-0 [&amp;&gt;svg]:text-sidebar-accent-foreground data-[active=true]:bg-sidebar-accent data-[active=true]:text-sidebar-accent-foreground text-sm group-data-[collapsible=icon]:hidden">x</a>`
	if md != wantMd {
		t.Errorf("size=md active pinned render mismatch\n got: %s\nwant: %s", md, wantMd)
	}
}

func TestSidebarMenuButtonActiveAndTooltip(t *testing.T) {
	got := render(t, ui.SidebarMenuButton(true, "", "", "Inbox", gsx.Raw("Inbox"), nil))
	for _, want := range []string{
		`data-slot="sidebar-menu-button"`,
		`data-active="true"`,
		// A non-empty tooltip must wrap in ui.Tooltip (data-slot="tooltip")
		// and put data-gsxui-tooltip-trigger directly on the button itself
		// — no nested ui.TooltipTrigger, the same button-in-button trap as
		// DialogTrigger (see ui/sidebar.gsx's own MECHANISM doc comment).
		`data-slot="tooltip"`,
		`data-gsxui-tooltip-trigger`,
		`data-slot="tooltip-content"`,
		">Inbox<",
		// The collapsed-AND-desktop-only gating has no state/isMobile param
		// here, so it's pure CSS: only visible while an ancestor .group
		// carries data-collapsible=icon — true for the desktop copy of
		// this button when icon-collapsed, never true for the mobile
		// Sheet copy (review round 1, IMPORTANT 4).
		// FIX (review round 2, BLOCKING): gated on :open too, matching the
		// `## dialog` ADAPT idiom (`open:grid`) — TooltipContent stays in
		// the DOM while closed, so an ungated display utility would beat
		// the UA's closed-popover display:none and render a ghost hit box.
		"hidden group-data-[collapsible=icon]:open:block",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestSidebarMenuSkeletonShowIcon(t *testing.T) {
	with := render(t, ui.SidebarMenuSkeleton(true, nil))
	without := render(t, ui.SidebarMenuSkeleton(false, nil))
	if strings.Count(with, `data-slot="skeleton"`) <= strings.Count(without, `data-slot="skeleton"`) {
		t.Errorf("showIcon must add a skeleton\n with: %s\nwithout: %s", with, without)
	}
}

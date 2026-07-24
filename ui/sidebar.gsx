package ui

import (
	"math/rand"
	"strconv"

	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Sidebar and its ~25 parts are the shadcn/ui Sidebar (registry/new-york-v4/
// ui/sidebar.tsx, 727 lines) — the widest single component in gsxui, but
// shallow: almost every part below is a plain markup wrapper carrying one
// verbatim class string. The two places with real design work are STATE and
// MOBILE (2026-07-24 tier4 source map, `## sidebar`); everything else is a
// straight port.
//
// 1. STATE. SidebarProvider takes `open bool` and stamps
// data-state="expanded"|"collapsed" server-side on its own sidebar-wrapper
// div — no flash, no hydration step, and no cookie/localStorage inside the
// component (shadcn's own React version mirrors state into a
// non-HttpOnly `sidebar_state` cookie on every toggle; gsxui deliberately
// does not). sidebar.js flips that data-state on click and emits
// gsxui:change ({ open }) via ui/gsxui.js's shared emit() helper — it never
// touches document.cookie or localStorage. Persisting across reloads is the
// CALLER's job: site/examples/sidebar/persisted.gsx is the documented
// cookie round-trip (a Go handler reading sidebar_state into `open`, plus a
// gsxui:change listener writing it back — target-guarded against
// [data-slot="sidebar-wrapper"], since gsxui:change is also emitted by
// tabs/toggle/toggle-group/resizable; see that example's own doc comment).
//
// FIX (review round 1, CRITICAL — supersedes the task brief's own Sidebar
// signature, which omitted `open`): every state-driven visual in this
// component — the gap's width, the container's offcanvas slide, icon-
// collapsed sizing, SidebarInset's peer-data-[state=collapsed]:ml-2, the
// rail's cursor, and all eight group-data-[collapsible=...] consumers —
// roots at the DESKTOP ROOT div, never at SidebarProvider's own wrapper
// (confirmed: no selector anywhere in this component reads
// [data-slot=sidebar-wrapper]'s own data-state). A Sidebar with no `open`
// parameter therefore had no way to ever render collapsed — not a flash,
// a PERMANENT bug: SidebarProvider(false, ...) rendered a fully expanded
// sidebar, full stop, and the first click/rail-press/Cmd+B had to fight
// its way past a `data-state="expanded"` it was never told to start from
// (sidebar.js's toggle reads the CURRENT data-state to decide the new one,
// so a wrongly-expanded root swallowed the first press entirely). `Sidebar`
// now takes `open bool` as its first parameter — the same value the
// caller already passes to `SidebarProvider`, following this codebase's
// established explicit-param pattern for values React context would
// otherwise thread (see ui/toggle-group.gsx's own doc comment on
// `groupType`/`variant`/`size`/`spacing` — "the caller, which already has
// the value in scope, passes it explicitly"). `data-state` renders from
// `open` directly; `data-collapsible` is the configured mode when
// collapsed, `""` when expanded — exactly the reference's own
// `state === "collapsed" ? collapsible : ""` ternary. The Go-invisible
// configured collapse MODE ("offcanvas"/"icon"/"none") is separately
// stamped as data-gsxui-sidebar-collapsible so sidebar.js can restore it
// on a later collapse after the user has expanded it client-side (see
// MECHANISM below) — that attribute is unaffected by this fix.
//
// 2. MOBILE. React swaps the whole sidebar for a Sheet via a client-only
// useIsMobile() media query; a server render cannot branch on viewport.
// Sidebar instead renders BOTH trees unconditionally and lets CSS gate
// them: the desktop tree keeps shadcn's own `hidden … md:block` root and
// `hidden … md:flex` container, and the mobile tree is the Sheet branch,
// carrying the reference's own data-mobile="true", the
// --sidebar-width: 18rem override, and the sr-only header/title/
// description. No media-query JS anywhere. collapsible="none" still
// short-circuits to the reference's single flat div branch, Sheet and all.
//
// FIX (review round 1, mobile gate specificity): the `md:hidden` that
// CSS-gates the mobile tree (a gsxui ADAPT — upstream never needs it,
// since it only ever renders one branch) lives on the Sheet ROOT
// (`<Sheet class="md:hidden">`, a plain `class="contents"` div with no
// competing display utility of its own), NOT on the `<dialog>` element
// itself. It was originally on SheetContent's own class, where it lost to
// the dialog's own `open:flex` (compiled `:is([open],:popover-open,
// :open){display:flex}`, specificity (0,2,0) against `.md\:hidden`'s
// (0,1,0)) — opening the sheet at a narrow viewport and then resizing to
// >=md left BOTH trees visible at once, since a media query adds no
// specificity to win that fight. An ancestor `display:none` still hides a
// `<dialog>` promoted to the top layer via showModal() (top-layer
// promotion escapes stacking/overflow/clipping containment, not ordinary
// CSS box generation from an ancestor's own `display`), and the Sheet
// root's plain `contents`/`md:hidden` pair has no such specificity fight
// to lose, since neither `[open]`/`:popover-open`/`:open` pseudo-classes
// ever match a plain wrapper div.
//
// CONSEQUENCE (children renders TWICE — read this before nesting anything
// with an id inside a Sidebar): every non-"none" Sidebar renders `children`
// once inside the mobile Sheet and once inside the desktop tree. This is a
// known, accepted cost of the CSS-gated dual-tree design, not a bug. DO NOT
// put `id` attributes anywhere inside a Sidebar's children (or in attrs
// passed to Sidebar/SidebarProvider themselves, which forward onto both the
// Sheet root and the desktop container) — they would be duplicated and the
// document invalid. Use classes and data-* attributes instead. Ledgered as
// a GAP in docs/jsx-parity.md `## sidebar`.
//
// 5. NO CONTEXT SUBSTITUTE. data-state/data-collapsible/data-variant/
// data-side are stamped on the desktop root and consumed by descendants
// entirely through group-data-*/peer-data-* selectors — already pure CSS in
// the reference, so nothing needs threading through Go params.
// SidebarTrigger and SidebarRail resolve "the provider" from the DOM
// (closest('[data-slot="sidebar-wrapper"]')) at click time, not an
// imperative handle — matching dialog.js's own proximity-wiring
// convention.
//
// 6. Cmd/Ctrl+B toggles (SIDEBAR_KEYBOARD_SHORTCUT), registered once at
// module scope in sidebar.js, guarded against key-repeat, preventDefault'd.
//
// MECHANISM (upstream's six constants): SIDEBAR_WIDTH/-_MOBILE/-_ICON
// become the unexported sidebarWidth/sidebarWidthMobile/sidebarWidthIcon
// consts below; SIDEBAR_COOKIE_NAME/-_MAX_AGE are not gsxui constants at
// all (no cookie in the component, see decision 1 above) — the cookie name
// only appears, as a literal, in site/examples/sidebar/persisted.gsx's own
// Go handler snippet. SIDEBAR_KEYBOARD_SHORTCUT lives in sidebar.js only.
//
// ADAPT (TooltipProvider dropped, not a gap in this port): the source map's
// own dependency list names Tooltip+TooltipContent+TooltipProvider+
// TooltipTrigger, but gsxui's ui/tooltip.gsx never ports TooltipProvider at
// all (see docs/jsx-parity.md `## tooltip` GAP: "shared delayDuration/
// skip-delay-group machinery across multiple tooltips is not ported —
// tooltip.js hard-codes a 300ms open delay per trigger"). SidebarProvider
// therefore does NOT wrap children in a TooltipProvider (there is nothing
// to wrap them in) — a correction to the source map's dependency framing,
// not a new gap this component introduces.
//
// asChild/Slot polymorphism (SidebarGroupLabel/SidebarGroupAction/
// SidebarMenuButton/SidebarMenuSubButton in the reference) is dropped
// throughout, same as every other gsxui port — each part always renders
// its own tag.
const (
	sidebarWidth       = "16rem"
	sidebarWidthMobile = "18rem"
	sidebarWidthIcon   = "3rem"
)

// SidebarProvider wraps a page's sidebar + inset content in the
// data-slot="sidebar-wrapper" div every group-data-*/peer-data-* selector
// in this component ultimately roots against — see the package doc
// comment's STATE section. No TooltipProvider wrapper (see the package doc
// comment's own ADAPT) — gsxui's Tooltip needs none.
component SidebarProvider(open bool, children gsx.Node, attrs gsx.Attrs) {
	{{
		state := "collapsed"
		if open {
			state = "expanded"
		}
	}}
	<div
		data-slot="sidebar-wrapper"
		data-state={state}
		style={ css`--sidebar-width:@{sidebarWidth}`, css`--sidebar-width-icon:@{sidebarWidthIcon}` }
		class="group/sidebar-wrapper flex min-h-svh w-full has-data-[variant=inset]:bg-sidebar"
		{ attrs... }
	>
		{ children }
	</div>
}

// Sidebar renders the reference's three branches — see the package doc
// comment's MOBILE section for why the offcanvas/icon branches below are
// TWO trees (mobile Sheet + desktop group/peer), not the upstream
// isMobile-gated either/or, and its own FIX entry for why `open` is now
// this component's first parameter, matching SidebarProvider's own value.
// collapsible="none" still renders the reference's single flat div,
// unconditionally, exactly once, and ignores `open` entirely (matching the
// reference, whose own collapsible="none" branch never reads `state`).
component Sidebar(open bool, side string, variant string, collapsible string, children gsx.Node, attrs gsx.Attrs) {
	{{
		s := side
		if s == "" {
			s = "left"
		}
		v := variant
		if v == "" {
			v = "sidebar"
		}
		c := collapsible
		if c == "" {
			c = "offcanvas"
		}
		state := "collapsed"
		collapsibleAttr := c
		if open {
			state = "expanded"
			collapsibleAttr = ""
		}
	}}
	{ if collapsible == "none" {
		<div
			data-slot="sidebar"
			class="flex h-full w-(--sidebar-width) flex-col bg-sidebar text-sidebar-foreground"
			{ attrs... }
		>
			{ children }
		</div>
	} else {
		<>
			<Sheet class="md:hidden" { attrs... }>
				<SheetContent
					side={s}
					data-slot="sidebar"
					data-sidebar="sidebar"
					data-mobile="true"
					style=css`--sidebar-width:@{sidebarWidthMobile}`
					class="w-(--sidebar-width) bg-sidebar p-0 text-sidebar-foreground [&>button]:hidden"
				>
					<SheetHeader class="sr-only">
						<SheetTitle>Sidebar</SheetTitle>
						<SheetDescription>Displays the mobile sidebar.</SheetDescription>
					</SheetHeader>
					<div class="flex h-full w-full flex-col">{ children }</div>
				</SheetContent>
			</Sheet>
			<div
				class="group peer hidden text-sidebar-foreground md:block"
				data-state={state}
				data-collapsible={collapsibleAttr}
				data-gsxui-sidebar-collapsible={c}
				data-variant={v}
				data-side={s}
				data-slot="sidebar"
			>
				<div
					data-slot="sidebar-gap"
					class={
						"relative w-(--sidebar-width) bg-transparent transition-[width] duration-200 ease-linear",
						"group-data-[collapsible=offcanvas]:w-0",
						"group-data-[side=right]:rotate-180",
						if v == "floating" || v == "inset" {
							"group-data-[collapsible=icon]:w-[calc(var(--sidebar-width-icon)+(--spacing(4)))]"
						} else {
							"group-data-[collapsible=icon]:w-(--sidebar-width-icon)"
						},
					}
				></div>
				<div
					data-slot="sidebar-container"
					class={
						"fixed inset-y-0 z-10 hidden h-svh w-(--sidebar-width) transition-[left,right,width] duration-200 ease-linear md:flex",
						if s == "left" {
							"left-0 group-data-[collapsible=offcanvas]:left-[calc(var(--sidebar-width)*-1)]"
						} else {
							"right-0 group-data-[collapsible=offcanvas]:right-[calc(var(--sidebar-width)*-1)]"
						},
						if v == "floating" || v == "inset" {
							"p-2 group-data-[collapsible=icon]:w-[calc(var(--sidebar-width-icon)+(--spacing(4))+2px)]"
						} else {
							"group-data-[collapsible=icon]:w-(--sidebar-width-icon) group-data-[side=left]:border-r group-data-[side=right]:border-l"
						},
					}
					{ attrs... }
				>
					<div
						data-sidebar="sidebar"
						data-slot="sidebar-inner"
						class="flex h-full w-full flex-col bg-sidebar group-data-[variant=floating]:rounded-lg group-data-[variant=floating]:border group-data-[variant=floating]:border-sidebar-border group-data-[variant=floating]:shadow-sm"
					>
						{ children }
					</div>
				</div>
			</div>
		</>
	} }
}

// SidebarTrigger composes ui.Button (variant="ghost" size="icon") — click
// handling (toggle desktop state or open the mobile Sheet, resolved by
// sidebar.js at click time) is entirely JS; see ui/sidebar.js.
component SidebarTrigger(attrs gsx.Attrs) {
	<Button
		data-sidebar="trigger"
		data-slot="sidebar-trigger"
		variant="ghost"
		size="icon"
		class="size-7"
		{ attrs... }
	>
		<icon.PanelLeft/>
		<span class="sr-only">Toggle Sidebar</span>
	</Button>
}

// SidebarRail is a pointer-only affordance (tabindex="-1", deliberately not
// a keyboard tab stop, matching the reference exactly — no drag is wired,
// only a click) that shares SidebarTrigger's toggle mechanism.
//
// ADAPT (review round 1, MINOR 6 — [[data-mobile]_&]:hidden added, not in
// upstream): callers place SidebarRail as a child of Sidebar, same as
// every shadcn block does (registry/new-york-v4/blocks/sidebar-*/
// components/app-sidebar.tsx all render it as Sidebar's own last child) —
// but `children` renders TWICE in this port (see the package doc
// comment's MOBILE section), so a rail placed that way also lands inside
// the mobile Sheet tree, where it has no desktop group/peer ancestor to
// read data-* off, its click is a no-op (isMobile() routes clicks inside
// the mobile tree to openMobile(), not toggleDesktop()), and its own
// `sm:flex` would make it visibly overlay Sheet content between the sm and
// md breakpoints. `[[data-mobile]_&]:hidden` (the same ancestor-selector
// arbitrary-variant idiom this class already uses for the data-side/
// data-state/data-collapsible combinators above) unconditionally hides it
// whenever nested under a data-mobile ancestor — its compiled selector
// (`[data-mobile] &`, specificity (0,2,0)) beats the plain `sm:flex`
// utility (0,1,0) at every viewport, no !important needed.
component SidebarRail(attrs gsx.Attrs) {
	<button
		type="button"
		data-sidebar="rail"
		data-slot="sidebar-rail"
		aria-label="Toggle Sidebar"
		tabindex="-1"
		title="Toggle Sidebar"
		class="absolute inset-y-0 z-20 hidden w-4 -translate-x-1/2 transition-all ease-linear group-data-[side=left]:-right-4 group-data-[side=right]:left-0 after:absolute after:inset-y-0 after:left-1/2 after:w-[2px] hover:after:bg-sidebar-border sm:flex in-data-[side=left]:cursor-w-resize in-data-[side=right]:cursor-e-resize [[data-side=left][data-state=collapsed]_&]:cursor-e-resize [[data-side=right][data-state=collapsed]_&]:cursor-w-resize group-data-[collapsible=offcanvas]:translate-x-0 group-data-[collapsible=offcanvas]:after:left-full hover:group-data-[collapsible=offcanvas]:bg-sidebar [[data-side=left][data-collapsible=offcanvas]_&]:-right-2 [[data-side=right][data-collapsible=offcanvas]_&]:-left-2 [[data-mobile]_&]:hidden"
		{ attrs... }
	></button>
}

component SidebarInset(children gsx.Node, attrs gsx.Attrs) {
	<main
		data-slot="sidebar-inset"
		class="relative flex w-full flex-1 flex-col bg-background md:peer-data-[variant=inset]:m-2 md:peer-data-[variant=inset]:ml-0 md:peer-data-[variant=inset]:rounded-xl md:peer-data-[variant=inset]:shadow-sm md:peer-data-[variant=inset]:peer-data-[state=collapsed]:ml-2"
		{ attrs... }
	>
		{ children }
	</main>
}

component SidebarInput(attrs gsx.Attrs) {
	<Input data-slot="sidebar-input" data-sidebar="input" class="h-8 w-full bg-background shadow-none" { attrs... }/>
}

component SidebarHeader(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="sidebar-header" data-sidebar="header" class="flex flex-col gap-2 p-2" { attrs... }>
		{ children }
	</div>
}

component SidebarFooter(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="sidebar-footer" data-sidebar="footer" class="flex flex-col gap-2 p-2" { attrs... }>
		{ children }
	</div>
}

component SidebarSeparator(attrs gsx.Attrs) {
	<Separator data-slot="sidebar-separator" data-sidebar="separator" class="mx-2 w-auto bg-sidebar-border" { attrs... }/>
}

// SidebarContent picks up nova's own gap-0 (new-york-v4 has gap-2) plus
// no-scrollbar (new, matching carousel/input-otp's own scrollbar-hiding
// convention) — the two real metric deltas nova's `.cn-sidebar-content`
// carries; see the 2026-07-24 tier4 source map `## sidebar` §6. Nova's
// border→ring swap on the floating inner surface is NOT adopted (house
// standing exception) — see SidebarMenuItem below... actually see the
// sidebar-inner div inside Sidebar, which keeps new-york-v4's `border`
// verbatim.
component SidebarContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="sidebar-content"
		data-sidebar="content"
		class="flex min-h-0 flex-1 flex-col gap-0 overflow-auto no-scrollbar group-data-[collapsible=icon]:overflow-hidden"
		{ attrs... }
	>
		{ children }
	</div>
}

component SidebarGroup(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="sidebar-group" data-sidebar="group" class="relative flex w-full min-w-0 flex-col p-2" { attrs... }>
		{ children }
	</div>
}

// SidebarGroupLabel drops asChild (see the package doc comment) — always a
// <div>.
component SidebarGroupLabel(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="sidebar-group-label"
		data-sidebar="group-label"
		class="flex h-8 shrink-0 items-center rounded-md px-2 text-xs font-medium text-sidebar-foreground/70 ring-sidebar-ring outline-hidden transition-[margin,opacity] duration-200 ease-linear focus-visible:ring-2 [&>svg]:size-4 [&>svg]:shrink-0 group-data-[collapsible=icon]:-mt-8 group-data-[collapsible=icon]:opacity-0"
		{ attrs... }
	>
		{ children }
	</div>
}

// SidebarGroupAction drops asChild (see the package doc comment) — always a
// <button>; type="button" added per this codebase's house convention for
// every native, non-form button (upstream's own plain "button" string tag
// carries no explicit type prop, unlike a Radix Primitive.button, which
// defaults to type="button" internally).
component SidebarGroupAction(children gsx.Node, attrs gsx.Attrs) {
	<button
		type="button"
		data-slot="sidebar-group-action"
		data-sidebar="group-action"
		class="absolute top-3.5 right-3 flex aspect-square w-5 items-center justify-center rounded-md p-0 text-sidebar-foreground ring-sidebar-ring outline-hidden transition-transform hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 [&>svg]:size-4 [&>svg]:shrink-0 after:absolute after:-inset-2 md:after:hidden group-data-[collapsible=icon]:hidden"
		{ attrs... }
	>
		{ children }
	</button>
}

component SidebarGroupContent(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="sidebar-group-content" data-sidebar="group-content" class="w-full text-sm" { attrs... }>
		{ children }
	</div>
}

component SidebarMenu(children gsx.Node, attrs gsx.Attrs) {
	<ul data-slot="sidebar-menu" data-sidebar="menu" class="flex w-full min-w-0 flex-col gap-0" { attrs... }>
		{ children }
	</ul>
}

component SidebarMenuItem(children gsx.Node, attrs gsx.Attrs) {
	<li data-slot="sidebar-menu-item" data-sidebar="menu-item" class="group/menu-item relative" { attrs... }>
		{ children }
	</li>
}

const sidebarMenuButtonBase = "peer/menu-button flex w-full items-center gap-2 overflow-hidden rounded-md p-2 text-left text-sm ring-sidebar-ring outline-hidden transition-[width,height,padding] group-has-data-[sidebar=menu-action]/menu-item:pr-8 group-data-[collapsible=icon]:size-8! group-data-[collapsible=icon]:p-2! hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 data-[active=true]:bg-sidebar-accent data-[active=true]:font-medium data-[active=true]:text-sidebar-accent-foreground data-[state=open]:hover:bg-sidebar-accent data-[state=open]:hover:text-sidebar-accent-foreground [&>span:last-child]:truncate [&>svg]:size-4 [&>svg]:shrink-0"

func sidebarMenuButtonVariantClass(variant string) string {
	switch variant {
	case "outline":
		return "bg-background shadow-[0_0_0_1px_var(--sidebar-border)] hover:bg-sidebar-accent hover:text-sidebar-accent-foreground hover:shadow-[0_0_0_1px_var(--sidebar-accent)]"
	default:
		return "hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
	}
}

func sidebarMenuButtonSizeClass(size string) string {
	switch size {
	case "sm":
		return "h-7 text-xs"
	case "lg":
		return "h-12 text-sm group-data-[collapsible=icon]:p-0!"
	default:
		return "h-8 text-sm"
	}
}

// SidebarMenuButton drops asChild (see the package doc comment) — always a
// <button>, type="button" added per house convention (see
// SidebarGroupAction's own doc comment).
//
// Tooltip gating (source map §5): `tooltip` absent (the Go zero value, "")
// renders a plain button, no Tooltip wrapper at all. A non-empty tooltip
// wraps in ui.Tooltip and puts data-gsxui-tooltip-trigger directly on THIS
// button — not a nested ui.TooltipTrigger, which would be the same
// button-in-button HTML trap documented on DialogTrigger/SheetTrigger (see
// docs/jsx-parity.md `## dialog` FINDING); ui.Tooltip's own doc comment
// establishes this exact idiom. TooltipContent always anchors ABOVE the
// trigger (gsxui's Tooltip has no side/align params at all — see
// tooltip.gsx's own doc comment — unlike the reference's side="right"),
// and gating "only ever shown when collapsed AND on desktop" has no
// isMobile/state param to read here either, so it is reproduced in pure
// CSS instead: `hidden group-data-[collapsible=icon]:block` only matches
// while SOME ancestor `.group` carries data-collapsible="icon" — which,
// per the package doc comment's MOBILE section, is true for the DESKTOP
// copy of this button exactly when the sidebar is icon-collapsed, and
// never true for the MOBILE Sheet copy (whose ancestor tree carries no
// data-collapsible at all) — mobile is suppressed for free, no separate
// isMobile check needed.
component SidebarMenuButton(isActive bool, variant string, size string, tooltip string, children gsx.Node, attrs gsx.Attrs) {
	{ if tooltip == "" {
		<button
			type="button"
			data-slot="sidebar-menu-button"
			data-sidebar="menu-button"
			data-size={size |> default("default")}
			data-active={isActive}
			class={ sidebarMenuButtonBase, sidebarMenuButtonVariantClass(variant), sidebarMenuButtonSizeClass(size) }
			{ attrs... }
		>
			{ children }
		</button>
	} else {
		<Tooltip>
			<button
				type="button"
				data-slot="sidebar-menu-button"
				data-sidebar="menu-button"
				data-size={size |> default("default")}
				data-active={isActive}
				data-gsxui-tooltip-trigger
				class={ sidebarMenuButtonBase, sidebarMenuButtonVariantClass(variant), sidebarMenuButtonSizeClass(size) }
				{ attrs... }
			>
				{ children }
			</button>
			<TooltipContent class="hidden group-data-[collapsible=icon]:block">{ tooltip }</TooltipContent>
		</Tooltip>
	} }
}

// SidebarMenuAction drops asChild (see the package doc comment) — always a
// <button>, type="button" added per house convention. showOnHover (default
// false) adds the reference's md:opacity-0 + focus/hover-visible class
// block as one more entry in the composable class list.
component SidebarMenuAction(showOnHover bool, children gsx.Node, attrs gsx.Attrs) {
	<button
		type="button"
		data-slot="sidebar-menu-action"
		data-sidebar="menu-action"
		class={
			"absolute top-1.5 right-1 flex aspect-square w-5 items-center justify-center rounded-md p-0 text-sidebar-foreground ring-sidebar-ring outline-hidden transition-transform peer-hover/menu-button:text-sidebar-accent-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 [&>svg]:size-4 [&>svg]:shrink-0",
			"after:absolute after:-inset-2 md:after:hidden",
			"peer-data-[size=sm]/menu-button:top-1 peer-data-[size=default]/menu-button:top-1.5 peer-data-[size=lg]/menu-button:top-2.5",
			"group-data-[collapsible=icon]:hidden",
			"group-focus-within/menu-item:opacity-100 group-hover/menu-item:opacity-100 peer-data-[active=true]/menu-button:text-sidebar-accent-foreground data-[state=open]:opacity-100 md:opacity-0": showOnHover,
		}
		{ attrs... }
	>
		{ children }
	</button>
}

component SidebarMenuBadge(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="sidebar-menu-badge"
		data-sidebar="menu-badge"
		class="pointer-events-none absolute right-1 flex h-5 min-w-5 items-center justify-center rounded-md px-1 text-xs font-medium text-sidebar-foreground tabular-nums select-none peer-hover/menu-button:text-sidebar-accent-foreground peer-data-[active=true]/menu-button:text-sidebar-accent-foreground peer-data-[size=sm]/menu-button:top-1 peer-data-[size=default]/menu-button:top-1.5 peer-data-[size=lg]/menu-button:top-2.5 group-data-[collapsible=icon]:hidden"
		{ attrs... }
	>
		{ children }
	</div>
}

// SidebarMenuSkeleton's text width is a random 50-90% in the reference,
// computed once via React.useMemo so it never changes across re-renders of
// the SAME mounted component. gsx has no client-side "compute once at
// mount" equivalent for a server-rendered page (design decision flagged,
// not resolved, by the source map's own §2 entry) — this port computes it
// server-side, once per render call, which is fixed for the page's
// lifetime exactly like useMemo(() => …, []) is, and (a difference from
// upstream worth naming) varies independently between multiple
// SidebarMenuSkeleton instances on the same page, which upstream's
// per-instance useMemo also does. No children param — matches the
// reference, whose own JSX never renders {children} either.
component SidebarMenuSkeleton(showIcon bool, attrs gsx.Attrs) {
	{{
		width := strconv.Itoa(rand.Intn(40)+50) + "%"
	}}
	<div data-slot="sidebar-menu-skeleton" data-sidebar="menu-skeleton" class="flex h-8 items-center gap-2 rounded-md px-2" { attrs... }>
		{ if showIcon {
			<Skeleton class="size-4 rounded-md" data-sidebar="menu-skeleton-icon"/>
		} }
		<Skeleton
			class="h-4 max-w-(--skeleton-width) flex-1"
			data-sidebar="menu-skeleton-text"
			style=css`--skeleton-width:@{width}`
		/>
	</div>
}

component SidebarMenuSub(children gsx.Node, attrs gsx.Attrs) {
	<ul
		data-slot="sidebar-menu-sub"
		data-sidebar="menu-sub"
		class="mx-3.5 flex min-w-0 translate-x-px flex-col gap-1 border-l border-sidebar-border px-2.5 py-0.5 group-data-[collapsible=icon]:hidden"
		{ attrs... }
	>
		{ children }
	</ul>
}

component SidebarMenuSubItem(children gsx.Node, attrs gsx.Attrs) {
	<li data-slot="sidebar-menu-sub-item" data-sidebar="menu-sub-item" class="group/menu-sub-item relative" { attrs... }>
		{ children }
	</li>
}

// SidebarMenuSubButton drops asChild (see the package doc comment) —
// always an <a> (matching the reference's own Comp = asChild ? Slot.Root :
// "a"), so a caller supplies href via attrs like any other native anchor.
component SidebarMenuSubButton(size string, isActive bool, children gsx.Node, attrs gsx.Attrs) {
	<a
		data-slot="sidebar-menu-sub-button"
		data-sidebar="menu-sub-button"
		data-size={size |> default("md")}
		data-active={isActive}
		class={
			"flex h-7 min-w-0 -translate-x-px items-center gap-2 overflow-hidden rounded-md px-2 text-sidebar-foreground ring-sidebar-ring outline-hidden hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 [&>span:last-child]:truncate [&>svg]:size-4 [&>svg]:shrink-0 [&>svg]:text-sidebar-accent-foreground",
			"data-[active=true]:bg-sidebar-accent data-[active=true]:text-sidebar-accent-foreground",
			switch size {
			case "sm":
				"text-xs"
			default:
				"text-sm"
			},
			"group-data-[collapsible=icon]:hidden",
		}
		{ attrs... }
	>
		{ children }
	</a>
}

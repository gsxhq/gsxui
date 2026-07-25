package ui

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// NavigationMenu is the shadcn/ui NavigationMenu root — a hover mega-menu,
// NOT a menu of items (contrast dropdown/context-menu/menubar, all sharing
// role="menu"/role="menuitem" machinery this component needs none of).
// Renders as <nav> (Radix's own NavigationMenuPrimitive.Root anatomy,
// derived-not-read — no @radix-ui/react-navigation-menu dist exists in the
// reference checkout to trace, per docs/superpowers/plans/2026-07-25-tier4-
// source-map-menus.md ## navigation-menu; a <nav> landmark is the
// well-established public Radix/WAI-ARIA convention for this component).
// viewport: "" (default, true — the shared measuring panel described
// below) | "false" (each Content becomes its own self-contained floating
// panel instead, via the group-data-[viewport=false]/navigation-menu:
// selectors baked into NavigationMenuContent's own class).
//
// ARCHITECTURE (ADAPT, load-bearing — no portal, and no DOM-node-moving
// either, a genuinely new shape among this codebase's menu family): Radix's
// own NavigationMenuContent, when a shared Viewport exists, portals into it
// so only the active panel's real DOM lives inside the one bordered/
// shadowed box, which is how Viewport gets to measure and resize itself
// off whatever child happens to be mounted — the CSS custom properties
// that drive that sizing (--radix-navigation-menu-viewport-width/height)
// are set by Radix's own Viewport runtime, absent from the reference
// checkout (derived-not-read, source map ## navigation-menu §2). gsx has
// no portal. This port does not simulate one by reparenting DOM nodes at
// runtime either — instead every NavigationMenuContent renders DOM-nested
// inside its own NavigationMenuItem (never moved), matching the
// DOM-nesting-not-portalling ADAPT every sibling submenu in this codebase
// already uses, and is independently popover="manual", escaping overflow/
// clipping via the native top layer exactly like every other popover here.
// NavigationMenuViewport, when data-viewport!="false", is a SEPARATE,
// second popover="manual" panel shown/positioned/sized in lockstep with
// whichever Content is currently active — coincident, not containing: both
// open at the identical fixed left/top rect (computed off the List's own
// boundingClientRect, so every item's panel appears at the same on-screen
// spot — the actual "shared viewport" visual contract), Viewport supplies
// the chrome (border/bg/shadow) and is resized to the active Content's own
// measured width/height (ui/navigation-menu.js), Content itself stays
// unchromed in this mode (its own group-data-[viewport=false] chrome block
// simply never matches, same as upstream). Two coincident top-layer
// popovers, not one containing the other — chosen over a DOM-move because
// every show/hide stays a single, independent, idempotent call, the same
// shape as every other popover pair in this codebase, with no
// node-reparenting bookkeeping to get wrong.
//
// SIZING (source map ## navigation-menu §3, decisive): discrete,
// event-driven state is sufficient — no CSS in any of the eight
// independently-authored nova-family stylesheets transitions width/height
// on the viewport, only a scale transform gated on open/close. Viewport is
// therefore sized ONCE on activation (a plain getBoundingClientRect read
// off the newly-active Content) plus a ResizeObserver on that same Content
// for late-arriving resizes (an image finishing loading, dynamically
// populated content) — explicitly NOT a requestAnimationFrame tweening
// loop between panel sizes, which the reference does not have.
//
// Hover and focus both open; popover="manual" gives the top layer without
// light dismiss/Esc — modeled on ui/hover-card.js's own deliberate choice
// (hover/focus alone drive open/close), not dropdown.js's own click+toggle-
// event model, since this is hover-driven the same way a hover card is.
// data-state drives the discrete-transition open/close, reused
// byte-identically from the popover family (ui/dropdown.gsx's own block).
// Requires the navigation-menu behavior module (ui/navigation-menu.js).
component NavigationMenu(viewport string, children gsx.Node, attrs gsx.Attrs) {
	<nav
		data-slot="navigation-menu"
		data-gsxui-navigation-menu
		{ if viewport == "false" {
			data-viewport="false"
		} else {
			data-viewport="true"
		} }
		class="group/navigation-menu relative flex max-w-max flex-1 items-center justify-center"
		{ attrs... }
	>
		{ children }
		{ if viewport != "false" {
			<NavigationMenuViewport/>
		} }
	</nav>
}

// NavigationMenuList is the <ul> holding NavigationMenuItems (and,
// optionally, a trailing NavigationMenuIndicator — ui/navigation-menu.js
// positions it relative to this element).
component NavigationMenuList(children gsx.Node, attrs gsx.Attrs) {
	<ul data-slot="navigation-menu-list" data-gsxui-navigation-menu-list class="group flex flex-1 list-none items-center justify-center gap-1" { attrs... }>
		{ children }
	</ul>
}

// NavigationMenuItem is the <li> pairing one NavigationMenuTrigger with its
// own NavigationMenuContent — data-gsxui-navigation-menu-item is the
// proximity anchor ui/navigation-menu.js uses to resolve "this trigger's
// own content" (closest("[data-gsxui-navigation-menu-item]")), the same
// role DropdownMenu's own root plays for its single trigger/content pair.
component NavigationMenuItem(children gsx.Node, attrs gsx.Attrs) {
	<li data-slot="navigation-menu-item" data-gsxui-navigation-menu-item class="relative" { attrs... }>
		{ children }
	</li>
}

// NavigationMenuTriggerStyle is the shadcn/ui navigationMenuTriggerStyle
// cva, exported standalone (source map ## navigation-menu §1) so a plain
// <a>/NavigationMenuLink styled as a top-level item that isn't a real
// Trigger — no dropdown attached — can reuse the identical visual, the same
// way shadcn's own demo composes it. ADAPT (nova): rounded-md -> rounded-lg
// (nova metric). bg-background is dropped entirely — nova's own
// .cn-navigation-menu-trigger carries no background-color utility at all, a
// real, not just metric, style choice per the source map's own explicit
// finding ("port nova's version as written... don't silently reintroduce
// bg-background"). Nova's FURTHER hover-model rewrite (hover:bg-muted,
// focus:bg-muted, data-open: shorthand replacing data-[state=open]:) is NOT
// adopted — a real interaction-model/color change, out of scope the same
// way ## menubar's own MenubarTrigger ADAPT already ruled for the identical
// situation; the standing house exception keeps data-[state=…]: everywhere
// in this codebase, and colors are untouched by a metric-tokens-only nova
// retarget except where a source-map finding explicitly says otherwise (as
// it does, narrowly, for bg-background alone).
func NavigationMenuTriggerStyle() string {
	return "group inline-flex h-9 w-max items-center justify-center rounded-lg px-4 py-2 text-sm font-medium transition-[color,box-shadow] outline-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1 disabled:pointer-events-none disabled:opacity-50 data-[state=open]:bg-accent/50 data-[state=open]:text-accent-foreground data-[state=open]:hover:bg-accent data-[state=open]:focus:bg-accent"
}

// NavigationMenuTrigger opens/closes its sibling NavigationMenuContent
// (ui/navigation-menu.js: pointerover/focusin open, pointerout/focusout
// schedule a hover-card-shaped grace-period close, click toggles). class is
// cn(NavigationMenuTriggerStyle(), "group", className) verbatim — the base
// string already starts with "group"; shadcn's own cn() call adds a second,
// redundant "group" token as a separate argument (source map ## navigation-
// menu §1: "verbatim, not a transcription error"). The trailing chevron
// keeps new-york-v4's own selector spelling (group-data-[state=open]:, the
// standing house exception) rather than nova's own group-data-open:/
// group-data-popup-open: shorthand pair — the second of which belongs to
// Base UI's differently-shaped primitive per the source map's own
// provenance note and has no Radix data-state equivalent to key off here.
component NavigationMenuTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		data-slot="navigation-menu-trigger"
		data-gsxui-navigation-menu-trigger
		type="button"
		aria-expanded="false"
		data-state="closed"
		class={ NavigationMenuTriggerStyle(), "group" }
		{ attrs... }
	>
		{ children }
		<icon.ChevronDown class="relative top-[1px] ml-1 size-3 transition duration-300 group-data-[state=open]:rotate-180"/>
	</button>
}

// NavigationMenuContent is the panel a NavigationMenuTrigger opens — see
// the file header's own ARCHITECTURE paragraph for why it is independently
// popover="manual", DOM-nested inside its own NavigationMenuItem rather
// than portalled or moved. data-side="bottom" is server-rendered
// statically — ui/navigation-menu.js always opens below the list, the same
// hand-rolled-fixed-position stopgap as every sibling popover in this
// codebase (see docs/jsx-parity.md ## dropdown NOTE). The six-token
// data-[motion=…]/animate-in/out/fade-in/out slide mechanism (new-york-v4's
// own direction-aware "which side did this panel enter from" animation,
// driven by Radix's own runtime index tracking) is NOT reproduced — GAP,
// ledgered per the task's own binding decision to reuse the popover
// family's discrete-transition block byte-identically instead, the same
// six-tokens-replaced-by-one-mechanism RULING ## menubar's own
// MenubarContent already made for its own (differently-shaped) omission.
// The group-data-[viewport=false]/navigation-menu: chrome block (border/
// bg/shadow/rounded) is ONLY live when the ancestor NavigationMenu's own
// data-viewport="false" — in the default (true) mode Content stays
// unchromed, its own coincident Viewport supplies the box instead (file
// header). border is kept, not swapped for nova's own ring-1 (standing
// house exception); rounded-md -> rounded-lg is nova's own metric.
component NavigationMenuContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="navigation-menu-content"
		data-gsxui-navigation-menu-content
		popover="manual"
		data-state="closed"
		data-side="bottom"
		class={
			"top-0 left-0 w-full p-2 pr-2.5 md:absolute md:w-auto",
			"group-data-[viewport=false]/navigation-menu:top-full group-data-[viewport=false]/navigation-menu:mt-1.5 group-data-[viewport=false]/navigation-menu:overflow-hidden group-data-[viewport=false]/navigation-menu:rounded-lg group-data-[viewport=false]/navigation-menu:border group-data-[viewport=false]/navigation-menu:bg-popover group-data-[viewport=false]/navigation-menu:text-popover-foreground group-data-[viewport=false]/navigation-menu:shadow **:data-[slot=navigation-menu-link]:focus:ring-0 **:data-[slot=navigation-menu-link]:focus:outline-none",
			// Discrete-transition enter/exit, reused byte-identically from the
			// popover family (ui/dropdown.gsx's own block) per the task's own
			// binding decision — see the doc comment above for the six-token
			// data-motion slide this replaces.
			"opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:opacity-0 starting:open:scale-95",
			"data-[side=bottom]:starting:open:-translate-y-2 data-[side=left]:starting:open:translate-x-2 data-[side=right]:starting:open:-translate-x-2 data-[side=top]:starting:open:translate-y-2"
		}
		{ attrs... }
	>
		{ children }
	</div>
}

// NavigationMenuViewport is the shared measuring panel — see the file
// header's own ARCHITECTURE paragraph for the coincident-popover mechanism
// this port uses in place of Radix's Content-portals-into-Viewport model.
// Two nested elements, matching shadcn's own structure verbatim (source map
// ## navigation-menu §1): an outer plain <div> (no data-slot) and the inner
// data-slot="navigation-menu-viewport" div, which is the one
// ui/navigation-menu.js actually positions/sizes/shows — the outer
// wrapper's own absolute/top-full/left-0 CSS-relative-to-<nav> positioning
// is superseded by the inner's own JS-driven position:fixed (dead weight,
// kept per this codebase's own "some selectors are inert until real
// placement logic lands" precedent, docs/jsx-parity.md ## dropdown NOTE),
// but isolate/z-50/flex/justify-center still do real work centering the
// inner viewport horizontally.
//
// The h-[var(--radix-navigation-menu-viewport-height)]/
// md:w-[var(--radix-navigation-menu-viewport-width)] CSS-custom-property
// sizing is replaced with direct inline width/height set by
// ui/navigation-menu.js (the source map's own explicitly-offered
// alternative to a custom property, ## navigation-menu §3) — no Radix
// runtime exists to set those two variables in this checkout
// (derived-not-read, ## navigation-menu §2), the same
// no-runtime-var-to-read substitution as every sibling popover's own
// origin-* ADAPT. rounded-md -> rounded-lg is nova's own metric; border is
// kept (standing house exception, not nova's own ring-1/ring-foreground).
// Discrete-transition enter/exit is reused byte-identically from the
// popover family, same as NavigationMenuContent's own.
component NavigationMenuViewport(attrs gsx.Attrs) {
	<div class="absolute top-full left-0 isolate z-50 flex justify-center">
		<div
			data-slot="navigation-menu-viewport"
			data-gsxui-navigation-menu-viewport
			popover="manual"
			data-state="closed"
			data-side="bottom"
			class={
				"origin-top relative mt-1.5 w-full overflow-hidden rounded-lg border bg-popover text-popover-foreground shadow md:w-auto",
				"opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:opacity-0 starting:open:scale-95",
				"data-[side=bottom]:starting:open:-translate-y-2 data-[side=left]:starting:open:translate-x-2 data-[side=right]:starting:open:-translate-x-2 data-[side=top]:starting:open:translate-y-2"
			}
			{ attrs... }
		></div>
	</div>
}

// NavigationMenuLink is a single item inside a NavigationMenuContent (or,
// styled with NavigationMenuTriggerStyle() instead, a plain top-level nav
// link with no dropdown at all). active mirrors Radix's own data-active
// convention — the caller computes whether this link's own destination
// matches the current page (gsx has no router to do it implicitly, the same
// no-context shape as ## toggle-group's group→item params). Clicking a
// Link inside an open Content closes that panel (ui/navigation-menu.js); a
// Link with no Content ancestor (the plain top-level case) has nothing to
// close and the click just navigates. rounded-sm -> rounded-lg is nova's
// own metric, WITH nova's own in-data-[slot=navigation-menu-content]:
// rounded-md override preserved (a real, if narrow, per-context nova value,
// same in-data-[slot=…] idiom already used by button.gsx's own
// in-data-[slot=button-group]: rounded bump) — new-york-v4's own flex-col
// gap-1 structure (a mega-menu tile: title stacked over description) is
// kept over nova's own flex items-center gap-2 rewrite, which reads as a
// plain single-line list-item shape for a differently-composed nav-menu
// variant, not a metric bump. hover:/focus:/data-[active=true]: colors stay
// on accent, not nova's own muted rewrite — same out-of-scope ruling as
// NavigationMenuTriggerStyle's own doc comment.
component NavigationMenuLink(active bool, children gsx.Node, attrs gsx.Attrs) {
	<a
		data-slot="navigation-menu-link"
		data-gsxui-navigation-menu-link
		{ if active {
			data-active="true"
		} else {
			data-active="false"
		} }
		class="flex flex-col gap-1 rounded-lg p-2 text-sm transition-all outline-none in-data-[slot=navigation-menu-content]:rounded-md hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1 data-[active=true]:bg-accent/50 data-[active=true]:text-accent-foreground data-[active=true]:hover:bg-accent data-[active=true]:focus:bg-accent [&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground"
		{ attrs... }
	>
		{ children }
	</a>
}

// NavigationMenuIndicator is the small rotated-square pointer tracking the
// active trigger, positioned under NavigationMenuList's own last child (the
// caller places it there, matching Radix's own composition) by
// ui/navigation-menu.js — translateX off the trigger's own rect, relative
// to the nearest positioned ancestor (this port injects position:relative
// onto the List at the same time, since neither shadcn's own class nor
// nova's carries a position token for either element — Radix's own runtime
// sets that inline, the same "computed position via JS, not baked into the
// class" idiom every fixed-positioned popover in this codebase already
// uses). ADAPT: the data-[state=hidden|visible]:animate-out/in fade-out/in
// pair (tw-animate-css, GAP category — see ## animations) is replaced with
// a plain CSS opacity transition instead of the popover family's discrete-
// transition block: Indicator is not popover-attached (no :popover-open
// state to key open: off), so the byte-identical block reused by
// Content/Viewport does not apply here.
component NavigationMenuIndicator(attrs gsx.Attrs) {
	<div
		data-slot="navigation-menu-indicator"
		data-gsxui-navigation-menu-indicator
		data-state="hidden"
		class="top-full z-[1] flex h-1.5 items-end justify-center overflow-hidden opacity-0 transition-opacity duration-200 data-[state=visible]:opacity-100"
		{ attrs... }
	>
		<div class="relative top-[60%] h-2 w-2 rotate-45 rounded-tl-sm bg-border shadow-md"></div>
	</div>
}

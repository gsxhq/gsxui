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
//
// GAP (shared-viewport mode not shipped — FIX ROUND 1, reverted from a
// broken v1 attempt): shadcn's own NavigationMenu takes a `viewport` prop,
// defaulting true, which portals every NavigationMenuContent into one
// shared NavigationMenuViewport panel that morphs between differently-sized
// menus. This port ships ONLY the `viewport={false}` configuration — the
// OTHER first-class mode new-york-v4's own navigation-menu.tsx supports,
// with its own real class strings, not an invented fallback. Each
// NavigationMenuContent below is unconditionally its own fully-chromed
// panel (border/bg/shadow/radius baked in, no group-data-[viewport=false]
// gate — there is no second mode to gate against). Two reasons, not one:
// (1) gsx has no portal, and reproducing Radix's Content-portals-into-
// Viewport model without one means moving a live popover between DOM
// parents at runtime — a v1 of this component tried the alternative (two
// COINCIDENT popovers, Content and a separate Viewport shown/positioned at
// the same rect instead of one containing the other) and it does not work:
// once both are promoted to the top layer, whichever is shown second
// paints an OPAQUE surface directly over the other — the viewport becomes
// a lid over the content, not a frame around it, and the mega-menu rendered
// as an empty box with the links completely unreachable (confirmed live,
// `elementFromPoint` returned the viewport at every sampled point inside
// the panel). (2) The shared viewport's only real visual advantage is that
// morph between differently-sized panels, and the source map's own decisive
// finding (`## navigation-menu` §3) already established that none of the
// eight nova-family stylesheets transitions width or height on the
// viewport — only a scale transform on open/close. The shared viewport
// buys essentially nothing this port doesn't already get from each
// Content's own independent open/close scale transition. Adopting the
// shared-viewport mode later would need genuine runtime relocation of the
// Content node into a Viewport node (an appendChild-based move, which does
// not lose the node's own listeners/state) — deliberately not attempted
// here. `NavigationMenuViewport` itself is DROPPED, not kept as a no-op:
// the viewport={false} path this port actually ships never renders it at
// all (`{viewport && <NavigationMenuViewport />}` in shadcn's own source),
// so there is nothing left for it to do — shipping it anyway would be
// exactly the "hook nothing binds to" mistake already caught and reverted
// once earlier in this batch.
component NavigationMenu(children gsx.Node, attrs gsx.Attrs) {
	<nav
		data-slot="navigation-menu"
		data-gsxui-navigation-menu
		class="group/navigation-menu relative flex max-w-max flex-1 items-center justify-center"
		{ attrs... }
	>
		{ children }
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
// menu §1: "verbatim, not a transcription error") — this port's own class
// merge does not dedupe it, so it renders twice, matching the literal
// source. The trailing chevron keeps new-york-v4's own selector spelling
// (group-data-[state=open]:, the standing house exception) rather than
// nova's own group-data-open:/group-data-popup-open: shorthand pair — the
// second of which belongs to Base UI's differently-shaped primitive per the
// source map's own provenance note and has no Radix data-state equivalent
// to key off here.
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
// the file header's own GAP paragraph for why this is the shadcn
// `viewport={false}` configuration, unconditionally: independently
// popover="manual", DOM-nested inside its own NavigationMenuItem (not
// portalled, not moved — the DOM-nesting-not-portalling ADAPT every
// sibling submenu in this codebase already uses), positioned by
// ui/navigation-menu.js under ITS OWN trigger's rect (top-full/mt-1.5,
// matching new-york-v4's own viewport=false positioning, which anchors to
// the trigger's own NavigationMenuItem — class="relative" — not a shared
// panel location). data-side="bottom" is server-rendered statically, the
// same hand-rolled-fixed-position stopgap as every sibling popover in this
// codebase (see docs/jsx-parity.md ## dropdown NOTE).
//
// The six-token data-[motion=…]/animate-in/out/fade-in/out slide mechanism
// (new-york-v4's own direction-aware "which side did this panel enter
// from" animation, driven by Radix's own runtime index tracking) is NOT
// reproduced — GAP, per the task's own binding decision to reuse the
// popover family's discrete-transition block byte-identically instead, the
// same six-tokens-replaced-by-one-mechanism RULING ## menubar's own
// MenubarContent already made for its own (differently-shaped) omission.
//
// The border/bg/shadow/rounded chrome — new-york-v4's own
// group-data-[viewport=false]/navigation-menu: block — is applied
// unconditionally, its group-data-[viewport=false]/navigation-menu: prefix
// stripped from every token: there is only one mode in this port, so the
// gate would never toggle and the selector would be permanently dead CSS,
// the same "drop the selector, don't ship dead CSS" call as dropdown's own
// inset ADAPT. border is kept, not swapped for nova's own ring-1 (standing
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
			"top-full mt-1.5 overflow-hidden rounded-lg border bg-popover text-popover-foreground shadow **:data-[slot=navigation-menu-link]:focus:ring-0 **:data-[slot=navigation-menu-link]:focus:outline-none",
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
// uses). Unaffected by this file's own viewport={false} GAP — the
// indicator tracks the active TRIGGER, not the (now-removed) shared
// viewport panel. ADAPT: the data-[state=hidden|visible]:animate-out/in
// fade-out/in pair (tw-animate-css, GAP category — see ## animations) is
// replaced with a plain CSS opacity transition instead of the popover
// family's discrete-transition block: Indicator is not popover-attached (no
// :popover-open state to key open: off), so the byte-identical block
// NavigationMenuContent reuses does not apply here.
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

package ui

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// DropdownMenu is the shadcn/ui DropdownMenu on the native popover API: the
// top layer replaces Radix's Portal, light dismiss and Esc are browser-
// native. Trigger and content are wired by proximity — DropdownMenuTrigger
// opens the popover inside the same DropdownMenu root, no ids. JS adds
// fixed-position anchoring to the trigger rect (CSS anchor positioning is
// not yet Baseline — see docs/jsx-parity.md), state/aria sync, arrow-key
// roving focus, and close-on-select. Requires the dropdown behavior module
// (ui/dropdown/dropdown.js).
//
// SUBMENUS — POPOVER NESTING, NOT PORTALLING (ADAPT, load-bearing): Radix
// portals DropdownMenuSubContent to document.body, same as top-level
// Content. That does not work on the native popover API. Measured in Chrome
// (docs/superpowers/plans/2026-07-25-tier4-source-map-menus.md ## shared-
// items §5, "spec §1" below):
//
//	child popover DOM-nested inside parent, opened with showPopover()      -> parent STAYS OPEN
//	child not nested, opened via a real popovertarget invoker click        -> parent STAYS OPEN
//	child not nested, opened programmatically with showPopover()           -> parent LIGHT-DISMISSES
//
// A popover="auto" establishes its ancestor chain by DOM containment OR by
// invoker relationship. dropdown.js opens submenus on pointerenter and on
// ArrowRight — not only on click — so every open is programmatic; DOM
// nesting is therefore the only robust option, not a style preference.
// DropdownMenuSubContent renders nested inside its DropdownMenuSub, which
// itself sits inside the parent DropdownMenuContent (pinned by
// TestDropdownMenuSubNestsContentInsideParentContent). The parent content's
// own overflow-x-hidden/overflow-y-auto (TestDropdownContentPinned) does NOT
// clip the nested submenu: once shown, a popover is promoted to the
// browser's top layer for painting/stacking purposes, exempting it from any
// DOM ancestor's overflow/clip/transform — the same mechanism that already
// lets DropdownMenuContent itself skip Radix's Portal. No overflow fix was
// needed on DropdownMenuContent's pinned class.
//
// SERVER-RENDERED CHECKED STATE (MECHANISM): DropdownMenuCheckboxItem/
// DropdownMenuRadioItem take a `checked bool` that stamps both aria-checked
// and data-state="checked"|"unchecked" on first paint — the menu is correct
// with JS disabled. The check/circle indicator icon is always rendered (no
// server-side conditional mount, unlike Radix's own ItemIndicator, which
// mounts/unmounts by React state) and its visibility is purely CSS, gated by
// an ancestor-selector arbitrary variant, `[[data-state=checked]_&]:flex`,
// keyed off the item's own data-state — the same `[[…]_&]` idiom already
// used by sidebar.gsx/kbd.gsx/collapsible.gsx, chosen over a named
// `group/…:` marker specifically so the item's own pinned base class string
// needs no extra token prepended to it (verbatim per the source map).
// dropdown.js's click handlers only ever flip aria-checked/data-state; they
// never touch the indicator's own markup, so no icon markup needs to be
// duplicated into dropdown.js (which would be an invisible-to-registry.Deps
// dependency — forbidden).
component DropdownMenu(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="dropdown-menu" data-gsxui-dropdown class="contents" { attrs... }>{ children }</div>
}

component DropdownMenuTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		data-slot="dropdown-menu-trigger"
		data-gsxui-dropdown-trigger
		type="button"
		aria-haspopup="menu"
		aria-expanded="false"
		{ attrs... }
	>
		{ children }
	</button>
}

// DropdownMenuContent renders the popover. popover="auto" gives top layer,
// light dismiss, and free Esc; data-state is server-rendered "closed" and
// kept in sync by dropdown.js on the toggle event. data-side="bottom" is
// server-rendered statically — dropdown.js always anchors below the
// trigger, so shadcn's data-[side=bottom]:slide-in-from-top-2 enter slide
// applies without Radix's runtime side tracking (same ADAPT as tooltip).
component DropdownMenuContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="dropdown-menu-content"
		data-gsxui-dropdown-content
		popover="auto"
		role="menu"
		tabindex="-1"
		data-state="closed"
		data-side="bottom"
		class={
			"z-50 max-h-96 min-w-[8rem] origin-top-left overflow-x-hidden overflow-y-auto rounded-lg border bg-popover p-1 text-popover-foreground shadow-md",
			// Discrete-transition enter/exit replacing the tw-animate keyframe
			// pair — a popover's exit keyframe never gets to play (hide is
			// instant display:none); see popover.gsx's ADAPT comment and
			// docs/jsx-parity.md ## animations for the full mechanism.
			"opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:opacity-0 starting:open:scale-95",
			"data-[side=bottom]:starting:open:-translate-y-2 data-[side=left]:starting:open:translate-x-2 data-[side=right]:starting:open:-translate-x-2 data-[side=top]:starting:open:translate-y-2"
		}
		{ attrs... }
	>
		{ children }
	</div>
}

// DropdownMenuItem is the shadcn/ui DropdownMenuItem, ported as a real menu
// item on a <div role="menuitem">: dropdown.js's arrow-key roving focus
// walks these. variant: "" (default) | "destructive". inset is dropped
// (see docs/jsx-parity.md) — the data-[inset]:pl-8 selector is removed with
// it.
component DropdownMenuItem(variant string, children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="dropdown-menu-item"
		data-gsxui-dropdown-item
		data-variant={variant |> default("default")}
		role="menuitem"
		tabindex="-1"
		class="relative flex cursor-default items-center gap-1.5 rounded-md px-1.5 py-1 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 data-[variant=destructive]:text-destructive data-[variant=destructive]:focus:bg-destructive/10 data-[variant=destructive]:focus:text-destructive dark:data-[variant=destructive]:focus:bg-destructive/20 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground data-[variant=destructive]:*:[svg]:text-destructive!"
		{ attrs... }
	>
		{ children }
	</div>
}

// DropdownMenuGroup wraps a set of items for a11y grouping. shadcn's own
// Group carries no class string at all (source map ## shared-items §1) —
// role="group" is added here, not in the .tsx: Radix's Group primitive
// stamps it at runtime (derived-not-read, well-known WAI-ARIA menu
// authoring practice — a set of role="menuitemradio" items in particular
// must be wrapped by a role="group" container).
component DropdownMenuGroup(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="dropdown-menu-group" data-gsxui-menu-group role="group" { attrs... }>{ children }</div>
}

// DropdownMenuCheckboxItem is the shadcn/ui DropdownMenuCheckboxItem.
// checked is server-rendered (see the file header MECHANISM); value is the
// item's own identity, stamped as data-value and echoed on dropdown.js's
// gsxui:change event. Selecting a checkbox item does NOT close the menu
// (Radix's own contract — confirmed against the source map's derived ARIA
// anatomy, ## shared-items §3 — a checkbox toggles in place).
component DropdownMenuCheckboxItem(checked bool, value string, children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="dropdown-menu-checkbox-item"
		data-gsxui-menu-checkbox-item
		role="menuitemcheckbox"
		data-value={value}
		{ if checked {
			aria-checked="true"
			data-state="checked"
		} else {
			aria-checked="false"
			data-state="unchecked"
		} }
		tabindex="-1"
		class="relative flex cursor-default items-center gap-2 rounded-sm py-1.5 pr-2 pl-8 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4"
		{ attrs... }
	>
		<span class="pointer-events-none absolute left-2 hidden size-3.5 items-center justify-center [[data-state=checked]_&]:flex">
			<icon.Check class="size-4"/>
		</span>
		{ children }
	</div>
}

// DropdownMenuRadioGroup wraps a set of DropdownMenuRadioItems. value is the
// server-rendered current value, stamped as data-value on the root — the
// same "checked state is server-rendered" contract as CheckboxItem, kept in
// sync by dropdown.js on selection and echoed on the group's own
// gsxui:change event. data-gsxui-menu-radio-group is the proximity anchor
// dropdown.js uses to scope "clear every OTHER item in this group" to this
// group alone, not every radio item on the page.
component DropdownMenuRadioGroup(value string, children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="dropdown-menu-radio-group" data-gsxui-menu-radio-group role="group" data-value={value} { attrs... }>
		{ children }
	</div>
}

// DropdownMenuRadioItem is the shadcn/ui DropdownMenuRadioItem — same
// shape/class as CheckboxItem (source map ## shared-items §1: byte-identical
// base style, own finding), swapping the check indicator for a filled dot.
// checked reflects whether THIS item's value equals its DropdownMenuRadioGroup's
// current value — the caller computes that comparison (gsx has no context to
// do it implicitly, same no-context shape as ## toggle-group's group→item
// params). Selecting a radio item DOES close the menu, same as a plain Item
// (only CheckboxItem stays open).
component DropdownMenuRadioItem(checked bool, value string, children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="dropdown-menu-radio-item"
		data-gsxui-menu-radio-item
		role="menuitemradio"
		data-value={value}
		{ if checked {
			aria-checked="true"
			data-state="checked"
		} else {
			aria-checked="false"
			data-state="unchecked"
		} }
		tabindex="-1"
		class="relative flex cursor-default items-center gap-2 rounded-sm py-1.5 pr-2 pl-8 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4"
		{ attrs... }
	>
		<span class="pointer-events-none absolute left-2 hidden size-3.5 items-center justify-center [[data-state=checked]_&]:flex">
			<icon.Circle class="size-2 fill-current"/>
		</span>
		{ children }
	</div>
}

// DropdownMenuLabel's inset prop is dropped along with DropdownMenuItem's
// (see docs/jsx-parity.md) — the data-[inset]:pl-8 selector is removed.
component DropdownMenuLabel(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="dropdown-menu-label" class="px-1.5 py-1 text-xs font-medium" { attrs... }>{ children }</div>
}

component DropdownMenuSeparator(attrs gsx.Attrs) {
	<div data-slot="dropdown-menu-separator" role="separator" class="-mx-1 my-1 h-px bg-border" { attrs... }></div>
}

component DropdownMenuShortcut(children gsx.Node, attrs gsx.Attrs) {
	<span data-slot="dropdown-menu-shortcut" class="ml-auto text-xs tracking-widest text-muted-foreground" { attrs... }>
		{ children }
	</span>
}

// DropdownMenuSub is the non-rendering submenu root — layout-neutral
// (class="contents", same idiom as DropdownMenu's own root) so its
// SubTrigger/SubContent children sit inline in the parent content's normal
// item flow. data-gsxui-menu-sub is the proximity anchor dropdown.js uses to
// pair a SubTrigger with its own SubContent (closest("[data-gsxui-menu-sub]")
// — same shape as DropdownMenu's own data-gsxui-dropdown root) and to scope
// the pointer-leave grace-period boundary check to "the whole sub."
component DropdownMenuSub(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="dropdown-menu-sub" data-gsxui-menu-sub class="contents" { attrs... }>{ children }</div>
}

// DropdownMenuSubTrigger opens/closes its sibling DropdownMenuSubContent
// (dropdown.js: pointerenter, ArrowRight, click). aria-haspopup/aria-expanded
// are server-rendered closed (derived-not-read ARIA anatomy, source map
// ## shared-items §3); dropdown.js keeps aria-expanded and data-state in step
// on every open/close, stamping data-state="open" BEFORE showPopover() (the
// same flash-avoidance rule DropdownMenuContent's own trigger click handler
// documents). ADAPT: the data-[inset]:pl-8 token is dropped from the pinned
// class, same call as DropdownMenuItem/DropdownMenuLabel's own ADAPT —
// gsxui's dropdown scope has no inset prop anywhere in this file, so the
// selector would be permanently dead CSS.
component DropdownMenuSubTrigger(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="dropdown-menu-sub-trigger"
		data-gsxui-menu-sub-trigger
		role="menuitem"
		aria-haspopup="menu"
		aria-expanded="false"
		data-state="closed"
		tabindex="-1"
		class="flex cursor-default items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[state=open]:bg-accent data-[state=open]:text-accent-foreground [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground"
		{ attrs... }
	>
		{ children }
		<icon.ChevronRight class="ml-auto size-4"/>
	</div>
}

// DropdownMenuSubContent is the submenu popover — see the file header's
// SUBMENUS comment for why it must render DOM-nested (not portalled) inside
// its DropdownMenuSub. ADAPT (mirrors DropdownMenuContent's own, verbatim
// class deviations from the source map, each already established elsewhere
// in this file):
//   - origin-(--radix-dropdown-menu-content-transform-origin) -> origin-top-left:
//     no Radix runtime CSS var to read; dropdown.js always anchors the
//     submenu to the right of its trigger (data-side="right", server-
//     rendered statically — same hand-rolled-fixed-position stopgap as the
//     top-level content's "always below," see docs/jsx-parity.md ## dropdown
//     NOTE; not Floating-UI collision-aware placement).
//   - the six data-[state=open|closed]:animate-in/out/fade/zoom tokens ->
//     the same discrete-transition block DropdownMenuContent uses (a
//     popover's exit keyframe never gets to play — see popover.gsx's ADAPT
//     and docs/jsx-parity.md ## animations).
//   - min-w-[8rem] -> min-w-[96px], rounded-md -> rounded-lg: nova's own
//     .cn-dropdown-menu-sub-content metrics win over new-york-v4's per the
//     visual-reference rule (nova also drops the min-w value new-york-v4
//     shares with the top-level Content, making the submenu narrower).
//     border is KEPT, not swapped for nova's ring-1 (standing house
//     exception — border->ring-1 swaps are not adopted anywhere in this
//     codebase).
component DropdownMenuSubContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="dropdown-menu-sub-content"
		data-gsxui-menu-sub-content
		popover="auto"
		role="menu"
		tabindex="-1"
		data-state="closed"
		data-side="right"
		class={
			"z-50 min-w-[96px] origin-top-left overflow-hidden rounded-lg border bg-popover p-1 text-popover-foreground shadow-lg",
			"opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:opacity-0 starting:open:scale-95",
			"data-[side=bottom]:starting:open:-translate-y-2 data-[side=left]:starting:open:translate-x-2 data-[side=right]:starting:open:-translate-x-2 data-[side=top]:starting:open:translate-y-2"
		}
		{ attrs... }
	>
		{ children }
	</div>
}

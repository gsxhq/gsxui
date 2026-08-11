package mira

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
// used by sidebar.gsx/kbd.gsx/collapsible.gsx, chosen over a named ancestor
// marker specifically so the item's own pinned base class string
// needs no extra token prepended to it (verbatim per the source map).
// dropdown.js's click handlers only ever flip aria-checked/data-state; they
// never touch the indicator's own markup, so no icon markup needs to be
// duplicated into dropdown.js (which would be an invisible-to-registry.Deps
// dependency — forbidden).
component DropdownMenu(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "contents" } { attrs... } data-gsxui-slot-dropdown-menu>{ children }</div>
}

component DropdownMenuTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		type="button"
		aria-haspopup="menu"
		aria-expanded="false"
		{ attrs... }
		data-gsxui-slot-dropdown-menu-trigger
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
		class={
			"transition-none data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 ring-foreground/10 bg-popover text-popover-foreground min-w-32 rounded-lg p-1 shadow-md ring-1 duration-100"
		}
		popover="auto"
		role="menu"
		tabindex="-1"
		data-state="closed"
		data-side="bottom"
		{ attrs... }
		data-gsxui-slot-dropdown-menu-content
	>
		{ children }
	</div>
}

// DropdownMenuItem is the shadcn/ui DropdownMenuItem, ported as a real menu
// item on a <div role="menuitem">: dropdown.js's arrow-key roving focus
// walks these. variant: "" (default) | "destructive". Callers reflect the
// CSS-only inset axis with data-inset through attrs.
component DropdownMenuItem(variant string, children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"group/dropdown-menu-item",
			"focus:bg-accent focus:text-accent-foreground dark:data-[variant=destructive]:focus:bg-destructive/20 not-data-[variant=destructive]:focus:**:text-accent-foreground min-h-7 gap-2 rounded-md px-2 py-1 text-xs/relaxed [&_svg:not([class*='size-'])]:size-3.5 flex",
			switch variant {
			case "destructive":
				"text-destructive focus:bg-destructive/10 focus:text-destructive *:[svg]:text-destructive"
			default:
				"text-foreground"
			}
		}
		data-variant={variant |> default("default")}
		role="menuitem"
		tabindex="-1"
		{ attrs... }
		data-gsxui-slot-dropdown-menu-item
	>
		{ children }
	</div>
}

// DropdownMenuGroup wraps a set of items for a11y grouping. shadcn's own
// Group carries no class string at all (source map ## shared-items §1) —
// role="group" is added here, not in the .tsx: Radix's Group primitive
// stamps it at runtime (derived-not-read, well-known WAI-ARIA menu
// authoring practice — a set of role="menuitemradio" items in particular
// must be wrapped by a role="group" container). No data-gsxui-* hook: unlike
// DropdownMenuRadioGroup/DropdownMenuSub, nothing in dropdown.js binds to or
// scopes by this element — it's purely a11y markup.
component DropdownMenuGroup(children gsx.Node, attrs gsx.Attrs) {
	<div role="group" { attrs... } data-gsxui-slot-dropdown-menu-group>{ children }</div>
}

// DropdownMenuCheckboxItem is the shadcn/ui DropdownMenuCheckboxItem.
// checked is server-rendered (see the file header MECHANISM); value is the
// item's own identity, stamped as data-value and echoed on dropdown.js's
// gsxui:change event. Selecting a checkbox item does NOT close the menu — a
// deliberate ADAPT per the task brief, not a Radix default (Radix's actual
// CheckboxItem onSelect closes unless preventDefault is called; this port
// simply never closes on a checkbox select — see dropdown.js).
//
// ADAPT (nova metrics + indicator side, matches DropdownMenuItem's own
// already-shipped nova pass): new-york-v4's gap-2/rounded-sm/py-1.5-pr-2-pl-8
// -> nova's gap-1.5/rounded-md/py-1-pr-8-pl-1.5 (source map lines 327/330) —
// otherwise a checkbox/radio row renders a different height and radius than
// a plain DropdownMenuItem in the same menu. The indicator side-swap
// (new-york-v4: left, pl-8 reserved; nova: right, pr-8 reserved) is real
// structure, not a metric bump (source map's own note on
// .cn-dropdown-menu-item-indicator) — resolved by following ui/select.gsx's
// SelectItem, the shipped house precedent for a nova-metric item carrying a
// check indicator: right-2, size-4 indicator span (not new-york-v4's
// size-3.5). Visibility still gates on the ancestor-selector arbitrary
// checked-state ancestor variant (not select.gsx's named ancestor marker)
// — keeps the item's own base class free of an
// extra prepended token; the two mechanisms are visually equivalent.
component DropdownMenuCheckboxItem(checked bool, value string, children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"data-disabled:pointer-events-none [&_svg:not([class*='text-'])]:text-muted-foreground focus:bg-accent focus:text-accent-foreground focus:**:text-accent-foreground min-h-7 gap-2 rounded-md py-1.5 pr-8 pl-2 text-xs [&_svg:not([class*='size-'])]:size-3.5 flex"
		}
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
		{ attrs... }
		data-gsxui-slot-dropdown-menu-checkbox-item
	>
		<span
			class={
				"pointer-events-none absolute end-2 hidden size-4 items-center justify-center [[data-state=checked]_&]:flex [&>svg]:size-4"
			}
			data-gsxui-slot-dropdown-menu-checkbox-item-indicator
		>
			<icon.Check/>
		</span>
		{ children }
	</div>
}

// DropdownMenuRadioGroup wraps a set of DropdownMenuRadioItems. value is the
// server-rendered current value, stamped as data-value on the root — the
// same "checked state is server-rendered" contract as CheckboxItem, kept in
// sync by dropdown.js on selection and echoed on the group's own
// gsxui:change event. data-gsxui-slot-dropdown-menu-radio-group is the proximity anchor
// dropdown.js uses to scope "clear every OTHER item in this group" to this
// group alone, not every radio item on the page.
component DropdownMenuRadioGroup(value string, children gsx.Node, attrs gsx.Attrs) {
	<div
		role="group"
		data-value={value}
		{ attrs... }
		data-gsxui-slot-dropdown-menu-radio-group
	>
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
// (only CheckboxItem stays open). Nova metrics + right-side indicator: same
// ADAPT as DropdownMenuCheckboxItem's own doc comment, not repeated here.
component DropdownMenuRadioItem(checked bool, value string, children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"data-disabled:pointer-events-none [&_svg:not([class*='text-'])]:text-muted-foreground focus:bg-accent focus:text-accent-foreground focus:**:text-accent-foreground min-h-7 gap-2 rounded-md py-1.5 pr-8 pl-2 text-xs [&_svg:not([class*='size-'])]:size-3.5 flex"
		}
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
		{ attrs... }
		data-gsxui-slot-dropdown-menu-radio-item
	>
		<span
			class={
				"pointer-events-none absolute end-2 hidden size-4 items-center justify-center [[data-state=checked]_&]:flex [&>svg]:size-2 [&>svg]:fill-current"
			}
			data-gsxui-slot-dropdown-menu-radio-item-indicator
		>
			<icon.Circle/>
		</span>
		{ children }
	</div>
}

// DropdownMenuLabel supports the same caller-reflected data-inset axis as
// DropdownMenuItem.
component DropdownMenuLabel(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "text-muted-foreground px-2 py-1.5 text-xs" } { attrs... } data-gsxui-slot-dropdown-menu-label>
		{ children }
	</div>
}

component DropdownMenuSeparator(attrs gsx.Attrs) {
	<div
		class={ "bg-border/50 -mx-1 my-1 h-px" }
		role="separator"
		{ attrs... }
		data-gsxui-slot-dropdown-menu-separator
	></div>
}

component DropdownMenuShortcut(children gsx.Node, attrs gsx.Attrs) {
	<span
		class={
			"text-muted-foreground group-focus/dropdown-menu-item:text-accent-foreground ml-auto text-[0.625rem] tracking-widest"
		}
		{ attrs... }
		data-gsxui-slot-dropdown-menu-shortcut
	>
		{ children }
	</span>
}

// DropdownMenuSub is the non-rendering submenu root — layout-neutral
// (class="contents", same idiom as DropdownMenu's own root) so its
// SubTrigger/SubContent children sit inline in the parent content's normal
// item flow. data-gsxui-slot-dropdown-menu-sub is the proximity anchor dropdown.js uses to
// pair a SubTrigger with its own SubContent (closest("[data-gsxui-slot-dropdown-menu-sub]")
// — same shape as DropdownMenu's own data-gsxui-slot-dropdown-menu root) and to scope
// the pointer-leave grace-period boundary check to "the whole sub."
component DropdownMenuSub(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "contents" } { attrs... } data-gsxui-slot-dropdown-menu-sub>{ children }</div>
}

// DropdownMenuSubTrigger opens/closes its sibling DropdownMenuSubContent
// (dropdown.js: pointerenter, ArrowRight, click). aria-haspopup/aria-expanded
// are server-rendered closed (derived-not-read ARIA anatomy, source map
// ## shared-items §3); dropdown.js keeps aria-expanded and data-state in step
// on every open/close, stamping data-state="open" BEFORE showPopover() (the
// same flash-avoidance rule DropdownMenuContent's own trigger click handler
// documents). The CSS-only data-inset axis remains caller-reflectable
// through attrs. ADAPT (nova metrics): gap-2/
// rounded-sm/px-2-py-1.5 -> nova's gap-1.5/rounded-md/px-1.5-py-1 (source map
// line 345), matching DropdownMenuItem's own already-shipped nova pass —
// otherwise a sub-trigger row renders taller/more-rounded than a plain
// DropdownMenuItem in the same menu. data-[state=open]: kept, not nova's
// data-open: (standing house exception).
component DropdownMenuSubTrigger(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"focus:bg-accent focus:text-accent-foreground data-[state=open]:bg-accent data-[state=open]:text-accent-foreground not-data-[variant=destructive]:focus:**:text-accent-foreground min-h-7 gap-2 rounded-md px-2 py-1 text-xs [&_svg:not([class*='size-'])]:size-3.5 flex"
		}
		role="menuitem"
		aria-haspopup="menu"
		aria-expanded="false"
		data-state="closed"
		tabindex="-1"
		{ attrs... }
		data-gsxui-slot-dropdown-menu-sub-trigger
	>
		{ children }
		<icon.ChevronRight/>
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
		class={
			"transition-none data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 ring-foreground/10 bg-popover text-popover-foreground min-w-32 rounded-lg p-1 shadow-md ring-1 duration-100"
		}
		popover="auto"
		role="menu"
		tabindex="-1"
		data-state="closed"
		data-side="right"
		{ attrs... }
		data-gsxui-slot-dropdown-menu-sub-content
	>
		{ children }
	</div>
}

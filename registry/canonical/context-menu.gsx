package canonical

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// ContextMenu is the shadcn/ui ContextMenu on the native popover API: the
// top layer replaces Radix's Portal, light dismiss and Esc are browser-
// native. This is dropdown-menu.gsx's own mechanism — reused verbatim for menu
// semantics (role="menu" content, role="menuitem" items, arrow-key roving
// focus, close-on-select) — with one structural swap: ContextMenuTrigger is
// an AREA (a plain <div>, not a <button>) that opens the menu on right-
// click rather than left-click, positioned at the cursor instead of
// anchored to the trigger's own rect. Trigger and content are wired by
// proximity — closest("[data-gsxui-slot-context-menu]") — no ids. Requires the
// context-menu behavior module (ui/context-menu.js).
//
// SUBMENUS — POPOVER NESTING, NOT PORTALLING (ADAPT, load-bearing): same
// mechanism, same measurement, as dropdown-menu.gsx's own file header (source
// map ## shared-items §5, "spec §1" below) — reused verbatim, not
// re-derived, since both components sit on the identical native-popover
// substrate:
//
//	child popover DOM-nested inside parent, opened with showPopover()      -> parent STAYS OPEN
//	child not nested, opened via a real popovertarget invoker click        -> parent STAYS OPEN
//	child not nested, opened programmatically with showPopover()           -> parent LIGHT-DISMISSES
//
// context-menu.js opens submenus on pointerenter and on ArrowRight — not
// only on click — so every open is programmatic; DOM nesting is therefore
// the only robust option here too. ContextMenuSubContent renders nested
// inside its ContextMenuSub, which itself sits inside the parent
// ContextMenuContent (pinned by
// TestContextMenuSubNestsContentInsideParentContentPinned).
//
// SERVER-RENDERED CHECKED STATE (MECHANISM): same contract as dropdown-menu.gsx's
// own — ContextMenuCheckboxItem/ContextMenuRadioItem take a `checked bool`
// that stamps both aria-checked and data-state="checked"|"unchecked" on
// first paint. The check/circle indicator icon is always rendered (no
// server-side conditional mount) and its visibility is purely CSS via the
// ancestor-selector arbitrary variant `[[data-state=checked]_&]:flex`,
// keyed off the item's own data-state — see dropdown-menu.gsx's own doc comment
// for the full rationale (same idiom, not re-derived here).
//
// THE ONE REAL SHARED-ITEMS DELTA (source map ## shared-items §2): every
// other shared part below (Group, CheckboxItem, RadioGroup, RadioItem, Sub,
// SubContent) is byte-identical, modulo the dropdown-menu/context-menu
// prefix substitution, to its dropdown-menu.gsx counterpart — ported
// mechanically. SubTrigger alone has two real, deliberate-looking-but-
// verbatim deltas in shadcn's own two sources: dropdown-menu's SubTrigger
// class carries `gap-2`, context-menu's lacks it; dropdown-menu's trailing
// ChevronRightIcon carries an explicit `size-4`, context-menu's does not.
// nova's own style-nova.css harmonizes the gap on BOTH to `gap-1.5`
// (nova wins on visual disagreement, the adjudicated house rule), so
// ContextMenuSubTrigger below takes the gap. The missing `size-4` on the
// icon is preserved as a no-op divergence from dropdown's own version — the
// class's own `[&_svg:not([class*='size-'])]:size-4` descendant selector
// already stamps size-4 onto it regardless, so there is no visual
// difference to fix and no reason to deviate from context-menu's own
// source here.
component ContextMenu(children gsx.Node, attrs gsx.Attrs) {
	<div class={ contextMenu.Root() } { attrs... } data-gsxui-slot-context-menu>{ children }</div>
}

// ContextMenuTrigger renders the drop-zone AREA a user right-clicks inside
// — shadcn's own demo renders a dashed-border div, not a button (unlike
// DropdownMenuTrigger, which IS the clickable control). context-menu.js
// listens for the `contextmenu` event within this area (event delegation
// via closest(), so any descendant right-click counts) rather than a click
// on the element itself.
component ContextMenuTrigger(children gsx.Node, attrs gsx.Attrs) {
	<div { attrs... } data-gsxui-slot-context-menu-trigger>{ children }</div>
}

// ContextMenuContent renders the popover. popover="auto" gives top layer,
// light dismiss, and free Esc; data-state is server-rendered "closed" and
// kept in sync by context-menu.js on the toggle event. Unlike
// DropdownMenuContent/PopoverContent/HoverCardContent, NO data-side is
// server-rendered here — those siblings anchor to a fixed side of their
// trigger (dropdown/popover always below, hover-card always below), so a
// static stamp matches every open; a context menu opens wherever the
// cursor was, direction never fixed, so there is no single side to stamp.
// The class string's data-[side=*]:slide-in-from-* selectors are therefore
// permanently dead weight (unlike dropdown's, where the "bottom" one is
// live) — kept for future-proofing per dropdown's own precedent, see
// docs/jsx-parity.md's ## context-menu ledger.
component ContextMenuContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={ contextMenu.Content() }
		popover="auto"
		role="menu"
		tabindex="-1"
		data-state="closed"
		{ attrs... }
		data-gsxui-slot-context-menu-content
	>
		{ children }
	</div>
}

// ContextMenuItem is the shadcn/ui ContextMenuItem, ported as a real menu
// item on a <div role="menuitem">, identical shape to
// DropdownMenuItem — context-menu.js's arrow-key roving focus walks these
// the same way dropdown-menu.js's does. variant: "" (default) | "destructive".
// Callers reflect the CSS-only inset axis with data-inset through attrs.
component ContextMenuItem(variant string, children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"group/context-menu-item",
			contextMenu.Item(),
			contextMenu.ItemVariant(variant),
		}
		data-variant={variant |> default("default")}
		role="menuitem"
		tabindex="-1"
		{ attrs... }
		data-gsxui-slot-context-menu-item
	>
		{ children }
	</div>
}

// ContextMenuGroup wraps a set of items for a11y grouping. shadcn's own
// Group carries no class string at all in either source (source map ##
// shared-items §1) — role="group" is added here, not in the .tsx (same
// derived-not-read WAI-ARIA menu authoring practice as DropdownMenuGroup's
// own doc comment). No data-gsxui-* hook: nothing in context-menu.js binds
// to or scopes by this element — it's purely a11y markup, same call as
// DropdownMenuGroup.
component ContextMenuGroup(children gsx.Node, attrs gsx.Attrs) {
	<div role="group" { attrs... } data-gsxui-slot-context-menu-group>{ children }</div>
}

// ContextMenuCheckboxItem is the shadcn/ui ContextMenuCheckboxItem. checked
// is server-rendered (see the file header MECHANISM); value is the item's
// own identity, stamped as data-value and echoed on context-menu.js's
// gsxui:change event. Selecting a checkbox item does NOT close the menu —
// same deliberate ADAPT as DropdownMenuCheckboxItem's own doc comment, not
// a Radix default.
//
// Class string, nova metrics, and right-side indicator are byte-identical
// to DropdownMenuCheckboxItem's own (source map ## shared-items §1: verbatim
// modulo the component-name prefix) — see that component's own doc comment
// for the full ADAPT rationale (nova's gap-1.5/rounded-md/py-1-pr-8-pl-1.5
// over new-york-v4's gap-2/rounded-sm/py-1.5-pr-2-pl-8, right-side indicator
// following ui/select.gsx's SelectItem precedent), not re-derived here.
component ContextMenuCheckboxItem(checked bool, value string, children gsx.Node, attrs gsx.Attrs) {
	<div
		class={ contextMenu.CheckboxItem() }
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
		data-gsxui-slot-context-menu-checkbox-item
	>
		<span class={ contextMenu.CheckboxItemIndicator() } data-gsxui-slot-context-menu-checkbox-item-indicator>
			<icon.Check/>
		</span>
		{ children }
	</div>
}

// ContextMenuRadioGroup wraps a set of ContextMenuRadioItems. value is the
// server-rendered current value, stamped as data-value on the root — same
// server-rendered-checked contract as CheckboxItem, kept in sync by
// context-menu.js on selection and echoed on the group's own gsxui:change
// event. data-gsxui-slot-context-menu-radio-group is the proximity anchor context-menu.js
// uses to scope "clear every OTHER item in this group" to this group alone.
component ContextMenuRadioGroup(value string, children gsx.Node, attrs gsx.Attrs) {
	<div
		role="group"
		data-value={value}
		{ attrs... }
		data-gsxui-slot-context-menu-radio-group
	>
		{ children }
	</div>
}

// ContextMenuRadioItem is the shadcn/ui ContextMenuRadioItem — same shape/
// class as CheckboxItem (source map ## shared-items §1: byte-identical base
// style, own finding, true in both files), swapping the check indicator for
// a filled dot. checked reflects whether THIS item's value equals its
// ContextMenuRadioGroup's current value — the caller computes that
// comparison (gsx has no context to do it implicitly). Selecting a radio
// item DOES close the menu, same as a plain Item (only CheckboxItem stays
// open). Nova metrics + right-side indicator: same ADAPT as
// ContextMenuCheckboxItem's own doc comment, not repeated here.
component ContextMenuRadioItem(checked bool, value string, children gsx.Node, attrs gsx.Attrs) {
	<div
		class={ contextMenu.RadioItem() }
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
		data-gsxui-slot-context-menu-radio-item
	>
		<span class={ contextMenu.RadioItemIndicator() } data-gsxui-slot-context-menu-radio-item-indicator>
			<icon.Circle/>
		</span>
		{ children }
	</div>
}

// ContextMenuLabel supports the same caller-reflected data-inset axis as
// ContextMenuItem.
// Unlike DropdownMenuLabel, shadcn's own context-menu.tsx class carries
// text-foreground — ported verbatim, a genuine per-component difference in
// the shadcn source, not a copy error (see docs/jsx-parity.md ##
// context-menu).
component ContextMenuLabel(children gsx.Node, attrs gsx.Attrs) {
	<div class={ contextMenu.Label() } { attrs... } data-gsxui-slot-context-menu-label>{ children }</div>
}

component ContextMenuSeparator(attrs gsx.Attrs) {
	<div class={ contextMenu.Separator() } role="separator" { attrs... } data-gsxui-slot-context-menu-separator></div>
}

component ContextMenuShortcut(children gsx.Node, attrs gsx.Attrs) {
	<span class={ contextMenu.Shortcut() } { attrs... } data-gsxui-slot-context-menu-shortcut>
		{ children }
	</span>
}

// ContextMenuSub is the non-rendering submenu root — layout-neutral
// (class="contents", same idiom as ContextMenu's own root and
// DropdownMenuSub) so its SubTrigger/SubContent children sit inline in the
// parent content's normal item flow. data-gsxui-slot-context-menu-sub is the proximity
// anchor context-menu.js uses to pair a SubTrigger with its own SubContent
// and to scope the pointer-leave grace-period boundary check to "the whole
// sub" — same shape as DropdownMenuSub's own doc comment.
component ContextMenuSub(children gsx.Node, attrs gsx.Attrs) {
	<div class={ contextMenu.Sub() } { attrs... } data-gsxui-slot-context-menu-sub>{ children }</div>
}

// ContextMenuSubTrigger opens/closes its sibling ContextMenuSubContent
// (context-menu.js: pointerenter, ArrowRight, click). aria-haspopup/
// aria-expanded are server-rendered closed (derived-not-read ARIA anatomy,
// source map ## shared-items §3); context-menu.js keeps aria-expanded and
// data-state in step on every open/close, stamping data-state="open" BEFORE
// showPopover() (same flash-avoidance rule as ContextMenuContent's own
// opener). The CSS-only data-inset axis remains caller-reflectable through
// attrs. ADAPT (nova metrics + the one
// real shared-items delta): gap-1.5/rounded-md/px-1.5-py-1 — see the file
// header's "THE ONE REAL SHARED-ITEMS DELTA" paragraph for the gap-2->
// gap-1.5 harmonization and the deliberately-preserved missing icon size-4.
// data-[state=open]: kept, not nova's data-open: (standing house exception).
component ContextMenuSubTrigger(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={ contextMenu.SubTrigger() }
		role="menuitem"
		aria-haspopup="menu"
		aria-expanded="false"
		data-state="closed"
		tabindex="-1"
		{ attrs... }
		data-gsxui-slot-context-menu-sub-trigger
	>
		{ children }
		<icon.ChevronRight/>
	</div>
}

// ContextMenuSubContent is the submenu popover — see the file header's
// SUBMENUS comment for why it must render DOM-nested (not portalled) inside
// its ContextMenuSub. Class string is byte-identical, modulo the
// dropdown-menu/context-menu CSS-var-name prefix in
// origin-(--radix-*-content-transform-origin) (source map ## shared-items
// §1), to DropdownMenuSubContent's own — same ADAPT list applies verbatim,
// not repeated here: origin-top-left (no Radix runtime CSS var),
// data-side="right" server-rendered statically (the SubTrigger, unlike the
// cursor-positioned top-level ContextMenuContent, IS a fixed anchor —
// context-menu.js always opens its submenus to the right of their
// trigger), the discrete-transition enter/exit block replacing the six
// data-[state=…]:animate-*/fade/zoom tokens, and nova's min-w-[96px]/
// rounded-lg metrics (border kept, not nova's ring-1 — standing house
// exception).
component ContextMenuSubContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={ contextMenu.SubContent() }
		popover="auto"
		role="menu"
		tabindex="-1"
		data-state="closed"
		data-side="right"
		{ attrs... }
		data-gsxui-slot-context-menu-sub-content
	>
		{ children }
	</div>
}

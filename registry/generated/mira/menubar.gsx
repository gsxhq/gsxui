package ui

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Menubar is the shadcn/ui Menubar on the native popover API — the same
// dropdown-menu.gsx/context-menu.gsx mechanism (top layer replaces Radix's
// Portal, light dismiss and Esc are browser-native, role="menu" content,
// role="menuitem" items, submenu popover nesting) reused verbatim for the
// seven item parts this component shares with Task 1/2's own (Item,
// CheckboxItem, RadioGroup, RadioItem, Sub, SubTrigger, SubContent — see
// each part's own doc comment for anything genuinely menubar-specific).
// What's new here is the BAR: Menubar renders a row of MenubarTrigger
// buttons, each paired 1:1 with its own MenubarContent inside a
// MenubarMenu wrapper (the "root" proximity scope for that one menu, the
// same role DropdownMenu's own root plays for its single trigger/content
// pair) — and menubar.js layers two behaviors on top that neither dropdown
// nor context-menu need: roving tabindex across the triggers (the whole bar
// is one tab stop; ArrowLeft/ArrowRight move it) and open-follows-hover
// (once one menu is open, hovering a sibling trigger switches to it with no
// click). Both are DERIVED-NOT-READ from the public Radix Menubar contract
// per the source map (`docs/superpowers/plans/2026-07-25-tier4-source-map-
// menus.md` `## menubar` §2/§3) — no `@radix-ui/react-menubar` dist exists
// in the reference checkout to trace, so this is built to the documented
// public behavior, not verified against Radix's own runtime. Requires the
// menubar behavior module (ui/menubar.js).
//
// HOOK NAMESPACING (load-bearing, THE fix from the prior round): every
// selector menubar.js registers is `data-gsxui-slot-menubar-*`, never a prefix
// shared with dropdown-menu.js's `data-gsxui-slot-dropdown-menu-*` or context-menu.js's
// `data-gsxui-slot-context-menu-*`. ui/gsxui.js's delegation registry is keyed
// only by `${type}:${capture}` and dispatches to EVERY handler whose
// selector matches the event target, regardless of which module registered
// it — ui/index.js imports all three menu modules, so an identical selector
// across two of them double-fires on one event (measured live on an
// earlier round: one click on a checkbox item fired two gsxui:change events
// and net-zeroed the toggle back to its starting state; click-to-open
// submenus opened then immediately closed; arrow keys skipped every other
// item). See docs/jsx-parity.md's ## dropdown ledger MECHANISM entry for
// the full incident writeup.
//
// SUBMENUS — POPOVER NESTING, NOT PORTALLING (ADAPT, load-bearing): same
// mechanism, same measurement, as dropdown-menu.gsx's own file header — reused
// verbatim, not re-derived:
//
//	child popover DOM-nested inside parent, opened with showPopover()      -> parent STAYS OPEN
//	child not nested, opened via a real popovertarget invoker click        -> parent STAYS OPEN
//	child not nested, opened programmatically with showPopover()           -> parent LIGHT-DISMISSES
//
// menubar.js opens submenus on pointerenter and on ArrowRight — not only on
// click — so every open is programmatic; DOM nesting is therefore the only
// robust option here too. MenubarSubContent renders nested inside its
// MenubarSub, which itself sits inside the parent MenubarContent (pinned by
// TestMenubarSubNestsContentInsideParentContentPinned). Unlike dropdown/
// context-menu's own top-level Content, MenubarContent is NOT a scrollable
// clipped container (overflow-hidden only, no max-h/overflow-y-auto — see
// MenubarContent's own doc comment and the source map's `## shared-items`
// §5), so the nested-submenu-gets-clipped-by-its-scrolling-ancestor risk
// that section flags for dropdown/context-menu doesn't apply here at all.
//
// SERVER-RENDERED CHECKED STATE (MECHANISM): same contract as dropdown-menu.gsx's
// own — MenubarCheckboxItem/MenubarRadioItem take a `checked bool` that
// stamps both aria-checked and data-state="checked"|"unchecked" on first
// paint. The check/circle indicator icon is always rendered (no server-side
// conditional mount) and its visibility is purely CSS via the
// ancestor-selector arbitrary variant `[[data-state=checked]_&]:flex`,
// keyed off the item's own data-state — see dropdown-menu.gsx's own doc comment
// for the full rationale, not re-derived here. Selecting a checkbox item
// does NOT close the menu — the same deliberate ADAPT as dropdown/
// context-menu's own CheckboxItem, not a Radix default.
//
// NOVA METRICS: Item and CheckboxItem/RadioItem below end up BYTE-IDENTICAL,
// modulo the component-name prefix, to their already-shipped
// DropdownMenuItem/CheckboxItem/RadioItem counterparts — new-york-v4's own
// menubar.tsx carries a real per-component delta from dropdown-menu/
// context-menu (CheckboxItem/RadioItem's rounded-xs vs their rounded-sm),
// which nova's own style-nova.css erases: `.cn-menubar-checkbox-item`/
// `.cn-menubar-radio-item` land on the same rounded-md every other
// item-shaped part in this whole menu family uses.
//
// CORRECTION (indicator side is a DELIBERATE OVERRIDE, not agreement with
// nova — an earlier version of this comment claimed nova made menubar's
// CheckboxItem/RadioItem identical to dropdown's own; that was wrong on the
// indicator geometry specifically): nova's own `.cn-menubar-checkbox-item`/
// `.cn-menubar-radio-item` are `py-1 pr-1.5 pl-7`, and
// `.cn-menubar-checkbox-item-indicator`/`-radio-item-indicator` are
// `left-1.5 size-4` — nova genuinely keeps MENUBAR's own indicator on the
// LEFT, unlike its dropdown/context-menu rules (both `right-2`). This port
// does NOT follow nova here: the indicator stays at the RIGHT edge
// (`right-2`, `size-4`, `pr-8 pl-1.5` rows), matching DropdownMenuCheckboxItem/
// ContextMenuCheckboxItem and following ui/select.gsx's SelectItem — a
// deliberate ADAPT for cross-menu-family visual consistency (all three menu
// components' checkbox/radio rows read the same way), not nova's own
// per-component value for menubar specifically. Ledgered so a future reader
// diffing against style-nova.css doesn't find this port "wrong" by nova's
// own menubar-specific rule.
//
// SubTrigger does NOT join the byte-identical set above: new-york-v4's own
// menubar.tsx SubTrigger carries no gap-* token AND, unlike dropdown-menu.tsx/
// context-menu.tsx's own SubTrigger, no `[&_svg]:pointer-events-none`/
// `[&_svg]:shrink-0`/`[&_svg:not([class*='text-'])]:text-muted-foreground`
// at all — confirmed by re-reading the source map's own `## menubar` §1
// quote (no svg selectors present). nova's own `.cn-menubar-sub-trigger`
// adds `gap-1.5` (ported, harmonizing the same upstream asymmetry already
// fixed for context-menu's own SubTrigger) but is SILENT on the
// pointer-events/shrink/muted-color triple — as is EVERY nova SubTrigger
// rule in this menu family, `.cn-dropdown-menu-sub-trigger`/
// `.cn-context-menu-sub-trigger` included, so nova is not the source of
// dropdown/context-menu's own three tokens either; those come from
// dropdown-menu.tsx/context-menu.tsx's own (non-menubar) source, which
// menubar.tsx's own source never had to begin with. RULING: ported
// source-faithfully — MenubarSubTrigger's chevron renders in the default
// foreground color, NOT muted like DropdownMenuSubTrigger's/
// ContextMenuSubTrigger's own chevrons — a genuine, deliberate-per-source
// divergence, not an oversight or a nova disagreement to fix.
//
// TWO DELIBERATE DIVERGENCES FROM new-york-v4's OWN menubar.tsx (both
// flagged by the source map, both ruled on here rather than silently
// ported or silently fixed):
//   - MenubarSubTrigger's class carries `outline-hidden`, NOT new-york-v4's
//     own `outline-none` (menubar.tsx's one-off spelling — every other
//     Item/SubTrigger class in this whole codebase, across all three menu
//     families, uses outline-hidden). The source map calls this a real,
//     if narrow, accessibility delta: `outline-none` removes the focus
//     ring in forced-colors/high-contrast modes too, `outline-hidden` is
//     Tailwind v4's own "hidden but still present for forced-colors" idiom.
//     RULING: port outline-hidden — matching this port's own existing
//     house-wide convention beats reproducing a one-off, less-accessible
//     upstream spelling, and every other menubar part (Item, CheckboxItem,
//     etc.) already uses outline-hidden a few lines away in the very same
//     source file, so outline-none reads as menubar.tsx's own copy-paste
//     drift, not a deliberate per-component choice worth preserving.
//   - MenubarSubTrigger's trailing chevron is ported as `size-4`, not
//     new-york-v4's own literal `h-4 w-4` (menubar's own third distinct
//     spelling for the same 16px icon, alongside dropdown's `size-4` and
//     context-menu's bare/unstyled one). Same no-op-divergence call already
//     ledgered for context-menu's own SubTrigger icon: the class's own
//     `[&_svg:not([class*='size-'])]:size-4` descendant selector already
//     stamps size-4 onto any child svg lacking its own size-* class, so
//     h-4/w-4 and size-4 render pixel-identical either way — porting
//     size-4 directly keeps the markup consistent with every other icon
//     literal in this codebase instead of introducing Tailwind v3-style
//     `h-* w-*` pairs nowhere else in gsxui.
//
// MenubarContent's animate-out RULING: the source map flags that
// MenubarContent's own class string (unlike every sibling Content/
// SubContent in this whole menu family, including menubar's own
// SubContent) omits `data-[state=closed]:animate-out` from its otherwise-
// standard six-token animate-in/out/fade/zoom set — corroborated as
// deliberate, not a transcription slip, since nova's independently-authored
// `.cn-menubar-content` has the identical omission. This port's own ADAPT
// (already established by dropdown-menu.gsx/context-menu.gsx, see popover.gsx
// and docs/jsx-parity.md ## animations) replaces that ENTIRE six-token
// tw-animate-css family, on every popover-based Content/SubContent in this
// codebase without exception, with one discrete-transition CSS block keyed
// off Tailwind's `open:`/`starting:open:` variants instead — a mechanism
// that has no equivalent to a standalone "closed-state token" to omit in
// the first place: the SAME `open:opacity-100`/base-`opacity-0` pair drives
// both directions, there is no separate switch for "animate the close" a
// port could selectively leave off. RULING: MenubarContent gets the
// identical discrete-transition block every other Content part gets, same
// as MenubarSubContent — new-york-v4's own omission is recorded here as
// what the source says, per the task's own instruction, but has zero
// representation in this port's markup, because the token it omits was
// already replaced site-wide before this component existed.
component Menubar(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={ "h-9 rounded-lg border p-1 flex" }
		role="menubar"
		{ attrs... }
		data-gsxui-slot-menubar
	>
		{ children }
	</div>
}

// MenubarMenu is the non-rendering root pairing ONE MenubarTrigger with its
// own MenubarContent — layout-neutral (class="contents", same idiom as
// DropdownMenu's own root) so the pair sits inline in the bar's normal flex
// row. data-gsxui-slot-menubar-menu is the proximity anchor menubar.js uses to
// resolve "this trigger's own content" (closest("[data-gsxui-slot-menubar-
// menu]")), the same shape DropdownMenu's own root plays for its single
// trigger/content pair — Menubar itself (the outer bar) is the SEPARATE
// scope roving tabindex and open-follows-hover coordinate across every
// MenubarMenu's trigger.
component MenubarMenu(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "contents" } { attrs... } data-gsxui-slot-menubar-menu>{ children }</div>
}

// MenubarTrigger is one pill in the bar. Unlike DropdownMenuTrigger (which
// carries no class at all — callers style their own button), MenubarTrigger
// IS the visible, styled control shadcn's own menubar.tsx renders directly
// (no asChild wrapping in the reference demo), so its class string is
// ported for real. No tabindex is server-rendered — roving tabindex is
// JS-normalized-at-init (ui/toggle-group.js's own MECHANISM, reused here):
// every trigger is its own tab stop until menubar.js collapses the bar to
// exactly one on module load, the graceful no-JS fallback. data-state is
// server-rendered "closed" and kept in sync by menubar.js on the toggle
// event — MenubarTrigger's own class keys :open highlighting off
// data-state (nova's rounded-sm/px-1.5/py-[2px] metrics, NOT bumped to
// rounded-md — nova's own .cn-menubar-trigger literally keeps rounded-sm,
// unlike every item-shaped part in this file), unlike DropdownMenuTrigger,
// which has no such selector to key at all.
component MenubarTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		class={
			"hover:bg-muted aria-expanded:bg-muted rounded-[calc(var(--radius-md)-2px)] px-2 py-[calc(--spacing(0.85))] text-xs/relaxed font-medium flex"
		}
		type="button"
		aria-haspopup="menu"
		aria-expanded="false"
		data-state="closed"
		{ attrs... }
		data-gsxui-slot-menubar-trigger
	>
		{ children }
	</button>
}

// MenubarContent renders the popover for one MenubarMenu. popover="auto"
// gives top layer, light dismiss, and free Esc; data-state is
// server-rendered "closed" and kept in sync by menubar.js on the toggle
// event. data-side="bottom" is server-rendered statically — menubar.js
// always anchors below the trigger, same hand-rolled-fixed-position
// stopgap as dropdown's own (see docs/jsx-parity.md ## dropdown NOTE).
// overflow-hidden only (NOT dropdown's own overflow-x-hidden/
// overflow-y-auto): new-york-v4's own MenubarContent carries no max-h/
// overflow-y-auto at all, unlike DropdownMenuContent/ContextMenuContent —
// it is not a scrollable clipped container, so the nested-submenu-clipped-
// by-a-scrolling-ancestor risk the source map flags for dropdown/
// context-menu's own SubContent (`## shared-items` §5) never arises here.
// min-w-36 is nova's own genuine delta from new-york-v4's min-w-[12rem]
// (12rem -> 9rem, a real narrowing, not a spelling variance — contrast
// MenubarSubContent's own min-w-[8rem], where nova's value and
// new-york-v4's coincide exactly and no delta applies). See the file
// header's own RULING for why this class carries the identical discrete-
// transition block every sibling Content/SubContent in this codebase
// carries, despite new-york-v4's own class omitting
// data-[state=closed]:animate-out.
component MenubarContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"transition-none bg-popover text-popover-foreground data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 ring-foreground/10 min-w-32 rounded-lg p-1 shadow-md ring-1 duration-100"
		}
		popover="auto"
		role="menu"
		tabindex="-1"
		data-state="closed"
		data-side="bottom"
		{ attrs... }
		data-gsxui-slot-menubar-content
	>
		{ children }
	</div>
}

// MenubarItem is byte-identical, modulo the slot-attribute prefix,
// to DropdownMenuItem's own already-shipped pinned class — the source
// map's own finding that plain Item is identical across all three menu
// families, and this port's own nova retarget already applied uniformly to
// dropdown/context-menu's own Item. variant: "" (default) | "destructive".
// Callers reflect the CSS-only inset axis with data-inset through attrs.
component MenubarItem(variant string, children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"group/menubar-item",
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
		data-gsxui-slot-menubar-item
	>
		{ children }
	</div>
}

// MenubarGroup wraps a set of items for a11y grouping. shadcn's own Group
// carries no class string at all (source map `## menubar` §1) — role="group"
// is added here, not in the .tsx, same derived-not-read WAI-ARIA menu
// authoring practice as DropdownMenuGroup's own doc comment. No
// data-gsxui-* hook: nothing in menubar.js binds to or scopes by this
// element, same call as DropdownMenuGroup/ContextMenuGroup.
component MenubarGroup(children gsx.Node, attrs gsx.Attrs) {
	<div role="group" { attrs... } data-gsxui-slot-menubar-group>{ children }</div>
}

// MenubarCheckboxItem is the shadcn/ui MenubarCheckboxItem. checked is
// server-rendered (see the file header MECHANISM); value is the item's own
// identity, stamped as data-value and echoed on menubar.js's gsxui:change
// event. Selecting a checkbox item does NOT close the menu — the same
// deliberate ADAPT as dropdown/context-menu's own CheckboxItem.
//
// Ends up byte-identical, modulo the prefix, to DropdownMenuCheckboxItem's
// own pinned class: new-york-v4's one real per-component delta here
// (rounded-xs vs dropdown/context's rounded-sm) is erased by nova, which
// unifies checkbox/radio/item/sub-trigger all onto rounded-md (source
// map's own "Metric deltas worth naming" finding) — see
// DropdownMenuCheckboxItem's own doc comment for the full nova-metrics
// rationale, not re-derived here. The right-side indicator is NOT what
// nova's own menubar rule specifies, though — see the file header's own
// CORRECTION paragraph: nova's `.cn-menubar-checkbox-item`/`-indicator`
// keep the indicator on the LEFT (`pl-7`/`left-1.5`); this port overrides
// that with dropdown/context-menu/select's own right-side geometry
// deliberately, for cross-menu-family consistency, not because nova agrees.
component MenubarCheckboxItem(checked bool, value string, children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"data-disabled:pointer-events-none [&_svg:not([class*='text-'])]:text-muted-foreground focus:bg-accent focus:text-accent-foreground focus:**:text-accent-foreground min-h-7 gap-2 rounded-md py-1.5 pr-2 pl-7.5 text-xs flex"
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
		data-gsxui-slot-menubar-checkbox-item
	>
		<span
			class={
				"pointer-events-none absolute hidden left-2 size-4 items-center justify-center [[data-state=checked]_&]:flex [&_svg:not([class*='size-'])]:size-4"
			}
			data-gsxui-slot-menubar-checkbox-item-indicator
		>
			<icon.Check/>
		</span>
		{ children }
	</div>
}

// MenubarRadioGroup wraps a set of MenubarRadioItems. value is the
// server-rendered current value, stamped as data-value on the root — the
// same server-rendered-checked contract as CheckboxItem, kept in sync by
// menubar.js on selection and echoed on the group's own gsxui:change event.
// data-gsxui-slot-menubar-radio-group is the proximity anchor menubar.js uses to
// scope "clear every OTHER item in this group" to this group alone.
component MenubarRadioGroup(value string, children gsx.Node, attrs gsx.Attrs) {
	<div role="group" data-value={value} { attrs... } data-gsxui-slot-menubar-radio-group>
		{ children }
	</div>
}

// MenubarRadioItem is the shadcn/ui MenubarRadioItem — same shape/class as
// CheckboxItem, swapping the check indicator for a filled dot. checked
// reflects whether THIS item's value equals its MenubarRadioGroup's current
// value — the caller computes that comparison (gsx has no context to do it
// implicitly). Selecting a radio item DOES close the menu, same as a plain
// Item (only CheckboxItem stays open). Same nova metrics + deliberate
// right-side-indicator override (NOT nova's own left-side menubar value) as
// MenubarCheckboxItem's own doc comment.
component MenubarRadioItem(checked bool, value string, children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"data-disabled:pointer-events-none [&_svg:not([class*='text-'])]:text-muted-foreground focus:bg-accent focus:text-accent-foreground focus:**:text-accent-foreground min-h-7 gap-2 rounded-md py-1.5 pr-2 pl-7.5 text-xs [&_svg:not([class*='size-'])]:size-3.5 flex"
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
		data-gsxui-slot-menubar-radio-item
	>
		<span
			class={
				"pointer-events-none absolute hidden left-2 size-4 items-center justify-center [[data-state=checked]_&]:flex [&>svg]:size-2 [&>svg]:fill-current"
			}
			data-gsxui-slot-menubar-radio-item-indicator
		>
			<icon.Circle/>
		</span>
		{ children }
	</div>
}

// MenubarLabel supports the same caller-reflected data-inset axis as
// MenubarItem.
// text-sm, NOT DropdownMenuLabel's own text-xs: nova's own
// .cn-menubar-label rule is genuinely text-sm, a real per-component nova
// value, not a copy of dropdown's own (already-shipped) Label metrics.
component MenubarLabel(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "text-muted-foreground px-2 py-1.5 text-xs" } { attrs... } data-gsxui-slot-menubar-label>
		{ children }
	</div>
}

component MenubarSeparator(attrs gsx.Attrs) {
	<div class={ "bg-border/50 -mx-1 my-1 h-px" } role="separator" { attrs... } data-gsxui-slot-menubar-separator></div>
}

// The automatic leading margin is KEPT even though nova's own shortcut rule
// omits it (its focus color, text size, and tracking rule has no automatic
// margin) — undisclosed until now, not
// a metric token so not itself in scope for the nova retarget, and dropping
// it would break the shortcut's own push-to-the-right layout inside the
// item row (every sibling Shortcut in this codebase keeps it). The
// The named-ancestor focus color addition nova carries is NOT ported either
// — no corresponding named ancestor marker exists anywhere
// in this file for it to key off, the same scope call as every other
// dropdown/context-menu Shortcut in this codebase.
component MenubarShortcut(children gsx.Node, attrs gsx.Attrs) {
	<span
		class={ "text-muted-foreground group-focus/menubar-item:text-accent-foreground text-[0.625rem] tracking-widest" }
		{ attrs... }
		data-gsxui-slot-menubar-shortcut
	>
		{ children }
	</span>
}

// MenubarSub is the non-rendering submenu root — layout-neutral
// (class="contents", same idiom as MenubarMenu's own root) so its
// SubTrigger/SubContent children sit inline in the parent content's normal
// item flow. data-gsxui-slot-menubar-sub is the proximity anchor menubar.js uses
// to pair a SubTrigger with its own SubContent (closest("[data-gsxui-
// slot-menubar-sub]")) and to scope the pointer-leave grace-period boundary
// check to "the whole sub" — same shape as DropdownMenuSub's own doc
// comment.
component MenubarSub(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "contents" } { attrs... } data-gsxui-slot-menubar-sub>{ children }</div>
}

// MenubarSubTrigger opens/closes its sibling MenubarSubContent (menubar.js:
// pointerenter, ArrowRight, click). aria-haspopup/aria-expanded are
// server-rendered closed (derived-not-read ARIA anatomy, source map `##
// shared-items` §3); menubar.js keeps aria-expanded and data-state in step
// on every open/close, stamping data-state="open" BEFORE showPopover() (the
// same flash-avoidance rule DropdownMenuContent's own trigger click handler
// documents). The CSS-only data-inset axis remains caller-reflectable
// through attrs. See the file header's own "TWO
// DELIBERATE DIVERGENCES" paragraph for the outline-hidden (not
// new-york-v4's own outline-none) and size-4 (not h-4 w-4) rulings — both
// unique to menubar.tsx among the three menu family sources, both decided
// here rather than silently ported or silently fixed. gap-1.5 is nova's own
// addition (new-york-v4's own MenubarSubTrigger carries no gap-* token at
// all, the "third" version of the same asymmetry context-menu's own
// SubTrigger already got harmonized for — see the file header). nova's
// rounded-md/px-1.5/py-1 metrics, same as MenubarItem's own.
// data-[state=open]: kept, not nova's data-open: (standing house
// exception).
//
// NOT byte-identical to DropdownMenuSubTrigger/ContextMenuSubTrigger (an
// earlier version of this comment implied it was, via the file header's own
// now-corrected "NOVA METRICS" paragraph — see that CORRECTION): this class
// carries no `[&_svg]:pointer-events-none`, `[&_svg]:shrink-0`, or
// `[&_svg:not([class*='text-'])]:text-muted-foreground` — new-york-v4's own
// menubar.tsx SubTrigger never had those three tokens to begin with (unlike
// dropdown-menu.tsx/context-menu.tsx's own), and nova's `.cn-menubar-sub-
// trigger` rule is silent on them too (so is EVERY nova SubTrigger rule in
// this menu family — nova isn't where dropdown/context-menu's own three
// tokens came from either). Ported source-faithfully: this chevron renders
// in the default foreground color, not muted like its dropdown/context-menu
// siblings — a real, deliberate divergence, not an oversight.
component MenubarSubTrigger(children gsx.Node, attrs gsx.Attrs) {
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
		data-gsxui-slot-menubar-sub-trigger
	>
		{ children }
		<icon.ChevronRight/>
	</div>
}

// MenubarSubContent is the submenu popover — see the file header's
// SUBMENUS comment for why it must render DOM-nested (not portalled) inside
// its MenubarSub. Has the FULL data-[state=…]:animate-in/out/fade/zoom
// six-token set in new-york-v4's own source (unlike MenubarContent's own,
// see the file header's animate-out RULING) — moot for this port either
// way, since both get the identical discrete-transition block. min-w-[8rem]
// is KEPT (not narrowed): nova's own min-w-32 is exactly 8rem, the same
// value new-york-v4 already carries — no genuine delta here, unlike
// MenubarContent's own min-w-36 (a real 12rem -> 9rem narrowing) or
// DropdownMenuSubContent's own min-w-[96px] (a real 8rem -> 6rem
// narrowing). rounded-lg (nova), shadow-lg (matches both), border kept
// (house exception, not nova's ring-1) — same ADAPT list as
// DropdownMenuSubContent's own doc comment, not repeated in full here.
component MenubarSubContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"transition-none bg-popover text-popover-foreground data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 ring-foreground/10 min-w-32 rounded-lg p-1 shadow-md ring-1 duration-100"
		}
		popover="auto"
		role="menu"
		tabindex="-1"
		data-state="closed"
		data-side="right"
		{ attrs... }
		data-gsxui-slot-menubar-sub-content
	>
		{ children }
	</div>
}

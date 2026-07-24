package ui

import "github.com/gsxhq/gsx"

// ContextMenu is the shadcn/ui ContextMenu on the native popover API: the
// top layer replaces Radix's Portal, light dismiss and Esc are browser-
// native. This is dropdown.gsx's own mechanism — reused verbatim for menu
// semantics (role="menu" content, role="menuitem" items, arrow-key roving
// focus, close-on-select) — with one structural swap: ContextMenuTrigger is
// an AREA (a plain <div>, not a <button>) that opens the menu on right-
// click rather than left-click, positioned at the cursor instead of
// anchored to the trigger's own rect. Trigger and content are wired by
// proximity — closest("[data-gsxui-contextmenu]") — no ids. Requires the
// context-menu behavior module (ui/context-menu.js).
component ContextMenu(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="context-menu" data-gsxui-contextmenu class="contents" { attrs... }>{ children }</div>
}

// ContextMenuTrigger renders the drop-zone AREA a user right-clicks inside
// — shadcn's own demo renders a dashed-border div, not a button (unlike
// DropdownMenuTrigger, which IS the clickable control). context-menu.js
// listens for the `contextmenu` event within this area (event delegation
// via closest(), so any descendant right-click counts) rather than a click
// on the element itself.
component ContextMenuTrigger(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="context-menu-trigger" data-gsxui-contextmenu-trigger { attrs... }>{ children }</div>
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
		data-slot="context-menu-content"
		data-gsxui-contextmenu-content
		popover="auto"
		role="menu"
		tabindex="-1"
		data-state="closed"
		class={
			"z-50 max-h-96 min-w-36 origin-top-left overflow-x-hidden overflow-y-auto rounded-lg border bg-popover p-1 text-popover-foreground shadow-md",
			// Discrete-transition enter/exit replacing the tw-animate keyframe
			// pair — a popover's exit keyframe never gets to play (hide is
			// instant display:none); see popover.gsx's ADAPT comment and
			// docs/jsx-parity.md ## animations for the full mechanism. The
			// data-[side=…] starting slides stay inert here exactly as the
			// slide-in-from-* tokens they replace were: no data-side is ever
			// stamped (see this file's header comment).
			"opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:opacity-0 starting:open:scale-95",
			"data-[side=bottom]:starting:open:-translate-y-2 data-[side=left]:starting:open:translate-x-2 data-[side=right]:starting:open:-translate-x-2 data-[side=top]:starting:open:translate-y-2"
		}
		{ attrs... }
	>
		{ children }
	</div>
}

// ContextMenuItem is the shadcn/ui ContextMenuItem, ported as a real menu
// item on a <div role="menuitem">, identical shape to
// DropdownMenuItem — context-menu.js's arrow-key roving focus walks these
// the same way dropdown.js's does. variant: "" (default) | "destructive".
// inset is dropped (see docs/jsx-parity.md), same call as dropdown's own —
// the data-[inset]:pl-8 selector is removed with it.
component ContextMenuItem(variant string, children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="context-menu-item"
		data-gsxui-contextmenu-item
		data-variant={variant |> default("default")}
		role="menuitem"
		tabindex="-1"
		class="relative flex cursor-default items-center gap-1.5 rounded-md px-1.5 py-1 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 data-[variant=destructive]:text-destructive data-[variant=destructive]:focus:bg-destructive/10 data-[variant=destructive]:focus:text-destructive dark:data-[variant=destructive]:focus:bg-destructive/20 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground data-[variant=destructive]:*:[svg]:text-destructive!"
		{ attrs... }
	>
		{ children }
	</div>
}

// ContextMenuLabel's inset prop is dropped along with ContextMenuItem's
// (see docs/jsx-parity.md) — the data-[inset]:pl-8 selector is removed.
// Unlike DropdownMenuLabel, shadcn's own context-menu.tsx class carries
// text-foreground — ported verbatim, a genuine per-component difference in
// the shadcn source, not a copy error (see docs/jsx-parity.md ##
// context-menu).
component ContextMenuLabel(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="context-menu-label" class="px-1.5 py-1 text-xs font-medium text-foreground" { attrs... }>{ children }</div>
}

component ContextMenuSeparator(attrs gsx.Attrs) {
	<div data-slot="context-menu-separator" role="separator" class="-mx-1 my-1 h-px bg-border" { attrs... }></div>
}

component ContextMenuShortcut(children gsx.Node, attrs gsx.Attrs) {
	<span data-slot="context-menu-shortcut" class="ml-auto text-xs tracking-widest text-muted-foreground" { attrs... }>
		{ children }
	</span>
}

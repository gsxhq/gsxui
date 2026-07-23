package ui

import "github.com/gsxhq/gsx"

// DropdownMenu is the shadcn/ui DropdownMenu on the native popover API: the
// top layer replaces Radix's Portal, light dismiss and Esc are browser-
// native. Trigger and content are wired by proximity — DropdownMenuTrigger
// opens the popover inside the same DropdownMenu root, no ids. JS adds
// fixed-position anchoring to the trigger rect (CSS anchor positioning is
// not yet Baseline — see docs/jsx-parity.md), state/aria sync, arrow-key
// roving focus, and close-on-select. Requires the dropdown behavior module
// (ui/dropdown/dropdown.js).
component DropdownMenu(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="dropdown-menu" data-gsxui-dropdown class="contents" { attrs... }>{ children }</div>
}

component DropdownMenuTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button data-slot="dropdown-menu-trigger" data-gsxui-dropdown-trigger type="button" aria-haspopup="menu" aria-expanded="false" { attrs... }>{ children }</button>
}

// DropdownMenuContent renders the popover. popover="auto" gives top layer,
// light dismiss, and free Esc; data-state is server-rendered "closed" and
// kept in sync by dropdown.js on the toggle event.
component DropdownMenuContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="dropdown-menu-content"
		data-gsxui-dropdown-content
		popover="auto"
		role="menu"
		tabindex="-1"
		data-state="closed"
		class="z-50 max-h-96 min-w-[8rem] origin-top-left overflow-x-hidden overflow-y-auto rounded-md border bg-popover p-1 text-popover-foreground shadow-md data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95"
		{ attrs... }
	>{ children }</div>
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
		data-variant={ variant |> default("default") }
		role="menuitem"
		tabindex="-1"
		class="relative flex cursor-default items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 data-[variant=destructive]:text-destructive data-[variant=destructive]:focus:bg-destructive/10 data-[variant=destructive]:focus:text-destructive dark:data-[variant=destructive]:focus:bg-destructive/20 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground data-[variant=destructive]:*:[svg]:text-destructive!"
		{ attrs... }
	>{ children }</div>
}

// DropdownMenuLabel's inset prop is dropped along with DropdownMenuItem's
// (see docs/jsx-parity.md) — the data-[inset]:pl-8 selector is removed.
component DropdownMenuLabel(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="dropdown-menu-label" class="px-2 py-1.5 text-sm font-medium" { attrs... }>{ children }</div>
}

component DropdownMenuSeparator(attrs gsx.Attrs) {
	<div data-slot="dropdown-menu-separator" role="separator" class="-mx-1 my-1 h-px bg-border" { attrs... }></div>
}

component DropdownMenuShortcut(children gsx.Node, attrs gsx.Attrs) {
	<span data-slot="dropdown-menu-shortcut" class="ml-auto text-xs tracking-widest text-muted-foreground" { attrs... }>{ children }</span>
}

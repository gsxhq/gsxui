package ui

import "github.com/gsxhq/gsx"

// Popover is the shadcn/ui Popover on the native popover API: the top layer
// replaces Radix's Portal, light dismiss and Esc are browser-native.
// Trigger and content are wired by proximity — PopoverTrigger opens the
// popover inside the same Popover root, no ids. This is dropdown.gsx's
// mechanism with the menu semantics stripped: no role="menu", no arrow-key
// roving focus, no close-on-select — a popover holds arbitrary content (a
// form, free text), not a list of selectable items. JS adds fixed-position
// anchoring CENTERED below the trigger rect (Radix's own Popover default is
// side=bottom align=center, unlike DropdownMenuContent's align=start) and
// state/aria sync. Requires the popover behavior module (ui/popover.js).
component Popover(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="popover" data-gsxui-popover class="contents" { attrs... }>{ children }</div>
}

component PopoverTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		data-slot="popover-trigger"
		data-gsxui-popover-trigger
		type="button"
		aria-expanded="false"
		{ attrs... }
	>
		{ children }
	</button>
}

// PopoverContent renders the popover. popover="auto" gives top layer, light
// dismiss, and free Esc; data-state is server-rendered "closed" and kept in
// sync by popover.js on the toggle event. data-side="bottom" is
// server-rendered statically — popover.js always anchors below the
// trigger, so shadcn's data-[side=bottom]:slide-in-from-top-2 enter slide
// applies without Radix's runtime side tracking (same ADAPT as dropdown/
// tooltip). origin-top replaces shadcn's Radix runtime transform-origin var
// (--radix-popover-content-transform-origin) — the content is always
// centered below the trigger, so its scale/fade animation always
// originates from top-center (same substitution shape as dropdown's
// origin-top-left, adjusted for centered rather than start alignment).
component PopoverContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="popover-content"
		data-gsxui-popover-content
		popover="auto"
		data-state="closed"
		data-side="bottom"
		tabindex="-1"
		class="z-50 w-72 origin-top rounded-md border bg-popover p-4 text-popover-foreground shadow-md outline-hidden data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95"
		{ attrs... }
	>
		{ children }
	</div>
}

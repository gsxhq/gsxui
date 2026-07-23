package ui

import "github.com/gsxhq/gsx"

// Tooltip is the shadcn/ui Tooltip on the native popover API: popover="manual"
// puts the content in the top layer without light dismiss (hover/focus drive
// it, not outside clicks or Esc). Trigger and content are wired by proximity
// — TooltipTrigger shows the popover inside the same Tooltip root, no ids.
// JS adds fixed-position anchoring above the trigger rect (CSS anchor
// positioning is not yet Baseline — see docs/jsx-parity.md), a 300ms open
// delay, and state/event sync. Radix's TooltipProvider delay-group
// machinery is not ported (see docs/jsx-parity.md); the Arrow ports as a
// static child span in TooltipContent (the tooltip is always anchored
// above the trigger, so the diamond always sits bottom-center — no Radix
// side-tracking slot needed). Requires the tooltip behavior module
// (ui/tooltip.js).
component Tooltip(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="tooltip" data-gsxui-tooltip class="contents" { attrs... }>{ children }</div>
}

component TooltipTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button data-slot="tooltip-trigger" data-gsxui-tooltip-trigger type="button" { attrs... }>{ children }</button>
}

// TooltipContent renders the popover. popover="manual" is load-bearing:
// "auto" popovers light-dismiss on outside pointerdown, which would race
// tooltip.js's own pointerout/focusout hide logic. data-state is
// server-rendered "closed" and kept in sync by tooltip.js. data-side="top"
// is server-rendered statically — tooltip.js always anchors above the
// trigger, so shadcn's data-[side=top]:slide-in-from-bottom-2 enter slide
// applies without Radix's runtime side tracking. overflow-visible is a
// popover-port ADAPT (same family as inset:auto): the UA styles popovers
// overflow:auto, which would clip the protruding arrow span and grow a
// scrollbar instead of showing the diamond.
component TooltipContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="tooltip-content"
		data-gsxui-tooltip-content
		popover="manual"
		role="tooltip"
		data-state="closed"
		data-side="top"
		class="z-50 w-fit origin-bottom animate-in rounded-md bg-foreground px-3 py-1.5 text-xs text-balance text-background overflow-visible fade-in-0 zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95"
		{ attrs... }
	>
		{ children }
		<span data-slot="tooltip-arrow" class="absolute top-full left-1/2 z-50 size-2.5 -translate-x-1/2 -translate-y-[calc(50%+2px)] rotate-45 rounded-[2px] bg-foreground"></span>
	</div>
}

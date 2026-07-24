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
// scrollbar instead of showing the diamond. has-data-[slot=kbd]:pr-1.5 and
// **:data-[slot=kbd]:rounded-sm (nova density retarget, 2026-07-24) extend
// the existing tooltip-nested Kbd integration (see kbd.gsx's own
// [[data-slot=tooltip-content]_&] color tokens) with metric-only padding/
// radius adjustments for a Kbd child — inert unless a ui.Kbd is nested
// inside TooltipContent.
component TooltipContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="tooltip-content"
		data-gsxui-tooltip-content
		popover="manual"
		role="tooltip"
		data-state="closed"
		data-side="top"
		class={
			"z-50 w-fit origin-bottom gap-1.5 rounded-md bg-foreground px-3 py-1.5 text-xs has-data-[slot=kbd]:pr-1.5 **:data-[slot=kbd]:rounded-sm text-balance text-background overflow-visible",
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
		<span
			data-slot="tooltip-arrow"
			class="absolute top-full left-1/2 z-50 size-2.5 -translate-x-1/2 -translate-y-[calc(50%+2px)] rotate-45 rounded-[2px] bg-foreground"
		></span>
	</div>
}

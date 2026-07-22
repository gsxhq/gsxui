package tooltip

import "github.com/gsxhq/gsx"

// Tooltip is the shadcn/ui Tooltip on the native popover API: popover="manual"
// puts the content in the top layer without light dismiss (hover/focus drive
// it, not outside clicks or Esc). Trigger and content are wired by proximity
// — TooltipTrigger shows the popover inside the same Tooltip root, no ids.
// JS adds fixed-position anchoring above the trigger rect (CSS anchor
// positioning is not yet Baseline — see docs/jsx-parity.md), a 300ms open
// delay, and state/event sync. Radix's TooltipProvider delay-group
// machinery and the Arrow part are not ported (see docs/jsx-parity.md).
// Requires the tooltip behavior module (ui/tooltip/tooltip.js).
component Tooltip(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="tooltip" data-gsxui-tooltip class="contents" { attrs... }>{ children }</div>
}

component TooltipTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button data-slot="tooltip-trigger" data-gsxui-tooltip-trigger type="button" { attrs... }>{ children }</button>
}

// TooltipContent renders the popover. popover="manual" is load-bearing:
// "auto" popovers light-dismiss on outside pointerdown, which would race
// tooltip.js's own pointerout/focusout hide logic. data-state is
// server-rendered "closed" and kept in sync by tooltip.js.
component TooltipContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="tooltip-content"
		data-gsxui-tooltip-content
		popover="manual"
		role="tooltip"
		data-state="closed"
		class="z-50 w-fit origin-(--radix-tooltip-content-transform-origin) animate-in rounded-md bg-foreground px-3 py-1.5 text-xs text-balance text-background fade-in-0 zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95"
		{ attrs... }
	>{ children }</div>
}

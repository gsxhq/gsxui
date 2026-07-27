package ui

import "github.com/gsxhq/gsx"

// Tooltip uses a manual native popover so hover/focus behavior controls its
// top-layer lifetime. Its arrow is static because placement is always top.
component Tooltip(children gsx.Node, attrs gsx.Attrs) {
	<div data-gsxui-tooltip { withSlot("tooltip", attrs)... }>{ children }</div>
}

component TooltipTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button data-gsxui-tooltip-trigger type="button" { withSlot("tooltip-trigger", attrs)... }>{ children }</button>
}

component TooltipContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-gsxui-tooltip-content
		popover="manual"
		role="tooltip"
		data-state="closed"
		data-side="top"
		{ withSlot("tooltip-content", attrs)... }
	>
		{ children }
		<span { withSlot("tooltip-arrow", nil)... }></span>
	</div>
}

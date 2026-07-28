package ui

import "github.com/gsxhq/gsx"

// Tooltip uses a manual native popover so hover/focus behavior controls its
// top-layer lifetime. Its arrow is static because placement is always top.
component Tooltip(children gsx.Node, attrs gsx.Attrs) {
	<div data-gsxui-tooltip { attrs... } data-gsxui-slot-tooltip>{ children }</div>
}

component TooltipTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button data-gsxui-tooltip-trigger type="button" { attrs... } data-gsxui-slot-tooltip-trigger>{ children }</button>
}

component TooltipContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-gsxui-tooltip-content
		popover="manual"
		role="tooltip"
		data-state="closed"
		data-side="top"
		{ attrs... }
		data-gsxui-slot-tooltip-content
	>
		{ children }
		<span data-gsxui-slot-tooltip-arrow></span>
	</div>
}

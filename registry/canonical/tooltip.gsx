package canonical

import "github.com/gsxhq/gsx"

// Tooltip uses a manual native popover so hover/focus behavior controls its
// top-layer lifetime. Its arrow is static because placement is always top.
component Tooltip(children gsx.Node, attrs gsx.Attrs) {
	<div class={ tooltip.Root() } { attrs... } data-gsxui-slot-tooltip>{ children }</div>
}

component TooltipTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button type="button" { attrs... } data-gsxui-slot-tooltip-trigger>{ children }</button>
}

component TooltipContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		popover="manual"
		role="tooltip"
		data-state="closed"
		data-side="top"
		class={ tooltip.Content() }
		{ attrs... }
		data-gsxui-slot-tooltip-content
	>
		{ children }
		<span class={ tooltip.Arrow() } data-gsxui-slot-tooltip-arrow></span>
	</div>
}

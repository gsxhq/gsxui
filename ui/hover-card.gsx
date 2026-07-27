package ui

import "github.com/gsxhq/gsx"

// HoverCard uses a manual native popover so pointer/focus behavior, rather
// than light dismissal, controls its top-layer lifetime.
component HoverCard(children gsx.Node, attrs gsx.Attrs) {
	<div data-gsxui-hovercard { withSlot("hover-card", attrs)... }>{ children }</div>
}

component HoverCardTrigger(children gsx.Node, attrs gsx.Attrs) {
	<span data-gsxui-hovercard-trigger { withSlot("hover-card-trigger", attrs)... }>{ children }</span>
}

component HoverCardContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-gsxui-hovercard-content
		popover="manual"
		data-state="closed"
		data-side="bottom"
		{ withSlot("hover-card-content", attrs)... }
	>
		{ children }
	</div>
}

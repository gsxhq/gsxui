package ui

import "github.com/gsxhq/gsx"

// HoverCard uses a manual native popover so pointer/focus behavior, rather
// than light dismissal, controls its top-layer lifetime.
component HoverCard(children gsx.Node, attrs gsx.Attrs) {
	<div data-gsxui-hovercard { attrs... } data-gsxui-slot-hover-card>{ children }</div>
}

component HoverCardTrigger(children gsx.Node, attrs gsx.Attrs) {
	<span data-gsxui-hovercard-trigger { attrs... } data-gsxui-slot-hover-card-trigger>{ children }</span>
}

component HoverCardContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-gsxui-hovercard-content
		popover="manual"
		data-state="closed"
		data-side="bottom"
		{ attrs... } data-gsxui-slot-hover-card-content
	>
		{ children }
	</div>
}

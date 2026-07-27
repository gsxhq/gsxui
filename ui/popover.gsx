package ui

import "github.com/gsxhq/gsx"

// Popover uses the native auto-popover top layer with light dismissal and
// proximity-scoped behavior in ui/popover.js.
component Popover(children gsx.Node, attrs gsx.Attrs) {
	<div data-gsxui-popover { attrs... } data-gsxui-slot-popover>{ children }</div>
}

component PopoverTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		data-gsxui-popover-trigger
		type="button"
		aria-expanded="false"
		{ attrs... } data-gsxui-slot-popover-trigger
	>
		{ children }
	</button>
}

component PopoverContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-gsxui-popover-content
		popover="auto"
		data-state="closed"
		data-side="bottom"
		tabindex="-1"
		{ attrs... } data-gsxui-slot-popover-content
	>
		{ children }
	</div>
}

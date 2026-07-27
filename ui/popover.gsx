package ui

import "github.com/gsxhq/gsx"

// Popover uses the native auto-popover top layer with light dismissal and
// proximity-scoped behavior in ui/popover.js.
component Popover(children gsx.Node, attrs gsx.Attrs) {
	<div data-gsxui-popover { withSlot("popover", attrs)... }>{ children }</div>
}

component PopoverTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		data-gsxui-popover-trigger
		type="button"
		aria-expanded="false"
		{ withSlot("popover-trigger", attrs)... }
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
		{ withSlot("popover-content", attrs)... }
	>
		{ children }
	</div>
}

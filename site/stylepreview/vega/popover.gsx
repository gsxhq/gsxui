package vega

import "github.com/gsxhq/gsx"

// Popover uses the native auto-popover top layer with light dismissal and
// proximity-scoped behavior in ui/popover.js.
component Popover(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "contents" } { attrs... } data-gsxui-slot-popover>{ children }</div>
}

component PopoverTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		type="button"
		aria-expanded="false"
		{ attrs... }
		data-gsxui-slot-popover-trigger
	>
		{ children }
	</button>
}

component PopoverContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		popover="auto"
		data-state="closed"
		data-side="bottom"
		tabindex="-1"
		class={
			"bg-popover text-popover-foreground data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 ring-foreground/10 [&:popover-open]:flex flex-col gap-4 rounded-md p-4 text-sm shadow-md ring-1 duration-100"
		}
		{ attrs... }
		data-gsxui-slot-popover-content
	>
		{ children }
	</div>
}

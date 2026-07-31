package maia

import "github.com/gsxhq/gsx"

// Popover uses the native auto-popover top layer with light dismissal and
// proximity-scoped behavior in ui/popover.js.
component Popover(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "contents" } { attrs... } data-gsxui-slot-popover>{ children }</div>
}

component PopoverTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		data-gsxui-popover-trigger
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
			"z-50 w-72 origin-top gap-2.5 rounded-lg border bg-popover p-2.5 text-sm text-popover-foreground shadow-md outline-hidden opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 [&:popover-open]:opacity-100 [&:popover-open]:scale-100 starting:[&:popover-open]:opacity-0 starting:[&:popover-open]:scale-95 starting:[&:popover-open]:data-[side=bottom]:-translate-y-2 starting:[&:popover-open]:data-[side=left]:translate-x-2 starting:[&:popover-open]:data-[side=right]:-translate-x-2 starting:[&:popover-open]:data-[side=top]:translate-y-2"
		}
		{ attrs... }
		data-gsxui-slot-popover-content
	>
		{ children }
	</div>
}

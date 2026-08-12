package mira

import "github.com/gsxhq/gsx"

// Tooltip uses a manual native popover so hover/focus behavior controls its
// top-layer lifetime. Its arrow is static because placement is always top.
component Tooltip(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "contents" } { attrs... } data-gsxui-slot-tooltip>{ children }</div>
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
		class={
			"data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95 data-[state=delayed-open]:animate-in data-[state=delayed-open]:fade-in-0 data-[state=delayed-open]:zoom-in-95 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 [&:popover-open]:inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-xs has-[[data-gsxui-slot-kbd]]:pr-1.5 [&_[data-gsxui-slot-kbd]]:relative [&_[data-gsxui-slot-kbd]]:isolate [&_[data-gsxui-slot-kbd]]:z-50 [&_[data-gsxui-slot-kbd]]:rounded-sm bg-foreground text-background w-fit max-w-xs origin-bottom"
		}
		{ attrs... }
		data-gsxui-slot-tooltip-content
	>
		{ children }
		<span class={ "size-2.5 translate-y-[calc(-50%-2px)] rotate-45 rounded-[2px] bg-foreground top-full left-1/2 -translate-x-1/2" } data-gsxui-slot-tooltip-arrow></span>
	</div>
}

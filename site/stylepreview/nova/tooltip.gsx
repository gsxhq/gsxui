package nova

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
			"z-50 w-fit origin-bottom gap-1.5 overflow-visible rounded-md bg-foreground px-3 py-1.5 text-xs text-balance text-background has-[[data-gsxui-slot-kbd]]:pe-1.5 opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 [&:popover-open]:opacity-100 [&:popover-open]:scale-100 starting:[&:popover-open]:opacity-0 starting:[&:popover-open]:scale-95 starting:[&:popover-open]:data-[side=bottom]:-translate-y-2 starting:[&:popover-open]:data-[side=left]:translate-x-2 starting:[&:popover-open]:data-[side=right]:-translate-x-2 starting:[&:popover-open]:data-[side=top]:translate-y-2"
		}
		{ attrs... }
		data-gsxui-slot-tooltip-content
	>
		{ children }
		<span class={ "absolute top-full left-1/2 z-50 size-2.5 -translate-x-1/2 -translate-y-[calc(50%+2px)] rotate-45 rounded-[2px] bg-foreground" } data-gsxui-slot-tooltip-arrow></span>
	</div>
}

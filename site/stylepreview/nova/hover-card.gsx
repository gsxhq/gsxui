package nova

import "github.com/gsxhq/gsx"

// HoverCard uses a manual native popover so pointer/focus behavior, rather
// than light dismissal, controls its top-layer lifetime.
component HoverCard(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "contents" } { attrs... } data-gsxui-slot-hover-card>{ children }</div>
}

component HoverCardTrigger(children gsx.Node, attrs gsx.Attrs) {
	<span { attrs... } data-gsxui-slot-hover-card-trigger>{ children }</span>
}

component HoverCardContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		popover="manual"
		data-state="closed"
		data-side="bottom"
		class={
			"z-50 w-64 origin-top rounded-lg border bg-popover p-2.5 text-sm text-popover-foreground shadow-md outline-hidden opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 [&:popover-open]:opacity-100 [&:popover-open]:scale-100 starting:[&:popover-open]:opacity-0 starting:[&:popover-open]:scale-95 starting:[&:popover-open]:data-[side=bottom]:-translate-y-2 starting:[&:popover-open]:data-[side=left]:translate-x-2 starting:[&:popover-open]:data-[side=right]:-translate-x-2 starting:[&:popover-open]:data-[side=top]:translate-y-2"
		}
		{ attrs... }
		data-gsxui-slot-hover-card-content
	>
		{ children }
	</div>
}

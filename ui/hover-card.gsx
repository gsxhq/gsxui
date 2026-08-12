package ui

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
			"transition-none data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 ring-foreground/10 bg-popover text-popover-foreground w-64 rounded-lg p-2.5 text-sm shadow-md ring-1 duration-100 origin-top outline-hidden"
		}
		{ attrs... }
		data-gsxui-slot-hover-card-content
	>
		{ children }
	</div>
}

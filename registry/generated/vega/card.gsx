package ui

import "github.com/gsxhq/gsx"

// Card and its parts are the shadcn/ui Card compound set. Parts are plain
// sibling components — compose them in markup; no shared state, no context.

component Card(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"group/card",
			"ring-foreground/10 bg-card text-card-foreground gap-(--card-spacing) overflow-hidden rounded-xl py-(--card-spacing) text-sm shadow-xs ring-1 [--card-spacing:--spacing(6)] has-[>img:first-child]:pt-0 data-[size=sm]:[--card-spacing:--spacing(4)] *:[img:first-child]:rounded-t-xl *:[img:last-child]:rounded-b-xl flex flex-col"
		}
		{ attrs... }
		data-gsxui-slot-card
	>
		{ children }
	</div>
}

component CardHeader(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={ "gap-1 rounded-t-xl px-(--card-spacing) grid auto-rows-min grid-rows-[auto_auto]" }
		{ attrs... }
		data-gsxui-slot-card-header
	>
		{ children }
	</div>
}

component CardTitle(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={ "text-base leading-normal font-medium group-data-[size=sm]/card:text-sm" }
		{ attrs... }
		data-gsxui-slot-card-title
	>
		{ children }
	</div>
}

component CardDescription(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "text-muted-foreground text-sm" } { attrs... } data-gsxui-slot-card-description>{ children }</div>
}

component CardAction(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={ "col-start-2 row-span-2 row-start-1 self-start justify-self-end" }
		{ attrs... }
		data-gsxui-slot-card-action
	>
		{ children }
	</div>
}

component CardContent(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "px-(--card-spacing)" } { attrs... } data-gsxui-slot-card-content>{ children }</div>
}

component CardFooter(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "rounded-b-xl px-(--card-spacing) flex" } { attrs... } data-gsxui-slot-card-footer>{ children }</div>
}

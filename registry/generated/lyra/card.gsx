package ui

import "github.com/gsxhq/gsx"

// Card and its parts are the shadcn/ui Card compound set. Parts are plain
// sibling components — compose them in markup; no shared state, no context.

component Card(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"group/card",
			"ring-foreground/10 bg-card text-card-foreground gap-(--card-spacing) overflow-hidden rounded-none py-(--card-spacing) text-xs/relaxed ring-1 [--card-spacing:--spacing(4)] has-[[data-gsxui-slot-card-footer]]:pb-0 has-[>img:first-child]:pt-0 data-[size=sm]:[--card-spacing:--spacing(3)] data-[size=sm]:has-[[data-gsxui-slot-card-footer]]:pb-0 *:[img:first-child]:rounded-none *:[img:last-child]:rounded-none flex flex-col"
		}
		{ attrs... }
		data-gsxui-slot-card
	>
		{ children }
	</div>
}

component CardHeader(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"gap-1 rounded-none px-(--card-spacing) grid auto-rows-min items-start has-[[data-gsxui-slot-card-action]]:grid-cols-[1fr_auto] has-[[data-gsxui-slot-card-description]]:grid-rows-[auto_auto]"
		}
		{ attrs... }
		data-gsxui-slot-card-header
	>
		{ children }
	</div>
}

component CardTitle(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "text-sm font-medium group-data-[size=sm]/card:text-sm" } { attrs... } data-gsxui-slot-card-title>
		{ children }
	</div>
}

component CardDescription(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "text-muted-foreground text-xs/relaxed" } { attrs... } data-gsxui-slot-card-description>
		{ children }
	</div>
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
	<div class={ "rounded-none border-t p-(--card-spacing) flex items-center" } { attrs... } data-gsxui-slot-card-footer>
		{ children }
	</div>
}

package sera

import "github.com/gsxhq/gsx"

// Card and its parts are the shadcn/ui Card compound set. Parts are plain
// sibling components — compose them in markup; no shared state, no context.

component Card(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"group/card",
			"bg-card text-card-foreground ring-foreground/5 gap-(--card-spacing) overflow-hidden py-(--card-spacing) text-sm shadow-sm ring-1 [--card-spacing:--spacing(8)] has-[>img:first-child]:pt-0 data-[size=sm]:[--card-spacing:--spacing(5)] *:[img:first-child]:rounded-none *:[img:last-child]:rounded-none flex flex-col"
		}
		{ attrs... }
		data-gsxui-slot-card
	>
		{ children }
	</div>
}

component CardHeader(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={ "gap-1.5 rounded-none px-(--card-spacing) grid auto-rows-min grid-rows-[auto_auto]" }
		{ attrs... }
		data-gsxui-slot-card-header
	>
		{ children }
	</div>
}

component CardTitle(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "text-lg font-semibold tracking-wider uppercase" } { attrs... } data-gsxui-slot-card-title>
		{ children }
	</div>
}

component CardDescription(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "text-muted-foreground text-sm leading-relaxed" } { attrs... } data-gsxui-slot-card-description>
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
	<div class={ "px-(--card-spacing) flex" } { attrs... } data-gsxui-slot-card-footer>{ children }</div>
}

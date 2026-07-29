package nova

import "github.com/gsxhq/gsx"

// Card and its parts are the shadcn/ui Card compound set. Parts are plain
// sibling components — compose them in markup; no shared state, no context.

component Card(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"flex flex-col gap-4 rounded-xl border bg-card py-4 text-sm text-card-foreground has-[[data-gsxui-slot-card-footer]]:pb-0"
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
			"grid auto-rows-min grid-rows-[auto_auto] items-start gap-1 px-4 has-[[data-gsxui-slot-card-action]]:grid-cols-[1fr_auto] [&.border-b]:pb-4"
		}
		{ attrs... }
		data-gsxui-slot-card-header
	>
		{ children }
	</div>
}

component CardTitle(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "text-base leading-snug font-medium" } { attrs... } data-gsxui-slot-card-title>{ children }</div>
}

component CardDescription(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "text-sm text-muted-foreground" } { attrs... } data-gsxui-slot-card-description>{ children }</div>
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
	<div class={ "px-4" } { attrs... } data-gsxui-slot-card-content>{ children }</div>
}

component CardFooter(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "flex items-center rounded-b-xl border-t p-4" } { attrs... } data-gsxui-slot-card-footer>
		{ children }
	</div>
}

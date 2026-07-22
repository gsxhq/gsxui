package card

import "github.com/gsxhq/gsx"

// Card and its parts are the shadcn/ui Card compound set. Parts are plain
// sibling components — compose them in markup; no shared state, no context.

component Card(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="card" class="flex flex-col gap-6 rounded-xl border bg-card py-6 text-card-foreground shadow-sm" { attrs... }>{ children }</div>
}

component CardHeader(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="card-header" class="@container/card-header grid auto-rows-min grid-rows-[auto_auto] items-start gap-2 px-6 has-data-[slot=card-action]:grid-cols-[1fr_auto] [.border-b]:pb-6" { attrs... }>{ children }</div>
}

component CardTitle(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="card-title" class="leading-none font-semibold" { attrs... }>{ children }</div>
}

component CardDescription(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="card-description" class="text-sm text-muted-foreground" { attrs... }>{ children }</div>
}

component CardAction(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="card-action" class="col-start-2 row-span-2 row-start-1 self-start justify-self-end" { attrs... }>{ children }</div>
}

component CardContent(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="card-content" class="px-6" { attrs... }>{ children }</div>
}

component CardFooter(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="card-footer" class="flex items-center px-6 [.border-t]:pt-6" { attrs... }>{ children }</div>
}

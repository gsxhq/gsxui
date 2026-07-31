package canonical

import "github.com/gsxhq/gsx"

// Card and its parts are the shadcn/ui Card compound set. Parts are plain
// sibling components — compose them in markup; no shared state, no context.

component Card(children gsx.Node, attrs gsx.Attrs) {
	<div class={ card.Root() } { attrs... } data-gsxui-slot-card>
		{ children }
	</div>
}

component CardHeader(children gsx.Node, attrs gsx.Attrs) {
	<div class={ card.Header() } { attrs... } data-gsxui-slot-card-header>
		{ children }
	</div>
}

component CardTitle(children gsx.Node, attrs gsx.Attrs) {
	<div class={ card.Title() } { attrs... } data-gsxui-slot-card-title>{ children }</div>
}

component CardDescription(children gsx.Node, attrs gsx.Attrs) {
	<div class={ card.Description() } { attrs... } data-gsxui-slot-card-description>{ children }</div>
}

component CardAction(children gsx.Node, attrs gsx.Attrs) {
	<div class={ card.Action() } { attrs... } data-gsxui-slot-card-action>
		{ children }
	</div>
}

component CardContent(children gsx.Node, attrs gsx.Attrs) {
	<div class={ card.Content() } { attrs... } data-gsxui-slot-card-content>{ children }</div>
}

component CardFooter(children gsx.Node, attrs gsx.Attrs) {
	<div class={ card.Footer() } { attrs... } data-gsxui-slot-card-footer>{ children }</div>
}

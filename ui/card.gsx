package ui

import "github.com/gsxhq/gsx"

// Card and its parts are the shadcn/ui Card compound set. Parts are plain
// sibling components — compose them in markup; no shared state, no context.

component Card(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("card", attrs)... }>
		{ children }
	</div>
}

component CardHeader(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("card-header", attrs)... }>
		{ children }
	</div>
}

component CardTitle(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("card-title", attrs)... }>{ children }</div>
}

component CardDescription(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("card-description", attrs)... }>{ children }</div>
}

component CardAction(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("card-action", attrs)... }>
		{ children }
	</div>
}

component CardContent(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("card-content", attrs)... }>{ children }</div>
}

component CardFooter(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("card-footer", attrs)... }>{ children }</div>
}

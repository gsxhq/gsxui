package ui

import "github.com/gsxhq/gsx"

// Alert and its parts are the shadcn/ui Alert. variant: "" (default) |
// "destructive". Extra attributes fall through to the <div>; caller classes
// merge (caller wins per property).
component Alert(variant string, children gsx.Node, attrs gsx.Attrs) {
	<div
		role="alert"
		data-variant={variant |> default("default")}
		{ withSlot("alert", attrs)... }
	>
		{ children }
	</div>
}

component AlertTitle(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("alert-title", attrs)... }>
		{ children }
	</div>
}

component AlertDescription(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("alert-description", attrs)... }>
		{ children }
	</div>
}

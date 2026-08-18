package canonical

import "github.com/gsxhq/gsx"

// Alert and its parts are the shadcn/ui Alert. variant: "" (default) |
// "destructive". Extra attributes fall through to the <div>; caller classes
// merge (caller wins per property).
component Alert(variant string, children gsx.Node, attrs gsx.Attrs) {
	<div
		role="alert"
		data-variant={variant |> default("default")}
		class={ "group/alert", alert.Root(), alert.Variant(variant) }
		{ attrs... }
		data-gsxui-slot-alert
	>
		{ children }
	</div>
}

component AlertTitle(children gsx.Node, attrs gsx.Attrs) {
	<div class={ alert.Title() } { attrs... } data-gsxui-slot-alert-title>
		{ children }
	</div>
}

component AlertDescription(children gsx.Node, attrs gsx.Attrs) {
	<div class={ alert.Description() } { attrs... } data-gsxui-slot-alert-description>
		{ children }
	</div>
}

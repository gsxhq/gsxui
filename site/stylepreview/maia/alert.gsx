package maia

import "github.com/gsxhq/gsx"

// Alert and its parts are the shadcn/ui Alert. variant: "" (default) |
// "destructive". Extra attributes fall through to the <div>; caller classes
// merge (caller wins per property).
component Alert(variant string, children gsx.Node, attrs gsx.Attrs) {
	<div
		role="alert"
		data-variant={variant |> default("default")}
		class={
			"relative grid w-full items-start gap-y-0.5 rounded-lg border bg-card px-2.5 py-2 text-sm text-card-foreground has-[>svg]:grid-cols-[auto_1fr] has-[>svg]:gap-x-2 [&>svg]:row-span-2 [&>svg]:translate-y-0.5 [&>svg]:text-current [&>svg:not([class*='size-'])]:size-4",
			switch variant {
			case "destructive":
				"bg-card text-destructive [&_[data-gsxui-slot-alert-description]]:text-destructive/90"
			default:
				"bg-card text-card-foreground"
			}
		}
		{ attrs... }
		data-gsxui-slot-alert
	>
		{ children }
	</div>
}

component AlertTitle(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "col-start-2 font-medium" } { attrs... } data-gsxui-slot-alert-title>
		{ children }
	</div>
}

component AlertDescription(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"col-start-2 grid justify-items-start text-sm text-muted-foreground [&_p]:leading-relaxed [&_p:not(:last-child)]:mb-4"
		}
		{ attrs... }
		data-gsxui-slot-alert-description
	>
		{ children }
	</div>
}

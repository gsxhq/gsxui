package ui

import "github.com/gsxhq/gsx"

// Alert and its parts are the shadcn/ui Alert. variant: "" (default) |
// "destructive". Extra attributes fall through to the <div>; caller classes
// merge (caller wins per property).
component Alert(variant string, children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="alert"
		role="alert"
		class={
			"relative grid w-full grid-cols-[0_1fr] items-start gap-y-0.5 rounded-lg border px-4 py-3 text-sm has-[>svg]:grid-cols-[calc(var(--spacing)*4)_1fr] has-[>svg]:gap-x-3 [&>svg]:size-4 [&>svg]:translate-y-0.5 [&>svg]:text-current",
			switch variant {
			case "destructive":
				"bg-card text-destructive *:data-[slot=alert-description]:text-destructive/90 [&>svg]:text-current"
			default:
				"bg-card text-card-foreground"
			}
		}
		{ attrs... }
	>
		{ children }
	</div>
}

component AlertTitle(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="alert-title" class="col-start-2 line-clamp-1 min-h-4 font-medium tracking-tight" { attrs... }>
		{ children }
	</div>
}

component AlertDescription(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="alert-description"
		class="col-start-2 grid justify-items-start gap-1 text-sm text-muted-foreground [&_p]:leading-relaxed"
		{ attrs... }
	>
		{ children }
	</div>
}

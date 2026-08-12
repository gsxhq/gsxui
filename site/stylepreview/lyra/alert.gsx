package lyra

import "github.com/gsxhq/gsx"

// Alert and its parts are the shadcn/ui Alert. variant: "" (default) |
// "destructive". Extra attributes fall through to the <div>; caller classes
// merge (caller wins per property).
component Alert(variant string, children gsx.Node, attrs gsx.Attrs) {
	<div
		role="alert"
		data-variant={variant |> default("default")}
		class={
			"group/alert",
			"grid gap-0.5 rounded-none border px-2.5 py-2 text-left text-xs has-[>svg]:grid-cols-[auto_1fr] has-[>svg]:gap-x-2 *:[svg]:row-span-2 *:[svg]:translate-y-0 *:[svg]:text-current *:[svg:not([class*='size-'])]:size-4 w-full",
			switch variant {
			case "destructive":
				"text-destructive bg-card [&>[data-gsxui-slot-alert-description]]:text-destructive/90 *:[svg]:text-current"
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
	<div
		class={
			"font-medium group-has-[>svg]/alert:col-start-2 [&_a]:underline [&_a]:underline-offset-3 [&_a]:hover:text-foreground"
		}
		{ attrs... }
		data-gsxui-slot-alert-title
	>
		{ children }
	</div>
}

component AlertDescription(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"text-muted-foreground text-xs/relaxed text-balance md:text-pretty [&_p:not(:last-child)]:mb-2 [&_a]:underline [&_a]:underline-offset-3 [&_a]:hover:text-foreground"
		}
		{ attrs... }
		data-gsxui-slot-alert-description
	>
		{ children }
	</div>
}

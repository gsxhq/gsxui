package ui

import "github.com/gsxhq/gsx"

// Badge is the shadcn/ui Badge. variant: "" (default) | "secondary" |
// "destructive" | "outline" | "ghost" | "link". Extra attributes fall
// through to the <span>; caller classes merge (caller wins per property).
//
// Retargeted to nova density (2026-07-24 nova density map, `## badge`).
// ADAPT: nova keys directional icon padding off `data-icon="inline-start|
// inline-end"` stamps this component doesn't emit; ported instead onto
// gsxui's existing has-[>svg]:px-* selector mechanism (the same
// substitution button.gsx/toggle.gsx make — see their own doc comments),
// collapsing nova's matching inline-start/inline-end value (both px-1.5)
// into one has-[>svg]:px-1.5.
component Badge(variant string, children gsx.Node, attrs gsx.Attrs) {
	<span
		data-variant={variant |> default("default")}
		class={
			"h-5 gap-1 rounded-4xl border border-transparent px-2 py-0.5 text-xs font-medium transition-all has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 [&>svg]:size-3!",
			switch variant {
			case "secondary":
				"bg-secondary text-secondary-foreground [a]:hover:bg-secondary/80"
			case "destructive":
				"bg-destructive/10 [a]:hover:bg-destructive/20 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40 text-destructive dark:bg-destructive/20"
			case "outline":
				"border-border text-foreground [a]:hover:bg-muted [a]:hover:text-muted-foreground bg-input/30"
			case "ghost":
				"hover:bg-muted hover:text-muted-foreground dark:hover:bg-muted/50"
			case "link":
				"text-primary underline-offset-4 hover:underline"
			default:
				"bg-primary text-primary-foreground [a]:hover:bg-primary/80"
			}
		}
		{ attrs... }
		data-gsxui-slot-badge
	>
		{ children }
	</span>
}

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
			"aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 focus-visible:ring-3 gap-1.5 rounded-none border-0 bg-transparent px-0 py-0 text-[0.625rem] font-semibold uppercase tracking-widest transition-colors has-[>svg]:px-0 [&>svg]:size-3 inline-flex items-center justify-center overflow-hidden whitespace-nowrap w-fit shrink-0 focus-visible:border-ring focus-visible:ring-ring/50 [&>svg]:pointer-events-none",
			switch variant {
			case "secondary":
				"text-muted-foreground [a]:hover:text-foreground"
			case "destructive":
				"text-destructive [a]:hover:text-destructive/70 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40"
			case "outline":
				"text-foreground [a]:hover:text-foreground/70"
			case "ghost":
				"text-muted-foreground hover:text-foreground"
			case "link":
				"text-foreground underline-offset-4 hover:underline"
			default:
				"text-foreground [a]:hover:text-foreground/70"
			}
		}
		{ attrs... }
		data-gsxui-slot-badge
	>
		{ children }
	</span>
}

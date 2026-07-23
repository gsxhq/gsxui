package ui

import "github.com/gsxhq/gsx"

// Badge is the shadcn/ui Badge. variant: "" (default) | "secondary" |
// "destructive" | "outline" | "ghost" | "link". Extra attributes fall
// through to the <span>; caller classes merge (caller wins per property).
component Badge(variant string, children gsx.Node, attrs gsx.Attrs) {
	<span
		data-slot="badge"
		data-variant={variant |> default("default")}
		class={
			"inline-flex w-fit shrink-0 items-center justify-center gap-1 overflow-hidden rounded-full border border-transparent px-2 py-0.5 text-xs font-medium whitespace-nowrap transition-[color,box-shadow] focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&>svg]:pointer-events-none [&>svg]:size-3",
			switch variant {
			case "secondary":
				"bg-secondary text-secondary-foreground [a&]:hover:bg-secondary/90"
			case "destructive":
				"bg-destructive text-white focus-visible:ring-destructive/20 dark:bg-destructive/60 dark:focus-visible:ring-destructive/40 [a&]:hover:bg-destructive/90"
			case "outline":
				"border-border text-foreground [a&]:hover:bg-accent [a&]:hover:text-accent-foreground"
			case "ghost":
				"[a&]:hover:bg-accent [a&]:hover:text-accent-foreground"
			case "link":
				"text-primary underline-offset-4 [a&]:hover:underline"
			default:
				"bg-primary text-primary-foreground [a&]:hover:bg-primary/90"
			}
		}
		{ attrs... }
	>
		{ children }
	</span>
}

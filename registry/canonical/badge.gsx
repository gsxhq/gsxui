package canonical

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
		class={ badge.Root(), badge.Variant(variant) }
		{ attrs... }
		data-gsxui-slot-badge
	>
		{ children }
	</span>
}

package ui

import "github.com/gsxhq/gsx"

// Empty and its parts are the shadcn/ui Empty (registry/new-york-v4/ui/
// empty.tsx) — no Radix primitive underneath either; every part is already a
// plain styled <div>, the same "package-namespaced compound parts" shape as
// card/breadcrumb.
component Empty(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "gap-4 rounded-xl border-dashed p-6 flex flex-col" } { attrs... } data-gsxui-slot-empty>
		{ children }
	</div>
}

component EmptyHeader(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "gap-2 flex flex-col" } { attrs... } data-gsxui-slot-empty-header>
		{ children }
	</div>
}

// EmptyMedia's variant cva map (registry's emptyMediaVariants) picks between
// two entirely static presentation blocks by the JS-resolved variant value.
// The CSS-only contract reflects that value through data-variant.
// The stable token is "empty-icon", matching shadcn's own semantic name.
component EmptyMedia(variant string, children gsx.Node, attrs gsx.Attrs) {
	<div
		data-variant={variant |> default("default")}
		class={
			"mb-2 flex",
			switch variant {
			case "icon":
				"bg-muted text-foreground flex size-8 shrink-0 items-center justify-center rounded-lg [&_svg:not([class*='size-'])]:size-4"
			default:
				"bg-transparent"
			}
		}
		{ attrs... }
		data-gsxui-slot-empty-icon
	>
		{ children }
	</div>
}

component EmptyTitle(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "text-sm font-medium tracking-tight" } { attrs... } data-gsxui-slot-empty-title>
		{ children }
	</div>
}

// EmptyDescription renders a <div>, matching shadcn's own actual element —
// its TypeScript prop type reads React.ComponentProps<"p"> but the JSX it
// returns is a <div>, the same shipped-type/element mismatch already noted
// for Kbd/KbdGroup (see docs/jsx-parity.md ## kbd); ported verbatim, tag
// included, per the token-for-token rule.
component EmptyDescription(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "text-sm/relaxed" } { attrs... } data-gsxui-slot-empty-description>
		{ children }
	</div>
}

component EmptyContent(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "gap-2.5 text-sm flex flex-col" } { attrs... } data-gsxui-slot-empty-content>
		{ children }
	</div>
}

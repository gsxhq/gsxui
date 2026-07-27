package ui

import "github.com/gsxhq/gsx"

// Empty and its parts are the shadcn/ui Empty (registry/new-york-v4/ui/
// empty.tsx) — no Radix primitive underneath either; every part is already a
// plain styled <div>, the same "package-namespaced compound parts" shape as
// card/breadcrumb.
component Empty(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("empty", attrs)... }>
		{ children }
	</div>
}

component EmptyHeader(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("empty-header", attrs)... }>
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
		{ withSlot("empty-icon", attrs)... }
	>
		{ children }
	</div>
}

component EmptyTitle(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("empty-title", attrs)... }>
		{ children }
	</div>
}

// EmptyDescription renders a <div>, matching shadcn's own actual element —
// its TypeScript prop type reads React.ComponentProps<"p"> but the JSX it
// returns is a <div>, the same shipped-type/element mismatch already noted
// for Kbd/KbdGroup (see docs/jsx-parity.md ## kbd); ported verbatim, tag
// included, per the token-for-token rule.
component EmptyDescription(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("empty-description", attrs)... }>
		{ children }
	</div>
}

component EmptyContent(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("empty-content", attrs)... }>
		{ children }
	</div>
}

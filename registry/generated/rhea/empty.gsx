package ui

import "github.com/gsxhq/gsx"

// Empty and its parts are the shadcn/ui Empty (registry/new-york-v4/ui/
// empty.tsx) — no Radix primitive underneath either; every part is already a
// plain styled <div>, the same "package-namespaced compound parts" shape as
// card/breadcrumb.
component Empty(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"gap-4 rounded-3xl border-dashed p-12 flex flex-col w-full min-w-0 flex-1 items-center justify-center text-center text-balance"
		}
		{ attrs... }
		data-gsxui-slot-empty
	>
		{ children }
	</div>
}

component EmptyHeader(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "gap-2 flex flex-col max-w-sm items-center text-center" } { attrs... } data-gsxui-slot-empty-header>
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
			"mb-2 flex shrink-0 items-center justify-center [&_svg]:shrink-0 [&_svg]:pointer-events-none",
			switch variant {
			case "icon":
				"bg-muted text-foreground flex size-10 shrink-0 items-center justify-center rounded-xl [&_svg:not([class*='size-'])]:size-5"
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
	<div class={ "text-lg font-medium tracking-tight" } { attrs... } data-gsxui-slot-empty-title>
		{ children }
	</div>
}

// EmptyDescription renders a <div>, matching shadcn's own actual element —
// its TypeScript prop type reads React.ComponentProps<"p"> but the JSX it
// returns is a <div>, the same shipped-type/element mismatch already noted
// for Kbd/KbdGroup (see docs/jsx-parity.md ## kbd); ported verbatim, tag
// included, per the token-for-token rule.
component EmptyDescription(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={ "text-sm/relaxed text-muted-foreground [&>a]:underline [&>a]:underline-offset-4 [&>a:hover]:text-primary" }
		{ attrs... }
		data-gsxui-slot-empty-description
	>
		{ children }
	</div>
}

component EmptyContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={ "gap-4 text-sm flex flex-col w-full max-w-sm min-w-0 items-center text-balance" }
		{ attrs... }
		data-gsxui-slot-empty-content
	>
		{ children }
	</div>
}

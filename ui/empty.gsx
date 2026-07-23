package ui

import "github.com/gsxhq/gsx"

// Empty and its parts are the shadcn/ui Empty (registry/new-york-v4/ui/
// empty.tsx) — no Radix primitive underneath either; every part is already a
// plain styled <div>, the same "package-namespaced compound parts" shape as
// card/breadcrumb.
component Empty(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="empty"
		class="flex min-w-0 flex-1 flex-col items-center justify-center gap-6 rounded-lg border-dashed p-6 text-center text-balance md:p-12"
		{ attrs... }
	>
		{ children }
	</div>
}

component EmptyHeader(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="empty-header"
		class="flex max-w-sm flex-col items-center gap-2 text-center"
		{ attrs... }
	>
		{ children }
	</div>
}

// EmptyMedia's variant cva map (registry's emptyMediaVariants) picks between
// two entirely static class blocks by the JS-resolved variant value — no
// data-[variant=...] selectors to preserve, so this ports as a switch inside
// class={}, the same idiom as badge/button-group. data-slot is "empty-icon"
// in shadcn's own source, not "empty-media" — ported as-is (token-for-token),
// same call as ButtonGroupText's unmatched data-slot (see docs/jsx-parity.md).
component EmptyMedia(variant string, children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="empty-icon"
		data-variant={variant |> default("default")}
		class={
			"mb-2 flex shrink-0 items-center justify-center [&_svg]:pointer-events-none [&_svg]:shrink-0",
			switch variant {
			case "icon":
				"flex size-10 shrink-0 items-center justify-center rounded-lg bg-muted text-foreground [&_svg:not([class*='size-'])]:size-6"
			default:
				"bg-transparent"
			}
		}
		{ attrs... }
	>
		{ children }
	</div>
}

component EmptyTitle(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="empty-title" class="text-lg font-medium tracking-tight" { attrs... }>
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
		data-slot="empty-description"
		class="text-sm/relaxed text-muted-foreground [&>a]:underline [&>a]:underline-offset-4 [&>a:hover]:text-primary"
		{ attrs... }
	>
		{ children }
	</div>
}

component EmptyContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="empty-content"
		class="flex w-full max-w-sm min-w-0 flex-col items-center gap-4 text-sm text-balance"
		{ attrs... }
	>
		{ children }
	</div>
}

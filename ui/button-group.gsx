package ui

import "github.com/gsxhq/gsx"

// ButtonGroup and its parts are the shadcn/ui ButtonGroup
// (registry/new-york-v4/ui/button-group.tsx). shadcn's `buttonGroupVariants`
// cva map picks between two entirely static class blocks by the JS-resolved
// `orientation` value — there are no `data-[orientation=...]:` selectors in
// the source to preserve, so the port mirrors Badge's `switch` inside
// `class={}` (WIN, same idiom), not a CSS-side data-attribute selector.
// data-orientation is still stamped via the house `|> default` pattern (see
// button.gsx/dropdown.gsx) for consistency with every other data-variant
// stamp in this codebase — an ADAPT: shadcn leaves the attribute entirely
// unset when `orientation` is undefined (see docs/jsx-parity.md).
component ButtonGroup(orientation string, children gsx.Node, attrs gsx.Attrs) {
	<div
		role="group"
		data-slot="button-group"
		data-orientation={orientation |> default("horizontal")}
		class={
			"flex w-fit items-stretch has-[>[data-slot=button-group]]:gap-2 [&>*]:focus-visible:relative [&>*]:focus-visible:z-10 has-[select[aria-hidden=true]:last-child]:[&>[data-slot=select-trigger]:last-of-type]:rounded-r-md [&>[data-slot=select-trigger]:not([class*='w-'])]:w-fit [&>input]:flex-1",
			switch orientation {
			case "vertical":
				"flex-col [&>*:not(:first-child)]:rounded-t-none [&>*:not(:first-child)]:border-t-0 [&>*:not(:last-child)]:rounded-b-none"
			default:
				"[&>*:not(:first-child)]:rounded-l-none [&>*:not(:first-child)]:border-l-0 [&>*:not(:last-child)]:rounded-r-none"
			}
		}
		{ attrs... }
	>
		{ children }
	</div>
}

// ButtonGroupText's asChild tag-swap is dropped (GAP, always a <div>) — same
// narrow gap as Button's own asChild. Note this element carries no
// data-slot in shadcn's own source either (unlike every other button-group
// part); ported as-is rather than "fixed", per the token-for-token rule.
component ButtonGroupText(children gsx.Node, attrs gsx.Attrs) {
	<div
		class="flex items-center gap-2 rounded-md border bg-muted px-4 text-sm font-medium shadow-xs [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4"
		{ attrs... }
	>
		{ children }
	</div>
}

// ButtonGroupSeparator wraps ui.Separator directly (flat package, no
// re-implementation) — the button-group -> separator dependency
// internal/registry derives and registry_test.go pins. orientation defaults
// to "vertical" here, the opposite of Separator's own "horizontal" default,
// matching shadcn's `orientation = "vertical"` override for this call site
// (a button group's own axis is usually horizontal, so its separator is
// usually vertical). data-[orientation=vertical]:h-auto and bg-input both
// win their respective conflicts against Separator's own base classes via
// the ordinary caller-class-merge position (attrs after base, see
// docs/jsx-parity.md styling notes).
component ButtonGroupSeparator(orientation string, attrs gsx.Attrs) {
	<Separator
		data-slot="button-group-separator"
		orientation={orientation |> default("vertical")}
		class="relative m-0! self-stretch bg-input data-[orientation=vertical]:h-auto"
		{ attrs... }
	/>
}

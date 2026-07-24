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
//
// Retargeted to nova density (2026-07-24 nova density map, `## button-group`).
// DEVIATION from the map's own notes: the map frames nova's corner mechanism
// as inner-corner zeroing REPLACED by outer-corner restoration
// (`[&>[data-slot]:not(:has(~[data-slot]))]:rounded-r-lg!` /
// `rounded-b-lg!`). Checked against the actual nova source
// (shadcn-ui/apps/v4/registry/bases/radix/ui/button-group.tsx +
// styles/style-nova.css): the radix base's `buttonGroupVariants` — shared by
// every style, nova included — still carries the zero-inner-corner classes
// (`[&>*:not(:first-child)]:rounded-l-none/border-l-0
// [&>*:not(:last-child)]:rounded-r-none`) verbatim; nova's stylesheet only
// ADDS the restore rule as a supplementary `!important` override for the one
// case the zero rule gets wrong — a trailing non-slotted element (e.g. a
// visually-hidden `<select aria-hidden>`, see the root class's own
// `has-[select[aria-hidden=true]:last-child]` rule) that makes the true last
// *visible* child fail `:last-child`. Dropping the zero rule outright (a
// literal read of "replace") would leave every button at full `rounded-lg`
// on all four corners — no flush seam between group members, a real visual
// regression, not nova's actual behavior. Ported as ADD: the zero-corner
// selectors are kept unchanged and the outer-corner restore is layered on
// top, matching what nova really ships.

component ButtonGroup(orientation string, children gsx.Node, attrs gsx.Attrs) {
	<div
		role="group"
		data-slot="button-group"
		data-orientation={orientation |> default("horizontal")}
		class={
			"flex w-fit items-stretch has-[>[data-slot=button-group]]:gap-2 [&>*]:focus-visible:relative [&>*]:focus-visible:z-10 has-[select[aria-hidden=true]:last-child]:[&>[data-slot=select-trigger]:last-of-type]:rounded-r-lg [&>[data-slot=select-trigger]:not([class*='w-'])]:w-fit [&>input]:flex-1",
			switch orientation {
			case "vertical":
				"flex-col [&>*:not(:first-child)]:rounded-t-none [&>*:not(:first-child)]:border-t-0 [&>*:not(:last-child)]:rounded-b-none [&>[data-slot]:not(:has(~[data-slot]))]:rounded-b-lg!"
			default:
				"[&>*:not(:first-child)]:rounded-l-none [&>*:not(:first-child)]:border-l-0 [&>*:not(:last-child)]:rounded-r-none [&>[data-slot]:not(:has(~[data-slot]))]:rounded-r-lg!"
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
		class="flex items-center gap-2 rounded-lg border bg-muted px-2.5 text-sm font-medium [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4"
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

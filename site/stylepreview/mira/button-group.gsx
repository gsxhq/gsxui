package mira

import "github.com/gsxhq/gsx"

// ButtonGroup and its parts are the shadcn/ui ButtonGroup
// (registry/new-york-v4/ui/button-group.tsx). Its public orientation is
// reflected for the style pack and defaults via the house `|> default`
// pattern (see
// button.gsx/dropdown-menu.gsx) for consistency with every other data-variant
// stamp in this codebase — an ADAPT: shadcn leaves the attribute entirely
// unset when `orientation` is undefined (see docs/jsx-parity.md).
//
// Retargeted to nova density (2026-07-24 nova density map, `## button-group`).
// DEVIATION from the map's own notes: the map frames nova's corner mechanism
// as inner-corner zeroing REPLACED by priority outer-corner restoration
// based on the last visible child. Checked against the actual nova source
// (shadcn-ui/apps/v4/registry/bases/radix/ui/button-group.tsx +
// styles/style-nova.css): the radix base's `buttonGroupVariants` — shared by
// every style, nova included — still carries the zero-inner-corner classes
// (`[&>*:not(:first-child)]:rounded-l-none/border-l-0
// [&>*:not(:last-child)]:rounded-r-none`) verbatim; nova's stylesheet only
// ADDS a restore rule for the one
// case the zero rule gets wrong — a trailing visually-hidden
// `<select aria-hidden>` (see the root class's own
// `has-[select[aria-hidden=true]:last-child]` rule) that makes the true last
// visible child fail `:last-child`. Dropping the zero rule outright (a
// literal read of "replace") would leave every button at full `rounded-lg`
// on all four corners — no flush seam between group members, a real visual
// regression, not nova's actual behavior. Ported as ADD: the zero-corner
// selectors are kept unchanged and the outer-corner restore is layered on
// top, matching what nova really ships.

component ButtonGroup(orientation string, children gsx.Node, attrs gsx.Attrs) {
	<div
		role="group"
		data-orientation={orientation |> default("horizontal")}
		class={
			"flex w-fit items-stretch [&>*:focus-visible]:relative [&>*:focus-visible]:z-10 [&>input]:flex-1",
			switch orientation { case "vertical": "flex-col" default: "flex-row" }
		}
		{ attrs... }
		data-gsxui-slot-button-group
	>
		{ children }
	</div>
}

// ButtonGroupText's asChild tag-swap is dropped (GAP, always a <div>) — same
// narrow gap as Button's own asChild. Note this element carries no generic
// shadcn slot hook either (unlike every other button-group part); ported
// as-is rather than "fixed", per the token-for-token rule.
component ButtonGroupText(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"bg-muted gap-2 rounded-md border px-2.5 text-xs/relaxed font-medium [&_svg:not([class*='size-'])]:size-4 flex"
		}
		{ attrs... }
		data-gsxui-slot-button-group-text
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
		orientation={orientation |> default("vertical")}
		class={ "relative m-0 self-stretch bg-input data-[orientation=vertical]:h-auto" }
		{ attrs... }
		data-gsxui-slot-button-group-separator
	/>
}

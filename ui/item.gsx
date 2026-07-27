// Item and its parts are the shadcn/ui Item family
// (registry/new-york-v4/ui/item.tsx): a generic media + content + actions
// row, with ItemGroup/ItemSeparator for stacked lists and
// ItemHeader/ItemFooter for framing. All cva variant maps are static class
// blocks (no data-keyed selectors in the source). The CSS-only contract
// reflects those axes through data attributes; see docs/jsx-parity.md
// `## item` for the drop list (asChild) and mechanisms.
package ui

import "github.com/gsxhq/gsx"

component ItemGroup(children gsx.Node, attrs gsx.Attrs) {
	<div role="list" { withSlot("item-group", attrs)... }>
		{ children }
	</div>
}

// ItemSeparator composes ui.Separator directly (flat package, no
// re-implementation) — the item -> separator dependency internal/registry
// derives and registry_test.go pins. shadcn's own version types its props as
// React.ComponentProps<typeof Separator> and spreads {...props} after its
// explicit orientation="horizontal", so a caller-supplied orientation prop
// wins there. `orientation` here is a real Go param (not left to attrs) for
// exactly that reason: attrs is untyped fallthrough onto Separator's own
// rendered <div>, not a hook into the orientation argument compiled into the
// call to Separator, so only an explicit param can actually override it —
// same competing-defaults mechanism as ButtonGroupSeparator's own
// orientation |> default("vertical").
component ItemSeparator(orientation string, attrs gsx.Attrs) {
	<Separator
		orientation={orientation |> default("horizontal")}
		{ withSlot("item-separator", attrs)... }
	/>
}

// Item's variant/size cva map (itemVariants) picks between static class
// blocks by the JS-resolved prop values — no data-[variant=...]/
// data-[size=...] selectors in registry/new-york-v4/ui/item.tsx to preserve,
// so both are reflected as explicit data attributes for the style pack.
//
// asChild tag-swapping is dropped — always renders a <div> — the same
// narrow gap as button's own asChild (see docs/jsx-parity.md ## button).
component Item(variant string, size string, children gsx.Node, attrs gsx.Attrs) {
	<div
		data-variant={variant |> default("default")}
		data-size={size |> default("default")}
		{ withSlot("item", attrs)... }
	>
		{ children }
	</div>
}

// ItemMedia's variant cva map (itemMediaVariants) is reflected through
// data-variant for the style pack.
//
// Retargeted to nova density (2026-07-24 nova density map, `## item`).
// DEVIATION: nova's icon-media variant drops the bordered/muted size-8 box
// entirely (bare 1rem svg, no container) and the image variant gains
// ancestor-size responsive sizing tied to a `size=xs` axis this
// package doesn't ship — both left as-is here: the box drop bundles a
// color/border removal (border, bg-muted) the retarget is scoped to leave
// alone, and the responsive image sizing is half dead weight without a real
// xs size param (Item's `size` axis stays sm/default only, per task scope).
component ItemMedia(variant string, children gsx.Node, attrs gsx.Attrs) {
	<div
		data-variant={variant |> default("default")}
		{ withSlot("item-media", attrs)... }
	>
		{ children }
	</div>
}

component ItemContent(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("item-content", attrs)... }>
		{ children }
	</div>
}

component ItemTitle(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("item-title", attrs)... }>
		{ children }
	</div>
}

// ItemDescription renders a real <p>, matching shadcn's own source exactly
// (unlike EmptyDescription, whose type says "p" but whose element is a
// <div> — see empty.gsx).
component ItemDescription(children gsx.Node, attrs gsx.Attrs) {
	<p { withSlot("item-description", attrs)... }>
		{ children }
	</p>
}

component ItemActions(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("item-actions", attrs)... }>
		{ children }
	</div>
}

component ItemHeader(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("item-header", attrs)... }>
		{ children }
	</div>
}

component ItemFooter(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("item-footer", attrs)... }>
		{ children }
	</div>
}

// Item and its parts are the shadcn/ui Item family
// (registry/new-york-v4/ui/item.tsx): a generic media + content + actions
// row, with ItemGroup/ItemSeparator for stacked lists and
// ItemHeader/ItemFooter for framing. All cva variant maps are static class
// blocks (no data-keyed selectors in the source). The CSS-only contract
// reflects those axes through data attributes; see docs/jsx-parity.md
// `## item` for the drop list (asChild) and mechanisms.
package maia

import "github.com/gsxhq/gsx"

component ItemGroup(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={ "gap-4 has-data-[size=sm]:gap-2.5 has-data-[size=xs]:gap-2 flex flex-col" }
		role="list"
		{ attrs... }
		data-gsxui-slot-item-group
	>
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
		class={ "my-2" }
		{ attrs... }
		data-gsxui-slot-item-separator
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
		class={
			"group/item",
			"[a]:hover:bg-muted rounded-2xl border text-sm flex flex-wrap",
			switch variant {
			case "outline":
				"border-border"
			case "muted":
				"bg-muted/50 border-transparent"
			default:
				"border-transparent"
			}
		}
		{ attrs... }
		data-gsxui-slot-item
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
		class={
			"gap-2 flex",
			switch variant {
			case "icon":
				"[&_svg:not([class*='size-'])]:size-4"
			case "image":
				"size-10 overflow-hidden rounded-lg group-data-[size=sm]/item:size-8 group-data-[size=xs]/item:size-6 group-data-[size=xs]/item:rounded-md [&_img]:size-full [&_img]:object-cover"
			default:
				"bg-transparent"
			}
		}
		{ attrs... }
		data-gsxui-slot-item-media
	>
		{ children }
	</div>
}

component ItemContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"flex-1 [&+[data-gsxui-slot-item-content]]:flex-none gap-1 group-data-[size=xs]/item:gap-0.5 flex flex-col"
		}
		{ attrs... }
		data-gsxui-slot-item-content
	>
		{ children }
	</div>
}

component ItemTitle(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={ "gap-2 text-sm leading-snug font-medium underline-offset-4 flex" }
		{ attrs... }
		data-gsxui-slot-item-title
	>
		{ children }
	</div>
}

// ItemDescription renders a real <p>, matching shadcn's own source exactly
// (unlike EmptyDescription, whose type says "p" but whose element is a
// <div> — see empty.gsx).
component ItemDescription(children gsx.Node, attrs gsx.Attrs) {
	<p class={ "text-muted-foreground text-left text-sm" } { attrs... } data-gsxui-slot-item-description>
		{ children }
	</p>
}

component ItemActions(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "gap-2 flex" } { attrs... } data-gsxui-slot-item-actions>
		{ children }
	</div>
}

component ItemHeader(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "gap-2 flex" } { attrs... } data-gsxui-slot-item-header>
		{ children }
	</div>
}

component ItemFooter(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "gap-2 flex" } { attrs... } data-gsxui-slot-item-footer>
		{ children }
	</div>
}

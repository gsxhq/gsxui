// Item and its parts are the shadcn/ui Item family
// (registry/new-york-v4/ui/item.tsx): a generic media + content + actions
// row, with ItemGroup/ItemSeparator for stacked lists and
// ItemHeader/ItemFooter for framing. All cva variant maps are static class
// blocks (no data-keyed selectors in the source), so they port as switches
// inside class={} — see docs/jsx-parity.md `## item` for the drop list
// (asChild) and mechanisms.
package ui

import "github.com/gsxhq/gsx"

const itemBase = "group/item flex flex-wrap items-center rounded-md border border-transparent text-sm transition-colors duration-100 outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 [a]:transition-colors [a]:hover:bg-accent/50"

component ItemGroup(children gsx.Node, attrs gsx.Attrs) {
	<div role="list" data-slot="item-group" class="group/item-group flex flex-col" { attrs... }>
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
		data-slot="item-separator"
		orientation={orientation |> default("horizontal")}
		class="my-0"
		{ attrs... }
	/>
}

// Item's variant/size cva map (itemVariants) picks between static class
// blocks by the JS-resolved prop values — no data-[variant=...]/
// data-[size=...] selectors in registry/new-york-v4/ui/item.tsx to preserve,
// so both port as switches inside class={}, the same idiom as
// button.gsx's variantClass/sizeClass pair (here inlined, since only Item
// itself uses this pair — no sibling component shares it the way
// pagination.gsx shares button.gsx's helpers).
//
// asChild tag-swapping is dropped — always renders a <div> — the same
// narrow gap as button's own asChild (see docs/jsx-parity.md ## button).
component Item(variant string, size string, children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="item"
		data-variant={variant |> default("default")}
		data-size={size |> default("default")}
		class={
			itemBase,
			switch variant { case "outline": "border-border" case "muted": "bg-muted/50" default: "bg-transparent" },
			switch size { case "sm": "gap-2.5 px-4 py-3" default: "gap-4 p-4" }
		}
		{ attrs... }
	>
		{ children }
	</div>
}

// ItemMedia's variant cva map (itemMediaVariants), same static-block shape
// as Item's own — ported as a switch inside class={}.
component ItemMedia(variant string, children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="item-media"
		data-variant={variant |> default("default")}
		class={
			"flex shrink-0 items-center justify-center gap-2 group-has-[[data-slot=item-description]]/item:translate-y-0.5 group-has-[[data-slot=item-description]]/item:self-start [&_svg]:pointer-events-none",
			switch variant {
			case "icon":
				"size-8 rounded-sm border bg-muted [&_svg:not([class*='size-'])]:size-4"
			case "image":
				"size-10 overflow-hidden rounded-sm [&_img]:size-full [&_img]:object-cover"
			default:
				"bg-transparent"
			}
		}
		{ attrs... }
	>
		{ children }
	</div>
}

component ItemContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="item-content"
		class="flex flex-1 flex-col gap-1 [&+[data-slot=item-content]]:flex-none"
		{ attrs... }
	>
		{ children }
	</div>
}

component ItemTitle(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="item-title"
		class="flex w-fit items-center gap-2 text-sm leading-snug font-medium"
		{ attrs... }
	>
		{ children }
	</div>
}

// ItemDescription renders a real <p>, matching shadcn's own source exactly
// (unlike EmptyDescription, whose type says "p" but whose element is a
// <div> — see empty.gsx).
component ItemDescription(children gsx.Node, attrs gsx.Attrs) {
	<p
		data-slot="item-description"
		class="line-clamp-2 text-sm leading-normal font-normal text-balance text-muted-foreground [&>a]:underline [&>a]:underline-offset-4 [&>a:hover]:text-primary"
		{ attrs... }
	>
		{ children }
	</p>
}

component ItemActions(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="item-actions" class="flex items-center gap-2" { attrs... }>
		{ children }
	</div>
}

component ItemHeader(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="item-header"
		class="flex basis-full items-center justify-between gap-2"
		{ attrs... }
	>
		{ children }
	</div>
}

component ItemFooter(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="item-footer"
		class="flex basis-full items-center justify-between gap-2"
		{ attrs... }
	>
		{ children }
	</div>
}

package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// ToggleGroup declares the root and item slots default.css styles. Every
// rule in default.css is authored as the ROOT's own data-attribute reaching
// down into the item as a descendant (`[data-size="sm"] > [...item]`), but
// ToggleGroupItem has no context/broadcast mechanism (see ui/toggle-group.gsx's
// package doc comment, GAP) — the caller passes variant/size/spacing to BOTH
// ToggleGroup and every ToggleGroupItem explicitly, so the item already
// carries its own mirrored data-variant/data-size/data-spacing/
// data-orientation. This shape re-keys every item rule onto the ITEM's OWN
// dimension value instead of reaching up into the parent (which Tailwind's
// arbitrary-variant syntax has no idiom for anyway) — computed style is
// identical whenever caller passes matching values to both (the only
// supported usage; verified against the sweep), and is the same "each slot
// styled from its own accessor call, using params already in scope at that
// callsite" pattern every other multi-slot component uses.
//
// orientation is a real shapeable dimension (`data-orientation`, not
// `aria-orientation` — contrast Resizable's handle, which could not be a
// dimension because its axis is `aria-orientation`), but the contract and
// default.css both only ever exercise one value ("horizontal" — vertical is
// out of scope GAP per the package doc comment), so Values has exactly one
// member. disabled/data-state/aria-invalid are runtime state on the ITEM,
// handled entirely by Toggle's own shape/recipe (ToggleGroupItem also
// carries data-gsxui-slot-toggle and composes Toggle's presentation) — not
// this shape's concern.
var ToggleGroup = recipe.Shape{
	Component: "toggle-group",
	Slots: []recipe.Slot{
		{
			Name: "", Base: true,
			Dimensions: []recipe.Dimension{
				{Name: "size", Default: "default", Values: []string{"default", "sm", "lg"}},
			},
		},
		{
			Name: "item", Base: true,
			Dimensions: []recipe.Dimension{
				{Name: "size", Default: "default", Values: []string{"default", "sm", "lg"}},
				{Name: "spacing", Default: "0", Values: []string{"0", "2"}},
				{Name: "orientation", Default: "horizontal", Values: []string{"horizontal"}},
			},
		},
	},
}

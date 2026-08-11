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
// member.
//
// disabled/data-state/aria-invalid ARE this shape's concern, in the sense
// that matters: the ITEM slot's recipe (registry/styles/<style>/
// toggle-group.css) composes Toggle's own per-style state rules directly
// onto `.gsxui-recipe-toggle-group-item`, using the values from that SAME
// style's registry/styles/<style>/toggle.css. This used to be handled by one
// shared, style-invariant fallback in assets/css/styles/default.css instead
// (a single hardcoded reference look for every style) — the 8-style port
// retired it once 8 materially different Toggle looks made that fallback
// visibly wrong for most installed styles, and folded Toggle's presentation
// into this shape's own per-style recipe instead. See
// registry/canonical/toggle-group.gsx's package doc comment and
// registry/styles/nova/toggle-group.css's own comment on the item rule.
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

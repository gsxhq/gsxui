package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Separator is a single-element component: one root slot carrying an
// orientation dimension, matching default.css's data-orientation axis and
// the style contract exactly (horizontal, default; vertical) — no
// disagreement to resolve here, unlike Badge's variant list.
//
// default.css also carries three same-element compound rules where
// Separator's own marker co-occurs with ANOTHER component's marker on one
// element (ButtonGroup's, Item's and Field's own "-separator" slots,
// e.g. `[data-gsxui-slot-separator][data-gsxui-slot-button-group-separator]`).
// Those rules are left untouched in default.css: the styled element carries
// both markers, but the properties they set (m-0, self-stretch, bg-input,
// absolute positioning) are specific to how ButtonGroup/Item/Field use a
// composed Separator, not to Separator itself, and none of those three
// components are migrated yet.
var Separator = recipe.Shape{
	Component: "separator",
	Slots: []recipe.Slot{{
		Name: "", Base: true,
		Dimensions: []recipe.Dimension{
			{Name: "orientation", Default: "horizontal", Values: []string{"horizontal", "vertical"}},
		},
	}},
}

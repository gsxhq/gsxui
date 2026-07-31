package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Toast is the toast card — a compound component of nine plain slots and no
// dimensions at all.
//
// The style contract declares a data-type axis on the root slot, and the
// default style does react to it (the type glyph's colour, and the spin on
// a loading toast). That axis is deliberately NOT a recipe dimension: the
// attribute is mutated at RUNTIME — ui/toaster.js sets data-type on every
// cloned card and swaps it in place when toast.promise settles (its morph()
// turns a loading toast into a success/error one without re-rendering) — so
// a class baked in at render time would go stale the moment the type
// changed. The rules keyed on it stay attribute-keyed, folded into the icon
// slot's own recipe rule as ancestor-conditional arbitrary variants (see
// registry/styles/nova/toast.css). A shape declaring less than the contract
// is allowed; declaring a dimension that cannot follow the DOM is not.
var Toast = recipe.Shape{
	Component: "toast",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "icon", Base: true},
		{Name: "content", Base: true},
		{Name: "title", Base: true},
		{Name: "description", Base: true},
		{Name: "action", Base: true},
		{Name: "cancel", Base: true},
		{Name: "close", Base: true},
		{Name: "close-icon", Base: true},
	},
}

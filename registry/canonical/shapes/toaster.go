package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Toaster is the positioned toast region: one slot, no dimensions. The
// style contract declares a data-expanded axis on it, but that is pure
// runtime state (ui/toaster.js sets it while the stack is hover-expanded)
// with no rule of its own in the default style, so the shape declares no
// dimension for it — a shape may declare less than the contract.
var Toaster = recipe.Shape{
	Component: "toaster",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
	},
}

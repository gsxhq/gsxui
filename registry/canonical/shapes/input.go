package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Input is a single-slot form control, no dimensions. disabled and
// aria-invalid are runtime state, not shape dimensions — see
// internal/stylecontract's Input contract, same shape as Textarea's.
var Input = recipe.Shape{
	Component: "input",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
	},
}

package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Spinner is a single-element component: one root slot, no dimensions, one
// rule in default.css.
var Spinner = recipe.Shape{
	Component: "spinner",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
	},
}

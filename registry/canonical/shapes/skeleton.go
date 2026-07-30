package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Skeleton is a single-element loading placeholder: one root slot, no
// dimensions, one rule in default.css.
var Skeleton = recipe.Shape{
	Component: "skeleton",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
	},
}

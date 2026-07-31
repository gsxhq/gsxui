package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Progress is a two-slot container: the track (root) and the indicator
// child. Neither slot carries a dimension.
var Progress = recipe.Shape{
	Component: "progress",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "indicator", Base: true},
	},
}

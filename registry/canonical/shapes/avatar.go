package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Avatar is a three-slot compound component: root, image and fallback.
// None of the three slots carries a dimension.
var Avatar = recipe.Shape{
	Component: "avatar",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "image", Base: true},
		{Name: "fallback", Base: true},
	},
}

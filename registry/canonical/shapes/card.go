package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Card is a pure container: seven slots, no dimensions. Every slot carries only
// a base rule, which is the common shape across the catalogue — most components
// vary structurally rather than by variant.
var Card = recipe.Shape{
	Component: "card",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "header", Base: true},
		{Name: "title", Base: true},
		{Name: "description", Base: true},
		{Name: "action", Base: true},
		{Name: "content", Base: true},
		{Name: "footer", Base: true},
	},
}

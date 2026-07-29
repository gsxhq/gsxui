package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Alert is a compound component: a root slot carrying a variant dimension,
// plus two plain child slots (title, description) with no dimensions of
// their own. The variant values match both the style contract's
// `data-variant` axis and default.css's own [data-variant="…"] rules exactly
// — no discrepancy here, unlike Badge.
var Alert = recipe.Shape{
	Component: "alert",
	Slots: []recipe.Slot{
		{
			Name: "", Base: true,
			Dimensions: []recipe.Dimension{
				{Name: "variant", Default: "default", Values: []string{
					"default", "destructive"}},
			},
		},
		{Name: "title", Base: true},
		{Name: "description", Base: true},
	},
}

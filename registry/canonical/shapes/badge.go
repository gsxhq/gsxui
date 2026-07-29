package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Badge is a single-element component: one root slot carrying a variant
// dimension. Its values match the style contract's `data-variant` axis
// (internal/stylecontract's Badge slot) rather than the smaller illustrative
// list this migration's brief quoted — "ghost" and "link" are both real,
// upstream-sourced variants even though "ghost" carries only hover utilities
// (see registry/styles/nova/badge.css's doc comment for why that's not a gap).
var Badge = recipe.Shape{
	Component: "badge",
	Slots: []recipe.Slot{{
		Name: "", Base: true,
		Dimensions: []recipe.Dimension{
			{Name: "variant", Default: "default", Values: []string{
				"default", "secondary", "destructive", "outline", "ghost", "link"}},
		},
	}},
}

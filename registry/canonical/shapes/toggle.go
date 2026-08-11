package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Toggle is a single-element component (a real <button type="button">): one
// root slot carrying the public variant/size parameters as dimensions.
// pressed/data-state is NOT a dimension — ui/toggle.js re-stamps
// data-state="on"/"off" (and aria-pressed) on every click without a server
// round trip, so it must stay a runtime Tailwind state variant
// (`data-[state=on]:`) folded onto the one recipe rule, the same treatment
// disabled/aria-invalid already get on Button and Switch. See
// registry/styles/nova/toggle.css for the full fold.
//
// ui/toggle-group.gsx's ToggleGroupItem stamps data-gsxui-slot-toggle
// alongside its own data-gsxui-slot-toggle-group-item marker, WITHOUT
// calling any of this shape's accessors — see
// registry/canonical/shapes/toggle-group.go for why ToggleGroupItem's own
// recipe carries its own copy of this component's per-style presentation
// instead (the 8-style port retired the single shared, style-invariant
// marker fallback this comment used to describe: with 8 materially
// different Toggle looks, one hardcoded reference rendering was no longer a
// reasonable compromise).
var Toggle = recipe.Shape{
	Component: "toggle",
	Slots: []recipe.Slot{{
		Name: "",
		Base: true,
		Dimensions: []recipe.Dimension{
			{Name: "variant", Default: "default", Values: []string{"default", "outline"}},
			{Name: "size", Default: "default", Values: []string{"default", "sm", "lg"}},
		},
	}},
}

package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Radio is a native <input type="radio">: one root slot, no dimensions.
// checked/disabled/aria-invalid are runtime state (matching the style
// contract's own axes on it), not recipe dimensions — same treatment as
// Badge/Alert's aria-invalid and Resizable's aria-orientation.
var Radio = recipe.Shape{
	Component: "radio",
	Slots:     []recipe.Slot{{Name: "", Base: true}},
}

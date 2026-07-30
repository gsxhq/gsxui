package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Checkbox is a single-element component: one root slot, no dimensions. Its
// checked-glyph background-image (a base64 data-URI, light and dark) has no
// Tailwind utility form and stays in assets/css/styles/default.css's
// @layer utilities escape hatch — see registry/styles/nova/checkbox.css's
// doc comment.
var Checkbox = recipe.Shape{
	Component: "checkbox",
	Slots:     []recipe.Slot{{Name: "", Base: true}},
}

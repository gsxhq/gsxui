package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Label is a single-element component: one root slot, no dimensions. Its
// default.css block also carries two relational rules that react to
// context OUTSIDE the component's own markers rather than to a dimension —
// an ancestor with data-disabled="true", and a preceding sibling that is
// :disabled — so nothing about them belongs on the shape as a Dimension;
// they translate directly as arbitrary variants on the root slot's recipe
// rule (see registry/styles/nova/label.css).
var Label = recipe.Shape{
	Component: "label",
	Slots:     []recipe.Slot{{Name: "", Base: true}},
}

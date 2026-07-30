package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// ScrollArea is a single-element component: one root slot, no dimensions.
// The Radix four-part tree (Root/Viewport/ScrollBar/Thumb) collapses onto one
// native `overflow-auto` div (see ui/scroll-area.gsx's doc comment for the
// full rationale). `orientation` is a plain data-attribute switch handled by
// assets/css/foundation.css (`overflow`/`overflow-x`, shared across every
// style and untouched by the slot axis), not a styled dimension of this
// shape — the recipe layer (assets/css/styles/default/scroll-area.css) has
// no rule keyed on data-orientation at all, only on the plain marker. This
// is the second HYPHENATED component migrated ("scroll-area" kebab identity,
// derived camel Go identifier `scrollAreaRecipe`/`scrollArea`, same
// mechanism as AspectRatio).
var ScrollArea = recipe.Shape{
	Component: "scroll-area",
	Slots:     []recipe.Slot{{Name: "", Base: true}},
}

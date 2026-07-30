package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Slider is a single-element component (a real <input type="range">): one
// root slot, no dimensions. Its default.css block carries no author-facing
// variant/size axis at all — every relational rule is either a pseudo-class
// state (:disabled, :hover, :focus-visible, :active) or a pseudo-ELEMENT
// selector (::-webkit-slider-runnable-track, ::-moz-range-track,
// ::-webkit-slider-thumb, ::-moz-range-thumb) targeting the browser's own
// generated parts of the range input. None of these are dimensions: nothing
// about them is chosen by the caller at render time, so they translate as
// arbitrary Tailwind variants stacked on the ONE recipe rule
// (`disabled:`, `hover:`, `[&::-webkit-slider-thumb]:`, …), the same
// mechanism Label's ancestor/sibling variants and Kbd's foreign-ancestor
// variant already proved — see registry/styles/nova/slider.css.
var Slider = recipe.Shape{
	Component: "slider",
	Slots:     []recipe.Slot{{Name: "", Base: true}},
}

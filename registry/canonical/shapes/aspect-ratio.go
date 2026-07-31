package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// AspectRatio is a single-element component: one root slot, no dimensions,
// one utility (block). It is the first HYPHENATED component migrated onto
// the slot axis — its registry/CSS identity stays kebab-case
// ("aspect-ratio", emitted class "gsxui-recipe-aspect-ratio"), while
// stylegen derives the Go type/receiver as camel (aspectRatioRecipe /
// aspectRatio) via componentIdentifier, because a hyphen is illegal in a Go
// identifier and "aspect-ratio.Root()" would otherwise parse as subtraction.
var AspectRatio = recipe.Shape{
	Component: "aspect-ratio",
	Slots:     []recipe.Slot{{Name: "", Base: true}},
}

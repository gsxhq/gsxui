package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// slider binds Slider's shape to the accessor calls slider.gsx authors.
var slider = sliderRecipe{c: recipe.Component{Shape: shapes.Slider}}

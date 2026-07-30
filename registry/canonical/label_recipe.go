package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// label binds Label's shape to the accessor calls label.gsx authors.
var label = labelRecipe{c: recipe.Component{Shape: shapes.Label}}

package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// chart binds Chart's shape to the accessor calls chart.gsx authors.
var chart = chartRecipe{c: recipe.Component{Shape: shapes.Chart}}

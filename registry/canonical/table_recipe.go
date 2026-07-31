package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// table binds Table's shape to the accessor calls table.gsx authors.
var table = tableRecipe{c: recipe.Component{Shape: shapes.Table}}

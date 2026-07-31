package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// input binds Input's shape to the accessor calls input.gsx authors.
var input = inputRecipe{c: recipe.Component{Shape: shapes.Input}}

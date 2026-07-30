package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// field binds Field's shape to the accessor calls field.gsx authors.
var field = fieldRecipe{c: recipe.Component{Shape: shapes.Field}}

package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// empty binds Empty's shape to the accessor calls empty.gsx authors.
var empty = emptyRecipe{c: recipe.Component{Shape: shapes.Empty}}

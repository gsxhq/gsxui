package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// resizable binds Resizable's shape to the accessor calls resizable.gsx
// authors.
var resizable = resizableRecipe{c: recipe.Component{Shape: shapes.Resizable}}

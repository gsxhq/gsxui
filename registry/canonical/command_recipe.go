package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// command binds Command's shape to the accessor calls command.gsx authors.
var command = commandRecipe{c: recipe.Component{Shape: shapes.Command}}

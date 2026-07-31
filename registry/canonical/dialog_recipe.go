package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// dialog binds Dialog's shape to the accessor calls dialog.gsx authors.
var dialog = dialogRecipe{c: recipe.Component{Shape: shapes.Dialog}}

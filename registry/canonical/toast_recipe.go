package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// toast binds Toast's shape to the accessor calls toast.gsx authors.
var toast = toastRecipe{c: recipe.Component{Shape: shapes.Toast}}

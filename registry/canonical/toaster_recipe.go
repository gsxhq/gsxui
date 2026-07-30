package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// toaster binds Toaster's shape to the accessor calls toaster.gsx authors.
var toaster = toasterRecipe{c: recipe.Component{Shape: shapes.Toaster}}

package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// alertDialog binds AlertDialog's shape to the accessor calls
// alert-dialog.gsx authors.
var alertDialog = alertDialogRecipe{c: recipe.Component{Shape: shapes.AlertDialog}}

package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// spinner binds Spinner's shape to the accessor calls spinner.gsx authors.
// The variable name is the component name: stylegen resolves spinner.Root()
// by looking "spinner" up in shapes.All().
var spinner = spinnerRecipe{c: recipe.Component{Shape: shapes.Spinner}}

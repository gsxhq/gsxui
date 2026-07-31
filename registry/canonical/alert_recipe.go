package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// alert binds Alert's shape to the accessor calls alert.gsx authors. The
// variable name is the component name: stylegen resolves alert.Root() by
// looking "alert" up in shapes.All().
var alert = alertRecipe{c: recipe.Component{Shape: shapes.Alert}}

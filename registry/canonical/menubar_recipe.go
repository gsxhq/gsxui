package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// menubar binds Menubar's shape to the accessor calls menubar.gsx authors.
var menubar = menubarRecipe{c: recipe.Component{Shape: shapes.Menubar}}

package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// drawer binds Drawer's shape to the accessor calls drawer.gsx authors.
var drawer = drawerRecipe{c: recipe.Component{Shape: shapes.Drawer}}

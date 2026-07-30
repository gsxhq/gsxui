package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// sidebar binds Sidebar's shape to the accessor calls sidebar.gsx authors.
var sidebar = sidebarRecipe{c: recipe.Component{Shape: shapes.Sidebar}}

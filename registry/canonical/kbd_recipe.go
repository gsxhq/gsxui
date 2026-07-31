package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// kbd binds Kbd's shape to the accessor calls kbd.gsx authors.
var kbd = kbdRecipe{c: recipe.Component{Shape: shapes.Kbd}}

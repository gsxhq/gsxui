package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// popover binds Popover's shape to the accessor calls popover.gsx authors.
var popover = popoverRecipe{c: recipe.Component{Shape: shapes.Popover}}

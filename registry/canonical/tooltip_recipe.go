package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// tooltip binds Tooltip's shape to the accessor calls tooltip.gsx authors.
var tooltip = tooltipRecipe{c: recipe.Component{Shape: shapes.Tooltip}}

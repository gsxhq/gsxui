package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// collapsible binds Collapsible's shape to the accessor calls
// collapsible.gsx authors.
var collapsible = collapsibleRecipe{c: recipe.Component{Shape: shapes.Collapsible}}

package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// separator binds Separator's shape to the accessor calls separator.gsx
// authors.
var separator = separatorRecipe{c: recipe.Component{Shape: shapes.Separator}}

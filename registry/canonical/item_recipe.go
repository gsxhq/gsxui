package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// item binds Item's shape to the accessor calls item.gsx authors.
var item = itemRecipe{c: recipe.Component{Shape: shapes.Item}}

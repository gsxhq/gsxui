package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// sheet binds Sheet's shape to the accessor calls sheet.gsx authors.
var sheet = sheetRecipe{c: recipe.Component{Shape: shapes.Sheet}}

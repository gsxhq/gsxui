package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// combobox binds Combobox's shape to the accessor calls combobox.gsx authors.
var combobox = comboboxRecipe{c: recipe.Component{Shape: shapes.Combobox}}

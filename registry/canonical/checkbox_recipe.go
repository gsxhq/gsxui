package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// checkbox binds Checkbox's shape to the accessor calls checkbox.gsx
// authors.
var checkbox = checkboxRecipe{c: recipe.Component{Shape: shapes.Checkbox}}

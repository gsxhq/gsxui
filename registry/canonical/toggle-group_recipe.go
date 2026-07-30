package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// toggleGroup binds ToggleGroup's shape to the accessor calls
// toggle-group.gsx authors.
var toggleGroup = toggleGroupRecipe{c: recipe.Component{Shape: shapes.ToggleGroup}}

package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// navigationMenu binds NavigationMenu's shape to the accessor calls
// navigation-menu.gsx authors. The variable name is the derived Go
// identifier for the "navigation-menu" component (see
// internal/stylegen/identifier.go).
var navigationMenu = navigationMenuRecipe{c: recipe.Component{Shape: shapes.NavigationMenu}}

package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// dropdownMenu binds DropdownMenu's shape to the accessor calls
// dropdown-menu.gsx authors. The variable name is the derived Go identifier
// for the "dropdown-menu" component (see internal/stylegen/identifier.go):
// stylegen resolves dropdownMenu.Item() by looking "dropdown-menu" up in
// shapes.All().
var dropdownMenu = dropdownMenuRecipe{c: recipe.Component{Shape: shapes.DropdownMenu}}

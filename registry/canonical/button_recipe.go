package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// button binds Button's shape to the accessor calls button.gsx authors. The
// variable name is the component name: stylegen resolves button.Root() by
// looking "button" up in shapes.All().
var button = buttonRecipe{c: recipe.Component{Shape: shapes.Button}}

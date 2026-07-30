package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// toggle binds Toggle's shape to the accessor calls toggle.gsx authors. The
// variable name is the component name: stylegen resolves toggle.Root() by
// looking "toggle" up in shapes.All().
var toggle = toggleRecipe{c: recipe.Component{Shape: shapes.Toggle}}

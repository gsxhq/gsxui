package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// avatar binds Avatar's shape to the accessor calls avatar.gsx authors. The
// variable name is the component name: stylegen resolves avatar.Root() by
// looking "avatar" up in shapes.All().
var avatar = avatarRecipe{c: recipe.Component{Shape: shapes.Avatar}}

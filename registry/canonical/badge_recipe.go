package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// badge binds Badge's shape to the accessor calls badge.gsx authors. The
// variable name is the component name: stylegen resolves badge.Root() by
// looking "badge" up in shapes.All().
var badge = badgeRecipe{c: recipe.Component{Shape: shapes.Badge}}

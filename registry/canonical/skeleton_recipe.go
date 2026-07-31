package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// skeleton binds Skeleton's shape to the accessor calls skeleton.gsx
// authors. The variable name is the component name: stylegen resolves
// skeleton.Root() by looking "skeleton" up in shapes.All().
var skeleton = skeletonRecipe{c: recipe.Component{Shape: shapes.Skeleton}}

package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// progress binds Progress's shape to the accessor calls progress.gsx
// authors. The variable name is the component name: stylegen resolves
// progress.Root() by looking "progress" up in shapes.All().
var progress = progressRecipe{c: recipe.Component{Shape: shapes.Progress}}

package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// accordion binds Accordion's shape to the accessor calls accordion.gsx
// authors.
var accordion = accordionRecipe{c: recipe.Component{Shape: shapes.Accordion}}

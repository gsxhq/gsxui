package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// radio binds Radio's shape to the accessor calls radio.gsx authors.
var radio = radioRecipe{c: recipe.Component{Shape: shapes.Radio}}

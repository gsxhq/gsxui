package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// aspectRatio binds AspectRatio's shape to the accessor calls
// aspect-ratio.gsx authors. The receiver is camel (aspectRatio), derived
// from the hyphenated registry name "aspect-ratio" by componentIdentifier —
// "aspect-ratio.Root()" is not legal Go (a hyphen parses as subtraction).
var aspectRatio = aspectRatioRecipe{c: recipe.Component{Shape: shapes.AspectRatio}}

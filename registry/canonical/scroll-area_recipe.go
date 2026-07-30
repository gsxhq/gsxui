package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// scrollArea binds ScrollArea's shape to the accessor calls
// scroll-area.gsx authors. The receiver is camel (scrollArea), derived from
// the hyphenated registry name "scroll-area" by componentIdentifier —
// "scroll-area.Root()" is not legal Go (a hyphen parses as subtraction).
var scrollArea = scrollAreaRecipe{c: recipe.Component{Shape: shapes.ScrollArea}}

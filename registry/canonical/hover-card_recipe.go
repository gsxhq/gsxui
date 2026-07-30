package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// hoverCard binds HoverCard's shape to the accessor calls hover-card.gsx
// authors. The receiver is camel (hoverCard), derived from the hyphenated
// registry name "hover-card" by componentIdentifier.
var hoverCard = hoverCardRecipe{c: recipe.Component{Shape: shapes.HoverCard}}

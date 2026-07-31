package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// card binds Card's shape to the accessor calls card.gsx authors. The
// variable name is the component name: stylegen resolves card.Header() by
// looking "card" up in shapes.All().
var card = cardRecipe{c: recipe.Component{Shape: shapes.Card}}

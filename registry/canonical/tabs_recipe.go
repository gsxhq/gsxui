package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// tabs binds Tabs' shape to the accessor calls tabs.gsx authors.
var tabs = tabsRecipe{c: recipe.Component{Shape: shapes.Tabs}}

package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// contextMenu binds ContextMenu's shape to the accessor calls
// context-menu.gsx authors. The variable name is the derived Go identifier
// for the "context-menu" component (see internal/stylegen/identifier.go).
var contextMenu = contextMenuRecipe{c: recipe.Component{Shape: shapes.ContextMenu}}

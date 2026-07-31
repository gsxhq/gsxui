package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// buttonGroup binds ButtonGroup's shape to the accessor calls
// button-group.gsx authors. The registry/CSS identity stays kebab
// ("button-group", .gsxui-recipe-button-group); the Go identifiers are the
// camel forms internal/stylegen/identifier.go derives from it.
var buttonGroup = buttonGroupRecipe{c: recipe.Component{Shape: shapes.ButtonGroup}}

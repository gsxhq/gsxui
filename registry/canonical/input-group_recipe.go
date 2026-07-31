package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// inputGroup binds InputGroup's shape to the accessor calls input-group.gsx
// authors. The registry/CSS identity stays kebab ("input-group",
// .gsxui-recipe-input-group); the Go identifiers are the camel forms
// internal/stylegen/identifier.go derives from it.
var inputGroup = inputGroupRecipe{c: recipe.Component{Shape: shapes.InputGroup}}

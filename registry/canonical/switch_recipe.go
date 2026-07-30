package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// switch_ binds Switch's shape to the accessor calls switch.gsx authors.
// The trailing underscore is componentIdentifier's suffix for "switch"
// colliding with the Go keyword — the CSS/registry identity is unaffected
// (gsxui-recipe-switch, data-gsxui-slot-switch).
var switch_ = switch_Recipe{c: recipe.Component{Shape: shapes.Switch}}

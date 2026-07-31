package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// select_ binds Select's shape to the accessor calls select.gsx authors.
// The trailing underscore is componentIdentifier's suffix for "select"
// colliding with the Go keyword — the CSS/registry identity is unaffected
// (gsxui-recipe-select, data-gsxui-slot-select).
var select_ = select_Recipe{c: recipe.Component{Shape: shapes.Select}}

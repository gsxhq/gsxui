package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// nativeSelect binds NativeSelect's shape to the accessor calls
// native-select.gsx authors. The receiver name is the camel derivation of
// the hyphenated registry name "native-select" — see
// internal/stylegen/identifier.go.
var nativeSelect = nativeSelectRecipe{c: recipe.Component{Shape: shapes.NativeSelect}}

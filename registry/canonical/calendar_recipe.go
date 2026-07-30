package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// calendar binds Calendar's shape to the accessor calls calendar.gsx authors.
var calendar = calendarRecipe{c: recipe.Component{Shape: shapes.Calendar}}

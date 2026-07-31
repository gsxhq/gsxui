package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// breadcrumb binds Breadcrumb's shape to the accessor calls breadcrumb.gsx
// authors.
var breadcrumb = breadcrumbRecipe{c: recipe.Component{Shape: shapes.Breadcrumb}}

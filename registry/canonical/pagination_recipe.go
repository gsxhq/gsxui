package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// pagination binds Pagination's shape to the accessor calls pagination.gsx
// authors. There is no Link() accessor — see shapes.Pagination's doc comment
// for why PaginationLink's own rules stay behind in default.css instead.
var pagination = paginationRecipe{c: recipe.Component{Shape: shapes.Pagination}}

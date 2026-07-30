package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// carousel binds Carousel's shape to the accessor calls carousel.gsx authors.
var carousel = carouselRecipe{c: recipe.Component{Shape: shapes.Carousel}}

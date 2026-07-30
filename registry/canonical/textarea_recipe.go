package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// textarea binds Textarea's shape to the accessor calls textarea.gsx authors.
var textarea = textareaRecipe{c: recipe.Component{Shape: shapes.Textarea}}

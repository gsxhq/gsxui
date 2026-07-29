// Package canonical holds the structural component sources and the shapes they
// declare. It is compiled so that it type-checks and so that style-independent
// structure and behavior tests can run against the authoritative source, but it
// is never shipped: consumers receive the generated (canonical x recipe) output
// in package ui. Nothing outside internal/stylegen may import it.
package canonical

import (
	"maps"

	"github.com/gsxhq/gsxui/internal/recipe"
)

var shapes = map[string]recipe.Shape{
	buttonShape.Component: buttonShape,
}

// Shapes returns every declared component shape, keyed by component name.
// internal/stylegen reads this instead of parsing Go declarations as data.
func Shapes() map[string]recipe.Shape {
	out := make(map[string]recipe.Shape, len(shapes))
	maps.Copy(out, shapes)
	return out
}

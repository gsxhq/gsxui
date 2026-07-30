// Package shapes holds every component's style interface as pure data.
//
// It is a leaf package on purpose. internal/stylegen imports it to read shapes
// and to GENERATE the per-slot accessors in registry/canonical — which cannot
// compile until those accessors exist. Keeping shapes here means the generator
// always has something that compiles to read from.
//
// Nothing here may import registry/canonical, and internal/recipe must never
// import this package.
package shapes

import (
	"maps"

	"github.com/gsxhq/gsxui/internal/recipe"
)

var all = map[string]recipe.Shape{
	Button.Component:   Button,
	Card.Component:     Card,
	Badge.Component:    Badge,
	Alert.Component:    Alert,
	Skeleton.Component: Skeleton,
	Spinner.Component:  Spinner,
	Progress.Component: Progress,
	Avatar.Component:   Avatar,
}

// All returns every declared component shape, keyed by component name.
func All() map[string]recipe.Shape {
	out := make(map[string]recipe.Shape, len(all))
	maps.Copy(out, all)
	return out
}

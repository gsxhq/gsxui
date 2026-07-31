package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Carousel is the native-scroll-snap port of shadcn/ui's embla-backed
// Carousel. Six styled slots; four of them vary over one dimension,
// `orientation` (horizontal|vertical), which is a real author-facing param on
// every part (this port has no Radix context, so orientation is passed
// explicitly — see registry/canonical/carousel.gsx) and appears as
// data-orientation on the element, exactly the axis default.css already keyed
// its rules on.
//
// track and item declare Base: false: default.css gives them no
// orientation-independent rule at all, only the two per-orientation ones.
// Their orientation-independent mechanics (display:flex, scroll-snap-align,
// flex basis, overflow/scroll-snap-type on content) live in
// assets/css/foundation.css's @layer components as FOUNDATION, not style —
// they are what makes the native-scroll algorithm work rather than what makes
// it look like nova, so they do not migrate and are not slots' base rules.
//
// carousel-control-label carries no rule in any style — its visually-hidden
// treatment is a foundation @layer base rule shared with every other
// screen-reader label in the catalogue — so it declares no slot, the same
// call select-bridge/select-group already made.
//
// previous/next also carry `size-8 rounded-full`, which lived in
// assets/css/styles/default.css's @layer utilities escape hatch rather than in
// default/carousel.css: both parts compose <Button variant="outline"
// size="icon">, and once Button migrated, a components-layer rule could no
// longer beat Button's own resolved rounded-lg. Migrating Carousel retires
// that escape hatch — the round shape is now a plain utility on the same
// element as Button's, appended after it, so merge.Merge drops rounded-lg.
// The horizontal press token
// (active:not-aria-[haspopup]:translate-y-[calc(1px_-_50%)]) is the same
// mechanism: it repeats Button's own modifier chain so tailwind-merge replaces
// Button's translate-y-px rather than stacking with it, preserving both the
// -50% centering and the 1px press dip.
var Carousel = recipe.Shape{
	Component: "carousel",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "content", Base: true},
		{Name: "track", Dimensions: []recipe.Dimension{carouselOrientation}},
		{Name: "item", Dimensions: []recipe.Dimension{carouselOrientation}},
		{Name: "previous", Base: true, Dimensions: []recipe.Dimension{carouselOrientation}},
		{Name: "next", Base: true, Dimensions: []recipe.Dimension{carouselOrientation}},
	},
}

var carouselOrientation = recipe.Dimension{
	Name:    "orientation",
	Default: "horizontal",
	Values:  []string{"horizontal", "vertical"},
}

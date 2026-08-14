package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Chart is the shadcn/ui chart family's server-rendered container
// (registry/new-york-v4/ui/chart.tsx's ChartContainer). It has no
// dimensions: the container's own class is a fixed shape. Per-series color
// is not a shape dimension either: it travels as CSS custom properties
// (--color-<key>) generated at render time by Chart's own styleBlock,
// scoped by the container's data-chart attribute, never as a class variant
// a style pack could enumerate.
//
// The root and legend slots are declared here. Task 3 (cartesian model
// builder) added ChartLegend, which renders a real
// data-gsxui-slot-chart-legend element server-side (see chartWriteLegend
// in chart.gsx), so its slot lands now. ChartTooltip, added in the same
// task, stays data-only — it registers tooltip config into the chart
// model for the client renderer, it renders no element of its own — so
// "tooltip" is still NOT declared here: stylegen's checkShapeCoverage
// rejects a declared slot no canonical component renders yet ("shape
// declares slot %q but the component never renders it" —
// internal/stylegen/generate.go), by design (see that check's own doc
// comment on why a declared-but-dead slot is a hard error, not a
// warning). Whichever later task renders real tooltip markup adds
// "tooltip" here (and its contract entry, and its pack CSS) alongside the
// component that makes it real, the same way "legend" arrives now.
var Chart = recipe.Shape{
	Component: "chart",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "legend", Base: true},
	},
}

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
// Only the root slot is declared here. The plan's own shape sketch also
// names "tooltip" and "legend" parts (the recharts-compat tooltip card and
// the legend row), but stylegen's checkShapeCoverage rejects a declared
// slot no canonical component renders yet ("shape declares slot %q but the
// component never renders it" — internal/stylegen/generate.go) by design:
// see that check's own doc comment on why a declared-but-dead slot is
// deliberately a hard error, not a warning. ChartTooltip and ChartLegend —
// the components that will render data-gsxui-slot-chart-tooltip and
// data-gsxui-slot-chart-legend — land in a later task; that task adds
// those two slots here (and their contract entries, and their pack CSS)
// alongside the components that make them real, the same way every other
// shape in this package grew its slots alongside their consumers.
var Chart = recipe.Shape{
	Component: "chart",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
	},
}

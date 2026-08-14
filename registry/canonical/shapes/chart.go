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
// The root and legend slots were declared first. Task 3 (cartesian model
// builder) added ChartLegend, which renders a real
// data-gsxui-slot-chart-legend element server-side (see chartWriteLegend
// in chart.gsx), so its slot landed then. ChartTooltip itself still
// registers only data (tooltip config into the chart model for the client
// renderer) — it renders no element of its own — but Task 6 adds a
// SEPARATE component, ChartTooltipTemplate, that renders the tooltip's
// hidden server-authored chrome (a <template> ui/chart.render.js clones on
// first hover) whenever a ChartTooltip is registered, so "tooltip" lands
// as a real slot now: its class={ chart.Tooltip() } accessor call sits in
// that template's real markup, satisfying stylegen's checkShapeCoverage
// ("shape declares slot %q but the component never renders it" —
// internal/stylegen/generate.go — no longer applies once something
// renders it).
var Chart = recipe.Shape{
	Component: "chart",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "legend", Base: true},
		{Name: "tooltip", Base: true},
	},
}

// Package chart holds the site's example gsx components for ui/chart.
package chart

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

// Basic renders a two-series bar chart: Chart's container, BarChart's
// model builder, ChartTooltip's hover chrome and ChartLegend's
// server-rendered legend row together — this example keeps the style
// contract's examples-coverage gate satisfied for the chart-legend and
// chart-tooltip slots ChartLegend and ChartTooltip's own
// ChartTooltipTemplate render (see registry/canonical/shapes/chart.go).
component Basic() {
	{{
		cfg := ui.ChartConfig{
			{Key: "desktop", Label: "Desktop", Color: "var(--chart-1)"},
			{Key: "mobile", Label: "Mobile", Color: "var(--chart-2)"},
		}
		data := []ui.ChartDatum{
			{"month": "Jan", "desktop": 186.0, "mobile": 80.0},
			{"month": "Feb", "desktop": 305.0, "mobile": 200.0},
			{"month": "Mar", "desktop": 237.0, "mobile": 120.0},
		}
	}}
	<ui.Chart config={cfg}>
		{ ui.BarChart(data, nil, gsx.Fragment(
			ui.ChartCartesianGrid(nil),
			ui.ChartXAxis("month", nil),
			ui.ChartTooltip(nil),
			ui.ChartBar("desktop", nil),
			ui.ChartBar("mobile", nil),
			ui.ChartLegend(nil),
		), nil) }
	</ui.Chart>
}

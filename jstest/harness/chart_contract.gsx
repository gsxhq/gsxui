package main

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

// ChartContractFixture exercises ui.Chart's tooltip/legend/dark-mode/swap
// behavior (Task 6's spec suite, jstest/specs/chart.spec.ts) without
// touching the public docs example (site/examples/chart/basic.gsx), which
// intentionally keeps to the plain Color mechanism most callers reach for
// first. desktop's ChartSeriesTheme gives its bar a genuinely different
// computed fill under .dark: gsxui's shared placeholder theme
// (assets/css/themes/default.css) keeps every --chart-N token identical
// across both color schemes on purpose (it is a neutral grayscale ramp,
// not yet a themed palette), so a Color: "var(--chart-N)" series (like
// mobile here) can never exercise a real dark-mode flip — a per-scheme
// ChartSeriesTheme is the one config shape that can, and is exactly the
// mechanism Chart.styleBlock (registry/canonical/chart.gsx) exists to
// serve.
component ChartContractFixture() {
	{{
		cfg := ui.ChartConfig{
			{Key: "desktop", Label: "Desktop", Theme: &ui.ChartSeriesTheme{Light: "#2563eb", Dark: "#f97316"}},
			{Key: "mobile", Label: "Mobile", Color: "var(--chart-2)"},
		}
		data := []ui.ChartDatum{
			{"month": "Jan", "desktop": 186.0, "mobile": 80.0},
			{"month": "Feb", "desktop": 305.0, "mobile": 200.0},
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

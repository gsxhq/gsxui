// Package chart holds the site's example gsx components for ui/chart.
package chart

import "github.com/gsxhq/gsxui/ui"

// Basic renders the Chart container with a two-series config around a
// placeholder body. Chart ships its container in this task — the data
// model builder and client renderer (BarChart/AreaChart and friends) land
// in later tasks — so this example exists to keep the style contract's
// examples-coverage gate satisfied from Chart's first commit, the same way
// every other registered component has a real example from birth.
component Basic() {
	{{
		cfg := ui.ChartConfig{
			{Key: "desktop", Label: "Desktop", Color: "var(--chart-1)"},
			{Key: "mobile", Label: "Mobile", Color: "var(--chart-2)"},
		}
	}}
	<ui.Chart config={cfg}>
		<div>Chart placeholder</div>
	</ui.Chart>
}

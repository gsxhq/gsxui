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
// serve. desktop also carries an Icon: ChartSeries.Icon replaces the
// tooltip row's color-swatch indicator with raw SVG markup
// (chart.render.js's tooltipHTML, the s.icon branch), sized entirely by
// the tooltip row's own SVG-child sizing utilities (see
// ui.ChartTooltipTemplate's own doc comment for exactly which ones) — an
// icon with no intrinsic width/height of its own (viewBox only) is
// jstest/specs/chart.spec.ts's probe for whether those utilities actually
// compiled. This comment deliberately does NOT spell out the utility
// literal: this .gsx file is itself inside full.css's harness @source
// glob, so writing it here would compile the class from the comment text
// alone and mask a real regression.
component ChartContractFixture() {
	{{
		cfg := ui.ChartConfig{
			{
				Key: "desktop", Label: "Desktop",
				Theme: &ui.ChartSeriesTheme{Light: "#2563eb", Dark: "#f97316"},
				Icon:  gsx.Raw(`<svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="10"/></svg>`),
			},
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

// ChartNarrowFlexFixture reproduces the flex-column runaway-height defect
// jstest/specs/chart.spec.ts's own "chart height stays..." spec pins: a
// narrow (340px) flex column parent, ui.Chart mounted with NO caller class
// override (the component's own default classes only, chart.Root() /
// .gsxui-recipe-chart in registry/styles/<style>/chart.css). Every animation
// frame, ui/chart.render.js's renderCartesian re-measures panel.clientHeight
// and draws the SVG at that height; Chart's own div is a flex item of this
// column, and a flex item's default min-height:auto lets its own drawn
// content override an aspect-ratio-derived height, so each frame's
// slightly-taller draw feeds the next frame's measurement — a real runaway
// (190px growing past 900px over the ~400ms bar entrance animation) in any
// sufficiently narrow flex parent. site/examples/chart's own wide,
// unconstrained layout never triggers this (aspect-video already gives it
// more height than its content could ever need), which is why this fixture
// exists: nothing else in the harness reproduces a narrow flex parent.
component ChartNarrowFlexFixture() {
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
	<div style="width: 340px" class="flex flex-col gap-4" data-chart-narrow-flex>
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
	</div>
}

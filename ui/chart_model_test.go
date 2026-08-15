package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

// cfg2Series is the two-series ChartConfig every cartesian model pin test
// shares: neither series carries a Color/Label, so styleBlock stays silent
// (no <style> tag to skip over) and ChartConfig.label falls through to each
// series' own (empty) Label field rather than the key, exercising the same
// path a real, labelled config would take.
func cfg2Series(t *testing.T) ui.ChartConfig {
	t.Helper()
	return ui.ChartConfig{
		{Key: "desktop"},
		{Key: "mobile"},
	}
}

// extractModelJSON returns the contents of the <script
// type="application/json" data-gsxui-chart-model> tag in got, failing the
// test if the tag is missing.
func extractModelJSON(t *testing.T, got string) string {
	t.Helper()
	const marker = `<script type="application/json" data-gsxui-chart-model>`
	start := strings.Index(got, marker)
	if start < 0 {
		t.Fatalf("render missing %q\nin: %s", marker, got)
	}
	start += len(marker)
	end := strings.Index(got[start:], "</script>")
	if end < 0 {
		t.Fatalf("render missing closing </script> after model marker\nin: %s", got)
	}
	return got[start : start+end]
}

func TestBarChartModelPinned(t *testing.T) {
	data := []ui.ChartDatum{{"month": "Jan", "desktop": 186.0}, {"month": "Feb", "desktop": 305.0}}
	got := render(t, ui.Chart(cfg2Series(t),
		ui.BarChart(data, 0, 0, 0, 0,
			gsx.Fragment(
				ui.ChartXAxis("month", false, false, false, 0, 0),
				ui.ChartBar("desktop", "", "", nil, 0),
			), nil), nil))
	model := extractModelJSON(t, got)
	// Pin the FULL canonical JSON: a byte-faithful port of templui's own
	// flat Model (chart.templ:1645) — field order, json tags and
	// omitempty-or-not all match upstream, so Task 5's adapted chart.js
	// (which reads these same flat fields) needs no rewriting. marginTop/
	// Right/Bottom/Left default to 5 (no margin registered, chartMarginOf's
	// own zero-sentinel default); xAxisHeight defaults to 30 and
	// minTickGap to 5 once an axis is registered; categoryGap 0.1 is bar's
	// own unconditional Recharts default; cursor and tooltip (empty
	// object) are always present, no ChartTooltip was registered here.
	want := `{"kind":"bar","marginTop":5,"marginRight":5,"marginBottom":5,"marginLeft":5,"xAxisHeight":30,"minTickGap":5,"categoryGap":0.1,"cursor":false,"labels":["Jan","Feb"],"tooltipLabels":["Jan","Feb"],"series":[{"key":"desktop","label":"","color":"var(--color-desktop)","values":[186,305]}],"tooltip":{}}`
	if model != want {
		t.Errorf("model drift\n got: %s\nwant: %s", model, want)
	}
	if strings.Contains(model, "chart-1") || strings.Contains(model, `"id"`) {
		t.Errorf("model must not carry the volatile data-chart id: %s", model)
	}
}

// TestBarChartModelCellsFromFillColumn mirrors
// TestRadialBarChartModelCellsFromFillColumn for cartesian Bar: upstream's
// hasFillColumn fallback (chart.templ:1242, an `else if` alongside Bar's
// own <Cell>-child form we don't port) applies to Bar too, not just
// radial-bar — a "fill" column on the raw rows names each bar's own color
// with no <Cell> child needed, identical data-access pattern.
func TestBarChartModelCellsFromFillColumn(t *testing.T) {
	data := []ui.ChartDatum{
		{"month": "Jan", "desktop": 186.0, "fill": "var(--color-chrome)"},
		{"month": "Feb", "desktop": 305.0, "fill": "var(--color-safari)"},
	}
	got := render(t, ui.Chart(cfg2Series(t),
		ui.BarChart(data, 0, 0, 0, 0,
			gsx.Fragment(
				ui.ChartXAxis("month", false, false, false, 0, 0),
				ui.ChartBar("desktop", "", "", nil, 0),
			), nil), nil))
	model := extractModelJSON(t, got)
	want := `{"kind":"bar","marginTop":5,"marginRight":5,"marginBottom":5,"marginLeft":5,"xAxisHeight":30,"minTickGap":5,"categoryGap":0.1,"cursor":false,"labels":["Jan","Feb"],"tooltipLabels":["Jan","Feb"],"series":[{"key":"desktop","label":"","color":"var(--color-desktop)","values":[186,305],"cells":["var(--color-chrome)","var(--color-safari)"]}],"tooltip":{}}`
	if model != want {
		t.Errorf("model drift\n got: %s\nwant: %s", model, want)
	}
}

func TestLineChartModelPinned(t *testing.T) {
	data := []ui.ChartDatum{{"month": "Jan", "desktop": 186.0}, {"month": "Feb", "desktop": 305.0}}
	got := render(t, ui.Chart(cfg2Series(t),
		ui.LineChart(data, 0, 0, 0, 0,
			gsx.Fragment(
				ui.ChartXAxis("month", false, false, false, 0, 0),
				ui.ChartLine("desktop", "", "", 0, false, 0, 0, "", 0, 0),
			), nil), nil))
	model := extractModelJSON(t, got)
	// Same shape as the bar pin, minus categoryGap: templui's buildModel
	// only sets that field inside its "bar" branch.
	want := `{"kind":"line","marginTop":5,"marginRight":5,"marginBottom":5,"marginLeft":5,"xAxisHeight":30,"minTickGap":5,"cursor":false,"labels":["Jan","Feb"],"tooltipLabels":["Jan","Feb"],"series":[{"key":"desktop","label":"","color":"var(--color-desktop)","values":[186,305]}],"tooltip":{}}`
	if model != want {
		t.Errorf("model drift\n got: %s\nwant: %s", model, want)
	}
}

func TestAreaChartModelPinned(t *testing.T) {
	data := []ui.ChartDatum{{"month": "Jan", "desktop": 186.0}, {"month": "Feb", "desktop": 305.0}}
	got := render(t, ui.Chart(cfg2Series(t),
		ui.AreaChart(data, 0, 0, 0, 0, "",
			gsx.Fragment(
				ui.ChartXAxis("month", false, false, false, 0, 0),
				ui.ChartArea("desktop", "", "", "", "", 0),
			), nil), nil))
	model := extractModelJSON(t, got)
	want := `{"kind":"area","marginTop":5,"marginRight":5,"marginBottom":5,"marginLeft":5,"xAxisHeight":30,"minTickGap":5,"cursor":false,"labels":["Jan","Feb"],"tooltipLabels":["Jan","Feb"],"series":[{"key":"desktop","label":"","color":"var(--color-desktop)","values":[186,305]}],"tooltip":{}}`
	if model != want {
		t.Errorf("model drift\n got: %s\nwant: %s", model, want)
	}
}

// TestAreaChartMarginPartialOverrideMatchesUpstream pins the margin
// semantics review flagged as a parity regression: upstream (Recharts'
// own margin prop, ported byte-faithful by templui and, before this pin,
// by chartMarginOf) replaces its built-in 5/5/5/5 default AS A SET the
// moment a caller supplies a margin object at all — an unset side inside
// that object stays a literal 0, it is never re-defaulted to 5
// independently. Passing only marginLeft/marginRight (site/examples/chart/
// area.gsx's and line.gsx's own shape) must leave marginTop/marginBottom
// at 0, not silently upgrade them to 5.
func TestAreaChartMarginPartialOverrideMatchesUpstream(t *testing.T) {
	data := []ui.ChartDatum{{"month": "Jan", "desktop": 186.0}}
	got := render(t, ui.Chart(cfg2Series(t),
		ui.AreaChart(data, 0, 12, 0, 12, "",
			gsx.Fragment(
				ui.ChartXAxis("month", false, false, false, 0, 0),
				ui.ChartArea("desktop", "", "", "", "", 0),
			), nil), nil))
	model := extractModelJSON(t, got)
	want := `{"kind":"area","marginTop":0,"marginRight":12,"marginBottom":0,"marginLeft":12,"xAxisHeight":30,"minTickGap":5,"cursor":false,"labels":["Jan"],"tooltipLabels":["Jan"],"series":[{"key":"desktop","label":"","color":"var(--color-desktop)","values":[186]}],"tooltip":{}}`
	if model != want {
		t.Errorf("model drift\n got: %s\nwant: %s", model, want)
	}
}

// TestAreaChartMarginAllZeroDefaultsToFive is the bare-tag twin of the
// partial-override pin above: every margin param left at its Go zero value
// (the shape site/examples/chart/basic.gsx and every other margin-less
// demo call use) must still get Recharts' own 5/5/5/5 default as a whole
// set, not 0/0/0/0.
func TestAreaChartMarginAllZeroDefaultsToFive(t *testing.T) {
	data := []ui.ChartDatum{{"month": "Jan", "desktop": 186.0}}
	got := render(t, ui.Chart(cfg2Series(t),
		ui.AreaChart(data, 0, 0, 0, 0, "",
			gsx.Fragment(
				ui.ChartXAxis("month", false, false, false, 0, 0),
				ui.ChartArea("desktop", "", "", "", "", 0),
			), nil), nil))
	model := extractModelJSON(t, got)
	want := `{"kind":"area","marginTop":5,"marginRight":5,"marginBottom":5,"marginLeft":5,"xAxisHeight":30,"minTickGap":5,"cursor":false,"labels":["Jan"],"tooltipLabels":["Jan"],"series":[{"key":"desktop","label":"","color":"var(--color-desktop)","values":[186]}],"tooltip":{}}`
	if model != want {
		t.Errorf("model drift\n got: %s\nwant: %s", model, want)
	}
}

func TestPieChartModelPinned(t *testing.T) {
	data := []ui.ChartDatum{{"browser": "desktop", "visitors": 186.0}, {"browser": "mobile", "visitors": 305.0}}
	got := render(t, ui.Chart(cfg2Series(t),
		ui.PieChart(
			ui.ChartPie(data, "visitors", "browser", 0, 0, 0, "", false, "", false), nil), nil))
	model := extractModelJSON(t, got)
	// Pie carries its own data on the ChartPie child (Recharts' own Pie
	// takes its own data prop, unlike Bar/Line/Area reading the root's
	// rows), so the root itself contributes no margin/labels/series at
	// all — buildPieModel never touches those top-level fields upstream,
	// which is why marginTop/Right/Bottom/Left stay 0 (PieChart takes no
	// margin params at all, unlike every other kind) and labels/series
	// marshal as literal null (both lack omitempty and are never
	// initialized for this kind). labelLine is false, not the shipped
	// bar/line/area pins' own true: this task flips ChartPie's LabelLine
	// default from Recharts' true to gsxui's own false convention (`##
	// chart` ADAPT), and this call leaves it at that default — a legitimate
	// pin change, not drift (see this task's report for the byte diff).
	want := `{"kind":"pie","marginTop":0,"marginRight":0,"marginBottom":0,"marginLeft":0,"cursor":false,"labels":null,"series":null,"tooltip":{},"pies":[{"key":"visitors","seriesLabel":"visitors","values":[186,305],"labels":["",""],"tooltipNames":["",""],"nameKey":"browser","nameValues":["desktop","mobile"],"colors":["var(--color-desktop)","var(--color-mobile)"]}]}`
	if model != want {
		t.Errorf("model drift\n got: %s\nwant: %s", model, want)
	}
}

func TestRadarChartModelPinned(t *testing.T) {
	data := []ui.ChartDatum{{"month": "Jan", "desktop": 186.0}, {"month": "Feb", "desktop": 305.0}}
	got := render(t, ui.Chart(cfg2Series(t),
		ui.RadarChart(data, 0, 0, 0, 0,
			gsx.Fragment(
				ui.ChartPolarAngleAxis("month"),
				ui.ChartRadar("desktop", "", 0, "", 0, false, 0, "", 0),
			), nil), nil))
	model := extractModelJSON(t, got)
	want := `{"kind":"radar","marginTop":5,"marginRight":5,"marginBottom":5,"marginLeft":5,"cursor":false,"labels":["Jan","Feb"],"tooltipLabels":["Jan","Feb"],"series":[{"key":"desktop","label":"","color":"var(--color-desktop)","values":[186,305]}],"tooltip":{},"polar":{"radialLines":true,"hasAngleAxis":true}}`
	if model != want {
		t.Errorf("model drift\n got: %s\nwant: %s", model, want)
	}
}

func TestRadialBarChartModelPinned(t *testing.T) {
	data := []ui.ChartDatum{{"month": "Jan", "desktop": 186.0}, {"month": "Feb", "desktop": 305.0}}
	got := render(t, ui.Chart(cfg2Series(t),
		ui.RadialBarChart(data, 0, 0, 0, 0, 0, 0, 0, 0,
			ui.ChartRadialBar("desktop", "", false, 0, "", ""), nil), nil))
	model := extractModelJSON(t, got)
	// Radial's labels come from the row index, not a registered axis key
	// (there is no polar-angle axis on a radial bar chart), and it never
	// sets the top-level tooltipLabels field upstream either.
	want := `{"kind":"radial","marginTop":5,"marginRight":5,"marginBottom":5,"marginLeft":5,"cursor":false,"labels":["0","1"],"series":[{"key":"desktop","label":"","color":"var(--color-desktop)","values":[186,305],"tooltipNames":["",""]}],"tooltip":{},"radial":{"startAngle":0,"endAngle":360}}`
	if model != want {
		t.Errorf("model drift\n got: %s\nwant: %s", model, want)
	}
}

// TestRadialBarChartModelCellsFromFillColumn pins hasFillColumn's per-row
// fallback (chart.templ's own hasFillColumn/s.Cells wiring in
// buildRadialModel): a "fill" column on the raw rows names each segment's
// own color, the same pattern ChartPie's own per-slice Colors already
// uses — no <Cell> child is needed, this is pure data access.
func TestRadialBarChartModelCellsFromFillColumn(t *testing.T) {
	data := []ui.ChartDatum{
		{"month": "Jan", "desktop": 186.0, "fill": "var(--color-chrome)"},
		{"month": "Feb", "desktop": 305.0, "fill": "var(--color-safari)"},
	}
	got := render(t, ui.Chart(cfg2Series(t),
		ui.RadialBarChart(data, 0, 0, 0, 0, 0, 0, 0, 0,
			ui.ChartRadialBar("desktop", "", false, 0, "", ""), nil), nil))
	model := extractModelJSON(t, got)
	want := `{"kind":"radial","marginTop":5,"marginRight":5,"marginBottom":5,"marginLeft":5,"cursor":false,"labels":["0","1"],"series":[{"key":"desktop","label":"","color":"var(--color-desktop)","values":[186,305],"cells":["var(--color-chrome)","var(--color-safari)"],"tooltipNames":["",""]}],"tooltip":{},"radial":{"startAngle":0,"endAngle":360}}`
	if model != want {
		t.Errorf("model drift\n got: %s\nwant: %s", model, want)
	}
}

// TestChartXAxisBareTagMatchesExplicitFalse proves the API-surface refactor
// changed nothing about the DEFAULT model output for ChartXAxis: a bare
// `<ui.ChartXAxis key="month"/>` (every bool/float64 param left at its Go
// zero value) yields byte-identical xAxis-derived fields to the old
// explicit-false options path this pin already covers above
// (TestBarChartModelPinned's own xAxisHeight/minTickGap/xTickLine/
// xAxisLine/tickMargin fields) — required by this task's brief (#5).
func TestChartXAxisBareTagMatchesExplicitFalse(t *testing.T) {
	data := []ui.ChartDatum{{"month": "Jan", "desktop": 186.0}, {"month": "Feb", "desktop": 305.0}}
	bare := render(t, ui.Chart(cfg2Series(t),
		ui.BarChart(data, 0, 0, 0, 0,
			gsx.Fragment(
				ui.ChartXAxis("month", false, false, false, 0, 0),
				ui.ChartBar("desktop", "", "", nil, 0),
			), nil), nil))
	explicit := render(t, ui.Chart(cfg2Series(t),
		ui.BarChart(data, 0, 0, 0, 0,
			gsx.Fragment(
				ui.ChartXAxis("month", false /* hide */, false /* tickLine */, false /* axisLine */, 0 /* tickMargin */, 0 /* minTickGap */),
				ui.ChartBar("desktop", "", "", nil, 0),
			), nil), nil))
	bareModel := extractModelJSON(t, bare)
	explicitModel := extractModelJSON(t, explicit)
	if bareModel != explicitModel {
		t.Errorf("bare ChartXAxis diverged from explicit-false ChartXAxis\n bare: %s\nexplicit: %s", bareModel, explicitModel)
	}
}

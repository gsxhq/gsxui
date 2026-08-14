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
		ui.BarChart(data, nil,
			gsx.Fragment(
				ui.ChartXAxis("month", nil),
				ui.ChartBar("desktop", nil),
			), nil), nil))
	model := extractModelJSON(t, got)
	// Pin the FULL canonical JSON: a byte-faithful port of templui's own
	// flat Model (chart.templ:1645) — field order, json tags and
	// omitempty-or-not all match upstream, so Task 5's adapted chart.js
	// (which reads these same flat fields) needs no rewriting. marginTop/
	// Right/Bottom/Left default to 5 (no ChartMargin registered);
	// xAxisHeight defaults to 30 and minTickGap to 5 once an axis is
	// registered; categoryGap 0.1 is bar's own unconditional Recharts
	// default; cursor and tooltip (empty object) are always present, no
	// ChartTooltip was registered here.
	want := `{"kind":"bar","marginTop":5,"marginRight":5,"marginBottom":5,"marginLeft":5,"xAxisHeight":30,"minTickGap":5,"categoryGap":0.1,"cursor":false,"labels":["Jan","Feb"],"tooltipLabels":["Jan","Feb"],"series":[{"key":"desktop","label":"","color":"var(--color-desktop)","values":[186,305]}],"tooltip":{}}`
	if model != want {
		t.Errorf("model drift\n got: %s\nwant: %s", model, want)
	}
	if strings.Contains(model, "chart-1") || strings.Contains(model, `"id"`) {
		t.Errorf("model must not carry the volatile data-chart id: %s", model)
	}
}

func TestLineChartModelPinned(t *testing.T) {
	data := []ui.ChartDatum{{"month": "Jan", "desktop": 186.0}, {"month": "Feb", "desktop": 305.0}}
	got := render(t, ui.Chart(cfg2Series(t),
		ui.LineChart(data, nil,
			gsx.Fragment(
				ui.ChartXAxis("month", nil),
				ui.ChartLine("desktop", nil),
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
		ui.AreaChart(data, nil,
			gsx.Fragment(
				ui.ChartXAxis("month", nil),
				ui.ChartArea("desktop", nil),
			), nil), nil))
	model := extractModelJSON(t, got)
	want := `{"kind":"area","marginTop":5,"marginRight":5,"marginBottom":5,"marginLeft":5,"xAxisHeight":30,"minTickGap":5,"cursor":false,"labels":["Jan","Feb"],"tooltipLabels":["Jan","Feb"],"series":[{"key":"desktop","label":"","color":"var(--color-desktop)","values":[186,305]}],"tooltip":{}}`
	if model != want {
		t.Errorf("model drift\n got: %s\nwant: %s", model, want)
	}
}

func TestPieChartModelPinned(t *testing.T) {
	data := []ui.ChartDatum{{"browser": "desktop", "visitors": 186.0}, {"browser": "mobile", "visitors": 305.0}}
	got := render(t, ui.Chart(cfg2Series(t),
		ui.PieChart(
			ui.ChartPie(data, "visitors", &ui.ChartPieOptions{NameKey: "browser"}), nil), nil))
	model := extractModelJSON(t, got)
	// Pie carries its own data on the ChartPie child (Recharts' own Pie
	// takes its own data prop, unlike Bar/Line/Area reading the root's
	// rows), so the root itself contributes no margin/labels/series at
	// all — buildPieModel never touches those top-level fields upstream,
	// which is why marginTop/Right/Bottom/Left stay 0 (no 5/5/5/5 default,
	// unlike every other kind) and labels/series marshal as literal null
	// (both lack omitempty and are never initialized for this kind).
	want := `{"kind":"pie","marginTop":0,"marginRight":0,"marginBottom":0,"marginLeft":0,"cursor":false,"labels":null,"series":null,"tooltip":{},"pies":[{"key":"visitors","seriesLabel":"visitors","values":[186,305],"labels":["",""],"tooltipNames":["",""],"nameKey":"browser","nameValues":["desktop","mobile"],"colors":["var(--color-desktop)","var(--color-mobile)"],"labelLine":true}]}`
	if model != want {
		t.Errorf("model drift\n got: %s\nwant: %s", model, want)
	}
}

func TestRadarChartModelPinned(t *testing.T) {
	data := []ui.ChartDatum{{"month": "Jan", "desktop": 186.0}, {"month": "Feb", "desktop": 305.0}}
	got := render(t, ui.Chart(cfg2Series(t),
		ui.RadarChart(data, nil,
			gsx.Fragment(
				ui.ChartPolarAngleAxis("month"),
				ui.ChartRadar("desktop", nil),
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
		ui.RadialBarChart(data, nil,
			ui.ChartRadialBar("desktop", nil), nil), nil))
	model := extractModelJSON(t, got)
	// Radial's labels come from the row index, not a registered axis key
	// (there is no polar-angle axis on a radial bar chart), and it never
	// sets the top-level tooltipLabels field upstream either.
	want := `{"kind":"radial","marginTop":5,"marginRight":5,"marginBottom":5,"marginLeft":5,"cursor":false,"labels":["0","1"],"series":[{"key":"desktop","label":"","color":"var(--color-desktop)","values":[186,305],"tooltipNames":["",""]}],"tooltip":{},"radial":{"startAngle":0,"endAngle":360}}`
	if model != want {
		t.Errorf("model drift\n got: %s\nwant: %s", model, want)
	}
}

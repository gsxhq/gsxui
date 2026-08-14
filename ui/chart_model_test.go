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
	// Pin the FULL canonical JSON (indent-free, struct-ordered fields, no
	// top-level map — Data rows are the one map[string]any exception,
	// relying on encoding/json's own sorted-key map marshaling).
	want := `{"kind":"bar","data":[{"desktop":186,"month":"Jan"},{"desktop":305,"month":"Feb"}],"series":[{"key":"desktop","color":"var(--color-desktop)"}],"xAxis":{"key":"month","tickLine":false,"axisLine":false,"tickMargin":0,"minTickGap":5,"height":30,"hide":false}}`
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
	want := `{"kind":"line","data":[{"desktop":186,"month":"Jan"},{"desktop":305,"month":"Feb"}],"series":[{"key":"desktop","color":"var(--color-desktop)"}],"xAxis":{"key":"month","tickLine":false,"axisLine":false,"tickMargin":0,"minTickGap":5,"height":30,"hide":false}}`
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
	want := `{"kind":"area","data":[{"desktop":186,"month":"Jan"},{"desktop":305,"month":"Feb"}],"series":[{"key":"desktop","color":"var(--color-desktop)"}],"xAxis":{"key":"month","tickLine":false,"axisLine":false,"tickMargin":0,"minTickGap":5,"height":30,"hide":false}}`
	if model != want {
		t.Errorf("model drift\n got: %s\nwant: %s", model, want)
	}
}

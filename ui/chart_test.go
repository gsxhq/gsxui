package ui_test

import (
	"context"
	"errors"
	"io"
	"regexp"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestChartStyleBlockPrecedence(t *testing.T) {
	cfg := ui.ChartConfig{
		{Key: "desktop", Color: "var(--chart-1)"},
		{Key: "mobile", Color: "red", Theme: &ui.ChartSeriesTheme{Light: "blue", Dark: "green"}},
	}
	got := render(t, ui.Chart(cfg, gsx.Raw("x"), nil))
	// Theme wins over Color; light block unprefixed, dark under .dark.
	for _, want := range []string{
		"--color-desktop: var(--chart-1);",
		"--color-mobile: blue;",  // light: theme wins over red
		"--color-mobile: green;", // dark
		"[data-chart=",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in %s", want, got)
		}
	}
}

func TestChartNoColorsNoStyle(t *testing.T) {
	got := render(t, ui.Chart(ui.ChartConfig{{Key: "a", Label: "A"}}, gsx.Raw("x"), nil))
	if strings.Contains(got, "<style>") {
		t.Errorf("uncolored config must emit no style block: %s", got)
	}
}

// TestChartStyleBlockFiltersHostileKeys pins the trust boundary of the
// <style> block: Key is a plain string and goes through gsx's CSS value
// filter, so a key that tries to close its declaration and inject a rule
// is neutralised (ZgotmplZ) — while Color is gsx.RawCSS (the author's
// vouch) and passes verbatim, which is what lets every shadcn demo's
// var(--chart-N) — a value the filter would otherwise reject — render.
func TestChartStyleBlockFiltersHostileKeys(t *testing.T) {
	cfg := ui.ChartConfig{
		{Key: "x; } body { display:none", Color: "red"},
	}
	got := render(t, ui.Chart(cfg, gsx.Raw("x"), nil))
	if strings.Contains(got, "body { display:none") || strings.Contains(got, "body{display:none") {
		t.Errorf("hostile key reached the style block unfiltered: %s", got)
	}
	if !strings.Contains(got, "ZgotmplZ") {
		t.Errorf("expected the CSS filter's ZgotmplZ substitution in %s", got)
	}
}

// TestChartStyleBlockVouchedColorPassesVerbatim pins that a RawCSS Color
// using functional notation renders unchanged — the value the CSS filter
// rejects from strings, and the value every shadcn chart demo relies on.
func TestChartStyleBlockVouchedColorPassesVerbatim(t *testing.T) {
	cfg := ui.ChartConfig{{Key: "desktop", Color: "var(--chart-1)"}}
	got := render(t, ui.Chart(cfg, gsx.Raw("x"), nil))
	if !strings.Contains(got, "--color-desktop:var(--chart-1)") && !strings.Contains(got, "--color-desktop: var(--chart-1)") {
		t.Errorf("vouched var(--chart-1) must pass verbatim: %s", got)
	}
	if strings.Contains(got, "ZgotmplZ") {
		t.Errorf("a RawCSS color must not be filtered: %s", got)
	}
}

// TestChartLegendRendersIconNodeDirectly pins that a series' Icon reaches
// the server-rendered legend as the gsx.Node it was given, rendered in
// place through gsx's writer — not pre-rendered to an HTML string and
// re-emitted with gsx.Raw. The behavioral difference: an in-place node's
// render error propagates and fails the whole render, whereas the string
// round-trip (chartRenderHTML) swallowed errors into an empty icon.
func TestChartLegendRendersIconNodeDirectly(t *testing.T) {
	boom := errors.New("icon render failed")
	icon := gsx.Func(func(_ context.Context, _ io.Writer) error { return boom })
	cfg := ui.ChartConfig{{Key: "desktop", Label: "Desktop", Color: "var(--chart-1)", Icon: icon}}
	data := []ui.ChartDatum{{"month": "Jan", "desktop": 1.0}}
	node := ui.Chart(cfg,
		ui.BarChart(data, 0, 0, 0, 0, gsx.Fragment(
			ui.ChartBar("desktop", "", "", nil, 0),
			ui.ChartLegend("", "", "", false),
		), nil), nil)
	var sb strings.Builder
	err := node.Render(context.Background(), &sb)
	if !errors.Is(err, boom) {
		t.Errorf("legend must render the Icon node in place so its error propagates; got err=%v, output: %s", err, sb.String())
	}
}

// TestChartLegendStylesRender pins that the legend's positioning and its
// swatch color survive gsx's CSS filtering: both are composed from css
// literals with vouched values. Regression pin — the string-concatenated
// originals were filtered to ZgotmplZ (a legend with no position and a
// colorless swatch), unseen by specs that only asserted presence.
func TestChartLegendStylesRender(t *testing.T) {
	cfg := ui.ChartConfig{{Key: "desktop", Label: "Desktop", Color: "var(--chart-1)"}}
	data := []ui.ChartDatum{{"month": "Jan", "desktop": 1.0}}
	got := render(t, ui.Chart(cfg,
		ui.BarChart(data, 0, 0, 0, 0, gsx.Fragment(
			ui.ChartBar("desktop", "", "", nil, 0),
			ui.ChartLegend("", "", "", false),
		), nil), nil))
	if strings.Contains(got, "ZgotmplZ") {
		t.Errorf("legend styles were filtered: %s", got)
	}
	// The wrapper's style attribute precedes the slot marker on the same
	// tag; the swatch carries the model's resolved series color variable.
	for _, want := range []string{
		`style="position:absolute;left:0;right:0; bottom:5px" data-gsxui-slot-chart-legend`,
		`style="background-color:var(--color-desktop)"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("legend missing %q in: %s", want, got)
		}
	}
}

// TestChartIDsAreContentDerived pins the id contract that survives htmx
// swaps, streaming and multi-node renders: the data-chart id scoping a
// chart's <style> block is a pure function of its config. Two renders of
// the same document are byte-identical (no process counter, no random
// source); charts with distinct configs get distinct ids (their rules
// differ, so scoping matters); charts with identical configs share one id
// (their rules are identical, so sharing is a dedup, not a collision).
func TestChartIDsAreContentDerived(t *testing.T) {
	a := ui.ChartConfig{{Key: "d", Color: "red"}}
	b := ui.ChartConfig{{Key: "d", Color: "blue"}}
	doc := func() gsx.Node {
		return gsx.Fragment(
			ui.Chart(a, gsx.Raw("1"), nil),
			ui.Chart(b, gsx.Raw("2"), nil),
			ui.Chart(a, gsx.Raw("3"), nil),
		)
	}
	first := render(t, doc())
	if second := render(t, doc()); first != second {
		t.Errorf("chart ids must be deterministic across renders\nfirst:  %s\nsecond: %s", first, second)
	}
	ids := regexp.MustCompile(`data-chart="([^"]+)"`).FindAllStringSubmatch(first, -1)
	if len(ids) != 3 {
		t.Fatalf("expected 3 chart ids, got %d in %s", len(ids), first)
	}
	if ids[0][1] == ids[1][1] {
		t.Errorf("distinct configs must get distinct ids: %s", first)
	}
	if ids[0][1] != ids[2][1] {
		t.Errorf("identical configs must share one id (identical rules): %s", first)
	}
}

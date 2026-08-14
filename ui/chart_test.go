package ui_test

import (
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

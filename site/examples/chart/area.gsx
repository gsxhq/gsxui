package chart

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
	uiicon "github.com/gsxhq/gsxui/ui/icon"
)

// Area ports shadcn's chart-area-gradient demo: a stacked, gradient-filled
// AreaChart wrapped in a Card, exercising ChartDefs/ChartLinearGradient
// (the gradient <defs> upstream's own chart-area-interactive demo also
// relies on) alongside the cartesian grid, x axis and tooltip. The
// interactive demo's client-side time-range Select has no gsxui pendant
// (no client state this task models), so this keeps the static six months
// of data every upstream chart-area-*.tsx demo shares instead.
component Area() {
	{{
		cfg := ui.ChartConfig{
			{Key: "desktop", Label: "Desktop", Color: "var(--chart-1)"},
			{Key: "mobile", Label: "Mobile", Color: "var(--chart-2)"},
		}
		data := []ui.ChartDatum{
			{"month": "January", "desktop": 186.0, "mobile": 80.0},
			{"month": "February", "desktop": 305.0, "mobile": 200.0},
			{"month": "March", "desktop": 237.0, "mobile": 120.0},
			{"month": "April", "desktop": 73.0, "mobile": 190.0},
			{"month": "May", "desktop": 209.0, "mobile": 130.0},
			{"month": "June", "desktop": 214.0, "mobile": 140.0},
		}
	}}
	<ui.Card>
		<ui.CardHeader>
			<ui.CardTitle>Area Chart - Gradient</ui.CardTitle>
			<ui.CardDescription>Showing total visitors for the last 6 months</ui.CardDescription>
		</ui.CardHeader>
		<ui.CardContent>
			<ui.Chart config={cfg}>
				{ ui.AreaChart(data, &ui.ChartAreaChartOptions{Margin: &ui.ChartMargin{Left: 12, Right: 12}}, gsx.Fragment(
					ui.ChartCartesianGrid(&ui.ChartCartesianGridOptions{Vertical: ui.ChartBool(false)}),
					ui.ChartXAxis("month", &ui.ChartXAxisOptions{TickMargin: 8}),
					ui.ChartTooltip(&ui.ChartTooltipOptions{Cursor: ui.ChartBool(false)}),
					ui.ChartDefs(gsx.Fragment(
						ui.ChartLinearGradient(
							ui.ChartLinearGradientOptions{ID: "fillDesktop", X1: "0", Y1: "0", X2: "0", Y2: "1"},
							gsx.Raw(`<stop offset="5%" stop-color="var(--color-desktop)" stop-opacity="0.8"></stop><stop offset="95%" stop-color="var(--color-desktop)" stop-opacity="0.1"></stop>`),
						),
						ui.ChartLinearGradient(
							ui.ChartLinearGradientOptions{ID: "fillMobile", X1: "0", Y1: "0", X2: "0", Y2: "1"},
							gsx.Raw(`<stop offset="5%" stop-color="var(--color-mobile)" stop-opacity="0.8"></stop><stop offset="95%" stop-color="var(--color-mobile)" stop-opacity="0.1"></stop>`),
						),
					)),
					ui.ChartArea("mobile", &ui.ChartAreaOptions{Type: ui.ChartCurveNatural, Fill: "url(#fillMobile)", FillOpacity: 0.4, Stroke: "var(--color-mobile)", StackID: "a"}),
					ui.ChartArea("desktop", &ui.ChartAreaOptions{Type: ui.ChartCurveNatural, Fill: "url(#fillDesktop)", FillOpacity: 0.4, Stroke: "var(--color-desktop)", StackID: "a"}),
				), nil) }
			</ui.Chart>
		</ui.CardContent>
		<ui.CardFooter>
			<div class="flex w-full items-start gap-2 text-sm">
				<div class="grid gap-2">
					<div class="flex items-center gap-2 leading-none font-medium">
						Trending up by 5.2% this month
						<uiicon.TrendingUp class="size-4"/>
					</div>
					<div class="flex items-center gap-2 leading-none text-muted-foreground">
						January - June 2024
					</div>
				</div>
			</div>
		</ui.CardFooter>
	</ui.Card>
}

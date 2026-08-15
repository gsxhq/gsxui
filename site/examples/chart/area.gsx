package chart

import (
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
				<ui.AreaChart data={data} marginLeft={12} marginRight={12}>
					<ui.ChartCartesianGrid horizontal/>
					<ui.ChartXAxis key="month" tickMargin={8}/>
					<ui.ChartTooltip/>
					<ui.ChartDefs>
						<ui.ChartLinearGradient id="fillDesktop" x1="0" y1="0" x2="0" y2="1">
							<stop offset="5%" stop-color="var(--color-desktop)" stop-opacity="0.8"/>
							<stop offset="95%" stop-color="var(--color-desktop)" stop-opacity="0.1"/>
						</ui.ChartLinearGradient>
						<ui.ChartLinearGradient id="fillMobile" x1="0" y1="0" x2="0" y2="1">
							<stop offset="5%" stop-color="var(--color-mobile)" stop-opacity="0.8"/>
							<stop offset="95%" stop-color="var(--color-mobile)" stop-opacity="0.1"/>
						</ui.ChartLinearGradient>
					</ui.ChartDefs>
					<ui.ChartArea key="mobile" curve="natural" fill="url(#fillMobile)" fillOpacity={0.4} stroke="var(--color-mobile)" stackId="a"/>
					<ui.ChartArea key="desktop" curve="natural" fill="url(#fillDesktop)" fillOpacity={0.4} stroke="var(--color-desktop)" stackId="a"/>
				</ui.AreaChart>
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

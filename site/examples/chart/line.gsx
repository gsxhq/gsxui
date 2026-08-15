package chart

import (
	"github.com/gsxhq/gsxui/ui"
	uiicon "github.com/gsxhq/gsxui/ui/icon"
)

// Line ports shadcn's chart-line-multiple demo: two monotone lines sharing
// one x axis, cartesian grid and tooltip, wrapped in a Card with a
// trending-up footer.
component Line() {
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
			<ui.CardTitle>Line Chart - Multiple</ui.CardTitle>
			<ui.CardDescription>January - June 2024</ui.CardDescription>
		</ui.CardHeader>
		<ui.CardContent>
			<ui.Chart config={cfg}>
				<ui.LineChart data={data} marginLeft={12} marginRight={12}>
					<ui.ChartCartesianGrid horizontal/>
					<ui.ChartXAxis key="month" tickMargin={8}/>
					<ui.ChartTooltip/>
					<ui.ChartLine key="desktop" curve="monotone" stroke="var(--color-desktop)" strokeWidth={2}/>
					<ui.ChartLine key="mobile" curve="monotone" stroke="var(--color-mobile)" strokeWidth={2}/>
				</ui.LineChart>
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
						Showing total visitors for the last 6 months
					</div>
				</div>
			</div>
		</ui.CardFooter>
	</ui.Card>
}

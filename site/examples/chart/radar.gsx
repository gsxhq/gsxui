package chart

import (
	"github.com/gsxhq/gsxui/ui"
	uiicon "github.com/gsxhq/gsxui/ui/icon"
)

// Radar ports shadcn's chart-radar-default demo: a single-series RadarChart
// with a polar grid and angle axis, wrapped in a Card.
component Radar() {
	{{
		cfg := ui.ChartConfig{
			{Key: "desktop", Label: "Desktop", Color: "var(--chart-1)"},
		}
		data := []ui.ChartDatum{
			{"month": "January", "desktop": 186.0},
			{"month": "February", "desktop": 305.0},
			{"month": "March", "desktop": 237.0},
			{"month": "April", "desktop": 273.0},
			{"month": "May", "desktop": 209.0},
			{"month": "June", "desktop": 214.0},
		}
	}}
	<ui.Card>
		<ui.CardHeader class="items-center pb-4">
			<ui.CardTitle>Radar Chart</ui.CardTitle>
			<ui.CardDescription>Showing total visitors for the last 6 months</ui.CardDescription>
		</ui.CardHeader>
		<ui.CardContent class="pb-0">
			<ui.Chart config={cfg} class="mx-auto aspect-square max-h-[250px]">
				<ui.RadarChart data={data}>
					<ui.ChartTooltip/>
					<ui.ChartPolarAngleAxis key="month"/>
					<ui.ChartPolarGrid radialLines/>
					<ui.ChartRadar key="desktop" fill="var(--color-desktop)" fillOpacity={0.6}/>
				</ui.RadarChart>
			</ui.Chart>
		</ui.CardContent>
		<ui.CardFooter class="flex-col gap-2 text-sm">
			<div class="flex items-center gap-2 leading-none font-medium">
				Trending up by 5.2% this month
				<uiicon.TrendingUp class="size-4"/>
			</div>
			<div class="flex items-center gap-2 leading-none text-muted-foreground">
				January - June 2024
			</div>
		</ui.CardFooter>
	</ui.Card>
}

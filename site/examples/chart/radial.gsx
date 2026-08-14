package chart

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
	uiicon "github.com/gsxhq/gsxui/ui/icon"
)

// Radial ports shadcn's chart-radial-simple demo: five browser-share radial
// bars over a background track, wrapped in a Card.
component Radial() {
	{{
		cfg := ui.ChartConfig{
			{Key: "visitors", Label: "Visitors"},
			{Key: "chrome", Label: "Chrome", Color: "var(--chart-1)"},
			{Key: "safari", Label: "Safari", Color: "var(--chart-2)"},
			{Key: "firefox", Label: "Firefox", Color: "var(--chart-3)"},
			{Key: "edge", Label: "Edge", Color: "var(--chart-4)"},
			{Key: "other", Label: "Other", Color: "var(--chart-5)"},
		}
		data := []ui.ChartDatum{
			{"browser": "chrome", "visitors": 275.0, "fill": "var(--color-chrome)"},
			{"browser": "safari", "visitors": 200.0, "fill": "var(--color-safari)"},
			{"browser": "firefox", "visitors": 187.0, "fill": "var(--color-firefox)"},
			{"browser": "edge", "visitors": 173.0, "fill": "var(--color-edge)"},
			{"browser": "other", "visitors": 90.0, "fill": "var(--color-other)"},
		}
	}}
	<ui.Card class="flex flex-col">
		<ui.CardHeader class="items-center pb-0">
			<ui.CardTitle>Radial Chart</ui.CardTitle>
			<ui.CardDescription>January - June 2024</ui.CardDescription>
		</ui.CardHeader>
		<ui.CardContent class="flex-1 pb-0">
			<ui.Chart config={cfg} class="mx-auto aspect-square max-h-[250px]">
				{ ui.RadialBarChart(data, &ui.ChartRadialBarChartOptions{InnerRadius: 30, OuterRadius: 110}, gsx.Fragment(
					ui.ChartTooltip(&ui.ChartTooltipOptions{Cursor: ui.ChartBool(false), HideLabel: true, NameKey: "browser"}),
					ui.ChartRadialBar("visitors", &ui.ChartRadialBarOptions{Background: true}),
				), nil) }
			</ui.Chart>
		</ui.CardContent>
		<ui.CardFooter class="flex-col gap-2 text-sm">
			<div class="flex items-center gap-2 leading-none font-medium">
				Trending up by 5.2% this month
				<uiicon.TrendingUp class="size-4"/>
			</div>
			<div class="leading-none text-muted-foreground">
				Showing total visitors for the last 6 months
			</div>
		</ui.CardFooter>
	</ui.Card>
}

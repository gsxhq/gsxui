package chart

import (
	"github.com/gsxhq/gsxui/ui"
	uiicon "github.com/gsxhq/gsxui/ui/icon"
)

// Pie ports shadcn's chart-pie-donut demo: a donut PieChart over five
// browser-share slices with a hidden-label tooltip, wrapped in a Card. The
// sibling chart-pie-donut-text demo's center total (a Recharts Label
// render-prop mounted through Pie) has no gsxui pendant — Label isn't
// ported (see ui/chart.gsx's polar-model-builder section header) — so this
// keeps the plain donut instead.
component Pie() {
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
			<ui.CardTitle>Pie Chart - Donut</ui.CardTitle>
			<ui.CardDescription>January - June 2024</ui.CardDescription>
		</ui.CardHeader>
		<ui.CardContent class="flex-1 pb-0">
			<ui.Chart config={cfg} class="mx-auto aspect-square max-h-[250px]">
				<ui.PieChart>
					<ui.ChartTooltip hideLabel/>
					<ui.ChartPie data={data} key="visitors" nameKey="browser" innerRadius={60}/>
				</ui.PieChart>
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

// Package showcase holds the composed demo cards rendered on the site's
// landing page. Unlike the per-component packages next door, these are
// not registered docs examples — they exist to show several ui components
// working together in one realistic piece of UI.
package showcase

import "github.com/gsxhq/gsxui/ui"

// Stand-in portrait for the usage card, same inline-SVG-to-data-URL
// technique as site/examples/avatar.
var showcaseAvatarSVG = []byte("<svg xmlns='http://www.w3.org/2000/svg' width='64' height='64'><rect width='64' height='64' fill='#6d28d9'/><text x='32' y='34' text-anchor='middle' dominant-baseline='central' font-family='sans-serif' font-weight='600' font-size='26' fill='#fff'>AL</text></svg>")

// StatsCard composes Avatar, Badge and Progress into a usage summary.
component StatsCard() {
	<ui.Card>
		<ui.CardHeader>
			<ui.CardTitle>Usage</ui.CardTitle>
			<ui.CardDescription>Your plan resets in 12 days.</ui.CardDescription>
			<ui.CardAction>
				<ui.Badge variant="secondary">Active</ui.Badge>
			</ui.CardAction>
		</ui.CardHeader>
		<ui.CardContent class="flex flex-col gap-4">
			<div class="flex items-center gap-3">
				<ui.Avatar>
					<ui.AvatarImage src={showcaseAvatarSVG |> dataURL("image/svg+xml")} alt="Ada Lovelace"/>
					<ui.AvatarFallback>AL</ui.AvatarFallback>
				</ui.Avatar>
				<div class="flex flex-col">
					<span class="text-sm font-medium">Ada Lovelace</span>
					<span class="text-sm text-muted-foreground">Pro plan</span>
				</div>
			</div>
			<div class="grid gap-2">
				<div class="flex items-center justify-between text-sm">
					<span>Storage</span>
					<span class="text-muted-foreground">72%</span>
				</div>
				<ui.Progress value={72}/>
			</div>
			<div class="grid gap-2">
				<div class="flex items-center justify-between text-sm">
					<span>Bandwidth</span>
					<span class="text-muted-foreground">31%</span>
				</div>
				<ui.Progress value={31}/>
			</div>
		</ui.CardContent>
	</ui.Card>
}

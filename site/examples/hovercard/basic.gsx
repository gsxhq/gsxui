// Package hovercard holds the site's example gsx components for
// ui/hover-card.
package hovercard

import "github.com/gsxhq/gsxui/ui"

// A small square avatar tile, same authoring pattern as
// site/examples/avatar/basic.gsx's own data-URI stand-in portrait — plain
// SVG bytes assembled into a data: URL via the `|> dataURL(mime)` std
// filter gsx's image-sink sanitizer accepts.
var avatarSVG = []byte("<svg xmlns='http://www.w3.org/2000/svg' width='64' height='64'><rect width='64' height='64' fill='#000'/><text x='32' y='34' text-anchor='middle' dominant-baseline='central' font-family='sans-serif' font-weight='600' font-size='22' fill='#fff'>VC</text></svg>")

// Basic mirrors shadcn's own hover-card-demo.tsx (registry/new-york-v4/
// examples/hover-card-demo.tsx): a link-styled @nextjs trigger and a
// profile-preview card (avatar, name, description, joined-date row).
// HoverCardTrigger's <span> wrapper poses none of DialogTrigger/
// TooltipTrigger's button-in-button trap, so — unlike the dropdown/tooltip
// examples, which attach data-gsxui-*-trigger straight onto a Button and
// skip the Trigger component — this composes ui.Button as a real child.
component Basic() {
	<ui.HoverCard>
		<ui.HoverCardTrigger>
			<ui.Button variant="link">@nextjs</ui.Button>
		</ui.HoverCardTrigger>
		<ui.HoverCardContent class="w-80">
			<div class="flex justify-between gap-4">
				<ui.Avatar>
					<ui.AvatarImage src={avatarSVG |> dataURL("image/svg+xml")} alt="@nextjs"/>
					<ui.AvatarFallback>VC</ui.AvatarFallback>
				</ui.Avatar>
				<div class="space-y-1">
					<h4 class="text-sm font-semibold">@nextjs</h4>
					<p class="text-sm">The React Framework – created and maintained by @vercel.</p>
					<div class="text-xs text-muted-foreground">Joined December 2021</div>
				</div>
			</div>
		</ui.HoverCardContent>
	</ui.HoverCard>
}

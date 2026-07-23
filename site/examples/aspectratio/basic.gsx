// Package aspectratio holds the site's example gsx components for
// ui/aspect-ratio.
package aspectratio

import "github.com/gsxhq/gsxui/ui"

// Basic renders a 16/9 box with a border and muted background — a visible
// stand-in for shadcn's docs demo (which fills the ratio box with a real
// <Image>) without pulling in an image asset.
component Basic() {
	<div class="w-full max-w-sm">
		<ui.AspectRatio
			ratio="16 / 9"
			class="flex items-center justify-center rounded-lg border bg-muted text-sm text-muted-foreground"
		>
			16 / 9
		</ui.AspectRatio>
	</div>
}

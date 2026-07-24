// Package popover holds the site's example gsx components for ui/popover.
package popover

import (
	"github.com/gsxhq/gsxui/ui"
)

// Basic mirrors shadcn's own popover-demo.tsx (registry/new-york-v4/
// examples/popover-demo.tsx): an Open-popover button trigger and a w-80
// dimensions form (Label+Input rows), the caller's w-80 winning over
// PopoverContent's own default w-72 via tailwind-merge.
component Basic() {
	<ui.Popover>
		<ui.Button variant="outline" data-gsxui-popover-trigger>Open popover</ui.Button>
		<ui.PopoverContent class="w-80">
			<div class="grid gap-4">
				<div class="space-y-2">
					<h4 class="leading-none font-medium">Dimensions</h4>
					<p class="text-sm text-muted-foreground">Set the dimensions for the layer.</p>
				</div>
				<div class="grid gap-2">
					<div class="grid grid-cols-3 items-center gap-4">
						<ui.Label for="width">Width</ui.Label>
						<ui.Input id="width" value="100%" class="col-span-2 h-8"/>
					</div>
					<div class="grid grid-cols-3 items-center gap-4">
						<ui.Label for="maxWidth">Max. width</ui.Label>
						<ui.Input id="maxWidth" value="300px" class="col-span-2 h-8"/>
					</div>
					<div class="grid grid-cols-3 items-center gap-4">
						<ui.Label for="height">Height</ui.Label>
						<ui.Input id="height" value="25px" class="col-span-2 h-8"/>
					</div>
					<div class="grid grid-cols-3 items-center gap-4">
						<ui.Label for="maxHeight">Max. height</ui.Label>
						<ui.Input id="maxHeight" value="none" class="col-span-2 h-8"/>
					</div>
				</div>
			</div>
		</ui.PopoverContent>
	</ui.Popover>
}

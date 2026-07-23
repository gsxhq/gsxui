// Package progress holds the site's example gsx components for ui/progress.
package progress

import "github.com/gsxhq/gsxui/ui"

// Basic renders a few static Progress values stacked, since v1 ships no
// client JS to animate the value over time.
component Basic() {
	<div class="flex w-full max-w-sm flex-col gap-4">
		<ui.Progress value={25}/>
		<ui.Progress value={60}/>
		<ui.Progress value={90}/>
	</div>
}

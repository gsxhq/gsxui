// Package separator holds the site's example gsx components for
// ui/separator.
package separator

import "github.com/gsxhq/gsxui/ui"

// Basic renders a horizontal Separator between two lines of text.
component Basic() {
	<div class="text-sm">
		<p>Above the line.</p>
		<ui.Separator class="my-4"/>
		<p>Below the line.</p>
	</div>
}

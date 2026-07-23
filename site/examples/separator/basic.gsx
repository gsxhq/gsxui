// Package separator holds the site's example gsx components for
// ui/separator.
package separator

import uiseparator "github.com/gsxhq/gsxui/ui/separator"

// Basic renders a horizontal Separator between two lines of text.
component Basic() {
	<div class="text-sm">
		<p>Above the line.</p>
		<uiseparator.Separator class="my-4"/>
		<p>Below the line.</p>
	</div>
}

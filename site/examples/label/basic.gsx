// Package label holds the site's example gsx components for ui/label.
package label

import (
	"github.com/gsxhq/gsxui/ui"
)

// Basic pairs a Label with an Input via matching for/id.
component Basic() {
	<div class="flex max-w-sm flex-col gap-2">
		<ui.Label for="label-basic-email">Email</ui.Label>
		<ui.Input id="label-basic-email" type="email" placeholder="you@example.com"/>
	</div>
}

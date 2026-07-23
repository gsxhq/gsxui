// Package label holds the site's example gsx components for ui/label.
package label

import (
	uiinput "github.com/gsxhq/gsxui/ui/input"
	uilabel "github.com/gsxhq/gsxui/ui/label"
)

// Basic pairs a Label with an Input via matching for/id.
component Basic() {
	<div class="flex max-w-sm flex-col gap-2">
		<uilabel.Label for="label-basic-email">Email</uilabel.Label>
		<uiinput.Input id="label-basic-email" type="email" placeholder="you@example.com"/>
	</div>
}

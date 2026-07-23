// Package radio holds the site's example gsx components for ui/radio.
package radio

import (
	uilabel "github.com/gsxhq/gsxui/ui/label"
	uiradio "github.com/gsxhq/gsxui/ui/radio"
)

// Basic groups three Radio inputs by a shared name — native radio grouping,
// no wrapper component needed.
component Basic() {
	<div class="flex flex-col gap-2">
		<div class="flex items-center gap-2">
			<uiradio.Radio id="radio-basic-card" name="radio-basic-plan" checked/>
			<uilabel.Label for="radio-basic-card">Card</uilabel.Label>
		</div>
		<div class="flex items-center gap-2">
			<uiradio.Radio id="radio-basic-paypal" name="radio-basic-plan"/>
			<uilabel.Label for="radio-basic-paypal">PayPal</uilabel.Label>
		</div>
	</div>
}

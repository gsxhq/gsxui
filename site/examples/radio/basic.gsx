// Package radio holds the site's example gsx components for ui/radio.
package radio

import (
	"github.com/gsxhq/gsxui/ui"
)

// Basic groups three Radio inputs by a shared name — native radio grouping,
// no wrapper component needed.
component Basic() {
	<div class="flex flex-col gap-2">
		<div class="flex items-center gap-2">
			<ui.Radio id="radio-basic-card" name="radio-basic-plan" checked/>
			<ui.Label for="radio-basic-card">Card</ui.Label>
		</div>
		<div class="flex items-center gap-2">
			<ui.Radio id="radio-basic-paypal" name="radio-basic-plan"/>
			<ui.Label for="radio-basic-paypal">PayPal</ui.Label>
		</div>
	</div>
}

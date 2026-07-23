package radio

import (
	"github.com/gsxhq/gsxui/ui"
)

// States adds a disabled option to the group — disabled is a bare boolean
// attribute forwarded through Radio's { attrs... } spread.
component States() {
	<div class="flex flex-col gap-2">
		<div class="flex items-center gap-2">
			<ui.Radio id="radio-states-monthly" name="radio-states-billing" checked/>
			<ui.Label for="radio-states-monthly">Monthly</ui.Label>
		</div>
		<div class="flex items-center gap-2">
			<ui.Radio id="radio-states-yearly" name="radio-states-billing"/>
			<ui.Label for="radio-states-yearly">Yearly</ui.Label>
		</div>
		<div class="flex items-center gap-2">
			<ui.Radio id="radio-states-lifetime" name="radio-states-billing" disabled/>
			<ui.Label for="radio-states-lifetime">Lifetime (unavailable)</ui.Label>
		</div>
	</div>
}

package radio

import (
	uilabel "github.com/gsxhq/gsxui/ui/label"
	uiradio "github.com/gsxhq/gsxui/ui/radio"
)

// States adds a disabled option to the group — disabled is a bare boolean
// attribute forwarded through Radio's { attrs... } spread.
component States() {
	<div class="flex flex-col gap-2">
		<div class="flex items-center gap-2">
			<uiradio.Radio id="radio-states-monthly" name="radio-states-billing" checked/>
			<uilabel.Label for="radio-states-monthly">Monthly</uilabel.Label>
		</div>
		<div class="flex items-center gap-2">
			<uiradio.Radio id="radio-states-yearly" name="radio-states-billing"/>
			<uilabel.Label for="radio-states-yearly">Yearly</uilabel.Label>
		</div>
		<div class="flex items-center gap-2">
			<uiradio.Radio id="radio-states-lifetime" name="radio-states-billing" disabled/>
			<uilabel.Label for="radio-states-lifetime">Lifetime (unavailable)</uilabel.Label>
		</div>
	</div>
}

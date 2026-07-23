package checkbox

import (
	"github.com/gsxhq/gsxui/ui"
)

// States renders Checkbox unchecked, checked, and disabled+checked —
// checked/disabled are bare boolean attributes forwarded through Checkbox's
// { attrs... } spread onto the native input, so browser :checked/:disabled
// truth drives the styling with no data-state plumbing.
component States() {
	<div class="flex flex-col gap-3">
		<div class="flex items-center gap-2">
			<ui.Checkbox id="checkbox-states-unchecked"/>
			<ui.Label for="checkbox-states-unchecked">Unchecked</ui.Label>
		</div>
		<div class="flex items-center gap-2">
			<ui.Checkbox id="checkbox-states-checked" checked/>
			<ui.Label for="checkbox-states-checked">Checked</ui.Label>
		</div>
		<div class="flex items-center gap-2">
			<ui.Checkbox id="checkbox-states-disabled" checked disabled/>
			<ui.Label for="checkbox-states-disabled">Checked and disabled</ui.Label>
		</div>
	</div>
}

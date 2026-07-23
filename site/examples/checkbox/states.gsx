package checkbox

import (
	uicheckbox "github.com/gsxhq/gsxui/ui/checkbox"
	uilabel "github.com/gsxhq/gsxui/ui/label"
)

// States renders Checkbox unchecked, checked, and disabled+checked —
// checked/disabled are bare boolean attributes forwarded through Checkbox's
// { attrs... } spread onto the native input, so browser :checked/:disabled
// truth drives the styling with no data-state plumbing.
component States() {
	<div class="flex flex-col gap-3">
		<div class="flex items-center gap-2">
			<uicheckbox.Checkbox id="checkbox-states-unchecked"/>
			<uilabel.Label for="checkbox-states-unchecked">Unchecked</uilabel.Label>
		</div>
		<div class="flex items-center gap-2">
			<uicheckbox.Checkbox id="checkbox-states-checked" checked/>
			<uilabel.Label for="checkbox-states-checked">Checked</uilabel.Label>
		</div>
		<div class="flex items-center gap-2">
			<uicheckbox.Checkbox id="checkbox-states-disabled" checked disabled/>
			<uilabel.Label for="checkbox-states-disabled">Checked and disabled</uilabel.Label>
		</div>
	</div>
}

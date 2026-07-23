package label

import (
	uicheckbox "github.com/gsxhq/gsxui/ui/checkbox"
	uilabel "github.com/gsxhq/gsxui/ui/label"
)

// Disabled shows Label's peer-disabled styling: Checkbox's base class
// carries "peer", so a disabled sibling Checkbox dims and cursor-blocks its
// Label automatically — no state plumbing needed on Label itself.
component Disabled() {
	<div class="flex items-center gap-2">
		<uicheckbox.Checkbox id="label-disabled-terms" disabled/>
		<uilabel.Label for="label-disabled-terms">Accept terms (disabled)</uilabel.Label>
	</div>
}

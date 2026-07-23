package label

import (
	"github.com/gsxhq/gsxui/ui"
)

// Disabled shows Label's peer-disabled styling: Checkbox's base class
// carries "peer", so a disabled sibling Checkbox dims and cursor-blocks its
// Label automatically — no state plumbing needed on Label itself.
component Disabled() {
	<div class="flex items-center gap-2">
		<ui.Checkbox id="label-disabled-terms" disabled/>
		<ui.Label for="label-disabled-terms">Accept terms (disabled)</ui.Label>
	</div>
}

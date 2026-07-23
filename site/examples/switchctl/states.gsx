package switchctl

import (
	"github.com/gsxhq/gsxui/ui"
)

// States adds disabled variants — disabled is a bare boolean attribute
// forwarded through Switch's { attrs... } spread onto the native
// checkbox+role=switch input.
component States() {
	<div class="flex flex-col gap-3">
		<div class="flex items-center gap-2">
			<ui.Switch id="switch-states-off-disabled" disabled/>
			<ui.Label for="switch-states-off-disabled">Off, disabled</ui.Label>
		</div>
		<div class="flex items-center gap-2">
			<ui.Switch id="switch-states-on-disabled" checked disabled/>
			<ui.Label for="switch-states-on-disabled">On, disabled</ui.Label>
		</div>
	</div>
}

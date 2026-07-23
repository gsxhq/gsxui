package switchctl

import (
	uilabel "github.com/gsxhq/gsxui/ui/label"
	uiswitch "github.com/gsxhq/gsxui/ui/switchctl"
)

// States adds disabled variants — disabled is a bare boolean attribute
// forwarded through Switch's { attrs... } spread onto the native
// checkbox+role=switch input.
component States() {
	<div class="flex flex-col gap-3">
		<div class="flex items-center gap-2">
			<uiswitch.Switch id="switch-states-off-disabled" disabled/>
			<uilabel.Label for="switch-states-off-disabled">Off, disabled</uilabel.Label>
		</div>
		<div class="flex items-center gap-2">
			<uiswitch.Switch id="switch-states-on-disabled" checked disabled/>
			<uilabel.Label for="switch-states-on-disabled">On, disabled</uilabel.Label>
		</div>
	</div>
}

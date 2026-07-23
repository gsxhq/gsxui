// Package switchctl holds the site's example gsx components for
// ui/switchctl. Each example is a real, compiled gsx component — the exact
// source below is what the component page displays AND what it renders, so
// source shown is source run; the examples_test.go drift test enforces
// they can't diverge.
package switchctl

import (
	uilabel "github.com/gsxhq/gsxui/ui/label"
	uiswitch "github.com/gsxhq/gsxui/ui/switchctl"
)

// Basic pairs an unchecked and a checked Switch, each with a Label.
component Basic() {
	<div class="flex flex-col gap-3">
		<div class="flex items-center gap-2">
			<uiswitch.Switch id="switch-basic-airplane"/>
			<uilabel.Label for="switch-basic-airplane">Airplane mode</uilabel.Label>
		</div>
		<div class="flex items-center gap-2">
			<uiswitch.Switch id="switch-basic-wifi" checked/>
			<uilabel.Label for="switch-basic-wifi">Wi-Fi</uilabel.Label>
		</div>
	</div>
}

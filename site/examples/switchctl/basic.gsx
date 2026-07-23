// Package switchctl holds the site's example gsx components for ui/switchctl.
package switchctl

import (
	"github.com/gsxhq/gsxui/ui"
)

// Basic pairs an unchecked and a checked Switch, each with a Label.
component Basic() {
	<div class="flex flex-col gap-3">
		<div class="flex items-center gap-2">
			<ui.Switch id="switch-basic-airplane"/>
			<ui.Label for="switch-basic-airplane">Airplane mode</ui.Label>
		</div>
		<div class="flex items-center gap-2">
			<ui.Switch id="switch-basic-wifi" checked/>
			<ui.Label for="switch-basic-wifi">Wi-Fi</ui.Label>
		</div>
	</div>
}

// Package buttongroup holds the site's example gsx components for
// ui/button-group.
package buttongroup

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Basic renders three realistic groups: a horizontal toolbar pair split by a
// separator (the button-group-separator demo shape), a quantity stepper
// showing ButtonGroupText between two icon buttons, and a vertical
// orientation group.
component Basic() {
	<div class="flex flex-wrap items-start gap-6">
		<ui.ButtonGroup>
			<ui.Button variant="outline">Archive</ui.Button>
			<ui.ButtonGroupSeparator/>
			<ui.Button variant="outline">Report</ui.Button>
		</ui.ButtonGroup>
		<ui.ButtonGroup aria-label="Quantity">
			<ui.Button variant="outline" size="icon" aria-label="Decrease quantity">
				<icon.Minus/>
			</ui.Button>
			<ui.ButtonGroupText>42</ui.ButtonGroupText>
			<ui.Button variant="outline" size="icon" aria-label="Increase quantity">
				<icon.Plus/>
			</ui.Button>
		</ui.ButtonGroup>
		<ui.ButtonGroup orientation="vertical" aria-label="Media controls" class="h-fit">
			<ui.Button variant="outline" size="icon">
				<icon.Plus/>
			</ui.Button>
			<ui.Button variant="outline" size="icon">
				<icon.Minus/>
			</ui.Button>
		</ui.ButtonGroup>
	</div>
}

// Package tooltip holds the site's example gsx components for ui/tooltip.
package tooltip

import (
	"github.com/gsxhq/gsxui/ui"
)

// Basic shows a tooltip on hover/focus of a Button trigger.
component Basic() {
	<ui.Tooltip>
		<ui.Button variant="outline" data-gsxui-tooltip-trigger>Hover me</ui.Button>
		<ui.TooltipContent>Add to library</ui.TooltipContent>
	</ui.Tooltip>
}

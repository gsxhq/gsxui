// Package tooltip holds the site's example gsx components for ui/tooltip.
package tooltip

import (
	uibutton "github.com/gsxhq/gsxui/ui/button"
	uitooltip "github.com/gsxhq/gsxui/ui/tooltip"
)

// Basic shows a tooltip on hover/focus of a Button trigger.
component Basic() {
	<uitooltip.Tooltip>
		<uibutton.Button variant="outline" data-gsxui-tooltip-trigger>Hover me</uibutton.Button>
		<uitooltip.TooltipContent>Add to library</uitooltip.TooltipContent>
	</uitooltip.Tooltip>
}

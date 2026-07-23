package tooltip

import (
	uibutton "github.com/gsxhq/gsxui/ui/button"
	uitooltip "github.com/gsxhq/gsxui/ui/tooltip"
)

// Wide overrides TooltipContent's default w-fit with a caller class,
// proving tailwind-merge wins over the base utility.
component Wide() {
	<uitooltip.Tooltip>
		<uibutton.Button variant="outline" data-gsxui-tooltip-trigger>Info</uibutton.Button>
		<uitooltip.TooltipContent class="w-64 text-center">Deploys automatically when you push to the main branch.</uitooltip.TooltipContent>
	</uitooltip.Tooltip>
}

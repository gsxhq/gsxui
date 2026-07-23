package tooltip

import (
	"github.com/gsxhq/gsxui/ui"
)

// Wide overrides TooltipContent's default w-fit with a caller class,
// proving tailwind-merge wins over the base utility.
component Wide() {
	<ui.Tooltip>
		<ui.Button variant="outline" data-gsxui-tooltip-trigger>Info</ui.Button>
		<ui.TooltipContent class="w-64 text-center">
			Deploys automatically when you push to the main branch.
		</ui.TooltipContent>
	</ui.Tooltip>
}

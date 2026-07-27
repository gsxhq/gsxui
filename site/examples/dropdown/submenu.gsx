package dropdown

import (
	"github.com/gsxhq/gsxui/ui"
)

// Submenu shows DropdownMenuSub/SubTrigger/SubContent nested inside a
// regular DropdownMenuContent — opened by hover, ArrowRight, or click; see
// ui/dropdown.gsx's file header for why SubContent must render DOM-nested
// rather than portalled.
component Submenu() {
	<ui.DropdownMenu>
		<ui.Button
			variant="outline"
			data-gsxui-dropdown-trigger
			data-gsxui-slot="dropdown-menu-trigger"
			aria-haspopup="menu"
			aria-expanded="false"
		>
			Actions
		</ui.Button>
		<ui.DropdownMenuContent>
			<ui.DropdownMenuItem>New File</ui.DropdownMenuItem>
			<ui.DropdownMenuSub>
				<ui.DropdownMenuSubTrigger>More Tools</ui.DropdownMenuSubTrigger>
				<ui.DropdownMenuSubContent>
					<ui.DropdownMenuItem>Save Page As...</ui.DropdownMenuItem>
					<ui.DropdownMenuItem>Create Shortcut...</ui.DropdownMenuItem>
					<ui.DropdownMenuItem>Name Window...</ui.DropdownMenuItem>
				</ui.DropdownMenuSubContent>
			</ui.DropdownMenuSub>
			<ui.DropdownMenuSeparator/>
			<ui.DropdownMenuItem>Exit</ui.DropdownMenuItem>
		</ui.DropdownMenuContent>
	</ui.DropdownMenu>
}

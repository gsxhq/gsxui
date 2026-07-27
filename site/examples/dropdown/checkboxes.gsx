package dropdown

import (
	"github.com/gsxhq/gsxui/ui"
)

// Checkboxes shows DropdownMenuCheckboxItem: one server-rendered checked,
// one unchecked. Selecting an item flips it in place (dropdown.js) and does
// not close the menu — a caller who needs the value elsewhere listens for
// gsxui:change on the item.
component Checkboxes() {
	<ui.DropdownMenu>
		<ui.Button
			variant="outline"
			data-gsxui-dropdown-trigger
			data-gsxui-slot="dropdown-menu-trigger"
			aria-haspopup="menu"
			aria-expanded="false"
		>
			View
		</ui.Button>
		<ui.DropdownMenuContent>
			<ui.DropdownMenuLabel>Appearance</ui.DropdownMenuLabel>
			<ui.DropdownMenuSeparator/>
			<ui.DropdownMenuCheckboxItem checked={true} value="toolbar">Show Toolbar</ui.DropdownMenuCheckboxItem>
			<ui.DropdownMenuCheckboxItem checked={false} value="statusbar">Show Status Bar</ui.DropdownMenuCheckboxItem>
			<ui.DropdownMenuCheckboxItem checked={false} value="sidebar" aria-disabled="true" data-disabled="true">Show Sidebar</ui.DropdownMenuCheckboxItem>
		</ui.DropdownMenuContent>
	</ui.DropdownMenu>
}

// Package dropdown holds the site's example gsx components for ui/dropdown.
package dropdown

import (
	"github.com/gsxhq/gsxui/ui"
)

// Basic renders a menu with a label, separator, and plain items. Its one
// trigger Button carries both Button styling and DropdownMenu's trigger
// behavior/state contract.
component Basic() {
	<ui.DropdownMenu>
		<ui.Button
			variant="outline"
			data-gsxui-dropdown-trigger
			data-gsxui-slot-dropdown-menu-trigger
			aria-haspopup="menu"
			aria-expanded="false"
		>
			Options
		</ui.Button>
		<ui.DropdownMenuContent>
			<ui.DropdownMenuGroup>
			<ui.DropdownMenuLabel>My Account</ui.DropdownMenuLabel>
			<ui.DropdownMenuSeparator/>
			<ui.DropdownMenuItem>Profile</ui.DropdownMenuItem>
			<ui.DropdownMenuItem>Billing</ui.DropdownMenuItem>
			<ui.DropdownMenuItem>Settings <ui.DropdownMenuShortcut>⌘,</ui.DropdownMenuShortcut></ui.DropdownMenuItem>
			</ui.DropdownMenuGroup>
		</ui.DropdownMenuContent>
	</ui.DropdownMenu>
}

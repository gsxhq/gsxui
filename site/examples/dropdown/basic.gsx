// Package dropdown holds the site's example gsx components for ui/dropdown.
package dropdown

import (
	"github.com/gsxhq/gsxui/ui"
)

// Basic renders a menu with a label, separator, and plain items.
component Basic() {
	<ui.DropdownMenu>
		<ui.DropdownMenuTrigger>Options</ui.DropdownMenuTrigger>
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

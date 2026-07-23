// Package dropdown holds the site's example gsx components for ui/dropdown.
package dropdown

import (
	uibutton "github.com/gsxhq/gsxui/ui/button"
	uidropdown "github.com/gsxhq/gsxui/ui/dropdown"
)

// Basic renders a menu with a label, separator, and plain items.
component Basic() {
	<uidropdown.DropdownMenu>
		<uibutton.Button variant="outline" data-gsxui-dropdown-trigger>Options</uibutton.Button>
		<uidropdown.DropdownMenuContent>
			<uidropdown.DropdownMenuLabel>My Account</uidropdown.DropdownMenuLabel>
			<uidropdown.DropdownMenuSeparator/>
			<uidropdown.DropdownMenuItem>Profile</uidropdown.DropdownMenuItem>
			<uidropdown.DropdownMenuItem>Billing</uidropdown.DropdownMenuItem>
			<uidropdown.DropdownMenuItem>Settings</uidropdown.DropdownMenuItem>
		</uidropdown.DropdownMenuContent>
	</uidropdown.DropdownMenu>
}

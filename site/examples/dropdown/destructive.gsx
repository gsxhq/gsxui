package dropdown

import (
	uibutton "github.com/gsxhq/gsxui/ui/button"
	uidropdown "github.com/gsxhq/gsxui/ui/dropdown"
)

// Destructive shows a destructive item alongside a disabled one — variant
// styles the row red; aria-disabled/data-disabled skip it in roving focus
// and the click handler (see ui/dropdown/dropdown.js).
component Destructive() {
	<uidropdown.DropdownMenu>
		<uibutton.Button variant="outline" data-gsxui-dropdown-trigger>Manage</uibutton.Button>
		<uidropdown.DropdownMenuContent>
			<uidropdown.DropdownMenuItem>Rename</uidropdown.DropdownMenuItem>
			<uidropdown.DropdownMenuItem aria-disabled="true" data-disabled="true">Archive (unavailable)</uidropdown.DropdownMenuItem>
			<uidropdown.DropdownMenuSeparator/>
			<uidropdown.DropdownMenuItem variant="destructive">Delete</uidropdown.DropdownMenuItem>
		</uidropdown.DropdownMenuContent>
	</uidropdown.DropdownMenu>
}

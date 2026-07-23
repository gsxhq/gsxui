package dropdown

import (
	"github.com/gsxhq/gsxui/ui"
)

// Destructive shows a destructive item alongside a disabled one — variant
// styles the row red; aria-disabled/data-disabled skip it in roving focus
// and the click handler (see ui/dropdown/dropdown.js).
component Destructive() {
	<ui.DropdownMenu>
		<ui.Button variant="outline" data-gsxui-dropdown-trigger>Manage</ui.Button>
		<ui.DropdownMenuContent>
			<ui.DropdownMenuItem>Rename</ui.DropdownMenuItem>
			<ui.DropdownMenuItem aria-disabled="true" data-disabled="true">Archive (unavailable)</ui.DropdownMenuItem>
			<ui.DropdownMenuSeparator/>
			<ui.DropdownMenuItem variant="destructive">Delete</ui.DropdownMenuItem>
		</ui.DropdownMenuContent>
	</ui.DropdownMenu>
}

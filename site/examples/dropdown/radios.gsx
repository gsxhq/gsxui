package dropdown

import (
	"github.com/gsxhq/gsxui/ui"
)

// Radios shows DropdownMenuRadioGroup/DropdownMenuRadioItem: "top" is the
// server-rendered current value. Selecting a different item closes the
// menu (unlike a checkbox item) and emits gsxui:change on the group with
// the newly picked value.
component Radios() {
	<ui.DropdownMenu>
		<ui.Button
			variant="outline"
			data-gsxui-dropdown-trigger
			data-gsxui-slot="dropdown-menu-trigger"
			aria-haspopup="menu"
			aria-expanded="false"
		>
			Panel Position
		</ui.Button>
		<ui.DropdownMenuContent>
			<ui.DropdownMenuLabel>Panel Position</ui.DropdownMenuLabel>
			<ui.DropdownMenuSeparator/>
			<ui.DropdownMenuRadioGroup value="top">
				<ui.DropdownMenuRadioItem checked={true} value="top">Top</ui.DropdownMenuRadioItem>
				<ui.DropdownMenuRadioItem checked={false} value="bottom">Bottom</ui.DropdownMenuRadioItem>
				<ui.DropdownMenuRadioItem checked={false} value="right">Right</ui.DropdownMenuRadioItem>
			</ui.DropdownMenuRadioGroup>
		</ui.DropdownMenuContent>
	</ui.DropdownMenu>
}

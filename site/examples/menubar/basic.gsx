// Package menubar holds the site's example gsx components for ui/menubar.
package menubar

import (
	"github.com/gsxhq/gsxui/ui"
)

// Basic mirrors shadcn's own menubar-demo.tsx (registry/new-york-v4/
// examples/menubar-demo.tsx) shape, trimmed to two menus: File (plain items
// with shortcuts, a separator) and Edit (plain items). Once one of the two
// triggers is open, hovering the other switches menus with no click — the
// open-follows-hover behavior this task's own ui/menubar.js adds.
component Basic() {
	<ui.Menubar>
		<ui.MenubarMenu>
			<ui.MenubarTrigger>File</ui.MenubarTrigger>
			<ui.MenubarContent>
				<ui.MenubarItem>
					New Tab
					<ui.MenubarShortcut>⌘T</ui.MenubarShortcut>
				</ui.MenubarItem>
				<ui.MenubarItem>
					New Window
					<ui.MenubarShortcut>⌘N</ui.MenubarShortcut>
				</ui.MenubarItem>
				<ui.MenubarSeparator/>
				<ui.MenubarItem>
					Print...
					<ui.MenubarShortcut>⌘P</ui.MenubarShortcut>
				</ui.MenubarItem>
			</ui.MenubarContent>
		</ui.MenubarMenu>
		<ui.MenubarMenu>
			<ui.MenubarTrigger>Edit</ui.MenubarTrigger>
			<ui.MenubarContent>
				<ui.MenubarItem>
					Undo
					<ui.MenubarShortcut>⌘Z</ui.MenubarShortcut>
				</ui.MenubarItem>
				<ui.MenubarItem>
					Redo
					<ui.MenubarShortcut>⇧⌘Z</ui.MenubarShortcut>
				</ui.MenubarItem>
				<ui.MenubarSeparator/>
				<ui.MenubarItem>Cut</ui.MenubarItem>
				<ui.MenubarItem>Copy</ui.MenubarItem>
				<ui.MenubarItem>Paste</ui.MenubarItem>
			</ui.MenubarContent>
		</ui.MenubarMenu>
	</ui.Menubar>
}

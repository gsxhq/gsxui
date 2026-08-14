package menubar

import (
	"github.com/gsxhq/gsxui/ui"
)

// Full mirrors shadcn's own menubar-demo.tsx (registry/new-york-v4/
// examples/menubar-demo.tsx) end to end: a File menu with a disabled item
// and a submenu, an Edit menu with checkbox items, a View menu with a
// radio group, and a Profiles-shaped menu — exercising every shared item
// type Task 1/2 already shipped (Item, CheckboxItem, RadioGroup/RadioItem,
// Sub/SubTrigger/SubContent) plus this task's own bar-level roving tabindex
// and open-follows-hover. Every `inset` prop shadcn's demo passes maps to
// the CSS-only data-inset attr axis (see docs/jsx-parity.md's
// ## dropdown-menu ledger).
component Full() {
	<ui.Menubar>
		<ui.MenubarMenu>
			<ui.MenubarTrigger>File</ui.MenubarTrigger>
			<ui.MenubarContent>
				<ui.MenubarGroup>
					<ui.MenubarLabel>File actions</ui.MenubarLabel>
					<ui.MenubarItem>
						New Tab
						<ui.MenubarShortcut>⌘T</ui.MenubarShortcut>
					</ui.MenubarItem>
					<ui.MenubarItem>
						New Window
						<ui.MenubarShortcut>⌘N</ui.MenubarShortcut>
					</ui.MenubarItem>
					<ui.MenubarItem variant="destructive" aria-disabled="true" data-disabled="true">
						New Incognito Window
					</ui.MenubarItem>
				</ui.MenubarGroup>
				<ui.MenubarSub>
					<ui.MenubarSubTrigger>Share</ui.MenubarSubTrigger>
					<ui.MenubarSubContent>
						<ui.MenubarItem>Email Link</ui.MenubarItem>
						<ui.MenubarItem>Messages</ui.MenubarItem>
						<ui.MenubarItem>Notes</ui.MenubarItem>
					</ui.MenubarSubContent>
				</ui.MenubarSub>
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
				<ui.MenubarSub>
					<ui.MenubarSubTrigger>Find</ui.MenubarSubTrigger>
					<ui.MenubarSubContent>
						<ui.MenubarItem>Search the web</ui.MenubarItem>
						<ui.MenubarSeparator/>
						<ui.MenubarItem>Find...</ui.MenubarItem>
						<ui.MenubarItem>Find Next</ui.MenubarItem>
						<ui.MenubarItem>Find Previous</ui.MenubarItem>
					</ui.MenubarSubContent>
				</ui.MenubarSub>
				<ui.MenubarSeparator/>
				<ui.MenubarItem>Cut</ui.MenubarItem>
				<ui.MenubarItem>Copy</ui.MenubarItem>
				<ui.MenubarItem>Paste</ui.MenubarItem>
			</ui.MenubarContent>
		</ui.MenubarMenu>
		<ui.MenubarMenu>
			<ui.MenubarTrigger>View</ui.MenubarTrigger>
			<ui.MenubarContent>
				<ui.MenubarCheckboxItem checked={true} value="bookmarks">Always Show Bookmarks Bar</ui.MenubarCheckboxItem>
				<ui.MenubarCheckboxItem checked={false} value="full-urls">Always Show Full URLs</ui.MenubarCheckboxItem>
				<ui.MenubarSeparator/>
				<ui.MenubarItem aria-disabled="true" data-disabled="true" data-inset="true">
					Reload
					<ui.MenubarShortcut>⌘R</ui.MenubarShortcut>
				</ui.MenubarItem>
				<ui.MenubarItem aria-disabled="true" data-disabled="true" data-inset="true">
					Force Reload
					<ui.MenubarShortcut>⇧⌘R</ui.MenubarShortcut>
				</ui.MenubarItem>
			</ui.MenubarContent>
		</ui.MenubarMenu>
		<ui.MenubarMenu>
			<ui.MenubarTrigger>Profiles</ui.MenubarTrigger>
			<ui.MenubarContent>
				<ui.MenubarRadioGroup value="benoit">
					<ui.MenubarRadioItem checked={false} value="andy">Andy</ui.MenubarRadioItem>
					<ui.MenubarRadioItem checked={true} value="benoit">Benoit</ui.MenubarRadioItem>
					<ui.MenubarRadioItem checked={false} value="luis">Luis</ui.MenubarRadioItem>
				</ui.MenubarRadioGroup>
				<ui.MenubarSeparator/>
				<ui.MenubarItem data-inset="true">Edit...</ui.MenubarItem>
				<ui.MenubarSeparator/>
				<ui.MenubarItem data-inset="true">Add Profile...</ui.MenubarItem>
			</ui.MenubarContent>
		</ui.MenubarMenu>
	</ui.Menubar>
}

package contextmenu

import (
	"github.com/gsxhq/gsxui/ui"
)

// Full mirrors shadcn's own context-menu-demo.tsx (registry/new-york-v4/
// examples/context-menu-demo.tsx) end to end, now that Tier 4 Batch B Task 2
// has shipped the Sub/checkbox/radio-group parts basic.gsx's own doc comment
// used to point at as dropped: a submenu (More Tools, itself holding a
// destructive-variant Delete), a pair of checkbox items, and a labeled radio
// group. Every `inset` prop shadcn's demo passes maps to the CSS-only
// data-inset attr axis (see docs/jsx-parity.md's ## context-menu ledger).
component Full() {
	<ui.ContextMenu>
		<ui.ContextMenuTrigger
			class="flex h-[150px] w-[300px] items-center justify-center rounded-md border border-dashed text-sm"
		>
			Right click here
		</ui.ContextMenuTrigger>
		<ui.ContextMenuContent class="w-52">
			<ui.ContextMenuGroup>
				<ui.ContextMenuItem data-inset="true">
					Back
					<ui.ContextMenuShortcut>⌘[</ui.ContextMenuShortcut>
				</ui.ContextMenuItem>
			</ui.ContextMenuGroup>
			<ui.ContextMenuItem aria-disabled="true" data-disabled="true" data-inset="true">
				Forward
				<ui.ContextMenuShortcut>⌘]</ui.ContextMenuShortcut>
			</ui.ContextMenuItem>
			<ui.ContextMenuItem data-inset="true">
				Reload
				<ui.ContextMenuShortcut>⌘R</ui.ContextMenuShortcut>
			</ui.ContextMenuItem>
			<ui.ContextMenuSub>
				<ui.ContextMenuSubTrigger data-inset="true">More Tools</ui.ContextMenuSubTrigger>
				<ui.ContextMenuSubContent class="w-44">
					<ui.ContextMenuItem>Save Page As...</ui.ContextMenuItem>
					<ui.ContextMenuItem>Create Shortcut...</ui.ContextMenuItem>
					<ui.ContextMenuItem>Name Window...</ui.ContextMenuItem>
					<ui.ContextMenuSeparator/>
					<ui.ContextMenuItem>Developer Tools</ui.ContextMenuItem>
					<ui.ContextMenuSeparator/>
					<ui.ContextMenuItem variant="destructive">Delete</ui.ContextMenuItem>
				</ui.ContextMenuSubContent>
			</ui.ContextMenuSub>
			<ui.ContextMenuSeparator/>
			<ui.ContextMenuCheckboxItem checked={true} value="bookmarks">Show Bookmarks</ui.ContextMenuCheckboxItem>
			<ui.ContextMenuCheckboxItem checked={false} value="full-urls">Show Full URLs</ui.ContextMenuCheckboxItem>
			<ui.ContextMenuSeparator/>
			<ui.ContextMenuRadioGroup value="pedro">
				<ui.ContextMenuLabel data-inset="true">People</ui.ContextMenuLabel>
				<ui.ContextMenuRadioItem checked={true} value="pedro">Pedro Duarte</ui.ContextMenuRadioItem>
				<ui.ContextMenuRadioItem checked={false} value="colm">Colm Tuite</ui.ContextMenuRadioItem>
			</ui.ContextMenuRadioGroup>
		</ui.ContextMenuContent>
	</ui.ContextMenu>
}

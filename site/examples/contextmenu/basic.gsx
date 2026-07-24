// Package contextmenu holds the site's example gsx components for
// ui/context-menu.
package contextmenu

import (
	"github.com/gsxhq/gsxui/ui"
)

// Basic mirrors shadcn's own context-menu-demo.tsx (registry/new-york-v4/
// examples/context-menu-demo.tsx): a dashed-border right-click zone, a few
// browser-chrome-shaped items (Back/Forward/Reload) each carrying a
// ContextMenuShortcut, a separator, and a destructive-variant item. The
// Sub/checkbox/radio-group parts of the real demo are dropped along with
// those components (see docs/jsx-parity.md's ## context-menu GAP entry).
component Basic() {
	<ui.ContextMenu>
		<ui.ContextMenuTrigger class="flex h-[150px] w-[300px] items-center justify-center rounded-md border border-dashed text-sm">
			Right click here
		</ui.ContextMenuTrigger>
		<ui.ContextMenuContent class="w-52">
			<ui.ContextMenuItem>
				Back
				<ui.ContextMenuShortcut>⌘[</ui.ContextMenuShortcut>
			</ui.ContextMenuItem>
			<ui.ContextMenuItem aria-disabled="true" data-disabled="true">
				Forward
				<ui.ContextMenuShortcut>⌘]</ui.ContextMenuShortcut>
			</ui.ContextMenuItem>
			<ui.ContextMenuItem>
				Reload
				<ui.ContextMenuShortcut>⌘R</ui.ContextMenuShortcut>
			</ui.ContextMenuItem>
			<ui.ContextMenuSeparator/>
			<ui.ContextMenuItem variant="destructive">Delete</ui.ContextMenuItem>
		</ui.ContextMenuContent>
	</ui.ContextMenu>
}

package contextmenu

import (
	"github.com/gsxhq/gsxui/ui"
)

// Rtl mirrors shadcn's own context-menu-rtl demo's Arabic strings for the
// browser-chrome items, kept in Basic's flat shape (no Sub/checkbox/radio —
// those parts are already dropped from this directory's Basic, per
// docs/jsx-parity.md's context-menu GAP entry, so Rtl doesn't reintroduce
// them either).
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.ContextMenu>
			<ui.ContextMenuTrigger
				class="flex h-[150px] w-[300px] items-center justify-center rounded-md border border-dashed text-sm"
			>
				انقر بزر الماوس الأيمن هنا
			</ui.ContextMenuTrigger>
			<ui.ContextMenuContent class="w-52">
				<ui.ContextMenuItem>
					رجوع
					<ui.ContextMenuShortcut>⌘[</ui.ContextMenuShortcut>
				</ui.ContextMenuItem>
				<ui.ContextMenuItem aria-disabled="true" data-disabled="true">
					تقدم
					<ui.ContextMenuShortcut>⌘]</ui.ContextMenuShortcut>
				</ui.ContextMenuItem>
				<ui.ContextMenuItem>
					إعادة تحميل
					<ui.ContextMenuShortcut>⌘R</ui.ContextMenuShortcut>
				</ui.ContextMenuItem>
				<ui.ContextMenuSeparator/>
				<ui.ContextMenuItem variant="destructive">حذف</ui.ContextMenuItem>
			</ui.ContextMenuContent>
		</ui.ContextMenu>
	</div>
}

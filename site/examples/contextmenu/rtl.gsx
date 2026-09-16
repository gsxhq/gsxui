package contextmenu

import (
	"github.com/gsxhq/gsxui/site/uirtl"
)

// Rtl mirrors shadcn's own context-menu-rtl demo's Arabic strings for the
// browser-chrome items, kept in Basic's flat shape (no Sub/checkbox/radio —
// those parts are already dropped from this directory's Basic, per
// docs/jsx-parity.md's context-menu GAP entry, so Rtl doesn't reintroduce
// them either).
component Rtl() {
	<div dir="rtl" lang="ar">
		<uirtl.ContextMenu>
			<uirtl.ContextMenuTrigger
				class="flex h-[150px] w-[300px] items-center justify-center rounded-md border border-dashed text-sm"
			>
				انقر بزر الماوس الأيمن هنا
			</uirtl.ContextMenuTrigger>
			<uirtl.ContextMenuContent class="w-52">
				<uirtl.ContextMenuItem>
					رجوع
					<uirtl.ContextMenuShortcut>⌘[</uirtl.ContextMenuShortcut>
				</uirtl.ContextMenuItem>
				<uirtl.ContextMenuItem aria-disabled="true" data-disabled="true">
					تقدم
					<uirtl.ContextMenuShortcut>⌘]</uirtl.ContextMenuShortcut>
				</uirtl.ContextMenuItem>
				<uirtl.ContextMenuItem>
					إعادة تحميل
					<uirtl.ContextMenuShortcut>⌘R</uirtl.ContextMenuShortcut>
				</uirtl.ContextMenuItem>
				<uirtl.ContextMenuSeparator/>
				<uirtl.ContextMenuItem variant="destructive">حذف</uirtl.ContextMenuItem>
			</uirtl.ContextMenuContent>
		</uirtl.ContextMenu>
	</div>
}

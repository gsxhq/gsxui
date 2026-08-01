package dropdown

import (
	"github.com/gsxhq/gsxui/ui"
)

// Rtl keeps Basic's flat menu shape (label, separator, plain items —
// shadcn's own dropdown-menu-rtl demo adds Sub/checkbox/radio groups this
// directory's Basic doesn't compose; DEVIATED from rather than invented)
// and translates its strings to Arabic, using dropdown-menu-rtl's own
// account/profile/billing/settings vocabulary.
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.DropdownMenu>
			<ui.Button
				variant="outline"
				data-gsxui-slot-dropdown-menu-trigger
				aria-haspopup="menu"
				aria-expanded="false"
			>
				فتح القائمة
			</ui.Button>
			<ui.DropdownMenuContent>
				<ui.DropdownMenuGroup>
					<ui.DropdownMenuLabel>الحساب</ui.DropdownMenuLabel>
					<ui.DropdownMenuSeparator/>
					<ui.DropdownMenuItem>الملف الشخصي</ui.DropdownMenuItem>
					<ui.DropdownMenuItem>الفوترة</ui.DropdownMenuItem>
					<ui.DropdownMenuItem>
						الإعدادات <ui.DropdownMenuShortcut>⌘,</ui.DropdownMenuShortcut>
					</ui.DropdownMenuItem>
				</ui.DropdownMenuGroup>
			</ui.DropdownMenuContent>
		</ui.DropdownMenu>
	</div>
}

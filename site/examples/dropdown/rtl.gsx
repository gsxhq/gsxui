package dropdown

import (
	"github.com/gsxhq/gsxui/site/uirtl"
)

// Rtl keeps Basic's flat menu shape (label, separator, plain items —
// shadcn's own dropdown-menu-rtl demo adds Sub/checkbox/radio groups this
// directory's Basic doesn't compose; DEVIATED from rather than invented)
// and translates its strings to Arabic, using dropdown-menu-rtl's own
// account/profile/billing/settings vocabulary.
component Rtl() {
	<div dir="rtl" lang="ar">
		<uirtl.DropdownMenu>
			<uirtl.Button
				variant="outline"
				data-gsxui-slot-dropdown-menu-trigger
				aria-haspopup="menu"
				aria-expanded="false"
			>
				فتح القائمة
			</uirtl.Button>
			<uirtl.DropdownMenuContent>
				<uirtl.DropdownMenuGroup>
					<uirtl.DropdownMenuLabel>الحساب</uirtl.DropdownMenuLabel>
					<uirtl.DropdownMenuSeparator/>
					<uirtl.DropdownMenuItem>الملف الشخصي</uirtl.DropdownMenuItem>
					<uirtl.DropdownMenuItem>الفوترة</uirtl.DropdownMenuItem>
					<uirtl.DropdownMenuItem>
						الإعدادات <uirtl.DropdownMenuShortcut>⌘,</uirtl.DropdownMenuShortcut>
					</uirtl.DropdownMenuItem>
				</uirtl.DropdownMenuGroup>
			</uirtl.DropdownMenuContent>
		</uirtl.DropdownMenu>
	</div>
}

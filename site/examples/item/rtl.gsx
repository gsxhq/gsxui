// Package item holds the site's example gsx components for ui/item.
package item

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Rtl mirrors shadcn's own item-rtl demo: an outline Item with a title and
// description plus an outline action button, and a second, size="sm" item
// (a verified-profile row with a leading badge icon and trailing chevron),
// translated to Arabic and wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar" class="flex w-full max-w-md flex-col gap-6">
		<ui.Item variant="outline">
			<ui.ItemContent>
				<ui.ItemTitle>عنصر أساسي</ui.ItemTitle>
				<ui.ItemDescription>عنصر بسيط يحتوي على عنوان ووصف.</ui.ItemDescription>
			</ui.ItemContent>
			<ui.ItemActions>
				<ui.Button variant="outline" size="sm">إجراء</ui.Button>
			</ui.ItemActions>
		</ui.Item>
		<ui.Item variant="outline" size="sm">
			<ui.ItemMedia>
				<icon.BadgeCheck class="size-5"/>
			</ui.ItemMedia>
			<ui.ItemContent>
				<ui.ItemTitle>تم التحقق من ملفك الشخصي.</ui.ItemTitle>
			</ui.ItemContent>
			<ui.ItemActions>
				<icon.ChevronRight class="size-4"/>
			</ui.ItemActions>
		</ui.Item>
	</div>
}

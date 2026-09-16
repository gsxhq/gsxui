// Package item holds the site's example gsx components for ui/item.
package item

import (
	"github.com/gsxhq/gsxui/site/uirtl"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Rtl mirrors shadcn's own item-rtl demo: an outline Item with a title and
// description plus an outline action button, and a second, size="sm" item
// (a verified-profile row with a leading badge icon and trailing chevron),
// translated to Arabic and wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar" class="flex w-full max-w-md flex-col gap-6">
		<uirtl.Item variant="outline">
			<uirtl.ItemContent>
				<uirtl.ItemTitle>عنصر أساسي</uirtl.ItemTitle>
				<uirtl.ItemDescription>عنصر بسيط يحتوي على عنوان ووصف.</uirtl.ItemDescription>
			</uirtl.ItemContent>
			<uirtl.ItemActions>
				<uirtl.Button variant="outline" size="sm">إجراء</uirtl.Button>
			</uirtl.ItemActions>
		</uirtl.Item>
		<uirtl.Item variant="outline" size="sm">
			<uirtl.ItemMedia>
				<icon.BadgeCheck class="size-5"/>
			</uirtl.ItemMedia>
			<uirtl.ItemContent>
				<uirtl.ItemTitle>تم التحقق من ملفك الشخصي.</uirtl.ItemTitle>
			</uirtl.ItemContent>
			<uirtl.ItemActions>
				<icon.ChevronRight class="size-4"/>
			</uirtl.ItemActions>
		</uirtl.Item>
	</div>
}

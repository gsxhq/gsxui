package badge

import (
	"github.com/gsxhq/gsxui/site/uirtl"
	uiicon "github.com/gsxhq/gsxui/ui/icon"
)

// Rtl mirrors shadcn's own badge-rtl demo: the four base variants plus an
// icon-leading and icon-trailing badge, translated to Arabic, wrapped in
// dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar" class="flex w-full flex-wrap justify-center gap-2">
		<uirtl.Badge>شارة</uirtl.Badge>
		<uirtl.Badge variant="secondary">ثانوي</uirtl.Badge>
		<uirtl.Badge variant="destructive">مدمر</uirtl.Badge>
		<uirtl.Badge variant="outline">مخطط</uirtl.Badge>
		<uirtl.Badge variant="secondary">
			<uiicon.BadgeCheck/>
			متحقق
		</uirtl.Badge>
		<uirtl.Badge variant="outline">
			إشارة مرجعية
			<uiicon.Bookmark/>
		</uirtl.Badge>
	</div>
}

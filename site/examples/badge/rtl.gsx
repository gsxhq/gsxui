package badge

import (
	"github.com/gsxhq/gsxui/ui"
	uiicon "github.com/gsxhq/gsxui/ui/icon"
)

// Rtl mirrors shadcn's own badge-rtl demo: the four base variants plus an
// icon-leading and icon-trailing badge, translated to Arabic, wrapped in
// dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar" class="flex w-full flex-wrap justify-center gap-2">
		<ui.Badge>شارة</ui.Badge>
		<ui.Badge variant="secondary">ثانوي</ui.Badge>
		<ui.Badge variant="destructive">مدمر</ui.Badge>
		<ui.Badge variant="outline">مخطط</ui.Badge>
		<ui.Badge variant="secondary">
			<uiicon.BadgeCheck/>
			متحقق
		</ui.Badge>
		<ui.Badge variant="outline">
			إشارة مرجعية
			<uiicon.Bookmark/>
		</ui.Badge>
	</div>
}

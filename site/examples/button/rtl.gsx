package button

import (
	"github.com/gsxhq/gsxui/ui"
	uiicon "github.com/gsxhq/gsxui/ui/icon"
)

// Rtl mirrors shadcn's own button-rtl demo: outline, destructive, an
// outline button with a trailing arrow that flips via rtl:rotate-180, an
// icon-only button, and a disabled secondary button with a loading
// spinner — translated to Arabic, wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar" class="flex flex-wrap items-center gap-2">
		<ui.Button variant="outline">زر</ui.Button>
		<ui.Button variant="destructive">حذف</ui.Button>
		<ui.Button variant="outline">
			إرسال
			<uiicon.ArrowRight class="rtl:rotate-180"/>
		</ui.Button>
		<ui.Button variant="outline" size="icon" aria-label="Add">
			<uiicon.Plus/>
		</ui.Button>
		<ui.Button variant="secondary" disabled>
			<ui.Spinner/>
			جاري التحميل
		</ui.Button>
	</div>
}

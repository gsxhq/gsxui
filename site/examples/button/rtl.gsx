package button

import (
	"github.com/gsxhq/gsxui/site/uirtl"
	uiicon "github.com/gsxhq/gsxui/ui/icon"
)

// Rtl mirrors shadcn's own button-rtl demo: outline, destructive, an
// outline button with a trailing arrow that flips via rtl:rotate-180, an
// icon-only button, and a disabled secondary button with a loading
// spinner — translated to Arabic, wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar" class="flex flex-wrap items-center gap-2">
		<uirtl.Button variant="outline">زر</uirtl.Button>
		<uirtl.Button variant="destructive">حذف</uirtl.Button>
		<uirtl.Button variant="outline">
			إرسال
			<uiicon.ArrowRight class="rtl:rotate-180"/>
		</uirtl.Button>
		<uirtl.Button variant="outline" size="icon" aria-label="Add">
			<uiicon.Plus/>
		</uirtl.Button>
		<uirtl.Button variant="secondary" disabled>
			<uirtl.Spinner/>
			جاري التحميل
		</uirtl.Button>
	</div>
}

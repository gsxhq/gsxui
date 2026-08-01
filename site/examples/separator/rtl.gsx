package separator

import "github.com/gsxhq/gsxui/ui"

// Rtl mirrors shadcn's own separator-rtl demo: a title/subtitle block, a
// Separator, then a description line, translated to Arabic and wrapped in
// dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar" class="flex max-w-sm flex-col gap-4 text-sm">
		<div class="flex flex-col gap-1.5">
			<div class="font-medium leading-none">shadcn/ui</div>
			<div class="text-muted-foreground">الأساس لنظام التصميم الخاص بك</div>
		</div>
		<ui.Separator/>
		<div>مجموعة من المكونات المصممة بشكل جميل يمكنك تخصيصها وتوسيعها والبناء عليها.</div>
	</div>
}

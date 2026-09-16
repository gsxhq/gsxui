package switchctl

import "github.com/gsxhq/gsxui/site/uirtl"

// Rtl mirrors shadcn's own switch-rtl demo: a horizontal Field pairing a
// label+description with a Switch, translated to Arabic and wrapped in
// dir="rtl".
component Rtl() {
	<uirtl.Field orientation="horizontal" dir="rtl" lang="ar" class="max-w-sm">
		<uirtl.FieldContent>
			<uirtl.FieldLabel for="switch-rtl-focus-mode">المشاركة عبر الأجهزة</uirtl.FieldLabel>
			<uirtl.FieldDescription>
				يتم مشاركة التركيز عبر الأجهزة، ويتم إيقاف تشغيله عند مغادرة التطبيق.
			</uirtl.FieldDescription>
		</uirtl.FieldContent>
		<uirtl.Switch id="switch-rtl-focus-mode"/>
	</uirtl.Field>
}

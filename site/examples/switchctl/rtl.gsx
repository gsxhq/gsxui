package switchctl

import "github.com/gsxhq/gsxui/ui"

// Rtl mirrors shadcn's own switch-rtl demo: a horizontal Field pairing a
// label+description with a Switch, translated to Arabic and wrapped in
// dir="rtl".
component Rtl() {
	<ui.Field orientation="horizontal" dir="rtl" lang="ar" class="max-w-sm">
		<ui.FieldContent>
			<ui.FieldLabel for="switch-rtl-focus-mode">المشاركة عبر الأجهزة</ui.FieldLabel>
			<ui.FieldDescription>يتم مشاركة التركيز عبر الأجهزة، ويتم إيقاف تشغيله عند مغادرة التطبيق.</ui.FieldDescription>
		</ui.FieldContent>
		<ui.Switch id="switch-rtl-focus-mode"/>
	</ui.Field>
}

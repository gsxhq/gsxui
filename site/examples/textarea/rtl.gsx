package textarea

import "github.com/gsxhq/gsxui/ui"

// Rtl mirrors shadcn's own textarea-rtl demo: a Field pairing a label,
// Textarea, and description, translated to Arabic and wrapped in dir="rtl".
component Rtl() {
	<ui.Field dir="rtl" lang="ar" class="w-full max-w-xs">
		<ui.FieldLabel for="textarea-rtl-feedback">التعليقات</ui.FieldLabel>
		<ui.Textarea id="textarea-rtl-feedback" value="" placeholder="تعليقاتك تساعدنا على التحسين..." rows="4"/>
		<ui.FieldDescription>شاركنا أفكارك حول خدمتنا.</ui.FieldDescription>
	</ui.Field>
}

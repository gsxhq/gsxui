package textarea

import "github.com/gsxhq/gsxui/site/uirtl"

// Rtl mirrors shadcn's own textarea-rtl demo: a Field pairing a label,
// Textarea, and description, translated to Arabic and wrapped in dir="rtl".
component Rtl() {
	<uirtl.Field dir="rtl" lang="ar" class="w-full max-w-xs">
		<uirtl.FieldLabel for="textarea-rtl-feedback">التعليقات</uirtl.FieldLabel>
		<uirtl.Textarea id="textarea-rtl-feedback" value="" placeholder="تعليقاتك تساعدنا على التحسين..." rows="4"/>
		<uirtl.FieldDescription>شاركنا أفكارك حول خدمتنا.</uirtl.FieldDescription>
	</uirtl.Field>
}

// Package radio holds the site's example gsx components for ui/radio.
package radio

import (
	"github.com/gsxhq/gsxui/ui"
)

// Rtl mirrors shadcn's own radio-group-rtl demo: three Field rows (default,
// comfortable — checked, compact), each pairing a Radio with a
// FieldLabel/FieldDescription via FieldContent, translated to Arabic and
// wrapped in dir="rtl".
component Rtl() {
	<ui.FieldGroup dir="rtl" lang="ar" class="w-fit">
		<ui.Field orientation="horizontal">
			<ui.Radio id="radio-rtl-default" name="radio-rtl-spacing"/>
			<ui.FieldContent>
				<ui.FieldLabel for="radio-rtl-default">افتراضي</ui.FieldLabel>
				<ui.FieldDescription>تباعد قياسي لمعظم حالات الاستخدام.</ui.FieldDescription>
			</ui.FieldContent>
		</ui.Field>
		<ui.Field orientation="horizontal">
			<ui.Radio id="radio-rtl-comfortable" name="radio-rtl-spacing" checked/>
			<ui.FieldContent>
				<ui.FieldLabel for="radio-rtl-comfortable">مريح</ui.FieldLabel>
				<ui.FieldDescription>مساحة أكبر بين العناصر.</ui.FieldDescription>
			</ui.FieldContent>
		</ui.Field>
		<ui.Field orientation="horizontal">
			<ui.Radio id="radio-rtl-compact" name="radio-rtl-spacing"/>
			<ui.FieldContent>
				<ui.FieldLabel for="radio-rtl-compact">مضغوط</ui.FieldLabel>
				<ui.FieldDescription>تباعد أدنى للتخطيطات الكثيفة.</ui.FieldDescription>
			</ui.FieldContent>
		</ui.Field>
	</ui.FieldGroup>
}

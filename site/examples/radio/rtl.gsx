// Package radio holds the site's example gsx components for ui/radio.
package radio

import (
	"github.com/gsxhq/gsxui/site/uirtl"
)

// Rtl mirrors shadcn's own radio-group-rtl demo: three Field rows (default,
// comfortable — checked, compact), each pairing a Radio with a
// FieldLabel/FieldDescription via FieldContent, translated to Arabic and
// wrapped in dir="rtl".
component Rtl() {
	<uirtl.FieldGroup dir="rtl" lang="ar" class="w-fit">
		<uirtl.Field orientation="horizontal">
			<uirtl.Radio id="radio-rtl-default" name="radio-rtl-spacing"/>
			<uirtl.FieldContent>
				<uirtl.FieldLabel for="radio-rtl-default">افتراضي</uirtl.FieldLabel>
				<uirtl.FieldDescription>تباعد قياسي لمعظم حالات الاستخدام.</uirtl.FieldDescription>
			</uirtl.FieldContent>
		</uirtl.Field>
		<uirtl.Field orientation="horizontal">
			<uirtl.Radio id="radio-rtl-comfortable" name="radio-rtl-spacing" checked/>
			<uirtl.FieldContent>
				<uirtl.FieldLabel for="radio-rtl-comfortable">مريح</uirtl.FieldLabel>
				<uirtl.FieldDescription>مساحة أكبر بين العناصر.</uirtl.FieldDescription>
			</uirtl.FieldContent>
		</uirtl.Field>
		<uirtl.Field orientation="horizontal">
			<uirtl.Radio id="radio-rtl-compact" name="radio-rtl-spacing"/>
			<uirtl.FieldContent>
				<uirtl.FieldLabel for="radio-rtl-compact">مضغوط</uirtl.FieldLabel>
				<uirtl.FieldDescription>تباعد أدنى للتخطيطات الكثيفة.</uirtl.FieldDescription>
			</uirtl.FieldContent>
		</uirtl.Field>
	</uirtl.FieldGroup>
}

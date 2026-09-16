package checkbox

import (
	"github.com/gsxhq/gsxui/site/uirtl"
)

// Rtl mirrors shadcn's own checkbox-rtl demo: Field-composed checkboxes
// covering plain, described, disabled, and title+description shapes,
// translated to Arabic and wrapped in dir="rtl".
component Rtl() {
	<uirtl.FieldGroup dir="rtl" lang="ar" class="max-w-sm">
		<uirtl.Field orientation="horizontal">
			<uirtl.Checkbox id="checkbox-rtl-terms"/>
			<uirtl.Label for="checkbox-rtl-terms">قبول الشروط والأحكام</uirtl.Label>
		</uirtl.Field>
		<uirtl.Field orientation="horizontal">
			<uirtl.Checkbox id="checkbox-rtl-terms-2" checked/>
			<uirtl.FieldContent>
				<uirtl.FieldLabel for="checkbox-rtl-terms-2">قبول الشروط والأحكام</uirtl.FieldLabel>
				<uirtl.FieldDescription>بالنقر على هذا المربع، فإنك توافق على الشروط.</uirtl.FieldDescription>
			</uirtl.FieldContent>
		</uirtl.Field>
		<uirtl.Field orientation="horizontal" data-disabled="true">
			<uirtl.Checkbox id="checkbox-rtl-notify" disabled/>
			<uirtl.FieldLabel for="checkbox-rtl-notify">تفعيل الإشعارات</uirtl.FieldLabel>
		</uirtl.Field>
		<uirtl.FieldLabel>
			<uirtl.Field orientation="horizontal">
				<uirtl.Checkbox id="checkbox-rtl-notify-2"/>
				<uirtl.FieldContent>
					<uirtl.FieldTitle>تفعيل الإشعارات</uirtl.FieldTitle>
					<uirtl.FieldDescription>يمكنك تفعيل أو إلغاء تفعيل الإشعارات في أي وقت.</uirtl.FieldDescription>
				</uirtl.FieldContent>
			</uirtl.Field>
		</uirtl.FieldLabel>
	</uirtl.FieldGroup>
}

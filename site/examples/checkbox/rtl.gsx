package checkbox

import (
	"github.com/gsxhq/gsxui/ui"
)

// Rtl mirrors shadcn's own checkbox-rtl demo: Field-composed checkboxes
// covering plain, described, disabled, and title+description shapes,
// translated to Arabic and wrapped in dir="rtl".
component Rtl() {
	<ui.FieldGroup dir="rtl" lang="ar" class="max-w-sm">
		<ui.Field orientation="horizontal">
			<ui.Checkbox id="checkbox-rtl-terms"/>
			<ui.Label for="checkbox-rtl-terms">قبول الشروط والأحكام</ui.Label>
		</ui.Field>
		<ui.Field orientation="horizontal">
			<ui.Checkbox id="checkbox-rtl-terms-2" checked/>
			<ui.FieldContent>
				<ui.FieldLabel for="checkbox-rtl-terms-2">قبول الشروط والأحكام</ui.FieldLabel>
				<ui.FieldDescription>بالنقر على هذا المربع، فإنك توافق على الشروط.</ui.FieldDescription>
			</ui.FieldContent>
		</ui.Field>
		<ui.Field orientation="horizontal" data-disabled="true">
			<ui.Checkbox id="checkbox-rtl-notify" disabled/>
			<ui.FieldLabel for="checkbox-rtl-notify">تفعيل الإشعارات</ui.FieldLabel>
		</ui.Field>
		<ui.FieldLabel>
			<ui.Field orientation="horizontal">
				<ui.Checkbox id="checkbox-rtl-notify-2"/>
				<ui.FieldContent>
					<ui.FieldTitle>تفعيل الإشعارات</ui.FieldTitle>
					<ui.FieldDescription>يمكنك تفعيل أو إلغاء تفعيل الإشعارات في أي وقت.</ui.FieldDescription>
				</ui.FieldContent>
			</ui.Field>
		</ui.FieldLabel>
	</ui.FieldGroup>
}

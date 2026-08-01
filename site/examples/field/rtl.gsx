package field

import "github.com/gsxhq/gsxui/ui"

// Rtl keeps Basic's own profile-form shape (FieldSet/FieldLegend/
// FieldGroup/Field/FieldLabel/Input/FieldDescription/FieldSeparator/
// FieldContent/FieldTitle/Textarea) translated to Arabic. shadcn's own
// field-rtl demo is a payment form using Select and Checkbox, neither of
// which this directory's Basic composes; DEVIATED from rather than
// invented.
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.FieldSet>
			<ui.FieldLegend variant="label">الملف الشخصي</ui.FieldLegend>
			<ui.FieldGroup class="group/field-group" data-variant="outline">
				<ui.Field>
					<ui.FieldLabel for="name-rtl">الاسم</ui.FieldLabel>
					<ui.Input id="name-rtl" placeholder="جميلة أحمد"/>
					<ui.FieldDescription>اسمك الكامل، كما يظهر في ملفك الشخصي.</ui.FieldDescription>
				</ui.Field>
				<ui.FieldSeparator/>
				<ui.FieldSeparator>أو</ui.FieldSeparator>
				<ui.Field orientation="horizontal">
					<ui.FieldContent>
						<ui.FieldTitle>حقل أفقي</ui.FieldTitle>
						<ui.FieldDescription>يستخدم المحتوى والعنوان نفس الفتحات المتعارف عليها.</ui.FieldDescription>
					</ui.FieldContent>
				</ui.Field>
				<ui.Field orientation="responsive" data-disabled="true">
					<ui.FieldContent>
						<ui.FieldTitle>حقل متجاوب معطل</ui.FieldTitle>
					</ui.FieldContent>
				</ui.Field>
				<ui.Field>
					<ui.FieldLabel for="bio-rtl">نبذة</ui.FieldLabel>
					<ui.Textarea id="bio-rtl" placeholder="أخبرنا عن نفسك"/>
					<ui.FieldDescription>تظهر للأعضاء الآخرين في ملفك الشخصي العام.</ui.FieldDescription>
				</ui.Field>
			</ui.FieldGroup>
		</ui.FieldSet>
	</div>
}

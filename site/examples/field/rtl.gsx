package field

import "github.com/gsxhq/gsxui/site/uirtl"

// Rtl keeps Basic's own profile-form shape (FieldSet/FieldLegend/
// FieldGroup/Field/FieldLabel/Input/FieldDescription/FieldSeparator/
// FieldContent/FieldTitle/Textarea) translated to Arabic. shadcn's own
// field-rtl demo is a payment form using Select and Checkbox, neither of
// which this directory's Basic composes; DEVIATED from rather than
// invented.
component Rtl() {
	<div dir="rtl" lang="ar">
		<uirtl.FieldSet>
			<uirtl.FieldLegend variant="label">الملف الشخصي</uirtl.FieldLegend>
			<uirtl.FieldGroup class="group/field-group" data-variant="outline">
				<uirtl.Field>
					<uirtl.FieldLabel for="name-rtl">الاسم</uirtl.FieldLabel>
					<uirtl.Input id="name-rtl" placeholder="جميلة أحمد"/>
					<uirtl.FieldDescription>اسمك الكامل، كما يظهر في ملفك الشخصي.</uirtl.FieldDescription>
				</uirtl.Field>
				<uirtl.FieldSeparator/>
				<uirtl.FieldSeparator>أو</uirtl.FieldSeparator>
				<uirtl.Field orientation="horizontal">
					<uirtl.FieldContent>
						<uirtl.FieldTitle>حقل أفقي</uirtl.FieldTitle>
						<uirtl.FieldDescription>يستخدم المحتوى والعنوان نفس الفتحات المتعارف عليها.</uirtl.FieldDescription>
					</uirtl.FieldContent>
				</uirtl.Field>
				<uirtl.Field orientation="responsive" data-disabled="true">
					<uirtl.FieldContent>
						<uirtl.FieldTitle>حقل متجاوب معطل</uirtl.FieldTitle>
					</uirtl.FieldContent>
				</uirtl.Field>
				<uirtl.Field>
					<uirtl.FieldLabel for="bio-rtl">نبذة</uirtl.FieldLabel>
					<uirtl.Textarea id="bio-rtl" placeholder="أخبرنا عن نفسك"/>
					<uirtl.FieldDescription>تظهر للأعضاء الآخرين في ملفك الشخصي العام.</uirtl.FieldDescription>
				</uirtl.Field>
			</uirtl.FieldGroup>
		</uirtl.FieldSet>
	</div>
}

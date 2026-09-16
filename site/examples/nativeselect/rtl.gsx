// Package nativeselect holds the site's example gsx components for
// uirtl.NativeSelect (hyphen stripped: Go package names can't contain one,
// same as button-group's buttongroup / context-menu's contextmenu).
package nativeselect

import "github.com/gsxhq/gsxui/site/uirtl"

// Rtl mirrors shadcn's own native-select-rtl demo: a status NativeSelect
// with a placeholder option and four status options, translated to Arabic
// and wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar">
		<uirtl.NativeSelect name="status-rtl">
			<uirtl.NativeSelectOption value="">اختر الحالة</uirtl.NativeSelectOption>
			<uirtl.NativeSelectOption value="todo">مهام</uirtl.NativeSelectOption>
			<uirtl.NativeSelectOption value="in-progress">قيد التنفيذ</uirtl.NativeSelectOption>
			<uirtl.NativeSelectOption value="done">منجز</uirtl.NativeSelectOption>
			<uirtl.NativeSelectOption value="cancelled">ملغي</uirtl.NativeSelectOption>
		</uirtl.NativeSelect>
	</div>
}

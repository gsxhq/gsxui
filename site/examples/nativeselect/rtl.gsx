// Package nativeselect holds the site's example gsx components for
// ui.NativeSelect (hyphen stripped: Go package names can't contain one,
// same as button-group's buttongroup / context-menu's contextmenu).
package nativeselect

import "github.com/gsxhq/gsxui/ui"

// Rtl mirrors shadcn's own native-select-rtl demo: a status NativeSelect
// with a placeholder option and four status options, translated to Arabic
// and wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.NativeSelect name="status-rtl">
			<ui.NativeSelectOption value="">اختر الحالة</ui.NativeSelectOption>
			<ui.NativeSelectOption value="todo">مهام</ui.NativeSelectOption>
			<ui.NativeSelectOption value="in-progress">قيد التنفيذ</ui.NativeSelectOption>
			<ui.NativeSelectOption value="done">منجز</ui.NativeSelectOption>
			<ui.NativeSelectOption value="cancelled">ملغي</ui.NativeSelectOption>
		</ui.NativeSelect>
	</div>
}

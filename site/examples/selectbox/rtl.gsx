// Package selectbox holds the site's example gsx components for the custom-
// listbox ui.Select (dir name "selectbox" — "select" is a Go keyword and
// can't name a package; the registered component key is still "select").
package selectbox

import "github.com/gsxhq/gsxui/ui"

// Rtl mirrors shadcn's own select-rtl demo: a fruits group plus a
// vegetables group separated by a SelectSeparator, translated to Arabic
// and wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.Select name="fruit-rtl">
			<ui.SelectTrigger size="sm" class="w-32">
				<ui.SelectValue placeholder="اختر فاكهة"/>
			</ui.SelectTrigger>
			<ui.SelectContent>
				<ui.SelectGroup>
					<ui.SelectLabel>الفواكه</ui.SelectLabel>
					<ui.SelectItem value="apple">تفاح</ui.SelectItem>
					<ui.SelectItem value="banana">موز</ui.SelectItem>
					<ui.SelectItem value="blueberry">توت أزرق</ui.SelectItem>
					<ui.SelectItem value="grapes">عنب</ui.SelectItem>
					<ui.SelectItem value="pineapple">أناناس</ui.SelectItem>
				</ui.SelectGroup>
				<ui.SelectSeparator/>
				<ui.SelectGroup>
					<ui.SelectLabel>الخضروات</ui.SelectLabel>
					<ui.SelectItem value="carrot">جزر</ui.SelectItem>
					<ui.SelectItem value="broccoli">بروكلي</ui.SelectItem>
					<ui.SelectItem value="spinach">سبانخ</ui.SelectItem>
				</ui.SelectGroup>
			</ui.SelectContent>
		</ui.Select>
	</div>
}

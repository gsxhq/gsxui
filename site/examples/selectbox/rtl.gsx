// Package selectbox holds the site's example gsx components for the custom-
// listbox uirtl.Select (dir name "selectbox" — "select" is a Go keyword and
// can't name a package; the registered component key is still "select").
package selectbox

import "github.com/gsxhq/gsxui/site/uirtl"

// Rtl mirrors shadcn's own select-rtl demo: a fruits group plus a
// vegetables group separated by a SelectSeparator, translated to Arabic
// and wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar">
		<uirtl.Select name="fruit-rtl">
			<uirtl.SelectTrigger size="sm" class="w-32">
				<uirtl.SelectValue placeholder="اختر فاكهة"/>
			</uirtl.SelectTrigger>
			<uirtl.SelectContent>
				<uirtl.SelectGroup>
					<uirtl.SelectLabel>الفواكه</uirtl.SelectLabel>
					<uirtl.SelectItem value="apple">تفاح</uirtl.SelectItem>
					<uirtl.SelectItem value="banana">موز</uirtl.SelectItem>
					<uirtl.SelectItem value="blueberry">توت أزرق</uirtl.SelectItem>
					<uirtl.SelectItem value="grapes">عنب</uirtl.SelectItem>
					<uirtl.SelectItem value="pineapple">أناناس</uirtl.SelectItem>
				</uirtl.SelectGroup>
				<uirtl.SelectSeparator/>
				<uirtl.SelectGroup>
					<uirtl.SelectLabel>الخضروات</uirtl.SelectLabel>
					<uirtl.SelectItem value="carrot">جزر</uirtl.SelectItem>
					<uirtl.SelectItem value="broccoli">بروكلي</uirtl.SelectItem>
					<uirtl.SelectItem value="spinach">سبانخ</uirtl.SelectItem>
				</uirtl.SelectGroup>
			</uirtl.SelectContent>
		</uirtl.Select>
	</div>
}

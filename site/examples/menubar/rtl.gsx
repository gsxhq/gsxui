// Package menubar holds the site's example gsx components for ui/menubar.
package menubar

import (
	"github.com/gsxhq/gsxui/site/uirtl"
)

// Rtl mirrors shadcn's own menubar-rtl demo (itself menubar-demo.tsx's
// shape, same as this dir's own Full): a File menu with a disabled item
// and a Share submenu, an Edit menu with a Find submenu, a View menu with
// checkbox items, and a Profiles menu with a radio group — translated to
// Arabic and wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar">
		<uirtl.Menubar>
			<uirtl.MenubarMenu>
				<uirtl.MenubarTrigger>ملف</uirtl.MenubarTrigger>
				<uirtl.MenubarContent>
					<uirtl.MenubarGroup>
						<uirtl.MenubarItem>
							علامة تبويب جديدة
							<uirtl.MenubarShortcut>⌘T</uirtl.MenubarShortcut>
						</uirtl.MenubarItem>
						<uirtl.MenubarItem>
							نافذة جديدة
							<uirtl.MenubarShortcut>⌘N</uirtl.MenubarShortcut>
						</uirtl.MenubarItem>
						<uirtl.MenubarItem aria-disabled="true" data-disabled="true">
							نافذة التصفح المتخفي الجديدة
						</uirtl.MenubarItem>
					</uirtl.MenubarGroup>
					<uirtl.MenubarSub>
						<uirtl.MenubarSubTrigger>مشاركة</uirtl.MenubarSubTrigger>
						<uirtl.MenubarSubContent>
							<uirtl.MenubarItem>رابط البريد الإلكتروني</uirtl.MenubarItem>
							<uirtl.MenubarItem>الرسائل</uirtl.MenubarItem>
							<uirtl.MenubarItem>الملاحظات</uirtl.MenubarItem>
						</uirtl.MenubarSubContent>
					</uirtl.MenubarSub>
					<uirtl.MenubarSeparator/>
					<uirtl.MenubarItem>
						طباعة...
						<uirtl.MenubarShortcut>⌘P</uirtl.MenubarShortcut>
					</uirtl.MenubarItem>
				</uirtl.MenubarContent>
			</uirtl.MenubarMenu>
			<uirtl.MenubarMenu>
				<uirtl.MenubarTrigger>تعديل</uirtl.MenubarTrigger>
				<uirtl.MenubarContent>
					<uirtl.MenubarItem>
						تراجع
						<uirtl.MenubarShortcut>⌘Z</uirtl.MenubarShortcut>
					</uirtl.MenubarItem>
					<uirtl.MenubarItem>
						إعادة
						<uirtl.MenubarShortcut>⇧⌘Z</uirtl.MenubarShortcut>
					</uirtl.MenubarItem>
					<uirtl.MenubarSeparator/>
					<uirtl.MenubarSub>
						<uirtl.MenubarSubTrigger>بحث</uirtl.MenubarSubTrigger>
						<uirtl.MenubarSubContent>
							<uirtl.MenubarItem>البحث على الويب</uirtl.MenubarItem>
							<uirtl.MenubarSeparator/>
							<uirtl.MenubarItem>بحث...</uirtl.MenubarItem>
							<uirtl.MenubarItem>البحث التالي</uirtl.MenubarItem>
							<uirtl.MenubarItem>البحث السابق</uirtl.MenubarItem>
						</uirtl.MenubarSubContent>
					</uirtl.MenubarSub>
					<uirtl.MenubarSeparator/>
					<uirtl.MenubarItem>قص</uirtl.MenubarItem>
					<uirtl.MenubarItem>نسخ</uirtl.MenubarItem>
					<uirtl.MenubarItem>لصق</uirtl.MenubarItem>
				</uirtl.MenubarContent>
			</uirtl.MenubarMenu>
			<uirtl.MenubarMenu>
				<uirtl.MenubarTrigger>عرض</uirtl.MenubarTrigger>
				<uirtl.MenubarContent>
					<uirtl.MenubarCheckboxItem checked={true} value="bookmarks">شريط الإشارات المرجعية</uirtl.MenubarCheckboxItem>
					<uirtl.MenubarCheckboxItem checked={false} value="full-urls">عناوين URL الكاملة</uirtl.MenubarCheckboxItem>
					<uirtl.MenubarSeparator/>
					<uirtl.MenubarItem aria-disabled="true" data-disabled="true">
						إعادة تحميل
						<uirtl.MenubarShortcut>⌘R</uirtl.MenubarShortcut>
					</uirtl.MenubarItem>
					<uirtl.MenubarItem aria-disabled="true" data-disabled="true">
						إعادة تحميل قسري
						<uirtl.MenubarShortcut>⇧⌘R</uirtl.MenubarShortcut>
					</uirtl.MenubarItem>
				</uirtl.MenubarContent>
			</uirtl.MenubarMenu>
			<uirtl.MenubarMenu>
				<uirtl.MenubarTrigger>الملفات الشخصية</uirtl.MenubarTrigger>
				<uirtl.MenubarContent>
					<uirtl.MenubarRadioGroup value="benoit">
						<uirtl.MenubarRadioItem checked={false} value="andy">Andy</uirtl.MenubarRadioItem>
						<uirtl.MenubarRadioItem checked={true} value="benoit">Benoit</uirtl.MenubarRadioItem>
						<uirtl.MenubarRadioItem checked={false} value="luis">Luis</uirtl.MenubarRadioItem>
					</uirtl.MenubarRadioGroup>
					<uirtl.MenubarSeparator/>
					<uirtl.MenubarItem>تعديل...</uirtl.MenubarItem>
					<uirtl.MenubarSeparator/>
					<uirtl.MenubarItem>إضافة ملف شخصي...</uirtl.MenubarItem>
				</uirtl.MenubarContent>
			</uirtl.MenubarMenu>
		</uirtl.Menubar>
	</div>
}

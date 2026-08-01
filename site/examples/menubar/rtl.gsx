// Package menubar holds the site's example gsx components for ui/menubar.
package menubar

import (
	"github.com/gsxhq/gsxui/ui"
)

// Rtl mirrors shadcn's own menubar-rtl demo (itself menubar-demo.tsx's
// shape, same as this dir's own Full): a File menu with a disabled item
// and a Share submenu, an Edit menu with a Find submenu, a View menu with
// checkbox items, and a Profiles menu with a radio group — translated to
// Arabic and wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.Menubar>
			<ui.MenubarMenu>
				<ui.MenubarTrigger>ملف</ui.MenubarTrigger>
				<ui.MenubarContent>
					<ui.MenubarGroup>
						<ui.MenubarItem>
							علامة تبويب جديدة
							<ui.MenubarShortcut>⌘T</ui.MenubarShortcut>
						</ui.MenubarItem>
						<ui.MenubarItem>
							نافذة جديدة
							<ui.MenubarShortcut>⌘N</ui.MenubarShortcut>
						</ui.MenubarItem>
						<ui.MenubarItem aria-disabled="true" data-disabled="true">نافذة التصفح المتخفي الجديدة</ui.MenubarItem>
					</ui.MenubarGroup>
					<ui.MenubarSub>
						<ui.MenubarSubTrigger>مشاركة</ui.MenubarSubTrigger>
						<ui.MenubarSubContent>
							<ui.MenubarItem>رابط البريد الإلكتروني</ui.MenubarItem>
							<ui.MenubarItem>الرسائل</ui.MenubarItem>
							<ui.MenubarItem>الملاحظات</ui.MenubarItem>
						</ui.MenubarSubContent>
					</ui.MenubarSub>
					<ui.MenubarSeparator/>
					<ui.MenubarItem>
						طباعة...
						<ui.MenubarShortcut>⌘P</ui.MenubarShortcut>
					</ui.MenubarItem>
				</ui.MenubarContent>
			</ui.MenubarMenu>
			<ui.MenubarMenu>
				<ui.MenubarTrigger>تعديل</ui.MenubarTrigger>
				<ui.MenubarContent>
					<ui.MenubarItem>
						تراجع
						<ui.MenubarShortcut>⌘Z</ui.MenubarShortcut>
					</ui.MenubarItem>
					<ui.MenubarItem>
						إعادة
						<ui.MenubarShortcut>⇧⌘Z</ui.MenubarShortcut>
					</ui.MenubarItem>
					<ui.MenubarSeparator/>
					<ui.MenubarSub>
						<ui.MenubarSubTrigger>بحث</ui.MenubarSubTrigger>
						<ui.MenubarSubContent>
							<ui.MenubarItem>البحث على الويب</ui.MenubarItem>
							<ui.MenubarSeparator/>
							<ui.MenubarItem>بحث...</ui.MenubarItem>
							<ui.MenubarItem>البحث التالي</ui.MenubarItem>
							<ui.MenubarItem>البحث السابق</ui.MenubarItem>
						</ui.MenubarSubContent>
					</ui.MenubarSub>
					<ui.MenubarSeparator/>
					<ui.MenubarItem>قص</ui.MenubarItem>
					<ui.MenubarItem>نسخ</ui.MenubarItem>
					<ui.MenubarItem>لصق</ui.MenubarItem>
				</ui.MenubarContent>
			</ui.MenubarMenu>
			<ui.MenubarMenu>
				<ui.MenubarTrigger>عرض</ui.MenubarTrigger>
				<ui.MenubarContent>
					<ui.MenubarCheckboxItem checked={true} value="bookmarks">شريط الإشارات المرجعية</ui.MenubarCheckboxItem>
					<ui.MenubarCheckboxItem checked={false} value="full-urls">عناوين URL الكاملة</ui.MenubarCheckboxItem>
					<ui.MenubarSeparator/>
					<ui.MenubarItem aria-disabled="true" data-disabled="true">
						إعادة تحميل
						<ui.MenubarShortcut>⌘R</ui.MenubarShortcut>
					</ui.MenubarItem>
					<ui.MenubarItem aria-disabled="true" data-disabled="true">
						إعادة تحميل قسري
						<ui.MenubarShortcut>⇧⌘R</ui.MenubarShortcut>
					</ui.MenubarItem>
				</ui.MenubarContent>
			</ui.MenubarMenu>
			<ui.MenubarMenu>
				<ui.MenubarTrigger>الملفات الشخصية</ui.MenubarTrigger>
				<ui.MenubarContent>
					<ui.MenubarRadioGroup value="benoit">
						<ui.MenubarRadioItem checked={false} value="andy">Andy</ui.MenubarRadioItem>
						<ui.MenubarRadioItem checked={true} value="benoit">Benoit</ui.MenubarRadioItem>
						<ui.MenubarRadioItem checked={false} value="luis">Luis</ui.MenubarRadioItem>
					</ui.MenubarRadioGroup>
					<ui.MenubarSeparator/>
					<ui.MenubarItem>تعديل...</ui.MenubarItem>
					<ui.MenubarSeparator/>
					<ui.MenubarItem>إضافة ملف شخصي...</ui.MenubarItem>
				</ui.MenubarContent>
			</ui.MenubarMenu>
		</ui.Menubar>
	</div>
}

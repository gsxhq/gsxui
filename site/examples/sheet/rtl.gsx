package sheet

import "github.com/gsxhq/gsxui/ui"

// Rtl mirrors shadcn's own sheet-rtl demo: the same profile-editor sheet as
// Basic, but flipped to open from the left (the RTL-appropriate side) and
// translated to Arabic, wrapped in dir="rtl". shadcn also swaps the footer
// close button for a dedicated SheetClose render prop; this port keeps the
// same data-gsxui-dialog-close idiom Basic and Directions already use for
// the identical button-in-button reason documented on Basic.
component Rtl() {
	<ui.Sheet>
		<ui.Button
			variant="outline"
			data-gsxui-slot-sheet-trigger
			aria-haspopup="dialog"
			aria-expanded="false"
		>
			فتح
		</ui.Button>
		<ui.SheetContent dir="rtl" lang="ar" side="left">
			<ui.SheetHeader>
				<ui.SheetTitle>تعديل الملف الشخصي</ui.SheetTitle>
				<ui.SheetDescription>
					قم بإجراء تغييرات على ملفك الشخصي هنا. انقر حفظ عند الانتهاء.
				</ui.SheetDescription>
			</ui.SheetHeader>
			<div class="grid gap-4 px-4">
				<div class="grid grid-cols-4 items-center gap-4">
					<ui.Label for="sheet-rtl-name" class="text-right">الاسم</ui.Label>
					<ui.Input id="sheet-rtl-name" value="Pedro Duarte" class="col-span-3"/>
				</div>
				<div class="grid grid-cols-4 items-center gap-4">
					<ui.Label for="sheet-rtl-username" class="text-right">اسم المستخدم</ui.Label>
					<ui.Input id="sheet-rtl-username" value="@peduarte" class="col-span-3"/>
				</div>
			</div>
			<ui.SheetFooter>
				<ui.Button data-gsxui-dialog-close>حفظ التغييرات</ui.Button>
			</ui.SheetFooter>
		</ui.SheetContent>
	</ui.Sheet>
}

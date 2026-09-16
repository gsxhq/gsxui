package sheet

import "github.com/gsxhq/gsxui/site/uirtl"

// Rtl mirrors shadcn's own sheet-rtl demo: the same profile-editor sheet as
// Basic, but flipped to open from the left (the RTL-appropriate side) and
// translated to Arabic, wrapped in dir="rtl". shadcn also swaps the footer
// close button for a dedicated SheetClose render prop; this port keeps the
// same data-gsxui-dialog-close idiom Basic and Directions already use for
// the identical button-in-button reason documented on Basic.
component Rtl() {
	<uirtl.Sheet>
		<uirtl.Button
			variant="outline"
			data-gsxui-slot-sheet-trigger
			aria-haspopup="dialog"
			aria-expanded="false"
		>
			فتح
		</uirtl.Button>
		<uirtl.SheetContent dir="rtl" lang="ar" side="left">
			<uirtl.SheetHeader>
				<uirtl.SheetTitle>تعديل الملف الشخصي</uirtl.SheetTitle>
				<uirtl.SheetDescription>
					قم بإجراء تغييرات على ملفك الشخصي هنا. انقر حفظ عند الانتهاء.
				</uirtl.SheetDescription>
			</uirtl.SheetHeader>
			<div class="grid gap-4 px-4">
				<div class="grid grid-cols-4 items-center gap-4">
					<uirtl.Label for="sheet-rtl-name" class="text-right">الاسم</uirtl.Label>
					<uirtl.Input id="sheet-rtl-name" value="Pedro Duarte" class="col-span-3"/>
				</div>
				<div class="grid grid-cols-4 items-center gap-4">
					<uirtl.Label for="sheet-rtl-username" class="text-right">اسم المستخدم</uirtl.Label>
					<uirtl.Input id="sheet-rtl-username" value="@peduarte" class="col-span-3"/>
				</div>
			</div>
			<uirtl.SheetFooter>
				<uirtl.Button data-gsxui-dialog-close>حفظ التغييرات</uirtl.Button>
			</uirtl.SheetFooter>
		</uirtl.SheetContent>
	</uirtl.Sheet>
}

package dialog

import (
	"github.com/gsxhq/gsxui/site/uirtl"
)

// Rtl keeps Basic's confirm-dialog shape (shadcn's own dialog-rtl demo is
// an edit-profile form using Field/Input, which this directory's Basic
// doesn't compose — DEVIATED from rather than invented) and translates its
// strings to Arabic, wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar">
		<uirtl.Dialog>
			<uirtl.Button
				variant="outline"
				data-gsxui-slot-dialog-trigger
				aria-haspopup="dialog"
				aria-expanded="false"
			>
				حذف الحساب
			</uirtl.Button>
			<uirtl.DialogContent>
				<uirtl.DialogHeader>
					<uirtl.DialogTitle>هل أنت متأكد تمامًا؟</uirtl.DialogTitle>
					<uirtl.DialogDescription>سيؤدي هذا إلى حذف حسابك نهائيًا.</uirtl.DialogDescription>
				</uirtl.DialogHeader>
				<uirtl.DialogFooter>
					<uirtl.Button variant="outline" data-gsxui-dialog-close>إلغاء</uirtl.Button>
					<uirtl.Button variant="destructive">متابعة</uirtl.Button>
				</uirtl.DialogFooter>
			</uirtl.DialogContent>
		</uirtl.Dialog>
	</div>
}

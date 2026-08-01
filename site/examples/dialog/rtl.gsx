package dialog

import (
	"github.com/gsxhq/gsxui/ui"
)

// Rtl keeps Basic's confirm-dialog shape (shadcn's own dialog-rtl demo is
// an edit-profile form using Field/Input, which this directory's Basic
// doesn't compose — DEVIATED from rather than invented) and translates its
// strings to Arabic, wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.Dialog>
			<ui.Button
				variant="outline"
				data-gsxui-slot-dialog-trigger
				aria-haspopup="dialog"
				aria-expanded="false"
			>
				حذف الحساب
			</ui.Button>
			<ui.DialogContent>
				<ui.DialogHeader>
					<ui.DialogTitle>هل أنت متأكد تمامًا؟</ui.DialogTitle>
					<ui.DialogDescription>سيؤدي هذا إلى حذف حسابك نهائيًا.</ui.DialogDescription>
				</ui.DialogHeader>
				<ui.DialogFooter>
					<ui.Button variant="outline" data-gsxui-dialog-close>إلغاء</ui.Button>
					<ui.Button variant="destructive">متابعة</ui.Button>
				</ui.DialogFooter>
			</ui.DialogContent>
		</ui.Dialog>
	</div>
}

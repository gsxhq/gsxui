package alertdialog

import (
	"github.com/gsxhq/gsxui/ui"
)

// Rtl mirrors Basic's destructive-confirm flow, translated to Arabic and
// wrapped in dir="rtl". shadcn's own alert-dialog-rtl demo also renders a
// second, smaller dialog with an AlertDialogMedia icon slot; gsxui's
// AlertDialogContent has no size variant or media slot (see
// ui/alert-dialog.gsx), so this example keeps the single dialog shape
// Basic already establishes rather than inventing new API surface.
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.AlertDialog>
			<ui.Button
				variant="outline"
				data-gsxui-slot-alert-dialog-trigger
				aria-haspopup="dialog"
				aria-expanded="false"
			>
				إظهار الحوار
			</ui.Button>
			<ui.AlertDialogContent dir="rtl" lang="ar">
				<ui.AlertDialogHeader>
					<ui.AlertDialogTitle>هل أنت متأكد تمامًا؟</ui.AlertDialogTitle>
					<ui.AlertDialogDescription>
						لا يمكن التراجع عن هذا الإجراء. سيؤدي هذا إلى حذف حسابك نهائيًا من خوادمنا.
					</ui.AlertDialogDescription>
				</ui.AlertDialogHeader>
				<ui.AlertDialogFooter>
					<ui.AlertDialogCancel>إلغاء</ui.AlertDialogCancel>
					<ui.AlertDialogAction>متابعة</ui.AlertDialogAction>
				</ui.AlertDialogFooter>
			</ui.AlertDialogContent>
		</ui.AlertDialog>
	</div>
}

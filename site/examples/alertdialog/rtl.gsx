package alertdialog

import (
	"github.com/gsxhq/gsxui/site/uirtl"
)

// Rtl mirrors Basic's destructive-confirm flow, translated to Arabic and
// wrapped in dir="rtl". shadcn's own alert-dialog-rtl demo also renders a
// second, smaller dialog with an AlertDialogMedia icon slot; gsxui's
// AlertDialogContent has no size variant or media slot (see
// ui/alert-dialog.gsx), so this example keeps the single dialog shape
// Basic already establishes rather than inventing new API surface.
component Rtl() {
	<div dir="rtl" lang="ar">
		<uirtl.AlertDialog>
			<uirtl.Button
				variant="outline"
				data-gsxui-slot-alert-dialog-trigger
				aria-haspopup="dialog"
				aria-expanded="false"
			>
				إظهار الحوار
			</uirtl.Button>
			<uirtl.AlertDialogContent dir="rtl" lang="ar">
				<uirtl.AlertDialogHeader>
					<uirtl.AlertDialogTitle>هل أنت متأكد تمامًا؟</uirtl.AlertDialogTitle>
					<uirtl.AlertDialogDescription>
						لا يمكن التراجع عن هذا الإجراء. سيؤدي هذا إلى حذف حسابك نهائيًا من خوادمنا.
					</uirtl.AlertDialogDescription>
				</uirtl.AlertDialogHeader>
				<uirtl.AlertDialogFooter>
					<uirtl.AlertDialogCancel>إلغاء</uirtl.AlertDialogCancel>
					<uirtl.AlertDialogAction>متابعة</uirtl.AlertDialogAction>
				</uirtl.AlertDialogFooter>
			</uirtl.AlertDialogContent>
		</uirtl.AlertDialog>
	</div>
}

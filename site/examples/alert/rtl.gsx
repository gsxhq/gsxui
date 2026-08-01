package alert

import (
	"github.com/gsxhq/gsxui/ui"
	uiicon "github.com/gsxhq/gsxui/ui/icon"
)

// Rtl mirrors shadcn's own alert-rtl demo: a payment-success Alert and a
// new-feature Alert, translated to Arabic, wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar" class="grid w-full max-w-md items-start gap-4">
		<ui.Alert>
			<uiicon.CircleCheck/>
			<ui.AlertTitle>تم الدفع بنجاح</ui.AlertTitle>
			<ui.AlertDescription>
				تمت معالجة دفعتك البالغة 29.99 دولارًا. تم إرسال إيصال إلى عنوان بريدك الإلكتروني.
			</ui.AlertDescription>
		</ui.Alert>
		<ui.Alert>
			<uiicon.Info/>
			<ui.AlertTitle>ميزة جديدة متاحة</ui.AlertTitle>
			<ui.AlertDescription>لقد أضفنا دعم الوضع الداكن. يمكنك تفعيله في إعدادات حسابك.</ui.AlertDescription>
		</ui.Alert>
	</div>
}

package alert

import (
	"github.com/gsxhq/gsxui/site/uirtl"
	uiicon "github.com/gsxhq/gsxui/ui/icon"
)

// Rtl mirrors shadcn's own alert-rtl demo: a payment-success Alert and a
// new-feature Alert, translated to Arabic, wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar" class="grid w-full max-w-md items-start gap-4">
		<uirtl.Alert>
			<uiicon.CircleCheck/>
			<uirtl.AlertTitle>تم الدفع بنجاح</uirtl.AlertTitle>
			<uirtl.AlertDescription>
				تمت معالجة دفعتك البالغة 29.99 دولارًا. تم إرسال إيصال إلى عنوان بريدك الإلكتروني.
			</uirtl.AlertDescription>
		</uirtl.Alert>
		<uirtl.Alert>
			<uiicon.Info/>
			<uirtl.AlertTitle>ميزة جديدة متاحة</uirtl.AlertTitle>
			<uirtl.AlertDescription>لقد أضفنا دعم الوضع الداكن. يمكنك تفعيله في إعدادات حسابك.</uirtl.AlertDescription>
		</uirtl.Alert>
	</div>
}

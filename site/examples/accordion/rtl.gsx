package accordion

import "github.com/gsxhq/gsxui/ui"

// Rtl mirrors shadcn's own accordion-rtl demo: the same three FAQ
// question/answer pairs, translated to Arabic, wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.Accordion name="faq-rtl" class="max-w-md">
			<ui.AccordionItem name="faq-rtl" open>
				<ui.AccordionTrigger>كيف يمكنني إعادة تعيين كلمة المرور؟</ui.AccordionTrigger>
				<ui.AccordionContent>
					انقر على 'نسيت كلمة المرور' في صفحة تسجيل الدخول، أدخل عنوان بريدك الإلكتروني، وسنرسل لك رابطًا لإعادة تعيين
					كلمة المرور. سينتهي صلاحية الرابط خلال 24 ساعة.
				</ui.AccordionContent>
			</ui.AccordionItem>
			<ui.AccordionItem name="faq-rtl">
				<ui.AccordionTrigger>هل يمكنني تغيير خطة الاشتراك الخاصة بي؟</ui.AccordionTrigger>
				<ui.AccordionContent>
					نعم، يمكنك ترقية أو تخفيض خطتك في أي وقت من إعدادات حسابك. ستظهر التغييرات في دورة الفوترة التالية.
				</ui.AccordionContent>
			</ui.AccordionItem>
			<ui.AccordionItem name="faq-rtl">
				<ui.AccordionTrigger>ما هي طرق الدفع التي تقبلونها؟</ui.AccordionTrigger>
				<ui.AccordionContent>
					نقبل جميع بطاقات الائتمان الرئيسية و PayPal والتحويلات المصرفية. تتم معالجة جميع المدفوعات بأمان من خلال شركاء
					الدفع لدينا.
				</ui.AccordionContent>
			</ui.AccordionItem>
		</ui.Accordion>
	</div>
}

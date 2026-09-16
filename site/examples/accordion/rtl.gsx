package accordion

import "github.com/gsxhq/gsxui/site/uirtl"

// Rtl mirrors shadcn's own accordion-rtl demo: the same three FAQ
// question/answer pairs, translated to Arabic, wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar">
		<uirtl.Accordion name="faq-rtl" class="max-w-md">
			<uirtl.AccordionItem name="faq-rtl" open>
				<uirtl.AccordionTrigger>كيف يمكنني إعادة تعيين كلمة المرور؟</uirtl.AccordionTrigger>
				<uirtl.AccordionContent>
					انقر على 'نسيت كلمة المرور' في صفحة تسجيل الدخول، أدخل عنوان بريدك الإلكتروني، وسنرسل لك رابطًا لإعادة تعيين
					كلمة المرور. سينتهي صلاحية الرابط خلال 24 ساعة.
				</uirtl.AccordionContent>
			</uirtl.AccordionItem>
			<uirtl.AccordionItem name="faq-rtl">
				<uirtl.AccordionTrigger>هل يمكنني تغيير خطة الاشتراك الخاصة بي؟</uirtl.AccordionTrigger>
				<uirtl.AccordionContent>
					نعم، يمكنك ترقية أو تخفيض خطتك في أي وقت من إعدادات حسابك. ستظهر التغييرات في دورة الفوترة التالية.
				</uirtl.AccordionContent>
			</uirtl.AccordionItem>
			<uirtl.AccordionItem name="faq-rtl">
				<uirtl.AccordionTrigger>ما هي طرق الدفع التي تقبلونها؟</uirtl.AccordionTrigger>
				<uirtl.AccordionContent>
					نقبل جميع بطاقات الائتمان الرئيسية و PayPal والتحويلات المصرفية. تتم معالجة جميع المدفوعات بأمان من خلال شركاء
					الدفع لدينا.
				</uirtl.AccordionContent>
			</uirtl.AccordionItem>
		</uirtl.Accordion>
	</div>
}

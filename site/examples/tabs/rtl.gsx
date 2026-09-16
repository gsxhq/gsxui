package tabs

import "github.com/gsxhq/gsxui/site/uirtl"

// Rtl mirrors shadcn's own tabs-rtl demo: a 4-tab group ("overview"
// selected), each panel a Card with title/description/content, translated
// to Arabic and wrapped in dir="rtl".
component Rtl() {
	<uirtl.Tabs value="overview" dir="rtl" lang="ar" class="w-full max-w-sm">
		<uirtl.TabsList>
			<uirtl.TabsTrigger value="overview" selected>نظرة عامة</uirtl.TabsTrigger>
			<uirtl.TabsTrigger value="analytics">التحليلات</uirtl.TabsTrigger>
			<uirtl.TabsTrigger value="reports">التقارير</uirtl.TabsTrigger>
			<uirtl.TabsTrigger value="settings">الإعدادات</uirtl.TabsTrigger>
		</uirtl.TabsList>
		<uirtl.TabsContent value="overview" selected>
			<uirtl.Card>
				<uirtl.CardHeader>
					<uirtl.CardTitle>نظرة عامة</uirtl.CardTitle>
					<uirtl.CardDescription>
						عرض مقاييسك الرئيسية وأنشطة المشروع الأخيرة. تتبع التقدم عبر جميع مشاريعك النشطة.
					</uirtl.CardDescription>
				</uirtl.CardHeader>
				<uirtl.CardContent class="text-sm text-muted-foreground">
					لديك ١٢ مشروعًا نشطًا و٣ مهام معلقة.
				</uirtl.CardContent>
			</uirtl.Card>
		</uirtl.TabsContent>
		<uirtl.TabsContent value="analytics">
			<uirtl.Card>
				<uirtl.CardHeader>
					<uirtl.CardTitle>التحليلات</uirtl.CardTitle>
					<uirtl.CardDescription>
						تتبع مقاييس الأداء ومشاركة المستخدمين. راقب الاتجاهات وحدد فرص النمو.
					</uirtl.CardDescription>
				</uirtl.CardHeader>
				<uirtl.CardContent class="text-sm text-muted-foreground">
					زادت مشاهدات الصفحة بنسبة ٢٥٪ مقارنة بالشهر الماضي.
				</uirtl.CardContent>
			</uirtl.Card>
		</uirtl.TabsContent>
		<uirtl.TabsContent value="reports">
			<uirtl.Card>
				<uirtl.CardHeader>
					<uirtl.CardTitle>التقارير</uirtl.CardTitle>
					<uirtl.CardDescription>
						إنشاء وتنزيل تقاريرك التفصيلية. تصدير البيانات بتنسيقات متعددة للتحليل.
					</uirtl.CardDescription>
				</uirtl.CardHeader>
				<uirtl.CardContent class="text-sm text-muted-foreground">لديك ٥ تقارير جاهزة ومتاحة للتصدير.</uirtl.CardContent>
			</uirtl.Card>
		</uirtl.TabsContent>
		<uirtl.TabsContent value="settings">
			<uirtl.Card>
				<uirtl.CardHeader>
					<uirtl.CardTitle>الإعدادات</uirtl.CardTitle>
					<uirtl.CardDescription>إدارة تفضيلات حسابك وخياراته. تخصيص تجربتك لتناسب احتياجاتك.</uirtl.CardDescription>
				</uirtl.CardHeader>
				<uirtl.CardContent class="text-sm text-muted-foreground">تكوين الإشعارات والأمان والسمات.</uirtl.CardContent>
			</uirtl.Card>
		</uirtl.TabsContent>
	</uirtl.Tabs>
}

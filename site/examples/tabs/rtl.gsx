package tabs

import "github.com/gsxhq/gsxui/ui"

// Rtl mirrors shadcn's own tabs-rtl demo: a 4-tab group ("overview"
// selected), each panel a Card with title/description/content, translated
// to Arabic and wrapped in dir="rtl".
component Rtl() {
	<ui.Tabs value="overview" dir="rtl" lang="ar" class="w-full max-w-sm">
		<ui.TabsList>
			<ui.TabsTrigger value="overview" selected>نظرة عامة</ui.TabsTrigger>
			<ui.TabsTrigger value="analytics">التحليلات</ui.TabsTrigger>
			<ui.TabsTrigger value="reports">التقارير</ui.TabsTrigger>
			<ui.TabsTrigger value="settings">الإعدادات</ui.TabsTrigger>
		</ui.TabsList>
		<ui.TabsContent value="overview" selected>
			<ui.Card>
				<ui.CardHeader>
					<ui.CardTitle>نظرة عامة</ui.CardTitle>
					<ui.CardDescription>
						عرض مقاييسك الرئيسية وأنشطة المشروع الأخيرة. تتبع التقدم عبر جميع مشاريعك النشطة.
					</ui.CardDescription>
				</ui.CardHeader>
				<ui.CardContent class="text-sm text-muted-foreground">لديك ١٢ مشروعًا نشطًا و٣ مهام معلقة.</ui.CardContent>
			</ui.Card>
		</ui.TabsContent>
		<ui.TabsContent value="analytics">
			<ui.Card>
				<ui.CardHeader>
					<ui.CardTitle>التحليلات</ui.CardTitle>
					<ui.CardDescription>تتبع مقاييس الأداء ومشاركة المستخدمين. راقب الاتجاهات وحدد فرص النمو.</ui.CardDescription>
				</ui.CardHeader>
				<ui.CardContent class="text-sm text-muted-foreground">
					زادت مشاهدات الصفحة بنسبة ٢٥٪ مقارنة بالشهر الماضي.
				</ui.CardContent>
			</ui.Card>
		</ui.TabsContent>
		<ui.TabsContent value="reports">
			<ui.Card>
				<ui.CardHeader>
					<ui.CardTitle>التقارير</ui.CardTitle>
					<ui.CardDescription>
						إنشاء وتنزيل تقاريرك التفصيلية. تصدير البيانات بتنسيقات متعددة للتحليل.
					</ui.CardDescription>
				</ui.CardHeader>
				<ui.CardContent class="text-sm text-muted-foreground">لديك ٥ تقارير جاهزة ومتاحة للتصدير.</ui.CardContent>
			</ui.Card>
		</ui.TabsContent>
		<ui.TabsContent value="settings">
			<ui.Card>
				<ui.CardHeader>
					<ui.CardTitle>الإعدادات</ui.CardTitle>
					<ui.CardDescription>إدارة تفضيلات حسابك وخياراته. تخصيص تجربتك لتناسب احتياجاتك.</ui.CardDescription>
				</ui.CardHeader>
				<ui.CardContent class="text-sm text-muted-foreground">تكوين الإشعارات والأمان والسمات.</ui.CardContent>
			</ui.Card>
		</ui.TabsContent>
	</ui.Tabs>
}

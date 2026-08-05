package card

import "github.com/gsxhq/gsxui/ui"

// Rtl mirrors shadcn's own card-rtl demo: a login card with a
// CardAction sign-up link, an email/password form, and stacked footer
// buttons — translated to Arabic, wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar" class="w-full max-w-sm">
		<ui.Card>
			<ui.CardHeader>
				<ui.CardTitle>تسجيل الدخول إلى حسابك</ui.CardTitle>
				<ui.CardDescription>أدخل بريدك الإلكتروني أدناه لتسجيل الدخول إلى حسابك</ui.CardDescription>
				<ui.CardAction>
					<ui.Button variant="link">إنشاء حساب</ui.Button>
				</ui.CardAction>
			</ui.CardHeader>
			<ui.CardContent>
				<form>
					<div class="flex flex-col gap-6">
						<div class="grid gap-2">
							<ui.Label for="card-rtl-email">البريد الإلكتروني</ui.Label>
							<ui.Input id="card-rtl-email" type="email" placeholder="m@example.com" autocomplete="email" required/>
						</div>
						<div class="grid gap-2">
							<div class="flex items-center">
								<ui.Label for="card-rtl-password">كلمة المرور</ui.Label>
								<a href="#" class="ms-auto inline-block text-sm underline-offset-4 hover:underline">نسيت كلمة المرور؟</a>
							</div>
							<ui.Input id="card-rtl-password" type="password" autocomplete="current-password" required/>
						</div>
					</div>
				</form>
			</ui.CardContent>
			<ui.CardFooter class="flex-col gap-2">
				<ui.Button type="submit" class="w-full">تسجيل الدخول</ui.Button>
				<ui.Button variant="outline" class="w-full">تسجيل الدخول باستخدام Google</ui.Button>
			</ui.CardFooter>
		</ui.Card>
	</div>
}

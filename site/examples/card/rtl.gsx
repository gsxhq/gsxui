package card

import "github.com/gsxhq/gsxui/site/uirtl"

// Rtl mirrors shadcn's own card-rtl demo: a login card with a
// CardAction sign-up link, an email/password form, and stacked footer
// buttons — translated to Arabic, wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar" class="w-full max-w-sm">
		<uirtl.Card>
			<uirtl.CardHeader>
				<uirtl.CardTitle>تسجيل الدخول إلى حسابك</uirtl.CardTitle>
				<uirtl.CardDescription>أدخل بريدك الإلكتروني أدناه لتسجيل الدخول إلى حسابك</uirtl.CardDescription>
				<uirtl.CardAction>
					<uirtl.Button variant="link">إنشاء حساب</uirtl.Button>
				</uirtl.CardAction>
			</uirtl.CardHeader>
			<uirtl.CardContent>
				<form>
					<div class="flex flex-col gap-6">
						<div class="grid gap-2">
							<uirtl.Label for="card-rtl-email">البريد الإلكتروني</uirtl.Label>
							<uirtl.Input id="card-rtl-email" type="email" placeholder="m@example.com" autocomplete="email" required/>
						</div>
						<div class="grid gap-2">
							<div class="flex items-center">
								<uirtl.Label for="card-rtl-password">كلمة المرور</uirtl.Label>
								<a href="#" class="ms-auto inline-block text-sm underline-offset-4 hover:underline">نسيت كلمة المرور؟</a>
							</div>
							<uirtl.Input id="card-rtl-password" type="password" autocomplete="current-password" required/>
						</div>
					</div>
				</form>
			</uirtl.CardContent>
			<uirtl.CardFooter class="flex-col gap-2">
				<uirtl.Button type="submit" class="w-full">تسجيل الدخول</uirtl.Button>
				<uirtl.Button variant="outline" class="w-full">تسجيل الدخول باستخدام Google</uirtl.Button>
			</uirtl.CardFooter>
		</uirtl.Card>
	</div>
}

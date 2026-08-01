// Package rtl holds the site's example gsx components demonstrating gsxui
// components under dir="rtl" — not tied to a single ui/ component the way
// the other site/examples packages are.
package rtl

import "github.com/gsxhq/gsxui/ui"

// Login is an Arabic sign-in card rendered under dir="rtl": every part —
// Card, Label, Input, Button — adapts through logical Tailwind properties
// alone, no RTL-specific markup of its own. This is the /docs/rtl page's
// "Try it out" demo.
component Login() {
	<div dir="rtl" lang="ar" class="max-w-sm">
		<ui.Card>
			<ui.CardHeader>
				<ui.CardTitle>تسجيل الدخول</ui.CardTitle>
				<ui.CardDescription>أدخل بيانات حسابك للمتابعة</ui.CardDescription>
			</ui.CardHeader>
			<ui.CardContent>
				<form class="flex flex-col gap-4">
					<div class="flex flex-col gap-2">
						<ui.Label for="rtl-login-email">البريد الإلكتروني</ui.Label>
						<ui.Input id="rtl-login-email" type="email" name="email" placeholder="you@example.com" required/>
					</div>
					<div class="flex flex-col gap-2">
						<ui.Label for="rtl-login-password">كلمة المرور</ui.Label>
						<ui.Input id="rtl-login-password" type="password" name="password" required/>
					</div>
					<ui.Button type="submit" class="w-full">دخول</ui.Button>
				</form>
			</ui.CardContent>
			<ui.CardFooter>
				<p class="text-sm text-muted-foreground">ليس لديك حساب؟ إنشاء حساب</p>
			</ui.CardFooter>
		</ui.Card>
	</div>
}

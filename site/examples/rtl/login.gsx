// Package rtl holds the site's example gsx components demonstrating gsxui
// components under dir="rtl" — not tied to a single ui/ component the way
// the other site/examples packages are.
package rtl

import "github.com/gsxhq/gsxui/site/uirtl"

// Login is an Arabic sign-in card rendered under dir="rtl": every part —
// Card, Label, Input, NativeSelect, Button — adapts because the demo
// renders the transformed site/uirtl package, the same output gsxui add
// produces for a project with "rtl": true. This is the /docs/rtl page's
// "Try it out" demo.
component Login() {
	<div dir="rtl" lang="ar" class="max-w-sm">
		<uirtl.Card>
			<uirtl.CardHeader>
				<uirtl.CardTitle>تسجيل الدخول</uirtl.CardTitle>
				<uirtl.CardDescription>أدخل بيانات حسابك للمتابعة</uirtl.CardDescription>
			</uirtl.CardHeader>
			<uirtl.CardContent>
				<form class="flex flex-col gap-4">
					<div class="flex flex-col gap-2">
						<uirtl.Label for="rtl-login-email">البريد الإلكتروني</uirtl.Label>
						<uirtl.Input
							id="rtl-login-email"
							type="email"
							name="email"
							placeholder="you@example.com"
							autocomplete="email"
							required
						/>
					</div>
					<div class="flex flex-col gap-2">
						<uirtl.Label for="rtl-login-password">كلمة المرور</uirtl.Label>
						<uirtl.Input
							id="rtl-login-password"
							type="password"
							name="password"
							autocomplete="current-password"
							required
						/>
					</div>
					<div class="flex flex-col gap-2">
						<uirtl.Label for="rtl-login-language">اللغة</uirtl.Label>
						<uirtl.NativeSelect id="rtl-login-language" name="language">
							<uirtl.NativeSelectOption value="ar" selected={true}>العربية</uirtl.NativeSelectOption>
							<uirtl.NativeSelectOption value="en">English</uirtl.NativeSelectOption>
						</uirtl.NativeSelect>
					</div>
					<uirtl.Button type="submit" class="w-full">دخول</uirtl.Button>
				</form>
			</uirtl.CardContent>
			<uirtl.CardFooter>
				<p class="text-sm text-muted-foreground">ليس لديك حساب؟ إنشاء حساب</p>
			</uirtl.CardFooter>
		</uirtl.Card>
	</div>
}

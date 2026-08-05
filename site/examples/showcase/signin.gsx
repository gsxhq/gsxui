// Package showcase holds the composed demo cards rendered on the site's
// landing page. Unlike the per-component packages next door, these are
// not registered docs examples — they exist to show several ui components
// working together in one realistic piece of UI.
package showcase

import "github.com/gsxhq/gsxui/ui"

// SignInCard composes Card, Label, Input, Checkbox and Button into the
// classic email/password sign-in form.
component SignInCard() {
	<ui.Card>
		<ui.CardHeader>
			<ui.CardTitle>Sign in</ui.CardTitle>
			<ui.CardDescription>Enter your email below to sign in to your account.</ui.CardDescription>
		</ui.CardHeader>
		<ui.CardContent>
			<form>
				<div class="flex flex-col gap-6">
					<div class="grid gap-2">
						<ui.Label for="home-showcase-email">Email</ui.Label>
						<ui.Input id="home-showcase-email" type="email" placeholder="m@example.com" autocomplete="email" required/>
					</div>
					<div class="grid gap-2">
						<div class="flex items-center">
							<ui.Label for="home-showcase-password">Password</ui.Label>
							<a href="#" class="ms-auto inline-block text-sm underline-offset-4 hover:underline">Forgot password?</a>
						</div>
						<ui.Input id="home-showcase-password" type="password" autocomplete="current-password" required/>
					</div>
					<div class="flex items-center gap-2">
						<ui.Checkbox id="home-showcase-remember"/>
						<ui.Label for="home-showcase-remember">Remember me</ui.Label>
					</div>
				</div>
			</form>
		</ui.CardContent>
		<ui.CardFooter class="flex-col gap-2">
			<ui.Button class="w-full">Sign in</ui.Button>
			<ui.Button variant="ghost" class="w-full">Create an account</ui.Button>
		</ui.CardFooter>
	</ui.Card>
}

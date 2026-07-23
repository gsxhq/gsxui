package input

import (
	"github.com/gsxhq/gsxui/ui"
)

// FormRow composes Label, Input, Checkbox, and a submit Button into a
// realistic form row — Label/for pairs with Input/Checkbox by id, the
// pattern most real forms actually reach for.
component FormRow() {
	<form class="flex max-w-sm flex-col gap-4">
		<div class="flex flex-col gap-2">
			<ui.Label for="form-row-email">Email</ui.Label>
			<ui.Input id="form-row-email" type="email" name="email" placeholder="you@example.com" required/>
		</div>
		<div class="flex items-center gap-2">
			<ui.Checkbox id="form-row-remember" name="remember"/>
			<ui.Label for="form-row-remember">Remember me</ui.Label>
		</div>
		<ui.Button type="submit">Sign in</ui.Button>
	</form>
}

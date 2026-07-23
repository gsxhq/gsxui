package input

import (
	uibutton "github.com/gsxhq/gsxui/ui/button"
	uicheckbox "github.com/gsxhq/gsxui/ui/checkbox"
	uiinput "github.com/gsxhq/gsxui/ui/input"
	uilabel "github.com/gsxhq/gsxui/ui/label"
)

// FormRow composes Label, Input, Checkbox, and a submit Button into a
// realistic form row — Label/for pairs with Input/Checkbox by id, the
// pattern most real forms actually reach for.
component FormRow() {
	<form class="flex max-w-sm flex-col gap-4">
		<div class="flex flex-col gap-2">
			<uilabel.Label for="form-row-email">Email</uilabel.Label>
			<uiinput.Input id="form-row-email" type="email" name="email" placeholder="you@example.com" required/>
		</div>
		<div class="flex items-center gap-2">
			<uicheckbox.Checkbox id="form-row-remember" name="remember"/>
			<uilabel.Label for="form-row-remember">Remember me</uilabel.Label>
		</div>
		<uibutton.Button type="submit">Sign in</uibutton.Button>
	</form>
}

package field

import "github.com/gsxhq/gsxui/ui"

// Invalid renders a Field in its error state: aria-invalid on the control
// plus a FieldError message, the same shape as ui/input's own "states"
// example (see site/examples/input/states.gsx).
component Invalid() {
	<ui.FieldSet>
		<ui.FieldLegend>Account</ui.FieldLegend>
		<ui.FieldGroup>
			<ui.Field data-invalid="true">
				<ui.FieldLabel for="email">Email</ui.FieldLabel>
				<ui.Input id="email" type="email" aria-invalid="true" value="not-an-email"/>
				<ui.FieldError>Enter a valid email address.</ui.FieldError>
			</ui.Field>
		</ui.FieldGroup>
	</ui.FieldSet>
}

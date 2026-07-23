// Package field holds the site's example gsx components for ui/field.
package field

import "github.com/gsxhq/gsxui/ui"

// Basic renders a realistic small form: a FieldSet with a legend, wrapping a
// FieldGroup of two vertical Fields (label + control + description) split
// by a FieldSeparator.
component Basic() {
	<ui.FieldSet>
		<ui.FieldLegend>Profile</ui.FieldLegend>
		<ui.FieldGroup>
			<ui.Field>
				<ui.FieldLabel for="name">Name</ui.FieldLabel>
				<ui.Input id="name" placeholder="Jamie Lee"/>
				<ui.FieldDescription>Your full name, as it appears on your profile.</ui.FieldDescription>
			</ui.Field>
			<ui.FieldSeparator/>
			<ui.Field>
				<ui.FieldLabel for="bio">Bio</ui.FieldLabel>
				<ui.Textarea id="bio" placeholder="Tell us about yourself"/>
				<ui.FieldDescription>Shown to other members on your public profile.</ui.FieldDescription>
			</ui.Field>
		</ui.FieldGroup>
	</ui.FieldSet>
}

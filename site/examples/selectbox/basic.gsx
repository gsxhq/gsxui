// Package selectbox holds the site's example gsx components for ui/selectbox.
package selectbox

import "github.com/gsxhq/gsxui/ui"

// Basic renders a Select with plain SelectOption children, one
// pre-selected.
component Basic() {
	<ui.Select name="fruit">
		<ui.SelectOption value="apple">Apple</ui.SelectOption>
		<ui.SelectOption value="banana" selected>Banana</ui.SelectOption>
		<ui.SelectOption value="cherry">Cherry</ui.SelectOption>
	</ui.Select>
}

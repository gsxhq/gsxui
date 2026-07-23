// Package selectbox holds the site's example gsx components for ui/selectbox.
package selectbox

import uiselect "github.com/gsxhq/gsxui/ui/selectbox"

// Basic renders a Select with plain SelectOption children, one
// pre-selected.
component Basic() {
	<uiselect.Select name="fruit">
		<uiselect.SelectOption value="apple">Apple</uiselect.SelectOption>
		<uiselect.SelectOption value="banana" selected>Banana</uiselect.SelectOption>
		<uiselect.SelectOption value="cherry">Cherry</uiselect.SelectOption>
	</uiselect.Select>
}

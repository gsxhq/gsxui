// Package selectbox holds the site's example gsx components for
// ui/selectbox. Each example is a real, compiled gsx component — the exact
// source below is what the component page displays AND what it renders, so
// source shown is source run; the examples_test.go drift test enforces
// they can't diverge.
package selectbox

import uiselect "github.com/gsxhq/gsxui/ui/selectbox"

// Basic renders a Select with plain SelectOption children, one
// pre-selected.
component Basic() {
	<uiselect.Select name="fruit">
		<uiselect.SelectOption value="apple" selected={false} disabled={false}>Apple</uiselect.SelectOption>
		<uiselect.SelectOption value="banana" selected={true} disabled={false}>Banana</uiselect.SelectOption>
		<uiselect.SelectOption value="cherry" selected={false} disabled={false}>Cherry</uiselect.SelectOption>
	</uiselect.Select>
}

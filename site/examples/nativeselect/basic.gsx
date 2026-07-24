// Package nativeselect holds the site's example gsx components for
// ui.NativeSelect (hyphen stripped: Go package names can't contain one,
// same as button-group's buttongroup / context-menu's contextmenu).
package nativeselect

import "github.com/gsxhq/gsxui/ui"

// Basic renders a NativeSelect with plain NativeSelectOption children, one
// pre-selected.
component Basic() {
	<ui.NativeSelect name="fruit">
		<ui.NativeSelectOption value="apple">Apple</ui.NativeSelectOption>
		<ui.NativeSelectOption value="banana" selected>Banana</ui.NativeSelectOption>
		<ui.NativeSelectOption value="cherry">Cherry</ui.NativeSelectOption>
	</ui.NativeSelect>
}

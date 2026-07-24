package nativeselect

import "github.com/gsxhq/gsxui/ui"

// Groups renders a NativeSelect whose options are organized under
// NativeSelectGroup (native <optgroup>), plus one disabled option.
component Groups() {
	<ui.NativeSelect name="timezone">
		<ui.NativeSelectGroup label="Americas">
			<ui.NativeSelectOption value="est" selected>Eastern</ui.NativeSelectOption>
			<ui.NativeSelectOption value="pst">Pacific</ui.NativeSelectOption>
		</ui.NativeSelectGroup>
		<ui.NativeSelectGroup label="Europe">
			<ui.NativeSelectOption value="cet">Central</ui.NativeSelectOption>
			<ui.NativeSelectOption value="gmt" disabled>Greenwich (unavailable)</ui.NativeSelectOption>
		</ui.NativeSelectGroup>
	</ui.NativeSelect>
}

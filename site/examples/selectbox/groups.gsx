package selectbox

import "github.com/gsxhq/gsxui/ui"

// Groups renders a Select whose options are organized under SelectGroup
// (native <optgroup>), plus one disabled option.
component Groups() {
	<ui.Select name="timezone">
		<ui.SelectGroup label="Americas">
			<ui.SelectOption value="est" selected>Eastern</ui.SelectOption>
			<ui.SelectOption value="pst">Pacific</ui.SelectOption>
		</ui.SelectGroup>
		<ui.SelectGroup label="Europe">
			<ui.SelectOption value="cet">Central</ui.SelectOption>
			<ui.SelectOption value="gmt" disabled>Greenwich (unavailable)</ui.SelectOption>
		</ui.SelectGroup>
	</ui.Select>
}

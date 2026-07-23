package selectbox

import uiselect "github.com/gsxhq/gsxui/ui/selectbox"

// Groups renders a Select whose options are organized under SelectGroup
// (native <optgroup>), plus one disabled option.
component Groups() {
	<uiselect.Select name="timezone">
		<uiselect.SelectGroup label="Americas">
			<uiselect.SelectOption value="est" selected={true} disabled={false}>Eastern</uiselect.SelectOption>
			<uiselect.SelectOption value="pst" selected={false} disabled={false}>Pacific</uiselect.SelectOption>
		</uiselect.SelectGroup>
		<uiselect.SelectGroup label="Europe">
			<uiselect.SelectOption value="cet" selected={false} disabled={false}>Central</uiselect.SelectOption>
			<uiselect.SelectOption value="gmt" selected={false} disabled={true}>Greenwich (unavailable)</uiselect.SelectOption>
		</uiselect.SelectGroup>
	</uiselect.Select>
}

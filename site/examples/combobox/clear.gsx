package combobox

import "github.com/gsxhq/gsxui/ui"

// Clear mirrors shadcn's own combobox-clear.tsx: showClear renders the X
// button instead of the chevron trigger (ComboboxInput's own
// group-has-data-[slot=combobox-clear]/input-group:hidden rule hides the
// trigger whenever a clear button is present), and a pre-selected value
// server-renders the checked item — combobox.js reflects it into the
// input's displayed text at init (see ui/combobox.js's own init()).
component Clear() {
	<ui.Combobox name="language" value="go">
		<ui.ComboboxInput placeholder="Search language..." showClear class="w-[220px]"/>
		<ui.ComboboxContent>
			<ui.ComboboxList>
				<ui.ComboboxEmpty>No language found.</ui.ComboboxEmpty>
				<ui.ComboboxItem value="go" selected={true}>Go</ui.ComboboxItem>
				<ui.ComboboxItem value="rust" selected={false}>Rust</ui.ComboboxItem>
				<ui.ComboboxItem value="typescript" selected={false}>TypeScript</ui.ComboboxItem>
				<ui.ComboboxItem value="python" selected={false}>Python</ui.ComboboxItem>
			</ui.ComboboxList>
		</ui.ComboboxContent>
	</ui.Combobox>
}

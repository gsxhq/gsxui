package combobox

import "github.com/gsxhq/gsxui/ui"

// Clear mirrors shadcn's own combobox-clear.tsx: showClear renders the X
// button instead of the chevron trigger (ComboboxInput's own
// group-has-data-[slot=combobox-clear]/input-group:hidden rule hides the
// trigger whenever a clear button is present), and a pre-selected value
// server-renders the checked item.
//
// FIX (review round 2, IMPORTANT): a server-selected value must REFLECT in
// the DOM on first paint (docs/superpowers/specs/2026-07-24-tier4-batch-a-
// design.md §4) — combobox.js's init() only SEEDS the input's displayed
// label from the checked item once JS has run (docs/jsx-parity.md
// `## combobox`'s own "reopening filtered to the committed label" FIX
// entry), so a no-JS request rendered an empty input despite the hidden
// bridge already posting "go" correctly, and a JS request flashed
// placeholder-to-label on load. ComboboxInput's attrs already reach the
// inner <input> (attrs.Without("class") — see its own ADAPT doc comment),
// so passing the picked item's own label through as `value` here needs no
// signature change; combobox.js's init() is now a no-op fallback for
// callers who don't supply it (it only seeds when the input is still
// empty).
component Clear() {
	<ui.Combobox name="language" value="go">
		<ui.ComboboxInput placeholder="Search language..." showClear value="Go" class="w-[220px]"/>
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

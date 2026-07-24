// Package selectbox holds the site's example gsx components for the custom-
// listbox ui.Select (dir name "selectbox" — "select" is a Go keyword and
// can't name a package; the registered component key is still "select").
package selectbox

import "github.com/gsxhq/gsxui/ui"

// Basic is a single-group Select with a label, five items and a placeholder;
// the trigger is width-pinned (w-[180px]) the way shadcn callers size the
// Trigger. name="fruit" renders the hidden native <select> form bridge.
component Basic() {
	<ui.Select name="fruit">
		<ui.SelectTrigger class="w-[180px]">
			<ui.SelectValue placeholder="Select a fruit"/>
		</ui.SelectTrigger>
		<ui.SelectContent>
			<ui.SelectGroup>
				<ui.SelectLabel>Fruits</ui.SelectLabel>
				<ui.SelectItem value="apple">Apple</ui.SelectItem>
				<ui.SelectItem value="banana">Banana</ui.SelectItem>
				<ui.SelectItem value="blueberry">Blueberry</ui.SelectItem>
				<ui.SelectItem value="grapes">Grapes</ui.SelectItem>
				<ui.SelectItem value="pineapple">Pineapple</ui.SelectItem>
			</ui.SelectGroup>
		</ui.SelectContent>
	</ui.Select>
}

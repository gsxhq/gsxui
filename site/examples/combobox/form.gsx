package combobox

import "github.com/gsxhq/gsxui/ui"

// Form demonstrates the hidden sr-only <input> form bridge (Combobox's own
// name/value params — ui/combobox.gsx's package doc comment MECHANISM
// entry): submitting this form, JS on or off, carries "framework" in the
// POST body, the same non-JS-safe guarantee ## select's own bridge gives
// ui.Select — see site/examples/input/form-row.gsx for the plain-input
// version of this pattern.
component Form() {
	<form class="flex max-w-xs flex-col gap-4">
		<div class="flex flex-col gap-2">
			<ui.Label for="combobox-form-framework">Framework</ui.Label>
			<ui.Combobox name="framework" value="">
				<ui.ComboboxInput id="combobox-form-framework" placeholder="Search framework..." showTrigger/>
				<ui.ComboboxContent>
					<ui.ComboboxList>
						<ui.ComboboxEmpty>No framework found.</ui.ComboboxEmpty>
						<ui.ComboboxItem value="next.js" selected={false}>Next.js</ui.ComboboxItem>
						<ui.ComboboxItem value="sveltekit" selected={false}>SvelteKit</ui.ComboboxItem>
						<ui.ComboboxItem value="nuxt.js" selected={false}>Nuxt.js</ui.ComboboxItem>
						<ui.ComboboxItem value="remix" selected={false}>Remix</ui.ComboboxItem>
					</ui.ComboboxList>
				</ui.ComboboxContent>
			</ui.Combobox>
		</div>
		<ui.Button type="submit">Continue</ui.Button>
	</form>
}

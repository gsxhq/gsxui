package combobox

import "github.com/gsxhq/gsxui/ui"

var frameworks = []struct{ Value, Label string }{
	{"next.js", "Next.js"},
	{"sveltekit", "SvelteKit"},
	{"nuxt.js", "Nuxt.js"},
	{"remix", "Remix"},
	{"astro", "Astro"},
}

// Basic mirrors shadcn's own combobox-basic.tsx: a five-item single-select
// with no groups and no clear button. Typing filters via ui/combobox.js's
// contains() matcher (hide/show only, no reordering — see its own header);
// name="framework" renders the hidden sr-only <input> form bridge.
component Basic() {
	<ui.Combobox name="framework" value="">
		<ui.ComboboxInput placeholder="Search framework..." showTrigger class="w-[220px]"/>
		<ui.ComboboxContent>
			<ui.ComboboxList>
				<ui.ComboboxEmpty>No framework found.</ui.ComboboxEmpty>
				{ for _, f := range frameworks {
					<ui.ComboboxItem value={f.Value} selected={false}>{ f.Label }</ui.ComboboxItem>
				} }
			</ui.ComboboxList>
		</ui.ComboboxContent>
	</ui.Combobox>
}

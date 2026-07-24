package combobox

import "github.com/gsxhq/gsxui/ui"

var fruitGroups = []struct {
	Label string
	Items []struct{ Value, Label string }
}{
	{"Citrus", []struct{ Value, Label string }{
		{"orange", "Orange"},
		{"lemon", "Lemon"},
		{"lime", "Lime"},
	}},
	{"Berries", []struct{ Value, Label string }{
		{"strawberry", "Strawberry"},
		{"blueberry", "Blueberry"},
		{"raspberry", "Raspberry"},
	}},
}

// Groups mirrors shadcn's own combobox-groups.tsx: nested groups, each with
// a ComboboxLabel heading, separated by a ComboboxSeparator between groups
// (not after the last one) — the same shape ## select's own scrollable
// example uses for SelectGroup/SelectSeparator.
component Groups() {
	<ui.Combobox name="fruit" value="">
		<ui.ComboboxInput placeholder="Search fruit..." showTrigger class="w-[220px]"/>
		<ui.ComboboxContent>
			<ui.ComboboxList>
				<ui.ComboboxEmpty>No fruit found.</ui.ComboboxEmpty>
				{ for i, g := range fruitGroups {
					<ui.ComboboxGroup>
						<ui.ComboboxLabel>{ g.Label }</ui.ComboboxLabel>
						{ for _, f := range g.Items {
							<ui.ComboboxItem value={f.Value} selected={false}>{ f.Label }</ui.ComboboxItem>
						} }
					</ui.ComboboxGroup>
					{ if i < len(fruitGroups)-1 {
						<ui.ComboboxSeparator/>
					} }
				} }
			</ui.ComboboxList>
		</ui.ComboboxContent>
	</ui.Combobox>
}

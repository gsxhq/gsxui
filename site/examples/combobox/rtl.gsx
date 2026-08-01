package combobox

import "github.com/gsxhq/gsxui/ui"

var categories = []struct{ Value, Label string }{
	{"technology", "التكنولوجيا"},
	{"design", "التصميم"},
	{"business", "الأعمال"},
	{"marketing", "التسويق"},
	{"education", "التعليم"},
	{"health", "الصحة"},
}

// Rtl mirrors shadcn's own combobox-rtl demo's Arabic category list, but
// keeps Basic's single-select shape (shadcn's version is a multi-select
// chips combobox — ComboboxChips/ComboboxChip/ComboboxValue aren't part of
// this directory's Basic and are DEVIATED from here rather than invented).
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.Combobox name="category" value="">
			<ui.ComboboxInput placeholder="أضف فئات" showTrigger class="w-[220px]"/>
			<ui.ComboboxContent>
				<ui.ComboboxList>
					<ui.ComboboxEmpty>لم يتم العثور على فئات.</ui.ComboboxEmpty>
					{ for _, c := range categories {
						<ui.ComboboxItem value={c.Value} selected={false}>{ c.Label }</ui.ComboboxItem>
					} }
				</ui.ComboboxList>
			</ui.ComboboxContent>
		</ui.Combobox>
	</div>
}

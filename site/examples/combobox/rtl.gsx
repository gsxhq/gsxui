package combobox

import "github.com/gsxhq/gsxui/site/uirtl"

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
		<uirtl.Combobox name="category" value="">
			<uirtl.ComboboxInput placeholder="أضف فئات" showTrigger class="w-[220px]"/>
			<uirtl.ComboboxContent>
				<uirtl.ComboboxList>
					<uirtl.ComboboxEmpty>لم يتم العثور على فئات.</uirtl.ComboboxEmpty>
					{ for _, c := range categories {
						<uirtl.ComboboxItem value={c.Value} selected={false}>{ c.Label }</uirtl.ComboboxItem>
					} }
				</uirtl.ComboboxList>
			</uirtl.ComboboxContent>
		</uirtl.Combobox>
	</div>
}

package togglegroup

import "github.com/gsxhq/gsxui/site/uirtl"

// Rtl mirrors shadcn's own toggle-group-rtl demo: an outline single-select
// group of three text items ("list" selected), translated to Arabic and
// wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar">
		<uirtl.ToggleGroup groupType="single" variant="outline" aria-label="View mode">
			<uirtl.ToggleGroupItem groupType="single" variant="outline" value="list" pressed aria-label="قائمة">
				قائمة
			</uirtl.ToggleGroupItem>
			<uirtl.ToggleGroupItem groupType="single" variant="outline" value="grid" aria-label="شبكة">
				شبكة
			</uirtl.ToggleGroupItem>
			<uirtl.ToggleGroupItem groupType="single" variant="outline" value="cards" aria-label="بطاقات">
				بطاقات
			</uirtl.ToggleGroupItem>
		</uirtl.ToggleGroup>
	</div>
}

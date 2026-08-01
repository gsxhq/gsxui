package togglegroup

import "github.com/gsxhq/gsxui/ui"

// Rtl mirrors shadcn's own toggle-group-rtl demo: an outline single-select
// group of three text items ("list" selected), translated to Arabic and
// wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.ToggleGroup groupType="single" variant="outline" aria-label="View mode">
			<ui.ToggleGroupItem groupType="single" variant="outline" value="list" pressed aria-label="قائمة">
				قائمة
			</ui.ToggleGroupItem>
			<ui.ToggleGroupItem groupType="single" variant="outline" value="grid" aria-label="شبكة">شبكة</ui.ToggleGroupItem>
			<ui.ToggleGroupItem groupType="single" variant="outline" value="cards" aria-label="بطاقات">
				بطاقات
			</ui.ToggleGroupItem>
		</ui.ToggleGroup>
	</div>
}

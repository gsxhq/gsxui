package buttongroup

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Rtl mirrors Basic's outline-pair and quantity-stepper groups, translated
// to Arabic and wrapped in dir="rtl" — the back-arrow icon flips via
// rtl:rotate-180 like shadcn's own button-group-rtl demo. That upstream
// demo also nests a DropdownMenu (with a radio-group submenu) inside a
// ButtonGroup; gsxui's button-group example set doesn't compose
// dropdown-menu (see Basic), so this example keeps the two established
// group shapes instead of introducing that composition here.
component Rtl() {
	<div dir="rtl" lang="ar" class="flex flex-wrap items-start gap-6">
		<ui.ButtonGroup>
			<ui.Button variant="outline" size="icon" aria-label="Go Back">
				<icon.ArrowLeft class="rtl:rotate-180"/>
			</ui.Button>
		</ui.ButtonGroup>
		<ui.ButtonGroup>
			<ui.Button variant="outline">أرشفة</ui.Button>
			<ui.Button variant="outline">تقرير</ui.Button>
		</ui.ButtonGroup>
		<ui.ButtonGroup aria-label="Quantity">
			<ui.Button variant="outline" size="icon" aria-label="Decrease quantity">
				<icon.Minus/>
			</ui.Button>
			<ui.ButtonGroupText>42</ui.ButtonGroupText>
			<ui.Button variant="outline" size="icon" aria-label="Increase quantity">
				<icon.Plus/>
			</ui.Button>
		</ui.ButtonGroup>
	</div>
}

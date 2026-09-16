package buttongroup

import (
	"github.com/gsxhq/gsxui/site/uirtl"
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
		<uirtl.ButtonGroup>
			<uirtl.Button variant="outline" size="icon" aria-label="Go Back">
				<icon.ArrowLeft class="rtl:rotate-180"/>
			</uirtl.Button>
		</uirtl.ButtonGroup>
		<uirtl.ButtonGroup>
			<uirtl.Button variant="outline">أرشفة</uirtl.Button>
			<uirtl.Button variant="outline">تقرير</uirtl.Button>
		</uirtl.ButtonGroup>
		<uirtl.ButtonGroup aria-label="Quantity">
			<uirtl.Button variant="outline" size="icon" aria-label="Decrease quantity">
				<icon.Minus/>
			</uirtl.Button>
			<uirtl.ButtonGroupText>42</uirtl.ButtonGroupText>
			<uirtl.Button variant="outline" size="icon" aria-label="Increase quantity">
				<icon.Plus/>
			</uirtl.Button>
		</uirtl.ButtonGroup>
	</div>
}

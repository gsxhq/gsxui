// Package label holds the site's example gsx components for ui/label.
package label

import (
	"github.com/gsxhq/gsxui/ui"
)

// Rtl mirrors shadcn's own label-rtl demo: a Checkbox paired with a Label
// via matching id/for, translated to Arabic and wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar" class="flex items-center gap-2">
		<ui.Checkbox id="label-rtl-terms"/>
		<ui.Label for="label-rtl-terms">قبول الشروط والأحكام</ui.Label>
	</div>
}

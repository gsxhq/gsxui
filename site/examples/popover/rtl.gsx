// Package popover holds the site's example gsx components for ui/popover.
package popover

import (
	"github.com/gsxhq/gsxui/ui"
)

// Rtl mirrors this dir's own Basic (shadcn's popover-demo.tsx shape — an
// Open-popover trigger and a w-80 dimensions form), translated to Arabic
// and wrapped in dir="rtl". shadcn's own popover-rtl demo instead sweeps
// every side/logical-side popover placement across a row of triggers;
// this dir's ui.PopoverContent has no side prop to sweep (see basic.gsx),
// so the Basic shape is kept and translated rather than invented.
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.Popover>
			<ui.Button
				variant="outline"
				data-gsxui-slot-popover-trigger
				aria-expanded="false"
			>
				فتح النافذة المنبثقة
			</ui.Button>
			<ui.PopoverContent class="w-80">
				<div class="grid gap-4">
					<div class="space-y-2">
						<h4 class="leading-none font-medium">الأبعاد</h4>
						<p class="text-sm text-muted-foreground">تعيين الأبعاد للطبقة.</p>
					</div>
					<div class="grid gap-2">
						<div class="grid grid-cols-3 items-center gap-4">
							<ui.Label for="width-rtl">العرض</ui.Label>
							<ui.Input id="width-rtl" value="100%" class="col-span-2"/>
						</div>
						<div class="grid grid-cols-3 items-center gap-4">
							<ui.Label for="maxWidth-rtl">أقصى عرض</ui.Label>
							<ui.Input id="maxWidth-rtl" value="300px" class="col-span-2"/>
						</div>
						<div class="grid grid-cols-3 items-center gap-4">
							<ui.Label for="height-rtl">الارتفاع</ui.Label>
							<ui.Input id="height-rtl" value="25px" class="col-span-2"/>
						</div>
						<div class="grid grid-cols-3 items-center gap-4">
							<ui.Label for="maxHeight-rtl">أقصى ارتفاع</ui.Label>
							<ui.Input id="maxHeight-rtl" value="none" class="col-span-2"/>
						</div>
					</div>
				</div>
			</ui.PopoverContent>
		</ui.Popover>
	</div>
}

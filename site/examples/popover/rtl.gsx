// Package popover holds the site's example gsx components for ui/popover.
package popover

import (
	"github.com/gsxhq/gsxui/site/uirtl"
)

// Rtl mirrors this dir's own Basic (shadcn's popover-demo.tsx shape — an
// Open-popover trigger and a w-80 dimensions form), translated to Arabic
// and wrapped in dir="rtl". shadcn's own popover-rtl demo instead sweeps
// every side/logical-side popover placement across a row of triggers;
// this dir's uirtl.PopoverContent has no side prop to sweep (see basic.gsx),
// so the Basic shape is kept and translated rather than invented.
component Rtl() {
	<div dir="rtl" lang="ar">
		<uirtl.Popover>
			<uirtl.Button
				variant="outline"
				data-gsxui-slot-popover-trigger
				aria-expanded="false"
			>
				فتح النافذة المنبثقة
			</uirtl.Button>
			<uirtl.PopoverContent class="w-80">
				<div class="grid gap-4">
					<div class="space-y-2">
						<h4 class="leading-none font-medium">الأبعاد</h4>
						<p class="text-sm text-muted-foreground">تعيين الأبعاد للطبقة.</p>
					</div>
					<div class="grid gap-2">
						<div class="grid grid-cols-3 items-center gap-4">
							<uirtl.Label for="width-rtl">العرض</uirtl.Label>
							<uirtl.Input id="width-rtl" value="100%" class="col-span-2"/>
						</div>
						<div class="grid grid-cols-3 items-center gap-4">
							<uirtl.Label for="maxWidth-rtl">أقصى عرض</uirtl.Label>
							<uirtl.Input id="maxWidth-rtl" value="300px" class="col-span-2"/>
						</div>
						<div class="grid grid-cols-3 items-center gap-4">
							<uirtl.Label for="height-rtl">الارتفاع</uirtl.Label>
							<uirtl.Input id="height-rtl" value="25px" class="col-span-2"/>
						</div>
						<div class="grid grid-cols-3 items-center gap-4">
							<uirtl.Label for="maxHeight-rtl">أقصى ارتفاع</uirtl.Label>
							<uirtl.Input id="maxHeight-rtl" value="none" class="col-span-2"/>
						</div>
					</div>
				</div>
			</uirtl.PopoverContent>
		</uirtl.Popover>
	</div>
}

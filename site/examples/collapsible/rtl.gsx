package collapsible

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Rtl mirrors shadcn's own collapsible-rtl demo (an order-status panel
// rather than the starred-repos list Basic uses), translated to Arabic and
// wrapped in dir="rtl". Structure stays Basic's: chevron trigger row, one
// always-visible line, two more revealed by expanding.
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.Collapsible class="flex w-[350px] flex-col gap-2">
			<ui.CollapsibleTrigger class="flex cursor-default items-center justify-between gap-4 px-4">
				<h4 class="text-sm font-semibold">الطلب #4189</h4>
				<span
					aria-hidden="true"
					class="inline-flex size-8 shrink-0 items-center justify-center rounded-md transition-colors hover:bg-accent hover:text-accent-foreground [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4"
				>
					<icon.ChevronsUpDown/>
				</span>
			</ui.CollapsibleTrigger>
			<div class="flex items-center justify-between rounded-md border px-4 py-2 text-sm">
				<span class="text-muted-foreground">الحالة</span>
				<span class="font-medium">تم الشحن</span>
			</div>
			<ui.CollapsibleContent class="flex flex-col gap-2">
				<div class="rounded-md border px-4 py-2 text-sm">
					<p class="font-medium">عنوان الشحن</p>
					<p class="text-muted-foreground">100 Market St, San Francisco</p>
				</div>
				<div class="rounded-md border px-4 py-2 text-sm">
					<p class="font-medium">العناصر</p>
					<p class="text-muted-foreground">2x سماعات الاستوديو</p>
				</div>
			</ui.CollapsibleContent>
		</ui.Collapsible>
	</div>
}

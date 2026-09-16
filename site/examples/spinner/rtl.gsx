package spinner

import "github.com/gsxhq/gsxui/site/uirtl"

// Rtl mirrors shadcn's own spinner-rtl demo: a muted Item row pairing a
// Spinner with a processing-payment title and a trailing amount, translated
// to Arabic and wrapped in dir="rtl". shadcn's rtl demo also formats the
// amount with Arabic-Indic numerals via Intl; this port keeps the literal
// "١٠٠٫٠٠ دولار" string shadcn's own Arabic translation table supplies.
component Rtl() {
	<div dir="rtl" lang="ar" class="flex w-full max-w-xs flex-col gap-4">
		<uirtl.Item variant="muted">
			<uirtl.ItemMedia>
				<uirtl.Spinner/>
			</uirtl.ItemMedia>
			<uirtl.ItemContent>
				<uirtl.ItemTitle class="line-clamp-1">جاري معالجة الدفع...</uirtl.ItemTitle>
			</uirtl.ItemContent>
			<uirtl.ItemContent class="flex-none justify-end">
				<span class="text-sm tabular-nums">١٠٠.٠٠ دولار</span>
			</uirtl.ItemContent>
		</uirtl.Item>
	</div>
}

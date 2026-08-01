package aspectratio

import "github.com/gsxhq/gsxui/ui"

// Rtl mirrors shadcn's own aspect-ratio-rtl demo: a figure/figcaption
// pair with an Arabic caption, wrapped in dir="rtl". Like Basic, it uses a
// bordered muted-background box as a stand-in for shadcn's real <Image>.
component Rtl() {
	<figure dir="rtl" lang="ar" class="w-full max-w-sm">
		<ui.AspectRatio
			ratio="16 / 9"
			class="flex items-center justify-center rounded-lg border bg-muted text-sm text-muted-foreground"
		>
			16 / 9
		</ui.AspectRatio>
		<figcaption class="mt-2 text-center text-sm text-muted-foreground">منظر طبيعي جميل</figcaption>
	</figure>
}

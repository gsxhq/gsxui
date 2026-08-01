package slider

import "github.com/gsxhq/gsxui/ui"

// Rtl mirrors shadcn's own slider-rtl demo: the same single-thumb slider as
// Basic (value 75 here, matching shadcn's rtl demo rather than Basic's 50),
// wrapped in dir="rtl" so the native range input's own RTL handling takes
// over the thumb's visual direction.
component Rtl() {
	<div dir="rtl" lang="ar" class="mx-auto w-full max-w-xs">
		<ui.Slider value={75} min={0} max={100} step={1} class="w-full" aria-label="Volume" dir="rtl"/>
	</div>
}

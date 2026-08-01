// Package progress holds the site's example gsx components for ui/progress.
package progress

import "github.com/gsxhq/gsxui/ui"

// Rtl mirrors shadcn's own progress-rtl demo's single upload-progress bar
// (translated label, value 56), wrapped in dir="rtl". shadcn's demo also
// composes a ProgressLabel/ProgressValue pair and Arabic-Indic numeral
// formatting; this dir's ui.Progress exposes neither (see basic.gsx),
// so a single bare Progress at the same value stands in.
component Rtl() {
	<div dir="rtl" lang="ar" class="w-full max-w-sm">
		<ui.Progress value={56}/>
	</div>
}

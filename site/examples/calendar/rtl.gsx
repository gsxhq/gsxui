package calendar

import (
	"time"

	"github.com/gsxhq/gsxui/ui"
)

// RtlDefaultMonth mirrors Basic's own DefaultMonth (2026-01, never
// time.Now()).
var RtlDefaultMonth = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// Rtl renders the same single-mode grid as Basic, wrapped in dir="rtl":
// the week runs right-to-left and the caption's prev/next icons flip via
// rtl:rotate-180, with no RTL-specific props of its own.
component Rtl() {
	<div dir="rtl">
		<ui.Calendar
			mode="single"
			month={RtlDefaultMonth}
			weekStartsOn={time.Sunday}
			showOutsideDays={true}
			captionLayout="label"
		/>
	</div>
}

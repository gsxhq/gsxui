package calendar

import (
	"time"

	"github.com/gsxhq/gsxui/ui"
)

// Bounded pins fromYear/toYear to the SAME single year (2026) as the
// displayed month — the narrowest possible navigation range. Task 4 review
// finding 1: calendar.js's own nav bounds were never recomputed after a
// client-side navigation, so with the wide ~110-year default range nothing
// in the browser suite ever reached a bound and the gap went unnoticed.
// This example exists purely to make that bound reachable in a handful of
// clicks (jstest/specs/calendar.spec.ts's own bounded-navigation tests).
// captionLayout="dropdown" additionally exercises the year <select> at the
// same narrow bound — Task 4 review's other observation, that
// syncDropdowns could otherwise write a year with no matching <option>.
component Bounded() {
	<ui.Calendar
		mode="single"
		month={time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
		weekStartsOn={time.Sunday}
		showOutsideDays={true}
		captionLayout="dropdown"
		fromYear={2026}
		toYear={2026}
	/>
}

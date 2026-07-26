package calendar

import (
	"time"

	"github.com/gsxhq/gsxui/ui"
)

// RangeDefaultMonth mirrors Basic's own DefaultMonth (2026-01, never
// time.Now()) — same reason: Task 5's browser tests assert exact grid
// contents and click-driven selection, which can't depend on the day the
// suite runs.
var RangeDefaultMonth = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// Range renders a range-mode calendar with no initial selection — Task 5's
// selection tests build the range entirely client-side: the two-click
// machine (first click sets from, second sets to, swapping if it precedes
// from, starting over once both are set) and the hover preview that fills in
// data-range-middle while from is set and to is not.
component Range(month time.Time) {
	{{ if month.IsZero() {
		month = RangeDefaultMonth
	} }}
	<ui.Calendar mode="range" month={month} weekStartsOn={time.Sunday} showOutsideDays={true} captionLayout="label"/>
}

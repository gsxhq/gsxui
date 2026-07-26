package calendar

import (
	"time"

	"github.com/gsxhq/gsxui/ui"
)

// MultipleDefaultMonth mirrors Basic's own DefaultMonth (2026-01, never
// time.Now()) — same reason as Range's own RangeDefaultMonth.
var MultipleDefaultMonth = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// Multiple renders a multiple-mode calendar with no initial selection —
// Task 5's browser tests toggle days on and off directly and check that
// each click only ever flips the ONE clicked day, never the rest of an
// existing selection.
component Multiple(month time.Time) {
	{{ if month.IsZero() {
		month = MultipleDefaultMonth
	} }}
	<ui.Calendar mode="multiple" month={month} weekStartsOn={time.Sunday} showOutsideDays={true} captionLayout="label"/>
}

package calendar

import (
	"time"

	"github.com/gsxhq/gsxui/ui"
)

// LoadedDefaultMonth mirrors Basic's own DefaultMonth (2026-01, never
// time.Now()) — same reason: a test asserting exact grid contents can't
// depend on the day it runs.
var LoadedDefaultMonth = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// Loaded exists solely to make calendar.js's ported dayDisabled, single-mode
// daySelected, ariaLabel and outside computations LIVE in the browser suite
// — Task 4 review's Important finding 2: basic.gsx has no selection, no
// disabled rules and captionLayout="label", so every one of those code
// paths in ui/calendar.js was zero-coverage; a bug in any of them (e.g. `<`
// vs `<=` against disabledBefore) left all seven original browser tests
// green. weekStartsOn=Monday is deliberate too (Task 4 review's Minor
// finding 4): basic.gsx's own Sunday week start makes the JS twin's
// `(x - weekStartsOn + 7) % 7` wrap-around term dead code, since `x - 0` is
// already non-negative.
//
// disabledBefore is 2026-01-05 (a Monday), deliberately NOT the same weekday
// as disabledWeekdays' own Saturday: an earlier draft used 2026-01-03 (also
// a Saturday), which meant a `<=` vs `<` bug at the before-bound was
// invisible — that date was already disabled by the independent weekday
// rule, so the OR-combination masked the regression. Every boundary date
// here (disabledBefore, disabledAfter, the disabledDates entry) is chosen to
// fall on a distinct, non-Saturday weekday, so each rule's own boundary is
// independently observable in a diff instead of being covered by another
// rule's coincidence.
component Loaded(month time.Time) {
	{{
		if month.IsZero() {
			month = LoadedDefaultMonth
		}
		selected := []time.Time{time.Date(2026, 1, 8, 0, 0, 0, 0, time.UTC)}
		disabledDates := []time.Time{time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)}
		disabledWeekdays := []time.Weekday{time.Saturday}
	}}
	<ui.Calendar
		mode="single"
		month={month}
		selected={selected}
		weekStartsOn={time.Monday}
		showOutsideDays={true}
		captionLayout="label"
		disabledBefore={time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC)}
		disabledAfter={time.Date(2026, 1, 28, 0, 0, 0, 0, time.UTC)}
		disabledDates={disabledDates}
		disabledWeekdays={disabledWeekdays}
	/>
}

// LoadedRangeDefaultMonth: same fixed-month discipline as LoadedDefaultMonth.
var LoadedRangeDefaultMonth = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// LoadedRange exists to exercise calendar.js's range-flag computation
// (rangeStart/rangeMiddle/rangeEnd and the cell's own cellSelected, which is
// true across the WHOLE range, not just the two ends) — Task 4 review's
// Important finding 2 again: basic.gsx and Loaded above are both
// single-mode, so mode === "range" was never reached client-side.
// weekStartsOn=Monday for the same reason as Loaded.
component LoadedRange(month time.Time) {
	{{ if month.IsZero() {
		month = LoadedRangeDefaultMonth
	} }}
	<ui.Calendar
		mode="range"
		month={month}
		from={time.Date(2026, 1, 9, 0, 0, 0, 0, time.UTC)}
		to={time.Date(2026, 1, 14, 0, 0, 0, 0, time.UTC)}
		weekStartsOn={time.Monday}
		showOutsideDays={true}
		captionLayout="label"
	/>
}

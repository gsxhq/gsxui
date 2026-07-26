package calendar

import (
	"time"

	"github.com/gsxhq/gsxui/ui"
)

// DropdownDefaultMonth mirrors Basic's own DefaultMonth (2026-01, never
// time.Now()) — same reason as every other fixture in this package: a
// component the browser invariants sweep asserts against can't depend on
// the day the suite runs.
var DropdownDefaultMonth = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// Dropdown renders captionLayout="dropdown" with a realistic, wide
// fromYear/toYear span — a birthdate-picker-shaped range, not Bounded's own
// deliberately narrow single-year one. Bounded.gsx exists to make the
// navigation-bound edge reachable in a handful of clicks; this example
// exists to show the month/year <select> pair the way a caller actually
// uses captionLayout="dropdown" day to day, with fromYear/toYear supplied
// explicitly rather than left to calendarNavBounds's own current−100/
// current+10 default (ui/calendar.gsx).
//
// No Query hook: like Bounded and Form, this example carries no ?month=
// browser test of its own — the invariants sweep exercises it as-is.
component Dropdown() {
	<ui.Calendar
		mode="single"
		month={DropdownDefaultMonth}
		weekStartsOn={time.Sunday}
		showOutsideDays={true}
		captionLayout="dropdown"
		fromYear={1980}
		toYear={2030}
	/>
}

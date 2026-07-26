package calendar

import (
	"time"

	"github.com/gsxhq/gsxui/ui"
)

// HiddenOutsideDefaultMonth mirrors Basic's own DefaultMonth (2026-01, never
// time.Now()) — same reason: a test asserting exact grid contents can't
// depend on the day it runs.
var HiddenOutsideDefaultMonth = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// HiddenOutside is the only example that passes showOutsideDays={false}, and
// it exists for exactly that reason. The final review found the parameter
// declared and read nowhere: every cell rendered regardless of its value, so
// passing false had no observable effect at all — and no example ever passed
// false, which is why it went unnoticed for seven tasks.
//
// The upstream semantics this exercises are NOT "drop the cell"
// (react-day-picker 9.14.0, helpers/createGetModifiers.js): showOutsideDays
// false sets modifiers.hidden, the <td> stays for layout with
// classNames.hidden = "invisible", and no DayButton renders inside it. gsxui
// keeps the button element (calendar.js may never create or destroy one) and
// blanks it instead — empty text, tabindex="-1", aria-hidden="true", native
// disabled — so a hidden day is out of the tab order entirely, unlike a
// merely DISABLED day, which stays in it.
//
// A ?month= Query hook (site/examples/calendar.go) makes this example
// reachable at an arbitrary month, so jstest/specs/calendar.spec.ts can run
// the same client-navigate-then-diff-against-the-server agreement test the
// other examples do — the one thing that keeps calendar.js's own `hidden`
// computation honest against calendar.gsx's.
component HiddenOutside(month time.Time) {
	{{ if month.IsZero() {
		month = HiddenOutsideDefaultMonth
	} }}
	<ui.Calendar mode="single" month={month} weekStartsOn={time.Sunday} showOutsideDays={false} captionLayout="label"/>
}

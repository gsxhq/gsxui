package calendar

import (
	"time"

	"github.com/gsxhq/gsxui/ui"
)

// MultipleDefaultMonth mirrors Basic's own DefaultMonth (2026-01, never
// time.Now()) — same reason as Range's own RangeDefaultMonth.
var MultipleDefaultMonth = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// Multiple renders a multiple-mode calendar with no initial selection.
// The form and name make the repeated-value bridge observable as ordinary
// FormData: every selected date is submitted under the same "dates" key.
component Multiple(month time.Time) {
	{{
		if month.IsZero() {
			month = MultipleDefaultMonth
		}
	}}
	<form>
		<ui.Calendar
			mode="multiple"
			month={month}
			name="dates"
			weekStartsOn={time.Sunday}
			showOutsideDays={true}
			captionLayout="label"
		/>
	</form>
}

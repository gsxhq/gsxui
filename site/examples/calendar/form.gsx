package calendar

import (
	"time"

	"github.com/gsxhq/gsxui/ui"
)

// FormDefaultMonth mirrors Basic's own DefaultMonth — pinned so Task 5's
// form-reset browser test gets a stable grid regardless of the day the
// suite runs.
var FormDefaultMonth = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// Form is real corpus coverage for the hidden-input form bridge Task 5
// review's Critical fix restored: name="date" with NOTHING preselected
// still renders an (empty-valued) hidden input from first paint — the
// single most likely way a caller wires ui.Calendar into a form — and
// ui/calendar.js's own reset handler (registered on
// "form:has([data-gsxui-calendar])") restores both the visible selection
// and that input when Reset is clicked. Before this example existed, that
// selector matched nothing anywhere in the static corpus and
// jstest/support/selector-allowlist.ts carried it as an allowedUnmatched
// entry (Task 5 review, Important 1) — this example is what removes it.
component Form() {
	<form class="flex max-w-xs flex-col gap-4">
		<ui.Calendar mode="single" month={FormDefaultMonth} name="date" weekStartsOn={time.Sunday} showOutsideDays={true} captionLayout="label"/>
		<div class="flex gap-2">
			<ui.Button type="submit">Continue</ui.Button>
			<ui.Button type="reset" variant="outline">Reset</ui.Button>
		</div>
	</form>
}

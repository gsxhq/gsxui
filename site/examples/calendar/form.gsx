package calendar

import (
	"time"

	"github.com/gsxhq/gsxui/ui"
)

// FormDefaultMonth mirrors Basic's own DefaultMonth — pinned so Task 5's
// form-reset browser test gets a stable grid regardless of the day the
// suite runs.
var FormDefaultMonth = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// Form is real corpus coverage for Calendar's hidden-input bridge and for a
// mixed form whose reset event is intentionally handled by two modules.
// Calendar and Combobox restore disjoint descendants of the same form; the
// selector invariant records that overlap explicitly and this composition
// proves both restorations survive one native reset.
component Form() {
	<form class="flex max-w-xs flex-col gap-4">
		<ui.Calendar
			mode="single"
			month={FormDefaultMonth}
			name="date"
			weekStartsOn={time.Sunday}
			showOutsideDays={true}
			captionLayout="label"
		/>
		<div class="flex flex-col gap-2">
			<ui.Label for="calendar-form-framework">Framework</ui.Label>
			<ui.Combobox name="framework" value="">
				<ui.ComboboxInput id="calendar-form-framework" placeholder="Search framework..." showTrigger/>
				<ui.ComboboxContent>
					<ui.ComboboxList>
						<ui.ComboboxEmpty>No framework found.</ui.ComboboxEmpty>
						<ui.ComboboxItem value="next.js" selected={false}>Next.js</ui.ComboboxItem>
						<ui.ComboboxItem value="sveltekit" selected={false}>SvelteKit</ui.ComboboxItem>
					</ui.ComboboxList>
				</ui.ComboboxContent>
			</ui.Combobox>
		</div>
		<div class="flex gap-2">
			<ui.Button type="submit">Continue</ui.Button>
			<ui.Button type="reset" variant="outline">Reset</ui.Button>
		</div>
	</form>
}

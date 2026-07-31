package calendar

import (
	"time"

	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// DatePickerDefaultMonth mirrors Basic's own DefaultMonth (2026-01, never
// time.Now()) — same reason as every other fixture in this package.
var DatePickerDefaultMonth = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// DatePicker is the composition shadcn documents as "Date Picker"
// (new-york-v4's own date-picker-demo.tsx): NOT its own component — a
// ui.Popover whose trigger is a ui.Button (showing "Pick a date" until a
// day is chosen, then the ISO date) and whose content is a ui.Calendar in
// single mode. This is the whole point of the example: gsxui has no
// DatePicker export, because Popover + Calendar already IS one, and
// composing them at the call site is strictly cheaper than threading a
// thirteenth Popover-specific parameter (open state, trigger label)
// through ui.Calendar itself. data-gsxui-popover-trigger goes straight on
// ui.Button, the same direct-attribute idiom popover/basic.gsx's own
// "Open popover" trigger already uses, rather than wrapping it in a
// separate ui.PopoverTrigger.
//
// The inline <script> is carousel/api.gsx's own established idiom: a
// document-level listener for a component's own CustomEvent, scoped by id
// so it never reacts to some OTHER example's calendar on the same
// /components/calendar page (every example on that page renders at once —
// site/pages/component.gsx — so a bare "did a calendar change" guard would
// also fire for Range's, Multiple's, and every other calendar example's own
// clicks). ui.Calendar emits gsxui:change with { mode: "single", selected }
// on every click (docs/superpowers/specs/2026-07-25-calendar-design.md
// §4); the listener writes selected[0] (or "Pick a date" once more, if the
// click cleared the selection — commitSingle in ui/calendar.js clears on a
// re-click of the already-selected day) onto the trigger's own label, and
// clears the muted "no date yet" styling to match.
component DatePicker() {
	<div>
		<ui.Popover>
			<ui.Button
				variant="outline"
				data-gsxui-popover-trigger
				data-gsxui-slot-popover-trigger
				aria-expanded="false"
				class="w-[240px] justify-start text-left font-normal text-muted-foreground"
			>
				<icon.Calendar/>
				<span data-datepicker-label>Pick a date</span>
			</ui.Button>
			<ui.PopoverContent class="w-auto p-0">
				<ui.Calendar
					id="datepicker-calendar"
					mode="single"
					month={DatePickerDefaultMonth}
					weekStartsOn={time.Sunday}
					showOutsideDays={true}
					captionLayout="label"
				/>
			</ui.PopoverContent>
		</ui.Popover>
		<script>
			document.addEventListener("gsxui:change", (e) => {
				if (e.target.id !== "datepicker-calendar") return;
				const button = e.target.closest("[data-gsxui-slot-popover]")?.querySelector("[data-gsxui-popover-trigger]");
				const label = button?.querySelector("[data-datepicker-label]");
				if (!label) return;
				const picked = e.detail.selected[0];
				label.textContent = picked ?? "Pick a date";
				button.classList.toggle("text-muted-foreground", !picked);
			});
		</script>
	</div>
}

package ui

import (
	"time"

	"github.com/gsxhq/gsx"
)

// calendarRootClass is the Calendar root slot ("w-fit", `## calendar` source
// map §2) plus the outer DayPicker className new-york-v4 sets alongside it
// ("group/calendar bg-background p-3 [--cell-size:--spacing(8)]
// [[data-slot=card-content]_&]:bg-transparent
// [[data-slot=popover-content]_&]:bg-transparent", source map §1) — gsxui
// collapses both onto the single root <div> since there is no separate
// DayPicker wrapper (`## calendar` element-structure decision). Retargeted to
// nova density (source map §6, `.cn-calendar`): p-3→p-2,
// --cell-size:--spacing(8)→--spacing(7). DROPPED, not carried: new-york-v4's
// two `rtl:**:[.rdp-button\_next>svg]:rotate-180`-shaped selectors — they
// target `.rdp-button_next`/`.rdp-button_previous`, and gsxui does not port
// react-day-picker's bare `rdp-*` hook classes at all (source map §2, "not
// resolved here" — resolved here as: don't port them), so the selectors would
// never match anything in this port. Ledgered for Task 7.
const calendarRootClass = "w-fit group/calendar bg-background p-2 [--cell-size:--spacing(7)] [[data-slot=card-content]_&]:bg-transparent [[data-slot=popover-content]_&]:bg-transparent"

// calendarGridClass is the `month_grid` slot, verbatim (source map §2).
const calendarGridClass = "w-full border-collapse"

// calendarWeekdaysClass is the `weekdays` slot, verbatim (source map §2) —
// applied to the `<tr>` inside `<thead aria-hidden>` (Weekdays.js renders
// both; classNames.weekdays lands on the inner `<tr>`, confirmed by reading
// Weekdays.js itself: `<thead aria-hidden><tr {...props}/></thead>`).
const calendarWeekdaysClass = "flex"

// calendarWeekdayClass is the `weekday` slot, verbatim (source map §2).
const calendarWeekdayClass = "flex-1 rounded-md text-[0.8rem] font-normal text-muted-foreground select-none"

// calendarWeekClass is the `week` slot, verbatim (source map §2).
const calendarWeekClass = "mt-2 flex w-full"

// calendarDayClass is the `<td>`'s always-on class: the `day` slot verbatim
// (source map §2, showWeekNumber's false branch — gsxui doesn't port
// week-number, source map §9), plus the `outside`/`today` slots (§2)
// rewritten as data-attribute Tailwind variants instead of conditional Go
// class-building, per this task's brief. `disabled`/`hidden` and the
// selection-driven `range_start`/`range_middle`/`range_end` slots are NOT
// folded in here — Task 1 computes no disabled/selection state at all, and
// pre-baking their data-variant classes now (before any task ever emits the
// matching data attribute) is speculative rather than settled; Tasks 2/3 add
// those tokens to this same constant when they wire up the attributes they
// gate on. Ledgered as an incremental (not full) class-string carryover for
// Task 7.
const calendarDayClass = "group/day relative aspect-square h-full w-full p-0 text-center select-none [&:last-child[data-selected=true]_button]:rounded-r-md [&:first-child[data-selected=true]_button]:rounded-l-md data-outside:text-muted-foreground data-outside:aria-selected:text-muted-foreground data-today:rounded-md data-today:bg-accent data-today:text-accent-foreground data-today:data-[selected=true]:rounded-none"

// calendarDayButtonClass is CalendarDayButton's own class string, verbatim
// (source map §3) — already expressed entirely as data-attribute variants
// upstream (`data-[range-start=true]:…`, `group-data-[focused=true]/day:…`),
// so no rewriting is needed here; Tasks 2/3 make it live by emitting the
// matching data attributes on this button. Composed onto button.gsx's own
// base/variantClass("ghost")/sizeClass("icon") — CalendarDayButton wraps
// gsxui's Button-equivalent with variant="ghost" size="icon" (source map §3)
// — the same package-private-helper reuse pagination.gsx and toggle-group.gsx
// already establish (internal/registry's declIndex derives a calendar ->
// button dependency from this, the same shape as pagination -> button).
const calendarDayButtonClass = "flex aspect-square size-auto w-full min-w-(--cell-size) flex-col gap-1 leading-none font-normal group-data-[focused=true]/day:relative group-data-[focused=true]/day:z-10 group-data-[focused=true]/day:border-ring group-data-[focused=true]/day:ring-[3px] group-data-[focused=true]/day:ring-ring/50 data-[range-end=true]:rounded-md data-[range-end=true]:rounded-r-md data-[range-end=true]:bg-primary data-[range-end=true]:text-primary-foreground data-[range-middle=true]:rounded-none data-[range-middle=true]:bg-accent data-[range-middle=true]:text-accent-foreground data-[range-start=true]:rounded-md data-[range-start=true]:rounded-l-md data-[range-start=true]:bg-primary data-[range-start=true]:text-primary-foreground data-[selected-single=true]:bg-primary data-[selected-single=true]:text-primary-foreground dark:hover:text-accent-foreground [&>span]:text-xs [&>span]:opacity-70"

// monthGrid returns the 42 dates of the six-week grid that displays the
// given month: seven columns starting at weekStartsOn, six rows, beginning
// at the week start on or before the 1st.
//
// Six rows always, never four or five. react-day-picker sizes the grid to
// the month and changes height between them; a fixed cell count is what lets
// calendar.js update days in place instead of rebuilding the grid, which in
// turn is what keeps every class string in this file and preserves focus
// across navigation. Ledgered as a deviation in docs/jsx-parity.md.
//
// UTC throughout. Local-midnight arithmetic crosses DST boundaries and
// produces off-by-one-day grids twice a year.
//
// This function has a twin in calendar.js. They must agree; the agreement is
// tested in jstest/specs/calendar.spec.ts.
func monthGrid(year int, month time.Month, weekStartsOn time.Weekday) [42]time.Time {
	first := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	offset := (int(first.Weekday()) - int(weekStartsOn) + 7) % 7
	start := first.AddDate(0, 0, -offset)

	var grid [42]time.Time
	for i := range grid {
		grid[i] = start.AddDate(0, 0, i)
	}
	return grid
}

// dayOutside reports whether d falls outside the displayed month (year,
// month) — every grid cell whose year or month doesn't match is a
// leading/trailing padding day borrowed from a neighboring month.
func dayOutside(d time.Time, year int, month time.Month) bool {
	return d.Year() != year || d.Month() != month
}

// firstFocusableIndex returns the grid index (0..41) of the first day
// belonging to the displayed month — the single day this task's simplified
// roving-tabindex model gives tabindex="0". This is deliberately narrower
// than react-day-picker's own priority order (explicit focus > lastFocused >
// selected > today > first focusable, source map §7.3); Tasks 2/3 widen it
// as selection and keyboard navigation land.
func firstFocusableIndex(grid [42]time.Time, year int, month time.Month) int {
	for i, d := range grid {
		if !dayOutside(d, year, month) {
			return i
		}
	}
	return 0
}

// Calendar is the shadcn/ui Calendar (registry/new-york-v4/ui/calendar.tsx),
// wrapping react-day-picker with no Radix primitive underneath — see
// docs/superpowers/plans/2026-07-25-calendar-source-map.md for the
// class-string and DOM/ARIA authority this port cites throughout.
//
// Task 1 renders the month grid only: the root, the aria-hidden weekday
// header row, and the six week rows (always six — source map §10, a settled
// deviation from upstream's 4-6 row sizing). mode/selected/from/to/
// captionLayout/fromYear/toYear/disabledBefore/disabledAfter/disabledDates/
// disabledWeekdays/name are all part of the sixteen-parameter signature every
// later task and example depends on (call-site positional order is load-
// bearing), but stay unused until Tasks 2/3 wire up selection, the caption,
// disabled-day computation, and the hidden form inputs — no placeholder
// markup for any of them here.
//
// The weekday header row is `aria-hidden`, matching upstream: every day
// button's own aria-label already leads with the weekday name (source map
// §8 correction 2), so hiding the abbreviated header loses no information
// and spares a screen-reader user "column 3, Tuesday" chatter on every move.
//
// Grid dates are UTC throughout (monthGrid's own doc comment); "today" is
// time.Now().UTC()'s calendar date, computed once per render so a render
// straddling local midnight is internally consistent.
component Calendar(
	mode string,
	month time.Time,
	selected []time.Time,
	from time.Time,
	to time.Time,
	weekStartsOn time.Weekday,
	showOutsideDays bool,
	captionLayout string,
	fromYear int,
	toYear int,
	disabledBefore time.Time,
	disabledAfter time.Time,
	disabledDates []time.Time,
	disabledWeekdays []time.Weekday,
	name string,
	attrs gsx.Attrs,
) {
	{{
		year := month.Year()
		monthOfYear := month.Month()
		grid := monthGrid(year, monthOfYear, weekStartsOn)
		today := time.Now().UTC()
		focusIdx := firstFocusableIndex(grid, year, monthOfYear)
		multiselectable := mode == "range" || mode == "multiple"
	}}
	<div
		data-slot="calendar"
		data-gsxui-calendar
		data-gsxui-calendar-month={ month.Format("2006-01") }
		data-gsxui-calendar-mode={ mode |> default("single") }
		data-gsxui-calendar-week-start={ int(weekStartsOn) }
		class={ calendarRootClass }
		{ attrs... }
	>
		<table
			data-gsxui-calendar-grid
			role="grid"
			{ if multiselectable {
				aria-multiselectable="true"
			} }
			class={ calendarGridClass }
		>
			<thead aria-hidden="true">
				<tr class={ calendarWeekdaysClass }>
					{ for i := 0; i < 7; i++ {
						{{ wd := time.Weekday((int(weekStartsOn) + i) % 7) }}
						<th scope="col" class={ calendarWeekdayClass }>{ wd.String()[:2] }</th>
					} }
				</tr>
			</thead>
			<tbody>
				{ for week := 0; week < 6; week++ {
					<tr class={ calendarWeekClass }>
						{ for day := 0; day < 7; day++ {
							{{
								idx := week*7 + day
								d := grid[idx]
								outside := dayOutside(d, year, monthOfYear)
								isToday := d.Year() == today.Year() && d.Month() == today.Month() && d.Day() == today.Day()
								tabindex := "-1"
								if idx == focusIdx {
									tabindex = "0"
								}
							}}
							<td
								role="gridcell"
								data-date={ d.Format("2006-01-02") }
								data-outside={ gsx.Toggle(outside) }
								data-today={ gsx.Toggle(isToday) }
								aria-selected="false"
								class={ calendarDayClass }
							>
								<button
									type="button"
									data-gsxui-calendar-day
									tabindex={ tabindex }
									aria-label={ d.Format("Monday, January 2, 2006") }
									class={ base, variantClass("ghost"), sizeClass("icon"), calendarDayButtonClass }
								>
									{ d.Day() }
								</button>
							</td>
						} }
					</tr>
				} }
			</tbody>
		</table>
	</div>
}

package ui

import (
	"strconv"
	"strings"
	"time"

	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// calendarRootClass is the Calendar `root` slot ("w-fit") plus the outer
// DayPicker top-level `className` new-york-v4 sets alongside it — both
// quoted byte-verbatim, as the two `root` rows, in source map §2 (the second
// row and its explanatory note were added in the Task 1 review-fix pass;
// §1 covers only the `--cell-size` token of that second string). gsxui
// collapses both onto the single root <div> since there is no separate
// DayPicker wrapper (`## calendar` element-structure decision). Retargeted to
// nova density (source map §6, `.cn-calendar`): p-3→p-2,
// --cell-size:--spacing(8)→--spacing(7). DROPPED, not carried: new-york-v4's
// two `rtl:**:[.rdp-button\_next>svg]:rotate-180`-shaped selectors (quoted in
// full in source map §2's note) — they target `.rdp-button_next`/
// `.rdp-button_previous`, and gsxui does not port react-day-picker's bare
// `rdp-*` hook classes at all, so the selectors would never match anything
// in this port. See source map §2's note for the full reasoning.
//
// `relative` (Task 3 review-fix pass, added here — not in either upstream
// `root` row above): the `nav` slot (§2) is `absolute inset-x-0 top-0 ...`,
// positioned against the nearest `relative` ancestor. Upstream's own
// positioned ancestor for this is the `months` slot's own `<div>`
// (§2: `relative flex flex-col gap-4 md:flex-row`) — but gsxui doesn't
// render `months`/`month` wrapper elements at all (§9: out of scope, single-
// month only), so this root `<div>` is the only element left standing in
// for that responsibility. Without this token the nav's two buttons would
// position against whatever ancestor happens to be `relative` further up
// the CONSUMING page (or the viewport itself) instead of sitting over the
// caption row, on every render — not a cosmetic gap. Ledgered as: root
// absorbs `months`'s own `relative` token, since `months` itself is not
// rendered.
const calendarRootClass = "relative w-fit group/calendar bg-background p-2 [--cell-size:--spacing(7)] [[data-slot=card-content]_&]:bg-transparent [[data-slot=popover-content]_&]:bg-transparent"

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
// week-number, source map §9), plus the `outside`/`today`/`disabled` slots
// (§2) rewritten as data-attribute Tailwind variants instead of conditional
// Go class-building, per this task's brief. This task adds the `disabled`
// slot (`text-muted-foreground opacity-50`, §2) as `data-disabled:…`, gated
// by the bare-presence `data-disabled` attribute this task starts emitting
// on the cell (the same presence-only pattern Task 1 used for
// `data-outside`/`data-today` — there is no `data-[disabled=true]` value-based
// selector anywhere in this string, so presence alone is sufficient).
//
// `range_start`/`range_middle`/`range_end` (§2) are deliberately NOT folded
// in here, even though Task 1's own comment flagged them as pending: the
// brief's DOM contract (task-1-brief.md) gives the cell no `data-range-*`
// attribute of its own — only the button carries those (source map §8
// finding 3) — so there is no cell-level attribute left to gate a
// cell-level range variant on. Porting §2's `range_start`/`range_middle`/
// `range_end` classNames values onto the cell would require adding a
// `data-range-*` attribute to the cell that contradicts the fixed DOM
// contract. Recorded as a deviation for Task 7's ledger: cell-level range
// coloring is dropped, the button's own range pill classes (§3, already
// live via `calendarDayButtonClass`) carry all the range-day styling
// instead. `hidden` (§2, `invisible`) is likewise not ported — no
// `data-hidden` attribute exists in the DOM contract either, and
// `showOutsideDays`/hiding is not in this task's scope.
//
// The pre-existing `[data-selected=true]` tokens below (both the
// `[&:last-child…]`/`[&:first-child…]` row-edge selectors and
// `data-today:data-[selected=true]:rounded-none`) were already present from
// Task 1, unconditionally dead until this task starts emitting the cell's
// own `data-selected="true"` (see the component body) — a value-based
// attribute, not a bare Toggle, precisely because these selectors require
// the literal value "true" to match (a bare `data-selected` attribute's
// value is the empty string for `[attr=value]` matching purposes, which
// would never satisfy `[data-selected=true]`).
const calendarDayClass = "group/day relative aspect-square h-full w-full p-0 text-center select-none [&:last-child[data-selected=true]_button]:rounded-r-md [&:first-child[data-selected=true]_button]:rounded-l-md data-outside:text-muted-foreground data-outside:aria-selected:text-muted-foreground data-today:rounded-md data-today:bg-accent data-today:text-accent-foreground data-today:data-[selected=true]:rounded-none data-disabled:text-muted-foreground data-disabled:opacity-50"

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

// calendarNavClass is the `nav` slot, verbatim (source map §2).
const calendarNavClass = "absolute inset-x-0 top-0 flex w-full items-center justify-between gap-1"

// calendarNavButtonClass is button_previous/button_next's own trailing
// classes, verbatim (source map §2: `size-(--cell-size) p-0 select-none
// aria-disabled:opacity-50`) — the `buttonVariants({variant: buttonVariant})`
// portion of that slot is supplied separately at the call site via this
// file's own base/variantClass("ghost") (upstream's buttonVariant default,
// `calendar.tsx` line 23), not duplicated into this constant. Deliberately
// NOT composed with button.gsx's sizeClass(...): upstream's own size-(--cell-
// size) here overrides buttonVariants' default size class entirely (cn()'s
// tailwind-merge resolves the conflict at render time); gsxui has no
// runtime class-merge, so the port resolves the same conflict at author
// time instead, by never emitting a sizeClass(...) token for this button in
// the first place — the button.gsx nav buttons are Button-shaped but not
// literal Button calls (see the Calendar component's own doc comment for why
// ui.Button itself is not used here).
const calendarNavButtonClass = "size-(--cell-size) p-0 select-none aria-disabled:opacity-50"

// calendarMonthCaptionClass is the `month_caption` slot, verbatim (source map §2).
const calendarMonthCaptionClass = "flex h-(--cell-size) w-full items-center justify-center px-(--cell-size)"

// calendarCaptionLabelClass is the `caption_label` slot's plain-"label"-
// layout branch only, verbatim (source map §2: "font-medium select-none" +
// the captionLayout==="label" ternary arm "text-sm"). The dropdown-layout
// ternary arm ("flex h-8 items-center gap-1 rounded-md pr-1 pl-2 text-sm
// [&>svg]:size-3.5 [&>svg]:text-muted-foreground") styled upstream's own
// semi-transparent-overlay trick (`Dropdown.js`: a visible, aria-hidden
// `<span>` showing the selected option's label + a chevron, layered over a
// second, fully invisible `<select>` — `classNames[UI.Dropdown]` is
// literally `absolute inset-0 bg-popover opacity-0`, source map §2's
// `dropdown` row). gsxui's captionLayout="dropdown" instead uses
// ui.NativeSelect as-is: a real, always-visible, styled `<select>` — there
// is no second invisible control and no faux label span in this port's DOM,
// so neither the `caption_label` dropdown arm nor the `dropdown` slot
// (`absolute inset-0 bg-popover opacity-0`) has anywhere to attach. Both are
// DROPPED, not carried — ledgered in the Task 3 report, not silently
// dropped. Porting them here would either do nothing (no matching element)
// or actively break the select's own visibility, which defeats the reason
// gsxui chose the simpler real-select design in the first place.
const calendarCaptionLabelClass = "font-medium select-none text-sm"

// calendarDropdownsClass is the `dropdowns` slot, verbatim (source map §2).
const calendarDropdownsClass = "flex h-(--cell-size) w-full items-center justify-center gap-1.5 text-sm font-medium"

// calendarDropdownRootClass is the `dropdown_root` slot, retargeted to nova
// density (source map §6, `.cn-calendar-dropdown-root`) the same way Task 1's
// calendarRootClass retargeted the `root`/top-level className (p-3→p-2,
// --spacing(8)→(7)): new-york-v4's own value (source map §2) is `relative
// rounded-md border border-input shadow-xs has-focus:border-ring
// has-focus:ring-[3px] has-focus:ring-ring/50`; nova's rule drops
// `shadow-xs` entirely (shadow-presence removal, not a value swap) and
// respells `ring-[3px]`→`ring-3` (same value, nova's own spelling — the same
// substitution the menus/combobox source maps already made for this exact
// token). `rounded-md` is unchanged: nova's rule has no `rounded-*` token at
// all, so there is nothing to retarget there (§1/§6 already reject porting
// `bases/base`'s own `--cell-radius` substitution). Every token nova speaks
// to has a nova counterpart here — none silently left at the new-york-v4
// value. Passed as ui.NativeSelect's own `class` attr at each call site,
// merging onto NativeSelect's wrapper `<div>`'s hardcoded "relative w-fit"
// (ui/native-select.gsx) — the wrapper is the analogue of upstream's
// `Dropdown.js` `<span data-disabled className={DropdownRoot}>`.
const calendarDropdownRootClass = "relative rounded-md border border-input has-focus:border-ring has-focus:ring-3 has-focus:ring-ring/50"

// calendarMonthNames are the twelve month names for the dropdown
// captionLayout's month <select>, matching upstream's own default
// formatMonthDropdown formatter (`date.toLocaleString("default", {month:
// "short"})`, `calendar.tsx` lines 43-45): three-letter abbreviations — a
// different formatting choice from the day button's own full-month
// aria-label ("Monday, January 2, 2006").
var calendarMonthNames = [12]string{
	"Jan", "Feb", "Mar", "Apr", "May", "Jun",
	"Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
}

// calendarNavBounds resolves fromYear/toYear into the effective navigation
// bounds, defaulting when the caller passes zero. react-day-picker's own
// getYearOptions (source map §7.4) requires an explicit navStart/navEnd for
// captionLayout="dropdown" and renders no year options at all without one;
// gsxui defaults instead (source map §10 item 2) so callers aren't forced to
// pick a range just to get a dropdown calendar. currentYear anchors the
// default to "now", not a fixed year, so the default keeps working as time
// passes. Source map §7.5 establishes this is the ONE bounds contract behind
// both the dropdown year list and the prev/next nav buttons' aria-disabled
// state — not two separate mechanisms.
func calendarNavBounds(fromYear, toYear, currentYear int) (int, int) {
	if fromYear == 0 {
		fromYear = currentYear - 100
	}
	if toYear == 0 {
		toYear = currentYear + 10
	}
	return fromYear, toYear
}

// calendarPrevDisabled/calendarNextDisabled report whether the previous/next
// nav button is at its navigation bound for the displayed (year, month) —
// source map §7.5's getPreviousMonth/getNextMonth contract: previousMonth/
// nextMonth become undefined once navigation would cross navStart/navEnd,
// which is what drives Nav.js's aria-disabled/tabIndex-only branch (never a
// native disabled attribute, at either button, under any configuration).
func calendarPrevDisabled(year int, month time.Month, fromYear int) bool {
	return year < fromYear || (year == fromYear && month <= time.January)
}

func calendarNextDisabled(year int, month time.Month, toYear int) bool {
	return year > toYear || (year == toYear && month >= time.December)
}

// sameDay reports whether a and b fall on the same UTC calendar day.
func sameDay(a, b time.Time) bool {
	ay, am, ad := a.UTC().Date()
	by, bm, bd := b.UTC().Date()
	return ay == by && am == bm && ad == bd
}

// dayOnly normalizes t to a UTC midnight time.Time, so calendar days can be
// ordered with Before/After without a non-midnight or non-UTC caller's
// time-of-day leaking into the comparison — the same UTC-calendar-day
// discipline sameDay applies to equality, applied to ordering instead.
func dayOnly(t time.Time) time.Time {
	y, m, d := t.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// boolStr renders b as the literal string "true"/"false". Required wherever
// a Tailwind selector matches on the attribute's exact value rather than
// its mere presence — CalendarDayButton's data-selected-single/data-range-*
// attributes (source map §3, `data-[range-start=true]:…`) always carry one
// of these two literal strings, never bare presence and never omitted.
func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// daySelected reports whether d is one of the caller's selected dates.
// Meaningful in single/multiple mode only (source map §7.1): range mode's
// selection is carried by from/to instead, via rangeFlags. daySelected is
// combined with rangeFlags' three booleans at the call site to get the
// cell's real selected state across all three modes (source map §4:
// modifiers.selected is true for every day from/to inclusive in range mode
// too, layered under range_start/range_middle/range_end rather than
// replaced by them) — daySelected alone is deliberately narrower than that,
// it only ever answers the single/multiple half.
func daySelected(mode string, d time.Time, selected []time.Time) bool {
	if mode != "single" && mode != "multiple" {
		return false
	}
	for _, s := range selected {
		if sameDay(d, s) {
			return true
		}
	}
	return false
}

// dayDisabled reports whether d is disabled under any of the four rules
// this task ports as typed Calendar parameters: strictly before
// disabledBefore, strictly after disabledAfter, an exact match in
// disabledDates, or a weekday in disabledWeekdays (source map §7.2's full
// upstream Matcher taxonomy covers more shapes — DateInterval, arbitrary
// predicate, etc. — none of which this Calendar signature exposes). A zero
// time.Time for disabledBefore/disabledAfter means that bound is unset.
func dayDisabled(d, disabledBefore, disabledAfter time.Time, disabledDates []time.Time, disabledWeekdays []time.Weekday) bool {
	if !disabledBefore.IsZero() && dayOnly(d).Before(dayOnly(disabledBefore)) {
		return true
	}
	if !disabledAfter.IsZero() && dayOnly(d).After(dayOnly(disabledAfter)) {
		return true
	}
	for _, dd := range disabledDates {
		if sameDay(d, dd) {
			return true
		}
	}
	for _, wd := range disabledWeekdays {
		if d.Weekday() == wd {
			return true
		}
	}
	return false
}

// rangeFlags computes the range-start/middle/end modifiers for d — range
// mode only (source map §7.1's useRange.js selection model; source map §3
// for how the three flags drive CalendarDayButton's own pill classes). A
// zero from or to means that bound isn't set yet (an in-progress range);
// rangeMiddle requires both bounds set and d strictly between them.
func rangeFlags(mode string, d, from, to time.Time) (start, middle, end bool) {
	if mode != "range" {
		return false, false, false
	}
	if !from.IsZero() && sameDay(d, from) {
		start = true
	}
	if !to.IsZero() && sameDay(d, to) {
		end = true
	}
	if !from.IsZero() && !to.IsZero() && dayOnly(d).After(dayOnly(from)) && dayOnly(d).Before(dayOnly(to)) {
		middle = true
	}
	return
}

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
	// Unreachable in practice: monthGrid always straddles the 1st of (year,
	// month) with six full weeks either side, so the grid always contains
	// every day of that month — at least 28 of the 42 cells satisfy
	// !dayOutside. Kept only so the function has a defensive, well-typed
	// return for the compiler; it is not a live code path.
	return 0
}

// Calendar is the shadcn/ui Calendar (registry/new-york-v4/ui/calendar.tsx),
// wrapping react-day-picker with no Radix primitive underneath — see
// docs/superpowers/plans/2026-07-25-calendar-source-map.md for the
// class-string and DOM/ARIA authority this port cites throughout.
//
// Task 1 rendered the month grid only. Task 2 added all server-computed day
// state: selection (`single`/`multiple` via `selected`, `range` via
// `from`/`to`) and the four disabled rules
// (`disabledBefore`/`disabledAfter`/`disabledDates`/`disabledWeekdays`),
// entirely as data attributes — never as conditional Go class-building, so
// Tasks 4-6's JavaScript can update the same attributes without knowing
// about classes.
//
// Task 3 (this pass) adds the caption and the prev/next nav buttons above
// the grid, plus the hidden-input form bridge:
//
//   - captionLayout="label": the month/year as text ("January 2006"),
//     flanked by two icon-only nav buttons. captionLayout="dropdown": two
//     ui.NativeSelect month/year pickers instead of the text, still flanked
//     by the same two nav buttons. Any other captionLayout value (including
//     "") renders the label layout, matching upstream's own
//     captionLayout="label" default.
//   - The nav buttons are NOT ui.Button calls: ui.Button's `disabled bool`
//     param always renders a native `disabled` attribute, but source map
//     §7.5 traces that upstream's own nav buttons NEVER take native
//     `disabled`, at any navigation bound — only `aria-disabled="true"` +
//     `tabindex="-1"`, which is what keeps a disabled-at-the-bound button
//     reachable by Tab. So the buttons here are hand-built `<button>`
//     elements composed from button.gsx's own base/variantClass("ghost")
//     package-private helpers (the same flat-package reuse pagination.gsx
//     and toggle-group.gsx already establish) plus calendarNavButtonClass —
//     never through ui.Button, and never emitting `disabled`.
//   - fromYear/toYear default when zero (calendarNavBounds:
//     currentYear-100/currentYear+10) — source map §7.4/§10 item 2's already
//     -decided deviation from upstream's own requirement that the caller
//     supply an explicit range for captionLayout="dropdown". The resolved
//     bounds also drive the nav buttons' aria-disabled state in EVERY
//     captionLayout, not just "dropdown" — source map §7.5 is explicit that
//     this is one bounds contract, not two.
//   - The caption is a live region (source map §8 finding 4): `role="status"
//     aria-live="polite"`, so a month change gets announced. In "label"
//     layout that's the visible caption `<span>` itself. In "dropdown"
//     layout the visible control is a `<select>`, not a natural announcement
//     target, so — matching upstream's own second, separate mechanism for
//     this exact case (`DayPicker.js`'s dropdown branch renders an
//     additional visually-hidden `<span role="status" aria-live="polite">`
//     carrying the same formatted caption text, alongside the dropdowns, not
//     instead of them) — gsxui renders a `class="sr-only"` span carrying the
//     identical text. Both the label span and the dropdown sr-only span also
//     carry `data-gsxui-calendar-caption`, the one hook Task 4's calendar.js
//     needs regardless of layout to update the announced text on navigation.
//   - Hidden inputs render whenever `name != ""` — unconditionally, NOT
//     gated on there being a selection yet (Task 5 review, Critical: an
//     earlier revision also required `len(selected) > 0`/`!from.IsZero()`,
//     which meant `<ui.Calendar mode="single" name="date"/>` with nothing
//     preselected emitted no input at all; the user would click a day,
//     calendar.js's own syncHiddenInputs would find nothing to update, and
//     the form would submit with no `date` field whatsoever — a real hole
//     in the single most likely way a caller wires this component into a
//     form). Single/multiple mode gets one `<input type="hidden"
//     name={name}>`, value empty until something is selected. Range mode
//     gets that plus a second `name="{name}-to"` input for `to`, marked
//     with a bare `data-gsxui-calendar-hidden-to` so calendar.js can tell
//     the two apart by something other than the `name` string itself (a
//     caller passing `name="valid-to"` in range mode would otherwise make
//     an `endsWith("-to")` check misidentify its OWN `from` input as the
//     `to` input — Task 5 review, Minor 1). Both live inside the root
//     element so Tasks 4-6's JavaScript can find them by DOM scope
//     (`closest("[data-gsxui-calendar]")`), the same proximity idiom
//     ui.Select/ui.Combobox already use for their own hidden-input bridges.
//
// The weekday header row is `aria-hidden`, matching upstream: every day
// button's own aria-label already leads with the weekday name (source map
// §8 correction 2), so hiding the abbreviated header loses no information
// and spares a screen-reader user "column 3, Tuesday" chatter on every move.
//
// Grid dates are UTC throughout (monthGrid's own doc comment); "today" is
// time.Now().UTC()'s calendar date, computed once per render so a render
// straddling local midnight is internally consistent. Selection/disabled
// comparisons are UTC-calendar-day comparisons too (sameDay/dayOnly), never
// time.Time equality — a caller may pass a non-midnight or non-UTC time.
//
// The root carries the selection as data attributes so calendar.js (Tasks
// 4-6) can read the server's initial state without re-deriving it:
// data-gsxui-calendar-selected (comma-separated ISO dates, single/multiple),
// data-gsxui-calendar-from/-to (range mode). Each is omitted entirely when
// empty, never emitted as ="".
//
// Task 4 adds the four disabled rules to the root the same way, for the same
// reason: data-gsxui-calendar-disabled-before/-after (ISO), -disabled-dates
// (comma-separated ISO), -disabled-weekdays (comma-separated
// time.Weekday ints, Sunday=0). Task 3 baked dayDisabled's result into each
// of the initial month's 42 cells and stopped there; that snapshot goes
// stale the moment calendar.js navigates to a month the server never
// rendered, so the RULES themselves — not just their initial-month
// application — have to reach the client too. Same omit-when-empty
// convention as -selected/-from/-to.
//
// data-gsxui-calendar-nav-from-year/-nav-to-year carry calendarNavBounds'
// own resolved (navFromYear, navToYear) — unconditionally, never omitted,
// since a bound always has a value (calendarNavBounds defaults it when the
// caller passes zero). Task 4 review found that without this, calendar.js
// had no way to recompute calendarPrevDisabled/calendarNextDisabled after a
// client-side navigation: prevDisabled/nextDisabled above are only ever
// true for the SERVER's initial (year, month), so a caller who configures a
// narrow fromYear/toYear range (e.g. both 2026) could click past the
// declared bound into a month the server can never render for this
// component, and — in captionLayout="dropdown" — past the year <select>'s
// own last <option>, leaving it holding a value with no matching option.
// calendar.js ports calendarPrevDisabled/calendarNextDisabled themselves and
// re-applies the aria-disabled/tabindex pair to the nav buttons on every
// navigation, exactly mirroring the { if prevDisabled/nextDisabled { ... } }
// omit-otherwise shape below.
//
// Two distinct data-attribute sets, on two elements, per source map §8
// finding 3 and the brief's DOM contract — not collapsed onto one element:
// the cell carries data-disabled/data-selected/aria-selected (alongside
// Task 1's data-date/data-outside/data-today), in that fixed order; the
// button carries shadcn's own data-selected-single/data-range-start/
// data-range-middle/data-range-end, always present as the literal strings
// "true"/"false" (source map §3 — the class selectors match on the literal
// value, e.g. data-[range-start=true]:…). The cell's data-selected is
// likewise a literal "true" (omitted when false) rather than a bare Toggle:
// calendarDayClass already carries [data-selected=true]-based selectors
// from Task 1 (the row-edge rounding and the today+selected interaction)
// that only fire on an exact value match, not mere attribute presence.
//
// Task 4 adds one exception to "the cell owns data-date, the button doesn't":
// the button ALSO carries its own data-date (identical value, additive —
// the cell's copy is untouched). NOT because [data-gsxui-calendar-day] is
// the only one-per-cell selector — td[role="gridcell"] already was, just as
// uniquely (the header row's cells are <th>, never <td>) — but because
// Task 5's click handler (day selection) matches the BUTTON, the thing the
// user actually clicks, and needs that day's date to act on. Reading it
// off the button directly is one property access; reading it off the
// enclosing cell is a walk (`closest("td")`) for information the clicked
// element could have carried itself. calendar.js's own in-place update and
// jstest/specs/calendar.spec.ts's gridDates/gridCells helpers follow the
// same button-first convention for consistency, not because any of them
// individually required it.
//
// This makes data-date one fact recorded in two places, not two facts —
// keep them equal on every write. calendar.js's repaint already writes both
// (cell.dataset.date and button.dataset.date) from the same loop-local
// dateISO; a future change that updates only one of the two copies (on
// either the Go or the JS side) desyncs them silently, since nothing else
// here cross-checks the pair.
//
// The cell's data-selected/aria-selected fire in EVERY mode, including
// range: react-day-picker's own modifiers.selected (source map §4) is true
// for every day from from through to inclusive, both ends included
// (confirmed directly against useRange.js's own isSelected, which calls
// rangeIncludesDate(selected, date, false, dateLib) — the third argument is
// excludeEnds, and it is false). range_start/range_middle/range_end layer
// on top of selected; they do not replace it. Getting this wrong silences
// aria-selected for an entire range (a real accessibility regression) and
// leaves calendarDayClass's three [data-selected=true]-gated row-wrap
// rounding selectors permanently dead for any range that wraps a <tr>
// boundary. The button's own data-selected-single stays mode-sensitive
// exactly as before — it means "selected and not any range flag," so it
// must stay "false" for every day of a range even though the cell right
// above it is "selected".
//
// data-focused is not emitted here at all: nothing holds live focus on the
// server, so first paint never has a focused day. Task 6 owns setting it
// client-side. Likewise, source map §8 finding 5's degraded-to-aria-disabled
// case (a disabled day that already holds focus keeps the native disabled
// attribute off) never applies on first paint — every disabled day gets the
// native `disabled` attribute here, unconditionally.
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

		var selectedISOParts []string
		for _, s := range selected {
			selectedISOParts = append(selectedISOParts, s.Format("2006-01-02"))
		}
		selectedISO := strings.Join(selectedISOParts, ",")

		navFromYear, navToYear := calendarNavBounds(fromYear, toYear, today.Year())
		prevDisabled := calendarPrevDisabled(year, monthOfYear, navFromYear)
		nextDisabled := calendarNextDisabled(year, monthOfYear, navToYear)
		dropdownLayout := captionLayout == "dropdown"
		captionText := month.Format("January 2006")

		// Unconditional on there being a selection yet (Task 5 review,
		// Critical) — only on `name` and `mode`. hiddenSingleValue/
		// hiddenFromValue/hiddenToValue default to "" and stay "" until
		// there's something real to put in them; the input itself still
		// renders either way, so a form submits an (empty) field rather
		// than silently dropping it.
		showHiddenSingle := name != "" && mode != "range"
		showHiddenFrom := name != "" && mode == "range"
		showHiddenTo := name != "" && mode == "range"

		hiddenSingleValue := ""
		if len(selected) > 0 {
			hiddenSingleValue = selected[0].Format("2006-01-02")
		}
		hiddenFromValue := ""
		if !from.IsZero() {
			hiddenFromValue = from.Format("2006-01-02")
		}
		hiddenToValue := ""
		if !to.IsZero() {
			hiddenToValue = to.Format("2006-01-02")
		}

		// The four disabled rules, serialized onto the root so calendar.js
		// (Task 4) can re-derive dayDisabled for a client-navigated month
		// without the server rendering it — Task 3 only ever rendered these
		// rules baked into the initial month's 42 cells, which goes stale
		// the instant JS navigates away from that month. Comma-separated,
		// same ISO/int encoding as the per-day comparisons above; omitted
		// entirely when the rule is unset, never emitted empty (matching
		// -selected/-from/-to's own omit-when-empty convention just above).
		var disabledDatesISO []string
		for _, dd := range disabledDates {
			disabledDatesISO = append(disabledDatesISO, dd.Format("2006-01-02"))
		}
		disabledDatesAttr := strings.Join(disabledDatesISO, ",")

		var disabledWeekdaysStr []string
		for _, wd := range disabledWeekdays {
			disabledWeekdaysStr = append(disabledWeekdaysStr, strconv.Itoa(int(wd)))
		}
		disabledWeekdaysAttr := strings.Join(disabledWeekdaysStr, ",")
	}}
	<div
		data-slot="calendar"
		data-gsxui-calendar
		data-gsxui-calendar-month={ month.Format("2006-01") }
		data-gsxui-calendar-mode={ mode |> default("single") }
		data-gsxui-calendar-week-start={ int(weekStartsOn) }
		{ if selectedISO != "" {
			data-gsxui-calendar-selected={ selectedISO }
		} }
		{ if !from.IsZero() {
			data-gsxui-calendar-from={ from.Format("2006-01-02") }
		} }
		{ if !to.IsZero() {
			data-gsxui-calendar-to={ to.Format("2006-01-02") }
		} }
		{ if !disabledBefore.IsZero() {
			data-gsxui-calendar-disabled-before={ disabledBefore.Format("2006-01-02") }
		} }
		{ if !disabledAfter.IsZero() {
			data-gsxui-calendar-disabled-after={ disabledAfter.Format("2006-01-02") }
		} }
		{ if disabledDatesAttr != "" {
			data-gsxui-calendar-disabled-dates={ disabledDatesAttr }
		} }
		{ if disabledWeekdaysAttr != "" {
			data-gsxui-calendar-disabled-weekdays={ disabledWeekdaysAttr }
		} }
		data-gsxui-calendar-nav-from-year={ strconv.Itoa(navFromYear) }
		data-gsxui-calendar-nav-to-year={ strconv.Itoa(navToYear) }
		class={ calendarRootClass }
		{ attrs... }
	>
		<nav class={ calendarNavClass }>
			<button
				type="button"
				data-gsxui-calendar-prev
				aria-label="Previous month"
				{ if prevDisabled {
					aria-disabled="true"
					tabindex="-1"
				} }
				class={ base, variantClass("ghost"), calendarNavButtonClass }
			>
				<icon.ChevronLeft/>
			</button>
			<button
				type="button"
				data-gsxui-calendar-next
				aria-label="Next month"
				{ if nextDisabled {
					aria-disabled="true"
					tabindex="-1"
				} }
				class={ base, variantClass("ghost"), calendarNavButtonClass }
			>
				<icon.ChevronRight/>
			</button>
		</nav>
		<div class={ calendarMonthCaptionClass }>
			{ if dropdownLayout {
				<div class={ calendarDropdownsClass }>
					<NativeSelect data-gsxui-calendar-month-select aria-label="Month" class={ calendarDropdownRootClass }>
						{ for i := 0; i < 12; i++ {
							<NativeSelectOption value={ strconv.Itoa(i) } selected={ i == int(monthOfYear)-1 } data-gsxui-calendar-month-option>
								{ calendarMonthNames[i] }
							</NativeSelectOption>
						} }
					</NativeSelect>
					<NativeSelect data-gsxui-calendar-year-select aria-label="Year" class={ calendarDropdownRootClass }>
						{ for y := navFromYear; y <= navToYear; y++ {
							<NativeSelectOption value={ strconv.Itoa(y) } selected={ y == year } data-gsxui-calendar-year-option>
								{ strconv.Itoa(y) }
							</NativeSelectOption>
						} }
					</NativeSelect>
				</div>
				<span class="sr-only" data-gsxui-calendar-caption role="status" aria-live="polite">{ captionText }</span>
			} else {
				<span data-gsxui-calendar-caption role="status" aria-live="polite" class={ calendarCaptionLabelClass }>{ captionText }</span>
			} }
		</div>
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
								dayDis := dayDisabled(d, disabledBefore, disabledAfter, disabledDates, disabledWeekdays)
								daySel := daySelected(mode, d, selected)
								rStart, rMiddle, rEnd := rangeFlags(mode, d, from, to)
								selSingle := daySel && !rStart && !rMiddle && !rEnd
								cellSel := daySel || rStart || rMiddle || rEnd
							}}
							<td
								role="gridcell"
								data-date={ d.Format("2006-01-02") }
								data-outside={ gsx.Toggle(outside) }
								data-today={ gsx.Toggle(isToday) }
								data-disabled={ gsx.Toggle(dayDis) }
								{ if cellSel {
									data-selected="true"
								} }
								aria-selected={ boolStr(cellSel) }
								class={ calendarDayClass }
							>
								<button
									type="button"
									data-gsxui-calendar-day
									data-date={ d.Format("2006-01-02") }
									tabindex={ tabindex }
									aria-label={ d.Format("Monday, January 2, 2006") }
									data-selected-single={ boolStr(selSingle) }
									data-range-start={ boolStr(rStart) }
									data-range-middle={ boolStr(rMiddle) }
									data-range-end={ boolStr(rEnd) }
									disabled={ dayDis }
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
		{ if showHiddenSingle {
			<input type="hidden" name={ name } value={ hiddenSingleValue }/>
		} }
		{ if showHiddenFrom {
			<input type="hidden" name={ name } value={ hiddenFromValue }/>
		} }
		{ if showHiddenTo {
			<input type="hidden" name={ name + "-to" } value={ hiddenToValue } data-gsxui-calendar-hidden-to/>
		} }
	</div>
}

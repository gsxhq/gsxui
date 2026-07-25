package ui_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gsxhq/gsxui/ui"
)

// grid returns the data-date values of every day button, in DOM order.
func gridDates(t *testing.T, html string) []string {
	t.Helper()
	var out []string
	for _, part := range strings.Split(html, `data-date="`)[1:] {
		end := strings.Index(part, `"`)
		if end < 0 {
			t.Fatalf("unterminated data-date in %s", html)
		}
		out = append(out, part[:end])
	}
	return out
}

func TestCalendarGridIs42CellsStartingOnTheWeekStart(t *testing.T) {
	// 2026-01-01 is a Thursday; with weekStartsOn=Sunday the grid opens on
	// 2025-12-28 and runs 42 days to 2026-02-07.
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	dates := gridDates(t, got)
	if len(dates) != 42 {
		t.Fatalf("got %d day cells, want 42", len(dates))
	}
	if dates[0] != "2025-12-28" {
		t.Errorf("first cell = %q, want 2025-12-28", dates[0])
	}
	if dates[41] != "2026-02-07" {
		t.Errorf("last cell = %q, want 2026-02-07", dates[41])
	}
}

func TestCalendarWeekStartMondayShiftsTheGrid(t *testing.T) {
	// 2026-02-01 is a Sunday. weekStartsOn=Monday pushes the grid back six
	// days to 2026-01-26.
	got := render(t, ui.Calendar("single", time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Monday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	if first := gridDates(t, got)[0]; first != "2026-01-26" {
		t.Errorf("first cell = %q, want 2026-01-26", first)
	}
}

func TestCalendarLeapFebruary(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	var found bool
	for _, d := range gridDates(t, got) {
		if d == "2024-02-29" {
			found = true
		}
	}
	if !found {
		t.Error("2024-02-29 missing from the February 2024 grid")
	}
}

func TestCalendarMarksOutsideDays(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	// 2025-12-28 precedes January, 2026-01-01 does not.
	if !strings.Contains(got, `data-date="2025-12-28" data-outside`) {
		t.Error("2025-12-28 is not marked data-outside")
	}
	if strings.Contains(got, `data-date="2026-01-15" data-outside`) {
		t.Error("2026-01-15 is wrongly marked data-outside")
	}
}

func TestCalendarRovingTabindexHasExactlyOneStop(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	if n := strings.Count(got, `tabindex="0"`); n != 1 {
		t.Errorf("got %d tabindex=\"0\" day buttons, want exactly 1", n)
	}
}

// Upstream renders a real <table>, so row and column-header roles are
// implicit. Writing them explicitly is redundant ARIA, which the authoring
// practices warn against. Source map §8 corrections 1 and 2.
func TestCalendarGridAria(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	for want, count := range map[string]int{
		`role="grid"`:                1,
		`<thead aria-hidden="true">`: 1,
		`<th scope="col"`:            7,
		`<tr`:                        7, // one weekday header row + six week rows
		`role="gridcell"`:            42,
	} {
		if n := strings.Count(got, want); n != count {
			t.Errorf("got %d %s, want %d", n, want, count)
		}
	}
	for _, unwanted := range []string{`role="row"`, `role="columnheader"`} {
		if strings.Contains(got, unwanted) {
			t.Errorf("%s is implicit on a real table — do not write it", unwanted)
		}
	}
}

// aria-multiselectable is set only where more than one date can be chosen.
func TestCalendarAriaMultiselectable(t *testing.T) {
	month := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		mode string
		want bool
	}{
		{"single", false},
		{"multiple", true},
		{"range", true},
	} {
		got := render(t, ui.Calendar(tc.mode, month, nil, time.Time{}, time.Time{},
			time.Sunday, true, "label", 0, 0, time.Time{}, time.Time{}, nil, nil, "", nil))
		has := strings.Contains(got, `aria-multiselectable="true"`)
		if has != tc.want {
			t.Errorf("mode %q: aria-multiselectable present = %v, want %v", tc.mode, has, tc.want)
		}
	}
}

// Every day's accessible name begins with the weekday. This is what makes
// hiding the weekday header row safe — source map §8 correction 2.
func TestCalendarDayLabelCarriesTheWeekday(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	// 2026-01-15 is a Thursday.
	if !strings.Contains(got, `aria-label="Thursday, January 15, 2026"`) {
		t.Errorf("day label missing or wrong format\nin: %s", got)
	}
}

func TestCalendarRootAttributes(t *testing.T) {
	got := render(t, ui.Calendar("range", time.Date(2026, 7, 4, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Monday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	for _, want := range []string{
		`data-slot="calendar"`,
		`data-gsxui-calendar`,
		`data-gsxui-calendar-month="2026-07"`,
		`data-gsxui-calendar-mode="range"`,
		`data-gsxui-calendar-week-start="1"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("root missing %q", want)
		}
	}
}

// todayCellMarker is the exact substring a day cell renders when it carries
// data-today: the bare attribute immediately followed by aria-selected, per
// the cell's fixed attribute order (data-date, data-outside?, data-today?,
// aria-selected). data-today also appears inside every cell's own (always-
// present) data-today:* Tailwind-variant class tokens, so a bare
// strings.Count(got, "data-today") would over-count; anchoring on the
// aria-selected neighbor isolates the real attribute.
func todayCellMarker(dateISO string) string {
	return fmt.Sprintf(`data-date="%s" data-today aria-selected="false"`, dateISO)
}

// data-today is driven by time.Now().UTC(), not a fixed date (source map's
// "today" modifier; Task 6 builds on this attribute) — asserting against a
// literal date would freeze the moment this test was written and go stale
// immediately after. Rendering the CURRENT UTC month and checking against
// time.Now().UTC() keeps this test correct on every day it runs.
func TestCalendarTodayIsMarkedExactlyOnce(t *testing.T) {
	now := time.Now().UTC()
	got := render(t, ui.Calendar("single", now,
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	if n := strings.Count(got, "data-today aria-selected"); n != 1 {
		t.Fatalf("got %d data-today cells in the current month's grid, want exactly 1", n)
	}
	want := todayCellMarker(now.Format("2006-01-02"))
	if !strings.Contains(got, want) {
		t.Errorf("today's own cell (%s) is not the one marked data-today\nin: %s", now.Format("2006-01-02"), got)
	}
}

// A grid nowhere near today's date carries no data-today cell at all — the
// half of the contract a same-month-only comparison bug (e.g. matching on
// day-of-month while ignoring year/month) would not catch on its own.
func TestCalendarTodayAbsentFarFromToday(t *testing.T) {
	farMonth := time.Now().UTC().AddDate(5, 0, 0)
	got := render(t, ui.Calendar("single", farMonth,
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	if strings.Contains(got, "data-today aria-selected") {
		t.Errorf("grid for %s wrongly contains a data-today cell\nin: %s", farMonth.Format("2006-01"), got)
	}
}

// cellFor returns the <td>…</td> substring for the given ISO date, so a test
// can assert on both the cell's own state attributes and the day button
// nested inside it.
func cellFor(t *testing.T, html, isoDate string) string {
	t.Helper()
	anchor := `data-date="` + isoDate + `"`
	i := strings.Index(html, anchor)
	if i < 0 {
		t.Fatalf("no cell for %s", isoDate)
	}
	start := strings.LastIndex(html[:i], "<td")
	if start < 0 {
		t.Fatalf("cell for %s has no opening <td", isoDate)
	}
	end := strings.Index(html[start:], "</td>")
	if end < 0 {
		t.Fatalf("cell for %s has no closing </td>", isoDate)
	}
	return html[start : start+end+len("</td>")]
}

func TestCalendarSingleSelection(t *testing.T) {
	sel := []time.Time{time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)}
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		sel, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	if !strings.Contains(got, `data-gsxui-calendar-selected="2026-01-15"`) {
		t.Error("root does not carry the selection")
	}
	if n := strings.Count(got, `aria-selected="true"`); n != 1 {
		t.Errorf("got %d aria-selected=true, want 1", n)
	}

	cell := cellFor(t, got, "2026-01-15")
	if !strings.Contains(cell, `data-selected`) {
		t.Error("selected cell missing data-selected")
	}
	// Single selection is neither an end nor a middle of a range.
	if !strings.Contains(cell, `data-selected-single="true"`) {
		t.Errorf("day button missing data-selected-single=\"true\"\nin: %s", cell)
	}
}

func TestCalendarMultipleSelection(t *testing.T) {
	sel := []time.Time{
		time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 9, 0, 0, 0, 0, time.UTC),
	}
	got := render(t, ui.Calendar("multiple", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		sel, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	if !strings.Contains(got, `data-gsxui-calendar-selected="2026-01-05,2026-01-09"`) {
		t.Error("root does not carry both selections comma-separated")
	}
	if n := strings.Count(got, `aria-selected="true"`); n != 2 {
		t.Errorf("got %d aria-selected=true, want 2", n)
	}
	// Multiple selection is not a range, so each selected day is "single".
	if n := strings.Count(got, `data-selected-single="true"`); n != 2 {
		t.Errorf("got %d data-selected-single=true, want 2", n)
	}
}

func TestCalendarRangeMarksStartMiddleEnd(t *testing.T) {
	got := render(t, ui.Calendar("range", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil,
		time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 8, 0, 0, 0, 0, time.UTC),
		time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	for _, want := range []string{
		`data-gsxui-calendar-from="2026-01-05"`,
		`data-gsxui-calendar-to="2026-01-08"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("root missing %q", want)
		}
	}

	// The range attributes are on the day button and carry "true"/"false" —
	// shadcn styles them with data-[range-start=true]: selectors.
	for date, want := range map[string]string{
		"2026-01-05": `data-range-start="true"`,
		"2026-01-06": `data-range-middle="true"`,
		"2026-01-07": `data-range-middle="true"`,
		"2026-01-08": `data-range-end="true"`,
	} {
		if cell := cellFor(t, got, date); !strings.Contains(cell, want) {
			t.Errorf("%s: missing %s\nin: %s", date, want, cell)
		}
	}
	// The ends are ends, not middles.
	if n := strings.Count(got, `data-range-middle="true"`); n != 2 {
		t.Errorf("got %d range-middle days, want 2", n)
	}
	// A range selection is never "selected-single".
	if strings.Contains(got, `data-selected-single="true"`) {
		t.Error("range mode must not mark any day selected-single")
	}

	// modifiers.selected (source map §4) is true for every day from/to
	// inclusive in range mode too — range_start/range_middle/range_end
	// layer on top of selected, they don't replace it (useRange.js's
	// isSelected calls rangeIncludesDate(selected, date, false, dateLib),
	// excludeEnds=false). Both ends and the two middle days must all carry
	// the cell's own data-selected/aria-selected="true".
	for _, date := range []string{"2026-01-05", "2026-01-06", "2026-01-07", "2026-01-08"} {
		cell := cellFor(t, got, date)
		// The exact value, not a bare Contains(cell, "data-selected"): that
		// substring is also a prefix of the button's own always-present
		// data-selected-single="true"/"false", so a bare check would pass
		// against any rendered cell whatsoever and prove nothing.
		if !strings.Contains(cell, `data-selected="true"`) {
			t.Errorf("%s: cell missing data-selected=\"true\"\nin: %s", date, cell)
		}
		if !strings.Contains(cell, `aria-selected="true"`) {
			t.Errorf("%s: cell missing aria-selected=\"true\"\nin: %s", date, cell)
		}
	}
	if n := strings.Count(got, `aria-selected="true"`); n != 4 {
		t.Errorf("got %d aria-selected=true cells, want 4 (the whole range, both ends included)", n)
	}
}

// TestCalendarRangeWrappingARowMarksBothRowEdgesSelected exercises
// calendarDayClass's three [data-selected=true]-gated row-wrap rounding
// selectors ([&:last-child…]/[&:first-child…]/today+selected), which only
// fire when the cell itself (not just the button) carries data-selected.
// TestCalendarRangeMarksStartMiddleEnd alone never reaches a row wrap (its
// range sits inside a single week), so it can't catch a regression here.
//
// With weekStartsOn=Sunday, the January 2026 grid's second row runs
// 2026-01-04 (Sun) through 2026-01-10 (Sat); the third row runs 2026-01-11
// (Sun) through 2026-01-17 (Sat). A range from 2026-01-09 to 2026-01-12
// straddles that row boundary: 2026-01-10 is the last cell of its row,
// 2026-01-11 is the first cell of the next.
func TestCalendarRangeWrappingARowMarksBothRowEdgesSelected(t *testing.T) {
	got := render(t, ui.Calendar("range", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil,
		time.Date(2026, 1, 9, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 12, 0, 0, 0, 0, time.UTC),
		time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	for _, date := range []string{"2026-01-09", "2026-01-10", "2026-01-11", "2026-01-12"} {
		cell := cellFor(t, got, date)
		// The exact value, not a bare Contains(cell, "data-selected"): every
		// cell's button unconditionally emits data-selected-single="true"/
		// "false", and "data-selected" is a substring of "data-selected-
		// single" too, so a bare presence check passes against any rendered
		// cell whatsoever regardless of whether the cell itself is marked —
		// it would have passed identically against the mode-scoped code
		// this test exists to catch. Anchor on the cell's own literal value.
		if !strings.Contains(cell, `data-selected="true"`) {
			t.Errorf("%s: cell missing data-selected=\"true\" — the row-wrap rounding selectors need this\nin: %s", date, cell)
		}
	}
}

func TestCalendarDisabledBounds(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC), // disabledBefore
		time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC), // disabledAfter
		nil, nil, "", nil))

	if !strings.Contains(got, `data-date="2026-01-09" data-disabled`) {
		t.Error("2026-01-09 should be disabled (before the bound)")
	}
	if strings.Contains(got, `data-date="2026-01-15" data-disabled`) {
		t.Error("2026-01-15 is inside the bounds and must not be disabled")
	}
	if !strings.Contains(got, `data-date="2026-01-21" data-disabled`) {
		t.Error("2026-01-21 should be disabled (after the bound)")
	}
}

func TestCalendarDisabledWeekdaysAndDates(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{},
		[]time.Time{time.Date(2026, 1, 14, 0, 0, 0, 0, time.UTC)},
		[]time.Weekday{time.Saturday},
		"", nil))

	if !strings.Contains(got, `data-date="2026-01-14" data-disabled`) {
		t.Error("the explicitly disabled date is not marked")
	}
	// 2026-01-03 is a Saturday.
	if !strings.Contains(got, `data-date="2026-01-03" data-disabled`) {
		t.Error("Saturdays should be disabled")
	}
}

// Disabled days stay in the grid so keyboard users are not stranded either
// side of a disabled span.
//
// Upstream's rule (source map §8 correction 5) is that the NATIVE disabled
// attribute is set only while the day is not the live-focused one; a
// disabled day that holds focus degrades to aria-disabled so focus is never
// yanked out of the grid. On the server nothing holds focus yet, so first
// paint always uses the native attribute. Task 6 owns the focused case.
func TestCalendarDisabledDaysStayInTheGrid(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC), time.Time{},
		nil, nil, "", nil))

	if n := len(gridDates(t, got)); n != 42 {
		t.Errorf("got %d cells, want 42 — disabled days must still render", n)
	}
}

func TestCalendarLabelCaption(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	if !strings.Contains(got, "January 2026") {
		t.Error("label caption missing the month and year")
	}
	for _, want := range []string{
		`data-gsxui-calendar-prev`,
		`data-gsxui-calendar-next`,
		`aria-label="Previous month"`,
		`aria-label="Next month"`,
		// The caption announces month changes. Source map §8 correction 4.
		`role="status"`,
		`aria-live="polite"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q", want)
		}
	}
	if strings.Contains(got, "<select") {
		t.Error("label caption must not render selects")
	}
}

// hasBareDisabledAttr reports whether tag contains a real, bare HTML
// `disabled` attribute (rendered by gsx as literal " disabled" with nothing
// else — see gsx's Writer.BoolAttr) as opposed to a `disabled:`-prefixed
// Tailwind variant token such as "disabled:pointer-events-none", which also
// contains the substring " disabled" but is always immediately followed by
// ':', never by a space or '>'. button.gsx's own `base` class constant (used
// by every nav button, since these compose base/variantClass("ghost")
// directly) unconditionally contains "disabled:pointer-events-none
// disabled:opacity-50" — a bare strings.Contains(tag, " disabled") check
// would match that class token on every render regardless of whether a real
// disabled attribute is present, which is exactly the "assertion string is a
// substring of a different, always-present attribute" failure mode this
// project has already rejected once (see this file's own todayCellMarker
// comment for the same class of bug). Anchoring on the character
// immediately after "disabled" is what makes this check able to fail.
func hasBareDisabledAttr(tag string) bool {
	for i := 0; ; {
		idx := strings.Index(tag[i:], " disabled")
		if idx < 0 {
			return false
		}
		pos := i + idx + len(" disabled")
		if pos >= len(tag) || tag[pos] != ':' {
			return true
		}
		i = pos
	}
}

// Nav buttons never take a native disabled attribute, at any bound, under
// any configuration — only aria-disabled + tabindex="-1", which keeps them
// reachable by keyboard. Source map §7.5 and §8 correction 6.
func TestCalendarNavBoundsUseAriaDisabledNotDisabled(t *testing.T) {
	// fromYear == toYear == 2026 pins navigation to a single year, so
	// January 2026 has no previous month.
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 2026, 2026,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	prev := got[strings.Index(got, "data-gsxui-calendar-prev"):]
	prev = prev[:strings.Index(prev, ">")]
	if !strings.Contains(prev, `aria-disabled="true"`) {
		t.Errorf("prev at the navigation bound should be aria-disabled\ngot: %s", prev)
	}
	if !strings.Contains(prev, `tabindex="-1"`) {
		t.Errorf("prev at the navigation bound should be tabindex=-1\ngot: %s", prev)
	}
	if hasBareDisabledAttr(prev) {
		t.Errorf("prev must not take the native disabled attribute\ngot: %s", prev)
	}
}

func TestCalendarDropdownCaption(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "dropdown", 2020, 2030,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	for _, want := range []string{
		`data-gsxui-calendar-month-select`,
		`data-gsxui-calendar-year-select`,
		`<option value="0"`,  // January
		`<option value="11"`, // December
		`<option value="2020"`,
		`<option value="2030"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q", want)
		}
	}
	// 11 years inclusive, not 10.
	if n := strings.Count(got, `data-gsxui-calendar-year-option`); n != 11 {
		t.Errorf("got %d year options, want 11", n)
	}
	if n := strings.Count(got, `data-gsxui-calendar-month-option`); n != 12 {
		t.Errorf("got %d month options, want 12", n)
	}
}

func TestCalendarHiddenInputSingle(t *testing.T) {
	sel := []time.Time{time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)}
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		sel, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "booking", nil))

	if !strings.Contains(got, `<input type="hidden" name="booking" value="2026-01-15"`) {
		t.Errorf("hidden input missing or wrong\nin: %s", got)
	}
}

func TestCalendarHiddenInputsRange(t *testing.T) {
	got := render(t, ui.Calendar("range", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil,
		time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 8, 0, 0, 0, 0, time.UTC),
		time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "stay", nil))

	for _, want := range []string{
		`name="stay" value="2026-01-05"`,
		`name="stay-to" value="2026-01-08"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q", want)
		}
	}
}

func TestCalendarNoNameRendersNoHiddenInput(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	if strings.Contains(got, `type="hidden"`) {
		t.Error("no name given, so no hidden input should render")
	}
}

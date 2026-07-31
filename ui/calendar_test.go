package ui_test

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/internal/stylegen"
	"github.com/gsxhq/gsxui/merge"
	"github.com/gsxhq/gsxui/ui"
)

// novaCalendarRecipe loads the default style's Calendar recipe once. ui.Calendar
// is generated output now, so the classes it renders are the default style's
// concrete utilities — same pattern as novaButtonRecipe in button_test.go.
var novaCalendarRecipe = sync.OnceValue(func() recipe.Style {
	path := filepath.Join("..", "registry", "styles", stylegen.DefaultStyle, "calendar.css")
	src, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	style, err := recipe.ParseStyle(path, src)
	if err != nil {
		panic(err)
	}
	return style
})

// canonicalCalendarClass is the class attribute one Calendar slot renders,
// plus any caller classes, merged through the same merger gsx uses.
func canonicalCalendarClass(slot string, caller ...string) string {
	class := "gsxui-recipe-calendar"
	if slot != "" {
		class += "-" + slot
	}
	rule, ok := novaCalendarRecipe().Lookup(class)
	if !ok {
		panic("default style declares no recipe " + class)
	}
	classes := append(append([]string(nil), rule.Utilities...), caller...)
	return `class="` + html.EscapeString(merge.Merge(classes)) + `"`
}

// grid returns the data-date values of every day BUTTON, in DOM order —
// scoped to the data-gsxui-calendar-day tag specifically, not a bare
// "every data-date in the document" split. The <td> cell carries its own
// data-date too (Task 1), so a naive split would double-count once Task 4
// adds the button's own copy (see calendar.gsx's own comment on why the
// button needs it — Task 5's click handler matches the button and needs
// its date without walking up to the cell, not any uniqueness argument):
// jstest/specs/calendar.spec.ts's browser-side gridDates/gridCells helpers
// read it directly off [data-gsxui-calendar-day] for the same reason.
func gridDates(t *testing.T, html string) []string {
	t.Helper()
	var out []string
	for _, part := range strings.Split(html, `data-gsxui-calendar-day`)[1:] {
		tagEnd := strings.Index(part, ">")
		if tagEnd < 0 {
			t.Fatalf("unterminated day button tag in %s", html)
		}
		tag := part[:tagEnd]
		i := strings.Index(tag, `data-date="`)
		if i < 0 {
			t.Fatalf("day button missing data-date\nbutton tag: %s", tag)
		}
		rest := tag[i+len(`data-date="`):]
		end := strings.Index(rest, `"`)
		if end < 0 {
			t.Fatalf("unterminated data-date in button tag: %s", tag)
		}
		out = append(out, rest[:end])
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
		`data-gsxui-slot-calendar`,
		`data-gsxui-calendar`,
		`data-gsxui-calendar-month="2026-07"`,
		`data-gsxui-calendar-mode="range"`,
		`data-gsxui-calendar-week-start="1"`,
		`data-caption-layout="label"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("root missing %q", want)
		}
	}
}

// TestCalendarStyleContract catches any styled anonymous table/grid part
// falling back to undocumented DOM depth or a library presentation class.
// The 42 day cells/buttons are also the fixed structure calendar.js repaints
// in place, so their styling and behavior contracts must stay one-to-one.
func TestCalendarStyleContract(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", gsx.Attrs{{Key: "class", Value: "rounded-none"}}))

	for token, count := range map[string]int{
		`data-gsxui-slot-calendar`:        1,
		`data-gsxui-slot-calendar-months`: 1,
		`data-gsxui-slot-calendar-nav`:    1,
		`data-gsxui-slot-calendar-previous data-gsxui-slot-calendar-nav-button data-gsxui-slot-button`: 1,
		`data-gsxui-slot-calendar-next data-gsxui-slot-calendar-nav-button data-gsxui-slot-button`:     1,
		`data-gsxui-slot-calendar-month-caption`:                                                       1,
		`data-gsxui-slot-calendar-caption`:                                                             1,
		`data-gsxui-slot-calendar-grid`:                                                                1,
		`data-gsxui-slot-calendar-weekdays`:                                                            1,
		`data-gsxui-slot-calendar-weekday`:                                                             7,
		`data-gsxui-slot-calendar-week`:                                                                6,
		`data-gsxui-slot-calendar-day`:                                                                 42,
		`data-gsxui-slot-calendar-day-button data-gsxui-slot-button`:                                   42,
	} {
		if n := countPresenceAttribute(got, token); n != count {
			t.Errorf("%s count = %d, want %d", token, n, count)
		}
	}
	rootStart := strings.Index(got, `data-gsxui-slot-calendar`)
	if rootStart < 0 {
		t.Fatal("calendar root has no style slot")
	}
	tagStart := strings.LastIndex(got[:rootStart], "<div")
	rootEnd := strings.Index(got[rootStart:], ">")
	rootTag := got[tagStart : rootStart+rootEnd]
	// Calendar is on the slot axis now, so the root renders its own resolved
	// utilities with the caller's class merged in after them — the caller's
	// rounded-none is still what wins, but it is no longer the whole
	// attribute (playbook step 9).
	wantRoot := canonicalCalendarClass("", "rounded-none")
	if !strings.Contains(rootTag, wantRoot) {
		t.Errorf("root caller class is not forwarded\nwant: %s\nroot: %s", wantRoot, rootTag)
	}
	// Not a bare Count of "rounded-none": day cells and day buttons legitimately
	// carry data-[today]:data-[selected=true]:rounded-none and
	// data-[range-middle=true]:rounded-none as variant tokens. The root's whole
	// merged attribute is what must appear exactly once.
	if strings.Count(got, wantRoot) != 1 {
		t.Errorf("caller class must merge onto the root exactly once\nwant: %s\nin: %s", wantRoot, got)
	}
}

func TestCalendarPreviousComposesAllPresenceMarkersOnButton(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))
	requirePresenceAttributesOnSameTag(t, got, "data-gsxui-calendar-prev",
		"data-gsxui-slot-calendar-previous",
		"data-gsxui-slot-calendar-nav-button",
		"data-gsxui-slot-button",
	)
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
	// The exact value, not a bare Contains(cell, "data-selected"): the
	// button inside every cell unconditionally emits
	// data-selected-single="true"/"false", and "data-selected" is a prefix
	// of that, so the bare check passed against any rendered cell whatsoever
	// and could never go red. (The Task 2 fix round tightened the identical
	// pattern in the two range tests below and missed this one.)
	if !strings.Contains(cell, `data-selected="true"`) {
		t.Errorf("selected cell missing data-selected=\"true\"\nin: %s", cell)
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

// The day button carries its own data-date, additive to (not instead of) the
// enclosing <td>'s copy — Task 4's addition, so Task 5's click handler and
// calendar.js's own in-place update can read a day's date straight off the
// button the user actually interacts with, without walking up to the
// enclosing cell first. Anchored on the button's own opening tag
// specifically, so this can't pass merely because the (different) enclosing
// <td> element happens to carry the same attribute already.
func TestCalendarDayButtonCarriesItsOwnDate(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	cell := cellFor(t, got, "2026-01-15")
	btnStart := strings.Index(cell, "<button")
	if btnStart < 0 {
		t.Fatal("cell has no button")
	}
	btnTagEnd := strings.Index(cell[btnStart:], ">")
	if btnTagEnd < 0 {
		t.Fatal("button's opening tag is never closed")
	}
	btnTag := cell[btnStart : btnStart+btnTagEnd]
	if !strings.Contains(btnTag, `data-date="2026-01-15"`) {
		t.Errorf("day button missing its own data-date\nbutton tag: %s", btnTag)
	}
}

// The four disabled rules land on the root as data attributes (Task 4) so
// calendar.js can re-derive dayDisabled for a month the server never
// rendered — Task 3's per-cell data-disabled only ever reflected the
// initial month.
func TestCalendarRootCarriesDisabledRules(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC),
		[]time.Time{time.Date(2026, 1, 14, 0, 0, 0, 0, time.UTC)},
		[]time.Weekday{time.Saturday, time.Sunday},
		"", nil))

	for _, want := range []string{
		`data-gsxui-calendar-disabled-before="2026-01-10"`,
		`data-gsxui-calendar-disabled-after="2026-01-20"`,
		`data-gsxui-calendar-disabled-dates="2026-01-14"`,
		`data-gsxui-calendar-disabled-weekdays="6,0"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("root missing %q\nin: %s", want, got)
		}
	}
}

// Each disabled-rule attribute is omitted entirely when unset, matching
// -selected/-from/-to's own convention — never emitted as ="".
func TestCalendarRootOmitsDisabledRulesWhenUnset(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	for _, unwanted := range []string{
		"data-gsxui-calendar-disabled-before",
		"data-gsxui-calendar-disabled-after",
		"data-gsxui-calendar-disabled-dates",
		"data-gsxui-calendar-disabled-weekdays",
	} {
		if strings.Contains(got, unwanted) {
			t.Errorf("root should omit %s when unset\nin: %s", unwanted, got)
		}
	}
}

// Nav bounds are unconditionally on the root — Task 4 review's Important
// finding 1: without them, calendar.js has no way to recompute
// calendarPrevDisabled/calendarNextDisabled after a client-side navigation,
// so a caller with a narrow fromYear/toYear range could be clicked straight
// past the declared bound into a month the server can never render for
// this component (and, in captionLayout="dropdown", past the year
// <select>'s own last <option>).
func TestCalendarRootCarriesResolvedNavBounds(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 2020, 2030,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	for _, want := range []string{
		`data-gsxui-calendar-nav-from-year="2020"`,
		`data-gsxui-calendar-nav-to-year="2030"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("root missing %q\nin: %s", want, got)
		}
	}
}

// A zero fromYear/toYear resolves through calendarNavBounds's own default
// (currentYear-100/currentYear+10) — the root must carry that RESOLVED
// value, not a bare "0" or an omitted attribute (unlike -selected/-from/-to/
// the disabled rules, a nav bound always has a value, so it is never
// omitted), since calendar.js has no other way to learn the default bound.
func TestCalendarRootCarriesDefaultNavBoundsWhenUnset(t *testing.T) {
	currentYear := time.Now().UTC().Year()
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	wantFrom := fmt.Sprintf(`data-gsxui-calendar-nav-from-year="%d"`, currentYear-100)
	wantTo := fmt.Sprintf(`data-gsxui-calendar-nav-to-year="%d"`, currentYear+10)
	if !strings.Contains(got, wantFrom) {
		t.Errorf("root missing %q\nin: %s", wantFrom, got)
	}
	if !strings.Contains(got, wantTo) {
		t.Errorf("root missing %q\nin: %s", wantTo, got)
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
// yanked out of the grid. The initial roving tab stop uses the same
// aria-disabled degradation so Tab can enter the grid; every other disabled
// day keeps the native attribute.
// The earlier version of this test asserted only `len(gridDates) == 42`,
// which is true of EVERY possible render of this component (the grid is a
// fixed 42 cells by construction — see monthGrid) and so could never go red
// for the property the name claims. It now checks the actual contract: the
// days the disabledBefore bound rules out are still present as cells, still
// carry their date, and are marked disabled rather than omitted.
func TestCalendarDisabledDaysStayInTheGrid(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC), time.Time{},
		nil, nil, "", nil))

	dates := gridDates(t, got)
	if len(dates) != 42 {
		t.Fatalf("got %d cells, want 42", len(dates))
	}

	// disabledBefore=2026-01-10 rules out the grid's first thirteen days:
	// 2025-12-28..2025-12-31 (four padding days) plus 2026-01-01..2026-01-09.
	var wantDisabled []string
	for d := time.Date(2025, 12, 28, 0, 0, 0, 0, time.UTC); d.Before(time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)); d = d.AddDate(0, 0, 1) {
		wantDisabled = append(wantDisabled, d.Format("2006-01-02"))
	}
	if len(wantDisabled) != 13 {
		t.Fatalf("test setup: expected 13 disabled days, computed %d", len(wantDisabled))
	}

	present := make(map[string]bool, len(dates))
	for _, d := range dates {
		present[d] = true
	}
	for _, d := range wantDisabled {
		if !present[d] {
			t.Errorf("%s was dropped from the grid instead of being rendered disabled", d)
			continue
		}
		if !hasBareAttr(cellTag(t, got, d), "data-disabled") {
			t.Errorf("%s: cell is not marked data-disabled\ntag: %s", d, cellTag(t, got, d))
		}
		button := buttonTag(t, got, d)
		if d == "2026-01-01" {
			if hasBareDisabledAttr(button) || !strings.Contains(button, `aria-disabled="true"`) {
				t.Errorf("%s: disabled roving stop must use aria-disabled without native disabled\ntag: %s", d, button)
			}
		} else if !hasBareDisabledAttr(button) {
			t.Errorf("%s: non-tab-stop disabled day carries no native disabled attribute\ntag: %s", d, button)
		}
	}

	// Exactly those thirteen, and no more — a dayDisabled that disabled the
	// whole grid would satisfy every assertion above.
	var disabledCount int
	for _, d := range dates {
		if hasBareAttr(cellTag(t, got, d), "data-disabled") {
			disabledCount++
		}
	}
	if disabledCount != len(wantDisabled) {
		t.Errorf("got %d disabled cells, want %d", disabledCount, len(wantDisabled))
	}

	// The first day inside the bound is untouched — the `<` vs `<=` boundary.
	if hasBareAttr(cellTag(t, got, "2026-01-10"), "data-disabled") {
		t.Error("2026-01-10 is the bound itself and must not be disabled")
	}
}

// --- final review, Important 1: showOutsideDays -----------------------------
//
// The parameter was declared in the fourteen-parameter signature and read
// nowhere: every cell rendered regardless, so `showOutsideDays={false}` had
// no observable effect at all. Upstream's semantics (react-day-picker
// 9.14.0, helpers/createGetModifiers.js) are not "drop the cell" —
// showOutsideDays:false sets modifiers.hidden, the <td> stays for layout
// (carrying its own data-hidden), classNames.hidden is "invisible", and no
// DayButton renders inside it. gsxui keeps the button element (calendar.js
// may never create or destroy one) and blanks it instead.

func TestCalendarShowOutsideDaysFalseHidesPaddingDaysKeepingTheirCells(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, false, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	dates := gridDates(t, got)
	if len(dates) != 42 {
		t.Fatalf("got %d cells, want 42 — hiding must not remove cells", len(dates))
	}

	// January 2026 with a Sunday week start spans 2025-12-28..2026-02-07:
	// four leading padding days and seven trailing ones, eleven in total.
	var hiddenCount int
	for _, d := range dates {
		tag := cellTag(t, got, d)
		outside := hasBareAttr(tag, "data-outside")
		hidden := hasBareAttr(tag, "data-hidden")
		if outside != hidden {
			t.Errorf("%s: data-outside=%v but data-hidden=%v — with showOutsideDays=false the two must coincide", d, outside, hidden)
		}
		if hidden {
			hiddenCount++
		}
	}
	if hiddenCount != 11 {
		t.Errorf("got %d hidden cells, want 11", hiddenCount)
	}

	// The hidden day's button is blanked, not removed: no text, out of the
	// tab order, hidden from assistive tech, and not activatable.
	btn := buttonTag(t, got, "2025-12-28")
	if !strings.Contains(btn, `tabindex="-1"`) {
		t.Errorf("hidden day's button is not tabindex=-1\ntag: %s", btn)
	}
	if !strings.Contains(btn, `aria-hidden="true"`) {
		t.Errorf("hidden day's button is not aria-hidden\ntag: %s", btn)
	}
	if !hasBareDisabledAttr(btn) {
		t.Errorf("hidden day's button carries no native disabled attribute\ntag: %s", btn)
	}
	if cell := cellFor(t, got, "2025-12-28"); !strings.Contains(cell, "></button>") {
		t.Errorf("hidden day's button still renders its day number\nin: %s", cell)
	}

	// An in-month day is entirely untouched by any of this.
	inMonth := buttonTag(t, got, "2026-01-15")
	if strings.Contains(inMonth, "aria-hidden") {
		t.Errorf("an in-month day must not be aria-hidden\ntag: %s", inMonth)
	}
	if hasBareDisabledAttr(inMonth) {
		t.Errorf("an in-month day must not be disabled by hiding\ntag: %s", inMonth)
	}
	if cell := cellFor(t, got, "2026-01-15"); !strings.Contains(cell, ">15</button>") {
		t.Errorf("an in-month day must still render its day number\nin: %s", cell)
	}

	// Hidden days are OUT of the roving sequence (unlike disabled days,
	// which stay in it) — exactly one tab stop remains, and it is an
	// in-month day.
	if n := strings.Count(got, `tabindex="0"`); n != 1 {
		t.Errorf("got %d tabindex=\"0\" buttons, want exactly 1", n)
	}
	if !strings.Contains(buttonTag(t, got, "2026-01-01"), `tabindex="0"`) {
		t.Error("the sole tab stop is not the first day of the displayed month")
	}
}

func TestCalendarShowOutsideDaysTrueMarksNothingHidden(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	for _, d := range gridDates(t, got) {
		if hasBareAttr(cellTag(t, got, d), "data-hidden") {
			t.Errorf("%s: marked data-hidden even though showOutsideDays is true", d)
		}
	}
	if !strings.Contains(cellFor(t, got, "2025-12-28"), ">28</button>") {
		t.Error("an outside day must still render its day number when showOutsideDays is true")
	}
}

// The resolved flag reaches the client on the root, unconditionally (like the
// nav bounds, and unlike the omit-when-unset selection/disabled attributes —
// a bool always has a value). calendar.js re-derives `hidden` for every month
// the server never rendered and has no other source for it.
func TestCalendarRootCarriesShowOutsideDays(t *testing.T) {
	for _, tc := range []struct {
		show bool
		want string
	}{
		{true, `data-gsxui-calendar-show-outside-days="true"`},
		{false, `data-gsxui-calendar-show-outside-days="false"`},
	} {
		got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			nil, time.Time{}, time.Time{}, time.Sunday, tc.show, "label", 0, 0,
			time.Time{}, time.Time{}, nil, nil, "", nil))
		if !strings.Contains(got, tc.want) {
			t.Errorf("showOutsideDays=%v: root missing %q", tc.show, tc.want)
		}
	}
}

// --- final review, Important 2: mode's zero value ---------------------------
//
// data-gsxui-calendar-mode was defaulted to "single" on the way out while the
// body kept reading the bare parameter, and daySelected returns false for any
// mode that is neither "single" nor "multiple". So a caller who omitted
// `mode` (Go's zero value "") got the selected day rendered UNSELECTED while
// the root attribute and the hidden input both carried that same date — and
// calendar.js's own `|| "single"` default then flipped it on the first
// repaint. Server and client disagreed on first paint.
func TestCalendarZeroModeRendersAsSingleEverywhereNotJustOnTheRoot(t *testing.T) {
	sel := []time.Time{time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)}
	got := render(t, ui.Calendar("", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		sel, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "d", nil))

	if !strings.Contains(got, `data-gsxui-calendar-mode="single"`) {
		t.Error("root does not carry the defaulted mode")
	}
	cell := cellFor(t, got, "2026-01-15")
	if !strings.Contains(cell, `data-selected="true"`) {
		t.Errorf("the selected day's cell is not marked selected\nin: %s", cell)
	}
	if !strings.Contains(cell, `aria-selected="true"`) {
		t.Errorf("the selected day's cell is not aria-selected\nin: %s", cell)
	}
	if !strings.Contains(cell, `data-selected-single="true"`) {
		t.Errorf("the selected day's button is not data-selected-single\nin: %s", cell)
	}
	// The three carriers of the same fact agree.
	if !strings.Contains(got, `data-gsxui-calendar-selected="2026-01-15"`) {
		t.Error("root does not carry the selection")
	}
	if !strings.Contains(got, `<input type="hidden" name="d" value="2026-01-15"`) {
		t.Errorf("hidden input does not carry the selection\nin: %s", got)
	}
	// aria-multiselectable is not affected by the same bug either way ("" and
	// "single" both yield false), pinned here so a future refactor that reads
	// mode before defaulting it can't silently regress this one too.
	if strings.Contains(got, `aria-multiselectable`) {
		t.Error("single mode must not be aria-multiselectable")
	}
}

// --- final review, Important 3: UTC comparison, local-zone serialization ----
//
// sameDay/dayOnly normalize through .UTC(), but every serialization used to
// be t.Format("2006-01-02") in the VALUE's own location. A server in
// Asia/Tokyo handed selected = 2026-01-15 07:00 JST (= 2026-01-14 22:00 UTC)
// marked 2026-01-14 in the grid while the root attribute and the hidden
// input both said 2026-01-15: highlight on one day, form value on another,
// and the first client repaint (which re-derives selection from the root
// attribute) silently moved the highlight. The component's own doc comment
// promises non-UTC callers are handled; before this they were handled for
// comparison only.
func TestCalendarSerializesTheSameUTCDayItCompares(t *testing.T) {
	jst := time.FixedZone("JST", 9*60*60)
	// Each of these is a UTC-calendar-day BEHIND its own local date.
	sel := time.Date(2026, 1, 15, 7, 0, 0, 0, jst)   // 2026-01-14 UTC
	from := time.Date(2026, 1, 9, 3, 0, 0, 0, jst)   // 2026-01-08 UTC
	to := time.Date(2026, 1, 13, 5, 0, 0, 0, jst)    // 2026-01-12 UTC
	before := time.Date(2026, 1, 6, 2, 0, 0, 0, jst) // 2026-01-05 UTC
	after := time.Date(2026, 1, 26, 8, 0, 0, 0, jst) // 2026-01-25 UTC
	dd := time.Date(2026, 1, 21, 4, 0, 0, 0, jst)    // 2026-01-20 UTC

	// Single mode first: the marked day and both serializations must agree.
	single := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		[]time.Time{sel}, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		before, after, []time.Time{dd}, nil, "d", nil))

	if !strings.Contains(cellFor(t, single, "2026-01-14"), `data-selected="true"`) {
		t.Error("2026-01-14 (the UTC day sameDay compares) is not the marked day")
	}
	for _, want := range []string{
		`data-gsxui-calendar-selected="2026-01-14"`,
		`data-gsxui-calendar-disabled-before="2026-01-05"`,
		`data-gsxui-calendar-disabled-after="2026-01-25"`,
		`data-gsxui-calendar-disabled-dates="2026-01-20"`,
		`<input type="hidden" name="d" value="2026-01-14"`,
	} {
		if !strings.Contains(single, want) {
			t.Errorf("missing %q — serialized in the caller's own zone instead of through dayOnly", want)
		}
	}
	// The local-zone dates must appear nowhere as a serialized value.
	for _, unwanted := range []string{`="2026-01-15"`, `="2026-01-06"`, `="2026-01-26"`, `="2026-01-21"`} {
		if strings.Contains(single, `data-gsxui-calendar-selected`+unwanted) ||
			strings.Contains(single, `data-gsxui-calendar-disabled-before`+unwanted) ||
			strings.Contains(single, `data-gsxui-calendar-disabled-after`+unwanted) ||
			strings.Contains(single, `data-gsxui-calendar-disabled-dates`+unwanted) {
			t.Errorf("a root attribute was serialized in the caller's own zone: %s", unwanted)
		}
	}
	// The disabled RULES the client re-derives must rule out the same days
	// the server itself marked.
	if !hasBareAttr(cellTag(t, single, "2026-01-20"), "data-disabled") {
		t.Error("the disabled date's own UTC day is not marked disabled")
	}

	// Range mode: from/to on the root and in both hidden inputs.
	rng := render(t, ui.Calendar("range", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, from, to, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "stay", nil))

	for _, want := range []string{
		`data-gsxui-calendar-from="2026-01-08"`,
		`data-gsxui-calendar-to="2026-01-12"`,
		`name="stay" value="2026-01-08"`,
		`name="stay-to" value="2026-01-12"`,
	} {
		if !strings.Contains(rng, want) {
			t.Errorf("missing %q — serialized in the caller's own zone instead of through dayOnly", want)
		}
	}
	if !strings.Contains(cellFor(t, rng, "2026-01-08"), `data-range-start="true"`) {
		t.Error("2026-01-08 (from's own UTC day) is not the range start")
	}
	if !strings.Contains(cellFor(t, rng, "2026-01-12"), `data-range-end="true"`) {
		t.Error("2026-01-12 (to's own UTC day) is not the range end")
	}
}

// --- final review, minor: the grid's accessible name ------------------------
//
// role="grid" suppresses a <table>'s own implicit naming, so without this the
// grid is announced unnamed on entry. Upstream sets aria-label={labelGrid(…)},
// which defaults to the formatted month/year — the same string the caption
// already shows.
func TestCalendarGridIsNamedAfterTheDisplayedMonth(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 7, 4, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	i := strings.Index(got, "<table")
	if i < 0 {
		t.Fatal("no <table> rendered")
	}
	end := strings.Index(got[i:], ">")
	if end < 0 {
		t.Fatal("the <table>'s opening tag is never closed")
	}
	tag := got[i : i+end]
	if !strings.Contains(tag, `aria-label="July 2026"`) {
		t.Errorf("the grid has no accessible name naming the displayed month\ntag: %s", tag)
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

// hasBareAttr reports whether tag contains a real, bare HTML attribute named
// `name` (rendered by gsx as a literal " name" with nothing else — see gsx's
// Writer.BoolAttr), as opposed to a same-named Tailwind variant token such as
// "disabled:pointer-events-none" or "data-hidden:invisible" (always followed
// by ':'), a longer attribute that merely starts with the same characters
// ("data-selected-single", followed by '-'), or a valued attribute of the
// same name (followed by '='). button.gsx's own `base` class constant, on
// every button here, unconditionally contains "disabled:pointer-events-none
// disabled:opacity-50", and calendarDayClass contains "data-disabled:…" and
// "data-hidden:invisible" — a bare strings.Contains(tag, " disabled") check
// would match those class tokens on every render regardless of whether a
// real attribute is present, which is exactly the "assertion string is a
// substring of a different, always-present attribute" failure mode this
// project has already rejected repeatedly (see this file's own
// todayCellMarker comment for the same class of bug). Anchoring on the
// character immediately after the name is what makes this check able to fail.
func hasBareAttr(tag, name string) bool {
	for i := 0; ; {
		idx := strings.Index(tag[i:], " "+name)
		if idx < 0 {
			return false
		}
		pos := i + idx + len(" ") + len(name)
		if pos >= len(tag) {
			return true
		}
		switch tag[pos] {
		case ':', '-', '=':
			i = pos
		default:
			return true
		}
	}
}

func hasBareDisabledAttr(tag string) bool { return hasBareAttr(tag, "disabled") }

// cellTag returns just the opening `<td …>` tag for the given ISO date, so a
// test can check the CELL's own attributes without the nested button's (or
// the class attribute's Tailwind variant tokens) coming along.
func cellTag(t *testing.T, html, isoDate string) string {
	t.Helper()
	cell := cellFor(t, html, isoDate)
	end := strings.Index(cell, ">")
	if end < 0 {
		t.Fatalf("cell for %s has no closing '>' on its opening tag", isoDate)
	}
	return cell[:end]
}

// buttonTag returns just the opening `<button …>` tag of the day button
// inside the given ISO date's cell.
func buttonTag(t *testing.T, html, isoDate string) string {
	t.Helper()
	cell := cellFor(t, html, isoDate)
	start := strings.Index(cell, "<button")
	if start < 0 {
		t.Fatalf("cell for %s has no button", isoDate)
	}
	end := strings.Index(cell[start:], ">")
	if end < 0 {
		t.Fatalf("button in cell for %s has no closing '>' on its opening tag", isoDate)
	}
	return cell[start : start+end]
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
		`data-gsxui-slot-native-select-wrapper`,
		`data-gsxui-slot-calendar-dropdowns`,
		`data-caption-layout="dropdown" ` + canonicalCalendarClass("caption") + ` data-gsxui-slot-calendar-caption`,
		`<option value="0"`,  // January
		`<option value="11"`, // December
		`<option value="2020"`,
		`<option value="2030"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q", want)
		}
	}
	if strings.Contains(got, "calendar-dropdown-root") {
		t.Errorf("redundant calendar dropdown root marker remains")
	}
	if n := strings.Count(got, `data-gsxui-slot-native-select-wrapper`); n != 2 {
		t.Errorf("got %d native select wrappers, want 2", n)
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

// TestCalendarHiddenInputSingleRendersEmptyWhenUnselected is Task 5 review's
// Critical fix, pinned: `<ui.Calendar mode="single" name="date"/>` with
// nothing preselected — the single most likely way a caller wires this
// component into a form — used to render NO hidden input at all (the
// earlier showHiddenSingle also required len(selected) > 0), so the form
// would submit with no "date" field whatsoever regardless of what the user
// clicked. The input must exist from first paint, with an empty value,
// exactly like ui.Select/ui.Combobox's own hidden bridges always exist.
func TestCalendarHiddenInputSingleRendersEmptyWhenUnselected(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "date", nil))

	if !strings.Contains(got, `<input type="hidden" name="date" value=""`) {
		t.Errorf("expected an empty-valued hidden input to render unconditionally\nin: %s", got)
	}
}

// TestCalendarHiddenInputsRangeRenderBothEvenWithOnlyFromSet is the range
// half of the same Critical fix: `to` unset used to mean no "date-to" input
// at all (showHiddenTo also required !to.IsZero()), so a range-mode form
// with only a `from` chosen submitted with the "-to" field entirely absent
// instead of present-but-empty.
func TestCalendarHiddenInputsRangeRenderBothEvenWithOnlyFromSet(t *testing.T) {
	got := render(t, ui.Calendar("range", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil,
		time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
		time.Time{},
		time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "stay", nil))

	for _, want := range []string{
		`name="stay" value="2026-01-05"`,
		`name="stay-to" value=""`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// TestCalendarHiddenToInputCarriesItsOwnMarker is Task 5 review's Minor 1:
// the "-to" input must be identifiable by something other than its own
// `name` string, since a caller's own `name` could legitimately end in
// "-to" (e.g. name="valid-to") and calendar.js's own endsWith("-to") check
// would then misidentify the FROM input as the TO input.
func TestCalendarHiddenToInputCarriesItsOwnMarker(t *testing.T) {
	got := render(t, ui.Calendar("range", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil,
		time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 8, 0, 0, 0, 0, time.UTC),
		time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "valid-to", nil))

	if !strings.Contains(got, `<input type="hidden" name="valid-to" value="2026-01-05">`) {
		t.Errorf("the FROM input (name=%q, an unlucky collision with the -to naming convention) "+
			"must not carry the -to marker\nin: %s", "valid-to", got)
	}
	if !strings.Contains(got, `<input type="hidden" name="valid-to-to" value="2026-01-08" data-gsxui-calendar-hidden-to>`) {
		t.Errorf("the TO input must carry the -to marker regardless of what name collides with it\nin: %s", got)
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

func TestCalendarZeroMonthDefaultsFromSelection(t *testing.T) {
	jst := time.FixedZone("JST", 9*60*60)
	selected := []time.Time{time.Date(2031, 4, 1, 2, 0, 0, 0, jst)} // 2031-03-31 UTC

	got := render(t, ui.Calendar("single", time.Time{},
		selected, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	if !strings.Contains(got, `data-gsxui-calendar-month="2031-03"`) {
		t.Errorf("zero month did not default to the selected UTC month\nin: %s", got)
	}
	if !strings.Contains(cellFor(t, got, "2031-03-31"), `data-selected="true"`) {
		t.Error("selected UTC day is not present and marked in the defaulted month")
	}
}

func TestCalendarZeroMonthDefaultsFromRangeStart(t *testing.T) {
	from := time.Date(2042, 8, 17, 0, 0, 0, 0, time.UTC)

	got := render(t, ui.Calendar("range", time.Time{},
		nil, from, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	if !strings.Contains(got, `data-gsxui-calendar-month="2042-08"`) {
		t.Errorf("zero month did not default to the range start\nin: %s", got)
	}
}

func TestCalendarZeroMonthDefaultsToToday(t *testing.T) {
	before := time.Now().UTC()
	got := render(t, ui.Calendar("single", time.Time{},
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))
	after := time.Now().UTC()

	beforeMonth := before.Format("2006-01")
	afterMonth := after.Format("2006-01")
	if !strings.Contains(got, `data-gsxui-calendar-month="`+beforeMonth+`"`) &&
		!strings.Contains(got, `data-gsxui-calendar-month="`+afterMonth+`"`) {
		t.Errorf("zero month did not default to today's UTC month (%s or %s)\nin: %s",
			beforeMonth, afterMonth, got)
	}
}

func TestCalendarDisabledInitialTabStopUsesAriaDisabled(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC), time.Time{}, nil, nil, "", nil))

	tabStop := buttonTag(t, got, "2026-01-01")
	if !strings.Contains(tabStop, `tabindex="0"`) {
		t.Fatalf("first in-month day is not the initial roving tab stop\nbutton: %s", tabStop)
	}
	if hasBareDisabledAttr(tabStop) {
		t.Errorf("disabled roving tab stop is natively disabled and cannot receive focus\nbutton: %s", tabStop)
	}
	if !strings.Contains(tabStop, `aria-disabled="true"`) {
		t.Errorf("disabled roving tab stop does not expose aria-disabled\nbutton: %s", tabStop)
	}

	nonTabStop := buttonTag(t, got, "2026-01-02")
	if !hasBareDisabledAttr(nonTabStop) {
		t.Errorf("a disabled day outside the roving tab stop should remain natively disabled\nbutton: %s", nonTabStop)
	}
}

func TestCalendarMultipleHiddenInputsCarryEveryValue(t *testing.T) {
	selected := []time.Time{
		time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 8, 0, 0, 0, 0, time.UTC),
	}
	got := render(t, ui.Calendar("multiple", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		selected, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "dates", nil))

	if n := strings.Count(got, `<input type="hidden" name="dates"`); n != 2 {
		t.Fatalf("got %d multiple-mode hidden inputs, want one per selected date\nin: %s", n, got)
	}
	for _, value := range []string{"2026-01-05", "2026-01-08"} {
		if !strings.Contains(got, `name="dates" value="`+value+`" data-gsxui-calendar-hidden-multiple`) {
			t.Errorf("missing repeated multiple-mode input for %s\nin: %s", value, got)
		}
	}
}

func TestCalendarMultipleHiddenInputKeepsEmptyPlaceholder(t *testing.T) {
	got := render(t, ui.Calendar("multiple", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "dates", nil))

	if !strings.Contains(got, `<input type="hidden" name="dates" value="" data-gsxui-calendar-hidden-multiple`) {
		t.Errorf("unselected multiple calendar must preserve an empty form field\nin: %s", got)
	}
}

// The `nav` slot (source map §2) is `absolute inset-x-0 top-0 ...`, meant to
// position against a `relative` ancestor — upstream's own `months` wrapper
// div (§2: "relative flex flex-col gap-4 md:flex-row"). gsxui doesn't render
// `months`/`month` wrapper elements at all (§9, single-month only), so the
// port collapses both into one wrapper div, and that wrapper — NOT the root —
// carries `relative`.
//
// Which element carries it is load-bearing, not incidental. An absolutely
// positioned element resolves against its containing block's PADDING box, and
// the root carries `p-2`. With `relative` on the root, nav's `top-0` pins 8px
// above the caption row it labels (measured: buttons at y=40 vs row at y=54
// centre). The wrapper sits inside the padding, so `top-0` there coincides
// with the caption. Drop the token entirely and the buttons position against
// whatever `relative` ancestor the CONSUMING page happens to have, or the
// viewport.
//
// Render tests pin the structural slot ancestry; the actual positioning is
// pinned by jstest/specs/calendar.spec.ts's computed geometry.
func TestCalendarNavHasARelativePositioningAncestor(t *testing.T) {
	got := render(t, ui.Calendar("single", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		nil, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))

	wrapperIdx := strings.Index(got, `data-gsxui-slot-calendar-months`)
	if wrapperIdx < 0 {
		t.Fatal("no calendar-months wrapper rendered inside the root")
	}
	navIdx := strings.Index(got[wrapperIdx:], `data-gsxui-slot-calendar-nav`)
	if navIdx < 0 {
		t.Fatal("calendar-nav is not nested under calendar-months")
	}
	// The positioning contract itself: calendar-months carries `relative`, so
	// calendar-nav's `absolute` resolves against it. Before Calendar migrated
	// this was asserted negatively (no class attribute at all, the rule lived
	// in default.css); now it is asserted on the resolved class directly.
	if !strings.Contains(got, canonicalCalendarClass("months")+` data-gsxui-slot-calendar-months`) {
		t.Errorf("calendar-months lost its relative positioning class\nin: %s", got)
	}
	if !strings.Contains(got, canonicalCalendarClass("nav")+` data-gsxui-slot-calendar-nav`) {
		t.Errorf("calendar-nav lost its absolute positioning class\nin: %s", got)
	}
}

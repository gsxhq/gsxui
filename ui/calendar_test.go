package ui_test

import (
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

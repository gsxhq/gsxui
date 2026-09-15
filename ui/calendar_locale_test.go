package ui_test

import (
	"encoding/json"
	"html"
	"strings"
	"testing"
	"time"

	"github.com/gsxhq/gsxui/ui"
)

var jan2026 = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func renderCalendar(t *testing.T, month time.Time, layout string, loc ui.CalendarLocale) string {
	t.Helper()
	return render(t, ui.Calendar("single", month, nil, time.Time{}, time.Time{}, time.Sunday, true, layout, 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", loc, nil))
}

// rootLocale decodes the data-gsxui-calendar-locale attribute the root carries.
func rootLocale(t *testing.T, got string) map[string]any {
	t.Helper()
	const key = `data-gsxui-calendar-locale="`
	i := strings.Index(got, key)
	if i < 0 {
		t.Fatalf("root missing %s\nin: %s", key, got)
	}
	rest := got[i+len(key):]
	end := strings.Index(rest, `"`)
	if end < 0 {
		t.Fatalf("unterminated locale attribute")
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(html.UnescapeString(rest[:end])), &out); err != nil {
		t.Fatalf("locale attribute is not JSON: %v\n%s", err, rest[:end])
	}
	return out
}

var arabic = ui.CalendarLocale{
	Months:        [12]string{"يناير", "فبراير", "مارس", "أبريل", "مايو", "يونيو", "يوليو", "أغسطس", "سبتمبر", "أكتوبر", "نوفمبر", "ديسمبر"},
	MonthsShort:   [12]string{"يناير", "فبراير", "مارس", "أبريل", "مايو", "يونيو", "يوليو", "أغسطس", "سبتمبر", "أكتوبر", "نوفمبر", "ديسمبر"},
	Weekdays:      [7]string{"الأحد", "الاثنين", "الثلاثاء", "الأربعاء", "الخميس", "الجمعة", "السبت"},
	WeekdaysShort: [7]string{"ح", "ن", "ث", "ر", "خ", "ج", "س"},
	Caption:       "{month} {year}",
	DayLabel:      "{weekday}، {day} {month} {year}",
	Digits:        "٠١٢٣٤٥٦٧٨٩",
}

func TestCalendarZeroLocaleIsEnglish(t *testing.T) {
	got := renderCalendar(t, jan2026, "dropdown", ui.CalendarLocale{})
	for _, want := range []string{
		`aria-label="Thursday, January 15, 2026"`,
		`aria-label="January 2026"`,
		`>Su<`, `>Mo<`,
		`>Jan<`, `>Dec<`,
		`>2026<`,
		`aria-label="Previous month"`, `aria-label="Next month"`,
		`aria-label="Month"`, `aria-label="Year"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("zero locale missing %s", want)
		}
	}
	loc := rootLocale(t, got)
	if loc["caption"] != "{month} {year}" || loc["dayLabel"] != "{weekday}, {month} {day}, {year}" || loc["digits"] != "" {
		t.Errorf("zero locale JSON = %v", loc)
	}
	months := loc["months"].([]any)
	if len(months) != 12 || months[0] != "January" || months[11] != "December" {
		t.Errorf("months = %v", months)
	}
	weekdays := loc["weekdays"].([]any)
	if len(weekdays) != 7 || weekdays[0] != "Sunday" || weekdays[6] != "Saturday" {
		t.Errorf("weekdays = %v", weekdays)
	}
}

func TestCalendarArabicLocale(t *testing.T) {
	got := renderCalendar(t, jan2026, "dropdown", arabic)
	for _, want := range []string{
		`aria-label="الخميس، ١٥ يناير ٢٠٢٦"`, // 2026-01-15, Thursday
		`aria-label="يناير ٢٠٢٦"`,            // grid label = caption
		`>يناير ٢٠٢٦<`,                       // caption span
		`>١٥<`,                               // day cell text
		`>ح<`, `>س<`,                         // header row
		`>٢٠٢٦<`,                 // year option text
		`data-date="2026-01-15"`, // never localised
		`value="2026"`,           // year option value never localised
	} {
		if !strings.Contains(got, want) {
			t.Errorf("arabic locale missing %s\nin: %s", want, got)
		}
	}
	if strings.Contains(got, `>15<`) {
		t.Errorf("day cell text still ASCII under Arabic digits")
	}
	loc := rootLocale(t, got)
	if loc["digits"] != "٠١٢٣٤٥٦٧٨٩" || loc["dayLabel"] != "{weekday}، {day} {month} {year}" {
		t.Errorf("arabic locale JSON = %v", loc)
	}
}

func TestCalendarYearFirstPattern(t *testing.T) {
	loc := ui.CalendarLocale{
		Months:        [12]string{"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"},
		MonthsShort:   [12]string{"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"},
		Weekdays:      [7]string{"日曜日", "月曜日", "火曜日", "水曜日", "木曜日", "金曜日", "土曜日"},
		WeekdaysShort: [7]string{"日", "月", "火", "水", "木", "金", "土"},
		Caption:       "{year}年{month}",
		DayLabel:      "{year}年{month}{day}日{weekday}",
	}
	got := renderCalendar(t, jan2026, "label", loc)
	for _, want := range []string{
		`>2026年1月<`,
		`aria-label="2026年1月15日木曜日"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("year-first pattern missing %s", want)
		}
	}
}

func TestCalendarPatternIgnoresUnknownPlaceholdersAndBlanksAbsentParts(t *testing.T) {
	loc := arabic
	loc.Caption = "{weekday}[{month} {year}]{nope}"
	got := renderCalendar(t, jan2026, "label", loc)
	if !strings.Contains(got, `>[يناير ٢٠٢٦]{nope}<`) {
		t.Errorf("caption pattern handling wrong\nin: %s", got)
	}
}

func TestCalendarDigitsWithWrongLengthStayASCII(t *testing.T) {
	loc := arabic
	loc.Digits = "٠١٢"
	got := renderCalendar(t, jan2026, "label", loc)
	if !strings.Contains(got, `>15<`) || !strings.Contains(got, `aria-label="الخميس، 15 يناير 2026"`) {
		t.Errorf("a Digits value that is not ten runes must render ASCII\nin: %s", got)
	}
}

func TestCalendarDigitsNeverReachFormValues(t *testing.T) {
	selected := []time.Time{time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)}
	got := render(t, ui.Calendar("single", jan2026, selected, time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "d", arabic, nil))
	for _, want := range []string{
		`name="d" value="2026-01-15"`,
		`data-gsxui-calendar-month="2026-01"`,
		`data-gsxui-calendar-selected="2026-01-15"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("form/data value localised or missing: %s\nin: %s", want, got)
		}
	}
}

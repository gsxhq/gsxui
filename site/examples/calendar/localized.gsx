package calendar

import (
	"time"

	"github.com/gsxhq/gsxui/ui"
)

// LocalizedDefaultMonth mirrors Basic's own DefaultMonth (2026-01).
var LocalizedDefaultMonth = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// Arabic is CLDR's ar gregorian data: wide month names (the abbreviated set
// is identical), wide and narrow weekday names, and Arabic-Indic digits.
// The day label puts the day before the month, and the week starts on
// Saturday.
var Arabic = ui.CalendarLocale{
	Months: [12]string{
		"يناير",
		"فبراير",
		"مارس",
		"أبريل",
		"مايو",
		"يونيو",
		"يوليو",
		"أغسطس",
		"سبتمبر",
		"أكتوبر",
		"نوفمبر",
		"ديسمبر",
	},
	MonthsShort: [12]string{
		"يناير",
		"فبراير",
		"مارس",
		"أبريل",
		"مايو",
		"يونيو",
		"يوليو",
		"أغسطس",
		"سبتمبر",
		"أكتوبر",
		"نوفمبر",
		"ديسمبر",
	},
	Weekdays:      [7]string{"الأحد", "الاثنين", "الثلاثاء", "الأربعاء", "الخميس", "الجمعة", "السبت"},
	WeekdaysShort: [7]string{"ح", "ن", "ث", "ر", "خ", "ج", "س"},
	Caption:       "{month} {year}",
	DayLabel:      "{weekday}، {day} {month} {year}",
	Digits:        "٠١٢٣٤٥٦٧٨٩",
}

// Localized renders the dropdown-caption grid in Arabic under dir="rtl".
// The month is a parameter so the harness's ?month= override can render any
// month for the Go/JS agreement diff.
component Localized(month time.Time) {
	<div dir="rtl" lang="ar">
		<ui.Calendar
			mode="single"
			month={month}
			weekStartsOn={time.Saturday}
			showOutsideDays={true}
			captionLayout="dropdown"
			locale={Arabic}
		/>
	</div>
}

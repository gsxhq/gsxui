package calendar

import (
	"time"

	"github.com/gsxhq/gsxui/ui"
)

// RussianDefaultMonth mirrors Basic's own DefaultMonth (2026-01).
var RussianDefaultMonth = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// Russian is CLDR's ru gregorian data. A month standing alone ("январь",
// the caption) differs from a month inside a date ("января", each day's
// label), so MonthsFormat carries the second form. The week starts on
// Monday.
var Russian = ui.CalendarLocale{
	Months: [12]string{
		"январь",
		"февраль",
		"март",
		"апрель",
		"май",
		"июнь",
		"июль",
		"август",
		"сентябрь",
		"октябрь",
		"ноябрь",
		"декабрь",
	},
	MonthsFormat: [12]string{
		"января",
		"февраля",
		"марта",
		"апреля",
		"мая",
		"июня",
		"июля",
		"августа",
		"сентября",
		"октября",
		"ноября",
		"декабря",
	},
	MonthsShort: [12]string{
		"янв.",
		"февр.",
		"март",
		"апр.",
		"май",
		"июнь",
		"июль",
		"авг.",
		"сент.",
		"окт.",
		"нояб.",
		"дек.",
	},
	Weekdays:      [7]string{"воскресенье", "понедельник", "вторник", "среда", "четверг", "пятница", "суббота"},
	WeekdaysShort: [7]string{"вс", "пн", "вт", "ср", "чт", "пт", "сб"},
	Caption:       "{month} {year}",
	DayLabel:      "{weekday}, {day} {month} {year} г.",
}

// LocalizedRussian renders the dropdown-caption grid in Russian. The month is a
// parameter so the harness's ?month= override can render any month for the
// Go/JS agreement diff.
component LocalizedRussian(month time.Time) {
	<div lang="ru">
		<ui.Calendar
			mode="single"
			month={month}
			weekStartsOn={time.Monday}
			showOutsideDays={true}
			captionLayout="dropdown"
			locale={Russian}
		/>
	</div>
}

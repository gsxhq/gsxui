package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Calendar is the largest shape in the catalogue: thirteen slots, no
// dimensions. Every data-* the rules key on (data-outside, data-today,
// data-selected, data-disabled, data-hidden, data-focused, data-range-start/
// -middle/-end, data-selected-single, aria-selected, aria-disabled) is
// SERVER- or JS-computed day state, not a caller-chosen variant — the same
// call Select made for data-state/data-placeholder — so all of it folds onto
// its slot's one rule as native/arbitrary Tailwind variants.
//
// EVERY slot Calendar stamps that carries a rule is declared here, including
// the compound ones. That is load-bearing, not tidiness: the class decoder
// resolves a recipe class by LONGEST match against this list, so leaving
// "day-button", "nav-button", "month-caption" or "weekdays" undeclared would
// not fail loudly — gsxui-recipe-calendar-day-button would silently decode as
// slot "day" and the rest would be read as a dimension. The four compound
// names that shadow a shorter declared sibling are exactly:
//
//	nav / nav-button        month-caption (vs months)
//	day / day-button        weekday / weekdays (vs week)
//
// Three markers Calendar stamps declare NO slot, because no style carries a
// rule for them — calendar-previous and calendar-next (they sit on the same
// two <button>s as calendar-nav-button, which is where all of the presentation
// lives, and exist for calendar.js's wiring), and the elements carrying
// data-gsxui-slot-button. Same call select-bridge/select-group already made.
//
// Calendar's dropdowns compose <ui.NativeSelect>; the rules whose styled
// SUBJECT is a NativeSelect element are NOT here. They live, unwrapped, in
// assets/css/styles/default.css's @layer utilities — put there when
// NativeSelect migrated, because they now have to outrank NativeSelect's own
// resolved utilities. Migrating them onto Calendar's recipe would attach
// foreign behaviour to Calendar's classes; the styled-subject test says they
// are Calendar's rules about someone else's element, which is exactly the
// escape hatch's job.
var Calendar = recipe.Shape{
	Component: "calendar",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "months", Base: true},
		{Name: "nav", Base: true},
		{Name: "nav-button", Base: true},
		{Name: "month-caption", Base: true},
		{Name: "dropdowns", Base: true},
		{Name: "caption", Base: true},
		{Name: "grid", Base: true},
		{Name: "weekdays", Base: true},
		{Name: "weekday", Base: true},
		{Name: "week", Base: true},
		{Name: "day", Base: true},
		{Name: "day-button", Base: true},
	},
}

package ui

import "github.com/gsxhq/gsx"

// Radio is the shadcn/ui RadioGroupItem, ported as a real native
// <input type="radio">: form-native, zero JS, browser :checked/:disabled
// truth replaces Radix's button-role + hidden-input + Indicator part
// (ledger ADAPT). shadcn's RadioGroup container (a styled grid wrapper) is
// not ported — native `name` grouping on sibling <input type="radio">
// elements already gives you the group; the layout wrapper is the caller's
// concern, same as any other flex/grid container (ledger ADAPT).
//
// Tokens are carried verbatim from RadioGroupItem (aspect-square size-4
// shrink-0 rounded-full border border-input text-primary shadow-xs
// transition-[color,box-shadow] outline-none + focus-visible/disabled/
// aria-invalid/dark tokens); appearance-none is added for the same
// mechanical reason as checkbox.
//
// The checked paint follows the nova style (the live site's default, per
// the density-retarget decision): the whole circle fills with primary
// (checked:bg-primary checked:border-primary) and the indicator is a
// primary-FOREGROUND dot punched into it — nova's .cn-radio-group-item /
// .cn-radio-group-indicator-icon (`data-checked:bg-primary` + a size-2
// bg-primary-foreground dot), reading as a bold donut. (new-york-v4's
// older outlined-circle-with-primary-dot recipe is superseded.) The dot is
// still a checked:bg-[radial-gradient(...)] painted in currentColor — a
// data-URI can't reference the caller's CSS custom properties, but a
// currentColor gradient can, and checked:text-primary-foreground is what
// makes currentColor resolve to the dot's color; it is load-bearing, the
// same role text-primary played for the old recipe. background-color and
// background-image are distinct properties, so checked:bg-primary and the
// gradient coexist (and tailwind-merge classifies them into different
// conflict groups). The radial-gradient's embedded spaces are Tailwind's
// underscore escape for whitespace inside a bracketed arbitrary value, not
// literal spaces — see checkbox.gsx and docs/jsx-parity.md.
component Radio(attrs gsx.Attrs) {
	<input
		type="radio"
		data-slot="radio"
		class="peer aspect-square size-4 shrink-0 appearance-none rounded-full border border-input shadow-xs transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 dark:bg-input/30 checked:border-primary checked:bg-primary checked:text-primary-foreground dark:checked:bg-primary checked:bg-[radial-gradient(circle_closest-side,currentColor_45%,transparent_50%)]"
		{ attrs... }
	/>
}

package lyra

import "github.com/gsxhq/gsx"

// Radio is the shadcn/ui RadioGroupItem, ported as a real native
// <input type="radio">: form-native, zero JS, browser :checked/:disabled
// truth replaces Radix's button-role + hidden-input + Indicator part
// (ledger ADAPT). shadcn's RadioGroup container (a styled grid wrapper) is
// not ported — native `name` grouping on sibling <input type="radio">
// elements already gives you the group; the layout wrapper is the caller's
// concern, same as any other flex/grid container (ledger ADAPT).
//
// The checked paint follows the nova style (the live site's default, per
// the density-retarget decision): the whole circle fills with primary
// (checked:bg-primary checked:border-primary) and the indicator is a
// primary-FOREGROUND dot punched into it — nova's .cn-radio-group-item /
// .cn-radio-group-indicator-icon (`data-checked:bg-primary` + a size-2
// bg-primary-foreground dot), reading as a bold donut. (new-york-v4's
// older outlined-circle-with-primary-dot recipe is superseded.) The dot is
// still a radial-gradient painted in currentColor — a
// data-URI can't reference the caller's CSS custom properties, but a
// currentColor gradient can, and checked:text-primary-foreground is what
// makes currentColor resolve to the dot's color; it is load-bearing, the
// same role text-primary played for the old recipe. background-color and
// background-image are distinct properties, so the primary fill and
// gradient coexist.
component Radio(attrs gsx.Attrs) {
	<input
		type="radio"
		class={
			"border-input dark:bg-input/30 checked:bg-primary checked:text-primary-foreground dark:checked:bg-primary checked:border-primary aria-invalid:checked:border-primary aria-invalid:border-destructive focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 dark:aria-invalid:border-destructive/50 flex size-4 rounded-full focus-visible:ring-3 aria-invalid:ring-3 shrink-0 outline-none disabled:cursor-not-allowed disabled:opacity-50 border"
		}
		{ attrs... }
		data-gsxui-slot-radio
	/>
}

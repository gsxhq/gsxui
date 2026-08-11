package ui

import "github.com/gsxhq/gsx"

// Input is the shadcn/ui Input — a straight port of the native <input>.
// type="text" is an overridable default (the button type="button" pattern) —
// pass type="email" etc. at the call site. Void, childless element: the
// explicit { attrs... } spread is what opts it into fallthrough.
component Input(attrs gsx.Attrs) {
	<input
		type="text"
		class={
			"w-full min-w-0 outline-none file:inline-flex file:border-0 file:bg-transparent file:text-foreground placeholder:text-muted-foreground disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 border-transparent border-b-input bg-transparent focus-visible:border-b-ring aria-invalid:border-b-destructive dark:aria-invalid:border-b-destructive/50 h-10 border px-0 py-1 text-base transition-[color,border-color] file:h-7 file:text-sm file:font-medium md:text-sm"
		}
		{ attrs... }
		data-gsxui-slot-input
	/>
}

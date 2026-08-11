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
			"w-full min-w-0 outline-none file:inline-flex file:border-0 file:bg-transparent file:text-foreground placeholder:text-muted-foreground disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 bg-input/20 dark:bg-input/30 border-input focus-visible:border-ring focus-visible:ring-ring/30 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive dark:aria-invalid:border-destructive/50 h-7 rounded-md border px-2 py-0.5 text-sm transition-colors file:h-6 file:text-xs/relaxed file:font-medium focus-visible:ring-2 aria-invalid:ring-2 md:text-xs/relaxed"
		}
		{ attrs... }
		data-gsxui-slot-input
	/>
}

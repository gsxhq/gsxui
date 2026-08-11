package lyra

import "github.com/gsxhq/gsx"

// Input is the shadcn/ui Input — a straight port of the native <input>.
// type="text" is an overridable default (the button type="button" pattern) —
// pass type="email" etc. at the call site. Void, childless element: the
// explicit { attrs... } spread is what opts it into fallthrough.
component Input(attrs gsx.Attrs) {
	<input
		type="text"
		class={
			"dark:bg-input/30 border-input focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive dark:aria-invalid:border-destructive/50 disabled:bg-input/50 dark:disabled:bg-input/80 h-8 rounded-none border bg-transparent px-2.5 py-1 text-xs transition-colors file:h-6 file:text-xs file:font-medium focus-visible:ring-1 aria-invalid:ring-1 md:text-xs"
		}
		{ attrs... }
		data-gsxui-slot-input
	/>
}

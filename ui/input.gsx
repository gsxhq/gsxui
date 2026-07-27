package ui

import "github.com/gsxhq/gsx"

// Input is the shadcn/ui Input — a straight port of the native <input>.
// type="text" is an overridable default (the button type="button" pattern) —
// pass type="email" etc. at the call site. Void, childless element: the
// explicit { attrs... } spread is what opts it into fallthrough.
component Input(attrs gsx.Attrs) {
	<input
		type="text"
		{ withSlot("input", attrs)... }
	/>
}

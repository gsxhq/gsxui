package ui

import "github.com/gsxhq/gsx"

// Separator is the shadcn/ui Separator. orientation: "" (default,
// "horizontal") | "vertical" — stamped onto data-orientation, which one
// verbatim class string dispatches on via Tailwind's data-[orientation=...]
// selectors, so no variant switch is needed. Radix's decorative prop (its
// default true) is not ported — see docs/jsx-parity.md.
component Separator(orientation string, attrs gsx.Attrs) {
	<div
		role="none"
		data-orientation={orientation |> default("horizontal")}
		{ attrs... } data-gsxui-slot-separator
	></div>
}

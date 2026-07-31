package canonical

import "github.com/gsxhq/gsx"

// Separator is the shadcn/ui Separator. orientation: "" (default,
// "horizontal") | "vertical" — stamped onto data-orientation, and also
// drives the slot axis's own orientation dimension below, resolved to
// concrete utilities at generation time. Radix's decorative prop (its
// default true) is not ported — see docs/jsx-parity.md.
component Separator(orientation string, attrs gsx.Attrs) {
	<div
		role="none"
		data-orientation={orientation |> default("horizontal")}
		class={ separator.Root(), separator.Orientation(orientation) }
		{ attrs... }
		data-gsxui-slot-separator
	></div>
}

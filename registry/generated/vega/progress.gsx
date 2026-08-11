package ui

import "github.com/gsxhq/gsx"

// Progress is the shadcn/ui Progress. shadcn wraps Radix's
// ProgressPrimitive.Root/Indicator pair (registry/new-york-v4/ui/progress.tsx);
// this port replaces both with two plain divs, no client JS or component
// state (ADAPT — see docs/jsx-parity.md). role="progressbar" plus
// aria-valuemin/aria-valuemax/aria-valuenow replace what Radix's Root
// stamps internally. value is a Go zero-value float64 (0-100); 0 matches
// shadcn's own `value || 0` fallback, so an unset value renders the
// indicator fully translated off-screen exactly like the original.
//
// Radix's Indicator drives its fill via
// `style={{ transform: translateX(-${100 - (value || 0)}%) }}` — ported
// verbatim as the same translateX mechanism (not width, which would clip
// transition-all's animation differently).
//
// The declaration is a css`` literal (MECHANISM) rather than a Go
// concatenation: gw's CSS value filter — a port of html/template's —
// blocklists "(" and ")", and Go's + folds the static translateX(…) syntax
// and the dynamic percentage into one string the filter then judges as a
// unit, which is what previously forced gsx.RawCSS over the whole
// declaration. In a css`` literal the function call is static template text
// and only the percentage is a hole, so the filter still judges it and
// nothing is trusted that the component did not itself compute. The hole
// renders the float64 directly, the same way aria-valuenow does.
component Progress(value float64, attrs gsx.Attrs) {
	<div
		role="progressbar"
		aria-valuemin="0"
		aria-valuemax="100"
		aria-valuenow={value}
		class={ "bg-muted h-1.5 rounded-full" }
		{ attrs... }
		data-gsxui-slot-progress
	>
		<div
			style=css`transform: translateX(-@{100 - value}%)`
			class={ "bg-primary" }
			data-gsxui-slot-progress-indicator
		></div>
	</div>
}

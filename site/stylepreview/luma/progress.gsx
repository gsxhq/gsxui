package luma

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
//
// FIX (found reviewing theme-creator-parity Task 4's new preview cards):
// the indicator's own recipe accessor only ever resolved to a per-style
// colour, nothing else — no width or height anywhere. Every style's
// registry/styles/<style>/progress.css only ever sets the ROOT's
// bg-muted/height/radius and the indicator's bg-primary, never the
// indicator's size, because those are structural (invariant across every
// style), not presentational. Upstream's own source confirms this split:
// registry/bases/radix/ui/progress.tsx's Root carries "relative flex w-full
// items-center overflow-x-hidden" and Indicator carries "size-full flex-1
// transition-all" as literal, non-recipe classes — the same class-plus-
// recipe-accessor shape Card's own Root already uses for "group/card".
// Without this, translateX's math was correct but had nothing visible to
// translate: the indicator rendered at 0×0 in every style, so Progress
// silently never showed a fill.
component Progress(value float64, attrs gsx.Attrs) {
	<div
		role="progressbar"
		aria-valuemin="0"
		aria-valuemax="100"
		aria-valuenow={value}
		class={ "relative flex w-full items-center overflow-x-hidden", "bg-muted h-3 rounded-full" }
		{ attrs... }
		data-gsxui-slot-progress
	>
		<div
			style=css`transform: translateX(-@{100 - value}%)`
			class={ "size-full flex-1 transition-all", "bg-primary" }
			data-gsxui-slot-progress-indicator
		></div>
	</div>
}

package ui

import (
	"strconv"

	"github.com/gsxhq/gsx"
)

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
// transition-all's animation differently), computed here as
// strconv.FormatFloat(100-value, ...) since a float64 param can't itself
// concatenate into a string. The composed style value is wrapped in
// gsx.RawCSS (MECHANISM, same precedent as AspectRatio's ratio property —
// see ui/aspect-ratio.gsx and its docs/jsx-parity.md entry) to opt it out
// of gw's CSS value filter, which blocklists "(" and ")" — punctuation
// translateX(...)'s function-call syntax requires, not injected data. The
// percentage is trusted, developer-computed layout intent, the same trust
// boundary aspect-ratio's ratio already extends.
component Progress(value float64, attrs gsx.Attrs) {
	{{ remaining := strconv.FormatFloat(100-value, 'f', -1, 64) }}
	<div
		data-slot="progress"
		role="progressbar"
		aria-valuemin="0"
		aria-valuemax="100"
		aria-valuenow={value}
		class="relative h-1 w-full overflow-hidden rounded-full bg-primary/20"
		{ attrs... }
	>
		<div
			data-slot="progress-indicator"
			class="h-full w-full flex-1 bg-primary transition-all"
			style={ "transform: translateX(-" + gsx.RawCSS(remaining) + "%)" }
		></div>
	</div>
}

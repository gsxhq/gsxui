package ui

import (
	"strconv"

	"github.com/gsxhq/gsx"
)

// Slider is the shadcn/ui Slider (registry/new-york-v4/ui/slider.tsx),
// collapsed onto ONE native <input type="range"> — the four-part Radix
// tree (Root/Track/Range/Thumb) contributes nothing a real range input
// doesn't already carry for free: role="slider" plus
// aria-valuemin/-valuenow/-valuemax are implicit on <input type=range>,
// and the browser's own keyboard model (Left/Right/Down/Up adjust by
// step, PageUp/PageDown by a bigger jump, Home/End to min/max) already
// matches Radix's own traced behavior contract closely enough that no
// keyboard JS is needed at all — WIN, see docs/jsx-parity.md `## slider`.
//
// GAP (single scalar value only): a native range input is always ONE
// thumb, ONE value. Radix's multi-thumb/range-slider support (two or more
// thumbs, e.g. a min/max price range), orientation="vertical", inverted,
// and minStepsBetweenThumbs have no analog on this ADAPT and are not
// ported — see docs/jsx-parity.md.
//
// ADAPT (filled-range gradient, not accent-color): the Radix Range part
// (the primary-colored portion from min to the current value) has no free
// native equivalent. This port paints it as a linear-gradient() on
// ::-webkit-slider-runnable-track/::-moz-range-track (assets/gsxui.css),
// driven by a --fill custom property carrying the fill percentage — not
// accent-color, whose cross-browser track-fill behavior is unverified
// (see the 2026-07-24 controls source map's own open-decision writeup).
// value/min/max are float64s here, so the INITIAL --fill (first paint,
// zero JS) is exact-arithmetic Go server-side; ui/slider.js takes over
// from there, updating --fill on every `input` event while the user
// drags or keys the thumb.
//
// Retargeted to nova density (2026-07-24 controls source map, `## slider`):
// track h-1 (nova, was new-york-v4's h-1.5), thumb size-3 (nova, was
// size-4) with hover:ring-3/focus-visible:ring-3/active:ring-3 (nova adds
// active as a wholly new state). nova's own border-ring thumb recolor is
// NOT adopted — house policy keeps border-primary (color scope, not a
// metric, see docs/jsx-parity.md). nova's after:-inset-2 hit-target
// compensation for its own size-4->size-3 shrink is categorically
// unreachable here — a pseudo-element cannot itself carry a further
// pseudo-element (::-webkit-slider-thumb::after is not a legal selector)
// — ledgered as a GAP rather than silently dropped; see
// docs/jsx-parity.md.
//
// ADAPT (max/step zero-value defaulting): max and step fall back to 100/1
// when left at the Go zero value (0), matching shadcn's own slider-demo
// (max={100} step={1}); min's zero value (0) needs no such fallback since
// it is already Radix's own default. This is the same unset-vs-explicit-
// zero ambiguity every other zero-value-defaulted param in this codebase
// accepts (e.g. Toggle's variant/size via `|> default(...)`) — a caller
// who genuinely wants max=0 or step=0 cannot express it; not reachable in
// practice for a slider (a max of 0 with min<0 is an unusual range, and
// step=0 is invalid on a native range input regardless).
component Slider(value float64, min float64, max float64, step float64, attrs gsx.Attrs) {
	{{
		maxV := max
		if maxV == 0 {
			maxV = 100
		}
		stepV := step
		if stepV == 0 {
			stepV = 1
		}
		fillPct := 0.0
		if maxV > min {
			fillPct = (value - min) / (maxV - min) * 100
		}
		fill := strconv.FormatFloat(fillPct, 'f', -1, 64)
	}}
	<input
		type="range"
		data-slot="slider"
		data-gsxui-slider
		min={min}
		max={maxV}
		step={stepV}
		value={value}
		style=css`--fill: @{fill}%`
		class="appearance-none bg-transparent w-full cursor-pointer outline-none disabled:cursor-not-allowed disabled:opacity-50"
		{ attrs... }
	/>
}

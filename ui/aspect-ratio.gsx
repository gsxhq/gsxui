package ui

import "github.com/gsxhq/gsx"

// AspectRatio is the shadcn/ui AspectRatio. shadcn's version is a bare
// passthrough onto Radix's AspectRatioPrimitive.Root, which renders two
// nested divs: an outer one sized by the padding-bottom-percentage hack
// (position:relative; width:100%; padding-bottom:calc(100% / ratio), ratio
// set via inline style from a numeric `ratio` prop) and an inner one
// (position:absolute; inset:0) holding the actual children. This port
// replaces that whole mechanism with the CSS `aspect-ratio` property
// (ADAPT — see docs/jsx-parity.md) directly on a single div: no padding
// hack, no wrapper-within-wrapper, and no numeric-ratio-to-percentage
// arithmetic to reproduce. ratio is a string, not a float, so callers write
// the same expression the CSS property itself accepts, e.g. ratio="16 / 9"
// (aspect-ratio also accepts a bare number, ratio="1.5").
//
// gsx.RawCSS(ratio) opts the composed value out of gw's CSS value filter
// (MECHANISM): that filter is a conservative, character-blocklist port of
// html/template's CSS sanitizer and rejects "/" outright (also "(", ")",
// ";", and several others that never appear in a valid aspect-ratio value
// either) — but "/" is not incidental punctuation here, it is
// aspect-ratio's own required <width> "/" <height> separator, so no value
// most callers would ever write could pass the filter unmodified. ratio is
// trusted, developer-authored layout intent, the same trust boundary every
// ui/*.gsx component already extends to its own class strings (Tailwind's
// class attribute receives no injection filtering either) — not sanitized
// end-user request data.
component AspectRatio(ratio string, children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="aspect-ratio" style={ "aspect-ratio: " + gsx.RawCSS(ratio) } { attrs... }>
		{ children }
	</div>
}

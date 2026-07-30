package maia

import (
	"strconv"
	"strings"

	"github.com/gsxhq/gsx"
)

// aspectRatioForm is which shape of the CSS `aspect-ratio` grammar a caller's
// ratio string was recognised as. AspectRatio branches on it to pick a
// declaration whose separator and keyword are static template text, leaving
// only author-supplied numbers in interpolation holes.
type aspectRatioForm int

const (
	// aspectRatioOpaque is a ratio this component does not recognise. It is
	// emitted as one whole hole and left to gsx's CSS value filter to judge.
	aspectRatioOpaque aspectRatioForm = iota
	// aspectRatioPair is `<number> / <number>`.
	aspectRatioPair
	// aspectRatioAutoPair is `auto || <number> / <number>`.
	aspectRatioAutoPair
)

// parseAspectRatio recognises the CSS aspect-ratio grammar, `auto || <ratio>`,
// where <ratio> is `<number> [ / <number> ]?`. It returns the form and, for
// the two pair forms, the width and height numbers. A bare `<number>` and the
// bare `auto` keyword need no splitting — they carry no character the CSS
// value filter rejects — so they are left opaque and rendered as-is.
func parseAspectRatio(ratio string) (aspectRatioForm, string, string) {
	spec, auto := strings.TrimSpace(ratio), false
	// `auto || <ratio>` is order-free, so auto is accepted on either end.
	if rest, ok := strings.CutPrefix(spec, "auto "); ok {
		spec, auto = strings.TrimSpace(rest), true
	} else if rest, ok := strings.CutSuffix(spec, " auto"); ok {
		spec, auto = strings.TrimSpace(rest), true
	}
	w, h, ok := strings.Cut(spec, "/")
	w, h = strings.TrimSpace(w), strings.TrimSpace(h)
	if !ok || !aspectRatioNumber(w) || !aspectRatioNumber(h) {
		return aspectRatioOpaque, "", ""
	}
	if auto {
		return aspectRatioAutoPair, w, h
	}
	return aspectRatioPair, w, h
}

// aspectRatioNumber reports whether s is a CSS <number>. ParseFloat is the
// numeric authority; the alphabet check rejects the spellings Go accepts and
// CSS does not (inf, nan, 0x1p-2, digit separators).
func aspectRatioNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil && strings.Trim(s, "+-.0123456789eE") == ""
}

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
// (aspect-ratio also accepts a bare number, ratio="1.5", and the keyword
// forms "auto" and "auto 16 / 9").
//
// Branching on the parsed form (MECHANISM) is what replaces gsx.RawCSS here.
// gw's CSS value filter — a port of html/template's — rejects "/" in an
// interpolated value, and "/" is aspect-ratio's own required <width> "/"
// <height> separator, so the whole declaration used to be marked raw to get
// one character through. Each arm below writes the separator and the auto
// keyword as static template text and leaves each number in its own hole,
// where the filter still judges it: no hole can contribute a "/", so none
// can open a CSS comment or start a second declaration, and nothing is
// trusted that parseAspectRatio did not first recognise as a number. An
// opaque ratio is emitted unparsed and the filter alone decides — safe
// values render, hostile ones become the inert ZgotmplZ.
component AspectRatio(ratio string, children gsx.Node, attrs gsx.Attrs) {
	{{ form, w, h := parseAspectRatio(ratio) }}
	<div
		{ if form == aspectRatioPair {
			style=css`aspect-ratio: @{w} / @{h}`
		} else if form == aspectRatioAutoPair {
			style=css`aspect-ratio: auto @{w} / @{h}`
		} else {
			style=css`aspect-ratio: @{ratio}`
		} }
		class={ "block" }
		{ attrs... }
		data-gsxui-slot-aspect-ratio
	>
		{ children }
	</div>
}

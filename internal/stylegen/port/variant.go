package port

import (
	"regexp"
	"strings"
)

// Kind is what Classify decided a single @apply token means for the
// upstream→recipe transformation.
type Kind int

const (
	KindPlain     Kind = iota // travels verbatim onto its own rule
	KindDimension             // bare data-[dim=val]: — candidate for Rule 2
	KindSlot                  // **:data-[slot=x]: or *:data-[slot=x]: — Rule 3
)

func (k Kind) String() string {
	switch k {
	case KindPlain:
		return "Plain"
	case KindDimension:
		return "Dimension"
	case KindSlot:
		return "Slot"
	default:
		return "Kind(?)"
	}
}

// Classified is what Classify learned about one @apply token.
type Classified struct {
	Kind      Kind
	Dimension string // KindDimension only
	Value     string // KindDimension only
	Slot      string // KindSlot only
	Child     bool   // KindSlot only: true for *: (direct child), false for **:
	Rest      string // the utility with the matched prefix stripped
	Raw       string // the original token
}

// dimensionPrefix matches a whole leading variant segment that is exactly a
// bare `data-[<dim>=<value>]` bracket, with no group-/in-/has-/not-
// qualifier and no other prefix glued to it (e.g. a trailing `/<name>` group
// marker, which belongs to the `group-data-*` family and must stay
// KindPlain).
var dimensionPrefix = regexp.MustCompile(`^data-\[([^=\]]+)=([^\]]+)\]$`)

// slotPrefix matches a whole leading `data-[slot=<x>]` variant segment, as
// found immediately after a `*` or `**` child-combinator segment.
var slotPrefix = regexp.MustCompile(`^data-\[slot=([^\]]+)\]$`)

// Classify decides what a single @apply token means for the
// upstream→recipe transformation, per the normative rules:
//
//   - A leading, unqualified `data-[<dim>=<value>]:` (no group-/in-/has-/
//     not- prefix, and not itself preceded by another variant) is
//     KindDimension. Whether <dim> is a real shape dimension is decided
//     later, by the shape — not here.
//   - `**:data-[slot=<x>]:` is KindSlot with Child false (any descendant);
//     `*:data-[slot=<x>]:` is KindSlot with Child true (direct child).
//   - Everything else is KindPlain and travels verbatim in Rest — bare
//     boolean data variants (data-open:, data-closed:, ...), group-data-*/
//     group-has-data-*/in-data-*/has-data-*, aria-*, breakpoint/dark/
//     supports-*/not-*, and *:[...] arbitrary-selector children whose
//     bracket target isn't `data-[slot=x]`.
//   - A token may stack variants (data-[size=default]:sm:max-w-md); only a
//     leading data-[x=y]: counts as KindDimension, and the remainder stays
//     in Rest untouched (not re-classified).
func Classify(token string) Classified {
	segments := splitVariants(token)

	if len(segments) >= 2 && (segments[0] == "*" || segments[0] == "**") {
		if m := slotPrefix.FindStringSubmatch(segments[1]); m != nil {
			return Classified{
				Kind:  KindSlot,
				Slot:  m[1],
				Child: segments[0] == "*",
				Rest:  strings.Join(segments[2:], ":"),
				Raw:   token,
			}
		}
	}

	if len(segments) >= 2 {
		if m := dimensionPrefix.FindStringSubmatch(segments[0]); m != nil {
			return Classified{
				Kind:      KindDimension,
				Dimension: m[1],
				Value:     m[2],
				Rest:      strings.Join(segments[1:], ":"),
				Raw:       token,
			}
		}
	}

	return Classified{Kind: KindPlain, Rest: token, Raw: token}
}

// splitVariants splits a Tailwind utility token into its stacked variant
// segments (and, last, the base utility) on ':', tracking `[...]`/`(...)`
// nesting depth so a colon inside an arbitrary selector's own nested
// pseudo-class call (e.g. the `:` in `[svg:not([class*='size-'])]`) is
// never treated as a separator. This is the same forward-scan,
// depth-tracked technique the repo already depends on transitively via
// github.com/jackielii/tailwind-merge-go's class-name parser (used by
// merge.Merge, see internal/recipe's CheckConflicts): split on ':' wherever
// bracket depth and paren depth are both zero.
//
// A colon immediately after a bare `*`/`**` combinator (as in
// `*:[svg:not(...)]:size-8`) is still a split point — splitVariants does
// not fuse the combinator with the bracket that follows it into one
// segment. Classify is what treats a `*`/`**` segment together with an
// immediately-following `data-[slot=x]` segment as one compound prefix; for
// any other bracket target (like the `[svg:not(...)]` arbitrary selector)
// the segments simply don't match that compound pattern and the whole
// token falls through to KindPlain with Rest equal to the untouched Raw
// token, not a reassembly of the split segments.
func splitVariants(token string) []string {
	var segments []string
	bracketDepth := 0
	parenDepth := 0
	start := 0

	for i := 0; i < len(token); i++ {
		switch token[i] {
		case '[':
			bracketDepth++
		case ']':
			bracketDepth--
		case '(':
			parenDepth++
		case ')':
			parenDepth--
		case ':':
			if bracketDepth == 0 && parenDepth == 0 {
				segments = append(segments, token[start:i])
				start = i + 1
			}
		}
	}
	segments = append(segments, token[start:])
	return segments
}

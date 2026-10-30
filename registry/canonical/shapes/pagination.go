package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Pagination is a plain compound container, no Radix primitive and no
// dimension: nav/ul/li/a/span parts, all styled purely by their own marker.
//
// There is no "link" slot: PaginationLink renders through <ui.Button>, which
// brings Button's own recipe, as upstream's PaginationLink renders through
// Button. previous/next carry upstream's per-style inline padding
// (cn-pagination-previous/-next, pl-1.5 or pl-2). Upstream forces it with
// `!` to beat Button's `has-[>svg]:px-2`; gsxui ships no !important, so each
// style restates the padding under the same has-[>svg] variant by hand (the
// port strips the `!` and would drop the extra token). A caller overrides it
// in that form, `has-[>svg]:pl-12`, which merge.Merge lets win
// (jstest/specs/style-visual.spec.ts's "pagination-previous-caller").
var Pagination = recipe.Shape{
	Component: "pagination",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "content", Base: true},
		{Name: "previous", Base: true},
		{Name: "previous-label", Base: true},
		{Name: "next", Base: true},
		{Name: "next-label", Base: true},
		{Name: "ellipsis", Base: true},
	},
}

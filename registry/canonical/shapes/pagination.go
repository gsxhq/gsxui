package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Pagination is a plain compound container, no Radix primitive and no
// dimension: nav/ul/li/a/span parts, all styled purely by their own marker.
//
// There is deliberately NO "link" slot here. PaginationLink's only rules in
// default.css are the previous/next same-element compound rules
// (`:is([data-gsxui-slot-pagination-link][data-gsxui-slot-pagination-previous])
// { pl-1.5 }`, mirrored for -next/pr-1.5) — and those rules must stay
// authored as plain CSS inside @layer components, NOT migrate onto this
// recipe. See assets/css/styles/default/pagination.css's own comment for
// why: caller-overridability for this pair depends on the CSS LAYER
// boundary (@layer components always loses to @layer utilities regardless
// of specificity), which only exists pre-resolution. Once a relational rule
// resolves to a literal utility class scanned by Tailwind, it competes on
// ordinary specificity instead — and an arbitrary attribute-compound
// selector like `[&[data-gsxui-slot-pagination-previous]]:pl-1.5` outranks
// a caller's plain `pl-12` (0,2,0) vs (0,1,0)) regardless of source order,
// so the caller permanently loses the override guarantee
// jstest/specs/style-visual.spec.ts's "pagination-previous-caller" fixture
// pins. Nothing else in this component's own recipe sets padding-left/
// padding-right on PaginationLink in any state, so the layer gate has
// nothing to contest the retained rule against.
var Pagination = recipe.Shape{
	Component: "pagination",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "content", Base: true},
		{Name: "previous-label", Base: true},
		{Name: "next-label", Base: true},
		{Name: "ellipsis", Base: true},
	},
}

package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Empty is a plain compound container, no Radix primitive underneath —
// every part a styled <div>, the same shape as Card/Alert. Its "icon" slot
// carries the one dimension (data-variant: default|icon), matching both
// default.css and internal/stylecontract's Empty axis exactly — no
// disagreement to resolve here. The "description" slot's two relational
// rules key on a plain child <a> (not a marker, not a consumer class), so
// they translate as arbitrary child/child-state variants
// (`[&>a]:`, `[&>a:hover]:`) on description's own rule — see
// registry/styles/nova/empty.css.
var Empty = recipe.Shape{
	Component: "empty",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "header", Base: true},
		{Name: "icon", Base: true, Dimensions: []recipe.Dimension{
			{Name: "variant", Default: "default", Values: []string{"default", "icon"}},
		}},
		{Name: "title", Base: true},
		{Name: "description", Base: true},
		{Name: "content", Base: true},
	},
}

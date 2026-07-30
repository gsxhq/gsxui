package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Popover is the shadcn/ui native-popover machinery: light-dismissed
// auto-popover content anchored to a trigger. No dimensions — upstream and
// this catalogue's own default.css agree there is no variant/size axis.
//
// popover-trigger carries no rule in default.css at all (it's a plain
// <button> with no styling of its own), so it gets no accessor — nothing to
// migrate, the same "select-bridge/select-group" shape Select's own shape
// doc comment already established.
//
// popover-content's transition rules (opacity/scale/:popover-open/
// @starting-style/data-side translate arms) came from the file this
// component shared with HoverCard and Tooltip
// (assets/css/styles/default/overlay-popover-transition.css, now deleted —
// all three sharing components migrated in the same wave). See
// registry/styles/nova/popover.css and registry/styles/nova/select.css's own
// doc comment for the identical :popover-open/@starting-style translation
// this shape reuses verbatim.
var Popover = recipe.Shape{
	Component: "popover",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "content", Base: true},
	},
}

package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Tooltip is the shadcn/ui manual-popover machinery: hover/focus-driven
// content, always placed top (see ui/tooltip.gsx's doc comment: "its arrow
// is static because placement is always top"). No dimensions.
//
// tooltip-trigger carries no rule in default.css (a plain <button>), so it
// gets no accessor — same "select-bridge" shape Select's shape already
// established.
//
// tooltip-content has a same-shape relational rule to Card's proven
// descendant-marker :has() form: `:has([data-gsxui-slot-kbd])` reacts to a
// Kbd rendered inside it, not to a consumer-supplied class, so it translates
// directly as `has-[[data-gsxui-slot-kbd]]:pr-1.5` (see registry/styles/
// nova/tooltip.css) — no §10b caller-override hazard, since nothing a caller
// supplies is being detected.
//
// tooltip-content's transition rules were shared with Popover and HoverCard
// via assets/css/styles/default/overlay-popover-transition.css, now deleted
// now that all three sharing components are migrated together.
//
// tooltip-arrow is a plain base rule, no relational shape at all.
var Tooltip = recipe.Shape{
	Component: "tooltip",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "content", Base: true},
		{Name: "arrow", Base: true},
	},
}

package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// HoverCard is the shadcn/ui manual-popover machinery: pointer/focus-driven
// content, not light dismissal (see ui/hover-card.gsx's doc comment). No
// dimensions.
//
// hover-card-trigger carries no rule in default.css (a plain <span>), so it
// gets no accessor — same "select-bridge" shape Select's shape already
// established. "hover-card" is hyphenated: the registry/CSS identity stays
// "hover-card" (gsxui-recipe-hover-card, data-gsxui-slot-hover-card), and
// componentIdentifier derives the Go receiver "hoverCard" the same way
// "input-group" already does (see the slot-axis design doc §3.3).
//
// hover-card-content's transition rules were shared with Popover and Tooltip
// via assets/css/styles/default/overlay-popover-transition.css, now deleted
// now that all three sharing components are migrated together. See
// registry/styles/nova/hover-card.css.
var HoverCard = recipe.Shape{
	Component: "hover-card",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "content", Base: true},
	},
}

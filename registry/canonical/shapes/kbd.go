package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Kbd is a compound component: a root slot (Kbd itself) plus a sibling
// "group" slot (KbdGroup) that wraps multiple Kbds — mirroring Card's
// plain-sibling-slots shape, not a parent/child nesting. Neither slot has a
// dimension. The root slot's default.css block also carries a relational
// rule that reacts to an ANCESTOR tooltip-content, not to anything of Kbd's
// own — see registry/styles/nova/kbd.css for the three arbitrary variants
// that rule folds into (including one recovered from a stray match far from
// Kbd's main default.css block, near Tooltip's own rules).
var Kbd = recipe.Shape{
	Component: "kbd",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "group", Base: true},
	},
}

package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// ButtonGroup is a layout composite: it styles itself, a text part, and a
// composed Separator. Most of its presentation reaches into whatever it
// contains, so the shape is deliberately small and several relational rules
// stay behind in assets/css/ — see registry/styles/nova/button-group.css and
// assets/css/styles/default/button-group.css for the split and the reasons.
//
// "orientation" is a real dimension (data-orientation, stamped by
// button-group.gsx and declared by internal/stylecontract's ButtonGroup axis).
// Only "vertical" has a root-only rule in default.css; "horizontal" restates
// the base rule's own flex-row as an explicit no-op rather than inventing
// behaviour, the same treatment ToggleGroup's size arms already use. Keeping
// it a dimension (rather than a data-[orientation=vertical]: arbitrary
// variant) matters for § 10b: a dimension resolves to a single plain utility
// at generation time, so a caller's own flex-direction utility still wins on
// the ordinary merge, where a two-class variant selector would not.
//
// The separator slot has NO orientation dimension even though the contract
// declares the axis: default.css only ever keys button-group-separator on
// "vertical", and a "horizontal" arm would have no content of its own. It is
// carried as a data-[orientation=vertical]: arbitrary variant inside the base
// rule instead — which is also literally what upstream shadcn writes
// (registry/new-york-v4/ui/button-group.tsx's ButtonGroupSeparator className).
var ButtonGroup = recipe.Shape{
	Component: "button-group",
	Slots: []recipe.Slot{
		{Name: "", Base: true, Dimensions: []recipe.Dimension{
			{Name: "orientation", Default: "horizontal", Values: []string{"horizontal", "vertical"}},
		}},
		{Name: "text", Base: true},
		{Name: "separator", Base: true},
	},
}

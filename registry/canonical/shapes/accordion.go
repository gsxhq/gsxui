package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Accordion is a pure container: five styled slots, no dimensions. The root
// slot itself carries no rule in default.css (the root div is purely
// structural, data-name plus attrs), so it is not declared here at all —
// checkShapeCoverage rejects a declared slot the component never renders a
// class for, and Accordion's root never will. It sits on the native grouped
// <details name>/<summary> disclosure mechanism (no JS), so its [open]-state
// rotate on the trigger icon is an ancestor-attribute relational rule, and
// the trigger's ::-webkit-details-marker hiding has no Tailwind form
// (escape hatch — see registry/styles/nova/accordion.css's doc comment).
var Accordion = recipe.Shape{
	Component: "accordion",
	Slots: []recipe.Slot{
		{Name: "item", Base: true},
		{Name: "trigger", Base: true},
		{Name: "trigger-icon", Base: true},
		{Name: "content", Base: true},
		{Name: "content-inner", Base: true},
	},
}

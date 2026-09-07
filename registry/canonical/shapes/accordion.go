package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Accordion is a pure container: six styled slots, no dimensions. It sits on
// the native grouped <details name>/<summary> disclosure mechanism (no JS),
// so its [open]-state rotate on the trigger icon is an ancestor-attribute
// relational rule, and the trigger's ::-webkit-details-marker hiding has no
// Tailwind form (escape hatch — see registry/styles/nova/accordion.css's doc
// comment).
//
// The root slot IS declared, and was not through the 8-style port: the shape
// read the root div as purely structural (data-name plus attrs) because
// default.css — nova's presentation — carries no rule for it, which is true
// and stays true. It is not true of the other seven. Upstream's Accordion
// root carries `cn-accordion flex w-full flex-col`
// (apps/v4/registry/bases/radix/ui/accordion.tsx), and four of the eight
// style sheets give `.cn-accordion` real presentation on top of that
// style-invariant base: `overflow-hidden rounded-2xl border` in luma/maia/
// rhea, `overflow-hidden rounded-md border` in mira. Left undeclared, the
// whole outer frame those four styles draw around the item stack simply
// never rendered.
var Accordion = recipe.Shape{
	Component: "accordion",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "item", Base: true},
		{Name: "trigger", Base: true},
		{Name: "trigger-icon", Base: true},
		{Name: "content", Base: true},
		{Name: "content-inner", Base: true},
	},
}

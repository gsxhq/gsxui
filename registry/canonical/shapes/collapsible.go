package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Collapsible is a native <details>/<summary> disclosure. Only its trigger
// slot carries any rule in default.css (list-none, plus the ::-webkit-
// details-marker escape-hatch rule); root and content render with no styling
// of their own, so the shape declares only the trigger slot — a shape may
// declare less than the style contract (see agreement_test.go), and there is
// nothing for the other two slots to author.
var Collapsible = recipe.Shape{
	Component: "collapsible",
	Slots: []recipe.Slot{
		{Name: "trigger", Base: true},
	},
}

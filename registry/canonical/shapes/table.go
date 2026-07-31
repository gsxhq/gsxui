package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Table is a compound container: six styled slots, no dimensions.
// table-header and table-body are NOT declared here — neither carries any
// rule of its own in default.css. table-header's only mention
// (":where([data-gsxui-slot-table-header]) tr { border-b }") is a bare-tag
// descendant selector with no marker of its own; every <tr> under it is
// rendered by TableRow, which already applies border-b unconditionally as
// its OWN base rule, so the header rule is a verified no-op duplicate, not
// a rule this migration drops. table-body has no rule targeting itself at
// all — only ":where([data-gsxui-slot-table-body]) tr:last-child", which is
// TableRow reacting to a table-body ancestor (see table-row's Dimensions-
// free relational utilities in registry/styles/nova/table.css). Declaring
// either as a Base slot here would fail checkShapeCoverage (the component
// never renders a class for them) — same reasoning as Accordion's root.
var Table = recipe.Shape{
	Component: "table",
	Slots: []recipe.Slot{
		{Name: "container", Base: true},
		{Name: "", Base: true},
		{Name: "footer", Base: true},
		{Name: "row", Base: true},
		{Name: "head", Base: true},
		{Name: "cell", Base: true},
		{Name: "caption", Base: true},
	},
}

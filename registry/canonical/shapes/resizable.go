package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Resizable declares only the two slots default.css actually styles: the
// handle and its grip. resizable-panel-group and resizable-panel carry no
// CSS of their own (their layout is server-rendered inline styles, not
// recipe utilities — see ui/resizable.gsx's package doc comment), so the
// shape omits them entirely, same carve-out Collapsible's root/content used.
//
// The handle's aria-orientation is runtime state the handle reacts to, not a
// recipe dimension: the style contract's own axis on resizable-handle is
// exactly the same shape as Badge/Alert's aria-invalid axis — "runtime state
// nobody varies a recipe dimension on" (see agreement_test.go's authority
// note). The two orientation-conditional rules (width/cursor swap, and the
// grip's rotate-90) are authored directly as aria-[orientation=...]:
// arbitrary-variant utilities on the handle's own base rule, not as a
// Dimension/Values pair.
var Resizable = recipe.Shape{
	Component: "resizable",
	Slots: []recipe.Slot{
		{Name: "handle", Base: true},
		{Name: "handle-grip", Base: true},
	},
}

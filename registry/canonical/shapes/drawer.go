package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Drawer is a side-anchored dialog. Its root renders <ui.Dialog> and its
// trigger and close are plain unstyled buttons, so none of the three has a
// rule of its own to migrate — the same carve-out Dialog's own shape
// documents.
//
// content is the whole story. DrawerContent stamps
// data-gsxui-slot-dialog-content on its own <dialog> WITHOUT ever calling
// Dialog's content accessor, because it must not inherit Dialog's
// plain-modal box (top-1/2 left-1/2 max-w-* rounded-xl p-4 …). Before this
// migration that left it depending on two marker-keyed fallbacks — the
// Dialog-content chrome block in assets/css/styles/default.css and
// assets/css/styles/default/drawer-sheet-shared.css — for chrome that has no
// class of its own. Both are (0,0,0) rules the layer gate had to exempt.
// This shape ends that: content's base rule carries the shared chrome
// verbatim as its own utilities. That is also what upstream does — shadcn's
// drawer.tsx writes the whole fixed/z-50/flex-col/side chrome onto
// DrawerContent's own className rather than inheriting any dialog class.
//
// side is the one author-facing Dimension. Its four values match
// internal/stylecontract's drawer-content data-side axis exactly, and the
// default (bottom) matches ui/drawer.gsx's own `direction |> default`.
// data-state and the native open attribute stay runtime state, folded onto
// the rules as data-[state=…]: and open: variants — same treatment as
// Dialog's content slot.
var Drawer = recipe.Shape{
	Component: "drawer",
	Slots: []recipe.Slot{
		{Name: "content", Base: true, Dimensions: []recipe.Dimension{
			{Name: "side", Default: "bottom", Values: []string{"bottom", "left", "right", "top"}},
		}},
		{Name: "handle", Base: true},
		{Name: "header", Base: true},
		{Name: "footer", Base: true},
		{Name: "title", Base: true},
		{Name: "description", Base: true},
	},
}

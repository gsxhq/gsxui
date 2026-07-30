package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Sheet is Drawer's sibling: the same side-anchored <dialog>, with a
// background surface instead of a popover one, a close button, and no
// handle. Its root renders <ui.Dialog> and its trigger and close are plain
// unstyled buttons, so none of the three gets an accessor; close-label is
// fully styled by the shared sr-only rule in assets/css/foundation.css's
// @layer base, the same carve-out Dialog's own shape documents.
//
// SheetContent, like DrawerContent, stamps data-gsxui-slot-dialog-content
// without calling Dialog's content accessor — it must not inherit Dialog's
// plain-modal box. This shape is what lets it stop depending on marker
// fallbacks for the chrome it has no class for: content's base rule carries
// the shared <dialog> chrome verbatim. With Drawer already migrated, this
// retires BOTH fallbacks outright — the Dialog-content block in
// assets/css/styles/default.css and
// assets/css/styles/default/drawer-sheet-shared.css. Upstream agrees:
// shadcn's sheet.tsx writes the whole fixed/z-50/flex-col/side chrome onto
// SheetContent's own className.
//
// side is the one author-facing Dimension. Its four values match
// internal/stylecontract's sheet-content data-side axis exactly, and the
// default (right) matches ui/sheet.gsx's own `side |> default`. data-state
// and the native open attribute stay runtime state.
var Sheet = recipe.Shape{
	Component: "sheet",
	Slots: []recipe.Slot{
		{Name: "content", Base: true, Dimensions: []recipe.Dimension{
			{Name: "side", Default: "right", Values: []string{"bottom", "left", "right", "top"}},
		}},
		{Name: "close-button", Base: true},
		{Name: "close-icon", Base: true},
		{Name: "header", Base: true},
		{Name: "footer", Base: true},
		{Name: "title", Base: true},
		{Name: "description", Base: true},
	},
}

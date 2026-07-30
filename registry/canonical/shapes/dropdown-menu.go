package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// DropdownMenu shares assets/css/styles/default/menu.css with ContextMenu and
// Menubar (33 rules split three ways) and shared the now-deleted
// assets/css/styles/default/menu-popover-transition.css's discrete-transition
// block with ContextMenu, Menubar, NavigationMenu AND Combobox, which was that
// file's last unmigrated selector and has since migrated too (the file is
// gone; see registry/styles/nova/combobox.css) — this shape only claims the
// slots menu.css/menu-popover-transition.css actually styled. DropdownMenuTrigger, DropdownMenuGroup and
// DropdownMenuRadioGroup carry no class at all in either source CSS file
// (callers style their own trigger button; the two group wrappers are pure
// a11y markup) so they are deliberately NOT declared here — same "no rule,
// no slot" call the design doc's Card precedent already established for
// wholly-unstyled marked elements.
//
// Component is "dropdown-menu", matching both the style contract's
// RegistryName (internal/stylecontract/contracts_menus.go, fixed in a
// separate prior commit to agree with its own root slot and with upstream
// dropdown-menu.tsx) and every data-gsxui-slot-* marker here. The registry
// name used to be the irregular "dropdown" (disagreeing with its own slot
// names AND upstream); that was corrected ahead of this migration since it
// is pre-existing tech debt, not something this migration introduced.
var DropdownMenu = recipe.Shape{
	Component: "dropdown-menu",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "content", Base: true},
		{
			Name: "item", Base: true,
			Dimensions: []recipe.Dimension{
				{Name: "variant", Default: "default", Values: []string{"default", "destructive"}},
			},
		},
		{Name: "checkbox-item", Base: true},
		{Name: "checkbox-item-indicator", Base: true},
		{Name: "radio-item", Base: true},
		{Name: "radio-item-indicator", Base: true},
		{Name: "label", Base: true},
		{Name: "separator", Base: true},
		{Name: "shortcut", Base: true},
		{Name: "sub", Base: true},
		{Name: "sub-trigger", Base: true},
		{Name: "sub-content", Base: true},
	},
}

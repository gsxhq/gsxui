package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// ContextMenu shares assets/css/styles/default/menu.css with DropdownMenu and
// Menubar, and shared the now-deleted
// assets/css/styles/default/menu-popover-transition.css with DropdownMenu,
// Menubar, NavigationMenu AND Combobox, whose own migration deleted that file
// — see DropdownMenu's own shape file for the split rationale. ContextMenuTrigger
// and ContextMenuGroup/ContextMenuRadioGroup carry no class in either source
// CSS file (the trigger is a caller-styled drop-zone area; the two group
// wrappers are pure a11y markup) and are deliberately not declared here.
var ContextMenu = recipe.Shape{
	Component: "context-menu",
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

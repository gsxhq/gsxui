package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Menubar shares assets/css/styles/default/menu.css with DropdownMenu and
// ContextMenu, and shared the now-deleted
// assets/css/styles/default/menu-popover-transition.css with DropdownMenu,
// ContextMenu, NavigationMenu AND Combobox, whose own migration deleted that
// file — see DropdownMenu's own shape file for the split rationale.
// MenubarGroup and MenubarRadioGroup carry no class in either source CSS
// file (pure a11y markup, same call as DropdownMenuGroup/
// DropdownMenuRadioGroup) and are deliberately not declared here. Unlike
// DropdownMenu/ContextMenu, MenubarTrigger and MenubarMenu ARE styled (the
// bar itself is a real, visible control, not a caller-styled button).
var Menubar = recipe.Shape{
	Component: "menubar",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "menu", Base: true},
		{Name: "trigger", Base: true},
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

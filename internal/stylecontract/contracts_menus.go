package stylecontract

var menuOpenSideAxes = []Axis{
	{Attribute: "data-state", Values: []string{"closed", "open"}, RuntimeValues: []string{"open"}},
	{Attribute: "data-side", Values: []string{"bottom"}},
}

var menuSubOpenSideAxes = []Axis{
	{Attribute: "data-state", Values: []string{"closed", "open"}, RuntimeValues: []string{"open"}},
	{Attribute: "data-side", Values: []string{"right"}},
}

var menuCheckedAxes = []Axis{
	{Attribute: "data-state", Values: []string{"unchecked", "checked"}},
	{Attribute: "data-disabled"},
	{Attribute: "aria-checked", Values: []string{"false", "true"}},
}

var menuItemAxes = []Axis{
	{Attribute: "data-variant", Values: []string{"default", "destructive"}},
	{Attribute: "data-inset"},
	{Attribute: "data-disabled"},
}

var menuSubTriggerAxes = []Axis{
	{Attribute: "data-state", Values: []string{"closed", "open"}, RuntimeValues: []string{"open"}},
	{Attribute: "aria-expanded", Values: []string{"false", "true"}, RuntimeValues: []string{"true"}},
	{Attribute: "data-inset"},
	{Attribute: "data-disabled"},
}

var menuContracts = []Component{
	{
		Name:         "Command",
		RegistryName: "command",
		Slots: []Slot{
			{Name: "command"},
			{Name: "command-dialog"},
			{Name: "command-dialog-content"},
			{Name: "command-dialog-header"},
			{Name: "command-dialog-command"},
			{Name: "command-input-wrapper"},
			{Name: "command-input", Axes: []Axis{{Attribute: "disabled"}}},
			{Name: "command-list"},
			{Name: "command-empty"},
			{Name: "command-group"},
			{Name: "command-group-heading"},
			{Name: "command-separator"},
			{Name: "command-item", Axes: []Axis{
				{Attribute: "data-selected", Values: []string{"true"}, RuntimeValues: []string{"true"}},
				{Attribute: "data-disabled", Values: []string{"true"}},
				{Attribute: "aria-selected", Values: []string{"false", "true"}, RuntimeValues: []string{"true"}},
			}},
			{Name: "command-shortcut"},
		},
	},
	{
		Name:         "Combobox",
		RegistryName: "combobox",
		Slots: []Slot{
			{Name: "combobox"},
			{Name: "combobox-bridge"},
			{Name: "combobox-input-group"},
			{Name: "combobox-input", Axes: []Axis{
				{Attribute: "disabled"},
				{Attribute: "aria-expanded", Values: []string{"false", "true"}, RuntimeValues: []string{"true"}},
			}},
			{Name: "combobox-trigger", Axes: []Axis{{Attribute: "disabled"}}},
			{Name: "combobox-trigger-icon"},
			{Name: "combobox-clear", Axes: []Axis{{Attribute: "disabled"}}},
			{Name: "combobox-content", Axes: []Axis{
				{Attribute: "data-state", Values: []string{"closed", "open"}, RuntimeValues: []string{"open"}},
				{Attribute: "data-side", Values: []string{"bottom"}},
				{Attribute: "data-empty"},
			}},
			{Name: "combobox-list", Axes: []Axis{{Attribute: "data-empty"}}},
			{Name: "combobox-item", Axes: []Axis{
				{Attribute: "data-state", Values: []string{"unchecked", "checked"}},
				{Attribute: "data-highlighted", Values: []string{"true"}, RuntimeValues: []string{"true"}},
				{Attribute: "data-disabled"},
				{Attribute: "aria-selected", Values: []string{"false", "true"}, RuntimeValues: []string{"true"}},
			}},
			{Name: "combobox-item-indicator"},
			{Name: "combobox-group"},
			{Name: "combobox-label"},
			{Name: "combobox-empty"},
			{Name: "combobox-separator"},
			{Name: "combobox-value"},
		},
	},
	{
		Name:         "DropdownMenu",
		RegistryName: "dropdown",
		Slots: []Slot{
			{Name: "dropdown-menu"},
			{Name: "dropdown-menu-trigger", Axes: []Axis{
				{Attribute: "aria-expanded", Values: []string{"false", "true"}, RuntimeValues: []string{"true"}},
				{Attribute: "disabled"},
			}},
			{Name: "dropdown-menu-content", Axes: menuOpenSideAxes},
			{Name: "dropdown-menu-group"},
			{Name: "dropdown-menu-item", Axes: menuItemAxes},
			{Name: "dropdown-menu-checkbox-item", Axes: menuCheckedAxes},
			{Name: "dropdown-menu-checkbox-item-indicator"},
			{Name: "dropdown-menu-radio-group"},
			{Name: "dropdown-menu-radio-item", Axes: menuCheckedAxes},
			{Name: "dropdown-menu-radio-item-indicator"},
			{Name: "dropdown-menu-label", Axes: []Axis{{Attribute: "data-inset"}}},
			{Name: "dropdown-menu-separator"},
			{Name: "dropdown-menu-shortcut"},
			{Name: "dropdown-menu-sub"},
			{Name: "dropdown-menu-sub-trigger", Axes: menuSubTriggerAxes},
			{Name: "dropdown-menu-sub-content", Axes: menuSubOpenSideAxes},
		},
	},
	{
		Name:         "ContextMenu",
		RegistryName: "context-menu",
		Slots: []Slot{
			{Name: "context-menu"},
			{Name: "context-menu-trigger"},
			{Name: "context-menu-content", Axes: []Axis{
				{Attribute: "data-state", Values: []string{"closed", "open"}, RuntimeValues: []string{"open"}},
			}},
			{Name: "context-menu-group"},
			{Name: "context-menu-item", Axes: menuItemAxes},
			{Name: "context-menu-checkbox-item", Axes: menuCheckedAxes},
			{Name: "context-menu-checkbox-item-indicator"},
			{Name: "context-menu-radio-group"},
			{Name: "context-menu-radio-item", Axes: menuCheckedAxes},
			{Name: "context-menu-radio-item-indicator"},
			{Name: "context-menu-label", Axes: []Axis{{Attribute: "data-inset"}}},
			{Name: "context-menu-separator"},
			{Name: "context-menu-shortcut"},
			{Name: "context-menu-sub"},
			{Name: "context-menu-sub-trigger", Axes: menuSubTriggerAxes},
			{Name: "context-menu-sub-content", Axes: menuSubOpenSideAxes},
		},
	},
	{
		Name:         "Menubar",
		RegistryName: "menubar",
		Slots: []Slot{
			{Name: "menubar"},
			{Name: "menubar-menu"},
			{Name: "menubar-trigger", Axes: []Axis{
				{Attribute: "data-state", Values: []string{"closed", "open"}, RuntimeValues: []string{"open"}},
				{Attribute: "aria-expanded", Values: []string{"false", "true"}, RuntimeValues: []string{"true"}},
				{Attribute: "disabled"},
			}},
			{Name: "menubar-content", Axes: menuOpenSideAxes},
			{Name: "menubar-group"},
			{Name: "menubar-item", Axes: menuItemAxes},
			{Name: "menubar-checkbox-item", Axes: menuCheckedAxes},
			{Name: "menubar-checkbox-item-indicator"},
			{Name: "menubar-radio-group"},
			{Name: "menubar-radio-item", Axes: menuCheckedAxes},
			{Name: "menubar-radio-item-indicator"},
			{Name: "menubar-label", Axes: []Axis{{Attribute: "data-inset"}}},
			{Name: "menubar-separator"},
			{Name: "menubar-shortcut"},
			{Name: "menubar-sub"},
			{Name: "menubar-sub-trigger", Axes: menuSubTriggerAxes},
			{Name: "menubar-sub-content", Axes: menuSubOpenSideAxes},
		},
	},
	{
		Name:         "NavigationMenu",
		RegistryName: "navigation-menu",
		Slots: []Slot{
			{Name: "navigation-menu", Axes: []Axis{{Attribute: "data-viewport", Values: []string{"false"}}}},
			{Name: "navigation-menu-list"},
			{Name: "navigation-menu-item"},
			{Name: "navigation-menu-trigger", Axes: []Axis{
				{Attribute: "data-state", Values: []string{"closed", "open"}, RuntimeValues: []string{"open"}},
				{Attribute: "aria-expanded", Values: []string{"false", "true"}, RuntimeValues: []string{"true"}},
				{Attribute: "disabled"},
			}},
			{Name: "navigation-menu-trigger-icon"},
			{Name: "navigation-menu-content", Axes: menuOpenSideAxes},
			{Name: "navigation-menu-link", Axes: []Axis{
				{Attribute: "data-variant", Values: []string{"default", "trigger"}},
				{Attribute: "data-active", Values: []string{"false", "true"}},
			}},
			{Name: "navigation-menu-indicator", Axes: []Axis{{Attribute: "data-state", Values: []string{"hidden", "visible"}, RuntimeValues: []string{"hidden", "visible"}}}},
			{Name: "navigation-menu-indicator-arrow"},
		},
	},
}

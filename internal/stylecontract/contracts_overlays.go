package stylecontract

var overlayContracts = []Component{
	{
		Name:         "Accordion",
		RegistryName: "accordion",
		Slots: []Slot{
			{Name: "accordion"},
			{Name: "accordion-item", Axes: []Axis{{Attribute: "open"}}},
			{Name: "accordion-trigger"},
			{Name: "accordion-trigger-icon"},
			{Name: "accordion-content"},
			{Name: "accordion-content-inner"},
		},
	},
	{
		Name:         "AlertDialog",
		RegistryName: "alert-dialog",
		Slots: []Slot{
			{Name: "alert-dialog"},
			{Name: "alert-dialog-trigger", Axes: []Axis{
				{Attribute: "aria-haspopup", Values: []string{"dialog"}},
				{Attribute: "aria-expanded", Values: []string{"false", "true"}, RuntimeValues: []string{"true"}},
			}},
			{Name: "alert-dialog-content", Axes: []Axis{
				{Attribute: "role", Values: []string{"alertdialog"}},
				{Attribute: "data-state", Values: []string{"closed", "open"}, RuntimeValues: []string{"open"}},
				{Attribute: "open"},
			}},
			{Name: "alert-dialog-header"},
			{Name: "alert-dialog-footer"},
			{Name: "alert-dialog-title"},
			{Name: "alert-dialog-description"},
			{Name: "alert-dialog-action", Axes: []Axis{
				{Attribute: "data-variant", Values: []string{"default"}},
				{Attribute: "data-size", Values: []string{"default"}},
			}},
			{Name: "alert-dialog-cancel", Axes: []Axis{
				{Attribute: "data-variant", Values: []string{"outline"}},
				{Attribute: "data-size", Values: []string{"default"}},
			}},
		},
	},
	{
		Name:         "Collapsible",
		RegistryName: "collapsible",
		Slots: []Slot{
			{Name: "collapsible", Axes: []Axis{{Attribute: "open"}}},
			{Name: "collapsible-trigger"},
			{Name: "collapsible-content"},
		},
	},
	{
		Name:         "Dialog",
		RegistryName: "dialog",
		Slots: []Slot{
			{Name: "dialog"},
			{Name: "dialog-trigger", Axes: []Axis{
				{Attribute: "aria-haspopup", Values: []string{"dialog"}},
				{Attribute: "aria-expanded", Values: []string{"false", "true"}, RuntimeValues: []string{"true"}},
			}},
			{Name: "dialog-content", Axes: []Axis{
				{Attribute: "data-state", Values: []string{"closed", "open"}, RuntimeValues: []string{"open"}},
				{Attribute: "open"},
			}},
			{Name: "dialog-close"},
			{Name: "dialog-close-button"},
			{Name: "dialog-close-icon"},
			{Name: "dialog-close-label"},
			{Name: "dialog-header"},
			{Name: "dialog-footer"},
			{Name: "dialog-footer-close"},
			{Name: "dialog-title"},
			{Name: "dialog-description"},
		},
	},
	{
		Name:         "Drawer",
		RegistryName: "drawer",
		Slots: []Slot{
			{Name: "drawer"},
			{Name: "drawer-trigger", Axes: []Axis{
				{Attribute: "aria-haspopup", Values: []string{"dialog"}},
				{Attribute: "aria-expanded", Values: []string{"false", "true"}, RuntimeValues: []string{"true"}},
			}},
			{Name: "drawer-content", Axes: []Axis{
				{Attribute: "data-state", Values: []string{"closed", "open"}, RuntimeValues: []string{"open"}},
				{Attribute: "data-side", Values: []string{"bottom", "left", "right", "top"}},
				{Attribute: "open"},
			}},
			{Name: "drawer-handle"},
			{Name: "drawer-header"},
			{Name: "drawer-footer"},
			{Name: "drawer-title"},
			{Name: "drawer-description"},
			{Name: "drawer-close"},
		},
	},
	{
		Name:         "HoverCard",
		RegistryName: "hover-card",
		Slots: []Slot{
			{Name: "hover-card"},
			{Name: "hover-card-trigger"},
			{Name: "hover-card-content", Axes: []Axis{
				{Attribute: "popover", Values: []string{"manual"}},
				{Attribute: "data-state", Values: []string{"closed", "open"}, RuntimeValues: []string{"open"}},
				{Attribute: "data-side", Values: []string{"bottom"}},
			}},
		},
	},
	{
		Name:         "Popover",
		RegistryName: "popover",
		Slots: []Slot{
			{Name: "popover"},
			{Name: "popover-trigger", Axes: []Axis{{Attribute: "aria-expanded", Values: []string{"false", "true"}, RuntimeValues: []string{"true"}}}},
			{Name: "popover-content", Axes: []Axis{
				{Attribute: "popover", Values: []string{"auto"}},
				{Attribute: "data-state", Values: []string{"closed", "open"}, RuntimeValues: []string{"open"}},
				{Attribute: "data-side", Values: []string{"bottom"}},
			}},
		},
	},
	{
		Name:         "Sheet",
		RegistryName: "sheet",
		Slots: []Slot{
			{Name: "sheet"},
			{Name: "sheet-trigger", Axes: []Axis{
				{Attribute: "aria-haspopup", Values: []string{"dialog"}},
				{Attribute: "aria-expanded", Values: []string{"false", "true"}, RuntimeValues: []string{"true"}},
			}},
			{Name: "sheet-content", Axes: []Axis{
				{Attribute: "data-state", Values: []string{"closed", "open"}, RuntimeValues: []string{"open"}},
				{Attribute: "data-side", Values: []string{"bottom", "left", "right", "top"}},
				{Attribute: "open"},
			}},
			{Name: "sheet-close"},
			{Name: "sheet-close-button"},
			{Name: "sheet-close-icon"},
			{Name: "sheet-close-label"},
			{Name: "sheet-header"},
			{Name: "sheet-footer"},
			{Name: "sheet-title"},
			{Name: "sheet-description"},
		},
	},
	{
		Name:         "Tooltip",
		RegistryName: "tooltip",
		Slots: []Slot{
			{Name: "tooltip"},
			{Name: "tooltip-trigger"},
			{Name: "tooltip-content", Axes: []Axis{
				{Attribute: "popover", Values: []string{"manual"}},
				{Attribute: "role", Values: []string{"tooltip"}},
				{Attribute: "data-state", Values: []string{"closed", "open"}, RuntimeValues: []string{"open"}},
				{Attribute: "data-side", Values: []string{"top"}},
			}},
			{Name: "tooltip-arrow"},
		},
	},
}

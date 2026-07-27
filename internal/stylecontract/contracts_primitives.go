package stylecontract

var primitiveContracts = []Component{
	{
		Name: "Alert",
		Slots: []Slot{
			{
				Name: "alert",
				Axes: []Axis{{
					Attribute: "data-variant",
					Values:    []string{"default", "destructive"},
				}},
			},
			{Name: "alert-title"},
			{Name: "alert-description"},
		},
	},
	{
		Name:  "AspectRatio",
		Slots: []Slot{{Name: "aspect-ratio"}},
	},
	{
		Name: "Avatar",
		Slots: []Slot{
			{Name: "avatar"},
			{Name: "avatar-image"},
			{Name: "avatar-fallback"},
		},
	},
	{
		Name: "Badge",
		Slots: []Slot{{
			Name: "badge",
			Axes: []Axis{
				{
					Attribute: "data-variant",
					Values:    []string{"default", "secondary", "destructive", "outline", "ghost", "link"},
				},
				{
					Attribute: "aria-invalid",
					Values:    []string{"true"},
				},
			},
		}},
	},
	{
		Name: "Breadcrumb",
		Slots: []Slot{
			{Name: "breadcrumb"},
			{Name: "breadcrumb-list"},
			{Name: "breadcrumb-item"},
			{Name: "breadcrumb-link"},
			{
				Name: "breadcrumb-page",
				Axes: []Axis{
					{Attribute: "aria-disabled", Values: []string{"true"}},
					{Attribute: "aria-current", Values: []string{"page"}},
				},
			},
			{Name: "breadcrumb-separator"},
			{Name: "breadcrumb-ellipsis"},
			{Name: "breadcrumb-ellipsis-label"},
		},
	},
	{
		Name: "Button",
		Slots: []Slot{{
			Name: "button",
			Axes: []Axis{
				{
					Attribute: "data-variant",
					Values:    []string{"default", "destructive", "outline", "secondary", "ghost", "link"},
				},
				{
					Attribute: "data-size",
					Values:    []string{"default", "xs", "sm", "lg", "icon", "icon-xs", "icon-sm", "icon-lg"},
				},
				{Attribute: "disabled"},
				{Attribute: "aria-invalid", Values: []string{"true"}},
			},
		}},
	},
	{
		Name: "ButtonGroup",
		Slots: []Slot{
			{
				Name: "button-group",
				Axes: []Axis{{
					Attribute: "data-orientation",
					Values:    []string{"horizontal", "vertical"},
				}},
			},
			{Name: "button-group-text"},
			{
				Name: "button-group-separator",
				Axes: []Axis{{
					Attribute: "data-orientation",
					Values:    []string{"horizontal", "vertical"},
				}},
			},
		},
	},
	{
		Name: "Card",
		Slots: []Slot{
			{Name: "card"},
			{Name: "card-header"},
			{Name: "card-title"},
			{Name: "card-description"},
			{Name: "card-action"},
			{Name: "card-content"},
			{Name: "card-footer"},
		},
	},
	{
		Name: "Empty",
		Slots: []Slot{
			{Name: "empty"},
			{Name: "empty-header"},
			{
				Name: "empty-icon",
				Axes: []Axis{{
					Attribute: "data-variant",
					Values:    []string{"default", "icon"},
				}},
			},
			{Name: "empty-title"},
			{Name: "empty-description"},
			{Name: "empty-content"},
		},
	},
	{
		Name:  "Icon",
		Slots: []Slot{{Name: "icon"}},
	},
	{
		Name: "Item",
		Slots: []Slot{
			{Name: "item-group"},
			{
				Name: "item-separator",
				Axes: []Axis{{
					Attribute: "data-orientation",
					Values:    []string{"horizontal", "vertical"},
				}},
			},
			{
				Name: "item",
				Axes: []Axis{
					{
						Attribute: "data-variant",
						Values:    []string{"default", "outline", "muted"},
					},
					{
						Attribute: "data-size",
						Values:    []string{"default", "sm"},
					},
				},
			},
			{
				Name: "item-media",
				Axes: []Axis{{
					Attribute: "data-variant",
					Values:    []string{"default", "icon", "image"},
				}},
			},
			{Name: "item-content"},
			{Name: "item-title"},
			{Name: "item-description"},
			{Name: "item-actions"},
			{Name: "item-header"},
			{Name: "item-footer"},
		},
	},
	{
		Name: "Kbd",
		Slots: []Slot{
			{Name: "kbd"},
			{Name: "kbd-group"},
		},
	},
	{
		Name:  "Label",
		Slots: []Slot{{Name: "label"}},
	},
	{
		Name: "Pagination",
		Slots: []Slot{
			{Name: "pagination"},
			{Name: "pagination-content"},
			{Name: "pagination-item"},
			{
				Name: "pagination-link",
				Axes: []Axis{
					{
						Attribute: "data-active",
						Values:    []string{"false", "true"},
					},
					{
						Attribute: "data-variant",
						Values:    []string{"ghost", "outline"},
					},
					{
						Attribute: "data-size",
						Values:    []string{"default", "xs", "sm", "lg", "icon", "icon-xs", "icon-sm", "icon-lg"},
					},
					{
						Attribute: "aria-current",
						Values:    []string{"page"},
					},
				},
			},
			{Name: "pagination-previous"},
			{Name: "pagination-previous-label"},
			{Name: "pagination-next"},
			{Name: "pagination-next-label"},
			{Name: "pagination-ellipsis"},
			{Name: "pagination-ellipsis-label"},
		},
	},
	{
		Name: "Separator",
		Slots: []Slot{{
			Name: "separator",
			Axes: []Axis{{
				Attribute: "data-orientation",
				Values:    []string{"horizontal", "vertical"},
			}},
		}},
	},
	{
		Name:  "Skeleton",
		Slots: []Slot{{Name: "skeleton"}},
	},
	{
		Name:  "Spinner",
		Slots: []Slot{{Name: "spinner"}},
	},
	{
		Name: "Table",
		Slots: []Slot{
			{Name: "table-container"},
			{Name: "table"},
			{Name: "table-header"},
			{Name: "table-body"},
			{Name: "table-footer"},
			{
				Name: "table-row",
				Axes: []Axis{
					{Attribute: "data-state", Values: []string{"selected"}},
				},
			},
			{Name: "table-head"},
			{Name: "table-cell"},
			{Name: "table-caption"},
		},
	},
}

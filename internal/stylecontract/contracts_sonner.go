package stylecontract

var sonnerContracts = []Component{
	{
		Name: "Sonner",
		Slots: []Slot{
			{Name: "toaster", Axes: []Axis{{Attribute: "data-expanded", Values: []string{"false", "true"}}}},
			{Name: "toast", Axes: []Axis{
				{Attribute: "data-type", Values: []string{"default", "success", "info", "warning", "error", "loading"}},
				{Attribute: "data-state", Values: []string{"closed", "open"}},
				{Attribute: "data-visible", Values: []string{"false", "true"}},
			}},
			{Name: "toast-icon"},
			{Name: "toast-content"},
			{Name: "toast-title"},
			{Name: "toast-description"},
			{Name: "toast-action"},
			{Name: "toast-cancel"},
			{Name: "toast-close"},
			{Name: "toast-close-icon"},
		},
	},
}

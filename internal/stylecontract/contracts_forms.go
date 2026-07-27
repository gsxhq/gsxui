package stylecontract

var formContracts = []Component{
	{
		Name: "Checkbox",
		Slots: []Slot{{
			Name: "checkbox",
			Axes: []Axis{
				{Attribute: "checked"},
				{Attribute: "disabled"},
				{Attribute: "aria-invalid", Values: []string{"true"}},
			},
		}},
	},
	{
		Name: "Field",
		Slots: []Slot{
			{Name: "field-set"},
			{Name: "field-legend", Axes: []Axis{{Attribute: "data-variant", Values: []string{"legend", "label"}}}},
			{Name: "field-group"},
			{Name: "field", Axes: []Axis{
				{Attribute: "data-orientation", Values: []string{"vertical", "horizontal", "responsive"}},
				{Attribute: "data-invalid", Values: []string{"true"}},
			}},
			{Name: "field-content"},
			{Name: "field-label"},
			{Name: "field-title"},
			{Name: "field-description"},
			{Name: "field-separator-wrapper", Axes: []Axis{{Attribute: "data-content", Values: []string{"false", "true"}}}},
			{Name: "field-separator"},
			{Name: "field-separator-content"},
			{Name: "field-error"},
		},
	},
	{
		Name: "Input",
		Slots: []Slot{{
			Name: "input",
			Axes: []Axis{
				{Attribute: "disabled"},
				{Attribute: "aria-invalid", Values: []string{"true"}},
			},
		}},
	},
	{
		Name: "InputGroup",
		Slots: []Slot{
			{Name: "input-group", Axes: []Axis{{Attribute: "data-disabled", Values: []string{"true"}}}},
			{Name: "input-group-addon", Axes: []Axis{{Attribute: "data-align", Values: []string{"inline-start", "inline-end", "block-start", "block-end"}}}},
			{Name: "input-group-button", Axes: []Axis{
				{Attribute: "data-variant", Values: []string{"default", "destructive", "outline", "secondary", "ghost", "link"}},
				{Attribute: "data-size", Values: []string{"xs", "sm", "icon-xs", "icon-sm"}},
			}},
			{Name: "input-group-text"},
			{Name: "input-group-control", Axes: []Axis{{Attribute: "aria-invalid", Values: []string{"true"}}}},
		},
	},
	{
		Name: "InputOTP",
		Slots: []Slot{
			{Name: "input-otp"},
			{Name: "input-otp-input", Axes: []Axis{{Attribute: "disabled"}}},
			{Name: "input-otp-group"},
			{Name: "input-otp-slot", Axes: []Axis{
				{Attribute: "data-active", Values: []string{"false", "true"}},
				{Attribute: "aria-invalid", Values: []string{"true"}},
			}},
			{Name: "input-otp-separator"},
			{Name: "input-otp-caret-overlay"},
			{Name: "input-otp-caret"},
		},
	},
	{
		Name: "NativeSelect",
		Slots: []Slot{
			{Name: "native-select-wrapper"},
			{Name: "native-select", Axes: []Axis{
				{Attribute: "disabled"},
				{Attribute: "aria-invalid", Values: []string{"true"}},
			}},
		},
	},
	{
		Name: "Progress",
		Slots: []Slot{
			{Name: "progress"},
			{Name: "progress-indicator"},
		},
	},
	{
		Name: "Radio",
		Slots: []Slot{{
			Name: "radio",
			Axes: []Axis{
				{Attribute: "checked"},
				{Attribute: "disabled"},
				{Attribute: "aria-invalid", Values: []string{"true"}},
			},
		}},
	},
	{
		Name: "Select",
		Slots: []Slot{
			{Name: "select"},
			{Name: "select-bridge"},
			{Name: "select-trigger", Axes: []Axis{
				{Attribute: "data-size", Values: []string{"default", "sm"}},
				{Attribute: "data-state", Values: []string{"closed", "open"}},
				{Attribute: "data-placeholder"},
				{Attribute: "disabled"},
				{Attribute: "aria-invalid", Values: []string{"true"}},
			}},
			{Name: "select-value"},
			{Name: "select-content", Axes: []Axis{
				{Attribute: "data-state", Values: []string{"closed", "open"}},
				{Attribute: "data-side", Values: []string{"bottom", "left", "right", "top"}},
			}},
			{Name: "select-group"},
			{Name: "select-label"},
			{Name: "select-item", Axes: []Axis{
				{Attribute: "data-state", Values: []string{"unchecked", "checked"}},
				{Attribute: "data-disabled", Values: []string{"true"}},
				{Attribute: "aria-selected", Values: []string{"false", "true"}},
			}},
			{Name: "select-item-indicator"},
			{Name: "select-item-text"},
			{Name: "select-separator"},
		},
	},
	{
		Name:  "Slider",
		Slots: []Slot{{Name: "slider", Axes: []Axis{{Attribute: "disabled"}}}},
	},
	{
		Name: "Switch",
		Slots: []Slot{{
			Name: "switch",
			Axes: []Axis{
				{Attribute: "checked"},
				{Attribute: "disabled"},
			},
		}},
	},
	{
		Name: "Textarea",
		Slots: []Slot{{
			Name: "textarea",
			Axes: []Axis{
				{Attribute: "disabled"},
				{Attribute: "aria-invalid", Values: []string{"true"}},
			},
		}},
	},
	{
		Name: "Toggle",
		Slots: []Slot{{
			Name: "toggle",
			Axes: []Axis{
				{Attribute: "data-variant", Values: []string{"default", "outline"}},
				{Attribute: "data-size", Values: []string{"default", "sm", "lg"}},
				{Attribute: "data-state", Values: []string{"off", "on"}},
				{Attribute: "disabled"},
				{Attribute: "aria-invalid", Values: []string{"true"}},
			},
		}},
	},
	{
		Name: "ToggleGroup",
		Slots: []Slot{
			{Name: "toggle-group", Axes: []Axis{
				{Attribute: "data-variant", Values: []string{"default", "outline"}},
				{Attribute: "data-size", Values: []string{"default", "sm", "lg"}},
				{Attribute: "data-spacing", Values: []string{"0", "2"}},
				{Attribute: "data-orientation", Values: []string{"horizontal"}},
				{Attribute: "disabled"},
			}},
			{Name: "toggle-group-item", Axes: []Axis{
				{Attribute: "data-variant", Values: []string{"default", "outline"}},
				{Attribute: "data-size", Values: []string{"default", "sm", "lg"}},
				{Attribute: "data-spacing", Values: []string{"0", "2"}},
				{Attribute: "data-orientation", Values: []string{"horizontal"}},
				{Attribute: "data-state", Values: []string{"off", "on"}},
				{Attribute: "disabled"},
				{Attribute: "aria-invalid", Values: []string{"true"}},
			}},
		},
	},
}

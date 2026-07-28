package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestFormControlsExposeCSSOnlySlots(t *testing.T) {
	tests := []struct {
		name  string
		node  gsx.Node
		slots []string
	}{
		{"Checkbox", ui.Checkbox(nil), []string{"checkbox"}},
		{"FieldSet", ui.FieldSet(gsx.Raw("x"), nil), []string{"field-set"}},
		{"FieldLegend", ui.FieldLegend("", gsx.Raw("x"), nil), []string{"field-legend"}},
		{"FieldGroup", ui.FieldGroup(gsx.Raw("x"), nil), []string{"field-group"}},
		{"Field", ui.Field("", gsx.Raw("x"), nil), []string{"field"}},
		{"FieldContent", ui.FieldContent(gsx.Raw("x"), nil), []string{"field-content"}},
		{"FieldLabel", ui.FieldLabel(gsx.Raw("x"), nil), []string{"label field-label"}},
		{"FieldTitle", ui.FieldTitle(gsx.Raw("x"), nil), []string{"field-title"}},
		{"FieldDescription", ui.FieldDescription(gsx.Raw("x"), nil), []string{"field-description"}},
		{"FieldSeparator", ui.FieldSeparator(gsx.Raw("Or"), nil), []string{"field-separator-wrapper", "separator field-separator", "field-separator-content"}},
		{"FieldError", ui.FieldError(gsx.Raw("x"), nil), []string{"field-error"}},
		{"Input", ui.Input(nil), []string{"input"}},
		{"InputGroup", ui.InputGroup(gsx.Raw("x"), nil), []string{"input-group"}},
		{"InputGroupAddon", ui.InputGroupAddon("", gsx.Raw("x"), nil), []string{"input-group-addon"}},
		{"InputGroupButton", ui.InputGroupButton("", "", gsx.Raw("x"), nil), []string{"button input-group-button"}},
		{"InputGroupText", ui.InputGroupText(gsx.Raw("x"), nil), []string{"input-group-text"}},
		{"InputGroupInput", ui.InputGroupInput(nil), []string{"input input-group-control"}},
		{"InputGroupTextarea", ui.InputGroupTextarea("x", nil), []string{"textarea input-group-control"}},
		{"InputOTP", ui.InputOTP(gsx.Raw("x"), nil), []string{"input-otp", "input-otp-input"}},
		{"InputOTPGroup", ui.InputOTPGroup(gsx.Raw("x"), nil), []string{"input-otp-group"}},
		{"InputOTPSlot", ui.InputOTPSlot(nil), []string{"input-otp-slot"}},
		{"InputOTPSeparator", ui.InputOTPSeparator(nil), []string{"input-otp-separator", "icon"}},
		{"NativeSelect", ui.NativeSelect(gsx.Raw("<option>x</option>"), nil), []string{"native-select-wrapper", "native-select", "icon"}},
		{"Progress", ui.Progress(25, nil), []string{"progress", "progress-indicator"}},
		{"Radio", ui.Radio(nil), []string{"radio"}},
		{"Select", ui.Select("fruit", false, false, "", gsx.Raw("x"), nil), []string{"select", "select-bridge"}},
		{"SelectTrigger", ui.SelectTrigger("", ui.SelectValue("Pick", nil), nil), []string{"select-trigger", "select-value", "icon"}},
		{"SelectContent", ui.SelectContent(gsx.Raw("x"), nil), []string{"select-content"}},
		{"SelectGroup", ui.SelectGroup(gsx.Raw("x"), nil), []string{"select-group"}},
		{"SelectLabel", ui.SelectLabel(gsx.Raw("x"), nil), []string{"select-label"}},
		{"SelectItem", ui.SelectItem("x", false, false, gsx.Raw("x"), nil), []string{"select-item", "select-item-indicator", "icon", "select-item-text"}},
		{"SelectSeparator", ui.SelectSeparator(nil), []string{"select-separator"}},
		{"Slider", ui.Slider(50, 0, 100, 1, nil), []string{"slider"}},
		{"Switch", ui.Switch(nil), []string{"switch"}},
		{"Textarea", ui.Textarea("", nil), []string{"textarea"}},
		{"Toggle", ui.Toggle(false, "", "", gsx.Raw("x"), nil), []string{"toggle"}},
		{"ToggleGroup", ui.ToggleGroup("multiple", "", "", "", gsx.Raw("x"), nil), []string{"toggle-group"}},
		{"ToggleGroupItem", ui.ToggleGroupItem("multiple", "", "", "", false, "x", gsx.Raw("x"), nil), []string{"toggle toggle-group-item"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := render(t, tt.node)
			for _, slot := range tt.slots {
				for name := range strings.FieldsSeq(slot) {
					if !strings.Contains(got, `data-gsxui-slot-`+name) {
						t.Errorf("missing slot marker %q\nin: %s", name, got)
					}
				}
			}
			if strings.Contains(got, `data-slot=`) {
				t.Errorf("legacy data-slot must not render\nin: %s", got)
			}
			if tt.name == "InputGroupButton" {
				if !strings.Contains(got, canonicalButtonClass("ghost", "default")) {
					t.Errorf("InputGroupButton lost exact canonical Button roles\nin: %s", got)
				}
			} else if strings.Contains(got, ` class="`) {
				t.Errorf("non-Button built-in presentation class must not render\nin: %s", got)
			}
		})
	}
}

func TestFormControlCompositionTokenOrder(t *testing.T) {
	tests := []struct {
		name string
		node gsx.Node
		want string
	}{
		{
			name: "InputGroupInput",
			node: ui.InputGroupInput(nil),
			want: `<input type="text" data-gsxui-slot-input-group-control data-gsxui-slot-input>`,
		},
		{
			name: "InputGroupTextarea",
			node: ui.InputGroupTextarea("hello", nil),
			want: `<textarea data-gsxui-slot-input-group-control data-gsxui-slot-textarea>hello</textarea>`,
		},
		{
			name: "FieldLabel",
			node: ui.FieldLabel(gsx.Raw("Email"), nil),
			want: `<label data-gsxui-slot-field-label data-gsxui-slot-label>Email</label>`,
		},
		{
			name: "FieldSeparator",
			node: ui.FieldSeparator(gsx.Raw("Or"), nil),
			want: `<div data-content="true" data-gsxui-slot-field-separator-wrapper><div role="none" data-orientation="horizontal" data-gsxui-slot-field-separator data-gsxui-slot-separator></div><span data-gsxui-slot-field-separator-content>Or</span></div>`,
		},
		{
			name: "ToggleGroupItem",
			node: ui.ToggleGroupItem("multiple", "", "", "", false, "bold", gsx.Raw("B"), nil),
			want: `<button type="button" data-gsxui-toggle-group-item data-variant="default" data-size="default" data-spacing="0" data-orientation="horizontal" data-state="off" data-value="bold" aria-pressed="false" data-gsxui-slot-toggle-group-item data-gsxui-slot-toggle>B</button>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := render(t, tt.node); got != tt.want {
				t.Errorf("composition render mismatch\n got: %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestFormControlBehaviorHooksAreSeparateFromSlots(t *testing.T) {
	otp := render(t, ui.InputOTP(ui.InputOTPSlot(nil), nil))
	for _, want := range []string{
		`data-gsxui-input-otp`,
		`data-gsxui-input-otp-input`,
		`data-gsxui-input-otp-slot`,
	} {
		if !strings.Contains(otp, want) {
			t.Errorf("InputOTP missing behavior hook %q\nin: %s", want, otp)
		}
	}

	selectHTML := render(t, ui.Select(
		"fruit", false, false, "",
		ui.SelectContent(
			ui.SelectGroup(
				gsx.Fragment(
					ui.SelectLabel(gsx.Raw("Fruit"), nil),
					ui.SelectItem("apple", false, false, gsx.Raw("Apple"), nil),
				),
				nil,
			),
			nil,
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-select`,
		`data-gsxui-select-bridge`,
		`data-gsxui-select-content`,
		`data-gsxui-select-group`,
		`data-gsxui-select-label`,
		`data-gsxui-select-item`,
		`data-gsxui-select-item-text`,
	} {
		if !strings.Contains(selectHTML, want) {
			t.Errorf("Select missing behavior hook %q\nin: %s", want, selectHTML)
		}
	}
}

func TestFormControlDynamicInlineValuesRemain(t *testing.T) {
	tests := []struct {
		name string
		node gsx.Node
		want string
	}{
		{"Progress", ui.Progress(25, nil), `style="transform: translateX(-75%)"`},
		{"Slider", ui.Slider(25, 0, 100, 1, nil), `style="--fill: 25%"`},
		{"ToggleGroup", ui.ToggleGroup("multiple", "", "", "2", gsx.Raw("x"), nil), `style="--gap: 2"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := render(t, tt.node)
			if !strings.Contains(got, tt.want) {
				t.Errorf("missing dynamic style %q\nin: %s", tt.want, got)
			}
		})
	}
}

func TestFormControlCallerClassIsForwardedOnce(t *testing.T) {
	tests := []struct {
		name string
		node gsx.Node
	}{
		{"Checkbox", ui.Checkbox(gsx.Attrs{{Key: "class", Value: "task-5-caller"}})},
		{"Field", ui.Field("", nil, gsx.Attrs{{Key: "class", Value: "task-5-caller"}})},
		{"Input", ui.Input(gsx.Attrs{{Key: "class", Value: "task-5-caller"}})},
		{"InputGroup", ui.InputGroup(nil, gsx.Attrs{{Key: "class", Value: "task-5-caller"}})},
		{"InputOTP", ui.InputOTP(nil, gsx.Attrs{{Key: "class", Value: "task-5-caller"}})},
		{"NativeSelect", ui.NativeSelect(nil, gsx.Attrs{{Key: "class", Value: "task-5-caller"}})},
		{"Progress", ui.Progress(0, gsx.Attrs{{Key: "class", Value: "task-5-caller"}})},
		{"Radio", ui.Radio(gsx.Attrs{{Key: "class", Value: "task-5-caller"}})},
		{"Select", ui.Select("", false, false, "", nil, gsx.Attrs{{Key: "class", Value: "task-5-caller"}})},
		{"Slider", ui.Slider(0, 0, 100, 1, gsx.Attrs{{Key: "class", Value: "task-5-caller"}})},
		{"Switch", ui.Switch(gsx.Attrs{{Key: "class", Value: "task-5-caller"}})},
		{"Textarea", ui.Textarea("", gsx.Attrs{{Key: "class", Value: "task-5-caller"}})},
		{"Toggle", ui.Toggle(false, "", "", nil, gsx.Attrs{{Key: "class", Value: "task-5-caller"}})},
		{"ToggleGroup", ui.ToggleGroup("multiple", "", "", "", nil, gsx.Attrs{{Key: "class", Value: "task-5-caller"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := render(t, tt.node)
			if strings.Count(got, `class="task-5-caller"`) != 1 {
				t.Errorf("caller class must be the only class and render once\nin: %s", got)
			}
		})
	}
}

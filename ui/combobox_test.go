package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestComboboxItemSelectedStampsIndicator(t *testing.T) {
	got := render(t, ui.ComboboxItem("next.js", true, gsx.Raw("Next.js"), nil))
	for _, want := range []string{
		`data-gsxui-slot-combobox-item`,
		`data-value="next.js"`,
		`aria-selected="true"`,
		`data-gsxui-slot-combobox-item-indicator`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestComboboxItemUnselected(t *testing.T) {
	got := render(t, ui.ComboboxItem("next.js", false, gsx.Raw("Next.js"), nil))
	if !strings.Contains(got, `aria-selected="false"`) {
		t.Errorf("want aria-selected=false\nin: %s", got)
	}
}

func TestComboboxRootRendersFormBridge(t *testing.T) {
	// Non-JS form posts must carry the value: a real named control is
	// server-rendered, mirroring ui/select.gsx's bridge.
	got := render(t, ui.Combobox("framework", "next.js", gsx.Raw("x"), nil))
	for _, want := range []string{
		`data-gsxui-slot-combobox`,
		`data-gsxui-slot-combobox-bridge`,
		`name="framework"`,
		`value="next.js"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestComboboxInputComposesInputGroup(t *testing.T) {
	got := render(t, ui.ComboboxInput("Search framework...", true, false, false, nil, nil))
	for _, want := range []string{
		`data-gsxui-combobox-input-group`,
		// InputGroupInput composes both styling tokens; the group keys its
		// focus ring off input-group-control.
		`data-gsxui-slot-combobox-input-group data-gsxui-slot-input-group`,
		`data-gsxui-slot-combobox-input data-gsxui-slot-input-group-control data-gsxui-slot-input`,
		`data-gsxui-slot-combobox-trigger data-gsxui-slot-input-group-button data-gsxui-slot-button`,
		`data-gsxui-slot-combobox-trigger-icon data-gsxui-slot-icon`,
		`role="combobox"`,
		`aria-expanded="false"`,
		// Permanent regardless of open/closed state (APG expectation,
		// review round 1 finding 8).
		`aria-haspopup="listbox"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
	if strings.Contains(got, `data-gsxui-slot-combobox-clear data-gsxui-slot-input-group-button data-gsxui-slot-button`) {
		t.Errorf("showClear=false must not render the clear button\nin: %s", got)
	}
}

func TestComboboxInputShowClear(t *testing.T) {
	got := render(t, ui.ComboboxInput("", false, true, false, nil, nil))
	if !strings.Contains(got, `data-gsxui-slot-combobox-clear data-gsxui-slot-input-group-button data-gsxui-slot-button`) {
		t.Errorf("want the clear button\nin: %s", got)
	}
	if strings.Contains(got, `data-gsxui-slot-combobox-trigger data-gsxui-slot-input-group-button data-gsxui-slot-button`) {
		t.Errorf("showTrigger=false must not render the trigger\nin: %s", got)
	}
}

func TestComboboxInputBothTriggerAndClearRenderBothButtons(t *testing.T) {
	// showTrigger and showClear are independent booleans — both true must
	// render both buttons (CSS alone, group-has-data-[slot=combobox-clear]/
	// input-group:hidden, hides the trigger visually; the reviewer's own
	// "TESTS" list asks this be pinned server-side too).
	got := render(t, ui.ComboboxInput("", true, true, false, nil, nil))
	for _, want := range []string{
		`data-gsxui-slot-combobox-trigger data-gsxui-slot-input-group-button data-gsxui-slot-button`,
		`data-gsxui-slot-combobox-clear data-gsxui-slot-input-group-button data-gsxui-slot-button`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestComboboxInputAttrsSplitClassToWrapperIdToControl(t *testing.T) {
	// The exact bug fixed pre-ship (and re-broken, then re-fixed, in review
	// round 1): id must land on the actual <input> control, class must
	// land on the InputGroup wrapper as its only class, matching source's own
	// {className, ...props} destructure (## combobox ledger entry).
	got := render(t, ui.ComboboxInput("", false, false, false, nil, gsx.Attrs{
		{Key: "id", Value: "combo-input"},
		{Key: "class", Value: "w-[300px]"},
	}))
	inputIdx := strings.Index(got, "<input")
	if inputIdx < 0 {
		t.Fatalf("no <input> in output: %s", got)
	}
	wrapper, control := got[:inputIdx], got[inputIdx:]
	if strings.Contains(wrapper, `id="combo-input"`) {
		t.Errorf("id must not land on the InputGroup wrapper\nwrapper: %s", wrapper)
	}
	if !strings.Contains(control, `id="combo-input"`) {
		t.Errorf("want id on the <input> control\ncontrol: %s", control)
	}
	if !strings.Contains(wrapper, "w-[300px]") {
		t.Errorf("want caller class merged onto the InputGroup wrapper\nwrapper: %s", wrapper)
	}
	if strings.Contains(control, "w-[300px]") {
		t.Errorf("caller class must not also land on the <input> control\ncontrol: %s", control)
	}
	if strings.Count(got, `class=`) != 1 {
		t.Errorf("caller class must be the only rendered class\nin: %s", got)
	}
}

func TestComboboxContentCarriesPopoverStateContract(t *testing.T) {
	got := render(t, ui.ComboboxContent(gsx.Raw("x"), nil))
	for _, want := range []string{
		`popover="auto"`,
		`data-gsxui-slot-combobox-content`,
		`data-state="closed"`,
		`data-side="bottom"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestComboboxEmptyIsHiddenUntilListIsEmpty(t *testing.T) {
	got := render(t, ui.ComboboxEmpty(gsx.Raw("No framework found."), nil))
	for _, want := range []string{
		`data-gsxui-slot-combobox-empty`,
		`data-gsxui-combobox-empty`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestComboboxListHasListboxRole(t *testing.T) {
	got := render(t, ui.ComboboxList(gsx.Raw("x"), nil))
	if !strings.Contains(got, `role="listbox"`) {
		t.Errorf("want role=listbox\nin: %s", got)
	}
	if !strings.Contains(got, `data-gsxui-slot-combobox-list`) {
		t.Errorf("want data-gsxui-slot-combobox-list\nin: %s", got)
	}
}

func TestComboboxBehaviorPartsHaveDedicatedHooks(t *testing.T) {
	got := render(t, ui.ComboboxGroup(
		gsx.Fragment(
			ui.ComboboxLabel(gsx.Raw("Framework"), nil),
			ui.ComboboxEmpty(gsx.Raw("Empty"), nil),
			ui.ComboboxSeparator(nil),
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-combobox-group`,
		`data-gsxui-combobox-label`,
		`data-gsxui-combobox-separator`,
		`data-gsxui-combobox-empty`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestComboboxNameEmptyOmitsFormBridge(t *testing.T) {
	got := render(t, ui.Combobox("", "", gsx.Raw("x"), nil))
	if strings.Contains(got, "data-gsxui-combobox-bridge") {
		t.Errorf("name=\"\" must not render the hidden form bridge\nin: %s", got)
	}
	if !strings.Contains(got, `data-gsxui-slot-combobox`) {
		t.Errorf("want data-gsxui-slot-combobox\nin: %s", got)
	}
}

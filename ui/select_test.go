package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

// TestSelectTriggerPinnedDefault covers the structural and behavioral
// attributes; presentation is owned by the stylesheet.
func TestSelectTriggerPinnedDefault(t *testing.T) {
	got := render(t, ui.SelectTrigger("", ui.SelectValue("Select a fruit", nil), nil))
	want := `<button type="button" data-gsxui-select-trigger role="combobox" aria-expanded="false" aria-autocomplete="none" data-state="closed" data-size="default" data-placeholder data-gsxui-slot="select-trigger"><span data-gsxui-select-value data-gsxui-slot="select-value">Select a fruit</span><svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" data-gsxui-slot="icon"><path d="m6 9 6 6 6-6"/></svg></button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestSelectTriggerSm pins the sm variant's distinguishing stamp — only
// data-size flips (the h-7 / sm-radius tokens are unconditional in the class
// string, keyed off the data-size value).
func TestSelectTriggerSm(t *testing.T) {
	got := render(t, ui.SelectTrigger("sm", ui.SelectValue("Pick", nil), nil))
	if !strings.Contains(got, `data-size="sm"`) {
		t.Errorf("sm trigger should stamp data-size=\"sm\"\nin: %s", got)
	}
}

// TestSelectTriggerDefaultSize proves the zero-value size stamps "default".
func TestSelectTriggerDefaultSize(t *testing.T) {
	got := render(t, ui.SelectTrigger("", gsx.Raw("x"), nil))
	if !strings.Contains(got, `data-size="default"`) {
		t.Errorf("zero-value size should stamp data-size=\"default\"\nin: %s", got)
	}
	if !strings.Contains(got, `data-placeholder`) {
		t.Errorf("trigger should server-render data-placeholder\nin: %s", got)
	}
}

func TestSelectValuePinned(t *testing.T) {
	got := render(t, ui.SelectValue("Select a fruit", nil))
	want := `<span data-gsxui-select-value data-gsxui-slot="select-value">Select a fruit</span>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestSelectContentPinned covers the listbox popover structure.
func TestSelectContentPinned(t *testing.T) {
	got := render(t, ui.SelectContent(gsx.Raw("x"), nil))
	want := `<div data-gsxui-select-content popover="auto" role="listbox" tabindex="-1" data-state="closed" data-side="bottom" data-gsxui-slot="select-content">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestSelectGroupPinned(t *testing.T) {
	got := render(t, ui.SelectGroup(gsx.Raw("x"), nil))
	want := `<div data-gsxui-select-group role="group" data-gsxui-slot="select-group">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestSelectLabelPinned(t *testing.T) {
	got := render(t, ui.SelectLabel(gsx.Raw("Fruits"), nil))
	want := `<div data-gsxui-select-label data-gsxui-slot="select-label">Fruits</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestSelectItemUnchecked pins an unselected, enabled item and its
// behavior-only text hook.
func TestSelectItemUnchecked(t *testing.T) {
	got := render(t, ui.SelectItem("apple", false, false, gsx.Raw("Apple"), nil))
	want := `<div data-gsxui-select-item role="option" data-value="apple" data-state="unchecked" aria-selected="false" tabindex="-1" data-gsxui-slot="select-item"><span data-gsxui-slot="select-item-indicator"><svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" data-gsxui-slot="icon"><path d="M20 6 9 17l-5-5"/></svg></span><span data-gsxui-select-item-text data-gsxui-slot="select-item-text">Apple</span></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestSelectItemChecked pins a server-selected item: only data-state flips to
// "checked" (which the CSS keys the checkmark visibility off), aria-selected
// stays "false" (not focused at server render).
func TestSelectItemChecked(t *testing.T) {
	got := render(t, ui.SelectItem("banana", true, false, gsx.Raw("Banana"), nil))
	want := `<div data-gsxui-select-item role="option" data-value="banana" data-state="checked" aria-selected="false" tabindex="-1" data-gsxui-slot="select-item"><span data-gsxui-slot="select-item-indicator"><svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" data-gsxui-slot="icon"><path d="M20 6 9 17l-5-5"/></svg></span><span data-gsxui-select-item-text data-gsxui-slot="select-item-text">Banana</span></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestSelectItemDisabled proves disabled stamps both the presence-selector
// data-disabled (matched by data-[disabled]:pointer-events-none/opacity-50)
// and aria-disabled="true"; enabled items carry neither.
func TestSelectItemDisabled(t *testing.T) {
	got := render(t, ui.SelectItem("cherry", false, true, gsx.Raw("Cherry"), nil))
	for _, want := range []string{`data-disabled="true"`, `aria-disabled="true"`, `data-state="unchecked"`} {
		if !strings.Contains(got, want) {
			t.Errorf("disabled item missing %q\nin: %s", want, got)
		}
	}
	enabled := render(t, ui.SelectItem("cherry", false, false, gsx.Raw("Cherry"), nil))
	if strings.Contains(enabled, "data-disabled") || strings.Contains(enabled, "aria-disabled") {
		t.Errorf("enabled item should carry neither data-disabled nor aria-disabled\nin: %s", enabled)
	}
}

func TestSelectSeparatorPinned(t *testing.T) {
	got := render(t, ui.SelectSeparator(nil))
	want := `<div aria-hidden="true" data-gsxui-slot="select-separator"></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestSelectHiddenBridgePresent proves that when name != "" the root renders
// the real hidden native <select> form bridge with the synthetic empty option
// (select.js fills the real options at init); required/disabled/form forward.
func TestSelectHiddenBridgePresent(t *testing.T) {
	got := render(t, ui.Select("fruit", true, false, "myform", gsx.Raw("K"), nil))
	for _, want := range []string{
		`data-gsxui-select-bridge`,
		`aria-hidden="true"`,
		`tabindex="-1"`,
		`data-gsxui-slot="select-bridge"`,
		`name="fruit"`,
		` required`,
		`form="myform"`,
		`<option value=""></option>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("bridge missing %q\nin: %s", want, got)
		}
	}
	if strings.Contains(got, ` disabled`) {
		t.Errorf("disabled=false should omit the bridge disabled attr\nin: %s", got)
	}
}

// TestSelectHiddenBridgeDisabled proves disabled forwards as a bare boolean attr.
func TestSelectHiddenBridgeDisabled(t *testing.T) {
	got := render(t, ui.Select("fruit", false, true, "", gsx.Raw("K"), nil))
	if !strings.Contains(got, ` disabled`) || strings.Contains(got, `disabled="`) {
		t.Errorf("disabled should render as a bare boolean attr\nin: %s", got)
	}
	// form="" is omitted entirely (empty form id would mis-associate).
	if strings.Contains(got, `form=`) {
		t.Errorf("empty form should be omitted\nin: %s", got)
	}
}

// TestSelectNoBridgeWithoutName proves the bridge is absent when name is empty
// (a display-only select carries no form control) — the honest no-JS GAP.
func TestSelectNoBridgeWithoutName(t *testing.T) {
	got := render(t, ui.Select("", false, false, "", gsx.Raw("K"), nil))
	if strings.Contains(got, "data-gsxui-select-bridge") {
		t.Errorf("no name should render no bridge\nin: %s", got)
	}
	if !strings.Contains(got, `<div data-gsxui-select data-gsxui-slot="select">K</div>`) {
		t.Errorf("root without a bridge should be just the wrapper\nin: %s", got)
	}
}

// TestSelectContentCallerClassMerges proves caller classes are forwarded.
func TestSelectContentCallerClassMerges(t *testing.T) {
	got := render(t, ui.SelectContent(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "min-w-72"}}))
	if strings.Count(got, `class="min-w-72"`) != 1 {
		t.Errorf("caller class must be forwarded exactly once\nin: %s", got)
	}
}

func TestSelectTriggerCallerClassMerges(t *testing.T) {
	got := render(t, ui.SelectTrigger("", gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "w-[180px]"}}))
	if strings.Count(got, `class="w-[180px]"`) != 1 {
		t.Errorf("caller class must be forwarded exactly once\nin: %s", got)
	}
}

func TestSelectAttrsFallThrough(t *testing.T) {
	got := render(t, ui.SelectTrigger("", gsx.Raw("x"), gsx.Attrs{{Key: "id", Value: "t1"}, {Key: "aria-label", Value: "Fruit"}}))
	for _, want := range []string{`id="t1"`, `aria-label="Fruit"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// TestSelectItemIconDependency proves the Check indicator is actually wired in
// (the select -> icon edge internal/registry derives), not merely imported.
func TestSelectItemIconDependency(t *testing.T) {
	got := render(t, ui.SelectItem("x", false, false, gsx.Raw("X"), nil))
	if !strings.Contains(got, `data-gsxui-slot="icon"`) || !strings.Contains(got, `d="M20 6 9 17l-5-5"`) {
		t.Errorf("expected the Check icon svg in the item render\nin: %s", got)
	}
}

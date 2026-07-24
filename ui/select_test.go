package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

// TestSelectTriggerPinnedDefault pins the default-size trigger token-by-token:
// the merged nova metrics (gap-1.5, rounded-lg, pr-2/pl-2.5 directional split,
// no shadow-xs, h-8 default / h-7 sm + the sm radius override,
// focus-visible:ring-[3px] house syntax), the combobox ARIA, the server-
// rendered data-placeholder (initial no-value state), and the internally-
// owned ChevronDown.
func TestSelectTriggerPinnedDefault(t *testing.T) {
	got := render(t, ui.SelectTrigger("", ui.SelectValue("Select a fruit", nil), nil))
	want := `<button type="button" data-slot="select-trigger" data-gsxui-select-trigger role="combobox" aria-expanded="false" aria-autocomplete="none" data-state="closed" data-size="default" data-placeholder class="flex w-fit items-center justify-between gap-1.5 rounded-lg border border-input bg-transparent pr-2 pl-2.5 py-2 text-sm whitespace-nowrap transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 data-[placeholder]:text-muted-foreground data-[size=default]:h-8 data-[size=sm]:h-7 data-[size=sm]:rounded-[min(var(--radius-md),10px)] *:data-[slot=select-value]:line-clamp-1 *:data-[slot=select-value]:flex *:data-[slot=select-value]:items-center *:data-[slot=select-value]:gap-1.5 dark:bg-input/30 dark:hover:bg-input/50 dark:aria-invalid:ring-destructive/40 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 [&amp;_svg:not([class*=&#39;text-&#39;])]:text-muted-foreground"><span data-slot="select-value" style="pointer-events: none">Select a fruit</span><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4 opacity-50"><path d="m6 9 6 6 6-6"/></svg></button>`
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
	for _, want := range []string{
		`data-[size=sm]:h-7`,
		`data-[size=sm]:rounded-[min(var(--radius-md),10px)]`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
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
	want := `<span data-slot="select-value" style="pointer-events: none">Select a fruit</span>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestSelectContentPinned pins the listbox popover: the merged nova content
// metrics (min-w-36, rounded-lg, border kept — no ring swap, shadow-md) plus
// the discrete-transition block ported byte-for-byte from DropdownMenuContent
// (replacing Radix's tw-animate keyframes) and the popover/role/data-state
// hooks.
func TestSelectContentPinned(t *testing.T) {
	got := render(t, ui.SelectContent(gsx.Raw("x"), nil))
	want := `<div data-slot="select-content" data-gsxui-select-content popover="auto" role="listbox" tabindex="-1" data-state="closed" data-side="bottom" class="z-50 max-h-96 min-w-36 origin-top-left overflow-x-hidden overflow-y-auto rounded-lg border bg-popover p-1 text-popover-foreground shadow-md opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:opacity-0 starting:open:scale-95 data-[side=bottom]:starting:open:-translate-y-2 data-[side=left]:starting:open:translate-x-2 data-[side=right]:starting:open:-translate-x-2 data-[side=top]:starting:open:translate-y-2">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestSelectGroupPinned(t *testing.T) {
	got := render(t, ui.SelectGroup(gsx.Raw("x"), nil))
	want := `<div data-slot="select-group" data-gsxui-select-group role="group">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestSelectLabelPinned(t *testing.T) {
	got := render(t, ui.SelectLabel(gsx.Raw("Fruits"), nil))
	want := `<div data-slot="select-label" class="px-1.5 py-1 text-xs text-muted-foreground">Fruits</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestSelectItemUnchecked pins an unselected, enabled item: data-state=
// "unchecked" (the value-tracking attribute driving the checkmark's CSS
// gating), aria-selected="false" (recomputed to (isValue AND isFocused) by
// select.js — always false at server render, nothing is focused), no
// data-disabled, and the merged nova item metrics. The check indicator is
// server-rendered but CSS-gated hidden until
// group-data-[state=checked]/select-item:flex.
func TestSelectItemUnchecked(t *testing.T) {
	got := render(t, ui.SelectItem("apple", false, false, gsx.Raw("Apple"), nil))
	want := `<div data-slot="select-item" data-gsxui-select-item role="option" data-value="apple" data-state="unchecked" aria-selected="false" tabindex="-1" class="group/select-item relative flex w-full cursor-default items-center gap-1.5 rounded-md py-1 pr-8 pl-1.5 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 [&amp;_svg:not([class*=&#39;text-&#39;])]:text-muted-foreground *:[span]:last:flex *:[span]:last:items-center *:[span]:last:gap-2"><span data-slot="select-item-indicator" class="pointer-events-none absolute right-2 hidden size-4 items-center justify-center group-data-[state=checked]/select-item:flex"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4"><path d="M20 6 9 17l-5-5"/></svg></span><span data-slot="select-item-text">Apple</span></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestSelectItemChecked pins a server-selected item: only data-state flips to
// "checked" (which the CSS keys the checkmark visibility off), aria-selected
// stays "false" (not focused at server render).
func TestSelectItemChecked(t *testing.T) {
	got := render(t, ui.SelectItem("banana", true, false, gsx.Raw("Banana"), nil))
	want := `<div data-slot="select-item" data-gsxui-select-item role="option" data-value="banana" data-state="checked" aria-selected="false" tabindex="-1" class="group/select-item relative flex w-full cursor-default items-center gap-1.5 rounded-md py-1 pr-8 pl-1.5 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 [&amp;_svg:not([class*=&#39;text-&#39;])]:text-muted-foreground *:[span]:last:flex *:[span]:last:items-center *:[span]:last:gap-2"><span data-slot="select-item-indicator" class="pointer-events-none absolute right-2 hidden size-4 items-center justify-center group-data-[state=checked]/select-item:flex"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4"><path d="M20 6 9 17l-5-5"/></svg></span><span data-slot="select-item-text">Banana</span></div>`
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
	want := `<div data-slot="select-separator" aria-hidden="true" class="bg-border -mx-1 my-1 h-px"></div>`
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
		`class="sr-only"`,
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
	if !strings.Contains(got, `<div data-slot="select" data-gsxui-select class="contents">K</div>`) {
		t.Errorf("root without a bridge should be just the wrapper\nin: %s", got)
	}
}

// TestSelectContentCallerClassMerges proves a caller override wins over the
// base via tailwind-merge while structural base classes survive.
func TestSelectContentCallerClassMerges(t *testing.T) {
	got := render(t, ui.SelectContent(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "min-w-72"}}))
	if strings.Contains(got, "min-w-36") {
		t.Errorf("base min-w-36 should be dropped by caller min-w-72\nin: %s", got)
	}
	for _, want := range []string{"min-w-72", "rounded-lg", "bg-popover"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestSelectTriggerCallerClassMerges(t *testing.T) {
	got := render(t, ui.SelectTrigger("", gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "w-[180px]"}}))
	if strings.Contains(got, "w-fit") {
		t.Errorf("base w-fit should be dropped by caller w-[180px]\nin: %s", got)
	}
	if !strings.Contains(got, "w-[180px]") {
		t.Errorf("caller width should land on the trigger\nin: %s", got)
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
	if !strings.Contains(got, `data-slot="icon"`) || !strings.Contains(got, `d="M20 6 9 17l-5-5"`) {
		t.Errorf("expected the Check icon svg in the item render\nin: %s", got)
	}
}

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
		`data-slot="combobox-item"`,
		`data-value="next.js"`,
		`aria-selected="true"`,
		`data-slot="combobox-item-indicator"`,
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
		`data-slot="combobox"`,
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
		`data-slot="input-group"`,
		`data-slot="input-group-input"`,
		`data-slot="combobox-trigger"`,
		`role="combobox"`,
		`aria-expanded="false"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
	if strings.Contains(got, `data-slot="combobox-clear"`) {
		t.Errorf("showClear=false must not render the clear button\nin: %s", got)
	}
}

func TestComboboxInputShowClear(t *testing.T) {
	got := render(t, ui.ComboboxInput("", false, true, false, nil, nil))
	if !strings.Contains(got, `data-slot="combobox-clear"`) {
		t.Errorf("want the clear button\nin: %s", got)
	}
	if strings.Contains(got, `data-slot="combobox-trigger"`) {
		t.Errorf("showTrigger=false must not render the trigger\nin: %s", got)
	}
}

func TestComboboxContentCarriesDiscreteTransitionBlock(t *testing.T) {
	// The popover family's shared exit-animation mechanism — must be
	// byte-identical to ui/select.gsx's content, not re-derived.
	got := render(t, ui.ComboboxContent(gsx.Raw("x"), nil))
	for _, want := range []string{
		`popover="auto"`,
		`data-slot="combobox-content"`,
		"transition-discrete",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

func TestComboboxEmptyIsHiddenUntilListIsEmpty(t *testing.T) {
	got := render(t, ui.ComboboxEmpty(gsx.Raw("No framework found."), nil))
	for _, want := range []string{
		`data-slot="combobox-empty"`,
		"group-data-empty/combobox-content:flex",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q\nin: %s", want, got)
		}
	}
}

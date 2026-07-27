package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestTooltipPinned(t *testing.T) {
	got := render(t, ui.Tooltip(gsx.Fragment(
		ui.TooltipTrigger(gsx.Raw("Open"), nil),
		ui.TooltipContent(gsx.Raw("Body"), nil),
	), nil))
	want := `<div data-gsxui-tooltip data-gsxui-slot-tooltip><button data-gsxui-tooltip-trigger type="button" data-gsxui-slot-tooltip-trigger>Open</button><div data-gsxui-tooltip-content popover="manual" role="tooltip" data-state="closed" data-side="top" data-gsxui-slot-tooltip-content>Body<span data-gsxui-slot-tooltip-arrow></span></div></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestTooltipAttrsOverrideDefaults(t *testing.T) {
	got := render(t, ui.TooltipContent(gsx.Raw("x"), gsx.Attrs{
		{Key: "popover", Value: "auto"},
		{Key: "class", Value: "caller"},
	}))
	if !strings.Contains(got, `popover="auto"`) || strings.Contains(got, `popover="manual"`) {
		t.Errorf("caller popover value must override default\nin: %s", got)
	}
	if strings.Count(got, `class="caller"`) != 1 {
		t.Errorf("caller class must render once\nin: %s", got)
	}
}

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
	want := `<div class="contents" data-gsxui-slot-tooltip><button type="button" data-gsxui-slot-tooltip-trigger>Open</button><div popover="manual" role="tooltip" data-state="closed" data-side="top" class="data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95 data-[state=delayed-open]:animate-in data-[state=delayed-open]:fade-in-0 data-[state=delayed-open]:zoom-in-95 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 [&amp;:popover-open]:inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-xs has-[[data-gsxui-slot-kbd]]:pr-1.5 [&amp;_[data-gsxui-slot-kbd]]:relative [&amp;_[data-gsxui-slot-kbd]]:isolate [&amp;_[data-gsxui-slot-kbd]]:z-50 [&amp;_[data-gsxui-slot-kbd]]:rounded-sm" data-gsxui-slot-tooltip-content>Body<span class="size-2.5 translate-y-[calc(-50%-2px)] rotate-45 rounded-[2px]" data-gsxui-slot-tooltip-arrow></span></div></div>`
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
	if strings.Count(got, "caller") != 1 {
		t.Errorf("caller class must render once\nin: %s", got)
	}
}

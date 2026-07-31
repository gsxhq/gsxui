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
	want := `<div class="contents" data-gsxui-slot-tooltip><button data-gsxui-tooltip-trigger type="button" data-gsxui-slot-tooltip-trigger>Open</button><div popover="manual" role="tooltip" data-state="closed" data-side="top" class="z-50 w-fit origin-bottom gap-1.5 overflow-visible rounded-md bg-foreground px-3 py-1.5 text-xs text-balance text-background has-[[data-gsxui-slot-kbd]]:pr-1.5 opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 [&amp;:popover-open]:opacity-100 [&amp;:popover-open]:scale-100 starting:[&amp;:popover-open]:opacity-0 starting:[&amp;:popover-open]:scale-95 starting:[&amp;:popover-open]:data-[side=bottom]:-translate-y-2 starting:[&amp;:popover-open]:data-[side=left]:translate-x-2 starting:[&amp;:popover-open]:data-[side=right]:-translate-x-2 starting:[&amp;:popover-open]:data-[side=top]:translate-y-2" data-gsxui-slot-tooltip-content>Body<span class="absolute top-full left-1/2 z-50 size-2.5 -translate-x-1/2 -translate-y-[calc(50%+2px)] rotate-45 rounded-[2px] bg-foreground" data-gsxui-slot-tooltip-arrow></span></div></div>`
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

package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestPopoverPinned(t *testing.T) {
	got := render(t, ui.Popover(gsx.Fragment(
		ui.PopoverTrigger(gsx.Raw("Open"), nil),
		ui.PopoverContent(gsx.Raw("Body"), nil),
	), nil))
	want := `<div class="contents" data-gsxui-slot-popover><button type="button" aria-expanded="false" data-gsxui-slot-popover-trigger>Open</button><div popover="auto" data-state="closed" data-side="bottom" tabindex="-1" class="z-50 w-72 origin-top gap-2.5 rounded-lg border bg-popover p-2.5 text-sm text-popover-foreground shadow-md outline-hidden opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 [&amp;:popover-open]:opacity-100 [&amp;:popover-open]:scale-100 starting:[&amp;:popover-open]:opacity-0 starting:[&amp;:popover-open]:scale-95 starting:[&amp;:popover-open]:data-[side=bottom]:-translate-y-2 starting:[&amp;:popover-open]:data-[side=left]:translate-x-2 starting:[&amp;:popover-open]:data-[side=right]:-translate-x-2 starting:[&amp;:popover-open]:data-[side=top]:translate-y-2" data-gsxui-slot-popover-content>Body</div></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestPopoverAttrsOverrideDefaults(t *testing.T) {
	got := render(t, ui.PopoverContent(gsx.Raw("x"), gsx.Attrs{
		{Key: "popover", Value: "manual"},
		{Key: "class", Value: "caller"},
	}))
	if !strings.Contains(got, `popover="manual"`) || strings.Contains(got, `popover="auto"`) {
		t.Errorf("caller popover value must override default\nin: %s", got)
	}
	if strings.Count(got, "caller") != 1 || strings.Count(got, "class=") != 1 {
		t.Errorf("caller class must render once\nin: %s", got)
	}
}

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
	want := `<div class="contents" data-gsxui-slot-popover><button type="button" aria-expanded="false" data-gsxui-slot-popover-trigger>Open</button><div popover="auto" data-state="closed" data-side="bottom" tabindex="-1" class="bg-popover text-popover-foreground data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 ring-foreground/10 [&amp;:popover-open]:flex flex-col gap-2.5 rounded-lg p-2.5 text-sm shadow-md ring-1 duration-100" data-gsxui-slot-popover-content>Body</div></div>`
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

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
	want := `<div data-gsxui-popover data-gsxui-slot-popover><button data-gsxui-popover-trigger type="button" aria-expanded="false" data-gsxui-slot-popover-trigger>Open</button><div data-gsxui-popover-content popover="auto" data-state="closed" data-side="bottom" tabindex="-1" data-gsxui-slot-popover-content>Body</div></div>`
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
	if strings.Count(got, `class="caller"`) != 1 {
		t.Errorf("caller class must render once\nin: %s", got)
	}
}

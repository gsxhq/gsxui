package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestHoverCardPinned(t *testing.T) {
	got := render(t, ui.HoverCard(gsx.Fragment(
		ui.HoverCardTrigger(gsx.Raw("Open"), nil),
		ui.HoverCardContent(gsx.Raw("Body"), nil),
	), nil))
	want := `<div data-gsxui-hovercard class="contents" data-gsxui-slot-hover-card><span data-gsxui-hovercard-trigger data-gsxui-slot-hover-card-trigger>Open</span><div data-gsxui-hovercard-content popover="manual" data-state="closed" data-side="bottom" class="z-50 w-64 origin-top rounded-lg border bg-popover p-2.5 text-sm text-popover-foreground shadow-md outline-hidden opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 [&amp;:popover-open]:opacity-100 [&amp;:popover-open]:scale-100 starting:[&amp;:popover-open]:opacity-0 starting:[&amp;:popover-open]:scale-95 starting:[&amp;:popover-open]:data-[side=bottom]:-translate-y-2 starting:[&amp;:popover-open]:data-[side=left]:translate-x-2 starting:[&amp;:popover-open]:data-[side=right]:-translate-x-2 starting:[&amp;:popover-open]:data-[side=top]:translate-y-2" data-gsxui-slot-hover-card-content>Body</div></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestHoverCardAttrsOverrideDefaults(t *testing.T) {
	got := render(t, ui.HoverCardContent(gsx.Raw("x"), gsx.Attrs{
		{Key: "popover", Value: "auto"},
		{Key: "class", Value: "caller"},
	}))
	if !strings.Contains(got, `popover="auto"`) || strings.Contains(got, `popover="manual"`) {
		t.Errorf("caller popover value must override default\nin: %s", got)
	}
	if strings.Count(got, "caller") != 1 || strings.Count(got, "class=") != 1 {
		t.Errorf("caller class must render once\nin: %s", got)
	}
}

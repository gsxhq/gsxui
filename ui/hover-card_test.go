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
	want := `<div data-gsxui-hovercard data-gsxui-slot-hover-card><span data-gsxui-hovercard-trigger data-gsxui-slot-hover-card-trigger>Open</span><div data-gsxui-hovercard-content popover="manual" data-state="closed" data-side="bottom" data-gsxui-slot-hover-card-content>Body</div></div>`
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
	if strings.Count(got, `class="caller"`) != 1 {
		t.Errorf("caller class must render once\nin: %s", got)
	}
}

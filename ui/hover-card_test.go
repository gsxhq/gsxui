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
	// transition-none suppresses an implicit transition:all — see the
	// style-porter report's "duration-N alone" entry.
	want := `<div class="contents" data-gsxui-slot-hover-card><span data-gsxui-slot-hover-card-trigger>Open</span><div popover="manual" data-state="closed" data-side="bottom" class="transition-none data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 ring-foreground/10 bg-popover text-popover-foreground w-64 rounded-lg p-2.5 text-sm shadow-md ring-1 duration-100" data-gsxui-slot-hover-card-content>Body</div></div>`
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

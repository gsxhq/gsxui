package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

// TestScrollAreaVerticalPinned pins the zero-value (vertical, default)
// orientation reflected explicitly for style and mechanism selectors.
func TestScrollAreaVerticalPinned(t *testing.T) {
	got := render(t, ui.ScrollArea("", gsx.Raw("x"), nil))
	want := `<div data-orientation="vertical" class="relative rounded-[inherit] outline-none transition-[color,box-shadow] focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1" data-gsxui-slot-scroll-area>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestScrollAreaHorizontalPinned pins orientation="horizontal".
func TestScrollAreaHorizontalPinned(t *testing.T) {
	got := render(t, ui.ScrollArea("horizontal", gsx.Raw("x"), nil))
	want := `<div data-orientation="horizontal" class="relative rounded-[inherit] outline-none transition-[color,box-shadow] focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1" data-gsxui-slot-scroll-area>x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestScrollAreaCallerClassMerges verifies caller-supplied sizing
// (h-72 w-48 rounded-md border, the scroll-area-demo shape) survives
// fall-through as the only rendered class.
func TestScrollAreaCallerClassMerges(t *testing.T) {
	got := render(t, ui.ScrollArea("", nil, gsx.Attrs{{Key: "class", Value: "h-72 w-48 rounded-md border"}}))
	if !strings.Contains(got, "h-72 w-48 rounded-md border") || strings.Count(got, "class=") != 1 {
		t.Errorf("caller class is not forwarded exactly\nin: %s", got)
	}
}

// TestScrollAreaCallerClassOverridesRounded verifies a caller-supplied
// rounded-* utility remains caller-owned.
func TestScrollAreaCallerClassOverridesRounded(t *testing.T) {
	got := render(t, ui.ScrollArea("", nil, gsx.Attrs{{Key: "class", Value: "rounded-full"}}))
	if !strings.Contains(got, "rounded-full") || strings.Contains(got, "rounded-[inherit]") {
		t.Errorf("missing caller class rounded-full, or recipe's rounded-[inherit] was not overridden\nin: %s", got)
	}
}

func TestScrollAreaAttrsFallThrough(t *testing.T) {
	got := render(t, ui.ScrollArea("", nil, gsx.Attrs{{Key: "id", Value: "sa1"}, {Key: "aria-label", Value: "Log"}}))
	for _, want := range []string{`id="sa1"`, `aria-label="Log"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestScrollAreaChildrenRenderInside(t *testing.T) {
	got := render(t, ui.ScrollArea("", gsx.Raw("<p>content</p>"), nil))
	if !strings.Contains(got, "<p>content</p>") {
		t.Errorf("missing rendered children\nin: %s", got)
	}
}

package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

// TestScrollAreaVerticalPinned pins the zero-value (vertical, default)
// orientation: no data-orientation attribute (Root never carried one in
// shadcn's own source — only the dropped ScrollBar part did), overflow-auto
// picked by the switch's default case.
func TestScrollAreaVerticalPinned(t *testing.T) {
	got := render(t, ui.ScrollArea("", gsx.Raw("x"), nil))
	want := `<div data-slot="scroll-area" class="relative rounded-[inherit] outline-none transition-[color,box-shadow] focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1 [scrollbar-width:thin] [scrollbar-color:var(--border)_transparent] overflow-auto">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestScrollAreaHorizontalPinned pins orientation="horizontal": overflow-x-auto
// replaces overflow-auto, matching Radix's opt-in
// <ScrollBar orientation="horizontal"/> shape with no second component.
func TestScrollAreaHorizontalPinned(t *testing.T) {
	got := render(t, ui.ScrollArea("horizontal", gsx.Raw("x"), nil))
	want := `<div data-slot="scroll-area" class="relative rounded-[inherit] outline-none transition-[color,box-shadow] focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1 [scrollbar-width:thin] [scrollbar-color:var(--border)_transparent] overflow-x-auto">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
	if strings.Contains(got, "overflow-auto\"") || strings.Contains(got, "overflow-auto ") {
		t.Errorf("horizontal must not carry overflow-auto\nin: %s", got)
	}
}

// TestScrollAreaCallerClassMerges verifies caller-supplied sizing
// (h-72 w-48 rounded-md border, the scroll-area-demo shape) survives
// fall-through and merges alongside the structural classes rather than
// replacing them (no conflicting utility to drop, unlike slider's w-full).
func TestScrollAreaCallerClassMerges(t *testing.T) {
	got := render(t, ui.ScrollArea("", nil, gsx.Attrs{{Key: "class", Value: "h-72 w-48 rounded-md border"}}))
	for _, want := range []string{"h-72", "w-48", "rounded-md", "border", "overflow-auto", "relative"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// TestScrollAreaCallerClassOverridesRounded verifies a caller-supplied
// rounded-* utility wins over the base rounded-[inherit], the ordinary
// caller-class-merge convention (attrs after base).
func TestScrollAreaCallerClassOverridesRounded(t *testing.T) {
	got := render(t, ui.ScrollArea("", nil, gsx.Attrs{{Key: "class", Value: "rounded-full"}}))
	if strings.Contains(got, "rounded-[inherit]") {
		t.Errorf("caller rounded-full must drop base rounded-[inherit]\nin: %s", got)
	}
	if !strings.Contains(got, "rounded-full") {
		t.Errorf("missing caller class rounded-full\nin: %s", got)
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

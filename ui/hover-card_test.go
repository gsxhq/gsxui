package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestHoverCardStructure(t *testing.T) {
	got := render(t, ui.HoverCard(
		gsx.Fragment(
			ui.HoverCardTrigger(gsx.Raw("@nextjs"), nil),
			ui.HoverCardContent(gsx.Raw("The React Framework"), nil),
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-hovercard`, // root hook
		`class="contents"`,     // root is layout-neutral
		`data-slot="hover-card-trigger"`,
		`data-gsxui-hovercard-trigger`, // trigger hook
		">@nextjs<",
		`data-slot="hover-card-content"`,
		`data-gsxui-hovercard-content`,
		`popover="manual"`,    // top-layer, no light dismiss (hover/focus driven)
		`data-state="closed"`, // server-rendered initial state
		`data-side="bottom"`,  // Radix HoverCard's own default side, unlike tooltip's top
		">The React Framework<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestHoverCardTriggerIsSpan(t *testing.T) {
	// The trigger is a phrasing <span> wrapper, not a <button> — children
	// carry the real interactive element (an <a>, or a
	// <ui.Button variant="link">), see docs/jsx-parity.md ## hover-card.
	got := render(t, ui.HoverCardTrigger(gsx.Raw("x"), nil))
	if !strings.HasPrefix(got, "<span") {
		t.Errorf("HoverCardTrigger should render a <span>\nin: %s", got)
	}
	if strings.Contains(got, "<button") {
		t.Errorf("HoverCardTrigger should not render its own <button>\nin: %s", got)
	}
}

func TestHoverCardContentCallerClassMerges(t *testing.T) {
	// Caller w-80 must WIN over base w-64 via tailwind-merge (the shadcn
	// @nextjs demo does exactly this) — and base structural classes must
	// survive.
	got := render(t, ui.HoverCardContent(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "w-80"}}))
	if strings.Contains(got, "w-64") {
		t.Errorf("base w-64 should be dropped by caller w-80\nin: %s", got)
	}
	for _, want := range []string{"w-80", "rounded-lg", "bg-popover"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestHoverCardAttrsFallThrough(t *testing.T) {
	got := render(t, ui.HoverCardContent(gsx.Raw("x"), gsx.Attrs{
		{Key: "id", Value: "my-hover-card"},
		{Key: "aria-label", Value: "Profile preview"},
	}))
	for _, want := range []string{`id="my-hover-card"`, `aria-label="Profile preview"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestHoverCardPopoverAttrOverridable(t *testing.T) {
	// popover is a regular attribute with a value — attrs fallthrough can
	// override it like any other attribute.
	got := render(t, ui.HoverCardContent(gsx.Raw("x"), gsx.Attrs{{Key: "popover", Value: "auto"}}))
	if strings.Contains(got, `popover="manual"`) {
		t.Errorf("caller popover=auto should replace the default manual\nin: %s", got)
	}
	if !strings.Contains(got, `popover="auto"`) {
		t.Errorf("missing overridden popover=\"auto\"\nin: %s", got)
	}
}

func TestHoverCardPinned(t *testing.T) {
	// Exact full-render pin for HoverCardContent, verified token-by-token
	// against shadcn's HoverCardContent classes (registry/new-york-v4/ui/
	// hover-card.tsx) and docs/jsx-parity.md's ADAPT: popover="manual"/
	// data-state/data-side replace Radix's Portal+Content wiring, origin-top
	// replaces the Radix runtime transform-origin var.
	got := render(t, ui.HoverCardContent(gsx.Raw("x"), nil))
	want := `<div data-slot="hover-card-content" data-gsxui-hovercard-content popover="manual" data-state="closed" data-side="bottom" class="z-50 w-64 origin-top rounded-lg border bg-popover p-2.5 text-sm text-popover-foreground shadow-md outline-hidden opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:opacity-0 starting:open:scale-95 data-[side=bottom]:starting:open:-translate-y-2 data-[side=left]:starting:open:translate-x-2 data-[side=right]:starting:open:-translate-x-2 data-[side=top]:starting:open:translate-y-2">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestHoverCardTriggerPinned(t *testing.T) {
	got := render(t, ui.HoverCardTrigger(gsx.Raw("@nextjs"), nil))
	want := `<span data-slot="hover-card-trigger" data-gsxui-hovercard-trigger>@nextjs</span>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestTooltipStructure(t *testing.T) {
	got := render(t, ui.Tooltip(
		gsx.Fragment(
			ui.TooltipTrigger(gsx.Raw("Hover me"), nil),
			ui.TooltipContent(gsx.Raw("Add to library"), nil),
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-tooltip`,         // root hook
		`class="contents"`,           // root is layout-neutral
		`data-gsxui-tooltip-trigger`, // trigger hook
		">Hover me<",
		`data-slot="tooltip-content"`,
		`data-gsxui-tooltip-content`,
		`popover="manual"`,    // top-layer, no light dismiss (hover/focus driven)
		`role="tooltip"`,      // content a11y
		`data-state="closed"`, // server-rendered initial state
		">Add to library<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestTooltipTriggerType(t *testing.T) {
	got := render(t, ui.TooltipTrigger(gsx.Raw("x"), nil))
	if !strings.Contains(got, `type="button"`) {
		t.Errorf("missing type=\"button\"\nin: %s", got)
	}
}

func TestTooltipContentCallerClassMerges(t *testing.T) {
	// Caller z-10 must WIN over base z-50 via tailwind-merge — and base
	// structural classes must survive.
	got := render(t, ui.TooltipContent(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "z-10"}}))
	if strings.Contains(got, "z-50") {
		t.Errorf("base z-50 should be dropped by caller z-10\nin: %s", got)
	}
	for _, want := range []string{"z-10", "rounded-md", "bg-foreground"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestTooltipAttrsFallThrough(t *testing.T) {
	got := render(t, ui.TooltipContent(gsx.Raw("x"), gsx.Attrs{
		{Key: "id", Value: "my-tip"},
		{Key: "aria-label", Value: "Tip"},
	}))
	for _, want := range []string{`id="my-tip"`, `aria-label="Tip"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestTooltipPopoverAttrOverridable(t *testing.T) {
	// popover is a regular attribute with a value — attrs fallthrough can
	// override it like any other attribute.
	got := render(t, ui.TooltipContent(gsx.Raw("x"), gsx.Attrs{{Key: "popover", Value: "auto"}}))
	if strings.Contains(got, `popover="manual"`) {
		t.Errorf("caller popover=auto should replace the default manual\nin: %s", got)
	}
	if !strings.Contains(got, `popover="auto"`) {
		t.Errorf("missing overridden popover=\"auto\"\nin: %s", got)
	}
}

func TestTooltipPinned(t *testing.T) {
	// Exact full-render pin for TooltipContent, verified token-by-token
	// against shadcn's TooltipContent classes (registry/new-york-v4/ui/
	// tooltip.tsx) and docs/jsx-parity.md's ADAPT: popover="manual"/role=
	// "tooltip"/data-state replace Radix's Portal+Content wiring; the Arrow
	// part is dropped.
	got := render(t, ui.TooltipContent(gsx.Raw("Add to library"), nil))
	want := `<div data-slot="tooltip-content" data-gsxui-tooltip-content popover="manual" role="tooltip" data-state="closed" class="z-50 w-fit origin-bottom animate-in rounded-md bg-foreground px-3 py-1.5 text-xs text-balance text-background fade-in-0 zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95">Add to library</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

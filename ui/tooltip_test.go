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
	// structural classes must survive. Scope the z-50 check to the CONTENT
	// element's class attribute (the first one): the arrow child span
	// legitimately carries its own z-50 (shadcn's Arrow classes verbatim),
	// which the caller's class on the content must not touch.
	got := render(t, ui.TooltipContent(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "z-10"}}))
	contentClass := got[strings.Index(got, `class="`)+len(`class="`):]
	contentClass = contentClass[:strings.Index(contentClass, `"`)]
	if strings.Contains(contentClass, "z-50") {
		t.Errorf("base z-50 should be dropped by caller z-10\ncontent class: %s", contentClass)
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
	// "tooltip"/data-state replace Radix's Portal+Content wiring. The Arrow
	// ports as a static child span: our tooltip is always JS-anchored above
	// the trigger, so the diamond always straddles the bubble's
	// bottom-center — no Radix side-tracking slot needed.
	got := render(t, ui.TooltipContent(gsx.Raw("Add to library"), nil))
	want := `<div data-slot="tooltip-content" data-gsxui-tooltip-content popover="manual" role="tooltip" data-state="closed" data-side="top" class="z-50 w-fit origin-bottom gap-1.5 rounded-md bg-foreground px-3 py-1.5 text-xs has-data-[slot=kbd]:pr-1.5 **:data-[slot=kbd]:rounded-sm text-balance text-background overflow-visible opacity-0 scale-95 transition-[opacity,scale,translate,display,overlay] transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:opacity-0 starting:open:scale-95 data-[side=bottom]:starting:open:-translate-y-2 data-[side=left]:starting:open:translate-x-2 data-[side=right]:starting:open:-translate-x-2 data-[side=top]:starting:open:translate-y-2">Add to library<span data-slot="tooltip-arrow" class="absolute top-full left-1/2 z-50 size-2.5 -translate-x-1/2 -translate-y-[calc(50%+2px)] rotate-45 rounded-[2px] bg-foreground"></span></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

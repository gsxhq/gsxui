package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestPopoverStructure(t *testing.T) {
	got := render(t, ui.Popover(
		gsx.Fragment(
			ui.PopoverTrigger(gsx.Raw("Open"), nil),
			ui.PopoverContent(gsx.Raw("Set the dimensions"), nil),
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-popover`,         // root hook
		`class="contents"`,           // root is layout-neutral
		`data-gsxui-popover-trigger`, // trigger hook
		`aria-expanded="false"`,      // trigger a11y: initial state
		`type="button"`,
		">Open<",
		`data-slot="popover-content"`,
		`data-gsxui-popover-content`,
		`popover="auto"`,      // top-layer, light-dismiss, free Esc
		`data-state="closed"`, // server-rendered initial state
		`data-side="bottom"`,  // popover.js always anchors below the trigger
		">Set the dimensions<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestPopoverTriggerType(t *testing.T) {
	got := render(t, ui.PopoverTrigger(gsx.Raw("x"), nil))
	if !strings.Contains(got, `type="button"`) {
		t.Errorf("missing type=\"button\"\nin: %s", got)
	}
	if !strings.Contains(got, `aria-expanded="false"`) {
		t.Errorf("missing aria-expanded=\"false\"\nin: %s", got)
	}
}

func TestPopoverContentCallerClassMerges(t *testing.T) {
	// Caller w-80 must WIN over base w-72 via tailwind-merge (the shadcn
	// dimensions-form demo does exactly this) — and base structural classes
	// must survive.
	got := render(t, ui.PopoverContent(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "w-80"}}))
	if strings.Contains(got, "w-72") {
		t.Errorf("base w-72 should be dropped by caller w-80\nin: %s", got)
	}
	for _, want := range []string{"w-80", "rounded-md", "bg-popover"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestPopoverAttrsFallThrough(t *testing.T) {
	got := render(t, ui.PopoverContent(gsx.Raw("x"), gsx.Attrs{
		{Key: "id", Value: "my-popover"},
		{Key: "aria-label", Value: "Dimensions"},
	}))
	for _, want := range []string{`id="my-popover"`, `aria-label="Dimensions"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestPopoverPopoverAttrOverridable(t *testing.T) {
	// popover is a regular attribute with a value — attrs fallthrough can
	// override it like any other attribute (e.g. a caller opting a specific
	// popover out of the default "auto" light-dismiss behavior).
	got := render(t, ui.PopoverContent(gsx.Raw("x"), gsx.Attrs{{Key: "popover", Value: "manual"}}))
	if strings.Contains(got, `popover="auto"`) {
		t.Errorf("caller popover=manual should replace the default auto\nin: %s", got)
	}
	if !strings.Contains(got, `popover="manual"`) {
		t.Errorf("missing overridden popover=\"manual\"\nin: %s", got)
	}
}

func TestPopoverPinned(t *testing.T) {
	// Exact full-render pin for PopoverContent, verified token-by-token
	// against shadcn's PopoverContent classes (registry/new-york-v4/ui/
	// popover.tsx) and docs/jsx-parity.md's ADAPT: popover="auto"/data-state/
	// data-side replace Radix's Portal+Content wiring, origin-top replaces
	// the Radix runtime transform-origin var.
	got := render(t, ui.PopoverContent(gsx.Raw("x"), nil))
	want := `<div data-slot="popover-content" data-gsxui-popover-content popover="auto" data-state="closed" data-side="bottom" class="z-50 w-72 origin-top rounded-md border bg-popover p-4 text-popover-foreground shadow-md outline-hidden data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestPopoverTriggerPinned(t *testing.T) {
	got := render(t, ui.PopoverTrigger(gsx.Raw("Open popover"), nil))
	want := `<button data-slot="popover-trigger" data-gsxui-popover-trigger type="button" aria-expanded="false">Open popover</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

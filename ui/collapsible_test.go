package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

// TestCollapsibleClosedDefault covers the Go zero value: open bool's zero
// value (false) must render <details> WITHOUT the open attribute at all —
// matching shadcn's Radix default (defaultOpen unset == collapsed).
func TestCollapsibleClosedDefault(t *testing.T) {
	got := render(t, ui.Collapsible(false, gsx.Raw("x"), nil))
	if !strings.Contains(got, `<details data-slot="collapsible">`) {
		t.Errorf("closed collapsible must not render the open attribute\nin: %s", got)
	}
	if strings.Contains(got, "open") {
		t.Errorf("closed collapsible must not mention open at all\nin: %s", got)
	}
}

// TestCollapsibleOpenStamping covers the true branch: bare `open` boolean
// attribute (not open="true"/open="false"), same native-<details> mechanism
// as AccordionItem.
func TestCollapsibleOpenStamping(t *testing.T) {
	got := render(t, ui.Collapsible(true, gsx.Raw("x"), nil))
	if !strings.Contains(got, `<details data-slot="collapsible" open>`) {
		t.Errorf("open collapsible must render a bare open attribute\nin: %s", got)
	}
}

func TestCollapsibleAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Collapsible(false, gsx.Raw("x"), gsx.Attrs{{Key: "id", Value: "c1"}, {Key: "class", Value: "flex flex-col gap-2"}}))
	for _, want := range []string{`id="c1"`, "flex flex-col gap-2"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestCollapsiblePinned(t *testing.T) {
	// Exact full-render pin, closed. shadcn's Collapsible root carries no
	// classes of its own (registry/new-york-v4/ui/collapsible.tsx is a bare
	// data-slot passthrough) — no ADAPT applies here.
	got := render(t, ui.Collapsible(false, gsx.Raw("x"), nil))
	want := `<details data-slot="collapsible">x</details>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestCollapsibleTriggerMarkerSuppressed covers the ADAPT: shadcn ships no
// classes at all for CollapsibleTrigger (it's a Radix button, no native
// disclosure triangle to suppress); porting onto <summary> introduces one,
// so list-none + the webkit marker selector are added here — exactly
// accordion's marker-suppression story (see accordion.gsx / docs/
// jsx-parity.md's `## accordion` ADAPT entry), ledgered under `## collapsible`.
func TestCollapsibleTriggerMarkerSuppressed(t *testing.T) {
	got := render(t, ui.CollapsibleTrigger(gsx.Raw("Toggle"), nil))
	for _, want := range []string{
		`data-slot="collapsible-trigger"`,
		"list-none",
		"[&amp;::-webkit-details-marker]:hidden",
		">Toggle<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// TestCollapsibleTriggerCallerClassMerges covers the tailwind-merge path on
// CollapsibleTrigger's base class, same shape as
// TestAccordionCallerClassMerges: a caller-supplied list-disc conflicts
// with the base list-none (both set the list-style-type/list utility
// category) and must win, while the unrelated webkit-marker-hiding
// arbitrary variant is untouched (different bucket, no collision).
func TestCollapsibleTriggerCallerClassMerges(t *testing.T) {
	got := render(t, ui.CollapsibleTrigger(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "list-disc"}}))
	if strings.Contains(got, "list-none") {
		t.Errorf("base list-none should be dropped by caller list-disc\nin: %s", got)
	}
	for _, want := range []string{"list-disc", "[&amp;::-webkit-details-marker]:hidden"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestCollapsibleTriggerAttrsFallThrough(t *testing.T) {
	got := render(t, ui.CollapsibleTrigger(nil, gsx.Attrs{{Key: "id", Value: "t1"}, {Key: "aria-label", Value: "toggle"}}))
	for _, want := range []string{`id="t1"`, `aria-label="toggle"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestCollapsibleTriggerPinned(t *testing.T) {
	got := render(t, ui.CollapsibleTrigger(gsx.Raw("Toggle"), nil))
	want := `<summary data-slot="collapsible-trigger" class="list-none [&amp;::-webkit-details-marker]:hidden">Toggle</summary>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestCollapsibleContentPinned covers the design's plain-div ADAPT: no
// data-state, no animation classes (shadcn ships none for Collapsible) —
// CSS consumers wanting open/closed styling target the ancestor <details>'
// [open] attribute instead (ledgered GAP: no data-state without JS).
func TestCollapsibleContentPinned(t *testing.T) {
	got := render(t, ui.CollapsibleContent(gsx.Raw("x"), nil))
	want := `<div data-slot="collapsible-content">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCollapsibleContentAttrsFallThrough(t *testing.T) {
	got := render(t, ui.CollapsibleContent(gsx.Raw("x"), gsx.Attrs{{Key: "id", Value: "panel-1"}, {Key: "class", Value: "flex flex-col gap-2"}}))
	for _, want := range []string{`id="panel-1"`, "flex flex-col gap-2"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// TestCollapsibleStructure is an integration-shaped smoke test composing all
// three parts together, mirroring TestAccordionStructure.
func TestCollapsibleStructure(t *testing.T) {
	got := render(t, ui.Collapsible(true, gsx.Fragment(
		ui.CollapsibleTrigger(gsx.Raw("Toggle"), nil),
		ui.CollapsibleContent(gsx.Raw("Body"), nil),
	), nil))
	for _, want := range []string{
		`<details data-slot="collapsible" open>`,
		`<summary data-slot="collapsible-trigger"`,
		`<div data-slot="collapsible-content">`,
		">Toggle<",
		">Body<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

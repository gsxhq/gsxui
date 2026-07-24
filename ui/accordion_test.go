package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestAccordionStructure(t *testing.T) {
	got := render(t, ui.Accordion("grp", gsx.Fragment(
		ui.AccordionItem("grp", true, gsx.Fragment(
			ui.AccordionTrigger(gsx.Raw("Item 1"), nil),
			ui.AccordionContent(gsx.Raw("Body 1"), nil),
		), nil),
		ui.AccordionItem("grp", false, gsx.Fragment(
			ui.AccordionTrigger(gsx.Raw("Item 2"), nil),
			ui.AccordionContent(gsx.Raw("Body 2"), nil),
		), nil),
	), nil))
	for _, want := range []string{
		`data-slot="accordion"`,
		`<details`,
		`data-slot="accordion-item"`,
		`<summary`,
		`data-slot="accordion-trigger"`,
		`data-slot="accordion-content"`,
		">Item 1<", ">Item 2<",
		">Body 1<", ">Body 2<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// TestAccordionItemNameGroups covers native <details name> exclusive
// grouping (ledger WIN replacing Radix's state machine): every item in a
// group must carry the SAME name attribute — that's the entire mechanism,
// no JS reconciling them.
func TestAccordionItemNameGroups(t *testing.T) {
	got := render(t, gsx.Fragment(
		ui.AccordionItem("faq", false, gsx.Raw("a"), nil),
		ui.AccordionItem("faq", false, gsx.Raw("b"), nil),
	))
	if n := strings.Count(got, `name="faq"`); n != 2 {
		t.Errorf("want 2 items sharing name=\"faq\", got %d\nin: %s", n, got)
	}
}

func TestAccordionItemOpenStamping(t *testing.T) {
	// The bare `open` attribute must render as a standalone token (native
	// boolean attribute), not open="false"/open="true" — spot-check its
	// exact position between name and class.
	open := render(t, ui.AccordionItem("g", true, gsx.Raw("x"), nil))
	if !strings.Contains(open, `name="g" open class=`) {
		t.Errorf("open item's open attribute did not render bare before class\nin: %s", open)
	}

	closed := render(t, ui.AccordionItem("g", false, gsx.Raw("x"), nil))
	if strings.Contains(closed, "open") {
		t.Errorf("closed item must not render the open attribute\nin: %s", closed)
	}
}

func TestAccordionTriggerChevron(t *testing.T) {
	got := render(t, ui.AccordionTrigger(gsx.Raw("Section"), nil))
	for _, want := range []string{
		`data-slot="accordion-trigger"`,
		"list-none",
		"[&amp;::-webkit-details-marker]:hidden", // native marker suppressed both engines
		`data-slot="icon"`,                       // the chevron
		"[[open]&gt;summary_&amp;]:rotate-180",   // ancestor-[open] arbitrary variant, HTML-escaped
		">Section<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestAccordionCallerClassMerges(t *testing.T) {
	// border-b-4 conflicts with the base border-b (both set
	// border-bottom-width) and must win; last:border-b-0 targets a
	// different pseudo-class bucket and survives untouched.
	got := render(t, ui.AccordionItem("g", false, gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "border-b-4"}}))
	if strings.Contains(got, "border-b last") {
		t.Errorf("base border-b should be dropped by caller border-b-4\nin: %s", got)
	}
	for _, want := range []string{"border-b-4", "last:border-b-0"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestAccordionAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Accordion("g", gsx.Raw("x"), gsx.Attrs{{Key: "id", Value: "acc1"}, {Key: "aria-label", Value: "faq"}}))
	for _, want := range []string{`id="acc1"`, `aria-label="faq"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// TestAccordionContentCallerClassRoutesToInner covers the shadcn parity fix:
// rest props land on the OUTER Content element (class hardcoded to
// "overflow-hidden text-sm"), while a caller's className merges onto the
// INNER `pt-0 pb-2.5` wrapper div — matching shadcn's cn() split between the
// two divs, expressed here via attrs.Without("class") / attrs.Class().
func TestAccordionContentCallerClassRoutesToInner(t *testing.T) {
	got := render(t, ui.AccordionContent(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "pb-8"}}))

	// pb-8 conflicts with the inner pt-0 pb-2.5's pb-2.5 (both set
	// padding-bottom) and must win there; pt-0 survives untouched.
	if !strings.Contains(got, `<div class="pt-0 pb-8">x</div>`) {
		t.Errorf("caller class must merge onto the INNER wrapper, dropping pb-2.5\nin: %s", got)
	}
	if strings.Contains(got, "pb-2.5") {
		t.Errorf("pb-2.5 should have been dropped by caller's pb-8\nin: %s", got)
	}
	// The outer div's class is the hardcoded shadcn constant only — no
	// caller class landed there.
	if !strings.Contains(got, `<div data-slot="accordion-content" class="overflow-hidden text-sm">`) {
		t.Errorf("outer div's class must stay the hardcoded constant, unmerged with caller class\nin: %s", got)
	}
}

// TestAccordionContentNonClassAttrsRouteToOuter covers the other half of the
// routing split: non-class rest props (id, aria-*, data-*, …) land on the
// OUTER Content element, matching Radix's rest-prop spread target.
func TestAccordionContentNonClassAttrsRouteToOuter(t *testing.T) {
	got := render(t, ui.AccordionContent(gsx.Raw("x"), gsx.Attrs{{Key: "id", Value: "panel-1"}}))

	if !strings.Contains(got, `<div data-slot="accordion-content" class="overflow-hidden text-sm" id="panel-1">`) {
		t.Errorf("id must land on the OUTER div\nin: %s", got)
	}
	if !strings.Contains(got, `<div class="pt-0 pb-2.5">x</div>`) {
		t.Errorf("inner wrapper must be unaffected by non-class attrs\nin: %s", got)
	}
}

func TestAccordionItemPinned(t *testing.T) {
	// Exact full-render pin for a closed AccordionItem, verified token-by-
	// token against shadcn's AccordionItem (registry/new-york-v4/ui/
	// accordion.tsx) — border-b last:border-b-0 carries over verbatim, no
	// ADAPT applies to this part.
	got := render(t, ui.AccordionItem("g", false, gsx.Raw("x"), nil))
	want := `<details data-slot="accordion-item" name="g" class="border-b last:border-b-0">x</details>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

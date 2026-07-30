package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func assertCSSOnlyMarkup(t *testing.T, got string, slots ...string) {
	t.Helper()
	if strings.Contains(got, `data-slot=`) {
		t.Errorf("legacy data-slot must not render\nin: %s", got)
	}
	for _, slot := range slots {
		for name := range strings.FieldsSeq(slot) {
			if !strings.Contains(got, `data-gsxui-slot-`+name) {
				t.Errorf("missing slot marker %q\nin: %s", name, got)
			}
		}
	}
}

func assertCallerAttrsOnce(t *testing.T, got string) {
	t.Helper()
	for _, attr := range []string{`class="caller-only"`, `id="caller-id"`} {
		if strings.Count(got, attr) != 1 {
			t.Errorf("%s must render exactly once\nin: %s", attr, got)
		}
	}
	if strings.Count(got, `class=`) != 1 {
		t.Errorf("caller class must be the only rendered class\nin: %s", got)
	}
}

func callerAttrs() gsx.Attrs {
	return gsx.Attrs{
		{Key: "class", Value: "caller-only"},
		{Key: "id", Value: "caller-id"},
	}
}

func TestAccordionCSSOnlyContract(t *testing.T) {
	// Accordion migrated to the slot axis: unlike the other CSS-only parts
	// this test covers, every one of its slots (item, trigger, trigger-icon,
	// content, content-inner) now legitimately renders a class="..." attribute
	// carrying its recipe's resolved utilities — assertCallerAttrsOnce's
	// blanket "class= appears exactly once in the whole render" assumption
	// no longer holds. Assert caller attrs merge onto their own element
	// (item) exactly once instead, mirroring InputGroupButton's precedent
	// for a migrated composed part.
	got := render(t, ui.AccordionItem("faq", true, gsx.Fragment(
		ui.AccordionTrigger(gsx.Raw("Question"), nil),
		ui.AccordionContent(gsx.Raw("Answer"), nil),
	), callerAttrs()))
	assertCSSOnlyMarkup(t, got,
		"accordion-item",
		"accordion-trigger",
		"icon accordion-trigger-icon",
		"accordion-content",
		"accordion-content-inner",
	)
	if strings.Count(got, `id="caller-id"`) != 1 {
		t.Errorf(`id="caller-id" must render exactly once\nin: %s`, got)
	}
	if !strings.Contains(got, "border-b last:border-b-0 caller-only") {
		t.Errorf("caller class must merge onto the item element's own recipe class\nin: %s", got)
	}
	if !strings.Contains(got, `name="faq" open`) {
		t.Errorf("open item must retain native grouped disclosure state\nin: %s", got)
	}
}

func TestCollapsibleCSSOnlyContract(t *testing.T) {
	got := render(t, ui.Collapsible(true, gsx.Fragment(
		ui.CollapsibleTrigger(gsx.Raw("Toggle"), nil),
		ui.CollapsibleContent(gsx.Raw("Body"), nil),
	), callerAttrs()))
	assertCSSOnlyMarkup(t, got, "collapsible", "collapsible-trigger", "collapsible-content")
	// Not assertCallerAttrsOnce: CollapsibleTrigger now carries its own
	// recipe class (list-none) in addition to the root's caller-supplied
	// class, so class= legitimately renders twice — same carve-out as
	// FieldLabel/InputGroupButton for a composed primitive's own class.
	for _, attr := range []string{`class="caller-only"`, `id="caller-id"`} {
		if strings.Count(got, attr) != 1 {
			t.Errorf("%s must render exactly once\nin: %s", attr, got)
		}
	}
	if !strings.Contains(got, `class="list-none"`) {
		t.Errorf("collapsible trigger must keep its own recipe class\nin: %s", got)
	}
	if strings.Count(got, `class=`) != 2 {
		t.Errorf("expected exactly 2 class= attributes (caller root + trigger)\nin: %s", got)
	}
	if !strings.Contains(got, `<details`) || !strings.Contains(got, ` open`) {
		t.Errorf("open collapsible must retain native disclosure state\nin: %s", got)
	}
}

func TestDialogCSSOnlyContract(t *testing.T) {
	got := render(t, ui.Dialog(
		gsx.Fragment(
			ui.DialogTrigger(gsx.Raw("Open"), nil),
			ui.DialogContent(false, gsx.Fragment(
				ui.DialogTitle(gsx.Raw("Title"), nil),
				ui.DialogDescription(gsx.Raw("Description"), nil),
			), callerAttrs()),
		),
		nil,
	))
	assertCSSOnlyMarkup(t, got,
		"dialog",
		"dialog-trigger",
		"dialog-content",
		"dialog-title",
		"dialog-description",
		"dialog-close dialog-close-button",
		"dialog-close-icon",
		"dialog-close-label",
	)
	assertCallerAttrsOnce(t, got)
	for _, hook := range []string{
		`data-gsxui-dialog-title`,
		`data-gsxui-dialog-description`,
		`data-state="closed"`,
	} {
		if !strings.Contains(got, hook) {
			t.Errorf("missing dialog behavior/state hook %q\nin: %s", hook, got)
		}
	}
}

func TestAlertDialogCSSOnlyContract(t *testing.T) {
	root := render(t, ui.AlertDialog(gsx.Raw("x"), callerAttrs()))
	assertCSSOnlyMarkup(t, root, "dialog alert-dialog")
	assertCallerAttrsOnce(t, root)

	content := render(t, ui.AlertDialogContent(gsx.Fragment(
		ui.AlertDialogTitle(gsx.Raw("Title"), nil),
		ui.AlertDialogDescription(gsx.Raw("Description"), nil),
	), callerAttrs()))
	assertCSSOnlyMarkup(t, content,
		"dialog-content alert-dialog-content",
		"alert-dialog-title",
		"alert-dialog-description",
	)
	assertCallerAttrsOnce(t, content)
	for _, want := range []string{
		`role="alertdialog"`,
		`data-state="closed"`,
		`data-gsxui-dialog-static`,
		`data-gsxui-dialog-title`,
		`data-gsxui-dialog-description`,
	} {
		if !strings.Contains(content, want) {
			t.Errorf("missing alert dialog semantic/state value %q\nin: %s", want, content)
		}
	}

	action := render(t, ui.AlertDialogAction(gsx.Raw("Continue"), callerAttrs()))
	assertCSSOnlyMarkup(t, action, "button alert-dialog-action")
	assertButtonCallerAttrsOnce(t, action, "default", "default")
	cancel := render(t, ui.AlertDialogCancel(gsx.Raw("Cancel"), callerAttrs()))
	assertCSSOnlyMarkup(t, cancel, "button alert-dialog-cancel")
	assertButtonCallerAttrsOnce(t, cancel, "outline", "default")
	if !strings.Contains(cancel, `data-variant="outline"`) {
		t.Errorf("cancel must retain the outline Button variant\nin: %s", cancel)
	}
}

func TestDrawerCSSOnlyContract(t *testing.T) {
	root := render(t, ui.Drawer(gsx.Raw("x"), callerAttrs()))
	assertCSSOnlyMarkup(t, root, "dialog drawer")
	assertCallerAttrsOnce(t, root)

	for _, side := range []string{"bottom", "top", "left", "right"} {
		input := side
		if side == "bottom" {
			input = ""
		}
		got := render(t, ui.DrawerContent(input, gsx.Raw("x"), callerAttrs()))
		assertCSSOnlyMarkup(t, got, "dialog-content drawer-content")
		assertCallerAttrsOnce(t, got)
		for _, want := range []string{`data-state="closed"`, `data-side="` + side + `"`} {
			if !strings.Contains(got, want) {
				t.Errorf("%s drawer missing %q\nin: %s", side, want, got)
			}
		}
	}
}

func TestSheetCSSOnlyContract(t *testing.T) {
	root := render(t, ui.Sheet(gsx.Raw("x"), callerAttrs()))
	assertCSSOnlyMarkup(t, root, "dialog sheet")
	assertCallerAttrsOnce(t, root)

	for _, side := range []string{"right", "left", "top", "bottom"} {
		input := side
		if side == "right" {
			input = ""
		}
		got := render(t, ui.SheetContent(input, true, gsx.Raw("x"), callerAttrs()))
		assertCSSOnlyMarkup(t, got, "dialog-content sheet-content")
		assertCallerAttrsOnce(t, got)
		for _, want := range []string{`data-state="closed"`, `data-side="` + side + `"`} {
			if !strings.Contains(got, want) {
				t.Errorf("%s sheet missing %q\nin: %s", side, want, got)
			}
		}
	}

	injected := render(t, ui.SheetContent("", false, gsx.Raw("x"), nil))
	assertCSSOnlyMarkup(t, injected, "sheet-close sheet-close-button", "sheet-close-icon", "sheet-close-label")
	if strings.Contains(injected, `class=`) {
		t.Errorf("injected sheet close parts must not carry presentation classes\nin: %s", injected)
	}
}

func TestPopoverCSSOnlyContract(t *testing.T) {
	// Popover migrated to the slot axis: popover-content now legitimately
	// renders a class="..." attribute carrying its recipe's resolved
	// utilities, so assertCallerAttrsOnce's blanket "class=\"caller-only\"
	// verbatim" assumption no longer holds. Assert the caller class merges
	// onto the content element's own recipe class instead, mirroring
	// Accordion's precedent for a migrated CSS-only part.
	got := render(t, ui.PopoverContent(gsx.Raw("x"), callerAttrs()))
	assertCSSOnlyMarkup(t, got, "popover-content")
	if strings.Count(got, `id="caller-id"`) != 1 {
		t.Errorf(`id="caller-id" must render exactly once\nin: %s`, got)
	}
	if strings.Count(got, "caller-only") != 1 || strings.Count(got, "class=") != 1 {
		t.Errorf("caller class must merge onto content's own recipe class exactly once\nin: %s", got)
	}
	for _, want := range []string{
		`popover="auto"`,
		`data-state="closed"`,
		`data-side="bottom"`,
		`tabindex="-1"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing popover mechanism/state value %q\nin: %s", want, got)
		}
	}
}

func TestHoverCardCSSOnlyContract(t *testing.T) {
	// HoverCard migrated: same carve-out as Popover above.
	got := render(t, ui.HoverCardContent(gsx.Raw("x"), callerAttrs()))
	assertCSSOnlyMarkup(t, got, "hover-card-content")
	if strings.Count(got, `id="caller-id"`) != 1 {
		t.Errorf(`id="caller-id" must render exactly once\nin: %s`, got)
	}
	if strings.Count(got, "caller-only") != 1 || strings.Count(got, "class=") != 1 {
		t.Errorf("caller class must merge onto content's own recipe class exactly once\nin: %s", got)
	}
	for _, want := range []string{
		`popover="manual"`,
		`data-state="closed"`,
		`data-side="bottom"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing hover-card mechanism/state value %q\nin: %s", want, got)
		}
	}
}

func TestTooltipCSSOnlyContract(t *testing.T) {
	// Tooltip migrated: same carve-out as Popover/HoverCard above, plus a
	// second class= for the arrow span's own recipe class (same shape as
	// CollapsibleTrigger's own class alongside the root's caller class).
	got := render(t, ui.TooltipContent(gsx.Raw("x"), callerAttrs()))
	assertCSSOnlyMarkup(t, got, "tooltip-content", "tooltip-arrow")
	if strings.Count(got, `id="caller-id"`) != 1 {
		t.Errorf(`id="caller-id" must render exactly once\nin: %s`, got)
	}
	if strings.Count(got, "caller-only") != 1 {
		t.Errorf("caller class must merge onto content's own recipe class exactly once\nin: %s", got)
	}
	if strings.Count(got, `class=`) != 2 {
		t.Errorf("expected exactly 2 class= attributes (content + arrow)\nin: %s", got)
	}
	for _, want := range []string{
		`popover="manual"`,
		`role="tooltip"`,
		`data-state="closed"`,
		`data-side="top"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing tooltip semantic/state value %q\nin: %s", want, got)
		}
	}
}

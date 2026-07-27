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
		if !strings.Contains(got, `data-gsxui-slot="`+slot+`"`) {
			t.Errorf("missing slot order %q\nin: %s", slot, got)
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
	assertCallerAttrsOnce(t, got)
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
	assertCallerAttrsOnce(t, got)
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
	assertCallerAttrsOnce(t, action)
	cancel := render(t, ui.AlertDialogCancel(gsx.Raw("Cancel"), callerAttrs()))
	assertCSSOnlyMarkup(t, cancel, "button alert-dialog-cancel")
	assertCallerAttrsOnce(t, cancel)
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
	got := render(t, ui.PopoverContent(gsx.Raw("x"), callerAttrs()))
	assertCSSOnlyMarkup(t, got, "popover-content")
	assertCallerAttrsOnce(t, got)
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
	got := render(t, ui.HoverCardContent(gsx.Raw("x"), callerAttrs()))
	assertCSSOnlyMarkup(t, got, "hover-card-content")
	assertCallerAttrsOnce(t, got)
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
	got := render(t, ui.TooltipContent(gsx.Raw("x"), callerAttrs()))
	assertCSSOnlyMarkup(t, got, "tooltip-content", "tooltip-arrow")
	assertCallerAttrsOnce(t, got)
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

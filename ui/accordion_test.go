package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestAccordionPinned(t *testing.T) {
	got := render(t, ui.Accordion("faq", ui.AccordionItem("faq", true, gsx.Fragment(
		ui.AccordionTrigger(gsx.Raw("Question"), nil),
		ui.AccordionContent(gsx.Raw("Answer"), nil),
	), nil), nil))
	want := `<div data-name="faq" data-gsxui-slot-accordion><details name="faq" open data-gsxui-slot-accordion-item><summary data-gsxui-slot-accordion-trigger>Question<svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" data-gsxui-slot-accordion-trigger-icon data-gsxui-slot-icon><path d="m6 9 6 6 6-6"/></svg></summary><div data-gsxui-slot-accordion-content><div data-gsxui-slot-accordion-content-inner>Answer</div></div></details></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestAccordionItemNativeState(t *testing.T) {
	open := render(t, ui.AccordionItem("faq", true, gsx.Raw("x"), nil))
	if !strings.Contains(open, `name="faq" open`) {
		t.Errorf("open item must carry grouped native state\nin: %s", open)
	}
	closed := render(t, ui.AccordionItem("faq", false, gsx.Raw("x"), nil))
	if strings.Contains(closed, " open") {
		t.Errorf("closed item must omit open\nin: %s", closed)
	}
}

func TestAccordionCallerClassRoutesToInnerContent(t *testing.T) {
	got := render(t, ui.AccordionContent(gsx.Raw("x"), gsx.Attrs{
		{Key: "id", Value: "panel"},
		{Key: "class", Value: "pb-8"},
	}))
	want := `<div id="panel" data-gsxui-slot-accordion-content><div class="pb-8" data-gsxui-slot-accordion-content-inner>x</div></div>`
	if got != want {
		t.Errorf("non-class attrs must stay outer and caller class must render once on inner content\n got: %s\nwant: %s", got, want)
	}
}

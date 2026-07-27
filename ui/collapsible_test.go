package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestCollapsiblePinned(t *testing.T) {
	got := render(t, ui.Collapsible(true, gsx.Fragment(
		ui.CollapsibleTrigger(gsx.Raw("Toggle"), nil),
		ui.CollapsibleContent(gsx.Raw("Body"), nil),
	), nil))
	want := `<details open data-gsxui-slot="collapsible"><summary data-gsxui-slot="collapsible-trigger">Toggle</summary><div data-gsxui-slot="collapsible-content">Body</div></details>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCollapsibleClosedDefault(t *testing.T) {
	got := render(t, ui.Collapsible(false, gsx.Raw("x"), nil))
	if strings.Contains(got, " open") {
		t.Errorf("closed collapsible must omit open\nin: %s", got)
	}
}

func TestCollapsibleCallerAttrs(t *testing.T) {
	got := render(t, ui.CollapsibleContent(gsx.Raw("x"), gsx.Attrs{
		{Key: "id", Value: "panel"},
		{Key: "class", Value: "caller"},
	}))
	want := `<div class="caller" data-gsxui-slot="collapsible-content" id="panel">x</div>`
	if got != want {
		t.Errorf("attrs mismatch\n got: %s\nwant: %s", got, want)
	}
}

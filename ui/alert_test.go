package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestAlertStructure(t *testing.T) {
	got := render(t, ui.Alert("", gsx.Fragment(
		ui.AlertTitle(gsx.Raw("Heads up"), nil),
		ui.AlertDescription(gsx.Raw("You can add components here."), nil),
	), nil))
	for _, want := range []string{
		`data-slot="alert"`,
		`role="alert"`,
		`data-slot="alert-title"`, ">Heads up<",
		`data-slot="alert-description"`, ">You can add components here.<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestAlertVariants(t *testing.T) {
	cases := map[string]string{
		"":            "bg-card text-card-foreground",
		"destructive": "bg-card text-destructive",
	}
	for variant, wantClass := range cases {
		got := render(t, ui.Alert(variant, gsx.Raw("x"), nil))
		if !strings.Contains(got, wantClass) {
			t.Errorf("variant %q: missing %q\nin: %s", variant, wantClass, got)
		}
	}
}

func TestAlertCallerClassMerges(t *testing.T) {
	got := render(t, ui.Alert("", gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "px-8"}}))
	if strings.Contains(got, "px-2.5") {
		t.Errorf("base px-2.5 should be dropped by caller px-8\nin: %s", got)
	}
	if !strings.Contains(got, "px-8") {
		t.Errorf("missing caller class px-8\nin: %s", got)
	}
}

func TestAlertPinned(t *testing.T) {
	// Exact full-render pin for the default variant, verified token-by-token
	// against shadcn's alertVariants base + default variant
	// (registry/new-york-v4/ui/alert.tsx) — straight port, cva() replaced by
	// a switch (see docs/jsx-parity.md).
	got := render(t, ui.Alert("", gsx.Raw("Heads up"), nil))
	want := `<div data-slot="alert" role="alert" class="relative grid w-full items-start gap-y-0.5 rounded-lg border px-2.5 py-2 text-sm has-[&gt;svg]:grid-cols-[auto_1fr] has-[&gt;svg]:gap-x-2 *:[svg]:row-span-2 *:[svg:not([class*=&#39;size-&#39;])]:size-4 [&amp;&gt;svg]:translate-y-0.5 [&amp;&gt;svg]:text-current bg-card text-card-foreground">Heads up</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestAlertAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Alert("", gsx.Raw("x"), gsx.Attrs{{Key: "id", Value: "a1"}, {Key: "aria-label", Value: "notice"}}))
	for _, want := range []string{`id="a1"`, `aria-label="notice"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestAlertTitleAndDescriptionAttrsFallThrough(t *testing.T) {
	got := render(t, ui.AlertTitle(gsx.Raw("x"), gsx.Attrs{{Key: "id", Value: "t1"}}))
	if !strings.Contains(got, `id="t1"`) {
		t.Errorf("AlertTitle: missing id fallthrough\nin: %s", got)
	}
	got = render(t, ui.AlertDescription(gsx.Raw("x"), gsx.Attrs{{Key: "id", Value: "d1"}}))
	if !strings.Contains(got, `id="d1"`) {
		t.Errorf("AlertDescription: missing id fallthrough\nin: %s", got)
	}
}

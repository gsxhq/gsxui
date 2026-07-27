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
		`data-gsxui-slot-alert`, `data-variant="default"`,
		`role="alert"`,
		`data-gsxui-slot-alert-title`, ">Heads up<",
		`data-gsxui-slot-alert-description`, ">You can add components here.<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestAlertVariants(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  string
	}{
		{want: "default"},
		{input: "destructive", want: "destructive"},
	} {
		got := render(t, ui.Alert(tc.input, gsx.Raw("x"), nil))
		if !strings.Contains(got, `data-variant="`+tc.want+`"`) {
			t.Errorf("variant %q: missing reflected value %q\nin: %s", tc.input, tc.want, got)
		}
	}
}

func TestAlertCallerClassIsForwardedOnce(t *testing.T) {
	got := render(t, ui.Alert("", gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "px-8"}}))
	if strings.Count(got, `class="px-8"`) != 1 {
		t.Errorf("caller class must be the only class and render once\nin: %s", got)
	}
}

func TestAlertPinned(t *testing.T) {
	// Exact full-render pin for the default variant, verified token-by-token
	// against shadcn's alertVariants base + default variant
	// (registry/new-york-v4/ui/alert.tsx) — straight port, cva() replaced by
	// a switch (see docs/jsx-parity.md).
	got := render(t, ui.Alert("", gsx.Raw("Heads up"), nil))
	want := `<div role="alert" data-variant="default" data-gsxui-slot-alert>Heads up</div>`
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

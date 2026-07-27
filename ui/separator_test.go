package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestSeparatorDefault(t *testing.T) {
	got := render(t, ui.Separator("", nil))
	for _, want := range []string{
		`data-gsxui-slot="separator"`,
		`role="none"`,
		`data-orientation="horizontal"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestSeparatorVertical(t *testing.T) {
	got := render(t, ui.Separator("vertical", nil))
	if !strings.Contains(got, `data-orientation="vertical"`) {
		t.Errorf("missing vertical data-orientation stamp\nin: %s", got)
	}
}

func TestSeparatorCallerClassIsForwardedOnce(t *testing.T) {
	got := render(t, ui.Separator("", gsx.Attrs{{Key: "class", Value: "bg-red-500"}}))
	if strings.Count(got, `class="bg-red-500"`) != 1 {
		t.Errorf("caller class must be the only class and render once\nin: %s", got)
	}
}

func TestSeparatorPinned(t *testing.T) {
	// Exact full-render pin, verified token-by-token against shadcn's
	// Separator (registry/new-york-v4/ui/separator.tsx) and
	// docs/jsx-parity.md — decorative is dropped (ADAPT), role="none" always.
	got := render(t, ui.Separator("", nil))
	want := `<div role="none" data-orientation="horizontal" data-gsxui-slot="separator"></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestSeparatorAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Separator("", gsx.Attrs{{Key: "id", Value: "s1"}, {Key: "aria-hidden", Value: "true"}}))
	for _, want := range []string{`id="s1"`, `aria-hidden="true"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

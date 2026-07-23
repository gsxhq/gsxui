package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestLabelDefault(t *testing.T) {
	got := render(t, ui.Label(gsx.Raw("Email"), nil))
	for _, want := range []string{
		"<label", `data-slot="label"`,
		"flex items-center gap-2 text-sm leading-none font-medium select-none",
		"peer-disabled:cursor-not-allowed peer-disabled:opacity-50",
		">Email</label>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestLabelPinned(t *testing.T) {
	// Exact full-render pin, verified token-by-token against shadcn's Label
	// (registry/new-york-v4/ui/label.tsx) and docs/jsx-parity.md — straight
	// port, no ADAPT deviations.
	got := render(t, ui.Label(gsx.Raw("Email"), nil))
	want := `<label data-slot="label" class="flex items-center gap-2 text-sm leading-none font-medium select-none group-data-[disabled=true]:pointer-events-none group-data-[disabled=true]:opacity-50 peer-disabled:cursor-not-allowed peer-disabled:opacity-50">Email</label>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestLabelCallerClassMerges(t *testing.T) {
	got := render(t, ui.Label(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "text-lg"}}))
	if strings.Contains(got, "text-sm") {
		t.Errorf("caller text-lg must drop default text-sm\nin: %s", got)
	}
	for _, want := range []string{"text-lg", "flex items-center"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestLabelAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Label(gsx.Raw("x"), gsx.Attrs{{Key: "for", Value: "email"}, {Key: "id", Value: "email-label"}}))
	for _, want := range []string{`for="email"`, `id="email-label"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

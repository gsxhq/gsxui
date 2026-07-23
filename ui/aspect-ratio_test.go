package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestAspectRatioDefault(t *testing.T) {
	got := render(t, ui.AspectRatio("16 / 9", gsx.Raw(`<img src="x.png"/>`), nil))
	for _, want := range []string{
		`data-slot="aspect-ratio"`,
		`style="aspect-ratio: 16 / 9"`,
		`<img src="x.png"/>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestAspectRatioNumericRatio(t *testing.T) {
	// aspect-ratio also accepts a bare number, not only the "w / h" form.
	got := render(t, ui.AspectRatio("1.5", nil, nil))
	if !strings.Contains(got, `style="aspect-ratio: 1.5"`) {
		t.Errorf("missing numeric ratio style\nin: %s", got)
	}
}

func TestAspectRatioAttrsFallThrough(t *testing.T) {
	got := render(t, ui.AspectRatio("16 / 9", nil, gsx.Attrs{{Key: "id", Value: "ar1"}, {Key: "class", Value: "bg-muted"}}))
	if !strings.Contains(got, `id="ar1"`) {
		t.Errorf("missing fallthrough id\nin: %s", got)
	}
	if !strings.Contains(got, `class="bg-muted"`) {
		t.Errorf("missing fallthrough class\nin: %s", got)
	}
}

func TestAspectRatioPinned(t *testing.T) {
	// Exact full-render pin. shadcn's AspectRatio is a bare passthrough onto
	// Radix's padding-hack Root (registry/new-york-v4/ui/aspect-ratio.tsx);
	// this port replaces the two-div padding-percentage mechanism with the
	// CSS aspect-ratio property directly on a single div (ADAPT, see
	// docs/jsx-parity.md).
	got := render(t, ui.AspectRatio("16 / 9", gsx.Raw(`<img src="x.png"/>`), nil))
	want := `<div data-slot="aspect-ratio" style="aspect-ratio: 16 / 9"><img src="x.png"/></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

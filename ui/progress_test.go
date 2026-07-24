package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestProgressDefault(t *testing.T) {
	got := render(t, ui.Progress(60, nil))
	for _, want := range []string{
		`<div data-slot="progress"`,
		`role="progressbar"`,
		`aria-valuemin="0"`,
		`aria-valuemax="100"`,
		`aria-valuenow="60"`,
		"relative h-1 w-full overflow-hidden rounded-full bg-primary/20",
		`data-slot="progress-indicator"`,
		"h-full w-full flex-1 bg-primary transition-all",
		`style="transform: translateX(-40%)"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// The Go zero value for value (0) matches shadcn's own `value || 0`
// fallback: an unset Progress renders its indicator fully translated
// off-screen, not a JS-undefined NaN transform.
func TestProgressZeroValue(t *testing.T) {
	got := render(t, ui.Progress(0, nil))
	if !strings.Contains(got, `style="transform: translateX(-100%)"`) {
		t.Errorf("missing zero-value indicator transform\nin: %s", got)
	}
	if !strings.Contains(got, `aria-valuenow="0"`) {
		t.Errorf("missing aria-valuenow=0\nin: %s", got)
	}
}

func TestProgressCallerClassMerges(t *testing.T) {
	got := render(t, ui.Progress(50, gsx.Attrs{{Key: "class", Value: "h-4"}}))
	if strings.Contains(got, "h-1 w-full") {
		t.Errorf("base h-1 should be dropped by caller h-4\nin: %s", got)
	}
	if !strings.Contains(got, "h-4") {
		t.Errorf("missing caller class h-4\nin: %s", got)
	}
}

func TestProgressAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Progress(50, gsx.Attrs{{Key: "id", Value: "p1"}, {Key: "aria-label", Value: "Upload progress"}}))
	for _, want := range []string{`id="p1"`, `aria-label="Upload progress"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestProgressPinned(t *testing.T) {
	// Exact full-render pin. Token-for-token against shadcn's Progress
	// (registry/new-york-v4/ui/progress.tsx): the class strings on both divs
	// are carried verbatim; Radix's Root/Indicator primitives are replaced by
	// role="progressbar" + aria-valuemin/max/now (ADAPT, see
	// docs/jsx-parity.md) and the translateX(-{100-value}%) transform is
	// ported faithfully via gsx.RawCSS.
	got := render(t, ui.Progress(25, nil))
	want := `<div data-slot="progress" role="progressbar" aria-valuemin="0" aria-valuemax="100" aria-valuenow="25" class="relative h-1 w-full overflow-hidden rounded-full bg-primary/20"><div data-slot="progress-indicator" class="h-full w-full flex-1 bg-primary transition-all" style="transform: translateX(-75%)"></div></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

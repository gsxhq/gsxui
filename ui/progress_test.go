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
		`<div role="progressbar"`,
		`role="progressbar"`,
		`aria-valuemin="0"`,
		`aria-valuemax="100"`,
		`aria-valuenow="60"`,
		`data-gsxui-slot-progress`,
		`data-gsxui-slot-progress-indicator`,
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
	if strings.Count(got, `class="`) != 2 || !strings.Contains(got, "h-4") || !strings.Contains(got, "bg-muted") {
		t.Errorf("caller class must merge into the root's class attribute and render once\nin: %s", got)
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
	// Presentation is expressed as the recipe's resolved class (slot axis
	// migration); only the dynamic transform remains inline.
	got := render(t, ui.Progress(25, nil))
	want := `<div role="progressbar" aria-valuemin="0" aria-valuemax="100" aria-valuenow="25" class="bg-muted h-1 rounded-full" data-gsxui-slot-progress><div style="transform: translateX(-75%)" class="bg-primary" data-gsxui-slot-progress-indicator></div></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

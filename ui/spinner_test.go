package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestSpinnerDefault(t *testing.T) {
	got := render(t, ui.Spinner(nil))
	for _, want := range []string{
		"<svg",
		`role="status"`,
		`aria-label="Loading"`,
		"size-4 animate-spin",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// The spinner must NOT be aria-hidden — it is an announced loading status,
// unlike every other icon.* usage in this package (chevrons etc., which
// keep svgIcon's aria-hidden="true" default).
func TestSpinnerNotAriaHidden(t *testing.T) {
	got := render(t, ui.Spinner(nil))
	if strings.Contains(got, `aria-hidden="true"`) {
		t.Errorf("spinner should not default to aria-hidden=true\nin: %s", got)
	}
	if !strings.Contains(got, `aria-hidden="false"`) {
		t.Errorf("missing aria-hidden=false\nin: %s", got)
	}
}

func TestSpinnerCallerClassMerges(t *testing.T) {
	got := render(t, ui.Spinner(gsx.Attrs{{Key: "class", Value: "size-6"}}))
	if strings.Contains(got, "size-4") {
		t.Errorf("base size-4 should be dropped by caller size-6\nin: %s", got)
	}
	if !strings.Contains(got, "animate-spin") {
		t.Errorf("animate-spin should survive the merge\nin: %s", got)
	}
	if !strings.Contains(got, "size-6") {
		t.Errorf("missing caller class size-6\nin: %s", got)
	}
}

func TestSpinnerAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Spinner(gsx.Attrs{{Key: "id", Value: "s1"}}))
	if !strings.Contains(got, `id="s1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// Caller-supplied role/aria-label must win over Spinner's own defaults —
// the same overridable-default idiom as button's type="button".
func TestSpinnerRoleAndLabelOverridable(t *testing.T) {
	got := render(t, ui.Spinner(gsx.Attrs{{Key: "role", Value: "presentation"}, {Key: "aria-label", Value: "Working"}}))
	if !strings.Contains(got, `role="presentation"`) {
		t.Errorf("caller role should win\nin: %s", got)
	}
	if strings.Contains(got, `role="status"`) {
		t.Errorf("default role=status should not also render\nin: %s", got)
	}
	if !strings.Contains(got, `aria-label="Working"`) {
		t.Errorf("caller aria-label should win\nin: %s", got)
	}
}

func TestSpinnerPinned(t *testing.T) {
	// Exact full-render pin. Verified against shadcn's Spinner
	// (registry/new-york-v4/ui/spinner.tsx: Loader2Icon role="status"
	// aria-label="Loading" class="size-4 animate-spin") rendered through
	// ui/icon's svgIcon wrapper (data-slot/viewBox/stroke defaults), with
	// aria-hidden explicitly overridden to "false" — see spinner.gsx.
	// aria-hidden overrides svgIcon's default "true" via the spread bag
	// rather than in svgIcon's own literal position — per gsx's documented
	// rule ("Scalar attributes from a spread keep the spread's source
	// position", docs/guide/syntax/styling.md), the winning value renders
	// once, at the spread's position (after class, alongside role and
	// aria-label — the other attrs threaded through the same bag), not
	// where svgIcon's own aria-hidden="true" default was authored.
	got := render(t, ui.Spinner(nil))
	want := `<svg data-slot="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4 animate-spin" role="status" aria-label="Loading" aria-hidden="false"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

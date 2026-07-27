package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestToggleOffPinned(t *testing.T) {
	// Exact full-render pin for the zero-value (unpressed, default variant,
	// default size) render — token-by-token against toggleVariants' base +
	// default variant + default size (registry/new-york-v4/ui/toggle.tsx).
	got := render(t, ui.Toggle(false, "", "", gsx.Raw("Bold"), nil))
	want := `<button type="button" data-gsxui-toggle data-variant="default" data-size="default" data-state="off" aria-pressed="false" data-gsxui-slot-toggle>Bold</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestTogglePressedPinned(t *testing.T) {
	// Exact full-render pin for pressed={true} — the server-visible initial
	// "on" state (aria-pressed="true" data-state="on"), no click required.
	got := render(t, ui.Toggle(true, "", "", gsx.Raw("Bold"), nil))
	want := `<button type="button" data-gsxui-toggle data-variant="default" data-size="default" data-state="on" aria-pressed="true" data-gsxui-slot-toggle>Bold</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestToggleOutlineVariant(t *testing.T) {
	got := render(t, ui.Toggle(false, "outline", "", gsx.Raw("x"), nil))
	for _, want := range []string{
		`data-variant="outline"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestToggleSizes(t *testing.T) {
	sm := render(t, ui.Toggle(false, "", "sm", gsx.Raw("x"), nil))
	for _, want := range []string{`data-size="sm"`} {
		if !strings.Contains(sm, want) {
			t.Errorf("sm missing %q\nin: %s", want, sm)
		}
	}

	lg := render(t, ui.Toggle(false, "", "lg", gsx.Raw("x"), nil))
	for _, want := range []string{`data-size="lg"`} {
		if !strings.Contains(lg, want) {
			t.Errorf("lg missing %q\nin: %s", want, lg)
		}
	}

	def := render(t, ui.Toggle(false, "", "default", gsx.Raw("x"), nil))
	for _, want := range []string{`data-size="default"`} {
		if !strings.Contains(def, want) {
			t.Errorf("default missing %q\nin: %s", want, def)
		}
	}
}

func TestToggleDisabledFallsThrough(t *testing.T) {
	// disabled is not a declared param — it flows through attrs like any
	// other plain boolean HTML attribute (no href/disabled interplay to
	// resolve server-side the way Button has).
	got := render(t, ui.Toggle(false, "", "", gsx.Raw("x"), gsx.Attrs{{Key: "disabled", Value: true}}))
	if !strings.Contains(got, "disabled") {
		t.Errorf("want disabled attribute\nin: %s", got)
	}
}

func TestToggleAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Toggle(false, "", "", gsx.Raw("x"), gsx.Attrs{{Key: "aria-label", Value: "Toggle bold"}}))
	if !strings.Contains(got, `aria-label="Toggle bold"`) {
		t.Errorf("missing aria-label\nin: %s", got)
	}
}

func TestToggleCallerClassMerges(t *testing.T) {
	got := render(t, ui.Toggle(false, "", "", gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "h-12"}}))
	if strings.Count(got, `class="h-12"`) != 1 {
		t.Errorf("caller class must be forwarded exactly once\nin: %s", got)
	}
}

package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

// canonicalSliderClass is Slider's fully-resolved recipe class, copied
// verbatim from the generated ui/slider.gsx output — Slider migrated onto
// the slot axis, so its root class is no longer empty.
const canonicalSliderClass = `w-full cursor-pointer bg-transparent disabled:cursor-not-allowed disabled:opacity-50 [&::-webkit-slider-runnable-track]:h-1 [&::-webkit-slider-runnable-track]:rounded-full [&::-webkit-slider-runnable-track]:bg-[linear-gradient(to_right,var(--primary)_0%,var(--primary)_var(--fill,0%),var(--muted)_var(--fill,0%),var(--muted)_100%)] rtl:[&::-webkit-slider-runnable-track]:bg-[linear-gradient(to_left,var(--primary)_0%,var(--primary)_var(--fill,0%),var(--muted)_var(--fill,0%),var(--muted)_100%)] [&::-moz-range-track]:h-1 [&::-moz-range-track]:rounded-full [&::-moz-range-track]:bg-[linear-gradient(to_right,var(--primary)_0%,var(--primary)_var(--fill,0%),var(--muted)_var(--fill,0%),var(--muted)_100%)] rtl:[&::-moz-range-track]:bg-[linear-gradient(to_left,var(--primary)_0%,var(--primary)_var(--fill,0%),var(--muted)_var(--fill,0%),var(--muted)_100%)] [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:size-3 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:border [&::-webkit-slider-thumb]:border-primary [&::-webkit-slider-thumb]:bg-contrast [&::-webkit-slider-thumb]:-mt-1 [&::-webkit-slider-thumb]:transition-shadow [&::-webkit-slider-thumb]:duration-150 [&::-moz-range-thumb]:size-3 [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:border [&::-moz-range-thumb]:border-primary [&::-moz-range-thumb]:bg-contrast [&::-moz-range-thumb]:transition-shadow [&::-moz-range-thumb]:duration-150 hover:[&::-webkit-slider-thumb]:shadow-[0_0_0_3px_color-mix(in_oklab,var(--ring)_50%,transparent)] focus-visible:[&::-webkit-slider-thumb]:shadow-[0_0_0_3px_color-mix(in_oklab,var(--ring)_50%,transparent)] active:[&::-webkit-slider-thumb]:shadow-[0_0_0_3px_color-mix(in_oklab,var(--ring)_50%,transparent)] hover:[&::-moz-range-thumb]:shadow-[0_0_0_3px_color-mix(in_oklab,var(--ring)_50%,transparent)] focus-visible:[&::-moz-range-thumb]:shadow-[0_0_0_3px_color-mix(in_oklab,var(--ring)_50%,transparent)] active:[&::-moz-range-thumb]:shadow-[0_0_0_3px_color-mix(in_oklab,var(--ring)_50%,transparent)]`

func TestSliderPinned(t *testing.T) {
	// Exact full-render pin: value=50 min=0(zero-value) max=100 step=1 —
	// the shadcn slider-demo shape (defaultValue={[50]} max={100} step={1}).
	// --fill is server-computed exact arithmetic: (50-0)/(100-0)*100 = 50.
	got := render(t, ui.Slider(50, 0, 100, 1, nil))
	// render() HTML-escapes attribute values, so "&" in the arbitrary-variant
	// class list becomes "&amp;" in the rendered output.
	want := `<input type="range" min="0" max="100" step="1" value="50" style="--fill: 50%" class="` +
		strings.ReplaceAll(canonicalSliderClass, "&", "&amp;") + `" data-gsxui-slot-slider>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestSliderMinMaxStepAttrs(t *testing.T) {
	got := render(t, ui.Slider(25, 0, 100, 1, nil))
	for _, want := range []string{`min="0"`, `max="100"`, `step="1"`, `value="25"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestSliderFillNonZeroMin(t *testing.T) {
	// min != 0: exact-arithmetic fill must account for the offset, not
	// treat min as 0. (25-20)/(40-20)*100 = 25%.
	got := render(t, ui.Slider(25, 20, 40, 1, nil))
	for _, want := range []string{`min="20"`, `max="40"`, `style="--fill: 25%"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestSliderFillAtMin(t *testing.T) {
	got := render(t, ui.Slider(0, 0, 100, 1, nil))
	if !strings.Contains(got, `style="--fill: 0%"`) {
		t.Errorf("missing 0%% fill\nin: %s", got)
	}
}

func TestSliderFillAtMax(t *testing.T) {
	got := render(t, ui.Slider(100, 0, 100, 1, nil))
	if !strings.Contains(got, `style="--fill: 100%"`) {
		t.Errorf("missing 100%% fill\nin: %s", got)
	}
}

func TestSliderMaxStepZeroValueDefaults(t *testing.T) {
	// max/step left at the Go zero value fall back to shadcn's own
	// defaults (100/1) — see slider.gsx's own doc comment on this
	// unset-vs-explicit-zero ambiguity.
	got := render(t, ui.Slider(50, 0, 0, 0, nil))
	for _, want := range []string{`max="100"`, `step="1"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestSliderDisabledFallsThrough(t *testing.T) {
	got := render(t, ui.Slider(50, 0, 100, 1, gsx.Attrs{{Key: "disabled", Value: true}}))
	if !strings.Contains(got, "disabled") {
		t.Errorf("want disabled attribute\nin: %s", got)
	}
}

func TestSliderAttrsFallThrough(t *testing.T) {
	got := render(t, ui.Slider(50, 0, 100, 1, gsx.Attrs{{Key: "aria-label", Value: "Volume"}, {Key: "name", Value: "volume"}}))
	for _, want := range []string{`aria-label="Volume"`, `name="volume"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestSliderCallerClassMerges(t *testing.T) {
	got := render(t, ui.Slider(50, 0, 100, 1, gsx.Attrs{{Key: "class", Value: "w-[60%]"}}))
	if strings.Count(got, "w-[60%]") != 1 || strings.Count(got, `class="`) != 1 {
		t.Errorf("caller class must be forwarded exactly once, merged with the recipe class\nin: %s", got)
	}
}

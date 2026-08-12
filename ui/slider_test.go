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
// w-full/cursor-pointer/bg-transparent/disabled: are restored "carried: no
// upstream counterpart" structural chrome; the flat bg-muted track fill was
// replaced with the real bg-[linear-gradient(...)] fill + rtl: mirrored
// counterpart the recipe needs to render a filled-vs-unfilled track from a
// single <input type="range"> — see the style-porter report's "Slider —
// full structural rewrite" entry.
const canonicalSliderClass = `w-full cursor-pointer bg-transparent disabled:cursor-not-allowed disabled:opacity-50 data-vertical:min-h-40 [&::-webkit-slider-runnable-track]:bg-[linear-gradient(to_right,var(--primary)_0%,var(--primary)_var(--fill,0%),var(--muted)_var(--fill,0%),var(--muted)_100%)] rtl:[&::-webkit-slider-runnable-track]:bg-[linear-gradient(to_left,var(--primary)_0%,var(--primary)_var(--fill,0%),var(--muted)_var(--fill,0%),var(--muted)_100%)] [&::-webkit-slider-runnable-track]:rounded-full data-horizontal:[&::-webkit-slider-runnable-track]:h-1 data-horizontal:[&::-webkit-slider-runnable-track]:w-full data-vertical:[&::-webkit-slider-runnable-track]:h-full data-vertical:[&::-webkit-slider-runnable-track]:w-1 [&::-moz-range-track]:bg-[linear-gradient(to_right,var(--primary)_0%,var(--primary)_var(--fill,0%),var(--muted)_var(--fill,0%),var(--muted)_100%)] rtl:[&::-moz-range-track]:bg-[linear-gradient(to_left,var(--primary)_0%,var(--primary)_var(--fill,0%),var(--muted)_var(--fill,0%),var(--muted)_100%)] [&::-moz-range-track]:rounded-full data-horizontal:[&::-moz-range-track]:h-1 data-horizontal:[&::-moz-range-track]:w-full data-vertical:[&::-moz-range-track]:h-full data-vertical:[&::-moz-range-track]:w-1 [&::-webkit-slider-thumb]:border-ring [&::-webkit-slider-thumb]:ring-ring/50 [&::-webkit-slider-thumb]:relative [&::-webkit-slider-thumb]:size-3 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:border [&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:transition-[color,box-shadow] after:[&::-webkit-slider-thumb]:absolute after:[&::-webkit-slider-thumb]:-inset-2 hover:[&::-webkit-slider-thumb]:ring-3 focus-visible:[&::-webkit-slider-thumb]:ring-3 focus-visible:[&::-webkit-slider-thumb]:outline-hidden active:[&::-webkit-slider-thumb]:ring-3 [&::-moz-range-thumb]:border-ring [&::-moz-range-thumb]:ring-ring/50 [&::-moz-range-thumb]:relative [&::-moz-range-thumb]:size-3 [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:border [&::-moz-range-thumb]:bg-white [&::-moz-range-thumb]:transition-[color,box-shadow] after:[&::-moz-range-thumb]:absolute after:[&::-moz-range-thumb]:-inset-2 hover:[&::-moz-range-thumb]:ring-3 focus-visible:[&::-moz-range-thumb]:ring-3 focus-visible:[&::-moz-range-thumb]:outline-hidden active:[&::-moz-range-thumb]:ring-3`

func TestSliderPinned(t *testing.T) {
	// Exact full-render pin: value=50 min=0(zero-value) max=100 step=1 —
	// the shadcn slider-demo shape (defaultValue={[50]} max={100} step={1}).
	// --fill is server-computed exact arithmetic: (50-0)/(100-0)*100 = 50.
	got := render(t, ui.Slider(50, 0, 100, 1, nil))
	// render() HTML-escapes attribute values, so "&" in the arbitrary-variant
	// class list becomes "&amp;" in the rendered output.
	want := `<input type="range" min="0" max="100" step="1" value="50" style="--fill: 50%" class="` +
		strings.ReplaceAll(canonicalSliderClass, "&", "&amp;") + `" data-horizontal data-gsxui-slot-slider>`
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

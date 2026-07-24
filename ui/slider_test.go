package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestSliderPinned(t *testing.T) {
	// Exact full-render pin: value=50 min=0(zero-value) max=100 step=1 —
	// the shadcn slider-demo shape (defaultValue={[50]} max={100} step={1}).
	// --fill is server-computed exact arithmetic: (50-0)/(100-0)*100 = 50.
	got := render(t, ui.Slider(50, 0, 100, 1, nil))
	want := `<input type="range" data-slot="slider" data-gsxui-slider min="0" max="100" step="1" value="50" style="--fill: 50%" class="appearance-none bg-transparent w-full cursor-pointer outline-none disabled:cursor-not-allowed disabled:opacity-50">`
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
	if strings.Contains(got, "w-full") {
		t.Errorf("caller w-[60%%] must drop default w-full\nin: %s", got)
	}
	if !strings.Contains(got, "w-[60%]") || !strings.Contains(got, "appearance-none") {
		t.Errorf("want w-[60%%] plus surviving structural classes\nin: %s", got)
	}
}

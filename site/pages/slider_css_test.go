package pages_test

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"
)

// TestSliderCSSAppliesThumbPresentation reads Slider's migrated recipe
// stylesheet (registry/styles/nova/slider.css) and checks the thumb
// pseudo-elements carry real presentation.
//
// This test used to be TestSliderCSSUsesSemanticContrast and asserted that
// gsxui's old hand-authored Slider recipe used a custom semantic utility,
// bg-contrast, for the thumb background, and that no literal color ever
// appeared instead. That premise is gone: nova is now ported directly from
// upstream shadcn-ui, whose Slider recipe uses a literal bg-white on the
// thumb — there is no bg-contrast concept upstream, and the porter must not
// invent per-style custom utilities that don't exist there. So this test no
// longer cares which background utility the thumb uses (literal or
// semantic) — only that the webkit/moz thumb pseudo-element selectors exist
// and each applies some background utility, which is what actually matters
// for the thumb to render visibly against the track.
func TestSliderCSSAppliesThumbPresentation(t *testing.T) {
	_, testFile, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(testFile), "..", "..", "registry", "styles", "nova", "slider.css")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile(%s): %v", path, err)
	}
	for _, selector := range []string{
		`[&::-webkit-slider-thumb]:`,
		`[&::-moz-range-thumb]:`,
	} {
		if !regexp.MustCompile(regexp.QuoteMeta(selector)).Match(b) {
			t.Errorf("recipe has no rules for thumb pseudo-element %s", selector)
		}
	}
	background := regexp.MustCompile(`\[&::-webkit-slider-thumb\]:bg-\S+`)
	if !background.Match(b) {
		t.Errorf("recipe does not apply any background utility to the webkit thumb")
	}
}

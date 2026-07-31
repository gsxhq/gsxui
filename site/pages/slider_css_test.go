package pages_test

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"
)

// TestSliderCSSUsesSemanticContrast now reads Slider's migrated recipe
// stylesheet (registry/styles/nova/slider.css) rather than the deleted
// assets/css/styles/default/slider.css block — the same bg-contrast utility
// carries the semantic token forward through the arbitrary-variant form
// (`[&::-webkit-slider-thumb]:bg-contrast`), not a raw `background:` property.
func TestSliderCSSUsesSemanticContrast(t *testing.T) {
	_, testFile, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(testFile), "..", "..", "registry", "styles", "nova", "slider.css")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile(%s): %v", path, err)
	}
	for _, selector := range []string{
		`[&::-webkit-slider-thumb]:bg-contrast`,
		`[&::-moz-range-thumb]:bg-contrast`,
	} {
		if !regexp.MustCompile(regexp.QuoteMeta(selector)).Match(b) {
			t.Errorf("recipe does not apply %s", selector)
		}
	}
	literalWhite := regexp.MustCompile(`bg-white\b|#fff`)
	if literalWhite.Match(b) {
		t.Errorf("recipe bypasses the contrast token with a literal white background")
	}
}

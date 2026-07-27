package pages_test

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"
)

var sliderThumbContrastRe = regexp.MustCompile(`(?s)input\[type="range"\]::-(?:webkit-slider-thumb|moz-range-thumb) \{[^}]*background: var\(--contrast\);`)

func TestSliderCSSUsesSemanticContrast(t *testing.T) {
	_, testFile, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(testFile), "..", "..", "assets", "css", "styles", "default.css")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile(%s): %v", path, err)
	}
	if got := len(sliderThumbContrastRe.FindAll(b, -1)); got != 2 {
		t.Errorf("semantic Slider thumb backgrounds = %d, want 2", got)
	}
}

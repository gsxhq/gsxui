package pages_test

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"
)

func TestSliderCSSUsesSemanticContrast(t *testing.T) {
	_, testFile, _, _ := runtime.Caller(0)
	// Slider's authored presentation is its own file: the default style's
	// components layer is split one file per component so that migrating a
	// component is a file deletion rather than an edit to one shared 3000-line
	// sheet. When Slider migrates this file goes away and this test fails loudly
	// rather than passing against a stale block.
	path := filepath.Join(filepath.Dir(testFile), "..", "..", "assets", "css", "styles", "default", "slider.css")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile(%s): %v", path, err)
	}
	for _, selector := range []string{
		`:where([data-gsxui-slot-slider])::-webkit-slider-thumb`,
		`:where([data-gsxui-slot-slider])::-moz-range-thumb`,
	} {
		contrast := regexp.MustCompile(regexp.QuoteMeta(selector) + `\s*\{[^}]*background:\s*var\(--contrast\);`)
		if !contrast.Match(b) {
			t.Errorf("%s does not use background: var(--contrast)", selector)
		}
		literalWhite := regexp.MustCompile(regexp.QuoteMeta(selector) + `\s*\{[^}]*background:\s*#fff;`)
		if literalWhite.Match(b) {
			t.Errorf("%s bypasses the contrast token with a literal white background", selector)
		}
	}
}

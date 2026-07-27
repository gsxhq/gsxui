package pages_test

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"
)

// accordionCSSRe extracts the accordion animation block: from its marker
// comment through the namespaced accordion-content rule that closes it.
var accordionCSSRe = regexp.MustCompile(`(?s)/\* Accordion open/close animation.*?\[data-gsxui-slot~="accordion-content"\]\) \{\n  min-height: 0;\n\}`)

// TestAccordionAnimationCSSDriftPin ensures the library's single authored
// foundation contains the native <details> open/close animation once.
func TestAccordionAnimationCSSDriftPin(t *testing.T) {
	_, testFile, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(testFile)
	path := filepath.Join(testDir, "..", "..", "assets", "css", "foundation.css")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile(%s): %v", path, err)
	}
	if got := len(accordionCSSRe.FindAll(b, -1)); got != 1 {
		t.Errorf("accordion animation block count = %d, want 1", got)
	}
}

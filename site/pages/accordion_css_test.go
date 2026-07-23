package pages_test

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"
)

// accordionCSSRe extracts the accordion animation block: from its marker
// comment through the [data-slot="accordion-content"] rule that closes it.
var accordionCSSRe = regexp.MustCompile(`(?s)/\* Accordion open/close animation.*?\[data-slot="accordion-content"\] \{\n  min-height: 0;\n\}`)

// TestAccordionAnimationCSSDriftPin ensures the accordion animation block is
// present in BOTH assets/gsxui.css (what `gsxui init` vendors to consumers)
// and web/site.css (the site's copied-not-imported twin) and byte-identical
// between them — the same copied-block contract TestThemeDefaultsDriftPin
// pins for the theme tokens. The block is what animates the native
// <details> open/close via ::details-content grid-row 0fr->1fr; losing it
// from either file silently reverts that file's consumers to instant
// toggling.
func TestAccordionAnimationCSSDriftPin(t *testing.T) {
	_, testFile, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(testFile)

	extract := func(rel string) string {
		t.Helper()
		path := filepath.Join(testDir, "..", "..", rel)
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("os.ReadFile(%s): %v", path, err)
		}
		m := accordionCSSRe.Find(b)
		if m == nil {
			t.Fatalf("accordion animation block missing from %s", rel)
		}
		return string(m)
	}

	site := extract("web/site.css")
	assets := extract("assets/gsxui.css")
	if site != assets {
		t.Errorf("accordion animation block drifted between web/site.css and assets/gsxui.css\nsite.css:\n%s\n\ngsxui.css:\n%s", site, assets)
	}
}

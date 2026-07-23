package pages_test

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/site/pages"
)

// TestThemeDefaultsDriftPin ensures the Go themeGroups defaults stay in sync
// with web/site.css's :root and .dark blocks byte-for-byte. This pins the
// Go-map ↔ CSS sync in CI.
func TestThemeDefaultsDriftPin(t *testing.T) {
	// Get the directory of this test file
	_, testFile, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(testFile)

	// Read web/site.css relative to this test file (../../web/site.css)
	cssPath := filepath.Join(testDir, "..", "..", "web", "site.css")

	cssBytes, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("os.ReadFile(%s): %v", cssPath, err)
	}
	cssText := string(cssBytes)

	// Extract :root and .dark blocks using regex
	rootBlock := extractCSSBlock(cssText, ":root")
	darkBlock := extractCSSBlock(cssText, ".dark")

	if rootBlock == "" {
		t.Fatal("failed to extract :root block from site.css")
	}
	if darkBlock == "" {
		t.Fatal("failed to extract .dark block from site.css")
	}

	// Parse the CSS blocks into maps of var -> value
	cssVars := map[string]map[string]string{
		"light": parseCSSVars(rootBlock),
		"dark":  parseCSSVars(darkBlock),
	}

	// Build the Go defaults map from themeGroups
	goDefaults := map[string]map[string]string{
		"light": {},
		"dark":  {},
	}

	// Iterate through all theme vars in the Go definitions
	for _, g := range pages.ThemeGroups() {
		for _, v := range g.Vars {
			goDefaults["light"][v.Name] = v.Light
			goDefaults["dark"][v.Name] = v.Dark
		}
	}

	// Assert byte-for-byte match for light mode
	for varName, goValue := range goDefaults["light"] {
		cssValue, ok := cssVars["light"][varName]
		if !ok {
			t.Errorf("light mode: %s missing in CSS :root block", varName)
			continue
		}
		if cssValue != goValue {
			t.Errorf("light mode: %s = %q (Go) vs %q (CSS) mismatch", varName, goValue, cssValue)
		}
	}

	// Assert byte-for-byte match for dark mode
	// Note: --radius is only in :root, not in .dark
	for varName, goValue := range goDefaults["dark"] {
		if varName == "--radius" {
			// --radius is theme-invariant, so it's only in :root
			continue
		}
		cssValue, ok := cssVars["dark"][varName]
		if !ok {
			t.Errorf("dark mode: %s missing in CSS .dark block", varName)
			continue
		}
		if cssValue != goValue {
			t.Errorf("dark mode: %s = %q (Go) vs %q (CSS) mismatch", varName, goValue, cssValue)
		}
	}

	// Also verify that all CSS vars are accounted for in the Go defaults
	for varName := range cssVars["light"] {
		if _, ok := goDefaults["light"][varName]; !ok && varName != "--radius" {
			t.Errorf("CSS :root contains %s which is not in Go defaults", varName)
		}
	}
	for varName := range cssVars["dark"] {
		if _, ok := goDefaults["dark"][varName]; !ok {
			t.Errorf("CSS .dark contains %s which is not in Go defaults", varName)
		}
	}
}

// extractCSSBlock extracts the content of a CSS block (e.g., ":root { ... }" or ".dark { ... }")
func extractCSSBlock(cssText, selector string) string {
	// Match the selector followed by { ... } (greedy match for the closing brace)
	pattern := regexp.MustCompile(regexp.QuoteMeta(selector) + `\s*\{([^}]*)\}`)
	match := pattern.FindStringSubmatch(cssText)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

// parseCSSVars parses CSS variable declarations from a block of text
func parseCSSVars(blockText string) map[string]string {
	result := make(map[string]string)
	if blockText == "" {
		return result
	}

	// Match --var-name: value; pairs
	// The value can span multiple tokens and may contain slashes, parentheses, etc.
	pattern := regexp.MustCompile(`--([a-zA-Z0-9-]+)\s*:\s*([^;]+);`)
	matches := pattern.FindAllStringSubmatch(blockText, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			varName := "--" + match[1]
			value := strings.TrimSpace(match[2])
			result[varName] = value
		}
	}

	return result
}

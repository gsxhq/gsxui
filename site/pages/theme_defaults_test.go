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

// TestThemeDefaultsDriftPin ensures the Go themeGroups defaults and the
// library-only theme tokens stay in sync with the authored default theme.
func TestThemeDefaultsDriftPin(t *testing.T) {
	// Get the directory of this test file
	_, testFile, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(testFile)

	cssPath := filepath.Join(testDir, "..", "..", "assets", "css", "themes", "default.css")

	cssBytes, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("os.ReadFile(%s): %v", cssPath, err)
	}
	cssText := string(cssBytes)

	// Extract :root and .dark blocks using regex
	rootBlock := extractCSSBlock(cssText, ":root")
	darkBlock := extractCSSBlock(cssText, ".dark")

	if rootBlock == "" {
		t.Fatal("failed to extract :root block from default.css")
	}
	if darkBlock == "" {
		t.Fatal("failed to extract .dark block from default.css")
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

	libraryOnly := map[string]map[string]string{
		"light": {
			"--success":  "oklch(69.6% 0.17 162.48)",
			"--info":     "oklch(68.5% 0.169 237.323)",
			"--warning":  "oklch(76.9% 0.188 70.08)",
			"--overlay":  "oklch(0% 0 0 / 10%)",
			"--contrast": "oklch(100% 0 0)",
		},
		"dark": {
			"--success":  "oklch(69.6% 0.17 162.48)",
			"--info":     "oklch(68.5% 0.169 237.323)",
			"--warning":  "oklch(76.9% 0.188 70.08)",
			"--overlay":  "oklch(0% 0 0 / 10%)",
			"--contrast": "oklch(100% 0 0)",
		},
	}
	for mode, expected := range libraryOnly {
		for name, want := range expected {
			if got := cssVars[mode][name]; got != want {
				t.Errorf("%s mode: %s = %q, want %q", mode, name, got, want)
			}
		}
	}

}

func TestThemeDefaultsExposeEverySidebarTokenInTheEditor(t *testing.T) {
	want := map[string][2]string{
		"--sidebar":                    {"oklch(0.985 0 0)", "oklch(0.205 0 0)"},
		"--sidebar-foreground":         {"oklch(0% 0 0)", "oklch(0.985 0 0)"},
		"--sidebar-primary":            {"oklch(0.205 0 0)", "oklch(0.488 0.243 264.376)"},
		"--sidebar-primary-foreground": {"oklch(0.985 0 0)", "oklch(0.985 0 0)"},
		"--sidebar-accent":             {"oklch(0.97 0 0)", "oklch(0.269 0 0)"},
		"--sidebar-accent-foreground":  {"oklch(0.205 0 0)", "oklch(0.985 0 0)"},
		"--sidebar-border":             {"oklch(0.922 0 0)", "oklch(1 0 0 / 10%)"},
		"--sidebar-ring":               {"oklch(0.708 0 0)", "oklch(0.439 0 0)"},
	}

	got := make(map[string][2]string)
	for _, group := range pages.ThemeGroups() {
		for _, variable := range group.Vars {
			if _, ok := want[variable.Name]; ok {
				got[variable.Name] = [2]string{variable.Light, variable.Dark}
			}
		}
	}
	if len(got) != len(want) {
		t.Fatalf("editor sidebar tokens = %#v, want %#v", got, want)
	}
	for name, values := range want {
		if got[name] != values {
			t.Errorf("%s editor defaults = %#v, want %#v", name, got[name], values)
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

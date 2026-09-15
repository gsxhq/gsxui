package ui_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// sweptMessages is every English string a component writes itself and a
// caller cannot reach through children, props or attrs. Each must appear in
// the canonical source only as the argument of T(...). Issue #31.
var sweptMessages = []struct {
	file string
	text string
}{
	{"dialog", "Close"},
	{"sheet", "Close"},
	{"toast", "Close"},
	{"toaster", "Notifications"},
	{"carousel", "Previous slide"},
	{"carousel", "Next slide"},
	{"sidebar", "Sidebar"},
	{"sidebar", "Displays the mobile sidebar."},
	{"sidebar", "Toggle Sidebar"},
	{"breadcrumb", "More"},
	{"pagination", "Previous"},
	{"pagination", "Next"},
	{"pagination", "More pages"},
}

// overridableLiterals are English literals that precede `{ attrs... }` in
// their element, so a caller's attribute of the same name replaces them.
// They stay plain strings. The test pins the ordering so a reorder cannot
// silently make one unreachable.
var overridableLiterals = []struct {
	file string
	text string
}{
	{"spinner", `aria-label="Loading"`},
	{"breadcrumb", `aria-label="breadcrumb"`},
	{"pagination", `aria-label="pagination"`},
	{"pagination", `aria-label="Go to previous page"`},
	{"pagination", `aria-label="Go to next page"`},
	{"sidebar", `aria-label="Toggle Sidebar"`},
	{"sidebar", `title="Toggle Sidebar"`},
	{"carousel", `aria-roledescription="carousel"`},
	{"carousel", `aria-roledescription="slide"`},
}

func canonicalSource(t *testing.T, name string) []string {
	t.Helper()
	src, err := os.ReadFile(filepath.Join("..", "registry", "canonical", name+".gsx"))
	if err != nil {
		t.Fatal(err)
	}
	return strings.Split(string(src), "\n")
}

func TestSweptMessagesOnlyAppearInsideT(t *testing.T) {
	for _, m := range sweptMessages {
		lines := canonicalSource(t, m.file)
		quoted := `"` + m.text + `"`
		asText := ">" + m.text + "<"
		wrapped := `T("` + m.text + `")`
		bare := 0
		for i, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "//") {
				continue
			}
			if strings.Contains(line, asText) || strings.TrimSpace(line) == m.text {
				t.Errorf("%s.gsx:%d writes %q as bare text", m.file, i+1, m.text)
			}
			if isOverridableLine(m.file, line) {
				continue
			}
			bare += strings.Count(line, quoted) - strings.Count(line, wrapped)
		}
		if bare != 0 {
			t.Errorf("%s.gsx: %q appears %d time(s) outside T(...)", m.file, m.text, bare)
		}
	}
}

func isOverridableLine(file, line string) bool {
	for _, l := range overridableLiterals {
		if l.file == file && strings.Contains(line, l.text) {
			return true
		}
	}
	return false
}

func TestOverridableLiteralsPrecedeAttrs(t *testing.T) {
	for _, l := range overridableLiterals {
		lines := canonicalSource(t, l.file)
		found := false
		for i, line := range lines {
			if !strings.Contains(line, l.text) {
				continue
			}
			found = true
			// The spread may sit on the literal's own line (a one-line tag)
			// or on a later line of the same element.
			ok := strings.Contains(line, "{ attrs... }")
			for j := i + 1; j < len(lines) && !ok; j++ {
				trimmed := strings.TrimSpace(lines[j])
				if strings.Contains(trimmed, "{ attrs... }") {
					ok = true
					break
				}
				if strings.HasSuffix(trimmed, ">") || strings.HasPrefix(trimmed, "component ") {
					break
				}
			}
			if !ok {
				t.Errorf("%s.gsx:%d: %s is no longer followed by { attrs... } in its element; a caller cannot override it — move it before the spread or add it to sweptMessages", l.file, i+1, l.text)
			}
		}
		if !found {
			t.Errorf("%s.gsx: %s not found; update overridableLiterals", l.file, l.text)
		}
	}
}

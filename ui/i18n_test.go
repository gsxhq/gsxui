package ui_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// sweptMessages is every English string a component writes itself,
// including defaults a caller can override through attrs. Each must appear
// in the canonical source only as the argument of i18n.T(...). Issues #31
// and #44.
var sweptMessages = []struct {
	file string
	text string
}{
	{"dialog", "Close"},
	{"sheet", "Close"},
	{"toast", "Close"},
	{"toaster", "Notifications"},
	{"spinner", "Loading"},
	{"breadcrumb", "breadcrumb"},
	{"pagination", "pagination"},
	{"pagination", "Go to previous page"},
	{"pagination", "Go to next page"},
	{"carousel", "carousel"},
	{"carousel", "slide"},
	{"carousel", "Previous slide"},
	{"carousel", "Next slide"},
	{"sidebar", "Sidebar"},
	{"sidebar", "Displays the mobile sidebar."},
	{"sidebar", "Toggle Sidebar"},
	{"breadcrumb", "More"},
	{"pagination", "Previous"},
	{"pagination", "Next"},
	{"pagination", "More pages"},
	{"calendar", "Previous month"},
	{"calendar", "Next month"},
	{"calendar", "Month"},
	{"calendar", "Year"},
	{"calendar", "Today, {date}"},
	{"calendar", "{date}, selected"},
}

// overridableMessages are swept messages that precede `{ attrs... }` in
// their element, so a caller's attribute of the same name replaces the
// translated default. The test pins the ordering so a reorder cannot
// silently make one unreachable. Issue #44.
var overridableMessages = []struct {
	file string
	text string
}{
	{"spinner", `aria-label={i18n.T("Loading")}`},
	{"breadcrumb", `aria-label={i18n.T("breadcrumb")}`},
	{"pagination", `aria-label={i18n.T("pagination")}`},
	{"pagination", `aria-label={i18n.T("Go to previous page")}`},
	{"pagination", `aria-label={i18n.T("Go to next page")}`},
	{"sidebar", `aria-label={i18n.T("Toggle Sidebar")}`},
	{"sidebar", `title={i18n.T("Toggle Sidebar")}`},
	{"carousel", `aria-roledescription={i18n.T("carousel")}`},
	{"carousel", `aria-roledescription={i18n.T("slide")}`},
}

func canonicalSource(t *testing.T, name string) []string {
	t.Helper()
	src, err := os.ReadFile(filepath.Join("..", "registry", "canonical", name+".gsx"))
	if err != nil {
		t.Fatal(err)
	}
	return strings.Split(string(src), "\n")
}

// A deleted message nets to zero here too; the end-to-end gate in internal/cli/e2e_test.go asserts each swept string is actually rendered.
func TestSweptMessagesOnlyAppearInsideT(t *testing.T) {
	for _, m := range sweptMessages {
		lines := canonicalSource(t, m.file)
		quoted := `"` + m.text + `"`
		asText := ">" + m.text + "<"
		wrapped := `i18n.T("` + m.text + `")`
		bare := 0
		for i, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "//") {
				continue
			}
			if strings.Contains(line, asText) || strings.TrimSpace(line) == m.text {
				t.Errorf("%s.gsx:%d writes %q as bare text", m.file, i+1, m.text)
			}
			bare += strings.Count(line, quoted) - strings.Count(line, wrapped)
		}
		if bare != 0 {
			t.Errorf("%s.gsx: %q appears %d time(s) outside i18n.T(...)", m.file, m.text, bare)
		}
	}
}

func TestOverridableMessagesPrecedeAttrs(t *testing.T) {
	for _, l := range overridableMessages {
		lines := canonicalSource(t, l.file)
		found := false
		for i, line := range lines {
			if !strings.Contains(line, l.text) {
				continue
			}
			found = true
			// The spread may sit after the message on its own line (a
			// one-line tag) or on a later line of the same element. A spread
			// earlier on the same line wins over the message.
			spread := strings.Index(line, "{ attrs... }")
			ok := spread > strings.Index(line, l.text)
			for j := i + 1; j < len(lines) && !ok && spread < 0; j++ {
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
				t.Errorf("%s.gsx:%d: %s is no longer followed by { attrs... } in its element; a caller cannot override it — move it before the spread", l.file, i+1, l.text)
			}
		}
		if !found {
			t.Errorf("%s.gsx: %s not found; update overridableMessages", l.file, l.text)
		}
	}
}

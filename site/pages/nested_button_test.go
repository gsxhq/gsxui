package pages_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/site/examples"
)

// nestedButtonDepth scans rendered HTML for <button> elements nested inside
// other <button> elements. Interactive content inside a button is invalid
// HTML — the parser hoists the inner button out as a SIBLING, splitting the
// markup: the outer (wired) button ends up empty and unclickable while the
// visible inner one is orphaned from its behavior wiring. That is exactly
// how the landing page's dialog trigger silently broke (a ui.Button nested
// inside ui.DialogTrigger), so this guards every rendered page against the
// whole class. Attribute/text occurrences are not a concern: page bodies
// HTML-escape displayed source (&lt;button), so raw "<button" only appears
// as real markup.
func nestedButtonDepth(body string) (maxDepth int) {
	depth := 0
	for i := 0; i < len(body); i++ {
		if body[i] != '<' {
			continue
		}
		rest := body[i:]
		switch {
		case strings.HasPrefix(rest, "<button"):
			depth++
			if depth > maxDepth {
				maxDepth = depth
			}
		case strings.HasPrefix(rest, "</button"):
			depth--
		}
	}
	return maxDepth
}

// TestNoNestedButtons renders the landing page and every component page
// through the real route tree and fails on any button-in-button markup.
func TestNoNestedButtons(t *testing.T) {
	handler := newTestHandler(t)

	paths := []string{"/", "/components/", "/docs/getting-started", "/docs/theming", "/theme"}
	for _, name := range examples.Components() {
		paths = append(paths, "/components/"+name)
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Fatalf("GET %s = %d, want %d", path, rec.Code, http.StatusOK)
			}
			if d := nestedButtonDepth(rec.Body.String()); d > 1 {
				t.Errorf("GET %s renders a <button> nested inside a <button> (depth %d) — the HTML parser will split them and orphan the visible button from its wiring", path, d)
			}
		})
	}
}

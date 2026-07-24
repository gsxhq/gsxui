package pages_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestBranding asserts Layout wires the site branding: the favicon set in
// <head> (SVG + PNG fallback + apple-touch-icon, served by site/main.go at
// root paths) and the inline gsxui logo as the header home link's accessible
// content (replacing the old plain-text "gsxui").
func TestBranding(t *testing.T) {
	handler := newTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /: status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()

	for _, want := range []string{
		`rel="icon" href="/favicon.svg" type="image/svg+xml"`,
		`href="/favicon-32.png"`,
		`rel="apple-touch-icon" href="/apple-touch-icon.png"`,
		`aria-label="gsxui"`, // the inline header logo's accessible name
	} {
		if !strings.Contains(body, want) {
			t.Errorf("GET /: body missing %q", want)
		}
	}
}

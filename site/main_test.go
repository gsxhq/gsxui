package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestIconRoutes verifies the embedded favicon set serves at stable root
// paths with the right content types — these URLs are referenced from
// Layout's <head> and must exist in dev and prod alike.
func TestIconRoutes(t *testing.T) {
	mux := http.NewServeMux()
	registerIcons(mux)

	cases := []struct {
		path        string
		contentType string
	}{
		{"/favicon.svg", "image/svg+xml"},
		{"/favicon-32.png", "image/png"},
		{"/apple-touch-icon.png", "image/png"},
	}
	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("%s: status = %d, want 200", tc.path, rec.Code)
		}
		if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, tc.contentType) {
			t.Errorf("%s: Content-Type = %q, want prefix %q", tc.path, ct, tc.contentType)
		}
		if cc := rec.Header().Get("Cache-Control"); cc == "" {
			t.Errorf("%s: missing Cache-Control header", tc.path)
		}
		if rec.Body.Len() == 0 {
			t.Errorf("%s: empty body", tc.path)
		}
	}
}

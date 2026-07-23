package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gsxhq/gsxui/site/pages"
)

func TestPageCacheHeader(t *testing.T) {
	h := pageCache(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if got, want := rec.Header().Get("Cache-Control"), "public, max-age=0, s-maxage=300"; got != want {
		t.Fatalf("Cache-Control = %q, want %q", got, want)
	}
}

// Only content-hashed filenames may carry the forever header — a stale
// cached copy of anything else would be unfixable short of a URL change.
// Vite fingerprints exactly the files it emits under assets/ (publicDir is
// off), so the path prefix is the fingerprint guarantee.
func TestImmutableAssetsHeader(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/static/assets/main-abc123.css", "public, max-age=31536000, immutable"},
		{"/static/assets/main-abc123.js", "public, max-age=31536000, immutable"},
		{"/static/.vite/manifest.json", ""}, // not fingerprinted — no forever header
	}
	h := immutableAssets(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	}))
	for _, tc := range tests {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest("GET", tc.path, nil))
		if got := rec.Header().Get("Cache-Control"); got != tc.want {
			t.Fatalf("Cache-Control for %s = %q, want %q", tc.path, got, tc.want)
		}
	}
}

// Error responses must never carry the page cache header pageCache stamped
// before the router ran — a CDN-cached 404/500 would outlive the fault.
func TestErrorHandlerClearsCacheControl(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"typed status error", pages.ErrorWithStatus{Status: http.StatusNotFound, Message: "no such component"}, http.StatusNotFound},
		{"untyped error", errors.New("boom"), http.StatusInternalServerError},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			rec.Header().Set("Cache-Control", "public, max-age=0, s-maxage=300")
			errorHandler(rec, httptest.NewRequest("GET", "/x", nil), tc.err)
			if got := rec.Header().Get("Cache-Control"); got != "" {
				t.Fatalf("Cache-Control on error response = %q, want cleared", got)
			}
			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
		})
	}
}

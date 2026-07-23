package pages_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/internal/registry"
)

// The /components/ catalog must list every registry component — derived
// from the same source as the sidebar, so a newly added component appears
// without touching the page.
func TestComponentsIndexListsAllComponents(t *testing.T) {
	h := newTestHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/components/", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /components/ = %d, want 200", rec.Code)
	}
	body, _ := io.ReadAll(rec.Body)
	names, err := registry.Components()
	if err != nil {
		t.Fatal(err)
	}
	if len(names) == 0 {
		t.Fatal("registry returned no components")
	}
	for _, name := range names {
		if !strings.Contains(string(body), `href="/components/`+name+`"`) {
			t.Errorf("index missing link to %s", name)
		}
	}
}

// /components without the trailing slash must reach the catalog too (the
// ServeMux subtree redirect), not the 404 fallback.
func TestComponentsIndexNoTrailingSlash(t *testing.T) {
	h := newTestHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/components", nil))
	if rec.Code < 300 || rec.Code > 399 {
		t.Fatalf("GET /components = %d, want a redirect", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/components/" {
		t.Fatalf("Location = %q, want /components/", loc)
	}
}

// The home hero's "Browse components" button must land on the catalog —
// its old target was an in-page anchor that fit above the fold on desktop
// viewports and visibly did nothing.
func TestHomeBrowseButtonLinksToCatalog(t *testing.T) {
	h := newTestHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	body, _ := io.ReadAll(rec.Body)
	if !strings.Contains(string(body), `href="/components/"`) {
		t.Error(`home page has no href="/components/" link`)
	}
	if strings.Contains(string(body), `href="#components"`) {
		t.Error(`home page still links the dead #components anchor`)
	}
}

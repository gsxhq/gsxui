package pages_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/site/pages"
	"github.com/gsxhq/vite"
	"github.com/jackielii/structpages"
)

// TestSiteRoutes is the Task 1 integration smoke test: mounts the real page
// tree on a real mux (as site/main.go does) and asserts "/" renders the
// landing page with at least one live gsxui component on it.
func TestSiteRoutes(t *testing.T) {
	mux := http.NewServeMux()
	if _, err := structpages.Mount(mux, pages.Pages{}, "/", "gsxui"); err != nil {
		t.Fatalf("structpages.Mount: %v", err)
	}

	// Layout reads *vite.Vite from the request context (vite.FromContext);
	// wire dev-mode Vite the same way site/main.go's v.Middleware does. Dev
	// mode does no I/O, so no real Vite server is needed for this test.
	v, err := vite.New(vite.Config{DevURL: "http://localhost:5173", DevBase: "/__vite/"})
	if err != nil {
		t.Fatalf("vite.New: %v", err)
	}
	handler := v.Middleware(mux)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET / = %d, want %d; body:\n%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `data-slot="button"`) {
		t.Errorf(`response missing data-slot="button"; body:\n%s`, rec.Body.String())
	}
}

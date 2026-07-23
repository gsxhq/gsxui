package pages_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/site/pages"
	"github.com/gsxhq/vite"
	"github.com/jackielii/structpages"
)

// newTestHandler mounts the real page tree on a real mux the same way
// site/main.go does — same WithErrorHandler wiring (so ErrorWithStatus
// Props errors produce the right status here too, not the framework's
// default 500), same dev-mode Vite middleware (no I/O, so no real Vite
// server is needed for tests; Layout reads *vite.Vite via
// vite.FromContext).
func newTestHandler(t *testing.T) http.Handler {
	t.Helper()
	mux := http.NewServeMux()
	if _, err := structpages.Mount(mux, pages.Pages{}, "/", "gsxui",
		structpages.WithErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
			var se pages.ErrorWithStatus
			if errors.As(err, &se) {
				http.Error(w, se.Message, se.Status)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}),
	); err != nil {
		t.Fatalf("structpages.Mount: %v", err)
	}

	v, err := vite.New(vite.Config{DevURL: "http://localhost:5173", DevBase: "/__vite/"})
	if err != nil {
		t.Fatalf("vite.New: %v", err)
	}
	return v.Middleware(mux)
}

// TestSiteRoutes is the Task 1 integration smoke test: mounts the real page
// tree on a real mux (as site/main.go does) and asserts "/" renders the
// landing page with at least one live gsxui component on it.
func TestSiteRoutes(t *testing.T) {
	handler := newTestHandler(t)

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

// TestComponentPageRoute is the Task 2 integration smoke test for
// /components/{name}: a registered component renders the preview panel
// (live component) next to its literal, unescaped-by-identifier source
// text; an unregistered name 404s via ErrorWithStatus rather than the
// framework's default 500 or a rendered-but-empty page.
func TestComponentPageRoute(t *testing.T) {
	handler := newTestHandler(t)

	t.Run("registered component", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/components/button", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /components/button = %d, want %d; body:\n%s", rec.Code, http.StatusOK, rec.Body.String())
		}
		body := rec.Body.String()
		if !strings.Contains(body, "border rounded-lg p-8 bg-background") {
			t.Errorf("response missing preview panel marker; body:\n%s", body)
		}
		if !strings.Contains(body, `data-slot="button"`) {
			t.Errorf(`response missing rendered example (data-slot="button"); body:\n%s`, body)
		}
		// A distinctive identifier from basic.gsx's literal source — proves
		// the displayed source is the exact embedded file, not paraphrased.
		if !strings.Contains(body, "uibutton.Button") {
			t.Errorf("response missing literal source text %q; body:\n%s", "uibutton.Button", body)
		}
		if !strings.Contains(body, `data-site-copy`) {
			t.Errorf(`response missing copy button (data-site-copy); body:\n%s`, body)
		}
		if !strings.Contains(body, "gsxui add button") {
			t.Errorf(`response missing install snippet "gsxui add button"; body:\n%s`, body)
		}
	})

	// Task 3a's representative: proves a form-control component (input)
	// is wired through the same registry/page harness as button, not
	// just button itself.
	t.Run("registered component (input)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/components/input", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /components/input = %d, want %d; body:\n%s", rec.Code, http.StatusOK, rec.Body.String())
		}
		body := rec.Body.String()
		if !strings.Contains(body, "border rounded-lg p-8 bg-background") {
			t.Errorf("response missing preview panel marker; body:\n%s", body)
		}
		if !strings.Contains(body, `data-slot="input"`) {
			t.Errorf(`response missing rendered example (data-slot="input"); body:\n%s`, body)
		}
		// A distinctive identifier from basic.gsx's literal source — proves
		// the displayed source is the exact embedded file, not paraphrased.
		if !strings.Contains(body, "uiinput.Input") {
			t.Errorf("response missing literal source text %q; body:\n%s", "uiinput.Input", body)
		}
		if !strings.Contains(body, `data-site-copy`) {
			t.Errorf(`response missing copy button (data-site-copy); body:\n%s`, body)
		}
		if !strings.Contains(body, "gsxui add input") {
			t.Errorf(`response missing install snippet "gsxui add input"; body:\n%s`, body)
		}
	})

	t.Run("interactive component (dialog)", func(t *testing.T) {
		// Task 3c representative: proves an interactive component (behavior
		// wired via web/main.js's ui/index.js import, not just markup) also
		// renders through the same harness, including its gsxui:* events
		// example.
		req := httptest.NewRequest(http.MethodGet, "/components/dialog", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /components/dialog = %d, want %d; body:\n%s", rec.Code, http.StatusOK, rec.Body.String())
		}
		body := rec.Body.String()
		if !strings.Contains(body, `data-gsxui-dialog`) {
			t.Errorf(`response missing data-gsxui-dialog; body:\n%s`, body)
		}
		if !strings.Contains(body, "gsxui:open") {
			t.Errorf(`response missing events example source (gsxui:open); body:\n%s`, body)
		}
	})

	t.Run("unknown component", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/components/nope", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("GET /components/nope = %d, want %d; body:\n%s", rec.Code, http.StatusNotFound, rec.Body.String())
		}
	})
}

// TestThemePageRoute is the Task 5 smoke test for /theme: the page renders
// (JS-less) with the token-editing controls and a default-themed preview
// panel present in the markup — no JS assertions, since the live restyling
// itself is web/theme.js's job, not the server's.
func TestThemePageRoute(t *testing.T) {
	handler := newTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/theme", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /theme = %d, want %d; body:\n%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	body := rec.Body.String()
	if !strings.Contains(body, `data-theme-var="--primary"`) {
		t.Errorf(`response missing data-theme-var="--primary"; body:\n%s`, body)
	}
	if !strings.Contains(body, `data-theme-preview`) {
		t.Errorf(`response missing data-theme-preview; body:\n%s`, body)
	}
	if !strings.Contains(body, `data-theme-tab="light"`) {
		t.Errorf(`response missing data-theme-tab="light"; body:\n%s`, body)
	}
	if !strings.Contains(body, `data-theme-import`) {
		t.Errorf(`response missing data-theme-import; body:\n%s`, body)
	}
	// Preview panel renders the representative component set: button
	// variants, badges, a Card+Label+Input+Checkbox form row, and both
	// Alert variants — all live gsxui components, not static markup.
	if !strings.Contains(body, `data-slot="button"`) {
		t.Errorf(`response missing data-slot="button" in preview; body:\n%s`, body)
	}
	if !strings.Contains(body, `data-slot="checkbox"`) {
		t.Errorf(`response missing data-slot="checkbox" in preview; body:\n%s`, body)
	}
	if !strings.Contains(body, `role="alert"`) {
		t.Errorf(`response missing role="alert" in preview; body:\n%s`, body)
	}
}

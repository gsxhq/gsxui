package pages_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/site/pages"
	"github.com/gsxhq/vite"
	"github.com/jackielii/structpages"
	"golang.org/x/net/html"
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
	if !strings.Contains(rec.Body.String(), `data-gsxui-slot-button`) {
		t.Errorf(`response missing data-gsxui-slot-button; body:\n%s`, rec.Body.String())
	}
	document, err := html.Parse(strings.NewReader(rec.Body.String()))
	if err != nil {
		t.Fatalf("parse home response: %v", err)
	}
	if err := validateDialogTriggerSlotMarkers(document, 2); err != nil {
		t.Errorf("home response has invalid dialog trigger slot markers: %v", err)
	}
}

func TestDialogTriggerSlotMarkers(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		{
			name: "two exact empty markers",
			html: `<button data-gsxui-dialog-trigger data-gsxui-slot-dialog-trigger></button><button data-gsxui-dialog-trigger data-gsxui-slot-dialog-trigger></button>`,
		},
		{
			name: "suffix marker does not satisfy the contract",
			html: `<button data-gsxui-dialog-trigger data-gsxui-slot-dialog-trigger-extra></button><button data-gsxui-dialog-trigger data-gsxui-slot-dialog-trigger></button>`,
			want: "missing data-gsxui-slot-dialog-trigger",
		},
		{
			name: "valued marker does not satisfy the contract",
			html: `<button data-gsxui-dialog-trigger data-gsxui-slot-dialog-trigger = "true"></button><button data-gsxui-dialog-trigger data-gsxui-slot-dialog-trigger></button>`,
			want: `data-gsxui-slot-dialog-trigger has value "true"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			document, err := html.Parse(strings.NewReader(tt.html))
			if err != nil {
				t.Fatalf("parse fixture: %v", err)
			}
			err = validateDialogTriggerSlotMarkers(document, 2)
			if tt.want == "" {
				if err != nil {
					t.Fatalf("validateDialogTriggerSlotMarkers: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("validateDialogTriggerSlotMarkers error = %v, want %q", err, tt.want)
			}
		})
	}
}

func validateDialogTriggerSlotMarkers(document *html.Node, want int) error {
	var triggers []*html.Node
	var visit func(*html.Node)
	visit = func(node *html.Node) {
		if node.Type == html.ElementNode && hasHTMLAttribute(node, "data-gsxui-dialog-trigger") {
			triggers = append(triggers, node)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			visit(child)
		}
	}
	visit(document)

	if len(triggers) != want {
		return fmt.Errorf("found %d dialog triggers, want %d", len(triggers), want)
	}
	for index, trigger := range triggers {
		value, ok := htmlAttribute(trigger, "data-gsxui-slot-dialog-trigger")
		if !ok {
			return fmt.Errorf("dialog trigger %d missing data-gsxui-slot-dialog-trigger", index+1)
		}
		if value != "" {
			return fmt.Errorf("dialog trigger %d data-gsxui-slot-dialog-trigger has value %q", index+1, value)
		}
	}
	return nil
}

func hasHTMLAttribute(node *html.Node, name string) bool {
	_, ok := htmlAttribute(node, name)
	return ok
}

func htmlAttribute(node *html.Node, name string) (string, bool) {
	for _, attr := range node.Attr {
		if attr.Key == name {
			return attr.Val, true
		}
	}
	return "", false
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
		if !strings.Contains(body, `data-gsxui-slot-button`) {
			t.Errorf(`response missing rendered example (data-gsxui-slot-button); body:\n%s`, body)
		}
		// A distinctive identifier from basic.gsx's literal source — proves
		// the displayed source is the exact embedded file, not paraphrased.
		if !strings.Contains(body, "ui.Button") {
			t.Errorf("response missing literal source text %q; body:\n%s", "ui.Button", body)
		}
		if !strings.Contains(body, `data-site-copy`) {
			t.Errorf(`response missing copy button (data-site-copy); body:\n%s`, body)
		}
		if !strings.Contains(body, "gsxui add button") {
			t.Errorf(`response missing install snippet "gsxui add button"; body:\n%s`, body)
		}
		if strings.Contains(body, `data-site-isolated-preview`) {
			t.Errorf("ordinary button examples should render inline; body:\n%s", body)
		}
	})

	t.Run("viewport-owned component uses isolated previews", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/components/sidebar", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /components/sidebar = %d, want %d; body:\n%s", rec.Code, http.StatusOK, rec.Body.String())
		}
		body := rec.Body.String()
		if got := strings.Count(body, `data-site-isolated-preview`); got != 10 {
			t.Errorf("sidebar page has %d isolated previews, want 10; body:\n%s", got, body)
		}
		for _, marker := range []string{
			`title="Basic preview"`,
			`src="/examples/sidebar/basic"`,
			`title="variant=sidebar (default) preview"`,
			`src="/examples/sidebar/variants?_preview=sidebar"`,
			`title="variant=floating preview"`,
			`src="/examples/sidebar/variants?_preview=floating"`,
			`title="icon collapsed preview"`,
			`src="/examples/sidebar/variants?_preview=icon-collapsed"`,
			`title="Persisted (cookie round-trip) preview"`,
			`src="/examples/sidebar/persisted"`,
			`ui.SidebarProvider`,
			`data-site-copy`,
		} {
			if !strings.Contains(body, marker) {
				t.Errorf("sidebar page missing %q; body:\n%s", marker, body)
			}
		}
		if strings.Contains(body, `data-gsxui-slot-sidebar-container`) {
			t.Errorf("sidebar implementation leaked into the parent document; body:\n%s", body)
		}
		if got := strings.Count(body, `data-site-copy`); got != 3 {
			t.Errorf("sidebar page has %d source copy controls, want one per registered example (3)", got)
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
		if !strings.Contains(body, `data-gsxui-slot-input`) {
			t.Errorf(`response missing rendered example (data-gsxui-slot-input); body:\n%s`, body)
		}
		// A distinctive identifier from basic.gsx's literal source — proves
		// the displayed source is the exact embedded file, not paraphrased.
		if !strings.Contains(body, "ui.Input") {
			t.Errorf("response missing literal source text %q; body:\n%s", "ui.Input", body)
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

	// Task-6-review representative: proves a renamed component's footer
	// link points at shadcn's actual slug (dropdown→dropdown-menu), not the
	// gsxui component name verbatim, which would 404 on ui.shadcn.com. (Prior
	// to the single-package sweep this used switchctl, whose registered name
	// only differed from its shadcn slug because "switch" is a Go keyword —
	// now that the registry-facing name is "switch" itself, name and slug
	// coincide and it no longer exercises shadcnSlug's rename path.)
	t.Run("renamed component (dropdown) links to shadcn's dropdown-menu slug", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/components/dropdown", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /components/dropdown = %d, want %d; body:\n%s", rec.Code, http.StatusOK, rec.Body.String())
		}
		body := rec.Body.String()
		if !strings.Contains(body, `ui.shadcn.com/docs/components/dropdown-menu"`) {
			t.Errorf(`response missing shadcn link to renamed slug "dropdown-menu"; body:\n%s`, body)
		}
	})

	// icon has no shadcn/ui counterpart (it ports Lucide directly), so the
	// footer must not link ui.shadcn.com/docs/components/icon (404) and
	// should point at lucide.dev instead.
	t.Run("icon links to lucide.dev instead of shadcn", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/components/icon", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /components/icon = %d, want %d; body:\n%s", rec.Code, http.StatusOK, rec.Body.String())
		}
		body := rec.Body.String()
		if !strings.Contains(body, "lucide.dev") {
			t.Errorf("response missing lucide.dev link; body:\n%s", body)
		}
		if strings.Contains(body, "ui.shadcn.com/docs/components/icon") {
			t.Errorf("response should not link ui.shadcn.com/docs/components/icon (404); body:\n%s", body)
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

func TestExamplePreviewRoute(t *testing.T) {
	handler := newTestHandler(t)

	t.Run("registered example", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/examples/sidebar/basic", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /examples/sidebar/basic = %d, want %d; body:\n%s", rec.Code, http.StatusOK, rec.Body.String())
		}
		body := rec.Body.String()
		for _, marker := range []string{
			`data-site-isolated-document`,
			`data-gsxui-slot-sidebar-wrapper`,
			`data-gsxui-slot-sidebar-container`,
			`src="/__vite/@vite/client"`,
			`src="/__vite/web/main.js"`,
		} {
			if !strings.Contains(body, marker) {
				t.Errorf("preview response missing %q; body:\n%s", marker, body)
			}
		}
		for _, forbidden := range []string{
			`<iframe`,
			`Search docs...`,
			`View the original on shadcn/ui`,
		} {
			if strings.Contains(body, forbidden) {
				t.Errorf("preview response unexpectedly contains %q; body:\n%s", forbidden, body)
			}
		}
	})

	for _, path := range []string{
		"/examples/sidebar/missing",
		"/examples/missing/basic",
		"/examples/sidebar/variants?_preview=missing",
		"/examples/sidebar/variants?_preview=",
		"/examples/sidebar/variants?_preview=floating&_preview=inset",
	} {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code != http.StatusNotFound {
				t.Fatalf("GET %s = %d, want %d; body:\n%s", path, rec.Code, http.StatusNotFound, rec.Body.String())
			}
		})
	}
}

func TestThemePreviewButtonRoute(t *testing.T) {
	handler := newTestHandler(t)

	t.Run("isolated exact-style matrix", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/theme/preview/button", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /theme/preview/button = %d, want %d; body:\n%s", rec.Code, http.StatusOK, rec.Body.String())
		}
		body := rec.Body.String()
		for _, marker := range []string{
			`<!DOCTYPE html>`,
			`data-theme-button-preview`,
			`src="/__vite/@vite/client"`,
			`src="/__vite/web/preview.js"`,
			`data-theme-preview-style="nova"`,
			`data-theme-preview-style="maia" hidden`,
			`data-theme-preview-case="text"`,
			`data-theme-preview-case="icon"`,
			`data-theme-preview-case="disabled"`,
			`data-theme-preview-case="link"`,
			`data-theme-preview-case="focus-visible"`,
			`data-theme-preview-case="invalid"`,
			`data-theme-preview-case="caller-composition"`,
			`aria-invalid="true"`,
			`href="/docs/getting-started"`,
		} {
			if !strings.Contains(body, marker) {
				t.Errorf("theme preview response missing %q; body:\n%s", marker, body)
			}
		}
		for _, variant := range []string{"default", "secondary", "destructive", "outline", "ghost", "link"} {
			if got := strings.Count(body, `data-variant="`+variant+`"`); got < 2 {
				t.Errorf("theme preview has %d %q variants, want one per style; body:\n%s", got, variant, body)
			}
		}
		for _, size := range []string{"default", "xs", "sm", "lg", "icon", "icon-xs", "icon-sm", "icon-lg"} {
			if got := strings.Count(body, `data-size="`+size+`"`); got < 2 {
				t.Errorf("theme preview has %d %q sizes, want one per style; body:\n%s", got, size, body)
			}
		}
		for _, forbidden := range []string{
			`<iframe`,
			`Search docs...`,
			`Theme editor`,
			`data-theme-var=`,
			`src="/__vite/web/main.js"`,
		} {
			if strings.Contains(body, forbidden) {
				t.Errorf("theme preview response unexpectedly contains %q; body:\n%s", forbidden, body)
			}
		}
	})

	t.Run("extra route segment is not accepted", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/theme/preview/button/extra", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("GET /theme/preview/button/extra = %d, want %d; body:\n%s", rec.Code, http.StatusNotFound, rec.Body.String())
		}
	})
}

// TestDocsRoutes is the Task 4 smoke test for the two standalone docs
// pages: both render 200 with a distinctive marker, and each links to the
// other via the sidebar's Docs section so the active-state wiring is
// exercised too.
func TestDocsRoutes(t *testing.T) {
	handler := newTestHandler(t)

	t.Run("getting started", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/docs/getting-started", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /docs/getting-started = %d, want %d; body:\n%s", rec.Code, http.StatusOK, rec.Body.String())
		}
		body := rec.Body.String()
		if !strings.Contains(body, `data-doc="getting-started"`) {
			t.Errorf(`response missing data-doc="getting-started"; body:\n%s`, body)
		}
		// The CLI's actual printed init output (internal/cli/init.go), not a
		// paraphrase — proves the walkthrough shows real output.
		if !strings.Contains(body, "gsxui initialized.") {
			t.Errorf(`response missing real "gsxui initialized." CLI output; body:\n%s`, body)
		}
		if !strings.Contains(body, "adding: button card") {
			t.Errorf(`response missing real "adding: button card" CLI output; body:\n%s`, body)
		}
		if !strings.Contains(body, "/docs/theming") {
			t.Errorf("response missing sidebar/inline link to /docs/theming; body:\n%s", body)
		}
	})

	t.Run("theming", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/docs/theming", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /docs/theming = %d, want %d; body:\n%s", rec.Code, http.StatusOK, rec.Body.String())
		}
		body := rec.Body.String()
		if !strings.Contains(body, `data-doc="theming"`) {
			t.Errorf(`response missing data-doc="theming"; body:\n%s`, body)
		}
		if !strings.Contains(body, "--primary-foreground") {
			t.Errorf("response missing token model content (--primary-foreground); body:\n%s", body)
		}
		if !strings.Contains(body, "data-gsxui-dialog-trigger") {
			t.Errorf("response missing data-attribute idiom example (data-gsxui-dialog-trigger); body:\n%s", body)
		}
		if !strings.Contains(body, `[<span class="ts-attribute">data-gsxui-slot-button</span>]`) {
			t.Errorf("response missing exact slot presence selector; body:\n%s", body)
		}
		if !strings.Contains(body, "/docs/getting-started") {
			t.Errorf("response missing sidebar link to /docs/getting-started; body:\n%s", body)
		}
	})
}

// TestThemePageRoute is the theme editor's JS-less picker contract: catalog
// choices arrive from the server and raw per-token editing never leaks into
// the completed editor markup.
func TestThemePageRoute(t *testing.T) {
	handler := newTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/theme", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /theme = %d, want %d; body:\n%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	body := rec.Body.String()
	for _, picker := range []string{"baseColor", "theme", "radius"} {
		if got := strings.Count(body, `data-theme-picker="`+picker+`"`); got != 1 {
			t.Errorf("theme page has %d %s pickers, want 1; body:\n%s", got, picker, body)
		}
	}
	for _, marker := range []string{
		`data-theme-picker-trigger`,
		`data-theme-choice`,
		`data-theme-choice-swatch`,
		`data-theme-selection-value`,
		`data-theme-selection-swatch`,
		`type="radio"`,
		`aria-label=`,
		`checked`,
		`value="neutral"`,
		`value="blue"`,
		`value="medium"`,
		`Neutral`,
		`Blue`,
		`Medium`,
	} {
		if !strings.Contains(body, marker) {
			t.Errorf("response missing picker marker %q; body:\n%s", marker, body)
		}
	}
	for _, forbidden := range []string{
		`data-theme-var=`,
		`data-theme-field="light.`,
		`data-theme-field="dark.`,
	} {
		if strings.Contains(body, forbidden) {
			t.Errorf("theme page still exposes raw token input %q; body:\n%s", forbidden, body)
		}
	}
	if !strings.Contains(body, `data-theme-preview-frame`) {
		t.Errorf(`response missing data-theme-preview-frame; body:\n%s`, body)
	}
	if got := strings.Count(body, `<iframe`); got != 1 {
		t.Errorf("theme page has %d iframes, want 1; body:\n%s", got, body)
	}
	if !strings.Contains(body, `src="/theme/preview/button"`) {
		t.Errorf(`response missing Button preview iframe route; body:\n%s`, body)
	}
	if !strings.Contains(body, `data-theme-style="maia"`) {
		t.Errorf(`response missing Maia style picker; body:\n%s`, body)
	}
	if !strings.Contains(body, `data-theme-mode-tab="light"`) {
		t.Errorf(`response missing data-theme-mode-tab="light"; body:\n%s`, body)
	}
	if !strings.Contains(body, `data-theme-import="json"`) {
		t.Errorf(`response missing preset JSON import; body:\n%s`, body)
	}
	if !strings.Contains(body, `data-theme-import="css"`) {
		t.Errorf(`response missing theme CSS import; body:\n%s`, body)
	}
	if !strings.Contains(body, `data-theme-command="init"`) {
		t.Errorf(`response missing init command; body:\n%s`, body)
	}
}

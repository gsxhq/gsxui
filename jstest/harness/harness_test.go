package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// repoRoot is two levels up from jstest/harness.
func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("resolving repo root: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Fatalf("repo root %q has no go.mod: %v", root, err)
	}
	return root
}

func TestBuildManifestCoversRegisteredExamples(t *testing.T) {
	m := buildManifest()
	if len(m) < 100 {
		t.Fatalf("manifest has %d entries, want at least 100", len(m))
	}

	var got *entry
	for i := range m {
		if m[i].Component == "dropdown" && m[i].Example == "checkboxes" {
			got = &m[i]
		}
	}
	if got == nil {
		t.Fatal("dropdown/checkboxes missing from manifest")
	}
	if got.URL != "/x/dropdown/checkboxes" {
		t.Errorf("URL = %q, want /x/dropdown/checkboxes", got.URL)
	}
}

// Hyphenated component names must survive into the manifest — the example
// directory is navigationmenu (Go package names can't contain hyphens) but
// examples.For is keyed by the registered name.
func TestManifestUsesRegisteredNamesNotDirectoryNames(t *testing.T) {
	for _, e := range buildManifest() {
		if e.Component == "navigationmenu" {
			t.Fatalf("manifest used directory name %q; want registered name navigation-menu", e.Component)
		}
	}
	var found bool
	for _, e := range buildManifest() {
		if e.Component == "navigation-menu" {
			found = true
		}
	}
	if !found {
		t.Error("navigation-menu missing from manifest")
	}
}

func TestExampleRouteRendersTheExample(t *testing.T) {
	srv := httptest.NewServer(newMux(repoRoot(t)))
	defer srv.Close()

	res, err := http.Get(srv.URL + "/x/toggle/basic")
	if err != nil {
		t.Fatalf("GET /x/toggle/basic: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	body, _ := io.ReadAll(res.Body)
	page := string(body)

	for _, want := range []string{
		`<!DOCTYPE html>`,
		`<link rel="stylesheet" href="/static/jstest/.tmp/site.css">`,
		`<script type="module" src="/ui/index.js"></script>`,
		`class="min-h-svh bg-background text-foreground antialiased"`,
		`data-gsxui-toggle`,
	} {
		if !strings.Contains(page, want) {
			t.Errorf("page missing %q", want)
		}
	}
}

func TestUnknownExampleIs404(t *testing.T) {
	srv := httptest.NewServer(newMux(repoRoot(t)))
	defer srv.Close()

	for _, path := range []string{"/x/toggle/nope", "/x/nosuchcomponent/basic", "/x/toggle", "/x/"} {
		res, err := http.Get(srv.URL + path)
		if err != nil {
			t.Fatalf("GET %s: %v", path, err)
		}
		res.Body.Close()
		if res.StatusCode != http.StatusNotFound {
			t.Errorf("GET %s status = %d, want 404", path, res.StatusCode)
		}
	}
}

func TestHealthz(t *testing.T) {
	srv := httptest.NewServer(newMux(repoRoot(t)))
	defer srv.Close()

	res, err := http.Get(srv.URL + "/healthz")
	if err != nil {
		t.Fatalf("GET /healthz: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", res.StatusCode)
	}
}

func TestManifestFlagWritesJSON(t *testing.T) {
	out := filepath.Join(t.TempDir(), "examples.json")
	if err := writeManifest(out); err != nil {
		t.Fatalf("writeManifest: %v", err)
	}
	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("reading manifest: %v", err)
	}
	var entries []entry
	if err := json.Unmarshal(b, &entries); err != nil {
		t.Fatalf("unmarshalling manifest: %v", err)
	}
	if len(entries) != len(buildManifest()) {
		t.Errorf("wrote %d entries, buildManifest has %d", len(entries), len(buildManifest()))
	}
}

func TestServesRealModuleSourceByteForByte(t *testing.T) {
	root := repoRoot(t)
	srv := httptest.NewServer(newMux(root))
	defer srv.Close()

	for _, name := range []string{"index.js", "dropdown.js", "gsxui.js"} {
		want, err := os.ReadFile(filepath.Join(root, "ui", name))
		if err != nil {
			t.Fatalf("reading ui/%s: %v", name, err)
		}
		res, err := http.Get(srv.URL + "/ui/" + name)
		if err != nil {
			t.Fatalf("GET /ui/%s: %v", name, err)
		}
		got, _ := io.ReadAll(res.Body)
		res.Body.Close()
		if !bytes.Equal(got, want) {
			t.Errorf("/ui/%s differs from ui/%s on disk", name, name)
		}
		if ct := res.Header.Get("Content-Type"); !strings.Contains(ct, "javascript") {
			t.Errorf("/ui/%s Content-Type = %q, want a javascript type", name, ct)
		}
	}
}

// The shim tree is the real one except for gsxui.js. Every sibling must be
// byte-identical, or the disjointness check would be recording selectors
// from code that isn't shipped.
func TestShimPassesThroughEverythingButGsxui(t *testing.T) {
	root := repoRoot(t)
	srv := httptest.NewServer(newMux(root))
	defer srv.Close()

	want, err := os.ReadFile(filepath.Join(root, "ui", "dropdown.js"))
	if err != nil {
		t.Fatalf("reading ui/dropdown.js: %v", err)
	}
	res, err := http.Get(srv.URL + "/shim/dropdown.js")
	if err != nil {
		t.Fatalf("GET /shim/dropdown.js: %v", err)
	}
	got, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if !bytes.Equal(got, want) {
		t.Error("/shim/dropdown.js is not the real ui/dropdown.js")
	}
}

func TestShimReplacesGsxui(t *testing.T) {
	root := repoRoot(t)
	srv := httptest.NewServer(newMux(root))
	defer srv.Close()

	real, err := os.ReadFile(filepath.Join(root, "ui", "gsxui.js"))
	if err != nil {
		t.Fatalf("reading ui/gsxui.js: %v", err)
	}
	res, err := http.Get(srv.URL + "/shim/gsxui.js")
	if err != nil {
		t.Fatalf("GET /shim/gsxui.js: %v", err)
	}
	got, _ := io.ReadAll(res.Body)
	res.Body.Close()

	if bytes.Equal(got, real) {
		t.Fatal("/shim/gsxui.js served the real module; the shim did not take effect")
	}
	for _, want := range []string{"__gsxuiRegistrations", "export function on", "export function emit"} {
		if !strings.Contains(string(got), want) {
			t.Errorf("shim missing %q", want)
		}
	}
}

func TestModuleRoutesRejectNonJSAndEscapes(t *testing.T) {
	srv := httptest.NewServer(newMux(repoRoot(t)))
	defer srv.Close()

	for _, path := range []string{
		"/ui/../go.mod",
		"/ui/icon",
		"/ui/nosuch.js",
		"/ui/button.gsx",
		"/shim/../go.mod",
		"/shim/nosuch.js",
	} {
		res, err := http.Get(srv.URL + path)
		if err != nil {
			t.Fatalf("GET %s: %v", path, err)
		}
		body, _ := io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode == http.StatusOK {
			t.Errorf("GET %s returned 200 with %d bytes; want a non-200", path, len(body))
		}
	}
}

func TestRegistrationsPageLoadsTheShim(t *testing.T) {
	srv := httptest.NewServer(newMux(repoRoot(t)))
	defer srv.Close()

	res, err := http.Get(srv.URL + "/registrations")
	if err != nil {
		t.Fatalf("GET /registrations: %v", err)
	}
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	page := string(body)
	if !strings.Contains(page, `<script type="module" src="/shim/index.js"></script>`) {
		t.Error("/registrations does not load /shim/index.js")
	}
	if strings.Contains(page, "/ui/index.js") {
		t.Error("/registrations also loads the real /ui/index.js; the shim would not be the module actually recording")
	}
}

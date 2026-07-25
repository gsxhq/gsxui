package main

import (
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

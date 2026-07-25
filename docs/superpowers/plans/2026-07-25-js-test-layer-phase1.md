# JS Test Layer — Phase 1 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Stand up a real-browser test layer for `ui/*.js` and gate CI on it, with four invariant sweeps that already cover all 103 examples before any per-component spec exists.

**Architecture:** A Go harness server (`jstest/harness`) imports `site/examples` and renders one example per page, serving the real `ui/*.js` modules as native ES modules with no bundler in the path. Playwright Test drives real Chromium against it. A `/shim/` module tree — identical except `gsxui.js` records `on()` calls instead of registering them — makes the delegation registry observable, which is what turns the hook-prefix collision into a mechanical check.

**Tech Stack:** Go 1.26.1, `@playwright/test` (Chromium only), `@tailwindcss/cli`, TypeScript specs (Playwright transpiles, no tsconfig).

**Source spec:** `docs/superpowers/specs/2026-07-25-js-test-layer-design.md`

**Phase scope:** This plan is Phase 1 of three. It delivers the harness, the CSS build, the Playwright rig, the four invariants, and CI wiring. Phase 2 (per-component specs for the high-risk eight, plus the `/f/<name>` fixture route and `jstest/fixtures/`) and Phase 3 (the remaining thirteen modules) get their own plans, written once the harness API is real rather than predicted. Nothing in the spec is dropped — only sequenced.

## Global Constraints

- **Real browser only.** No jsdom, no happy-dom, no `--experimental-vm-modules` DOM shims. jsdom resolves `:open` and `:focus-visible` to `false` without throwing, so a jsdom assertion can pass while testing nothing.
- **No bundler in the test path.** Pages load `ui/*.js` verbatim over HTTP as native ES modules. Never add Vite, esbuild, or a rollup step between a test and the code under test.
- **Routes use registered component names, not directory names.** `/x/navigation-menu/basic`, never `/x/navigationmenu/basic`. Directory names drop hyphens for Go package naming (`navigationmenu`, `contextmenu`, `nativeselect`, `switchctl`); `examples.For` is keyed by the registered name.
- **No global animation disabling.** Do not inject `* { transition: none !important }`. `dialog.js` and `sonner.js` wait on `transitionend` with a 600ms fallback cap; disabling transitions pushes both onto the cap and hides real timing bugs. Use Playwright's retrying assertions instead.
- **Never start or kill processes outside the test rig.** A `gsx dev` server on :7777 and a Vite server on :5173 may be running and belong to the user. The harness binds **127.0.0.1:7799**. Do not run `make site-dev`, do not `pkill`, do not free a port.
- **Every new test must be proven able to fail.** Each task that adds an assertion has an explicit red-validation step: break the thing, watch the test fail with the expected message, restore. A test that cannot go red is not a test.
- **`make test` stays Go-only.** JS tests go in `make test-js`. Do not add browser work to the fast inner loop.
- Chromium only. Firefox and WebKit are deferred by the spec.

---

## File Structure

```
jstest/
  harness/
    main.go            flags, routes, server startup, manifest mode
    shell.go           the HTML page shell around a rendered example
    manifest.go        manifest entry type + builder + JSON writer
    modules.go         /ui/ and /shim/ module serving
    shim.js            the recording gsxui.js replacement (go:embed'd)
    harness_test.go    Go tests for all of the above
  support/
    paths.ts           shared absolute paths (repo root, tmp dir, manifest, css)
    manifest.ts        reads the generated manifest synchronously
    fixtures.ts        Playwright fixture exposing recorded on() registrations
    selector-allowlist.ts   reviewed, justified cross-module selector overlaps
  specs/
    smoke.spec.ts      the rig works: markup + CSS + JS all arrive
    invariants.spec.ts one test per example, four soft-asserted invariants
  global-setup.ts      generates the manifest and compiles the stylesheet
  playwright.config.ts
  .tmp/                gitignored: examples.json, site.css
```

Responsibilities are split so each file stays small enough to hold in context: routing is separate from page rendering, page rendering separate from module serving, and the TypeScript support modules each expose one thing.

---

### Task 1: Harness skeleton — example routes and the manifest

**Files:**
- Create: `jstest/harness/main.go`
- Create: `jstest/harness/manifest.go`
- Create: `jstest/harness/shell.go`
- Create: `jstest/harness/harness_test.go`

**Interfaces:**
- Consumes: `site/examples`'s `examples.Components() []string`, `examples.For(component) []Example`, `Example{Name, Title string; Node gsx.Node; SourcePath string}`.
- Produces: `entry{Component, Example, URL string}` with JSON tags `component`/`example`/`url`; `buildManifest() []entry`; `newMux(root string) http.Handler`; `renderShell(w io.Writer, title string, script string, body template.HTML) error`. Task 2 adds routes to `newMux` and reuses `renderShell`. Task 3's TypeScript reads the manifest JSON shape.

- [ ] **Step 1: Write the failing test**

Create `jstest/harness/harness_test.go`:

```go
package main

import (
	"encoding/json"
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
```

Add `"io"` to the import block.

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test ./jstest/harness/`
Expected: FAIL — the package does not compile (`undefined: buildManifest`, `undefined: entry`, `undefined: newMux`, `undefined: writeManifest`).

- [ ] **Step 3: Write the manifest**

Create `jstest/harness/manifest.go`:

```go
package main

import (
	"encoding/json"
	"os"

	"github.com/gsxhq/gsxui/site/examples"
)

// entry is one addressable example page. Component is the REGISTERED
// component name ("navigation-menu"), not the example directory name
// ("navigationmenu") — Go package names can't contain hyphens, but
// examples.For is keyed by the registered name and that is what tests read.
type entry struct {
	Component string `json:"component"`
	Example   string `json:"example"`
	URL       string `json:"url"`
}

// buildManifest enumerates every registered example, in registration order.
func buildManifest() []entry {
	var out []entry
	for _, component := range examples.Components() {
		for _, ex := range examples.For(component) {
			out = append(out, entry{
				Component: component,
				Example:   ex.Name,
				URL:       "/x/" + component + "/" + ex.Name,
			})
		}
	}
	return out
}

// writeManifest serialises buildManifest to path. Playwright's globalSetup
// runs this before workers import spec files, so the specs can read the
// example list synchronously and generate one test per example.
func writeManifest(path string) error {
	b, err := json.MarshalIndent(buildManifest(), "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}
```

- [ ] **Step 4: Write the page shell**

Create `jstest/harness/shell.go`:

```go
package main

import (
	"html/template"
	"io"
)

// shellTmpl is the minimal page every harness route renders into. It
// deliberately carries only what a component needs: the compiled stylesheet
// and one module script. No site chrome, no web/site.js, no theme script —
// those are site code, not library code, and loading them would put
// untested JS in the way of the JS under test.
//
// The body classes match site/pages/layout.gsx so theme tokens resolve the
// same way they do in production.
var shellTmpl = template.Must(template.New("shell").Parse(
	`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Title}}</title>
<link rel="stylesheet" href="/static/jstest/.tmp/site.css">
<script type="module" src="{{.Script}}"></script>
</head>
<body class="min-h-svh bg-background text-foreground antialiased">
<main data-harness-root class="p-8">{{.Body}}</main>
</body>
</html>
`))

type shellData struct {
	Title  string
	Script string
	Body   template.HTML
}

// renderShell writes the shell around already-rendered markup. body is
// trusted: it comes from a gsx component's own Render, which escapes its
// own interpolations.
func renderShell(w io.Writer, title, script string, body template.HTML) error {
	return shellTmpl.Execute(w, shellData{Title: title, Script: script, Body: body})
}
```

- [ ] **Step 5: Write the server**

Create `jstest/harness/main.go`:

```go
// Command harness serves gsxui's component examples one per page, for the
// Playwright suite in jstest/. It is a test tool: it binds loopback only and
// is never built into the shipped site.
//
//	go run ./jstest/harness -addr 127.0.0.1:7799 -root .
//	go run ./jstest/harness -manifest jstest/.tmp/examples.json
package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gsxhq/gsxui/site/examples"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:7799", "loopback address to listen on")
	root := flag.String("root", ".", "repo root, served read-only under /static/")
	manifest := flag.String("manifest", "", "write the example manifest to this path and exit")
	flag.Parse()

	if *manifest != "" {
		if err := writeManifest(*manifest); err != nil {
			log.Fatalf("writing manifest: %v", err)
		}
		return
	}

	if !strings.HasPrefix(*addr, "127.0.0.1:") && !strings.HasPrefix(*addr, "localhost:") {
		log.Fatalf("-addr %q is not loopback; the harness serves the repo tree and must not be reachable off-host", *addr)
	}

	fmt.Fprintf(os.Stderr, "harness listening on http://%s\n", *addr)
	if err := http.ListenAndServe(*addr, newMux(*root)); err != nil {
		log.Fatal(err)
	}
}

func newMux(root string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// /static/ serves the repo tree read-only. The compiled stylesheet lands
	// at jstest/.tmp/site.css, and Tailwind's bundled @fontsource imports
	// carry url() references relative to that output file — serving from the
	// repo root is what lets those resolve instead of 404ing into the
	// clean-load invariant.
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir(root))))

	mux.HandleFunc("GET /x/{component}/{example}", func(w http.ResponseWriter, r *http.Request) {
		component := r.PathValue("component")
		name := r.PathValue("example")

		for _, ex := range examples.For(component) {
			if ex.Name != name {
				continue
			}
			var buf bytes.Buffer
			if err := ex.Node.Render(r.Context(), &buf); err != nil {
				http.Error(w, "render: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			title := component + "/" + name
			if err := renderShell(w, title, "/ui/index.js", template.HTML(buf.String())); err != nil {
				log.Printf("rendering shell for %s: %v", title, err)
			}
			return
		}
		http.NotFound(w, r)
	})

	return mux
}
```

- [ ] **Step 6: Run the tests to verify they pass**

Run: `go test ./jstest/harness/ -v`
Expected: PASS — all six tests.

- [ ] **Step 7: Verify the manifest by hand**

Run: `go run ./jstest/harness -manifest /tmp/m.json && head -20 /tmp/m.json && python3 -c "import json;print(len(json.load(open('/tmp/m.json'))))"`
Expected: valid JSON, first entry has `component`/`example`/`url`, count is 103.

If the count is not 103, do not adjust the test's `>= 100` floor to match — investigate. A changed count means examples were added or removed, which is fine, but confirm that is actually what happened.

- [ ] **Step 8: Commit**

```bash
git add jstest/harness/
git commit -m "test(js): harness serves one example per page

A Go test server over site/examples: /x/<component>/<example> renders the
real gsx.Node into a minimal shell, and -manifest writes the component ×
example index Playwright needs to generate one test per example."
```

---

### Task 2: Serve the real modules, and a recording shim

**Files:**
- Create: `jstest/harness/modules.go`
- Create: `jstest/harness/shim.js`
- Modify: `jstest/harness/main.go` (register the new routes in `newMux`)
- Modify: `jstest/harness/harness_test.go` (append tests)

**Interfaces:**
- Consumes: `newMux(root string)` and `renderShell` from Task 1.
- Produces: routes `GET /ui/{file}`, `GET /shim/{file}`, `GET /registrations`. The `/registrations` page loads `/shim/index.js` with an empty body and exposes `window.__gsxuiRegistrations`, an array of `{type: string, capture: boolean, selector: string, module: string}` where `module` is the bare module filename (`"dropdown.js"`). Task 6 consumes that array.

- [ ] **Step 1: Write the failing tests**

Append to `jstest/harness/harness_test.go`:

```go
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
```

Add `"bytes"` to the import block if it is not already there.

- [ ] **Step 2: Run the tests to verify they fail**

Run: `go test ./jstest/harness/ -run 'Module|Shim|Registrations|Serves' -v`
Expected: FAIL — every new test 404s or the package fails to compile.

- [ ] **Step 3: Write the shim**

Create `jstest/harness/shim.js`:

```js
// Recording replacement for ui/gsxui.js, served at /shim/gsxui.js.
//
// It has gsxui.js's exact export surface but registers nothing: on() records
// what a module WOULD have bound, so the Playwright suite can check that no
// two modules can claim the same element for the same (type, capture) pair.
// That is the mechanical form of the hook-prefix collision that shipped in
// Tier 4 Batch B, where dropdown.js and context-menu.js both matched
// data-gsxui-menu-* and one click ran both handlers.
//
// emit() is the real implementation, unchanged — nothing depends on it here,
// but a module that calls emit at import time must not throw.

const registrations = [];
window.__gsxuiRegistrations = registrations;

export function on(type, selector, handler, { capture = false } = {}) {
  registrations.push({ type, capture, selector, module: callerModule() });
}

export function emit(el, name, detail) {
  return el.dispatchEvent(
    new CustomEvent(name, { bubbles: true, cancelable: true, detail }),
  );
}

// callerModule walks the V8 stack to the first frame outside this file. Frame
// 0 is "Error", frame 1 is on() itself, frame 2 is the module that called it.
// Returns the bare filename ("dropdown.js") so failure output is readable.
function callerModule() {
  const frames = (new Error().stack || "").split("\n");
  for (const frame of frames) {
    const match = frame.match(/\/shim\/([A-Za-z0-9._-]+\.js)/);
    if (match && match[1] !== "gsxui.js") return match[1];
  }
  return "unknown";
}
```

- [ ] **Step 4: Write the module routes**

Create `jstest/harness/modules.go`:

```go
package main

import (
	_ "embed"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//go:embed shim.js
var shimSource []byte

// registerModuleRoutes wires the two module trees:
//
//   /ui/*.js    the real ui/ modules, byte-for-byte, as native ES modules
//   /shim/*.js  the same, except gsxui.js is shim.js
//
// Serving the real source with no bundler is deliberate: nothing sits
// between a test and the code under test, so there is no transform to keep
// in sync and no bundle cache to invalidate.
func registerModuleRoutes(mux *http.ServeMux, root string) {
	uiDir := filepath.Join(root, "ui")

	serve := func(w http.ResponseWriter, r *http.Request, shim bool) {
		name := r.PathValue("file")
		// Only bare .js filenames. filepath.Base collapses any traversal
		// attempt, and the equality check rejects anything that changed.
		if filepath.Base(name) != name || !strings.HasSuffix(name, ".js") {
			http.NotFound(w, r)
			return
		}
		if shim && name == "gsxui.js" {
			w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
			w.Write(shimSource)
			return
		}
		b, err := os.ReadFile(filepath.Join(uiDir, name))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		w.Write(b)
	}

	mux.HandleFunc("GET /ui/{file}", func(w http.ResponseWriter, r *http.Request) {
		serve(w, r, false)
	})
	mux.HandleFunc("GET /shim/{file}", func(w http.ResponseWriter, r *http.Request) {
		serve(w, r, true)
	})

	// A blank page whose only job is to import every behavior module through
	// the shim, so window.__gsxuiRegistrations holds the full registry.
	mux.HandleFunc("GET /registrations", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		renderShell(w, "registrations", "/shim/index.js", template.HTML(""))
	})
}
```

- [ ] **Step 5: Register the routes**

In `jstest/harness/main.go`, inside `newMux`, immediately before `return mux`:

```go
	registerModuleRoutes(mux, root)

	return mux
```

- [ ] **Step 6: Run the tests to verify they pass**

Run: `go test ./jstest/harness/ -v`
Expected: PASS — all eleven tests.

- [ ] **Step 7: Check the shim parses**

Run: `node --check jstest/harness/shim.js`
Expected: no output, exit 0.

- [ ] **Step 8: Commit**

```bash
git add jstest/harness/
git commit -m "test(js): serve ui modules raw, plus a recording shim

/ui/*.js is the shipped source byte-for-byte — native ES modules, no
bundler between the test and the code. /shim/*.js is the same tree with a
gsxui.js that records on() calls instead of registering them, which is what
makes the delegation registry observable."
```

---

### Task 3: Playwright rig — config, global setup, CSS, and a smoke test

**Files:**
- Modify: `package.json` (add `@playwright/test` and `@tailwindcss/cli` to devDependencies)
- Modify: `.gitignore`
- Create: `jstest/support/paths.ts`
- Create: `jstest/support/manifest.ts`
- Create: `jstest/global-setup.ts`
- Create: `jstest/playwright.config.ts`
- Create: `jstest/specs/smoke.spec.ts`

**Interfaces:**
- Consumes: the harness's `-manifest` flag, `/healthz`, `/x/<component>/<example>`, and the `/static/jstest/.tmp/site.css` stylesheet path baked into the shell by Task 1.
- Produces: `paths.ts` exporting `repoRoot`, `tmpDir`, `manifestPath`, `cssPath`, `baseURL`; `manifest.ts` exporting `type ExampleEntry = { component: string; example: string; url: string }` and `examples(): ExampleEntry[]`. Tasks 4–6 import `examples()`.

- [ ] **Step 1: Install the dependencies**

Run:

```bash
npm install --save-dev @playwright/test@^1.55.0 @tailwindcss/cli@^4.3.3
npx playwright install chromium
```

`@playwright/test` does not download browsers on install — `npx playwright install` does, into `~/.cache/ms-playwright`. That is why `site/Dockerfile`'s `npm ci` is unaffected by this change.

- [ ] **Step 2: Ignore the generated artifacts**

Append to `.gitignore`:

```
jstest/.tmp/
jstest/test-results/
jstest/playwright-report/
```

- [ ] **Step 3: Write the shared paths module**

Create `jstest/support/paths.ts`:

```ts
import path from "node:path";
import { fileURLToPath } from "node:url";

const here = path.dirname(fileURLToPath(import.meta.url));

/** Repo root — jstest/support is two levels down. */
export const repoRoot = path.resolve(here, "..", "..");
export const jstestDir = path.join(repoRoot, "jstest");

/** Generated, gitignored: the example manifest and the compiled stylesheet. */
export const tmpDir = path.join(jstestDir, ".tmp");
export const manifestPath = path.join(tmpDir, "examples.json");
export const cssPath = path.join(tmpDir, "site.css");

/**
 * 7799 is deliberately clear of the dev loop's 7777 and Vite's 5173, so a
 * running `make site-dev` never collides with a test run.
 */
export const harnessPort = 7799;
export const baseURL = `http://127.0.0.1:${harnessPort}`;
```

- [ ] **Step 4: Write the manifest reader**

Create `jstest/support/manifest.ts`:

```ts
import { readFileSync } from "node:fs";
import { manifestPath } from "./paths";

export type ExampleEntry = {
  component: string;
  example: string;
  url: string;
};

/**
 * Reads the manifest Playwright's globalSetup generated. Synchronous on
 * purpose: spec files call this at module scope to generate one test per
 * example, and Playwright's test declaration is not async.
 */
export function examples(): ExampleEntry[] {
  let raw: string;
  try {
    raw = readFileSync(manifestPath, "utf8");
  } catch (err) {
    throw new Error(
      `example manifest missing at ${manifestPath}. It is written by ` +
        `jstest/global-setup.ts, which only runs under Playwright — use ` +
        `\`make test-js\` or \`npx playwright test --config jstest/playwright.config.ts\`. ` +
        `(${err})`,
    );
  }
  const entries = JSON.parse(raw) as ExampleEntry[];
  if (!Array.isArray(entries) || entries.length === 0) {
    throw new Error(`example manifest at ${manifestPath} is empty`);
  }
  return entries;
}
```

- [ ] **Step 5: Write global setup**

Create `jstest/global-setup.ts`:

```ts
import { execFileSync } from "node:child_process";
import { mkdirSync } from "node:fs";
import { cssPath, manifestPath, repoRoot, tmpDir } from "./support/paths";

/**
 * Runs once before any worker imports a spec file, which is what lets specs
 * read the manifest synchronously and still work under a bare
 * `npx playwright test`.
 *
 * The stylesheet is compiled from web/site.css — the production entry — so
 * every class a component can emit is present (site.css already declares
 * `@source "../ui/**\/*.gsx"`). Real CSS is not optional here: the ghost-box
 * invariant is a computed-style assertion, and it caught shipped defects in
 * both dialog and sidebar.
 */
export default function globalSetup() {
  mkdirSync(tmpDir, { recursive: true });

  execFileSync("go", ["run", "./jstest/harness", "-manifest", manifestPath], {
    cwd: repoRoot,
    stdio: "inherit",
  });

  execFileSync("npx", ["@tailwindcss/cli", "-i", "web/site.css", "-o", cssPath], {
    cwd: repoRoot,
    stdio: "inherit",
  });
}
```

Note the `@source` glob in the doc comment is escaped (`*.gsx` after `**\/`) so it does not close the block comment. Keep it that way.

- [ ] **Step 6: Write the Playwright config**

Create `jstest/playwright.config.ts`:

```ts
import { defineConfig, devices } from "@playwright/test";
import { baseURL, harnessPort, repoRoot } from "./support/paths";

export default defineConfig({
  testDir: "./specs",
  globalSetup: "./global-setup.ts",
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0,
  reporter: process.env.CI ? "github" : "list",
  use: {
    baseURL,
    trace: "on-first-retry",
  },
  projects: [{ name: "chromium", use: { ...devices["Desktop Chrome"] } }],
  webServer: {
    command: `go run ./jstest/harness -addr 127.0.0.1:${harnessPort} -root .`,
    cwd: repoRoot,
    url: `${baseURL}/healthz`,
    reuseExistingServer: !process.env.CI,
    timeout: 120_000,
    stdout: "pipe",
    stderr: "pipe",
  },
});
```

- [ ] **Step 7: Write the smoke test**

Create `jstest/specs/smoke.spec.ts`:

```ts
import { expect, test } from "@playwright/test";
import { examples } from "../support/manifest";

test("the manifest reaches the specs", () => {
  const all = examples();
  expect(all.length).toBeGreaterThan(100);
  expect(all).toContainEqual({
    component: "dropdown",
    example: "checkboxes",
    url: "/x/dropdown/checkboxes",
  });
});

test("markup, stylesheet and behavior JS all arrive", async ({ page }) => {
  await page.goto("/x/toggle/basic");

  // Markup: rendered through the real gsx component.
  const toggle = page.locator("[data-gsxui-toggle]").first();
  await expect(toggle).toBeVisible();

  // CSS: bg-background resolves to an opaque colour. Without the stylesheet
  // the body background is rgba(0, 0, 0, 0).
  const background = await page.evaluate(
    () => getComputedStyle(document.body).backgroundColor,
  );
  expect(background).not.toBe("rgba(0, 0, 0, 0)");

  // JS: ui/toggle.js is bound and flips both attributes on click.
  await expect(toggle).toHaveAttribute("aria-pressed", "false");
  await toggle.click();
  await expect(toggle).toHaveAttribute("aria-pressed", "true");
  await expect(toggle).toHaveAttribute("data-state", "on");
});

test("the shim records what the real modules would have registered", async ({ page }) => {
  await page.goto("/registrations");

  const registrations = await page.evaluate(() => window.__gsxuiRegistrations);
  expect(registrations.length).toBeGreaterThan(50);

  // Every entry is attributed to a real module, not "unknown" — the stack
  // walk in shim.js is the only thing that can break this.
  const modules = new Set(registrations.map((r) => r.module));
  expect(modules).not.toContain("unknown");
  expect(modules).toContain("dropdown.js");
  expect(modules).toContain("toggle.js");
});
```

Create `jstest/support/globals.d.ts` so `window.__gsxuiRegistrations` type-checks:

```ts
export type Registration = {
  type: string;
  capture: boolean;
  selector: string;
  module: string;
};

declare global {
  interface Window {
    __gsxuiRegistrations: Registration[];
  }
}
```

Import the type where needed with `import type { Registration } from "../support/globals";`.

- [ ] **Step 8: Run the smoke test**

Run: `npx playwright test --config jstest/playwright.config.ts`
Expected: PASS — 3 tests.

- [ ] **Step 9: Verify the stylesheet loads with no 404s**

This is the step that confirms the `/static/` root choice in Task 1. Tailwind bundles `@fontsource-variable/geist`, whose CSS carries `url()` references; if those resolve to paths the harness does not serve, every page load logs a failed request and Task 4's clean-load invariant will be noisy.

Run:

```bash
npx playwright test --config jstest/playwright.config.ts --debug
```

or, without the debugger, add a temporary spec that records failed requests:

```ts
test("no failed subresource requests", async ({ page }) => {
  const failed: string[] = [];
  page.on("requestfailed", (r) => failed.push(r.url()));
  page.on("response", (r) => {
    if (r.status() >= 400) failed.push(`${r.status()} ${r.url()}`);
  });
  await page.goto("/x/toggle/basic");
  expect(failed).toEqual([]);
});
```

Keep this test — it belongs in the suite permanently. Put it in `smoke.spec.ts`.

If it fails on font URLs, report the actual resolved URLs before changing anything. The two legitimate fixes, in order of preference: (a) adjust the `/static/` mount so the emitted relative paths resolve, or (b) pass `--optimize` / adjust the Tailwind input so font `url()`s are absolute. Do **not** fix it by filtering font 404s out of the clean-load invariant — that would blunt the invariant for every other resource.

- [ ] **Step 10: Commit**

```bash
git add package.json package-lock.json .gitignore jstest/
git commit -m "test(js): Playwright rig over the harness

globalSetup writes the example manifest and compiles web/site.css to a
gitignored temp file, so specs read the example list synchronously and a
bare \`npx playwright test\` works. Smoke test proves markup, stylesheet and
behavior JS all arrive."
```

---

### Task 4: Invariant — every example loads clean

**Files:**
- Create: `jstest/specs/invariants.spec.ts`

**Interfaces:**
- Consumes: `examples()` from `../support/manifest`.
- Produces: the per-example test skeleton. Tasks 5 and 6 add checks inside the same `test()` body, using `expect.soft` so one example reports every invariant it violates in a single run rather than one per re-run.

- [ ] **Step 1: Write the spec**

Create `jstest/specs/invariants.spec.ts`:

```ts
import { expect, test } from "@playwright/test";
import { examples } from "../support/manifest";

/**
 * One test per example, each asserting every invariant with expect.soft so a
 * single run reports everything an example violates.
 *
 * These four checks are where the leverage is: they cover all 103 examples
 * without a line of per-component test code, and two of them encode defect
 * classes that have actually shipped here.
 */
for (const example of examples()) {
  test(`${example.component}/${example.example}`, async ({ page }) => {
    const consoleErrors: string[] = [];
    const pageErrors: string[] = [];
    page.on("console", (msg) => {
      if (msg.type() === "error") consoleErrors.push(msg.text());
    });
    page.on("pageerror", (err) => pageErrors.push(String(err)));

    await page.goto(example.url);

    // Invariant 1: nothing threw and nothing logged an error while the
    // module graph loaded and the markup parsed.
    expect.soft(pageErrors, "uncaught exceptions").toEqual([]);
    expect.soft(consoleErrors, "console errors").toEqual([]);
  });
}
```

- [ ] **Step 2: Run it**

Run: `npx playwright test --config jstest/playwright.config.ts invariants`
Expected: 103 tests. They should pass; if any fail, the failure is a real defect in that example or component — report the console output verbatim and fix the cause, not the test.

- [ ] **Step 3: Prove the check can go red**

Temporarily append to `ui/toggle.js`:

```js
console.error("deliberate red-validation error");
```

Run: `npx playwright test --config jstest/playwright.config.ts invariants -g "toggle/basic"`
Expected: FAIL, with `console errors` in the message and the deliberate text in the diff.

Now change the line to `throw new Error("deliberate red-validation throw");` and re-run.
Expected: FAIL, this time on `uncaught exceptions`.

Remove both lines and re-run to confirm green. Do not commit either.

- [ ] **Step 4: Commit**

```bash
git add jstest/specs/invariants.spec.ts
git commit -m "test(js): clean-load invariant over every example

One test per example; zero console errors and zero uncaught exceptions.
Validated red by injecting each failure mode into ui/toggle.js."
```

---

### Task 5: Invariants — no ghost boxes, no duplicate ids

**Files:**
- Modify: `jstest/specs/invariants.spec.ts`

**Interfaces:**
- Consumes: the per-example `test()` body from Task 4.
- Produces: nothing new for later tasks.

- [ ] **Step 1: Add the checks**

In `jstest/specs/invariants.spec.ts`, after the Task 4 assertions and inside the same `test()` body:

```ts
    // Invariant 2: no ghost boxes. A closed popover must compute
    // `display: none`. An author display utility (a bare `block`/`grid`)
    // beats the UA stylesheet's closed-popover rule and leaves an invisible
    // but hit-testable box — this shipped in both dialog and sidebar, and
    // the fix is to gate the utility on `:open` (`open:grid`).
    const ghosts = await page.evaluate(() =>
      [...document.querySelectorAll("[popover]")]
        .filter((el) => !el.matches(":popover-open"))
        .filter((el) => getComputedStyle(el).display !== "none")
        .map((el) => ({
          slot: (el as HTMLElement).dataset.slot ?? null,
          id: el.id || null,
          display: getComputedStyle(el).display,
          classes: el.className,
        })),
    );
    expect.soft(ghosts, "closed popovers computing a display other than none").toEqual([]);

    // Invariant 3: no duplicate ids. One example renders alone on the page,
    // so any collision is within a single example's own markup — either the
    // example reuses an id, or a component generates one non-uniquely.
    const duplicateIds = await page.evaluate(() => {
      const counts = new Map<string, number>();
      for (const el of document.querySelectorAll("[id]")) {
        counts.set(el.id, (counts.get(el.id) ?? 0) + 1);
      }
      return [...counts].filter(([, n]) => n > 1).map(([id, n]) => ({ id, count: n }));
    });
    expect.soft(duplicateIds, "duplicate element ids").toEqual([]);
```

- [ ] **Step 2: Run it**

Run: `npx playwright test --config jstest/playwright.config.ts invariants`
Expected: 103 tests.

Failures here are findings, not noise. Triage each one:
- A **ghost box** is a component bug. Fix it by gating the offending `display` utility on `:open`, following the idiom already used in `ui/sidebar.gsx` (`group-data-[collapsible=icon]:open:block`).
- A **duplicate id inside one example's `.gsx`** is an example bug — fix the example, and remember `make highlight` in the same commit if you touch anything under `site/examples/`.
- A **component that generates duplicate ids** is out of scope for this plan. Do not paper over it: record it in `docs/jsx-parity.md` under that component and report it in the task summary.

Report what you found either way, including "nothing failed" if that is the outcome.

- [ ] **Step 3: Prove both checks can go red**

Ghost box — temporarily edit `ui/tooltip.gsx`'s content element, changing its `open:block` (or equivalent `:open`-gated display utility) to a bare `block`, then:

Run: `go tool gsx generate && npx playwright test --config jstest/playwright.config.ts invariants -g "tooltip"`
Expected: FAIL on `closed popovers computing a display other than none`, listing the tooltip content element.

Revert, `go tool gsx generate`, confirm green.

Duplicate id — temporarily add a second element with an existing id to `site/examples/toggle/basic.gsx`, then:

Run: `go tool gsx generate && npx playwright test --config jstest/playwright.config.ts invariants -g "toggle/basic"`
Expected: FAIL on `duplicate element ids`.

Revert, `go tool gsx generate`, confirm green. Do not commit either edit, and do not leave a stale `.x.go` behind — `git status --porcelain -- '*.x.go'` must be clean.

- [ ] **Step 4: Commit**

```bash
git add jstest/specs/invariants.spec.ts
git commit -m "test(js): ghost-box and duplicate-id invariants

A closed popover must compute display:none — an ungated display utility
beats the UA rule and leaves a hit-testable invisible box, which shipped in
both dialog and sidebar. Both checks validated red."
```

---

### Task 6: Invariant — selector disjointness across modules

**Files:**
- Create: `jstest/support/fixtures.ts`
- Create: `jstest/support/selector-allowlist.ts`
- Modify: `jstest/specs/invariants.spec.ts`
- Modify: `jstest/specs/smoke.spec.ts` (switch its `test` import to the extended one)

**Interfaces:**
- Consumes: `window.__gsxuiRegistrations` from Task 2's `/registrations` page; `Registration` from `../support/globals`.
- Produces: `jstest/support/fixtures.ts` exporting an extended `test` with a worker-scoped `registrations: Registration[]` fixture, and re-exporting `expect`. Phase 2's specs import `test` from here.

- [ ] **Step 1: Write the registrations fixture**

Create `jstest/support/fixtures.ts`:

```ts
import { test as base } from "@playwright/test";
import type { Registration } from "./globals";

/**
 * `registrations` is the full delegation registry, recorded once per worker
 * by loading /registrations — a blank page whose only script is
 * /shim/index.js, where on() records instead of binding.
 *
 * Worker-scoped rather than collected in globalSetup because globalSetup's
 * ordering against webServer startup is not something to depend on: a
 * fixture runs after the server is definitely up.
 */
export const test = base.extend<{}, { registrations: Registration[] }>({
  registrations: [
    async ({ browser }, use) => {
      const page = await browser.newPage();
      await page.goto("/registrations");
      const recorded = await page.evaluate(() => window.__gsxuiRegistrations);
      await page.close();
      if (!recorded || recorded.length === 0) {
        throw new Error("/registrations recorded nothing — the shim did not run");
      }
      await use(recorded);
    },
    { scope: "worker" },
  ],
});

export { expect } from "@playwright/test";
```

- [ ] **Step 2: Write the allowlist**

Create `jstest/support/selector-allowlist.ts`:

```ts
/**
 * Reviewed exceptions to selector disjointness.
 *
 * An entry says: these two modules may both claim the same element for the
 * same (type, capture) pair, and here is why that is correct. It is a
 * decision someone made and a reviewer can challenge — which is why this is
 * an explicit list rather than a tolerance threshold.
 *
 * Empty is the expected steady state. Adding an entry needs a reason a
 * reviewer would accept, not "the test was failing".
 */
export type AllowedOverlap = {
  /** Module filenames, sorted, e.g. ["dialog.js", "sheet.js"]. */
  modules: [string, string];
  /** "type:capture", e.g. "click:false". */
  key: string;
  reason: string;
};

export const allowedOverlaps: AllowedOverlap[] = [];
```

- [ ] **Step 3: Add the check**

In `jstest/specs/invariants.spec.ts`, replace the `@playwright/test` import:

```ts
import { expect, test } from "../support/fixtures";
```

Add to the imports:

```ts
import { allowedOverlaps } from "../support/selector-allowlist";
```

Change the test signature to take the fixture:

```ts
  test(`${example.component}/${example.example}`, async ({ page, registrations }) => {
```

Then append inside the same `test()` body:

```ts
    // Invariant 4: selector disjointness. ui/gsxui.js keys its registry by
    // `${type}:${capture}` alone and dispatches to EVERY handler whose
    // selector matches, so if two modules both match one element for one
    // (type, capture) pair, both handlers run on a single event. That is
    // exactly the hook-prefix collision that shipped in Tier 4 Batch B:
    // dropdown.js and context-menu.js both matched data-gsxui-menu-*, so
    // one click on a checkbox item fired two gsxui:change events and left
    // the state unchanged.
    //
    // The check is same-element, not ancestor-chain, and that is deliberate.
    // gsxui.js dispatches via target.closest(selector); when one element
    // matches two modules' selectors, an event on it resolves to that same
    // element for both. Nested elements matching different modules is
    // ordinary composition (a dialog inside a dropdown) and is not a defect.
    const overlaps = await page.evaluate((regs) => {
      const found: { key: string; tag: string; modules: string[]; selectors: string[] }[] = [];
      for (const el of document.querySelectorAll("*")) {
        const byKey = new Map<string, Map<string, string>>();
        for (const reg of regs) {
          let matches = false;
          try {
            matches = el.matches(reg.selector);
          } catch {
            continue; // an unparseable selector is Task 2's problem, not this one
          }
          if (!matches) continue;
          const key = `${reg.type}:${reg.capture}`;
          if (!byKey.has(key)) byKey.set(key, new Map());
          byKey.get(key)!.set(reg.module, reg.selector);
        }
        for (const [key, mods] of byKey) {
          if (mods.size > 1) {
            found.push({
              key,
              tag: el.tagName.toLowerCase(),
              modules: [...mods.keys()].sort(),
              selectors: [...mods.values()],
            });
          }
        }
      }
      return found;
    }, registrations);

    const unexpected = overlaps.filter(
      (o) =>
        !allowedOverlaps.some(
          (a) =>
            a.key === o.key &&
            o.modules.length === 2 &&
            a.modules[0] === o.modules[0] &&
            a.modules[1] === o.modules[1],
        ),
    );
    expect.soft(unexpected, "elements claimed by two modules for one event").toEqual([]);
```

- [ ] **Step 4: Update the smoke spec's imports**

In `jstest/specs/smoke.spec.ts`, change:

```ts
import { expect, test } from "@playwright/test";
```

to:

```ts
import { expect, test } from "../support/fixtures";
```

The existing smoke tests do not use the `registrations` fixture, so they are unaffected — worker fixtures are lazy and only initialise when a test destructures them.

- [ ] **Step 5: Run it**

Run: `npx playwright test --config jstest/playwright.config.ts`
Expected: PASS.

If real overlaps turn up, they are findings. Report each with its modules, key and selectors. An overlap is a **bug to fix by namespacing the hook attribute** unless you can articulate why both handlers firing is correct — in which case add an allowlist entry with that reasoning. Default to fixing.

- [ ] **Step 6: Prove the check can go red**

This is the most important red-validation in the plan, because it is the check that would have caught a defect that shipped.

Temporarily edit `ui/context-menu.js` and change one of its item selectors to `dropdown`'s prefix — e.g. change the constant matching `[data-gsxui-contextmenu-item]` to `[data-gsxui-dropdown-item]`.

Run: `npx playwright test --config jstest/playwright.config.ts invariants -g "dropdown"`
Expected: FAIL on `elements claimed by two modules for one event`, naming `context-menu.js` and `dropdown.js`, the shared key, and both selectors.

Now temporarily add the matching allowlist entry and re-run.
Expected: PASS — proving the allowlist path works.

Revert both edits and confirm green.

- [ ] **Step 7: Commit**

```bash
git add jstest/support/ jstest/specs/
git commit -m "test(js): selector-disjointness invariant

Records the delegation registry through a shimmed gsxui.js, then checks on
every example page that no element is claimed by two modules for the same
(type, capture) pair — the mechanical form of the hook-prefix collision.
Validated red by re-introducing that collision in context-menu.js."
```

---

### Task 7: Wire it into make, CI, and the docs

**Files:**
- Modify: `Makefile`
- Modify: `.github/workflows/fly-deploy.yml`
- Modify: `README.md`
- Modify: `CHANGELOG.md`
- Modify: `docs/jsx-parity.md`

**Interfaces:**
- Consumes: everything above.
- Produces: `make test-js`, and a CI step that gates deploys on it.

- [ ] **Step 1: Add the make targets**

In `Makefile`, add `test-js` to the `.PHONY` line, then add after the `test` target:

```makefile
# test-js runs the Playwright suite in jstest/ against real Chromium. It is
# deliberately NOT part of `make test` — a Go-only edit should not pay for a
# browser boot. Playwright's globalSetup writes the example manifest and
# compiles web/site.css into jstest/.tmp/ (gitignored), and its webServer
# block starts jstest/harness on 127.0.0.1:7799 — clear of the dev loop's
# 7777 and Vite's 5173, so `make site-dev` can keep running.
#
# First run on a new machine needs the browser: `npx playwright install chromium`.
test-js:
	npx playwright test --config jstest/playwright.config.ts
```

Change the `check` target's first line from `check: test` to:

```makefile
check: test test-js
```

and extend its `node --check` loop to cover the harness shim:

```makefile
	@for f in $$(find ui jstest -name '*.js'); do node --check $$f || exit 1; done
```

- [ ] **Step 2: Verify the targets**

Run: `make test-js`
Expected: PASS, all specs.

Run: `make check`
Expected: PASS — Go tests, JS tests, generated-file drift check, `node --check`, gofmt.

- [ ] **Step 3: Wire CI**

In `.github/workflows/fly-deploy.yml`, replace the `test` job with:

```yaml
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - uses: actions/setup-node@v4
        with:
          node-version: 22
          cache: npm
      - run: npm ci
      # Browser binaries live outside node_modules, so npm's cache doesn't
      # cover them. Key on the lockfile: a Playwright version bump invalidates.
      - uses: actions/cache@v4
        id: playwright-cache
        with:
          path: ~/.cache/ms-playwright
          key: playwright-${{ runner.os }}-${{ hashFiles('package-lock.json') }}
      - if: steps.playwright-cache.outputs.cache-hit != 'true'
        run: npx playwright install --with-deps chromium
      - run: go test ./...
      - run: make test-js
      - if: failure()
        uses: actions/upload-artifact@v4
        with:
          name: playwright-report
          path: jstest/playwright-report/
          retention-days: 7
```

`deploy` already declares `needs: test`, so it now gates on both suites with no new job.

Also update the workflow's header comment — it currently says the Go test suite is the gate — to say the Go suite *and* the Playwright suite.

- [ ] **Step 4: Document it**

In `README.md`, replace the line:

```
- `make test` regenerates and tests everything; `make check` adds JS syntax
  and gofmt checks.
```

with:

```
- `make test` regenerates and runs the Go suite; `make test-js` runs the
  browser suite (Playwright against real Chromium — `npx playwright install
  chromium` once per machine); `make check` runs both plus JS syntax and
  gofmt checks. CI gates deploys on both suites.
- Component behavior (`ui/*.js`) is tested in `jstest/`: a Go harness serves
  one example per page with the real modules loaded as native ES modules,
  and four invariants sweep every example — clean load, no ghost popovers,
  no duplicate ids, and no element claimed by two modules for one event.
```

In `CHANGELOG.md`, under a new `## 2026-07-25` heading if the existing one has already been released, otherwise into the existing one, add under `### Added`:

```
- **JS test layer** — Playwright suite in `jstest/` over a Go example harness; four invariants sweep every example, gated in CI.
```

Keep it one line. The changelog is terse by house style.

In `docs/jsx-parity.md`, add to the MECHANISM note about `gsxui.js`'s registry dispatching to all matching handlers:

```
This is now enforced: jstest/specs/invariants.spec.ts records the registry
through a shimmed gsxui.js and fails if any element on any example page is
matched by two modules' selectors for the same (type, capture) pair.
```

- [ ] **Step 5: Verify nothing drifted**

Run: `make check && git status --porcelain`
Expected: `make check` passes; the only dirty paths are the files this task edits. No `.x.go` changes, no untracked files outside `jstest/.tmp/` (which is gitignored).

- [ ] **Step 6: Commit**

```bash
git add Makefile .github/workflows/fly-deploy.yml README.md CHANGELOG.md docs/jsx-parity.md
git commit -m "test(js): gate CI on the browser suite

make test-js runs Playwright; make check runs both suites. CI's test job
gains Node, a cached Chromium, and the JS step — deploy already needed
test, so it now gates on both."
```

- [ ] **Step 7: Push and watch CI**

```bash
git push
gh run watch
```

Expected: both the `test` and `deploy` jobs green. If the Playwright step fails only in CI, download the uploaded `playwright-report` artifact before changing anything — a CI-only failure is usually a timing assumption, and the trace will show which.

---

## Self-Review

**Spec coverage**

| Spec section | Task |
|---|---|
| §1 real browser, Chromium only | Global Constraints, Task 3 config |
| §2 harness routes `/x/`, `/ui/`, `/shim/`, `/static/`, `-manifest` | Tasks 1–2 |
| §2 `/f/<name>` fixture route | Phase 2 — declared in Phase scope, not dropped |
| §3 CSS via `@tailwindcss/cli` to a temp file | Task 3 steps 1, 5, 9 |
| §4 `jstest/fixtures/` | Phase 2 — same |
| §5 invariant 1 clean load | Task 4 |
| §5 invariant 2 no ghost boxes | Task 5 |
| §5 invariant 3 selector disjointness | Task 6 |
| §5 invariant 4 no duplicate ids | Task 5 |
| §5 allowlist | Task 6 step 2 |
| §5 the two regression specs | Phase 2 |
| §6 no global animation disabling | Global Constraints |
| §7 config, globalSetup, webServer, make, CI, Dockerfile note, TypeScript | Tasks 3, 7 |
| §8 phases | Phase scope |
| §9 out of scope | not implemented, correctly |

**Type consistency:** `entry`/`ExampleEntry` share the JSON keys `component`/`example`/`url`. `Registration` is declared once in `support/globals.d.ts` and consumed by `fixtures.ts` and the disjointness check. `newMux(root string)` has one signature, used by `main` and by every Go test. `renderShell(w, title, script, body)` is called from two places with matching argument order.

**Known risk, flagged rather than hidden:** Task 3 step 9 exists because Tailwind's `url()` rewriting for bundled `@fontsource` imports has not been verified against this output path. The step makes the verification explicit and names the two acceptable fixes; it forbids the tempting-but-wrong one (filtering 404s out of the clean-load invariant).

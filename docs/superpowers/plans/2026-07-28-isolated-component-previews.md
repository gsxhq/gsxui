# Isolated Component Previews Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Render Sidebar examples in same-origin live iframes so their fixed application-shell layout cannot escape the component gallery.

**Architecture:** Add exact registered-example/preview lookup and a minimal `/examples/{component}/{example}` document that reuses the site's canonical head/assets. Extend example presentation metadata with `Isolated` and optional named `Previews`; the component page renders only those examples as iframes while ordinary examples remain inline.

**Tech Stack:** Go 1.26, GSX, structpages, standard `net/http`, Vite, vanilla JavaScript, Playwright.

## Global Constraints

- Do not change `ui.Sidebar` or add a component preview/embedded API.
- Do not use CSS containment to alter Sidebar's production fixed-position semantics.
- Preview routes resolve registered component/example pairs exactly and return 404 for unknown pairs.
- The iframe document uses the same `web/main.js` entry and theme bootstrap as the main site.
- Copyable highlighted GSX source stays on the parent component page.
- Ordinary examples keep their existing inline rendering.
- Authored `.gsx` files are the source of truth; regenerate `.x.go` files with the reviewed local GSX compiler and never hand-edit generated files.

---

### Task 1: Add the isolated example document

**Files:**
- Modify: `site/examples/registry.go`
- Modify: `site/examples/examples_test.go`
- Modify: `site/pages/pages.go`
- Modify: `site/pages/layout.gsx`
- Create: `site/pages/example_preview.gsx`
- Modify: `site/pages/pages_test.go`
- Generated: `site/pages/layout.x.go`
- Generated: `site/pages/example_preview.x.go`

**Interfaces:**
- Produces: `examples.Find(component string, name string, query url.Values) (string, gsx.Node, bool)`.
- Produces: `ExamplePreview` at `GET /examples/{component}/{example}`.
- Produces: package-private `siteHead(title string)` used by both `Layout` and `ExamplePreview`.
- Consumes: existing `Example.Node` and optional `Example.Query`.

- [ ] **Step 1: Write failing registry lookup tests**

Add table-driven tests proving exact registered lookup, query-driven rendering,
and rejection of unknown component/example pairs:

```go
func TestFindResolvesExactRegisteredExample(t *testing.T) {
    title, node, ok := Find("button", "basic", nil)
    if !ok || title != "Basic" || node == nil {
        t.Fatalf("Find(button, basic) = %q, %#v, %t", title, node, ok)
    }
    if _, _, ok := Find("button", "missing", nil); ok {
        t.Fatal("Find accepted an unregistered example")
    }
}
```

Include a temporary registered example whose `Query` returns a distinctive
node and render that returned node to prove `Find` chooses the request hook
rather than the static node.

- [ ] **Step 2: Write failing preview route tests**

In `site/pages/pages_test.go`, request:

```text
/examples/sidebar/basic
/examples/sidebar/missing
/examples/missing/basic
```

Assert that the valid response is 200 and contains the Sidebar wrapper,
desktop container, `web/main.js` development entry, viewport-sized preview
body marker, and no docs navigation or iframe. Assert both invalid responses
are 404.

- [ ] **Step 3: Run the focused tests and verify RED**

Run:

```bash
go test ./site/examples ./site/pages -run 'Test(Find|ExamplePreview)' -count=1
```

Expected: compilation/test failures because `Find`, `ExamplePreview`, and its
route do not exist.

- [ ] **Step 4: Implement exact example lookup**

Add:

```go
func Find(component string, name string, query url.Values) (string, gsx.Node, bool) {
    for _, example := range registry[component] {
        if example.Name != name {
            continue
        }
        node := example.Node
        if example.Query != nil {
            node = example.Query(query)
        }
        return example.Title, node, true
    }
    return "", nil, false
}
```

The lookup must inspect only `registry[component]`; do not normalize, infer,
or reconstruct keys.

- [ ] **Step 5: Extract the canonical document head**

Move the complete existing `<head>` from `Layout` into:

```gsx
component siteHead(title string) {
    <head>
        ...
    </head>
}
```

Keep the current theme bootstrap, development FOUC gate, favicon links, and
`vite.FromContext(ctx).Entry("web/main.js")` asset loops byte-for-byte in
behavior. Replace `Layout`'s head with `<siteHead title={title}/>` and make no
other layout change.

- [ ] **Step 6: Implement the isolated page**

Add `ExamplePreview` with Props that reads the two path values and calls
`examples.Find`. Return `ErrorWithStatus{Status: http.StatusNotFound}` when
the pair is unknown.

Render:

```gsx
<!DOCTYPE html>
<html lang="en">
    <siteHead title={props.Title}/>
    <body
        data-site-isolated-document
        class="h-svh overflow-auto bg-background text-foreground antialiased"
    >
        { props.Node }
    </body>
</html>
```

Register it on `Pages` as:

```go
ExamplePreview `route:"/examples/{component}/{example} Example Preview"`
```

- [ ] **Step 7: Regenerate and verify GREEN**

Create an untracked temporary Go workspace selecting this gsxui worktree and
`/Users/jackieli/personal/gsxhq/gsx/.worktrees/preserve-bare-fallthrough`,
then run:

```bash
GOWORK="$preview_workspace/go.work" go run github.com/gsxhq/gsx/cmd/gsx generate
GOWORK="$preview_workspace/go.work" go test ./site/examples ./site/pages -run 'Test(Find|ExamplePreview)' -count=1
GOWORK="$preview_workspace/go.work" gopls check -severity=hint \
  site/examples/registry.go site/pages/pages.go site/pages/layout.x.go \
  site/pages/example_preview.x.go
```

Expected: focused tests pass and gopls reports no diagnostics.

- [ ] **Step 8: Commit Task 1**

```bash
git add site/examples/registry.go site/examples/examples_test.go \
  site/pages/pages.go site/pages/layout.gsx site/pages/layout.x.go \
  site/pages/example_preview.gsx site/pages/example_preview.x.go \
  site/pages/pages_test.go
git commit -m "feat: add isolated example documents"
```

---

### Task 2: Present Sidebar examples in live iframes

**Files:**
- Modify: `site/examples/registry.go`
- Modify: `site/examples/sidebar.go`
- Modify: `site/pages/component.gsx`
- Modify: `site/pages/pages_test.go`
- Modify: `web/site.js`
- Modify: `jstest/harness/main.go`
- Modify: `jstest/harness/manifest.go`
- Modify: `jstest/specs/style-visual.spec.ts`
- Create: `jstest/specs/sidebar-page.spec.ts`
- Generated: `site/pages/component.x.go`

**Interfaces:**
- Consumes: `ExamplePreview` route from Task 1.
- Produces: `Example.Isolated bool`.
- Produces: optional exact named `Example.Previews`.
- Produces: iframe marker `data-site-isolated-preview`.
- Produces: iframe URLs `/examples/{component}/{example}` and
  `/examples/{component}/{example}?_preview={preview}`.

- [ ] **Step 1: Write the failing parent-page tests**

Request `/components/sidebar`, parse the returned HTML, and assert:

- exactly ten `iframe[data-site-isolated-preview]` elements;
- exact sources for Basic, eight named Variants cases, and Persisted;
- each iframe has a non-empty title;
- the parent document contains no
  `data-gsxui-slot-sidebar-container` element;
- each existing highlighted source block and Copy control remains.

Retain the existing Button route assertion proving an ordinary example still
renders `data-gsxui-slot-button` inline and contains no isolated iframe.

- [ ] **Step 2: Write the failing browser regression**

Create `jstest/specs/sidebar-page.spec.ts`. The test must:

1. open `/components/sidebar`;
2. assert zero parent-document Sidebar containers;
3. scope to the iframe with title `Basic preview`;
4. assert its wrapper rect stays inside the iframe viewport;
5. assert its desktop container computes to `position: fixed` with `left: 0`;
6. click its unique `Toggle Sidebar` trigger;
7. assert the iframe wrapper changes from `data-state="expanded"` to
   `data-state="collapsed"`;
8. click the parent page's unique theme toggle;
9. assert the iframe document's root receives the parent's `dark` state.

Also assert the parent docs navigation still contains the Sidebar component
link and its bounding box remains inside the parent page's normal sidebar
column.

- [ ] **Step 3: Run the focused tests and verify RED**

Run:

```bash
go test ./site/pages -run TestComponentPageRoute -count=1
npx playwright test jstest/specs/sidebar-page.spec.ts --config jstest/playwright.config.ts
```

Expected: the Go test sees inline Sidebar containers and the browser test
sees no isolated iframe.

- [ ] **Step 4: Add isolation metadata and iframe rendering**

Add `Isolated bool` and optional named `Previews` to `examples.Example`; set
Sidebar's three registrations to isolated. Register the eight Variants cases
as exact named previews sharing the existing Variants source block.

Extend `exampleProps` with `Name` and `Isolated`. In `Component.Props`, copy
those fields exactly from the registered example.

In `Component.Page`, render a single iframe when `Previews` is empty, or one
labelled full-width iframe per named preview when it is non-empty:

```gsx
{ if ex.Isolated {
    <iframe
        data-site-isolated-preview
        title={ex.Title + " preview"}
        src={"/examples/" + props.Name + "/" + ex.Name}
        loading="lazy"
        class="block h-[32rem] w-full rounded-lg border bg-background"
    ></iframe>
} else {
    <div class="border rounded-lg p-8 bg-background">
        { ex.Node }
    </div>
} }
```

Do not wrap the iframe in the existing padded preview panel.

Refactor `site/examples/sidebar/variants.gsx` so every named preview node
contains exactly one SidebarProvider. Do not mount multiple fixed Sidebar
containers in a single isolated document.

Expand named previews into unique browser-manifest entries and route the
harness through `examples.Find`. Replace the old composite Sidebar visual
snapshot with one desktop snapshot per named case; keep one meaningful mobile
snapshot by opening the default case's native Sheet branch.

- [ ] **Step 5: Synchronize live theme changes**

After toggling and persisting the parent theme, update loaded same-origin
preview documents:

```js
for (const frame of document.querySelectorAll(
  "iframe[data-site-isolated-preview]",
)) {
  try {
    frame.contentDocument?.documentElement.classList.toggle("dark", dark);
  } catch {
    // A frame that is not loaded or not same-origin applies storage on load.
  }
}
```

Do not add cross-window messages or another storage abstraction.

- [ ] **Step 6: Regenerate and verify GREEN**

Using the same temporary exact-core workspace:

```bash
GOWORK="$preview_workspace/go.work" go run github.com/gsxhq/gsx/cmd/gsx generate
GOWORK="$preview_workspace/go.work" go test ./site/pages -run TestComponentPageRoute -count=1
npx playwright test jstest/specs/sidebar-page.spec.ts --config jstest/playwright.config.ts
GOWORK="$preview_workspace/go.work" gopls check -severity=hint \
  site/examples/registry.go site/examples/sidebar.go \
  site/pages/component.x.go site/pages/pages_test.go
```

Expected: route and browser regressions pass with no diagnostics.

- [ ] **Step 7: Run authoritative verification**

Run:

```bash
GOWORK="$preview_workspace/go.work" make ci
git diff --check
git status --short
```

Expected: generated drift, Go/vet/format checks, CSS audit, and every
discovered Chromium test pass. Only intended task files are modified.

- [ ] **Step 8: Commit Task 2**

```bash
git add site/examples/registry.go site/examples/sidebar.go \
  site/pages/component.gsx site/pages/component.x.go site/pages/pages_test.go \
  web/site.js jstest/specs/sidebar-page.spec.ts
git commit -m "fix: isolate sidebar component previews"
```

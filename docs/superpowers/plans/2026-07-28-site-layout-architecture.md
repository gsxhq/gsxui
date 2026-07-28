# Site Layout Architecture Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Give gsxui distinct marketing, documentation, and theme-editor workspace shells modeled on the current shadcn site.

**Architecture:** A package-private `siteLayout` accepts a named `layoutMode` and optional server-rendered TOC items. Documentation uses a wide three-column shell with a 640px article; `/theme` uses a viewport workspace whose compact palette customizer and one iframe preview manage their own overflow.

**Tech Stack:** Go 1.26.1, GSX, Tailwind CSS 4.3.3, browser-native JavaScript, Playwright, structpages, Vite.

## Global Constraints

- Execute after the palette-model plan so the final customizer structure is known.
- Work only in `/Users/jackieli/personal/gsxhq/gsxui/.worktrees/css-only-theme-architecture`.
- Do not stop or replace the user's Vite process on port 5173.
- Do not hand-edit generated `.x.go`; run `go tool gsx generate`.
- Keep `site/dist/.gitkeep`; restore it after production builds.
- Layout modes are explicit named values, never inferred from the route.
- TOC IDs are explicit source identifiers; do not slugify display text.
- The TOC is server-rendered and usable without JavaScript.
- `/theme` contains exactly one preview iframe at every viewport width.
- Existing palette, transport, handshake, Retry, and manual-copy behavior must not change.
- Every task uses red-green TDD and ends in its own commit.

---

## File map

- `site/pages/layout.gsx`: shared chrome, named modes, wide docs rails, workspace footer/overflow.
- `site/pages/toc.gsx`: TOC model, headings, and right rail.
- Page `.gsx` files: select modes and supply shared heading data.
- `web/site.js`: active TOC link enhancement.
- `site/pages/pages_test.go`: route structure, modes, TOC integrity, and iframe invariants.
- `jstest/specs/site-layout.spec.ts`: responsive geometry and navigation.
- Matching `.x.go`: generated artifacts only.

### Task 1: Named layout modes and wide docs foundation

**Files:**
- Modify: `site/pages/layout.gsx`
- Modify: `site/pages/{home,components_index,component,getting_started,theming,theme}.gsx`
- Modify generated: matching `site/pages/*.x.go`
- Test: `site/pages/pages_test.go`

**Interfaces:**
- Produces: `type layoutMode string`
- Produces: `layoutMarketing`, `layoutDocs`, `layoutWorkspace`
- Produces: `component siteLayout(title string, active string, mode layoutMode, toc []docTOCItem, children gsx.Node)`

- [ ] **Step 1: Write failing route-mode tests**

Add `TestSiteLayoutModes` requiring:

```go
tests := []struct{ path, mode string }{
	{"/", "marketing"},
	{"/components", "docs"},
	{"/components/button", "docs"},
	{"/docs/getting-started", "docs"},
	{"/theme", "workspace"},
}
```

Require `data-site-layout="<mode>"`; require docs sidebar/article markers only
on docs routes; require no docs rails or `data-site-footer` on `/theme`.

- [ ] **Step 2: Run the test and verify missing markers fail**

```bash
go test ./site/pages -run TestSiteLayoutModes -count=1
```

Expected: FAIL against the one-size `Layout`.

- [ ] **Step 3: Implement the named layout API**

Define:

```go
type layoutMode string

const (
	layoutMarketing layoutMode = "marketing"
	layoutDocs      layoutMode = "docs"
	layoutWorkspace layoutMode = "workspace"
)
```

Rename `Layout` to package-private `siteLayout`. Resolve mode-specific class
variables before tags. Emit `data-site-layout` on the outer shell.

Marketing is centered without docs rails. Docs uses a full-width
`max-w-[1568px]` container, sticky independently scrollable left navigation,
and a grid with a 288px left column from `lg`, a 640px-centered flexible
article column, and a 288px right column from `xl`. At 1280px this yields the
audited 288 / 640 / 288 structure with 32px outer gutters. The article is:

```gsx
<main data-site-docs-article class="mx-auto w-full min-w-0 max-w-[640px]">
	{ children }
</main>
```

Workspace omits both navigation rails and the footer.

- [ ] **Step 4: Update all page call sites**

Use `layoutMarketing` for Home, `layoutDocs` for component and guide routes,
and `layoutWorkspace` for Theme. Pass `nil` TOC until Task 2.

- [ ] **Step 5: Generate and run focused route checks**

```bash
go tool gsx generate
go test ./site/pages -run 'Test(SiteLayoutModes|SiteRoutes|ComponentPageRoute|DocsRoutes|ThemePageRoute)$' -count=1
go run ./internal/generatedcheck/cmd
```

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add site/pages
git commit -m "refactor: add explicit site layout modes"
```

### Task 2: Server-rendered TOC and active-section enhancement

**Files:**
- Create: `site/pages/toc.gsx`
- Create generated: `site/pages/toc.x.go`
- Modify: `site/pages/layout.gsx`
- Modify: `site/pages/component.gsx`
- Modify: `site/pages/getting_started.gsx`
- Modify: `site/pages/theming.gsx`
- Modify generated: matching `.x.go`
- Modify: `web/site.js`
- Test: `site/pages/pages_test.go`
- Create: `jstest/specs/site-layout.spec.ts`

**Interfaces:**
- Produces: `type docTOCItem struct { ID string; Title string; Depth int }`
- Produces: `component docHeading(item docTOCItem)`
- Produces: `component docTableOfContents(items []docTOCItem)`
- Produces: `func componentTOCItems(examples []exampleProps) []docTOCItem`
- Consumes: `siteLayout(... toc []docTOCItem ...)`

- [ ] **Step 1: Write failing server integrity tests**

Add `TestDocsTableOfContents`. Parse `/components/button`,
`/docs/getting-started`, and `/docs/theming` with `x/net/html`. For every
`a[data-site-toc-link][href^="#"]`, require exactly one
`[data-doc-heading][id]` with matching normalized text. Reject duplicate IDs.
Assert `/components` has no `data-site-docs-toc`.

- [ ] **Step 2: Run and verify the missing TOC failure**

```bash
go test ./site/pages -run TestDocsTableOfContents -count=1
```

Expected: FAIL because headings have no shared IDs or rail.

- [ ] **Step 3: Implement shared server-rendered headings and rail**

In `toc.gsx`:

```gsx
component docHeading(item docTOCItem) {
	{ if item.Depth == 3 {
		<h3 id={item.ID} data-doc-heading tabindex="-1">{ item.Title }</h3>
	} else {
		<h2 id={item.ID} data-doc-heading tabindex="-1">{ item.Title }</h2>
	} }
}
```

`docTableOfContents` renders a labelled navigation landmark and
`a[data-site-toc-link]` links. In `layout.gsx`, show a sticky right rail from
`xl` only when `len(toc) > 0`.

- [ ] **Step 4: Supply one heading source per page**

Derive component items from registry examples:

```go
func componentTOCItems(examples []exampleProps) []docTOCItem {
	items := make([]docTOCItem, len(examples))
	for i, example := range examples {
		items[i] = docTOCItem{
			ID: "example-" + example.Name, Title: example.Title, Depth: 2,
		}
	}
	return items
}
```

Static guides define ordered item slices with explicit IDs and render both
heading and rail from the same indexed values.

- [ ] **Step 5: Add active-link behavior**

In `web/site.js`, map existing hashes to existing heading IDs, initialize the
first/hash target, and use `IntersectionObserver` with
`rootMargin: "-80px 0px -70% 0px"`. The active link receives
`data-active` and `aria-current="location"`; all others lose both.
Recompute from an existing target on `hashchange`; do not create or rewrite
hashes in JavaScript.

- [ ] **Step 6: Write responsive Playwright tests**

Create `site-layout.spec.ts` requiring at 1280px:

```ts
await expect(page.locator("[data-site-docs-sidebar]")).toBeVisible();
await expect(page.locator("[data-site-docs-toc]")).toBeVisible();
expect(await page.locator("[data-site-docs-article]")
  .evaluate((el) => el.getBoundingClientRect().width)).toBe(640);
```

Require TOC hidden at 1100px, both rails hidden at 900px, no TOC on
`/components`, and hash navigation/active state on a component section.

- [ ] **Step 7: Generate, test, and commit**

```bash
go tool gsx generate
go test ./site/pages -run 'Test(SiteLayoutModes|DocsTableOfContents|ComponentPageRoute|DocsRoutes)$' -count=1
npx playwright test --config jstest/playwright.config.ts jstest/specs/site-layout.spec.ts
node --check web/site.js
go run ./internal/generatedcheck/cmd
git add site/pages web/site.js jstest/specs/site-layout.spec.ts
git commit -m "feat: add responsive documentation navigation"
```

### Task 3: Theme editor viewport workspace

**Files:**
- Modify: `site/pages/layout.gsx`
- Modify: `site/pages/theme.gsx`
- Modify generated: `site/pages/layout.x.go`, `site/pages/theme.x.go`
- Modify: `site/pages/pages_test.go`
- Modify: `jstest/specs/site-layout.spec.ts`
- Test existing: `jstest/specs/theme-editor.spec.ts`

**Interfaces:**
- Consumes: palette picker regions from the palette-model plan
- Produces: `data-theme-style-panel`, `data-theme-preview-panel`,
  `data-theme-controls-panel`
- Preserves every existing `data-theme-*` behavior hook

- [ ] **Step 1: Write failing workspace structure tests**

Require the three panel markers, exactly one preview frame, workspace mode, no
docs rails, and no footer in `TestThemePageRoute`.

- [ ] **Step 2: Write failing desktop and narrow geometry tests**

At 1280x900 require no document overflow, controls narrower than preview, and
independent controls scrolling. At 900px require normal document scrolling,
Y-order style `<` preview `<` controls, and one iframe.

- [ ] **Step 3: Run and verify current geometry failures**

```bash
go test ./site/pages -run TestThemePageRoute -count=1
npx playwright test --config jstest/playwright.config.ts jstest/specs/site-layout.spec.ts
```

Expected: FAIL against the long page and `xl` sticky preview.

- [ ] **Step 4: Implement the workspace grid**

Make workspace mode `lg:h-svh lg:overflow-hidden`; give its main
`lg:min-h-0 lg:flex-1`; omit footer.

Use DOM order style, preview, controls and desktop placement:

```gsx
<section data-theme-style-panel>...</section>
<div data-theme-preview-panel class="lg:col-start-2 lg:row-span-2 lg:row-start-1 lg:min-h-0">...</div>
<div data-theme-controls-panel class="lg:col-start-1 lg:row-start-2 lg:min-h-0 lg:overflow-y-auto">...</div>
```

Use `lg:grid-cols-[minmax(22rem,28rem)_minmax(0,1fr)]`. Give the iframe a
bounded mobile height and `lg:min-h-0 lg:flex-1`, removing its fixed desktop
640px minimum.

- [ ] **Step 5: Generate and verify behavior plus geometry**

```bash
go tool gsx generate
go test ./site/pages -run 'Test(SiteLayoutModes|ThemePageRoute)$' -count=1
npx playwright test --config jstest/playwright.config.ts jstest/specs/site-layout.spec.ts jstest/specs/theme-editor.spec.ts
go run ./internal/generatedcheck/cmd
```

Expected: PASS.

- [ ] **Step 6: Inspect the user's existing live Vite page**

At desktop and narrow widths, verify workspace overflow, one iframe, palette
selection, Custom import, light/dark mode, and no console/failed requests.
Do not start a second dev watcher.

- [ ] **Step 7: Commit**

```bash
git add site/pages jstest/specs/site-layout.spec.ts
git commit -m "feat: make theme editor a viewport workspace"
```

### Task 4: Site-layout integration gate and adversarial review

**Files:**
- Review all Task 1–3 files
- Modify only for verified defects

**Interfaces:**
- Produces a clean, independently reviewed site-layout subsystem

- [ ] **Step 1: Run static and authoritative checks**

```bash
gopls check -severity=hint site/pages/layout.x.go site/pages/toc.x.go site/pages/component.x.go site/pages/getting_started.x.go site/pages/theming.x.go site/pages/theme.x.go site/pages/pages_test.go
node --check web/site.js
git diff --check
make check
```

Expected: no new diagnostics and the full gate passes.

- [ ] **Step 2: Build production assets and restore the sentinel**

```bash
npm run build
git restore -- site/dist/.gitkeep
git status --short
```

Expected: build succeeds and no generated distribution output remains.

- [ ] **Step 3: Request independent adversarial review**

Probe 1280/1100/900 docs geometry, duplicate/missing heading IDs, hash and
active navigation, short-height workspace overflow, narrow order, palette and
import behavior, iframe Retry, home, components index, and isolated Sidebar
previews. Pass only with no Critical or Important findings.

- [ ] **Step 4: Confirm cleanup**

```bash
lsof -nP -iTCP:7799 -sTCP:LISTEN || true
git status --short --branch
```

Expected: no harness on 7799, port 5173 untouched, and a clean worktree.

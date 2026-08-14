# Chart Component Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Ship gsxui's chart component — a Recharts-shaped gsx composition API over one vendored client SVG renderer, adapted from templui v2, covering area/bar/line/pie/radar/radial-bar with grid/axes/tooltip/legend.

**Architecture:** Server-side gsx components assemble a chart model through render `context.Context` and emit it as JSON beside an empty container; a vendored JS renderer (adapted from templui's `chart.js`, itself a port of the needed Recharts/recharts-scale/d3-shape/react-smooth slices) draws Recharts-compatible SVG client-side. Colors flow through `--chart-1..5` theme tokens and per-chart `--color-<key>` variables, so theming and dark mode need no JS participation.

**Tech Stack:** gsx components (`registry/canonical/chart.gsx` → generated per style), vanilla ES module JS (`ui/chart.js` stub + `ui/chart.render.js` body), stylegen recipe pipeline, Playwright jstest, Go render tests.

**Reference:** templui v2 chart at pinned SHA `9ec720c03909` (github.com/templui/templui, MIT):
- `components/chart/chart.templ` (Go API + model builder + legend/tooltip markup)
- `components/chart/chart.js` (client renderer)
Fetch with `gh api repos/templui/templui/contents/components/chart/<file>?ref=9ec720c03909 --jq .content | base64 -d`. Do not commit the raw reference files; adapt into gsxui files with credit headers.

## Global Constraints

- **Credit templui everywhere it applies** (Jackie's explicit requirement): `NOTICE.md` gains a shadcn-templ/templui section; `registry/canonical/chart.gsx` and `ui/chart.render.js` carry a header comment naming templui, the pinned SHA, and MIT. A task does not merge without its credits.
- gsxui never ships `!important` (make audit gates it); the port strips upstream importants.
- No style token strings hardcoded in JS — presentation classes travel via the model JSON from the server (same discipline as the rest of `ui/`).
- No htmx strings in `ui/` — swap re-init rides the `init()` contract only.
- Slot markers are `data-gsxui-slot-chart[-*]`, never `data-tui-*`/`data-slot`.
- gsx comment rules: line-start `//` between tags is a comment; `{/* */}` among attributes or touching text.
- Unexported Go identifiers unless serialized or part of the consumer API.
- Gates after every task: `go build ./...`, `go test ./... -count=1`, `go run ./cmd/stylegen --check` + `--check-layers` + `--check-authoring`, `make audit`, `make verify-generated`, `make verify-generated-styles`, `npx playwright test --config jstest/playwright.config.ts` (never without the config), `gofmt -l`. Any `site/examples/**` edit requires `make highlight` and committing `site/hl/blocks.gen.go`.
- Commit per task; never `git add -A`.

---

### Task 1: NOTICE credit and reference capture

**Files:**
- Modify: `NOTICE.md`
- Create: none committed (reference downloads go to a temp dir)

**Interfaces:**
- Produces: the credit text later file headers cite; the pinned reference files at `/tmp`-equivalent path for Tasks 2-7.

- [ ] **Step 1: Download the two reference files at the pinned SHA**

```bash
mkdir -p /tmp/templui-ref
gh api "repos/templui/templui/contents/components/chart/chart.templ?ref=9ec720c03909" --jq .content | base64 -d > /tmp/templui-ref/chart.templ
gh api "repos/templui/templui/contents/components/chart/chart.js?ref=9ec720c03909" --jq .content | base64 -d > /tmp/templui-ref/chart.js
wc -l /tmp/templui-ref/chart.templ /tmp/templui-ref/chart.js
```
Expected: both files non-empty (~1.4k and ~2.5k lines).

- [ ] **Step 2: Add the NOTICE section**

Append to `NOTICE.md` after the existing shadcn section, matching its prose register:

```markdown
The chart component (`ui/chart.gsx`, `ui/chart.js`, `ui/chart.render.js`)
is adapted from **templui / shadcn-templ v2**
([templui/templui](https://github.com/templui/templui), MIT © templui
contributors, at 9ec720c03909): the Recharts-shaped Go composition API,
the server-side chart model, and the client renderer — itself a literal
port of the slices of Recharts, recharts-scale, d3-shape and react-smooth
those components need — originate there, translated from templ/React
idioms to gsx.
```

- [ ] **Step 3: Verify the CLI ships the NOTICE** — `internal/cli/add.go` already copies `NOTICE.md` on add (line ~167); run `go test ./internal/cli/ -count=1` to confirm nothing pins the old byte length.

- [ ] **Step 4: Commit**

```bash
git add NOTICE.md
git commit -m "docs(notice): credit templui/shadcn-templ v2 for the chart adaptation"
```

---

### Task 2: Shape, canonical container, per-style recipes

**Files:**
- Create: `registry/canonical/shapes/chart.go`, `registry/canonical/chart.gsx`, `registry/styles/<all 8>/chart.css`
- Modify: `registry/canonical/shapes/shapes.go` (register in `All()`), `internal/stylecontract/contracts_primitives.go` (Chart entry)
- Test: `ui/chart_test.go`

**Interfaces:**
- Produces: `ui.Chart(config ChartConfig, children gsx.Node, attrs gsx.Attrs)` container; package-level types in chart.gsx (flat package ui — names carry the Chart prefix to avoid generic collisions): `type ChartConfig []ChartSeries`, `type ChartSeries struct { Key, Label, Color string; Theme *ChartSeriesTheme; Icon gsx.Node }`, `type ChartSeriesTheme struct { Light, Dark string }`; unexported `func (c ChartConfig) styleBlock(id string) string`.
- Consumes: Task 1's reference `chart.templ` lines 1-160 (Config/styleBlock/Container).

- [ ] **Step 1: Write failing styleBlock tests** in `ui/chart_test.go`, porting templui's semantics exactly:

```go
func TestChartStyleBlockPrecedence(t *testing.T) {
	cfg := ui.ChartConfig{
		{Key: "desktop", Color: "var(--chart-1)"},
		{Key: "mobile", Color: "red", Theme: &ui.ChartSeriesTheme{Light: "blue", Dark: "green"}},
	}
	got := render(t, ui.Chart(cfg, gsx.Raw("x"), nil))
	// Theme wins over Color; light block unprefixed, dark under .dark.
	for _, want := range []string{
		"--color-desktop: var(--chart-1);",
		"--color-mobile: blue;",   // light: theme wins over red
		"--color-mobile: green;",  // dark
		"[data-chart=",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in %s", want, got)
		}
	}
}

func TestChartNoColorsNoStyle(t *testing.T) {
	got := render(t, ui.Chart(ui.ChartConfig{{Key: "a", Label: "A"}}, gsx.Raw("x"), nil))
	if strings.Contains(got, "<style>") {
		t.Errorf("uncolored config must emit no style block: %s", got)
	}
}
```

- [ ] **Step 2: Run, verify FAIL** (`go test ./ui/ -run TestChartStyleBlock -count=1` — undefined: ui.Chart).

- [ ] **Step 3: Author the canonical component.** `registry/canonical/chart.gsx` header comment: `// Adapted from templui/shadcn-templ v2 (github.com/templui/templui @ 9ec720c03909, MIT) — see NOTICE.md.` Port from reference lines 1-160: `Config`, `Series`, `SeriesTheme`, `styleBlock` (byte-faithful semantics: per-scheme blocks, `.dark ` prefix, theme-over-color, fallback-to-color, no-colors→empty), and the container component:

```
component Chart(config ChartConfig, children gsx.Node, attrs gsx.Attrs) {
	// id: deterministic per render via the existing uniqueid helper the
	// theme preview uses (TestPreviewDocumentIdsAreUnique pins collisions).
	<div
		data-chart={id}
		class={ chart.Root() {/* recipe accessor per the *_recipe.go binding pattern, e.g. registry/canonical/dropdown-menu_recipe.go */} }
		{ attrs... }
		data-gsxui-slot-chart
	>
		if block := config.styleBlock(id); block != "" {
			{ gsx.Raw("<style>" + block + "</style>") }
		}
		{ children }
	</div>
}
```

Drop templui's `[&_.recharts-*]` selector wall from the class: our renderer stamps gsxui classes, and those selectors' jobs (axis tick fill, grid stroke, cursor fill) move into the renderer-emitted attributes reading CSS variables (Task 5). Shape file `shapes/chart.go`: slots `""` (base), `tooltip`, `legend` — no dimensions.

- [ ] **Step 4: Author the 8 pack css files.** Each `registry/styles/<style>/chart.css`: header `/* Source: shadcn-ui/ui@41bbc12c... style-<style>.css Chart section + templui chart.templ container; adapted per NOTICE.md. */`, then `@layer components` with `gsxui-recipe-chart` (container base) and `gsxui-recipe-chart-tooltip` — the tooltip @apply copied from that style's upstream `.cn-chart-tooltip` rule (all 8 exist; lyra's is at style-lyra.css:257), important markers stripped.

- [ ] **Step 5: Regenerate and re-run** — `make generate`, then Step 1 tests PASS; `go run ./cmd/stylegen --check-authoring` green (chart is a migrated component from birth).

- [ ] **Step 6: Register the stylecontract entry** (Name "Chart", slots chart/chart-tooltip/chart-legend) and run `go test ./internal/stylecontract/ ./internal/registry/ -count=1`. The examples-coverage gate will fail until Task 8 — if `TestRegisteredExamplesCoverStyleContract` trips already, add the minimal `/x/chart` basic example in THIS task instead of waiting.

- [ ] **Step 7: Commit** — `git commit -m "feat(chart): container, config->color variables, per-style recipes"` (add the created/modified paths explicitly).

---

### Task 3: Cartesian model builder (bar, line, area)

**Files:**
- Modify: `registry/canonical/chart.gsx` (append model section)
- Test: `ui/chart_model_test.go`

**Interfaces:**
- Produces: `type Datum map[string]any`; roots `BarChart`/`LineChart`/`AreaChart(data []Datum, children gsx.Node, attrs gsx.Attrs)` (options via typed optional params following templui's props structs — port `Margin`, `CurveType`, `StackOffset`, pointer-default helpers `Bool`/`Float`); children `XAxis`, `YAxis`, `CartesianGrid`, `Tooltip`, `Legend`, `Bar`, `Line`, `Area`, `Defs`, `LinearGradient` registering into the ctx builder; serialized model `<script type="application/json" data-gsxui-chart-model>`.
- Consumes: reference `chart.templ` lines 160-1100 (chartState, buildModel, per-child registration, ModelScript) — port function-for-function, translating templ ctx (`ctx = context.WithValue(...)` inside templates) to a root `gsx.Func` that seeds the builder before rendering children and serializes after (gsx renders children lazily through `Render(ctx, w)`, so the templui order — children first, then `chartOutput` — maps to: root writes open tag, renders `children` with builder-ctx, then emits the model script).
- Legend/tooltip content: server renders legend HTML from config exactly as templui's `legendContent`; the tooltip's compiled recipe class string is embedded in the model (`tooltipClass` field) for the renderer.

- [ ] **Step 1: Write the failing byte-pin test** — fixed dataset, full expected JSON:

```go
func TestBarChartModelPinned(t *testing.T) {
	data := []ui.ChartDatum{{"month": "Jan", "desktop": 186.0}, {"month": "Feb", "desktop": 305.0}}
	got := render(t, ui.Chart(cfg2Series(t),
		ui.BarChart(data,
			gsx.Nodes(
				ui.ChartXAxis("month", nil),
				ui.ChartBar("desktop", nil),
			), nil), nil))
	model := extractModelJSON(t, got) // helper: contents of data-gsxui-chart-model script
	// Pin the FULL canonical JSON (indent-free, sorted keys — the builder
	// must marshal deterministically; sort series by registration order,
	// map keys via a struct model, never map[string]any at the top level).
	want := `{"kind":"bar","data":[{"desktop":186,"month":"Jan"},{"desktop":305,"month":"Feb"}],"series":[{"key":"desktop","color":"var(--color-desktop)"}],"xAxis":{"key":"month","tickLine":false,...}}`
	if model != want {
		t.Errorf("model drift\n got: %s\nwant: %s", model, want)
	}
}
```

The `...` above is written out in full when the builder lands — the test author fills the exact defaults from the reference's model fields (xAxisHeight 0, minTickGap 5, tickCount 5 etc., documented in the reference chart.js header) and the pin is byte-exact from then on.

- [ ] **Step 2: FAIL run** (`undefined: ui.BarChart`).
- [ ] **Step 3: Port the builder** — chartState struct, ctx keys (unexported), each child component appends to the state; `BarChart` renders `<div style="position:relative;width:100%;height:100%">` + children + model script, per the reference.
- [ ] **Step 4: PASS + regen** (`make generate`, `go test ./ui/ -count=1`).
- [ ] **Step 5: Commit** `feat(chart): cartesian model builder (bar, line, area)`.

---

### Task 4: Polar model builder (pie, radar, radial-bar)

Same shape as Task 3, porting the remaining roots and children (`PieChart`, `RadarChart`, `RadialBarChart`, `Pie`, `Radar`, `RadialBar`, `PolarGrid`, `PolarAngleAxis`, `PolarRadiusAxis`) from reference lines 729-1100, with one byte-pinned model test per kind (`TestPieChartModelPinned`, `TestRadarChartModelPinned`, `TestRadialBarChartModelPinned`) written first and watched failing. Commit `feat(chart): polar model builder (pie, radar, radial-bar)`.

---

### Task 5: Renderer body + init stub + CLI companion vendoring

**Files:**
- Create: `ui/chart.js` (stub), `ui/chart.render.js` (body)
- Modify: `internal/cli/add.go` (companion JS), `internal/registry/registry.go` only if HasJS needs no change (it doesn't — stub is `ui/chart.js`)
- Test: `internal/cli/add_test.go` (or the existing add tests file), `jstest/specs/chart.spec.ts`

**Interfaces:**
- Consumes: model JSON per Tasks 3-4; reference `chart.js` in full.
- Produces: `init("[data-gsxui-slot-chart]", initRoot)` stub that `await import("./chart.render.js")` once and calls `renderChart(el)`; body exports `renderChart(container)` reading the sibling model script and drawing SVG.

- [ ] **Step 1: CLI test first** — extend the add tests: adding chart vendors `chart.js` AND `chart.render.js` into `cfg.JS`; adding button does not. Watch fail.
- [ ] **Step 2: Extend `addArtifacts`** — after the existing HasJS copy, glob companion files `ui/<name>.*.js` (only chart has one today) and copy each with the same artifact shape. PASS.
- [ ] **Step 3: Write the stub** (`ui/chart.js`, full file):

```js
// Chart bootstrap: the renderer body (chart.render.js, ~90KB — adapted
// from templui, see NOTICE.md) loads on the first chart in the DOM, so
// pages without charts never parse it.
import { init } from "./gsxui.js";

let bodyPromise;
function initRoot(el) {
	bodyPromise ??= import("./chart.render.js");
	bodyPromise.then(({ renderChart }) => renderChart(el));
}
init("[data-gsxui-slot-chart]", initRoot);
```

- [ ] **Step 4: Port the renderer body.** `ui/chart.render.js` header: `/* Adapted from templui/shadcn-templ v2 chart.js (github.com/templui/templui @ 9ec720c03909, MIT) — see NOTICE.md. ... */` keeping their provenance list (Recharts ResponsiveContainer, recharts-scale getNiceTickValues, CartesianAxis preserveEnd, d3-shape curves, d3-array ticks, Sector, react-smooth). Adaptations, each grepped to completion: `data-tui-chart`→`data-gsxui-slot-chart`; model script lookup by `data-gsxui-chart-model`; tooltip container class read from `model.tooltipClass` (server-supplied) instead of their `TOOLTIP_CLASS` constant + recipe classes; export `renderChart` instead of their auto-boot; delete their loader/observer (the stub owns lifecycle); keep geometry/animation code unmodified otherwise — unmodified code is re-verifiable against the reference by diff.
- [ ] **Step 5: jstest chart spec (write BEFORE wiring the example page if TDD ordering allows a 404 fixture check first):**

```ts
test("BarChart renders themed SVG bars", async ({ page }) => {
	const response = await page.goto("/x/chart/basic");
	expect(response?.status()).toBe(200);
	const chart = page.locator("[data-gsxui-slot-chart]").first();
	const bars = chart.locator("svg rect[data-gsxui-chart-bar]");
	await expect(bars).toHaveCount(6); // basic example: 6 months x 1 series
	const fill = await bars.first().evaluate((n) => getComputedStyle(n).fill);
	const chart1 = await page.evaluate(() =>
		getComputedStyle(document.documentElement).getPropertyValue("--chart-1"),
	);
	expect(fill).not.toBe("");
	// fill is var(--color-desktop) -> config maps it to var(--chart-1):
	// compare resolved rgb of both.
});
```

(The minimal `/x/chart/basic` example lands here or in Task 2's Step 6 — whichever executed first; this task owns making the spec pass.)

- [ ] **Step 6: Gates + commit** `feat(chart): vendored renderer, lazy stub, CLI companion vendoring`.

---

### Task 6: Tooltip, legend, animations, htmx/dark behavior specs

**Files:**
- Modify: `ui/chart.render.js` (only if Task 5 deferred tooltip wiring), `jstest/specs/chart.spec.ts`

Specs written first, watched fail where behavior is missing, then fixed:

- [ ] Tooltip appears on hover with the per-style recipe class (`[data-gsxui-slot-chart-tooltip]` visible, class string equals the compiled `gsxui-recipe-chart-tooltip` utilities for the default style).
- [ ] Legend renders server-side (present in initial HTML before JS boots — assert via `page.route` blocking chart.render.js and still seeing legend text).
- [ ] Dark-mode flip: toggle `.dark` on `<html>`, assert a bar's computed fill changes with NO re-render (same node identity via `evaluateHandle`).
- [ ] htmx-style swap: replace the chart container's parent subtree via `innerHTML`, assert the observer re-inits and SVG reappears.
- [ ] Commit `test(chart): tooltip, legend, dark-mode and swap behavior`.

---

### Task 7: Lazy-load network gate

**Files:**
- Test: `jstest/specs/chart.spec.ts`

- [ ] **Step 1:** Spec: navigate to a no-chart page (`/x/button/variants`), collect requests, assert none matches `chart.render`; navigate to `/x/chart/basic`, assert exactly one `chart.render` fetch, then a second chart page visit in the same context still renders (module cache path).
- [ ] **Step 2:** Watch it fail only if the stub is wrong; fix; commit `test(chart): renderer body loads only on chart pages`.

---

### Task 8: Examples and docs

**Files:**
- Create: `site/examples/chart/{basic,area,line,pie,radar,radial}.gsx`, docs page under `site/` following the existing component-page pattern (copy the structure of the command docs page)
- Modify: `site/hl/blocks.gen.go` (via `make highlight`), navigation/registry of docs pages wherever the site enumerates components

- [ ] Port shadcn's six chart demos 1:1 (bar-interactive simplified to static props where it needs client state gsxui doesn't model). Each example compiles, `TestRegisteredExamplesCoverStyleContract` passes, `make highlight` regenerated and committed.
- [ ] Full Playwright suite green.
- [ ] Commit `feat(chart): docs page and six examples`.

---

### Task 9: Gallery chart card + theme editor picker restoration

**Files:**
- Modify: `site/stylepreview/gallery.gsx.src` (replace `galleryChartCard` bars with real `BarChart` + `AreaChart`), `site/pages/theme.gsx` (restore picker — revert the UI hunk of cb9a2379), `web/preview.js` (boot chart renderer for the preview document), `jstest/specs/theme-editor.spec.ts`
- Baselines: refresh only legitimately-changed pngs (never `--update-snapshots=all` blanket)

- [ ] Editor spec first: re-pin picker presence + chart-1 swatch drives a rendered bar's computed fill in the preview iframe; watch fail.
- [ ] Restore picker (the axis underneath is live — schema/state/share codes untouched since cb9a2379).
- [ ] Gallery card: real charts; if booting the renderer in `web/preview.js` costs more than importing the stub + body there, fall back to the static frame and record the choice in the card comment (spec allowed either; the live render is preferred since the picker demo depends on it).
- [ ] Full suite; refresh failing chart-card baselines per style; commit `feat(theme): real gallery charts; restore chart-color picker`.

---

### Task 10: Ledger, roadmap, final gates

**Files:**
- Modify: `docs/component-roadmap.md` (chart row: deferred → shipped, pointing at the spec and NOTICE credit), `docs/jsx-parity.md` (new `## chart` section: WIN/ADAPT/GAP entries — at minimum: ADAPT engine substitution recharts→templui-derived renderer with credit; GAP any interactive demo behavior not ported; MECHANISM ctx-builder registration), `CHANGELOG.md`

- [ ] Write the entries (docs style: state each fact once).
- [ ] Run the complete gate list from Global Constraints one final time.
- [ ] Commit `docs(chart): ledger, roadmap, changelog`.

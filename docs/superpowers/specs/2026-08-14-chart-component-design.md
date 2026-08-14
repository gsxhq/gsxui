# Chart component — design

2026-08-14. Approved in brainstorming (engine, scope, provenance, and all
three design sections confirmed by Jackie).

## Decisions

- **Engine**: templui v2's approach — NOT Chart.js. A Recharts-shaped
  server composition API plus one vendored client renderer that ports the
  slices of Recharts, recharts-scale, d3-shape, and react-smooth the
  components need, drawing Recharts-compatible SVG. This resolves
  `component-roadmap.md`'s deferral reason (Recharts is un-vendorable npm;
  a single MIT JS file is not).
- **Scope**: full templui parity in one campaign — area, bar, line, pie,
  radar, radial-bar, with grid/axes/tooltip/legend/gradients.
- **Provenance**: adapt templui v2's `chart.templ`/`chart.js`
  (github.com/templui/templui, MIT) to gsxui idioms. `NOTICE.md` gains a
  shadcn-templ/templui credit.
- **Cost gate** (hard requirement): a consumer who never installs chart
  pays zero bytes; a consumer who installs it pays the renderer body only
  on pages that render a chart.

## API surface

`ui/chart.gsx` + support package, mirroring shadcn's `chart.tsx`
composition as templui does:

- `Chart` container: `chart.Config` (`[]chart.Series{Key, Label, Color,
  Theme{Light, Dark}, Icon}`), stamps `data-gsxui-slot-chart` + generated
  chart id, emits the scoped `<style>` block mapping config colors to
  `--color-<key>` (shadcn ChartStyle semantics: Theme wins over Color, a
  scheme without a themed value falls back to Color, no colored entry →
  no style block).
- Roots: `AreaChart` `BarChart` `LineChart` `PieChart` `RadarChart`
  `RadialBarChart`, taking `[]chart.Datum` (`map[string]any`) plus
  Recharts-default options (pointer fields for default-true props, per
  templui's `Bool`/`Float` helpers).
- Children render nothing; they register into the root's model builder
  through render `context.Context` (`gsx.Node.Render(ctx, w)` carries it:
  the root seeds a builder via `gsx.Func`, children mutate, the root
  serializes after children): `XAxis` `YAxis` `CartesianGrid` `Tooltip`
  `Legend` `Area` `Bar` `Line` `Pie` `Radar` `RadialBar` `PolarGrid`
  `PolarAngleAxis` `PolarRadiusAxis` `Defs`/`LinearGradient`.

House-rule deviations from templui: typed gsx component params over
variadic props structs where natural; `data-gsxui-slot-chart-*` markers,
not `data-tui-*`; gsxui's `init()` contract, not their loader; unexported
Go types unless serialized.

## Pipeline

Server serializes the model to one `<script type="application/json"
data-gsxui-chart-model>` beside an empty positioned div. `ui/chart.js`
(the adapted renderer, one file) draws SVG at the container's measured
size (ResizeObserver), with ported nice-tick math, tick culling, curves,
stack offsets, sector geometry, and react-smooth-style entrance/update
animations. Legend is server-rendered HTML. SVG fills reference
`var(--color-<key>)`/`var(--chart-N)`, so dark-mode flips repaint without
JS.

## Vendoring and lazy loading

- `registry.HasJS("chart") = true`; copy-in gates all chart bytes on
  `gsxui add chart`.
- `ui/chart.js` exposes the standard `init()` as a ~10-line stub that
  dynamic-`import()`s the renderer body (split file) the first time a
  scanned subtree contains `[data-gsxui-slot-chart]`. Pages without
  charts parse only the stub; htmx swaps re-init through the existing
  contract unchanged.

## Theming

1. Default palette: `--chart-1..5` theme tokens (already emitted, already
   in the editor's schema/state/share codes).
2. Per-chart: the Config → `--color-<key>` scoped style block.
3. Per-style: upstream ships exactly one per-style chart rule,
   `.cn-chart-tooltip` — ported to all 8 packs as
   `gsxui-recipe-chart-tooltip` via the normal recipe pipeline. The
   client-rendered tooltip stamps that compiled class string carried in
   the model JSON — no style tokens hardcoded in JS.

## Editor and gallery

- Restore the chart-color picker (reverts the UI half of cb9a2379; the
  axis — schema, state, share codes, JSON import, emitted variables —
  never left).
- Replace `galleryChartCard`'s hand-rolled div bars with a real BarChart
  plus one Area/Line chart driven by `--chart-1..5`, so the picker
  visibly changes the preview. The preview document loads no component
  JS; either `web/preview.js` boots just the chart renderer, or the
  gallery card keeps a static frame and `/x/chart` is the live demo —
  implementer decides by what preview.js carries cheaply, and records
  the choice in the plan.

## Docs, examples, registry

`/x/chart` docs page with examples ported from shadcn's demos (bar, area
interactive, line, pie donut, radar, radial), feeding the contract
coverage gates. Registry entry with deps; flip `component-roadmap.md`'s
chart row from deferred to done.

## Testing

- Go: pin container markup, style-block precedence, and the model JSON
  byte-for-byte for fixed datasets.
- jstest: SVG geometry basics (bar count, computed fill == resolved
  `--chart-1`), tooltip on hover, re-init after htmx swap, dark-mode
  color flip, and a network-level spec proving the renderer body is not
  fetched until a chart page loads.
- Editor specs re-pin the restored picker; visual baselines pick up the
  new gallery card.

## Out of scope

Chart kinds beyond the six; scatter/composed charts; brush/zoom
interactions; server-side SVG rendering.

## Addendum (2026-08-14, Task 10 — tooltip-template ruling)

Supersedes the Theming section's third paragraph ("The client-rendered
tooltip stamps that compiled class string carried in the model JSON"): that
mechanism (`ChartTooltipModel.TooltipClass`) is structurally unbuildable —
`checkShapeCoverage` rejects a recipe-accessor call made from plain Go code,
and `buildChartModel` is not a component body. The tooltip's chrome instead
renders as a real, hidden, server-rendered `<template>` element
(`ChartTooltipTemplate`, `registry/canonical/chart.gsx`) carrying `class={
chart.Tooltip() }` in genuine markup; `ui/chart.render.js` reads that
template's own class attributes off the DOM instead of receiving them as a
JSON string. Implemented in Task 6; see `docs/jsx-parity.md`'s `## chart`
MECHANISM entry.

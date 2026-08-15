# Changelog

Notable changes to gsxui's component set, newest first.

## 2026-08-15

### Changed

- **chart** — every root and child (`BarChart`/`LineChart`/`AreaChart`/`PieChart`/`RadarChart`/`RadialBarChart`, `ChartCartesianGrid`/`ChartXAxis`/`ChartYAxis`/`ChartTooltip`/`ChartLegend`/`ChartBar`/`ChartLine`/`ChartArea`/`ChartPie`/`ChartRadar`/`ChartRadialBar`/`ChartPolarGrid`/`ChartPolarAngleAxis`/`ChartPolarRadiusAxis`/`ChartDefs`/`ChartLinearGradient`) is now a tag-based `component` with flat typed params (e.g. `<ui.AreaChart data={data} marginLeft={12}><ui.ChartCartesianGrid horizontal/>...</ui.AreaChart>`) instead of an options-struct/function-call surface. `ChartBool`/`ChartFloat` and every public `*Options` struct are removed — pre-release, no compatibility shim. Boolean props that used to default to Recharts' own `true` (`horizontal`/`vertical`, `cursor`, `radialLines`, `labelLine`) now default `false`, enabling is explicit; see `docs/jsx-parity.md`'s `## chart` ADAPT entry.

## 2026-08-14

### Added

- **chart** — six kinds (bar, line, area, pie, radar, radial-bar) over a server-rendered model and a single lazy-loaded client renderer, so pages without a chart pay zero bytes for one. Ported from templui's `chart.templ`/`chart.js` (credited in `NOTICE.md`).

## 2026-07-29

### Added

- **contract** — `registry/generated/recipes.json` publishes the recipe model: every component's slots, their dimensions and values, and each style's utilities for them (`components.<c>.slots.<s>.dimensions.<d>`, `styles.<style>.<c>.slots.<s>`). It carries a `version` field so a consumer can check it understands the schema.

### Changed

- **button** — `gsxui init` ships a `ui/button.gsx` whose Button renders concrete Tailwind utilities compiled from its style recipe, so consumer-project buttons are styled out of the box with no extra stylesheet.
- **sidebar** — `SidebarTrigger` renders at 28px (`size-7`), matching both its own authored CSS and upstream shadcn.
- **button** — destructive buttons lighten on hover in dark mode (`bg-destructive/90`), matching the documented style contract. This diverges from upstream shadcn, which keeps `/60` through hover; see `docs/jsx-parity.md`.

## 2026-07-25

### Added

- **menubar** — application-style menu bar with nested submenus and full keyboard navigation.
- **navigation-menu** — hover-driven top nav with panel content.
- **calendar** — month grid with single, range and multiple selection, keyboard grid, and no react-day-picker.
- **JS test layer** — Playwright suite in `jstest/` over a Go example harness; four invariants sweep every example, gated in CI.

### Changed

- **dropdown** and **context-menu** gained checkbox items, radio groups, and submenus.

## 2026-07-24

### Added

- **resizable** — drag-resizable split panes with keyboard support.
- **combobox** — filterable input + listbox with filter-as-you-type and form binding.
- **sidebar** — collapsible app sidebar with desktop and mobile layouts.

## Earlier

- Initial 47 components (Tiers 1-3, plus `command`) shipped pre-changelog.

# Changelog

Notable changes to gsxui's component set, newest first.

## 2026-09-14

### Added

- **i18n** — `i18n.T`, a string type every component uses for the strings it writes itself (`Close`, `Notifications`, `Previous slide`, `Toggle Sidebar`, `More pages`, the calendar nav labels, …). It renders as English until the consuming module registers a renderer for it in `gsx.toml`; see `/docs/i18n`. Vendored automatically as `ui/i18n/` with any component that uses it (#31).
- **calendar** — `locale CalendarLocale`: month and weekday names, caption and day-label patterns, and a native digit set. The zero value is English. `calendar.js` composes navigation text from the same values the server wrote to the root (#30).

### Changed

- **site** — htmx pin `4.0.0-beta6` -> `4.0.0` (GA, released 2026-08-28) (#27).

### Fixed

- **jstest harness** — serves the favicon set (`/favicon.svg`, `/favicon-32.png`, `/apple-touch-icon.png`) the site layout links, via the new shared `site/icons` package; harness pages no longer log three 404s each (#27).
- **toaster** — docs note that under htmx 4 a response containing only `hx-swap-oob="beforeend:#gsxui-toaster"` content leaves the request's main target untouched (#27).
- **accordion** — the open item's muted tint in luma, maia, mira and rhea never rendered: the sheets carried `data-[state=open]:bg-muted/50`, and nothing stamps `data-state` on a native `<details>`. Now the native `open:` variant. The porter emits `open:`/`not-open:` for the `<details>` slot itself and the authoring gate rejects `data-[state=` in accordion and collapsible sheets (#28).
- **alert-dialog, dialog** — `AlertDialogFooter` in every non-nova style, and lyra's `DialogFooter`, rendered nova's full-bleed muted end-bar instead of upstream's plain right-aligned button row (#28).
- **alert-dialog** — the header was centred at every width in all styles; upstream's default size is left-aligned from `sm` up. `AlertDialogContent` now stamps `data-size="default"`, which the ported `sm:` alignment selectors were keyed on all along (#28).

## 2026-08-24

### Changed

- **field, checkbox, radio, switch** — upstream parity (pin `41bbc12c` -> `ac60ef5c`): when a checkbox, radio, or switch inside a `FieldLabel` card takes keyboard focus, the focus ring now renders on the card instead of the control, and the card gains a hover tint; `FieldLabel` declares the `group/field-label` marker these selectors scope to. See `docs/jsx-parity.md`'s `## field-label focus ring`.

### Fixed

- **field** — `FieldLabel`'s checked-state highlight (the tint behind a checked checkbox or radio card) never rendered in any style: it was authored as `has-data-checked:*` but gsxui's controls are native inputs and nothing stamps `data-checked`. Now `has-[input:checked]:*`, scoped to `input` so a `FieldLabel` wrapping a select (whose always-selected `<option>` also matches `:checked`) doesn't tint permanently.

## 2026-08-19

### Changed

- **sidebar** — `SidebarMenuButton` gains `href string` and `disabled bool` params (positional, after `tooltip`; direct Go callers must pass them). A non-empty href renders an `<a>` instead of a `<button>` — gsxui's `asChild` pendant, as on Button — with `aria-current="page"` when active; `disabled` always renders a real disabled `<button>`, even with href. The tooltip trigger marker moves onto whichever element is rendered.

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

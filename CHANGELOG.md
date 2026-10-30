# Changelog

Notable changes to gsxui's component set, newest first.

## 2026-10-30

### Fixed

- **pagination**, **calendar** — `PaginationLink` and Calendar's previous/next and day buttons render through `Button`, so an app gets Button's sizing, centering, radius and hover; before, they carried only the `data-gsxui-slot-button` marker and rendered as plain text and blocks (a pagination link measured 7×19 px). `PaginationPrevious`/`PaginationNext` take upstream's per-style edge padding (`pl-2`/`pr-2` in vega, maia, mira, luma and sera; `pl-1.5`/`pr-1.5` in nova, lyra and rhea). A caller overrides it with `has-[>svg]:pl-*`, since a plain `pl-*` loses to Button's `has-[>svg]:px-*` (#51).

- **cli** — `gsxui add` and `gsxui apply` format each file whose imports they rewrite, so `gsx fmt` and `gofmt` leave vendored components unchanged when the app's module path sorts before `github.com/gsxhq/gsx` (`example.com/…`, `codeberg.org/…`). Before, formatting such a file made the next `gsxui add` stop with "differs from the gsxui version" (#50).

## 2026-10-21

### Fixed

- **a11y** — Lighthouse/axe failures inside components (#47):
  - **item** — an `Item` inside an `ItemGroup` is `role="listitem"`; standalone Items are unchanged.
  - **menubar** — `MenubarTrigger` is `role="menuitem"`, as in Radix.
  - **command** — `CommandInput` gets `aria-controls` pointing at `CommandList` (set by `command.js`); `CommandList` is labelled `Suggestions` (cmdk's default, an `i18n.T` message); `CommandSeparator` is `aria-hidden` instead of `role="separator"`, which a listbox may not own.
  - **combobox** — the icon-only trigger and clear buttons are labelled `Open popup` and `Clear selection` (`i18n.T` messages); `combobox` now vendors `ui/i18n`.
  - **resizable** — `ResizablePanelGroup` carries `data-orientation` instead of `aria-orientation`, which a plain `<div>` may not have. CSS that targeted the group's `aria-orientation` must switch; the handles keep theirs.

### Added

- **calendar** — day buttons' `aria-label` now says `Today, …` for the client's today and `…, selected` for selected days (every day of a committed range), as upstream's react-day-picker does. Both patterns are `i18n.T` messages (`Today, {date}`, `{date}, selected`), so a translator can reorder them around the date. `calendar.js` applies them; the server-rendered label stays plain until it runs.

- **calendar** — `CalendarLocale.MonthsFormat`, the month names day labels use inside a date, where a language's stand-alone caption form differs (Russian caption `июнь 2026`, day label `понедельник, 29 июня 2026 г.`). Unset, day labels keep using `Months`, so existing locales render unchanged. Example: `/components/calendar` "Russian" (#43).

### Changed

- **i18n** — the defaults a caller can override through `attrs` are now `i18n.T` messages too: `Spinner`'s `Loading`, the `breadcrumb` and `pagination` nav labels, `Go to previous page`/`Go to next page`, `SidebarRail`'s `Toggle Sidebar` (`aria-label` and `title`), and `Carousel`/`CarouselItem`'s `carousel`/`slide` `aria-roledescription`. English output is unchanged; a registered translator now sees them, and a caller's attribute still replaces them. `spinner` now vendors `ui/i18n` (#44).

### Fixed

- **calendar** — a range passed with `from` after `to` now marks the days between as selected (`data-range-middle`, `aria-selected`) on the server render, as react-day-picker and `calendar.js` already did; before, they stayed unselected until the first client repaint.

## 2026-10-02

### Changed

- **toolchain** — gsx `v0.3.0`. Regenerated output: the site document's `<link>` tags now carry the CSP nonce under `gsx.WithNonce`, like `<script>` and `<style>`; output without a context nonce is unchanged. No gsxui component puts a Go value on an event-handler attribute, so v0.3.0's JS encoding of `on*` attributes changes nothing here; a consumer writing `onclick={expr}` should read gsx's release notes.

## 2026-10-01

### Changed

- **toolchain** — gsx `v0.2.0` (from the 2026-09-30 main pin). Two rendering changes reach consumers: a `style` value keeps its authored trailing `;` (`style="aspect-ratio: 16 / 9;"`), and a `switch`-form attribute follows the spread position, so `AspectRatio` now renders `class` before `style`. Style values passed in a `gsx.Attrs` bag are CSS-filtered by the new runtime; pass trusted CSS as `gsx.RawCSS`, or a multi-declaration string renders `ZgotmplZ`. Regenerated output: aspect-ratio and the style gallery.

## 2026-09-28

### Fixed

- **button-group** — `ButtonGroup` now carries shadcn's own `group/button-group` marker class on its root. Without it the `sera` style's `group-has-[>[data-variant=outline]]/button-group:border-border` rule on `ButtonGroupText` could never match, so a text part beside an `outline` `Button` kept `border-transparent` instead of taking `border-border`. Affects the `sera` style only; the other seven never reference the marker.

## 2026-09-16

### Added

- **rtl** — `"rtl": true` in `gsxui.json` (`gsxui init --rtl`) makes `gsxui add` and `gsxui apply` vendor components with logical direction classes, using shadcn's `transform-rtl` table; `gsxui migrate rtl` rewrites components vendored earlier in place, keeping your edits. One build then serves `dir="ltr"` and `dir="rtl"` (#32).
- **foundation.css** — `origin-top-start`, `origin-top-end`, `origin-bottom-start`, `origin-bottom-end`, `origin-start`, `origin-end` utilities, which the transform emits.

### Changed

- **site** — the RTL page describes the opt-in model; its demos render from `site/uirtl`, the transform's output. The page previously stated components used logical classes by default, which was not true of the registry.

## 2026-09-14

### Added

- **i18n** — `i18n.T`, a string type every component uses for the strings it writes itself (`Close`, `Notifications`, `Previous slide`, `Toggle Sidebar`, `More pages`, the calendar nav labels, …). It renders as English until the consuming module registers a renderer for it in `gsx.toml`; see `/docs/i18n`. Vendored automatically as `ui/i18n/` with any component that uses it (#31).
- **calendar** — `locale CalendarLocale`: month and weekday names, caption and day-label patterns, and a native digit set. The zero value is English. `calendar.js` composes navigation text from the same values the server wrote to the root (#30). `locale` is positional, before `attrs`; direct Go callers of `ui.Calendar(...)` must pass `ui.CalendarLocale{}`. Markup callers using named attributes are unaffected.

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

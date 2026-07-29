# Changelog

Notable changes to gsxui's component set, newest first.

## 2026-07-29

### Fixed

- **button** — `gsxui init` now ships a `ui/button.gsx` whose Button renders concrete Tailwind utilities compiled from its style recipe, instead of `.gsxui-recipe-button*` classes that no shipped stylesheet defined (those rules lived only in the site-only `web/site-button.css`, never copied by `gsxui init`). Consumer-project buttons were unstyled at merge base; they are now styled out of the box.

### Changed

- **button** — destructive buttons now lighten on hover in dark mode (`bg-destructive/90`), matching the documented style contract. This diverges from upstream shadcn, which keeps `/60` through hover; see `docs/jsx-parity.md`.

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

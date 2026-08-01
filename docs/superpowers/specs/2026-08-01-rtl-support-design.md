# RTL Support — Design

**Date:** 2026-08-01
**Status:** Approved

## Goal

Full RTL correctness for gsxui: every component renders and behaves correctly
under `dir="rtl"` (Arabic and similar languages), following shadcn/ui's
official RTL conventions (first-class since their January 2026 release), plus
RTL examples on the docs site wherever shadcn's docs have them.

## Background

- shadcn's approach: Tailwind **logical properties** everywhere
  (`start-`/`end-`, `ms-`/`me-`, `ps-`/`pe-`, `text-start`/`text-end`, logical
  border/rounding, `slide-in-from-start` animations), `rtl:rotate-180` on
  directional icons, and a Radix `DirectionProvider` so portalled floating
  content aligns correctly. Calendar, Pagination, and Sidebar have documented
  manual migration steps; the RTL docs page has an Arabic login-card demo.
- gsxui state: most vendored classes are already logical, but ~19 `ui/*.gsx`
  files retain physical classes (`pr-8`, `right-2`, `border-l`, `ml-auto`,
  `rounded-l/r-*`, `text-left`, `left-1/2`, `translate-x`, …). No JS is
  direction-aware today (`position()` in `ui/gsxui.js` resolves
  `align: start/end` assuming LTR; slider/carousel/resizable/arrow-key
  handling are all LTR-only). gsxui has no Radix — our own `position()` plays
  the DirectionProvider role.

## Design

### 1. CSS pass — logical classes in `ui/*.gsx`

Convert remaining physical-direction classes to shadcn's exact logical
conventions:

| Physical | Logical |
|---|---|
| `left-*` / `right-*` | `start-*` / `end-*` |
| `ml-*`/`mr-*`, `pl-*`/`pr-*` | `ms-*`/`me-*`, `ps-*`/`pe-*` |
| `text-left` / `text-right` | `text-start` / `text-end` |
| `border-l*` / `border-r*` | `border-s*` / `border-e*` |
| `rounded-l-*` / `rounded-r-*` | `rounded-s-*` / `rounded-e-*` |
| `slide-in-from-left/right` | `slide-in-from-start/end` |

Rules:
- Genuinely physical usage stays physical (e.g. `left-1/2` centering,
  symmetric `translate-x` animations).
- Directional icons (chevrons/arrows in accordion, breadcrumb, pagination,
  sidebar trigger, carousel prev/next, dropdown/menubar/context-menu submenu
  indicators, …) get `rtl:rotate-180`. `rtl:` variants are otherwise used
  sparingly — logical properties do the work.
- Diff against shadcn's post-RTL-migration sources where available so our
  markup stays diffable against upstream.

### 2. JS pass — direction awareness

One shared helper in `ui/gsxui.js`:

```js
isRTL(el) // getComputedStyle(el).direction === "rtl"
```

Computed style respects both `dir` attributes and CSS `direction`.

- **`position()`**: when the anchor is RTL, mirror cross-axis alignment
  (`start` ↔ `end`) and horizontal preferred sides (`left` ↔ `right`, used by
  submenus). All poppers inherit the fix: dropdown-menu, context-menu,
  menubar, select, combobox, popover, tooltip, hover-card, navigation-menu.
- **Arrow keys** flip where Radix flips them:
  - submenu open/close keys in dropdown-menu, context-menu, menubar
    (ArrowRight opens in LTR, closes in RTL; ArrowLeft mirrored);
  - horizontal roving focus in tabs, toggle-group, menubar, carousel;
  - calendar day navigation.
- **`slider.js`**: mirror pointer-to-value math and ArrowLeft/ArrowRight;
  track fill uses logical inset in CSS.
- **`carousel.js`**: correct horizontal scroll/snap math under RTL; prev/next
  buttons flip meaning.
- **`resizable.js`**: mirror horizontal drag delta and Home/End/Arrow keys.
- **`sidebar.js` + `sidebar.gsx`**: shadcn's documented manual migration —
  drive side via `data-side` attribute + CSS logical rules instead of JS
  ternaries; `rtl:rotate-180` on the trigger icon; `Sidebar` accepts
  `dir`/`side` combinations.
- **`input-otp`**: OTP digit groups render LTR even in RTL context (standard
  numeric-code convention) — force `dir="ltr"` on the slot group; confirm
  against shadcn's behavior during implementation.

### 3. Docs site

- New **RTL docs page** mirroring shadcn's: intro, how-it-works (logical
  properties + `dir`), live Arabic **login-card demo**, font recommendation
  note (Noto family).
- Dedicated `dir="rtl"` example blocks with Arabic sample content on the
  component pages where shadcn has per-component RTL sections: **Calendar,
  Pagination, Sidebar**. Examples follow the existing
  `site/examples/<component>/` structure.

### 4. Verification

Playwright specs asserting real geometry under `dir="rtl"`:
- dropdown content aligns to the trigger's logical start (visually right);
- submenu opens to the left of its parent;
- sheet's `side="left|right"` stays PHYSICAL under RTL (shadcn's data-side
  contract): a `side="left"` sheet opens from the viewport's left in both
  directions, while its interior spacing/text stays logical;
- slider fills right-to-left; pointer/keyboard values mirror;
- carousel advances in the correct direction.

Plus existing gates: full build, `gsx fmt` idempotency, Go tests, existing
Playwright suite.

## Out of scope

- Site-wide LTR/RTL toggle on the docs site.
- Go-level direction API (`DirectionProvider` component) — the DOM `dir`
  attribute is the source of truth; server-rendered markup needs no context.
- Localized docs content beyond the sample Arabic strings.
- Vertical writing modes.

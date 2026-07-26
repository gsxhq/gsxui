# Post-v1 backlog

Deferred out of v1 scope, tracked here rather than in the parity ledger's
per-component GAP notes (see `jsx-parity.md` for the detailed rationale).

- **Tooltip delay-groups** — `TooltipProvider`'s shared `delayDuration`/
  skip-delay-group coordination across multiple tooltips on a page. v1
  hard-codes a fixed per-trigger open delay, no cross-tooltip grouping.
- **CSS anchor positioning migration** — dropdown and tooltip currently
  position via a hand-rolled `getBoundingClientRect()` + `position: fixed`
  calculation in JS. Once CSS anchor positioning (`anchor()`/
  `position-anchor`) reaches Baseline across browsers, both can drop that
  JS for native, collision-aware placement.
- **Checkbox checkmark theming (currentColor mask)** — the check glyph is a
  data-URI with hard-coded `stroke="white"`; data-URIs are static text and
  can't reference CSS variables, so the mark doesn't follow
  `--primary-foreground` and is wrong/low-contrast for themes where that
  color isn't near-white. Swap to a `currentColor` CSS-mask
  (`mask-image`/`-webkit-mask-image` painted via `background-color:
  currentColor`) in the Plan 4 theming work.
- **Theme system and configurator** — the remote `/theme` token form is due
  for replacement, but a local editor server is not part of the design.
  The approved CSS-only architecture, safe CLI handoff, phased work, and
  explicit no-promise gate for any second component style are tracked in
  [`theme-system-roadmap.md`](theme-system-roadmap.md).
- **Icon search** — the icon gallery page (`site/examples/icon`) ships v1
  as a static grid of ~40 popular icons plus a "1,748 total" note; a
  searchable/filterable index over the full Lucide set is not built.
- **Copy-button success feedback** — the site's copy-to-clipboard button
  (`data-site-copy`, wired in `web/site.js`) copies but gives no
  success/failure affordance (e.g. a checkmark swap or toast); noted during
  Plan 4 Task 2 review as a minor deferred to here.

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
- **`gsxui theme` (local)** — the theme editor ships remote (on the site,
  `/theme`) for v1; a local `gsxui theme` command was deferred, open
  question: embed a built CSS artifact in the CLI binary (stays in sync
  with the site's editor, but bloats/staleness-risks the binary) vs. have
  it reuse the calling project's own Tailwind build (accurate to that
  project's tokens, but requires shelling out to its build tooling). Needs
  a decision before implementation.
- **Icon search** — the icon gallery page (`site/examples/icon`) ships v1
  as a static grid of ~40 popular icons plus a "1,748 total" note; a
  searchable/filterable index over the full Lucide set is not built.
- **Copy-button success feedback** — the site's copy-to-clipboard button
  (`data-site-copy`, wired in `web/site.js`) copies but gives no
  success/failure affordance (e.g. a checkmark swap or toast); noted during
  Plan 4 Task 2 review as a minor deferred to here.

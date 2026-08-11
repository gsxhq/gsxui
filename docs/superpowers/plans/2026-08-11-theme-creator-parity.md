# Theme Creator Parity (Phase 2) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development or superpowers:executing-plans. Steps use checkbox (`- [ ]`) syntax.

**Goal:** Make gsxui's `/theme` comparable to shadcn's `/create`. Phase 1 delivered style parity (all 8 upstream styles). This closes the remaining axes and the preview-breadth gap.

**Architecture:** Extend the existing preset model (`internal/preset`) with new axes, extend the immutable client state (`web/theme-state.js`) with new reducers, and author substantially more preview content in the existing `gallery.gsx.src` → per-style fan-out pipeline. No new frameworks; the preview stays server-rendered with no component JS.

**Dossier (line-cited research, read before starting):** `/private/tmp/claude-501/-Users-jackieli-personal-gsxhq/ed3f6834-0598-4db1-b893-3b2b28ce0688/scratchpad/creator-dossier.md`

## Global Constraints

- **npm-free path must keep working.** Whatever `gsxui init`/`add`/`apply` hands a consumer must need no npm, CDN or build step. This governs the font decision (Task 3).
- The preview renders **all 8 styles simultaneously** in one document, toggled by `hidden` (`site/pages/theme_preview.gsx:42-65`). Any content expansion multiplies by 8 — measure before committing.
- `compact.go`'s share-code codec is **append-only** (its own rule at `compact.go:26-27`); new axes append new bit fields, never reorder existing ones. Old share codes must keep decoding.
- Preserve our existing advantages — do not regress them: arbitrary foreign `theme.css` import, dual client+server preset validation, `init --preset`/`apply --preset` plumbing, the 17-hue accent palette.
- `gallery_test.go` enforces that every cataloged component appears in the gallery. Keep it passing as content grows.
- Commit format `<type>: <summary>`. Inner loop: the narrow package test. `make check` once at the end.
- **Deferred, do not build:** the icon-library axis. Rationale: ~14,000 generated lines for 4 more libraries plus a client-side icon-swap runtime the preview deliberately lacks, for the lowest-signal axis in a theme creator. Record this in the docs so it reads as a decision, not an omission.

---

### Task 1: Cheap parity wins — URL sync, undo/redo, menu accent

**Files:** `web/theme-state.js`, `web/theme.js`, `site/pages/theme.gsx`, `internal/preset/{preset,catalog,compact}.go`, plus their tests.

Three independent wins, one commit each.

- [ ] **1a. Live URL state sync.** shadcn updates the address bar in place on every change; we require clicking "Copy share URL". Push the compact code into `history.replaceState` on each state commit (never `pushState` — that would spam history and fight undo). Read the existing share-code path in `web/theme.js` first; the codec already exists, this is wiring. On load, an existing `?preset=` must still hydrate (it already does — verify, don't rebuild). Test: change an axis, assert `location.search` carries a code that round-trips back to the same state.
- [ ] **1b. Undo/redo.** `web/theme-state.js` is already immutable-reducer shaped, which makes this a history stack over committed states, not a rewrite. Add `undo()`/`redo()` with a bounded stack (50), wire ⌘Z/⌘⇧Z plus visible buttons. Do NOT push a history entry per keystroke on a text/colour input — debounce to the same commit boundary the existing state uses. Test the reducer directly: apply three changes, undo twice, redo once, assert state.
- [ ] **1c. Menu accent (the cheap half of the menu-chrome axis).** Per dossier §5(e), accent is a token swap and fits the existing model; colour/appearance need class-level re-theming machinery we don't have. Ship accent only: add the axis to `PaletteSelection`, resolve it like the existing theme overlay, expose a Subtle/Bold control. **Do not** build the colour/appearance halves — note them in the docs as deliberately out of scope with the reason.

Commits: `feat: live URL state sync in theme editor`, `feat: undo/redo in theme editor`, `feat: menu accent axis`

---

### Task 2: Chart-colour axis

**Files:** `internal/preset/{catalog,preset}.go`, `web/theme-state.js`, `site/pages/theme.gsx`, `site/stylepreview/gallery.gsx.src`, tests.

Tokens `chart-1..5` already exist (Phase 1, `preset.go:157-161`). Two pieces:

- [ ] **2a. The axis.** Extend `PaletteSelection` with `ChartColor`, resolving by overlaying **only** the five chart keys from the chosen hue onto the current selection — this mirrors upstream (`registry/config.ts:719-729`) and reuses the exact shape `resolvePalette` already has for the theme overlay. Append to the compact codec per the append-only rule. Control in `theme.gsx` reusing the existing hue-picker markup.
- [ ] **2b. A visible consumer, deliberately minimal.** The preview has no component JS and no chart component exists. Author **one static card** in `gallery.gsx.src` — a bar or donut painted with `var(--chart-1)`..`var(--chart-5)` in plain SVG/CSS. This is a colour swatch with a chart's shape, not gsxui's first charting component. **Explicitly do not** build a real chart; if that starts happening, stop and flag it.

Commit: `feat: chart-colour axis with static preview card`

---

### Task 3: Font pair picker

**Files:** `internal/preset/{preset,catalog,compact,css,json}.go`, `web/{theme-state,theme-preview}.js`, `site/pages/theme.gsx`, font assets, tests.

Per dossier §5(b) this is a **schema change**, not a picker: `ThemeValues` is colour-only and validated with `csscolorparser.Parse` per token, so font families cannot live in that map. Upstream agrees — font is a separate field from the colour theme (`preset.ts:170-181`), applied by setting `--font-sans`/`--font-heading` directly rather than through token swapping.

- [ ] **3a. Schema.** Add `FontSans`/`FontHeading` to `Preset` (not to `ThemeValues`), with a font-family validator distinct from the colour validator. Add `FontChoices()` to `catalog.go` parallel to `BaseColorChoices`. Append new bit fields to `compact.go`. Emit `--font-sans`/`--font-heading` from `css.go`/`json.go` — currently absent. **Reuse upstream's exact var names** (`--font-sans`, `--font-heading`); their shipped consumer artifact already uses them, so we stay CSS-compatible.
- [ ] **3b. Assets, npm-free.** Upstream delegates to `@fontsource-variable/*` npm packages and `next/font/google` — neither is available to us. **Self-host a curated set of variable WOFF2 files** (start with ~6-8 covering the 8 styles' character: Inter, Geist, Figtree, JetBrains Mono, Noto Sans, Playfair Display) with `@font-face` + `font-display: swap`, shipped alongside `ui/` so neither our preview nor a consumer project needs npm or a CDN. Document the licence of each font shipped.
- [ ] **3c. Wiring.** `theme-state.js` gains `selectFont`/`selectFontHeading`; `theme-preview.js` sets the two properties on the preview root the way upstream's provider does. Control in `theme.gsx`: two dropdowns, heading and body.

Commit: `feat: font pair picker with self-hosted variable fonts`

---

### Task 4: Preview breadth

**Files:** `site/stylepreview/gallery.gsx.src`, `site/pages/{theme,theme_preview}.gsx`, `web/theme-preview.js`, `site/stylepreview/gallery_test.go`.

The biggest first-impression gap: ours is a 15-card component gallery, theirs is ~29 and ~33 realistic product blocks across two tabbed pages. Per dossier §4 the generation pipeline is a pure package-clause rewrite, so authoring richer content is unconstrained.

- [ ] **4a. Measure first.** All 8 styles render simultaneously; doubling cards doubles an already-8× document. Before authoring, measure the current preview document's size and render time, and decide: lazy-render hidden style sections, or accept the growth. Record the numbers. This gates how far 4b goes.
- [ ] **4b. Author ~15 more cards** in `gallery.gsx.src` in the same style as the existing 15 (~60 lines each), but as *realistic product surfaces* rather than component showcases — an invoice table, a settings panel, a pricing row, a stat/KPI row, a chat thread, a team-members list, a booking form, an activity feed. Reuse the 54 canonical components; the point is realistic copy and composition, not new components.
- [ ] **4c. Two-page split.** Add a page grouping + tab control mirroring the existing per-style section toggle (`theme_preview.gsx:42-65`, `theme-preview.js:115-124`) — an established pattern, cheap.
- [ ] **4d.** Extend `gallery_test.go`'s exhaustiveness assertion to the new content; regenerate per-style visual snapshots and review all 8 by eye.

Commit: `feat: expand theme preview to two pages of realistic product blocks`

---

### Task 5: Docs and final gate

- [ ] **5a.** Update `/theme` copy and any docs describing the editor's axes. State the deferred items as decisions with reasons: icon-library axis (cost vs signal), menu colour/appearance (needs class-level re-theming machinery).
- [ ] **5b.** `make check` — exit 0. If any `site/examples` `.gsx` changed, run `make highlight` in the same commit (repo CI rule).
- [ ] **5c.** Screenshot the finished creator in light and dark, all 8 styles, and put the paths in the final report.

Commit: `docs: theme creator axes and deferred-by-design notes`

## Self-Review

- **Gap coverage:** audit gaps (1) preview breadth → Task 4; (2) fonts → Task 3; (3) icon library → deliberately deferred, recorded in Global Constraints and Task 5a; (4) chart-colour → Task 2; (5) menu-chrome → Task 1c (accent half; other halves deferred with reason); (6) polish → Tasks 1a/1b. Shuffle is dropped: per the dossier it is 25 curated codes upstream, and it has little payoff until more axes exist — not worth building now.
- **Preserved advantages** are named in Global Constraints so a task cannot regress them silently.
- **Placeholders:** none. Each task names files, the decision already made, and the commit.
- **Ordering rationale:** cheap wins first (Task 1) so the creator improves immediately; Task 4 last because 4a's measurement may change its own scope.

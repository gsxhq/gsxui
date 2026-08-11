# Style System Parity — Design

*2026-08-11 · Status: approved · Goal: gsxui's theme creator comparable to shadcn's `/create`*

## 1. Goal

Make gsxui's `/theme` editor comparable to shadcn's theme creator
(https://ui.shadcn.com/create). That is the success criterion; everything
below serves it.

The creator cannot be comparable today for one dominant reason: it offers
**2 styles, one of which (`maia`) is Button-gated** in the CLI. shadcn's
offers 8, all real. Style count is therefore the critical path, and the work
splits into two phases with a hard ordering.

## 2. Decisions taken

| Decision | Choice | Rationale |
|---|---|---|
| Style set | **All 8 upstream styles**: vega, nova, maia, lyra, mira, luma, sera, rhea | Exact parity with shadcn's creator; upstream declares them in `apps/v4/registry/styles.tsx` |
| Our house `nova` | **Retired — replaced by upstream's `style-nova.css`** | One lineage for the whole set. Completes the density-retarget work (`docs/superpowers/plans/2026-07-24-nova-density-*.md`) by porting wholesale instead of hand-matching metrics |
| Port method | **Mechanical transformation tool**, not hand-authoring | 8 × 54 = 432 recipe files. Hand-authoring is ~100+ days; the formats map deterministically (§4) |
| Upstream source | Local checkout `/Users/jackieli/personal/shadcn-ui` (`apps/v4/registry/styles/style-*.css`), pinned by commit SHA in every generated header | Reproducible, reviewable provenance; no network at build time |
| Chart tokens | **Add `chart-1..5` to the preset schema in Phase 1** | Upstream styles reference them; they cost nothing without a chart component and their absence would block Phase 2's chart-colour axis |
| Monetisation | **None.** Everything free | Explicit user direction |

**Consequence to accept knowingly:** replacing our house nova with a true
upstream port changes colours, shadows and radii — not only density, which
is all the retarget plan touched. The default appearance of every component
shifts. This is intended, but it is a visible break for anyone already on
`ui/*.gsx`, and the CHANGELOG must say so plainly.

## 3. Phase split (hard ordering)

**Phase 1 — style engine and 8 ports.** The porter, the 432 recipe files,
style-aware CLI vendoring, a per-style visual gate, chart tokens in the
schema. Ends when `gsxui add <any component>` under any of the 8 styles
vendors that style's real output and the `/theme` picker offers 8 working
styles.

**Phase 2 — creator parity.** Preview breadth, font pair picker,
icon-library picker, chart-colour axis, menu-chrome axis, URL state sync /
shuffle / undo-redo. Scoped from the audit in §7.

Phase 1 is the critical path; Phase 2 is cosmetic until it lands.

## 4. Phase 1 architecture

### 4.1 The transformation

Upstream format (`style-maia.css`):

```css
.style-maia {
  /* MARK: Accordion */
  .cn-accordion-trigger {
    @apply **:data-[slot=accordion-trigger-icon]:text-muted-foreground gap-6 p-4 …;
  }
}
```

gsxui format (`registry/styles/<style>/<component>.css`): one
`@layer components` rule per role, named
`.gsxui-recipe-<component>[-<slot>][-<dim>-<value>]`, `@apply`-only.

Three deterministic rules convert one to the other:

1. **Slot rename.** `.cn-<component>[-<slot>]` →
   `.gsxui-recipe-<component>[-<slot>]`. Identity for the large majority; a
   declared override table (§4.3) handles the rest.
2. **Dimension re-split.** Upstream folds variants into attribute utilities
   on one class (`data-[size=sm]:max-w-xs`); we split them into per-dimension
   rules. `registry/canonical/shapes/<c>.go` already declares which
   dimensions and values exist, so the split is checkable in both
   directions: every `data-[<dim>=<val>]:` prefix whose `<dim>` is a declared
   dimension moves its utility into the `-<dim>-<val>` rule; a prefix naming
   an *undeclared* dimension is an error, not a silent drop.
3. **Descendant decomposition.** `**:data-[slot=<x>]:<util>` on a parent rule
   moves `<util>` into slot `<x>`'s own rule when `<x>` is one of our slots.
   This is what makes the port tractable: most of what a naive survey calls
   "no upstream counterpart" (trigger-icon, close-icon, …) is simply
   expressed as a descendant selector upstream.

Anything that survives all three rules unmapped is **reported, never
silently dropped** — the porter exits non-zero with a per-component list.

### 4.2 The porter

`internal/stylegen/port/` with a CLI entry (`go run ./cmd/stylegen port
--upstream <path> --style <name>|all`). Pure function of (upstream CSS,
shapes, mapping table) → recipe files. It writes
`registry/styles/<style>/<component>.css` with a provenance header carrying
the upstream commit SHA, source file and line range, matching the existing
header convention in `registry/styles/maia/button.css`.

The porter is **re-runnable and idempotent**: re-porting after an upstream
bump produces a reviewable diff. It never edits generated `.gsx`/`.x.go` —
those stay the output of the existing `stylegen generate` pipeline, which is
unchanged.

### 4.3 The mapping table

A single declared Go table (`internal/stylegen/port/mapping.go`) holding
only the divergences, each with a one-line reason:

- **Structural roots we have and upstream doesn't** — our native-first
  markup needs a co-location wrapper (`<div class="contents">`) where
  Radix's Root renders nothing. These take no upstream utilities.
- **`dialog` overlay fusion** — upstream's separate `.cn-dialog-overlay`
  must fold into our content rule's `::backdrop`, because native `<dialog>`
  fuses what Radix keeps apart. Same for alert-dialog, drawer, sheet.
- **Upstream sub-classes we have no slot for** — e.g. select's
  `scroll-up-button`/`scroll-down-button` (native popover needs no JS scroll
  buttons), the `*-aria` React-Aria duplicates. Explicitly ignored, listed.

### 4.4 Hand-authored residue

Four cases the porter cannot serve, authored once per style and marked as
such in their headers:

- **aspect-ratio, collapsible, spinner** — no upstream section exists in any
  of the 8 styles. Decision: these stay **style-invariant** (one recipe
  shared by all 8), because inventing 8 divergent looks for components
  upstream doesn't style would be fabrication, not parity.
- **toaster** — upstream has only a 1-rule `Toast` section and no Toaster.
  Derive from toast; style-invariant beyond what toast gives.

### 4.5 CLI: close the correctness trap

`internal/cli/add.go:155-160` and `apply.go:155-162` currently style-switch
for `button` only; every other component vendors from `ui/<name>.gsx`
(nova's output). The `maia` gate is what stops that being wrong. Make both
read `registry/generated/<style>/<name>.gsx` for **every** component, then
remove the gate. Test: vendoring any component under any style yields that
style's bytes.

`ui/*.gsx` remains the `DefaultStyle` copy, and `DefaultStyle` stays `nova`
— now meaning upstream nova.

### 4.6 Gates

Existing machinery proves a style is **complete**: `recipe.Conform` checks
every slot × dimension × value against the shape, and `--check-layers`
proves nothing is dead in the cascade. Nothing proves a style is
**correct**.

Phase 1 adds a **style axis to `jstest/specs/style-visual.spec.ts`**, so
each style has its own committed snapshot set rendered from the existing
`site/stylepreview/<style>/` gallery. Without this, 432 ported files would
be reviewed by eye. This gate is a Phase 1 deliverable, not a follow-up, and
it lands **before** the bulk port so the port's diffs are reviewable.

### 4.7 Documentation correctness

`site/pages/theme.gsx:281` claims "Both styles render the full component
catalogue", which is false today. It becomes true when §4.5 lands; the copy
is corrected in the same change either way.

## 5. Testing

- **Porter unit tests** — each transformation rule in isolation, plus
  golden-file tests for three representative components (accordion: simple
  slots; alert-dialog: dimensions + descendant selectors; dialog: the
  overlay fusion).
- **Round-trip** — every ported recipe must satisfy `recipe.Conform`
  against its shape; the porter fails loudly otherwise.
- **Unmapped reporting** — a test asserting the porter reports rather than
  drops an unrecognised `cn-*` class.
- **Visual** — per-style snapshots (§4.6).
- **CLI** — vendoring under each style yields that style's bytes.

## 6. Risks

1. **The port is checkable for completeness, not correctness.** Mitigated by
   §4.6's visual gate, but a snapshot only catches *change*, not *wrongness*
   on first capture. First capture per style needs human review.
2. **Default appearance shifts** for existing users (§2). Mitigated by
   CHANGELOG prominence; gsxui is pre-release with no tags, so no
   compatibility promise is broken.
3. **Upstream drift.** The porter pins a SHA; re-porting is a reviewable
   diff. Accepted.
4. **`make highlight` gate** — per MEMORY.md, any `site/examples` `.gsx`
   change needs `make highlight` in the same commit or CI blocks the deploy.
   Style preview galleries are generated; the port must not silently stale
   the highlight output.

## 7. Phase 2 scope (from the audit, not yet designed in detail)

Gaps versus shadcn's `/create`, ranked:

1. **Preview breadth** — theirs is ~15-20 realistic product blocks across 2
   tabbed pages; ours is a 15-card component gallery. Biggest
   first-impression gap.
2. **Font pair picker** (heading + body) — we have none.
3. **Icon-library picker** — they offer 5; we ship Lucide only.
4. **Chart-colour axis** — unblocked by Phase 1's chart tokens.
5. **Menu-chrome axis** (colour × appearance × accent).
6. **Polish** — live URL state sync, shuffle, undo/redo.

**Advantages to preserve, not regress:** we import arbitrary foreign
`theme.css` (their "Open Preset" only reloads their own codes); we validate
the preset schema client- *and* server-side; `init --preset` / `apply
--preset` is real, shipped and tested; our accent palette is larger (17 hues
versus the 11 confirmed in theirs).

## Self-review

- **Placeholders:** none. Phase 2 is deliberately scoped-not-designed, and
  says so.
- **Consistency:** §2's "replace nova" and §4.5's "`DefaultStyle` stays
  nova" are consistent — the name persists, the lineage changes.
- **Scope:** Phase 1 is one plan's worth. Phase 2 gets its own spec.
- **Ambiguity:** §4.4's style-invariant decision is stated explicitly rather
  than left to the implementer.

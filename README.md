# gsxui

shadcn-style components for [gsx](https://github.com/gsxhq/gsx) — copy-in,
type-checked, server-rendered.

## Install

    go install github.com/gsxhq/gsxui/cmd/gsxui@latest

In your project (a Go module):

    gsxui init          # tokens css, js runtime, class merger, gsx.toml wiring
    gsxui add dialog    # vendors dialog + its deps (button), regenerates
    gsxui list          # what's available

You own the vendored code. `gsxui add` never touches a modified file unless
you pass `--overwrite`.

After upgrading the gsxui binary, re-run `gsxui add <name> --overwrite` to
refresh vendored components — this discards local edits to those files.

**Status: pre-release.** v1 component set (20 components + icon) is complete.
The showcase site and theme editor are in progress.

- Components live in `ui/<name>/` — a `.gsx` source (JSX-style, named
  parameters, fallthrough attrs) plus a behavior `.js` when interactive.
- `make test` regenerates and tests everything; `make check` adds JS syntax
  and gofmt checks.
- Divergences from the shadcn/ui reference: `docs/jsx-parity.md`.

## Components

**Form controls:** button, checkbox, input, label, radio, selectbox,
switchctl, textarea

**Display:** alert, avatar, badge, card, separator, skeleton, table

**Overlay:** dialog, dropdown, tooltip

**Navigation:** accordion, tabs

**Primitives:** icon (Lucide, generated — a dependency of selectbox,
accordion, and dropdown's own future items, not usually added directly)

Some native-first components (checkbox, radio, switchctl, selectbox,
accordion) trade a slice of shadcn's Radix-driven behavior for a real
`<input>`/`<select>`/`<details>` element — zero client JS, browser-native
`:checked`/`:disabled`/exclusivity semantics. dropdown and tooltip trade
Radix's Portal for the native popover API. Every divergence, with its
rationale, is ledgered in `docs/jsx-parity.md`.

## Running the site

The showcase site (`site/`) is a gsx-built, server-rendered Go app that
imports `ui/` directly — it's the proof the components work, plus the
theme editor. Two commands from a fresh clone:

    npm install

then, for the dev loop (Vite HMR proxied in front of a live-reloading Go
server):

    make site-dev

or, to build the production bundle and run it (Vite assets embedded into
the binary, served without a dev proxy):

    make site

Either way, open the printed URL (`make site-dev` prints Vite's dev URL;
`make site` serves directly on `$GO_PORT`, falling back to `$PORT`, then
8080, if neither is set).

### Deploying

`site/Dockerfile` is a three-stage build (Vite production build → `go
build` against the committed `.x.go` + built `dist/` → distroless static
run stage) that binds `$PORT` (Cloud Run's convention; falls back to 8080).
Build it from the **repo root** as context:

    docker build -f site/Dockerfile -t gsxui-site .
    docker run -p 8080:8080 gsxui-site

`cloudbuild.yaml` (repo root) builds, pushes to Artifact Registry, and
deploys to Cloud Run in one Cloud Build config — pattern-matched from
[gsx's playground deploy](https://github.com/gsxhq/gsx/blob/main/cloudbuild.yaml).

## Post-v1 backlog

Deferred out of v1 scope, tracked here rather than in the parity ledger's
per-component GAP notes (see those for the detailed rationale):

- **Custom listbox select** — shadcn's full Radix `Select` (floating panel,
  check-mark item indicator, keyboard typeahead). v1 ships a styled native
  `<select>` instead (`ui/selectbox`); the Radix-equivalent listbox visual
  is not built.
- **Dropdown checkbox/radio items + submenus** — `DropdownMenuCheckboxItem`,
  `DropdownMenuRadioGroup`/`DropdownMenuRadioItem`,
  `DropdownMenuSub`/`SubTrigger`/`SubContent`, `DropdownMenuGroup`. Only the
  base item/label/separator/shortcut set ships in v1.
- **Tooltip delay-groups** — `TooltipProvider`'s shared `delayDuration`/
  skip-delay-group coordination across multiple tooltips on a page. v1
  hard-codes a fixed per-trigger open delay, no cross-tooltip grouping.
- **CSS anchor positioning migration** — dropdown and tooltip currently
  position via a hand-rolled `getBoundingClientRect()` + `position: fixed`
  calculation in JS. Once CSS anchor positioning (`anchor()`/
  `position-anchor`) reaches Baseline across browsers, both can drop that
  JS for native, collision-aware placement.
- **Popover exit-animation strategy** — verified inert, strategy TBD:
  dropdown/tooltip's native popover is already `display: none` at the
  moment the `toggle` (newState=closed) handler stamps
  `data-state="closed"`, so `data-[state=closed]:animate-out` never runs
  and closing snaps (open-path `animate-in` unaffected); accepted for v1.
  An animated-close strategy (`beforetoggle`/`allow-discrete`, or
  dialog-style `requestClose`) remains adoptable once designed.
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
- **gsx syntax highlighting for source blocks** — the component pages'
  `<pre><code>` source panels (`site/pages/component.gsx`) render escaped
  plain text; no `.gsx`-aware highlighter exists yet. Revisit once/if one
  does (or a generic JS-family highlighter proves good enough for gsx's
  JSX-in-Go syntax).
- **Icon search** — the icon gallery page (`site/examples/icon`) ships v1
  as a static grid of ~40 popular icons plus a "1,748 total" note; a
  searchable/filterable index over the full Lucide set is not built.
- **Copy-button success feedback** — the site's copy-to-clipboard button
  (`data-site-copy`, wired in `web/site.js`) copies but gives no
  success/failure affordance (e.g. a checkmark swap or toast); noted during
  Plan 4 Task 2 review as a minor deferred to here.

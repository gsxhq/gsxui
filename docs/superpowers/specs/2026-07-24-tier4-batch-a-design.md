# Tier 4 Batch A — resizable, combobox, sidebar

Design for the first Tier 4 batch (`docs/component-roadmap.md` § Tier 4).
Scope decided 2026-07-24: the three components with the best value-per-effort
and no new machinery debt. Deferred to later rounds: the menu family
(dropdown checkbox/radio/submenu → menubar → navigation-menu, which share one
body of submenu/positioning work), calendar (its own spec), chart (needs a
Go/JS charting answer first).

Coverage after this batch: 50 of shadcn's 61 shipped.

## 0. Reference sources — what we track and why

shadcn's registry forked since our 2026-07-23 audit. Three facts, all
verified in the local `shadcn-ui` checkout at `3f47b9113`:

1. `apps/v4/next.config.mjs` redirects `/docs/components/:name` →
   `/docs/components/base/:name`. **The live site defaults to the Base UI
   variant.**
2. The registry split into `registry/bases/{radix,base,aria}/ui/`. Those
   directories carry the newest commits; `registry/new-york-v4/ui` has not
   been touched since before the split.
3. In `bases/*`, components no longer carry Tailwind utilities. They emit
   named classes — `bases/base/ui/button.tsx` renders
   `cn-button cn-button-variant-default cn-button-size-default` — with the
   utilities extracted into `registry/styles/style-*.css`.

Consequence: all 47 shipped gsxui components differ from their `bases/base`
twins, but overwhelmingly because of that class extraction, not a design
change.

**Decision: the reference method does not change.** `new-york-v4/ui` stays
the markup reference — it is the only place the markup exists in
inline-utility form, which is gsxui's authoring model — and
`registry/styles/style-nova.css` stays the visual source of truth. That is
the same two-source method the nova density retarget established
(`docs/jsx-parity.md` `## nova density`).

Two riders:

- Where `new-york-v4` and nova disagree, **nova wins**; `bases/base/ui/` is
  the tiebreak for *markup structure* questions nova's CSS cannot answer.
- **No re-sync of the 47 shipped components is implied by this finding.** If
  a re-sync is ever wanted it is its own scoped project, not a Tier 4 task.

This section is the record of that adjudication so a future reader does not
re-litigate it from the raw diff.

## 1. resizable

Reference: `new-york-v4/ui/resizable.tsx` (53 lines, wraps
`react-resizable-panels`), demos `resizable-demo{,-with-handle}.tsx`,
`resizable-vertical.tsx`. Nova adds `.cn-resizable-handle-icon` only.

### Parts

| gsx | shadcn | element |
|---|---|---|
| `ResizablePanelGroup` | `ResizablePanelGroup` | `div` |
| `ResizablePanel` | `ResizablePanel` | `div` |
| `ResizableHandle` | `ResizableHandle` | `div[role=separator]` |

### Layout model

The group is a flex container; `aria-orientation=vertical` flips it to
`flex-col` (shadcn's own selector: `aria-[orientation=vertical]:flex-col`).
Each panel's size is a `flex-basis` percentage written as an inline style by
the **server** from the `defaultSize` param, so the split is correct on first
paint with no JS and no layout shift. `minSize`/`maxSize` are server-stamped
as data attributes and enforced by JS during drag.

The handle is 1px (`w-px`, or `h-px` when horizontal-oriented per shadcn's
inverted naming) with a wider `::after` hit target — verbatim from the
reference class string.

### Behavior — `ui/resizable.js`

- **Pointer drag**: `pointerdown` on the separator → `setPointerCapture` →
  `pointermove` recomputes the flex-basis of the two adjacent panels from the
  group's content box, clamped to each panel's min/max → `pointerup`
  releases. Capture is what makes the drag survive the cursor leaving the
  1px handle.
- **Keyboard**: with the separator focused, Arrow keys nudge by a fixed step,
  Home/End drive the adjacent panel to its min/max. Direction pairs with
  orientation (Left/Right horizontal, Up/Down vertical).
- **ARIA**: `role="separator"`, `aria-orientation`, `aria-valuenow` /
  `aria-valuemin` / `aria-valuemax` tracking the panel before the handle, and
  `tabindex="0"`. Values are server-rendered and JS-updated, matching the
  slider precedent.
- Nested groups work with no extra machinery — each group resolves sizes
  against its own content box.
- **Announce, don't persist**: on drag end / keyboard commit, emit
  `gsxui:change` with `{ sizes }` on the group. Consumers who want
  `autoSaveId`-style persistence write those three lines themselves against
  whatever store they use. Same principle as sidebar (§3, §4).

### Deferred (ledger as GAP)

`autoSaveId` localStorage persistence — **deliberately not ported**, replaced
by the `gsxui:change` event above. Also deferred: collapsible panels
(`collapsible`/`collapsedSize`/`onCollapse`); the imperative panel API
(`ImperativePanelHandle`).

## 2. combobox

Reference: `new-york-v4/ui/combobox.tsx` (310 lines, wraps
`@base-ui/react` Combobox). Nova ships the full `.cn-combobox-*` set.

**Port the component, not the old recipe.** `examples/combobox-demo.tsx` is
the legacy Popover+Command composition; the live site's combobox page renders
`combobox-basic`, `-clear`, `-groups`, `-invalid`, `-disabled`,
`-input-group`, i.e. the component. Nova having `.cn-combobox-*` confirms it.

### Parts (v1)

`Combobox` (root), `ComboboxInput`, `ComboboxTrigger`, `ComboboxClear`,
`ComboboxContent`, `ComboboxList`, `ComboboxItem`, `ComboboxGroup`,
`ComboboxLabel`, `ComboboxEmpty`, `ComboboxSeparator`, `ComboboxValue`.

`ComboboxInput` wraps `InputGroup` + `InputGroupInput` +
`InputGroupAddon align="inline-end"` holding the trigger and clear buttons,
gated by `showTrigger` / `showClear` params — verbatim from the reference.
That makes `input-group` a real dependency, so it lands in `Deps`.

### Reused machinery

- **Anchoring / open-close**: the shipped popover discrete-transition block
  (`docs/jsx-parity.md ## animations`), byte-identical to dropdown's, as
  `select` already does.
- **Value model + hidden form bridge**: the shipped `select` model — a
  `sr-only` real control carrying the value so non-JS form posts and
  `FormData` both work.
- **Item indicator**: `data-slot="combobox-item-indicator"` check glyph from
  `ui/icon`, so `icon` is a derived dep (the carousel precedent).
- **Events**: `gsxui:select` with `{ value }` on the root plus a native
  `change` on the hidden bridge — exactly what `select.js` already does, so
  the two are interchangeable to a consumer.

### Filtering — the one open question the source map must settle

Base UI's default combobox filter is a collator-based **contains** matcher,
not cmdk's `command-score` ranking. The source-map task must read
`@base-ui/react`'s `dist` to confirm (the package is not in the local
checkout — mark `derived-not-read` until it is).

**Decision, pending that confirmation: match Base UI.** `combobox` and
`command` stay independent modules with no shared JS — `combobox.js` gets its
own matcher rather than importing `command.js`'s scorer. Rationale: it is
what the reference does, and a JS-only cross-module import would be invisible
to `registry.Deps` (derived by go/parser over the generated `.x.go`), so it
would silently break vendoring. Ledger the divergence from cmdk ranking.

If the source map finds Base UI actually ranks rather than filters, revisit —
but resolve it by duplicating the needed scorer into `combobox.js`, not by
introducing an undeclarable dependency.

### Deferred (ledger as GAP)

`ComboboxChips` / `ComboboxChip` / `ComboboxChipsInput` / `ChipRemove` — the
multi-select chips UI. `ComboboxCollection` — a React data-mapping helper
with no gsx equivalent (gsx callers write a `for` loop). `useComboboxAnchor`
— a React ref hook; gsx callers pass an anchor selector.

## 3. sidebar

Reference: `new-york-v4/ui/sidebar.tsx` (726 lines). Nova ships the full
`.cn-sidebar-*` set. All dependencies — sheet, tooltip, collapsible, button,
input, separator, skeleton — are shipped.

Wide, not deep: ~25 parts, nearly all of them plain markup wrappers. The
design work is confined to state and mobile.

### State — server-rendered input, persistence left to the consumer

React holds `open` in context and mirrors it into a `sidebar_state` cookie
(`max-age` 7 days) which the framework reads back on the server.

**gsxui does not ship the cookie.** Persistence is the consumer's call —
they may use a cookie read in a Go handler, an Alpine store, htmx, their own
JS, or nothing at all. A component that writes `document.cookie` dictates one
of those and fights the other three. What the component owns is rendering the
state it is given and announcing changes:

- `ui.SidebarProvider` takes `open bool` and stamps
  `data-state="expanded"|"collapsed"` server-side. Zero flash, no hydration
  step. Where that `bool` comes from is entirely the consumer's business.
- `sidebar.js` flips `data-state` on toggle (so the component works
  standalone, ephemerally) and emits `gsxui:change` with `{ open }` on the
  provider element — the house convention already used by `toggle`, `tabs`
  and `toggle-group`, via the shared `emit()` helper in `ui/gsxui.js`. That
  event is the whole integration surface: persist in a cookie, POST it, push
  it into an Alpine store — the component neither knows nor cares.
- The provider div carries `--sidebar-width` / `--sidebar-width-icon` inline,
  as the reference does, so consumers can override per-instance.
- The cookie round-trip ships as a **documented site example**, not as
  component code: a Go handler reading `sidebar_state` into `open`, plus the
  three-line `gsxui:change` listener that writes it back. That gives
  readers the shadcn behavior verbatim while keeping it opt-in.

No context substitute is needed anywhere else: `data-state`,
`data-collapsible`, `data-variant`, `data-side` are all stamped on the root
and consumed by descendants through `group-data-*` / `peer-data-*`
selectors — already pure CSS in the reference. This is the toggle-group
precedent (explicit params in, CSS/JS derivation down).

### Mobile — both trees, CSS-gated (decided)

React swaps the entire sidebar for a `Sheet` when `useIsMobile()` is true. A
server render cannot branch on viewport. This is the same obstruction that
kept `drawer-dialog` off the roadmap.

**Decision: render both trees and gate them with CSS.**

- The desktop tree keeps shadcn's own `hidden … md:block` root and
  `hidden … md:flex` container — the reference *already* CSS-hides desktop
  below `md`.
- The mobile tree is the `Sheet` branch, marked `md:hidden`, carrying the
  reference's `data-mobile="true"`, `--sidebar-width: 18rem` override, and
  `sr-only` header/title/description.
- No media-query JS. Both trees exist in the DOM; CSS picks one.

**Known cost, to be ledgered as a GAP**: `children` renders twice, so any
`id` inside sidebar content is duplicated and the document is invalid.
Document it in the component doc comment and the parity ledger, with the
guidance that sidebar content should use classes and `data-*`, not `id`.
Rejected alternatives: dropping the Sheet (deviates from the markup we exist
to mirror); JS relocation of the subtree (stateful DOM shuffling, and the
sidebar disappears entirely on mobile with JS off).

### Parts

Provider/root/gap/container/inner, Trigger, Rail, Inset, Input, Header,
Footer, Separator, Content, Group, GroupLabel, GroupAction, GroupContent,
Menu, MenuItem, MenuButton, MenuAction, MenuBadge, MenuSkeleton, MenuSub,
MenuSubItem, MenuSubButton.

`SidebarTrigger` and `SidebarRail` toggle via JS resolving the provider from
the DOM rather than a context handle. `SidebarMenuButton`'s collapsed-state
tooltip is CSS-gated on `group-data-[collapsible=icon]`, as the reference
does.

`Cmd/Ctrl+B` toggles, per `SIDEBAR_KEYBOARD_SHORTCUT`.

### Deferred

Nothing structural. `variant` (`sidebar` | `floating` | `inset`),
`collapsible` (`offcanvas` | `icon` | `none`), and `side` (`left` | `right`)
all ship — they are class-selector work, not new machinery.

## 4. Cross-cutting constraints

### Components render state; they do not own it

A standing principle, made explicit here because all three components in this
batch brush against it (sidebar's open/collapsed, resizable's panel sizes,
combobox's value):

> A gsxui component takes its state as a server-rendered parameter, reflects
> it in the DOM, and **emits an event when it changes**. It never chooses a
> persistence mechanism.

Consumers reach for cookies read in a Go handler, Alpine stores, htmx,
hand-written JS, or nothing — a component that writes `document.cookie` or
`localStorage` picks one and fights the rest. Where shadcn's React version
bakes in a store (sidebar's `sidebar_state` cookie, resizable's
`autoSaveId`), gsxui ports the *behavior* and drops the *storage*, ledgering
it as a deliberate divergence and showing the round-trip in a site example
instead.

Event naming follows the house convention already in `ui/*.js`:
`gsxui:change` with a detail payload for state-carrying components,
`gsxui:open`/`gsxui:close` for overlays, `gsxui:select` for selection, all
through the shared `emit()` helper in `ui/gsxui.js`.

### Process

Every task follows the established process, unchanged from Tier 3:

- Token-for-token class carryover from `new-york-v4`; nova metric and
  interaction tokens adopted from the start (new ports adopt nova directly).
- Standing house exceptions hold: border→`ring-1` swaps **not** adopted;
  house selector syntax kept (`data-[orientation=…]:` not
  `data-horizontal:`; `focus-visible:ring-[3px]` not `ring-3`;
  `data-[state=…]` not `data-open:`).
- Every drop or deviation ledgered in `docs/jsx-parity.md` under its
  component heading.
- TDD render pins in `ui/<name>_test.go`.
- Site example pages under `site/examples/<name>/`. **Any `site/examples/**/*.gsx`
  change requires `make highlight` in the same commit** — `site/hl` pin tests
  block CI otherwise.
- `registry.HasJS` requires `ui/<name>.js` named exactly; `Deps` is derived
  by go/parser over the generated `.x.go`, so a dependency that exists only
  in JS is invisible to vendoring and is not allowed.
- gsx authoring gotchas apply: `//` inside markup renders as page text;
  markup comments containing `<tag>` break the parser; `{{ x := }}` scopes
  function-wide; `<script>`/`<style>` bodies are raw text; inline `if/else`
  is not valid as a whole attribute value; example directory names cannot
  contain hyphens.
- Subagents must never start or kill processes.

## 5. Execution

Subagent pipeline, as Tiers 1–3: task brief → implementer → `review-package
BASE HEAD` → reviewer → fix round → progress ledger line in
`.superpowers/sdd/progress.md`, then a final whole-batch review on the
strongest model to catch between-task orphans.

Task order, dependency-driven:

0. **Source map** — one doc covering all three, tracing every ARIA, keyboard
   and class claim to real source. Must resolve the Base UI filter question
   (§2) by reading `@base-ui/react` dist.
1. **resizable** — self-contained, no dependencies on the other two.
2. **combobox** — depends on shipped popover/select/input-group/icon.
3. **sidebar** — largest; depends on shipped sheet/tooltip/collapsible.
4. **Roll-up** — whole-batch review, roadmap and README updates.

## 6. Verification

Beyond render pins and the reviewer pipeline, a **full live browser pass
side-by-side against the matching ui.shadcn.com pages** before the batch is
called shipped: resizable drag and keyboard resize, combobox typing/filter/
selection/clear and form submission, sidebar collapse transitions in all
three `variant`s and both `collapsible` modes, and the mobile tree at a
narrow viewport.

This is not optional polish. Static review missed the carousel snap geometry
and the button press-dip in Tier 3; both classes of bug — drag math and
collapse transitions — are present again here.

Known environment limit: hidden or occluded Chrome tabs freeze animation
clocks, leave CSS transitions permanently `pending`, and never run
smooth-scroll or rAF callbacks. DOM-state assertions still hold. Do not
report those artifacts as bugs.

## 7. Out of scope

- Re-syncing the 47 shipped components against `bases/base` (§0).
- The menu family, calendar, chart (later Tier 4 rounds).
- Nova follow-ups already listed in `docs/component-roadmap.md`.
- README post-v1 backlog pruning — housekeeping, tracked separately.

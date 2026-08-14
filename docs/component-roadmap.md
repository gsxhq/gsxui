# Component roadmap — shadcn coverage audit

Audited 2026-07-23 against `shadcn-ui/apps/v4/registry/new-york-v4/ui`
(57 registry components). gsxui ships 20 (naming deltas: `radio` =
their `radio-group`; our native `<select>` port was renamed `select` →
`native-select` on 2026-07-24 to match shadcn's own name for the same
div-wrapped-native-select design — no naming delta there anymore —
freeing `Select`/`select` for the Tier 3 custom listbox below. Our
registry entry for DropdownMenu was itself renamed `dropdown` →
`dropdown-menu` ahead of its slot-axis migration, matching its own
markers and shadcn's `dropdown-menu` — no naming delta there either
anymore).

Ordering is easy → hard **for this codebase**: difficulty is judged
against the machinery gsxui already has (native popover API anchoring
from tooltip/dropdown, `<dialog>` from dialog, `<details>` from
accordion, form-native controls), not against the React source's size.

Work one tier at a time; within a tier, order is roughly as listed.
Every port follows the established process: token-for-token class
carryover with drops ledgered in docs/jsx-parity.md, TDD render pins,
site example pages, browser verification against the shadcn docs.

## Tier 1 — pure markup/CSS, zero JS

| component | approach |
|---|---|
| kbd | styled `<kbd>` (28-line source) |
| spinner | animated svg (16-line source) |
| aspect-ratio | CSS `aspect-ratio` on a div; no Radix needed |
| breadcrumb | `<nav>`/`<ol>` markup + chevron separators from ui/icon |
| pagination | markup over Button styles; prev/next icons |
| progress | styled div pair (shadcn's Radix Progress is just two divs); `value` param sets width |
| empty | empty-state layout block (media/title/description/actions) |
| item | generic media+content+actions row |
| button-group | flex wrapper collapsing child Button borders |
| input-group | input + leading/trailing addon layout |
| field | form-field layout: label + control + description + error (shadcn's non-RHF form primitive — this, not `form`, is our form story) |

## Tier 2 — existing machinery reused, little to no new JS

| component | approach |
|---|---|
| collapsible | native `<details>`, same mechanism as accordion (incl. the `::details-content` animation) |
| alert-dialog | our `<dialog>` machinery + alert-dialog classes; `role="alertdialog"`, no X button, action/cancel buttons |
| sheet | `<dialog>` + side-anchored positioning + slide animations (`data-side` stamping like tooltip/dropdown) |
| toggle | `<button aria-pressed>` + a few lines of JS (or checkbox-based, decide at spec time) |
| hover-card | tooltip machinery (hover-triggered popover) with popover-sized content classes |
| popover | click-triggered `popover="auto"`; anchoring/animation ADAPTs already solved in tooltip/dropdown |
| context-menu | dropdown content machinery, opened at cursor from `contextmenu` event |

## Tier 3 — new interactive machinery, moderate JS — SHIPPED 2026-07-24

All eight ported per `docs/superpowers/plans/2026-07-24-tier3-components.md`
(source maps `…-tier3-source-map-{controls,wrapped}.md`), nova density from
the start; per-component ledger entries in `docs/jsx-parity.md`.

| component | shipped as | deferred sub-features (v1 gaps, ledgered) |
|---|---|---|
| toggle-group | toggle composition + `toggle-group.js` roving focus | vertical orientation; RTL arrow swap |
| slider | native `<input type=range>` + `slider.js` gradient fill | multi-thumb/range, vertical, inverted, minStepsBetweenThumbs; nova's ::after hit-target compensation (unreachable on thumb pseudo-elements) |
| scroll-area | one collapsed div, standard + webkit scrollbar CSS | Corner (no Firefox concept); `type` visibility timing (OS policy wins) |
| select (custom listbox) | dropdown machinery + `select.js` value model / prefix typeahead / hidden form bridge; `NativeSelect` stays alongside | scroll up/down buttons (popper-equivalent anchoring); no-JS form submit carries empty value (bridge is JS-filled) |
| sonner (toasts) | `Toaster` + `Toast` (separately registered) with the lifecycle in `toaster.js`; `toast()` via barrel export | positions other than bottom-right; swipe-to-dismiss |
| drawer | sheet-style `<dialog>` variant, four directions | vaul drag-to-dismiss, snap points, background scaling |
| carousel | CSS scroll-snap + `carousel.js` | embla `loop` wraparound; plugin ecosystem (autoplay ships as a data attribute) |
| input-otp | one hidden real input + presentational slots (`input-otp.js`) | first-paint slot values (JS-populated); anchored-regex pattern constants (per-char class instead) |

## Tier 4 — hard / composite

**command SHIPPED out-of-band 2026-07-24** (built for the site's ⌘K doc
search, ahead of its tier: verbatim port of cmdk's `command-score` ranking
+ cmdk's selection model, `ui/command.{gsx,js}`). It was listed here as
depending on combobox; in practice the dependency ran the other way —
combobox can now build on `command`'s filtering/keyboard model plus the
custom `select`'s value/typeahead machinery, both shipped.

**resizable, combobox and sidebar SHIPPED 2026-07-24**, per
`docs/superpowers/plans/2026-07-24-tier4-batch-a.md` (source map
`docs/superpowers/plans/2026-07-24-tier4-source-map.md`); nova density from
the start, per-component ledger entries in `docs/jsx-parity.md`.

| component | shipped as | deferred sub-features (v1 gaps, ledgered) |
|---|---|---|
| resizable | `ResizablePanelGroup`/`ResizablePanel`/`ResizableHandle`; server-rendered `flex: <n> 1 0px` panels (correct split pre-JS, any container) + `resizable.js` pointer/keyboard drag with ARIA `separator` sync | `autoSaveId` persistence (emits `gsxui:change` instead, caller owns storage); collapsible/collapsedSize panels; imperative panel resize/collapse/expand API |
| combobox | `Combobox` + 11 parts: popover-anchored input/listbox composing `command`'s highlight-cursor model and `select`'s value/form-bridge machinery; `contains()` filter (boolean match, no fuzzy ranking) | multi-select chips (`ComboboxChips`/`Chip`/`ChipsInput`); no shared JS with `command.js` (an import edge invisible to `registry.Deps`'s Go-only parse, ledgered) |
| sidebar | `SidebarProvider` + 22 parts: `open`-driven collapsible rail/icon modes + always-rendered mobile `Sheet` tree, CSS-gated; `sidebar.js` click/rail/`Cmd`/`Ctrl+B` toggle | persistence is the caller's (no cookie shipped; `persisted.gsx` is the recipe); `children` renders twice (dual-tree ADAPT — no `id`s inside a `Sidebar`'s children); tooltip `side`/delay-group params not carried |

**menubar and navigation-menu SHIPPED 2026-07-25**, per
`docs/superpowers/plans/2026-07-25-tier4-batch-b.md` (source map
`docs/superpowers/plans/2026-07-25-tier4-source-map-menus.md`); nova
density from the start, per-component ledger entries in
`docs/jsx-parity.md`. The same batch also closed out `dropdown` and
`context-menu`: both gained `Group`, `CheckboxItem`, `RadioGroup`/
`RadioItem`, and `Sub`/`SubTrigger`/`SubContent` — the seven parts the
README's former "Dropdown checkbox/radio items + submenus" Post-v1
backlog entry tracked — reaching **full parity** with shadcn's
`dropdown-menu`/`context-menu`.

| component | shipped as | deferred sub-features (v1 gaps, ledgered) |
|---|---|---|
| menubar | `Menubar`/`MenubarMenu`/`MenubarTrigger` + the same seven shared item parts, reusing `dropdown-menu.gsx`/`dropdown-menu.js`'s popover machinery a third time; `menubar.js` roving-tabindex bar, open-follows-hover between sibling menus, DOM-nested submenus | no RTL arrow-key swap (codebase-wide gap, not menubar-specific); menu and submenu positioning are both a hand-rolled fixed anchor, not collision-aware flip/shift |
| navigation-menu | `NavigationMenu` + `Item`/`Trigger`/`Content`/`Link`/`Indicator`/`List`; ships `viewport={false}` only — each `Content` renders as its own chromed panel (gsx has no portal, so upstream's shared-viewport architecture doesn't fit) | shared-viewport mode (`viewport=true`, one portalled panel morphing between panel sizes); the six-token `data-[motion=…]` direction-aware slide animation, replaced by the codebase's standard discrete-transition block |

**calendar SHIPPED 2026-07-26**, per `docs/superpowers/plans/2026-07-25-calendar.md`
(source map `docs/superpowers/plans/2026-07-25-calendar-source-map.md`); nova
density from the start, per-component ledger entry in `docs/jsx-parity.md`.

| component | shipped as | deferred sub-features (v1 gaps, ledgered) |
|---|---|---|
| calendar | `Calendar` — a from-scratch month-grid/date-math engine (react-day-picker is not vendored), fixed six-row grid with a Go/JS `monthGrid` twin pair and exact-year agreement tests, single/multiple/range selection (`resetOnSelect` semantics), disabled-date rules, label and dropdown caption layouts, lossless hidden-input form binding including repeated multiple values, reset-to-client-month behavior, full keyboard grid (arrows/Home/End/PageUp/PageDown/Shift-variants), client-side today reconciliation | multi-month (`numberOfMonths`), week numbers, custom modifiers, locales/i18n (needs a real locale parameter threaded through both the Go and JS twins, not a client-only `Intl` call — the server-rendered month must match what the client repaints), the footer slot; cell-level `range_start`/`range_middle`/`range_end` className slots (the button carries range/selection styling instead); the `dropdown` className slot (a real `<select>` replaces upstream's invisible-select-over-a-label trick, so the dropdown caption sizes as `NativeSelect` generic, not nova's compact caption scale) |

**chart SHIPPED 2026-08-14**, per `docs/superpowers/specs/2026-08-14-chart-component-design.md`
(templui/shadcn-templ v2 credited in `NOTICE.md`); per-component ledger entry in
`docs/jsx-parity.md`.

| component | shipped as | deferred sub-features (v1 gaps, ledgered) |
|---|---|---|
| chart | `Chart` container + six root kinds (`BarChart`/`LineChart`/`AreaChart`/`PieChart`/`RadarChart`/`RadialBarChart`) over a templui-v2-derived server model builder and one lazy-loaded client renderer (`ui/chart.render.js`, ported from templui's `chart.js`) — resolves this file's own former deferral reason (Recharts is un-vendorable npm; a single MIT JS file is not) | `Formatter`/`LabelFormatter`/`TickFormatter`/`Tick` render-prop closures and upstream's `TooltipModel.Rows` field (none can cross the JSON boundary); `Label`/`LabelList`/`Cell` render props; interactive demo controls (time-range select, series toggle) shipped as static compositions instead; `BarChart` `layout:"vertical"` (backs shadcn's `bar-horizontal`/`bar-mixed` demos — `ChartModel.Layout` always marshals absent, no Options field can set it); the Recharts keyboard layer's opt-OUT (`accessibilityLayer={false}` — the layer itself is ported and always on, matching Recharts' own default, but nothing can turn it off) |

Nothing remains in this "Remaining" bucket — every tier and chart are shipped.

**Registry-fork note (§0 finding, Task 0 of the 2026-07-24 tier4-batch-a
plan):** shadcn's own docs now default to a Base UI variant, and the
registry itself has split into `registry/bases/{radix,base,aria}/`.
`new-york-v4` stays this project's markup reference regardless — it's the
only remaining registry family still expressed as inline Tailwind
utilities, the form every ADAPT/MECHANISM entry in `docs/jsx-parity.md`
already assumes. This is a naming/structure observation only; no re-sync
of any already-shipped component is implied.

## Not ported (deliberate)

- **form** — react-hook-form bindings; meaningless server-side. `field`
  (Tier 1) is the layout half; validation is Go handler + `aria-invalid`
  patterns, to be shown in the future patterns/page-examples phase.
- **drawer-dialog** — shadcn's own `drawer-dialog.tsx` responsive pattern
  (`Dialog` on desktop / `Drawer` on mobile, swapped via a
  `useMediaQuery("(min-width: 768px)")` hook) has no gsxui equivalent for a
  JS media-query-driven component swap and is out of scope for the
  `drawer` component task itself (`docs/jsx-parity.md` `## drawer` GAP).
  Worth a future patterns/pages-phase example once that phase exists —
  either a small vanilla-JS media-query toggle between a server-rendered
  `Dialog` and `Drawer`, or a CSS-only approach if one turns out to cover
  it.
- **direction** — RTL context provider; HTML `dir` attribute serves gsx.
- **attachment, bubble, message, message-scroller, marker** — the new AI
  chat primitives; defer as a dedicated batch if gsx targets chat UIs.

## Nova follow-ups

Deferred out of the 2026-07-24 nova density retarget (`docs/superpowers/plans/2026-07-24-nova-density-map.md`
"Markup prerequisites" items 3–4; retarget itself was metric-tokens-only,
see `docs/jsx-parity.md` `## nova density`). Each item below needs new
markup/parts/params, not just a class-string edit — that's why they're
roadmap items rather than density-map deltas.

New parts nova ships that gsxui doesn't yet have a slot for:

- **AlertAction** — `absolute top-2 right-2` action slot on `Alert`, plus the
  root's `has-data-[slot=alert-action]:pr-18` reservation.
- **AlertDialogMedia** — `mb-2 size-10 rounded-md` icon/media slot in
  `AlertDialogHeader`, needs the header's own media-conditional grid-rows
  variant too.
- **PopoverHeader / PopoverTitle / PopoverDescription** — structured content
  parts for `Popover`, mirroring dialog's header/title/description split.
- **ProgressLabel / ProgressValue** — labeled-progress parts (current
  `progress` ships the bar only).
- **AvatarBadge / AvatarGroup / AvatarGroupCount** — status-badge overlay,
  stacked-avatar wrapper, and the `+N` overflow-count part (`.cn-avatar-group-count`,
  `size-8 rounded-full text-sm` + lg/sm variants).
- **Item size-xs** — a fourth `Item` density below the current smallest.
- **ItemGroup gap changes** — nova's `ItemGroup` gap tokens differ from
  what shipped; needs its own map entry once scoped.

New size axes (a `size`/`data-size` variant prop plus its class-selector
plumbing) nova adds that aren't ported:

- **alert-dialog** `size` (`default` | `sm`) — content max-width axis.
- **avatar** `size` (`sm` | `lg`) — `data-[size=lg]:size-10` / `data-[size=sm]:size-6`,
  plus `AvatarGroupCount`'s matching lg/sm.
- **card** `size` (`sm`) — `--card-spacing` swaps from `--spacing(4)` to
  `--spacing(3)` (p-3/gap-3, title drops to `text-sm`).
- **native-select** `size` (`sm`) — `data-[size=sm]:h-8` (the `size=default`
  axis is already baked in unconditionally, see `docs/jsx-parity.md` `## select`).
- **switch** `size` (`sm`) — smaller switch/thumb metrics.

Cross-cutting, larger design-system change (not a single component):

- **border → ring hairline system** — nova recolors several surfaces'
  `border` to `ring-1 ring-foreground/10` (dropdown/context-menu/popover/
  hover-card content, dialog, alert-dialog, card). Ring sits outside the
  border box, so adopting it shifts every affected surface's inner content
  box by 1px/side — explicitly NOT adopted in the density retarget (color/
  box-model scope, not metric-token scope; see `docs/jsx-parity.md`
  `## nova density`). Worth its own scoped pass rather than folding into
  a future density-only change.

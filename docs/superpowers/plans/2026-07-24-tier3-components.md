# Tier 3 Components Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Port the eight Tier 3 components from `docs/component-roadmap.md` — toggle-group, slider, scroll-area, drawer, carousel, input-otp, select (custom listbox, with the native port renamed NativeSelect), sonner — with nova density/visuals baked in from the start.

**Architecture:** Every class string, ARIA/state contract, keyboard model, and adaptation decision comes from the two committed source maps: **`docs/superpowers/plans/2026-07-24-tier3-source-map-controls.md`** (toggle-group, slider, scroll-area, select) and **`docs/superpowers/plans/2026-07-24-tier3-source-map-wrapped.md`** (sonner, drawer, carousel, input-otp). Each task names its map section; the implementer follows that section's merged/recommended class strings and behavior contract exactly, plus the per-task Decisions block below (which resolves every point the maps flag "for the planner"). House process is unchanged from Tiers 1–2: TDD render pins, `go tool gsx generate`, site example pages, `docs/jsx-parity.md` ledger entries, one commit per task.

**Tech Stack:** gsx components (`ui/*.gsx`), vanilla-JS behavior modules (`ui/*.js`, barrel-imported via `ui/index.js`), Go pin tests, Tailwind v4 + tw-animate-css, tailwind-merge-go, `internal/registry` (derived deps/HasJS).

## Global Constraints

- **Nova policy for NEW ports** (differs from the 27-component retarget, deliberately): these components have no existing gsxui baseline to preserve, so they adopt nova's live visuals directly — metric tokens AND interaction-state tokens (e.g. slider `hover:ring-3`/`active:ring-3`) AND surface colors where a map recommends them (e.g. drawer `bg-popover`). Two standing exceptions, identical to the retarget: (1) nova's `border → ring-1 ring-foreground/10` hairline swaps are NOT adopted — keep `border` (select content, etc.); (2) nova's custom selector syntax is NOT adopted — keep house shapes (`data-[orientation=…]:` not `data-horizontal:`, `focus-visible:ring-[3px]` for the focus ring convention, `data-[state=…]` not `data-open:`). The maps' merged class strings already encode all of this — follow them byte-for-byte.
- Popover-family open/close animation is the house discrete-transition mechanism (`docs/jsx-parity.md` `## animations`): base closed-state classes + `open:` + `starting:open:` + `transition-[…,display,overlay] transition-discrete duration-150`, `data-state` stamped synchronously before `showPopover()`. Dialog-family (drawer) keeps the keyframe + JS-intercepted-close architecture from `ui/dialog.js` at nova's 200ms.
- JS module naming is binding: behavior JS must be `ui/<gsx-basename>.js` or `registry.HasJS` never finds it (`ui/sonner.js`, NOT `ui/toast.js`).
- Every new `.js` module is imported in `ui/index.js` (side-effect import, alphabetical with the existing ones).
- Every task: pin tests written first and failing, then implement, `go tool gsx generate && go test ./...`, site examples per the map's Demo inventory recommendation, `internal/registry/registry_test.go` want-list updated, `docs/jsx-parity.md` component entry (MECHANISM/ADAPT/GAP/WIN ledger style, matching existing entries), `make check` green, one commit.
- `scripts/repin.py` exists for pin splicing (first-class-per-line only; beware byte-identical literals across files — its docstring documents both limits).
- gsx gotchas that have bitten before: `//` inside markup renders as literal text (use Go doc comments); markup comment text containing `<tag>` breaks the parser; `{{ x := }}` scoping is function-wide; example packages can't have hyphens in dir names (`togglegroup`, `scrollarea`, `inputotp`, `nativeselect` — the `Register` key is the hyphenated component name and is what must match).

---

### Task 1: toggle-group

**Files:**
- Create: `ui/toggle-group.gsx`, `ui/toggle-group.js`, `ui/toggle-group_test.go`, `site/examples/togglegroup/{basic,single,spacing}.gsx`, `site/examples/togglegroup.go`
- Modify: `ui/index.js`, `ui/toggle.gsx` (export its `variantClass`/`sizeClass` helpers if not already reachable), `internal/registry/registry_test.go`, `docs/jsx-parity.md`

**Map section:** controls map `## toggle-group` — markup tree, merged root/item class strings, full behavior contract (roving tabindex, key map, Shift+Tab exit, single-vs-multiple ARIA split).

**Decisions (binding):**
- DROP the dead `data-[spacing=default]:data-[variant=outline]:shadow-xs` root selector — nova's own `.cn-toggle-group` drops it; ledger the upstream-dead-selector FINDING in `docs/jsx-parity.md` instead of porting dead weight.
- v1 is **horizontal-only**: stamp `data-orientation="horizontal"` on root and items, include only the horizontal corner selectors (`data-[orientation=horizontal]:data-[spacing=0]:first:rounded-l-lg` / `last:rounded-r-lg`), ledger vertical as GAP. toggle-group.js gates arrows accordingly (ArrowLeft/Up = prev, ArrowRight/Down = next, Home/PageUp = first, End/PageDown = last, loop at ends, modifier keys suppress).
- Group→item inheritance is **explicit params** (no context in gsx): `ToggleGroup(type, variant, size, spacing, …)` and `ToggleGroupItem(type, variant, size, spacing, pressed, value, …)` — `type` decides `role="radio"`+`aria-checked` (single) vs bare `aria-pressed` (multiple); both stamp `data-state="on"|"off"`. Item composes toggle.gsx's existing nova-retargeted variant/size class helpers, not a re-derivation.
- Roving tabindex is **JS-normalized at init**: server renders items naturally tabbable (no tabindex attr — graceful no-JS fallback where every item is a tab stop); `toggle-group.js` on init sets the pressed item (single) / first non-disabled item (multiple or none pressed) to `tabindex="0"` and the rest to `-1`, and maintains it on arrow moves and clicks. Shift+Tab exits the group in one press (port the onItemShiftTab idiom).
- Single-type click behavior: clicking a new item replaces the value (exactly one `data-state="on"`); clicking the pressed item toggles it off (Radix single allows empty unless a caller opts otherwise — port the replace-on-activate mechanic from the map).

- [ ] **Step 1:** Write failing pin tests: root pinned render (role/data attrs/class), item single-type (`role="radio"` + `aria-checked`), item multiple-type (`aria-pressed`, no role), disabled cascade, caller class merge.
- [ ] **Step 2:** Implement `ui/toggle-group.gsx` + `ui/toggle-group.js` per map + decisions. `go tool gsx generate && go test ./...` → PASS.
- [ ] **Step 3:** Site examples per map Demo inventory: `basic` (type=multiple, outline, 3 icon items), `single` (type=single default variant), `spacing` (spacing=2, size=sm, outline). Register in `site/examples/togglegroup.go` key `"toggle-group"`.
- [ ] **Step 4:** Registry want-list + `docs/jsx-parity.md` `## toggle-group` entry (dead-selector FINDING, vertical GAP, explicit-params ADAPT, JS-roving ADAPT).
- [ ] **Step 5:** `make check`; commit `feat(ui): add toggle-group component`.

### Task 2: slider

**Files:**
- Create: `ui/slider.gsx`, `ui/slider.js`, `ui/slider_test.go`, `site/examples/slider/basic.gsx`, `site/examples/slider.go`
- Modify: `ui/index.js`, `assets/gsxui.css` (or the established shared-CSS location — follow the `::details-content` precedent) for the range pseudo-element block, `internal/registry/registry_test.go`, `docs/jsx-parity.md`

**Map section:** controls map `## slider` — native `<input type=range>` collapse, track/thumb metrics, cross-browser pseudo-element CSS (including the WebKit `margin-top` centering fix), focus-ring selector pairs.

**Decisions (binding):**
- Filled-range visual is **option 2 — CSS custom property + gradient track, synced by `ui/slider.js`** (not `accent-color`; guaranteed cross-browser two-tone fill, no "simple heuristic"). Server computes the initial `--fill` percentage from `value`/`min`/`max` at render (correct first paint, zero JS needed until first drag); `slider.js` is a delegated `input` listener updating `--fill` inline.
- Thumb adopts nova `size-3` + `hover:ring-3 focus-visible:ring-3 active:ring-3` (new-port nova policy). The `after:-inset-2` hit-target compensation is categorically unreachable on thumb pseudo-elements — ledger the hit-target regression explicitly as GAP.
- Track `h-1` (nova), `bg-muted rounded-full`; fill `bg-primary` via the gradient; thumb `border-primary bg-white` (nova's `border-ring` recolor NOT adopted — color scope).
- Single scalar value only: multi-thumb/range, vertical, `inverted`, `minStepsBetweenThumbs` are GAPs to ledger (native input has no analog). Keyboard/ARIA come free from the native input — state that as WIN in the ledger.

- [ ] **Step 1:** Failing pins: rendered input with `min`/`max`/`step`/`value` attrs, class string, inline `--fill` style computed from value, disabled state.
- [ ] **Step 2:** Implement `ui/slider.gsx` (+ shared CSS block) + `ui/slider.js`. Generate + test → PASS.
- [ ] **Step 3:** Example `basic` (value 50, max 100, step 1, `w-[60%]` wrapper — the one in-scope shadcn demo). Register key `"slider"`.
- [ ] **Step 4:** Registry want-list + jsx-parity `## slider` (gradient ADAPT, hit-target GAP, multi-thumb GAP, native-keyboard WIN, WebKit margin-top MECHANISM).
- [ ] **Step 5:** `make check`; commit `feat(ui): add slider component`.

### Task 3: scroll-area

**Files:**
- Create: `ui/scroll-area.gsx`, `ui/scroll-area_test.go`, `site/examples/scrollarea/{basic,horizontal}.gsx`, `site/examples/scrollarea.go`
- Modify: shared CSS (WebKit scrollbar pseudo block), `internal/registry/registry_test.go`, `docs/jsx-parity.md`

**Map section:** controls map `## scroll-area` — collapsed single-div design, standard `scrollbar-width`/`scrollbar-color` + layered `::-webkit-scrollbar-*` block, enumerated CANNOT list. Zero nova delta (confirmed).

**Decisions (binding):**
- ONE `ui.ScrollArea` component, one div (`overflow-auto` / `overflow-x-auto` for `orientation="horizontal"`), `rounded-[inherit]`, viewport focus ring per map. No `ScrollBar`/`Thumb`/`Corner` parts (GAP), no `type` visibility timing (GAP — OS policy wins), per the map's recommended v1 shape. No JS module (`HasJS` false).

- [ ] **Step 1:** Failing pins: vertical render, horizontal render, caller class merge.
- [ ] **Step 2:** Implement component + shared-CSS scrollbar block. Generate + test → PASS.
- [ ] **Step 3:** Examples `basic` (h-72 w-48 bordered tag list with separators) and `horizontal` (w-96 image-row equivalent). Register key `"scroll-area"`.
- [ ] **Step 4:** Registry want-list + jsx-parity `## scroll-area` (CSS-first ADAPT, dual-CSS-surface MECHANISM, Corner/type GAPs).
- [ ] **Step 5:** `make check`; commit `feat(ui): add scroll-area component`.

### Task 4: drawer

**Files:**
- Create: `ui/drawer.gsx`, `ui/drawer_test.go`, `site/examples/drawer/basic.gsx`, `site/examples/drawer.go`
- Modify: `internal/registry/registry_test.go`, `docs/jsx-parity.md`, `docs/component-roadmap.md` (patterns-phase note for the responsive drawer-dialog demo)

**Map section:** wrapped map `## drawer` — full part list, per-direction class strings (built on sheet's solved UA fights), nova deltas table, handle-bar rule.

**Decisions (binding):**
- Composes `ui.Dialog` machinery exactly like Sheet (`data-gsxui-dialog-content`, `HasJS("drawer") == false`, zero new dialog.js code). `DrawerContent` renders its own `<dialog>`, not a composition of Dialog/SheetContent.
- Go param is `direction` (vaul's vocabulary, default `"bottom"`), stamped as `data-side` (reusing sheet's internal attribute).
- Adopt nova: `rounded-*-xl` all four directions (incl. the new left/right free-edge rounding), `bg-popover text-popover-foreground` (map's flagged exception — confirmed), `h-1` handle, `gap-0.5` header at all breakpoints, `font-medium` title. Backdrop identical to sheet's (`bg-black/10` + `backdrop-blur-xs`, 200ms both ways).
- Handle bar rendered for bottom only, kept as decoration; drag-to-dismiss, snap points, background scaling all GAPs (roadmap-mandated cut).
- The per-direction class strings are transcribed-not-verified — the CONTROLLER (not the implementer subagent) runs a live-browser verification of all four directions on the dev server before this task is marked complete, mirroring how sheet's six ADAPTs were found.

- [ ] **Step 1:** Failing pins: content render per direction (all four class strings), header/footer/title/description, trigger/close buttons, handle-bar visibility rule.
- [ ] **Step 2:** Implement `ui/drawer.gsx`. Generate + test → PASS.
- [ ] **Step 3:** Example `basic`: adapted goal-counter demo — +/- buttons and a static decorative bar row of plain divs in place of the recharts sparkline (map's recommendation). Register key `"drawer"`.
- [ ] **Step 4:** Registry want-list (deps `["dialog"]`) + jsx-parity `## drawer` + roadmap patterns-phase note for `drawer-dialog`.
- [ ] **Step 5:** `make check`; commit `feat(ui): add drawer component`. Controller browser pass follows before the task closes.

### Task 5: carousel

**Files:**
- Create: `ui/carousel.gsx`, `ui/carousel.js`, `ui/carousel_test.go`, `site/examples/carousel/{basic,sizes,api}.gsx`, `site/examples/carousel.go`
- Modify: `ui/index.js`, `internal/registry/registry_test.go`, `docs/jsx-parity.md`

**Map section:** wrapped map `## carousel` — two-div content structure, scroll-snap substitution, prev/next disabled logic, `ui/carousel.js` contract, API-surface reduction.

**Decisions (binding):**
- Native scroll-snap replaces embla per roadmap: viewport gets `overflow-x-auto snap-x snap-mandatory` + scrollbar hiding, items get `snap-start`; `loop` is a ledgered GAP (no demo needs it); touch/trackpad momentum is a WIN.
- Prev/next scroll by one measured item width; disabled state from scrollLeft vs 0 / scrollWidth−clientWidth with epsilon, recomputed on rAF-throttled scroll + ResizeObserver, set as native `disabled` on the composed `ui.Button` (variant outline, size icon, `size-8 rounded-full` absolute positioning per source).
- JS surface: delegated clicks, `gsxui:carousel-select` CustomEvent `{index, count}`, `data-current-index` stamp, ArrowLeft/Right keydown within the root, `el.gsxuiCarousel = {scrollTo,next,prev}` imperative handle, optional `data-gsxui-carousel-autoplay="<ms>"` (pause on pointerenter/focusin, resume on leave/out).
- `orientation="vertical"` supported (snap-y + vertical prev/next positioning per source); spacing stays caller-controlled via class merge (`-ml-*`/`pl-*`), no spacing prop.

- [ ] **Step 1:** Failing pins: root (role/aria-roledescription/data attrs), content two-div structure horizontal + vertical, item, prev/next buttons (composed Button classes, sr-only labels).
- [ ] **Step 2:** Implement `ui/carousel.gsx` + `ui/carousel.js`. Generate + test → PASS.
- [ ] **Step 3:** Examples: `basic` (single per view), `sizes` (`md:basis-1/2 lg:basis-1/3`, align start), `api` (Slide X of Y indicator driven by the CustomEvent via a small inline example script). Register key `"carousel"`.
- [ ] **Step 4:** Registry want-list (deps `["button"]`, HasJS true) + jsx-parity `## carousel` (scroll-snap ADAPT, loop GAP, momentum WIN, API-reduction ledger).
- [ ] **Step 5:** `make check`; commit `feat(ui): add carousel component`.

### Task 6: input-otp

**Files:**
- Create: `ui/input-otp.gsx`, `ui/input-otp.js`, `ui/input-otp_test.go`, `site/examples/inputotp/{basic,pattern,separator}.gsx`, `site/examples/inputotp.go`
- Modify: `ui/index.js`, `internal/registry/registry_test.go`, `docs/jsx-parity.md`

**Map section:** wrapped map `## input-otp` — hidden-single-input architecture (the whole mechanism), slot class string with nova deltas table (size-8, rounded-l/r-lg, no shadow-xs, group invalid ring), fake-caret markup, JS contract.

**Decisions (binding):**
- Architecture is the map's: ONE real `<input>` (opacity-0, z-10, focusable, `autocomplete="one-time-code"`), slots purely presentational, populated by JS from `value` + selection on `input`/`selectionchange`/`focus`/`blur`; click-to-position via `setSelectionRange`. Slots render empty server-side (ledger the first-paint divergence).
- Index binding is **option (b) — DOM-order stamping**, no `index` param on `InputOTPSlot` (deliberate API departure from shadcn, command.js source-order-stamp precedent; ledger it).
- Pattern attr `data-gsxui-input-otp-pattern` takes an **unanchored per-character class** (e.g. `[0-9]`), documented in the component comment and the pattern example; do NOT copy shadcn's anchored REGEXP constants.
- Adopt nova group-level invalid ring (`has-aria-invalid:ring-3` block) and all slot nova deltas; keep `data-[active=true]:z-10`.
- Does not compose `ui.Input`; deps `["icon"]` (Minus separator); HasJS true.

- [ ] **Step 1:** Failing pins: container + hidden input (maxlength/pattern/name fallthrough), group, slot (empty, nova classes), separator, caret markup absent server-side.
- [ ] **Step 2:** Implement `ui/input-otp.gsx` + `ui/input-otp.js`. Generate + test → PASS.
- [ ] **Step 3:** Examples: `basic` (two 3-slot groups + separator), `pattern` (digits-only), `separator` (three 2-slot groups — the index-stamping stress case). Register key `"input-otp"`.
- [ ] **Step 4:** Registry want-list + jsx-parity `## input-otp` (single-input MECHANISM, index-stamping ADAPT, pattern-shape decision, first-paint GAP).
- [ ] **Step 5:** `make check`; commit `feat(ui): add input-otp component`.

### Task 7: NativeSelect rename

**Files:**
- Create: `ui/native-select.gsx` (moved), `ui/native-select_test.go` (moved)
- Delete: `ui/select.gsx`, `ui/select.x.go`, `ui/select_test.go` (contents move; `git mv` where possible)
- Modify: `site/examples/selectbox/` → `site/examples/nativeselect/` (package + Register key `"native-select"`), every `ui.Select`/`ui.SelectOption`/`ui.SelectGroup` call site (`grep -rn "ui\.Select" site/`), `internal/registry/registry_test.go`, `docs/jsx-parity.md` (`## select` entry retitled `## native-select` with a rename note), `docs/component-roadmap.md` naming-delta line

**Interfaces:**
- Produces: `ui.NativeSelect`, `ui.NativeSelectOption`, `ui.NativeSelectGroup` — signatures unchanged except the names. Frees the `Select*` identifier space for Task 8.

- [ ] **Step 1:** `git mv` the three files; rename components `Select→NativeSelect`, `SelectOption→NativeSelectOption`, `SelectGroup→NativeSelectGroup`; update all call sites and tests (pure rename, zero class-string changes).
- [ ] **Step 2:** `go tool gsx generate && go test ./...` → PASS (registry list now says `native-select`; want-list updated).
- [ ] **Step 3:** Docs: retitle the jsx-parity entry, note the rename + date; roadmap naming line updated.
- [ ] **Step 4:** `make check`; commit `refactor(ui): rename Select → NativeSelect ahead of the custom listbox`.

### Task 8: select (custom listbox)

**Files:**
- Create: `ui/select.gsx`, `ui/select.js`, `ui/select_test.go`, `site/examples/selectbox/{basic,scrollable}.gsx`, `site/examples/selectbox.go`
- Modify: `ui/index.js`, `internal/registry/registry_test.go`, `docs/jsx-parity.md`

**Map section:** controls map `## select (custom listbox)` — full tree with runtime ARIA, merged nova class strings for trigger/content/label/item/indicator/separator, the complete behavior contract (OPEN_KEYS, closed-trigger typeahead with the Space rule, `aria-selected = isSelected && isFocused`, pointer-type-aware activation, 1000ms prefix typeahead with repeat-char cycling, hidden native `<select>` form bridge), and the machinery-reuse table.

**Decisions (binding, all per the map's recommendations):**
- Parts: `Select` (root div `data-gsxui-select`), `SelectTrigger` (button, `role="combobox"`, `data-size` default|sm), `SelectValue`, `SelectContent` (`popover="auto"`, `role="listbox"`, discrete transitions duration-150, `border` kept — no ring swap), `SelectGroup`/`SelectLabel` (real parts — div `role="group"` + styled label; the name no longer collides after Task 7), `SelectItem` (`role="option"`, `data-value`, check indicator via `icon.Check` size-4), `SelectSeparator`. NO scroll up/down buttons (GAP — popper-equivalent anchoring, native overflow scroll).
- Focus model: **real DOM focus** (dropdown.js model — hover-is-focus, arrow `.focus()` walk); `aria-selected` recomputed on every focus change (value-match AND focus-match); `data-state="checked|unchecked"` tracks the value alone and drives the checkmark.
- Reuse dropdown.js idioms as-is: popover anchoring, `closest()` proximity wiring, sync `data-state` stamp before `showPopover()`, wasOpen pointerdown/click guard, pointerover/pointerout focus handlers. New in select.js: value model (one checked item, trigger text update), bespoke prefix typeahead (buffer + 1000ms reset + startsWith + same-char cycling; works on the CLOSED trigger too, Space swallowed only mid-search), hidden `<select>` sync (`value` assignment + bubbling `change` event).
- Form bridge: server renders a real `<select aria-hidden="true" tabindex="-1" class="sr-only">` sibling with one `<option>` per item value, `name`/`required`/`disabled`/`form` forwarded — mirroring NativeSelect's option-authoring shape.
- Trigger: nova metrics (`h-8` default / `h-7` sm + sm radius override, `gap-1.5`, `pr-2 pl-2.5`, `rounded-lg`, no shadow-xs); `focus-visible:ring-[3px]` house syntax.
- Deps: `["icon"]` (+ whatever the registry derives); HasJS true.

- [ ] **Step 1:** Failing pins: trigger (both sizes, placeholder state), content, group+label, item (checked and unchecked ARIA/state), separator, hidden native select rendering with options, caller class merges.
- [ ] **Step 2:** Implement `ui/select.gsx` + `ui/select.js`. Generate + test → PASS.
- [ ] **Step 3:** Examples: `basic` (group+label+5 items+placeholder, w-[180px]) and `scrollable` (5 groups / 27 timezone items). Register key `"select"`.
- [ ] **Step 4:** Registry want-list + jsx-parity `## select` (aria-selected FINDING, typeahead MECHANISM, hidden-select MECHANISM, scroll-button GAP, focus-model note).
- [ ] **Step 5:** `make check`; commit `feat(ui): add select — custom listbox on the popover machinery`.

### Task 9: sonner

**Files:**
- Create: `ui/sonner.gsx`, `ui/sonner.js`, `ui/sonner_test.go`, `site/examples/sonner/{basic,types}.gsx`, `site/examples/sonner.go`
- Modify: `ui/index.js` (side-effect import AND `export { toast }`), `site/pages/layout.gsx` (mount `<ui.Toaster/>` once), `internal/registry/registry_test.go`, `docs/jsx-parity.md`

**Map section:** wrapped map `## sonner` — toaster/toast markup (synthesized class strings), stacking model, timers, morph-in-place promise behavior, API shape, naming constraint.

**Decisions (binding):**
- `ui.Toaster`: the only server-rendered part (section + `<ol data-gsxui-toaster>`, bottom-right, mounted once in the site layout). Toast `<li>` DOM is constructed entirely by `ui/sonner.js` (the codebase's first client-constructed-markup module — the map's class strings for the card are the spec). Type-tinted icons (emerald/sky/amber/destructive) — confirmed, matches the live nova site; icon SVG paths hand-copied from `ui/icon/icon_data.go` into the JS with a provenance comment (maintenance seam, ledger it).
- Behavior: max 3 visible collapsed-stack (scale/offset via inline styles from a plain toast-record array — do NOT reproduce sonner's CSS-var machine), hover-to-expand with 14px gaps, 4000ms default duration, pause-on-hover with remaining-time resume, enter/exit via the house discrete-transition architecture, promise toasts morph the SAME node in place (no re-animation). No swipe gestures, bottom-right only (positions GAP).
- API: `toast(msg, opts)` + `.success/.info/.warning/.error/.loading/.promise/.dismiss` exported through the barrel (`import { toast } from "gsxui"` — first public imperative API precedent, ledger it), plus the declarative `data-gsxui-toast` / `-description` / `-type` / `-action` click-delegated trigger for zero-JS demo pages (action button emits `gsxui:toast-action` and dismisses).
- Registry: component name `sonner`, JS must be `ui/sonner.js` (HasJS contract).

- [ ] **Step 1:** Failing pins: Toaster render (region + ol, classes, aria). JS-side behavior is exercised by the examples + browser pass (no DOM-emulation test framework exists in this repo — match command.js's precedent of Go-pinning only the server-rendered surface).
- [ ] **Step 2:** Implement `ui/sonner.gsx` + `ui/sonner.js`. Generate + test → PASS.
- [ ] **Step 3:** Examples: `basic` (declarative trigger with description+action) and `types` (six buttons: default/success/info/warning/error/promise — promise via a tiny inline script using the public `toast.promise`). Register key `"sonner"`. Mount Toaster in `site/pages/layout.gsx`.
- [ ] **Step 4:** Registry want-list + jsx-parity `## sonner` (client-constructed-DOM MECHANISM, icon-copy seam, positions GAP, barrel-export precedent).
- [ ] **Step 5:** `make check`; commit `feat(ui): add sonner toasts — own module, no dependency`.

### Task 10: Sweep, ledger roll-up, deploy

**Files:**
- Modify: `docs/component-roadmap.md` (move the eight Tier 3 rows to shipped status / annotate; record deferred sub-features: toggle-group vertical, slider range/vertical, carousel loop, drawer drag+snap+scale, sonner positions, select scroll buttons, input-otp controlled-demo pattern), `docs/jsx-parity.md` (verify all eight entries landed), `.superpowers/sdd/progress.md`

- [ ] **Step 1:** Full-suite verify: `go tool gsx generate && go test ./... && make check` → green.
- [ ] **Step 2:** Roadmap + ledger roll-up as above.
- [ ] **Step 3:** Commit `docs: Tier 3 roll-up`, `git push origin main`, confirm `gh run list --repo gsxhq/gsxui --limit 1` reaches success.
- [ ] **Step 4:** Controller browser pass on ui.gsxhq.dev (or dev server): drawer all four directions, slider fill/thumb cross-check, scroll-area thumb fidelity, select keyboard walk + typeahead + form value, sonner stack/expand/promise-morph, input-otp typing/paste/click-to-position, carousel snap + disabled states, toggle-group roving focus. User performs the final visual pass against ui.shadcn.com.

## Self-Review

- Spec coverage: all 8 Tier 3 roadmap rows have a task (1 toggle-group, 2 slider, 3 scroll-area, 4 drawer, 5 carousel, 6 input-otp, 7+8 select incl. the roadmap's rename decision, 9 sonner) + roll-up (10) ✓. Every "flag for planner" item in both maps is resolved in a Decisions block: toggle-group dead selector (drop), vertical (GAP), roving-tabindex placement (JS init); slider fill (JS gradient), thumb size (nova size-3, ledgered hit-target); drawer naming (`direction`→`data-side`), bg-popover (adopt), verification pass (controller); sonner file naming (sonner.js), icon tinting (adopt), API shape; input-otp index binding (DOM-order), pattern shape (unanchored); select group naming (freed by Task 7), focus model, typeahead, scroll buttons ✓.
- Placeholders: none — exact values live in the two committed maps each task names; steps carry exact commands, file paths, and Register keys.
- Type consistency: Task 7 produces `NativeSelect*` before Task 8 consumes the freed `Select*` names; toggle-group consumes toggle.gsx helpers (exported in Task 1); no other cross-task Go surface.

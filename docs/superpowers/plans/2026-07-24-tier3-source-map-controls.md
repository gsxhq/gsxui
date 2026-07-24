# Tier 3 source map — toggle-group, slider, scroll-area, select (custom listbox)

Source-analysis map for porting these four `docs/component-roadmap.md` Tier 3
rows. Structural/markup reference: `shadcn-ui/apps/v4/registry/new-york-v4/ui/{toggle-group,toggle,slider,scroll-area,select}.tsx`.
Runtime-behavior reference: the actual Radix packages on disk at
`/Users/jackieli/work/his-project/node_modules/.bun/@radix-ui+react-{toggle-group,toggle,roving-focus,select,scroll-area}@*/node_modules/@radix-ui/react-*/dist/index.mjs`
(read directly — every ARIA attribute, keyboard key, and state machine claim
below is traced to that source, not recalled from docs). Density reference:
`shadcn-ui/apps/v4/registry/styles/style-nova.css` (`.cn-toggle-group*`,
`.cn-slider*`, `.cn-scroll-area*`, `.cn-select*`). House conventions:
`docs/jsx-parity.md` (`## nova density`, `## animations`, `## toggle`,
`## dropdown`, `## select`), `docs/component-roadmap.md` Tier 3 table.

## Legend

- `token→token (nova)` — nova changes this exact token's numeric value; both
  ends shown, nova's value is what gets ported.
- `(nova, new)` — nova adds a token new-york-v4 doesn't have.
- `(nova, drop)` — nova removes a token new-york-v4 has.
- `(NOT ADOPTED — color)` / `(NOT ADOPTED — border→ring)` / `(NOT ADOPTED —
  display-model)` — nova changed this token but the change falls in one of
  the categories `docs/jsx-parity.md` `## nova density` explicitly excludes
  (colors, the border-to-ring-1 box-model swap, display-model rewrites);
  new-york-v4's token is kept verbatim, nova's alternative is named for the
  record only.
- `(SHADOW-PRESENCE, drop)` — nova drops a shadow outright (no replacement);
  per house convention this IS adopted (removals only, never additions).
- `(nova, flagged — outside stated scope)` — nova's delta doesn't cleanly
  fall inside "heights/paddings/gaps/text/radii/svg sizes" (ring width,
  z-index, background-color additions, animation-duration) but is plausibly
  size-adjacent; presented for the planner to decide, not pre-applied.
- "current" gsxui reference components (`ui/toggle.gsx`, `ui/dropdown.gsx`,
  `ui/select.gsx`) are quoted where they already encode a nova-retargeted
  token, since Tier 3 ports should match siblings, not rediscover the value.

---

## toggle-group

### Markup structure

new-york-v4 (`toggle-group.tsx` + `toggle.tsx`'s `toggleVariants`) plus the
exact runtime attributes Radix's `@radix-ui/react-toggle-group` /
`react-roving-focus` / `react-toggle` stamp (traced in their `dist/index.mjs`,
not guessed from the `.tsx` alone — the `.tsx` never renders `role`/
`data-orientation`/`aria-checked` explicitly; Radix's primitives add them at
runtime):

```
<div data-slot="toggle-group" data-variant data-size data-spacing
     role="radiogroup"|"toolbar"      (type=single | type=multiple)
     data-orientation="horizontal"|"vertical"   (stamped by RovingFocusGroup's
       own Primitive.div, which IS this node via asChild — not a separate
       wrapper)
     style="--gap: {spacing}"
     class="...">
  <button data-slot="toggle-group-item" data-variant data-size data-spacing
          type="button" data-state="on"|"off"
          role="radio" aria-checked="true|false"           (type=single only)
          | aria-pressed="true|false"                       (type=multiple only)
          data-orientation="horizontal"|"vertical"
          tabindex="0"|"-1"                (roving — exactly one item is 0)
          disabled?
          class="...">
    {children}
  </button>
  ...
</div>
```

Class strings, new-york-v4 base → nova overlay:

**Root** (`toggleVariants` doesn't apply here; this is `toggle-group.tsx`'s
own root class):
```
group/toggle-group flex w-fit items-center gap-[--spacing(var(--gap))]
rounded-md→rounded-lg (nova)
data-[spacing=0]:rounded-[min(var(--radius-md),10px)] data-[size=sm]:... (nova, new — see FINDING below on the "default" typo this replaces)
data-[spacing=default]:data-[variant=outline]:shadow-xs   ← see FINDING, dead selector in the SOURCE itself, port verbatim as dead weight (house convention: don't silently fix upstream)
```
Nova's own `.cn-toggle-group` (`style-nova.css` L1371-1373):
`rounded-lg data-[size=sm]:rounded-[min(var(--radius-md),10px)]` — nova adds
the `data-[size=sm]` radius override, drops the (dead) shadow selector
entirely rather than fixing its typo. Recommended merged root class:
```
group/toggle-group flex w-fit items-center gap-[--spacing(var(--gap))] rounded-lg data-[size=sm]:rounded-[min(var(--radius-md),10px)] data-[spacing=default]:data-[variant=outline]:shadow-xs
```
(keep the dead shadow selector per house "port dead weight, ledger it"
convention used throughout `docs/jsx-parity.md`, e.g. dialog's dropped
`data-[state=open]:` pair — OR drop it and cite nova's own precedent of
dropping it; flagged for planner decision, not resolved here).

**Item** — starts from `toggleVariants(variant, size)` (gsxui's OWN already
nova-retargeted version, `ui/toggle.gsx`, not new-york-v4's raw
`toggleVariants` — reuse `ui.Toggle`'s existing `variantClass`/`sizeClass`
switches rather than re-deriving them):
```
base (ui/toggle.gsx L67): inline-flex items-center justify-center gap-1 rounded-lg text-sm font-medium whitespace-nowrap transition-[color,box-shadow] outline-none hover:bg-muted hover:text-muted-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 data-[state=on]:bg-accent data-[state=on]:text-accent-foreground dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4
variant outline (ui/toggle.gsx L68): border border-input bg-transparent hover:bg-accent hover:text-accent-foreground   (shadow-xs already dropped, SHADOW-PRESENCE)
size default/sm/lg: ui/toggle.gsx L69 verbatim
```
plus toggle-group.tsx's item-specific additions, nova-overlaid:
```
new-york-v4: w-auto min-w-0 shrink-0 px-3 focus:z-10 focus-visible:z-10
             data-[spacing=0]:rounded-none data-[spacing=0]:shadow-none
             data-[spacing=0]:first:rounded-l-md→rounded-l-lg (nova)
             data-[spacing=0]:last:rounded-r-md→rounded-r-lg (nova)
             data-[spacing=0]:data-[variant=outline]:border-l-0
             data-[spacing=0]:data-[variant=outline]:first:border-l
nova adds (`.cn-toggle-group-item`, L1375-1377, structural not just metric):
  group-data-[spacing=0]/toggle-group:rounded-none                              (same as new-york-v4's data-[spacing=0]:rounded-none, re-scoped through the /toggle-group named group instead of the item's own data attr — SAME semantics, gsxui should keep new-york-v4's simpler item-local `data-[spacing=0]:` form, not switch selector shape)
  group-data-[spacing=0]/toggle-group:px-2                                       (nova, new — spacing=0 items get a smaller inner px than the ungrouped px-2.5 default)
  group-data-[spacing=0]/toggle-group:has-data-[icon=inline-end]:pr-1.5          (NOT ADOPTED — requires data-icon stamps, out of scope, same call as toggle's own icon-padding ADAPT)
  group-data-[spacing=0]/toggle-group:has-data-[icon=inline-start]:pl-1.5        (NOT ADOPTED, same reason)
  group-data-horizontal/toggle-group:data-[spacing=0]:first:rounded-l-lg         (nova, NEW — vertical-orientation support; requires data-orientation on root, confirmed present via Radix's RovingFocusGroup — see markup tree above)
  group-data-vertical/toggle-group:data-[spacing=0]:first:rounded-t-lg           (nova, NEW — vertical orientation)
  group-data-horizontal/toggle-group:data-[spacing=0]:last:rounded-r-lg
  group-data-vertical/toggle-group:data-[spacing=0]:last:rounded-b-lg
```
Recommended merged item class (keeping new-york-v4's item-local
`data-[spacing=0]:` selector shape rather than switching to nova's
`group-data-[spacing=0]/toggle-group:` — same match, this codebase already
uses the item-local shape and switching serves no purpose without also
porting nova's Tailwind config for the boolean `data-horizontal`/`data-vertical`
variants):
```
w-auto min-w-0 shrink-0 px-3 focus:z-10 focus-visible:z-10
data-[spacing=0]:rounded-none data-[spacing=0]:shadow-none
data-[spacing=0]:data-[variant=outline]:border-l-0 data-[spacing=0]:data-[variant=outline]:first:border-l
data-[orientation=horizontal]:data-[spacing=0]:first:rounded-l-lg data-[orientation=horizontal]:data-[spacing=0]:last:rounded-r-lg
data-[orientation=vertical]:data-[spacing=0]:first:rounded-t-lg data-[orientation=vertical]:data-[spacing=0]:last:rounded-b-lg
```
(the corner-rounding selectors must gate on `data-orientation` if vertical
ships at all — new-york-v4's own `first:rounded-l-md last:rounded-r-md` is
silently WRONG for a vertical group since it'd round left/right corners on a
column layout; this is nova adding real missing functionality, not just a
metric bump — flag prominently for the planner, it's the one place in this
component where "apply nova's numbers" also means "port nova's selector
shape," not just its values).

**FINDING (shadcn upstream, dead selector in `toggle-group.tsx` itself)**:
`ToggleGroup`'s root class carries `data-[spacing=default]:data-[variant=outline]:shadow-xs`,
but `data-spacing` is stamped as the literal prop value (`data-spacing={spacing}`,
default `0`, a JS **number**) — never the string `"default"`. `data-[spacing=default]`
can therefore never match anything the component itself ever renders; it is
dead CSS in the shadcn source, not just a port-time casualty. (Circumstantial
evidence for what was probably intended: the ITEM's own `data-[spacing=0]:shadow-none`
suppresses the per-item outline shadow when items are joined into one pill —
the root's shadow was very likely meant to substitute a single shadow around
the whole pill in that same `spacing=0` case, i.e. the selector should read
`data-[spacing=0]`, not `data-[spacing=default]`.) Ledger this and port
verbatim as dead weight per this codebase's established "don't silently fix
upstream, ledger the dead selector" convention (`## dialog`'s dropped
`data-[state=open]:` pair is the precedent) — do not silently "fix" it to
`data-[spacing=0]` without calling it out as a deviation from the binding
source.

### Behavior contract

Traced directly from `@radix-ui/react-toggle-group@1.1.13`,
`@radix-ui/react-roving-focus@1.1.13`, `@radix-ui/react-toggle@1.1.12`
`dist/index.mjs` (not the `.tsx`, which only sets `data-slot`/classes):

- **Root role**: `type="single"` → `role="radiogroup"`. `type="multiple"` →
  `role="toolbar"`. (`ToggleGroup` throws if `type` is omitted — it is a
  required prop upstream, not defaulted.)
- **Item semantics differ completely by type** — this is the single most
  important behavior fact for the port:
  - `type="single"`: item is `role="radio" aria-checked={pressed}` — **`aria-pressed`
    is explicitly unset** (Radix's own `ToggleGroupItemImpl` passes
    `"aria-pressed": void 0` to override `Toggle`'s own default). Exactly one
    item in the group can be checked at a time (enforced by the value
    context, not by ARIA alone).
  - `type="multiple"`: item is a plain button, `aria-pressed={pressed}`, no
    `role` override (default implicit button role). Any subset of items can
    be pressed independently.
  - **Both types** stamp `data-state="on"|"off"` on the item (from the shared
    underlying `Toggle` primitive) — this is what the shared CSS selector
    `data-[state=on]:bg-accent` keys off, so the single class string works
    unmodified for both types; only the ARIA attribute pair differs.
- **Roving tabindex** (`RovingFocusGroupImpl`/`RovingFocusGroupItem`,
  `rovingFocus` defaults `true`, `loop` defaults `true`): exactly one item
  has `tabindex="0"` (the "current tab stop"), every other item
  `tabindex="-1"`. Entry-focus candidate priority when Tab first lands in the
  group: the item with `active===true` (i.e. the pressed one, if any) wins
  first, else the last-focused item, else the first item.
- **Keyboard model** (`getFocusIntent`, exact key map):
  ```
  ArrowLeft, ArrowUp   → prev
  ArrowRight, ArrowDown → next
  Home, PageUp         → first
  End, PageDown        → last
  ```
  **Orientation gating is inverted from what you'd guess**: if
  `orientation==="vertical"`, ArrowLeft/ArrowRight are ignored (return
  `undefined`, no-op); if `orientation==="horizontal"`, ArrowUp/ArrowDown are
  ignored. **If `orientation` is left unset** (new-york-v4's `ToggleGroup`
  never passes one) **neither condition is true, so all four arrow keys work
  regardless of the CSS layout direction** — this is the actual shipped
  behavior for every new-york-v4 toggle-group demo today (none of them pass
  `orientation`). `loop=true` (default) wraps at the ends; a caller could
  pass `loop={false}` but no example does.
  RTL: `ArrowLeft`/`ArrowRight` are swapped by `getDirectionAwareKey` when
  `dir==="rtl"`.
  A modifier key held (Ctrl/Alt/Meta/Shift) suppresses the focus-move
  entirely (`if (event.metaKey || ...) return`) — arrow+modifier does nothing,
  not even scroll.
  Shift+Tab specifically is intercepted (`onItemShiftTab`) to mark
  "tabbing back out," which changes the group container's own `tabIndex` to
  `-1` for that one blur/focus cycle so Shift+Tab exits the whole group in
  one press rather than landing back on the roving item.
- **Activation**: items are real `<button>`s, so Enter/Space activate via
  native button `click` semantics — Radix adds no `keydown` handler for
  activation itself, only for arrow/Home/End navigation. A click (or
  Enter/Space) on `ToggleGroupItemImpl` flips `pressed` via the underlying
  `Toggle`'s own `onClick`, which calls `valueContext.onItemActivate(value)`
  or `onItemDeactivate(value)`.
- **Single-type exclusivity mechanic**: `onItemActivate` in
  `ToggleGroupImplSingle` is literally `setValue` (replace); there is no
  "onItemDeactivate wins" race — clicking a different item in `type="single"`
  simply sets a new single value, which is why exactly one item shows
  `data-state="on"` at a time without any group-level "uncheck the others"
  loop needed in the port.
- **Disabled**: `context.disabled` (from `ToggleGroup`'s own `disabled` prop)
  ORs with each item's own `disabled` prop (`context.disabled || props.disabled`)
  — group-level disable cascades to every item unless the group itself isn't
  disabled and only a specific item is.

### gsxui adaptation notes

- **Not a native-input port** (same reasoning as `## toggle`'s own GAP
  entry, restated because it's the natural first question for a "grouped
  radio-like" control): native `<input type="radio">` was considered and
  rejected for the same reason plain `Toggle` rejected it — items render
  arbitrary child content (icon, text, or both; see every
  `toggle-group-*.tsx` example) as the pressable surface itself, not a label
  sibling to a hidden input. `type="multiple"` additionally has no radio
  analog at all (independent booleans, not mutual exclusion). Both types stay
  `<button>` + JS, same shape as `## toggle`.
- **Reuse `ui.Toggle`'s `variantClass`/`sizeClass`, not a re-derivation**:
  `ui/toggle.gsx` already carries the nova-retargeted variant/size switches
  token-for-token; `ToggleGroupItem` should call into the same functions
  (exported, or duplicated verbatim with a comment pointing at the toggle.gsx
  origin — matching how `## button-group`'s `ButtonGroupSeparator` composes
  `ui.Separator` directly rather than reimplementing it) plus the
  toggle-group-only additions above.
- **`type` must be a real, required param** — mirroring Radix's own
  `throw new Error` if `type` is omitted, `ToggleGroup`'s Go signature should
  make `type string` unavoidable to reason about (Go can't throw at
  render-construction time the way React can at mount, so this is a doc-
  comment/API-design note, not a runtime check to port).
  `ToggleGroupItem`'s own `role`/`aria-checked`-vs-`aria-pressed` split (see
  Behavior contract) means the item component needs to know the group's
  `type` to decide which ARIA attribute pair to stamp — same shape as
  `variant`/`size` inheriting from group to item (see below), `type` is a
  THIRD piece of group-to-item inherited state.
- **Group→item inheritance is a real gap, same shape as `## tabs`'
  `selected`**: shadcn does this via `ToggleGroupContext.Provider` (React
  context) — `variant`/`size`/`spacing` set on `ToggleGroup` are read by
  every `ToggleGroupItem` automatically; a caller only needs to override
  variant/size on an individual item to depart from the group. gsx has no
  context; the caller must pass `variant`/`size`/`spacing`/`type` explicitly
  to both `ToggleGroup` AND every `ToggleGroupItem` (same "caller already has
  both values in scope, resolve and pass down explicitly" shape as
  `## tabs`' `selected bool` comparison) — this is real ergonomic friction
  worth flagging (every item call site repeats 3-4 params the group already
  has) but not a new pattern; consistent with the rest of this codebase's
  no-context stance.
- **Selector-shape note for the group→item CSS relationship**: shadcn's own
  `group/toggle-group` + nova's `group-data-[variant=outline]/toggle-group:`-style
  selectors are the CSS-side version of the same inheritance — since gsxui's
  item ALREADY receives `data-variant`/`data-size`/`data-spacing` directly
  (passed explicitly per the point above, not inherited via a parent-scoped
  selector), gsxui's class strings should read the item's OWN
  `data-[variant=]`/`data-[spacing=]` attributes directly rather than a
  `group-data-[...]/toggle-group:` parent-selector chain — the CSS
  consequence of gsx's "no context" design is that the item-local selector
  shape (which new-york-v4 mostly already uses) is the right one to keep,
  and nova's few `group-data-.../toggle-group:` tokens (root shadow,
  orientation-aware corners) should be REWRITTEN to item-local
  `data-[orientation=]:`/root-shadow-stays-on-root-only equivalents rather
  than ported as literal parent-scoped selectors — already reflected in the
  merged class strings above.
- **Roving focus reuses `dropdown.js`'s idiom, not its code**: dropdown's
  `[role="menuitem"]:not([aria-disabled])` arrow-key walk (`ui/dropdown.js`)
  is architecturally the closest existing precedent (real DOM focus moves
  between real elements on ArrowDown/Up, `.focus()`-driven, no
  `aria-activedescendant`) — but toggle-group's items are ALWAYS in the
  page's tab order (roving tabindex, not "only reachable once a popover is
  open"), so the initial per-render tabindex assignment must be computed at
  SSR time, the same way `## tabs` server-renders `tabindex="0"` only on the
  `selected` trigger. Recommend: `ToggleGroupItem`'s Go component stamps
  `tabindex="0"` on the item that is BOTH pressed (if `type="single"`, the
  one matching value; if `type="multiple"`, none are privileged this way —
  see below) and `tabindex="-1"` on the rest; a small `toggle-group.js`
  reuses dropdown.js's arrow-key-walk shape but writes real `tabindex`
  toggling (0/-1) instead of dropdown's fixed `-1`/`.focus()`-only pattern,
  since these items must remain normal Tab stops when the group itself isn't
  focused.
  For `type="multiple"` with nothing pressed at SSR time, Radix's own
  fallback is "first focusable item" (no `active` item exists) — port the
  same default: first non-disabled item gets `tabindex="0"`.
- **Vertical orientation is new functionality, not free** — see the corner-
  rounding selector note above; if `orientation="vertical"` ships at all
  (roadmap doesn't explicitly ask for it, but nova's own `.cn-toggle-group-item`
  already assumes it exists), `data-orientation` must be a real stamped
  attribute on the root (driving both the corner-rounding CSS above and
  `toggle-group.js`'s own ArrowUp/Down-vs-Left/Right key gating) — recommend
  scoping v1 to horizontal-only (matches the Tier 3 roadmap line's own
  wording, "roving focus") and ledgering vertical as a GAP, consistent with
  how `## tabs` scoped out `orientation` entirely.

### Demo inventory

Present in `registry/new-york-v4/examples/` (this task's stated source):
`toggle-group-demo.tsx` (type=multiple, variant=outline, 3 icon items — the
canonical top-of-page demo), `toggle-group-single.tsx` (type=single, default
variant, exercises the `role="radio"`/`aria-checked` path),
`toggle-group-outline.tsx` (byte-identical to `-demo.tsx`, a separate section
on the docs page), `toggle-group-lg.tsx` / `toggle-group-sm.tsx` (size axis),
`toggle-group-disabled.tsx` (group-level `disabled`), `toggle-group-spacing.tsx`
(type=multiple, variant=outline, `spacing={2}`, `size=sm`, plus per-item
`data-[state=on]:*:[svg]:fill-*` color overrides — the most demanding single
demo for CSS correctness: exercises spacing≠0 rounding, size=sm radius, AND
child-selector state overrides together).

Recommended 3 for the plan: `toggle-group-demo` (baseline), `toggle-group-single`
(the role/aria-checked branch — otherwise untested by the multiple-only demos
above), `toggle-group-spacing` (spacing/rounding/size stress test).

Note: `ui.shadcn.com`'s live docs page additionally shows `Sizes`, `Vertical`,
`Font Weight Selector` (a real-world composed example using `size-16`
custom items), `Filter`/`Sort`/`Date Range` (application-style single-select
groups), sourced from `registry/bases/radix/examples/toggle-group-example.tsx`
— a newer example tree outside this task's stated `new-york-v4/examples`
scope. Named here for completeness; not read as a binding source.

---

## slider

### Markup structure

new-york-v4 (`slider.tsx`) collapses cleanly onto ONE native
`<input type="range">` per the roadmap decision — there is no multi-part tree
to preserve, so this section documents what the FOUR Radix parts
(`Root`/`Track`/`Range`/`Thumb`, one `Thumb` per value) contribute, since each
one maps to a specific piece of native-range styling or a specific gap:

```
Root  <div data-slot="slider" data-disabled? data-orientation="horizontal"|"vertical" style="--radix-slider-thumb-transform:..." class="relative flex w-full touch-none items-center select-none ...">
  Track <span data-slot="slider-track" class="relative grow overflow-hidden rounded-full bg-muted data-[orientation=horizontal]:h-1.5 ...">
    Range <span data-slot="slider-range" class="absolute bg-primary data-[orientation=horizontal]:h-full ...">
  Thumb <span data-slot="slider-thumb" role="slider" aria-valuemin aria-valuenow aria-valuemax aria-orientation="horizontal"|"vertical" tabindex="0" data-orientation data-disabled? class="block size-4 shrink-0 rounded-full border border-primary bg-white shadow-sm ring-ring/50 ...">
  (one more Thumb per value in a multi-thumb/range slider)
```

Class-string merge, new-york-v4 → nova (`.cn-slider*`, `style-nova.css`
L1228-1242):

**Root**: new-york-v4 `data-[orientation=vertical]:min-h-44` → nova
`data-vertical:min-h-40` — `min-h-44 → min-h-40 (nova)` (keep new-york-v4's
`data-[orientation=vertical]:` selector syntax, nova's `data-vertical:` is a
custom Tailwind boolean-attribute variant gsxui doesn't have configured, same
non-adoption call as scroll-area's own scrollbar selector below). Everything
else on Root (`relative flex w-full touch-none items-center select-none
data-[disabled]:opacity-50 data-[orientation=vertical]:h-full
data-[orientation=vertical]:w-auto data-[orientation=vertical]:flex-col`) is
Root-only layout concerned with the 4-part tree (grow track + centered
thumbs) and has **no meaning on a single `<input>`** — none of it ports; see
gsxui adaptation notes.

**Track**: `bg-muted rounded-full`, `h-1.5→h-1 (nova)` `w-1.5→w-1 (nova)`
(horizontal values shown; vertical mirrors).

**Range** (the primary-colored filled portion): `bg-primary`, no delta.

**Thumb**: new-york-v4 `block size-4 shrink-0 rounded-full border
border-primary bg-white shadow-sm ring-ring/50 transition-[color,box-shadow]
hover:ring-4 focus-visible:ring-4 focus-visible:outline-hidden
disabled:pointer-events-none disabled:opacity-50` → nova `border-ring
ring-ring/50 relative size-3 rounded-full border bg-white
transition-[color,box-shadow] after:absolute after:-inset-2 hover:ring-3
focus-visible:ring-3 focus-visible:outline-hidden active:ring-3`:
- `size-4 → size-3 (nova)`
- `shadow-sm` → **(SHADOW-PRESENCE, drop)** — nova drops it outright.
- `border-primary` → nova's `border-ring` is a **(NOT ADOPTED — color)**:
  keep `border-primary`, nova's swap recolors the border, it doesn't resize
  it.
- `hover:ring-4`/`focus-visible:ring-4` → nova's `ring-3` and the entirely
  NEW `active:ring-3` state are **(nova, flagged — outside stated scope)**:
  ring WIDTH isn't in this task's adoptable-category list
  (heights/paddings/gaps/text/radii/svg sizes), and `active:ring-3` is a new
  interaction state, not a value substitution on an existing one. Presented,
  not pre-applied — planner call.
- `after:absolute after:-inset-2` **(nova, new)** is a hit-target
  compensation for the `size-4→size-3` shrink (an invisible `::after`
  enlarging the hoverable/draggable area beyond the now-smaller visible
  thumb) — see the cross-browser gap below, this specific mechanism has **no
  native-range equivalent at all**.

### Behavior contract

Traced from `@radix-ui/react-slider@1.4.1` `dist/index.mjs`:

- Thumb: `role="slider"`, `aria-valuemin={min}` `aria-valuenow={value}`
  `aria-valuemax={max}` `aria-orientation={orientation}`,
  `data-orientation` `data-disabled`.
- Keyboard (`BACK_KEYS`/`PAGE_KEYS`/`ARROW_KEYS` maps, per-orientation):
  Arrow key in the "forward" direction for the slider's orientation
  increases by `step`; the "back" direction decreases. `PageUp`/`PageDown`
  jump by `step * 10`. Holding Shift with an arrow key ALSO multiplies the
  step by 10 (`isSkipKey = isPageKey || (event.shiftKey && ARROW_KEYS.includes(...))`).
  Home/End are handled by the generic keydown branch
  (`PAGE_KEYS.concat(ARROW_KEYS)`, confirmed present) to jump to min/max.
- `minStepsBetweenThumbs` (default `0`): for multi-thumb sliders, prevents
  one thumb's drag/keyboard move from crossing closer than N steps to its
  neighbor.
- `inverted` (default `false`): flips which visual end is "low."

### gsxui adaptation notes

- **The whole 4-part tree collapses onto ARIA-and-keyboard-for-free**: a
  native `<input type="range" min max step value>` already implicitly
  carries `role="slider"` and auto-derives `aria-valuenow`/`aria-valuemin`/
  `aria-valuemax` from its own attributes — zero markup needed to reproduce
  the Thumb's ARIA contract. The browser's native keyboard model
  (Left/Right/Down/Up adjust by `step`, PageUp/PageDown by a larger jump,
  Home/End to min/max) already matches Radix's own model closely enough that
  **no keyboard JS is needed either** — this is the single biggest reason
  the roadmap called this ADAPT "styled native `<input type=range>`."
  `minStepsBetweenThumbs`/multi-thumb sliders and `orientation="vertical"`
  have NO native-range equivalent (`<input type=range>` is always a single
  scalar value, single thumb) — GAP, ledger multi-thumb/range-slider support
  as out of scope for a native-input port; a range slider (two thumbs, e.g.
  price-range) would need Radix's actual multi-part architecture or a
  from-scratch custom-JS thumb pair, not this ADAPT.
- **Filled-range visual (`bg-primary` Range) has NO free native
  equivalent** — this is the real remaining work the roadmap line correctly
  flags as "the work." A bare styled `<input type=range>` has ONE uniform
  track; there is no native way to paint "primary color from min to
  current-value, muted color from current-value to max" without either:
  1. **`accent-color`** (`accent-color: var(--primary)` on the input) —
     colors the thumb AND, in Chromium, automatically fills the track
     portion left of the thumb; Firefox colors the thumb but needs
     verification for whether it also fills the track (uncertain from static
     source reading, no live browser available this session — **flag as
     needing browser verification**, matching this codebase's own
     "ADAPT (verified in-browser)" convention elsewhere); Safari's
     `accent-color` support for range inputs specifically has historically
     lagged checkbox/radio — also needs live verification. Zero JS either
     way, but the visual fill is NOT guaranteed cross-browser-identical.
  2. **CSS custom-property + `linear-gradient()`** on
     `::-webkit-slider-runnable-track`/`::-moz-range-track`: e.g.
     `background: linear-gradient(to right, var(--primary) 0%, var(--primary) var(--fill), var(--muted) var(--fill), var(--muted) 100%)`
     where `--fill` is a percentage custom property. The INITIAL percentage
     is computable server-side from `value`/`min`/`max` at render time (zero
     JS for first paint), but keeping it in sync while the user DRAGS the
     thumb needs a small `input` event listener updating `--fill` inline —
     this is genuinely new JS the roadmap's "zero JS" framing doesn't
     obviously anticipate; flag as the open decision for this component
     (accept option 1's cross-browser uncertainty for a truly zero-JS ADAPT,
     or ship a small `slider.js` for option 2's guaranteed-consistent fill).
- **The `after:absolute after:-inset-2` hit-target trick has NO
  translation** — `::-webkit-slider-thumb`/`::-moz-range-thumb` are
  themselves pseudo-elements; CSS doesn't allow a pseudo-element on a
  pseudo-element (`::-webkit-slider-thumb::after` is not a legal/supported
  selector), so nova's compensation for its own `size-4→size-3` shrink is
  categorically unreachable on a native-input port. Concrete consequence:
  either don't adopt `size-3` (keep `size-4`, sidesteps the whole problem —
  simplest), or adopt `size-3` and accept a smaller real hit target than
  nova's own React/Radix version has (visual parity without the interaction
  parity nova's own extra markup buys it) — ledger explicitly, don't silently
  drop the compensation without noting the resulting hit-target regression.
- **Cross-browser thumb/track pseudo-elements needed** (both required, no
  single selector reaches every engine):
  ```css
  /* track */
  input[type=range]::-webkit-slider-runnable-track { @apply h-1 rounded-full bg-muted; }
  input[type=range]::-moz-range-track { @apply h-1 rounded-full bg-muted; }
  /* thumb */
  input[type=range]::-webkit-slider-thumb { @apply appearance-none size-4 rounded-full border border-primary bg-white; margin-top: calc((0.25rem - 1rem) / 2); /* WebKit doesn't vertically center the thumb on the track automatically */ }
  input[type=range]::-moz-range-thumb { @apply size-4 rounded-full border border-primary bg-white; }
  ```
  plus `input[type=range] { @apply appearance-none bg-transparent w-full; }` on
  the input itself to suppress UA chrome entirely first. The WebKit
  `margin-top` compensation (thumb height minus track height, halved) is a
  well-known cross-browser quirk (Firefox auto-centers via `::-moz-range-track`'s
  own box model, WebKit does not) worth calling out explicitly since it's
  the kind of thing that looks fine in isolation and misaligned next to a
  real shadcn slider without it.
- **Focus ring**: `focus-visible:ring-4`/`ring-3` CANNOT target
  `::-webkit-slider-thumb`/`::-moz-range-thumb` with a `:focus-visible`
  pseudo-class chain the normal way (`input:focus-visible::-webkit-slider-thumb`
  IS valid and works in Chromium/Safari; Firefox requires
  `input:focus-visible::-moz-range-thumb` — both must be written, mirroring
  the track/thumb pair above).

### Demo inventory

Present in `registry/new-york-v4/examples/`: only `slider-demo.tsx`
(`defaultValue={[50]} max={100} step={1}`, single thumb, `w-[60%]` wrapper) —
this is the ONLY slider example in the stated source directory.

For a richer 2-4 pick, `registry/bases/radix/examples/slider-example.tsx`
(newer tree, outside this task's stated scope, read here only to inform this
list) shows: Basic, **Range** (`defaultValue={[25,50]}`, two-thumb — the
multi-thumb GAP flagged above), **Multiple Thumbs** (`[10,20,70]`, three),
**Vertical** (`orientation="vertical" className="h-40"`, two side-by-side
verticals), **Controlled** (`value`/`onValueChange` + a live-updating label,
`min={0} max={1} step={0.1}`), **Disabled**.

Recommended for the plan: `slider-demo` (the only one within scope AND the
only one a single-`<input>` ADAPT can fully reproduce), noting Range/Multiple/
Vertical are demo-inventory items the roadmap's chosen ADAPT (single native
input) cannot represent at all — a real scoping decision for the plan, not
just a missing example to add later.

---

## scroll-area

### Markup structure

```
Root     <div data-slot="scroll-area" class="relative">
  Viewport <div data-slot="scroll-area-viewport" data-radix-scroll-area-viewport class="size-full rounded-[inherit] transition-[color,box-shadow] outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1">{children}</div>
  ScrollBar (vertical, default) <div data-slot="scroll-area-scrollbar" data-orientation="vertical" data-state="visible"|"hidden" class="flex touch-none p-px transition-colors select-none h-full w-2.5 border-l border-l-transparent">
    Thumb <div data-slot="scroll-area-thumb" data-state style="width/height: var(--radix-scroll-area-thumb-*)" class="relative flex-1 rounded-full bg-border">
  ScrollBar (horizontal, opt-in via <ScrollBar orientation="horizontal"/>) — h-2.5 flex-col border-t border-t-transparent
  Corner   <div data-slot="scroll-area-corner">   (only rendered when type!=="scroll" AND both scrollbars visible)
```

Class-string merge — **zero nova metric delta on every part**:
- `.cn-scroll-area-scrollbar` (`style-nova.css` L994-996):
  `data-horizontal:h-2.5 data-horizontal:flex-col data-horizontal:border-t
  data-horizontal:border-t-transparent data-vertical:h-full data-vertical:w-2.5
  data-vertical:border-l data-vertical:border-l-transparent` — the NUMBERS
  (`h-2.5`/`w-2.5`) are identical to new-york-v4's own; only the selector
  SYNTAX differs (`data-horizontal:`/`data-vertical:`, a custom Tailwind
  boolean-attribute variant vs. `data-[orientation=horizontal]:` — a
  Tailwind-config-level difference, not a metric one). Keep new-york-v4's
  `data-[orientation=...]:` selector shape verbatim; no delta to apply.
- `.cn-scroll-area-thumb` (L998-1000): `rounded-full` — byte-identical to
  new-york-v4's own. No delta.

Consistent with the `## nova density` retarget's own "8 components, no
delta" bucket (scroll-area wasn't in that batch since it didn't exist yet,
but lands in the same bucket): **port new-york-v4's class strings for
scroll-area verbatim, no nova overlay needed anywhere.**

### Behavior contract

Traced from `@radix-ui/react-scroll-area@1.2.12` `dist/index.mjs`:

- `type` prop (`"hover"` default | `"scroll"` | `"auto"` | `"always"`)
  controls scrollbar visibility TIMING:
  - `hover`: scrollbar shows on `pointerenter` over the scroll-area root,
    hides after leaving (with a hide delay, `ScrollAreaScrollbarHover`).
  - `scroll`: scrollbar shows only while actively scrolling and fades
    `SCROLL_TIMEOUT = 600`ms after the last scroll event
    (`ScrollAreaScrollbarScroll`).
  - `auto`: visibility follows whether the content actually overflows
    (`ScrollAreaScrollbarAuto`).
  - `always`: `data-state="visible"` permanently (`ScrollAreaScrollbarVisible`).
- Each scrollbar computes its own thumb SIZE via `getThumbSize` (a ratio of
  `viewport size / content size` against the scrollbar track's own length,
  driven by `ResizeObserver`s on both the viewport and the content) and its
  thumb POSITION via a hand-rolled drag implementation (`onThumbPointerDown`
  captures the pointer, `onDragScroll` maps pointer movement back to
  `viewport.scrollLeft`/`scrollTop`) — a full custom scrollbar, not a styled
  native one.
- `Corner` renders ONLY when `type !== "scroll"` AND both the horizontal and
  vertical scrollbars are simultaneously visible (content overflows both
  axes) — the small square gap-filler where the two scrollbars would
  otherwise leave an empty corner; `context.dir` flips which side it sits on
  in RTL via `--radix-scroll-area-corner-width`/`-height` CSS vars feeding
  the OTHER scrollbar's own inset.
- The Viewport hides the browser's OWN scrollbar entirely (`scrollbar-width:none`
  + injected `<style>` for `[data-radix-scroll-area-viewport]::-webkit-scrollbar{display:none}`)
  while remaining a real `overflow: auto/scroll` element underneath — wheel,
  trackpad, touch-momentum, and keyboard (arrow keys / PageUp/PageDown when
  the viewport itself has focus) scrolling all keep working natively; only
  the OS-drawn scrollbar CHROME is replaced by Radix's own custom thumb.

### gsxui adaptation notes — what CSS-first can and cannot reproduce

Per the roadmap's own ADAPT framing ("CSS `scrollbar-width`/`scrollbar-color`
styling first; Radix-style custom thumbs only if that falls short"):

**What a single native-scrolling `<div class="overflow-auto ...">` CAN
reproduce, no JS**:
- Thumb/track FLAT COLOR via the standardized `scrollbar-width: thin|auto|none`
  + `scrollbar-color: <thumb-color> <track-color>` pair (CSS Scrollbars
  Module Level 1, Baseline-available in Firefox and Chromium for some years,
  Safari added support in 18.2/Nov 2024) — this alone gets a themed,
  thin(ner) scrollbar close to shadcn's visual weight with ZERO JS, matching
  the roadmap's stated preference.
- All native scroll INPUT modalities (wheel/trackpad/touch/keyboard) for
  free, since the div stays a real `overflow: auto` element the whole time —
  strictly BETTER than Radix's own approach needing no compensating work at
  all (Radix keeps native scrolling too, this is a non-gap either way, named
  for completeness).
- RTL layout of the scrollbar itself — native scrollbars already respect
  document direction; Radix's own `--radix-scroll-area-corner-width` RTL
  offset math exists ONLY because Radix draws its own scrollbar from
  scratch. Non-gap; native wins for free here.

**What CSS-first genuinely CANNOT reproduce — enumerate as the ledger**:
1. **No independent thumb SHAPE control** (`rounded-full`, the `p-px` inset
   gap from the track edge shadcn's own `ScrollAreaScrollbar` class carries)
   via the standard `scrollbar-color` property — it paints the browser's OWN
   thumb geometry (squared-off, platform-dependent width), no radius/inset/
   border control at all. The LEGACY, WebKit-proprietary (but Chromium-
   inherited) `::-webkit-scrollbar-thumb`/`::-webkit-scrollbar-track`
   pseudo-elements DO allow full shape control (`border-radius`, `background`,
   `border`) — but **Firefox has never implemented them**, so full visual
   fidelity needs BOTH the standard properties (Firefox, newer Safari,
   Chromium) AND the legacy pseudo-elements (older Safari, Chromium) layered
   together, roughly doubling the CSS surface for one visual result — worth
   ledgering as ongoing maintenance surface, not a one-time cost.
2. **No `Corner` element/styling at all** — `::-webkit-scrollbar-corner`
   exists in WebKit/Blink only; the standard `scrollbar-*` properties have
   NO corner concept whatsoever; Firefox has no way to style the
   intersection square at all, period. GAP: drop `ScrollAreaCorner` for v1
   (matches the "drop the part, don't ship dead CSS" convention used
   throughout `docs/jsx-parity.md` for parts with no reachable native
   equivalent, e.g. dropdown's never-ported scroll buttons).
3. **No `type="hover"/"scroll"/"auto"/"always"` visibility-timing control**
   — native scrollbar visibility is entirely OS/browser policy (macOS's own
   "Show scroll bars: Always / When scrolling / Automatically" System
   Settings preference OVERRIDES any page-level attempt to mimic Radix's
   own show-on-hover/fade-after-600ms timing; a user with "Always" set will
   never see gsxui's scrollbar "hide" no matter what CSS says, and vice
   versa) — GAP, ledger as an accepted platform-policy override, not a
   fixable styling gap.
4. **No independent thumb SIZE control divorced from actual content ratio**
   — native thumb size genuinely IS the viewport/content ratio (same as
   Radix's own `getThumbSize` math), so this is actually a NON-gap — named
   only to close the loop, not a real parity risk.
5. **`focus-visible:ring-[3px]` on the viewport itself** ports unchanged
   either way — the viewport is a real, potentially-focusable div in BOTH
   versions (native `overflow:auto` divs become keyboard-focusable/scrollable
   once they receive `tabindex` or contain no other focusable content and a
   user tabs to them per browser heuristics) — trivial, no gap.

**Recommended v1 shape**: ONE `ui.ScrollArea` component wrapping children in
a single `<div class="overflow-auto [scrollbar-width:thin] [scrollbar-color:var(--border)_transparent] dark:[scrollbar-color:...] ...">`
(Root+Viewport+Scrollbar+Thumb collapsed into one node, no separate DOM parts
to keep in sync — unlike Radix's four) plus a supplementary
`::-webkit-scrollbar`/`-thumb`/`-track` block in `assets/gsxui.css` for
full-fidelity WebKit/Chromium thumb shape (mirroring the `::details-content`
precedent of shipping hand-authored CSS for a mechanism Tailwind utilities
can't reach, `## accordion` MECHANISM). `ScrollBar`/`ScrollAreaThumb`/
`ScrollAreaCorner` as separate gsx components do not exist in v1 — GAP,
matches this doc's own point 2 above. An `orientation="horizontal"` variant
(the `scroll-area-horizontal-demo` shape) is just `overflow-x-auto` on the
same single div, no second component needed (unlike Radix, which needs a
second literal `<ScrollBar orientation="horizontal"/>` element).

### Demo inventory

Present in `registry/new-york-v4/examples/`: `scroll-area-demo.tsx`
(vertical, `h-72 w-48 rounded-md border` wrapper, 50 tag rows +
`Separator` between each — a good test of the CSS-first thumb showing up
against a border-radius-clipped container, `rounded-[inherit]` on the
viewport matters here), `scroll-area-horizontal-demo.tsx` (`w-96 rounded-md
border whitespace-nowrap`, a `flex w-max` image gallery, explicit
`<ScrollBar orientation="horizontal"/>`).

Both are in scope and both are needed (only 2 exist in the stated source —
vertical is the default/implicit case, horizontal requires the explicit
opt-in every time under this ADAPT's collapsed-single-div design, so both
must be demoed to prove the axis switch actually works). Recommend both for
the plan.

---

## select (custom listbox)

### Markup structure

new-york-v4 (`select.tsx`) tree, annotated with the runtime attributes
`@radix-ui/react-select@2.3.1` actually stamps (traced in its `dist/index.mjs`
— several of these, especially the hidden native select and the
`aria-selected` nuance below, are NOT visible in the `.tsx` file at all):

```
Select (context only, no DOM)                                data-slot="select"
  SelectTrigger  <button type="button" data-slot="select-trigger" data-size="default"|"sm"
                   role="combobox" aria-controls={open?contentId:undefined} aria-expanded
                   aria-required aria-autocomplete="none" data-state="open"|"closed"
                   data-disabled? data-placeholder?>
    SelectValue    <span data-slot="select-value" style="pointer-events:none">{value|placeholder}</span>
    SelectIcon     <span aria-hidden><ChevronDownIcon class="size-4 opacity-50"/></span>
  SelectPortal
    SelectContent  <div data-slot="select-content" role="listbox" id={contentId} data-state
                     style="--radix-select-content-transform-origin...">
      SelectScrollUpButton   <div data-slot="select-scroll-up-button"><ChevronUpIcon/></div>
      SelectViewport <div class="p-1" ...>
        SelectGroup  <div data-slot="select-group" role="group" aria-labelledby={groupId}>
          SelectLabel  <div data-slot="select-label" id={groupId}>
          SelectItem   <div data-slot="select-item" role="option" aria-labelledby={textId}
                         aria-selected={isSelected && isFocused}     ← see behavior contract, NOT simply "is this the value"
                         data-state="checked"|"unchecked"            ← THIS one simply tracks "is this the value"
                         data-highlighted={isFocused?"":undefined}
                         aria-disabled data-disabled tabindex="-1">
            <span data-slot="select-item-indicator"><ItemIndicator (mount-gated on isSelected)><CheckIcon/></ItemIndicator></span>
            SelectItemText <span id={textId}>{children}</span>
        SelectSeparator <div data-slot="select-separator" aria-hidden>
      SelectScrollDownButton (mirror)
  SelectBubbleInput  <select aria-hidden="true" tabindex="-1" name required disabled form
                        autoComplete style="visually-hidden" defaultValue={value}>
                        <option value="">  (only if placeholder-shown AND no item itself has value="")
                        {...one real <option> per mounted SelectItemText, side-channel-collected}
                      </select>
```

Class-string merge, new-york-v4 → nova (`.cn-select-*`, `style-nova.css`
L1002-1065). **Provenance caveat, stated up front**: nova's `.cn-select-*`
classes target **react-aria-components' `Select`**, not Radix's — a
DIFFERENT underlying primitive with a different state-attribute vocabulary
(`data-open`/`data-closed`/`data-focused`/`data-selected`/`data-placeholder`
booleans, vs. Radix's `data-state=open|closed` + `aria-selected`/
`data-highlighted`/`data-[disabled]`). Only the NUMERIC/metric values below
transplant; the selector SHAPES do not — every merged class string keeps
new-york-v4's own Radix-shaped selectors, substituting only nova's numbers.

**SelectTrigger**: new-york-v4 `flex w-fit items-center justify-between
gap-2 rounded-md border border-input bg-transparent px-3 py-2 text-sm
whitespace-nowrap shadow-xs transition-[color,box-shadow] outline-none
focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50
disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive
aria-invalid:ring-destructive/20 data-[placeholder]:text-muted-foreground
data-[size=default]:h-9 data-[size=sm]:h-8
*:data-[slot=select-value]:line-clamp-1 *:data-[slot=select-value]:flex
*:data-[slot=select-value]:items-center *:data-[slot=select-value]:gap-2
dark:bg-input/30 dark:hover:bg-input/50 dark:aria-invalid:ring-destructive/40
[&_svg]:pointer-events-none [&_svg]:shrink-0
[&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground`
→ merged:
```
flex w-fit items-center justify-between gap-2→gap-1.5 (nova) rounded-md→rounded-lg (nova) border border-input bg-transparent
px-3→pr-2 pl-2.5 (nova, directional split) py-2 text-sm whitespace-nowrap
transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50
disabled:cursor-not-allowed disabled:opacity-50
aria-invalid:border-destructive aria-invalid:ring-destructive/20
data-[placeholder]:text-muted-foreground
data-[size=default]:h-9→h-8 (nova) data-[size=sm]:h-8→h-7 (nova)
data-[size=sm]:rounded-[min(var(--radius-md),10px)] (nova, new)
*:data-[slot=select-value]:line-clamp-1 *:data-[slot=select-value]:flex *:data-[slot=select-value]:items-center *:data-[slot=select-value]:gap-2→gap-1.5 (nova)
dark:bg-input/30 dark:hover:bg-input/50 dark:aria-invalid:ring-destructive/40
[&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground
```
Dropped: `shadow-xs` **(SHADOW-PRESENCE, drop)** — matches native-select's own
already-approved drop. Kept `line-clamp-1`/`items-center` on the
`*:data-[slot=select-value]:` chain even though nova's `.cn-select-value`
doesn't repeat them — nova's Select-Value is react-aria's own differently-shaped
primitive (see provenance caveat), so their absence there isn't evidence
new-york-v4's tokens should drop; only `gap-2→gap-1.5` is confidently a metric
delta. `focus-visible:ring-[3px]` kept verbatim rather than switched to
nova's `ring-3` syntax — matches this codebase's existing universal
focus-ring convention (dropdown/tooltip/button all already use
`ring-[3px]`, confirmed by direct read); switching one component's syntax
for no visual difference would be inconsistent. `aria-invalid:ring-3` (nova,
new explicit width for the invalid state, absent from new-york-v4 entirely)
— **(nova, flagged — outside stated scope)**, ring width again.

**SelectContent**: new-york-v4 `relative z-50 max-h-(--radix-select-content-available-height)
min-w-[8rem] origin-(--radix-select-content-transform-origin) overflow-x-hidden
overflow-y-auto rounded-md border bg-popover text-popover-foreground shadow-md
data-[side=...]:slide-in-from-* data-[state=closed]:animate-out ...` → nova
`bg-popover text-popover-foreground data-open:animate-in data-closed:animate-out
... ring-foreground/10 min-w-36 rounded-lg shadow-md ring-1 duration-100`:
- `rounded-md → rounded-lg (nova)`
- `min-w-[8rem]` (128px) `→ min-w-36` (144px) `(nova)`
- `border` → nova's `ring-1 ring-foreground/10` is **(NOT ADOPTED —
  border→ring)**: keep `border` per house convention, matching every other
  popover-family surface's own "NOT ADOPTED" entry (`## nova density`).
- `shadow-md` — kept both places, no SHADOW-PRESENCE removal here.
- `duration-100` (nova, new explicit token) — **(nova, flagged — outside
  stated scope)**: not adopted, this codebase's popover family has already
  standardized on 150ms (`## dropdown`'s `duration-150`); a one-off 100ms on
  select alone would be an inconsistency, not a parity win.
- The whole `data-[state=...]:animate-*`/`fade-*`/`zoom-*`/`slide-in-from-*`
  block should be REPLACED, not ported, by this codebase's own discrete-
  transition mechanism (see gsxui adaptation notes below) — same substitution
  `## dropdown`'s `DropdownMenuContent` already made, ledgered in
  `docs/jsx-parity.md`'s `## animations` FINDING.

**SelectLabel**: new-york-v4 `px-2 py-1.5 text-xs text-muted-foreground` →
nova `text-muted-foreground px-1.5 py-1 text-xs` — `px-2→px-1.5 (nova)`
`py-1.5→py-1 (nova)`.

**SelectItem**: new-york-v4 `relative flex w-full cursor-default items-center
gap-2 rounded-sm py-1.5 pr-8 pl-2 text-sm outline-hidden select-none
focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none
data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0
[&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground
*:[span]:last:flex *:[span]:last:items-center *:[span]:last:gap-2` → nova
`gap-1.5 rounded-md py-1 pr-8 pl-1.5 text-sm [&_svg:not([class*='size-'])]:size-4
*:[span]:last:flex *:[span]:last:items-center *:[span]:last:gap-2`:
- `gap-2 → gap-1.5 (nova)`
- `rounded-sm → rounded-md (nova)` — note this is the one place radius goes
  UP one step from `sm`, not the usual `md→lg`; `select-item` nests inside
  the already-`rounded-lg` content box, matching e.g. `dropdown-item`'s own
  `rounded-md` (a nested/smaller element keeps a smaller radius than its
  container) — consistent with the rest of this codebase's radius scale,
  not an anomaly.
- `py-1.5 → py-1 (nova)`
- `pl-2 → pl-1.5 (nova)`
- `pr-8` unchanged (reserves the absolute-positioned indicator's slot).

**Item indicator**: new-york-v4 `absolute right-2 flex size-3.5 items-center
justify-center` → nova `pointer-events-none absolute right-2 flex size-4
items-center justify-center` — `size-3.5 → size-4 (nova)` — **the one place
nova's icon size GROWS rather than shrinks**, worth flagging since it breaks
the otherwise-consistent "nova shrinks" pattern every other part in this map
follows; `pointer-events-none` (nova, new) — a defensive addition since the
indicator span sits inside the clickable item, adopt (harmless, on-mechanism
CSS, not a color/box-model change).

**SelectSeparator**: `bg-border -mx-1 my-1 h-px` — byte-identical in nova.
No delta.

**Scroll up/down buttons**: new-york-v4 `flex cursor-default items-center
justify-center py-1` → nova adds `bg-popover z-10` — **(nova, flagged —
outside stated scope: color + z-index, neither in the adoptable list)** —
moot regardless per the recommendation to drop these two parts entirely, see
below.

### Behavior contract

Traced from `@radix-ui/react-select@2.3.1` `dist/index.mjs`:

- **Trigger**: `role="combobox"`, `aria-controls` (content id, only while
  open), `aria-expanded`, `aria-required`, `aria-autocomplete="none"`.
  `OPEN_KEYS = [" ", "Enter", "ArrowUp", "ArrowDown"]` — any of these open
  the listbox. **Typeahead works on the CLOSED trigger too**: every
  single-character keypress calls the typeahead handler regardless of open
  state; the exact rule for Space is `if (isTypingAhead && event.key === " ") return`
  — Space is swallowed as a search character (not treated as an open-request)
  ONLY while a typeahead search is already in progress (buffer non-empty),
  so a user can type multi-word searches like "Space Gray" into a closed
  trigger without every space re-triggering the open path.
- **Content**: `role="listbox"`. `onKeyDown`: `Tab` is `preventDefault`'d
  entirely (Tab never moves focus out of an open listbox — matches native
  `<select>`'s own open-dropdown behavior). Printable chars → typeahead.
  `ArrowUp`/`ArrowDown`/`Home`/`End` move focus among enabled items only
  (`getItems().filter(item => !item.disabled)`); `Home`/`ArrowUp` reverse the
  candidate list first. Escape/outside-pointerdown close via the standard
  Radix `DismissableLayer` (`onDismiss → onOpenChange(false)`). `contextmenu`
  inside the content is suppressed (`preventDefault`).
- **Item**: `role="option"`, `aria-labelledby` (its own generated text-id).
  **`aria-selected = isSelected && isFocused` — NOT simply "is this the
  current value"**: an item that IS the select's current value but is not
  ALSO the currently keyboard/pointer-highlighted item reports
  `aria-selected="false"`. This is easy to get wrong in a straight port —
  `data-state="checked"|"unchecked"` is the separate, simpler attribute that
  DOES simply track "is this the value" (and is what the visual checkmark
  CSS keys off); `aria-selected` tracks the AND of value-match and
  focus-match. `data-highlighted` (focus only, drives `focus:bg-accent`).
  `tabindex="-1"` always — items are NEVER real tab stops; the listbox
  manages focus entirely via `.focus()` calls (see the roving-vs-activedescendant
  recommendation below).
  Activation is pointer-type-aware: mouse uses `onPointerUp` to select+close
  (not `onClick`); touch/keyboard use `onClick`; `SELECTION_KEYS = [" ", "Enter"]`
  on `keydown` also select. Hovering an item with the mouse
  (`onPointerMove`, gated on `pointerType==="mouse"`) calls
  `.focus({preventScroll:true})` — **hover IS focus, the exact same idiom
  `ui/dropdown.js`'s own `pointerover→item.focus()` already implements**.
  Leaving the currently-focused item (`onPointerLeave`, only if it's still
  `document.activeElement`) clears the highlight via `onItemLeave` — same
  shape as `dropdown.js`'s content-level `pointerout` "leaving parks focus
  on the container" handler.
- **Typeahead** (independent instances on Trigger AND Content):
  `useTypeaheadSearch` buffers keystrokes for **1000ms** (`window.setTimeout`),
  resets the buffer to `""` after that idle window. `findNextItem` does
  **prefix `startsWith` matching** (NOT fuzzy scoring), wraps from the
  current item forward, and has a specific same-character-repeat rule: if
  every character typed so far is identical (e.g. "b","b","b"), the search
  is normalized to a SINGLE character and cycles through all items starting
  with that letter (repeated presses of the same key step through matches
  one at a time) rather than searching literally for "bbb".
- **Form association — the exact answer to "how does the value surface for
  forms"**: Radix renders a SECOND, entirely real, visually-hidden native
  `<select aria-hidden="true" tabindex="-1" name required disabled form
  autoComplete>` (`SelectBubbleInput`), populated with one real `<option>`
  per mounted `SelectItemText` (collected via a side-channel
  `SelectNativeOptionsProvider` context as each item mounts) plus a
  synthetic `<option value="">` when the current value is empty/placeholder
  AND no real item already has `value=""`. This hidden select's `.value` is
  set via `Object.getOwnPropertyDescriptor(HTMLSelectElement.prototype, "value").set`
  (bypassing React's own property tracking) followed by a real
  `dispatchEvent(new Event("change", {bubbles:true}))` — specifically so
  that native form submission, `FormData`, an associated `<label>`'s
  click-through, browser autofill, or any plain-DOM script reading
  `form.elements.fieldname.value` all see a completely ordinary, working
  `<select>`. **It is not a hidden `<input type="hidden">` — it is a second,
  fully-optioned, kept-in-sync `<select>`.**

### gsxui adaptation notes

- **`dropdown.js` machinery that applies AS-IS**:
  - Native popover API (`popover="auto"` on `SelectContent`) — top layer,
    light dismiss, free Esc, exactly `DropdownMenuContent`'s own mechanism.
  - `closest("[data-gsxui-select]")` trigger↔content proximity wiring —
    same shape as `closest("[data-gsxui-dropdown]")`.
  - The discrete-transition open/close animation block (`opacity-0
    scale-95 ... transition-discrete duration-150 open:opacity-100
    open:scale-100 starting:open:... data-[side=bottom]:starting:open:-translate-y-2`)
    ports byte-for-byte — this REPLACES the class-string block ledgered
    "not ported" under SelectContent above, per `## animations`' own
    established substitution.
  - `data-state="open"` stamped synchronously BEFORE `showPopover()` (same
    flash-avoidance fix `## dropdown`'s MECHANISM entry documents).
  - The trigger `pointerdown`-records-`wasOpen` / `click`-converges guard
    (`## dropdown` MECHANISM) — `SelectTrigger` is a real `<button>` too, so
    the exact same outside-pointerdown-vs-light-dismiss race applies
    verbatim.
  - Item hover-is-focus (`pointerover → item.focus()`) and content-leave
    clears highlight (`pointerout` + `relatedTarget` containment check) —
    confirmed above to be the SAME idiom Radix's own Select source
    implements, not just a convenient analogy.
- **`dropdown.js` machinery that does NOT apply / needs new logic**:
  - **Value model**: dropdown items are stateless fire-and-forget
    (`gsxui:select` then close). Select items must track ONE current value
    across the whole listbox, toggle exactly one item's `data-state="checked"`,
    update the trigger's displayed text, and (see below) sync a hidden
    native `<select>`. This is new state-management code, not a port.
  - **`aria-selected` needs its own focus-aware recompute**: since
    `aria-selected = isSelected && isFocused` (traced above, NOT simply
    "is this the value"), the SAME focus handler that ports dropdown's
    `pointerover`/arrow-key `.focus()` logic must ALSO recompute
    `aria-selected` on both the newly-focused item (true if it's also the
    value) and whichever item previously held it (now false) — dropdown's
    own item model has no "isSelected" concept at all to build this on top
    of; new logic layered onto the reused focus mechanism, not a copy.
  - **Focus model recommendation — real DOM focus (dropdown's model), NOT
    `aria-activedescendant` (command's model)**: gsxui has two existing
    precedents and they diverge here. `dropdown.js` keeps items permanently
    `tabindex="-1"` and moves REAL DOM focus onto them one at a time
    (`.focus()`-driven, confirmed this is ALSO exactly what Radix's actual
    Select source does — traced above, `event.currentTarget.focus({preventScroll:true})`
    on hover, `focusFirst(candidateNodes)` on arrow keys). `command.js`
    instead keeps focus permanently on the `<input>` and simulates a moving
    "current item" via `aria-activedescendant` + `data-selected`, because
    command's whole UX is built around never losing focus out of the search
    box while typing. Select's trigger is a `<button>`, not a text input —
    there's no input to keep focus pinned to, and no typeahead-into-a-
    visible-text-box UX to protect (Select's typeahead is invisible, buffer-
    only, per the behavior contract above). Recommend dropdown's real-focus
    model: it's both the more convenient fit for this codebase's existing
    code AND the behaviorally faithful port of what Radix's own Select
    source actually does (not a coincidence — traced, not assumed).
  - **Typeahead**: neither existing module is the right template.
    `command.js`'s `commandScore` fuzzy-ranking engine is the WRONG model —
    Select's typeahead is plain `startsWith` prefix matching with a 1000ms
    reset buffer and same-character-repeat cycling (`findNextItem`, traced
    above), a materially different algorithm. Recommend a small bespoke
    helper local to `select.js` (roughly 15-20 lines: buffer + `setTimeout`
    reset + `startsWith` filter + wrap-and-cycle-on-repeat) rather than
    reusing or generalizing either existing module.
  - **Hidden native `<select>` for form association — genuinely new
    mechanism, no existing precedent**: none of gsxui's native-first
    components (checkbox/radio/switch/`ui.NativeSelect` itself) pair
    themselves with a SECOND hidden control — they simply ARE the real
    control. `ui.Select` (custom listbox) needs to server-render a real
    `<select aria-hidden="true" tabindex="-1" class="sr-only">` sibling
    containing one real `<option>` per authored item value, with
    `name`/`required`/`disabled`/`form` forwarded from the same params
    `ui.NativeSelect` already exposes — meaning `ui.Select` and
    `ui.NativeSelect` should share their OPTION-AUTHORING shape (same
    `SelectOption`-equivalent param list) even though only one of them is
    visually a real `<select>`. `select.js` keeps this hidden select's
    `.value` in sync on every selection change and dispatches a `change`
    event on it — translated to plain DOM this is actually SIMPLER than
    Radix's own mechanism: `selectEl.value = value; selectEl.dispatchEvent(new Event("change", {bubbles:true}))`
    works directly with no property-descriptor bypass needed, since there's
    no virtual-DOM layer fighting a plain assignment the way React's own
    reconciler would.
- **Scroll buttons — recommend dropping both**: Radix needs
  `SelectScrollUpButton`/`DownButton` primarily because its DEFAULT
  `position="item-aligned"` sizes/positions the listbox to align the
  SELECTED item directly under the trigger (can clip at the viewport edge,
  needing incremental "scroll a bit more" affordances); the alternate
  `position="popper"` mode (anchors like an ordinary dropdown, sized to
  `--radix-select-trigger-width`) makes them far less necessary since the
  whole viewport just scrolls normally. gsxui's `dropdown.js` ALREADY only
  implements one fixed-below-trigger anchoring mode (no `item-aligned`
  concept exists anywhere in this codebase) — i.e. `ui.Select`'s anchoring
  is structurally always "popper-equivalent." Recommend dropping both scroll
  buttons entirely and letting the viewport's own `overflow-y-auto` scroll
  natively, exactly matching `DropdownMenuContent`'s own precedent (it has
  no scroll buttons either, for the same underlying reason) — GAP to
  ledger, not a missing feature so much as a consequence of an anchoring
  decision already made elsewhere in this codebase.
- **Groups/labels are a genuine structural port, not a collapse**: unlike
  `ui.NativeSelect`'s `SelectGroup`, which collapses onto `<optgroup
  label="...">` because a native optgroup can't hold an arbitrary styled
  label child (`## select`'s own MECHANISM entry), `ui.Select`'s custom
  listbox CAN hold arbitrary content, so `SelectGroup`→`<div role="group"
  aria-labelledby>` + `SelectLabel`→its own styled child both port as real,
  separate parts here — `ui.Select.SelectGroup` and `ui.NativeSelect.SelectGroup`
  are genuinely different shapes under the same-ish name, not aliases of
  each other; naming them distinctly (e.g. keep `NativeSelect`'s `SelectGroup`
  as-is, since `select`/`select` are both keywords-adjacent already per the
  packaging note, and give the new one an unambiguous name) is worth a
  planner decision.
- **Icon deps**: trigger chevron reuses `icon.ChevronDown` (already imported
  by `ui.NativeSelect`); item check indicator needs `icon.Check` — confirmed
  present in `ui/icon/icon_defs.go` (`Check = New("check")`) but not yet
  imported by any `ui/*.gsx` file — this IS the new `select → icon`
  dependency edge `internal/registry` will derive (alongside the
  pre-existing `native-select → icon` edge from `ui.ChevronDown`).

### Demo inventory

Present in `registry/new-york-v4/examples/`: `select-demo.tsx` (basic,
`SelectGroup` + `SelectLabel` + 5 `SelectItem`s, placeholder text, `w-[180px]`
trigger), `select-scrollable.tsx` (5 `SelectGroup`s / 27 total items spanning
world timezones — the long-list, multi-group, scroll-affordance stress test;
directly exercises the "drop scroll buttons, rely on native overflow" call
above under real long-content conditions).

Both are in scope and cover materially different ground (basic single-group
vs. long multi-group) — recommend both for the plan.

Note: `ui.shadcn.com`'s live docs additionally show Groups & Labels (with a
`SelectSeparator` between two groups — the one part-combination NEITHER
in-scope example exercises), Sizes (`sm`/`default` trigger heights), Disabled,
Invalid (`aria-invalid`), and a "With Icons" item-content variant, sourced
from `registry/bases/radix/examples/select-example.tsx` (newer tree, outside
this task's stated scope, read here only to inform this list — also notably
shows `NativeSelect` and the custom `Select` composed SIDE BY SIDE in one
"Inline" example, i.e. shadcn's own current docs treat them as two
permanently-coexisting components, not one superseding the other —
corroborates this port's own `NativeSelect`-stays-alongside decision from
`docs/component-roadmap.md`).

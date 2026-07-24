# Tier 4 Batch A source map — resizable, combobox, sidebar

Source analysis for Task 0 of the Tier 4 Batch A plan
(`.superpowers/sdd/2026-07-24-tier4-batch-a/task-0-brief.md`). Tasks 1–3 cite
this document by section heading (`## resizable`, `## combobox`,
`## sidebar`) for every class string and behavioral claim — do not re-derive
either from the `.tsx` files directly; cite here.

**Reference-source rule** (adjudicated, not re-litigated): markup reference
is `apps/v4/registry/new-york-v4/ui/<name>.tsx` — the only remaining
inline-Tailwind-utility form, gsxui's own authoring model. `registry/bases/
{radix,base,aria}/ui/<name>.tsx` (`.cn-*` named classes, utilities extracted
into `registry/styles/style-*.css`) are NOT ported from directly, but are
available as a structure tiebreak wherever nova's CSS changes a class name
that doesn't exist in the DOM at all — a markup-structure question nova's
stylesheet alone cannot answer. In this pass that tiebreak was actually
needed, and cited inline, only once: resizable's handle-icon shape (§6 of
`## resizable`, via `bases/base/ui/resizable.tsx`). Combobox and sidebar
raised no structure question nova's CSS couldn't answer on its own, so
their `bases/base/ui/*.tsx` counterparts were read for corroboration
(confirmed structurally identical to new-york-v4, differences noted where
found) but are not load-bearing citations anywhere in those two sections.
Visual reference is `registry/styles/style-nova.css`; where it and
new-york-v4 disagree, nova wins.

**Inputs read** (byte-read in full unless noted): `registry/new-york-v4/ui/
{resizable,combobox,sidebar}.tsx`; `registry/new-york-v4/examples/
resizable-*.tsx` (all 4); `registry/new-york-v4/examples/combobox-*.tsx`
(all 4 — turned out to be a *different, unrelated* component, see `##
combobox` FINDING); `apps/v4/examples/base/combobox-*.tsx` (all 11 — the
actual demos for `ui/combobox.tsx`); `registry/styles/style-nova.css`
(`.cn-resizable*`, `.cn-combobox*`, `.cn-sidebar*` sections); `registry/
bases/base/ui/{resizable,combobox,sidebar}.tsx` (structure tiebreak — read
for all three, but only load-bearing/cited inline for resizable's
handle-icon shape; see the Reference-source rule paragraph above);
`content/docs/components/base/{resizable,combobox}.mdx`
(shadcn's own prose docs, including a `react-resizable-panels` v3→v4
migration table — real documentation, still not the library's own built
output, cited and hedged accordingly); `docs/jsx-parity.md` (`## nova
density`, `## animations`, `## sheet`); `docs/component-roadmap.md`.

**Library presence check** (Step 1):
```
ls ~/personal/shadcn-ui/node_modules/react-resizable-panels/dist/  → no such directory (node_modules absent entirely)
ls ~/personal/shadcn-ui/node_modules/@base-ui/react/               → no such directory
find … -name 'react-resizable-panels' -o -iname '*base-ui*'        → no hits under node_modules
```
This checkout has **no `node_modules` at all** (confirmed: `ls -d
node_modules` → not found). Both `react-resizable-panels@4.5.8` and
`@base-ui/react@1.6.0` are declared in `apps/v4/package.json` and
`pnpm-lock.yaml` (versions confirmed from the lockfile) but never installed,
and the local pnpm content-addressable store (`~/Library/pnpm/store/v10`)
is opaque without an install step. **Every claim about either library's
runtime behavior — ARIA attributes not literally written in the `.tsx`,
keyboard handling, filter algorithm — is `derived-not-read`**, marked at
first use per component. sidebar.tsx has no comparable gap: it wraps no
third-party behavior library (only Radix's `Slot`, for `asChild`
polymorphism, plus gsxui's own already-shipped `sheet`/`tooltip`/
`collapsible`/`button`/`input`/`separator`/`skeleton`), so its behavior is
fully readable from the `.tsx` source itself — no `derived-not-read` tag
needed anywhere in `## sidebar`.

## Legend

- `derived-not-read` — reconstructed from the wrapper `.tsx`'s own
  props/classNames/data-attribute selectors, prose documentation, or general
  public-API knowledge, **not** from reading the library's own built
  source — unavailable in this checkout. Where I could not corroborate a
  claim even indirectly, I say so explicitly rather than filling the gap.
- `token→token (nova)` — nova changes this token's value; nova's value is
  what's recommended.
- `(nova, no delta)` — nova has a `.cn-*` entry for this part but its value
  equals new-york-v4's own.
- `(no nova counterpart)` — nova's stylesheet has no `.cn-*` class for this
  part at all; new-york-v4's value stands unreviewed by the retarget.
- `(NOT ADOPTED — color)` / `(NOT ADOPTED — border→ring)` — nova changed
  the token but the change is colors or the border→ring box-model swap,
  both out of scope per `docs/jsx-parity.md` `## nova density`; new-york-v4's
  token is kept, nova's alternative named for the record only.
- `(SHADOW-PRESENCE, drop)` — nova drops a shadow outright (no
  replacement); per house convention this IS adopted (removals only).
- Structure tiebreak citations (`bases/base/ui/<name>.tsx`) are called out
  by name every time they're used — never silently folded into "nova says."

---

## resizable

### 1. Markup structure — three parts, class strings verbatim

`resizable.tsx` (54 lines, full file read):

```tsx
function ResizablePanelGroup({ className, ...props }: ResizablePrimitive.GroupProps) {
  return (
    <ResizablePrimitive.Group
      data-slot="resizable-panel-group"
      className={cn(
        "flex h-full w-full aria-[orientation=vertical]:flex-col",
        className
      )}
      {...props}
    />
  )
}

function ResizablePanel({ ...props }: ResizablePrimitive.PanelProps) {
  return <ResizablePrimitive.Panel data-slot="resizable-panel" {...props} />
}

function ResizableHandle({ withHandle, className, ...props }: ResizablePrimitive.SeparatorProps & { withHandle?: boolean }) {
  return (
    <ResizablePrimitive.Separator
      data-slot="resizable-handle"
      className={cn(
        "relative flex w-px items-center justify-center bg-border after:absolute after:inset-y-0 after:left-1/2 after:w-1 after:-translate-x-1/2 focus-visible:ring-1 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:outline-hidden aria-[orientation=horizontal]:h-px aria-[orientation=horizontal]:w-full aria-[orientation=horizontal]:after:left-0 aria-[orientation=horizontal]:after:h-1 aria-[orientation=horizontal]:after:w-full aria-[orientation=horizontal]:after:translate-x-0 aria-[orientation=horizontal]:after:-translate-y-1/2 [&[aria-orientation=horizontal]>div]:rotate-90",
        className
      )}
      {...props}
    >
      {withHandle && (
        <div className="z-10 flex h-4 w-3 items-center justify-center rounded-xs border bg-border">
          <GripVerticalIcon className="size-2.5" />
        </div>
      )}
    </ResizablePrimitive.Separator>
  )
}
```

Three parts, three `data-slot`s: `resizable-panel-group`, `resizable-panel`,
`resizable-handle`. **`ResizablePanel` carries no class string of its own at
all** — no `cn()` call, no default className; any `className` a caller
passes flows straight through to `ResizablePrimitive.Panel`'s own `className`
prop, unmerged with anything. This is a genuine, notable asymmetry from the
other two parts, not an omission in this map.

### 2. Group orientation mechanism

`ResizablePanelGroup`'s own class carries `aria-[orientation=vertical]:flex-col`
— it reads the `aria-orientation` attribute the library's `Group` element
stamps on itself (confirmed by name only, not by reading the library — see
§5 migration table, which independently confirms `Group` gets
`aria-orientation`, replacing v3's `data-panel-group-direction`). The
`orientation` **prop** passed to `ResizablePanelGroup` (`"horizontal"` |
`"vertical"`, confirmed from all 4 demo files below) is what the library
turns into that `aria-orientation` attribute at runtime — `derived-not-read`
for the exact mechanism, but corroborated by the doc's own migration table.

### 3. The handle's inverted orientation — mapping table

The mapping between the **group's** `orientation` prop and the **handle's**
`aria-orientation` is inverted — and this is fully derivable from sources
already read in full (the two class strings + the demos), without needing
the library's own attribute-stamping code:

1. **ARIA semantics of `role="separator"`**: `aria-orientation` on a
   separator describes the separator itself, not the layout it sits in. A
   separator dividing two side-by-side (left/right) panels is, by that
   convention, a **vertical** separator (it's the panels that are arranged
   horizontally; the dividing line between them runs vertically).
2. **Group class**, read directly: `flex h-full w-full aria-[orientation=
   vertical]:flex-col` — the group is a `flex` **row** (panels side-by-side)
   unless its own `aria-orientation` is `vertical`, in which case it becomes
   `flex-col` (panels stacked). `resizable-demo-with-handle.tsx` passes
   `orientation="horizontal"` for its side-by-side "One"/"Two"+"Three"
   layout — confirming group `orientation="horizontal"` ⇒ `flex-row` ⇒
   panels side-by-side.
3. **Handle's base class**, read directly: `w-px` (no `aria-[orientation=
   horizontal]:` match) — a full-height, 1px-**wide** rule: a **vertical**
   line, which is exactly what divides side-by-side panels. Base classes
   apply whenever the `aria-[orientation=horizontal]:` variant does **not**
   match.
4. Chaining 1–3: side-by-side panels ⇒ group `orientation="horizontal"` ⇒
   the divider between them must be a vertical rule ⇒ the handle's base
   (unmatched) classes must be in effect ⇒ the handle's `aria-orientation`
   is **not** `"horizontal"` in this case — by ARIA semantics (point 1) it
   is `"vertical"`. Symmetrically, a `vertical`-orientation group
   (`flex-col`, panels stacked) needs a horizontal divider between them,
   which is exactly the `aria-[orientation=horizontal]:h-px
   aria-[orientation=horizontal]:w-full` branch — so the handle's
   `aria-orientation` is `"horizontal"` precisely when the group's own
   `orientation` is `"vertical"`.
5. **Corroborating detail**, read directly: `[&[aria-orientation=horizontal]
   >div]:rotate-90` — the `withHandle` grip icon rotates 90° exactly when
   the handle's own `aria-orientation` is `"horizontal"`, i.e. exactly when
   the rule itself has turned into a horizontal bar. A grip that was
   vertical-dots-shaped (three dots in a vertical line, `GripVerticalIcon`)
   rotated 90° becomes horizontal-dots-shaped — consistent with sitting on
   a horizontal rule, not a vertical one.

| Group's `orientation` prop | Group's own layout | Handle's `aria-orientation` | Handle renders as |
|---|---|---|---|
| `"horizontal"` (default — panels side-by-side) | `flex-row` | `"vertical"` | base classes: `w-px` — a full-height **vertical** rule |
| `"vertical"` (panels stacked) | `flex-col` | `"horizontal"` | `aria-[orientation=horizontal]:*` branch: `h-px w-full` — a full-width **horizontal** rule |

`derived-not-read`, stated narrowly: `react-resizable-panels`' own code
that actually sets the `aria-orientation` attribute on the `Separator` was
not read (library absent from this checkout, § Library presence check).
The table above is derived entirely from the ARIA `separator` role
convention plus the class strings and demos, all read directly — it assumes
the library follows that convention (which the corroborating rotate-90
detail in point 5 supports) rather than confirming it from the library's
own source.

### 4. `defaultSize` format

All 4 demo files (`resizable-demo.tsx`, `resizable-demo-with-handle.tsx`,
`resizable-handle.tsx`, `resizable-vertical.tsx`) pass `defaultSize` as a
**percentage string**: `"50%"`, `"25%"`, `"75%"`. Confirmed independently by
`content/docs/components/base/resizable.mdx`'s own v3→v4 changelog table:
`defaultSize={50}` (v3, numeric) → `defaultSize="50%"` (v4, string) — this
checkout is on v4 (`react-resizable-panels@4.5.8` per `package.json`/
`pnpm-lock.yaml`), so the string-percentage form is the only one to port.

### 5. Keyboard behavior of the library's `Separator` — `derived-not-read`

The library package is absent from this checkout (§ Library presence
check) — nothing below is read from `dist/index.mjs`. What I *can*
corroborate from `content/docs/components/base/resizable.mdx`'s own
v3→v4 migration table (real prose documentation, still not source):

```
PanelResizeHandle (v3) → Separator (v4)
data-panel-group-direction (v3 attr) → aria-orientation (v4 attr)
```

This confirms the v4 `Separator` is a focusable element (the class string
independently confirms this: `focus-visible:ring-1 focus-visible:ring-ring
focus-visible:ring-offset-1 focus-visible:outline-hidden` — dead CSS on a
non-focusable element) and that it exposes `aria-orientation`, consistent
with the WAI-ARIA "window splitter" `separator` role convention. **I could
not verify, and will not guess at**: the exact key set (commonly Arrow
keys nudge, Home/End jump to extremes, in the public `react-resizable-
panels` docs/changelog, but I have not read that prose in this checkout,
only inferred it's plausible from the ARIA role convention), the resize
step size per keypress, or whether Enter/Space do anything. Task 1–3
should treat resizable's keyboard contract as an open question requiring
either its own doc lookup or a deliberate gsxui-authored keyboard scheme,
not a value pinned here.

### 6. Handle icon — nova structure, NOT just a token swap

Nova's own `.cn-resizable-handle-icon` rule (`style-nova.css` lines
989–991, quoted verbatim):

```css
.cn-resizable-handle-icon {
  @apply bg-border h-6 w-1 rounded-lg;
}
```

That class name alone doesn't tell you what DOM it applies to — `registry/
bases/base/ui/resizable.tsx` (structure tiebreak, read in full) answers it:

```tsx
{withHandle && (
  <div className="cn-resizable-handle-icon z-10 flex shrink-0" />
)}
```

This is a **structural** change, not a metric one: nova's handle-icon `div`
is **completely empty** — no `<GripVerticalIcon/>` child, no border, no
`bg-border`/`rounded-xs` box around it, no fixed `h-4 w-3`. new-york-v4's
version is a bordered `h-4 w-3` box containing a `size-2.5` grip-dots icon;
nova's version is a solid `h-6 w-1 rounded-lg` bar with no icon glyph at
all — the "handle" affordance becomes a plain pill, not an icon-in-a-box.
Since nova wins on visual disagreement (adjudicated rule) and this is
exactly the kind of markup-structure question nova's CSS alone can't answer
(confirmed via the `bases/base` tiebreak, cited above), **recommend the
nova shape**: drop the `GripVerticalIcon` import/render and the
bordered-box wrapper entirely; render a single `<div class="z-10 h-6 w-1
shrink-0 rounded-lg bg-border">` (or equivalent) in its place. Flag this
prominently for Task 1–3 — it's a bigger delta than a number substitution,
easy to miss if only `style-nova.css`'s isolated rule is read without
cross-checking `bases/base`'s actual markup.

No other `.cn-resizable-*` entries exist in `style-nova.css` — confirmed by
`grep -n "\.cn-resizable" style-nova.css`, one hit total. `resizable-panel-
group` and `resizable-handle` themselves are `(no nova counterpart)` —
new-york-v4's own values for those two class strings stand unreviewed.

### 7. Demo inventory (all 4 `registry/new-york-v4/examples/resizable-*.tsx`, byte-read)

- `resizable-demo.tsx` — horizontal outer group (`max-w-md rounded-lg
  border md:min-w-[450px]`, panels `50%`/`50%`), right panel nests a
  vertical group (`25%`/`75%`). No handle (`withHandle` omitted on both
  `ResizableHandle`s). Exercises nested groups + the orientation-inversion
  case directly (§3).
- `resizable-demo-with-handle.tsx` — byte-identical structure, both
  `ResizableHandle`s pass `withHandle`. The doc's canonical "with visible
  handle" example.
- `resizable-handle.tsx` — single horizontal group, `min-h-[200px]
  max-w-md rounded-lg border md:min-w-[450px]`, panels `25%`("Sidebar")/
  `75%`("Content"), `withHandle`.
- `resizable-vertical.tsx` — single vertical group, same container class as
  `resizable-handle.tsx`, panels `25%`("Header")/`75%`("Content"), no
  `withHandle`.

Recommend porting all 4 as site examples — they're small and together
cover: nested groups, both orientations, and both handle-visibility states.

---

## combobox

### FINDING (before anything else): the `registry/new-york-v4/examples/combobox-*.tsx` demos do not exercise `ui/combobox.tsx` at all

Checked first because it would have derailed the whole section otherwise.
The four files at `registry/new-york-v4/examples/combobox-{demo,dropdown-
menu,popover,responsive}.tsx` (all byte-read) are the **old** Command+Popover
"combobox pattern" — they import `Command`/`CommandInput`/`CommandItem`/
etc. from `ui/command` and `Popover`/`PopoverTrigger`/`PopoverContent` from
`ui/popover`; **none of them import anything from `ui/combobox`** (confirmed
by `grep -rl "ui/combobox" registry/new-york-v4/examples/` → zero hits). The
actual demos for the `Combobox`/`ComboboxInput`/`ComboboxContent`/etc.
components in `ui/combobox.tsx` live at `apps/v4/examples/base/combobox-
*.tsx` (11 files: `demo`, `basic`, `multiple`, `clear`, `groups`, `custom`,
`invalid`, `disabled`, `auto-highlight`, `popup`, `input-group`, `rtl`),
importing from `@/styles/base-nova/ui/combobox` — a build-time-generated
path not physically present in this checkout, but confirmed structurally
identical to `registry/new-york-v4/ui/combobox.tsx` (same exported names,
same `data-slot`s, same primitive composition; only the class-string
authoring differs, `.cn-*` vs. inline utilities — the same base/new-york-v4
relationship as everywhere else in this map). Confirmed via `content/docs/
components/base/combobox.mdx`, which documents `ui/combobox.tsx`'s actual
API (`items`, `itemToStringValue`, `multiple`, `autoHighlight`, `showTrigger`,
`showClear`) and links every one of these 11 demos by name.

**Consequence for Task 1–3**: if you go looking for "the combobox demo" under
`registry/new-york-v4/examples/`, you will find working code that compiles
and looks plausible — and it is the **wrong component's** demo. Use
`apps/v4/examples/base/combobox-*.tsx` (§5 below) instead.

### 1. Every part's class string, verbatim (`registry/new-york-v4/ui/combobox.tsx`, 311 lines, full file read)

15 exported parts plus one hook (`useComboboxAnchor`). `data-slot` shown
inline; classes quoted exactly as in source (the `cn(...)` merge input, i.e.
before any caller `className` override).

- **`Combobox`** — bare alias: `const Combobox = ComboboxPrimitive.Root`. No
  class, no `data-slot` of its own (the primitive owns whatever it stamps).
- **`ComboboxValue`** — `data-slot="combobox-value"`, no class.
- **`ComboboxTrigger`** — `data-slot="combobox-trigger"`:
  ```
  [&_svg:not([class*='size-'])]:size-4
  ```
  Renders `{children}` then an unconditional `<ChevronDownIcon
  data-slot="combobox-trigger-icon" className="pointer-events-none size-4
  text-muted-foreground" />`.
- **`ComboboxClear`** — `data-slot="combobox-clear"`, composes
  `InputGroupButton` via `render={<InputGroupButton variant="ghost"
  size="icon-xs" />}` (Base UI's `render`-prop polymorphism, not
  Radix `asChild`); own class: `cn(className)` — i.e. **no base class of its
  own at all**, purely a merge passthrough. Renders `<XIcon
  className="pointer-events-none" />`.
- **`ComboboxInput`** — composition, not a single class string (see §2
  below for the full breakdown). Root wrapper is `<InputGroup
  className={cn("w-auto", className)}>`.
- **`ComboboxContent`** — `data-slot="combobox-content"`,
  `data-chips={!!anchor}`. The long string, quoted verbatim in full:
  ```
  group/combobox-content relative max-h-96 w-(--anchor-width) max-w-(--available-width) min-w-[calc(var(--anchor-width)+--spacing(7))] origin-(--transform-origin) overflow-hidden rounded-md bg-popover text-popover-foreground shadow-md ring-1 ring-foreground/10 duration-100 data-[chips=true]:min-w-(--anchor-width) data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 *:data-[slot=input-group]:m-1 *:data-[slot=input-group]:mb-0 *:data-[slot=input-group]:h-8 *:data-[slot=input-group]:border-input/30 *:data-[slot=input-group]:bg-input/30 *:data-[slot=input-group]:shadow-none data-open:animate-in data-open:fade-in-0 data-open:zoom-in-95 data-closed:animate-out data-closed:fade-out-0 data-closed:zoom-out-95
  ```
  This lives inside `ComboboxPrimitive.Portal` > `ComboboxPrimitive.
  Positioner` (props: `side="bottom"`, `sideOffset={6}`, `align="start"`,
  `alignOffset={0}`, `anchor`, own class `"isolate z-50"`) > `ComboboxPrimitive.
  Popup` (the element the string above is actually on).
  Note: new-york-v4's own value here already uses `ring-1 ring-foreground/10`
  (not `border`) — this predates nova, so the usual "NOT ADOPTED —
  border→ring" caveat (`docs/jsx-parity.md` `## nova density`) does not
  apply here; there's no border→ring swap to reject because shadcn's own
  source never had a border on this part.
- **`ComboboxList`** — `data-slot="combobox-list"`:
  ```
  max-h-[min(calc(--spacing(96)---spacing(9)),calc(var(--available-height)---spacing(9)))] scroll-py-1 overflow-y-auto p-1 data-empty:p-0
  ```
- **`ComboboxItem`** — `data-slot="combobox-item"`:
  ```
  relative flex w-full cursor-default items-center gap-2 rounded-sm py-1.5 pr-8 pl-2 text-sm outline-hidden select-none data-highlighted:bg-accent data-highlighted:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4
  ```
  Renders `{children}` then `<ComboboxPrimitive.ItemIndicator
  data-slot="combobox-item-indicator" render={<span className="pointer-events-none
  absolute right-2 flex size-4 items-center justify-center" />}><CheckIcon
  className="pointer-events-none size-4 pointer-coarse:size-5" /></...>`.
- **`ComboboxGroup`** — `data-slot="combobox-group"`, `cn(className)` only
  (no base class of its own).
- **`ComboboxLabel`** — `data-slot="combobox-label"`:
  ```
  px-2 py-1.5 text-xs text-muted-foreground pointer-coarse:px-3 pointer-coarse:py-2 pointer-coarse:text-sm
  ```
- **`ComboboxCollection`** — `data-slot="combobox-collection"`, no class.
- **`ComboboxEmpty`** — `data-slot="combobox-empty"`:
  ```
  hidden w-full justify-center py-2 text-center text-sm text-muted-foreground group-data-empty/combobox-content:flex
  ```
  (the `group-data-empty/combobox-content:` selector reads the named group
  from `ComboboxContent`'s own `group/combobox-content` class above.)
- **`ComboboxSeparator`** — `data-slot="combobox-separator"`:
  ```
  -mx-1 my-1 h-px bg-border
  ```
- **`ComboboxChips`** — `data-slot="combobox-chips"`:
  ```
  flex min-h-9 flex-wrap items-center gap-1.5 rounded-md border border-input bg-transparent bg-clip-padding px-2.5 py-1.5 text-sm shadow-xs transition-[color,box-shadow] focus-within:border-ring focus-within:ring-[3px] focus-within:ring-ring/50 has-aria-invalid:border-destructive has-aria-invalid:ring-[3px] has-aria-invalid:ring-destructive/20 has-data-[slot=combobox-chip]:px-1.5 dark:bg-input/30 dark:has-aria-invalid:border-destructive/50 dark:has-aria-invalid:ring-destructive/40
  ```
- **`ComboboxChip`** — `data-slot="combobox-chip"`:
  ```
  flex h-[calc(--spacing(5.5))] w-fit items-center justify-center gap-1 rounded-sm bg-muted px-1.5 text-xs font-medium whitespace-nowrap text-foreground has-disabled:pointer-events-none has-disabled:cursor-not-allowed has-disabled:opacity-50 has-data-[slot=combobox-chip-remove]:pr-0
  ```
  `showRemove` (default `true`) renders `<ComboboxPrimitive.ChipRemove
  render={<Button variant="ghost" size="icon-xs" />} className="-ml-1
  opacity-50 hover:opacity-100" data-slot="combobox-chip-remove"><XIcon
  className="pointer-events-none" /></...>`.
- **`ComboboxChipsInput`** — `data-slot="combobox-chip-input"`:
  ```
  min-w-16 flex-1 outline-none
  ```
- **`useComboboxAnchor`** — not a component: `() =>
  React.useRef<HTMLDivElement | null>(null)`. Exists purely so multi-select
  callers can share one `ref` between `ComboboxChips` (`ref={anchor}`) and
  `ComboboxContent` (`anchor={anchor}`) — confirmed by `combobox-multiple.tsx`
  / `combobox-rtl.tsx` demos (§5), both of which use it this exact way.

### 2. `ComboboxInput` composition

```tsx
function ComboboxInput({ className, children, disabled = false, showTrigger = true, showClear = false, ...props }) {
  return (
    <InputGroup className={cn("w-auto", className)}>
      <ComboboxPrimitive.Input render={<InputGroupInput disabled={disabled} />} {...props} />
      <InputGroupAddon align="inline-end">
        {showTrigger && (
          <InputGroupButton
            size="icon-xs" variant="ghost" asChild
            data-slot="input-group-button"
            className="group-has-data-[slot=combobox-clear]/input-group:hidden data-pressed:bg-transparent"
            disabled={disabled}
          >
            <ComboboxTrigger />
          </InputGroupButton>
        )}
        {showClear && <ComboboxClear disabled={disabled} />}
      </InputGroupAddon>
      {children}
    </InputGroup>
  )
}
```

Structure: `InputGroup` (gsxui already has `ui/input-group.gsx` — reuse,
not reimplement) wrapping `ComboboxPrimitive.Input` rendered *as* an
`InputGroupInput` (Base UI `render`-prop swap: the actual DOM node is
`InputGroupInput`'s `<input>`, with the combobox's own input behavior
layered onto it — not two nested elements), plus an `InputGroupAddon
align="inline-end"` holding **either** the trigger button (wrapping
`ComboboxTrigger`, itself wrapped as an `InputGroupButton`) **or** the
clear button, gated by two independent boolean props:

| `showTrigger` | `showClear` | Addon contents |
|---|---|---|
| `true` (default) | `false` (default) | trigger chevron button only |
| `false` | `true` | clear (X) button only |
| `true` | `true` | **both** render, but `group-has-data-[slot=combobox-clear]/input-group:hidden` on the trigger button hides it whenever a `combobox-clear` sibling is present in the same named `input-group` — i.e. clear visually wins when both are requested |
| `false` | `false` | addon renders empty (no icon at all) |

`{children}` is appended **after** the addon, inside `InputGroup` — this is
how `combobox-input-group.tsx` (§5) adds its own leading `InputGroupAddon`
(a `GlobeIcon`) without touching `ComboboxInput`'s own internals: it passes
that addon as `children`.

### 3. Filter question — **not resolved**, marked `derived-not-read`

The brief asks specifically: does `@base-ui/react`'s Combobox filter rank
(return a score) or filter (return a boolean)? Is it collator-based? Case/
accent handling?

**I could not answer this.** `@base-ui/react` is not present anywhere in
this checkout (§ Library presence check) — no `dist/`, no built output, no
type declarations. I also searched every `content/docs/**/*.mdx` file in
this checkout for filter-behavior prose (`grep -rn -i filter
content/docs/components/{base,radix,aria}/combobox.mdx`) and found only one
unrelated hit per file (`autoHighlight`'s doc line, which describes
highlighting the first *filtered* result, not the filter algorithm itself).
There is no local documentation describing the default filter's algorithm,
collator use, or case/accent sensitivity.

What I can state with confidence, from the `.tsx` source alone: `Combobox`
accepts an `items` prop (array — confirmed by every demo in §5, e.g.
`<Combobox items={frameworks}>`) and an optional `itemToStringValue` prop
(confirmed by `combobox-custom.tsx`/`combobox-popup.tsx`, used when `items`
are objects rather than plain strings) — so *some* built-in filtering
against `items` exists and is driven by that string projection, gated on
whatever the user types into `ComboboxInput`. That is the extent of what's
verifiable here. **Do not port a specific filter algorithm (fuzzy score vs.
substring boolean, `Intl.Collator` sensitivity, diacritic-folding) into
gsxui's own combobox based on this document** — that decision needs either
a real read of `@base-ui/react`'s source (not available in this checkout)
or a deliberate, independently-justified gsxui design choice, ledgered as
such, not a silent inference from this gap.

### 4. ARIA anatomy — `derived-not-read`

None of `role`, `aria-expanded`, `aria-controls`, `aria-activedescendant`,
or an item's `role="option"` appear anywhere in `combobox.tsx`'s own JSX —
every one of Base UI's ARIA stamps happens inside the primitive's own
runtime, which is not present in this checkout. What the general WAI-ARIA
"combobox" pattern would predict (input owns `role="combobox"` +
`aria-expanded` + `aria-controls` pointing at the listbox + `aria-
activedescendant` tracking the highlighted item; `ComboboxList`'s rendered
root as `role="listbox"`; `ComboboxItem` as `role="option"` with `aria-
selected`) is **plausible but unverified** — flagged as `derived-not-read`
in its strongest form: I am naming the standard pattern only as a
hypothesis for Task 1–3 to confirm against Base UI's actual published docs
(`base-ui.com/react/components/combobox#api-reference`, linked from
`content/docs/components/base/combobox.mdx`'s frontmatter, not fetched as
part of this pass) or source, not as a verified fact to implement against
directly.

### 5. Nova's `.cn-combobox-*` rules, quoted (`style-nova.css` lines 296–354, byte-read)

```css
.cn-combobox-content {
  @apply bg-popover text-popover-foreground data-open:animate-in data-closed:animate-out data-closed:fade-out-0 data-open:fade-in-0 data-closed:zoom-out-95 data-open:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 ring-foreground/10 *:data-[slot=input-group]:bg-input/30 *:data-[slot=input-group]:border-input/30 max-h-72 min-w-36 overflow-hidden rounded-lg shadow-md ring-1 duration-100 *:data-[slot=input-group]:m-1 *:data-[slot=input-group]:mb-0 *:data-[slot=input-group]:h-8 *:data-[slot=input-group]:shadow-none;
}
.cn-combobox-content-logical {
  @apply data-[side=inline-start]:slide-in-from-right-2 data-[side=inline-end]:slide-in-from-left-2;
}
.cn-combobox-label {
  @apply text-muted-foreground px-2 py-1.5 text-xs;
}
.cn-combobox-item {
  @apply data-highlighted:bg-accent data-highlighted:text-accent-foreground not-data-[variant=destructive]:data-highlighted:**:text-accent-foreground gap-2 rounded-md py-1 pr-8 pl-1.5 text-sm [&_svg:not([class*='size-'])]:size-4;
}
.cn-combobox-item-indicator {
  @apply pointer-events-none absolute right-2 flex size-4 items-center justify-center;
}
.cn-combobox-empty {
  @apply text-muted-foreground hidden w-full justify-center py-2 text-center text-sm group-data-empty/combobox-content:flex;
}
.cn-combobox-list {
  @apply no-scrollbar max-h-[min(calc(--spacing(72)---spacing(9)),calc(var(--available-height)---spacing(9)))] scroll-py-1 overflow-y-auto p-1 data-empty:p-0;
}
.cn-combobox-item-text {
  @apply flex flex-1 gap-2;
}
.cn-combobox-separator {
  @apply bg-border -mx-1 my-1 h-px;
}
.cn-combobox-trigger {
  @apply [&_svg:not([class*='size-'])]:size-4;
}
.cn-combobox-trigger-icon {
  @apply text-muted-foreground size-4;
}
.cn-combobox-chips {
  @apply dark:bg-input/30 border-input focus-within:border-ring focus-within:ring-ring/50 has-aria-invalid:ring-destructive/20 dark:has-aria-invalid:ring-destructive/40 has-aria-invalid:border-destructive dark:has-aria-invalid:border-destructive/50 flex min-h-8 flex-wrap items-center gap-1 rounded-lg border bg-transparent bg-clip-padding px-2.5 py-1 text-sm transition-colors focus-within:ring-3 has-aria-invalid:ring-3 has-data-[slot=combobox-chip]:px-1;
}
.cn-combobox-chip {
  @apply bg-muted text-foreground flex h-[calc(--spacing(5.25))] w-fit items-center justify-center gap-1 rounded-sm px-1.5 text-xs font-medium whitespace-nowrap has-data-[slot=combobox-chip-remove]:pr-0;
}
.cn-combobox-chip-remove {
  @apply -ml-1 opacity-50 hover:opacity-100;
}
```

(Nova also carries `.cn-combobox-item-aria` and `.cn-combobox-content-aria`,
lines 312–314 and 1636 — these target `react-aria-components`' Combobox, a
DIFFERENT underlying primitive from Base UI's, used by the `aria` style
family, not `base`. Excluded here — same provenance caveat the tier3
`## select (custom listbox)` map applies to `.cn-select-*` vs. Radix; do not
mix the two families' selector shapes.)

Metric deltas worth naming explicitly (new-york-v4 → nova):
`rounded-md→rounded-lg` (content), `max-h-96→max-h-72` (content, 384px→
288px) with `.cn-combobox-list`'s own inner `max-h-[...96...]→max-h-[...72...]`
matching, `rounded-sm→rounded-md` + `py-1.5→py-1` + `pl-2→pl-1.5` (item),
`rounded-md→rounded-lg` + `px-2.5 py-1.5→px-2.5 py-1` (chips), `shadow-xs`
on `ComboboxChips` **(SHADOW-PRESENCE, drop — not present in `.cn-combobox-
chips`)**, `ring-[3px]→ring-3` (chips focus/invalid rings — same value,
Tailwind v4 spelling only, no visual delta, matching the `## select`
precedent for this exact token). `.cn-combobox-item`'s `not-data-[variant=
destructive]:` selector references a `data-variant` combobox items don't
actually get in `combobox.tsx` (no `variant` prop on `ComboboxItem`
anywhere in the source) — likely forward-shared CSS with `command`'s own
item variant (`CommandItem` does support a destructive variant elsewhere in
this registry), harmless dead selector on this component, port or drop at
Task 1–3's discretion, not load-bearing either way.

### 6. Demo inventory (`apps/v4/examples/base/combobox-*.tsx`, all 11 byte-read; the `registry/new-york-v4/examples/combobox-*.tsx` files are NOT these — see FINDING above)

- `combobox-demo.tsx` / `combobox-basic.tsx` — byte-identical single-select
  baseline, 5-item string array, no extra props.
- `combobox-multiple.tsx` — `multiple autoHighlight`, `defaultValue={[frameworks[0]]}`,
  `ComboboxChips`+`ComboboxValue`+`ComboboxChip`+`ComboboxChipsInput`,
  `useComboboxAnchor()` shared via `ref`/`anchor`.
  Exercises the whole chips composition.
- `combobox-clear.tsx` — `showClear` on `ComboboxInput`, `defaultValue` set.
- `combobox-groups.tsx` — nested `items: [{value, items}]` shape,
  `ComboboxGroup items={group.items}` + `ComboboxLabel` +
  `ComboboxCollection` + conditional `ComboboxSeparator` between groups.
- `combobox-custom.tsx` — object items, `itemToStringValue`, custom
  `ComboboxItem` children (composes `ui/item`'s `Item`/`ItemContent`/
  `ItemTitle`/`ItemDescription` — gsxui has no `item` port status confirmed
  in this pass, check `docs/component-roadmap.md` before assuming it's
  available).
- `combobox-invalid.tsx` — `aria-invalid="true"` on `ComboboxInput`.
- `combobox-disabled.tsx` — `disabled` on `ComboboxInput`.
- `combobox-auto-highlight.tsx` — `autoHighlight` prop alone, no other
  variation from basic.
- `combobox-popup.tsx` — `ComboboxTrigger` composed with `render={<Button
  variant="outline".../>}`, `ComboboxValue` as trigger content,
  `ComboboxInput showTrigger={false}` moved **inside** `ComboboxContent`
  (the "trigger a popup from a button, search happens inside the popup"
  pattern — a real structural variant, not just a prop toggle).
- `combobox-input-group.tsx` — `ComboboxInput` given a leading `children`
  addon (`InputGroupAddon` with a `GlobeIcon`), grouped items, custom
  `alignOffset={-28}` and `className="w-60"` on `ComboboxContent`.
- `combobox-rtl.tsx` — RTL demo, out of scope for this pass (no RTL story
  in gsxui per this map's own scope).

Recommend `combobox-basic`, `combobox-multiple`, `combobox-groups`, and
`combobox-clear` as the initial site examples — together they exercise
single-select, multi-select-with-chips, groups+separators, and the clear
button, without needing the popup/input-group/custom-item variants for a
first pass.

---

## sidebar

No third-party behavior library — everything below is read directly from
`registry/new-york-v4/ui/sidebar.tsx` (727 lines, full file read). The only
non-gsxui-owned dependency is Radix's `Slot` (from the `radix-ui` npm
package, used only for `asChild` polymorphism on `SidebarGroupLabel`/
`SidebarGroupAction`/`SidebarMenuButton`/`SidebarMenuSubButton` — gsxui has
no `asChild` equivalent; every part below should just render its own tag
directly, dropping the `asChild`/`Slot` branch entirely, matching how gsxui
already handles Radix `asChild` elsewhere per `docs/jsx-parity.md`'s WIN
entries). All other composed parts (`Button`, `Input`, `Separator`, `Sheet`+
`SheetContent`+`SheetHeader`+`SheetTitle`+`SheetDescription`, `Skeleton`,
`Tooltip`+`TooltipContent`+`TooltipProvider`+`TooltipTrigger`,
`useIsMobile`) are gsxui components already shipped (confirmed: `ui/
{button,input,separator,sheet,skeleton,tooltip}.gsx` all present in `ui/`)
— per `docs/component-roadmap.md`'s own framing, sidebar "depends on sheet +
tooltip + collapsible (all shipped)" (collapsible is not actually composed
by `sidebar.tsx` itself — `useIsMobile` and the rail/gap mechanism handle
the collapse visually via CSS, not via `ui/collapsible`; the roadmap note
may be anticipating a `SidebarGroupCollapsible` pattern not present in this
`.tsx` file, which the docs site's block examples layer on top separately —
worth a note for Task 1–3, not a contradiction to resolve here).

### 1. Six constants

```
SIDEBAR_COOKIE_NAME       = "sidebar_state"
SIDEBAR_COOKIE_MAX_AGE    = 60 * 60 * 24 * 7          (604800, seconds — 7 days)
SIDEBAR_WIDTH             = "16rem"
SIDEBAR_WIDTH_MOBILE      = "18rem"
SIDEBAR_WIDTH_ICON        = "3rem"
SIDEBAR_KEYBOARD_SHORTCUT = "b"
```

`SidebarProvider` persists open/collapsed state via `document.cookie =
"${SIDEBAR_COOKIE_NAME}=${openState}; path=/; max-age=${SIDEBAR_COOKIE_MAX_AGE}"`
on every `setOpen` call — a plain non-HttpOnly cookie, readable/writable
from client JS, not a server round-trip. The keyboard shortcut is wired via
a `window` `keydown` listener checking `event.key === SIDEBAR_KEYBOARD_
SHORTCUT && (event.metaKey || event.ctrlKey)` (i.e. Cmd/Ctrl+B), calling
`event.preventDefault()` then `toggleSidebar()`.

### 2. Every part — element, `data-slot`, class string, verbatim

**`SidebarProvider`** — no `data-slot` on itself; wraps children in
`TooltipProvider delayDuration={0}` then a `<div data-slot="sidebar-wrapper"
style="--sidebar-width:{W}; --sidebar-width-icon:{WI}">`:
```
group/sidebar-wrapper flex min-h-svh w-full has-data-[variant=inset]:bg-sidebar
```

**`Sidebar`** — three render branches, see §3 below. Common prop defaults:
`side="left"`, `variant="sidebar"`, `collapsible="offcanvas"`.

**`SidebarTrigger`** (composes `Button variant="ghost" size="icon"`) —
`data-sidebar="trigger"`, `data-slot="sidebar-trigger"`:
```
size-7
```
Renders `<PanelLeftIcon/>` + `<span class="sr-only">Toggle Sidebar</span>`.
`onClick` calls the caller's own handler (if any) then `toggleSidebar()`.

**`SidebarRail`** — `<button data-sidebar="rail" data-slot="sidebar-rail"
aria-label="Toggle Sidebar" tabindex="-1" title="Toggle Sidebar">`:
```
absolute inset-y-0 z-20 hidden w-4 -translate-x-1/2 transition-all ease-linear group-data-[side=left]:-right-4 group-data-[side=right]:left-0 after:absolute after:inset-y-0 after:left-1/2 after:w-[2px] hover:after:bg-sidebar-border sm:flex
in-data-[side=left]:cursor-w-resize in-data-[side=right]:cursor-e-resize
[[data-side=left][data-state=collapsed]_&]:cursor-e-resize [[data-side=right][data-state=collapsed]_&]:cursor-w-resize
group-data-[collapsible=offcanvas]:translate-x-0 group-data-[collapsible=offcanvas]:after:left-full hover:group-data-[collapsible=offcanvas]:bg-sidebar
[[data-side=left][data-collapsible=offcanvas]_&]:-right-2
[[data-side=right][data-collapsible=offcanvas]_&]:-left-2
```
`onClick={toggleSidebar}`. `tabIndex={-1}` — deliberately not
keyboard-focusable; it's a purely pointer-driven affordance (a drag-to-resize
rail visual, though this component only wires a click, not an actual drag —
confirmed: no pointerdown/pointermove handlers anywhere in this file, only
`onClick`).

**`SidebarInset`** — `<main data-slot="sidebar-inset">`:
```
relative flex w-full flex-1 flex-col bg-background
md:peer-data-[variant=inset]:m-2 md:peer-data-[variant=inset]:ml-0 md:peer-data-[variant=inset]:rounded-xl md:peer-data-[variant=inset]:shadow-sm md:peer-data-[variant=inset]:peer-data-[state=collapsed]:ml-2
```

**`SidebarInput`** (composes `Input`) — `data-slot="sidebar-input"`,
`data-sidebar="input"`:
```
h-8 w-full bg-background shadow-none
```

**`SidebarHeader`** — `<div data-slot="sidebar-header" data-sidebar="header">`:
```
flex flex-col gap-2 p-2
```

**`SidebarFooter`** — `<div data-slot="sidebar-footer" data-sidebar="footer">`:
```
flex flex-col gap-2 p-2
```

**`SidebarSeparator`** (composes `Separator`) — `data-slot="sidebar-separator"`,
`data-sidebar="separator"`:
```
mx-2 w-auto bg-sidebar-border
```

**`SidebarContent`** — `<div data-slot="sidebar-content" data-sidebar="content">`:
```
flex min-h-0 flex-1 flex-col gap-2 overflow-auto group-data-[collapsible=icon]:overflow-hidden
```

**`SidebarGroup`** — `<div data-slot="sidebar-group" data-sidebar="group">`:
```
relative flex w-full min-w-0 flex-col p-2
```

**`SidebarGroupLabel`** — `asChild` picks `Slot.Root` else `"div"` —
`data-slot="sidebar-group-label"`, `data-sidebar="group-label"`:
```
flex h-8 shrink-0 items-center rounded-md px-2 text-xs font-medium text-sidebar-foreground/70 ring-sidebar-ring outline-hidden transition-[margin,opacity] duration-200 ease-linear focus-visible:ring-2 [&>svg]:size-4 [&>svg]:shrink-0
group-data-[collapsible=icon]:-mt-8 group-data-[collapsible=icon]:opacity-0
```

**`SidebarGroupAction`** — `asChild` picks `Slot.Root` else `"button"` —
`data-slot="sidebar-group-action"`, `data-sidebar="group-action"`:
```
absolute top-3.5 right-3 flex aspect-square w-5 items-center justify-center rounded-md p-0 text-sidebar-foreground ring-sidebar-ring outline-hidden transition-transform hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 [&>svg]:size-4 [&>svg]:shrink-0
after:absolute after:-inset-2 md:after:hidden
group-data-[collapsible=icon]:hidden
```

**`SidebarGroupContent`** — `<div data-slot="sidebar-group-content" data-sidebar="group-content">`:
```
w-full text-sm
```

**`SidebarMenu`** — `<ul data-slot="sidebar-menu" data-sidebar="menu">`:
```
flex w-full min-w-0 flex-col gap-1
```

**`SidebarMenuItem`** — `<li data-slot="sidebar-menu-item" data-sidebar="menu-item">`:
```
group/menu-item relative
```

**`sidebarMenuButtonVariants`** (`cva`, base + 2 variants × 3 sizes — see §5
for the axes):
```
base: peer/menu-button flex w-full items-center gap-2 overflow-hidden rounded-md p-2 text-left text-sm ring-sidebar-ring outline-hidden transition-[width,height,padding] group-has-data-[sidebar=menu-action]/menu-item:pr-8 group-data-[collapsible=icon]:size-8! group-data-[collapsible=icon]:p-2! hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 data-[active=true]:bg-sidebar-accent data-[active=true]:font-medium data-[active=true]:text-sidebar-accent-foreground data-[state=open]:hover:bg-sidebar-accent data-[state=open]:hover:text-sidebar-accent-foreground [&>span:last-child]:truncate [&>svg]:size-4 [&>svg]:shrink-0
variant=default: hover:bg-sidebar-accent hover:text-sidebar-accent-foreground
variant=outline: bg-background shadow-[0_0_0_1px_var(--sidebar-border)] hover:bg-sidebar-accent hover:text-sidebar-accent-foreground hover:shadow-[0_0_0_1px_var(--sidebar-accent)]
size=default: h-8 text-sm
size=sm: h-7 text-xs
size=lg: h-12 text-sm group-data-[collapsible=icon]:p-0!
```

**`SidebarMenuButton`** — `asChild` picks `Slot.Root` else `"button"` —
`data-slot="sidebar-menu-button"`, `data-sidebar="menu-button"`,
`data-size={size}`, `data-active={isActive}`, class =
`sidebarMenuButtonVariants({variant, size})`. Tooltip gating: §4.

**`SidebarMenuAction`** — `asChild` picks `Slot.Root` else `"button"` —
`data-slot="sidebar-menu-action"`, `data-sidebar="menu-action"`:
```
absolute top-1.5 right-1 flex aspect-square w-5 items-center justify-center rounded-md p-0 text-sidebar-foreground ring-sidebar-ring outline-hidden transition-transform peer-hover/menu-button:text-sidebar-accent-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 [&>svg]:size-4 [&>svg]:shrink-0
after:absolute after:-inset-2 md:after:hidden
peer-data-[size=sm]/menu-button:top-1
peer-data-[size=default]/menu-button:top-1.5
peer-data-[size=lg]/menu-button:top-2.5
group-data-[collapsible=icon]:hidden
```
`showOnHover` (default `false`) ADDS, when true:
```
group-focus-within/menu-item:opacity-100 group-hover/menu-item:opacity-100 peer-data-[active=true]/menu-button:text-sidebar-accent-foreground data-[state=open]:opacity-100 md:opacity-0
```

**`SidebarMenuBadge`** — `<div data-slot="sidebar-menu-badge" data-sidebar="menu-badge">`:
```
pointer-events-none absolute right-1 flex h-5 min-w-5 items-center justify-center rounded-md px-1 text-xs font-medium text-sidebar-foreground tabular-nums select-none
peer-hover/menu-button:text-sidebar-accent-foreground peer-data-[active=true]/menu-button:text-sidebar-accent-foreground
peer-data-[size=sm]/menu-button:top-1
peer-data-[size=default]/menu-button:top-1.5
peer-data-[size=lg]/menu-button:top-2.5
group-data-[collapsible=icon]:hidden
```

**`SidebarMenuSkeleton`** (composes `Skeleton`) — `<div data-slot="sidebar-menu-skeleton"
data-sidebar="menu-skeleton">`:
```
flex h-8 items-center gap-2 rounded-md px-2
```
`showIcon` (default `false`) prepends `<Skeleton class="size-4 rounded-md"
data-sidebar="menu-skeleton-icon" />`. Always renders `<Skeleton
class="h-4 max-w-(--skeleton-width) flex-1" data-sidebar="menu-skeleton-text"
style="--skeleton-width:{W}">` where `W` is `React.useMemo(() =>
Math.floor(Math.random()*40)+50 + "%", [])` — **a random width in [50%,
90%) fixed once per mount**, not re-randomized on re-render. gsx has no
client-side `useMemo`-once equivalent for a server-rendered page; Task 1–3
should either compute the random width server-side at render time (fixed
per page load, matching the "once per mount" semantics closely enough for a
skeleton placeholder) or accept a static default width — flag as a design
decision, not resolved here.

**`SidebarMenuSub`** — `<ul data-slot="sidebar-menu-sub" data-sidebar="menu-sub">`:
```
mx-3.5 flex min-w-0 translate-x-px flex-col gap-1 border-l border-sidebar-border px-2.5 py-0.5
group-data-[collapsible=icon]:hidden
```

**`SidebarMenuSubItem`** — `<li data-slot="sidebar-menu-sub-item" data-sidebar="menu-sub-item">`:
```
group/menu-sub-item relative
```

**`SidebarMenuSubButton`** — `asChild` picks `Slot.Root` else `"a"` —
`data-slot="sidebar-menu-sub-button"`, `data-sidebar="menu-sub-button"`,
`data-size={size}` (`"sm"|"md"`, default `"md"`), `data-active={isActive}`:
```
flex h-7 min-w-0 -translate-x-px items-center gap-2 overflow-hidden rounded-md px-2 text-sidebar-foreground ring-sidebar-ring outline-hidden hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 [&>span:last-child]:truncate [&>svg]:size-4 [&>svg]:shrink-0 [&>svg]:text-sidebar-accent-foreground
data-[active=true]:bg-sidebar-accent data-[active=true]:text-sidebar-accent-foreground
{size==="sm" && "text-xs"} {size==="md" && "text-sm"}
group-data-[collapsible=icon]:hidden
```

That's 25 distinct `data-slot`-bearing markup parts (`sidebar-wrapper`,
`sidebar` ×3 branches, `sidebar-gap`, `sidebar-container`, `sidebar-inner`,
`sidebar-trigger`, `sidebar-rail`, `sidebar-inset`, `sidebar-input`,
`sidebar-header`, `sidebar-footer`, `sidebar-separator`, `sidebar-content`,
`sidebar-group`, `sidebar-group-label`, `sidebar-group-action`,
`sidebar-group-content`, `sidebar-menu`, `sidebar-menu-item`,
`sidebar-menu-button`, `sidebar-menu-action`, `sidebar-menu-badge`,
`sidebar-menu-skeleton` (+2 sub-parts `menu-skeleton-icon`/`-text`,
`data-sidebar` only, no own `data-slot`), `sidebar-menu-sub`,
`sidebar-menu-sub-item`, `sidebar-menu-sub-button`) — matching the brief's
"~25 parts" estimate.

### 3. `Sidebar`'s three render branches

```tsx
function Sidebar({ side = "left", variant = "sidebar", collapsible = "offcanvas", className, children, ...props }) {
  const { isMobile, state, openMobile, setOpenMobile } = useSidebar()

  if (collapsible === "none") {
    return (
      <div data-slot="sidebar" className={cn("flex h-full w-(--sidebar-width) flex-col bg-sidebar text-sidebar-foreground", className)} {...props}>
        {children}
      </div>
    )
  }

  if (isMobile) {
    return (
      <Sheet open={openMobile} onOpenChange={setOpenMobile} {...props}>
        <SheetContent
          data-sidebar="sidebar" data-slot="sidebar" data-mobile="true"
          className="w-(--sidebar-width) bg-sidebar p-0 text-sidebar-foreground [&>button]:hidden"
          style={{"--sidebar-width": SIDEBAR_WIDTH_MOBILE}}
          side={side}
        >
          <SheetHeader className="sr-only">
            <SheetTitle>Sidebar</SheetTitle>
            <SheetDescription>Displays the mobile sidebar.</SheetDescription>
          </SheetHeader>
          <div className="flex h-full w-full flex-col">{children}</div>
        </SheetContent>
      </Sheet>
    )
  }

  return (
    <div className="group peer hidden text-sidebar-foreground md:block"
         data-state={state} data-collapsible={state === "collapsed" ? collapsible : ""}
         data-variant={variant} data-side={side} data-slot="sidebar">
      <div data-slot="sidebar-gap" className={cn(
        "relative w-(--sidebar-width) bg-transparent transition-[width] duration-200 ease-linear",
        "group-data-[collapsible=offcanvas]:w-0",
        "group-data-[side=right]:rotate-180",
        variant === "floating" || variant === "inset"
          ? "group-data-[collapsible=icon]:w-[calc(var(--sidebar-width-icon)+(--spacing(4)))]"
          : "group-data-[collapsible=icon]:w-(--sidebar-width-icon)"
      )} />
      <div data-slot="sidebar-container" className={cn(
        "fixed inset-y-0 z-10 hidden h-svh w-(--sidebar-width) transition-[left,right,width] duration-200 ease-linear md:flex",
        side === "left"
          ? "left-0 group-data-[collapsible=offcanvas]:left-[calc(var(--sidebar-width)*-1)]"
          : "right-0 group-data-[collapsible=offcanvas]:right-[calc(var(--sidebar-width)*-1)]",
        variant === "floating" || variant === "inset"
          ? "p-2 group-data-[collapsible=icon]:w-[calc(var(--sidebar-width-icon)+(--spacing(4))+2px)]"
          : "group-data-[collapsible=icon]:w-(--sidebar-width-icon) group-data-[side=left]:border-r group-data-[side=right]:border-l",
        className
      )} {...props}>
        <div data-sidebar="sidebar" data-slot="sidebar-inner"
             className="flex h-full w-full flex-col bg-sidebar group-data-[variant=floating]:rounded-lg group-data-[variant=floating]:border group-data-[variant=floating]:border-sidebar-border group-data-[variant=floating]:shadow-sm">
          {children}
        </div>
      </div>
    </div>
  )
}
```

Branch 1 (`collapsible="none"`) is the simplest — a static flex column, no
state-driven classes at all, meant for e.g. a sidebar that's always
rendered in one context and never collapses. Branch 2 (mobile) is a full
`Sheet` reuse — `Sidebar` on mobile IS a `Sheet` (side-anchored, controlled
by `openMobile`/`setOpenMobile` from context, not the desktop `open`/
`setOpen` state at all — two independent booleans). Branch 3 (desktop) is
the complex one: three nested `div`s — an invisible `sidebar-gap`
spacer (reserves layout width, itself width-animated so surrounding content
reflows smoothly), a `fixed` `sidebar-container` (the actual visually
positioned sidebar, slid off-screen via negative `left`/`right` when
`collapsible=offcanvas` + collapsed), and inside that a `sidebar-inner` (the
bordered/rounded surface for the `floating` variant). All three read
`group-data-[...]` off the outermost `div`, which is both `group` (for this
internal 3-level nesting) AND `peer` (for `SidebarInset`, a sibling outside
this whole tree, to react to).

### 4. `data-*` the desktop root stamps, and who consumes it — table

The desktop-branch root `<div class="group peer ...">` stamps 4 data
attributes simultaneously: `data-state`, `data-collapsible`, `data-variant`,
`data-side` (plus `data-slot="sidebar"`, not consumed by any selector).

| Attribute | Values | Consumed by (selector) | Consuming part |
|---|---|---|---|
| `data-state` | `"expanded"` \| `"collapsed"` | `data-[state=open]:hover:...` *(sic — actually reads `SidebarMenuButton`'s OWN `data-state`, not the root's; see note)* | — |
| `data-state` (root, via `peer`) | `"expanded"` \| `"collapsed"` | `peer-data-[state=collapsed]:ml-2` | `SidebarInset` (only when ALSO `md:` + `variant=inset`) |
| `data-state` (root, via ancestor combo) | `"collapsed"` | `[[data-side=left][data-state=collapsed]_&]:cursor-e-resize` / the mirrored right variant | `SidebarRail` |
| `data-collapsible` | `""` \| `"offcanvas"` \| `"icon"` \| `"none"` (only ever non-empty when `state==="collapsed"`, else forced `""` — see Sidebar's own ternary) | `group-data-[collapsible=offcanvas]:*` (gap, container, rail) / `group-data-[collapsible=icon]:*` (gap, container, group-label, group-action, menu-button, menu-action, menu-badge, menu-sub) | many — the single most-consumed attribute in the component |
| `data-variant` | `"sidebar"` \| `"floating"` \| `"inset"` | `group-data-[variant=floating]:*` (inner: rounded/border/shadow) / ternaries in gap+container (`variant==="floating"||"inset"` branches, evaluated server-side in the `.tsx`, not a CSS selector) / `peer-data-[variant=inset]:*` (SidebarInset) | `sidebar-inner`, `sidebar-gap`, `sidebar-container`, `sidebar-inset` |
| `data-side` | `"left"` \| `"right"` | `group-data-[side=right]:rotate-180` (gap) / `in-data-[side=left]:cursor-w-resize` + mirrored (rail) / `[[data-side=left]...]` ancestor-combo selectors (rail) / `group-data-[side=left]:border-r` + mirrored (container, non-floating/inset only) | `sidebar-gap`, `sidebar-rail`, `sidebar-container` |

Note on the first row: `SidebarMenuButton`'s own class carries
`data-[state=open]:hover:bg-sidebar-accent` — this is **not** reading the
root Sidebar's `data-state` at all; it's an unrelated `data-state` that
`SidebarMenuButton` itself never stamps in this file (no `data-state=` in
its own JSX) — it must be intended for a caller-supplied `data-state`
(e.g. wrapping a `DropdownMenuTrigger asChild` whose own open/closed state
bubbles a `data-state` onto the button) or is dead CSS as authored. Flagged,
not resolved — port as dead weight per the house "port dead weight, ledger
it" convention (`docs/superpowers/plans/2026-07-24-tier3-source-map-
controls.md`'s own toggle-group FINDING sets this precedent).

Two more consumption patterns worth naming explicitly since they're easy to
misread:
- `in-data-[side=left]:cursor-w-resize` (on `SidebarRail`) — Tailwind v4's
  `in-*` variant matches when **some ancestor** matches the bracketed
  selector; this is NOT the same mechanism as `group-data-*` (which
  requires an explicit `.group` class on the ancestor) — it works on any
  ancestor regardless of whether it's marked `group`, which is why it can
  reach the same root div `.group.peer` carries both mechanisms on.
- `[[data-side=left][data-state=collapsed]_&]:cursor-e-resize` — an
  arbitrary-selector variant meaning "match `&` (this element) when some
  ancestor has BOTH `data-side=left` AND `data-state=collapsed` at once" —
  only the desktop root div ever carries both attributes simultaneously, so
  this resolves to exactly that one ancestor, not two separate ancestors
  each satisfying one condition.

### 5. `SidebarMenuButton` tooltip gating + variant/size axes

```tsx
function SidebarMenuButton({ asChild = false, isActive = false, variant = "default", size = "default", tooltip, className, ...props }) {
  const Comp = asChild ? Slot.Root : "button"
  const { isMobile, state } = useSidebar()
  const button = <Comp data-slot="sidebar-menu-button" data-sidebar="menu-button" data-size={size} data-active={isActive}
                        className={cn(sidebarMenuButtonVariants({variant, size}), className)} {...props} />
  if (!tooltip) return button
  if (typeof tooltip === "string") tooltip = { children: tooltip }
  return (
    <Tooltip>
      <TooltipTrigger asChild>{button}</TooltipTrigger>
      <TooltipContent side="right" align="center" hidden={state !== "collapsed" || isMobile} {...tooltip} />
    </Tooltip>
  )
}
```

- `tooltip` prop: absent → plain button, no `Tooltip` wrapper at all (not
  just hidden — the whole Radix Tooltip subtree is skipped). A `string` is
  sugar for `{children: tooltipString}`; an object is spread directly onto
  `TooltipContent` (any `TooltipContent` prop, e.g. `side`/`align`
  overrides, can be passed through this way, though `side="right"
  align="center"` are the button's own forced defaults unless the caller's
  object overrides them).
- Gating: `hidden={state !== "collapsed" || isMobile}` — the tooltip is
  **only ever shown when the sidebar is both collapsed AND on desktop**.
  Expanded desktop sidebars never show the tooltip (labels are visible
  inline instead); mobile never shows it regardless of state (mobile is a
  `Sheet`, which doesn't have this collapsed-icon-only mode at all).
- **Variant axis** (`variant` prop, `cva` default `"default"`): `"default"`
  | `"outline"` — see §2's `sidebarMenuButtonVariants` block for both
  strings.
- **Size axis** (`size` prop, `cva` default `"default"`): `"default"` |
  `"sm"` | `"lg"` — `h-8`/`h-7`/`h-12` respectively, `lg` additionally forces
  `group-data-[collapsible=icon]:p-0!` (an icon-only collapsed `lg` button
  gets zero padding, vs. the shared `group-data-[collapsible=icon]:size-8!
  p-2!` every size gets from the base string — `lg`'s own size override
  wins via the `!important` + later-in-source-order combination).
- `data-size={size}` and `data-active={isActive}` are stamped on the button
  itself (not derived from `sidebarMenuButtonVariants`, which only produces
  the class string) — these are what `SidebarMenuAction`/`SidebarMenuBadge`/
  `SidebarMenuSubButton`'s own `peer-data-[size=...]/menu-button:` and
  similar selectors read, per §4's `peer/menu-button` mechanism.

### 6. Nova's `.cn-sidebar-*` rules, quoted (`style-nova.css` lines 1110–1220, byte-read)

```css
.cn-sidebar-gap {
  @apply transition-[width] duration-200 ease-linear;
}
.cn-sidebar-inner {
  @apply bg-sidebar group-data-[variant=floating]:ring-sidebar-border group-data-[variant=floating]:rounded-lg group-data-[variant=floating]:shadow-sm group-data-[variant=floating]:ring-1;
}
.cn-sidebar-rail {
  @apply hover:after:bg-sidebar-border;
}
.cn-sidebar-inset {
  @apply bg-background md:peer-data-[variant=inset]:m-2 md:peer-data-[variant=inset]:ml-0 md:peer-data-[variant=inset]:rounded-xl md:peer-data-[variant=inset]:shadow-sm md:peer-data-[variant=inset]:peer-data-[state=collapsed]:ml-2;
}
.cn-sidebar-input {
  @apply bg-background h-8 w-full shadow-none;
}
.cn-sidebar-header {
  @apply gap-2 p-2;
}
.cn-sidebar-content {
  @apply no-scrollbar gap-0;
}
.cn-sidebar-footer {
  @apply gap-2 p-2;
}
.cn-sidebar-separator {
  @apply bg-sidebar-border mx-2;
}
.cn-sidebar-group {
  @apply p-2;
}
.cn-sidebar-menu {
  @apply gap-0;
}
.cn-sidebar-group-content {
  @apply text-sm;
}
.cn-sidebar-group-label {
  @apply text-sidebar-foreground/70 ring-sidebar-ring h-8 rounded-md px-2 text-xs font-medium transition-[margin,opacity] duration-200 ease-linear group-data-[collapsible=icon]:-mt-8 group-data-[collapsible=icon]:opacity-0 focus-visible:ring-2 [&>svg]:size-4;
}
.cn-sidebar-group-action {
  @apply text-sidebar-foreground ring-sidebar-ring hover:bg-sidebar-accent hover:text-sidebar-accent-foreground absolute top-3.5 right-3 w-5 rounded-md p-0 focus-visible:ring-2 [&>svg]:size-4;
}
.cn-sidebar-menu-button {
  @apply ring-sidebar-ring hover:bg-sidebar-accent hover:text-sidebar-accent-foreground active:bg-sidebar-accent active:text-sidebar-accent-foreground data-active:bg-sidebar-accent data-active:text-sidebar-accent-foreground data-open:hover:bg-sidebar-accent data-open:hover:text-sidebar-accent-foreground gap-2 rounded-md p-2 text-left text-sm transition-[width,height,padding] group-has-data-[sidebar=menu-action]/menu-item:pr-8 group-data-[collapsible=icon]:size-8! group-data-[collapsible=icon]:p-2! focus-visible:ring-2 data-active:font-medium;
}
.cn-sidebar-menu-button-variant-default {
  @apply hover:bg-sidebar-accent hover:text-sidebar-accent-foreground;
}
.cn-sidebar-menu-button-variant-outline {
  @apply bg-background hover:bg-sidebar-accent hover:text-sidebar-accent-foreground shadow-[0_0_0_1px_var(--sidebar-border)] hover:shadow-[0_0_0_1px_var(--sidebar-accent)];
}
.cn-sidebar-menu-button-size-default {
  @apply h-8 text-sm;
}
.cn-sidebar-menu-button-size-sm {
  @apply h-7 text-xs;
}
.cn-sidebar-menu-button-size-lg {
  @apply h-12 text-sm group-data-[collapsible=icon]:p-0!;
}
.cn-sidebar-menu-action {
  @apply text-sidebar-foreground ring-sidebar-ring hover:bg-sidebar-accent hover:text-sidebar-accent-foreground peer-hover/menu-button:text-sidebar-accent-foreground absolute top-1.5 right-1 aspect-square w-5 rounded-md p-0 peer-data-[size=default]/menu-button:top-1.5 peer-data-[size=lg]/menu-button:top-2.5 peer-data-[size=sm]/menu-button:top-1 focus-visible:ring-2 [&>svg]:size-4;
}
.cn-sidebar-menu-badge {
  @apply text-sidebar-foreground peer-hover/menu-button:text-sidebar-accent-foreground peer-data-active/menu-button:text-sidebar-accent-foreground pointer-events-none absolute right-1 flex h-5 min-w-5 rounded-md px-1 text-xs font-medium peer-data-[size=default]/menu-button:top-1.5 peer-data-[size=lg]/menu-button:top-2.5 peer-data-[size=sm]/menu-button:top-1;
}
.cn-sidebar-menu-skeleton {
  @apply h-8 gap-2 rounded-md px-2;
}
.cn-sidebar-menu-skeleton-icon {
  @apply size-4 rounded-md;
}
.cn-sidebar-menu-skeleton-text {
  @apply h-4;
}
.cn-sidebar-menu-sub {
  @apply border-sidebar-border mx-3.5 translate-x-px gap-1 border-l px-2.5 py-0.5 group-data-[collapsible=icon]:hidden;
}
.cn-sidebar-menu-sub-button {
  @apply text-sidebar-foreground ring-sidebar-ring hover:bg-sidebar-accent hover:text-sidebar-accent-foreground active:bg-sidebar-accent active:text-sidebar-accent-foreground [&>svg]:text-sidebar-accent-foreground data-active:bg-sidebar-accent data-active:text-sidebar-accent-foreground h-7 gap-2 rounded-md px-2 focus-visible:ring-2 data-[size=md]:text-sm data-[size=sm]:text-xs [&>svg]:size-4;
}
```

A full `grep -c "\.cn-sidebar"` pass on `style-nova.css` returns **28**
hits, not 27: the block above quotes 27 of them, deliberately omitting
`.cn-sidebar-menu-button-aria` (line 1170, sitting directly between
`.cn-sidebar-menu-button` and `.cn-sidebar-menu-button-variant-default`):

```css
.cn-sidebar-menu-button-aria {
  @apply aria-expanded:hover:bg-sidebar-accent aria-expanded:hover:text-sidebar-accent-foreground;
}
```

Excluded for the same provenance reason `## combobox` §5 excludes
`.cn-combobox-item-aria`/`.cn-combobox-content-aria`: confirmed via
`registry/bases/{radix,base,aria}/ui/sidebar.tsx` that only the
`aria`-flavored base (`bases/aria/ui/sidebar.tsx:479`) appends the
`cn-sidebar-menu-button-aria` token to `SidebarMenuButton`'s class list —
`bases/radix` and `bases/base` (the two families structurally matching
new-york-v4, our actual reference) never do. It targets react-aria-
components' own `aria-expanded` convention for an open submenu trigger, not
anything Radix's `SidebarMenuButton` (or new-york-v4's, which has no
built-in submenu-trigger concept at all) stamps. Not adopted; named here
only so the 28-vs-27 count is accounted for, not silently short.

**No nova counterpart at all** for: `sidebar` (the root/wrapper itself —
neither `sidebar-wrapper` nor the desktop root's own class), `sidebar-
trigger`, `sidebar-container`, `sidebar-menu-item`, `sidebar-menu-sub-item`
— confirmed by the same 28-hit pass: 27 quoted in the main block, 1
(`-aria`) quoted and excluded just above, none of the 28 matching those 5
names. These 5 parts' new-york-v4 values stand unreviewed by the nova
retarget.

Metric deltas worth naming: `rounded-lg + ring-1 ring-sidebar-border` (nova,
inner floating variant) replaces new-york-v4's `rounded-lg border border-
sidebar-border` — a **border→ring swap, `(NOT ADOPTED — border→ring)`**
per `docs/jsx-parity.md` `## nova density`; keep `border` here. `.cn-sidebar-
content`/`.cn-sidebar-menu` both add `gap-0` (nova drops the `gap-2`/`gap-1`
new-york-v4 has on these two — a real metric delta, adopt) plus `no-
scrollbar` on content (new, adopt — matches the scrollbar-hiding convention
`## carousel`/`## input-otp` also use elsewhere in this codebase's own
nova-retargeted parts). `SidebarMenuButton`'s nova base drops the `[&>span:
last-child]:truncate` and `disabled:*`/`aria-disabled:*` selectors present
in new-york-v4's own — **not** a metric change, just class-list reshuffling
around the same `cn-sidebar-menu-button*` bucket; recommend keeping
new-york-v4's fuller selector set as the base and only substituting nova's
metric deltas (`data-active:` bare-boolean spelling is `.cn-*`-family
authoring style, not a value change from new-york-v4's `data-[active=true]:`
— keep new-york-v4's bracket form, this codebase's existing convention per
every other Tier 3 doc's own "keep new-york-v4's selector shape" calls).

### 7. Dependencies — all shipped

`Sidebar` composes `Button`, `Input`, `Separator`, `Sheet`+`SheetContent`+
`SheetHeader`+`SheetTitle`+`SheetDescription`, `Skeleton`, `Tooltip`+
`TooltipContent`+`TooltipProvider`+`TooltipTrigger`. Confirmed present in
`ui/`: `button.gsx`, `input.gsx`, `separator.gsx`, `sheet.gsx` (+`sheet.x.go`,
no `sheet.js` — matches `docs/superpowers/plans/2026-07-24-tier3-source-
map-wrapped.md`'s own drawer finding that `sheet → dialog` derives its JS,
no standalone `sheet.js`), `skeleton.gsx`, `tooltip.gsx` (+`tooltip.js`).
Icons needed: `PanelLeftIcon` (gsxui `icon.PanelLeft`, confirmed present in
`ui/icon/icon_data.go` line 1127 as `"panel-left"`). No new icon or
dependency gaps for this component.

---

## Summary of corrections to the brief

- **Combobox demos**: the brief doesn't name specific demo files for
  combobox, which is good, because `registry/new-york-v4/examples/
  combobox-*.tsx` (the file naming pattern used for every other component in
  this registry, including resizable) does **not** contain `ui/combobox.tsx`'s
  demos — it's a leftover, unrelated Command+Popover pattern. The real
  demos are at `apps/v4/examples/base/combobox-*.tsx`. Anyone reaching for
  "the combobox examples" by the usual path convention will find the wrong
  files silently (they compile, they render, they're just not this
  component). See `## combobox` FINDING for the full trace.
- **Both `react-resizable-panels` and `@base-ui/react`**: the brief's Step 1
  anticipated one or both might be absent and asked to mark accordingly. In
  fact `node_modules` doesn't exist in this checkout at all — nothing under
  it for either library. Both are `derived-not-read` in full for every
  runtime-behavior claim (keyboard, ARIA, filter algorithm), not just
  partially.
- Everything else in the brief (three resizable parts, the orientation-
  inversion table request, `defaultSize` format, sidebar's six constants,
  ~25 parts, three render branches, `data-*` table, tooltip gating,
  variant/size axes, nova rules) matched what the source actually contains
  — no other corrections.

## Judgment calls made in this pass

- Resolved the resizable handle's `.cn-resizable-handle-icon` question by
  reading `bases/base/ui/resizable.tsx` as the structure tiebreak (per the
  adjudicated rule) rather than treating nova's isolated CSS rule as
  sufficient on its own — surfaced a real structural delta (empty div, no
  icon) that a CSS-only read would have missed entirely.
- Declined to guess at Base UI Combobox's filter algorithm or exact ARIA
  attribute set even under hedged language with a plausible-pattern
  fallback — stated the WAI-ARIA combobox pattern as a named hypothesis for
  Task 1–3 to confirm, not as a fact to build against, per the brief's
  explicit "do NOT guess" instruction for this specific question.
- Kept `SidebarRail`'s dead-looking `data-[state=open]:hover:` selector on
  `SidebarMenuButton` and the `data-[spacing=default]`-style class of
  upstream oddities as "port as dead weight, ledger it" rather than
  silently fixing or silently dropping, matching the precedent already set
  in `docs/superpowers/plans/2026-07-24-tier3-source-map-controls.md`'s own
  toggle-group FINDING.

# Tier 3 "wrapped" source map — sonner, drawer, carousel, input-otp

Source analysis for porting the four Tier 3 components that wrap a
third-party JS library in shadcn (sonner, vaul, embla-carousel, input-otp).
gsxui replaces each library with its own vanilla-JS module — this map's
center of gravity is therefore the **behavior contract** section per
component, not just the markup.

Inputs read: `registry/new-york-v4/ui/{sonner,drawer,carousel,input-otp}.tsx`
(byte-read, quoted verbatim below), `registry/styles/style-nova.css`
(byte-read, `.cn-toast*`/`.cn-drawer*`/`.cn-carousel*`/`.cn-input-otp*`
sections), `registry/new-york-v4/examples/{sonner,drawer,carousel,input-otp}-*.tsx`
(all 14 files, byte-read), `docs/jsx-parity.md`, `docs/component-roadmap.md`,
`ui/sheet.gsx`, `ui/dialog.js`, `ui/input-group.gsx`, `ui/input.gsx`,
`ui/command.js`, `ui/gsxui.js`, `ui/index.js`, `internal/registry/registry.go`.

`/Users/jackieli/personal/shadcn-ui/node_modules/` does **not** contain
`sonner`, `vaul`, `embla-carousel-react`, or `input-otp` — none of the four
libraries' own source was available to read. Every library-internals claim
below is `derived-not-read`: reconstructed from the shadcn wrapper's own
props/classNames/data-attribute selectors plus general knowledge of each
library's public contract, not from reading the library's source in this
repo. Marked explicitly at first use per component, not repeated on every
line.

## Legend

- `A → B (nova)` — new-york-v4 token `A` is replaced by nova token `B`; the
  class string shown elsewhere already has the nova value substituted.
- `(nova, no delta)` — nova ships a `.cn-*` entry for this part but its
  value equals new-york-v4's own; called out so "no nova entry found" and
  "nova entry confirms no change" aren't confused.
- `(no nova counterpart)` — nova's stylesheet has no `.cn-*` class for this
  part at all; new-york-v4's value is kept verbatim, unreviewed by the
  retarget.
- `derived-not-read` — reconstructed from the wrapper source / public docs
  knowledge, not from reading the library's own source (unavailable in this
  checkout's `node_modules`).
- Border→ring hairline swaps and other color/box-model changes are **not**
  adopted anywhere in this map, matching `docs/jsx-parity.md` `## nova
  density`'s NOT ADOPTED entry — nova density is metric-tokens-only policy,
  applied consistently to these four new ports even though they weren't
  part of the original 27-component retarget batch.

---

## sonner

### 1. Markup structure

**Wrapper source** (`sonner.tsx`, full file, 41 lines): a bare passthrough
— `<Toaster>` renders `<Sonner theme={...} className="toaster group"
icons={{...}} style={{--normal-bg, --normal-text, --normal-border,
--border-radius}} {...props}/>`. There is no toast-card JSX anywhere in
shadcn's own source — sonner (the library) owns 100% of the toast DOM and
its default look ships as a separate, non-Tailwind stylesheet
(`sonner/dist/styles.css`) that shadcn only re-themes via the four CSS
custom properties above (`--normal-bg` → `var(--popover)`, `--normal-text`
→ `var(--popover-foreground)`, `--normal-border` → `var(--border)`,
`--border-radius` → `var(--radius)`). The icon set is overridden per-type
via the `icons` prop: `success` = `CircleCheckIcon`, `info` = `InfoIcon`,
`warning` = `TriangleAlertIcon`, `error` = `OctagonXIcon`, `loading` =
`Loader2Icon` (spinning), all `size-4`.

**Library's actual rendered DOM** (`derived-not-read`, sonner's well-known
public structure):

```html
<section aria-label="Notifications alt+T" tabindex="-1" aria-live="polite" aria-relevant="additions text" aria-atomic="false">
  <ol data-sonner-toaster data-theme="light" data-y-position="bottom" data-x-position="right"
      style="--front-toast-height:…px; --offset:32px; --width:356px; --gap:14px;">
    <li data-sonner-toast data-styled="true" data-mounted="true" data-promise="false" data-removed="false"
        data-visible="true" data-y-position="bottom" data-x-position="right" data-index="0" data-front="true"
        data-swiping="false" data-dismissible="true" data-type="default" data-expanded="false"
        style="--index:0; --toasts-before:0; --z-index:…; --offset:0px; --initial-height:…px;" tabindex="0">
      <div data-icon>…</div>
      <div data-content>
        <div data-title>Event has been created</div>
        <div data-description>Sunday, December 03, 2023 at 9:00 AM</div>
      </div>
      <button data-button data-action tabindex="0">Undo</button>
      <button data-button data-cancel tabindex="0">Cancel</button>
      <button data-close-button aria-label="Close toast" tabindex="0"><svg>…</svg></button>
    </li>
  </ol>
</section>
```

**gsxui's toast markup — synthesized, no literal source to port** (task
explicitly calls for this: reconstruct sonner's default look as Tailwind
classes consistent with our popover/card surfaces; every class below is a
*recommendation*, not a transcription):

```html
<section aria-label="Notifications" tabindex="-1" class="sr-only:not">
  <ol data-slot="toaster" data-gsxui-toaster
      class="pointer-events-none fixed z-100 flex flex-col gap-2 p-6 bottom-0 right-0"
  ></ol>
</section>

<li data-slot="toast" data-gsxui-toast data-type="default" role="status" aria-live="polite" aria-atomic="true"
    class="pointer-events-auto relative flex w-[356px] items-start gap-3 rounded-2xl border border-border bg-popover p-4 text-sm text-popover-foreground shadow-lg
           data-[type=success]:[&>[data-icon]]:text-emerald-500 data-[type=info]:[&>[data-icon]]:text-sky-500
           data-[type=warning]:[&>[data-icon]]:text-amber-500 data-[type=error]:[&>[data-icon]]:text-destructive">
  <div data-icon class="mt-0.5 shrink-0 [&>svg]:size-4"><!-- type icon or spinning loader --></div>
  <div data-content class="flex flex-1 flex-col gap-1">
    <div data-title class="font-medium text-foreground">Event has been created</div>
    <div data-description class="text-muted-foreground">Sunday, December 03, 2023 at 9:00 AM</div>
  </div>
  <button data-action type="button" class="shrink-0 self-center text-sm font-medium underline-offset-4 hover:underline">Undo</button>
  <button data-cancel type="button" class="shrink-0 self-center text-sm text-muted-foreground underline-offset-4 hover:underline">Cancel</button>
  <button data-close-button type="button" aria-label="Close"
          class="absolute -top-1.5 -right-1.5 flex size-5 items-center justify-center rounded-full border border-border bg-background text-foreground shadow-sm">
    <!-- icon.X, size-3 -->
  </button>
</li>
```

`rounded-2xl` on the card is nova's one contribution (`.cn-toast { rounded-2xl }`,
`style-nova.css` line 1245) — new-york-v4's own radius is CSS-var-driven
(`--border-radius: var(--radius)`, ≈ `rounded-lg`/`rounded-xl` visually
depending on the theme's `--radius`). Since our module is plain Tailwind
classes, not sonner's own CSS-var-consuming stylesheet, this is a genuine
simplification (WIN): no `--normal-bg`/`--normal-text`/`--border-radius`
indirection to reproduce at all — `bg-popover`/`text-popover-foreground`/
`rounded-2xl` are hardcoded classes, matching our popover/card surface
convention directly. All type-colored icon tints (`emerald`/`sky`/`amber`/
`destructive`) are a house choice, not from either source — flag for
review; a monochrome-icon alternative (all `text-foreground`, matching
new-york-v4's actual unthemed default before `richColors` is turned on) is
equally defensible and closer to shadcn's literal default (shadcn's own
`sonner.tsx` does **not** pass `richColors`, so type-tinted icons are not
actually its default look — the recommendation above intentionally departs
from strict new-york-v4 parity toward the more legible tinted-icon
convention seen on ui.shadcn.com's live nova site; ledger as a deliberate
choice for the planner to confirm, not silent scope creep).

Icons needed all exist in `ui/icon`: `circle-check`, `info`, `triangle-alert`,
`octagon-x`, `loader-circle` (verified present in `ui/icon/icon_data.go`).
Because the toast `<li>` is built by client JS (see §3), these icon paths
must be duplicated as literal SVG-path strings inside `ui/sonner.js` — the
same "verbatim port of static data into a JS module" precedent as
`command.js`'s hand-ported `commandScore` — not `icon.CircleCheck` Go calls
(those only exist server-side).

### 2. Behavior contract

**Stacking model** (`derived-not-read`): default position `bottom-right`.
Only the front (most recent) toast is fully visible; older toasts stack
behind it, each progressively scaled down and offset upward, driven by CSS
custom properties the library maintains per toast (`--index`,
`--toasts-before`, `--z-index`, `--offset`, `--initial-height`) plus one
container-level `--front-toast-height` (used to compute the collapsed
stack's total footprint). Sonner shows at most **3** stacked toasts
unexpanded (its `VISIBLE_TOASTS_AMOUNT` constant) — a 4th+ toast is queued
and only promoted into the visible stack once an older one is dismissed.
Hovering the toaster region sets `data-expanded="true"` on the toaster,
which separates every toast to its full height with `--gap` (14px) between
them instead of the collapsed scale/offset stack; un-hovering collapses
back.

**Timers**: default toast duration is **4000ms**. Hovering an individual
toast (or the toaster region generally) pauses its dismiss timer; leaving
resumes it with the remaining time. A `loading`-type toast (from
`toast.promise`) has no duration until the promise settles. Toasts are also
swipe-dismissible on touch (a drag past a threshold sets
`data-swipe-out="true"` and removes it) — not relevant to gsxui's v1 (no
gesture layer recommended, see §3).

**Enter/exit animation** (`derived-not-read`, approximate): enter is a
translate+fade from the position-appropriate edge (`bottom` position slides
up+fades in), roughly 350ms with an ease-out curve. Exit sets
`data-removed="true"`, which plays a translate+fade(+scale, on swipe) exit,
then the library removes the node from the DOM after the animation
duration elapses (either via an `animationend`-style listener or a matched
`setTimeout`). gsxui's module should use the exact discrete-transition
architecture already established in `docs/jsx-parity.md` `## animations`
(closed-state base classes + `open:`/`data-[state=open]` + `starting:open:`
+ `transition-discrete`) rather than porting sonner's own animation timing
values verbatim — there is no literal source to port those values from
(non-Tailwind stylesheet), so this is a fresh design following house
convention, not a token transcription.

**Promise/loading API** (from `sonner-types.tsx`, byte-read): the docs
exercise `toast(msg)`, `.success(msg)`, `.info(msg)`, `.warning(msg)`,
`.error(msg)`, and `.promise(fn, {loading, success, error})`. The promise
variant **morphs the same toast node in place** — same DOM element, same
position in the stack, no re-animation — swapping icon + message from
`data-type="loading"` (spinner) to `success`/`error` on settle. This
morph-in-place behavior is the one signature sonner UX worth explicitly
replicating; a naive dismiss-old/spawn-new implementation would visibly
reflow the stack instead.

**Recommended `gsxui.toast(...)` API shape** (a server-rendered library has
no meaningful "declarative toast markup" — a toast is definitionally a
client-triggered response to some JS event, e.g. a fetch resolving, so the
primary surface should stay a plain JS call, matching `toast()`'s own
shape almost exactly):

```js
import { toast } from "gsxui"; // re-exported through ui/index.js, see §3
toast("Event has been created", { description: "…", action: { label: "Undo", onClick: fn } });
toast.success(msg, opts); toast.info(msg, opts); toast.warning(msg, opts); toast.error(msg, opts); toast.loading(msg, opts);
const id = toast.promise(promiseOrFn, { loading: "…", success: (v) => `…`, error: "…" });
toast.dismiss(id); // omit id to dismiss all
```

Plus one **declarative trigger** for zero-JS demo/doc pages, mirroring the
established `data-gsxui-dialog-trigger` idiom (`docs/jsx-parity.md`
`## dialog` MECHANISM): a click-delegated
`data-gsxui-toast="message" data-gsxui-toast-description="…" data-gsxui-toast-type="success"`
attribute set on any element (button, link) that calls the same internal
`show()` the imperative API uses — covers the `sonner-demo.tsx`/
`sonner-types.tsx` docs pages without a page-specific `<script>` block.

### 3. gsxui adaptation notes

- **Naming constraint** (binding, from `internal/registry/registry.go`
  `HasJS`): a component's behavior JS is discovered by exact filename match
  — `ui/<name>.js` where `<name>` is the `.gsx` file's basename. If the gsx
  file is `ui/sonner.gsx` (matching the roadmap's own component name
  "sonner"), the JS module **must** be `ui/sonner.js`, not `ui/toast.js` —
  a `toast.js` file would silently never be detected as sonner's companion
  JS by `registry.HasJS("sonner")`/the CLI vendoring path. Corrects the
  task brief's own suggested filename.
- `ui.Toaster` component: renders the empty, always-present positioned
  region shown in §1 (the `<section>`+`<ol data-gsxui-toaster>` pair) — one
  instance per page, typically placed once in a site's root layout, same
  placement convention as shadcn's own `<Toaster/>` in `app/layout.tsx`.
  Ships **only** the default `bottom-right` position for v1 — the other 5
  sonner positions (`top-left`, `top-center`, `top-right`, `bottom-left`,
  `bottom-center`) are a real feature reduction worth ledgering as a v1
  gap, but none of the docs demos (`sonner-demo.tsx`, `sonner-types.tsx`)
  exercise anything but the default, so it's a safe cut.
- `ui/sonner.js` is a **new JS-authoring shape** relative to every other
  gsxui behavior module: `dialog.js`/`command.js`/`dropdown.js` etc. all
  operate on server-rendered markup already in the DOM (delegated
  attribute-driven behavior over static HTML). Sonner has no server-side
  markup to attach to *for the toast itself* — the JS must **construct**
  each `<li data-slot="toast">` element from scratch on every `toast()`
  call (template-literal or `document.createElement` tree), append it into
  the one static `<ol data-gsxui-toaster>` region, and manage its own
  lifecycle (mount → timer → dismiss/remove). Flag this prominently as the
  single largest new-architecture item across all four components — every
  other Tier 3 port (drawer, carousel, input-otp) still fits the
  "delegated behavior over static server markup" shape; sonner does not.
- Icon SVG paths (`circle-check`/`info`/`triangle-alert`/`octagon-x`/
  `loader-circle`) must be hand-copied as literal path-data strings into
  `ui/sonner.js` (mirrored from `ui/icon/icon_data.go`, verified present) —
  a maintenance seam: if `ui/icon`'s glyphs are ever regenerated from a
  newer Lucide version, this module's copies won't update automatically.
- Barrel wiring: `import "./sonner.js"` in `ui/index.js` (alongside the
  other 9 behavior imports already there), **plus** `export { toast } from
  "./sonner.js";` so `import { toast } from "gsxui"` works through the
  barrel — establishes the first public-imperative-API-via-barrel-export
  precedent in this codebase (every existing module only exports
  internals like `requestClose` for sibling-module use, not for page
  authors).
- Stacking/timer/pause-on-hover mechanics (§2) are real, non-trivial state
  machine work — recommend implementing against a plain array of toast
  records (`{id, el, type, timer, pausedAt}`) rather than trying to encode
  the whole thing in CSS custom properties the way sonner itself does;
  gsxui doesn't need sonner's own CSS-var-driven stacking technique since
  we're not shipping a themeable third-party stylesheet, just fixed
  Tailwind classes recomputed per toast via inline `style` for the
  scale/translate stack values.

### 4. Demo inventory

- `sonner-demo.tsx` — single outline button, `toast(msg, {description, action})`.
  Exercises: title+description+action button.
- `sonner-types.tsx` — 6 buttons: default/success/info/warning/error/promise.
  Exercises: every type variant + the promise morph-in-place behavior.

Recommend both as the site examples — together they cover the full type
surface and the action-button + promise API without needing a 3rd/4th demo.

---

## drawer

### 1. Markup structure

**shadcn source** (`drawer.tsx`, full file, 136 lines, byte-read): a bare
`vaul` passthrough, structurally identical in shape to `dialog.tsx`/
`sheet.tsx` — `Drawer`→`DrawerPrimitive.Root`, `DrawerTrigger`, `DrawerPortal`,
`DrawerClose`, `DrawerOverlay` (`fixed inset-0 z-50 bg-black/50
data-[state=closed]:animate-out data-[state=closed]:fade-out-0
data-[state=open]:animate-in data-[state=open]:fade-in-0`),
`DrawerContent`, `DrawerHeader`, `DrawerFooter`, `DrawerTitle`,
`DrawerDescription`. `DrawerContent`'s class string (quoted exactly):

```
group/drawer-content fixed z-50 flex h-auto flex-col bg-background
data-[vaul-drawer-direction=top]:inset-x-0 data-[vaul-drawer-direction=top]:top-0 data-[vaul-drawer-direction=top]:mb-24 data-[vaul-drawer-direction=top]:max-h-[80vh] data-[vaul-drawer-direction=top]:rounded-b-lg data-[vaul-drawer-direction=top]:border-b
data-[vaul-drawer-direction=bottom]:inset-x-0 data-[vaul-drawer-direction=bottom]:bottom-0 data-[vaul-drawer-direction=bottom]:mt-24 data-[vaul-drawer-direction=bottom]:max-h-[80vh] data-[vaul-drawer-direction=bottom]:rounded-t-lg data-[vaul-drawer-direction=bottom]:border-t
data-[vaul-drawer-direction=right]:inset-y-0 data-[vaul-drawer-direction=right]:right-0 data-[vaul-drawer-direction=right]:w-3/4 data-[vaul-drawer-direction=right]:border-l data-[vaul-drawer-direction=right]:sm:max-w-sm
data-[vaul-drawer-direction=left]:inset-y-0 data-[vaul-drawer-direction=left]:left-0 data-[vaul-drawer-direction=left]:w-3/4 data-[vaul-drawer-direction=left]:border-r data-[vaul-drawer-direction=left]:sm:max-w-sm
```

Important asymmetry in new-york-v4: **only top/bottom get rounded
corners** (`rounded-b-lg`/`rounded-t-lg`) — left/right have `border-l`/
`border-r` only, no rounding at all. Nova changes this (see nova deltas
below).

The drag-handle bar is an unconditional inline `<div>`, not a named
sub-component: `<div class="mx-auto mt-4 hidden h-2 w-[100px] shrink-0
rounded-full bg-muted group-data-[vaul-drawer-direction=bottom]/drawer-content:block" />`
— visible **only** for the bottom variant (the mobile-sheet convention),
hidden (`hidden`, overridden to `block` only under the bottom-direction
group selector) for top/left/right.

`DrawerHeader`: `flex flex-col gap-0.5 p-4
group-data-[vaul-drawer-direction=bottom]/drawer-content:text-center
group-data-[vaul-drawer-direction=top]/drawer-content:text-center
md:gap-1.5 md:text-left` (centered for bottom/top, left-aligned at `md`+
for left/right — direction-conditional text alignment, not present in
Sheet at all). `DrawerFooter`: `mt-auto flex flex-col gap-2 p-4` (byte-
identical to `SheetFooter`'s own). `DrawerTitle`: `font-semibold
text-foreground`. `DrawerDescription`: `text-sm text-muted-foreground`.

**Nova deltas** (`.cn-drawer-content` etc., `style-nova.css` lines 492–535,
byte-read):

| part | new-york-v4 | nova | delta |
|---|---|---|---|
| content radius (top/bottom) | `rounded-b-lg` / `rounded-t-lg` | `rounded-b-xl` / `rounded-t-xl` | `lg→xl (nova)`, matches the global overlay radius bump (`docs/…/2026-07-24-nova-density-map.md` "radius scale shifts up one step") |
| content radius (left/right) | none | `rounded-r-xl` (left) / `rounded-l-xl` (right) | **new** — nova rounds the free (non-anchored) edge on left/right too, a markup addition, not present in new-york-v4 at all |
| content background | `bg-background` | `bg-popover text-popover-foreground` | color-scope, **NOT ADOPTED** (see Legend) |
| overlay | `bg-black/50` | `bg-black/10 supports-backdrop-filter:backdrop-blur-xs` | color/blur-scope, **NOT ADOPTED** for the base opacity value; the `supports-backdrop-filter:backdrop-blur-xs` addition mirrors what `ui/sheet.gsx` already carries on its own `backdrop:` pseudo (`sheet.gsx` line 166) — recommend porting *that* utility onto drawer's own backdrop for consistency with sheet, independent of the nova opacity question |
| handle bar height | `h-2` | `h-1` | `h-2→h-1 (nova)` — thinner handle |
| header gap at `md`+ | `md:gap-1.5` | `md:gap-0.5` | nova drops the responsive gap bump; stays `gap-0.5` at every breakpoint |
| title weight | `font-semibold` | `font-medium` | `font-semibold→font-medium (nova)`, matches the general title de-emphasis pattern seen in dialog/card/sheet's own nova entries |
| footer, description | — | — | `(nova, no delta)` — both byte-identical to new-york-v4's own |

### 2. Behavior contract (vaul)

`derived-not-read` throughout this subsection — vaul is not in
`node_modules`; reconstructed from `drawer.tsx`'s own data-attribute
selectors plus general knowledge of vaul's public contract.

- **Data attributes vaul stamps**: `data-vaul-drawer` (root content
  element), `data-vaul-drawer-direction` (`top`|`bottom`|`left`|`right`,
  the one selector `drawer.tsx` itself keys every conditional class off
  of), `data-vaul-drawer-visible`, `data-vaul-overlay`, `data-vaul-handle`,
  `data-vaul-no-drag` (opt-out region for interactive children, e.g. a
  button inside the drawer that shouldn't start a drag), `data-vaul-
  snap-points` / `data-vaul-delayed-snap-points` (snap-point mode only).
- **Drag-to-dismiss gesture**: pointerdown/pointermove/pointerup on the
  content (or a designated handle element) tracks a drag delta and applies
  a live `translate3d()` transform following the pointer; releasing past a
  velocity/distance threshold completes the dismiss (continues the
  transform off-screen, ~500ms cubic-bezier-ish easing, `derived-not-read`
  on the exact curve/duration), releasing short of the threshold springs
  back to `translate3d(0,0,0)`.
- **Background scaling**: an opt-in prop, `shouldScaleBackground` on
  vaul's `Drawer.Root`. shadcn's own `drawer.tsx` `Drawer` component does
  **not** set it (`<DrawerPrimitive.Root data-slot="drawer" {...props}/>`,
  no default injected) and vaul's own library default is `false` — so
  **none of shadcn's docs demos get background scaling**. Skip entirely
  for v1; there's nothing to replicate that the actual demos exercise.
- **Snap points**: opt-in via a `snapPoints` prop, not used in
  `drawer-demo.tsx` or `drawer-dialog.tsx`. Recommend skip for v1 —
  consistent with the roadmap's own drag-gesture cut (snap points are
  meaningless without drag).
- **Direction variants**: `top`/`bottom`/`left`/`right`, same four as
  `Sheet`'s `side`. Drawer's conventional **default anchor is `bottom`**
  (the mobile bottom-sheet pattern the handle-bar visibility rule and both
  docs demos assume) — a real behavioral difference from `Sheet`'s own
  default of `right`.
- **Animation values**: irrelevant to port literally — gsxui replaces
  vaul's live-transform drag physics entirely with the `<dialog>` +
  Tailwind-keyframe slide-in/out architecture already solved by
  `ui/sheet.gsx` (nova 200ms duration-200 ease-in-out both directions, per
  `docs/jsx-parity.md` `## sheet`). Recommend reusing sheet's exact timing
  values, not vaul's own.

### 3. gsxui adaptation notes

Recommend a **new** `ui/drawer.gsx`, following `ui/sheet.gsx`'s exact
composition pattern (`Sheet`/`SheetTrigger`/`SheetContent` → `Drawer`/
`DrawerTrigger`/`DrawerContent`), not a shared/aliased file — there's
enough real divergence (default direction, rounded corners, handle bar,
`max-h-[80vh]` on top/bottom vs `sm:max-w-sm` on left/right) to warrant its
own component, matching the roadmap's own Tier 3 listing of `drawer` as a
distinct entry from Tier 2's `sheet`.

- `Drawer` composes `ui.Dialog` directly (`data-slot` override
  `"dialog"`→`"drawer"`), identical mechanism to `Sheet`/`AlertDialog` —
  `drawer → dialog` derives, `ui/dialog.js` is pulled in transitively, zero
  new dialog.js code needed (same conclusion as `## sheet`'s own Registry
  note: `HasJS("drawer")` should be `false`).
- `DrawerTrigger` renders its own `<button data-gsxui-dialog-trigger>`,
  same button-in-button-avoidance reasoning as `SheetTrigger`/
  `DialogTrigger` (`docs/jsx-parity.md` `## dialog` FINDING).
- `DrawerContent` does **not** compose `ui.DialogContent` or `ui.SheetContent`
  — same reasoning as `## sheet`'s own ADAPT (the centered-card recipe and
  the side-anchored recipe target the same CSS properties with materially
  different values; no single class string merges them). Renders its own
  `<dialog data-slot="drawer-content" data-gsxui-dialog-content
  data-state="closed" data-side="...">`.
- **Naming decision (flag for planner)**: shadcn/vaul's own attribute is
  `data-vaul-drawer-direction`; gsxui has no vaul underneath to key off of,
  and (mirroring sheet's own MECHANISM note) direction selection happens
  server-side via a Go `switch` in `class={}`, not a client CSS attribute
  selector — the stamped attribute is decorative only. Recommend a `direction`
  Go param (matching drawer/vaul's own vocabulary, distinct from Sheet's
  `side`) that both selects the class-string case and stamps `data-side`
  (reusing the same internal attribute name `Sheet` already uses, for any
  future shared tooling/CSS that keys off `data-side` generically) rather
  than inventing a second attribute name. Either choice is defensible;
  ledger as an explicit decision point.
- **Recommended per-direction class strings** — built by applying
  `ui/sheet.gsx`'s own already-solved `<dialog>`-vs-UA-defaults fixes
  (`m-0`, opposite-edge `-auto`, `open:flex`, `text-foreground`) to
  `drawer.tsx`'s values, plus the nova deltas from §1. This is a
  **recommendation**, not a verified-in-browser render — sheet's own six
  ADAPTs were only found by rendering in a real tab (`## sheet`'s own
  three-part fix note), so the planner's implementation task must repeat
  that verification pass for drawer, not assume these strings are correct
  by inspection alone:

  Base (shared): `fixed z-50 m-0 open:flex flex-col gap-4 bg-popover text-popover-foreground text-sm shadow-lg transition ease-in-out duration-200 data-[state=closed]:animate-out data-[state=open]:animate-in backdrop:bg-black/10 backdrop:duration-200 supports-backdrop-filter:backdrop:backdrop-blur-xs data-[state=open]:backdrop:animate-in data-[state=open]:backdrop:fade-in-0 data-[state=closed]:backdrop:animate-out data-[state=closed]:backdrop:fade-out-0`
  (using nova's `bg-popover text-popover-foreground` here despite the
  Legend's default NOT-ADOPTED stance — flag this as the one deliberate
  exception: unlike the retargeted-27 components, drawer has no prior
  new-york-v4-based gsxui version to stay consistent with, so there's no
  "existing bg-background baseline" to preserve; recommend nova's color
  here specifically, planner should confirm)

  - `bottom` (default): `inset-x-0 bottom-0 top-auto w-full max-w-none h-auto mt-24 max-h-[80vh] rounded-t-xl border-t data-[state=closed]:slide-out-to-bottom data-[state=open]:slide-in-from-bottom` + handle bar visible
  - `top`: `inset-x-0 top-0 bottom-auto w-full max-w-none h-auto mb-24 max-h-[80vh] rounded-b-xl border-b data-[state=closed]:slide-out-to-top data-[state=open]:slide-in-from-top` + handle bar hidden
  - `left`: `inset-y-0 left-0 right-auto h-full max-h-none w-3/4 rounded-r-xl border-r sm:max-w-sm data-[state=closed]:slide-out-to-left data-[state=open]:slide-in-from-left` + handle bar hidden
  - `right`: `inset-y-0 right-0 left-auto h-full max-h-none w-3/4 rounded-l-xl border-l sm:max-w-sm data-[state=closed]:slide-out-to-right data-[state=open]:slide-in-from-right` + handle bar hidden

  Note top/bottom use their own explicit `max-h-[80vh]` (author-origin,
  already beats Chrome's UA `max-height` safety net per sheet's own ADAPT
  reasoning — "author-origin ALWAYS beats UA-origin regardless of value")
  so they do **not** need sheet's `max-h-none` escape hatch; left/right
  (using `h-full` like sheet's left/right) do need `max-h-none`, copied
  verbatim from sheet's fix.

- **Handle bar as decoration**: v1 has no drag gesture, so rendering the
  handle bar (a visual "drag me" affordance) while dragging doesn't
  actually work is a real, if minor, UX mismatch — flag as an accepted GAP
  rather than silently shipping it or silently dropping it. Recommend
  keeping it (visual parity with the nova/new-york reference), ledgered.
- Registry: `Deps("drawer") == ["dialog"]`, `HasJS("drawer") == false` — no
  new JS module, matching sheet's own shape exactly.

### 4. Demo inventory

- `drawer-demo.tsx` — goal-counter drawer with a `recharts` `BarChart`
  sparkline. **Adaptation needed**: gsxui has no charting story (roadmap:
  "chart — defer until demanded"). Recommend a simplified stand-in that
  keeps the goal counter +/- buttons and drops the bar chart, or swaps it
  for a static decorative bar row (plain divs) — not a literal port.
- `drawer-dialog.tsx` — responsive pattern: renders `Dialog` on desktop /
  `Drawer` on mobile via a `useMediaQuery("(min-width: 768px)")` hook. No
  gsxui equivalent exists for a JS media-query-driven component swap
  (`docs/component-roadmap.md` has no such primitive). Recommend as a
  **patterns/pages-phase** example rather than a base-component demo — flag
  as out of scope for the drawer component task itself, worth a follow-up
  note in the roadmap's patterns section.

Recommend the adapted `drawer-demo` as the primary (only mandatory) site
example; `drawer-dialog`'s responsive-swap pattern deferred.

---

## carousel

### 1. Markup structure

**shadcn source** (`carousel.tsx`, full file, 242 lines, byte-read):
`Carousel` (root `<div role="region" aria-roledescription="carousel"
data-slot="carousel" class="relative">`, `onKeyDownCapture` wiring
ArrowLeft/ArrowRight to `scrollPrev`/`scrollNext`), `CarouselContent`
(**two** nested divs: outer `<div ref={carouselRef} class="overflow-hidden"
data-slot="carousel-content">` — embla's viewport ref target — wrapping
inner `<div class="flex -ml-4">` horizontal / `<div class="flex -mt-4
flex-col">` vertical — the actual flex track), `CarouselItem` (`<div
role="group" aria-roledescription="slide" data-slot="carousel-item"
class="min-w-0 shrink-0 grow-0 basis-full pl-4">` horizontal / `pt-4`
vertical), `CarouselPrevious`/`CarouselNext` (compose `Button
variant="outline" size="icon"`, `absolute size-8 rounded-full`, horizontal:
`top-1/2 -left-12 -translate-y-1/2` / `-right-12`, vertical: `-top-12
left-1/2 -translate-x-1/2 rotate-90` / `-bottom-12 …`, `disabled={!canScrollPrev}`
/ `!canScrollNext`, `ArrowLeft`/`ArrowRight` icon + `sr-only` label).

**Nova**: `.cn-carousel-previous`/`.cn-carousel-next { rounded-full }`
(`style-nova.css` lines 269–275) — **no delta**. New-york-v4's own class
already carries `rounded-full`/`size-8`; nova's entry is a confirming
no-op, not a change. `(nova, no delta)` for every carousel part — nova has
no other carousel entries.

### 2. Behavior contract (embla)

`derived-not-read` — `embla-carousel-react` not in `node_modules`;
reconstructed from `carousel.tsx`'s API surface plus general knowledge of
embla's public contract.

**Fundamental architecture difference to flag first**: embla is a
**JS-transform carousel**, not a native-scroll one — it does not use CSS
`scroll-snap` or native `overflow-x: auto` at all. It applies
`transform: translate3d(Xpx, 0, 0)` to the flex track div directly, driven
by its own drag-physics engine (pointer events, momentum, rubber-banding
at the ends). gsxui's roadmap decision replaces this with real native CSS
scroll-snap (`overflow-x: auto; scroll-snap-type: x mandatory` on the
viewport div, `scroll-snap-align: start` on each item) + JS only for the
prev/next buttons and disabled-state/current-index bookkeeping. This is a
genuine behavior substitution: native scroll gets touch/trackpad momentum
and rubber-banding for free (a WIN, not a gap), but loses embla's `loop`
(infinite wraparound) feature, which has no clean native-scroll analog
(embla's own loop mode works by literally cloning slides at both ends —
reproducible, but real extra work, and *none* of the docs demos below
request it).

**Embla features exercised by the docs demos** (all 6 example files under
`registry/new-york-v4/examples/carousel-*.tsx`, byte-read):

| file | opts / feature exercised |
|---|---|
| `carousel-demo.tsx` | no `opts` at all — embla defaults, single item per view (`basis-full`), horizontal, prev/next with disabled-state at the ends |
| `carousel-api.tsx` | `setApi`, subscribes to embla's `'select'` event, reads `api.selectedScrollSnap()+1` / `api.scrollSnapList().length` to render a live "Slide X of Y" indicator |
| `carousel-orientation.tsx` | `orientation="vertical"`, `opts={{align:"start"}}`, `CarouselContent` gets `h-[200px] -mt-1`, item `pt-1 md:basis-1/2` (2-per-view at `md`+) |
| `carousel-plugin.tsx` | `embla-carousel-autoplay` plugin (`delay: 2000, stopOnInteraction: true`), plus explicit `onMouseEnter={plugin.stop}` / `onMouseLeave={plugin.reset}` — hover pause/resume is userland code in the demo, not an embla-native feature |
| `carousel-size.tsx` | `opts={{align:"start"}}`, item `md:basis-1/2 lg:basis-1/3` — 2/3-per-view via responsive `basis-*` fractions |
| `carousel-spacing.tsx` | `CarouselContent` gets `-ml-1` override (default is `-ml-4`), item `pl-1` (default `pl-4`) — proves gap is **entirely** caller-controlled via the existing class-merge mechanism, no separate `spacing` prop exists or is needed |

`loop` mode is not exercised by any docs demo — safe to ledger as an
accepted v1 gap per the roadmap's own framing ("snap covers the docs
demos"). `align: "center"` vs `"start"` — embla's own default alignment is
commonly documented as `"start"`; `carousel-demo.tsx`'s single-per-view
case makes align mostly moot (one slide always fills the viewport), so this
detail matters only for `-orientation`/`-size`/`-spacing`, all three of
which explicitly pass `align: "start"` — **native `scroll-snap-align:
start` on each item is therefore sufficient for every demo that cares**,
no `center` alignment mode needs porting.

**Prev/Next disabled logic**: embla computes `canScrollPrev()`/
`canScrollNext()` from internal scroll-progress/edge state, re-evaluated on
its own `'select'`/`'reInit'` events. Native-scroll equivalent: track the
viewport's `scrollLeft` (or `scrollTop`, vertical) against `0` and against
`scrollWidth - clientWidth` (with a small epsilon for sub-pixel rounding),
recomputed on a throttled `scroll` listener and on `ResizeObserver` (layout
changes, e.g. responsive `basis-*` breakpoints changing how many slides
fit) — set directly as the native `disabled` attribute on the prev/next
`<button>`s (Button's own `disabled:` classes then apply for free, no
extra CSS needed).

### 3. gsxui adaptation notes

- **Parts**: `Carousel` (root `<div role="region" aria-roledescription="carousel"
  data-slot="carousel" data-gsxui-carousel data-orientation="horizontal|vertical">`,
  stamped for JS proximity lookup, same `closest()` pattern as
  `data-gsxui-dialog`/`data-gsxui-command`), `CarouselContent` (renders
  **both** divs from the source — outer becomes the real scroll container:
  add `overflow-x-auto snap-x snap-mandatory [scrollbar-width:none]
  [&::-webkit-scrollbar]:hidden` horizontal / `overflow-y-auto snap-y
  snap-mandatory …` vertical — the scrollbar-hiding utility pair is **new**,
  not present in either source, needed only because we're now on a real
  native scroll container that would otherwise show a visible scrollbar
  embla's transform-based approach never had), `CarouselItem` (adds
  `snap-start` — also **new**, not in shadcn's source, required for native
  scroll-snap to have any snap points at all), `CarouselPrevious`/
  `CarouselNext` (compose `ui.Button`, same classes as source, `data-gsxui-carousel-prev`/
  `-next` for delegated click wiring).
- **Prev/Next scroll amount**: scroll by **one item's measured width**
  (`item.getBoundingClientRect().width` + the item's own gap, i.e. embla's
  default `slidesToScroll: 1` behavior), not by one full viewport width —
  a viewport-width scroll would visibly skip slides whenever more than one
  slide is visible per view (the `-size`/`-orientation` demos).
- **`ui/carousel.js` contract**:
  - delegated `click` on `[data-gsxui-carousel-prev]`/`[data-gsxui-carousel-next]`
    → `viewport.scrollBy({left/top: ±itemWidth, behavior: "smooth"})`.
  - `scroll` listener on the viewport (rAF-throttled) recomputing
    prev/next `disabled` + emitting `gsxui:carousel-select` with
    `{index, count}` detail on the carousel root, and stamping
    `data-current-index` on the root for CSS-only dot-indicator styling
    (`[data-index="N"]` on a caller-authored dot list, no JS required on
    the consumer side for the common "row of dots" indicator pattern).
  - `ResizeObserver` on the viewport, recomputing disabled-state on
    layout change (covers responsive `basis-*` breakpoint changes).
  - keyboard: `ArrowLeft`/`ArrowRight` when focus is within
    `[data-gsxui-carousel]` (delegated `keydown`, `preventDefault`),
    mirroring the source's `onKeyDownCapture` on the root div.
  - optional bespoke autoplay: `data-gsxui-carousel-autoplay="2000"` on
    the root (ms interval) → `setInterval` calling next(), paused on
    `pointerenter`/`focusin` within the carousel, resumed on
    `pointerleave`/`focusout` — reproduces `carousel-plugin.tsx`'s actual
    demo behavior (explicit hover-pause/resume, not embla's
    `stopOnInteraction` semantics, which trigger on drag/click, not hover)
    without porting embla's plugin system.
- **API surface reduction (deliberate, ledger it)**: embla's `CarouselApi`
  (`scrollTo(i)`, `scrollPrev()`, `scrollNext()`, `canScrollPrev()`,
  `canScrollNext()`, `selectedScrollSnap()`, `scrollSnapList()`,
  `on(event, cb)`, `off()`, plus the whole plugin-extension mechanism) is
  not reproduced 1:1 — there's no plugin ecosystem to support here.
  Recommend a much smaller surface: the `gsxui:carousel-select` CustomEvent
  (covers `carousel-api.tsx`'s "Slide X of Y" indicator use case
  completely) plus a couple of methods stashed directly on the DOM node
  (`carouselEl.gsxuiCarousel = { scrollTo(i), next(), prev() }`) for any
  script that needs imperative control. Embla's full plugin system
  (Autoplay, ClassNames, Fade, WheelGesture, …) is not ported at all; only
  the one bespoke `data-gsxui-carousel-autoplay` attribute stands in for
  the single plugin the docs demos actually use.
- Registry: `CarouselPrevious`/`CarouselNext` compose `ui.Button` →
  `Deps("carousel") == ["button"]`. `HasJS("carousel") == true` (real new
  interactive JS, unlike sheet/alert-dialog/drawer).

### 4. Demo inventory

Recommend 3–4 site examples:

- `carousel-demo.tsx` — basic single-item, default opts. Baseline case.
- `carousel-orientation.tsx` — vertical axis + 2-per-view via `md:basis-1/2`.
  Exercises the vertical scroll-snap axis and multi-per-view together.
- `carousel-api.tsx` — current-slide indicator. Directly exercises the
  `gsxui:carousel-select` CustomEvent contract — good proof case, worth
  keeping even though it needs a small inline script on the demo page.
- `carousel-size.tsx` or `carousel-spacing.tsx` (pick one; near-duplicates)
  — multi-per-view responsive breakpoints (`md:basis-1/2 lg:basis-1/3`).

`carousel-plugin.tsx` (autoplay) is a good stretch 5th example once the
`data-gsxui-carousel-autoplay` attribute exists, but depends on that new
surface landing first — not required for the initial 3–4.

---

## input-otp

### 1. Markup structure

**shadcn source** (`input-otp.tsx`, full file, 78 lines, byte-read):
`InputOTP` (thin wrapper over the library's `<OTPInput data-slot="input-otp"
containerClassName="flex items-center gap-2 has-disabled:opacity-50"
className="disabled:cursor-not-allowed" {...props}/>`), `InputOTPGroup`
(`<div data-slot="input-otp-group" class="flex items-center">`),
`InputOTPSlot` (reads `OTPInputContext`'s `slots[index]` →
`{char, hasFakeCaret, isActive}`; renders):

```
relative flex h-9 w-9 items-center justify-center border-y border-r border-input text-sm shadow-xs
transition-all outline-none first:rounded-l-md first:border-l last:rounded-r-md
aria-invalid:border-destructive
data-[active=true]:z-10 data-[active=true]:border-ring data-[active=true]:ring-[3px] data-[active=true]:ring-ring/50
data-[active=true]:aria-invalid:border-destructive data-[active=true]:aria-invalid:ring-destructive/20
dark:bg-input/30 dark:data-[active=true]:aria-invalid:ring-destructive/40
```

— containing `{char}` and, when `hasFakeCaret`, a centered blinking-caret
overlay: `<div class="pointer-events-none absolute inset-0 flex items-center
justify-center"><div class="h-4 w-px animate-caret-blink bg-foreground
duration-1000" /></div>`. `InputOTPSeparator`: `<div data-slot="input-otp-separator"
role="separator"><MinusIcon/></div>`, no class of its own.

**Library's actual architecture** (`derived-not-read`, input-otp's
well-known small public contract, ~200 lines upstream): **one single real
`<input>`**, absolutely positioned to cover the whole slots row, made
visually invisible while remaining the actual focus/caret/paste owner (NOT
`display:none`/`visibility:hidden`, which would break focusability — an
opacity/transparent-text technique instead). The "slots" rendered by
`InputOTPGroup`/`InputOTPSlot` above are **purely presentational**, driven
entirely by `OTPInputContext`, which the library computes from the one
real input's live `value` + `selectionStart`/`selectionEnd` on every
`input`/`selectionchange`/`focus`/`blur` event.

**The single most important architectural fact, worth stating plainly**:
this is emphatically **not** N separate `<input>` elements with manual
focus-advance-on-keypress JS (the common naive reimplementation of an OTP
field). It is one input; native browser text-editing — typing advances the
native text cursor, Backspace deletes-and-moves-back, arrow keys move the
native selection — supplies focus-advance, backspace-across-slots, and
arrow navigation **entirely for free**, no JS required for any of it. The
library's JS exists only to (a) compute the per-index visual projection
(`char`/`isActive`/`hasFakeCaret`) from the real input's value+selection on
every relevant event, (b) keep that real input visually hidden but
interaction-capable, and (c) route a click on a specific visual slot to
`input.focus(); input.setSelectionRange(i, i)` so clicking "slot 3" moves
the real caret there.

**Paste handling** (`derived-not-read`): native paste lands in the real
input's value as normal browser behavior; the library's `input`/`onChange`
handler then strips characters not matching `pattern` and truncates to
`maxLength`, redistributing the resulting string across slot positions
starting from the paste's insertion point (or index 0 for an empty field)
— this is nothing more than "the input's value is now this string, recompute
every slot's `char` from `value[i]`," not a special paste-specific code
path in principle.

**Pattern/maxLength**: `maxLength` sets `<input maxlength>` — the caller
must render exactly that many `InputOTPSlot`s themselves; nothing
auto-generates slot count from `maxLength`. `pattern` (a regex string, e.g.
`REGEXP_ONLY_DIGITS_AND_CHARS` from `input-otp-pattern.tsx`) does two
things: (a) live keystroke filtering — reject/strip a typed or pasted
character that doesn't match before it commits to the input's value, and
(b) sets the native HTML5 `pattern` attribute for native form validity.

**Form association**: because it's a real `<input>` (with a `name`
attribute like any other form field), it participates in native `<form>`
submission directly — no separate hidden mirror input is needed (unlike
some hand-rolled "N-input" OTP components that keep a decorative fake
front-end plus one real hidden field). The same element is both the actual
form value **and** the interaction surface.

### 2. Behavior contract — summary for the gsxui port

- Hidden-input architecture (above) is the entire mechanism; no per-slot
  focus JS at all.
- Fake caret: rendered only inside the currently-active slot, only when
  that slot's own character is empty (i.e. the caret shows where the
  *next* typed character will land), only while the real input has focus
  and a collapsed selection at that index.
- Click-to-position: clicking a slot moves the real input's caret there.
- Paste: native, filtered by `pattern` if set.
- Keyboard: entirely native single-input text editing (Backspace, Delete,
  Home/End, arrow keys, Shift+arrow selection) — zero bespoke keyboard JS
  required, a genuine simplification relative to what a naive multi-input
  reimplementation would need.

### 3. gsxui adaptation notes

- **Parts**: `InputOTP` (renders the container `<div data-slot="input-otp"
  data-gsxui-input-otp class="relative flex items-center gap-2
  has-[input:disabled]:opacity-50">` wrapping the real, visually-hidden
  `<input data-gsxui-input-otp-input maxlength={maxLength} pattern={pattern}
  inputmode="numeric" autocomplete="one-time-code" class="absolute inset-0
  z-10 h-full w-full cursor-text opacity-0 disabled:cursor-not-allowed"
  {attrs...}/>` — `maxLength`/`pattern`/`name`/`value`/`disabled`/
  `aria-invalid` all fall through via `attrs`, matching `ui/input.gsx`'s own
  convention — plus `{children}` (the slot markup)), `InputOTPGroup`
  (unchanged plain flex wrapper, zero behavior, straight port), `InputOTPSlot`
  (renders `<div data-slot="input-otp-slot" data-active="false"></div>`
  **empty** at server-render time — gsx has no client React context to
  pre-populate `char`/`isActive` from an initial `value` the way shadcn's
  SSR-capable React version can; the slot starts inert and is populated
  entirely by client JS after mount, a small but real first-paint
  divergence from shadcn worth ledgering), `InputOTPSeparator` (`icon.Minus`,
  static, unchanged).
  - Hidden-input CSS: `opacity-0` (not `sr-only`/`hidden`/`display:none` —
    those break focusability/paste-target-ability in most browsers) and
    **not** `pointer-events-none` — kept clickable and layered on top
    (`z-10`) so native click-to-position-caret still works as a fallback
    even without the slot-click handler below.
- **Index binding — the open design question the task flags explicitly**:
  shadcn's `InputOTPSlot` takes an explicit `index` prop
  (`<InputOTPSlot index={0}/>`) keying into React context's `slots[index]`.
  gsx has no equivalent shared context. Two options:
  - **(a) Explicit `index` param**, byte-parity with shadcn's own call
    sites (`<ui.InputOTPSlot index={0}>`) — simplest 1:1 API match, but
    error-prone: the caller must count correctly by hand across
    `InputOTPSeparator`-split groups (`input-otp-separator.tsx`'s 3-group
    layout needs indices `0,1 | 2,3 | 4,5` spanning group boundaries —
    trivial to get wrong if slots are ever reordered/edited).
  - **(b) DOM-order stamping** (recommended): `InputOTPSlot` takes **no**
    index param at all; `ui/input-otp.js` walks
    `root.querySelectorAll('[data-slot="input-otp-slot"]')` in DOM order
    at module-init (and on any HTMX-swap re-init) and stamps `data-index`
    positionally — the exact same "stamp source order once, JS-computed
    identity thereafter" idiom `command.js` already establishes for its
    own `data-gsxui-index` (`ui/command.js` lines 130–141, byte-read).
    This is a **genuine departure from shadcn's authored API shape** —
    dropping the `index` param changes `InputOTPSlot`'s call signature
    from shadcn's own — flagged prominently since it's a real, deliberate
    choice, not an oversight. Recommended specifically because it reuses
    an established in-codebase precedent and removes an entire class of
    off-by-one caller bugs across separator-split groups.
- **`ui/input-otp.js` contract**:
  - On mount: for each `[data-gsxui-input-otp]` root, stamp `data-index`
    on its slots in DOM order (option b above).
  - Delegated `input`/`selectionchange`/`focus`/`blur` on
    `[data-gsxui-input-otp-input]`: recompute every slot's `char` (from
    `input.value[i]`), `data-active` (true only for the slot at
    `input.selectionStart` when the selection is collapsed and the input
    has focus), and fake-caret visibility (active slot + empty char only)
    — re-render each slot's text/caret markup accordingly.
  - Delegated `click` on `[data-slot="input-otp-slot"]`: `input.focus();
    input.setSelectionRange(index, index)`.
  - Pattern filtering: an `input` handler reads
    `data-gsxui-input-otp-pattern` (a regex-source string attribute on the
    real input) and strips non-matching characters from `input.value`
    before recomputing slot display. **Signature mismatch to resolve**:
    shadcn's own patterns (e.g. `REGEXP_ONLY_DIGITS_AND_CHARS`) are
    *whole-string-anchored* regexes (`^[0-9A-Za-z]*$`-shaped) fed to the
    library for a different internal use; a literal per-character
    `pattern.test(char)` test against an anchored `^...$` pattern would
    reject every single non-empty character (a one-char string never
    matches a `*`-quantified full-string anchor meaningfully across
    iterative single-char tests the way one might assume). Recommend
    `data-gsxui-input-otp-pattern` instead take an **unanchored,
    per-character class** (e.g. `[0-9]` not `^[0-9]*$`) — a real design
    decision the planner needs to make explicit in the component's
    doc/example, not something a literal copy of shadcn's exported
    pattern constants would get right by accident.
- InputOTP does **not** compose `ui.Input` — the per-slot visual chrome
  (borders, `first:`/`last:` corner rounding, `data-[active=true]` ring)
  is different enough from Input's own single-box recipe that composing
  would fight rather than help, the same reasoning `## sheet` gives for
  not composing `DialogContent`.
- Registry: only dep is `icon` (from `InputOTPSeparator`'s `icon.Minus`) →
  `Deps("input-otp") == ["icon"]`. `HasJS("input-otp") == true`.

**Nova deltas** (`.cn-input-otp*`, `style-nova.css` lines 682–701,
byte-read):

| part | new-york-v4 | nova | delta |
|---|---|---|---|
| slot size | `h-9 w-9` | `size-8` | `h-9 w-9 (36px) → size-8 (32px) (nova)` |
| slot shadow | `shadow-xs` | (absent) | removed — nova shadow-presence drop, consistent with input/checkbox/etc.'s own nova entries |
| slot corner radius | `first:rounded-l-md … last:rounded-r-md` | `first:rounded-l-lg … last:rounded-r-lg` | `md→lg (nova)`, matches the global radius-scale bump |
| slot active ring | `ring-[3px]` | `ring-3` | same value, Tailwind v4 token spelling only, no visual delta |
| slot `z-10` on active | present | not present in nova's excerpt | **keep `z-10` regardless** — functionally necessary (keeps the active ring from being occluded by a neighboring slot's border), not a token nova is deliberately dropping; likely just omitted from the generated utility-only reference sheet, not a real removal |
| container gap | `gap-2` (via `InputOTP`'s `containerClassName`) | `.cn-input-otp { gap-2 }` | `(nova, no delta)` |
| group | `flex items-center` only | `.cn-input-otp-group` adds `has-aria-invalid:ring-destructive/20 dark:has-aria-invalid:ring-destructive/40 has-aria-invalid:border-destructive rounded-lg has-aria-invalid:ring-3` | **new** group-level invalid-state ring, absent from new-york-v4's `InputOTPGroup` entirely (which has no validation styling of its own, only per-slot `aria-invalid:border-destructive`). Borderline color-vs-metric: recommend adopting — the `ring-*` width is metric, and the destructive-ring-on-invalid color pattern already exists elsewhere in gsxui (checkbox/radio/input's own `aria-invalid:ring-destructive*`), so it's precedent-consistent, not a novel color introduction |
| caret line | `h-4 w-px animate-caret-blink bg-foreground duration-1000` | same | `(nova, no delta)` |
| separator | no class | `.cn-input-otp-separator { [&_svg:not([class*='size-'])]:size-4 }` | minor sizing safeguard; likely a no-op if `icon.Minus` already defaults to `size-4` internally — recommend carrying it regardless, harmless either way |

### 4. Demo inventory

- `input-otp-demo.tsx` — 6-digit, two 3-slot groups + one separator.
  Baseline, exercises option (b)'s index-stamping across a group split.
- `input-otp-controlled.tsx` — single 6-slot group, live-updates a "You
  entered: …" text via React `useState`. **Adaptation needed**: gsx has no
  live two-way client↔server value binding without a page reload or an
  HTMX round trip; recommend showing this instead as an initial-`value`-set
  example (server renders a pre-filled OTP field) rather than a literal
  live-typing-echo port — note the gap explicitly rather than silently
  reinterpreting the demo's intent.
- `input-otp-pattern.tsx` — `REGEXP_ONLY_DIGITS_AND_CHARS`. Exercises the
  `data-gsxui-input-otp-pattern` filter attribute directly — good proof
  case for the per-character-vs-anchored-pattern design decision above.
- `input-otp-separator.tsx` — three 2-slot groups + two separators.
  Exercises DOM-order index-stamping across *three* group boundaries
  (indices `0,1 | 2,3 | 4,5`) — the strongest proof case for option (b)'s
  correctness, recommend keeping as a site example specifically for this
  reason.

Recommend `input-otp-demo`, `input-otp-pattern`, and `input-otp-separator`
as the three site examples; `input-otp-controlled` adapted-or-dropped per
the note above (its two-way-binding premise doesn't translate to SSR
directly).

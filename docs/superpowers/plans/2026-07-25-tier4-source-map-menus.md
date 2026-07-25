# Tier 4 Batch B source map — menu family (dropdown-menu / context-menu / menubar / navigation-menu)

Source analysis for Task 0 of the Tier 4 Batch B plan
(`.superpowers/sdd/2026-07-25-tier4-batch-b/task-0-brief.md`). Tasks 1–4 cite
this document by section heading (`## shared-items`, `## menubar`,
`## navigation-menu`) for every class string and behavioral claim — do not
re-derive either from the `.tsx` files directly; cite here.

**Note on deliverable location**: the brief (step 1 of its own task list)
names the canonical output path as
`docs/superpowers/plans/2026-07-25-tier4-source-map-menus.md`, with a final
`git add && git commit` step. The orchestrating instructions for this
research pass restricted me to read-only research plus **one** output file,
delivered via the report contract at this path
(`.superpowers/sdd/2026-07-25-tier4-batch-b/task-0-report.md`), and
explicitly did not ask for a commit. I followed the report contract as the
more specific, most-recent instruction and did not create or commit the
brief's separate canonical file. **Flagging this explicitly**: if Tasks 1–4
(or their orchestrator) expect the doc at
`docs/superpowers/plans/2026-07-25-tier4-source-map-menus.md`, someone needs
to copy/move this file's content there (and commit) before those tasks run
— content and structure below already match what that file is supposed to
contain, section-heading-for-section-heading.

**Reference-source rule** (adjudicated, not re-litigated — reproduced from
the calling context): markup reference is
`apps/v4/registry/new-york-v4/ui/<name>.tsx` — the only remaining
inline-Tailwind-utility form, gsxui's own authoring model. Do NOT port from
`registry/bases/{radix,base,aria}/ui/<name>.tsx` — those emit `.cn-*` named
classes with utilities extracted into `registry/styles/style-*.css`. Visual
reference is `registry/styles/style-nova.css`; where it and new-york-v4
disagree, nova wins. `bases/base/ui/<name>.tsx` is a structure tiebreak only,
cited inline wherever used.

**Inputs read** (byte-read in full unless noted): `registry/new-york-v4/ui/
{dropdown-menu,context-menu,menubar,navigation-menu}.tsx` (all four, full
files, 257/252/276/168 lines respectively); `registry/styles/style-nova.css`
(`.cn-{dropdown-menu,context-menu,menubar,navigation-menu}-*` sections, full
text of each, byte-read); `registry/bases/radix/ui/navigation-menu.tsx`
(first ~80 lines — structure/CSS-var-name corroboration only, see
`## navigation-menu` §2); `registry/bases/base/ui/{navigation-menu,menubar}.tsx`
(first ~20 lines each — read only to confirm these wrap a **different**
underlying library, `@base-ui/react`, not Radix at all, so they carry zero
evidentiary weight for this Radix-based port; see Library presence check);
`content/docs/components/{radix,base}/navigation-menu.mdx`,
`content/docs/components/radix/menubar.mdx` (grepped for
viewport/resize/observer/measure/roving/tabindex/hover prose — zero hits;
these are pure prop-table/usage docs, not implementation prose).

**Library presence check** (Step 1 of the brief):
```
ls ~/personal/shadcn-ui/node_modules                          → no such directory (node_modules absent entirely)
find ~/personal/shadcn-ui -iname '*radix-ui*' -not -path '*/.git/*' → no hits outside apps/v4/{examples,content/docs,registry/bases}/radix (source dirs, not the package)
grep radix-ui apps/v4/package.json                             → "radix-ui": "^1.4.3" (declared, not installed)
ls ~/Library/pnpm/store/v10                                    → exists but content-addressable/opaque, no unpacked radix-ui tree found under it
```
This checkout has **no `node_modules` at all**, matching the batch-A finding
exactly. `radix-ui@^1.4.3` (the umbrella package new-york-v4 imports
`DropdownMenu`/`ContextMenu`/`Menubar`/`NavigationMenu` from, per each
file's own `import { X as XPrimitive } from "radix-ui"` line) is declared in
`apps/v4/package.json` but never installed in this checkout, and no built
`dist/index.mjs` (or equivalent) for it exists anywhere reachable.
**Consequently, every ARIA-attribute claim, every keyboard-handling claim,
and every claim about roving tabindex, hover-follow, or viewport-measurement
mechanics in this document is `derived-not-read`**, tagged at first use —
none of it comes from Radix's own runtime source. Where I state such a claim
with confidence, it is because it is well-established public Radix API
behavior (documented on Radix's own public docs site / WAI-ARIA menu
pattern), not because I read it here — and I say so every time.

`registry/bases/{radix,base,aria}/ui/*.tsx` were **not** ported from
(reference-source rule). `bases/radix/ui/navigation-menu.tsx` was read only
as corroboration that the CSS custom-property *names*
(`--radix-navigation-menu-viewport-{height,width}`) are literal and
unchanged across both wrapper styles — it is the same Radix primitive
underneath, just with `.cn-*` classnames, so it adds no new behavioral
information, only naming corroboration (used in `## navigation-menu` §2).
`bases/base/ui/{navigation-menu,menubar}.tsx` import from `@base-ui/react/
navigation-menu` and `@base-ui/react/menu` + `@base-ui/react/menubar`
respectively — a **completely different underlying primitive library** from
Radix, confirmed from their own import lines. They were not used as a
structure tiebreak anywhere in this document because new-york-v4's own
markup raised no structural ambiguity nova's CSS couldn't answer — citing a
different library's wrapper would have been actively misleading here, not a
tiebreak.

## Legend

(Reused verbatim from the batch-A source map's own legend, same meanings.)

- `derived-not-read` — reconstructed from the wrapper `.tsx`'s own
  props/classNames/data-attribute selectors, prose documentation, or general
  public-API knowledge, **not** from reading the library's own built
  source — unavailable in this checkout. Where I could not corroborate a
  claim even indirectly, I say so explicitly rather than filling the gap.
- `(no nova counterpart)` — nova's stylesheet has no `.cn-*` class for this
  part at all; new-york-v4's value stands unreviewed by the retarget.
- Structure tiebreak citations (`bases/base/ui/<name>.tsx`,
  `bases/radix/ui/<name>.tsx`) are called out by name every time they're
  used — never silently folded into "nova says."

---

## shared-items

Covers the seven parts `dropdown-menu.tsx` and `context-menu.tsx` share
conceptually: `Group`, `CheckboxItem`, `RadioGroup`, `RadioItem`, `Sub`,
`SubTrigger`, `SubContent`.

### 1. Class strings verbatim, element, `data-slot`, indicator markup

**`Group`**
- Dropdown: `<DropdownMenuPrimitive.Group data-slot="dropdown-menu-group" {...props} />` — no class string at all.
- Context: `<ContextMenuPrimitive.Group data-slot="context-menu-group" {...props} />` — no class string at all.

**`CheckboxItem`** — both wrap `PrimitiveX.CheckboxItem`, same structure:
`{indicator span}{children}`.
- Class string (byte-identical in both, only the component name changes):
  ```
  relative flex cursor-default items-center gap-2 rounded-sm py-1.5 pr-2 pl-8 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4
  ```
- `data-slot`: `dropdown-menu-checkbox-item` / `context-menu-checkbox-item`.
- Indicator markup (identical in both, only primitive-namespace differs):
  `<span className="pointer-events-none absolute left-2 flex size-3.5 items-center justify-center"><XxxPrimitive.ItemIndicator><CheckIcon className="size-4" /></XxxPrimitive.ItemIndicator></span>`, rendered **before** `{children}`.
- Props: both forward `checked` explicitly (`checked={checked}` alongside `{...props}`); neither restates `data-inset` (CheckboxItem never gets an `inset` prop in either file).

**`RadioGroup`**
- Dropdown: `<DropdownMenuPrimitive.RadioGroup data-slot="dropdown-menu-radio-group" {...props} />` — no class string.
- Context: `<ContextMenuPrimitive.RadioGroup data-slot="context-menu-radio-group" {...props} />` — no class string.

**`RadioItem`** — same shape as CheckboxItem.
- Class string (byte-identical in both):
  ```
  relative flex cursor-default items-center gap-2 rounded-sm py-1.5 pr-2 pl-8 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4
  ```
  (identical to `CheckboxItem`'s own class string, character for character —
  own finding, not stated by shadcn anywhere: `CheckboxItem` and `RadioItem`
  share one base style verbatim in **both** files.)
- `data-slot`: `dropdown-menu-radio-item` / `context-menu-radio-item`.
- Indicator markup: `<span className="pointer-events-none absolute left-2 flex size-3.5 items-center justify-center"><XxxPrimitive.ItemIndicator><CircleIcon className="size-2 fill-current" /></XxxPrimitive.ItemIndicator></span>`, before `{children}`. Identical in both modulo namespace.

**`Sub`**
- Dropdown: `<DropdownMenuPrimitive.Sub data-slot="dropdown-menu-sub" {...props} />` — no class string.
- Context: `<ContextMenuPrimitive.Sub data-slot="context-menu-sub" {...props} />` — no class string.

**`SubTrigger`** — **not** identical modulo prefix; see §2 diff below for
the two real deltas. Shared skeleton: `data-inset={inset}`, renders
`{children}` then a trailing chevron icon.
- Dropdown class string:
  ```
  flex cursor-default items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[inset]:pl-8 data-[state=open]:bg-accent data-[state=open]:text-accent-foreground [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground
  ```
  Trailing icon: `<ChevronRightIcon className="ml-auto size-4" />`.
- Context class string:
  ```
  flex cursor-default items-center rounded-sm px-2 py-1.5 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[inset]:pl-8 data-[state=open]:bg-accent data-[state=open]:text-accent-foreground [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground
  ```
  Trailing icon: `<ChevronRightIcon className="ml-auto" />` (no `size-4`).
- `data-slot`: `dropdown-menu-sub-trigger` / `context-menu-sub-trigger`.

**`SubContent`**
- Dropdown class string:
  ```
  z-50 min-w-[8rem] origin-(--radix-dropdown-menu-content-transform-origin) overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-lg data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95
  ```
- Context class string:
  ```
  z-50 min-w-[8rem] origin-(--radix-context-menu-content-transform-origin) overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-lg data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95
  ```
  Only difference is the CSS custom-property *name* embedded in
  `origin-(...)` (`dropdown-menu` vs `context-menu` — the expected,
  necessary per-primitive substitution; not a real style delta).
- `data-slot`: `dropdown-menu-sub-content` / `context-menu-sub-content`.

### 2. Diff verdict — **not fully identical**; two real deltas, both in `SubTrigger`

Answering the brief's explicit deliverable (item a): `Group`, `CheckboxItem`,
`RadioGroup`, `RadioItem`, `Sub`, and `SubContent` are **byte-identical
modulo the component-name/CSS-var-name prefix substitution** — Task 2 can
treat those six mechanically, port dropdown-menu's version and s/dropdown/
context/ where the prefix appears, done.

`SubTrigger` is **not** identical. Two real, deliberate-looking-but-verbatim
deltas, both dropdown-menu-only:
1. **`gap-2`** is present in dropdown-menu's `SubTrigger` class string,
   **absent** in context-menu's (context-menu's `items-center` is followed
   directly by `rounded-sm`, no `gap-2` token in between).
2. Dropdown-menu's trailing `<ChevronRightIcon>` carries an explicit
   `size-4` class; context-menu's does not (just `ml-auto`).

Net visual effect of #2 is muted but not zero: both `SubTrigger` class
strings carry `[&_svg:not([class*='size-'])]:size-4`, a descendant selector
that stamps `size-4` onto any child `svg` that doesn't already have a
`size-*` class of its own. Since context-menu's `ChevronRightIcon` has no
`size-*` class, that selector *does* still apply `size-4` to it — so the
rendered icon size ends up the same in both. What does **not** get
compensated for is #1: nothing in context-menu's `SubTrigger` class
re-adds the `gap-2` spacing between the item's icon/content and the
trigger's own start, so context-menu's `SubTrigger` genuinely lays out with
less inter-element gap than dropdown-menu's, unless something else in a
particular usage supplies it. **Corroboration nova already "fixes" this
asymmetry**: `.cn-dropdown-menu-sub-trigger` and `.cn-context-menu-sub-trigger`
in `style-nova.css` (§`## shared-items` nova block below) **both** carry
`gap-1.5` — nova harmonizes the two, treating dropdown-menu's `gap-2`
presence as the norm and context-menu's absence as the thing to fix. Since
nova wins on visual disagreement (adjudicated rule), **recommend Task 2 add
the gap token to context-menu's ported `SubTrigger`** rather than preserving
the asymmetry, even though byte-for-byte new-york-v4 itself is inconsistent
here. Do preserve the missing `size-4` on the icon literally as no-op (the
selector covers it) — not worth deviating from source for.

### 3. ARIA anatomy — `derived-not-read`

None of `role`, `aria-checked`, `aria-haspopup`, or `aria-expanded` appear
literally anywhere in either `.tsx` file's own JSX — every one of Radix's
ARIA stamps happens inside the primitive's own runtime, not present in this
checkout (Library presence check above). What the public Radix Menu/
DropdownMenu/ContextMenu API contract (Radix's own published docs, not
fetched as part of this pass, and general WAI-ARIA menu-pattern knowledge)
predicts, stated as a hypothesis for Task 2 to implement against, **not** a
verified fact:
- `Content`/`SubContent` root → `role="menu"`.
- `Item` → `role="menuitem"`.
- `CheckboxItem` → `role="menuitemcheckbox"`, `aria-checked={checked}`.
- `RadioItem` → `role="menuitemradio"`, `aria-checked` reflecting whether
  this item's own `value` matches the enclosing `RadioGroup`'s `value`.
- `SubTrigger` → `aria-haspopup="menu"`, `aria-expanded` reflecting the
  sub's own open state, `aria-controls` pointing at the `SubContent`'s id.
- `Trigger` (top-level) → `aria-haspopup="menu"`, `aria-expanded`.

This is the same hedge batch-A applied to Base UI's combobox ARIA anatomy
(`## combobox` §4 of the batch-A map) — named as the standard pattern for
Task 2 to confirm independently, not pinned here as verified.

### 4. Keyboard model for submenus — `derived-not-read`

Not traceable to Radix's dist in this checkout (absent). What is directly
confirmed from the `.tsx` source itself: `SubContent`'s class carries
`data-[side=left]:slide-in-from-right-2` / `data-[side=right]:slide-in-from-
left-2` alongside the `top`/`bottom` variants — i.e. Radix's positioning
already accounts for a submenu opening to either the left or right of its
trigger (flip behavior), consistent with (but not proof of) the standard
Radix submenu keyboard contract: `ArrowRight` (in a horizontally-flowing
root, or the "into" direction generally) opens a highlighted `SubTrigger`'s
submenu and moves focus into it; `ArrowLeft` on a submenu's own `Content`
closes that submenu and returns focus to its `SubTrigger`; `Escape` closes
the innermost open menu/submenu only (not the whole stack) and returns
focus to that level's trigger; typeahead (single-character search jumping
focus to a matching item's label) is scoped per open menu level, not
global across the whole stack. **All of the above is named as the
well-known public Radix menu contract, not verified from source here** —
flagged in its strongest form per the brief's instruction not to guess:
Task 2 should treat this paragraph as "what to implement, pending its own
confirmation," not as a citation of something actually read.

### 5. Popover-nesting design vs. Radix's portal — what Task 1 gains and loses

The design (already settled, not mine to revisit): submenu content is a
`popover="auto"` element DOM-nested inside its parent's content, not
portalled. Two concrete things Radix's portal gives dropdown-menu/
context-menu content that DOM-nesting will not, both directly readable from
the class strings above, not inferred:

1. **Overflow escape.** `DropdownMenuContent`/`ContextMenuContent`'s own
   class string (not one of the seven shared parts, but the parent
   `SubContent` nests inside) is
   `max-h-(--radix-dropdown-menu-content-available-height) ... overflow-x-hidden overflow-y-auto` —
   the **top-level** content is a clipped, scrollable container by design
   (it clamps to the available viewport height and scrolls internally when
   the menu is taller than that). Radix's `SubContent` is portalled to
   `document.body` (or a designated container), so it is never a descendant
   of that scrolling box and is never clipped by it. A DOM-nested
   `SubContent` **would** sit inside that `overflow-y-auto` ancestor and
   risk being visually clipped/scrolled-away whenever the parent menu's own
   content is tall enough to need its scrollbar — a real, concrete risk
   specific to dropdown-menu/context-menu (menubar's own top-level
   `MenubarContent` is `overflow-hidden` only, not scrollable, so this
   particular risk doesn't apply there — see `## menubar` §1).
2. **Stacking/z-index independence.** Because the portalled subtree isn't a
   descendant of its trigger's own stacking context, its `z-50` is
   evaluated against `document.body`'s own stacking context, not whatever
   local stacking context the parent `Content` (itself `z-50`, also
   portalled) might establish. DOM-nesting collapses this: the submenu's
   effective stacking is now constrained by being a literal descendant of
   the parent popover's own box — mostly fine since native top-layer
   popovers already promote to the top layer regardless of DOM position,
   but worth Task 1 confirming empirically rather than assuming z-index
   alone still does the ordering work it did under Radix's portal model.
3. **Focus order / DOM order.** Radix manages focus entirely
   programmatically (roving `tabindex` + explicit `.focus()` calls) inside
   the portalled subtree, so the submenu's position in the page's actual
   DOM/tab order (end-of-`<body>`) is irrelevant to it — Radix's own state
   machine, not document order, drives what gets focused next. A DOM-nested
   submenu, by contrast, sits at its literal position in the document for
   any tooling that reads DOM order directly (not all a11y tooling does;
   most respects explicit `tabindex`/focus management over raw DOM order)
   — this is a neutral-to-positive difference, not a loss, but it's a
   *different* model from what Radix does, and Task 1 should manage focus
   explicitly rather than assume DOM nesting alone produces correct tab
   order, matching what Radix does today only by coincidence.

None of the above is `derived-not-read` — it follows directly from the
class strings (read in full, quoted above) plus the stated, already-measured
Chrome popover behavior in the calling context, not from Radix's own source.
The one thing that **is** `derived-not-read`: whether Radix's actual
`Positioner`/`Portal` implementation does anything beyond plain CSS
`position: fixed` + a portal target (e.g., its own internal `z-index`
management, or coordination with multiple simultaneously-open submenus at
different depths) that DOM-nesting could silently break in a way not
predictable from the class strings alone. Flagged as an open risk for Task
1's own testing, not resolved here.

### `.cn-{dropdown-menu,context-menu}-*` nova rules (verbatim, `style-nova.css`)

Dropdown (lines 537–583):
```css
.cn-dropdown-menu-content {
  @apply data-open:animate-in data-closed:animate-out data-closed:fade-out-0 data-open:fade-in-0 data-closed:zoom-out-95 data-open:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 ring-foreground/10 bg-popover text-popover-foreground min-w-32 rounded-lg p-1 shadow-md ring-1 duration-100;
}
.cn-dropdown-menu-content-logical {
  @apply data-[side=inline-start]:slide-in-from-right-2 data-[side=inline-end]:slide-in-from-left-2;
}
.cn-dropdown-menu-item {
  @apply focus:bg-accent focus:text-accent-foreground data-[variant=destructive]:text-destructive data-[variant=destructive]:focus:bg-destructive/10 dark:data-[variant=destructive]:focus:bg-destructive/20 data-[variant=destructive]:focus:text-destructive data-[variant=destructive]:*:[svg]:text-destructive not-data-[variant=destructive]:focus:**:text-accent-foreground gap-1.5 rounded-md px-1.5 py-1 text-sm data-inset:pl-7 [&_svg:not([class*='size-'])]:size-4;
}
.cn-dropdown-menu-checkbox-item {
  @apply focus:bg-accent focus:text-accent-foreground focus:**:text-accent-foreground gap-1.5 rounded-md py-1 pr-8 pl-1.5 text-sm data-inset:pl-7 [&_svg:not([class*='size-'])]:size-4;
}
.cn-dropdown-menu-radio-item {
  @apply focus:bg-accent focus:text-accent-foreground focus:**:text-accent-foreground gap-1.5 rounded-md py-1 pr-8 pl-1.5 text-sm data-inset:pl-7 [&_svg:not([class*='size-'])]:size-4;
}
.cn-dropdown-menu-item-indicator {
  @apply absolute right-2 flex items-center justify-center;
}
.cn-dropdown-menu-label {
  @apply text-muted-foreground px-1.5 py-1 text-xs font-medium data-inset:pl-7;
}
.cn-dropdown-menu-separator {
  @apply bg-border -mx-1 my-1 h-px;
}
.cn-dropdown-menu-shortcut {
  @apply text-muted-foreground group-focus/dropdown-menu-item:text-accent-foreground ml-auto text-xs tracking-widest;
}
.cn-dropdown-menu-sub-trigger {
  @apply focus:bg-accent focus:text-accent-foreground data-open:bg-accent data-open:text-accent-foreground not-data-[variant=destructive]:focus:**:text-accent-foreground gap-1.5 rounded-md px-1.5 py-1 text-sm data-inset:pl-7 [&_svg:not([class*='size-'])]:size-4;
}
.cn-dropdown-menu-sub-content {
  @apply data-open:animate-in data-closed:animate-out data-closed:fade-out-0 data-open:fade-in-0 data-closed:zoom-out-95 data-open:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 ring-foreground/10 bg-popover text-popover-foreground min-w-[96px] rounded-lg p-1 shadow-lg ring-1 duration-100;
}
.cn-dropdown-menu-subcontent {
  @apply shadow-lg;
}
```
Note `.cn-dropdown-menu-item-indicator` here is a **generic shared indicator
rule** — nova consolidates checkbox/radio indicator styling into one class
(`absolute right-2 flex items-center justify-center`) rather than new-york-v4's
per-part inline `<span>` classes; the indicator also moves from `left-2` (new-york-v4,
`pl-8` items reserve left space) to `right-2` (nova, `pr-8` items reserve
right space instead) — a real left↔right swap, not a typo, confirmed by
`.cn-dropdown-menu-checkbox-item`'s own `pr-8 pl-1.5` (right-padded) vs.
new-york-v4's `pr-2 pl-8` (left-padded). Nova moves the checkmark/dot from
the item's left edge to its right edge. Flag for Task 2 — this is a bigger
structural delta than a metric-only nova pass, in the same spirit as
batch-A's resizable handle-icon finding.

Context (lines 414–458), same shape, `right-2`-only indicator (no
`items-center justify-center` here, unlike dropdown's — minor, `.cn-context-
menu-item-indicator { @apply absolute right-2; }`):
```css
.cn-context-menu-content {
  @apply data-open:animate-in data-closed:animate-out data-closed:fade-out-0 data-open:fade-in-0 data-closed:zoom-out-95 data-open:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 ring-foreground/10 bg-popover text-popover-foreground min-w-36 rounded-lg p-1 shadow-md ring-1 duration-100;
}
.cn-context-menu-content-logical {
  @apply data-[side=inline-start]:slide-in-from-right-2 data-[side=inline-end]:slide-in-from-left-2;
}
.cn-context-menu-item {
  @apply focus:bg-accent focus:text-accent-foreground data-[variant=destructive]:text-destructive data-[variant=destructive]:focus:bg-destructive/10 dark:data-[variant=destructive]:focus:bg-destructive/20 data-[variant=destructive]:focus:text-destructive data-[variant=destructive]:*:[svg]:text-destructive focus:*:[svg]:text-accent-foreground gap-1.5 rounded-md px-1.5 py-1 text-sm data-inset:pl-7 [&_svg:not([class*='size-'])]:size-4;
}
.cn-context-menu-checkbox-item {
  @apply focus:bg-accent focus:text-accent-foreground gap-1.5 rounded-md py-1 pr-8 pl-1.5 text-sm data-inset:pl-7 [&_svg:not([class*='size-'])]:size-4;
}
.cn-context-menu-radio-item {
  @apply focus:bg-accent focus:text-accent-foreground gap-1.5 rounded-md py-1 pr-8 pl-1.5 text-sm data-inset:pl-7 [&_svg:not([class*='size-'])]:size-4;
}
.cn-context-menu-item-indicator {
  @apply absolute right-2;
}
.cn-context-menu-label {
  @apply text-muted-foreground px-1.5 py-1 text-xs font-medium data-inset:pl-7;
}
.cn-context-menu-separator {
  @apply bg-border -mx-1 my-1 h-px;
}
.cn-context-menu-shortcut {
  @apply text-muted-foreground group-focus/context-menu-item:text-accent-foreground ml-auto text-xs tracking-widest;
}
.cn-context-menu-sub-trigger {
  @apply focus:bg-accent focus:text-accent-foreground data-open:bg-accent data-open:text-accent-foreground gap-1.5 rounded-md px-1.5 py-1 text-sm data-inset:pl-7 [&_svg:not([class*='size-'])]:size-4;
}
.cn-context-menu-sub-content {
  @apply data-open:animate-in data-closed:animate-out data-closed:fade-out-0 data-open:fade-in-0 data-closed:zoom-out-95 data-open:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 bg-popover text-popover-foreground min-w-32 rounded-lg border p-1 shadow-lg duration-100;
}
```
Confirms the §2 finding independently: `.cn-context-menu-sub-trigger` here
**does** carry `gap-1.5`, same as `.cn-dropdown-menu-sub-trigger` — nova
harmonizes what new-york-v4 leaves asymmetric.

---

## menubar

### 1. Every part's class string, verbatim (`registry/new-york-v4/ui/menubar.tsx`, 276 lines, full file read)

**`Menubar`** — `<MenubarPrimitive.Root data-slot="menubar">`:
```
flex h-9 items-center gap-1 rounded-md border bg-background p-1 shadow-xs
```

**`MenubarMenu`** — `<MenubarPrimitive.Menu data-slot="menubar-menu" {...props} />`, no class.

**`MenubarGroup`** — `<MenubarPrimitive.Group data-slot="menubar-group" {...props} />`, no class.

**`MenubarPortal`** — `<MenubarPrimitive.Portal data-slot="menubar-portal" {...props} />`, no class.

**`MenubarRadioGroup`** — `<MenubarPrimitive.RadioGroup data-slot="menubar-radio-group" {...props} />`, no class.

**`MenubarTrigger`** — `data-slot="menubar-trigger"`:
```
flex items-center rounded-sm px-2 py-1 text-sm font-medium outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[state=open]:bg-accent data-[state=open]:text-accent-foreground
```

**`MenubarContent`** — wrapped in `MenubarPortal`; defaults `align="start"`,
`alignOffset={-4}`, `sideOffset={8}`; `data-slot="menubar-content"`:
```
z-50 min-w-[12rem] origin-(--radix-menubar-content-transform-origin) overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-md data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95
```
**Verbatim anomaly, confirmed by direct re-read of the source, not a
transcription slip**: this string has NO `data-[state=closed]:animate-out`
token — every sibling content class in this whole map (dropdown/context
menu's `Content`/`SubContent`, menubar's own `SubContent` below) carries the
triple `animate-out fade-out-0 zoom-out-95` together; `MenubarContent`
carries only `fade-out-0 zoom-out-95`, missing `animate-out`. Since
`tailwindcss-animate`'s `fade-out-0`/`zoom-out-95` utilities are themselves
keyframe-driven and gated by the `animate-out` class turning the animation
on, this reads as a probable upstream inconsistency in shadcn's own
menubar.tsx (top-level content's close animation may be inert while every
other close animation in this whole menu family works) rather than an
intentional difference. **Port it verbatim anyway per the house rule (copy
character for character)** but flag it prominently for Task 3 as a
candidate bug to decide on deliberately, not silently reproduce and forget.

**`MenubarItem`** — `data-slot="menubar-item"`, `data-inset={inset}`,
`data-variant={variant}` (`"default" | "destructive"`, default `"default"`):
```
relative flex cursor-default items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 data-[inset]:pl-8 data-[variant=destructive]:text-destructive data-[variant=destructive]:focus:bg-destructive/10 data-[variant=destructive]:focus:text-destructive dark:data-[variant=destructive]:focus:bg-destructive/20 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground data-[variant=destructive]:*:[svg]:text-destructive!
```
Byte-identical to `DropdownMenuItem`'s/`ContextMenuItem`'s own class string
(not one of the seven "shared-items" parts, but worth recording: `Item`
itself is also identical across all three menu families, modulo prefix —
only `SubTrigger` broke the pattern).

**`MenubarCheckboxItem`** — `data-slot="menubar-checkbox-item"`:
```
relative flex cursor-default items-center gap-2 rounded-xs py-1.5 pr-2 pl-8 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4
```
**Differs from dropdown/context-menu's `CheckboxItem`**: `rounded-xs` here
vs. `rounded-sm` there — the only token difference, everything else
byte-identical. Indicator markup identical in shape (`<span
class="pointer-events-none absolute left-2 flex size-3.5 items-center
justify-center"><MenubarPrimitive.ItemIndicator><CheckIcon
className="size-4" /></MenubarPrimitive.ItemIndicator></span>` before
`{children}`).

**`MenubarRadioItem`** — `data-slot="menubar-radio-item"`, same
`rounded-xs` delta vs. dropdown/context-menu's `RadioItem`:
```
relative flex cursor-default items-center gap-2 rounded-xs py-1.5 pr-2 pl-8 text-sm outline-hidden select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4
```
Indicator: `<CircleIcon className="size-2 fill-current" />`, same wrapper
span, same shape as dropdown/context-menu.

**`MenubarLabel`** — `data-slot="menubar-label"`, `data-inset={inset}`:
```
px-2 py-1.5 text-sm font-medium data-[inset]:pl-8
```
Byte-identical to `DropdownMenuLabel` (context-menu's own `Label` adds
`text-foreground` that dropdown/menubar's don't — a pre-existing,
out-of-scope-here delta between context-menu and the other two, noted for
completeness, not load-bearing for this map's three required sections).

**`MenubarSeparator`** — `data-slot="menubar-separator"`:
```
-mx-1 my-1 h-px bg-border
```
Byte-identical to dropdown/context-menu's own `Separator`.

**`MenubarShortcut`** — `<span data-slot="menubar-shortcut">`:
```
ml-auto text-xs tracking-widest text-muted-foreground
```
Byte-identical to dropdown/context-menu's own `Shortcut`.

**`MenubarSub`** — `<MenubarPrimitive.Sub data-slot="menubar-sub" {...props} />`, no class.

**`MenubarSubTrigger`** — `data-slot="menubar-sub-trigger"`,
`data-inset={inset}`:
```
flex cursor-default items-center rounded-sm px-2 py-1.5 text-sm outline-none select-none focus:bg-accent focus:text-accent-foreground data-[inset]:pl-8 data-[state=open]:bg-accent data-[state=open]:text-accent-foreground
```
Trailing icon: `<ChevronRightIcon className="ml-auto h-4 w-4" />` — note
`h-4 w-4` here, not `size-4` (dropdown) or bare (context) — a **third**
distinct spelling for what is visually the same 16px icon across the three
menu families. Also note `outline-none` here vs. `outline-hidden` in every
other `SubTrigger`/`Item` class string in this whole map — another
menubar-only token spelling difference (`outline-none` removes the focus
ring entirely across all browsers; `outline-hidden` is Tailwind v4's
"visually hidden but still present for forced-colors/high-contrast modes"
utility — these are **not** visually equivalent in forced-colors mode, a
real, if narrow, accessibility-relevant delta, not just a spelling
variance). Flag both for Task 3.

**`MenubarSubContent`** — `data-slot="menubar-sub-content"`:
```
z-50 min-w-[8rem] origin-(--radix-menubar-content-transform-origin) overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-lg data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95
```
Has the full `animate-out` triple (unlike `MenubarContent` above) —
confirming the `MenubarContent` anomaly is local to that one part, not a
menubar-wide pattern.

### 2. Roving tabindex across triggers — `derived-not-read`

Nothing in `menubar.tsx` itself stamps `tabIndex` anywhere — no
`tabIndex={0}` / `tabIndex={-1}` literal in the file, confirmed by reading
it in full. This mechanism lives entirely inside `@radix-ui/react-menubar`'s
own runtime (bundled into the `radix-ui` umbrella package), which is absent
from this checkout (Library presence check above) — **not traceable to
source in this pass**. What I can state as the well-known public contract
(Radix's own docs describe `Menubar` explicitly as implementing the
WAI-ARIA "menubar" pattern's roving-tabindex requirement), offered as
implementation guidance for Task 3, not a verified fact: exactly one
`MenubarTrigger` carries `tabIndex={0}` at any time (the "active" trigger —
initially the first, or whichever one last had a menu open), every other
trigger carries `tabIndex={-1}`; `ArrowLeft`/`ArrowRight` on a focused
trigger move which trigger is the `tabIndex={0}` one and moves DOM focus to
it (this is why pressing Tab from inside a menubar only ever stops once,
not once per trigger — Tab exits the whole menubar/enters it at the single
roving-focus point, arrow keys move within it).

### 3. Open-follows-hover once a menu is open — `derived-not-read`

Same source-absence caveat. Public-contract description, not verified:
while no `MenubarMenu` is open, hovering a `MenubarTrigger` does **not**
open its menu (a menubar, unlike a plain dropdown, requires an explicit
click/Enter/Space/ArrowDown to open the *first* menu) — but once any one
menu is open, moving the pointer over a **sibling** trigger closes the
currently-open menu and opens the hovered trigger's own menu instead,
without requiring another click (this is the behavior that makes a menubar
feel like one connected control rather than N independent dropdowns, and is
the crux of what a from-scratch port needs to reproduce deliberately — it
is not automatic if Task 3 wires each `MenubarMenu` as an independent
Popover with independent hover-intent timers, since that would default to
"nothing happens until you click the sibling trigger" instead). **I could
not verify the precise trigger condition** (is it *any* open menu among
this `Menubar`'s children, tracked via shared context — almost certainly
yes, since Radix's `Root` already needs one to coordinate roving tabindex —
or something narrower); Task 3 should treat "some menu among this menubar's
siblings is currently open" as the gating condition to implement against,
sourced from general Radix `Menubar` behavior, not from a read of its
source in this checkout.

### `.cn-menubar-*` nova rules (verbatim, `style-nova.css` lines 798–852)

```css
.cn-menubar {
  @apply h-8 gap-0.5 rounded-lg border p-[3px];
}
.cn-menubar-trigger {
  @apply hover:bg-muted aria-expanded:bg-muted rounded-sm px-1.5 py-[2px] text-sm font-medium;
}
.cn-menubar-content {
  @apply bg-popover text-popover-foreground data-open:animate-in data-open:fade-in-0 data-open:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 ring-foreground/10 min-w-36 rounded-lg p-1 shadow-md ring-1 duration-100;
}
.cn-menubar-content-logical {
  @apply data-[side=inline-start]:slide-in-from-right-2 data-[side=inline-end]:slide-in-from-left-2;
}
.cn-menubar-item {
  @apply focus:bg-accent focus:text-accent-foreground data-[variant=destructive]:text-destructive data-[variant=destructive]:focus:bg-destructive/10 dark:data-[variant=destructive]:focus:bg-destructive/20 data-[variant=destructive]:focus:text-destructive data-[variant=destructive]:*:[svg]:text-destructive! not-data-[variant=destructive]:focus:**:text-accent-foreground gap-1.5 rounded-md px-1.5 py-1 text-sm data-disabled:opacity-50 data-inset:pl-7 [&_svg:not([class*='size-'])]:size-4;
}
.cn-menubar-checkbox-item {
  @apply focus:bg-accent focus:text-accent-foreground focus:**:text-accent-foreground gap-1.5 rounded-md py-1 pr-1.5 pl-7 text-sm data-inset:pl-7;
}
.cn-menubar-checkbox-item-indicator {
  @apply left-1.5 size-4 [&_svg:not([class*='size-'])]:size-4;
}
.cn-menubar-radio-item {
  @apply focus:bg-accent focus:text-accent-foreground focus:**:text-accent-foreground gap-1.5 rounded-md py-1 pr-1.5 pl-7 text-sm data-disabled:opacity-50 data-inset:pl-7 [&_svg:not([class*='size-'])]:size-4;
}
.cn-menubar-radio-item-indicator {
  @apply left-1.5 size-4 [&_svg:not([class*='size-'])]:size-4;
}
.cn-menubar-label {
  @apply px-1.5 py-1 text-sm font-medium data-inset:pl-7;
}
.cn-menubar-separator {
  @apply bg-border;
}
.cn-menubar-shortcut {
  @apply text-muted-foreground group-focus/menubar-item:text-accent-foreground text-xs tracking-widest;
}
.cn-menubar-sub-trigger {
  @apply focus:bg-accent focus:text-accent-foreground data-open:bg-accent data-open:text-accent-foreground gap-1.5 rounded-md px-1.5 py-1 text-sm data-inset:pl-7 [&_svg:not([class*='size-'])]:size-4;
}
.cn-menubar-sub-content {
  @apply bg-popover text-popover-foreground data-open:animate-in data-closed:animate-out data-closed:fade-out-0 data-open:fade-in-0 data-closed:zoom-out-95 data-open:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 ring-foreground/10 min-w-32 rounded-lg p-1 shadow-lg ring-1 duration-100;
}
```
Metric deltas worth naming (new-york-v4 → nova): `h-9→h-8` (root),
`rounded-md→rounded-lg` (root, content, sub-content), `p-1→p-[3px]` (root),
`rounded-sm→rounded-md` (item/sub-trigger; nova's own `.cn-menubar-checkbox/
radio-item` use `rounded-md` too, vs. new-york-v4's `rounded-xs` for those
two specifically — nova unifies checkbox/radio/item/sub-trigger all onto
`rounded-md`, erasing new-york-v4's own `rounded-xs` outlier noted in §1).
`.cn-menubar-content`'s own close-state classes: nova's version **also**
drops `animate-out`/`fade-out-0`/`zoom-out-95` for the closed state
entirely (only `data-open:*` present) — the same missing-animate-out
pattern noted as an anomaly in new-york-v4's own `MenubarContent`, so this
is **not** unique to new-york-v4's authoring; nova's independently-written
version has the identical omission, which raises this from "possible
transcription artifact" to "very likely an intentional-but-undocumented
choice specific to the top-level `MenubarContent`'s close transition,
reproduced consistently across both class-authoring styles." Downgrade the
"probable bug" framing in §1 accordingly — still flag it, but as a
deliberate-looking cross-style pattern to decide on, not obviously a slip.

---

## navigation-menu

### 1. Every part's class string, verbatim (`registry/new-york-v4/ui/navigation-menu.tsx`, 168 lines, full file read)

**`NavigationMenu`** — `data-slot="navigation-menu"`, `data-viewport={viewport}`
(prop, default `true`):
```
group/navigation-menu relative flex max-w-max flex-1 items-center justify-center
```
Renders `{children}` then, only `{viewport && <NavigationMenuViewport />}`.

**`NavigationMenuList`** — `data-slot="navigation-menu-list"`:
```
group flex flex-1 list-none items-center justify-center gap-1
```

**`NavigationMenuItem`** — `data-slot="navigation-menu-item"`:
```
relative
```

**`navigationMenuTriggerStyle`** (a `cva`, no variants — just a fixed base
string, exported standalone so plain `<a>`/`<Link>` "trigger-styled" items
that aren't real `Trigger`s can reuse it, e.g. a top-level link item with no
dropdown):
```
group inline-flex h-9 w-max items-center justify-center rounded-md bg-background px-4 py-2 text-sm font-medium transition-[color,box-shadow] outline-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1 disabled:pointer-events-none disabled:opacity-50 data-[state=open]:bg-accent/50 data-[state=open]:text-accent-foreground data-[state=open]:hover:bg-accent data-[state=open]:focus:bg-accent
```

**`NavigationMenuTrigger`** — `data-slot="navigation-menu-trigger"`, class
= `cn(navigationMenuTriggerStyle(), "group", className)` (i.e. the string
above, literally re-concatenated with a second, redundant `"group"` token —
verbatim as written, not a transcription error: the base `cva` string
already starts with `group`, and this component's own `cn()` call adds
`"group"` again as a **separate** third argument). Renders `{children}` then
a space then:
```tsx
<ChevronDownIcon
  className="relative top-[1px] ml-1 size-3 transition duration-300 group-data-[state=open]:rotate-180"
  aria-hidden="true"
/>
```

**`NavigationMenuContent`** — `data-slot="navigation-menu-content"`. Two
string arguments to `cn()`, both quoted verbatim (they're one logical class
list, split across two source lines in the `.tsx`, not two separate
concerns):
```
top-0 left-0 w-full p-2 pr-2.5 data-[motion=from-end]:slide-in-from-right-52 data-[motion=from-start]:slide-in-from-left-52 data-[motion=to-end]:slide-out-to-right-52 data-[motion=to-start]:slide-out-to-left-52 data-[motion^=from-]:animate-in data-[motion^=from-]:fade-in data-[motion^=to-]:animate-out data-[motion^=to-]:fade-out md:absolute md:w-auto
```
```
group-data-[viewport=false]/navigation-menu:top-full group-data-[viewport=false]/navigation-menu:mt-1.5 group-data-[viewport=false]/navigation-menu:overflow-hidden group-data-[viewport=false]/navigation-menu:rounded-md group-data-[viewport=false]/navigation-menu:border group-data-[viewport=false]/navigation-menu:bg-popover group-data-[viewport=false]/navigation-menu:text-popover-foreground group-data-[viewport=false]/navigation-menu:shadow group-data-[viewport=false]/navigation-menu:duration-200 **:data-[slot=navigation-menu-link]:focus:ring-0 **:data-[slot=navigation-menu-link]:focus:outline-none group-data-[viewport=false]/navigation-menu:data-[state=closed]:animate-out group-data-[viewport=false]/navigation-menu:data-[state=closed]:fade-out-0 group-data-[viewport=false]/navigation-menu:data-[state=closed]:zoom-out-95 group-data-[viewport=false]/navigation-menu:data-[state=open]:animate-in group-data-[viewport=false]/navigation-menu:data-[state=open]:fade-in-0 group-data-[viewport=false]/navigation-menu:data-[state=open]:zoom-in-95
```
Note the whole second block is gated on `group-data-[viewport=false]/
navigation-menu:` — i.e. `Content`'s own bordered-popover-with-shadow
look, its own open/close zoom animation, and even the `Link`-inside-it
focus-ring suppression **only apply when `viewport={false}`**. When
`viewport` is `true` (the default), none of that block's classes match at
all — `Content` in that mode is unstyled-container-only (`top-0 left-0
w-full p-2 pr-2.5` plus the `data-[motion=...]` slide-animation classes from
the first block only), because the actual popover chrome (border, bg,
shadow, rounded corners) moves onto `Viewport` instead (see below) — the
viewport is the thing that looks like a popover; `Content` in viewport mode
is just the scrollable inner panel of arbitrary width/height that gets
measured.

**`NavigationMenuViewport`** — two nested elements. Outer plain `<div>`
(no `data-slot`):
```
absolute top-full left-0 isolate z-50 flex justify-center
```
Inner `<NavigationMenuPrimitive.Viewport data-slot="navigation-menu-viewport">`:
```
origin-top-center relative mt-1.5 h-[var(--radix-navigation-menu-viewport-height)] w-full overflow-hidden rounded-md border bg-popover text-popover-foreground shadow data-[state=closed]:animate-out data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:zoom-in-90 md:w-[var(--radix-navigation-menu-viewport-width)]
```

**`NavigationMenuLink`** — `data-slot="navigation-menu-link"`:
```
flex flex-col gap-1 rounded-sm p-2 text-sm transition-all outline-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1 data-[active=true]:bg-accent/50 data-[active=true]:text-accent-foreground data-[active=true]:hover:bg-accent data-[active=true]:focus:bg-accent [&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground
```

**`NavigationMenuIndicator`** — `data-slot="navigation-menu-indicator"`:
```
top-full z-[1] flex h-1.5 items-end justify-center overflow-hidden data-[state=hidden]:animate-out data-[state=hidden]:fade-out data-[state=visible]:animate-in data-[state=visible]:fade-in
```
Renders one child div, the little rotated-square "arrow" pointer:
```
relative top-[60%] h-2 w-2 rotate-45 rounded-tl-sm bg-border shadow-md
```

### 2. How the viewport gets its size — CSS custom properties, set by Radix's own runtime

Two CSS custom properties, read directly off `Viewport`'s own class string
above: `--radix-navigation-menu-viewport-height` (drives `h-[var(...)]`,
all breakpoints) and `--radix-navigation-menu-viewport-width` (drives
`md:w-[var(...)]`, i.e. only takes effect at `md:` and above — below `md`
the viewport is `w-full` instead, unconstrained by the measured width).
Neither variable is set, referenced, or computed **anywhere** in
`navigation-menu.tsx` itself — no inline `style={{ "--radix-navigation-menu-
viewport-height": ... }}` anywhere in the file (confirmed, full file read).
They are set by `@radix-ui/react-navigation-menu`'s own `Viewport`
component internals, which is **absent from this checkout** (Library
presence check above) — `derived-not-read` for the actual mechanism.
Corroboration, not new information: `registry/bases/radix/ui/navigation-
menu.tsx` (read, first ~80 lines) uses the **identical** variable names
(`h-(--radix-navigation-menu-viewport-height)`,
`md:w-(--radix-navigation-menu-viewport-width)`, Tailwind v4's parenthesis
arbitrary-value spelling of the same thing) — confirming this is the same
underlying Radix primitive under different class-authoring conventions, not
new-york-v4-specific naming, but adding no insight into *how* the values get
set, since that logic lives inside the Radix package either way.

### 3. Continuous measurement or discrete state — decisive answer

**Discrete, event-driven state is sufficient; no per-frame tweening or
continuous RAF loop is needed to match the visual contract.** Evidence,
not inference from vibes:

- **No CSS `transition` targets `width` or `height` anywhere on the
  viewport**, checked exhaustively: `Viewport`'s own new-york-v4 class
  string (quoted in full, §1) carries only `data-[state=closed]:animate-out
  data-[state=closed]:zoom-out-95 data-[state=open]:animate-in
  data-[state=open]:zoom-in-90` — a `scale` transform animation gated on
  open/close *state*, not on content-switch. I grepped every one of the
  eight style families that ship a `.cn-navigation-menu-viewport` rule
  (`nova`, `rhea`, `luma`, `vega`, `sera`, `mira`, `maia`, `lyra` —
  `grep -rn "navigation-menu-viewport" --include=*.css`, all eight hits
  inspected) and nova's own rule (quoted below) is representative: `bg-
  popover text-popover-foreground data-open:animate-in data-closed:animate-
  out data-closed:zoom-out-90 data-open:zoom-in-90 ring-foreground/10
  rounded-lg shadow ring-1 duration-100` — again only a scale-in/out on
  open/close, `duration-100` applies to that transform, not to width/height.
  **If the library relied on a continuously-tweened width/height (CSS
  `transition: width, height` or a per-frame JS animation loop), some CSS
  in at least one of these eight independently-authored style families
  would need to declare a matching `transition-[width,height]` (or
  equivalent) to make that tween visible — none do.**
- Consequence: whatever internal mechanism Radix uses to *measure* the
  active content (very plausibly a `ResizeObserver` on the currently-shown
  `Content`, since content is arbitrary and its size isn't known until
  rendered — this specific claim about the measurement mechanism itself
  *is* `derived-not-read`, offered as the most plausible explanation, not
  a verified fact), the **application** of that measurement to the DOM is a
  discrete value snap on each content-switch event, not an animated tween.
  Task 4 needs to: (a) know each `NavigationMenuContent`'s natural
  width/height once it's the active one, and (b) write that pair of numbers
  onto the viewport element as two CSS custom properties (or direct
  inline `width`/`height`) at the moment the active item changes — a
  **discrete recompute triggered by the open/switch event**, not a
  continuous per-frame measurement loop.
- The one place continuous measurement *would* still matter, called out
  explicitly so Task 4 doesn't quietly regress it: content that changes its
  own size **after** it has already become the active/open content (e.g. an
  image finishes loading, or content is itself dynamically populated) will
  leave the viewport sized to its stale first measurement unless something
  re-measures. A `ResizeObserver` on the currently-active content (read-only
  measurement, firing on actual size changes — not a per-frame RAF poll)
  covers that case correctly without needing continuous polling. **Recommend
  Task 4 scope this as: one measurement on activation + a `ResizeObserver`
  on the active content for late-arriving resizes, explicitly not a
  requestAnimationFrame tweening loop** — the CSS evidence above rules out
  needing the latter for visual parity with what shadcn ships today.

### `.cn-navigation-menu-*` nova rules (verbatim, `style-nova.css` lines 855–897)

```css
.cn-navigation-menu {
  @apply max-w-max;
}
.cn-navigation-menu-list {
  @apply gap-0;
}
.cn-navigation-menu-trigger {
  @apply hover:bg-muted focus:bg-muted data-open:hover:bg-muted data-open:focus:bg-muted data-open:bg-muted/50 focus-visible:ring-ring/50 data-popup-open:bg-muted/50 data-popup-open:hover:bg-muted rounded-lg px-2.5 py-1.5 text-sm font-medium transition-all focus-visible:ring-3 focus-visible:outline-1 disabled:opacity-50;
}
.cn-navigation-menu-link {
  @apply data-active:focus:bg-muted data-active:hover:bg-muted data-active:bg-muted/50 focus-visible:ring-ring/50 hover:bg-muted focus:bg-muted flex items-center gap-2 rounded-lg p-2 text-sm transition-all outline-none focus-visible:ring-3 focus-visible:outline-1 in-data-[slot=navigation-menu-content]:rounded-md [&_svg:not([class*='size-'])]:size-4;
}
.cn-navigation-menu-trigger-icon {
  @apply relative top-px ml-1 size-3 transition duration-300 group-data-open/navigation-menu-trigger:rotate-180 group-data-popup-open/navigation-menu-trigger:rotate-180;
}
.cn-navigation-menu-content {
  @apply data-[motion^=from-]:animate-in data-[motion^=to-]:animate-out data-[motion^=from-]:fade-in data-[motion^=to-]:fade-out data-[motion=from-end]:slide-in-from-right-52 data-[motion=from-start]:slide-in-from-left-52 data-[motion=to-end]:slide-out-to-right-52 data-[motion=to-start]:slide-out-to-left-52 group-data-[viewport=false]/navigation-menu:bg-popover group-data-[viewport=false]/navigation-menu:text-popover-foreground group-data-[viewport=false]/navigation-menu:data-open:animate-in group-data-[viewport=false]/navigation-menu:data-closed:animate-out group-data-[viewport=false]/navigation-menu:data-closed:zoom-out-95 group-data-[viewport=false]/navigation-menu:data-open:zoom-in-95 group-data-[viewport=false]/navigation-menu:data-open:fade-in-0 group-data-[viewport=false]/navigation-menu:data-closed:fade-out-0 group-data-[viewport=false]/navigation-menu:ring-foreground/10 p-1 ease-[cubic-bezier(0.22,1,0.36,1)] group-data-[viewport=false]/navigation-menu:rounded-lg group-data-[viewport=false]/navigation-menu:shadow group-data-[viewport=false]/navigation-menu:ring-1 group-data-[viewport=false]/navigation-menu:duration-300;
}
.cn-navigation-menu-viewport {
  @apply bg-popover text-popover-foreground data-open:animate-in data-closed:animate-out data-closed:zoom-out-90 data-open:zoom-in-90 ring-foreground/10 rounded-lg shadow ring-1 duration-100;
}
.cn-navigation-menu-indicator {
  @apply data-[state=visible]:animate-in data-[state=hidden]:animate-out data-[state=hidden]:fade-out data-[state=visible]:fade-in;
}
.cn-navigation-menu-indicator-arrow {
  @apply bg-border rounded-tl-sm shadow-md;
}
```
(`.cn-navigation-menu-positioner` and `.cn-navigation-menu-popup`, lines
891–897, exist in the same nova block but belong to Base UI's differently-
shaped navigation-menu primitive per the `bases/base` family — same
provenance caveat as `## shared-items`' combobox-aria note in batch-A;
excluded from the Radix-based port's citations.)

Metric deltas worth naming (new-york-v4 → nova): `rounded-md→rounded-lg`
(trigger, content-in-non-viewport-mode, viewport), `focus-visible:ring-
[3px]→focus-visible:ring-3` (spelling only, same value, matches the
`## select`/`## combobox` precedent), `duration-200→duration-300` (content,
non-viewport-mode close/open), `zoom-out-95→zoom-out-90` (viewport close,
now matching its own `zoom-in-90` open value exactly — new-york-v4's
viewport is asymmetric, `zoom-in-90`/`zoom-out-95`; nova makes both `90`).
`.cn-navigation-menu-trigger` drops `bg-background` entirely (new-york-v4's
`navigationMenuTriggerStyle` has it; nova's trigger rule has no
background-color utility at all, relying on the ambient page background
instead) — a real, not just metric, style choice; port nova's version as
written per the visual-reference rule, don't silently reintroduce
`bg-background`.

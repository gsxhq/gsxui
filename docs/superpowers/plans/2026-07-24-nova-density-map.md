# Nova density map (generated 2026-07-24)

Metric-only diff of every `ui/*.gsx` component against shadcn's nova style.
Nova source of truth: `/Users/jackieli/personal/shadcn-ui/apps/v4/registry/styles/style-nova.css`
(line numbers below refer to that file). Structural markup source:
`/Users/jackieli/personal/shadcn-ui/apps/v4/registry/bases/radix/ui/*.tsx`.

Scope: h-*, w-*, size-*, min-w-*, max-h-*, p*/px/py/pl/pr/pt/pb, gap-*, space-*,
font sizes, rounded-*, svg sizes, numeric arbitrary values. Colors/rings/shadows/
animations are excluded except shadow-PRESENCE changes, noted per component.

Global nova patterns worth internalizing before per-component work:

- **buttons gain `border border-transparent bg-clip-padding` in the base** (L150) —
  every variant now includes a 1px transparent border inside the stated height,
  so nova `h-8` is content-box-equivalent to old borderless `h-8` minus 2px.
  Outline variant just recolors that border; box size no longer changes per variant.
- **press feel**: `.cn-button` adds `active:not-aria-[haspopup]:translate-y-px` (L150).
- **radius scale shifts up one step**: controls that were `rounded-md` are now
  `rounded-lg` (button, input, textarea, native select, toggle, input-group,
  menus content); overlays that were `rounded-lg` are now `rounded-xl`
  (dialog, alert-dialog, card was already xl, empty).
- **icon-adjacent padding**: nova replaces gsxui's `has-[>svg]:px-*` with
  directional `has-data-[icon=inline-start]:pl-*` / `has-data-[icon=inline-end]:pr-*`
  — requires stamping `data-icon="inline-start|inline-end"` on icon children to
  activate (markup change, not just class change).
- **popover-family surfaces swap `border` for `ring-1 ring-foreground/10`**
  (dropdown/context/select/popover/hover-card content, dialog, card). Ring is
  outside the border-box, so removing `border` changes inner box by 1px per side.
- **footers become full-bleed bands** in dialog/alert-dialog/card:
  `bg-muted/50 -mx-4 -mb-4 rounded-b-xl border-t p-4` (dialog L479-481,
  alert-dialog L69-71) / `p-(--card-spacing)` (card L264-266).

Summary: **27 components with deltas**, **8 no delta**, **4 no nova counterpart**
(aspect-ratio, collapsible, spinner, icon).

---

## accordion

- current:
  - trigger (`AccordionTrigger`): `gap-4 rounded-md py-4 text-left text-sm font-medium`; chevron `size-4`
  - content (`AccordionContent` outer): `text-sm`; inner: `pt-0 pb-4`
  - item: no metric tokens (`border-b` only)
- nova:
  - `.cn-accordion-trigger` (L7-9): `rounded-lg py-2.5 text-left text-sm font-medium` + `**:data-[slot=accordion-trigger-icon]:ml-auto **:data-[slot=accordion-trigger-icon]:size-4`
  - `.cn-accordion-content` (L11-13): `text-sm`
  - `.cn-accordion-content-inner` (L15-17): `pt-0 pb-2.5`
  - `.cn-accordion-item` (L3-5): no metric tokens (`not-last:border-b`)
- delta:
  - trigger: `py-4 → py-2.5`
  - trigger: `rounded-md → rounded-lg`
  - trigger: `gap-4 → (remove)` — nova spaces the chevron with `ml-auto` on the icon instead of a flex gap
  - content inner: `pb-4 → pb-2.5`
- notes: nova targets the chevron via `data-[slot=accordion-trigger-icon]` — gsxui's chevron carries no such slot; either stamp it or keep sizing the icon directly. Border membership flips from `border-b last:border-b-0` to `not-last:border-b` (same visual). No shadow changes.

## alert

- current:
  - root: `grid-cols-[0_1fr] gap-y-0.5 rounded-lg px-4 py-3 text-sm has-[>svg]:grid-cols-[calc(var(--spacing)*4)_1fr] has-[>svg]:gap-x-3 [&>svg]:size-4`
  - title: `line-clamp-1 min-h-4 font-medium tracking-tight`
  - description: `gap-1 text-sm`
- nova:
  - `.cn-alert` (L20-22): `grid gap-0.5 rounded-lg px-2.5 py-2 text-sm has-data-[slot=alert-action]:pr-18 has-[>svg]:grid-cols-[auto_1fr] has-[>svg]:gap-x-2 *:[svg]:row-span-2 *:[svg]:translate-y-0.5 *:[svg:not([class*='size-'])]:size-4`
  - `.cn-alert-title` (L32-34): `font-medium` (no metric tokens)
  - `.cn-alert-description` (L36-38): `text-sm [&_p:not(:last-child)]:mb-4`
  - `.cn-alert-action` (L40-42): `absolute top-2 right-2` — NEW part
- delta:
  - root: `px-4 → px-2.5`, `py-3 → py-2`
  - root: `has-[>svg]:gap-x-3 → has-[>svg]:gap-x-2`
  - root: `grid-cols-[0_1fr]` + `has-[>svg]:grid-cols-[calc(var(--spacing)*4)_1fr] → plain grid` + `has-[>svg]:grid-cols-[auto_1fr]`
  - root: `[&>svg]:size-4 → *:[svg:not([class*='size-'])]:size-4` (size now caller-overridable); add `*:[svg]:row-span-2`
  - root: add `has-data-[slot=alert-action]:pr-18` (reserves space for the new action slot) and `text-left`
  - title: `min-h-4 → (remove)`; drop `tracking-tight`, `line-clamp-1`
  - description: `gap-1 → (remove)`; add `[&_p:not(:last-child)]:mb-4`
  - NEW part: `AlertAction` — `.cn-alert-action` `absolute top-2 right-2`
- notes: no shadow changes (neither has one). Nova keeps `bg-card` variants — colors out of scope.

## alert-dialog

- current (composes DialogContent):
  - content: `max-w-[calc(100%-2rem)] gap-4 rounded-lg p-6 sm:max-w-lg` (+ `border shadow-lg` from dialog)
  - header: `gap-1.5`
  - footer: `gap-2` (flex-col-reverse sm:flex-row)
  - title: `text-lg font-semibold`
  - description: `text-sm`
- nova:
  - `.cn-alert-dialog-content` (L49-51): `gap-4 rounded-xl p-4 ring-1 data-[size=default]:max-w-xs data-[size=sm]:max-w-xs data-[size=default]:sm:max-w-sm`
  - `.cn-alert-dialog-header` (L53-55): `gap-1.5` (+ media-conditional grid rows)
  - `.cn-alert-dialog-media` (L57-59): `mb-2 size-10 rounded-md *:[svg:not([class*='size-'])]:size-6` — NEW part
  - `.cn-alert-dialog-title` (L61-63): `text-base font-medium`
  - `.cn-alert-dialog-description` (L65-67): `text-sm`
  - `.cn-alert-dialog-footer` (L69-71): `-mx-4 -mb-4 rounded-b-xl border-t p-4`
- delta:
  - content: `p-6 → p-4`, `rounded-lg → rounded-xl`
  - content: `sm:max-w-lg → max-w-xs` + `sm:max-w-sm` (default size; alert dialogs get much narrower)
  - title: `text-lg → text-base` (weight `font-semibold → font-medium`, feel-relevant)
  - footer: `gap-2 → -mx-4 -mb-4 rounded-b-xl border-t p-4` full-bleed band (keep gap-2 for its buttons)
  - NEW part: `AlertDialogMedia` (`size-10 rounded-md mb-2`, icon `size-6`)
  - NEW: `size` variant axis (`default` | `sm`) on content
- notes: shadow presence change — gsxui inherits dialog's `shadow-lg`; nova has NO shadow on alert-dialog content (ring-1 only). Overlay becomes `bg-black/10` + blur (color, but scrim feel changes).

## aspect-ratio

- no nova counterpart — keep as is. No `.cn-aspect-ratio` entry in style-nova.css; base `aspect-ratio.tsx` is a classless passthrough, and gsxui's port carries no metric tokens.

## avatar

- current:
  - root: `size-8 rounded-full`
  - image: `size-full`
  - fallback: `size-full rounded-full text-sm`
- nova:
  - `.cn-avatar` (L74-76): `size-8 rounded-full data-[size=lg]:size-10 data-[size=sm]:size-6`
  - `.cn-avatar-fallback` (L78-80): `rounded-full` (no text-size token)
  - `.cn-avatar-image` (L82-84): `rounded-full` (no metric)
  - `.cn-avatar-group-count` (L90-92): `size-8 rounded-full text-sm [&>svg]:size-4` + lg `size-10 [&>svg]:size-5` / sm `size-6 [&>svg]:size-3` — NEW part
- delta: no delta on existing metrics (`size-8 rounded-full` matches).
  - NEW size variants: `data-[size=lg]:size-10`, `data-[size=sm]:size-6`
  - NEW parts: `AvatarBadge`, `AvatarGroupCount` (metrics above)
- notes: nova omits `text-sm` on fallback (it lives on group-count instead); gsxui's `text-sm` on fallback is a harmless holdover — decide whether to keep. No shadow changes.

## badge

- current: `gap-1 rounded-full border border-transparent px-2 py-0.5 text-xs font-medium [&>svg]:size-3`
- nova: `.cn-badge` (L95-97): `h-5 gap-1 rounded-4xl border border-transparent px-2 py-0.5 text-xs font-medium has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 [&>svg]:size-3!`
- delta:
  - add: `h-5` (explicit height; gsxui height is content-derived)
  - `rounded-full → rounded-4xl`
  - add: `has-data-[icon=inline-start]:pl-1.5`, `has-data-[icon=inline-end]:pr-1.5` (needs `data-icon` stamping on icon children)
  - `[&>svg]:size-3 → [&>svg]:size-3!` (now important)
- notes: nova destructive variant is tinted (`bg-destructive/10`) not solid — color, out of scope. No shadow changes.

## breadcrumb

- current:
  - list: `gap-1.5 text-sm sm:gap-2.5`
  - item: `gap-1.5`
  - separator: `[&>svg]:size-3.5`
  - ellipsis: `size-9` with inner icon `size-4`
- nova:
  - `.cn-breadcrumb-list` (L124-126): `gap-1.5 text-sm`
  - `.cn-breadcrumb-item` (L128-130): `gap-1`
  - `.cn-breadcrumb-separator` (L140-142): `[&>svg]:size-3.5`
  - `.cn-breadcrumb-ellipsis` (L144-146): `size-5 [&>svg]:size-4`
- delta:
  - list: `sm:gap-2.5 → (remove)` (flat `gap-1.5` at all breakpoints)
  - item: `gap-1.5 → gap-1`
  - ellipsis: `size-9 → size-5` (icon stays size-4)
- notes: link/page have no metric tokens either side. No shadow changes.

## button

- current:
  - base: `gap-2 rounded-md text-sm font-medium [&_svg:not([class*='size-'])]:size-4`
  - default: `h-9 px-4 py-2 has-[>svg]:px-3`
  - sm: `h-8 gap-1.5 rounded-md px-3 has-[>svg]:px-2.5`
  - xs: `h-6 gap-1 rounded-md px-2 text-xs has-[>svg]:px-1.5 [&_svg:not([class*='size-'])]:size-3`
  - lg: `h-10 rounded-md px-6 has-[>svg]:px-4`
  - icon: `size-9`; icon-xs: `size-6 rounded-md [&_svg...]:size-3`; icon-sm: `size-8`; icon-lg: `size-10`
- nova:
  - `.cn-button` (L149-151): `rounded-lg border border-transparent bg-clip-padding text-sm font-medium active:not-aria-[haspopup]:translate-y-px [&_svg:not([class*='size-'])]:size-4` — no base gap
  - `.cn-button-size-default` (L185-187): `h-8 gap-1.5 px-2.5 has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2`
  - `.cn-button-size-sm` (L181-183): `h-7 gap-1 rounded-[min(var(--radius-md),12px)] px-2.5 text-[0.8rem] has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 [&_svg:not([class*='size-'])]:size-3.5`
  - `.cn-button-size-xs` (L177-179): `h-6 gap-1 rounded-[min(var(--radius-md),10px)] px-2 text-xs in-data-[slot=button-group]:rounded-lg has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 [&_svg:not([class*='size-'])]:size-3`
  - `.cn-button-size-lg` (L189-191): `h-9 gap-1.5 px-2.5 has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2`
  - `.cn-button-size-icon` (L201-203): `size-8`
  - `.cn-button-size-icon-xs` (L193-195): `size-6 rounded-[min(var(--radius-md),10px)] in-data-[slot=button-group]:rounded-lg [&_svg:not([class*='size-'])]:size-3`
  - `.cn-button-size-icon-sm` (L197-199): `size-7 rounded-[min(var(--radius-md),12px)] in-data-[slot=button-group]:rounded-lg`
  - `.cn-button-size-icon-lg` (L205-207): `size-9`
- delta:
  - base: `rounded-md → rounded-lg`
  - base: add `border border-transparent bg-clip-padding` — ALL variants gain a transparent border; affects box size, and outline variant no longer changes the box
  - base: add `active:not-aria-[haspopup]:translate-y-px`
  - base: `gap-2 → (remove; gap moves into sizes)`
  - default: `h-9 → h-8`, `px-4 → px-2.5`, `py-2 → (remove)`, effective gap `2 → 1.5`, `has-[>svg]:px-3 → has-data-[icon=inline-start]:pl-2 / has-data-[icon=inline-end]:pr-2`
  - sm: `h-8 → h-7`, `px-3 → px-2.5`, `gap-1.5 → gap-1`, `rounded-md → rounded-[min(var(--radius-md),12px)]`, add `text-[0.8rem]`, add `[&_svg...]:size-3.5`, `has-[>svg]:px-2.5 → pl/pr-1.5 directional`
  - xs: `rounded-md → rounded-[min(var(--radius-md),10px)]` (+ `in-data-[slot=button-group]:rounded-lg`); h-6/gap-1/px-2/text-xs/svg-3 unchanged; `has-[>svg]:px-1.5 → pl/pr-1.5 directional`
  - lg: `h-10 → h-9`, `px-6 → px-2.5`, add `gap-1.5`
  - icon: `size-9 → size-8`; icon-sm: `size-8 → size-7` + `rounded-[min(var(--radius-md),12px)]`; icon-lg: `size-10 → size-9`; icon-xs: size-6 unchanged, `rounded-md → rounded-[min(var(--radius-md),10px)]`
- notes: shadow presence — gsxui outline variant has `shadow-xs`; nova outline has NO shadow. The directional icon paddings require `data-icon="inline-start|inline-end"` stamps on icon children (markup change). `rounded-[min(var(--radius-md),...)]` caps small-size radii under themable `--radius-md`.

## button-group

- current:
  - root: `has-[>[data-slot=button-group]]:gap-2 has-[select...]:...:rounded-r-md`; orientation via `[&>*:not(:first-child)]:rounded-l-none [&>*:not(:first-child)]:border-l-0 [&>*:not(:last-child)]:rounded-r-none` (vertical analog: rounded-t/b-none border-t-0)
  - text: `gap-2 rounded-md px-4 text-sm font-medium [&_svg:not([class*='size-'])]:size-4`
- nova:
  - `.cn-button-group` (L210-212): `has-[>[data-slot=button-group]]:gap-2 has-[select...]:[&>[data-slot=select-trigger]:last-of-type]:rounded-r-lg`
  - `.cn-button-group-orientation-horizontal` (L214-216): `[&>[data-slot]:not(:has(~[data-slot]))]:rounded-r-lg!`
  - `.cn-button-group-orientation-vertical` (L218-220): `[&>[data-slot]:not(:has(~[data-slot]))]:rounded-b-lg!`
  - `.cn-button-group-text` (L222-224): `gap-2 rounded-lg border px-2.5 text-sm font-medium [&_svg:not([class*='size-'])]:size-4`
- delta:
  - root: `rounded-r-md → rounded-r-lg` (select-trigger last-child rule)
  - text: `px-4 → px-2.5`, `rounded-md → rounded-lg`
- notes: orientation mechanism inverted — gsxui zeroes inner corners (`rounded-*-none` on non-edge children); nova instead force-restores the outer corner on the true last slotted child (`rounded-r-lg!` / `rounded-b-lg!`) and relies on button xs/sm sizes' `in-data-[slot=button-group]:rounded-lg` for group-context radii. Porting needs the nova pattern, not a token swap. Shadow presence: text loses `shadow-xs`.

## card

- current:
  - root: `gap-6 rounded-xl py-6` (+ `border shadow-sm`)
  - header: `gap-2 px-6 [.border-b]:pb-6`
  - title: `leading-none font-semibold` (no size)
  - description: `text-sm`
  - content: `px-6`
  - footer: `px-6 [.border-t]:pt-6`
- nova:
  - `.cn-card` (L244-246): `gap-(--card-spacing) rounded-xl py-(--card-spacing) text-sm ring-1 [--card-spacing:--spacing(4)] has-data-[slot=card-footer]:pb-0 has-[>img:first-child]:pt-0 data-[size=sm]:[--card-spacing:--spacing(3)] *:[img:first-child]:rounded-t-xl *:[img:last-child]:rounded-b-xl`
  - `.cn-card-header` (L248-250): `gap-1 rounded-t-xl px-(--card-spacing) [.border-b]:pb-(--card-spacing)`
  - `.cn-card-title` (L252-254): `text-base leading-snug font-medium group-data-[size=sm]/card:text-sm`
  - `.cn-card-description` (L256-258): `text-sm`
  - `.cn-card-content` (L260-262): `px-(--card-spacing)`
  - `.cn-card-footer` (L264-266): `rounded-b-xl border-t p-(--card-spacing)` (bg-muted/50)
- delta:
  - root: `py-6 → py-4` / `gap-6 → gap-4` (via `--card-spacing: --spacing(4)`); add `text-sm`; add `has-data-[slot=card-footer]:pb-0`
  - header: `px-6 → px-4`, `gap-2 → gap-1`, `[.border-b]:pb-6 → pb-4`
  - title: add `text-base`; `leading-none → leading-snug` (weight `font-semibold → font-medium`)
  - content: `px-6 → px-4`
  - footer: `px-6 [.border-t]:pt-6 → border-t p-4` full-width band (root's `pb-0` makes it flush-bottom, `rounded-b-xl`)
  - NEW: `data-[size=sm]` variant (`--card-spacing: --spacing(3)` = p-3/gap-3, title text-sm)
- notes: shadow presence — `shadow-sm` removed; `border` → `ring-1 ring-foreground/10`. Nova parameterizes all card spacing through `--card-spacing`; porting literally as p-4/gap-4 works but loses the sm variant lever.

## checkbox

- current: `size-4 rounded-[4px]`; check glyph painted at `bg-[length:12px_12px]`
- nova:
  - `.cn-checkbox` (L283-285): `size-4 rounded-[4px]`
  - `.cn-checkbox-indicator` (L291-293): `[&>svg]:size-3.5`
- delta:
  - check glyph: `bg-[length:12px_12px] → bg-[length:14px_14px]` (nova indicator svg is size-3.5 = 14px; gsxui paints a 12px = size-3 check)
- notes: box metrics identical. Shadow presence: gsxui `shadow-xs` — nova has none.

## collapsible

- no nova counterpart — keep as is. No `.cn-collapsible*` entries; base `collapsible.tsx` is a classless passthrough, matching gsxui's classless port (trigger's `list-none`/marker suppression is a `<summary>` mechanism, not density).

## context-menu

- current:
  - content: `max-h-96 min-w-[8rem] rounded-md p-1` (+ `border shadow-md`)
  - item: `gap-2 rounded-sm px-2 py-1.5 text-sm [&_svg:not([class*='size-'])]:size-4`
  - label: `px-2 py-1.5 text-sm font-medium`
  - separator: `-mx-1 my-1 h-px`
  - shortcut: `ml-auto text-xs tracking-widest`
- nova:
  - `.cn-context-menu-content` (L414-416): `min-w-36 rounded-lg p-1 shadow-md ring-1`
  - `.cn-context-menu-item` (L422-424): `gap-1.5 rounded-md px-1.5 py-1 text-sm data-inset:pl-7 [&_svg:not([class*='size-'])]:size-4`
  - `.cn-context-menu-label` (L438-440): `px-1.5 py-1 text-xs font-medium data-inset:pl-7`
  - `.cn-context-menu-separator` (L442-444): `-mx-1 my-1 h-px`
  - `.cn-context-menu-shortcut` (L446-448): `ml-auto text-xs tracking-widest`
  - (sub-trigger/sub-content L450-456 exist upstream; gsxui has no sub-menu parts)
- delta:
  - content: `min-w-[8rem] → min-w-36` (128px → 144px), `rounded-md → rounded-lg`
  - item: `gap-2 → gap-1.5`, `rounded-sm → rounded-md`, `px-2 → px-1.5`, `py-1.5 → py-1`
  - label: `px-2 → px-1.5`, `py-1.5 → py-1`, `text-sm → text-xs`
  - separator/shortcut: no delta
- notes: content `border → ring-1` (1px inner box change); shadow-md kept. gsxui's `max-h-96` is a popover-port addition with no nova counterpart — keep. Note context-menu content is `min-w-36` while dropdown stays `min-w-32` — they diverge in nova.

## dialog

- current:
  - content: `max-w-[calc(100%-2rem)] gap-4 rounded-lg p-6 sm:max-w-lg` (+ `border shadow-lg`)
  - close button: `top-4 right-4` (bare button, `[&_svg...]:size-4`)
  - header: `gap-2`
  - footer: `gap-2`
  - title: `text-lg leading-none font-semibold`
  - description: `text-sm`
- nova:
  - `.cn-dialog-content` (L467-469): `max-w-[calc(100%-2rem)] gap-4 rounded-xl p-4 text-sm ring-1 sm:max-w-sm`
  - `.cn-dialog-close` (L471-473): `absolute top-2 right-2`
  - `.cn-dialog-header` (L475-477): `gap-2`
  - `.cn-dialog-footer` (L479-481): `-mx-4 -mb-4 rounded-b-xl border-t p-4`
  - `.cn-dialog-title` (L483-485): `text-base leading-none font-medium`
  - `.cn-dialog-description` (L487-489): `text-sm`
- delta:
  - content: `p-6 → p-4`, `rounded-lg → rounded-xl`, `sm:max-w-lg → sm:max-w-sm`, add `text-sm`
  - close: `top-4 right-4 → top-2 right-2`
  - title: `text-lg → text-base` (weight `font-semibold → font-medium`)
  - footer: add `-mx-4 -mb-4 rounded-b-xl border-t p-4` full-bleed band (keep gap-2 inside)
- notes: shadow presence — `shadow-lg` removed; `border` → `ring-1`. In nova's base `dialog.tsx` the close X is a `Button variant="ghost" size="icon-sm"` (nova icon-sm = `size-7 rounded-[min(var(--radius-md),12px)]`) positioned by `.cn-dialog-close` — gsxui's bare unpadded button will need the icon-sm button treatment for the same hit target.

## dropdown

- current:
  - content: `max-h-96 min-w-[8rem] rounded-md p-1` (+ `border shadow-md`)
  - item: `gap-2 rounded-sm px-2 py-1.5 text-sm [&_svg:not([class*='size-'])]:size-4`
  - label: `px-2 py-1.5 text-sm font-medium`
  - separator: `-mx-1 my-1 h-px`
  - shortcut: `ml-auto text-xs tracking-widest`
- nova:
  - `.cn-dropdown-menu-content` (L537-539): `min-w-32 rounded-lg p-1 shadow-md ring-1`
  - `.cn-dropdown-menu-item` (L545-547): `gap-1.5 rounded-md px-1.5 py-1 text-sm data-inset:pl-7 [&_svg:not([class*='size-'])]:size-4`
  - `.cn-dropdown-menu-label` (L561-563): `px-1.5 py-1 text-xs font-medium data-inset:pl-7`
  - `.cn-dropdown-menu-separator` (L565-567): `-mx-1 my-1 h-px`
  - `.cn-dropdown-menu-shortcut` (L569-571): `ml-auto text-xs tracking-widest`
  - (checkbox-item/radio-item L549-555: `gap-1.5 rounded-md py-1 pr-8 pl-1.5 text-sm`; sub-content L577-579 `min-w-[96px] rounded-lg p-1` — parts gsxui doesn't ship)
- delta:
  - content: `rounded-md → rounded-lg` (`min-w-[8rem]` == `min-w-32`, no width delta)
  - item: `gap-2 → gap-1.5`, `rounded-sm → rounded-md`, `px-2 → px-1.5`, `py-1.5 → py-1`
  - label: `px-2 → px-1.5`, `py-1.5 → py-1`, `text-sm → text-xs`
  - separator/shortcut: no delta
- notes: content `border → ring-1`; shadow-md kept. `max-h-96` is a gsxui popover-port addition — keep. `data-inset:pl-7` replaces the dropped `data-[inset]:pl-8` (gsxui dropped inset entirely; if reinstated, use pl-7).

## empty

- current:
  - root: `gap-6 rounded-lg p-6 md:p-12`
  - header: `max-w-sm gap-2`
  - media (default): `mb-2`; media icon variant: `size-10 rounded-lg [&_svg:not([class*='size-'])]:size-6`
  - title: `text-lg font-medium tracking-tight`
  - description: `text-sm/relaxed`
  - content: `max-w-sm gap-4 text-sm`
- nova:
  - `.cn-empty` (L586-588): `gap-4 rounded-xl border-dashed p-6`
  - `.cn-empty-header` (L590-592): `gap-2`
  - `.cn-empty-media` (L594-596): `mb-2`
  - `.cn-empty-media-icon` (L602-604): `size-8 rounded-lg [&_svg:not([class*='size-'])]:size-4`
  - `.cn-empty-title` (L606-608): `text-sm font-medium tracking-tight`
  - `.cn-empty-description` (L610-612): `text-sm/relaxed`
  - `.cn-empty-content` (L614-616): `gap-2.5 text-sm`
- delta:
  - root: `gap-6 → gap-4`, `rounded-lg → rounded-xl`, `md:p-12 → (remove)` (flat p-6)
  - media icon: `size-10 → size-8`, svg `size-6 → size-4`
  - title: `text-lg → text-sm`
  - content: `gap-4 → gap-2.5`
- notes: nova css has no `max-w-sm` on header/content (widths live in base tsx if anywhere) — verify against `empty.tsx` before dropping. No shadow changes.

## field

- current:
  - fieldset: `gap-6 has-[>[data-slot=checkbox-group]]:gap-3 has-[>[data-slot=radio-group]]:gap-3`
  - legend: `mb-3 data-[variant=legend]:text-base data-[variant=label]:text-sm`
  - group: `gap-7 data-[slot=checkbox-group]:gap-3 [&>[data-slot=field-group]]:gap-4`
  - field: `gap-3`
  - content: `gap-1.5`
  - label: `gap-2 [&>*]:data-[slot=field]:p-4 has-[>[data-slot=field]]:rounded-md`
  - title: `gap-2 text-sm`
  - description: `text-sm nth-last-2:-mt-1 [[data-variant=legend]+&]:-mt-1.5`
  - separator: `-my-2 h-5 text-sm`; separator content: `px-2`
  - error: `text-sm`
- nova:
  - `.cn-field-set` (L619-621): `gap-4 has-[>[data-slot=checkbox-group]]:gap-3 has-[>[data-slot=radio-group]]:gap-3`
  - `.cn-field-legend` (L623-625): `mb-1.5 data-[variant=label]:text-sm data-[variant=legend]:text-base`
  - `.cn-field-group` (L627-629): `gap-5 data-[slot=checkbox-group]:gap-3 *:data-[slot=field-group]:gap-4`
  - `.cn-field` (L631-633): `gap-2`
  - `.cn-field-content` (L635-637): `gap-0.5`
  - `.cn-field-label` (L639-641): `gap-2 has-[>[data-slot=field]]:rounded-lg *:data-[slot=field]:p-2.5`
  - `.cn-field-title` (L647-649): `gap-2 text-sm`
  - `.cn-field-description` (L651-653): `text-sm [[data-variant=legend]+&]:-mt-1.5`
  - `.cn-field-separator` (L655-657): `-my-2 h-5 text-sm`; `.cn-field-separator-content` (L659-661): `px-2`
  - `.cn-field-error` (L663-665): `text-sm`
- delta:
  - fieldset: `gap-6 → gap-4`
  - legend: `mb-3 → mb-1.5`
  - group: `gap-7 → gap-5`
  - field: `gap-3 → gap-2`
  - content: `gap-1.5 → gap-0.5`
  - label: `p-4 → p-2.5` (nested field), `rounded-md → rounded-lg`
  - description: `nth-last-2:-mt-1 → (remove)` (nova keeps only the legend-adjacent `-mt-1.5`)
  - separator/title/error: no delta
- notes: no shadow changes.

## hover-card

- current: `w-64 rounded-md p-4` (+ `border shadow-md`)
- nova: `.cn-hover-card-content` (L668-670): `w-64 rounded-lg p-2.5 text-sm shadow-md ring-1`
- delta:
  - `rounded-md → rounded-lg`
  - `p-4 → p-2.5`
  - add: `text-sm`
- notes: `border → ring-1`; shadow-md kept (no presence change).

## input

- current: `h-9 rounded-md px-3 py-1 text-base file:h-7 file:text-sm md:text-sm` (+ `border shadow-xs`)
- nova: `.cn-input` (L677-679): `h-8 rounded-lg px-2.5 py-1 text-base file:h-6 file:text-sm md:text-sm` (border kept)
- delta:
  - `h-9 → h-8`
  - `rounded-md → rounded-lg`
  - `px-3 → px-2.5`
  - `file:h-7 → file:h-6`
- notes: shadow presence — `shadow-xs` removed in nova.

## input-group

- current:
  - root: `h-9 rounded-md` (+ `border shadow-xs`); `has-[>[data-align=inline-start]]:[&>input]:pl-2`, `inline-end...pr-2`, `block-start...[&>input]:pb-3`, `block-end...[&>input]:pt-3`
  - addon base: `py-1.5 text-sm [&>kbd]:rounded-[calc(var(--radius)-5px)] [&>svg:not([class*='size-'])]:size-4`
  - addon inline-start: `pl-3 has-[>button]:ml-[-0.45rem] has-[>kbd]:ml-[-0.35rem]`; inline-end mirror `pr-3 / -0.45 / -0.35`
  - addon block-start: `px-3 pt-3 group-has-[>input]:pt-2.5 [.border-b]:pb-3`; block-end mirror `px-3 pb-3 / pb-2.5 / pt-3`
  - button default(xs): `h-6 gap-1 rounded-[calc(var(--radius)-5px)] px-2 [&>svg...]:size-3.5`; sm: `h-8 gap-1.5 rounded-md px-2.5`; icon-xs: `size-6 rounded-[calc(var(--radius)-5px)] p-0`; icon-sm: `size-8 p-0`
  - text: `gap-2 text-sm [&_svg...]:size-4`
  - input child: `rounded-none` overrides; textarea child: `py-3`
- nova:
  - `.cn-input-group` (L1397-1399): `h-8 rounded-lg`; `has-[>[data-align=inline-start]]:[&>input]:pl-1.5`, `inline-end...pr-1.5`, `block-start...[&>input]:pb-3`, `block-end...[&>input]:pt-3`
  - `.cn-input-group-addon` (L1401-1403): `h-auto gap-2 py-1.5 text-sm [&>kbd]:rounded-[calc(var(--radius)-5px)] [&>svg:not([class*='size-'])]:size-4`
  - `.cn-input-group-addon-align-inline-start` (L1405-1407): `pl-2 has-[>button]:ml-[-0.3rem] has-[>kbd]:ml-[-0.15rem]`
  - `.cn-input-group-addon-align-inline-end` (L1409-1411): `pr-2 has-[>button]:mr-[-0.3rem] has-[>kbd]:mr-[-0.15rem]`
  - `.cn-input-group-addon-align-block-start` (L1413-1415): `px-2.5 pt-2 group-has-[>input]/input-group:pt-2 [.border-b]:pb-2`
  - `.cn-input-group-addon-align-block-end` (L1417-1419): `px-2.5 pb-2 group-has-[>input]/input-group:pb-2 [.border-t]:pt-2`
  - `.cn-input-group-button` (L1421-1423): `gap-2 text-sm`
  - `.cn-input-group-button-size-xs` (L1425-1427): `h-6 gap-1 rounded-[calc(var(--radius)-3px)] px-1.5 [&>svg:not([class*='size-'])]:size-3.5`
  - `.cn-input-group-button-size-icon-xs` (L1429-1431): `size-6 rounded-[calc(var(--radius)-3px)] p-0`
  - `.cn-input-group-button-size-icon-sm` (L1433-1435): `size-8 p-0`
  - `.cn-input-group-text` (L1437-1439): `gap-2 text-sm [&_svg:not([class*='size-'])]:size-4`
  - `.cn-input-group-textarea` (L1445-1447): `py-2`
- delta:
  - root: `h-9 → h-8`, `rounded-md → rounded-lg`
  - root input padding next to inline addons: `pl-2 → pl-1.5`, `pr-2 → pr-1.5`
  - addon inline-start: `pl-3 → pl-2`, `ml-[-0.45rem] → ml-[-0.3rem]`, kbd `ml-[-0.35rem] → ml-[-0.15rem]` (inline-end mirrors)
  - addon block-start: `px-3 → px-2.5`, `pt-3 → pt-2`, `group-has pt-2.5 → pt-2`, `[.border-b]:pb-3 → pb-2` (block-end mirrors)
  - button xs (default): `px-2 → px-1.5`, `rounded-[calc(var(--radius)-5px)] → rounded-[calc(var(--radius)-3px)]`
  - button icon-xs: `rounded-[calc(var(--radius)-5px)] → rounded-[calc(var(--radius)-3px)]`
  - textarea child: `py-3 → py-2`
- notes: gsxui's button `sm` case (`h-8 gap-1.5 rounded-md px-2.5`) has NO nova size entry — nova ships only xs/icon-xs/icon-sm; decide to drop or map to plain Button sm. Shadow presence: root `shadow-xs` removed. Addon `[&>kbd]:rounded-[calc(var(--radius)-5px)]` unchanged (only the BUTTON radius arithmetic changed).

## item

- current:
  - root base: `rounded-md text-sm`; size default: `gap-4 p-4`; size sm: `gap-2.5 px-4 py-3`
  - media: `gap-2`; icon variant: `size-8 rounded-sm border bg-muted [&_svg...]:size-4`; image variant: `size-10 rounded-sm`
  - content: `gap-1`; title: `text-sm`; description: `text-sm`
  - actions/header/footer: `gap-2`
  - separator (`ItemSeparator`): `my-0`
- nova:
  - `.cn-item` (L703-705): `rounded-lg text-sm`
  - `.cn-item-size-default` (L719-721): `gap-2.5 px-3 py-2.5`
  - `.cn-item-size-sm` (L723-725): `gap-2.5 px-3 py-2.5`
  - `.cn-item-size-xs` (L727-729): `gap-2 px-2.5 py-2 in-data-[slot=dropdown-menu-content]:p-0` — NEW
  - `.cn-item-media` (L731-733): `gap-2`
  - `.cn-item-media-variant-icon` (L739-741): `[&_svg:not([class*='size-'])]:size-4` (no box!)
  - `.cn-item-media-variant-image` (L743-745): `size-10 rounded-sm group-data-[size=sm]/item:size-8 group-data-[size=xs]/item:size-6`
  - `.cn-item-content` (L747-749): `gap-1 group-data-[size=xs]/item:gap-0`
  - `.cn-item-title` (L751-753): `gap-2 text-sm`
  - `.cn-item-description` (L755-757): `text-sm group-data-[size=xs]/item:text-xs`
  - `.cn-item-actions` / `-header` / `-footer` (L759-769): `gap-2`
  - `.cn-item-group` (L771-773): `gap-4 has-data-[size=sm]:gap-2.5 has-data-[size=xs]:gap-2` — NEW gaps
  - `.cn-item-separator` (L775-777): `my-2`
- delta:
  - root: `rounded-md → rounded-lg`
  - size default: `gap-4 p-4 → gap-2.5 px-3 py-2.5`
  - size sm: `px-4 py-3 → px-3 py-2.5` (now identical to default)
  - media icon: `size-8 → (remove box)` — nova's icon media has no sized/bordered container, just svg size-4
  - image variant: add responsive `sm:size-8` / `xs:size-6` group tokens
  - separator: `my-0 → my-2`
  - NEW: size `xs` (`gap-2 px-2.5 py-2`), ItemGroup gains gap tokens (`gap-4`, size-conditional 2.5/2)
- notes: the dropped icon-media box is a real visual change (bordered muted square → bare icon) — confirm intent before applying. No shadow changes.

## kbd

- current: `h-5 w-fit min-w-5 gap-1 rounded-sm px-1 text-xs [&_svg:not([class*='size-'])]:size-3`; group `gap-1`
- nova: `.cn-kbd` (L780-782): `h-5 w-fit min-w-5 gap-1 rounded-sm px-1 text-xs [&_svg:not([class*='size-'])]:size-3`; `.cn-kbd-group` (L784-786): `gap-1`
- delta: no delta.
- notes: tooltip-nested color tokens also match. No shadow changes.

## label

- current: `gap-2 text-sm leading-none font-medium`
- nova: `.cn-label` (L789-791): `gap-2 text-sm leading-none font-medium`
- delta: no delta.

## pagination

- current:
  - content: `gap-1`
  - ellipsis: `size-9` (icon size-4)
  - previous/next: `gap-1 px-2.5 sm:pl-2.5` / `sm:pr-2.5` on a size="default" button link
  - link: button sizes (default icon = old `size-9`)
- nova:
  - `.cn-pagination-content` (L909-911): `gap-0.5`
  - `.cn-pagination-ellipsis` (L913-915): `size-8 [&_svg:not([class*='size-'])]:size-4`
  - `.cn-pagination-previous` (L917-919): `pl-1.5!`
  - `.cn-pagination-next` (L921-923): `pr-1.5!`
- delta:
  - content: `gap-1 → gap-0.5`
  - ellipsis: `size-9 → size-8`
  - previous: `gap-1 px-2.5 sm:pl-2.5 → pl-1.5!` (rest of padding/gap comes from the button default size; the `sm:` responsive variant is gone)
  - next: `gap-1 px-2.5 sm:pr-2.5 → pr-1.5!`
  - link: inherits all button size deltas (icon `size-9 → size-8`, etc.)
- notes: no shadow changes.

## popover

- current: `w-72 rounded-md p-4` (+ `border shadow-md`)
- nova:
  - `.cn-popover-content` (L926-928): `flex flex-col gap-2.5 rounded-lg p-2.5 text-sm shadow-md ring-1` (base tsx keeps `w-72`)
  - `.cn-popover-header` (L934-936): `flex flex-col gap-0.5 text-sm` — NEW part
  - `.cn-popover-title` (L938-940): `font-medium` — NEW part
  - `.cn-popover-description` (L942-944): no metric — NEW part
- delta:
  - `rounded-md → rounded-lg`
  - `p-4 → p-2.5`
  - add: `gap-2.5` (content is now a flex column), `text-sm`
  - NEW parts: PopoverHeader (`gap-0.5 text-sm`), PopoverTitle, PopoverDescription
- notes: `w-72` unchanged (verified in base `popover.tsx`). `border → ring-1`; shadow-md kept.

## progress

- current: root `h-2 rounded-full`; indicator `h-full w-full`
- nova:
  - `.cn-progress` (L947-949): `h-1 rounded-full`
  - `.cn-progress-indicator` (L955-957): no metric (base tsx uses `size-full flex-1`)
  - `.cn-progress-label` (L959-961): `text-sm font-medium` — NEW part
  - `.cn-progress-value` (L963-965): `ml-auto text-sm` — NEW part
- delta:
  - root: `h-2 → h-1`
  - NEW parts: ProgressLabel / ProgressValue (text-sm)
- notes: no shadow changes.

## radio

- current: `size-4 rounded-full`; dot = radial-gradient at 45% of closest-side (~7.2px diameter)
- nova:
  - `.cn-radio-group-item` (L972-974): `size-4 rounded-full`
  - `.cn-radio-group-indicator-icon` (L984-986): `size-2` dot (8px)
  - `.cn-radio-group` (L968-970): `grid gap-2` (gsxui ships no group wrapper)
- delta: no delta on the control box. Dot: gsxui's gradient (~7.2px) vs nova `size-2` (8px) — optionally bump the gradient stop from 45% to 50% for exact parity.
- notes: shadow presence — gsxui `shadow-xs`; nova has none.

## select (native)

- current: select `h-9 rounded-md px-3 py-2 pr-8 text-sm appearance-none`; chevron `right-3 size-4`
- nova:
  - `.cn-native-select` (L900-902): `h-8 rounded-lg py-1 pr-8 pl-2.5 text-sm appearance-none data-[size=sm]:h-7 data-[size=sm]:rounded-[min(var(--radius-md),10px)] data-[size=sm]:py-0.5`
  - `.cn-native-select-icon` (L904-906): `right-2.5 size-4 top-1/2 -translate-y-1/2`
- delta:
  - `h-9 → h-8`
  - `rounded-md → rounded-lg`
  - `px-3 → pl-2.5` (pr-8 stays for the chevron gutter)
  - `py-2 → py-1`
  - chevron: `right-3 → right-2.5`
  - NEW: `data-[size=sm]` variant (`h-7 rounded-[min(var(--radius-md),10px)] py-0.5`)
- notes: gsxui maps to `.cn-native-select` (nova's styled native select), not `.cn-select-trigger` (the Radix listbox — that one is `data-[size=default]:h-8 py-2 pr-2 pl-2.5 gap-1.5 rounded-lg`, relevant only if the custom listbox gets built). Shadow presence: `shadow-xs` removed.

## separator

- current: `data-[orientation=horizontal]:h-px w-full`, `data-[orientation=vertical]:h-full w-px`
- nova: `.cn-separator-horizontal` (L1072-1074): `h-px w-full`; `.cn-separator-vertical` (L1076-1078): `h-full w-px`
- delta: no delta.
- notes: nova's base `separator.tsx` uses `data-vertical:self-stretch` instead of `h-full` — behavioral nicety in flex rows, optional.

## sheet

- current:
  - content: `gap-4` + side blocks `w-3/4 sm:max-w-sm` (left/right), `h-auto` (top/bottom); `shadow-lg`
  - close: `top-4 right-4` (icon size-4)
  - header: `gap-1.5 p-4`
  - footer: `gap-2 p-4`
  - title: `font-semibold` (no size)
  - description: `text-sm`
- nova:
  - `.cn-sheet-content` (L1085-1087): `gap-4 text-sm shadow-lg` + same side geometry (`w-3/4`, `sm:max-w-sm`, `h-auto`)
  - `.cn-sheet-close` (L1089-1091): `absolute top-3 right-3`
  - `.cn-sheet-header` (L1093-1095): `gap-0.5 p-4`
  - `.cn-sheet-footer` (L1097-1099): `gap-2 p-4`
  - `.cn-sheet-title` (L1101-1103): `text-base font-medium`
  - `.cn-sheet-description` (L1105-1107): `text-sm`
- delta:
  - content: add `text-sm`
  - close: `top-4 right-4 → top-3 right-3`
  - header: `gap-1.5 → gap-0.5`
  - title: add `text-base` (weight `font-semibold → font-medium`)
- notes: shadow-lg kept — no presence change. Side geometry (insets, w-3/4, sm:max-w-sm) already matches.

## skeleton

- current: `rounded-md`
- nova: `.cn-skeleton` (L1223-1225): `rounded-md`
- delta: no delta. (bg-accent vs bg-muted is color — out of scope.)

## spinner

- no nova counterpart — keep as is. No `.cn-spinner` in style-nova.css; base `spinner.tsx` carries `size-4 animate-spin`, which gsxui already matches exactly.

## switch

- current: `h-[1.15rem] w-8` (= 18.4px x 32px); thumb `size-4`, `translate-x-[calc(100%-2px)]`
- nova:
  - `.cn-switch` (L1250-1252): `data-[size=default]:h-[18.4px] data-[size=default]:w-[32px] data-[size=sm]:h-[14px] data-[size=sm]:w-[24px]`
  - `.cn-switch-thumb` (L1258-1260): default `size-4` / sm `size-3`, checked `translate-x-[calc(100%-2px)]`
- delta: no delta for the default size (`1.15rem` == 18.4px, `w-8` == 32px; thumb identical).
  - NEW: `data-[size=sm]` variant (`h-[14px] w-[24px]`, thumb `size-3`)
- notes: shadow presence — gsxui `shadow-xs`; nova has none.

## table

- current: head `h-10 px-2`; cell `p-2`; caption `mt-4 text-sm`; table `text-sm`
- nova: `.cn-table-head` (L1295-1297): `h-10 px-2`; `.cn-table-cell` (L1303-1305): `p-2`; `.cn-table-caption` (L1311-1313): `mt-4 text-sm`; `.cn-table` (L1271-1273): `text-sm`
- delta: no delta.

## tabs

- current:
  - root: `gap-2`
  - list: `h-9 rounded-lg p-[3px]`
  - trigger: `h-[calc(100%-1px)] gap-1.5 rounded-md border border-transparent px-2 py-1 text-sm [&_svg:not([class*='size-'])]:size-4`
  - content: no metric
- nova:
  - `.cn-tabs` (L1316-1318): `gap-2`
  - `.cn-tabs-list` (L1320-1322): `rounded-lg p-[3px] group-data-horizontal/tabs:h-8`
  - `.cn-tabs-trigger` (L1324-1326): `gap-1.5 rounded-md border border-transparent px-1.5 py-0.5 text-sm has-data-[icon=inline-end]:pr-1 has-data-[icon=inline-start]:pl-1 [&_svg:not([class*='size-'])]:size-4` (base tsx keeps `h-[calc(100%-1px)]`)
  - `.cn-tabs-content` (L1332-1334): `text-sm`
- delta:
  - list: `h-9 → h-8` (nova scopes it `group-data-horizontal/tabs:`; gsxui dropped orientation, so plain `h-8` is the faithful port)
  - trigger: `px-2 → px-1.5`, `py-1 → py-0.5`
  - trigger: add `has-data-[icon=inline-start]:pl-1`, `has-data-[icon=inline-end]:pr-1` (needs data-icon stamps)
  - content: add `text-sm`
- notes: active-trigger `shadow-sm` present in both (nova scopes it to the default variant) — no presence change. `h-[calc(100%-1px)]` on trigger unchanged (verified in base tabs.tsx).

## textarea

- current: `min-h-16 rounded-md px-3 py-2 text-base md:text-sm` (+ `border shadow-xs`)
- nova: `.cn-textarea` (L1337-1339): `rounded-lg px-2.5 py-2 text-base md:text-sm` (base tsx keeps `min-h-16`)
- delta:
  - `rounded-md → rounded-lg`
  - `px-3 → px-2.5`
- notes: `min-h-16` unchanged (lives in base textarea.tsx). Shadow presence: `shadow-xs` removed.

## toggle

- current:
  - base: `gap-2 rounded-md text-sm [&_svg:not([class*='size-'])]:size-4`
  - default: `h-9 min-w-9 px-2`; sm: `h-8 min-w-8 px-1.5`; lg: `h-10 min-w-10 px-2.5`
- nova:
  - `.cn-toggle` (L1342-1344): `gap-1 rounded-lg text-sm [&_svg:not([class*='size-'])]:size-4`
  - `.cn-toggle-size-default` (L1358-1360): `h-8 min-w-8 px-2.5 has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2`
  - `.cn-toggle-size-sm` (L1362-1364): `h-7 min-w-7 rounded-[min(var(--radius-md),12px)] px-2.5 text-[0.8rem] has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 [&_svg:not([class*='size-'])]:size-3.5`
  - `.cn-toggle-size-lg` (L1366-1368): `h-9 min-w-9 px-2.5 has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2`
- delta:
  - base: `gap-2 → gap-1`, `rounded-md → rounded-lg`
  - default: `h-9 → h-8`, `min-w-9 → min-w-8`, `px-2 → px-2.5`
  - sm: `h-8 → h-7`, `min-w-8 → min-w-7`, `px-1.5 → px-2.5`, add `rounded-[min(var(--radius-md),12px)]`, `text-[0.8rem]`, svg `size-3.5`
  - lg: `h-10 → h-9`, `min-w-10 → min-w-9` (px-2.5 unchanged)
  - all sizes: add directional `has-data-[icon=...]` paddings (needs data-icon stamps)
- notes: shadow presence — gsxui outline variant `shadow-xs`; nova outline has none.

## tooltip

- current: content `w-fit rounded-md px-3 py-1.5 text-xs`; arrow `size-2.5 rounded-[2px]` offset `calc(50%+2px)`
- nova:
  - `.cn-tooltip-content` (L1380-1382): `gap-1.5 rounded-md px-3 py-1.5 text-xs has-data-[slot=kbd]:pr-1.5 **:data-[slot=kbd]:rounded-sm`
  - `.cn-tooltip-arrow` (L1388-1390): `size-2.5 translate-y-[calc(-50%-2px)] rounded-[2px]`
- delta:
  - content: add `gap-1.5` (inline-flex row) and `has-data-[slot=kbd]:pr-1.5`, `**:data-[slot=kbd]:rounded-sm`
- notes: core paddings/radius/arrow geometry already match (gsxui's `-translate-y-[calc(50%+2px)]` equals nova's `translate-y-[calc(-50%-2px)]`). No shadow changes.

## icon (ui/icon/icon.gsx)

- no nova counterpart — keep as is. No `.cn-icon` entry; nova sizes icons per consuming part (`[&_svg:not([class*='size-'])]:size-4` almost everywhere, size-3/3.5 in xs/sm contexts, size-6 in alert-dialog media). gsxui's `size-4` default remains the right base.

---

## Shadow-presence footnote (roll-up)

| component | gsxui | nova |
|---|---|---|
| button (outline) | shadow-xs | none |
| input / textarea / select / input-group root | shadow-xs | none |
| checkbox / radio / switch | shadow-xs | none |
| button-group text | shadow-xs | none |
| card | shadow-sm + border | ring-1, no shadow |
| dialog / alert-dialog content | shadow-lg + border | ring-1, no shadow |
| dropdown / context-menu / popover / hover-card content | shadow-md + border | shadow-md + ring-1 (border removed) |
| sheet content | shadow-lg | shadow-lg (unchanged) |
| tabs active trigger | shadow-sm | shadow-sm (unchanged, variant-scoped) |

## Markup prerequisites the class swaps depend on

1. `data-icon="inline-start" | "inline-end"` on icon children of button / badge / toggle / tabs-trigger — nova's directional icon paddings key off it.
2. `data-slot="accordion-trigger-icon"` on the accordion chevron (nova sizes and `ml-auto`-positions it by slot).
3. New parts to add for full nova parity: AlertAction, AlertDialogMedia, PopoverHeader/Title/Description, ProgressLabel/Value, AvatarBadge/AvatarGroupCount, Item size-xs, ItemGroup gaps.
4. New size axes: alert-dialog `size` (default/sm), avatar `size` (sm/lg), card `size` (sm via `--card-spacing`), native-select `size` (sm), switch `size` (sm).

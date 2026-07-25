# Calendar source map

Source analysis for Task 0 of the 2026-07-25 calendar plan
(`.superpowers/sdd/2026-07-25-calendar/task-0-brief.md`). Tasks 1–7 cite this
document by section heading (e.g. "the `day` slot's classes, source map §2")
for every class string and behavioral claim — do not re-derive either from
the `.tsx`/library files directly; cite here.

**Reference-source rule** (adjudicated, not re-litigated): markup reference
is `apps/v4/registry/new-york-v4/ui/calendar.tsx` — the only remaining
inline-Tailwind-utility form, gsxui's own authoring model. `registry/bases/
{base,radix,aria}/ui/calendar.tsx` is NOT ported from directly, but is used
once below (§6) as a structure tiebreak to determine which literal DOM node
each `.cn-calendar-*` nova class name attaches to, since `style-nova.css`'s
selector names alone don't say that. Visual reference is `registry/styles/
style-nova.css`; where it and new-york-v4 disagree, nova wins.

**Element-structure decision** (settled before this pass, not re-litigated —
reproduced from the calling context, all citations independently confirmed
below): react-day-picker 9.14.0 renders the grid as a real HTML `<table>`
(`MonthGrid.js` → `<table>`, `Weekdays.js` → `<thead aria-hidden={true}><tr>`,
`Weekday.js` → `<th>`, `Weeks.js` → `<tbody>`, `Week.js` → `<tr>`, `Day.js` →
`<td>`, `DayButton.js` → `<button>`), not a div/ARIA-role grid. gsxui adopts
the table structure. Tasks 1–3 build markup against this.

**Inputs read** (byte-read in full unless noted):
- `apps/v4/registry/new-york-v4/ui/calendar.tsx` (221 lines, full file,
  local path `~/personal/shadcn-ui/apps/v4/registry/new-york-v4/ui/
  calendar.tsx`).
- `apps/v4/registry/styles/style-nova.css`, grepped case-insensitively for
  `calendar`/`rdp-` (zero hits outside the three `.cn-calendar*` rules named
  below — confirmed exhaustive, not partial).
- `apps/v4/registry/bases/base/ui/calendar.tsx` (246 lines, full file) —
  structure tiebreak only, per the reference-source rule above.
- react-day-picker 9.14.0's own shipped ESM build, extracted locally at
  `/Users/jackieli/.claude/jobs/6a93038d/tmp/package/dist/esm/` (this is a
  real npm-published build, not a checkout missing `node_modules` — every
  claim below sourced from it is a genuine source read, not
  `derived-not-read`, unless explicitly marked otherwise):
  `DayPicker.js` (full, 329 lines — the orchestrating component, the single
  most load-bearing file in this map), `UI.js` (full — the `UI`/`DayFlag`/
  `SelectionState`/`Animation` enums, i.e. the authoritative list of every
  classNames-slot key and modifier key), `useFocus.js`, `useSelection.js`
  (full), `helpers/{calculateFocusTarget,getFocusableDate,getNextFocus,
  createGetModifiers,getClassNamesForModifiers,getDefaultClassNames,
  getMonths,getYearOptions}.js` (full), `selection/{useSingle,useRange,
  useMulti}.js` (full), `utils/dateMatchModifiers.js` (full),
  `labels/{labelGridcell,labelDayButton}.js` (full), and every file under
  `components/` referenced in the element-structure decision above plus
  `Root.js`, `Months.js`, `Month.js`, `MonthCaption.js`, `Button.js`,
  `Nav.js`, `PreviousMonthButton.js`, `DropdownNav.js`, `Dropdown.js`,
  `Select.js`, `WeekNumberHeader.js`, `WeekNumber.js`, `Chevron.js`, `CaptionLabel.js`,
  `Footer.js` (each ≤30 lines, full).
- `useCalendar.js` (full, 110 lines), `helpers/getNavMonth.js`,
  `helpers/getPreviousMonth.js`, `helpers/getNextMonth.js` (all full) — added
  in the Task 0 review-fix pass, for the nav-button disabled contract in
  §7.5 and the `fromYear`/`toYear` bounds question it raises.
- Not read this pass, no attempt made: `useCalendar.js`'s `getInitialMonth`/
  `getDisplayMonths`/`getDates` internals, and `utils/addToRange.js` —
  nothing in this document depends on their internals beyond what's already
  stated inline wherever they're mentioned.

## Legend

- `derived-not-read` — reconstructed from props/class strings/general
  public-API knowledge, not from reading the library's own source.
- No such tag appears anywhere in this document outside §7.4's one
  cross-reference note — every class string and every behavioral claim
  below was read directly from the files listed above.

---

## 1. `--cell-size` custom property

Set once, on the `Root` element's own class string, in new-york-v4:
```
[--cell-size:--spacing(8)]
```
(`calendar.tsx` line 36, inside the `cn(...)` call building the `DayPicker`'s
own top-level `className` — quoted in full, byte-verbatim, as the second
`root` row in §2's table, added in the Task 1 review-fix pass; this paragraph
only ever carried the one `--cell-size` token out of that string, not the
whole thing). Nova overrides the value (not the mechanism) via
`.cn-calendar` (§6): `[--cell-size:--spacing(7)]`, and additionally
introduces a second custom property new-york-v4 does not have at all:
`--cell-radius: var(--radius-md)`, used by `bases/base/ui/calendar.tsx`
wherever new-york-v4 hardcodes `rounded-md`/`rounded-l-md`/`rounded-r-md`
literally (structure tiebreak, confirmed by direct comparison of the two
`.tsx` files — every `rounded-*-md` token in new-york-v4's `day`/
`range_start`/`range_middle`/`range_end`/`today` values has a `rounded-(--
cell-radius)` counterpart in `bases/base`'s version of the same slot). Since
new-york-v4 is the markup reference and it has no `--cell-radius` variable at
all, **recommend NOT introducing `--cell-radius`** — port `[--cell-size:
--spacing(7)]` (nova's value) as the one custom property, keep the literal
`rounded-md` tokens new-york-v4 already uses. `--cell-radius` is a
`bases/base`-family-only mechanism, out of scope per the reference-source
rule; do not silently absorb it while borrowing the tiebreak for §6.

## 2. Every className slot, exact string (`calendar.tsx`, byte-read in full)

The object literal in `calendar.tsx` itself has exactly 25 keys — every one
cross-checked against the authoritative `UI` enum in `UI.js` below — all
quoted below verbatim from the `cn(...)` call building each one (i.e. the value *before* `defaultClassNames[key]` and any caller
`classNames` override are appended — those two are handled generically by
every slot, not repeated per-slot below). **All 25 have real class strings,
read directly from source. Zero are `derived-not-read`.**

| Slot | Class string |
|---|---|
| `root` | `w-fit` |
| `root` (DayPicker top-level `className`, distinct — see note below) | `group/calendar bg-background p-3 [--cell-size:--spacing(8)] [[data-slot=card-content]_&]:bg-transparent [[data-slot=popover-content]_&]:bg-transparent rtl:**:[.rdp-button\_next>svg]:rotate-180 rtl:**:[.rdp-button\_previous>svg]:rotate-180` |
| `months` | `relative flex flex-col gap-4 md:flex-row` |
| `month` | `flex w-full flex-col gap-4` |
| `nav` | `absolute inset-x-0 top-0 flex w-full items-center justify-between gap-1` |
| `button_previous` | `buttonVariants({variant: buttonVariant})` + `size-(--cell-size) p-0 select-none aria-disabled:opacity-50` — disabled/focus **contract in §7.5**, not carried by the class string alone |
| `button_next` | `buttonVariants({variant: buttonVariant})` + `size-(--cell-size) p-0 select-none aria-disabled:opacity-50` — disabled/focus **contract in §7.5**, not carried by the class string alone |
| `month_caption` | `flex h-(--cell-size) w-full items-center justify-center px-(--cell-size)` |
| `dropdowns` | `flex h-(--cell-size) w-full items-center justify-center gap-1.5 text-sm font-medium` |
| `dropdown_root` | `relative rounded-md border border-input shadow-xs has-focus:border-ring has-focus:ring-[3px] has-focus:ring-ring/50` |
| `dropdown` | `absolute inset-0 bg-popover opacity-0` |
| `caption_label` | `font-medium select-none` **+** (ternary on `captionLayout==="label"`): `"text-sm"` **else** `"flex h-8 items-center gap-1 rounded-md pr-1 pl-2 text-sm [&>svg]:size-3.5 [&>svg]:text-muted-foreground"` — two genuinely different values depending on caption layout, not one string |
| `month_grid` | `w-full border-collapse` |
| `weekdays` | `flex` |
| `weekday` | `flex-1 rounded-md text-[0.8rem] font-normal text-muted-foreground select-none` |
| `week` | `mt-2 flex w-full` |
| `week_number_header` | `w-(--cell-size) select-none` |
| `week_number` | `text-[0.8rem] text-muted-foreground select-none` |
| `day` | `group/day relative aspect-square h-full w-full p-0 text-center select-none [&:last-child[data-selected=true]_button]:rounded-r-md` **+** (ternary on `props.showWeekNumber`): `"[&:nth-child(2)[data-selected=true]_button]:rounded-l-md"` **else** `"[&:first-child[data-selected=true]_button]:rounded-l-md"` |
| `range_start` | `rounded-l-md bg-accent` |
| `range_middle` | `rounded-none` |
| `range_end` | `rounded-r-md bg-accent` |
| `today` | `rounded-md bg-accent text-accent-foreground data-[selected=true]:rounded-none` |
| `outside` | `text-muted-foreground aria-selected:text-muted-foreground` |
| `disabled` | `text-muted-foreground opacity-50` |
| `hidden` | `invisible` |

**The two `root` rows above are two different strings on two different props,
byte-verified against `calendar.tsx` lines 33–40 (added in the Task 1
review-fix pass — this distinction existed in §1's prose but was never
quoted in full here, which is what this table now fixes).** `classNames.root`
(first row, `w-fit`) is the `DayPicker`'s `classNames` prop, consumed the same
way as every other slot in this table. Separately, `calendar.tsx`'s own
`<DayPicker className={cn(...)}>` call (line 35, the component's *own*
top-level `className` prop, a sibling of `classNames`, not a member of it)
passes:
```
cn(
  "group/calendar bg-background p-3 [--cell-size:--spacing(8)] [[data-slot=card-content]_&]:bg-transparent [[data-slot=popover-content]_&]:bg-transparent",
  String.raw`rtl:**:[.rdp-button\_next>svg]:rotate-180`,
  String.raw`rtl:**:[.rdp-button\_previous>svg]:rotate-180`,
  className
)
```
(`calendar.tsx` lines 35–40, byte-read). react-day-picker merges this
top-level `className` onto the same physical root `<div>` that
`classNames.root` also targets — there is exactly one root element in the
rendered DOM, so gsxui's single root `<div>` carries both strings' tokens
concatenated. §1's `[--cell-size:--spacing(8)]` citation is the one token
from this second string that already had its own paragraph; the rest
(`group/calendar`, `bg-background`, `p-3`, both `[[data-slot=…]_&]:` scoping
selectors, and the two `rtl:**:[...]:rotate-180` selectors) had never been
quoted anywhere in this document before this pass — this table row is now
the citation for all of it.

**gsxui does not port the two `rtl:**:[...]:rotate-180` selectors.** Both
target `.rdp-button_next`/`.rdp-button_previous` — react-day-picker's own
bare `rdp-*` hook classes (`getDefaultClassNames.js`, cited in the Inputs-read
list above) — and this document's own §2 intro already flags that gsxui does
not port any bare `rdp-*` token (open question there, resolved here): since
`.rdp-button_next`/`.rdp-button_previous` never appear in gsxui's markup, an
`rtl:**:[.rdp-button_next>svg]:rotate-180`-style selector can never match
anything in this port and is dead on arrival. Drop both, keep everything
else in the top-level-`className` row above.

Every slot's *final* class value at runtime is
`cn(<string above>, defaultClassNames[key], ...(caller classNames[key]))` —
`defaultClassNames[key]` is react-day-picker's own `rdp-<slot>` string
(confirmed from `getDefaultClassNames.js`, read in full: it iterates the
`UI`/`DayFlag`/`SelectionState`/`Animation` enums in `UI.js` and stamps
`` `rdp-${value}` `` for each). Whether gsxui ports the `rdp-*` bare class
tokens at all is a Task 1–3 decision, not resolved here — they carry no
styling of their own in this registry (no `.rdp-*` rule exists anywhere in
`style-nova.css`, confirmed by the same calendar/rdp- grep in Inputs read
above), they exist purely as stable hook classNames for library consumers.

**Full key list, from `UI.js`'s own enum** (authoritative — this *is* the
25-item list above, by value, confirmed 1:1): `root`, `chevron`, `day`,
`day_button`, `caption_label`, `dropdowns`, `dropdown`, `dropdown_root`,
`footer`, `month_grid`, `month_caption`, `months_dropdown`, `month`,
`months`, `nav`, `button_next`, `button_previous`, `week`, `weeks`,
`weekday`, `weekdays`, `week_number`, `week_number_header`,
`years_dropdown`. Plus `DayFlag`: `disabled`, `hidden`, `outside`,
`focused`, `today`. Plus `SelectionState`: `range_end`, `range_middle`,
`range_start`, `selected`. new-york-v4's `classNames={{...}}` object
customizes 25 of these; the rest (`chevron`, `day_button`, `footer`,
`weeks`, `months_dropdown`, `years_dropdown`, `focused`, `selected`) are left
at their library default (`rdp-<name>`) or, for `day_button`, styled entirely
through `CalendarDayButton`'s own `className` prop (§3) rather than through
the `classNames.day_button` slot at all — new-york-v4 never sets
`classNames.day_button`, confirmed by its absence from the object above.

## 3. `CalendarDayButton` — classes and data attributes (`calendar.tsx` lines 182–218)

Composes gsxui's `Button`-equivalent (`variant="ghost" size="icon"`), full
class string:
```
flex aspect-square size-auto w-full min-w-(--cell-size) flex-col gap-1 leading-none font-normal group-data-[focused=true]/day:relative group-data-[focused=true]/day:z-10 group-data-[focused=true]/day:border-ring group-data-[focused=true]/day:ring-[3px] group-data-[focused=true]/day:ring-ring/50 data-[range-end=true]:rounded-md data-[range-end=true]:rounded-r-md data-[range-end=true]:bg-primary data-[range-end=true]:text-primary-foreground data-[range-middle=true]:rounded-none data-[range-middle=true]:bg-accent data-[range-middle=true]:text-accent-foreground data-[range-start=true]:rounded-md data-[range-start=true]:rounded-l-md data-[range-start=true]:bg-primary data-[range-start=true]:text-primary-foreground data-[selected-single=true]:bg-primary data-[selected-single=true]:text-primary-foreground dark:hover:text-accent-foreground [&>span]:text-xs [&>span]:opacity-70
```
Note the `group-data-[focused=true]/day:*` selectors read a **named group**
established by the `day` slot's own `group/day` class (§2) — the button's
focus-ring styling is driven by the *cell's* `data-focused` attribute (§4),
not by a `data-focused` attribute on the button itself; the button has no
such attribute of its own.

Five data attributes, all computed by `CalendarDayButton` itself from the
`modifiers` prop it receives — **not** part of react-day-picker's own
`DayButton` contract (`DayButton.js`, read in full, sets none of these; it
only forwards whatever `buttonProps` `DayPicker.js` passes down, see §4):
```
data-day={day.date.toLocaleDateString()}
data-selected-single={modifiers.selected && !modifiers.range_start && !modifiers.range_end && !modifiers.range_middle}
data-range-start={modifiers.range_start}
data-range-end={modifiers.range_end}
data-range-middle={modifiers.range_middle}
```
`data-day` here is a **locale-formatted** string (`toLocaleDateString()`),
distinct from and unrelated to the `<td>`'s own `data-day` (§4), which is the
ISO date — two different attributes, same name, different element, both
real. `CalendarDayButton` also wires a `ref` + a `useEffect` that calls
`ref.current?.focus()` whenever `modifiers.focused` becomes true — this is
new-york-v4's own imperative-focus wiring, redundant with (but not the same
mechanism as) react-day-picker's own `DayButton.js` default implementation,
which does the *identical* `useEffect`/`ref.current?.focus()` pattern
already (confirmed, `DayButton.js` lines 10–14, byte-identical logic) —
`CalendarDayButton` re-implements what the default component already does,
because it replaces the default component wholesale rather than wrapping it.
gsxui's port needs this same imperative-focus-on-`data-focused`-change
behavior somewhere; it is not optional decoration, it is how roving-tabindex
focus actually lands in the DOM (see §7.3).

## 4. The `<td>` (`Day`) — data attributes and ARIA, from `DayPicker.js` itself

Read directly off the `React.createElement(components.Day, {...})` call
(`DayPicker.js`, the single call site building every day cell):
```
role="gridcell"
aria-selected={modifiers.selected || undefined}
aria-label={ariaLabel}          // only set when !isInteractive && !modifiers.hidden, via labelGridcell()
data-day={day.isoDate}           // ISO date string
data-month={day.outside ? day.dateMonthId : undefined}
data-selected={modifiers.selected || undefined}
data-disabled={modifiers.disabled || undefined}
data-hidden={modifiers.hidden || undefined}
data-outside={day.outside || undefined}
data-focused={modifiers.focused || undefined}
data-today={modifiers.today || undefined}
```
Confirms plan-correction #3 (§8) precisely: this whole attribute set is on
the **cell**, not the button. `Day.js` itself (the default component) is a
bare `<td {...tdProps}/>`, so nothing here is added by the default
component's own code — it all comes from the props `DayPicker.js` passes in.

The `<button>` inside the cell, by contrast (same call site, the
`components.DayButton` invocation, only rendered `if (!modifiers.hidden &&
isInteractive)`), receives from the core library:
```
type="button"
disabled={(!modifiers.focused && modifiers.disabled) || undefined}
aria-disabled={(modifiers.focused && modifiers.disabled) || undefined}
tabIndex={isFocusTarget(day) ? 0 : -1}
aria-label={labelDayButton(date, modifiers, ...)}
onClick / onBlur / onFocus / onKeyDown / onMouseEnter / onMouseLeave
```
— i.e. the core library's own contribution to the button is disabled/
tabIndex/aria-label/handlers only. `CalendarDayButton`'s five `data-*`
attributes (§3) are additional, shadcn-authored, layered on top of this by
the wrapper component — the two attribute sets (cell's 6 `data-*` + button's
5 `data-*`) are **not duplicates of each other** and both need porting; see
§8 finding 5 for why the `disabled`/`aria-disabled` split matters.

## 5. Structural components, element-by-element (all read in full, all confirming §"Element-structure decision")

| UI piece | Component file | Element |
|---|---|---|
| Root | `Root.js` | `<div ref={rootRef}>` |
| Months container | `Months.js` | `<div>` |
| Month wrapper | `Month.js` | `<div>` |
| Nav | `Nav.js` | `<nav>`, wrapping two `components.Button`-rendered `<button>`s |
| Nav button (generic) | `Button.js` | `<button>` (deprecated in favor of `PreviousMonthButton`/`NextMonthButton`, both of which just re-delegate to `components.Button`, confirmed via `PreviousMonthButton.js`) |
| Month caption | `MonthCaption.js` | `<div>` |
| Caption label (non-dropdown) | `CaptionLabel.js` | `<span>` |
| Dropdown nav container | `DropdownNav.js` | `<div>` |
| One dropdown (month or year) | `Dropdown.js` | `<span data-disabled className={DropdownRoot}>` wrapping a `<select>` (`Select.js`) plus a second `<span aria-hidden className={CaptionLabel}>` holding the selected option's label + a down `Chevron` |
| Month grid | `MonthGrid.js` | `<table>` |
| Weekday header row | `Weekdays.js` | `<thead aria-hidden={true}><tr>` |
| One weekday cell | `Weekday.js` | `<th>` |
| Weeks body | `Weeks.js` | `<tbody>` |
| One week row | `Week.js` | `<tr>` |
| Week-number header cell | `WeekNumberHeader.js` | `<th>` |
| Week-number cell (library default) | `WeekNumber.js` | `<th>` — **note**: new-york-v4's own `components.WeekNumber` override (`calendar.tsx` lines 166–174) replaces this with a `<td>` wrapping a centering `<div>`, i.e. shadcn itself already drops the row-header semantics react-day-picker's own default provides. Moot for this port since `week_number` is not being ported (§9), but recorded so nobody mistakes the library default for what shadcn ships. |
| Day cell | `Day.js` | `<td>` |
| Day button | `DayButton.js` | `<button ref={ref}>`, with its own `useEffect(() => { if (modifiers.focused) ref.current?.focus() }, [modifiers.focused])` |
| Footer | `Footer.js` | `<div>` |
| Chevron icon | `Chevron.js` | inline `<svg><polygon/></svg>` — moot, new-york-v4 replaces this wholesale with `lucide-react` icons (`calendar.tsx` lines 145–164) |

## 6. Nova overrides — exhaustive, only three rules exist

Grepped `style-nova.css` case-insensitively for `calendar` and `rdp-`: **zero
hits outside these three rules.** Every other slot in §2 is `(no nova
counterpart)` — new-york-v4's own value stands unreviewed by the nova
retarget; port it as-is.

```css
.cn-calendar {
  @apply p-2 [--cell-radius:var(--radius-md)] [--cell-size:--spacing(7)];
}
.cn-calendar-dropdown-root {
  @apply has-focus:border-ring border-input has-focus:ring-ring/50 border has-focus:ring-3;
}
.cn-calendar-caption-label {
  @apply h-6 pr-1 pl-1.5;
}
```

These class *names* don't appear anywhere in new-york-v4's markup (they're
`bases/base`-family-only literal classes, e.g. `"cn-calendar-dropdown-root
relative rounded-(--cell-radius)"` on `dropdown_root`) — which node each one
targets is not answerable from `style-nova.css` alone, hence the one
structure-tiebreak read of `bases/base/ui/calendar.tsx` (§ Inputs read).
Confirmed target + how each nova rule maps onto new-york-v4's own value:

- **`.cn-calendar` → `root`.** Additive/override on top of new-york-v4's own
  `root` value (`w-fit`) plus the `[--cell-size:...]` set on the outer
  `DayPicker` `className`, not `root` itself (§1) — nova changes `p-3→p-2`
  (this is set on the `DayPicker`'s own top-level `className` in
  new-york-v4, not on the `root` slot; nova's `.cn-calendar` folds it into
  one rule since `bases/base` puts `cn-calendar` directly on that same
  top-level className) and `--cell-size:--spacing(8)→--spacing(7)`.
  `--cell-radius` is `bases/base`-only; **do not port it** (§1).
- **`.cn-calendar-dropdown-root` → `dropdown_root`.** new-york-v4's own value
  is `relative rounded-md border border-input shadow-xs has-focus:border-ring
  has-focus:ring-[3px] has-focus:ring-ring/50`. Nova's rule
  (`has-focus:border-ring border-input has-focus:ring-ring/50 border
  has-focus:ring-3`) is the same tokens, `ring-[3px]→ring-3` (spelling only,
  same value — matches the precedent in the menus/combobox source maps for
  this exact token), and drops `shadow-xs` **(shadow-presence, drop — adopt
  the removal)**. `rounded-md` is unchanged (nova doesn't touch it here,
  despite `bases/base` itself using `rounded-(--cell-radius)` — that
  substitution is `bases/base`'s own authoring choice, not something
  `style-nova.css`'s rule for this class actually states, since the rule has
  no `rounded-*` token in it at all; keep `rounded-md` literal, matching
  §1's decision not to adopt `--cell-radius`).
- **`.cn-calendar-caption-label` → the *dropdown-layout* branch only of
  `caption_label`.** `bases/base`'s own markup (§ Inputs read, line 87)
  attaches `cn-calendar-caption-label` only to the
  `captionLayout !== "label"` branch (`"flex items-center gap-1
  rounded-(--cell-radius) text-sm [&>svg]:size-3.5 [&>svg]:text-muted-
  foreground"`), never to the plain-label branch (which gets a different,
  nova-rule-less class, `cn-calendar-caption`, confirmed absent from the
  three-rule grep above). new-york-v4's dropdown-branch value is `flex h-8
  items-center gap-1 rounded-md pr-1 pl-2 text-sm [&>svg]:size-3.5
  [&>svg]:text-muted-foreground`; nova's rule (`h-6 pr-1 pl-1.5`) overrides
  `h-8→h-6` and `pl-2→pl-1.5` (`pr-1` unchanged) — consistent scaling with
  the `--cell-size:8→7` reduction in `.cn-calendar`. **The plain-`"label"`
  branch (`"text-sm"`) is untouched by nova** — port it exactly as
  new-york-v4 has it.

**Flag for Task 1–3, not resolved here** — narrowed to what is actually true
after re-checking both files token-by-token (an earlier draft of this
section overstated this as a blanket "throughout" swap; corrected):
- At the `Calendar` component's own `classNames`-slot level (§2):
  `bases/base`'s `range_start` and `range_end` use `bg-muted` where
  new-york-v4 uses `bg-accent`; `bases/base`'s `today` uses `bg-muted
  text-foreground` where new-york-v4 uses `bg-accent text-accent-foreground`.
  `range_middle` is **not** part of this divergence — both files leave it
  colorless (`"rounded-none"`, no `bg-*`/`text-*` token, in both).
- At `CalendarDayButton`'s own data-attribute level (§3): **only**
  `data-range-middle` differs (`bg-muted text-foreground` in `bases/base` vs.
  `bg-accent text-accent-foreground` in new-york-v4), plus the trailing
  `dark:hover:text-foreground` (`bases/base`) vs. `dark:hover:text-accent-
  foreground` (new-york-v4). `data-range-start`, `data-range-end`, and
  `data-selected-single` are **byte-identical** in both files (`bg-primary
  text-primary-foreground` in every case) — no divergence there at all.

**This is `bases/base`'s own family-wide design language, not a
`style-nova.css` directive** — none of the three nova rules above touch
color, and the muted/accent swap is baked directly into `bases/base`'s
`.tsx` as plain Tailwind utilities, sitting alongside (not sourced from) its
`cn-calendar-*` class names. Per the reference-source rule, `bases/base` is
consulted here **only** as a structure tiebreak for where the three real
nova rules attach — its own color choices carry no authority and must not
leak into the port. Keep new-york-v4's `bg-accent`/`bg-primary` values
throughout; the port recommendation is unchanged by this correction, only
the size of the claim is.

## 7. Behavior traces

### 7.1 Selection semantics, all three modes (`selection/use{Single,Range,Multi}.js`, `useSelection.js`, full files)

`useSelection.js` is a pure dispatcher: `single`/`multiple`/`range` mode
picks the matching hook; no mode returns `undefined`.

- **`single`** (`useSingle.js`): one selected `Date | undefined`. Clicking
  the already-selected day clears the selection (`newDate = undefined`)
  **unless `required` is set**, in which case re-clicking the selected day
  is a no-op that still calls `onSelect` with the same value (the guard is
  `!required && selected && isSameDay(triggerDate, selected)`). Clicking any
  other day replaces the selection outright.
- **`range`** (`useRange.js`): selected value is `{from, to}`. Core logic
  delegates date-arithmetic to `addToRange` (not read this pass — not
  load-bearing for the class-string/behavior-trace deliverable, flagged as
  a gap: Task 1–3 implementing range-click sequencing beyond
  "first click sets `from`, later clicks extend/reset" should read
  `utils/addToRange.js` directly, not rely on this map). Two props alter the
  base behavior: **`resetOnSelect`** — when true, clicking while a full
  range is already selected (or nothing is selected yet) starts a **fresh**
  range at the clicked day (`{from: triggerDate, to: undefined}`) instead of
  extending/replacing via `addToRange`; clicking the single day that is
  currently both `from` and `to` of a one-day range clears the selection
  entirely unless `required` is set. **`excludeDisabled`** — after
  computing a candidate `{from, to}`, if `disabled` matchers are set and any
  day inside that range matches them (`rangeContainsModifiers`), the range
  is **reset to just the clicked day** (`from: triggerDate, to: undefined`)
  rather than allowing a range that spans over a disabled day. `min`/`max`
  props are passed straight into `addToRange`, not re-derived here.
- **`multiple`** (`useMulti.js`): selected value is `Date[]`. Clicking a
  day already in the selection removes it — **unless** removing it would
  drop the array below `min` (silent no-op) or below 1 while `required` is
  set and exactly 1 day is currently selected (silent no-op, distinct from
  the `min` check: this fires even if `min` is unset/0). Clicking a day not
  in the selection adds it — **unless** the array is already at `max`
  length, in which case the **entire selection is replaced** by the single
  newly-clicked day (`newDates = [triggerDate]`), not rejected as a no-op —
  a real, easy-to-miss-in-tests behavior: hitting the multi-select ceiling
  doesn't block further clicks, it resets to a fresh one-item selection.

### 7.2 Disabled-day model

`props.disabled` is a `Matcher | Matcher[]`, tested per-day by
`dateMatchModifiers` (`utils/dateMatchModifiers.js`, full file) inside
`createGetModifiers` (full file) when building each day's modifier flags.
`dateMatchModifiers` supports, in order of the checks in its own source:
boolean literal, a single `Date` (exact-day match via `isSameDay`), an array
of `Date`s (any exact-day match), a `DateRange` (`{from,to}`, via
`rangeIncludesDate`), a day-of-week matcher (`{dayOfWeek: number |
number[]}`), a `DateInterval` (`{before, after}` — note: if
`isAfter(before, after)` it's treated as a **closed** interval, disabling
days strictly between the two; otherwise it's an **open/inverted**
interval, disabling everything strictly outside the `[after, before]` span —
this inversion-by-comparing-the-two-bounds is a real, non-obvious rule, not
a typo), a `DateAfter` (`{after}`, disables strictly-after), a `DateBefore`
(`{before}`, disables strictly-before), or an arbitrary predicate function
`(date) => boolean`. `hidden` uses the exact same `dateMatchModifiers`
machinery, independently, for a semantically different flag (§7.3 — hidden
days are removed from the visual grid, e.g. `showOutsideDays=false` padding
days; disabled days remain visible but unselectable).

Once `modifiers.disabled` is true for a day, its concrete DOM expression
(read directly off `DayPicker.js`'s own `React.createElement`
call for `components.DayButton`, §4) is **not** simply "add the `disabled`
attribute" — it branches on whether the day also happens to be the live
interactively-focused day (`modifiers.focused`, i.e.
`Boolean(focused?.isEqualTo(day))`, a *different* flag from the
roving-tabindex `isFocusTarget(day)` check that drives `tabIndex`):
```
disabled={(!modifiers.focused && modifiers.disabled) || undefined}
aria-disabled={(modifiers.focused && modifiers.disabled) || undefined}
```
A disabled day that is **not** the live-focused day gets the real `disabled`
attribute (unclickable, unfocusable, removed from the tab sequence
entirely). A disabled day that **is** currently the live-focused day gets
`aria-disabled` instead, keeping the native `disabled` attribute off so the
element doesn't abruptly lose DOM focus. `getNextFocus`/`calculateFocusTarget`
(both full files) never deliberately land keyboard focus on a disabled day
in the first place (`isFocusableDay` explicitly excludes
`modifiers[DayFlag.disabled]`; `getNextFocus` recurses past disabled dates,
capped at 365 attempts) — the `aria-disabled` branch exists purely as a
defensive fallback for the case where a day's disabled-ness changes
*while* it already holds live focus (e.g. a controlled `disabled` prop
updates after the user has already focused that day), not as a path the
library's own navigation logic ever intentionally routes through.

### 7.3 Roving-tabindex / keyboard model (`useFocus.js`, `helpers/{calculateFocusTarget,getFocusableDate,getNextFocus}.js`, `DayPicker.js`'s `handleDayKeyDown`, all full files)

Two distinct concepts, easy to conflate:
- **`isFocusTarget(day)`** — computed once per render from
  `calculateFocusTarget`, drives `tabIndex={0}` (every other day gets `-1`).
  Priority order for which single day in the whole grid gets it (first match
  wins, only among days passing `isFocusableDay` — not disabled, not
  hidden, not outside): a day carrying an explicit `modifiers.focused` flag
  from a custom `modifiers` prop > the `lastFocused` day (the day that held
  focus before the most recent blur) > any currently-selected day > `today`
  > else the first focusable day in the whole `days` array.
- **`modifiers.focused`** — whether *this specific day* is the live,
  currently-interactively-focused day (`focusedDay` state), set via
  `onFocus`/cleared via `onBlur`/moved via arrow-key navigation. This is
  what drives the `group-data-[focused=true]/day:*` ring styling (§3) and
  the `disabled`/`aria-disabled` branch (§7.2) — a separate flag from
  `isFocusTarget`.

Keyboard map, read verbatim off `handleDayKeyDown`'s own `keyMap` object
(`DayPicker.js`):

| Key | `moveBy` | `moveDir` |
|---|---|---|
| `ArrowLeft` | `shiftKey ? "month" : "day"` | `dir==="rtl" ? "after" : "before"` |
| `ArrowRight` | `shiftKey ? "month" : "day"` | `dir==="rtl" ? "before" : "after"` |
| `ArrowDown` | `shiftKey ? "year" : "week"` | `"after"` |
| `ArrowUp` | `shiftKey ? "year" : "week"` | `"before"` |
| `PageUp` | `shiftKey ? "year" : "month"` | `"before"` |
| `PageDown` | `shiftKey ? "year" : "month"` | `"after"` |
| `Home` | `"startOfWeek"` | `"before"` |
| `End` | `"endOfWeek"` | `"after"` |

Every other key falls through untouched (no `preventDefault`, no handling —
confirmed by the `if (keyMap[e.key])` guard). `moveFocus(moveBy, moveDir)`
resolves the target date via `getFocusableDate` (adds days/weeks/months/
years, or jumps to `startOfWeek`/`endOfWeek`, clamped to `navStart`/`navEnd`
if the move would cross them) then `getNextFocus` recursively skips forward/
backward past any disabled or hidden date in that same direction (capped at
365 recursive attempts) before landing. `moveFocus` also calls
`calendar.goToDay(nextFocus)` — i.e. **arrow/Home/End/PageUp/PageDown
navigation can change which month is displayed**, not just which day is
focused, whenever the resolved target falls outside the currently-shown
month(s). Landing the actual DOM focus is not this hook's job — it happens
via `DayButton`'s own `useEffect(() => { if (modifiers.focused)
ref.current?.focus() }, [modifiers.focused])` (§3), reacting to the new
`data-focused`/`modifiers.focused` state this hook produces.

### 7.4 Dropdown year-range requirement (context only — the "default a range" deviation is already decided, not re-litigated here)

`getYearOptions` (`helpers/getYearOptions.js`, full file) returns
`undefined` outright if either `navStart` or `navEnd` is falsy — i.e.
upstream's `captionLayout="dropdown"` produces **no year options at all**
unless the caller supplies `startMonth`/`endMonth` (which `useCalendar.js`
turns into `navStart`/`navEnd` via `getNavMonth.js`, both now read in full —
see §7.5 for what else that same function does). This is the upstream
behavior the plan's already-adopted deviation (gsxui defaults a year range
instead of requiring the caller pass one) diverges from — recorded here so
Task 4/5 knows what upstream actually requires and why the default is a real
behavior change, not just a convenience wrapper around an optional prop.

### 7.5 Nav button (`button_previous`/`button_next`) disabled contract — added in the Task 0 review-fix pass

The same architectural pattern that governs the day button's `disabled`/
`aria-disabled` split (§7.2) also governs the nav buttons at the navigation
bounds, and it was previously missing from this document even though §2
lists both buttons' class strings. Read directly off `Nav.js` (full file,
the only call site that builds these two buttons — `DayPicker.js`'s own
`Nav` instantiation for the `navLayout==="around"`/`"after"` cases
duplicates the identical four props verbatim, confirmed at lines 255–256/
282–283/318–322 of `DayPicker.js`):

```
tabIndex: previousMonth ? undefined : -1,
"aria-disabled": previousMonth ? undefined : true,
```
(`button_previous`; the mirror image, keyed on `nextMonth`, applies to
`button_next`). **No native `disabled` prop is set anywhere in `Nav.js`, nor
at any of `DayPicker.js`'s three call sites for these buttons** — confirmed
by reading every prop passed at all four locations; there is no branch that
ever adds `disabled` to either button. This is a **stricter** contract than
the day button's (§7.2's `disabled`/`aria-disabled` split at least sometimes
sets native `disabled`): the nav buttons never get native `disabled` at all,
under any condition — only `aria-disabled` plus `tabIndex={-1}`.

**What gsxui should do** (not merely what upstream does, per the review
finding): render `button_previous`/`button_next` **without** a native
`disabled` attribute at a month boundary. Setting native `disabled` would
remove the button from the tab order and from the accessibility tree's
"disabled but present" state, which is exactly the behavior `aria-disabled`
+ `tabIndex={-1}` is chosen over — a `tabIndex={-1}` element without
`disabled` is skipped by sequential Tab navigation (matching upstream's
intent that it not be a stop in the tab sequence) while still being
present, focusable programmatically, and correctly announced as disabled to
assistive tech via `aria-disabled`, rather than vanishing from the
accessibility tree the way a native `disabled` button would. Task 3 should
implement this as: `aria-disabled={true}` + `tabIndex={-1}` at a bound,
`aria-disabled` omitted + `tabIndex` left at its natural value (`0`, or
unset) otherwise — never a native `disabled` attribute on either nav
button, at any bound.

**`fromYear`/`toYear` (and `startMonth`/`endMonth`) bounds do interact with
this — they're not a separate mechanism from the "data bounds."** Read
`helpers/getNavMonth.js` (full file) and `useCalendar.js` (full file)
together: `getNavMonths(props, dateLib)` resolves `navStart`/`navEnd` from
`startMonth`/`endMonth` (preferred) or the deprecated `fromMonth`/`toMonth`/
`fromYear`/`toYear` props, and — a detail with no equivalent anywhere else
in this document — **if `captionLayout` is `"dropdown"` or
`"dropdown-years"` and no explicit bound was given at all, `getNavMonth.js`
itself defaults `navStart`/`navEnd` to `today ± 100 years`**
(`startOfYear(addYears(today, -100))` through `endOfYear(today)`, lines
33–35/42–44 of `getNavMonth.js`). `useCalendar.js` then feeds that same
`navStart`/`navEnd` pair straight into `getPreviousMonth`/`getNextMonth`
(both full files), which return `undefined` once the currently-displayed
month reaches either bound — and `previousMonth`/`nextMonth` being
`undefined` is precisely what drives `tabIndex`/`aria-disabled` in `Nav.js`
above. So there is exactly **one** disabled mechanism, not two: `navStart`/
`navEnd` already incorporates whatever `startMonth`/`endMonth`/`fromYear`/
`toYear` config the caller supplies (or the dropdown-mode ±100-year default,
§7.4), and the nav buttons go `aria-disabled` whenever the calendar reaches
*that* bound — configured or defaulted — with no distinct "data bounds"
concept beyond it. (`getPreviousMonth`/`getNextMonth` also independently
return `undefined` whenever `props.disableNavigation` is set, and adjust
their step size under `pagedNavigation`, both unrelated to the
`fromYear`/`toYear` question but visible in the same two files.)

---

## 8. Plan corrections

Four supplied by the calling context, verified independently against the
sources above; a fifth found during this pass, not in the original four.

**1. Week rows are not `role="row"`.** Upstream (`Week.js`) renders a bare
`<tr {...trProps}>` — no `role` attribute anywhere in the file or at the
`DayPicker.js` call site building it. `<tr>` inside a `<table>`/`<tbody>`
already carries the implicit ARIA role `row`; writing `role="row"` explicitly
would be a redundant, do-nothing attribute. **Recommend: do not add
`role="row"`** — rely on the implicit table semantics from using a real
`<tr>` inside a real `<table>`/`<tbody>` (which the element-structure
decision already commits to). This isn't a style choice with two valid
answers; redundant explicit ARIA roles are the "bad ARIA" WAI-ARIA
authoring practices explicitly warn against, with zero compensating
benefit here.

**2. Weekday cells are not `role="columnheader"` — and the whole row is
hidden from assistive tech, not just un-roled.** Upstream renders `<th
scope="col">` (`Weekday.js` + the `scope: "col"` prop from `DayPicker.js`),
wrapped in `<thead aria-hidden={true}>` (`Weekdays.js`). Two separate
findings, not one: (a) `<th scope="col">` already carries the implicit role
`columnheader` — same redundant-ARIA reasoning as correction #1, don't write
it explicitly; (b) **the entire weekday header row is invisible to
assistive tech**, which correction #1's framing alone doesn't capture.
The brief asks explicitly whether reproducing `aria-hidden` is worth it:
**yes, reproduce it**, and the reason is legible from the rest of this map,
not a guess — `labelDayButton` (§ Inputs read, full file) already prepends
the day-of-week-bearing full formatted date (`"PPPP"`, e.g. "Sunday, July
26th, 2026") to every single day button's own `aria-label`, and
`labelGridcell` does the same for the non-interactive cell-label path. A
screen reader user tabbing/arrowing cell-by-cell already hears the weekday
name on *every* cell from that label; a `role="grid"` additionally
announcing "column 3, Tuesday" from an exposed header row on top of that
would be redundant chatter on every single navigation, not new information.
The header row's abbreviated visual labels ("Su", "Mo"...) exist for sighted
scanning only. **Recommend gsxui's day-button/cell labels reproduce the
same full-date-includes-weekday pattern** (not resolved as a class string
here, but load-bearing for Task 1–3's own label-formatting decision) — doing
so is what makes hiding the header row safe rather than a silent
information loss.

**3. Day-state data attributes are on the `<td>`, not the `<button>` — and
a second, different attribute set belongs on the button too.** Fully traced
in §4/§3 above: the core library's six attributes (`data-day` [ISO],
`data-month`, `data-selected`, `data-disabled`, `data-hidden`,
`data-outside`, plus `data-focused`/`data-today`) land on the cell via
`DayPicker.js`'s own `React.createElement(components.Day, {...})` call.
Separately, shadcn's `CalendarDayButton` computes its own five attributes
(`data-day` [locale-formatted, not ISO — a different value under the same
name], `data-selected-single`, `data-range-start`, `data-range-end`,
`data-range-middle`) and puts them on the **button**, driving the button's
own pill/color styling via `data-[range-end=true]:*` etc. selectors in its
class string (§3). **Recommend porting both sets, on their respective
elements** — the cell's set for the `today`/`outside`/`disabled` classNames
slots (§2, which are cell-level styling in this architecture) and the
button's set for `CalendarDayButton`'s own selected/range pill visuals.
Collapsing them onto one element would either lose the `group/day` ring-
styling hookup (§3, which needs the *cell's* `data-focused`, read via
`group-data-[focused=true]/day:`) or lose the button's own range-pill
coloring.

**4. Missing `aria-multiselectable` and `role="status" aria-live="polite"`
— and a third, previously-unflagged instance of the same caption pattern.**
Confirmed at the cited lines: `aria-multiselectable={mode === "multiple" ||
mode === "range"}` on the grid (`DayPicker.js`, the `components.MonthGrid`
call), `role="status" aria-live="polite"` on the non-dropdown caption label
(`components.CaptionLabel` call) and, for the dropdown-caption path, on an
equivalent visually-hidden `<span>` carrying the same formatted caption text
(both paths announce month/year changes; the dropdown path just does it
through a separately visually-hidden node instead of the visible label,
since the visible label in dropdown mode is a `<select>`, not naturally an
announcement target). **New, not in the original four**: `Footer`
(`DayPicker.js`'s `components.Footer` call) gets the identical `role="status"
aria-live="polite"` treatment. Not load-bearing for this port (`footer` is
explicitly not being ported, §9), but recorded so Task 7's ledger entry for
`footer` names the correct reason it's skippable-without-loss (it's an
announced region, same pattern as the caption, not a distinct pattern
needing its own decision if it's ever added later).

**5. New finding — the `disabled`/`aria-disabled` split on the day button
(§7.2) is a real behavioral rule the plan's spec/tests should account for,
not an incidental detail.** A disabled day does **not** always render with
the native `disabled` attribute — only when it is not simultaneously the
live-focused day, in which case it degrades to `aria-disabled` so the
element doesn't lose DOM focus out from under the user. If Task 1's tests
assert "every disabled day has `disabled` set" unconditionally, they will
be wrong for this one edge case (a day that becomes disabled while it
already holds live focus, e.g. via a controlled `disabled` prop update).
Flagging explicitly per the brief's instruction to record anything the plan
assumed that the sources contradict.

**6. New finding, added in the Task 0 review-fix pass — the nav buttons had
no written disabled contract at all.** Full trace in §7.5. Summary: neither
`button_previous` nor `button_next` ever receives a native `disabled`
attribute, at any navigation bound, under any configuration — only
`aria-disabled={true}` + `tabIndex={-1}`, confirmed by reading every prop at
all four places `Nav.js`/`DayPicker.js` construct these buttons. Task 3 had
no citation to build against and could plausibly have reached for native
`disabled` at a month boundary, which would remove the button from the tab
order entirely — exactly what upstream's `aria-disabled`-only approach is
designed to avoid. §7.5 also traces that `fromYear`/`toYear`/`startMonth`/
`endMonth` are not a separate "configured bounds" concept from the buttons'
"data bounds" — they all resolve into the same `navStart`/`navEnd` pair
that `getPreviousMonth`/`getNextMonth` test against, so the one contract
above covers both.

---

## 9. Slots NOT being ported

Per the plan's own scope (already decided, listed here so Task 7's ledger
entry is pre-written, not re-derived):

- **`week_number` / `week_number_header`** — the week-number gutter column
  (`showWeekNumber` prop). Note for the ledger: shadcn's own `calendar.tsx`
  already overrides react-day-picker's default `WeekNumber` component (a
  `<th scope="row" role="rowheader">`, §5) with a plain `<td>` — even if this
  slot is revisited later, the semantics to decide against are shadcn's
  `<td>`, not upstream's `<th>` rowheader, since shadcn is what a real port
  would start from.
- **`months` (multi-month display)** — `numberOfMonths > 1`. gsxui's port
  targets a single month grid; the `months`/`month` wrapper slots (§2) exist
  in the class-string list because they're unconditionally present even for
  one month, not because multi-month is in scope.
- **`footer`** — `props.footer`, rendered as a `role="status" aria-live=
  "polite"` div after the months container (§8 finding 4). Arbitrary
  caller-supplied content; no fixed markup/class contract to port beyond the
  wrapper attributes just named, which are cheap to add later if this slot
  is picked up in a future task.

---

## 10. Deviations already decided — recorded, not relitigated

Per the calling context, these two are settled and are **not** open
questions for this document or for Tasks 1–7 to revisit:

1. **The grid is always six rows.** Upstream defaults to sizing each
   month's grid to however many weeks that month actually spans (4–6,
   confirmed directly: `getMonths.js`, full file — `monthDates` is built
   from the real first/last week of the month, and is only padded out to a
   fixed 42-day/35-day (`broadcastCalendar`) grid when the caller passes
   `fixedWeeks: true`; six rows is not the default, it's an opt-in). gsxui
   always renders six weeks regardless of this prop.
2. **`captionLayout="dropdown"` defaults its year range** instead of
   requiring `startMonth`/`endMonth` from the caller, contra upstream's own
   requirement traced in §7.4 (`getYearOptions` returns `undefined` without
   them).

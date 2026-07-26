# calendar — a react-day-picker equivalent, server-rendered

The last component of Tier 4 that we are building. `chart` stays deferred:
shadcn's `chart.tsx` is a themed shell around Recharts — two of its six
exports are literal re-exports of Recharts components — and gsxui's copy-in
model has no way to vendor an npm dependency. Coverage after this ships is
53 of shadcn's 61 (`ls ui/*.gsx | wc -l`; an earlier draft of this line said
60, which was never right).

`calendar` is the largest single port in the project. shadcn's `calendar.tsx`
contains **no date logic at all**: it maps 24 className slots and overrides
four components of `react-day-picker`, which owns the month math, the
modifiers, the selection models and the keyboard grid. Porting it means
writing that engine.

## 1. What we looked at

Two references, used differently.

**shadcn's `new-york-v4/ui/calendar.tsx`** is the markup and class-string
authority, as for every other component. `styles/style-nova.css` is the
visual authority where the two disagree.

**A working templ + vanilla-JS calendar from another design system** was
read in full for its approach. Its rendering architecture is the opposite of
what we want — the templ file emits two empty divs and JavaScript fills the
weekday header and the day grid, so no-JS renders nothing — but its state
model is good and we take four things from it:

- **ISO date strings as the DOM state format.** `data-…-selected="2026-07-25"`
  and friends. Server-renderable, timezone-unambiguous, and every visual
  state derives from them.
- **UTC-only date arithmetic.** Local-midnight arithmetic is where
  off-by-one-day bugs breed. The JavaScript side uses one exact-year UTC
  constructor built from `setUTCFullYear`, not `Date.UTC(year, ...)`:
  ECMAScript deliberately remaps numeric years 0–99 to 1900–1999, while Go's
  `time.Time` represents those years literally.
- **Range hover preview as DOM state.** Hovering sets a `hover` attribute;
  the renderer computes the preview range from it and swaps the ends if the
  hover precedes the start. The two-click commit machine is the same shape.
- **Form `reset` handling** that restores the server-rendered values and
  returns the view to today — `ui/combobox.js` already reflects native reset
  state back into its visible control, so this is house style.

What we do not take: classes emitted from JavaScript, `innerHTML +=` in a
loop, a full grid rebuild on every `mouseover`, an Alpine dependency for the
month/year menus, 201 rendered year elements per calendar, a
`MutationObserver` on `document.body`, and five raw `document.addEventListener`
calls outside the delegation core. That reference also has **no keyboard grid
and no ARIA grid roles at all**, which is most of what react-day-picker is
for.

## 2. The rendering split

**Go renders a fixed 6-row, 42-cell grid. JavaScript never creates or
destroys a cell** — it updates `textContent` and `data-*` attributes in
place.

This is the load-bearing decision, and it follows from a constraint that is
easy to get wrong in the other direction. The alternative — JS building the
grid — would work: Tailwind v4 auto-detects `.js`, so classes emitted from
JavaScript do compile (verified: `animate-caret-blink` appears in
`ui/input-otp.js` and in no `.gsx` file, and is present in the compiled
stylesheet). `ui/input-otp.js` and `ui/sonner.js` already do it. So the
objection to a JS-built grid is not that it cannot be styled.

The objection is what it costs. A JS-built grid means: no calendar without
JS, a first-paint flash, class strings split across two files with no single
authority, and DOM churn that destroys focus on every navigation. A fixed
cell grid gives all of it back:

- Every class string stays in `calendar.gsx`, expressed as data-attribute
  variants (`data-selected:`, `data-outside:`, `data-range-middle:`) — how
  gsxui already styles state, and how shadcn's own `CalendarDayButton` marks
  days.
- JavaScript sets no classes. It sets `data-*`, `textContent`, `tabindex`,
  `aria-selected` and `disabled`.
- Focus survives navigation, because the focused element still exists.
- Height never jumps between a 5-row and a 6-row month.

The duplication this leaves is the date math, and only the date math:

```
Go:  monthGrid(year int, month time.Month, weekStartsOn time.Weekday) [42]time.Time
JS:  monthGrid(year, month, weekStartsOn) -> Date[42]
```

A pure function of three integers returning 42 dates, starting at the
week-start on or before the 1st. Both implementations use UTC. §7 covers how
we keep them honest.

Fixed 6 rows is a deliberate divergence from react-day-picker, which renders
4–6 rows and changes height between months. Constant height is what makes
the stable-cell design work, and it is the better behavior. Ledger it.

## 3. Parts and parameters

One `Calendar` component, matching shadcn's single-export shape (it exports
`Calendar` and `CalendarDayButton`; the latter exists only because
react-day-picker needs a component reference, and has no gsx equivalent).

| parameter | meaning |
|---|---|
| `mode` | `single` (default), `range`, `multiple` |
| `selected` | `[]time.Time` — one entry for single, N for multiple |
| `from`, `to` | range ends, `range` mode only |
| `month` | a `time.Time` whose year and month select the displayed grid; the day is ignored. Defaults to `from` in range mode, else the first `selected` date, else today |
| `weekStartsOn` | `time.Weekday`, default Sunday |
| `showOutsideDays` | default true |
| `captionLayout` | `label` (default) or `dropdown` |
| `fromYear`, `toYear` | bounds for the dropdown caption and for navigation |
| `disabledBefore`, `disabledAfter` | date bounds |
| `disabledDates` | specific dates |
| `disabledWeekdays` | e.g. weekends |
| `name` | renders the hidden form input(s) |

`captionLayout="dropdown"` renders `ui.NativeSelect` elements — real
`<select>`s, which is what react-day-picker renders too. Server-rendered,
zero JS, and it sidesteps the other reference's Alpine-filtered 201-item
menus entirely.

Dependencies are `button`, `icon` and `native-select`, all Go-level imports,
so `registry.Deps("calendar")` derives them. No JS-only dependency edge —
those are invisible to the registry's Go-only parse and are forbidden.

## 4. State, events, forms

Selection is a **server-rendered parameter reflected into the DOM**, and a
change emits `gsxui:change` with `{ mode, selected }`. The component never
chooses a persistence mechanism. This is the standing principle from
`2026-07-24-tier4-batch-a-design.md` §4.

Root state, all ISO strings:

```
data-gsxui-calendar-selected   comma-separated ISO dates (single, multiple)
data-gsxui-calendar-from       range start
data-gsxui-calendar-to         range end
data-gsxui-calendar-hover      range hover preview target
data-gsxui-calendar-month      the displayed month, "2026-07"
```

Per-cell state, all set by Go on first paint and by JS thereafter:
`data-date`, `data-outside`, `data-today`, `data-selected`,
`data-range-start`, `data-range-middle`, `data-range-end`, `data-disabled`.

With `name` set, hidden inputs carry the ISO value — the same bridge `select`
and `combobox` use. Single mode renders one input. Range mode renders a
`name` / `name + "-to"` pair. Multiple mode renders one input per selected
date, all with the same `name`, which is the native form representation of a
multi-valued field. JavaScript may create and remove these hidden form
controls as the selection changes; the fixed-DOM rule applies to the 42 day
cells, whose identity preserves focus, not to unrelated form bridge nodes.

Form `reset` restores the server-rendered selection values, clears live and
sticky focus bookkeeping, and returns the displayed view to the client's
current month.

**Hook attributes are namespaced `data-gsxui-calendar-*`**, per the Batch B
collision. `ui/gsxui.js` dispatches to every handler whose selector matches,
regardless of module, so a shared component hook prefix silently runs two
components' handlers on one event. The selector-disjointness invariant in
`jstest/specs/invariants.spec.ts` enforces this by default.

Form reset is one reviewed exception: a form composed from both Calendar and
Combobox intentionally matches both modules' reset handlers, and both must
run because each owns different descendant state. The exact
`["calendar.js", "combobox.js"]` / `"reset:false"` pair belongs in
`allowedOverlaps` with this reason; a mixed-form regression proves both
components reset correctly together.

## 5. Keyboard and ARIA

This is most of react-day-picker's value and the part the other reference
does not have.

`role="grid"` on the day table, `role="row"` per week, `role="gridcell"` per
cell, `aria-selected` per day, and a full accessible label on each day
button. Roving tabindex: exactly one day carries `tabindex="0"`.

| key | action |
|---|---|
| `←` `→` | ±1 day |
| `↑` `↓` | ±7 days |
| `Home` / `End` | first / last day of the displayed week |
| `PageUp` / `PageDown` | ∓1 month |
| `Shift+PageUp` / `Shift+PageDown` | ∓1 year |
| `Enter` / `Space` | select the focused day |

Movement that crosses a month boundary navigates the grid and lands focus on
the target day. Disabled days are focusable but not selectable, per the ARIA
grid pattern — skipping them strands keyboard users on either side of a
disabled span. Consequently, if the one roving `tabindex="0"` target is
disabled, it uses `aria-disabled="true"` rather than native `disabled`; this
keeps the grid reachable by Tab even when the first in-month day (or every
in-month day) is disabled.

## 6. Timezone

The server marks "today" in the server's timezone; a client in another zone
may be on a different date. On load, JavaScript re-marks `data-today` from
the client's local date. One attribute flip, and it is the only place the two
can legitimately disagree. Ledger it.

## 7. Testing

Calendar is the first component built after the JS test layer exists, so it
ships with its own Playwright specs rather than waiting for Phase 2. Three
things get tested that could not have been tested before:

**Go/JS grid agreement.** The single largest risk in §2. For a set of months
chosen to cover the awkward cases — 28-day February, leap February, a month
starting exactly on the week start, a month needing all 6 rows, a year
boundary — navigate client-side to month N and assert the resulting 42
`data-date` values equal those the server renders for month N directly. This
is exactly the kind of check the harness makes cheap: both pages come from
the same server.

**Keyboard grid.** Every row of §5's table, plus the month-boundary crossing
and focus landing.

**Selection semantics.** Single replaces; multiple toggles; range takes two
clicks, swaps when the second precedes the first, and previews on hover.
Each emits exactly one `gsxui:change`.

**Review regressions.** The public zero-value month follows the fallback
order in §3; years below 100 survive a Go → JS navigation round trip without
the ECMAScript 1900 offset; a disabled initial roving stop is reachable by
Tab; multiple-mode form data contains every selected date; reset returns the
view to the client's current month; and a form containing Calendar plus
Combobox resets both descendants.

The four existing invariants apply automatically once examples are
registered — including selector disjointness against all 21 other modules,
subject only to reviewed entries in `allowedOverlaps`.

## 8. Scope

**In:** everything above.

**Deferred, ledgered in `docs/jsx-parity.md`:** multi-month
(`numberOfMonths`), week numbers, custom modifiers, locales and i18n (month
and weekday names are English; both sides must agree, so `Intl` on the client
alone is not an option), and the footer slot.

**Not a component:** `DatePicker`. shadcn's is a docs composition of popover
+ button + calendar, not a registry entry. It ships as an example.

## 9. Cross-cutting constraints

Unchanged, restated because they bind every task:

- Markup reference `new-york-v4/ui/calendar.tsx`; visual reference
  `styles/style-nova.css`; nova wins on disagreement. Never port from
  `bases/*` — those emit `.cn-*` named classes.
- Class carryover is token-for-token; every drop or deviation ledgered.
- Behavior registers through `on()` in `ui/gsxui.js`. No raw
  `document.addEventListener`, no `MutationObserver`, no init scan.
- A dependency existing only in a JS import is invisible to `registry.Deps`
  and forbidden.
- Any `site/examples/**` change requires `make highlight` in the same commit.
- Adding a `ui/*.gsx` breaks `internal/registry/registry_test.go`'s pinned
  component list; updating it is expected.
- `make check` must pass: Go suite, Playwright suite, generated-file drift,
  `node --check`, gofmt.

# Component roadmap — shadcn coverage audit

Audited 2026-07-23 against `shadcn-ui/apps/v4/registry/new-york-v4/ui`
(57 registry components). gsxui ships 20 (naming deltas: our `dropdown` =
their `dropdown-menu`, `radio` = `radio-group`, `select` = their
`native-select` — shadcn now ships the same div-wrapped-native-select
design we ported, which is the alignment target for its class string).

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

## Tier 3 — new interactive machinery, moderate JS

| component | approach |
|---|---|
| toggle-group | toggle + single/multi selection state, roving focus (dropdown has the roving-focus pattern) |
| slider | styled native `<input type="range">` (form-native; cross-browser track/thumb CSS is the work) |
| scroll-area | ADAPT: CSS `scrollbar-width`/`scrollbar-color` styling first; Radix-style custom thumbs only if that falls short |
| select (custom listbox) | **user-promoted from backlog**: dropdown machinery + value selection, check indicator, keyboard typeahead; the native port stays alongside (rename decision: ours → `NativeSelect` to mirror shadcn naming) |
| sonner (toasts) | own toast module: stacking, timers, exit animations; no Radix/sonner dependency |
| drawer | sheet variant; v1 without vaul's drag-to-dismiss gesture (ledger the gap) |
| carousel | ADAPT: CSS scroll-snap + prev/next JS (shadcn wraps embla; snap covers the docs demos) |
| input-otp | segmented code input: paste splitting, focus advance, hidden real input |

## Tier 4 — hard / composite

| component | approach |
|---|---|
| combobox | input + filtered listbox on popover anchoring; builds on custom select |
| command | command palette (cmdk equivalent): filtering, groups, keyboard nav; builds on combobox |
| navigation-menu | hover mega-menu with viewport panel transitions |
| menubar | nested menus, submenu positioning, full keyboard model |
| calendar | month grid, range selection (react-day-picker equivalent — large); consider `<input type="date">` ADAPT as stopgap |
| resizable | drag-resized split panes + keyboard resize |
| sidebar | large composite (collapsible rail, mobile sheet mode, provider state); depends on sheet + tooltip + collapsible |
| chart | recharts wrapper in shadcn — needs a whole Go/JS charting answer; defer until demanded |

## Not ported (deliberate)

- **form** — react-hook-form bindings; meaningless server-side. `field`
  (Tier 1) is the layout half; validation is Go handler + `aria-invalid`
  patterns, to be shown in the future patterns/page-examples phase.
- **direction** — RTL context provider; HTML `dir` attribute serves gsx.
- **attachment, bubble, message, message-scroller, marker** — the new AI
  chat primitives; defer as a dedicated batch if gsx targets chat UIs.

# Tier 4 Batch B — the menu family

Design for the second Tier 4 batch. Batch A (resizable, combobox, sidebar)
shipped 2026-07-24 and took coverage to 50 of shadcn's 61.

Batch B closes the menu family in one pass, because all of it is the same
body of work: the missing item types and submenu machinery are shared by
`dropdown`, `context-menu` and the new `menubar`, and the README's
"dropdown checkbox/radio items + submenus" backlog entry IS this batch.
`navigation-menu` rides along as the family's odd member — a hover mega-menu
rather than a menu — because it is the last piece of shadcn's navigation
surface and shares the anchoring machinery.

Coverage after this batch: 52 of 61 (`menubar`, `navigation-menu`), with
`dropdown` and `context-menu` reaching full parity.

Remaining after it: `calendar` (its own spec — real date math, range
selection, keyboard grid) and `chart` (needs a Go/JS charting answer first).

## 1. The submenu mechanism — measured, not assumed

Radix portals submenu content to `document.body`. gsxui's menu family uses
the native popover API, where light-dismiss governs what closes what, so the
portal question had to be settled empirically. Measured in Chrome:

| child popover opened… | parent stays open? |
|---|---|
| DOM-nested inside the parent, via `showPopover()` | **yes** |
| not nested, via a real `popovertarget` invoker activation | **yes** |
| not nested, via programmatic `showPopover()` | **no — parent light-dismisses** |

A `popover="auto"` establishes an ancestor chain either by DOM containment or
by invoker relationship. Submenus open on hover and on `ArrowRight`, not only
on click, so they must be opened programmatically — which rules out relying
on the invoker chain.

**Decision: submenu content is DOM-nested inside its parent's content.** This
is an ADAPT away from Radix's portal, and it is load-bearing rather than
cosmetic — it is what keeps the parent menu open. Ledger it in
`docs/jsx-parity.md` under each affected component with the table above, so
nobody later "tidies" the submenu into a portal and silently breaks it.

Consequence to design around: nested content inherits the parent's stacking
and overflow context. The parent content must not clip (`overflow-visible`
where the reference allows), and the submenu is positioned against its own
trigger, not the parent.

## 2. Shared item types — dropdown and context-menu

Both ship only item/label/separator/shortcut today. Both are missing exactly
the same set, and shadcn's two sources are near-identical modulo the
component prefix:

| part | element / semantics |
|---|---|
| `…MenuGroup` | `role="group"` wrapper |
| `…MenuCheckboxItem` | `role="menuitemcheckbox"` + `aria-checked`, check indicator |
| `…MenuRadioGroup` | `role="group"` owning a radio set |
| `…MenuRadioItem` | `role="menuitemradio"` + `aria-checked`, dot indicator |
| `…MenuSub` | wrapper, no DOM of its own beyond grouping |
| `…MenuSubTrigger` | `role="menuitem"` + `aria-haspopup="menu"` + `aria-expanded`, chevron |
| `…MenuSubContent` | nested `popover="auto"`, `role="menu"` |

**State follows the batch-A principle** (`2026-07-24-tier4-batch-a-design.md`
§4, now a standing rule): checked state is a server-rendered parameter,
reflected in the DOM, and a change emits an event. No component stores it.
`CheckboxItem` takes `checked bool` and emits `gsxui:change` with
`{ checked, value }`; `RadioItem` takes `checked bool` and emits
`gsxui:change` with `{ value }` on the radio group. Matches `toggle`/`tabs`.

**No-JS behavior**: server-rendered `aria-checked` and the indicator are
correct on first paint; JS only handles interaction. This is strictly better
than Radix, which renders nothing until mount.

## 3. menubar

A horizontal bar of menu triggers, each owning a menu built from the same
item set as §2. What menubar adds over N independent dropdowns:

- **Roving tabindex** across triggers — one tab stop for the whole bar, with
  ArrowLeft/ArrowRight moving between triggers. Same model as the shipped
  `toggle-group.js`.
- **Open-follows-hover once open**: with one menu open, hovering a sibling
  trigger switches to it without a click. This is the behavior that makes a
  menubar feel native, and it is the only genuinely new interaction.
- Sibling menus are NOT nested in each other, so switching between top-level
  menus is a close-then-open, not a chain — the §1 constraint does not apply
  between siblings, only between a menu and its own submenus.

Reuses the shipped dropdown content machinery and the discrete-transition
block verbatim.

## 4. navigation-menu

The family's odd member: a hover-opened mega-menu whose panels share one
positioned viewport, with the viewport resizing between panels.

**v1 scope**: `NavigationMenu`, `List`, `Item`, `Trigger`, `Content`, `Link`,
`Viewport`, `Indicator`. Hover and focus open, `data-state` transitions,
viewport sizing from the active panel.

**Deferred, ledgered**: the animated width/height morph between panels of
different sizes (Radix measures each panel and animates
`--radix-navigation-menu-viewport-width/height`; ours will size to the active
panel without tweening between the two), and `Indicator`'s arrow-follows-
trigger animation if it proves to need per-frame measurement.

`navigation-menu` is the one part of this batch where the reference's
behavior depends on continuous measurement rather than discrete state. If
the source map finds the viewport morph needs a resize observer to look
right, ship it; if it needs per-frame tweening, defer it and say so.

## 5. Cross-cutting constraints

Unchanged from Batch A; restated because they bind every task:

- Markup reference `new-york-v4/ui/<name>.tsx`; visual reference
  `styles/style-nova.css`; nova wins on disagreement; `bases/base` is the
  markup-structure tiebreak only. Never port from `bases/*` directly — those
  emit `.cn-*` named classes.
- Class carryover is token-for-token; every drop or deviation ledgered.
- Components render state, never own persistence; events via `emit()` in
  `ui/gsxui.js` (`gsxui:change` with a detail payload, `gsxui:open`/`close`
  for overlays, `gsxui:select` for selection).
- Public JS hook attributes are `data-gsxui-*`; `data-slot` is markup
  identity, not a behavior contract. (Batch A's final review caught two
  components binding to `data-slot` — do not repeat it.)
- A dependency existing only in a JS import is invisible to `registry.Deps`
  and forbidden.
- `data-state="open"` is stamped synchronously BEFORE `showPopover()`.
- Author `display` utilities on a popover must be gated on `:open` — a bare
  `block`/`grid` beats the UA's closed-popover `display:none` and leaves
  hit-testable ghosts. This bit both `dialog` and `sidebar`.
- Any `site/examples/**/*.gsx` change requires `make highlight` in the same
  commit.
- Adding a `ui/*.gsx` breaks `internal/registry/registry_test.go`'s pinned
  component list; updating it is expected.

## 6. Task sequence

0. **Source map** — dropdown/context-menu shared parts, menubar,
   navigation-menu. Every ARIA and keyboard claim traced to Radix's real
   `dist/index.mjs`, or marked `derived-not-read`.
1. **Shared item types on `dropdown`** — group, checkbox, radio, sub. The
   submenu mechanism lands here first, since everything else reuses it.
2. **Same set on `context-menu`** — near-mechanical after task 1.
3. **`menubar`** — roving tabindex + open-follows-hover over task 1's parts.
4. **`navigation-menu`** — independent of 1–3; could run earlier if useful.
5. **Roll-up** — roadmap, README backlog entry retired, CHANGELOG, whole-batch
   review, live side-by-side browser pass.

## 7. Verification

As Batch A: per-task review, whole-batch review on the strongest model, then
a live side-by-side pass against ui.shadcn.com before calling it shipped.

Batch A's browser pass caught three plan-level errors that static review did
not. The specific things to exercise here, because they are the ones static
review cannot see:

- A submenu open while its parent stays open — the §1 mechanism, end to end.
- Hovering between top-level menubar menus with one open.
- Keyboard: ArrowRight into a submenu, ArrowLeft/Escape out, and that focus
  lands somewhere sane at every step.
- No ghost boxes: every closed menu and submenu must compute
  `display: none`.

Known environment limit: occluded Chrome tabs freeze animation clocks and
leave transitions permanently `pending`. Assert on resolved geometry with
transitions disabled; do not report frozen transitions as defects.

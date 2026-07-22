# gsxui

shadcn-style components for [gsx](https://github.com/gsxhq/gsx) — copy-in,
type-checked, server-rendered.

## Install

    go install github.com/gsxhq/gsxui/cmd/gsxui@latest

In your project (a Go module):

    gsxui init          # tokens css, js runtime, class merger, gsx.toml wiring
    gsxui add dialog    # vendors dialog + its deps (button), regenerates
    gsxui list          # what's available

You own the vendored code. `gsxui add` never touches a modified file unless
you pass `--overwrite`.

After upgrading the gsxui binary, re-run `gsxui add <name> --overwrite` to
refresh vendored components — this discards local edits to those files.

**Status: pre-release.** v1 component set (20 components + icon) is complete.
The showcase site and theme editor are in progress.

- Components live in `ui/<name>/` — a `.gsx` source (JSX-style, named
  parameters, fallthrough attrs) plus a behavior `.js` when interactive.
- `make test` regenerates and tests everything; `make check` adds JS syntax
  and gofmt checks.
- Divergences from the shadcn/ui reference: `docs/jsx-parity.md`.

## Components

**Form controls:** button, checkbox, input, label, radio, selectbox,
switchctl, textarea

**Display:** alert, avatar, badge, card, separator, skeleton, table

**Overlay:** dialog, dropdown, tooltip

**Navigation:** accordion, tabs

**Primitives:** icon (Lucide, generated — a dependency of selectbox,
accordion, and dropdown's own future items, not usually added directly)

Some native-first components (checkbox, radio, switchctl, selectbox,
accordion) trade a slice of shadcn's Radix-driven behavior for a real
`<input>`/`<select>`/`<details>` element — zero client JS, browser-native
`:checked`/`:disabled`/exclusivity semantics. dropdown and tooltip trade
Radix's Portal for the native popover API. Every divergence, with its
rationale, is ledgered in `docs/jsx-parity.md`.

## Post-v1 backlog

Deferred out of v1 scope, tracked here rather than in the parity ledger's
per-component GAP notes (see those for the detailed rationale):

- **Custom listbox select** — shadcn's full Radix `Select` (floating panel,
  check-mark item indicator, keyboard typeahead). v1 ships a styled native
  `<select>` instead (`ui/selectbox`); the Radix-equivalent listbox visual
  is not built.
- **Dropdown checkbox/radio items + submenus** — `DropdownMenuCheckboxItem`,
  `DropdownMenuRadioGroup`/`DropdownMenuRadioItem`,
  `DropdownMenuSub`/`SubTrigger`/`SubContent`, `DropdownMenuGroup`. Only the
  base item/label/separator/shortcut set ships in v1.
- **Tooltip delay-groups** — `TooltipProvider`'s shared `delayDuration`/
  skip-delay-group coordination across multiple tooltips on a page. v1
  hard-codes a fixed per-trigger open delay, no cross-tooltip grouping.
- **CSS anchor positioning migration** — dropdown and tooltip currently
  position via a hand-rolled `getBoundingClientRect()` + `position: fixed`
  calculation in JS. Once CSS anchor positioning (`anchor()`/
  `position-anchor`) reaches Baseline across browsers, both can drop that
  JS for native, collision-aware placement.
- **Popover exit-animation strategy** — verified inert, strategy TBD:
  dropdown/tooltip's native popover is already `display: none` at the
  moment the `toggle` (newState=closed) handler stamps
  `data-state="closed"`, so `data-[state=closed]:animate-out` never runs
  and closing snaps (open-path `animate-in` unaffected); accepted for v1.
  An animated-close strategy (`beforetoggle`/`allow-discrete`, or
  dialog-style `requestClose`) remains adoptable once designed.
- **Checkbox checkmark theming** — the check glyph is a data-URI with
  hard-coded `stroke="white"`; data-URIs are static text and can't reference
  CSS variables, so the mark doesn't follow `--primary-foreground` and is
  wrong/low-contrast for themes where that color isn't near-white. Swap to a
  `currentColor` CSS-mask approach in the Plan 4 theming work.

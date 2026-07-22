# JSX parity ledger

Divergences between gsxui components and their shadcn/ui reference, both
directions. Full audit: gsxhq docs repo, specs/2026-07-22-gsx-over-jsx-audit.md.

## badge
- WIN: `cva()` variant map replaced by `switch` inside `class={}`.
- GAP (narrow): shadcn's `asChild` tag-swapping (render the badge as an `<a>`)
  has no gsx equivalent — no dynamic tag. Behavior-attachment uses of
  `asChild` are covered by the data-attribute mechanism (see dialog).

## button
- GAP (narrow): `asChild` tag-swapping (no dynamic tag: `const Comp = asChild ? Slot : "button"`). Ported as `href` param rendering `<a>` — covers the dominant use. Behavior-attachment uses of `asChild` are covered by the data-attribute mechanism (see dialog).
- WIN: `type="button"` before `{ attrs... }` makes it an overridable default — positional spread precedence replaces prop-ordering conventions.
- WIN: `cva()` replaced by plain Go variant/size funcs shared by both branches.

## card
- Straight port; package-namespaced compound parts (`card.CardHeader`) replace module exports. No divergences.

## dialog
- WIN: Radix Portal/Overlay replaced by native <dialog> top layer + ::backdrop; Esc handling is browser-native.
- ADAPT: `text-foreground` added to DialogContent's classes — native <dialog> gets UA `color: CanvasText` and does not inherit the themed body color (Radix's <div> content does); without it dark mode renders wrong text color.
- ADAPT: close button drops shadcn's `data-[state=open]:bg-accent data-[state=open]:text-muted-foreground` — we stamp `data-state` on the <dialog> element, not the close button, so those selectors are dead in this DOM.
- CONVENTION (decided 2026-07-22): gsx keeps Go zero-value semantics — no default parameter values in the language (designs using exported consts and an init-time registry were considered and rejected: exported-symbol pollution vs. runtime lookup, neither worth it). Name parameters so the zero value is the shadcn default: bools invert (`hideCloseButton`), numerics document "0 means the default" (e.g. upcoming `sideOffset`), strings use `""` = default.
- MECHANISM (decided 2026-07-22): `asChild`/Slot is not ported and not needed for behavior attachment — the `data-gsxui-*` attributes are each interactive component's PUBLIC contract, and fallthrough attrs deliver them to any element or component: `<button.Button data-gsxui-dialog-trigger>Open</button.Button>` makes your styled Button the trigger, no cloning, no wrapper. Document this idiom prominently per interactive component; only tag-swapping remains a (narrow, accepted) gap.
- GAP: Radix client context (trigger↔content wiring) replaced by closest("[data-gsxui-dialog]") proximity in JS.
- NOTE: controlled open/onOpenChange not ported; JS CustomEvents (gsxui:open/close) + dialog.showModal()/close() are the programmatic API. State + events ride ToggleEvent (Baseline 2024) — the native close event proved undelivered in current Chrome, and toggle also covers programmatic open/close, which close-based wiring never could.
- A11Y: aria-labelledby/-describedby/-controls stamped by JS with lazily generated ids (authored ids/aria always win); aria-haspopup + initial aria-expanded server-rendered on DialogTrigger; aria-expanded synced on toggle.
- ADAPT: exit animations run by stamping data-state="closed" before close() (Esc included, via intercepted cancel); a direct programmatic dialog.close() skips the exit animation — native-immediacy divergence, accepted.
- WIN: DialogFooter's showCloseButton ports with zero friction — shadcn defaults it false, so the Go zero value IS the default; its Close button uses the data-attribute idiom (<button.Button data-gsxui-dialog-close>) where shadcn needs <DialogClose asChild><Button>>.

## checkbox
- ADAPT (native-first): Radix's `CheckboxPrimitive.Root` (button role, `aria-checked`, hidden real input, `data-[state=checked]`) is replaced by a real `<input type="checkbox">` — form-native, zero JS, browser `:checked`/`:disabled` truth. `data-slot="checkbox"` and every color/focus/disabled/aria-invalid token from `registry/new-york-v4/ui/checkbox.tsx` are carried over verbatim; `appearance-none` is added to suppress the UA checkbox box so the custom border/background show instead (mechanical necessity of the native-input swap, not present on the Radix version since it renders no native control at all).
- MECHANISM: the `Indicator`/`CheckIcon` child is not a child at all on a void `<input>` — its checkmark becomes a `checked:bg-[url('data:image/svg+xml...')]` data-URI background (lucide's `check` path, `M20 6 9 17l-5-5`, verified byte-for-byte against `ui/icon/icon_data.go`'s `"check"` entry) plus `checked:bg-center checked:bg-no-repeat checked:bg-[length:12px_12px]`, swapped in only under the `checked:` variant. `data-[state=checked]:text-primary-foreground` and the redundant `dark:data-[state=checked]:bg-primary` are dropped — both were about the (now nonexistent) Indicator child's icon color/duplicate dark override.
- FINDING (verified, corrected in place — see checkbox_test.go's `TestCheckboxNoStraySpaceInDataURI`): the data-URI's embedded SVG markup needs literal spaces (attribute and path-data separators). A literal space inside a `class="..."` value is a token boundary by the HTML spec itself — every whitespace-splitting class tool, including the configured `merge.Merge` (`tailwind-merge-go`), tears the token apart at each space. Reproduced directly against `merge.Merge`: with real spaces, `checked:bg-primary` was dropped entirely (its `bg` category lost to the accidentally-split `bg-[url(...` fragment) and interior tokens `stroke-width`/`stroke-linecap` were dead-code-eliminated as fake "conflicting" utilities against `stroke-linejoin`, corrupting the SVG. Fix: every space inside the bracketed arbitrary value is written as `_`, Tailwind's own documented escape for whitespace in arbitrary values (`bg-[length:12px_12px]` in the same string already used this convention — the plan's given string was inconsistent, using literal spaces only inside the URL). Verified the underscore form round-trips `merge.Merge` unchanged. This is Tailwind's real build-time requirement independent of gsxui — not a gsxui-specific workaround.

## icon
- WIN: shadcn/templUI ports wrap each Lucide React component (or a `<template>`-per-icon component) individually; gsx's tag-callable values (`func(attrs ...gsx.Attr) gsx.Node`) let a single generated `New(name)` factory back every icon var (`var ChevronDown = New("chevron-down")`), so `<icon.ChevronDown class="size-4"/>` is both markup-callable and a plain Go value, generated from one shared `svgIcon` component instead of 1,748 near-duplicate wrapper components.
- WIN: `aria-hidden="true"` is authored before `{ attrs... }` in `svgIcon` — positional spread precedence (the same idiom as badge's `data-variant` and dialog's `data-state`) makes it an overridable default: a caller's own `aria-hidden` (e.g. `aria-hidden="false"` alongside `aria-label`) wins with no conditional logic.
- MECHANISM: unknown icon names are a render-time error (`New("nope")` → `unknown icon "nope"`), never a silently empty `<svg>` — mirrors the hard-error idiom used elsewhere in gsxui for unrecognized identifiers, so a typo'd icon name fails loudly instead of shipping a blank glyph.

# JSX parity ledger

Divergences between gsxui components and their shadcn/ui reference, both
directions. Full audit: gsxhq docs repo, specs/2026-07-22-gsx-over-jsx-audit.md.

## alert
- WIN: `cva()` variant map replaced by `switch` inside `class={}` (default |
  destructive), the same idiom as badge/button. No `data-variant` stamp —
  shadcn's own Alert doesn't stamp one either (unlike Badge/Button), so
  there's nothing to port.

## avatar
- ADAPT: AvatarImage adds `absolute inset-0` (not in shadcn) — the image overlays the in-flow fallback, so the no-JS/pre-JS state renders correctly (fallback behind, image covers when loaded); `ui/avatar/avatar.js` handles only the error path (hide broken image). Radix's mount-gated rendering can't exist server-side.
- GAP: `AvatarBadge`, `AvatarGroup`, `AvatarGroupCount` (added to shadcn's
  registry after the base three parts) are not ported — out of scope for
  this task; only `Avatar`/`AvatarImage`/`AvatarFallback` per the task
  brief.
- ADAPT: `Avatar`'s `size` prop (default/sm/lg) is dropped along with it, so
  `data-size` is never stamped and the size-keyed selectors that depend on
  it are dead weight — `group/avatar` and `data-[size=lg]:size-10
  data-[size=sm]:size-6` are dropped from `Avatar`'s class, and
  `group-data-[size=sm]/avatar:text-xs` from `AvatarFallback`'s (the same
  "drop the selector, don't ship dead CSS" call as dialog's close-button
  `data-[state=open]:...` ADAPT). `size-*` stays fully overridable via the
  ordinary caller-class-merge mechanism on `Avatar`/`AvatarFallback`
  directly.

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

## radio
- ADAPT (native-first): Radix's `RadioGroupPrimitive.Item` (button role, `aria-checked`, Indicator) is replaced by a real `<input type="radio">` — same rationale as checkbox: form-native, zero JS, browser `:checked`/`:disabled` truth. Token set carried verbatim from `registry/new-york-v4/ui/radio-group.tsx`'s `RadioGroupItem` (`aspect-square size-4 shrink-0 rounded-full border border-input text-primary shadow-xs transition-[color,box-shadow] outline-none` plus focus-visible/disabled/aria-invalid/dark tokens); `appearance-none` is added for the same mechanical reason as checkbox.
- ADAPT (native-first): shadcn's `RadioGroup` root (a `grid gap-3` wrapper coordinating Radix's roving-tabindex/keyboard-nav group) is not ported. Native `<input type="radio" name="...">` siblings already form a group via the browser's own `name` attribute — arrow-key roving and single-selection are UA behavior, no JS required. The `grid gap-3` layout is not component behavior; it is the caller's own wrapper `<div>`.
- MECHANISM: the `Indicator`/`CircleIcon` child (`fill-primary`) is not a child at all on a void `<input>` — its filled dot becomes `checked:bg-[radial-gradient(circle_closest-side,currentColor_45%,transparent_50%)]`, painted in `currentColor` rather than checkbox's data-URI approach: a data-URI is static text and can't reference the caller's CSS custom properties (`--primary` can be themed/overridden per caller, per dark mode, etc.), but a `currentColor`-based gradient computed live on the element can. `text-primary` is kept (not dropped) specifically to make `currentColor` resolve to the primary color on this element — it is load-bearing for the dot's color, unlike checkbox's icon color which came from an explicit `fill`/`stroke` in the data-URI itself. (An earlier version of this port wrongly copied checkbox's fill-the-whole-circle-plus-white-dot recipe — `checked:bg-primary checked:border-primary` plus a white-circle data-URI — which does not match shadcn's actual outlined-circle-with-colored-dot visual; corrected here.) Same underscore-not-space escaping as checkbox's data-URI (see checkbox's FINDING entry above) — verified the radial-gradient string round-trips `merge.Merge` unchanged (`TestRadioNoStraySpaceInGradient`).

## switch
- NOTE (package/directory naming): the Go package cannot be named `switch` — it is a reserved keyword, illegal both as a package declaration (`package switch` is a parse error) and as an import alias (`import switch "..."` is also a parse error — there is no spelling of the directory name that makes `switch.Switch(...)` legal Go, unlike `select`, which at least parses as an identifier in non-statement position). The component lives in `ui/switchctl` (package `switchctl`, file `switch.gsx`), mirroring the plan's own established precedent for `select` → `ui/selectbox`/package `selectbox` (task-6-brief.md). Consequence accepted the same way task 6 accepts it: the registry/CLI-facing component name is the directory name, so `gsxui add switchctl` (not `gsxui add switch`) — `internal/registry`'s `Components()`/`Deps()` walk directory names, and there is no separate "logical name" layer to reconcile this without inventing one (out of scope here; a `registry` alias-name feature would need its own decision, not a Task 5 improvisation).
- ADAPT (native-first): Radix's `SwitchPrimitive.Root` + separate `SwitchPrimitive.Thumb` span replaced by one real `<input type="checkbox" role="switch">` — form-native, zero JS, browser `:checked`/`:disabled` truth. `role="switch"` is kept explicitly (a bare checkbox input has no switch semantics of its own) since it's the one piece of Radix's ARIA contract a native checkbox doesn't supply for free. The `data-size="sm"|"default"` variant and its `group/switch` + `group-data-[size]/switch:` machinery are dropped — Task 5 ships default size only (no size variant asked for), and `group/switch` existed solely to let a *sibling* Thumb element read the Root's size; with the thumb now the same element's own pseudo-element, there is no sibling to target and the group plumbing is dead weight.
- MECHANISM: the Thumb span becomes this input's own `::before` pseudo-element (`before:` variant, thumb-span→before:) — track and thumb are the same DOM node instead of parent/child. A pseudo-element renders nothing without an explicit `content` (default `content: normal` produces no box at all, unlike a real child element which always has one) — Tailwind's `before:content-['']` is therefore load-bearing here and has no analog on the Radix Thumb span. Native `checked:` replaces `data-[state=checked]`/`data-[state=unchecked]:` throughout; an unchecked-specific class is unnecessary (the bare, unprefixed utility already covers the unchecked default, exactly as for checkbox/radio). `ring-0` itself is dropped as a dead reset (nothing `before:` could have a ring here).

## icon
- WIN: shadcn/templUI ports wrap each Lucide React component (or a `<template>`-per-icon component) individually; gsx's tag-callable values (`func(attrs ...gsx.Attr) gsx.Node`) let a single generated `New(name)` factory back every icon var (`var ChevronDown = New("chevron-down")`), so `<icon.ChevronDown class="size-4"/>` is both markup-callable and a plain Go value, generated from one shared `svgIcon` component instead of 1,748 near-duplicate wrapper components.
- WIN: `aria-hidden="true"` is authored before `{ attrs... }` in `svgIcon` — positional spread precedence (the same idiom as badge's `data-variant` and dialog's `data-state`) makes it an overridable default: a caller's own `aria-hidden` (e.g. `aria-hidden="false"` alongside `aria-label`) wins with no conditional logic.
- MECHANISM: unknown icon names are a render-time error (`New("nope")` → `unknown icon "nope"`), never a silently empty `<svg>` — mirrors the hard-error idiom used elsewhere in gsxui for unrecognized identifiers, so a typo'd icon name fails loudly instead of shipping a blank glyph.

## input
- Straight port. `type="text"` is authored before `{ attrs... }` — the same overridable-default idiom as button's `type="button"`, so `type="email"` etc. at the call site replaces rather than duplicates it.

## label
- Straight port of the rendered markup: shadcn wraps Radix's `LabelPrimitive.Root`, which is itself a plain `<label>`.
- GAP (narrow, accepted): Radix's `onMouseDown` handler, which calls `preventDefault()` on multi-click to stop text selection inside the label, is not ported (no client JS for this component). Low impact — the base class already carries `select-none`, which suppresses text selection via CSS regardless.

## separator
- ADAPT: Radix's `decorative` prop (default `true`, flips `role="separator"` +
  `aria-orientation` when `false`) is not ported — `Separator` always renders
  `role="none"`, matching shadcn's default usage. No orientation param needed
  for a semantic separator variant; callers wanting a semantic (non-decorative)
  separator fall through `attrs` to set `role`/`aria-orientation` themselves.
- WIN: no variant switch — the single verbatim class string dispatches on
  `data-orientation` via Tailwind's `data-[orientation=...]` selectors, so
  `orientation` only needs to stamp the attribute.

## skeleton
- Straight port; no divergences.
## textarea
- ADAPT: shadcn's Textarea takes its content via a `value` prop forwarded through `...props` onto React's controlled `<textarea value={...}>`. Native HTML `<textarea>` has no `value` attribute — its initial content is a text child. Ported as a `value string` param rendered as the (escaped) text child instead: `textarea.Textarea("initial text", nil)`.

## table
- NOTE: `Table` renders a scroll-container `<div data-slot="table-container">` wrapping `<table data-slot="table">`, matching shadcn's structure exactly. Fallthrough `attrs` land on the `<table>` element (where shadcn's `{...props}` lands), not the container div — the container is purely structural scroll-wrapping and has no props of its own in the source either.

## dropdown
- WIN: Radix Portal/Content replaced by the native popover API (`popover="auto"` on `DropdownMenuContent`) — top layer, light dismiss, and Esc are all browser-native, same win as dialog's `<dialog>`. Trigger↔content wiring is `closest("[data-gsxui-dropdown]")` proximity in JS, no ids, same MECHANISM as dialog.
- NOTE: positioning is a hand-rolled `position: fixed` anchor to the trigger's `getBoundingClientRect()` in `dropdown.js` (open below, left-aligned, 4px gap), not Radix's Floating-UI collision-aware placement. CSS anchor positioning (`anchor()`/`position-anchor`) — the eventual native replacement — is not yet Baseline (Chrome-only as of this writing); JS positioning is a stopgap until it is. No `data-side` is ever stamped, so `dropdown-menu.tsx`'s `data-[side=*]:slide-in-from-*` selectors in the ported class string are currently dead weight, kept for future-proofing once real placement logic lands.
- GAP: `DropdownMenuCheckboxItem`, `DropdownMenuRadioGroup`/`DropdownMenuRadioItem`, `DropdownMenuSub`/`DropdownMenuSubTrigger`/`DropdownMenuSubContent`, and `DropdownMenuGroup` are not ported — post-v1 backlog per the task brief. Only `DropdownMenu`, `DropdownMenuTrigger`, `DropdownMenuContent`, `DropdownMenuItem`, `DropdownMenuLabel`, `DropdownMenuSeparator`, `DropdownMenuShortcut` ship.
- ADAPT: the `inset` prop on `DropdownMenuItem`/`DropdownMenuLabel` is dropped along with it — `data-[inset]:pl-8` is removed from both classes as dead weight, the same "drop the selector, don't ship dead CSS" call as dialog's close-button ADAPT and avatar's size ADAPT.
- MECHANISM: `DropdownMenuItem` is a real menu item on `<div role="menuitem" tabindex="-1">` (Radix's own `Item` is also a non-native, ARIA-role div, so this isn't a native-first swap like checkbox/radio/switch — it's a straight structural port). `dropdown.js`'s arrow-key roving focus walks `[role="menuitem"]:not([aria-disabled])` within the content; click on an item emits `gsxui:select` then closes via `hidePopover()`.
- A11Y: `aria-haspopup="menu"` and the initial `aria-expanded="false"` are server-rendered on `DropdownMenuTrigger`; `aria-expanded` is synced by `dropdown.js` on the `toggle` event, same non-bubbling-event/capture pattern as dialog. On open, the first enabled menu item receives focus.
- NOTE: state + events ride `ToggleEvent` (`toggle`, capture-delegated) — same rationale as dialog's ledger NOTE (covers every open/close path including light dismiss and Esc, which a `close`-only or click-only wiring would miss).

## tooltip
- WIN: Radix Portal/Content replaced by the native popover API — `popover="manual"` on `TooltipContent` puts it in the top layer without the light-dismiss/Esc behavior `popover="auto"` would add (which would race hover-driven show/hide and dismiss on an unrelated outside click while the pointer is still over the trigger).
- NOTE: positioning is a hand-rolled `position: fixed` anchor centered above the trigger's `getBoundingClientRect()` in `tooltip.js`, not Radix's Floating-UI collision-aware placement — same stopgap-until-CSS-anchor-positioning rationale as dropdown's NOTE. No `data-side` is ever stamped, so the ported class string's `data-[side=*]:slide-in-from-*` selectors are currently dead weight.
- GAP: `TooltipProvider` (shared `delayDuration`/skip-delay-group machinery across multiple tooltips) is not ported — `tooltip.js` hard-codes a 300ms open delay per trigger, no cross-tooltip grouping. `TooltipPrimitive.Arrow` (the small pointing triangle) is also not ported — no `Arrow` part or class exists in `tooltip.gsx`.
- MECHANISM: `pointerover`/`pointerout`/`focusin`/`focusout` are used (not `mouseover`/`mouseout`/`focus`/`blur`) specifically because they bubble — dialog/dropdown's non-bubbling events (`toggle`) need `{ capture: true }`, but these delegate the ordinary way via `on()`'s bubble-phase default. `showPopover()`/`hidePopover()` are called directly rather than `togglePopover()` since show and hide are driven by two independent event pairs (pointer vs. focus), not one toggling control.
- A11Y: `TooltipContent` server-renders `role="tooltip"`; no `aria-describedby` is wired from trigger to content (Radix's Tooltip does this internally) — narrow GAP, not ledgered separately since the task brief scoped a11y to `role="tooltip"` only.

# JSX parity ledger

Divergences between gsxui components and their shadcn/ui reference, both
directions. Full audit: gsxhq docs repo, specs/2026-07-22-gsx-over-jsx-audit.md.

## badge
- WIN: `cva()` variant map replaced by `switch` inside `class={}`.
- GAP: `asChild` not ported (no dynamic tag in gsx — ledger: render-child/Slot).

## button
- GAP: `asChild`/Slot not expressible (no dynamic tag: `const Comp = asChild ? Slot : "button"` has no gsx equivalent). Ported as `href` param rendering `<a>`. Raise: slot/render-child feature.
- WIN: `type="button"` before `{ attrs... }` makes it an overridable default — positional spread precedence replaces prop-ordering conventions.
- WIN: `cva()` replaced by plain Go variant/size funcs shared by both branches.

## card
- Straight port; package-namespaced compound parts (`card.CardHeader`) replace module exports. No divergences.

## dialog
- WIN: Radix Portal/Overlay replaced by native <dialog> top layer + ::backdrop; Esc handling is browser-native.
- ADAPT: `text-foreground` added to DialogContent's classes — native <dialog> gets UA `color: CanvasText` and does not inherit the themed body color (Radix's <div> content does); without it dark mode renders wrong text color.
- ADAPT: close button drops shadcn's `data-[state=open]:bg-accent data-[state=open]:text-muted-foreground` — we stamp `data-state` on the <dialog> element, not the close button, so those selectors are dead in this DOM.
- GAP: `showCloseButton` defaults true in shadcn; gsx named params default to zero — inverted to `hideCloseButton`. Raise: default parameter values.
- GAP: Radix client context (trigger↔content wiring) replaced by closest("[data-gsxui-dialog]") proximity in JS.
- NOTE: controlled open/onOpenChange not ported; JS CustomEvents (gsxui:open/close) + dialog.showModal()/close() are the programmatic API. State + events ride ToggleEvent (Baseline 2024) — the native close event proved undelivered in current Chrome, and toggle also covers programmatic open/close, which close-based wiring never could.

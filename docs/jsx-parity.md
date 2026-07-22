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

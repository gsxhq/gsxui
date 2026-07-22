# JSX parity ledger

Divergences between gsxui components and their shadcn/ui reference, both
directions. Full audit: gsxhq docs repo, specs/2026-07-22-gsx-over-jsx-audit.md.

## badge
- WIN: `cva()` variant map replaced by `switch` inside `class={}`.
- GAP: `asChild` not ported (no dynamic tag in gsx — ledger: render-child/Slot).

## card
- Straight port; package-namespaced compound parts (`card.CardHeader`) replace module exports. No divergences.

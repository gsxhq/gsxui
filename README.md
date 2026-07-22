# gsxui

shadcn-style components for [gsx](https://github.com/gsxhq/gsx) — copy-in,
type-checked, server-rendered.

**Status: pre-release.** Foundation in place: theme tokens
(`assets/gsxui.css`), the event-delegation JS runtime (`ui/core/`), and the
first components — badge, button, card, dialog — with pin tests.

- Components live in `ui/<name>/` — a `.gsx` source (JSX-style, named
  parameters, fallthrough attrs) plus a behavior `.js` when interactive.
- `make test` regenerates and tests everything; `make check` adds JS syntax
  and gofmt checks.
- Divergences from the shadcn/ui reference: `docs/jsx-parity.md`.

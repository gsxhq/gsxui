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

**Status: pre-release.** Components: badge, button, card, dialog. The
showcase site, theme editor, and the remaining shadcn set are in progress.

- Components live in `ui/<name>/` — a `.gsx` source (JSX-style, named
  parameters, fallthrough attrs) plus a behavior `.js` when interactive.
- `make test` regenerates and tests everything; `make check` adds JS syntax
  and gofmt checks.
- Divergences from the shadcn/ui reference: `docs/jsx-parity.md`.

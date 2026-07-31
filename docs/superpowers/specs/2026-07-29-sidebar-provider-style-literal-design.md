# SidebarProvider style literal correction

Date: 2026-07-29. Status: approved design, pre-implementation.

## Problem

`SidebarProvider` currently authors its two CSS custom properties as two
separate contributions in a composed style list:

```gsx
style={ css`--sidebar-width:@{sidebarWidth}`, css`--sidebar-width-icon:@{sidebarWidthIcon}` }
```

That is the wrong construct. The declarations form one inline CSS declaration
block, but the source models them as independent mergeable style contributions.
Generated code consequently builds two `gsx.ClassPart` values and relies on
the style merge path to insert punctuation between them.

## Decision

Author the declarations as one composed CSS literal with the declaration
separator written explicitly:

```gsx
style=css`--sidebar-width:@{sidebarWidth};--sidebar-width-icon:@{sidebarWidthIcon}`
```

The literal is one source-level CSS value and generates one style
contribution. The rendered attribute is:

```html
style="--sidebar-width:16rem;--sidebar-width-icon:3rem"
```

This is the only occurrence of a braced, multi-`css` style list in the
repository, so the change is scoped to `SidebarProvider`.

## Generated output

`ui/sidebar.x.go` remains generated source and must not be edited manually.
Run `go tool gsx generate` after changing `ui/sidebar.gsx` and commit the
resulting generated delta with the authored source.

## Testing

Strengthen `TestSidebarProviderCarriesWidthVars` to require the exact rendered
style attribute rather than checking the two variable names independently.
The exact assertion fails against the current merge-inserted `"; "` output
and passes only after the single-literal source is regenerated.

Run the focused Sidebar test first for the red/green cycle, then `make check`
to cover Go tests, Playwright behavior tests, generated-source drift,
JavaScript syntax, and formatting.

## Non-goals

- No GSX parser, code-generator, or runtime change.
- No change to style fallthrough/merge behavior.
- No audit or rewrite of unrelated single-declaration CSS literals.
- No hand edit of generated `.x.go` source.

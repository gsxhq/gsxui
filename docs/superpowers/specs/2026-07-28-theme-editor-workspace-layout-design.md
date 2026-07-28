# Theme editor workspace layout

**Date:** 2026-07-28

**Status:** Approved design

## Purpose

The theme editor is an application workspace, not a documentation article.
Rendering `/theme` through the centered documentation layout currently spends
width on a component-navigation sidebar and caps the remaining editor at
`max-w-6xl`. That makes the controls and live preview feel cramped and pushes
the preview too far down at intermediate viewport widths.

The theme editor must use the available viewport in the same spirit as
shadcn/create, while documentation and component pages retain their current
reading layout.

## Layout architecture

The shared site layout gains an explicit, package-private layout mode:

- **docs** keeps the existing `max-w-6xl` header, sidebar, content container,
  and footer treatment;
- **workspace** keeps the global header and footer, omits the documentation
  sidebar, and allows the main content to span the viewport with responsive
  outer padding.

The mode is explicit at the page call site rather than inferred from the route
or implemented with negative margins inside `ThemeEditor`. `/theme` is the
only workspace consumer in this slice. All other routes remain in docs mode.

The header uses the same workspace-width container so its edges align with the
editor. Navigation and behavior remain identical in both modes.

## Theme editor composition

From the `lg` breakpoint, the editor is a two-pane workspace:

- the controls occupy the narrower pane;
- the live Button preview occupies the wider pane;
- the preview is sticky below the site header; and
- both panes have an explicit minimum width so the layout never creates
  unusably narrow fields.

Below `lg`, the editor becomes one column. Its grid contains three semantic
regions in document order: style picker, preview, and detailed controls. At
`lg`, named grid areas place the style picker and detailed controls together in
the left pane while the preview occupies the right pane. This preserves one
iframe and places it immediately after the style picker on narrow screens.

## Visual direction

This is a spatial correction, not a restyle. Existing type, colors, control
appearance, and editor vocabulary remain unchanged. The workspace should feel
purpose-built through scale and placement:

- full available width;
- quiet responsive gutters (`px-4`, increasing on larger screens);
- a strong, persistent preview surface; and
- no decorative chrome that competes with the component being edited.

## Accessibility and behavior

The document contains one preview iframe at every viewport width. Its title,
status live region, Retry action, manual-copy fallback, handshake, and theme
state behavior remain unchanged. Responsive reordering must preserve logical
keyboard order: style selection, preview, then detailed theme controls.

## Verification

Tests must establish:

1. `/theme` selects workspace mode and does not render the docs sidebar.
2. Existing documentation/component routes retain docs mode and their sidebar.
3. The workspace container is not capped by the docs `max-w-6xl`.
4. Desktop rendering shows a two-pane editor with a sticky preview.
5. Narrow rendering places the single preview after the style picker.
6. Existing theme-editor behavior tests continue to pass.
7. Generated GSX output is regenerated from authored `.gsx` sources and the
   repository's full `make check` gate passes.

## Non-goals

This change does not redesign editor controls, add new theme capabilities,
alter the iframe protocol, change other site pages, or migrate another
component style.

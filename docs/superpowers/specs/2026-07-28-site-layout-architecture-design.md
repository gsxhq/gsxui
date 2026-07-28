# Site layout architecture

**Date:** 2026-07-28

**Status:** Approved design, revised after shadcn source audit

## Purpose

gsxui currently renders marketing, documentation, component examples, and the
theme editor through one centered `max-w-6xl` shell. Those pages have different
jobs:

- documentation needs stable navigation and a readable measure;
- component pages need the same documentation navigation plus live examples;
- the theme editor needs an application workspace with a persistent preview;
  and
- the home page remains a presentation page.

The site must express those differences as explicit layout modes instead of
making individual pages escape a shared container with local CSS.

## Reference architecture

This design follows the current shadcn site and its source in
`/Users/jackieli/personal/shadcn-ui/apps/v4`.

At a 1280px viewport, shadcn's component documentation renders:

- a full-width header;
- a 224px navigation menu inside a 288px left layout column;
- a centered 640px article;
- a 224px table of contents inside a 288px right layout column; and
- responsive disappearance of the right rail at `xl` and the left rail below
  `lg`.

The implementation is split between `app/(app)/docs/layout.tsx`,
`app/(app)/docs/[[...slug]]/page.tsx`, `components/docs-sidebar.tsx`, and
`components/docs-toc.tsx`.

shadcn/create is a separate application shell. It fixes the page to the
viewport below the global header, hides the footer, gives the customizer a
bounded width, and lets the iframe preview consume the remaining space. Its
implementation lives under `app/(app)/(create)`.

gsxui adopts these layout responsibilities, not shadcn's React component
structure or exact decorative styling.

## Page modes

The shared site shell has three package-private modes. The mode is explicit at
each page call site and is never inferred from the request path.

### Marketing

The home page keeps its existing centered presentation and footer. This slice
does not redesign it.

### Documentation

`/docs/**`, `/components`, and `/components/{name}` use a wide documentation
shell:

- the global header aligns to the wide site container;
- the left documentation navigation is sticky and independently scrollable;
- the article column is centered and capped at 640px;
- a sticky right-hand table of contents is shown from `xl`;
- the right rail is omitted when a page has no sections;
- the left rail moves into the existing mobile navigation below `lg`; and
- the footer remains part of the normal document flow.

Component examples remain inside the article column. Wide previews use their
existing isolated iframe boundary rather than escaping across navigation
rails.

### Workspace

`/theme` uses an application shell:

- the global header spans the viewport;
- the documentation sidebar and table of contents are omitted;
- the footer is omitted;
- the area below the header is constrained to the remaining viewport height
  at desktop widths; and
- the customizer and preview manage their own overflow.

This is the equivalent of shadcn/create, not a wider documentation article.

## Layout API

The current exported `Layout` component becomes a package-private site layout
with an explicit `layoutMode`. Its inputs are:

- document title;
- active left-navigation item;
- layout mode;
- optional table-of-contents items; and
- page content.

The layout mode is a named type with `marketing`, `docs`, and `workspace`
constants. A boolean such as `wide` is not used because it would not describe
sidebar, footer, overflow, and header behavior at the call site.

All shared header, navigation, footer, search, theme-toggle, and toaster markup
stays in one implementation. Mode-specific classes and conditional rails are
resolved before the tags so the generated output does not duplicate the site
shell.

## Documentation headings and table of contents

The table of contents is server-rendered from a small package-private model:

```text
docTOCItem
  id
  title
  depth
```

Headings and table-of-contents links share the same item values. Pages do not
maintain one list of labels for the article and another for the rail.

- Component-page items are derived from the existing example registry: the
  example name supplies the stable ID and the example title supplies the
  label.
- Static guide pages define their ordered section items once and render both
  their headings and rail links from those values.
- The components index has no table of contents until it gains real sections.
- IDs are explicit, stable source identifiers. The site does not guess slugs
  from display text.

A shared heading component renders the correct heading level, ID, and anchor
target from a `docTOCItem`. The right rail renders the same ordered slice.

Small site JavaScript observes those already-authored heading IDs and marks the
corresponding rail link active while scrolling. It does not discover headings,
generate IDs, or insert navigation after load. Without JavaScript, every link
remains present and usable; only active-section highlighting is absent.

## Theme editor composition

From `lg`, the theme editor fills the workspace below the header:

- the customizer is a bounded left column for style, palette, radius, and
  transport controls;
- its detailed controls scroll independently;
- the preview fills the remaining width and height; and
- the single preview iframe fills its preview surface instead of relying on a
  fixed 640px minimum height.

Below `lg`, the page returns to normal document scrolling. Its semantic regions
appear in this order:

1. style picker;
2. live preview;
3. the remaining theme and transport controls.

The desktop grid places the style picker and detailed controls in the left
column and the preview in the right column. It never duplicates the iframe,
status region, Retry action, or form controls.

## Visual direction

This refactor changes spatial hierarchy, not the design language. Existing
tokens, typography, control appearance, examples, and copy remain unchanged.
The distinguishing behavior is structural:

- documentation is wide around a deliberately narrow reading column;
- navigation rails remain quiet and stable;
- the workspace spends nearly all available space on the live result; and
- responsive gutters prevent either mode from touching viewport edges.

## Accessibility and behavior

- The table of contents is a labelled navigation landmark.
- Each link targets a unique heading ID.
- The active link uses `aria-current="location"` as well as its visual state.
- Sticky rails have bounded height and their own scrolling without trapping
  keyboard focus.
- The theme document contains exactly one titled preview iframe.
- Responsive layout changes preserve DOM and keyboard order.
- Existing search, theme toggle, preview handshake, Retry, import/export, and
  manual-copy behavior remain unchanged.

## Verification

Server and generated-source tests establish:

1. every page selects its intended layout mode;
2. docs pages render the left rail and workspace pages do not;
3. the workspace omits the site footer;
4. each table-of-contents link has exactly one matching heading ID;
5. component table-of-contents items match the example registry;
6. pages without sections omit the right rail; and
7. committed `.x.go` files match their authored `.gsx` sources.

Browser tests establish:

1. desktop docs show left navigation, a 640px article, and the right rail;
2. the right rail disappears below `xl` and the left rail below `lg`;
3. table-of-contents links navigate to their headings and active state follows
   scrolling;
4. desktop `/theme` occupies the viewport below the header, keeps its
   customizer usable, and gives remaining space to the iframe;
5. narrow `/theme` orders style, preview, then details with one iframe;
6. normal and silent-preview theme behavior still passes; and
7. the repository's full `make check` gate passes.

## Non-goals

This change does not redesign the home page, restyle component examples,
change theme-editor capabilities or iframe messaging, add collapsible desktop
navigation, migrate another component style, or copy shadcn's React-specific
implementation.

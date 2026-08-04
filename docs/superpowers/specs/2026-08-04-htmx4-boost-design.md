# htmx 4 boosted navigation for the gsxui site

Date: 2026-08-04
Status: approved

## Problem

Site navigation is full page loads. Interactive state dies on every click:
scroll the docs sidebar to the bottom, click a component link, and the
sidebar snaps back to the top. Theme and open UI state survive only via
re-initialization tricks.

## Decision

Use htmx 4 (beta) boosted navigation with core idiomorph morph swaps,
**site only** — the shipped component library (`ui/`, vendored files, the
CLI) is untouched and stays htmx-free.

- Dependency: `htmx.org@4.0.0-beta6`, **exact pin** (beta APIs may shift;
  upgrades are deliberate). Installed via npm into the site's existing
  toolchain, imported from `web/main.js` (the site's Vite entry).
- `siteLayout`'s `<body>` gains `hx-boost:inherited="true"` and
  `hx-swap="outerMorph"`. Every internal `<a>` — sidebar, header, TOC,
  in-content — becomes fetch + morph with history push. Layout-mode
  differences (docs/marketing/workspace body classes and `data-site-layout`)
  are attribute morphs on the same body node.
- Why morph solves the problem: the sidebar is the same DOM node across
  the swap, so its scroll offset and any open popovers survive; the
  delegated-listener core (`ui/gsxui.js`) means swapped-in content needs no
  re-initialization by design.
- Title/head: verify during implementation whether beta6's boost updates
  `document.title` from the response. If not, add a small listener on the
  htmx swap event that copies the response `<title>`. No other head content
  differs between pages.
- The ⌘K command palette keeps `window.location.assign` — that lives in
  vendored library code (`ui/command.js`) and stays htmx-free. Full load
  there is acceptable.
- Dev-mode FOUC gate (`html[data-loading]`, dev only) runs on initial load
  only; boosted swaps don't re-trigger it.

## Error handling

htmx falls back per its own semantics: non-2xx responses and cross-origin
links do full navigation; if htmx fails to load, links are plain anchors —
the site degrades to exactly today's behavior.

## Testing

- New Playwright spec (`jstest/specs/site-boost.spec.ts` or folded into
  site-layout.spec.ts per existing structure): scroll the docs sidebar to
  the bottom, click a late component link, assert (a) URL and main content
  changed, (b) sidebar scroll offset preserved, (c) `document.title`
  updated, (d) theme class survives.
- Existing site-layout, theme-editor, and home-showcase specs stay green —
  they are the regression net for boost breaking normal pages.

## Out of scope

- Any change to `ui/` (library JS, command palette navigation).
- Boosting theme-editor workspace forms beyond what body-level boost gives
  (verify the editor still works; exclude it with `hx-boost:inherited="false"`
  on its form if it misbehaves).
- View transitions, prefetching, htmx extensions.

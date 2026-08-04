# Non-Vite init support and docs

Date: 2026-08-04
Status: draft

## Problem

gsxui was designed so users may or may not use Vite, but today `gsxui init`
hard-refuses any project that isn't a pristine `gsx init` npm/Vite scaffold
(`internal/cli/scaffold_integration.go`), and its fallback instructions are
themselves Vite-only. Verified end-to-end (2026-08-04): the vendored component
JS is dependency-free native ESM and works served by a plain Go
`http.FileServer` with one `<script type="module">` tag — dialog behaviors run
with zero console errors, and CSS builds with any Tailwind v4 tool. No UMD
build is needed; the only real gaps are `init`'s refusal and missing docs.

A no-npm project has one genuine blocker: `web/gsxui/index.css` contains
`@import "tw-animate-css"`, an npm package no build tool can resolve without
`node_modules`. (`@import "tailwindcss"` is fine: every Tailwind v4 tool
resolves it — the standalone binary embeds it, `@tailwindcss/cli` depends on
it.)

## Principles

- We give examples and recommend tools, but never manage the user's tooling
  choice. Non-Vite mode runs **no npm commands** and writes **no
  package.json** — we assume npm may not exist.
- Copy-in philosophy: vendored output must be self-contained and buildable.

## Design

### CLI: `gsxui init` non-Vite mode (auto-detected)

Detection, evaluated in `planScaffoldIntegration`:

- **Neither `vite.config.ts` nor `web/main.js` exists** → non-Vite mode.
- Either exists (pristine or modified) → current behavior, except the refusal
  message for a modified scaffold now also includes the non-Vite guidance
  below, so a user who deleted Vite on purpose can see the path out.

In non-Vite mode init:

- contributes no scaffold-integration artifacts (no vite.config.ts /
  web/main.js rewriting, no `scripts.dev` validation);
- runs no npm commands and produces no package.json / package-lock.json
  artifacts (`packageMetadataArtifacts` skipped);
- vendors `tw-animate-css`'s compiled CSS as `web/gsxui/animate.css`
  (a single MIT-licensed Tailwind-v4 source file, committed under `assets/css/`
  so `go:embed` can ship it, copied from
  `node_modules/tw-animate-css/dist/tw-animate.css`), and the vendored `index.css` imports `"./animate.css"` instead of
  `"tw-animate-css"`;
- keeps everything else identical: CSS/JS/merge artifacts, preset handling,
  `go get` of gsx + tailwind-merge-go, `go get -tool` of the gsx CLI, and
  `gsx generate`;
- prints non-Vite next steps instead of the Vite summary:
  1. serve `web/gsxui/` statically and load
     `<script type="module" src=".../index.js">`;
  2. build CSS with the Tailwind tool of your choice, e.g.
     `npx @tailwindcss/cli -i web/gsxui/index.css -o dist.css` or the
     standalone `tailwindcss` binary — examples only, nothing is installed;
  3. link the built stylesheet.

Vite-mode init is unchanged (still installs `tw-animate-css` via npm; its
`index.css` keeps the npm import). `gsxui add` needs no changes — it already
runs only `gsx generate` — but the barrel regeneration must not disturb
`animate.css`.

The embedded `animate.css` copy carries an attribution header and a NOTICE.md
entry; it is refreshed from node_modules by a Makefile target so upgrades are
a dependency bump plus regeneration, and a test asserts the embedded copy
matches the installed package version.

### Docs: "Without Vite" on Getting Started

`site/pages/getting_started.gsx` gains a "Without Vite" subsection (own TOC
entry) after "Manual integration":

- states that `gsxui init` detects the absence of a Vite scaffold and what it
  does differently (no npm, vendored animate.css);
- snippets: a minimal Go static-serving + script-tag example (from the
  verified scenario), and an example CSS build command — worded as examples,
  not managed tooling;
- notes the component JS is dependency-free native ESM: any bundler works,
  none is required;
- the "Manual integration" snippet is updated to quote the CLI's new refusal
  output verbatim, per the page's real-output convention.

## Testing

- CLI tests: non-Vite init in a bare Go module — succeeds; no package.json
  created; no npm invoked; `animate.css` vendored; `index.css` imports
  `./animate.css`; summary text correct. Then `gsxui add dialog` on top
  succeeds.
- Modified-scaffold refusal keeps its test and additionally asserts the
  non-Vite guidance text.
- Embedded-copy-matches-npm-package test as above.
- Existing e2e (Vite path) stays green unchanged.

## Out of scope

- UMD/IIFE builds (no supported browser needs them).
- `<ui.Script for={Component}/>` — the add-time barrel already scopes scripts
  to vendored components; one static tag suffices.
- CSS build orchestration for the user (watch modes, minification, etc.).

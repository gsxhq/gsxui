# `gsx init` Scaffold Integration Design

**Date:** 2026-07-29

## Goal

Make `gsxui init --preset …` leave a fresh `gsx init` project ready to develop
and build without manual Tailwind installation, Vite configuration, or
entry-point imports.

## Proven Failure

The clean sequence

```sh
gsx init app --yes
cd app
gsxui init --preset '<code>'
```

currently copies `web/gsxui/index.css`, whose first imports are `tailwindcss`
and `tw-animate-css`, but does not install either package or configure
`@tailwindcss/vite`.

Following the printed CSS/JS import instructions then fails `vite build` because
`tailwindcss` is missing. Installing the packages without configuring the Vite
plugin makes Vite report success while leaving hundreds of literal `@apply`
rules in the output. This is a functional failure, not optional setup.

## Decision

`gsxui init` owns complete integration for an exact recognized `gsx init`
scaffold:

1. install `tailwindcss`, `@tailwindcss/vite`, and `tw-animate-css` as development
   dependencies using npm;
2. add the `@tailwindcss/vite` import and `tailwindcss()` plugin to
   `vite.config.ts`;
3. import the configured gsxui JS barrel and CSS entry from `web/main.js`;
4. copy the existing gsxui artifacts and configure the Go class merger;
5. generate GSX output;
6. print a completed summary whose only component next step is
   `gsxui add button`.

Rerunning init is idempotent.

## Recognition, Not Heuristics

The CLI does not parse or rewrite arbitrary JavaScript/TypeScript with string
searches. It maintains exact, versioned before/after scaffold documents for the
`gsx init` files it supports:

- `vite.config.ts`;
- `web/main.js`;

Each recognized source has:

- an exact SHA-256 fingerprint;
- a complete integrated replacement;
- an exact fingerprint for the already-integrated form.

Each file independently accepts only these transitions:

- recognized pristine scaffold → recognized integrated scaffold;
- recognized integrated scaffold → no-op.

This permits recovery from a partially followed manual setup when each file is
still one of the two known documents. It never accepts a third, approximately
matching form.

`package.json` is different because the project name is generated and npm owns
dependency metadata. The CLI parses it with Go's JSON decoder and requires the
known `gsx init` scripts and `@gsxhq/vite-plugin-gsx` dependency. It also
requires `package-lock.json` and rejects other package-manager lockfiles.

Formatting changes, user edits, different paths, a different package manager,
or an unrecognized Vite config cause preflight refusal before gsxui artifacts or
dependency commands run. The error lists the unresolved responsibilities and
the exact manual edits. `--overwrite` never authorizes overwriting these
consumer-owned integration files.

This narrow contract is deliberate: a real syntax-aware TypeScript rewriter can
be added later, but an approximate parser is not acceptable.

## Package Installation

The supported scaffold is npm-based and must contain `package.json` plus
`package-lock.json`, as produced by `gsx init --yes`.

The CLI runs:

```sh
npm install --save-dev tailwindcss@^4.3.3 @tailwindcss/vite@^4.3.3 tw-animate-css@^1.4.0
```

The version ranges are embedded release inputs and change only through a tested
gsxui release. npm owns `package.json` dependency insertion, version resolution,
the lockfile, and ensuring `node_modules` is populated. The command runs on
every init; with matching metadata and lockfile it is an idempotent npm no-op.
The CLI does not hand-edit `package-lock.json`.

Go dependency installation retains the existing commands.

## Planning and Failure Boundaries

All preset decoding, managed-artifact conflict checks, scaffold recognition,
and replacement planning finish before any external command or write.

The integration files join the existing artifact transaction but remain
consumer-owned: their original bytes are rollback snapshots, not managed files
recorded in `gsxui.json`.

Package installation runs only after preflight succeeds. If npm fails, the CLI
restores `package.json` and `package-lock.json` from snapshots and removes a
newly created lockfile. It does not attempt to roll back npm's cache. A newly
installed but now-unreferenced directory under `node_modules` is harmless and
may remain; authored project state is restored.

If a later Go dependency, artifact, or generation step fails, the same snapshots
restore npm metadata and the artifact transaction restores integration and
managed files.

Interrupted transactions use the existing recovery journal before a new init
begins.

## Output

Successful recognized integration prints:

```text
gsxui initialized.
  css:  web/gsxui/index.css
  js:   web/gsxui/index.js
  vite: Tailwind CSS configured
  next: gsxui add button
```

It no longer tells users to import files that were already imported.

Unknown projects receive no success summary. The refusal explains that automatic
integration currently supports an unmodified `gsx init` npm/Vite scaffold and
provides manual dependency, Vite plugin, and entry-import instructions.

## End-to-End Gate

An automated integration test builds the pinned `github.com/gsxhq/gsx/cmd/gsx`
tool from this module's `go.mod`, scaffolds a new project, runs `gsxui init` with
a compact Maia preset, adds Button, and performs a real npm production build.

The gate asserts:

- dependencies exist in `devDependencies` and the lockfile;
- `vite.config.ts` contains one Tailwind import and one plugin registration;
- `web/main.js` contains one gsxui JS import and one CSS import;
- rerunning `gsxui init` changes no integration file and duplicates nothing;
- generated CSS contains representative utilities used by the copied Button;
- generated CSS contains no `@apply`, `@theme`, or unresolved Tailwind directive;
- `go build ./...` succeeds.

Fast unit tests use the command seam and exact scaffold fixtures to cover
preflight refusal, rollback, idempotency, and output without installing npm
packages.

## Documentation

Getting Started describes the zero-manual-integration path for a fresh
`gsx init` app. A separate manual-integration section remains for unsupported
projects and exactly matches the refusal output.

## Rejected Alternatives

**Keep printing setup instructions.** Rejected because the normal scaffold path
is deterministic and the current instructions allow a silently uncompiled
build.

**Use substring insertion or regular expressions for TypeScript.** Rejected as
an unsafe heuristic that can duplicate plugins or corrupt valid custom config.

**Add Tailwind to every `gsx init` project.** Rejected because GSX itself does
not require Tailwind; the dependency belongs to gsxui initialization.

**Silently skip unknown files.** Rejected because init would again report
success for an unstyled project.

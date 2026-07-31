# Compact Presets and `gsx init` Scaffold Integration Plan

**Goal:** Ship the approved compact built-in preset transport and make
`gsxui init` completely configure an exact recognized `gsx init --yes`
npm/Vite project.

**Branch:** `recipe-model`

**Design sources:**

- `docs/superpowers/specs/2026-07-29-compact-preset-transport-design.md`
- `docs/superpowers/specs/2026-07-29-gsx-scaffold-integration-design.md`

**Implementation constraints:**

- Use red/green TDD for every behavioral change.
- Keep `gsxui:v1:` fully compatible and lossless for custom presets.
- Treat compact catalogue arrays as an explicit append-only transport ABI.
- Recognize JavaScript and TypeScript scaffold files by exact document
  fingerprints and complete replacements only; do not use substring or regex
  rewriting.
- Finish all preflight checks before writes or external commands.
- Do not touch the user's development server on port 5173. Browser tests use
  the repository harness on port 7799 and must leave it clean.

## Task 1: Pin the compact Go transport contract

**Files:**

- Modify: `internal/preset/transport_test.go`
- Create: `internal/preset/compact_test.go`
- Create: `internal/preset/compact.go`
- Modify: `internal/preset/transport.go`

**Steps:**

1. Add failing tests for exact default and boundary `p1` codes, every valid
   style/base/theme/radius combination, custom theme/radius `v1` fallback,
   unsupported versions, invalid characters, overflow, unused indexes,
   invalid base/theme combinations, and non-canonical base62 spellings.
2. Run `go test ./internal/preset -run 'Test(Compact|Share|InputResolver)' -count=1`
   and confirm the tests fail because `p1` is not implemented.
3. Add immutable style, base-color, theme, and radius ABI arrays; a strict
   base62 codec; bit packing/unpacking; and catalogue resolution.
4. Dispatch `EncodeShare` by exact `MatchPalette` result and dispatch
   `DecodeShare` by transport version while leaving the current `v1` decoder
   unchanged.
5. Rerun the focused tests and `gopls check -severity=hint` on changed Go
   files.

## Task 2: Pin input resolver compatibility

**Files:**

- Modify: `internal/preset/transport_test.go`

**Steps:**

1. Add failing cases proving a compact code resolves directly, through stdin,
   from a file, and from an HTTPS response.
2. Run the focused resolver tests and confirm the missing compact support is
   the failure.
3. Reuse `DecodeShare` through the existing resolver paths without changing
   their public syntax.
4. Rerun the resolver tests.

## Task 3: Add browser transport parity and editor schema

**Files:**

- Modify: `site/pages/theme.gsx`
- Regenerate: `site/pages/theme.x.go`
- Modify: `site/pages/theme_schema_test.go`
- Modify: `web/theme-state.js`
- Modify: `web/theme-state.test.js`
- Modify: `jstest/specs/theme-editor.spec.ts`

**Steps:**

1. Add failing schema tests for both prefixes and the ordered compact ABI.
2. Add failing Node tests for exact Go/JavaScript code parity across every
   valid catalogue combination, strict compact decoding, custom `v1`
   fallback, old full-code loading, and transient-preview invariance.
3. Add failing Playwright assertions that built-in selections produce compact
   share URLs/commands while custom imports remain full.
4. Run the focused Go, Node, and theme-editor Playwright tests and observe the
   expected failures.
5. Export the compact ABI into the server-authored schema, implement strict
   synchronous compact encode/decode in `theme-state.js`, and keep full
   canonical JSON export unchanged.
6. Regenerate authored GSX with `go tool gsx generate`.
7. Rerun all focused tests and confirm port 7799 is free afterward.

## Task 4: Pin exact scaffold recognition and planning

**Files:**

- Create: `internal/cli/scaffold_integration.go`
- Create: `internal/cli/scaffold/gsx-v1/vite.config.ts`
- Create: `internal/cli/scaffold/gsx-v1/vite.gsxui.config.ts`
- Create: `internal/cli/scaffold/gsx-v1/main.js`
- Create: `internal/cli/scaffold/gsx-v1/main.gsxui.js`
- Create: `internal/cli/scaffold_integration_test.go`
- Modify: `internal/cli/init_test.go`

**Steps:**

1. Add exact pristine and integrated fixtures copied from the pinned GSX
   scaffold, plus unit tests for pristine, integrated, partially integrated,
   custom-file refusal, package contract refusal, missing npm lockfile, and
   alternate lockfile refusal.
2. Run the focused CLI tests and confirm no scaffold planner exists.
3. Embed the versioned full documents, fingerprint them with SHA-256, validate
   `package.json` structurally, and return consumer-owned artifacts for only
   the pristine files.
4. Make `initTestModule` create a recognized scaffold so existing init/add/apply
   tests continue to exercise the supported path.
5. Rerun focused CLI tests.

## Task 5: Integrate npm and rollback into `gsxui init`

**Files:**

- Modify: `internal/cli/init.go`
- Modify: `internal/cli/run.go` if a narrower command seam is required
- Create or modify: `internal/cli/package_snapshot.go`
- Modify: `internal/cli/init_test.go`
- Modify: `internal/cli/scaffold_integration_test.go`

**Steps:**

1. Add failing tests proving preflight runs before commands/writes; npm runs
   before Go dependency commands; successful init writes exact integration
   files once; rerun is idempotent; npm failure restores package metadata; and
   later command/generation failure restores package metadata and all
   transactional files.
2. Run the focused tests and confirm the rollback/output expectations fail.
3. Add package metadata snapshots, run the pinned npm install command on every
   supported init, include scaffold integration files in the existing artifact
   transaction as unmanaged artifacts, and restore package metadata on every
   post-npm failure.
4. Change successful output to:

   ```text
   gsxui initialized.
     css:  web/gsxui/index.css
     js:   web/gsxui/index.js
     vite: Tailwind CSS configured
     next: gsxui add button
   ```

5. Rerun focused CLI tests and inspect `git diff --check`.

## Task 6: Add the real fresh-project gate

**Files:**

- Modify: `internal/cli/e2e_test.go`

**Steps:**

1. Add an end-to-end test that builds the pinned GSX tool from this module's
   dependency graph, runs `gsx init app --yes`, runs this checkout's
   `gsxui init` with a compact Maia preset, adds Button, reruns init, builds
   Vite, and builds Go.
2. Assert exact dependency metadata, one Tailwind import/plugin, one JS/CSS
   entry import, byte-stable integration files on rerun, representative Button
   utilities in emitted CSS, and absence of unresolved Tailwind directives.
3. Run the new E2E test without `-short` and retain its output as the primary
   proof of the user workflow.

## Task 7: Update user documentation and generated source

**Files:**

- Modify: `site/pages/getting_started.gsx`
- Modify: `site/snippets/init.output.txt`
- Regenerate: `site/pages/getting_started.x.go`
- Regenerate: `site/hl/blocks.gen.go`
- Modify other current docs only where they describe manual setup as the
  normal fresh-scaffold path.

**Steps:**

1. Change Getting Started to begin with `gsx init`, describe zero-manual setup
   for the recognized npm/Vite scaffold, and add a clearly separate exact
   manual-integration section for unsupported projects.
2. Keep refusal output and manual documentation byte-for-byte aligned through
   a focused test or shared constant.
3. Regenerate highlighted snippets and GSX output using repository generators.
4. Run documentation/page tests.

## Task 8: Authoritative verification and review

**Steps:**

1. Run focused Go tests for `internal/preset`, `internal/cli`, and `site/pages`.
2. Run `node --test web/theme-state.test.js`.
3. Run the focused theme-editor Playwright suite.
4. Run `make check`.
5. Run `make ci` as the uncached authoritative gate.
6. Confirm `site/dist/.gitkeep` exists, port 7799 is free, and no process
   started by this work remains.
7. Inspect `git status --short`, `git diff --check`, generated-source checks,
   and the complete diff against `4dc6cec`.
8. Commit the compact transport and scaffold integration in reviewable
   checkpoints without pushing or merging unless requested.

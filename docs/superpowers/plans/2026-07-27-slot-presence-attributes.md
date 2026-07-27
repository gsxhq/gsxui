# Slot Presence Attributes Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use
> superpowers:subagent-driven-development (recommended) or
> superpowers:executing-plans to implement this plan task-by-task. Steps use
> checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace packed `data-gsxui-slot` token lists with composable
`data-gsxui-slot-<name>` presence attributes and remove all private slot
merging code from source vendored by `gsxui add`.

**Architecture:** Every component forwards caller attributes and then emits
its own invariant slot marker as a distinct attribute. Composed components add
more distinct markers, which GSX's existing scalar attribute merge preserves
without custom parsing. CSS, contract validation, runtime manifests, examples,
and tests consume the same prefix-based public contract.

**Tech Stack:** Go 1.26.1, GSX, Tailwind CSS 4.3.3, PostCSS, vanilla ESM,
Playwright 1.62, existing typed style contract and generated-source checks.

## Global Constraints

- Work only in
  `/Users/jackieli/personal/gsxhq/gsxui/.worktrees/css-only-theme-architecture`
  on `codex/css-only-theme-architecture`.
- The public marker format is exactly `data-gsxui-slot-<kebab-case-name>`.
- Markers are valueless. Rendered markers with a value and CSS selectors
  using a value/token operator are contract violations.
- Component-owned markers are authored after the fallthrough spread so a
  caller cannot remove the component's styling identity.
- Composition uses distinct marker keys and GSX's normal attribute forwarding;
  do not add a replacement Go helper, core API, token parser, or compatibility
  shim.
- CSS uses presence selectors such as `[data-gsxui-slot-button]`; packed
  `[data-gsxui-slot]`, `[data-gsxui-slot="..."]`, and
  `[data-gsxui-slot~="..."]` are forbidden in current source and compiled CSS.
- `class` remains caller-owned and preserves the existing wrapper/control
  routing of `NativeSelect`.
- `NativeSelect` must not grow a slot-prefix filter. Calendar dropdown styling
  uses the documented `calendar-dropdowns` ancestor plus
  `native-select-wrapper`; the redundant `calendar-dropdown-root` slot is
  removed from the typed contract.
- Runtime JavaScript may create a slot marker with
  `setAttribute("data-gsxui-slot-<name>", "")`; it must not inspect styling
  markers for behavior.
- `gsxui add` vendors no slot support file or nested private slot package.
- The embedded file manifest and generated imports contain no dependency edge
  to the deleted helpers. Automatic deletion from an already-vendored external
  project is out of scope because it requires a separate destructive migration
  policy.
- Update authored `.gsx`, CSS, JavaScript, current docs, examples, tests, and
  committed generated output together. Historical implementation plans remain
  historical and are not rewritten.
- Use a real red-green-refactor cycle. Record the expected failure before each
  production change.
- Do not hand-edit generated `.x.go`, highlighted blocks, or runtime manifest
  output; run their owning generators.
- Run `gopls check -severity=hint` on every authored Go file changed.
- The authoritative final gate is `make ci`.

---

### Task 1: Make contract boundaries require presence attributes

**Files:**

- Modify: `site/examples/style_contract_test.go`
- Modify: `jstest/support/compiled-css-audit.ts`
- Modify: `jstest/support/compiled-css-audit.test.ts`
- Modify: `Makefile`

**Interfaces:**

- Produces: `slotAttribute(name string) string`, returning
  `"data-gsxui-slot-"+name`.
- Produces: contract extraction that treats every rendered attribute with the
  `data-gsxui-slot-` prefix as one declared slot and rejects non-empty values.
- Produces: runtime contract entries whose slot-presence row has
  `Attribute == slotAttribute(slot.Name)` and `Value == ""`.
- Produces: compiled CSS rejection of the obsolete packed gsxui slot selector
  in addition to generic `data-slot`.

- [ ] **Step 1: Change the rendered-contract test first.**

  In `TestRegisteredExamplesCoverStyleContract`, stop reading one
  `data-gsxui-slot` value. Collect every element attribute whose key starts
  with `data-gsxui-slot-`, trim the prefix, reject an empty suffix or non-empty
  value, and run the existing declaration/axis checks for the resulting slot
  slice.

  Change `TestRuntimeStyleContractManifestMatchesTypedContract` so runtime slot
  entries use:

  ```go
  Attribute: slotAttribute(slot.Name),
  Value:     "",
  ```

- [ ] **Step 2: Make the CSS audit test reject packed gsxui slots.**

  Extend the existing selector test with:

  ```css
  [data-gsxui-slot~="button"] { display: inline-flex; }
  [data-gsxui-slot-button] { display: inline-flex; }
  ```

  Require one new violation labeled `obsolete packed [data-gsxui-slot]
  selector`, while the presence selector remains accepted.

- [ ] **Step 3: Run both focused tests and record RED.**

  Run:

  ```bash
  go test ./site/examples -run 'TestRegisteredExamplesCoverStyleContract|TestRuntimeStyleContractManifestMatchesTypedContract' -count=1
  node --test jstest/support/compiled-css-audit.test.ts
  ```

  Expected: the Go test reports all declared slots as un-emitted because
  components still emit the packed scalar; the Node test reports no obsolete
  packed-slot violation.

- [ ] **Step 4: Implement the boundary checks.**

  Add `slotAttribute`, prefix extraction, and the new PostCSS selector audit.
  Extend `make audit` to reject:

  - `withSlot(` and `slotattr.With(` in authored component source;
  - bare `data-gsxui-slot` attributes in current `.gsx`, `.go`, `.js`, and
    CSS source;
  - packed gsxui slot selectors.

  Match the attribute/selector grammar, not arbitrary prose containing the
  migration name.

- [ ] **Step 5: Re-run the focused boundary tests.**

  The Node audit must pass. The rendered-contract test must remain red until
  Task 2 migrates production output; confirm its failure is now exclusively
  the missing presence markers.

- [ ] **Step 6: Commit the red contract boundary.**

  ```bash
  git add Makefile site/examples/style_contract_test.go \
    jstest/support/compiled-css-audit.ts \
    jstest/support/compiled-css-audit.test.ts
  git commit -m "test: require composable slot attributes"
  ```

### Task 2: Migrate component output and style selectors

**Files:**

- Modify: all `ui/*.gsx` files currently containing `withSlot(`
- Modify: `ui/icon/icon.gsx`
- Modify: `ui/input-otp.js`
- Modify: `assets/css/foundation.css`
- Modify: `assets/css/styles/default.css`
- Modify: `internal/stylecontract/contracts_composites.go`
- Modify: affected `ui/*_test.go` render and CSS-contract tests
- Modify: `site/pages/accordion_css_test.go`
- Modify: `site/pages/slider_css_test.go`
- Modify: `jstest/specs/runtime-style-contract.spec.ts`
- Modify: affected `jstest/specs/*-style-contract.spec.ts`
- Modify: `jstest/harness/style_contract.gsx`
- Modify: `jstest/harness/sidebar_contract.gsx`
- Modify: `jstest/support/selector-allowlist.ts` only if exact selector strings
  change
- Delete: `ui/slots.go`
- Delete: `ui/slots_test.go`
- Delete: `ui/internal/slotattr/slotattr.go`
- Delete: `ui/internal/slotattr/slotattr_test.go`

**Interfaces:**

- Consumes: `slotAttribute(name)` and prefix extraction from Task 1.
- Produces: every rendered styling role as an empty presence attribute.
- Produces: CSS selectors using `[data-gsxui-slot-<name>]`.
- Produces: runtime browser observation using exact presence-attribute
  selectors and `toHaveAttribute(attribute, "")`.

- [ ] **Step 1: Add focused composition expectations before migration.**

  Update the Button, AlertDialog, Sheet, Calendar, Pagination, Sidebar,
  NativeSelect, and Icon render tests to require separate markers. Examples:

  ```html
  data-gsxui-slot-alert-dialog-action data-gsxui-slot-button
  ```

  and:

  ```html
  data-gsxui-slot-calendar-previous
  data-gsxui-slot-calendar-nav-button
  data-gsxui-slot-button
  ```

  Assert marker presence independently where attribute ordering is not part of
  behavior.

  Add a collision case that passes a non-empty value for Button's own marker
  plus a distinct composed marker. Require the rendered Button marker to be
  valueless and the distinct marker to remain present. Add co-location checks
  proving the composed and primitive markers occur on the same opening tag,
  not merely somewhere in the same document.

- [ ] **Step 2: Run the focused render tests and record RED.**

  Run:

  ```bash
  go test ./ui -run \
    'Test(Button|AlertDialog|Sheet|Calendar|Pagination|Sidebar|NativeSelect|Icon)' \
    -count=1
  go test ./ui/icon -count=1
  ```

  Expected: failures show packed `data-gsxui-slot` output and the current
  helper-based composition.

- [ ] **Step 3: Migrate GSX component markers mechanically.**

  Apply these exact transformations, then run `go tool gsx fmt`:

  ```gsx
  { withSlot("button", attrs)... }
  ```

  becomes:

  ```gsx
  { attrs... }
  data-gsxui-slot-button
  ```

  `nil` emits only the marker. Nested calls forward the innermost attribute
  bag once and emit one marker for every literal role. Conditional roles are
  expressed as conditional GSX attributes rather than mutating an attribute
  bag.

  In `NavigationMenuLink`, replace the conditional `withSlot` assignment with
  a conditional `data-gsxui-slot-navigation-menu-trigger` attribute on the
  `<a>`.

  In `NativeSelect`, keep only caller `class` on the wrapper and forward the
  remaining attributes to `<select>`. Remove all handling of styling markers.
  Replace Calendar's `calendar-dropdown-root` selector with the relationship:

  ```css
  [data-gsxui-slot-calendar-dropdowns]
    [data-gsxui-slot-native-select-wrapper]
  ```

  Remove `calendar-dropdown-root` from the typed style contract.

  In `ui/icon/icon.gsx`, remove the `slotattr` import and emit
  `data-gsxui-slot-icon` after `{ attrs... }`.

- [ ] **Step 4: Migrate JavaScript and CSS.**

  Change InputOTP-created markers to exact empty attributes. Convert every
  foundation/default-style selector from token membership to presence
  selection without changing declarations, combinators, pseudo-classes,
  state axes, or cascade layers.

- [ ] **Step 5: Update the contract/runtime consumers and all render pins.**

  Use exact presence selectors in Playwright. Make `observe` treat an empty
  expected value as a real exact presence-attribute check, not token
  membership. Update Go render expectations and CSS-contract tests to inspect
  marker names.

- [ ] **Step 6: Regenerate and run the component gates.**

  Run:

  ```bash
  go tool gsx fmt ./ui ./jstest/harness
  make generate
  UPDATE_RUNTIME_STYLE_CONTRACT=1 go test ./site/examples \
    -run TestRuntimeStyleContractManifestMatchesTypedContract -count=1
  go test ./ui ./ui/icon ./internal/stylecontract ./site/examples -count=1
  node --test jstest/support/compiled-css-audit.test.ts
  npx playwright test --config jstest/playwright.config.ts \
    jstest/specs/runtime-style-contract.spec.ts \
    jstest/specs/foundation-contract.spec.ts \
    jstest/specs/menus-style-contract.spec.ts \
    jstest/specs/composites-style-contract.spec.ts \
    jstest/specs/sidebar-style-contract.spec.ts
  ```

  Expected: all focused gates pass and `make audit` finds no obsolete helper
  or packed selector in current source.

- [ ] **Step 7: Run diagnostics and commit.**

  Run `gopls check -severity=hint` on
  `internal/stylecontract/contracts_composites.go`,
  `site/examples/style_contract_test.go`, every changed authored UI test, and
  every changed authored page test.

  ```bash
  git add ui assets/css internal/stylecontract site/examples \
    site/pages jstest Makefile
  git commit -m "refactor: compose styling roles as attributes"
  ```

### Task 3: Update the public styling contract and site artifacts

**Files:**

- Modify: `README.md`
- Modify: `NOTICE.md`
- Modify: `docs/theme-system-roadmap.md`
- Modify: `docs/jsx-parity.md`
- Modify:
  `docs/superpowers/specs/2026-07-26-css-only-theme-system-design.md`
- Modify: `site/pages/theming.gsx`
- Modify: `site/snippets/theme-slot.css.txt`
- Modify: `site/pages/home.gsx`
- Modify: `site/pages/layout.gsx`
- Modify: `site/pages/pages_test.go`
- Modify: committed generated `.x.go` and `site/hl/blocks.gen.go` only through
  their owning generators

**Interfaces:**

- Consumes: component, registered-example, helper-deletion, and CLI vendoring
  work completed together in Task 2 because they form one green repository
  boundary.
- Produces: public theming guidance and site-owned trigger markup using exact
  presence attributes.
- Produces: regenerated page/example source and highlighted blocks.

- [ ] **Step 1: Record the remaining public-surface RED.**

  Run:

  ```bash
  make audit
  go test ./site/pages -count=1
  ```

  Expected: `make audit` reports the packed trigger markers in
  `site/pages/home.gsx` and `site/pages/layout.gsx`; page tests still expect
  the packed rendered contract.

- [ ] **Step 2: Migrate current documentation and site chrome.**

  Explain that composition is ordinary attribute forwarding:

  ```gsx
  <ui.Button data-gsxui-dialog-trigger data-gsxui-slot-dialog-trigger>
    Open
  </ui.Button>
  ```

  and:

  ```css
  [data-gsxui-slot-button] { ... }
  ```

  Update current public docs, migration instructions, and the two site-owned
  trigger call sites. Update page tests to require exact empty presence
  attributes. Mark the architecture spec's presence-marker revision
  implemented while keeping the configurator roadmap unscheduled. Keep
  historical implementation plans unchanged.

- [ ] **Step 3: Regenerate documentation artifacts.**

  Use the exact reviewed GSX compiler commit `5d3cbfb9` through an untracked
  Go workspace, then run:

  ```bash
  go run github.com/gsxhq/gsx/cmd/gsx generate
  make highlight
  ```

  Restore `site/dist/.gitkeep` if the build tooling removes it.

- [ ] **Step 4: Run focused public-surface gates and diagnostics.**

  Run:

  ```bash
  go test ./site/pages ./site/examples -count=1
  gopls check -severity=hint site/pages/pages_test.go
  make audit
  git diff --check
  ```

- [ ] **Step 5: Commit.**

  ```bash
  git add README.md NOTICE.md docs site
  git commit -m "docs: document composable slot markers"
  ```

### Task 4: Verify the whole migration

**Files:**

- Modify: only files required by failures that demonstrate a real migration
  gap

**Interfaces:**

- Consumes: all presence-marker output, CSS, tests, docs, and clean vendoring.
- Produces: a clean generated tree and authoritative CI evidence.
- Consumes: reviewed GSX compiler fix
  `5d3cbfb9129266405e4b390ce4f67e75c51ee9ee` through an untracked temporary
  Go workspace until a published version can be pinned.

- [ ] **Step 1: Prove no current packed contract remains.**

  Run:

  ```bash
  make audit
  rg -n 'withSlot|slotattr\.With|data-gsxui-slot(?:[=~\]])' \
    ui assets/css internal site jstest README.md NOTICE.md \
    -g '!*.x.go' -g '!blocks.gen.go'
  ```

  The second command may find prose in historical plan/spec context only when
  those paths are intentionally included; current source must have no match.

- [ ] **Step 2: Verify generated output is current.**

  With the exact reviewed GSX compiler selected through the untracked
  workspace, run:

  ```bash
  go run github.com/gsxhq/gsx/cmd/gsx generate
  make verify-generated
  git diff --exit-code -- '*.x.go' jstest/runtime-style-contract.json
  ! rg 'data-gsxui-slot-[a-z0-9-]+", Value: true' \
    ui site/examples site/pages jstest/harness -g '*.x.go'
  ```

- [ ] **Step 3: Run the authoritative gate.**

  Run:

  ```bash
  GOWORK=<temporary-workspace> make ci
  ```

  Expected: Go tests, generated checks, CSS audit, every discovered Chromium
  behavior test, JavaScript syntax checks, and formatting all pass. Record the
  actual Chromium count in the verification report rather than treating it as
  a fixed invariant.

  Record separately that external CI remains blocked until the reviewed GSX
  commit is published and `go.mod` pins that version; a tracked local replace
  is forbidden.

- [ ] **Step 4: Inspect representative rendered composition.**

  Against the live development site or Playwright harness, inspect Button,
  AlertDialog action/trigger, Calendar navigation/dropdowns, Sidebar menu
  button, and Icon. Each composed DOM node must carry every expected presence
  marker, computed styling must remain applied, and no helper source may be
  present in a fresh `gsxui add` consumer.

- [ ] **Step 5: Commit any verification-only corrections.**

  ```bash
  git add -A
  git commit -m "test: close slot marker migration gaps"
  ```

  Skip this commit when verification produces no corrections.

## Self-review

- Spec coverage: marker naming, forced ownership, composition, CSS,
  JavaScript-created nodes, typed contract, runtime manifest, examples,
  vendoring cleanup, generated output, and full verification are assigned.
- Placeholder scan: no deferred implementation placeholders remain.
- Type consistency: `slotAttribute(string) string` is the only new helper and
  is test-only contract code; production components use plain GSX attributes.
- The prior Phase 1 implementation plan remains historical evidence. The
  current architecture spec and this plan are the source of truth for the
  revised marker contract.

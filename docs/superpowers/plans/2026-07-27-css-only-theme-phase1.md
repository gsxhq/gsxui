# CSS-only theme system Phase 1 implementation plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use
> superpowers:subagent-driven-development (recommended) or
> superpowers:executing-plans to implement this plan task-by-task. Steps use
> checkbox (`- [ ]`) syntax for tracking.

**Goal:** Move every gsxui component's built-in presentation out of `.gsx`
and JavaScript into one default CSS style pack, leaving one canonical
semantic/behavior template and a caller-utility override seam.

**Architecture:** Components emit tokenized `data-gsxui-slot` styling hooks,
explicit state attributes, dedicated behavior hooks, dynamic mechanism
values, and caller attributes only. Tailwind v4 compiles a shared foundation,
default theme, and `@layer components` default style; the site, browser
harness, and CLI consume the same authored assets. A typed Go contract
records all components, slot tokens, and public presentation axes.

**Tech stack:** Go 1.26.1, GSX, Tailwind CSS 4.3.3, vanilla ESM, Playwright
1.62, `gopls`, existing gsxui registry/example harness.

## Global constraints

- Work only in the existing
  `/Users/jackieli/personal/gsxhq/gsxui/.worktrees/css-only-theme-architecture`
  worktree on `codex/css-only-theme-architecture`.
- Runtime component code remains standard Go/GSX and vanilla browser APIs;
  do not introduce a framework or a runtime CSS dependency.
- Keep exactly one `.gsx` implementation per semantic component. There are
  no style-specific templates and no second named style.
- Replace generic `data-slot` with space-separated
  `data-gsxui-slot` tokens. CSS selects tokens with `~=`, never scalar `=`.
- Composed components retain inner-to-outer tokens. For example,
  `AlertDialogAction` renders
  `data-gsxui-slot="button alert-dialog-action"`; it must not duplicate
  Button declarations in the alert-dialog rules.
- Components and component JavaScript contain no library-owned presentation
  class strings after migration. A rendered `class` attribute is
  caller-owned only.
- Do not preserve Tailwind `group/*`, `peer/*`, `data-slot`, or styling-only
  class markers in templates. Translate them to explicit relational CSS
  selectors over slot/state attributes.
- JavaScript may set reflected state, ARIA, text, dynamic custom-property
  values, transforms, and geometry needed by live behavior. It must not set
  presentation classes or query `data-gsxui-slot`; add dedicated
  `data-gsxui-<component>-<part>` behavior hooks where a query is required.
- Static utility declarations live in
  `assets/css/styles/default.css` inside `@layer components`. Rules use
  low-specificity selectors and no `!important`.
- Interaction/accessibility mechanics that every coherent style needs live
  in `assets/css/foundation.css`. Presentation remains in the default style.
- Caller Tailwind utilities must override the style pack through cascade
  layer order, including responsive and state variants.
- Component colors use registered semantic tokens, `currentColor`, or
  transparent. Replace hard-coded overlay/status/contrast palette utilities
  with semantic tokens while preserving their current computed colors.
- Add the standard eight shadcn sidebar tokens because the current
  `bg-sidebar`, `text-sidebar-foreground`, `border-sidebar-border`, and
  `ring-sidebar-ring` utilities do not compile at all without them. This is
  an intentional correction to previously inert styling, not a claimed
  pixel-equivalent port.
- Keep dynamic inline values only where behavior computes them:
  AspectRatio's ratio, Progress's transform, Slider's `--fill`,
  ResizablePanel's flex value, ToggleGroup's `--gap`, Sidebar widths and
  skeleton width.
- Use a real red-green-refactor cycle for every production change. Name the
  behavior each test protects and record the expected red failure in the
  task report.
- Do not hand-edit generated `.x.go`, icon data, highlighted blocks, or
  screenshot output. Run the owning generator/update command.
- Run `gopls check -severity=hint` on every authored Go file changed by a
  task.
- Historical plans/source maps may continue to describe upstream
  `data-slot`; current source, current docs, examples, NOTICE, tests, and
  generated output must use the new contract.
- Each task commits only after its focused Go and browser gates pass. The
  authoritative final gate is `make ci` if the Makefile adds it during this
  work; otherwise `make check` plus `make site`.

## File structure

The finished source tree has these responsibilities:

```text
assets/css/index.css                 consumer Tailwind entry/import order
assets/css/foundation.css            invariant mechanics and token mapping
assets/css/themes/default.css        default light/dark semantic values
assets/css/styles/default.css        current Nova presentation
internal/stylecontract/contract.go   contract types, ordering, validation
internal/stylecontract/*.go          per-family component/slot/axis data
ui/slots.go                          tokenized slot composition helper
ui/*.gsx                             semantic markup and reflected state
ui/*.js                              behavior only
web/site.css                         site-only styling over shared assets
```

The CLI maps the source assets to:

```text
web/gsxui/index.css
web/gsxui/foundation.css
web/gsxui/theme.css
web/gsxui/style.css
```

`Config.CSS` points to the entry file. The other three files are fixed
siblings of that entry, even when the user configures a non-default path.

---

### Task 1: Split the authored and vendored CSS assets

**Files:**

- Create: `assets/css/index.css`
- Create: `assets/css/foundation.css`
- Create: `assets/css/themes/default.css`
- Create: `assets/css/styles/default.css`
- Modify: `assets/gsxui.css` (delete after the split)
- Modify: `web/site.css`
- Modify: `internal/cli/config.go`
- Modify: `internal/cli/init.go`
- Modify: `internal/cli/init_test.go`
- Modify: `site/pages/theme_defaults_test.go`
- Modify: `site/pages/accordion_css_test.go`
- Test: `internal/cli/init_test.go`
- Test: `site/pages/theme_defaults_test.go`
- Test: `site/pages/accordion_css_test.go`

**Interfaces:**

- Produces: `Config.CSS == "web/gsxui/index.css"` by default.
- Produces: `cssAssetTargets(entry string) []cssAssetTarget`, mapping
  `assets/css/index.css`, `foundation.css`, `themes/default.css`, and
  `styles/default.css` to the entry and its three fixed siblings.
- Produces: one authored foundation/theme/style set consumed directly by
  `web/site.css`.
- Consumes: existing `writeVendored` conflict semantics unchanged.

- [ ] **Step 1: Change the CLI init test first.**

  Make `TestInitWritesEverything` require:

  ```go
  for _, p := range []string{
      "gsxui.json",
      "web/gsxui/index.css",
      "web/gsxui/foundation.css",
      "web/gsxui/theme.css",
      "web/gsxui/style.css",
      "web/gsxui/gsxui.js",
      "web/gsxui/index.js",
      "ui/merge/merge.go",
      "gsx.toml",
  } {
      // existing stat assertion
  }
  ```

  Assert `index.css` imports `./foundation.css`, `./theme.css`, and
  `./style.css`; assert `theme.css` contains `--primary`; assert
  `foundation.css` contains `@theme inline`; assert `style.css` contains the
  ScrollArea pseudo-element rule.

- [ ] **Step 2: Run the focused CLI test and record RED.**

  Run:

  ```bash
  go test ./internal/cli -run TestInitWritesEverything -count=1
  ```

  Expected: FAIL because the default is still `web/gsxui.css` and the four
  files do not exist.

- [ ] **Step 3: Move the existing CSS into owned files.**

  `assets/css/index.css` is exactly the consumer wiring:

  ```css
  @import "tailwindcss";
  @import "tw-animate-css";
  @import "./foundation.css";
  @import "./theme.css";
  @import "./style.css";
  ```

  Move `@custom-variant`, `@theme inline`, global `@layer base`, Accordion
  animation mechanics, and range-input UA normalization to
  `foundation.css`. Move `:root`/`.dark` values to
  `themes/default.css`. Move the Slider track/thumb presentation and
  ScrollArea scrollbar presentation to `styles/default.css`, wrapped in
  `@layer components`.

  Extend `@theme inline` and the default theme with:

  ```css
  --color-sidebar: var(--sidebar);
  --color-sidebar-foreground: var(--sidebar-foreground);
  --color-sidebar-primary: var(--sidebar-primary);
  --color-sidebar-primary-foreground: var(--sidebar-primary-foreground);
  --color-sidebar-accent: var(--sidebar-accent);
  --color-sidebar-accent-foreground: var(--sidebar-accent-foreground);
  --color-sidebar-border: var(--sidebar-border);
  --color-sidebar-ring: var(--sidebar-ring);
  --color-success: var(--success);
  --color-info: var(--info);
  --color-warning: var(--warning);
  --color-overlay: var(--overlay);
  --color-contrast: var(--contrast);
  ```

  Use shadcn's current neutral sidebar defaults:

  ```css
  /* light */
  --sidebar: oklch(0.985 0 0);
  --sidebar-foreground: oklch(0% 0 0);
  --sidebar-primary: oklch(0.205 0 0);
  --sidebar-primary-foreground: oklch(0.985 0 0);
  --sidebar-accent: oklch(0.97 0 0);
  --sidebar-accent-foreground: oklch(0.205 0 0);
  --sidebar-border: oklch(0.922 0 0);
  --sidebar-ring: oklch(0.708 0 0);

  /* dark */
  --sidebar: oklch(0.205 0 0);
  --sidebar-foreground: oklch(0.985 0 0);
  --sidebar-primary: oklch(0.488 0.243 264.376);
  --sidebar-primary-foreground: oklch(0.985 0 0);
  --sidebar-accent: oklch(0.269 0 0);
  --sidebar-accent-foreground: oklch(0.985 0 0);
  --sidebar-border: oklch(1 0 0 / 10%);
  --sidebar-ring: oklch(0.439 0 0);
  ```

  Preserve current status/overlay/contrast colors with semantic values:

  ```css
  --success: oklch(69.6% 0.17 162.48);
  --info: oklch(68.5% 0.169 237.323);
  --warning: oklch(76.9% 0.188 70.08);
  --overlay: oklch(0% 0 0 / 10%);
  --contrast: oklch(100% 0 0);
  ```

  Dark mode uses the same five values in Phase 1.

- [ ] **Step 4: Make the site import the shared sources.**

  Keep all `@import` statements at the beginning of `web/site.css`:

  ```css
  @import "tailwindcss";
  @import "tw-animate-css";
  @import "../assets/css/foundation.css";
  @import "../assets/css/themes/default.css";
  @import "../assets/css/styles/default.css";
  @import "@fontsource-variable/geist";
  @import "@fontsource-variable/geist-mono";
  ```

  Retain site-only font mappings, `@source` declarations, branding, docs
  typography, and syntax highlighting. Delete every copied library token,
  base, Accordion, Slider, and ScrollArea block.

- [ ] **Step 5: Implement deterministic CLI mapping.**

  Add:

  ```go
  type cssAssetTarget struct {
      source string
      target string
  }

  func cssAssetTargets(entry string) []cssAssetTarget {
      dir := filepath.Dir(entry)
      return []cssAssetTarget{
          {source: "assets/css/index.css", target: entry},
          {source: "assets/css/foundation.css", target: filepath.Join(dir, "foundation.css")},
          {source: "assets/css/themes/default.css", target: filepath.Join(dir, "theme.css")},
          {source: "assets/css/styles/default.css", target: filepath.Join(dir, "style.css")},
      }
  }
  ```

  `runInit` reads and writes all four through `writeVendored`. Change
  `DefaultConfig().CSS` to `web/gsxui/index.css`.

- [ ] **Step 6: Retarget drift tests to the single sources.**

  `TestThemeDefaultsDriftPin` reads
  `assets/css/themes/default.css`. Replace the copied-block Accordion test
  with a test that reads `assets/css/foundation.css` and asserts the
  Accordion mechanism is present once; it no longer compares two copies.

- [ ] **Step 7: Verify GREEN and the real CSS build.**

  Run:

  ```bash
  go test ./internal/cli ./site/pages -count=1
  npx @tailwindcss/cli -i web/site.css -o jstest/.tmp/phase1-task1.css
  go test ./...
  ```

  Expected: PASS; Tailwind completes without unknown `@apply` candidates.

- [ ] **Step 8: Check Go diagnostics and commit.**

  Run:

  ```bash
  gopls check -severity=hint internal/cli/config.go internal/cli/init.go
  git add assets web/site.css internal/cli site/pages
  git commit -m "refactor: split theme CSS assets"
  ```

---

### Task 2: Add tokenized slot composition and the typed contract

**Files:**

- Create: `ui/slots.go`
- Create: `ui/slots_test.go`
- Create: `internal/stylecontract/contract.go`
- Create: `internal/stylecontract/contract_test.go`
- Create: `internal/stylecontract/contracts_primitives.go`
- Create: `internal/stylecontract/contracts_forms.go`
- Create: `internal/stylecontract/contracts_overlays.go`
- Create: `internal/stylecontract/contracts_menus.go`
- Create: `internal/stylecontract/contracts_composites.go`
- Create: `internal/stylecontract/contracts_sidebar.go`
- Create: `internal/stylecontract/contracts_sonner.go`
- Test: `ui/slots_test.go`
- Test: `internal/stylecontract/contract_test.go`

**Interfaces:**

- Produces: `withSlot(name string, attrs gsx.Attrs) gsx.Attrs` in package
  `ui`.
- Produces: `stylecontract.Axis`, `stylecontract.Slot`,
  `stylecontract.Component`, `stylecontract.All() []Component`, and
  `stylecontract.Validate([]Component) error`.
- Produces: empty, fixed-order per-family contract slices for later tasks to
  fill.

- [ ] **Step 1: Write failing slot-composition tests.**

  In package `ui`, test these exact cases:

  ```go
  func TestWithSlotPrependsAndDeduplicates(t *testing.T) {
      attrs := gsx.Attrs{
          {Key: "data-gsxui-slot", Value: "alert-dialog-action button"},
          {Key: "class", Value: "h-12"},
          {Key: "id", Value: "save"},
      }
      got := withSlot("button", attrs)
      if v, _ := got.Get("data-gsxui-slot"); v != "button alert-dialog-action" {
          t.Fatalf("slot tokens = %q", v)
      }
      if got.Class() != "h-12" {
          t.Fatalf("class = %q", got.Class())
      }
      if v, _ := got.Get("id"); v != "save" {
          t.Fatalf("id = %q", v)
      }
  }
  ```

  Also test nil attrs and three nested calls produce
  `"button pagination-link pagination-previous"` in that order.

- [ ] **Step 2: Run the slot tests and record RED.**

  Run:

  ```bash
  go test ./ui -run TestWithSlot -count=1
  ```

  Expected: compile FAIL because `withSlot` does not exist.

- [ ] **Step 3: Implement the minimal helper.**

  `withSlot`:

  - starts with `name`;
  - reads the last scalar `data-gsxui-slot` value via `Attrs.Get`;
  - splits it with `strings.Fields`;
  - preserves first occurrence order and removes duplicates;
  - returns a new bag whose first entry is the normalized slot attribute,
    followed by `attrs.Without("data-gsxui-slot")`;
  - never mutates the input slice.

  An empty `name` is a programmer error and panics with
  `"gsxui: empty slot name"`; add a test for it.

- [ ] **Step 4: Write failing contract validation tests.**

  Use literal contracts and require errors for:

  - empty component name;
  - duplicate component name;
  - empty slot name;
  - duplicate slot token globally;
  - empty axis attribute;
  - duplicate axis value;
  - unsorted `All()` output.

  A presence-only axis has `Values == nil` and is valid.

- [ ] **Step 5: Run contract tests and record RED.**

  Run:

  ```bash
  go test ./internal/stylecontract -count=1
  ```

  Expected: compile FAIL because the package does not exist.

- [ ] **Step 6: Implement contract types and fixed ordering.**

  Use:

  ```go
  type Axis struct {
      Attribute string
      Values    []string
  }

  type Slot struct {
      Name string
      Axes []Axis
  }

  type Component struct {
      Name  string
      Slots []Slot
  }
  ```

  Each `contracts_*.go` declares one package-private slice, initially empty.
  `All` concatenates the slices in the file-structure order above, copies
  the result, and sorts by component name. `Validate` returns contextual
  errors naming the component, slot, attribute, or value.

- [ ] **Step 7: Verify, check diagnostics, and commit.**

  Run:

  ```bash
  go test ./ui -run TestWithSlot -count=1
  go test ./internal/stylecontract -count=1
  gopls check -severity=hint ui/slots.go internal/stylecontract/contract.go
  git add ui/slots.go ui/slots_test.go internal/stylecontract
  git commit -m "feat: define component style contract"
  ```

---

### Task 3: Add browser-level visual and cascade safety nets

**Files:**

- Create: `jstest/harness/style_contract.gsx`
- Create: generated `jstest/harness/style_contract.x.go`
- Create: `jstest/specs/style-visual.spec.ts`
- Create: `jstest/specs/style-visual.spec.ts-snapshots/*.png`
- Modify: `jstest/harness/main.go`
- Modify: `jstest/playwright.config.ts`
- Test: `jstest/specs/style-visual.spec.ts`

**Interfaces:**

- Produces: `/f/style-contract` fixture with default and caller-overridden
  Button, Input, Card, and Badge instances.
- Produces: fixed visual baselines for representative examples before their
  classes move.
- Produces: a computed-style assertion that a caller utility wins over the
  current library default; later tasks keep this passing through cascade
  layers.

- [ ] **Step 1: Register a real GSX fixture.**

  The fixture renders:

  ```gsx
  component StyleContractFixture() {
      <div class="grid gap-4 p-6">
          <Button>Default</Button>
          <Button class="h-12 rounded-none">Caller override</Button>
          <Input value="Ada"/>
          <Card><CardContent>Card</CardContent></Card>
          <Badge>Badge</Badge>
      </div>
  }
  ```

  Add a dedicated `GET /f/style-contract` handler in
  `jstest/harness/main.go`. Render `StyleContractFixture()` through the
  existing `renderShell` helper with `/ui/index.js`, then run
  `go tool gsx generate`. The route is test-only and does not enter the
  production example registry.

- [ ] **Step 2: Write computed-style assertions.**

  Assert the caller-overridden button computes:

  ```ts
  expect(await override.evaluate((el) => {
    const css = getComputedStyle(el);
    return {
      height: css.height,
      borderRadius: css.borderRadius,
      display: css.display,
    };
  })).toEqual({
    height: "48px",
    borderRadius: "0px",
    display: "inline-flex",
  });
  ```

  This protects the user-facing override contract, not source text.

- [ ] **Step 3: Add deterministic screenshots.**

  Use a 1280×900 viewport, disable animations/caret, wait for
  `document.fonts.ready`, and screenshot the body for light and dark on:

  ```text
  accordion/basic
  alert/variants
  button/variants
  card/compound
  calendar/loadedrange
  checkbox/states
  combobox/basic
  dialog/basic (opened before capture)
  dropdown/basic (opened before capture)
  field/invalid
  navigation-menu/mega
  sidebar/variants
  sonner/types
  tabs/basic
  ```

  Add 390×844 mobile captures for `button/variants`, `dialog/basic`,
  `sidebar/variants`, and `calendar/basic`. Set
  `maxDiffPixelRatio: 0.01` and `animations: "disabled"` in each assertion;
  do not loosen global Playwright expectations.

- [ ] **Step 4: Generate the baseline from the pre-migration rendering.**

  Run:

  ```bash
  npx playwright test --config jstest/playwright.config.ts \
    jstest/specs/style-visual.spec.ts --update-snapshots
  npx playwright test --config jstest/playwright.config.ts \
    jstest/specs/style-visual.spec.ts
  ```

  Expected: the update creates reviewed PNGs; the second run reports all
  tests PASS with zero unexpected diffs.

- [ ] **Step 5: Commit the safety net.**

  Run:

  ```bash
  git add jstest
  git commit -m "test: pin component visual contract"
  ```

---

### Task 4: Migrate primitives and composition seams

**Files:**

- Modify:
  `ui/{alert,aspect-ratio,avatar,badge,breadcrumb,button,button-group,card,empty,item,kbd,label,pagination,separator,skeleton,spinner,table}.gsx`
- Modify: `ui/icon/icon.gsx`
- Modify: `ui/avatar.js`
- Modify: corresponding `ui/*_test.go` and `ui/icon/icon_test.go`
- Modify: `assets/css/styles/default.css`
- Modify: `internal/stylecontract/contracts_primitives.go`
- Regenerate: corresponding `.x.go`

**Interfaces:**

- Consumes: `withSlot` and the default style asset.
- Produces: contracts and CSS for the 18 listed component families.
- Produces: token composition for Button-based Pagination links,
  Label-based Field use later, Separator-based uses later, and Icon.
- Produces: Avatar behavior hooks
  `data-gsxui-avatar`, `data-gsxui-avatar-image`, and
  `data-gsxui-avatar-fallback`.

- [ ] **Step 1: Change primitive render tests to the new contract first.**

  For each listed component:

  - replace `data-slot` expectations with `data-gsxui-slot`;
  - remove built-in `class` content from exact pins;
  - retain caller-class tests and require the caller class only;
  - assert every public variant/size/state attribute;
  - assert composed token lists. In particular:

  ```text
  PaginationLink:     button pagination-link
  PaginationPrevious: button pagination-link pagination-previous
  PaginationNext:     button pagination-link pagination-next
  ```

  Update Icon's default pin to have `data-gsxui-slot="icon"` and no default
  class; its caller `class="size-6"` test remains.

- [ ] **Step 2: Run primitive tests and record RED.**

  Run:

  ```bash
  go test ./ui ./ui/icon -run \
    'Test(Alert|AspectRatio|Avatar|Badge|Breadcrumb|Button|Card|Empty|Item|Kbd|Label|Pagination|Separator|Skeleton|Spinner|Table|Icon)' \
    -count=1
  ```

  Expected: FAIL on old `data-slot` and built-in classes.

- [ ] **Step 3: Migrate markup and fill the contract.**

  Every component root/part spreads `withSlot("<token>", attrs)`. Add slot
  tokens to styled anonymous parts such as screen-reader labels and item
  media indicators. Delete base/variant/size class functions after
  reflecting their branch values into existing or new `data-variant`,
  `data-size`, `data-orientation`, and presence attributes.

  Icon uses `data-gsxui-slot="icon"` and forwards caller class without a
  default class. Avatar JS queries dedicated behavior hooks instead of slot
  tokens.

- [ ] **Step 4: Port presentation to relational CSS.**

  Create one `/* Component */` section per component in
  `@layer components`. Port current static utilities via `@apply`.
  Translate composition selectors explicitly:

  - use `[data-gsxui-slot~="button"]` for Button base/axis rules;
  - use `:is(a[data-gsxui-slot~="badge"])` for link-only Badge states;
  - use `:has(...)`, child/sibling combinators, and slot tokens for
    button-group corner restoration;
  - use ancestor slot selectors for Kbd inside Tooltip and Item media with a
    description;
  - size nested icons from the owning component rule rather than passing
    classes to Icon;
  - keep accessibility-only `sr-only` mechanics in foundation under the new
    semantic slot.

  Replace Button/Badge `text-white` with `text-contrast`.

- [ ] **Step 5: Verify unit, visual, and full focused behavior.**

  Run:

  ```bash
  go tool gsx generate
  go test ./ui ./ui/icon -count=1
  npx playwright test --config jstest/playwright.config.ts \
    jstest/specs/style-visual.spec.ts
  npx playwright test --config jstest/playwright.config.ts \
    jstest/specs/invariants.spec.ts
  ```

  Expected: PASS and no screenshot changes.

- [ ] **Step 6: Check diagnostics and commit.**

  Run `gopls check -severity=hint` for every changed authored `.go` test and
  `internal/stylecontract/contracts_primitives.go`, then:

  ```bash
  git add ui assets/css/styles/default.css internal/stylecontract
  git commit -m "refactor: move primitive styling to CSS"
  ```

---

### Task 5: Migrate forms and controls

**Files:**

- Modify:
  `ui/{checkbox,field,input,input-group,input-otp,native-select,progress,radio,select,slider,switch,textarea,toggle,toggle-group}.gsx`
- Modify: `ui/{input-otp,select}.js`
- Modify: corresponding `ui/*_test.go`
- Modify: `assets/css/foundation.css`
- Modify: `assets/css/styles/default.css`
- Modify: `internal/stylecontract/contracts_forms.go`
- Regenerate: corresponding `.x.go`

**Interfaces:**

- Produces: form/control slots and all variant, size, orientation, invalid,
  checked, placeholder, and active axes.
- Produces: token composition:

  ```text
  InputGroupInput:    input input-group-control
  InputGroupTextarea: textarea input-group-control
  FieldLabel:         label field-label
  FieldSeparator:     separator field-separator
  ToggleGroupItem:    toggle toggle-group-item
  SidebarInput later: input sidebar-input
  ```

- Produces: dedicated Select and InputOTP behavior hooks for every part their
  JavaScript queries.

- [ ] **Step 1: Rewrite form/control render tests first.**

  Require namespaced slot tokens, preserved state axes, tokenized
  composition, caller-only classes, and the same dynamic inline values.
  Add exact tests for `InputGroupInput`, `FieldLabel`, `FieldSeparator`, and
  `ToggleGroupItem` token order.

- [ ] **Step 2: Run focused tests and record RED.**

  Run:

  ```bash
  go test ./ui -run \
    'Test(Checkbox|Field|Input|InputGroup|InputOTP|NativeSelect|Progress|Radio|Select|Slider|Switch|Textarea|Toggle)' \
    -count=1
  ```

  Expected: FAIL on old slots/classes.

- [ ] **Step 3: Migrate markup, JS-created caret markup, and behavior queries.**

  Remove class-producing functions/constants. Preserve dynamic styles listed
  in Global Constraints. InputOTP JS creates:

  ```html
  <div data-gsxui-slot="input-otp-caret-overlay"
       data-gsxui-input-otp-caret-overlay>
    <div data-gsxui-slot="input-otp-caret"
         data-gsxui-input-otp-caret></div>
  </div>
  ```

  It never assigns `className`. Select and InputOTP query dedicated
  `data-gsxui-*` hooks only.

- [ ] **Step 4: Port form/control CSS.**

  Translate named group/peer mechanics into relational selectors:

  - InputGroup focus/invalid uses `:has([data-gsxui-slot~="input-group-control"]:focus-visible)`
    and `:has([aria-invalid="true"])`;
  - Field layouts select explicit field/group/content/label tokens;
  - Select and Combobox indicators key off their owning item's
    `data-state`;
  - ToggleGroup layout keys off its own size/spacing and item tokens;
  - nested Icon sizing is owned by the parent slot.

  Keep only native range normalization in foundation; keep its track, thumb,
  colors, and focus treatment in the default style.

- [ ] **Step 5: Verify focused tests, JS syntax, and visuals.**

  Run:

  ```bash
  go tool gsx generate
  go test ./ui -count=1
  node --check ui/input-otp.js
  node --check ui/select.js
  npx playwright test --config jstest/playwright.config.ts \
    jstest/specs/style-visual.spec.ts \
    jstest/specs/invariants.spec.ts \
    jstest/specs/selector-coverage.spec.ts
  ```

- [ ] **Step 6: Check diagnostics and commit.**

  Run the required `gopls check` commands, then:

  ```bash
  git add ui assets/css internal/stylecontract
  git commit -m "refactor: move form control styling to CSS"
  ```

---

### Task 6: Migrate disclosure and overlay components

**Files:**

- Modify:
  `ui/{accordion,collapsible,dialog,alert-dialog,drawer,sheet,popover,hover-card,tooltip}.gsx`
- Modify: `ui/dialog.js`
- Modify: corresponding `ui/*_test.go`
- Modify: `assets/css/foundation.css`
- Modify: `assets/css/styles/default.css`
- Modify: `internal/stylecontract/contracts_overlays.go`
- Regenerate: corresponding `.x.go`

**Interfaces:**

- Produces: overlay token composition:

  ```text
  AlertDialog:        dialog alert-dialog
  AlertDialogContent: dialog-content alert-dialog-content
  AlertDialogAction:  button alert-dialog-action
  AlertDialogCancel:  button alert-dialog-cancel
  Drawer:             dialog drawer
  DrawerContent:      dialog-content drawer-content
  Sheet:              dialog sheet
  SheetContent:       dialog-content sheet-content
  ```

- Produces: Dialog title/description behavior hooks; Dialog JS no longer
  queries styling slots.

- [ ] **Step 1: Rewrite overlay render tests first.**

  Pin token order, semantic roles, open/closed state, side/direction values,
  caller-only classes, and composed Dialog/Button tokens.

- [ ] **Step 2: Run focused tests and record RED.**

  Run:

  ```bash
  go test ./ui -run \
    'Test(Accordion|Collapsible|Dialog|AlertDialog|Drawer|Sheet|Popover|HoverCard|Tooltip)' \
    -count=1
  ```

- [ ] **Step 3: Migrate markup and behavior hooks.**

  Remove class overrides passed into composed Dialog/Button components.
  Add semantic tokens for injected close buttons, icons, overlay
  descriptions, and screen-reader-only labels. Dialog JS queries
  `data-gsxui-dialog-title` and `data-gsxui-dialog-description`.

- [ ] **Step 4: Port overlay CSS without class markers.**

  Use slot/state/side ancestry for:

  - Dialog, Drawer, and Sheet positioning and top-layer/backdrop animation;
  - Drawer header alignment based on content `data-side`;
  - Tooltip Kbd treatment;
  - Accordion/Collapsible content mechanics;
  - AlertDialog Button composition.

  Replace `backdrop:bg-black/10` with a raw backdrop rule using
  `var(--overlay)`. Do not duplicate Dialog base declarations in Drawer,
  Sheet, or AlertDialog selectors; their tokenized elements already match
  Dialog rules.

- [ ] **Step 5: Verify behavior and visuals.**

  Run:

  ```bash
  go tool gsx generate
  go test ./ui -count=1
  node --check ui/dialog.js
  npx playwright test --config jstest/playwright.config.ts \
    jstest/specs/style-visual.spec.ts \
    jstest/specs/invariants.spec.ts
  ```

- [ ] **Step 6: Check diagnostics and commit.**

  ```bash
  git add ui assets/css internal/stylecontract
  git commit -m "refactor: move overlay styling to CSS"
  ```

---

### Task 7: Migrate menus, Command, and Combobox

**Files:**

- Modify:
  `ui/{command,combobox,dropdown,context-menu,menubar,navigation-menu}.gsx`
- Modify:
  `ui/{command,combobox,dropdown,context-menu,menubar,navigation-menu}.js`
- Modify: corresponding `ui/*_test.go`
- Modify: `assets/css/foundation.css`
- Modify: `assets/css/styles/default.css`
- Modify: `internal/stylecontract/contracts_menus.go`
- Regenerate: corresponding `.x.go`

**Interfaces:**

- Produces: complete menu/listbox slots and checked, disabled, highlighted,
  inset, open, side, active, and viewport axes.
- Produces: dedicated behavior hooks for every group, label, separator,
  value, and item query.
- Produces: `CommandDialog` composition through Dialog and Command tokens
  instead of a Command class override.

- [ ] **Step 1: Rewrite menu render tests first.**

  Remove exact class pins while retaining exact semantic markup, ARIA,
  behavior hooks, reflected state, caller classes, indicator slots, and
  composed token order. Add a `CommandDialog` pin proving the dialog-content
  and command tokens both reach their declared root and input slot elements
  without a library class argument.

- [ ] **Step 2: Run focused tests and record RED.**

  Run:

  ```bash
  go test ./ui -run \
    'Test(Command|Combobox|Dropdown|ContextMenu|Menubar|NavigationMenu)' \
    -count=1
  ```

- [ ] **Step 3: Migrate markup and JS queries.**

  Add behavior hooks for Command/Combobox group/heading/separator/empty and
  Select-like values. Replace all JS `querySelector('[data-slot=…]')`
  calls with dedicated hooks. Delete navigation `group` classes and all
  internal class overrides.

- [ ] **Step 4: Port relational menu CSS.**

  Express indicator visibility with item `data-state`, submenu positioning
  with the owning content/trigger slots, CommandDialog descendant sizing
  with dialog/command token ancestry, and navigation trigger-icon rotation
  with trigger `data-state`. Preserve popover open/closed ghost-box
  behavior and current side animations.

- [ ] **Step 5: Verify unit and real-browser menu behavior.**

  Run:

  ```bash
  go tool gsx generate
  go test ./ui -count=1
  for f in ui/{command,combobox,dropdown,context-menu,menubar,navigation-menu}.js; do
      node --check "$f"
  done
  npx playwright test --config jstest/playwright.config.ts \
    jstest/specs/style-visual.spec.ts \
    jstest/specs/dropdown.spec.ts \
    jstest/specs/navigation-menu.spec.ts \
    jstest/specs/invariants.spec.ts \
    jstest/specs/selector-coverage.spec.ts
  ```

- [ ] **Step 6: Check diagnostics and commit.**

  ```bash
  git add ui assets/css internal/stylecontract
  git commit -m "refactor: move menu styling to CSS"
  ```

---

### Task 8: Migrate Calendar, Carousel, Resizable, and ScrollArea

**Files:**

- Modify: `ui/{calendar,carousel,resizable,scroll-area}.gsx`
- Modify: `ui/{carousel,resizable}.js`
- Modify: corresponding `ui/*_test.go`
- Modify: `assets/css/foundation.css`
- Modify: `assets/css/styles/default.css`
- Modify: `internal/stylecontract/contracts_composites.go`
- Regenerate: corresponding `.x.go`

**Interfaces:**

- Produces: Calendar's complete root/caption/nav/grid/week/day/button slots
  and selection/range/focus/today/outside/disabled axes.
- Produces: dedicated Carousel and Resizable behavior hooks for every
  queried part.
- Preserves: Calendar's Go/JS date agreement and all dynamic values.

- [ ] **Step 1: Rewrite composite render tests first.**

  Add missing semantic tokens for Calendar's anonymous table/grid parts.
  Pin its 42-cell semantics and state attributes without class strings.
  Pin Carousel/Resizable/ScrollArea slots, roles, state, caller class, and
  dynamic flex/style values.

- [ ] **Step 2: Run focused tests and record RED.**

  Run:

  ```bash
  go test ./ui -run \
    'Test(Calendar|Carousel|Resizable|ScrollArea)' -count=1
  ```

- [ ] **Step 3: Migrate markup and behavior selectors.**

  Delete Calendar class constants. Replace named day/calendar groups with
  day and day-button slot/state relationships. Carousel and Resizable JS
  use dedicated hooks.

- [ ] **Step 4: Separate mechanics from presentation.**

  Foundation owns layout without which the live algorithms fail:

  - Carousel scroll viewport/item mechanics;
  - Resizable flex panel group, panel basis consumption, and operable handle
    hit target;
  - closed/native visibility guarantees.

  Default style owns gaps, colors, shape, typography, shadows, decorative
  scrollbar/slider-like pseudo-elements, and Calendar presentation.

- [ ] **Step 5: Verify all Calendar behavior and visual baselines.**

  Run:

  ```bash
  go tool gsx generate
  go test ./ui -count=1
  node --check ui/carousel.js
  node --check ui/resizable.js
  npx playwright test --config jstest/playwright.config.ts \
    jstest/specs/calendar.spec.ts \
    jstest/specs/style-visual.spec.ts \
    jstest/specs/invariants.spec.ts
  ```

- [ ] **Step 6: Check diagnostics and commit.**

  ```bash
  git add ui assets/css internal/stylecontract
  git commit -m "refactor: move composite styling to CSS"
  ```

---

### Task 9: Migrate Sidebar as an explicit relational state machine

**Files:**

- Modify: `ui/sidebar.gsx`
- Modify: `ui/sidebar.js`
- Modify: `ui/sidebar_test.go`
- Modify: `assets/css/foundation.css`
- Modify: `assets/css/styles/default.css`
- Modify: `assets/css/themes/default.css`
- Modify: `internal/stylecontract/contracts_sidebar.go`
- Modify: `site/pages/theme.gsx`
- Modify: `site/pages/theme_defaults_test.go`
- Regenerate: corresponding `.x.go`

**Interfaces:**

- Produces: all Sidebar slots and its side, state, variant, collapsible,
  mobile, size, active, and show-on-hover axes.
- Produces: Input/Separator/Dialog/Tooltip token composition for Sidebar's
  nested primitives.
- Produces: dedicated Sidebar behavior hooks; Sidebar JS contains no slot or
  class queries.
- Produces: live sidebar semantic tokens in the theme editor/default drift
  pin.

- [ ] **Step 1: Rewrite Sidebar tests first.**

  Pin the two-tree semantic structure, tokenized primitive composition,
  dynamic width custom properties, and state axes. Replace every class
  assertion with a slot/state/behavior assertion except caller-class tests.

  Add a theme-default test proving all eight sidebar tokens exist in both
  light and dark editor/default maps.

- [ ] **Step 2: Run focused tests and record RED.**

  Run:

  ```bash
  go test ./ui ./site/pages -run 'Test(Sidebar|ThemeDefaults)' -count=1
  ```

- [ ] **Step 3: Migrate Sidebar markup and JS.**

  Remove every named group/peer class and variant/size class function.
  Preserve all state on explicit attributes. Compose nested Input,
  Separator, Sheet/Dialog, Button, Skeleton, and Tooltip slots instead of
  passing library styling classes. Sidebar JS uses dedicated wrapper,
  desktop, trigger, and rail hooks only.

- [ ] **Step 4: Port Sidebar CSS as relational selectors.**

  Translate every current group/peer relationship explicitly:

  - wrapper state controls desktop width, offcanvas position, icon width,
    rail cursor, group label visibility, and sub-menu visibility;
  - side controls borders and rail position;
  - variant controls inset/floating treatment;
  - menu-item `:has(menu-action)` reserves space;
  - menu-button size/active state controls adjacent action/badge position;
  - show-on-hover uses menu-item hover/focus-within and menu-button active
    state;
  - mobile tree uses `data-mobile` rather than desktop ancestry.

  Every one of the standard sidebar tokens added in Task 1 must compile into
  at least one real declaration and have a computed-style assertion.

- [ ] **Step 5: Update the expected Sidebar screenshots only for the dead-token correction.**

  Run the Sidebar visual spec without update first. Inspect the diff and
  record that changed pixels are limited to the now-live sidebar colors.
  Then run:

  ```bash
  npx playwright test --config jstest/playwright.config.ts \
    jstest/specs/style-visual.spec.ts -g sidebar --update-snapshots
  npx playwright test --config jstest/playwright.config.ts \
    jstest/specs/style-visual.spec.ts -g sidebar
  ```

  No layout/density delta is accepted.

- [ ] **Step 6: Verify Sidebar behavior and commit.**

  Run:

  ```bash
  go tool gsx generate
  go test ./ui ./site/pages -count=1
  node --check ui/sidebar.js
  npx playwright test --config jstest/playwright.config.ts \
    jstest/specs/style-visual.spec.ts \
    jstest/specs/invariants.spec.ts
  gopls check -severity=hint site/pages/theme_defaults_test.go \
    internal/stylecontract/contracts_sidebar.go
  git add ui assets/css internal/stylecontract site/pages jstest
  git commit -m "refactor: move sidebar styling to CSS"
  ```

---

### Task 10: Migrate Sonner and all client-created component markup

**Files:**

- Modify: `ui/sonner.gsx`
- Modify: `ui/sonner.js`
- Modify: `ui/sonner_test.go`
- Modify: `assets/css/foundation.css`
- Modify: `assets/css/styles/default.css`
- Modify: `internal/stylecontract/contracts_sonner.go`
- Modify: relevant `jstest/specs/*.spec.ts`
- Regenerate: `ui/sonner.x.go`

**Interfaces:**

- Produces: explicit Toaster, toast, icon, content, title, description,
  action, cancel, and close slot tokens.
- Produces: fallback/client-cloned markup with slots and behavior hooks but
  no class assignment.
- Preserves: inline transform/opacity/z-index/pointer-event values because
  they are live stack geometry/state, not static presentation.

- [ ] **Step 1: Rewrite Sonner render/browser tests first.**

  Replace the enormous class-string pins with exact semantic/slot/ARIA
  structure. Add a browser test that removes the server-rendered Toaster,
  triggers fallback creation, and asserts the fallback has the same Toaster
  slot/hook contract and computed fixed positioning as the server path.

- [ ] **Step 2: Run focused tests and record RED.**

  Run:

  ```bash
  go test ./ui -run TestSonner -count=1
  npx playwright test --config jstest/playwright.config.ts -g 'sonner|toast'
  ```

- [ ] **Step 3: Migrate server and client markup.**

  `sonner.gsx` uses `withSlot` for every styled part and nested Icon. JS
  fallback creation sets `data-gsxui-slot="toaster"` and the same dedicated
  hook as the server path; it never assigns `className`. Adoption queries
  `li[data-gsxui-toast]`, not a styling slot.

- [ ] **Step 4: Port static Sonner CSS.**

  Foundation owns fixed region/stack mechanics required for adoption and
  live transforms. Default style owns toast width, spacing, surfaces,
  borders, type, icons, buttons, transitions, and shape. Use
  `text-success`, `text-info`, `text-warning`, and `text-destructive`; no
  Tailwind palette name remains in component CSS.

- [ ] **Step 5: Verify and commit.**

  Run:

  ```bash
  go tool gsx generate
  go test ./ui -count=1
  node --check ui/sonner.js
  npx playwright test --config jstest/playwright.config.ts \
    jstest/specs/style-visual.spec.ts \
    jstest/specs/invariants.spec.ts \
    jstest/specs/smoke.spec.ts
  gopls check -severity=hint internal/stylecontract/contracts_sonner.go
  git add ui assets/css internal/stylecontract jstest
  git commit -m "refactor: move toast styling to CSS"
  ```

---

### Task 11: Prove the boundary, update the current editor/handoff, and close Phase 1

**Files:**

- Create: `jstest/foundation.css`
- Create: `jstest/specs/foundation-contract.spec.ts`
- Modify: `jstest/global-setup.ts`
- Modify: `jstest/harness/shell.go`
- Modify: `jstest/support/paths.ts`
- Modify: `internal/stylecontract/contract_test.go`
- Modify: `internal/registry/registry_test.go`
- Create: `site/examples/style_contract_test.go`
- Modify: `Makefile`
- Modify: `web/theme.js`
- Modify: `site/pages/theme.gsx`
- Modify: `site/pages/theming.gsx`
- Modify: `site/pages/getting_started.gsx`
- Modify: `site/snippets/*` where CSS paths/output changed
- Modify: `README.md`
- Modify: `NOTICE.md`
- Modify: `dev/preview.html`
- Modify: current examples/comments that query `data-slot`
- Modify: `docs/theme-system-roadmap.md`
- Modify: `docs/backlog.md`
- Modify: `docs/jsx-parity.md`
- Regenerate: affected `.x.go` and `site/hl/blocks.gen.go`

**Interfaces:**

- Produces: a foundation-only compiled stylesheet and harness mode.
- Produces: exact registry ↔ typed-contract coverage.
- Produces: current editor export as `theme.css`, not a combined entry.
- Produces: documentation for the one-time breaking migration.
- Produces: `make ci` as the authoritative uncached gate if it does not
  already exist.

- [ ] **Step 1: Write failing whole-contract tests.**

  Add tests that:

  - compare `registry.Components()` with `stylecontract.All()` component
    names exactly, including `icon`;
  - call `stylecontract.Validate(stylecontract.All())`;
  - render every `examples.For(component)` node, parse the result with
    `golang.org/x/net/html`, split every `data-gsxui-slot` with
    `strings.Fields`, and compare the emitted slot set with the complete
    slot set in `stylecontract.All()` in both directions;
  - for every emitted axis attribute on a contracted slot, reject values not
    declared by that slot, and require every declared non-presence axis
    value to appear in at least one rendered example;
  - compile a foundation-only browser stylesheet;
  - prove dialog open/close, popover closed visibility, Carousel navigation,
    Resizable keyboard resizing, InputOTP entry, and Sonner adoption still
    function without `styles/default.css`;
  - prove the full stylesheet gives caller utilities precedence on the
    style-contract fixture.

- [ ] **Step 2: Run the new gates and record RED.**

  Run:

  ```bash
  go test ./internal/stylecontract ./internal/registry ./site/examples -count=1
  npx playwright test --config jstest/playwright.config.ts \
    jstest/specs/foundation-contract.spec.ts
  ```

  Expected: FAIL until coverage and the foundation-only harness mode exist.

- [ ] **Step 3: Add the foundation-only browser build and route mode.**

  `jstest/foundation.css` imports Tailwind, `tw-animate-css`, foundation,
  default theme, and fixture/example sources, but not the default style.
  Global setup compiles it beside the full CSS. The harness selects it only
  for an explicit `?css=foundation` request; ordinary tests continue loading
  the full site CSS.

- [ ] **Step 4: Move any behavior-critical declarations exposed by the tests.**

  Move only the smallest coherent mechanical rule from default style to
  foundation. Re-run the failing behavior test after each move. Do not add
  cosmetic fallback styles to make foundation-only screenshots attractive;
  foundation-only has no screenshot expectations.

- [ ] **Step 5: Make the current editor and docs describe the split artifacts.**

  `web/theme.js` generates/downloads `theme.css` containing only semantic
  variables. It does not repeat Tailwind imports, foundation, or style.
  Update the raw editor's field list for sidebar/status/overlay/contrast
  tokens and keep drift tests exact.

  Getting Started imports `web/gsxui/index.css`. Theming explains
  foundation/theme/style ownership, tokenized slots, caller utilities, and
  that only `theme.css` is replaced for a theme. Document the breaking
  migration:

  ```text
  1. Change the CSS entry from web/gsxui.css to web/gsxui/index.css.
  2. Re-run gsxui init with explicit overwrite after reviewing the four-file diff.
  3. Re-run gsxui add <component> --overwrite for vendored components.
  4. Replace project selectors that intentionally targeted data-slot with
     [data-gsxui-slot~="<token>"].
  ```

  Mark Phase 1 shipped in the roadmap only after every gate below passes;
  leave Phases 2–5 uncommitted.

- [ ] **Step 6: Regenerate authored outputs.**

  Run:

  ```bash
  go tool gsx generate
  make highlight
  ```

  Inspect the diff and confirm only generated `.x.go` and highlighted
  sources corresponding to authored changes moved.

- [ ] **Step 7: Run structural audits.**

  These are reviewer aids, not core validators:

  ```bash
  rg -n 'data-slot|\\[data-slot' ui site web dev NOTICE.md README.md \
    -g '!*.x.go' -g '!*.gen.go'
  rg -n 'className\\s*=|\\.classList|setAttribute\\([\"'\"']class' ui -g '*.js'
  rg -n 'group/|peer/|data-\\[slot|in-data-\\[slot' ui -g '*.gsx'
  rg -n 'class=' ui -g '*.gsx'
  ```

  Expected:

  - no current-source `data-slot` styling/behavior hook;
  - no JavaScript presentation class assignment;
  - no named Tailwind group/peer styling marker in `.gsx`;
  - no authored `class=` remains in `ui/*.gsx`; caller classes flow only
    through the caller attribute spread.

- [ ] **Step 8: Add and run the authoritative gate.**

  If missing, define `make ci` as the uncached equivalent of `make check`:
  Go tests use `-count=1`, browser tests run once, generated drift and all
  formatting/JS checks remain. Then run:

  ```bash
  make ci
  make site
  git diff --check
  git status --short
  ```

  Expected: every command exits 0; 163 existing browser tests plus the new
  style/foundation tests pass; no untracked generated artifact remains;
  restore `site/dist/.gitkeep` if the production build removed it.

- [ ] **Step 9: Check diagnostics and commit.**

  Run `gopls check -severity=hint` on every authored Go file changed in this
  task, then:

  ```bash
  git add Makefile README.md NOTICE.md assets dev docs internal jstest site ui web
  git commit -m "docs: complete CSS-only theme boundary"
  ```

## Final review requirements

After all task reviews pass:

1. Generate one review package from commit `637d089` through `HEAD`.
2. Dispatch the final reviewer on the most capable available model.
3. Require it to build throwaway consumer projects covering:
   - default and custom CSS entry paths;
   - caller utility overrides;
   - composed Button/Input/Label/Separator/Dialog tokens;
   - foundation-only interactive behavior;
   - light/dark sidebar tokens;
   - client-created InputOTP caret and Sonner fallback markup;
   - absence of legacy `data-slot` and presentation class dependencies.
4. Apply one consolidated fix wave for all Critical/Important findings and
   run one scoped re-review.
5. Run `make ci`, `make site`, `git diff --check`, and a clean-status check
   again before using the finishing-a-development-branch workflow.

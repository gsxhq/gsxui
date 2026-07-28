# Inline Style Button Pilot Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deliver the complete Button-only vertical slice approved in `docs/superpowers/specs/2026-07-28-inline-style-button-pilot-design.md`: Nova/Maia style selection, parser-generated inline Button source, safe CLI ownership, and an effective iframe-based web theme editor.

**Architecture:** Keep one tokenized canonical `ui/button.gsx`, parse strict Nova and Maia CSS recipes, and resolve them through the public GSX AST into committed concrete registry sources. A versioned preset is the single theme/style state shared by the CLI and browser. Consumer Button presentation is inline; all other components remain on the CSS-only pack until the explicit post-pilot review.

**Tech Stack:** Go 1.26.1, GSX public parser/AST/formatter, `github.com/tdewolff/parse/v2`, `github.com/mazznoer/csscolorparser`, Tailwind CSS 4, vanilla ESM, structpages, Vite, and Playwright.

## Global Constraints

- Stop after Button. Do not migrate another component or advertise Maia as catalogue-wide.
- Use `/Users/jackieli/personal/gsxhq/gsx` at merged commit `ef72f5eba066d7e87adf7dcadc2db62d00f22efe` and module version `v0.0.0-20260728095825-ef72f5eba066` so bare attribute presence, computed boolean presence semantics, and the public GSX AST/formatter are the tested implementation.
- Preserve `class_merger = "github.com/gsxhq/gsxui/merge.Merge"` in `gsx.toml`. `merge.Merge` is already the genuine Tailwind-aware runtime merger; do not add another merge helper or pre-merge caller classes.
- Do not use regular expressions, selector string splitting, or text replacement to parse CSS, GSX, JSON, or imported theme CSS.
- Do not add style-specific authored templates. Nova and Maia `.gsx` files are generated artifacts from one canonical structure.
- Generated consumer Button source must contain concrete Tailwind utilities, no `gsxui-recipe-*` token, and no private styling helper.
- Keep `data-gsxui-slot-button`, `data-variant`, and `data-size` in rendered output. Bare markers must remain bare through caller attribute bags.
- Remove Button presentation from the consumer `assets/css/styles/default.css` in the same task that installs inline generated Button source. The documentation site may have a site-only Nova fallback, but it must not be vendored by `gsxui init`.
- Preserve all non-Button component behavior and the current CSS-only presentation path.
- Resolve every preset and every intended output byte before the first filesystem mutation. Never silently overwrite a modified managed component.
- Authored `.gsx` files are the source of truth. Run GSX generation and commit generated `.x.go`; never hand-edit generated Go.
- Use test-first steps. Verify each RED failure is caused by the missing behavior, not a broken fixture.
- Run `gopls check -severity=hint` on changed Go and generated Go before each task commit.
- Keep commits task-sized and reviewable. Do not combine the preset model, resolver, CLI transaction, and browser editor into one commit.
- `make ci` with the exact merged GSX core is the authoritative final repository gate.

## Planned Source Layout

```text
internal/preset/
  preset.go             schema, ordered fields, defaults, validation
  json.go               duplicate-safe parse and canonical JSON
  css.go                parser-backed theme CSS import/export
  transport.go          file/stdin/share-code resolution

internal/stylegen/
  recipe.go             strict CSS recipe parser
  resolve.go            public GSX AST transform
  generate.go           deterministic Button artifact generation

registry/
  styles/{nova,maia}/button.css
  generated/{nova,maia}/button.gsx

site/stylepreview/{nova,maia}/
  button.gsx             generated package-adjusted preview fixtures
  button.x.go

web/
  theme.js               editor controller
  theme-state.js         pure resolved/draft state and transport
  theme-preview.js       iframe receiver and synchronization
  site-button.css        site-only canonical Button fallback
```

---

### Task 0: Migrate valued boolean data axes to explicit strings

GSX core `ef72f5eb` intentionally changes a Go `bool` attribute expression to
HTML presence semantics: false omits the attribute and true emits it bare,
except for native HTML attributes whose specification requires the strings
`"true"`/`"false"`. gsxui slot markers and native boolean attributes want that
new default. Five component-presentation axes do not: their style contracts
explicitly require valued `data-*="false|true"` output. Migrate only those
axes before building the pilot on the new core.

**Files:**
- Modify: `go.mod`
- Modify: `go.sum`
- Modify: `ui/field.gsx`
- Generated: `ui/field.x.go`
- Modify: `ui/pagination.gsx`
- Generated: `ui/pagination.x.go`
- Modify: `ui/sidebar.gsx`
- Generated: `ui/sidebar.x.go`

- [ ] **Step 1: Verify the dependency-change RED state**

Create a temporary `go.work` selecting this worktree and
`/Users/jackieli/personal/gsxhq/gsx`, then run generation and the focused
contracts with that workspace:

```bash
presence_workspace="$(mktemp -d)"
(cd "$presence_workspace" && go work init \
  /Users/jackieli/personal/gsxhq/gsxui/.worktrees/css-only-theme-architecture \
  /Users/jackieli/personal/gsxhq/gsx)
GOWORK="$presence_workspace/go.work" go tool gsx generate
GOWORK="$presence_workspace/go.work" go test \
  ./ui ./site/examples -run \
  'Test(FieldSeparator|FormControlCompositionTokenOrder|Pagination|Sidebar|RegisteredExamplesCoverStyleContract)' \
  -count=1
```

Expected RED: the generated diff changes only `ui/field.x.go`,
`ui/pagination.x.go`, and `ui/sidebar.x.go`; focused tests show that false
omits and true bares `data-content`, `data-active`, or
`data-show-on-hover`, conflicting with their declared valued axes.

- [ ] **Step 2: Pin the merged GSX core**

Run:

```bash
go get github.com/gsxhq/gsx@v0.0.0-20260728095825-ef72f5eba066
go get -tool github.com/gsxhq/gsx/cmd/gsx@v0.0.0-20260728095825-ef72f5eba066
go list -m -json github.com/gsxhq/gsx
```

Require `Origin.Hash` to equal
`ef72f5eba066d7e87adf7dcadc2db62d00f22efe`.

- [ ] **Step 3: Make valued axes string-typed at the authored source**

Import standard-library `strconv` in the three authored files and replace only:

```gsx
data-content={strconv.FormatBool(children != nil)}
data-active={strconv.FormatBool(isActive)}
data-show-on-hover={strconv.FormatBool(showOnHover)}
```

Keep every `data-gsxui-slot-*`, native `disabled`, and other boolean presence
attribute as a `bool`. Do not change style contracts or expected HTML: they
already pin the intended distinction.

- [ ] **Step 4: Regenerate and verify GREEN**

```bash
go tool gsx generate
go test ./ui ./site/examples -run \
  'Test(FieldSeparator|FormControlCompositionTokenOrder|Pagination|Sidebar|RegisteredExamplesCoverStyleContract)' \
  -count=1
gopls check -severity=hint \
  ui/field.x.go ui/pagination.x.go ui/sidebar.x.go
```

Expected: focused tests pass; generated output uses string attribute emission
for the five valued expressions while slot/native boolean output remains
presence-based.

- [ ] **Step 5: Run the authoritative gate**

```bash
make ci
git diff --check
```

- [ ] **Step 6: Commit Task 0**

```bash
git add go.mod go.sum \
  ui/field.gsx ui/field.x.go \
  ui/pagination.gsx ui/pagination.x.go \
  ui/sidebar.gsx ui/sidebar.x.go
git commit -m "fix: preserve valued boolean data axes"
```

---

### Task 1: Build the authoritative preset model and canonical transport

**Files:**
- Create: `internal/preset/preset.go`
- Create: `internal/preset/preset_test.go`
- Create: `internal/preset/json.go`
- Create: `internal/preset/json_test.go`
- Create: `internal/preset/css.go`
- Create: `internal/preset/css_test.go`
- Create: `internal/preset/transport.go`
- Create: `internal/preset/transport_test.go`
- Create: `internal/preset/testdata/default-nova.json`
- Create: `internal/preset/testdata/default-theme.css`
- Modify: `assets/css/themes/default.css`
- Modify: `go.mod`
- Modify: `go.sum`

**Interfaces:**

```go
package preset

const SchemaVersion = 1
const SchemaURL = "https://ui.gsxhq.dev/schemas/preset-v1.json"

type Style string

const (
    StyleNova Style = "nova"
    StyleMaia Style = "maia"
)

type ThemeValues map[string]string
type Theme struct {
    Light ThemeValues
    Dark  ThemeValues
}
type Preset struct {
    Schema       string
    SchemaVersion int
    Style        Style
    Radius       string
    Theme        Theme
}

func Styles() []Style
func TokenNames() []string
func Default(style Style) Preset
func Validate(Preset) error
func ParseJSON([]byte) (Preset, error)
func CanonicalJSON(Preset) ([]byte, error)
func EncodeShare(Preset) (string, error)
func DecodeShare(string) (Preset, error)
func ThemeCSS(Preset) ([]byte, error)
func ImportThemeCSS(base Preset, src []byte) (Preset, error)
type InputResolver struct {
    Stdin      io.Reader
    HTTPClient *http.Client
    MaxBytes   int64
}
func (r InputResolver) Resolve(ctx context.Context, argument string) (Preset, error)
```

`Preset` is immutable by convention: every exported function returning a
preset must deep-copy both theme maps. `TokenNames` returns a copy in canonical
CSS order. Typography remains outside schema version 1 for this pilot.

- [ ] **Step 1: Add the validation dependency and prove it accepts gsxui values**

Run:

```bash
go get github.com/mazznoer/csscolorparser@v0.1.8
go mod edit -require=github.com/tdewolff/parse/v2@v2.8.13
```

Write a focused test that passes every current color from
`assets/css/themes/default.css` through `csscolorparser.Parse`. Include
`oklch(1 0 0 / 10%)`, percentages, and the zero-chroma forms already shipped.
If the parser rejects a valid shipped value, stop and choose a standards-aware
parser; do not special-case strings.

- [ ] **Step 2: Write failing schema/default tests**

Pin:

- exact ordered token names from the current default theme;
- `Default(StyleNova)` and `Default(StyleMaia)` deep-copy behavior;
- unknown style;
- missing light or dark token;
- extra token;
- invalid color;
- negative radius;
- radius with multiple values;
- unsupported unit/function; and
- valid zero and non-negative CSS length values.

Radius validation must lex exactly one CSS component value and accept only zero
or a non-negative dimension with a CSS length unit. The schema v1 allowlist is:

```text
px cm mm Q in pc pt
em rem ex rex cap rcap ch rch ic ric lh rlh
vw svw lvw dvw vh svh lvh dvh
vi svi lvi dvi vb svb lvb dvb
vmin svmin lvmin dvmin vmax svmax lvmax dvmax
```

Unit matching is ASCII case-insensitive. Reject percentages, keywords,
multiple tokens, functions including `var()`/`calc()`, and trailing garbage
for schema v1.

- [ ] **Step 3: Implement ordered schema and parser-backed value validation**

Move the existing theme defaults out of `site/pages/theme.gsx` into
`internal/preset/preset.go`; this Go definition becomes the single authored
default. Keep presentation groups as metadata next to the ordered definitions:

```go
type TokenDefinition struct {
    Name  string
    Group string
    Light string
    Dark  string
}
```

Store names without the CSS `--` prefix in the preset and add it only during
CSS generation. Validate colors with `csscolorparser.Parse` and radius with the
CSS lexer. Return field paths such as `theme.dark.sidebar-border` in errors.
Add a repository drift test proving
`ThemeCSS(Default(StyleNova)) == assets/css/themes/default.css`; update the
asset to the canonical generated bytes in this task. Nova and Maia deliberately
share the same semantic defaults in the pilot.

- [ ] **Step 4: Write failing duplicate-safe JSON tests**

Use literal JSON fixtures to assert rejection of:

- duplicate top-level keys;
- duplicate `theme`, `light`, `dark`, and token keys;
- unknown keys at every level;
- missing `$schema`, `schemaVersion`, style, radius, or token;
- numeric/string type mismatches;
- future schema versions; and
- trailing JSON values.

Assert canonical output uses `$schema`, `schemaVersion`, `style`, `radius`,
`theme.light`, and `theme.dark` in that order, tokens in `TokenNames()` order,
two-space indentation, and one trailing newline.

- [ ] **Step 5: Implement a token-stream JSON decoder**

Use `json.Decoder.Token`, tracking seen keys for each object. Do not first
unmarshal into a map or struct because that loses duplicate-key evidence.
Decode only the declared schema, call `Validate`, and hand-write canonical
objects in schema order through `json.Encoder`/buffered helpers. Confirm:

```go
got, _ := ParseJSON(must(CanonicalJSON(Default(StyleNova))))
again, _ := CanonicalJSON(got)
```

is byte-identical to the golden.

- [ ] **Step 6: Write failing CSS import/export tests**

Pin exact `:root` and `.dark` output. CSS import is deliberately compatible:
recognized custom properties update a provided base preset; unrelated
declarations and unrelated custom properties are ignored. Reject the entire
import, without returning partial state, for:

- malformed CSS;
- duplicate recognized declarations in one mode;
- recognized properties outside `:root`/`.dark`;
- invalid color/radius values;
- conflicting radius declarations;
- comma selector lists or selectors that only resemble `.dark`; and
- nested syntax that changes declaration ownership.

At least one partial import must update only `primary` in both modes while
preserving every other base value.

- [ ] **Step 7: Implement CSS import/export with `tdewolff/parse/v2/css`**

Walk parser grammar events and accept declarations owned by exact `:root` or
`.dark` rules. Parse selectors using the CSS parser rather than substring
matching. Normalize property names only by removing the exact leading `--`.
Validate the fully merged candidate once parsing completes, then return it.

- [ ] **Step 8: Write and implement share/input resolution tests**

Canonical share codes are:

```text
gsxui:v1:<base64.RawURLEncoding(canonical-json)>
```

Test file paths, `-` stdin, raw JSON, share codes, and HTTPS URLs. URL loading
uses the injected client, a 10-second request context, a 1 MiB body limit,
2xx-only status handling, and redirect validation that never follows from
HTTPS to another scheme. Tests use a fake `RoundTripper`, not a live server.
Reject `http`, credentials in URLs, fragments, oversized responses, redirect
downgrades, and response bodies that do not resolve to valid JSON/share
content. Decode must reject padded base64, unknown transport versions, invalid
UTF-8, and non-canonical-but-valid payloads only after normalizing them through
`ParseJSON`.

- [ ] **Step 9: Run focused verification**

```bash
go test ./internal/preset -count=1
gopls check -severity=hint internal/preset/*.go
```

Expected: all preset, CSS, JSON, and share round trips pass.

- [ ] **Step 10: Commit Task 1**

```bash
git add go.mod go.sum internal/preset assets/css/themes/default.css
git commit -m "feat: add versioned theme presets"
```

---

### Task 2: Parse strict CSS Button recipes

**Files:**
- Create: `internal/stylegen/recipe.go`
- Create: `internal/stylegen/recipe_test.go`
- Create: `internal/stylegen/testdata/recipe-valid.css`

**Interfaces:**

```go
package stylegen

const RecipePrefix = "gsxui-recipe-"

type Recipe struct {
    Token     string
    Utilities []string
}
type Recipes struct {
    ordered []Recipe
    byToken map[string]Recipe
}

func ParseRecipes(filename string, src []byte) (Recipes, error)
func (r Recipes) Lookup(token string) (Recipe, bool)
func (r Recipes) Tokens() []string
```

- [ ] **Step 1: Write a valid parser test**

The fixture must contain `@layer components`, multiple exact one-class recipe
rules, more than one `@apply` in one rule, Tailwind arbitrary variants, escaped
brackets, and utilities containing colons/slashes. Assert utilities retain
authored order across `@apply` statements.

- [ ] **Step 2: Write a rejection table**

Cover:

- no `@layer components`;
- wrong/nested layer;
- non-recipe rule;
- empty rule or missing `@apply`;
- ordinary declaration;
- duplicate recipe token;
- element, ID, attribute, descendant, child, sibling, pseudo, or compound
  selectors;
- comma selector lists;
- dynamically constructed/escaped selector names;
- nested rule;
- invalid `@apply` syntax; and
- malformed CSS.

Every error assertion must include the fixture filename, line, column, and the
offending selector or at-rule when available.

- [ ] **Step 3: Verify RED**

```bash
go test ./internal/stylegen -run 'TestParseRecipes' -count=1
```

Expected: compilation failure because `ParseRecipes` is absent.

- [ ] **Step 4: Implement the grammar over CSS parser events**

Use `css.NewParser(parse.NewInputBytes(src), false)` and parser grammar
boundaries. Validate the selector token stream represents exactly:

```text
Delim(".") Ident("gsxui-recipe-...")
```

Do not validate selectors by trimming strings. Parse each `@apply` prelude as
CSS tokens, preserve its whitespace-separated utility order, and reject empty
or syntactically unterminated utilities. Return recipes in source order plus
an indexed lookup.

- [ ] **Step 5: Verify and commit**

```bash
go test ./internal/stylegen -run 'TestParseRecipes' -count=1
gopls check -severity=hint internal/stylegen/recipe.go internal/stylegen/recipe_test.go
git add internal/stylegen
git commit -m "feat: parse strict component style recipes"
```

---

### Task 3: Resolve recipe tokens through the public GSX AST

**Files:**
- Create: `internal/stylegen/resolve.go`
- Create: `internal/stylegen/resolve_test.go`
- Create: `internal/stylegen/testdata/resolve-input.gsx`
- Create: `internal/stylegen/testdata/resolve-nova.golden.gsx`

**Interfaces:**

```go
type ResolveReport struct {
    UsedTokens []string
}

func Resolve(filename string, src []byte, recipes Recipes) ([]byte, ResolveReport, error)
func RecipeTokens(filename string, src []byte) ([]string, error)
```

- [ ] **Step 1: Verify the merged GSX core pin**

Run:

```bash
go list -m -json github.com/gsxhq/gsx
```

Confirm:

shows `Origin.Hash` equal to
`ef72f5eba066d7e87adf7dcadc2db62d00f22efe`.

- [ ] **Step 2: Write the successful golden transform**

The canonical fixture must contain:

- invariant classes before/between/after recipe tokens;
- top-level class string parts;
- variant and size `switch` arm literals;
- both anchor and button branches;
- a caller `...attrs` expression; and
- comments/imports that must survive formatting.

Assert exact golden bytes, sorted first-use `ResolveReport`, idempotent GSX
formatting, successful reparsing, and zero `gsxui-recipe-` occurrences.

- [ ] **Step 3: Write all transform rejection tests**

Reject:

- used token missing from the recipe;
- recipe token never used by canonical source;
- duplicate recipe declaration;
- recipe-looking static `class="..."`;
- token embedded inside a larger class token;
- token assembled by concatenation or interpolation;
- recipe token in a non-class expression;
- non-string class expression containing a recipe identifier; and
- malformed GSX.

Also prove unrelated string literals and arbitrary Tailwind classes are
unchanged.

- [ ] **Step 4: Verify RED**

```bash
go test ./internal/stylegen -run 'TestResolve|TestRecipeTokens' -count=1
```

- [ ] **Step 5: Implement the AST transform**

Parse with `github.com/gsxhq/gsx/parser`. Walk with
`github.com/gsxhq/gsx/ast.Inspect`. Visit class `ClassPart.Expr` and every
`ValueSwitch` arm expression reachable from a class part. Import the standard
Go AST as `goast` and parse expressions with `go/parser.ParseExpr`; mutate only
`*goast.BasicLit` string literals after `strconv.Unquote`.

Split literal class content with HTML ASCII whitespace semantics. Replace a
whole recipe token with `strings.Join(recipe.Utilities, " ")`, preserving all
non-recipe tokens and order. Quote with `strconv.Quote`, then format the entire
file using public `gen.Format`.

Perform two independent set checks before output:

```text
used canonical tokens - declared recipe tokens = missing
declared recipe tokens - used canonical tokens = unused
```

Report both sets deterministically in one error.

- [ ] **Step 6: Pin anchor/button structural equivalence**

Add an AST-level helper in the test only that extracts the two resolved class
expressions and compares their formatted expression bytes. Do not compare
rendered class order, because that could conceal a structurally divergent
branch.

- [ ] **Step 7: Verify and commit**

```bash
go test ./internal/stylegen -count=1
gopls check -severity=hint internal/stylegen/*.go
git add internal/stylegen
git commit -m "feat: resolve inline styles through gsx ast"
```

---

### Task 4: Author canonical Button and Nova/Maia recipes

**Files:**
- Modify: `ui/button.gsx`
- Generated: `ui/button.x.go`
- Create: `registry/styles/nova/button.css`
- Create: `registry/styles/maia/button.css`
- Create: `registry/generated/nova/button.gsx`
- Create: `registry/generated/maia/button.gsx`
- Create: `internal/stylegen/generate.go`
- Create: `internal/stylegen/button_test.go`
- Create: `cmd/stylegen/main.go`
- Modify: `embed.go`
- Modify: `Makefile`

**Generated contract:**

```bash
go run ./cmd/stylegen
```

must deterministically regenerate both `registry/generated/*/button.gsx`.
`go run ./cmd/stylegen --check` must diff in memory and fail without writing.

- [ ] **Step 1: Write failing Button token-contract tests**

Before changing Button, test the desired canonical source for:

- the exact base, six variant, and eight size recipe tokens;
- an inline variant switch and inline size switch in each element branch;
- identical anchor/button class expressions;
- bare slot marker and persistent variant/size attributes;
- no concrete presentation utility owned by a recipe; and
- exact set equality with Nova and Maia recipes.

- [ ] **Step 2: Write failing generated-artifact tests**

For both styles assert:

- generated bytes equal the committed file;
- output parses and formats idempotently;
- output contains all existing public parameters/behavior;
- output contains no recipe tokens;
- output contains no `variantClass`, `sizeClass`, `buttonClass`, or new private
  styling helper;
- output includes recognizable concrete utilities for every role; and
- Nova and Maia outputs differ in geometry/radius/border treatment.

- [ ] **Step 3: Convert `ui/button.gsx` to semantic recipe roles**

Keep public signatures and all behavior unchanged. Put the base, variant, and
size switches directly in both `<a>` and `<button>` class expressions, followed
by caller attrs in the existing fallthrough position. Keep invariant
composition classes concrete only when they genuinely do not vary by style.

- [ ] **Step 4: Author the Nova recipe from current shipped Button CSS**

Translate the existing Button block from `assets/css/styles/default.css` into
the strict recipe grammar without changing computed presentation. Preserve
utility ordering needed by Tailwind merge semantics. Record a comment naming
the local source commit used for this translation.

- [ ] **Step 5: Author a real Maia recipe**

Use shadcn's current official Maia Button as the reference, pin the upstream
repository commit and source path in a recipe comment, and adapt only:

- gsxui's public variant/size names;
- gsxui semantic color variables;
- gsxui's existing disabled/link semantics; and
- utilities unavailable in the pinned Tailwind version.

The recipe must be observably different through height/padding, rounding, and
border/shadow treatment. Do not invent a vaguely “Maia-like” heuristic.

- [ ] **Step 6: Implement generator and embedded registry**

`internal/stylegen.GenerateButton(root string, check bool)` reads the canonical
source and both recipe files, calls `Resolve`, and either writes atomically or
reports drift. The CLI wrapper only locates repository root and reports errors.
Extend:

```go
//go:embed ui assets registry merge NOTICE.md
var Files embed.FS
```

Add `generate-styles` and `verify-generated-styles` Make targets. Make
`generate` run style generation before GSX generation; make `check`/`ci` call
the check-only target.

- [ ] **Step 7: Generate and verify**

Use a temporary `go.work` containing this worktree and the exact merged core:

```bash
pilot_workspace="$(mktemp -d)"
(cd "$pilot_workspace" && go work init \
  /Users/jackieli/personal/gsxhq/gsxui/.worktrees/css-only-theme-architecture \
  /Users/jackieli/personal/gsxhq/gsx)
GOWORK="$pilot_workspace/go.work" go run ./cmd/stylegen
GOWORK="$pilot_workspace/go.work" go tool gsx generate
GOWORK="$pilot_workspace/go.work" go test ./internal/stylegen ./ui -count=1
GOWORK="$pilot_workspace/go.work" gopls check -severity=hint \
  ui/button.x.go internal/stylegen/*.go cmd/stylegen/main.go
```

- [ ] **Step 8: Commit Task 4**

```bash
git add ui/button.gsx ui/button.x.go registry internal/stylegen/generate.go \
  internal/stylegen/button_test.go cmd/stylegen embed.go Makefile
git commit -m "feat: generate Nova and Maia Button sources"
```

---

### Task 5: Cut Button over from consumer CSS to inline source

**Files:**
- Modify: `assets/css/styles/default.css`
- Create: `web/site-button.css`
- Modify: `web/site.css`
- Modify: `jstest/support/compiled-css-audit.test.ts`
- Modify: `internal/cli/init_test.go`
- Modify: `jstest/specs/basic-demo-presentation.spec.ts`
- Create: `internal/stylegen/button_css_test.go`

- [ ] **Step 1: Add a failing consumer CSS ownership test**

Parse `assets/css/styles/default.css` with the CSS parser and assert no selector
targets `data-gsxui-slot-button`. Compile `web/site.css` through the existing
Tailwind test setup and assert the site-only fallback does target Button but is
scoped beneath:

```css
body:not([data-theme-button-preview])
```

Also assert `cssAssetTargets` still vendors `style.css` for other components
and never vendors `web/site-button.css`.

- [ ] **Step 2: Move current Button CSS to the site-only fallback**

Remove every Button presentation rule from
`assets/css/styles/default.css`. Re-express the same Nova presentation in
`web/site-button.css`, scoped so normal docs/site uses remain styled while the
theme preview body opts out. Import this file only from `web/site.css`.

Update `internal/cli/init_test.go` to require the vendored style pack to omit
Button, and keep `basic-demo-presentation.spec.ts` green through the site-only
fallback. Keep `internal/stylecontract` unchanged: Button's semantic slot and
axes remain part of the runtime contract.

- [ ] **Step 3: Verify and commit**

```bash
make generate
node --test jstest/support/compiled-css-audit.test.ts
go test ./ui ./internal/stylegen ./internal/cli -count=1
gopls check -severity=hint internal/stylegen/button_css_test.go
git add assets/css/styles/default.css web/site-button.css web/site.css \
  jstest/support/compiled-css-audit.test.ts jstest/specs/basic-demo-presentation.spec.ts \
  internal/stylegen/button_css_test.go internal/cli/init_test.go
git commit -m "refactor: move Button presentation inline"
```

---

### Task 6: Make init/add/diff preset-aware

**Files:**
- Modify: `internal/cli/config.go`
- Modify: `internal/cli/config_test.go`
- Create: `internal/cli/managed.go`
- Create: `internal/cli/managed_test.go`
- Modify: `internal/cli/init.go`
- Modify: `internal/cli/init_test.go`
- Modify: `internal/cli/add.go`
- Modify: `internal/cli/add_test.go`
- Modify: `internal/cli/run.go`
- Modify: `internal/cli/e2e_test.go`

**Config extension:**

```go
type Config struct {
    UI      string            `json:"ui"`
    JS      string            `json:"js"`
    CSS     string            `json:"css"`
    Managed map[string]string `json:"managed,omitempty"`
}
```

Managed keys are slash-normalized paths relative to module root; values are
lowercase SHA-256 hex of the last bytes written by gsxui. `gsxui.json` is not
listed in its own `Managed` map because a self-hash is circular; it is validated
as configuration metadata and written last.

- [ ] **Step 1: Pin config hash semantics**

Test canonical config JSON, deterministic managed-key ordering, path escape
rejection, invalid hash rejection, and backward-compatible loading of current
three-field `gsxui.json`. Loading an old config creates an empty managed map in
memory without mutating disk.

- [ ] **Step 2: Write failing `init --preset` tests**

Cover no preset (Nova default), file, stdin, raw share code, malformed preset,
existing different preset, and `--overwrite`. A successful init must write:

- canonical `gsxui.preset.json`;
- generated `theme.css`;
- the existing non-Button `style.css`;
- support assets and merger; and
- managed hashes for every CLI-owned artifact except `gsxui.json` itself.

All validation must happen before `gsxui.json` or any asset is written.

- [ ] **Step 3: Write failing style-aware add/diff tests**

`add button` must copy the exact embedded generated artifact for the selected
style after existing import/module rewriting. Pin:

- Nova and Maia exact bytes;
- dependency resolution that installs Button in the selected style;
- Maia rejection when the requested dependency closure contains any component
  other than Button;
- no preset file;
- invalid preset;
- absent, identical, modified, and overwrite target;
- managed hash update;
- `--diff` for absent/identical/different files;
- `--diff` does not write files, config, barrel, notice, or run GSX; and
- non-Button components continue to come from `ui/`.

- [ ] **Step 4: Implement shared artifact planning**

Introduce an in-memory artifact description:

```go
type artifact struct {
    RelativePath string
    Content      []byte
    Managed      bool
}
```

Build the complete list before writing. Button source comes from:

```text
registry/generated/<preset.Style>/button.gsx
```

All other components keep the current embedded path under Nova. Under Maia,
reject a resolved dependency set containing any component other than Button
before generating diffs or writes; this prevents `init --preset maia` followed
by `add alert` from creating a falsely mixed-style project. Apply `RewriteGsx`
only after selecting generated source.

- [ ] **Step 5: Implement read-only unified diff**

Use `github.com/pmezard/go-difflib/difflib`, promoting the already resolved
module to a direct dependency. Do not shell out to `git diff`. Exit
successfully when identical and print `(no changes)`. For an absent target,
diff against an empty file. `--diff` and `--overwrite` are mutually exclusive.

- [ ] **Step 6: Implement init/add behavior and update usage**

Accepted forms:

```text
gsxui init [--preset <file|code|->] [--overwrite]
gsxui add [--diff|--overwrite] <component>...
```

`runAdd` loads and validates `gsxui.preset.json` before registry resolution.
Write managed hashes only after successful content writes and GSX generation.
On generation failure, return an error; Task 7 will make the whole operation
transactional.

- [ ] **Step 7: Verify and commit**

```bash
go test ./internal/cli -run 'Test(Config|Managed|Init|Add|Diff|E2E)' -count=1
gopls check -severity=hint internal/cli/*.go
git add internal/cli
git commit -m "feat: make component install preset aware"
```

---

### Task 7: Add transactional apply and crash recovery

**Files:**
- Create: `internal/cli/transaction.go`
- Create: `internal/cli/transaction_test.go`
- Create: `internal/cli/apply.go`
- Create: `internal/cli/apply_test.go`
- Modify: `internal/cli/init.go`
- Modify: `internal/cli/add.go`
- Modify: `internal/cli/run.go`
- Modify: `internal/cli/e2e_test.go`

**CLI contract:**

```text
gsxui apply --preset <file|code|-> [--only theme|style] [--yes] [--overwrite]
```

- [ ] **Step 1: Write transaction failure/recovery tests**

Use an injected filesystem seam, not timing, to fail every rename index.
Snapshot the entire project tree before each case and assert byte-identical
restoration after:

- failure before first replacement;
- failure after preset replacement;
- failure after theme replacement;
- failure after Button replacement;
- config/hash replacement failure; and
- recovery from each journal phase on the next CLI mutation.

Temporary directories and rollback copies must be inside the project root so
renames stay on one filesystem. The journal path is
`.gsxui-transaction.json`; transaction directories are
`.gsxui-transaction-<nonce>/`.

- [ ] **Step 2: Implement a journaled replace transaction**

Journal exact relative target, staged file, backup file, previous existence,
previous hash, and phase. Fsync staged files, journal, and containing
directories at durability boundaries. Recovery always rolls back an incomplete
transaction before planning new work. Refuse symlinked target ancestors and
path escapes.

- [ ] **Step 3: Route init/add through the transaction**

Build all artifact bytes, conflicts, hashes, and external commands before
commit where possible. Stage every file, commit the filesystem transaction,
then run GSX generation. If generation fails, roll back written sources/config
and run generation once against restored sources; report both failures if
restoration generation also fails.

- [ ] **Step 4: Write failing apply behavior tests**

Pin:

- no-op apply;
- full Nova→Maia Button-only apply;
- full theme change;
- `--only theme` preserving Button bytes and style;
- `--only style` preserving theme/radius values;
- modified managed Button/theme refusal and invalid config refusal;
- `--overwrite`;
- interactive summary and declined confirmation;
- `--yes`;
- unsupported Maia when any installed styled component other than Button is
  discovered;
- Nova accepted with existing non-Button components;
- malformed preset and recovery before planning; and
- exact changed axes/file list.

Installed styled components are discovered by enumerating exact registry
targets under `cfg.UI` and checking those paths on disk, including components
installed before managed hashes existed. Managed hashes decide overwrite
ownership; they do not decide whether a component exists. Directory components
count as styled components when the registry says they have presentation. An
existing Button with no managed hash is unmanaged: a style-changing apply
refuses it unless `--overwrite`, even when its filename looks canonical.

- [ ] **Step 5: Implement apply planning**

The full apply output set is:

```text
gsxui.preset.json
<configured theme.css>
<configured ui>/button.gsx   only when installed and style changes
gsxui.json                   updated managed hashes; never self-hashed
```

`--only theme` merges incoming theme/radius into the current style.
`--only style` merges incoming style into current theme/radius. Refuse Maia
before prompting when unsupported installed components would remain Nova.

The confirmation summary must say that component source is replacement, list
modified managed conflicts, and recommend commit/stash. Do not implement
three-way merge.

- [ ] **Step 6: Verify and commit**

```bash
go test ./internal/cli -run 'Test(Transaction|Recover|Apply|Init|Add|E2E)' -count=1
gopls check -severity=hint internal/cli/*.go
git add internal/cli
git commit -m "feat: apply presets transactionally"
```

---

### Task 8: Build exact generated Button preview fixtures and iframe route

**Files:**
- Create: `site/stylepreview/nova/button.gsx`
- Generated: `site/stylepreview/nova/button.x.go`
- Create: `site/stylepreview/maia/button.gsx`
- Generated: `site/stylepreview/maia/button.x.go`
- Create: `site/stylepreview/matrix.gsx`
- Generated: `site/stylepreview/matrix.x.go`
- Create: `site/stylepreview/button_runtime_test.go`
- Create: `site/stylepreview/button_inline_bench_test.go`
- Modify: `internal/stylegen/generate.go`
- Modify: `internal/stylegen/button_test.go`
- Create: `site/pages/theme_preview.gsx`
- Generated: `site/pages/theme_preview.x.go`
- Modify: `site/pages/pages.go`
- Modify: `site/pages/pages_test.go`
- Modify: `site/pages/document.gsx`
- Generated: `site/pages/document.x.go`
- Create: `web/theme-preview.js`
- Modify: `web/main.js`

**Route:** `GET /theme/preview/button`

- [ ] **Step 1: Generate package-adjusted preview fixtures**

Extend style generation to derive both preview fixtures from the exact
committed registry artifact by changing only the package declaration from
`ui` to `nova`/`maia`. Do this with Go/GSX syntax parsing, not raw string
replacement. Check mode must include these files.

Tests compare the fixture AST with the registry AST after normalizing only the
package name. This prevents a hand-maintained preview from drifting from copied
consumer source.

- [ ] **Step 2: Write failing preview route tests**

Assert one minimal same-origin document:

- uses the canonical site head/assets and theme bootstrap;
- has `data-theme-button-preview` on body, disabling site-only Button CSS;
- renders both exact generated style packages;
- initially shows Nova and hides Maia;
- contains every variant and size;
- includes text/icon, disabled, link, focus-visible, invalid, Button Group
  adjacency, and caller-composition cases;
- includes no docs navigation/editor chrome; and
- returns 404 for extra route segments.

- [ ] **Step 3: Implement the Button matrix**

Render Nova and Maia matrices in sibling containers:

```gsx
<section data-theme-preview-style="nova">...</section>
<section data-theme-preview-style="maia" hidden>...</section>
```

Define:

```go
type ButtonRenderer func(
    variant, size, href string,
    disabled bool,
    children gsx.Node,
    attrs gsx.Attrs,
) gsx.Node
```

and one `stylepreview.Matrix(button ButtonRenderer)` component which renders
the complete matrix by calling the supplied constructor as a Go expression.
The route invokes that same matrix with `nova.Button` and `maia.Button`; do not
duplicate matrix markup, public behavior, or concrete classes.

- [ ] **Step 4: Implement the iframe message protocol**

Parent messages:

```js
{
  type: "gsxui:theme-preview:v1",
  preset: { style, radius, theme },
  mode: "light" | "dark"
}
```

Iframe replies:

```js
{ type: "gsxui:theme-preview-ready:v1" }
{ type: "gsxui:theme-preview-error:v1", message }
```

`web/theme-preview.js` must:

- require `event.origin === location.origin` and `event.source === parent`;
- validate all expected token names before mutation;
- stage variable values in a detached style declaration;
- atomically apply resolved values to `document.documentElement`;
- toggle `.dark`;
- show exactly one style section;
- preserve the last valid state on invalid input; and
- send ready after listeners are installed.

The parent will always resend current state on ready and iframe load, handling
both event orderings.

- [ ] **Step 5: Add runtime merge tests and benchmarks**

For both exact generated packages, render Buttons and assert caller classes win
through `merge.Merge`:

```text
rounded-none
h-20
px-12
bg-warning
text-warning-foreground
```

Composition classes such as `rounded-r-none` and `border-l-0` must survive when
they do not conflict. Pin anchor behavior, disabled button behavior, the bare
slot marker, and caller bare attrs.

Benchmark:

- canonical `ui.Button` with no inline defaults;
- generated Nova inline Button;
- generated Maia inline Button; and
- each inline style with conflicting caller classes.

Report `ns/op`, `B/op`, and `allocs/op`; do not add thresholds. Run with:

```bash
go test ./site/stylepreview -run '^$' -bench 'Button(Inline|CSS)' -benchmem -count=5
```

- [ ] **Step 6: Generate, verify, and commit**

```bash
make generate
go test ./internal/stylegen ./site/pages ./site/stylepreview -run 'Test(Button|ThemePreview)' -count=1
gopls check -severity=hint site/stylepreview/*/*.x.go site/pages/theme_preview.x.go
git add site/stylepreview site/pages internal/stylegen web/theme-preview.js web/main.js
git commit -m "feat: add exact Button theme preview"
```

---

### Task 9: Replace `/theme` with the complete preset editor

**Files:**
- Modify: `site/pages/theme.gsx`
- Generated: `site/pages/theme.x.go`
- Modify: `site/pages/theme_defaults_test.go`
- Create: `site/pages/theme_schema_test.go`
- Create: `web/theme-state.js`
- Create: `web/theme-state.test.js`
- Rewrite: `web/theme.js`
- Modify: `web/main.js`
- Modify: `package.json`
- Modify: `package-lock.json`
- Modify: `Makefile`

**Browser state model:**

```js
{
  resolved: canonicalPreset,
  drafts: Map<"light.primary" | "dark.primary" | "radius", string>,
  mode: "light" | "dark"
}
```

- [ ] **Step 1: Generate editor schema/default JSON from Go**

`ThemeEditor` serializes `preset.Default(StyleNova)`, ordered token/group
metadata, supported styles, and transport version into an
`<script type="application/json" data-theme-schema>` node using safe JSON
script escaping. Delete the hand-maintained `themeGroups`.

The server tests must parse that JSON and compare it with
`internal/preset.TokenNames`, both built-in defaults, and canonical JSON.

- [ ] **Step 2: Write pure state tests before controller code**

With `node:test`, pin:

- style change preserves theme/radius;
- mode change does not mutate preset;
- valid field commit updates resolved state;
- invalid draft leaves resolved state/export unchanged;
- reset restores selected style's built-in preset;
- canonical JSON/share/theme CSS round trips;
- share query load;
- JSON/CSS import atomicity;
- command strings for new and initialized projects; and
- clipboard fallback state.

Browser validation uses a hidden element:

```js
probe.style.setProperty("color", candidate)
```

and accepts the value only if the declaration remains non-empty. Validate
radius with `CSSNumericValue.parse`, accept only a non-negative `CSSUnitValue`
whose unit is the exact generated Go allowlist (or unitless zero), and reject
`CSSMathValue`, percentages, variables, and keywords. Do not accept arbitrary
`CSS.supports` values that Go will later reject.

- [ ] **Step 3: Implement `web/theme-state.js`**

Keep parsing/serialization pure and DOM-free except for injected validation
callbacks. The browser serializer follows the generated field order. It must
produce the same canonical bytes as Go for all built-in and test fixtures.
Share URL uses one `preset` query parameter containing the full share code.

Install `jsonc-parser@3.3.1` with `npm install --save-exact` and parse imported
preset JSON through `parseTree` with comments and trailing commas disabled.
Walk property nodes with a seen-key set at every object level so duplicate and
unknown keys are rejected before conversion; do not use `JSON.parse` for
imported text.

The DOM controller provides a compatible-theme-CSS parser callback implemented
with `CSSStyleSheet.replaceSync`. It walks `CSSStyleRule`/declaration objects
and accepts only normalized exact `:root` and `.dark` owners for recognized
properties. This is the browser counterpart to the authoritative Go parsers,
not regex extraction.

- [ ] **Step 4: Replace the editor UI**

The visible page must provide:

- Nova/Maia cards with Maia labeled “Button pilot”;
- light/dark tabs;
- exact text inputs for every semantic color;
- one radius input;
- inline field errors and an aggregate live status region;
- Reset;
- preset JSON copy/download/import;
- theme CSS copy/download/import;
- share-code and share-URL copy;
- `gsxui init --preset '<code>'` and
  `gsxui apply --preset '<code>'` commands; and
- a titled iframe pointing at `/theme/preview/button`.

Do not add typography, undo/redo, shuffle, viewport controls, or a catalogue
picker. Explain that Maia currently applies only to Button and that CLI refuses
unsafe mixed-style migration.

- [ ] **Step 5: Implement the DOM controller**

`web/theme.js` owns event delegation, draft lifecycle, downloads, clipboard,
import dialogs, URL initialization, status text, and iframe synchronization.
It must never mutate preview/export/commands from an invalid draft.

Clipboard failures reveal a selected manual-copy textarea and visible message.
Import errors identify JSON/CSS field or line where available and preserve all
current state. A failed iframe handshake displays a retry control; retry
reloads only the iframe.

- [ ] **Step 6: Add JS tests to repository gates**

Add:

```make
test-theme-state:
	node --test web/theme-state.test.js
```

Make `check` and `ci` depend on it. Keep `node --check` coverage for browser
modules.

- [ ] **Step 7: Verify and commit**

```bash
node --test web/theme-state.test.js
make generate
go test ./site/pages -run 'TestTheme' -count=1
gopls check -severity=hint site/pages/theme.x.go site/pages/theme_defaults_test.go
git add site/pages/theme.gsx site/pages/theme.x.go site/pages/theme_defaults_test.go \
  site/pages/theme_schema_test.go web/theme-state.js web/theme-state.test.js \
  web/theme.js web/main.js package.json package-lock.json Makefile
git commit -m "feat: complete the web preset editor"
```

---

### Task 10: Verify editor behavior and visual equivalence in Chromium

**Files:**
- Rewrite: `jstest/specs/theme-editor.spec.ts`
- Create: `jstest/specs/theme-preview.spec.ts`
- Create: `jstest/specs/button-inline.spec.ts`
- Create: `jstest/specs/theme-editor.spec.ts-snapshots/theme-nova-light-linux.png`
- Create: `jstest/specs/theme-editor.spec.ts-snapshots/theme-nova-dark-linux.png`
- Create: `jstest/specs/theme-editor.spec.ts-snapshots/theme-maia-light-linux.png`
- Create: `jstest/specs/theme-editor.spec.ts-snapshots/theme-maia-dark-linux.png`

- [ ] **Step 1: Pin the editor state contract**

Playwright must prove:

- initial Nova state and ready iframe;
- Nova→Maia changes style geometry but preserves all theme inputs;
- color/radius changes affect both styles;
- invalid focused text remains visible but does not change iframe/export/share;
- light/dark;
- reset;
- JSON and compatible CSS import/export;
- share URL reload round-trip;
- exact init/apply command;
- clipboard success and forced failure fallback; and
- iframe error/retry.

Do not assert only text/attributes; read representative computed styles inside
the iframe.

- [ ] **Step 2: Pin the complete preview matrix behavior**

Exercise keyboard focus, Enter/Space activation, enabled link navigation
prevention in the test harness, disabled Button non-activation, and Button
Group adjacency for both style sections. Confirm the hidden style section is
not focusable.

- [ ] **Step 3: Measure concrete source/recipe preview equivalence**

For each style, render a generated consumer fixture in a second harness route
and compare pinned computed properties against the editor matrix:

```text
display
height
padding-inline
border-radius
border-width/style
background-color
color
box-shadow
opacity when disabled
outline/ring under focus-visible
```

This is the guard against the editor becoming a second styling
implementation.

- [ ] **Step 4: Record four fixed visual baselines**

Use the same viewport, fonts, animation disablement, and ready markers for
Nova/Maia × light/dark. Snapshot only the iframe matrix, not the editor chrome,
so layout changes to controls do not mask Button drift.

- [ ] **Step 5: Run browser verification and commit**

```bash
npx playwright test --config jstest/playwright.config.ts \
  jstest/specs/theme-editor.spec.ts \
  jstest/specs/theme-preview.spec.ts \
  jstest/specs/button-inline.spec.ts
git add jstest/specs
git commit -m "test: verify Button style editor slice"
```

---

### Task 11: Prove CLI behavior in throwaway consumer projects

**Files:**
- Create: `internal/cli/pilot_consumer_test.go`
- Create: `internal/cli/testdata/consumer/main.gsx`
- Create: `internal/cli/testdata/consumer/main_test.go`
- Modify: `internal/cli/e2e_test.go`

- [ ] **Step 1: Create isolated Nova and Maia consumers**

Each test creates a temporary module, invokes the real CLI entry points with
the production embedded registry, supplies a local `go.work` selecting exact
GSX core, installs Button, generates, builds, and renders HTML.

Nova project:

```text
gsxui init
gsxui add button
```

Maia Button-only project:

```text
gsxui init --preset <maia-code>
gsxui add button
```

Assert source contains selected recognizable utilities, no recipe tokens,
rendered markers/behavior are correct, and `go test ./...` passes.

- [ ] **Step 2: Add ownership/migration consumers**

Pin real filesystem behavior for:

- edited Button refusal;
- `add --diff` read-only output;
- Nova→Maia full apply with `--yes`;
- theme-only preserving Button checksum;
- style-only preserving theme checksum;
- unsupported component causing Maia refusal before any write;
- overwrite replacing modified Button; and
- an injected interrupted transaction recovered on next invocation.

After every expected refusal/failure, hash the complete project tree and assert
it matches the pre-command hash.

- [ ] **Step 3: Run consumer verification**

```bash
go test ./internal/cli -run 'TestPilotConsumer' -count=1 -v
```

Expected: both clean consumers build/render and all unsafe mutations are
refused or rolled back.

- [ ] **Step 4: Commit Task 11**

```bash
git add internal/cli
git commit -m "test: prove inline Button consumer workflow"
```

---

### Task 12: Document the pilot honestly and run the review gate

**Files:**
- Modify: `README.md`
- Modify: `docs/ROADMAP.md`
- Create: `docs/reviews/2026-07-28-inline-style-button-pilot.md`
- Modify: `docs/superpowers/specs/2026-07-28-inline-style-button-pilot-design.md` only to link measured review results; do not rewrite approved decisions

- [ ] **Step 1: Document current user-facing commands**

Add concise examples for:

```bash
gsxui init --preset 'gsxui:v1:...'
gsxui add button
gsxui add --diff button
gsxui apply --preset ./gsxui.preset.json --only theme
```

State exactly:

- Button is the only inline multi-style pilot;
- Nova remains the complete catalogue style;
- Maia is Button-only;
- a copied Button owns its concrete classes;
- caller conflicts use configured `merge.Merge` at render time; and
- full style apply overwrites managed source only with explicit consent.

Do not promise other components or dates.

- [ ] **Step 2: Update the roadmap as a review-gated experiment**

Mark the completed Button slice and leave expansion as an undecided review
outcome. Do not convert the pilot into a catalogue migration commitment.

- [ ] **Step 3: Run all generation and focused checks**

```bash
make generate
make verify-generated-styles
make verify-generated
go test ./internal/preset ./internal/stylegen ./internal/cli ./ui ./site/pages -count=1
node --test web/theme-state.test.js
node --test jstest/support/compiled-css-audit.test.ts
npx playwright test --config jstest/playwright.config.ts \
  jstest/specs/theme-editor.spec.ts \
  jstest/specs/theme-preview.spec.ts \
  jstest/specs/button-inline.spec.ts
go test ./site/stylepreview -run '^$' -bench 'Button(Inline|CSS)' -benchmem -count=5
```

Capture benchmark output, generated source size, preset/share size, and browser
snapshot references in the review document.

- [ ] **Step 4: Run the authoritative clean gate**

Create a fresh temporary workspace selecting this gsxui worktree and exact
merged core, then run:

```bash
pilot_ci_workspace="$(mktemp -d)"
(cd "$pilot_ci_workspace" && go work init \
  /Users/jackieli/personal/gsxhq/gsxui/.worktrees/css-only-theme-architecture \
  /Users/jackieli/personal/gsxhq/gsx)
GOWORK="$pilot_ci_workspace/go.work" make ci
git status --short
```

Expected: `make ci` passes uncached and status contains only intended review
documentation changes.

- [ ] **Step 5: Conduct the required independent adversarial review**

The reviewer must read the approved spec and this plan, then independently:

- build Nova and Maia throwaway consumers;
- edit managed Button and theme files and attempt apply;
- inject malformed/duplicate JSON and hostile CSS imports;
- delete one recipe rule and add one unused rule;
- assemble a dynamic recipe token in a throwaway canonical file;
- inject failure at each transaction rename;
- pass conflicting Tailwind classes at runtime;
- exercise iframe synchronization before and after load; and
- inspect copied Button source for readability and private styling helpers.

Record exact commands, observed results, and any issue severity in
`docs/reviews/2026-07-28-inline-style-button-pilot.md`. Fix every accepted
finding with a focused regression test and rerun `make ci`.

- [ ] **Step 6: Answer the eight approved review questions**

The review document must answer, with evidence:

1. Is tokenized canonical GSX maintainable?
2. Are recipe files better than duplicated authored GSX?
3. Is copied concrete Button source desirable to own?
4. Is inline switch repetition preferable to private styling helpers?
5. Does the editor distinguish style from theme?
6. Are add/diff/apply ownership rules safe?
7. Is runtime Tailwind merge cost acceptable?
8. What CSS-only capability, if any, should be retained?

End with one of: `accept Button pilot`, `revise Button pilot`, or `reject
Button pilot`. Do not authorize another component in this implementation
branch.

- [ ] **Step 7: Commit documentation and review evidence**

```bash
git add README.md docs/ROADMAP.md docs/reviews \
  docs/superpowers/specs/2026-07-28-inline-style-button-pilot-design.md
git commit -m "docs: record inline Button pilot review"
```

---

## Completion Criteria

The plan is complete only when:

- `/theme` visibly and correctly edits Nova/Maia Button style, light/dark
  semantic colors, and radius through an isolated real iframe;
- the exact same recipes generate committed concrete consumer source;
- CLI init/add/diff/apply and rollback behavior is proven in throwaway
  projects;
- caller classes win through the existing configured `merge.Merge`;
- consumer `style.css` contains no Button presentation;
- generated output and all `.x.go` are drift-free;
- the browser matrix and generated consumer agree on pinned computed styles;
- `make ci` passes against GSX core
  `ef72f5eba066d7e87adf7dcadc2db62d00f22efe`; and
- the evidence-backed review pauses before any second component.

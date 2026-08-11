# Style Porter (Phase 1) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Port all 8 upstream shadcn styles into gsxui recipes with a re-runnable tool, so `gsxui add <any component>` works under any of the 8 styles and the `/theme` picker offers 8 real choices.

**Architecture:** A pure transformation `internal/stylegen/port/`: read upstream `style-<name>.css` → parse into (component, class, utilities) triples → apply four transformation rules against our `recipe.Shape` → emit `registry/styles/<style>/<component>.css`. Output is validated by the *existing* `recipe.ParseStyle` / `recipe.Conform` / `recipe.CheckConflicts` before being written, so a bad port fails loudly instead of shipping. The downstream `stylegen generate` pipeline is untouched — it already discovers styles from the filesystem.

**Tech Stack:** Go; the repo's existing `internal/recipe` (parser/shape/validator), `merge.Merge` (Tailwind-aware merger), `github.com/tdewolff/parse/v2/css` (already vendored, used by `internal/recipe/parse.go`).

**Spec:** `docs/superpowers/specs/2026-08-11-style-system-parity-design.md`
**Dossier (formats, line-cited):** `/private/tmp/claude-501/-Users-jackieli-personal-gsxhq/ed3f6834-0598-4db1-b893-3b2b28ce0688/scratchpad/porter-dossier.md` — read §1 (recipe output format), §2 (shape), §3 (validator contract), §5 (upstream input grammar), §6 (name mapping), §7 (test patterns) before Task 1.

## Global Constraints

- Upstream root: `/Users/jackieli/personal/shadcn-ui`, styles at `apps/v4/registry/styles/style-<name>.css`. Pin the commit SHA (`git -C /Users/jackieli/personal/shadcn-ui rev-parse HEAD`) into every generated header.
- The 8 styles, exactly: `vega, nova, maia, lyra, mira, luma, sera, rhea`.
- **Output must pass `recipe.ParseStyle` → `recipe.Conform` → `recipe.CheckConflicts` before being written to disk.** `Conform` requires *exact bidirectional coverage*: every `Base: true` slot emits a base rule, every declared `(slot, dimension, value)` emits a value rule, and no class exists that the shape doesn't declare.
- Recipe file format (dossier §1): one leading `/* Source: … */` block, then exactly one `@layer components { … }`; every rule body is `@apply` only; class grammar `.gsxui-recipe-<component>[-<slot>][-<dim>-<value>]` composed via `recipe.Shape.BaseClass` / `ValueClass` — never hand-concatenated.
- **Never** hand-edit `registry/generated/**` or `ui/*.gsx` — those are `stylegen generate` output.
- Inner test loop: `go test ./internal/stylegen/port -run 'TestX' -count=1`. Full `make check` once, at Task 10.
- Commit format `<type>: <summary>`.
- Any `site/examples` `.gsx` change requires `make highlight` in the same commit (repo CI gate).

---

### Task 1: Upstream reader

**Files:**
- Create: `internal/stylegen/port/upstream.go`, `internal/stylegen/port/upstream_test.go`
- Test fixture: `internal/stylegen/port/testdata/upstream-sample.css`

**Interfaces (produced):**
```go
// Rule is one upstream `.cn-*` rule.
type Rule struct {
	Class     string   // e.g. "cn-accordion-trigger" (no leading dot)
	Utilities []string // @apply tokens, in source order
	Line      int      // 1-based line of the selector, for provenance
}
// Section is one `/* MARK: <Name> */` block.
type Section struct {
	Name  string // e.g. "Radio Group"
	Rules []Rule
	Start, End int // 1-based line range, for the provenance header
}
func ReadUpstream(path string) (style string, sections []Section, err error)
```
`style` comes from the `.style-<name>` wrapper. Rules appearing before the first MARK belong to a synthetic section named `""`.

- [ ] **Step 1: Write the failing test.** Fixture `testdata/upstream-sample.css`:
```css
.style-demo {
  /* MARK: Accordion */
  .cn-accordion {
    @apply overflow-hidden rounded-2xl border;
  }

  .cn-accordion-trigger {
    @apply **:data-[slot=accordion-trigger-icon]:size-4 gap-6 p-4;
  }

  /* MARK: Radio Group */
  .cn-radio-group-item {
    @apply size-4 data-[size=sm]:size-3;
  }
}
```
Test asserts: `style == "demo"`; two sections named `"Accordion"` and `"Radio Group"`; the Accordion section has 2 rules with classes `cn-accordion`, `cn-accordion-trigger`; `cn-accordion` has exactly `["overflow-hidden","rounded-2xl","border"]`; the Radio Group section's line range covers its rule; `Line` on `cn-accordion` is 3.
- [ ] **Step 2: Run** `go test ./internal/stylegen/port -run 'TestReadUpstream' -count=1` — expect FAIL (package does not exist).
- [ ] **Step 3: Implement** using `github.com/tdewolff/parse/v2/css` the way `internal/recipe/parse.go` does (read that file first for the tokenizer idiom). MARK comments are `/* MARK: <Name> */` — capture via the comment token, trimming `MARK:` and surrounding space.
- [ ] **Step 4: Run** the same command — expect PASS.
- [ ] **Step 5: Add a real-file smoke test** `TestReadUpstreamRealMaia`: skip with `t.Skip` if `/Users/jackieli/personal/shadcn-ui` is absent; otherwise read `style-maia.css` and assert ≥55 sections and >600 rules total (dossier §5 measured 59 sections). Run it; expect PASS.
- [ ] **Step 6: Commit** `feat: upstream shadcn style-file reader`

---

### Task 2: The variant grammar

**Files:**
- Create: `internal/stylegen/port/variant.go`, `internal/stylegen/port/variant_test.go`

**Interfaces (produced):**
```go
// Classify decides what a single @apply token means for the transformation.
type Kind int
const (
	KindPlain     Kind = iota // travels verbatim onto its own rule
	KindDimension             // bare data-[dim=val]: — candidate for Rule 2
	KindSlot                  // **:data-[slot=x]: or *:data-[slot=x]: — Rule 3
)
type Classified struct {
	Kind      Kind
	Dimension string // KindDimension only
	Value     string // KindDimension only
	Slot      string // KindSlot only
	Child     bool   // KindSlot only: true for *: (direct child), false for **:
	Rest      string // the utility with the matched prefix stripped
	Raw       string // the original token
}
func Classify(token string) Classified
```

**Rules (normative — from spec §4.1 as amended, plus dossier §5.4):**
- `data-[<x>=<y>]:` with **no** `group-`/`in-`/`has-`/`not-` qualifier and not preceded by another variant → `KindDimension`. Whether `<x>` is a real dimension is decided later by the shape, not here.
- `**:data-[slot=<x>]:` → `KindSlot{Child:false}`; `*:data-[slot=<x>]:` → `KindSlot{Child:true}`.
- Everything else → `KindPlain`. That explicitly includes: bare boolean data variants (`data-open:`, `data-closed:`, `data-selected:`, `data-entering:`, …) which carry no `=value` and are runtime state, never a shape dimension; `group-data-*`, `group-has-data-*`, `in-data-*`, `has-data-*`; `aria-*`; `sm:`/`md:`/`dark:`/`supports-*`/`not-last:`/`not-data-*`; `*:[svg…]` arbitrary-selector children.
- A token may stack variants (`data-[size=default]:sm:max-w-md`). Only a **leading** `data-[x=y]:` counts as `KindDimension`; the remainder stays in `Rest` untouched.

- [ ] **Step 1: Write the failing table test** covering, at minimum:
```go
{"p-4", KindPlain, "", "", "", false, "p-4"},
{"data-[size=sm]:max-w-xs", KindDimension, "size", "sm", "", false, "max-w-xs"},
{"data-[size=default]:sm:max-w-md", KindDimension, "size", "default", "", false, "sm:max-w-md"},
{"**:data-[slot=accordion-trigger-icon]:size-4", KindSlot, "", "", "accordion-trigger-icon", false, "size-4"},
{"*:data-[slot=alert-description]:text-destructive/90", KindSlot, "", "", "alert-description", true, "text-destructive/90"},
{"data-open:bg-muted/50", KindPlain, "", "", "", false, "data-open:bg-muted/50"},
{"data-closed:animate-out", KindPlain, "", "", "", false, "data-closed:animate-out"},
{"group-data-[size=default]/alert-dialog-content:place-items-start", KindPlain, "", "", "", false, "group-data-[size=default]/alert-dialog-content:place-items-start"},
{"group-has-data-[slot=item-description]/item:translate-y-0.5", KindPlain, "", "", "", false, "group-has-data-[slot=item-description]/item:translate-y-0.5"},
{"in-data-[slot=dropdown-menu-content]:p-0", KindPlain, "", "", "", false, "in-data-[slot=dropdown-menu-content]:p-0"},
{"has-data-[icon=inline-end]:pr-1.5", KindPlain, "", "", "", false, "has-data-[icon=inline-end]:pr-1.5"},
{"dark:aria-invalid:ring-destructive/40", KindPlain, "", "", "", false, "dark:aria-invalid:ring-destructive/40"},
{"*:[svg:not([class*='size-'])]:size-8", KindPlain, "", "", "", false, "*:[svg:not([class*='size-'])]:size-8"},
{"not-last:border-b", KindPlain, "", "", "", false, "not-last:border-b"},
{"supports-backdrop-filter:backdrop-blur-xs", KindPlain, "", "", "", false, "supports-backdrop-filter:backdrop-blur-xs"},
```
Assert every field of `Classified` for each row.
- [ ] **Step 2: Run** `go test ./internal/stylegen/port -run 'TestClassify' -count=1` — expect FAIL.
- [ ] **Step 3: Implement.** Match on the token's leading variant segment. Beware: `data-[a=b]` contains `]` and `:` inside brackets — split variants on `:` only at bracket depth 0 (the repo has a precedent for this exact hazard in `lastTopLevelColon`-style scanning; write a small `splitVariants(token) []string` helper and unit-test it alongside).
- [ ] **Step 4: Run** — expect PASS.
- [ ] **Step 5: Commit** `feat: upstream variant classifier`

---

### Task 3: Name mapping and the fallback policy

**Files:**
- Create: `internal/stylegen/port/mapping.go`, `internal/stylegen/port/mapping_test.go`

**Interfaces (produced):**
```go
// SectionComponent maps an upstream MARK name to a gsxui component name.
// Returns ok=false for sections we deliberately ignore.
func SectionComponent(mark string) (component string, ok bool)
// SlotFor maps an upstream class within a component to a gsxui slot name.
// Returns ok=false when the class has no gsxui slot (reported, not dropped).
func SlotFor(component, class string) (slot string, ok bool)
// Ignored reports classes skipped by policy (e.g. React-Aria duplicates).
func Ignored(class string) bool
// StyleInvariant components take the same recipe in every style.
func StyleInvariant(component string) bool
```

**Content (from dossier §6 and spec §4.3/§4.4):**
- Default mapping is `strings.ToLower(strings.ReplaceAll(mark, " ", "-"))`. Declared overrides, each with a one-line reason comment: `"Radio Group" → "radio"` (the only name override across all 54).
- `Ignored`: any class whose name ends `-aria` (React-Aria duplicates appear both inline in a component's section and in the trailing `/* MARK: React Aria */` block — match by suffix, not by section); plus the whole `React Aria`, `Chart`, `Menu Translucent`, `Bubble`, `Attachment`, `Marker`, `Questionnaire`, `Message Scroller`, `Message` sections (components we do not have).
- `StyleInvariant`: `aspect-ratio`, `collapsible`, `spinner`, `toaster` — no upstream section exists (spec §4.4).
- Per-component slot overrides, at minimum: select's `scroll-up-button`/`scroll-down-button` → ignored (native popover needs no JS scroll buttons); dialog/alert-dialog/drawer/sheet `*-overlay` → mapped onto the content slot (native `<dialog>` fuses overlay and content; the utilities land on the content rule and the implementer adds a comment saying so).

**THE FALLBACK POLICY (load-bearing).** `Conform` demands exact coverage, but upstream does not style every slot/dimension/value we declare. Therefore: for any declared `(slot)` or `(slot, dim, value)` that receives **no** utilities from the transformation, the porter **carries over the corresponding rule from the current `registry/styles/nova/<component>.css`** (our existing hand-authored recipe) and marks it in the emitted file with a `/* carried: no upstream counterpart */` comment on the line above. This keeps behavior working, keeps `Conform` satisfied, and makes every unported rule visible in review. A carried rule that *also* has no source in today's nova recipe is a hard error — that means the shape declares something nothing has ever styled.

- [ ] **Step 1: Write the failing tests** — `SectionComponent("Radio Group") == ("radio", true)`; `SectionComponent("Native Select") == ("native-select", true)`; `SectionComponent("React Aria")` → `ok=false`; `Ignored("cn-select-value-aria") == true`; `Ignored("cn-select-value") == false`; `StyleInvariant("spinner") == true`; `StyleInvariant("button") == false`. Plus a coverage test: for all 54 components in `registry/canonical/shapes/`, either a MARK section maps to it or `StyleInvariant` is true — fail listing any component that is neither.
- [ ] **Step 2: Run** `go test ./internal/stylegen/port -run 'TestMapping' -count=1` — expect FAIL.
- [ ] **Step 3: Implement.**
- [ ] **Step 4: Run** — expect PASS. The coverage test failing is real signal: fix the table, not the test.
- [ ] **Step 5: Commit** `feat: upstream section/slot mapping with declared overrides`

---

### Task 4: The transformer

**Files:**
- Create: `internal/stylegen/port/transform.go`, `internal/stylegen/port/transform_test.go`

**Interfaces (produced):**
```go
// Ported is one component's recipe, ready to render.
type Ported struct {
	Component string
	Base      map[string][]string            // slot -> utilities
	Values    map[string]map[string]map[string][]string // slot -> dim -> value -> utilities
	Carried   map[string]bool                // class -> true when carried from nova
	SrcStart, SrcEnd int
}
// Transform applies Rules 1-3 plus the fallback policy.
// unmapped lists upstream classes/utilities with no gsxui destination.
func Transform(shape recipe.Shape, sec Section, fallback recipe.Style) (Ported, unmapped []string, err error)
```

- [ ] **Step 1: Write failing tests** for each rule in isolation, using a hand-built `recipe.Shape` (read `internal/recipe/shape.go` for the literal):
  - **Rule 1** (slot rename): `cn-accordion-item` with `["not-last:border-b"]` → `Base["item"] == ["not-last:border-b"]`.
  - **Rule 2** (dimension re-split): `cn-radio-group-item` with `["size-4","data-[size=sm]:size-3"]` against a shape declaring slot `item` with dimension `size` values `[default, sm]` → `Base["item"] == ["size-4"]` and `Values["item"]["size"]["sm"] == ["size-3"]`.
  - **Rule 2 negative**: a `data-[side=bottom]:` token where `side` is *not* a declared dimension → stays verbatim in `Base`, no error.
  - **Rule 3** (descendant decomposition): `cn-accordion-trigger` with `["**:data-[slot=accordion-trigger-icon]:size-4","p-4"]` → `Base["trigger"] == ["p-4"]` and `Base["trigger-icon"] == ["size-4"]`.
  - **Rule 3 unknown slot**: `**:data-[slot=nonesuch]:x` → appears in `unmapped`, not in any rule.
  - **Fallback**: a shape slot upstream never mentions gets the nova rule's utilities and `Carried[class] == true`.
  - **Fallback missing**: shape declares a slot that neither upstream nor nova has → `err != nil`.
- [ ] **Step 2: Run** `go test ./internal/stylegen/port -run 'TestTransform' -count=1` — expect FAIL.
- [ ] **Step 3: Implement.** Order matters: classify every token first, route `KindSlot` tokens to their target slot, route `KindDimension` tokens whose dimension is declared to their value rule, everything else stays on the rule it came from. Then apply the fallback for uncovered declared rules.
- [ ] **Step 4: Run** — expect PASS.
- [ ] **Step 5: Commit** `feat: upstream→recipe transformer with fallback policy`

---

### Task 5: Emitter + validation gate

**Files:**
- Create: `internal/stylegen/port/emit.go`, `internal/stylegen/port/emit_test.go`
- Test fixture: `internal/stylegen/port/testdata/golden-accordion-maia.css`

**Interfaces (produced):**
```go
// Render produces the recipe file bytes (header + @layer components block).
func Render(p Ported, style, upstreamSHA, upstreamPath string) ([]byte, error)
// Validate runs the repo's own three gates on rendered bytes.
func Validate(path string, shape recipe.Shape, src []byte) error
```
`Render` composes classes with `shape.BaseClass`/`ValueClass` (never string concat), emits the header in the dossier §1.1 grammar (`Source: shadcn-ui/ui@<sha>` / `<path>, <Component> lines <start>-<end>.`), one `@layer components` block, rules in shape-declaration order, `/* carried: no upstream counterpart */` above carried rules. `Validate` calls `recipe.ParseStyle`, then `recipe.Conform`, then `recipe.CheckConflicts(path, resolved, merge.Merge)` — the same trio `internal/stylegen/generate.go:122` runs.

- [ ] **Step 1: Write the failing golden test** — transform the accordion section from the real `style-maia.css` and assert `Render` output equals `testdata/golden-accordion-maia.css` byte-for-byte (follow the `internal/preset/testdata` golden pattern, dossier §7.1; support `-update` to regenerate). Plus `TestValidateRejectsUnknownClass`: hand-build bytes with a `.gsxui-recipe-accordion-bogus` rule and assert `Validate` errors.
- [ ] **Step 2: Run** `go test ./internal/stylegen/port -run 'TestRender|TestValidate' -count=1` — expect FAIL.
- [ ] **Step 3: Implement**, then generate the golden with `-update` and **read it yourself** — it is the first real evidence the port is sane. Check: header cites real line numbers; utilities look like Maia (rounded-2xl, generous padding), not like today's nova.
- [ ] **Step 4: Run** without `-update` — expect PASS.
- [ ] **Step 5: Commit** `feat: recipe emitter validated by parse/conform/conflicts`

---

### Task 6: `stylegen port` command, and port maia end-to-end

**Files:**
- Modify: `cmd/stylegen/main.go` (add the `port` subcommand)
- Create: `internal/stylegen/port/run.go`, `internal/stylegen/port/run_test.go`
- Modify: `registry/styles/maia/*.css` (regenerated output — 54 files)

**Interfaces (produced):**
```go
// Run ports one style (or all 8) from upstreamRoot into registry/styles/.
// dryRun writes nothing and returns the would-be diff summary.
func Run(repoRoot, upstreamRoot string, styles []string, dryRun bool) (Report, error)
type Report struct{ Written, Carried, Unmapped, Skipped int; UnmappedDetail []string }
```
CLI: `go run ./cmd/stylegen port --upstream <path> --style maia|all [--dry-run]`. Non-zero exit when `Unmapped > 0` unless `--allow-unmapped` (which prints the list).

- [ ] **Step 1: Write the failing test** — `TestRunDryRunMaia`: skip if upstream absent; run with `dryRun=true`; assert `Report.Written == 54`, `Unmapped == 0`, and that nothing on disk changed (compare a `registry/styles/maia/button.css` checksum before/after).
- [ ] **Step 2: Run** `go test ./internal/stylegen/port -run 'TestRun' -count=1` — expect FAIL.
- [ ] **Step 3: Implement** `Run` + the subcommand.
- [ ] **Step 4: Run** the test — expect PASS. Then port for real: `go run ./cmd/stylegen port --upstream /Users/jackieli/personal/shadcn-ui --style maia`. Resolve every unmapped report by extending the Task 3 mapping table until the run is clean.
- [ ] **Step 5: Regenerate downstream and inspect** — `go run ./cmd/stylegen generate` (or the repo's documented generate command; check `Makefile`), then `git diff --stat registry/generated/maia` should show all 54 components changed. Read `git diff registry/styles/maia/accordion.css` and confirm it now matches upstream Maia (e.g. `rounded-2xl`, `not-last:border-b`) rather than the old hand-authored spelling.
- [ ] **Step 6: Run** `go test ./internal/recipe ./internal/stylegen -count=1` — expect PASS (the existing suites validate every recipe).
- [ ] **Step 7: Commit** `feat: stylegen port subcommand; port maia from upstream`

---

### Task 7: Per-style visual gate

**Files:**
- Modify: `jstest/specs/style-visual.spec.ts`
- Modify: `jstest/harness/` route registration if the gallery is not already per-style addressable (check `site/stylepreview/<style>/` and how the harness serves it)

**Interfaces (produced):** a snapshot set per style under `jstest/specs/style-visual.spec.ts-snapshots/`, named `<style>-<card>.png`.

This lands **before** the bulk port (Task 8) so the remaining 6 styles' diffs are reviewable.

- [ ] **Step 1: Read** `jstest/specs/style-visual.spec.ts` and `site/stylepreview/` to learn how the gallery renders today and whether the style is a route parameter.
- [ ] **Step 2: Write the failing test** — parameterize the spec over `["nova","maia"]` (the two that exist now), asserting a screenshot per style per gallery card.
- [ ] **Step 3: Run** `npx playwright test -c jstest specs/style-visual.spec.ts` — expect FAIL (missing snapshots for the new names).
- [ ] **Step 4: Generate** snapshots with `--update-snapshots`, then **look at the maia ones** — this is the human review the spec's §6 risk 1 calls for. Confirm maia reads as rounded/generous versus nova's density.
- [ ] **Step 5: Run** again without the flag — expect PASS.
- [ ] **Step 6: Commit** `test: per-style visual snapshot gate`

---

### Task 8: Port the remaining 6 styles

**Files:**
- Create: `registry/styles/{vega,lyra,mira,luma,sera,rhea}/*.css` (6 × 54 = 324 files)
- Modify: `registry/styles/nova/*.css` (re-ported from upstream nova — replaces the house style, spec §2)
- Modify: `internal/preset/preset.go:24-29,53` (the `Style` enum), `internal/preset/catalog*.go` if styles are cataloged there
- Modify: `site/pages/theme_picker.gsx` (8 entries)

- [ ] **Step 1: Port them** — `go run ./cmd/stylegen port --upstream /Users/jackieli/personal/shadcn-ui --style all`. Resolve unmapped reports per style (lyra has no Carousel section — expect the fallback to carry it; confirm the report says so rather than erroring).
- [ ] **Step 2: Extend the preset enum** — add the 6 constants beside `StyleNova`/`StyleMaia`, with the upstream one-line descriptions from `apps/v4/registry/styles.tsx` (vega "Clean, neutral, and familiar"; nova "Reduced padding and margins"; maia "Rounded, with generous spacing."; lyra "Boxy and sharp. For mono fonts."; mira "Made for compact interfaces."; luma "Fluid, luminous, and soft."; sera "Editorial and typographic."; rhea "Like Luma but compact.").
- [ ] **Step 3: Wire the picker** — `site/pages/theme_picker.gsx` lists all 8 with their descriptions.
- [ ] **Step 4: Regenerate + snapshot** — run generate, then `npx playwright test -c jstest specs/style-visual.spec.ts --update-snapshots`, and **review all 8 styles' snapshots by eye**. They must look meaningfully different from each other; two styles rendering identically means the port silently fell back.
- [ ] **Step 5: Run** `go test ./... -count=1` — expect PASS.
- [ ] **Step 6: Commit** `feat: port remaining 6 upstream styles; nova now upstream-sourced`

---

### Task 9: Style-aware CLI vendoring + chart tokens + doc fix

**Files:**
- Modify: `internal/cli/add.go:155-160` and `apply.go:155-162` (style-aware for every component), `internal/cli/add.go:63-69` and `apply.go:115-127` (delete the maia gate)
- Modify: `internal/preset/preset.go` (`tokenDefinitions`: add `chart-1`..`chart-5`)
- Modify: `site/pages/theme.gsx:281` (the false "both styles render the full catalogue" copy)
- Test: `internal/cli/add_test.go`, `apply_test.go`, `internal/preset/preset_test.go`

- [ ] **Step 1: Write the failing tests** — `TestAddVendorsSelectedStyle`: for style `maia` and component `card`, the vendored `ui/card.gsx` bytes equal `registry/generated/maia/card.gsx`; repeat for `lyra`. Plus: the existing tests asserting the maia gate error must be *deleted* (they encode the old behavior) and replaced by one asserting `gsxui add card` under maia now succeeds. Plus `TestPresetHasChartTokens`: `chart-1`..`chart-5` are present with light and dark defaults.
- [ ] **Step 2: Run** `go test ./internal/cli ./internal/preset -run 'TestAdd|TestApply|TestPreset' -count=1` — expect FAIL.
- [ ] **Step 3: Implement** — replace the `if name == "button"` special case with an unconditional `registry/generated/<style>/<name>.gsx` read; delete the gate; add the chart tokens (values from upstream's `style-*.css` `:root`/`.dark` blocks — read them from the local checkout); correct the theme.gsx sentence.
- [ ] **Step 4: Run** — expect PASS.
- [ ] **Step 5: Commit** `feat: style-aware vendoring for all components, chart tokens, doc fix`

---

### Task 10: Full gate

- [ ] **Step 1:** `make check` — the authoritative gate. Expect exit 0. Fix anything it surfaces.
- [ ] **Step 2:** If any `site/examples` `.gsx` changed, run `make highlight` and include it in the same commit (repo CI rule).
- [ ] **Step 3: Commit** any fixes. Branch is then ready for the plan-level final review.

## Self-Review

- **Spec coverage:** §4.1 grammar → Tasks 2+4; §4.2 porter → Tasks 1-6; §4.3 mapping table → Task 3; §4.4 hand residue → Task 3's `StyleInvariant` + the fallback policy; §4.5 CLI trap → Task 9; §4.6 visual gate → Task 7 (deliberately before the bulk port); §4.7 doc fix → Task 9; §2 chart tokens → Task 9; §2 nova replacement → Task 8. §5 testing maps onto each task's own test step.
- **Placeholders:** none — every task names its files, test assertions, commands and commit message.
- **Type consistency:** `Rule`/`Section` (Task 1) feed `Transform` (Task 4); `Ported` feeds `Render` (Task 5); `Run`/`Report` (Task 6) wrap both. `recipe.Shape`, `recipe.ParseStyle`, `recipe.Conform`, `recipe.CheckConflicts` and `merge.Merge` are existing repo symbols, cited to their files in the dossier.
- **Known risk carried forward:** the fallback policy means a component upstream doesn't style renders identically across all 8 styles. That is intended (spec §4.4) but Task 8 Step 4's eye review must not mistake it for a porting bug — expect aspect-ratio, collapsible, spinner and toaster to be identical everywhere.

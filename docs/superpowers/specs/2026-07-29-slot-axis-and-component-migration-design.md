# Slot axis and full component migration

**Date:** 2026-07-29

**Status:** Approved design — **Stages 0 and 1 complete** (see §10)

**Extends:** [`2026-07-29-typed-recipe-model-design.md`](2026-07-29-typed-recipe-model-design.md).
That design's principles, validation model, contract artifact, layer invariant
(§9) and rejected alternatives (§8) all stand. This document specifies the slot
axis and plans the migration of the remaining catalogue.

## 1. Purpose

A component in this catalogue is rarely one styled element. It is a set of named
elements, and presentation dimensions hang off those elements individually:

| | slots | dimensions |
|---|---|---|
| Button | 1 | variant, size |
| Card | 7 | none |
| Sidebar | 38 | some, per slot |
| Whole catalogue | 291 markers over 1162 rules | 16 of 52 components have any |

`default.css` already carries per-slot dimensions such as
`[data-gsxui-slot-badge][data-variant="destructive"]`, so the model has to name
the slot to describe them at all. And the model has to reach `card.Header()` and
`sidebar.MenuButton()`, not just a component's root element.

Two thirds of the catalogue has no variant or size at all — for those components
the model buys validation, completeness and the contract, but no switch
generation.

## 2. The slot axis

```go
type Shape struct {
    Component string
    Slots     []Slot
}

type Slot struct {
    Name       string // "" is the component's root element
    Base       bool
    Dimensions []Dimension
}
```

Dimensions hang off a **slot**, not the component.

Class encoding names the component, then the slot, then the dimension and
value. The root slot contributes no segment, so a single-element component's
classes are as short as they would be without the axis:

```text
gsxui-recipe-card                        root, base
gsxui-recipe-card-header                 slot "header", base
gsxui-recipe-card-header-variant-muted   slot "header", dimension variant
gsxui-recipe-button-variant-outline      root, dimension variant
```

### 2.1 Decoding is longest-match, and that is mandatory

Sidebar declares `menu`, `menu-action`, `menu-badge`, `menu-button`,
`menu-button-tooltip` and `menu-button-tooltip-content`. A shortest-match or
first-match decode assigns `gsxui-recipe-sidebar-menu-button-tooltip` to slot
`menu` with a nonsense remainder.

`DecodeClass` therefore matches declared slot names **longest first**, then
resolves the remainder as either empty (the slot's base) or a declared
`dimension-value` pair. Both halves match against declarations, never against
dash positions — the same rule that already removes the `icon-lg` ambiguity.

An unresolvable class is an error naming the component, not a silent skip.

## 3. Slot access: generated methods

291 slots need 291 accessors. Hand-writing them is error-prone and keeps shape
and accessor in sync only by convention.

### 3.1 The bootstrap cycle, and why shapes move

Accessors are generated **into** `registry/canonical`, so the generator must be
buildable while that package is not. `cmd/stylegen` imports
`internal/stylegen`; Go compiles a package as a unit; so if `internal/stylegen`
read shapes from `registry/canonical` and `card.gsx` called a not-yet-generated
`card.Header()`, the package would fail to build, that import chain would fail
with it, and the generator binary could not be built to produce the very
methods it needs.

Shapes therefore live in a leaf package that compiles on its own:

```text
registry/canonical/shapes/         hand-written, pure data, no dependents
  card.go      var Card = recipe.Shape{Component: "card", Slots: …}
  button.go    var Button = recipe.Shape{…}
        │
        │  imported by stylegen (always compiles)
        ▼
registry/canonical/card_recipe.gen.go     GENERATED accessors
registry/canonical/card.gsx               authored structure
```

`internal/stylegen/generate.go` reads `shapes.All()`, and `internal/stylegen`
must not import `registry/canonical` in any form. Importing it for any reason
would revive the cycle through `cmd/stylegen -> stylegen -> canonical`, so the
leaf package only breaks the bootstrap while that import is absent entirely —
which the architecture test pins.

`registry/canonical` keeps its "never ships" property and its architecture test.
`shapes` is pure data and inherits the same rule.

### 3.2 Accessor naming

Slot name to method: kebab segments title-cased (`menu-button` → `MenuButton`).
The root slot's base accessor is `Root()`.

| declaration | accessor |
|---|---|
| root slot, base | `button.Root()` |
| root slot, dimension `variant` | `button.Variant(v)` |
| slot `header`, base | `card.Header()` |
| slot `menu-button`, dimension `size` | `sidebar.MenuButtonSize(v)` |

A generated accessor for a slot or dimension that the shape does not declare is
impossible by construction, and a typo at a call site is a compile error.

Accessors are generated as a per-component struct type (`buttonRecipe`) wrapping
a `recipe.Component`. `recipe.Component` itself carries only the untyped
primitives `SlotClass(slot)` and `SlotValueClass(slot, dimension, value)`; the
generated type is the typed API.

Because a method name cannot be split back into (slot, dimension) by string
manipulation — nothing in `MenuButtonSize` marks the boundary —
`stylegen.HelperCalls` takes the shape as a third argument
(`HelperCalls(filename string, src []byte, shape recipe.Shape)`) and resolves
each method name against that shape's own generated accessor names. A two-argument
form cannot work.

### 3.3 Registry name to Go identifier

A component's class expression is Go, so its receiver must be a Go
identifier. A hyphen is not legal in one, and in expression position
`input-group.Root()` is simply `input - group.Root()` — ordinary Go grammar,
not a quirk of gsx. Eleven components have kebab-case names, so the two
namespaces have to be bridged.

The registry and CSS identity stays kebab-case (`input-group`) — that is what
appears in class names, the contract, and `registry/generated/recipes.json`.
The generated Go type name and its receiver are *derived* from it by the same
title-casing rule already used for slot accessors (`menu-button` →
`MenuButton`): `input-group` gives type `inputGroupRecipe` and receiver
`inputGroup`. Nothing is renamed, and no second name is maintained.

The derivation must be validated for:

- names that title-case to a Go keyword (or otherwise are not valid Go
  identifiers), and
- two distinct components whose kebab-case names normalize to the same
  identifier.

A concurrent branch of work is implementing this generator; this section
states the rule the implementation must satisfy, not the implementation
itself.

## 4. The pipeline on the slot axis

Each unit of the pipeline is small and independently tested.

| Unit | Responsibility |
|---|---|
| `recipe.Shape` | carries `Slots []Slot`; `Base` and `Dimensions` live on a `Slot` |
| `Shape.Validate` | validates each slot; slot names are unique and non-empty except the root |
| `BaseClass`/`ValueClass` | encode a class for a given slot |
| `DecodeClass` | longest-match slot resolution (§2.1) |
| `Conform` | walks slots × dimensions × values, both directions; the error names the slot |
| `CheckConflicts` | checks each slot-and-value utility list; the error names the slot |
| `Contract` | emits `components[c].slots[s].dimensions[d]` |
| `stylegen` matcher | resolves a method name to `(slot, dimension)` against the shape |
| `recipe.Component` | exposes the untyped primitives `SlotClass`/`SlotValueClass`; the typed API is the generated per-component accessor type |
| generation | validates slot **coverage**: a shape declaring a slot, a base rule, or a dimension the component never renders is an error, reported all at once |

The contract emits each shape once: slots live under `components`, and styles
carry only utilities.

### 4.1 `internal/stylecontract` and `recipe.Shape`

`internal/stylecontract` (`Component`, `Slot`, `Axis`, `RegistryName`) already
declares Button's slot/variant/size vocabulary, independently of
`recipe.Shape`. Today the two duplicate that vocabulary with no executable
link between them — nothing checks that a `Shape` and its component's
`stylecontract.Component` agree.

The intended equivalences, as **requirements** for a later enforcement task,
not yet checked automatically:

- a shape's `Component` name equals its style-contract `RegistryName`;
- a shape's slot names map onto the contract's full slot names;
- a shape's dimensions map onto the contract's declared `data-*` axes
  (`Axis.Attribute`);
- a shape's dimension values are among the values that axis permits
  (`Axis.Values`).

No code enforces this today; it is recorded here as a gap and a future
requirement, not as a description of an existing check.

## 5. The layer hazard applies to every migration

§9 of the base design established the invariant:

> A rule overriding compiled component presentation must live in
> `@layer utilities` **and** carry specificity ≥ (0,1,0).

Migrating a component moves its rules out of `@layer components`, which silently
demotes every remaining components-layer rule that overrode it. This bug class
was found three times during Button alone, and twice it had already produced a
visible regression that a 313-test browser suite did not catch.

Every component migration therefore carries a mandatory, non-negotiable step:

1. Before: capture computed styles for every affected marker, light and dark,
   from a worktree at the pre-migration commit.
2. After: capture the same set.
3. Diff. Any difference is either fixed or explicitly justified in the commit
   message. "No test failed" is **not** evidence — the Carousel and
   InputGroupButton regressions passed every gate.

The sweep harness built during the Button work (3,283 elements × 22 properties ×
light/dark) is the tool for this and must be committed rather than rebuilt.

Making this a `make audit` gate is follow-up item 2 of the base design. It
should be built **before** the bulk migration, not after: it converts the
catalogue's riskiest property from vigilance into a build failure. **Done** —
`stylegen.CheckLayerPrecedence`, run from `make audit` as
`go run ./cmd/stylegen --check-layers`.

## 6. Migration order

Ordered by what each stage proves, not by component popularity.

**Stage 0 — the model.** Slot axis, generated accessors, shapes sub-package,
contract, matcher. Button is expressed through it (root slot only) and its
generated output is pinned byte-for-byte, proving that a root-only shape adds
nothing to what a single-element component emits.

**Stage 1 — the layer gate.** Build the `make audit` rule and commit the sweep
harness, before any bulk migration.

**Stage 2 — Card.** Multi-slot, zero dimensions. Exercises the slot axis in
isolation: 7 slots, no switch generation, small blast radius.

**Stage 3 — Badge and Alert.** Slots *with* dimensions, small. Proves per-slot
dimensions end to end.

**Stage 4 — the long tail**, in dependency order, composed components last.
Sidebar (38 slots) is last of all: it is the largest, and it composes Button.

A stage does not begin until the previous one's sweep is clean.

## 7. What this design does not do

- It does not convert marker-stamping components to compose `ui.Button`. That
  mixed composition model is follow-up item 3 of the base design and is separate
  work.
- It does not add variant-combination validation (follow-up item 4), though the
  bulk migration makes it more valuable and it may be promoted.
- It does not change theme values, presets, the CLI, or the editor.
- It does not pre-merge utilities. §8 of the base design still holds.

## 8. Risks

**The catalogue is large enough that a systematic error is expensive.** Stage 0
must land the model with Button byte-identical before anything else moves.

**The layer hazard scales with the migration.** Each migrated component makes
the next one's blast radius larger, because more rules sit in `@layer utilities`
competing on specificity alone. Stage 1 exists to bound this.

**291 accessors is a lot of generated code.** If generation is wrong it is
wrong everywhere at once. Stage 0's byte-identity check on Button is the guard.

**Composed components hide their dependencies.** Carousel renders through
`<Button>`; nothing in Carousel's own source says so. The sweep, not reading,
is what finds these.

## 9. Success criteria

1. `Shape` expresses a multi-slot component with per-slot dimensions.
2. Button, expressed on the slot axis with a root-only shape, has its
   generated output pinned byte-for-byte.
3. Slot decoding resolves `sidebar-menu-button-tooltip-content` correctly, and a
   test pins the longest-match rule against its shorter siblings.
4. Accessors are generated from shapes with no bootstrap cycle, and a typo'd
   slot or dimension is a compile error.
5. `recipes.json` groups slots under their component as a published contract
   a future editor or other consumer could enumerate, diff per axis, and
   report completeness against with no CSS access. No such consumer exists
   yet (§11.7).
6. A `make audit` rule fails the build on a components-layer rule that overrides
   a migrated component's presentation.
7. Every migrated component has a committed before/after computed-style sweep
   showing no unexplained difference.

## 10. Stage 0 and 1 complete

Shipped on this branch. Migration (Stages 2-4) is what remains.

**Stage 0 — the model.**

- `recipe.Shape` carries `Slots []Slot`; `Base`/`Dimensions` live on `Slot`.
  Validation, `Conform`, `CheckConflicts` and the error messages all name the
  slot.
- `BaseClass`/`ValueClass` take a slot. `DecodeClass` matches declared slot
  names longest-first (`TestShapeDecodeClassPrefersTheLongestSlot`).
- Shapes live in the leaf package `registry/canonical/shapes` (§3.1), and
  `internal/stylegen` does not import `registry/canonical` in any form.
- Per-slot accessors are generated into `registry/canonical/<c>_recipe.gen.go`
  as a per-component struct type. `HelperCalls` resolves method names against
  the shape, never by string splitting.
- Generation validates slot coverage in both directions: an unrendered slot,
  base rule, or dimension is an error.
- `registry/generated/recipes.json` groups rules under
  `components.<c>.slots.<s>` and carries a `version` field for consumers to
  check.
- Button is expressed through the axis (root slot only), reaching its root
  through `Root()`, and its generated output is pinned byte-for-byte.

**Stage 1 — the layer gate.**

- `stylegen.CheckLayerPrecedence` enforces both halves of the base design's §9
  invariant and runs from `make audit`. It found two latent instances on its
  first run.
- The computed-style sweep harness is committed:
  `jstest/specs/layer-precedence.spec.ts`, `jstest/support/computed-sweep.ts`,
  `jstest/support/sweep-diff.mjs`, driven by `make sweep-baseline` and
  `make sweep-compare` (3,283 elements x 22 properties x light/dark).

**Presentation pinned by the sweep:** `data-gsxui-slot-sidebar-trigger`
renders at 28px (`size-7`, its own authored value and upstream shadcn's), and
`data-gsxui-slot-input-group-button` at 28x24, radius 7px, font-size 14px.

**Not started:** Stage 2 (Card), Stage 3 (Badge, Alert), Stage 4 (the long
tail, Sidebar last). The rule in §6 still applies — a stage does not begin
until the previous one's sweep is clean.

## 11. Known gaps carried into the migration

Each is a current limitation, not a historical note. None blocks the migration.

1. **Sonner should be split into two components rather than exempted.**
   `ui/sonner.gsx` exports both `Toast` and `Toaster` and renders the toast DOM
   itself (ten `toast-*` slots), so `Toast` is a standalone component and
   registering it as one removes the naming mismatch entirely. Upstream shadcn
   has since added a separate `toast` component with the same slot vocabulary
   (`toast`, `toast-content`, `toast-title`, `toast-description`,
   `toast-action`, `toast-close`) plus `toast-portal` and `toast-viewport`,
   which gsxui lacks. Calendar needs no equivalent change: its contract is
   `RegistryName: "calendar"` with every slot `calendar-*`, and the date picker
   exists only as an example composing it, exactly as upstream does it.

   Superseded note — the original framing of this gap:
   **Sonner's slot names do not follow the registry-name prefix rule.** Every
   other component's style-contract slots are `<RegistryName>-<relative>`, which
   is what `registry/canonical/shapes/agreement_test.go` joins on. Sonner's
   registry name is `sonner` but its slots are `toaster`, `toast`, `toast-*`.
   It has no shape yet, so nothing is forced; when Sonner migrates, either its
   contract slots are renamed or the mapping gains a documented exception. The
   agreement check fails loudly rather than special-casing it silently.

1. **Upstream now ships eight styles; the model assumes two.** shadcn's
   `apps/v4/registry/styles/` carries luma, lyra, maia, mira, nova, rhea, sera
   and vega. gsxui models nova and maia. Three consequences if more are adopted:
   strict conformance means every style implements every `(slot, dimension,
   value)`, so the authoring burden is linear in styles × components; the layer
   gate performs one Tailwind compile per style and would need caching; and the
   `web/site-button.css` exemption list is already per-style (112 entries at two
   styles). Decide deliberately whether to adopt more styles before the
   migration multiplies the cost, not after.

1. **Responsive `@media` overrides are a blind spot in the layer gate.** The
   contest oracle compares an authored rule's enclosing at-rule prelude against
   the one Tailwind emits for a responsive variant, and the prelude text differs,
   so the two never match and never contest. `@media (hover: hover)` is handled
   — it is dropped on both sides, being Tailwind's `hover:` implementation — but
   width-based conditions are not. Close it by normalizing media conditions if
   the migration starts writing responsive overrides against migrated markers.

2. **The `web/site-button.css` exemption list is coupled to exact selector
   text.** 56 explicit `(selector, contested)` pairs, each carrying the
   docs/demo-fallback reason. A reflow of that file, or a semantically
   equivalent selector rewrite, fails the build as stale. That is the intended
   trade — a stale exemption must not survive silently — but expect churn while
   the components that depend on the fallback migrate.

3. **`checkAccessorNames` is unpinned at one of its three call sites.** It runs
   in `GenerateAccessors`, `Resolve` and `HelperCalls`, but only the first two
   are covered — neutering the `Resolve` call leaves the suite green.

4. **The computed-style sweep covers resting state at one viewport.** Interaction
   states are not swept. The layer gate catches such rules at source, so the
   sweep is the second net rather than the first; per-component regression pins
   in `jstest/specs/layer-precedence.spec.ts` are the pattern to follow. A static
   before/after diff of compiled CSS per marker would close this more cheaply
   than driving hover on every element.

5. **`shapes.All()` is a shallow copy.** The `Shape` structs are copied but the
   `Slots` and `Values` backing arrays are shared. No caller mutates today.

6. **From the first migration that puts a compiled component under a foundation
   mechanics rule, `make audit` and `go test ./internal/stylegen` require
   `npm install`**, because the contest oracle shells out to the Tailwind CLI.
   It errors loudly when absent rather than passing silently.

7. **`registry/generated/recipes.json` has no consumer yet.** It is a
   published contract, not one currently read by the theme editor or anything
   else — `grep -rn "recipes.json" --include='*.ts' --include='*.js'
   --include='*.go' .` finds only the writer in
   `internal/stylegen/generate.go`. There is also no version-rejection test: no
   test asserts that a consumer (real or fixture) rejects an unexpected
   `version` value. Both are gaps until an actual consumer is built.

8. **`recipe.Shape` and `internal/stylecontract.Component` duplicate the same
   slot/variant/size vocabulary with no executable link (§4.1).** The
   equivalences between them are stated as requirements, not enforced.

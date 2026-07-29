# Slot axis and full component migration

**Date:** 2026-07-29

**Status:** Approved design — **Stages 0 and 1 complete** (see §10)

**Extends:** [`2026-07-29-typed-recipe-model-design.md`](2026-07-29-typed-recipe-model-design.md).
That design's principles, validation model, contract artifact, layer invariant
(§9) and rejected alternatives (§8) all stand unchanged. This document adds the
one axis it lacks and plans the migration of the remaining catalogue.

## 1. Purpose

The typed recipe model shipped with Button. It does not generalize, because
`Shape` describes **one styled element**:

```go
type Shape struct {
    Component  string
    Base       bool
    Dimensions []Dimension
}
```

Button fit that. Nothing else in the catalogue does:

| | slots | dimensions |
|---|---|---|
| Button (migrated) | 1 | variant, size |
| Card | 7 | none |
| Sidebar | 38 | some, per slot |
| Remaining catalogue | 291 markers over 1162 rules | 16 of 52 components have any |

There is no way to express `card.Header()` or `sidebar.MenuButton()`. Two thirds
of the remaining components have no variant or size at all — for them the model
buys validation, completeness and the contract, but no switch generation.

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

Dimensions hang off a **slot**, not the component. This is required, not
cosmetic: `default.css` already carries per-slot dimensions such as
`[data-gsxui-slot-badge][data-variant="destructive"]`.

Class encoding extends the existing scheme; the root slot is unchanged, so
Button's emitted classes do not move:

```text
gsxui-recipe-card                        root, base
gsxui-recipe-card-header                 slot "header", base
gsxui-recipe-card-header-variant-muted   slot "header", dimension variant
gsxui-recipe-button-variant-outline      root, dimension variant  (unchanged)
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

Generating accessors **into** `registry/canonical` cannot work as the package
stands. `cmd/stylegen` imports `internal/stylegen`, which imported
`registry/canonical` to read `Shapes()`; Go compiles a package as a unit; so if
`card.gsx` calls a not-yet-generated `card.Header()`, the package fails to
build, that import chain fails with it, and the generator binary cannot be
built to produce the very methods it needs.

Shapes therefore move to a leaf package that compiles on its own:

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

As implemented, `internal/stylegen/generate.go` reads `shapes.All()` and
`internal/stylegen` no longer imports `registry/canonical` at all;
`canonical.Shapes()` was deleted. Leaving `stylegen` importing `canonical` for
any reason would keep the cycle alive through `cmd/stylegen -> stylegen ->
canonical`, so the leaf package only breaks the bootstrap if that import is
gone entirely — which the architecture test now pins.

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

Button's existing `Role()` is renamed to `Root()`. One line in one authored
file plus regeneration; "role" reads wrong once slots exist.

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

## 4. Consequences for the existing pipeline

Each is a localized change to code that already exists and is tested.

| Unit | Change |
|---|---|
| `recipe.Shape` | gains `Slots`; `Base`/`Dimensions` move onto `Slot` |
| `Shape.Validate` | validates each slot; slot names unique and non-empty except the root |
| `BaseClass`/`ValueClass` | take a slot |
| `DecodeClass` | longest-match slot resolution (§2.1) |
| `Conform` | walks slots × dimensions × values, both directions; error names the slot |
| `CheckConflicts` | per slot-and-value list; error names the slot |
| `Contract` | `components[c].slots[s].dimensions[d]`; **schema version 2** |
| `stylegen` matcher | method name resolves to `(slot, dimension)` via the shape, not `strings.ToLower` |
| `recipe.Component` | fixed `Role`/`Variant`/`Size` methods replaced by `SlotClass`/`SlotValueClass`, with the typed API on the generated per-component accessor type |
| generation | slot **coverage** is validated: a shape declaring a slot, a base rule, or a dimension the component never renders is an error, reported all at once |

The contract's shape-emitted-once property is preserved: slots live under
`components`, styles still carry only utilities.

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
contract v2, matcher. Button is re-expressed through it (root slot only) and
must emit byte-identical output. That byte-identity is the proof the axis is
backward compatible.

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
2. Button, re-expressed on the slot axis, emits byte-identical generated output.
3. Slot decoding resolves `sidebar-menu-button-tooltip-content` correctly, and a
   test pins the longest-match rule against its shorter siblings.
4. Accessors are generated from shapes with no bootstrap cycle, and a typo'd
   slot or dimension is a compile error.
5. `recipes.json` v2 groups slots under their component; the editor can still
   enumerate, diff per axis, and report completeness with no CSS access.
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
- Shapes moved to the leaf package `registry/canonical/shapes` (§3.1), and
  `internal/stylegen` no longer imports `registry/canonical` in any form.
- Per-slot accessors are generated into `registry/canonical/<c>_recipe.gen.go`
  as a per-component struct type. `HelperCalls` resolves method names against
  the shape, never by string splitting.
- Generation validates slot coverage in both directions: an unrendered slot,
  base rule, or dimension is an error.
- `registry/generated/recipes.json` is schema **version 2**: rules are grouped
  under `components.<c>.slots.<s>`.
- Button was re-expressed through the axis (root slot only) and emits
  byte-identical output. `Role()` was renamed `Root()`.

**Stage 1 — the layer gate.**

- `stylegen.CheckLayerPrecedence` enforces both halves of the base design's §9
  invariant and runs from `make audit`. It found two latent instances on its
  first run.
- The computed-style sweep harness is committed:
  `jstest/specs/layer-precedence.spec.ts`, `jstest/support/computed-sweep.ts`,
  `jstest/support/sweep-diff.mjs`, driven by `make sweep-baseline` and
  `make sweep-compare` (3,283 elements x 22 properties x light/dark).

**Presentation changes shipped alongside**, both justified by the sweep:
`data-gsxui-slot-sidebar-trigger` now renders at 28px (`size-7`, its own
authored value and upstream shadcn's) instead of silently losing to Button's
`size-8`; `data-gsxui-slot-input-group-button` is restored to 28x24, radius
7px, font-size 14px after the Button migration demoted its rules a layer.

**Not started:** Stage 2 (Card), Stage 3 (Badge, Alert), Stage 4 (the long
tail, Sidebar last). The rule in §6 still applies — a stage does not begin
until the previous one's sweep is clean.

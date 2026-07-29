# Typed recipe model and structural component compilation

**Date:** 2026-07-29

**Status:** Approved design

**Supersedes:** the flat-token half of
[`2026-07-28-inline-style-button-pilot-design.md`](2026-07-28-inline-style-button-pilot-design.md).
That pilot's compiler pass, provenance rules, and delivery model stand. Its
string-token vocabulary and hand-written variant switches do not.

## 1. Purpose

The inline-style pilot proved the hard part: a safe, AST-driven pass can compile
one canonical GSX component against an authored style recipe and emit concrete,
readable, recipe-free consumer source. What it did not build is a **model**.

Recipe tokens are flat strings. `gsxui-recipe-button-variant-outline` reads as
structured, but nothing parses it: the `-variant-` segment is a naming
convention. Four consequences follow, and every one of them is a missing data
model rather than a missing feature:

1. The component cannot be asked what variants it has, so no gallery.
2. Two styles can only be diffed as overlapping string sets, not per axis.
3. Completeness is one flat set comparison; it cannot say *"Maia is missing
   `size=icon-lg`"*.
4. The mapping from a public parameter to a token is hand-written as a `switch`
   in the component — duplicated per rendered element, and free to drift from
   the recipe.

This design introduces the typed model and makes the compiler generate what was
previously hand-maintained.

## 2. Principles

**The authored component is structural, not presentational.** Taken seriously,
a structural component has no correct presentation of its own, so it is not
something that ships. Every shipped artifact is `canonical × recipe`.

**The shape is an interface; a style is an implementation of it.** A component
has one set of dimensions and values. Styles supply utilities for that shape and
may not alter it.

**Prefer an existing parser to an invented one.** No new file format is
introduced. Shapes are Go values, recipes are CSS, the contract is JSON.

## 3. Architecture

```text
registry/canonical/            package canonical  (compiled, never shipped)
  button.gsx                     structure + behavior; calls role/variant/size
  button_recipe.go               var buttonShape = recipe.Shape{…}
  recipe.go                      role/variant/size helper bodies
        │
        │  imported for the shape ─────────────┐
        │  parsed as .gsx for class exprs      │
        ▼                                      │
internal/stylegen              the compiler pass
        │                                      │
        │  uses ─────────────────────────────► internal/recipe
        │                                        Shape, Dimension, Recipe
        ▼                                        CSS parse + validate
registry/styles/<style>/button.css               contract JSON schema
        │  utilities only
        ▼
registry/generated/<style>/button.gsx   per-style output, committed
        │  default style (nova) copied to ↓
ui/button.gsx                  package ui  ← ships
registry/generated/recipes.json            ← editor
```

### 3.1 `internal/recipe`

Pure data and validation. No gsx, no filesystem, no I/O. It owns `Shape` and
`Dimension`, parses a style's CSS *against* a shape, and defines the contract
JSON schema. Given a shape and a stylesheet it answers two questions: is this
recipe a complete implementation of the shape, and what are the utilities for
`(component, dimension, value)`?

### 3.2 `internal/stylegen`

The GSX pass. It parses the canonical `.gsx`, finds recipe helper calls in class
attributes, replaces each with generated source, reformats, reparses, and
verifies no recipe construct survives. It depends on `internal/recipe` and knows
nothing about CSS syntax.

### 3.3 `registry/canonical`

The authored components, their shapes, and the helper bodies. It is compiled so
that it type-checks and so that structure and behavior tests can run against the
authoritative source. It is never imported by `ui` or by consumers.

`stylegen` consumes this package two ways: it **imports** it to read the shapes
(no reading of Go declarations as data), and separately **parses** the `.gsx`
file as text to rewrite class expressions. These are different concerns and the
split is deliberate.

## 4. The recipe model

```go
// internal/recipe
type Shape struct {
    Component  string
    Base       bool // component has a role() base rule
    Dimensions []Dimension
}

type Dimension struct {
    Name    string
    Default string   // must be a member of Values
    Values  []string
}
```

Declared once, beside the canonical component whose public parameters it
mirrors:

```go
// registry/canonical/button_recipe.go
var buttonShape = recipe.Shape{
    Component: "button",
    Base:      true,
    Dimensions: []recipe.Dimension{
        {Name: "variant", Default: "default", Values: []string{
            "default", "destructive", "outline", "secondary", "ghost", "link"}},
        {Name: "size", Default: "default", Values: []string{
            "default", "xs", "sm", "lg", "icon", "icon-xs", "icon-sm", "icon-lg"}},
    },
}
```

Class names such as `.gsxui-recipe-button-variant-outline` remain the CSS
selector encoding, but they are now *derived from* the shape rather than being
the schema. The parser matches selectors against declared dimensions and values
instead of splitting on dashes, which removes the ambiguity between a dimension
name and a dashed value like `icon-lg`.

### 4.1 Style conformance

A component has one shape, and every style implements all of it. Validation runs
in both directions:

- every `(dimension, value)` in the shape has a rule in the style —
  `maia/button.css: dimension "size" missing value "icon-lg"`
- every rule in the style maps to a declared `(dimension, value)` —
  `maia/button.css: unknown recipe rule .gsxui-recipe-button-variant-plain`
- `Default` is a member of `Values`
- a base rule exists if and only if `Base` is true

Styles are therefore interchangeable: switching a style can never break a call
site, and a Nova/Maia diff is about utilities only.

### 4.2 Intra-list conflict detection

A single utility list that applies two conflicting Tailwind utilities is an
authoring mistake, and normalizing it silently would hide the bug. The build
fails instead:

```text
nova/button.css: .gsxui-recipe-button applies conflicting utilities
  rounded-lg and rounded-md
```

Conflicts *between* class parts (base against size, or either against caller
classes) are resolved at render time by `merge.Merge`. See §8.

## 5. Helpers and desugaring

The canonical authors semantic roles:

```gsx
class={
    "group/button",
    role("button"),
    variant("button", variant),
    size("button", size),
}
```

`role`, `variant`, and `size` are unexported functions in `registry/canonical`.
They have both a compile-time and a runtime meaning:

| | `role("button")` | `variant("button", v)` |
|---|---|---|
| **stylegen** | base utilities literal | a `switch` over `v` |
| **canonical runtime** | `"gsxui-recipe-button"` | token for `v`, falling back to the dimension's `Default` when `v` is empty or unknown |

The runtime fallback is not a convenience. It makes the canonical's semantics
match the generated `default:` arm, so a behavior test written against the
canonical asserts something true of every style.

Desugared output for `variant("button", variant)` under Nova:

```gsx
switch variant {
case "destructive": "bg-destructive text-contrast hover:bg-destructive/90 …"
case "outline":     "border-border bg-background hover:bg-accent …"
case "secondary":   "bg-secondary text-secondary-foreground …"
case "ghost":       "hover:bg-accent hover:text-accent-foreground …"
case "link":        "text-primary underline-offset-4 hover:underline"
default:            "bg-primary text-primary-foreground hover:bg-primary/90"
}
```

Two properties follow. The `default:` arm carries the declared default's
utilities rather than `""`, so a misspelled variant renders the default instead
of an unstyled element. And because the arms are generated from the shape, they
cannot drift from the recipe.

The canonical still writes `variant("button", variant)` once per rendered
element — the `<a>` and `<button>` branches are two call sites. Each is one line
rather than a sixteen-line switch, and neither is hand-maintained.

## 6. Rewrite mechanics

Today's resolver replaces string literals inside class parts. The new pass
replaces whole class parts: a `ClassPart` whose `Expr` is a recipe helper call
has its source span swapped for generated text.

The surrounding machinery is retained unchanged, because it is what makes the
pass trustworthy:

1. parse the canonical with `gsxparser.ParseFile`
2. collect edits as `(span, replacement)`, rejecting anything malformed at its
   own `file:line:col`
3. apply edits back-to-front, rejecting overlapping spans
4. `gen.Format`, reparse with `gsxparser.ParseFile`, then verify that no helper
   call and no `gsxui-recipe-` prefix survives

Rejection rules carry over from the current resolver: a helper call or recipe
token is illegal in a static class attribute, a non-class attribute, an
interpolation, a pipeline argument, a switch tag, and nested markup.

**Open implementation question.** Whether `ClassPart.End()` covers the trailing
comma of a part is unverified. A spike resolves this before the rewrite is built;
if it does not, the replaced span needs adjusting.

## 7. Generated contract

`registry/generated/recipes.json`, emitted by the same pass and drift-checked
alongside the generated `.gsx` files:

```json
{ "version": 1,
  "components": { "button": { "dimensions": {
      "variant": { "default": "default",
                   "values": ["default","destructive","outline","secondary","ghost","link"] },
      "size":    { "default": "default",
                   "values": ["default","xs","sm","lg","icon","icon-xs","icon-sm","icon-lg"] } } } },
  "styles": { "nova": { "button": {
      "base": ["inline-flex","shrink-0","items-center"],
      "variant": { "outline": ["border-border","bg-background","hover:bg-accent"] },
      "size":    { "xs": ["h-6","gap-1","px-2"] } } },
              "maia": {} } }
```

The shape is emitted once and shared by all styles, which is strict conformance
made visible in the artifact. The editor features this design exists to enable
follow directly: the gallery iterates `components`, a diff aligns two entries of
`styles` on identical axes, completeness is structural rather than textual, and
a preset export is a subtree.

No editor UI is in scope here. The artifact and its drift check are.

## 8. Rejected: build-time utility pre-merge

Pre-merging `base + variant + size` at generate time was considered and
rejected.

Those three are separate class parts, merged at render time together with caller
classes. A cross-part conflict — Nova's base `rounded-lg` against `size-xs`'s
`rounded-[min(var(--radius-md),10px)]` — can only be resolved statically by
emitting the **product** of the dimensions: six variants times eight sizes is 48
arms of pre-merged strings. That destroys the readability of the source the user
owns, which is the whole purpose of inline-style delivery. Trading the
deliverable's defining property for a per-render string merge is a bad trade.

Excluding the product leaves only merging within a single list, and those
conflicts are authoring mistakes that should be reported, not normalized.

The useful half is therefore kept as validation (§4.2), and cross-part merging
stays at render time where `merge.Merge` already handles it correctly.

## 9. Manifest-driven generation

The current generator hardcodes Button: a literal `buttonStyleSources` slice, a
`GenerateButton` entry point, and a `.button.gsx-*` temp file pattern.
Generation becomes discovery-driven:

```text
for each registry/canonical/<component>.gsx
  for each registry/styles/<style>/<component>.css
    emit registry/generated/<style>/<component>.gsx   package ui   (gsxui add)
    emit site/stylepreview/<style>/<component>.gsx    package <style> (editor)
  copy registry/generated/<defaultStyle>/<component>.gsx
    -> ui/<component>.gsx                             package ui   (ships)
```

`ui/` is the default style's generated output, not a fourth kind of artifact.
The default style is `nova`, named in one place and consumed by generation, by
`gsxui add` when no style is selected, and by this repository's own site. A
consumer who applies Maia overwrites `ui/<component>.gsx` from
`registry/generated/maia/`; the shape is identical, so the call sites do not
change.

The pilot remains Button-only in practice. The code stops being Button-shaped,
because doing this during the second component's migration is strictly worse
than doing it now.

Atomic write, `--check` drift mode, and the temp-file discipline are retained.

## 10. Error handling

Every diagnostic carries a source location and names the artifact at fault.

| Fault | Reported against |
|---|---|
| missing or unknown `(dimension, value)` rule | the style CSS, at the rule or file |
| `Default` not in `Values` | the shape, by component and dimension |
| conflicting utilities in one list | the style CSS, at the rule |
| helper call with an unknown component or dimension | the canonical `.gsx`, at the call |
| helper call in an illegal position | the canonical `.gsx`, at the position |
| recipe construct surviving resolution | internal invariant; fails the build |
| generated artifact drift | the artifact path, with the regeneration command |

Validation is complete before any artifact is written. A failing run mutates
nothing.

## 11. Testing

| Layer | What it proves | Where |
|---|---|---|
| Shape and validation | missing/unknown value, bad default, missing base, conflicting utilities — each with its exact message | `internal/recipe`, table tests |
| Desugaring | canonical to Nova and Maia, byte for byte | `internal/stylegen`, golden files |
| Rejection | helper in a non-class attribute, interpolation, static class, nested markup | `internal/stylegen` |
| Structure and behavior | `href` renders `<a>`, `disabled` forces `<button>`, `attrs` fall through, slot markers present — asserted once, style-independent | `registry/canonical` |
| Boundary | nothing outside `stylegen` and its tests imports `registry/canonical` | architecture test |
| Drift | generated `.gsx` and `recipes.json` match a fresh run | CI, `--check` |
| Tailwind reality | every recipe compiles with no unknown utilities | existing `compiled-css-audit`, retargeted |

The structure and behavior row is the return on compiling the canonical.
Behavior is asserted once against the structural source rather than N times
against generated styles, so the test pyramid never becomes parameterized by
presentation.

## 12. Scope

**In scope.** The typed model and validation; helper desugaring; the canonical
package and its boundary test; manifest-driven generation; the generated
contract JSON and its drift check; intra-list conflict detection; migrating
Button to the new pipeline.

**Out of scope.** Editor UI built on the contract; migrating any component other
than Button; compound variants and selection matrices; a third style; build-time
pre-merge (§8); changes to theme values, presets, or the CLI apply flow.

## 13. Success criteria

1. `registry/canonical/button.gsx` contains no concrete utility and no recipe
   token string — only `role`, `variant`, and `size` calls.
2. The shape is declared once; Nova and Maia supply utilities only.
3. Deleting a rule from either style fails the build naming the exact missing
   `(dimension, value)`.
4. Generated `ui/button.gsx` is byte-identical to a fresh run, contains readable
   Tailwind, no helper call, and no recipe prefix.
5. A misspelled variant at a call site renders the declared default.
6. `recipes.json` is sufficient to enumerate every variant, diff Nova against
   Maia per axis, and report completeness — with no access to the CSS.
7. Structure and behavior tests run once against the canonical and pass.
8. Nothing outside `stylegen` imports `registry/canonical`.

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
  button.gsx                     structure + behavior; calls button.Role/Variant/Size
  button_recipe.go               var buttonShape = recipe.Shape{…}; var button = recipe.Component{…}
  recipe.go                      shapes map and Shapes()
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

The authored components and their shapes. Each component binds its shape to a
package-level `recipe.Component` value named after the component (the helper
method bodies themselves live in `internal/recipe/component.go`). It is compiled so
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
    Base       bool // component has a Role() base rule
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
    button.Role(),
    button.Variant(variant),
    button.Size(size),
}
```

`button` is an unexported package-level `recipe.Component` value in
`registry/canonical`, and `Role`, `Variant`, and `Size` are its methods. They
have both a compile-time and a runtime meaning:

| | `button.Role()` | `button.Variant(v)` |
|---|---|---|
| **stylegen** | base utilities literal | a `switch` over `v` |
| **canonical runtime** | `"gsxui-recipe-button"` | token for `v`, falling back to the dimension's `Default` when `v` is empty or unknown |

The runtime fallback is not a convenience. It makes the canonical's semantics
match the generated `default:` arm, so a behavior test written against the
canonical asserts something true of every style.

Desugared output for `button.Variant(variant)` under Nova:

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

The canonical still writes `button.Variant(variant)` once per rendered
element — the `<a>` and `<button>` branches are two call sites. Each is one line
rather than a sixteen-line switch, and neither is hand-maintained.

### 5.1 Why methods, not free functions

The dimension names are meant to mirror the component's parameter names, so the
first design called free functions: `variant("button", variant)`. That does not
compile — a parameter named `variant` shadows the package-level function
`variant`, and Go reports "cannot call variant (variable of type string):
string is not a function." Every dimension of every component hits this
collision by construction, not by accident, because dimension and parameter
names are meant to match. Binding the shape to a method on a `recipe.Component`
value instead (`button.Variant(variant)`) makes the collision impossible: a
local variable cannot shadow a method, whatever it is named. It also states the
component name once, in the value's declaration, instead of at every call
site.

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

**Resolved.** `ClassPart.Pos()` includes the leading newline and indentation,
and `End()` excludes the trailing comma — so the edit span is
`[ExprPos, End())`, using the helper call expression's own position as the
start and the part's `End()` as the end. Implementation also found a case the
spike didn't anticipate: a helper call in a class part that also carries a
condition, pipeline, control flow, or CSS segments is rejected outright,
because `End()` extends past a trailing `: cond` guard and reusing the part's
span would silently delete it. See `recordPartEdit` and
`validateHelperPart` in `internal/stylegen/resolve.go`.

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

## 9. Compiled presentation vs. caller overrides: the layer invariant

Compiling Button's presentation into concrete utilities changes what a caller
override is competing against. Before this branch, every Button rule lived in
`@layer components` as a single low-specificity block, and any caller or
sibling-component rule placed in `@layer utilities` beat it automatically —
Tailwind's cascade orders whole `@layer`s before it ever looks at selector
specificity. After this branch, Button's own presentation arrives as plain
utility classes on the generated element, at ordinary utility specificity.
Anything that needs to override part of a Button's compiled presentation —
another component's CSS composing over a Button-rendered marker, not a
one-off caller class — is now competing against those utilities directly, and
losing by default if it is left in `@layer components`.

**The invariant:**

> A rule that overrides compiled Button presentation must live in
> `@layer utilities` **and** carry specificity >= (0,1,0).

Both halves are required, independently:

- **Layer.** The cascade orders whole layers before specificity. A
  components-layer rule against a Button-composed marker loses to Button's
  compiled utilities no matter how specific the selector is, because it never
  reaches a specificity comparison at all. It has to move to
  `@layer utilities` to even compete.
- **Specificity.** `:where(...)` is deliberately zero-specificity — it exists
  precisely so structural selector wrapping doesn't accidentally win fights
  it shouldn't. Once a rule is correctly in the utilities layer, a
  `:where()`-wrapped version of it still loses to a plain class like
  `.rounded-lg`, because within one layer the cascade falls back to ordinary
  specificity and `:where()` contributes none. The selector has to carry real
  specificity (a plain class or attribute selector, not `:where()`) to
  actually win.

Getting only one half right produces a rule that silently does nothing: it
parses, it compiles, `stylegen --check` and `make audit` are both green, and
in the running page the override never applies. This exact bug class was
found three separate times during implementation, by three different routes,
against three different components — which is why it is recorded here as an
explicit invariant rather than left implicit in the CSS.

**Worked examples**, all in `assets/css/styles/default.css`:

- **ButtonGroup corner geometry** (`@layer utilities`, near the Carousel
  arrow rules): a grouped child's corner radius and border must override the
  `rounded-lg` a Button now renders as a concrete utility. The rule lives in
  the utilities layer with a plain attribute-selector specificity, not
  wrapped in `:where()`.
- **Carousel arrows** (`@layer components` at `default.css:2610`,
  `@layer utilities` at `default.css:~3262`): geometry that doesn't compete
  with Button (`absolute`, positioning) stays in the components layer exactly
  as before; only the part that overrides Button's own compiled shape
  (`size-8 rounded-full`, against Button's `rounded-lg`) had to move to the
  utilities layer with unwrapped selectors. See the cross-reference comments
  at both locations.
- **InputGroupButton** (`@layer utilities`, `default.css:~3270`): retunes a
  grouped Button's whole size ramp — type scale and radius — which again
  means overriding utilities Button now renders concretely, not a
  low-specificity recipe class.

See item 1 in the Follow-up work section (§15) for the one case where
dropping `:where()` narrowed an existing caller-override guarantee, and item
2 for making this invariant a build-time check instead of a comment.

## 10. Manifest-driven generation

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

## 11. Error handling

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

## 12. Testing

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

## 13. Scope

**In scope.** The typed model and validation; helper desugaring; the canonical
package and its boundary test; manifest-driven generation; the generated
contract JSON and its drift check; intra-list conflict detection; migrating
Button to the new pipeline.

**Out of scope.** Editor UI built on the contract; migrating any component other
than Button; compound variants and selection matrices; a third style; build-time
pre-merge (§8); changes to theme values, presets, or the CLI apply flow.

## 14. Success criteria

1. `registry/canonical/button.gsx` contains no concrete utility and no recipe
   token string — only `button.Role()`, `button.Variant()`, and `button.Size()`
   calls.
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

## 15. Follow-up work

Working notes from implementation lived in a gitignored scratch directory
that will be deleted; the items below are captured here so they survive.

1. **Dropping `:where()` on the Carousel arrow rule is a real (untested)
   narrowing.** Promoted utilities-layer rules with real specificity now
   outrank equal-specificity caller classes — a caller passing
   `class="rounded-none"` to a Carousel arrow can no longer override it, a
   regression relative to merge base. Untested. Needs its own spec so the gap
   is known rather than assumed.
2. ~~**The components/utilities split in `default.css` is load-bearing but
   unenforced.**~~ Done: `stylegen.CheckLayerPrecedence` enforces both halves
   of the invariant and runs from `make audit` as
   `go run ./cmd/stylegen --check-layers`. It found two further latent
   instances on its first run (`data-gsxui-slot-sidebar-trigger`, in the wrong
   layer; `data-gsxui-slot-input-group-button`, right layer, zero
   specificity), both fixed in the same change.
3. **Mixed composition model.** Components that call `ui.Button` receive
   compiled utilities; components that merely stamp its `data-gsxui-slot-button`
   marker (site-only surfaces) still fall back to `web/site-button.css`.
   These two paths will keep drifting from each other until the
   marker-stampers are converted to call `ui.Button` directly.
4. **Nothing validates variant-combination completeness in recipes.** This is
   exactly how nova's `dark:hover:bg-destructive/90` went missing before this
   branch — the rule for `dark:` and the rule for `hover:` both existed, but
   their combination didn't. A style-authoring lint ("a rule declaring both
   `hover:X` and `dark:X` must also declare `dark:hover:X`") would catch this
   class of gap mechanically.
5. **`ui/button_test.go` derives its expected classes via `merge.Merge`**, so
   a regression in the merge logic itself cancels out between the expected
   and actual values and would not be caught by this test. Bounded risk: the
   real computed values are independently covered by
   `jstest/specs/style-visual.spec.ts`.
6. **`registry/canonical/recipe.go`'s `Shapes()` is a shallow copy** — the
   returned map is fresh, but each `Shape`'s `Dimensions`/`Values` slices are
   still shared with the package-level map, so a caller mutating them in
   place would corrupt shared state.
7. **`ui/button_test.go` imports `internal/stylegen` only for the
   `DefaultStyle` constant**, which pulls `registry/canonical` transitively
   into `ui`'s test binary. Moving `DefaultStyle` to `internal/recipe` would
   drop that import and tighten the canonical boundary (§13: nothing outside
   `stylegen` should depend on `registry/canonical`).
8. **Adding a dimension requires three uncoupled edits**: a method on
   `recipe.Component`, the residue list at `internal/stylegen/resolve.go:98`,
   and its duplicate at `internal/stylegen/generate_test.go:105`. Also note
   `strings.ToLower(Sel.Name)` maps a method `IconSize` to the dimension name
   `"iconsize"`, not `"icon-size"` — a naming trap for the next dimension
   added.
9. **The CSS-layer sweep covered resting state at one viewport.**
   Hover/focus/active states are covered only for the specific cases that
   `jstest/specs/style-visual.spec.ts` names, not exhaustively across
   components and breakpoints.

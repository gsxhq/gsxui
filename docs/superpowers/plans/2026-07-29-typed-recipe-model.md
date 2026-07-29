# Typed Recipe Model Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace flat recipe token strings with a typed shape model, and make the compiler generate the variant switches that are currently hand-written in the canonical component.

**Architecture:** A component's *shape* (dimensions, values, defaults) is a typed Go value declared once beside the canonical component. Style CSS files implement that shape and supply utilities only. `internal/recipe` owns the model, parsing, and validation as pure data with no gsx or I/O. `internal/stylegen` parses the canonical `.gsx`, replaces `role`/`variant`/`size` call expressions with generated literals and switches, and emits per-style sources plus a contract JSON. The canonical package compiles so it type-checks and hosts style-independent behavior tests, but never ships.

**Tech Stack:** Go 1.24+, `github.com/gsxhq/gsx` (parser/ast/gen), `github.com/tdewolff/parse/v2/css`, `github.com/gsxhq/gsxui/merge` (Tailwind-aware class merger), `encoding/json`.

**Spec:** `docs/superpowers/specs/2026-07-29-typed-recipe-model-design.md`

## Global Constraints

- Recipe class prefix is exactly `gsxui-recipe-`. It must never appear in any generated consumer artifact.
- The default style is `nova`. It is named in exactly one place (`stylegen.DefaultStyle`) and is what `ui/<component>.gsx` is generated from.
- Styles are strictly conformant: every style implements every `(dimension, value)` in the shape. No divergence, no fallback.
- Every diagnostic carries a source location and names the artifact at fault. Validation completes before any artifact is written; a failing run mutates nothing.
- Per the user's Go conventions: types, fields, and methods are unexported unless they need serialisation or cross-package use.
- Nothing outside `internal/stylegen` and its own tests may import `registry/canonical`.
- Generated files are committed and drift-checked with `go run ./cmd/stylegen --check`.
- **Every commit must leave `go build ./...` green and the test suite passing.** No task may hand the next one a broken tree.
- **`gofmt -l .` must print nothing.** The Makefile enforces `gofmt -l . | (! grep .)`; an unformatted file fails CI. Run `gofmt -w` on every file you touch before committing.
- **`make audit` must pass.** It is the first target of `make ci`. Several of its rules name `ui/button.gsx` explicitly; Task 7 updates them for the new file layout.
- Use `{/* */}` for comments inside gsx markup, never `//`.

## Resolved Spike

The spec left one mechanical question open. It is resolved — do not re-investigate:

For a `gsxast.ClassPart`, `Pos()` **includes** the leading newline and indentation, and `End()` **excludes** the trailing comma. The correct edit span for replacing a part's expression is therefore:

```go
start := fset.Position(part.ExprPos).Offset
end := fset.Position(part.End()).Offset
```

Verified against `class={ "group/button", role("button"), variant("button", variant), }`, which yields spans landing exactly on `role("button")` and `variant("button", variant)`. `gen.Format` re-indents afterwards, so replacement text needs no indentation of its own.

## File Structure

**New — `internal/recipe`** (pure model; no gsx, no filesystem):
- `shape.go` — `Shape`, `Dimension`, shape self-validation, class-name encoding/decoding
- `parse.go` — CSS parsing of a style file against a shape (adapted from today's `internal/stylegen/recipe.go`)
- `validate.go` — bidirectional conformance and intra-list conflict detection
- `contract.go` — contract JSON schema and marshalling

**New — `registry/canonical`** (compiled, never shipped):
- `recipe.go` — `role`/`variant`/`size` helper bodies and the shape registry
- `button_recipe.go` — `buttonShape`
- `button.gsx` — the canonical Button, moved from `ui/button.gsx`
- `button_test.go` — style-independent structure and behavior tests

**Modified — `internal/stylegen`**:
- `resolve.go` — desugar helper calls instead of substituting token literals
- `generate.go` — manifest-driven discovery; emit per-style, default-style copy, and contract JSON
- `recipe.go` — deleted; its parser moves to `internal/recipe/parse.go`

**Modified — elsewhere**:
- `registry/styles/{nova,maia}/button.css` — unchanged in content; validated against the shape
- `ui/button.gsx`, `ui/button.x.go` — become generated output
- `ui/button_test.go` — becomes a Nova-specific pin; behavior assertions move to canonical
- `cmd/stylegen/main.go` — calls the generalized entry point
- `internal/stylegen/architecture_test.go` — new boundary test

---

### Task 1: Shape type and self-validation

**Files:**
- Create: `internal/recipe/shape.go`
- Test: `internal/recipe/shape_test.go`

**Interfaces:**
- Consumes: nothing
- Produces: `recipe.Shape{Component string, Base bool, Dimensions []Dimension}`, `recipe.Dimension{Name, Default string, Values []string}`, `func (Shape) Validate() error`, `func (Shape) Dimension(name string) (Dimension, bool)`

- [ ] **Step 1: Write the failing test**

```go
package recipe

import "testing"

func validShape() Shape {
	return Shape{
		Component: "button",
		Base:      true,
		Dimensions: []Dimension{
			{Name: "variant", Default: "default", Values: []string{"default", "outline"}},
			{Name: "size", Default: "default", Values: []string{"default", "icon-lg"}},
		},
	}
}

func TestShapeValidateAcceptsWellFormed(t *testing.T) {
	t.Parallel()
	if err := validShape().Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}
}

func TestShapeValidateRejects(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		mutate  func(*Shape)
		wantErr string
	}{
		{"empty component", func(s *Shape) { s.Component = "" }, "component name is empty"},
		{"no dimensions", func(s *Shape) { s.Dimensions = nil }, `button: no dimensions declared`},
		{"empty dimension name", func(s *Shape) { s.Dimensions[0].Name = "" }, "button: dimension 0 has no name"},
		{"duplicate dimension", func(s *Shape) { s.Dimensions[1].Name = "variant" }, `button: duplicate dimension "variant"`},
		{"no values", func(s *Shape) { s.Dimensions[0].Values = nil }, `button: dimension "variant" declares no values`},
		{"duplicate value", func(s *Shape) { s.Dimensions[0].Values = []string{"outline", "outline"} }, `button: dimension "variant" duplicates value "outline"`},
		{"default not a value", func(s *Shape) { s.Dimensions[0].Default = "ghost" }, `button: dimension "variant" default "ghost" is not one of its values`},
		{"empty value", func(s *Shape) { s.Dimensions[0].Values = []string{"default", ""} }, `button: dimension "variant" has an empty value`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := validShape()
			tt.mutate(&s)
			err := s.Validate()
			if err == nil {
				t.Fatalf("Validate() = nil, want error containing %q", tt.wantErr)
			}
			if got := err.Error(); got != tt.wantErr {
				t.Errorf("Validate() = %q, want %q", got, tt.wantErr)
			}
		})
	}
}

func TestShapeDimensionLookup(t *testing.T) {
	t.Parallel()
	s := validShape()
	d, ok := s.Dimension("size")
	if !ok {
		t.Fatal("Dimension(size) = false, want true")
	}
	if d.Default != "default" {
		t.Errorf("Default = %q, want %q", d.Default, "default")
	}
	if _, ok := s.Dimension("tone"); ok {
		t.Error("Dimension(tone) = true, want false")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/recipe/ -run TestShape -v`
Expected: FAIL — package `internal/recipe` does not exist.

- [ ] **Step 3: Write minimal implementation**

```go
// Package recipe models component style recipes: the shape a component
// declares, and the utilities a style supplies for it.
package recipe

import "fmt"

// Shape is a component's style interface: the dimensions it varies over and
// the values each dimension admits. Every style must implement all of it.
type Shape struct {
	Component  string
	Base       bool // component has a role() base rule
	Dimensions []Dimension
}

// Dimension is one axis of a Shape. Default must be a member of Values and is
// what an empty or unrecognized value resolves to.
type Dimension struct {
	Name    string
	Default string
	Values  []string
}

func (s Shape) Validate() error {
	if s.Component == "" {
		return fmt.Errorf("component name is empty")
	}
	if len(s.Dimensions) == 0 {
		return fmt.Errorf("%s: no dimensions declared", s.Component)
	}
	seen := make(map[string]struct{}, len(s.Dimensions))
	for i, dimension := range s.Dimensions {
		if dimension.Name == "" {
			return fmt.Errorf("%s: dimension %d has no name", s.Component, i)
		}
		if _, exists := seen[dimension.Name]; exists {
			return fmt.Errorf("%s: duplicate dimension %q", s.Component, dimension.Name)
		}
		seen[dimension.Name] = struct{}{}
		if err := dimension.validate(s.Component); err != nil {
			return err
		}
	}
	return nil
}

func (d Dimension) validate(component string) error {
	if len(d.Values) == 0 {
		return fmt.Errorf("%s: dimension %q declares no values", component, d.Name)
	}
	seen := make(map[string]struct{}, len(d.Values))
	for _, value := range d.Values {
		if value == "" {
			return fmt.Errorf("%s: dimension %q has an empty value", component, d.Name)
		}
		if _, exists := seen[value]; exists {
			return fmt.Errorf("%s: dimension %q duplicates value %q", component, d.Name, value)
		}
		seen[value] = struct{}{}
	}
	if !d.Has(d.Default) {
		return fmt.Errorf("%s: dimension %q default %q is not one of its values",
			component, d.Name, d.Default)
	}
	return nil
}

// Has reports whether value is declared by this dimension.
func (d Dimension) Has(value string) bool {
	for _, candidate := range d.Values {
		if candidate == value {
			return true
		}
	}
	return false
}

// Dimension returns the named dimension.
func (s Shape) Dimension(name string) (Dimension, bool) {
	for _, dimension := range s.Dimensions {
		if dimension.Name == name {
			return dimension, true
		}
	}
	return Dimension{}, false
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/recipe/ -run TestShape -v`
Expected: PASS, all subtests.

- [ ] **Step 5: Commit**

```bash
git add internal/recipe/shape.go internal/recipe/shape_test.go
git commit -m "feat: add typed recipe shape model"
```

---

### Task 2: Class-name encoding and decoding

**Files:**
- Modify: `internal/recipe/shape.go`
- Test: `internal/recipe/shape_test.go`

**Interfaces:**
- Consumes: `Shape`, `Dimension` from Task 1
- Produces: `const Prefix = "gsxui-recipe-"`, `func (Shape) BaseClass() string`, `func (Shape) ValueClass(dimension, value string) string`, `func (Shape) DecodeClass(class string) (dimension, value string, kind ClassKind, err error)` with `ClassKind` one of `ClassBase`, `ClassValue`

Class names are a derived encoding, not the schema. Decoding matches against *declared* values, which is what removes the ambiguity between a dimension name and a dashed value like `icon-lg`.

- [ ] **Step 1: Write the failing test**

```go
func TestShapeClassEncoding(t *testing.T) {
	t.Parallel()
	s := validShape()
	if got, want := s.BaseClass(), "gsxui-recipe-button"; got != want {
		t.Errorf("BaseClass() = %q, want %q", got, want)
	}
	if got, want := s.ValueClass("size", "icon-lg"), "gsxui-recipe-button-size-icon-lg"; got != want {
		t.Errorf("ValueClass() = %q, want %q", got, want)
	}
}

func TestShapeDecodeClass(t *testing.T) {
	t.Parallel()
	s := validShape()
	tests := []struct {
		class         string
		wantKind      ClassKind
		wantDimension string
		wantValue     string
	}{
		{"gsxui-recipe-button", ClassBase, "", ""},
		{"gsxui-recipe-button-variant-outline", ClassValue, "variant", "outline"},
		// The dashed value must win over any shorter accidental split.
		{"gsxui-recipe-button-size-icon-lg", ClassValue, "size", "icon-lg"},
	}
	for _, tt := range tests {
		t.Run(tt.class, func(t *testing.T) {
			t.Parallel()
			dimension, value, kind, err := s.DecodeClass(tt.class)
			if err != nil {
				t.Fatalf("DecodeClass() error = %v", err)
			}
			if kind != tt.wantKind || dimension != tt.wantDimension || value != tt.wantValue {
				t.Errorf("DecodeClass() = (%q, %q, %v), want (%q, %q, %v)",
					dimension, value, kind, tt.wantDimension, tt.wantValue, tt.wantKind)
			}
		})
	}
}

func TestShapeDecodeClassRejects(t *testing.T) {
	t.Parallel()
	s := validShape()
	for _, class := range []string{
		"gsxui-recipe-card",                    // wrong component
		"gsxui-recipe-button-variant-plain",    // undeclared value
		"gsxui-recipe-button-tone-quiet",       // undeclared dimension
		"inline-flex",                          // not a recipe class
		"gsxui-recipe-",                        // prefix only
	} {
		t.Run(class, func(t *testing.T) {
			t.Parallel()
			if _, _, _, err := s.DecodeClass(class); err == nil {
				t.Errorf("DecodeClass(%q) = nil error, want error", class)
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/recipe/ -run TestShapeDecode -v`
Expected: FAIL — `undefined: ClassKind`.

- [ ] **Step 3: Write minimal implementation**

Append to `internal/recipe/shape.go`:

```go
import "strings" // add to the existing import block

// Prefix namespaces every recipe class. It must never survive into a
// generated consumer artifact.
const Prefix = "gsxui-recipe-"

// ClassKind distinguishes a component's base rule from a dimension value rule.
type ClassKind uint8

const (
	ClassBase ClassKind = iota
	ClassValue
)

func (s Shape) BaseClass() string { return Prefix + s.Component }

func (s Shape) ValueClass(dimension, value string) string {
	return Prefix + s.Component + "-" + dimension + "-" + value
}

// DecodeClass resolves a recipe class name against the shape. It matches
// declared dimensions and values rather than splitting on dashes, so a dashed
// value such as "icon-lg" is unambiguous.
func (s Shape) DecodeClass(class string) (dimension, value string, kind ClassKind, err error) {
	if !strings.HasPrefix(class, Prefix) {
		return "", "", 0, fmt.Errorf("%q is not a recipe class", class)
	}
	if class == s.BaseClass() {
		return "", "", ClassBase, nil
	}
	rest, ok := strings.CutPrefix(class, s.BaseClass()+"-")
	if !ok {
		return "", "", 0, fmt.Errorf("recipe class %q does not belong to component %q", class, s.Component)
	}
	for _, candidate := range s.Dimensions {
		suffix, ok := strings.CutPrefix(rest, candidate.Name+"-")
		if !ok {
			continue
		}
		if !candidate.Has(suffix) {
			return "", "", 0, fmt.Errorf("recipe class %q: dimension %q does not declare value %q",
				class, candidate.Name, suffix)
		}
		return candidate.Name, suffix, ClassValue, nil
	}
	return "", "", 0, fmt.Errorf("recipe class %q names no declared dimension of %q", class, s.Component)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/recipe/ -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/recipe/shape.go internal/recipe/shape_test.go
git commit -m "feat: encode and decode recipe class names against a shape"
```

---

### Task 2b: Move the CSS parser into internal/recipe

**Files:**
- Create: `internal/recipe/parse.go` (moved from `internal/stylegen/recipe.go`)
- Create: `internal/recipe/parse_test.go` (moved from `internal/stylegen/recipe_test.go`)
- Create: `internal/recipe/testdata/recipe-valid.css` (moved from `internal/stylegen/testdata/`)
- Delete: `internal/stylegen/recipe.go`, `internal/stylegen/recipe_test.go`, `internal/stylegen/testdata/recipe-valid.css`

**Interfaces:**
- Consumes: `Prefix` from Task 2
- Produces: `func ParseStyle(filename string, src []byte) (Style, error)`, `type Style` with `func (Style) Rules() []Rule`, `func (Style) Classes() []string`, `func (Style) Lookup(class string) (Rule, bool)`, `type Rule{Class string; Utilities []string}`

This is a **pure move with renames**. The parser's behavior is already correct and well tested; do not redesign it. The renames are: `ParseRecipes` → `ParseStyle`, `Recipes` → `Style`, `Recipe` → `Rule`, `Recipe.Token` → `Rule.Class`, `Recipes.Tokens()` → `Style.Classes()`, `RecipePrefix` → `Prefix` (already defined in Task 2, so delete the old constant).

One addition, not a rename: the existing type keeps an unexported `ordered []Rule` field but exposes only class names. `Conform` needs the rules themselves, so add:

```go
// Rules returns every parsed rule in source order.
func (s Style) Rules() []Rule {
	out := make([]Rule, len(s.ordered))
	for i, rule := range s.ordered {
		out[i] = Rule{Class: rule.Class, Utilities: slices.Clone(rule.Utilities)}
	}
	return out
}
```

- [ ] **Step 1: Move the files and apply renames**

```bash
git mv internal/stylegen/recipe.go internal/recipe/parse.go
git mv internal/stylegen/recipe_test.go internal/recipe/parse_test.go
git mv internal/stylegen/testdata/recipe-valid.css internal/recipe/testdata/recipe-valid.css
```

Then in both moved files: change `package stylegen` to `package recipe`; apply the renames listed above; delete the `const RecipePrefix = "gsxui-recipe-"` line (Task 2 defines `Prefix`); rename the error helpers `recipeError`/`parserRecipeError` to `styleError`/`parserStyleError` for consistency.

- [ ] **Step 2: Update stylegen's call sites in the same commit**

**Every commit must leave `go build ./...` green.** The move breaks
`internal/stylegen`, so repair it here rather than leaving the tree broken
until Task 6. These are mechanical substitutions in `internal/stylegen/resolve.go`
and `generate.go`:

- add the import `"github.com/gsxhq/gsxui/internal/recipe"`
- `ParseRecipes(` → `recipe.ParseStyle(`
- `Recipes` (the type) → `recipe.Style`
- `RecipePrefix` → `recipe.Prefix`
- `Recipe` (the type) → `recipe.Rule`
- `.Token` on a rule → `.Class`
- `.Tokens()` on a style → `.Classes()`

Do not change any behavior. Task 6 rewrites this logic; here it only has to
compile and keep passing its existing tests.

- [ ] **Step 3: Run tests to verify the move is clean**

Run: `go build ./... && go test ./internal/recipe/ ./internal/stylegen/ -v`
Expected: both packages PASS, build green.

- [ ] **Step 4: Commit**

```bash
git add -A internal/recipe internal/stylegen
git commit -m "refactor: move the recipe CSS parser into internal/recipe"
```

---

### Task 3: Conformance validation

**Files:**
- Create: `internal/recipe/validate.go`
- Test: `internal/recipe/validate_test.go`

**Interfaces:**
- Consumes: `Shape` (Task 1), `DecodeClass` (Task 2), `Style`/`Rule` (Task 2b)
- Produces: `type Resolved struct { Shape Shape; Base []string; Values map[string]map[string][]string }`, `func Conform(filename string, shape Shape, style Style) (Resolved, error)`, `func (Resolved) Utilities(dimension, value string) []string`

`Values` is keyed `[dimension][value]`. `Conform` validates in both directions and returns the resolved recipe only if the style is a complete, exact implementation of the shape.

- [ ] **Step 1: Write the failing test**

```go
package recipe

import (
	"reflect"
	"strings"
	"testing"
)

func conformShape() Shape {
	return Shape{
		Component: "button",
		Base:      true,
		Dimensions: []Dimension{
			{Name: "variant", Default: "default", Values: []string{"default", "outline"}},
		},
	}
}

const conformCSS = `@layer components {
  .gsxui-recipe-button { @apply inline-flex items-center; }
  .gsxui-recipe-button-variant-default { @apply bg-primary; }
  .gsxui-recipe-button-variant-outline { @apply border-border bg-background; }
}`

func mustParse(t *testing.T, src string) Style {
	t.Helper()
	style, err := ParseStyle("nova/button.css", []byte(src))
	if err != nil {
		t.Fatalf("ParseStyle() error = %v", err)
	}
	return style
}

func TestConformAcceptsCompleteStyle(t *testing.T) {
	t.Parallel()
	resolved, err := Conform("nova/button.css", conformShape(), mustParse(t, conformCSS))
	if err != nil {
		t.Fatalf("Conform() error = %v", err)
	}
	if got, want := resolved.Base, []string{"inline-flex", "items-center"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Base = %q, want %q", got, want)
	}
	if got, want := resolved.Utilities("variant", "outline"), []string{"border-border", "bg-background"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Utilities() = %q, want %q", got, want)
	}
}

func TestConformRejectsMissingValue(t *testing.T) {
	t.Parallel()
	src := strings.Replace(conformCSS,
		"  .gsxui-recipe-button-variant-outline { @apply border-border bg-background; }\n", "", 1)
	_, err := Conform("maia/button.css", conformShape(), mustParse(t, src))
	if err == nil {
		t.Fatal("Conform() = nil error, want error")
	}
	want := `maia/button.css: dimension "variant" missing value "outline"`
	if got := err.Error(); got != want {
		t.Errorf("Conform() = %q, want %q", got, want)
	}
}

func TestConformRejectsMissingBase(t *testing.T) {
	t.Parallel()
	src := strings.Replace(conformCSS,
		"  .gsxui-recipe-button { @apply inline-flex items-center; }\n", "", 1)
	_, err := Conform("maia/button.css", conformShape(), mustParse(t, src))
	if err == nil {
		t.Fatal("Conform() = nil error, want error")
	}
	want := `maia/button.css: missing base rule .gsxui-recipe-button`
	if got := err.Error(); got != want {
		t.Errorf("Conform() = %q, want %q", got, want)
	}
}

func TestConformRejectsUnknownRule(t *testing.T) {
	t.Parallel()
	src := strings.Replace(conformCSS, "}",
		"  .gsxui-recipe-button-variant-plain { @apply bg-muted; }\n}", 1)
	_, err := Conform("maia/button.css", conformShape(), mustParse(t, src))
	if err == nil {
		t.Fatal("Conform() = nil error, want error")
	}
	if got := err.Error(); !strings.Contains(got, `does not declare value "plain"`) {
		t.Errorf("Conform() = %q, want a message about the undeclared value", got)
	}
}

func TestConformRejectsInvalidShape(t *testing.T) {
	t.Parallel()
	shape := conformShape()
	shape.Dimensions[0].Default = "ghost"
	_, err := Conform("nova/button.css", shape, mustParse(t, conformCSS))
	if err == nil {
		t.Fatal("Conform() = nil error, want error")
	}
	if got, want := err.Error(), `button: dimension "variant" default "ghost" is not one of its values`; got != want {
		t.Errorf("Conform() = %q, want %q", got, want)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/recipe/ -run TestConform -v`
Expected: FAIL — `undefined: Conform`.

- [ ] **Step 3: Write minimal implementation**

```go
package recipe

import (
	"fmt"
	"slices"
)

// Resolved is a style proven to be a complete implementation of a shape.
type Resolved struct {
	Shape  Shape
	Base   []string
	Values map[string]map[string][]string // [dimension][value]
}

// Utilities returns the utilities a resolved recipe supplies for one value.
func (r Resolved) Utilities(dimension, value string) []string {
	return slices.Clone(r.Values[dimension][value])
}

// Conform checks that style implements shape exactly: every declared
// (dimension, value) has a rule, and every rule maps to a declared one.
func Conform(filename string, shape Shape, style Style) (Resolved, error) {
	if err := shape.Validate(); err != nil {
		return Resolved{}, err
	}

	resolved := Resolved{Shape: shape, Values: make(map[string]map[string][]string, len(shape.Dimensions))}
	for _, dimension := range shape.Dimensions {
		resolved.Values[dimension.Name] = make(map[string][]string, len(dimension.Values))
	}

	// Style to shape: every rule must be declared.
	for _, rule := range style.Rules() {
		dimension, value, kind, err := shape.DecodeClass(rule.Class)
		if err != nil {
			return Resolved{}, fmt.Errorf("%s: %w", filename, err)
		}
		if kind == ClassBase {
			if !shape.Base {
				return Resolved{}, fmt.Errorf("%s: component %q declares no base rule, found %s",
					filename, shape.Component, rule.Class)
			}
			resolved.Base = slices.Clone(rule.Utilities)
			continue
		}
		resolved.Values[dimension][value] = slices.Clone(rule.Utilities)
	}

	// Shape to style: every declaration must be supplied.
	if shape.Base && resolved.Base == nil {
		return Resolved{}, fmt.Errorf("%s: missing base rule .%s", filename, shape.BaseClass())
	}
	for _, dimension := range shape.Dimensions {
		for _, value := range dimension.Values {
			if _, ok := resolved.Values[dimension.Name][value]; !ok {
				return Resolved{}, fmt.Errorf("%s: dimension %q missing value %q",
					filename, dimension.Name, value)
			}
		}
	}
	return resolved, nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/recipe/ -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/recipe/validate.go internal/recipe/validate_test.go
git commit -m "feat: validate style conformance against a component shape"
```

---

### Task 4: Intra-list conflict detection

**Files:**
- Modify: `internal/recipe/validate.go`
- Test: `internal/recipe/validate_test.go`

**Interfaces:**
- Consumes: `Resolved` (Task 3)
- Produces: `func CheckConflicts(filename string, resolved Resolved, merge func([]string) string) error`

The Tailwind conflict model is injected as a function so `internal/recipe` stays free of a dependency on `merge`, and so the test can drive it with a stub. `merge.Merge` drops superseded classes, so a list that gets shorter when merged contained a conflict. This is the real Tailwind-aware model, not a heuristic.

- [ ] **Step 1: Write the failing test**

```go
func TestCheckConflictsAcceptsCleanLists(t *testing.T) {
	t.Parallel()
	resolved, err := Conform("nova/button.css", conformShape(), mustParse(t, conformCSS))
	if err != nil {
		t.Fatal(err)
	}
	if err := CheckConflicts("nova/button.css", resolved, noConflictMerger); err != nil {
		t.Fatalf("CheckConflicts() = %v, want nil", err)
	}
}

// noConflictMerger keeps every class, standing in for a merge with no conflicts.
func noConflictMerger(classes []string) string { return strings.Join(classes, " ") }

func TestCheckConflictsRejectsSupersededUtility(t *testing.T) {
	t.Parallel()
	src := strings.Replace(conformCSS,
		"@apply inline-flex items-center;", "@apply rounded-lg rounded-md;", 1)
	resolved, err := Conform("nova/button.css", conformShape(), mustParse(t, src))
	if err != nil {
		t.Fatal(err)
	}
	// Stub merger: drops "rounded-lg" when "rounded-md" follows it.
	stub := func(classes []string) string {
		var kept []string
		for _, class := range classes {
			if class == "rounded-lg" && slices.Contains(classes, "rounded-md") {
				continue
			}
			kept = append(kept, class)
		}
		return strings.Join(kept, " ")
	}
	err = CheckConflicts("nova/button.css", resolved, stub)
	if err == nil {
		t.Fatal("CheckConflicts() = nil, want error")
	}
	want := "nova/button.css: .gsxui-recipe-button applies conflicting utilities: rounded-lg is superseded"
	if got := err.Error(); got != want {
		t.Errorf("CheckConflicts() = %q, want %q", got, want)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/recipe/ -run TestCheckConflicts -v`
Expected: FAIL — `undefined: CheckConflicts`.

- [ ] **Step 3: Write minimal implementation**

Append to `internal/recipe/validate.go`:

```go
import "strings" // add to the existing import block

// CheckConflicts reports a utility list that contains a Tailwind conflict.
// merger is the Tailwind-aware class merger; a list that shortens when merged
// contained a utility superseded by a later one, which is an authoring error
// rather than something to normalize silently.
func CheckConflicts(filename string, resolved Resolved, merger func([]string) string) error {
	check := func(class string, utilities []string) error {
		kept := strings.Fields(merger(slices.Clone(utilities)))
		if len(kept) == len(utilities) {
			return nil
		}
		keptSet := make(map[string]struct{}, len(kept))
		for _, utility := range kept {
			keptSet[utility] = struct{}{}
		}
		for _, utility := range utilities {
			if _, ok := keptSet[utility]; !ok {
				return fmt.Errorf("%s: .%s applies conflicting utilities: %s is superseded",
					filename, class, utility)
			}
		}
		return nil
	}

	if resolved.Shape.Base {
		if err := check(resolved.Shape.BaseClass(), resolved.Base); err != nil {
			return err
		}
	}
	for _, dimension := range resolved.Shape.Dimensions {
		for _, value := range dimension.Values {
			class := resolved.Shape.ValueClass(dimension.Name, value)
			if err := check(class, resolved.Values[dimension.Name][value]); err != nil {
				return err
			}
		}
	}
	return nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/recipe/ -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/recipe/validate.go internal/recipe/validate_test.go
git commit -m "feat: reject conflicting utilities within one recipe rule"
```

---

### Task 5: The canonical package

**Files:**
- Create: `registry/canonical/recipe.go`
- Create: `registry/canonical/button_recipe.go`
- Create: `registry/canonical/button.gsx` (moved from `ui/button.gsx`, rewritten)
- Create: `registry/canonical/button_test.go`
- Test: `registry/canonical/recipe_test.go`

**Interfaces:**
- Consumes: `recipe.Shape`, `recipe.Dimension` (Task 1)
- Produces: `recipe.Component` (in `internal/recipe/component.go`) with `Role() string`, `Variant(value string) string`, `Size(value string) string`; `func Shapes() map[string]recipe.Shape` (exported for `stylegen`); and one unexported package-level `recipe.Component` var per component in `registry/canonical`, **named exactly after the component**

**Why methods, not free functions.** A component parameter named `variant`
shadows any package-level function named `variant`, so `variant("button", variant)`
does not compile — and dimension names are *supposed* to match parameter names,
so this would bite every dimension of every component. A method on a
package-level value cannot be shadowed by a local, so the collision is
impossible by construction rather than avoided by naming convention. It also
states the component name once instead of at every call site.

The var name is the component name — that is the contract stylegen resolves
against:

```go
// registry/canonical/button_recipe.go
var button = recipe.Component{Shape: buttonShape}
```

```gsx
class={
	"group/button",
	button.Role(),
	button.Variant(variant),
	button.Size(size),
}
```

Method name to dimension: `Role()` is the base rule; any other method maps to
the dimension whose name is the method name lowercased (`Variant` → `variant`,
`Size` → `size`). Adding a dimension later means adding one method to
`recipe.Component`, which keeps the mapping type-checked.

The helpers resolve an empty or unrecognized value to the dimension's declared default, so the canonical's runtime semantics match the generated `default:` arm.

- [ ] **Step 1: Write the failing test for the helpers**

`registry/canonical/recipe_test.go`:

```go
package canonical

import "testing"

func TestRoleClass(t *testing.T) {
	t.Parallel()
	if got, want := button.Role(), "gsxui-recipe-button"; got != want {
		t.Errorf("Role() = %q, want %q", got, want)
	}
}

func TestDimensionHelpersResolveDeclaredValues(t *testing.T) {
	t.Parallel()
	if got, want := button.Variant("outline"), "gsxui-recipe-button-variant-outline"; got != want {
		t.Errorf("Variant() = %q, want %q", got, want)
	}
	if got, want := button.Size("icon-lg"), "gsxui-recipe-button-size-icon-lg"; got != want {
		t.Errorf("Size() = %q, want %q", got, want)
	}
}

func TestDimensionHelpersFallBackToDefault(t *testing.T) {
	t.Parallel()
	// Empty and unrecognized values must both resolve to the declared default,
	// matching the generated switch's default arm.
	for _, value := range []string{"", "destructve"} {
		if got, want := button.Variant(value), "gsxui-recipe-button-variant-default"; got != want {
			t.Errorf("Variant(%q) = %q, want %q", value, got, want)
		}
	}
}

func TestShapesAreValid(t *testing.T) {
	t.Parallel()
	shapes := Shapes()
	if len(shapes) == 0 {
		t.Fatal("Shapes() is empty")
	}
	for name, shape := range shapes {
		if err := shape.Validate(); err != nil {
			t.Errorf("Shapes()[%q].Validate() = %v", name, err)
		}
		if shape.Component != name {
			t.Errorf("Shapes()[%q].Component = %q", name, shape.Component)
		}
	}
}

func TestButtonShapeMatchesPublicAPI(t *testing.T) {
	t.Parallel()
	shape := Shapes()["button"]
	variantDim, _ := shape.Dimension("variant")
	want := []string{"default", "destructive", "outline", "secondary", "ghost", "link"}
	if len(variantDim.Values) != len(want) {
		t.Fatalf("variant values = %q, want %q", variantDim.Values, want)
	}
	for i, value := range want {
		if variantDim.Values[i] != value {
			t.Errorf("variant value %d = %q, want %q", i, variantDim.Values[i], value)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./registry/canonical/ -v`
Expected: FAIL — package does not exist.

- [ ] **Step 3: Write the shape and helpers**

`registry/canonical/button_recipe.go`:

```go
package canonical

import "github.com/gsxhq/gsxui/internal/recipe"

// buttonShape is Button's style interface. Its dimensions mirror the public
// variant and size parameters of button.gsx; every style must implement all of
// it. Values are ordered as they appear in the generated switch.
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

`internal/recipe/component.go` (new file in the recipe package):

```go
package recipe

import "fmt"

// Component binds a Shape to the helper calls a canonical component authors.
// stylegen replaces every method call with concrete style source; the method
// bodies exist so the canonical type-checks and so its style-independent
// behavior tests can run.
type Component struct {
	Shape Shape
}

// Role is the component's base recipe class.
func (c Component) Role() string { return c.Shape.BaseClass() }

// Variant and Size are dimension helpers. An empty or unrecognized value
// resolves to the dimension's declared default, matching the generated
// switch's default arm — so a behavior test written against the canonical
// asserts something true of every style.
func (c Component) Variant(value string) string { return c.class("variant", value) }

func (c Component) Size(value string) string { return c.class("size", value) }

func (c Component) class(dimension, value string) string {
	declared, ok := c.Shape.Dimension(dimension)
	if !ok {
		panic(fmt.Sprintf("recipe: component %q declares no dimension %q", c.Shape.Component, dimension))
	}
	if !declared.Has(value) {
		value = declared.Default
	}
	return c.Shape.ValueClass(dimension, value)
}
```

`registry/canonical/button_recipe.go` gains the binding, below `buttonShape`:

```go
// button binds Button's shape to the helper calls button.gsx authors. The
// variable name is the component name: stylegen resolves button.Variant(v) by
// looking "button" up in Shapes().
var button = recipe.Component{Shape: buttonShape}
```

`registry/canonical/recipe.go`:

```go
// Package canonical holds the structural component sources and the shapes they
// declare. It is compiled so that it type-checks and so that style-independent
// structure and behavior tests can run against the authoritative source, but it
// is never shipped: consumers receive the generated (canonical x recipe) output
// in package ui. Nothing outside internal/stylegen may import it.
package canonical

import "github.com/gsxhq/gsxui/internal/recipe"

var shapes = map[string]recipe.Shape{
	buttonShape.Component: buttonShape,
}

// Shapes returns every declared component shape, keyed by component name.
// internal/stylegen reads this instead of parsing Go declarations as data.
func Shapes() map[string]recipe.Shape {
	out := make(map[string]recipe.Shape, len(shapes))
	maps.Copy(out, shapes)
	return out
}
```

(import `maps` alongside `recipe`.)

- [ ] **Step 4: Run the helper tests**

Run: `go test ./registry/canonical/ -v`
Expected: PASS.

- [ ] **Step 5: Move and rewrite the canonical Button**

```bash
cp ui/button.gsx registry/canonical/button.gsx
```

**Copy, do not move.** `ui/button.gsx` and `ui/button.x.go` stay exactly as they
are: `site/pages/home.gsx` calls `ui.Button`, and Task 7 is what replaces
`ui/button.gsx` with generated output. Removing it here would break
`go build ./...` for three tasks. Task 9 verifies that the final `ui/button.gsx`
is the generated Nova copy.

In `registry/canonical/button.gsx`, change `package ui` to `package canonical`, and replace **both** class attributes (the `<a>` branch and the `<button>` branch) so each reads:

```gsx
class={
	"group/button",
	button.Role(),
	button.Variant(variant),
	button.Size(size),
}
```

Everything else — the parameter list, the `href`/`disabled` branch, `data-variant`, `data-size`, `type="button"`, `{ attrs... }`, `data-gsxui-slot-button`, and the doc comment — stays exactly as it is.

- [ ] **Step 6: Write the style-independent behavior test**

`registry/canonical/button_test.go`. These assertions are true under every style, which is the point of hosting them here:

```go
package canonical_test

import (
	"regexp"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/registry/canonical"
)

var disabledAttr = regexp.MustCompile(`disabled(>|\s)`)

func TestButtonRendersButtonByDefault(t *testing.T) {
	got := render(t, canonical.Button("", "", "", false, gsx.Raw("Save"), nil))
	for _, want := range []string{
		"<button", "data-gsxui-slot-button", `type="button"`,
		`data-variant="default"`, `data-size="default"`, ">Save</button>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if disabledAttr.MatchString(got) {
		t.Errorf("unexpected disabled attr\nin: %s", got)
	}
}

func TestButtonWithHrefRendersAnchor(t *testing.T) {
	got := render(t, canonical.Button("", "", "/docs", false, gsx.Raw("Docs"), nil))
	if !strings.Contains(got, "<a ") || !strings.Contains(got, `href="/docs"`) {
		t.Errorf("want an anchor with href\nin: %s", got)
	}
}

func TestDisabledAlwaysRendersButtonEvenWithHref(t *testing.T) {
	got := render(t, canonical.Button("", "", "/docs", true, gsx.Raw("Docs"), nil))
	if !strings.Contains(got, "<button") {
		t.Errorf("disabled with href must render a button\nin: %s", got)
	}
	if !disabledAttr.MatchString(got) {
		t.Errorf("missing disabled attr\nin: %s", got)
	}
}

func TestCallerAttrsFallThrough(t *testing.T) {
	got := render(t, canonical.Button("", "", "", false, gsx.Raw("x"),
		gsx.Attrs{"data-testid": "save"}))
	if !strings.Contains(got, `data-testid="save"`) {
		t.Errorf("missing caller attr\nin: %s", got)
	}
}

func TestUnrecognizedVariantResolvesToDefault(t *testing.T) {
	got := render(t, canonical.Button("destructve", "", "", false, gsx.Raw("x"), nil))
	if !strings.Contains(got, "gsxui-recipe-button-variant-default") {
		t.Errorf("unrecognized variant must resolve to the default\nin: %s", got)
	}
}
```

Copy the `render` helper from `ui/button_test.go` into a new `registry/canonical/render_test.go`, changing only its package clause to `canonical_test`.

- [ ] **Step 7: Generate and run**

Run: `go generate ./... || npx gsx generate` (whichever the repo uses; check `Makefile` for the gsx codegen target), then `go test ./registry/canonical/ -v`
Expected: PASS.

- [ ] **Step 8: Commit**

```bash
git add -A registry/canonical
git commit -m "feat: add the canonical component package with recipe helpers"
```

---

### Task 6: Desugar helper calls in stylegen

**Files:**
- Modify: `internal/stylegen/resolve.go`
- Test: `internal/stylegen/resolve_test.go`
- Test: `internal/stylegen/testdata/resolve-input.gsx.txt`, `internal/stylegen/testdata/resolve-nova.golden.gsx.txt`

**Interfaces:**
- Consumes: `recipe.Resolved` (Task 3)
- Produces: `func Resolve(filename string, src []byte, resolved recipe.Resolved) ([]byte, error)`, `func HelperCalls(filename string, src []byte) ([]Call, error)` with `type Call struct{ Component, Dimension string }` (`Dimension` empty for `Role`)

Keep every existing rejection rule and the format/reparse/residue verification. The change is *what* gets substituted: a `ClassPart` whose `Expr` is a helper call is replaced over `[ExprPos, End())` (see Resolved Spike).

- [ ] **Step 1: Write the failing test**

```go
func TestResolveDesugarsHelperCalls(t *testing.T) {
	t.Parallel()
	src := []byte(`package canonical

import "github.com/gsxhq/gsx"

component B(variant string, children gsx.Node) {
	<button class={ "group/button", button.Role(), button.Variant(variant) }>
		{ children }
	</button>
}
`)
	resolved := testResolved(t) // shape: variant{default,outline}, base{inline-flex}
	got, err := Resolve("button.gsx", src, resolved)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	for _, want := range []string{
		`"inline-flex items-center"`,
		`case "outline":`,
		`"border-border bg-background"`,
		`default:`,
		`"bg-primary"`,
	} {
		if !strings.Contains(string(got), want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if strings.Contains(string(got), recipe.Prefix) {
		t.Errorf("resolved source still contains the recipe prefix:\n%s", got)
	}
	if strings.Contains(string(got), "button.Role()") {
		t.Errorf("resolved source still contains a helper call:\n%s", got)
	}
}

func TestResolveDefaultArmCarriesDeclaredDefault(t *testing.T) {
	t.Parallel()
	// The declared default's utilities must appear in the default arm, not "",
	// so a misspelled variant renders the default rather than nothing.
	src := []byte(`package canonical

import "github.com/gsxhq/gsx"

component B(variant string, children gsx.Node) {
	<button class={ button.Variant(variant) }>{ children }</button>
}
`)
	got, err := Resolve("button.gsx", src, testResolved(t))
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if !regexp.MustCompile(`default:\s*\n?\s*"bg-primary"`).Match(got) {
		t.Errorf("default arm must carry the declared default's utilities\nin: %s", got)
	}
	// The default value must not also get its own redundant case arm.
	if strings.Contains(string(got), `case "default":`) {
		t.Errorf("default value must not get a redundant case arm\nin: %s", got)
	}
}

func TestResolveRejectsUnknownDimension(t *testing.T) {
	t.Parallel()
	src := []byte(`package canonical

import "github.com/gsxhq/gsx"

component B(tone string, children gsx.Node) {
	<button class={ button.Tone(tone) }>{ children }</button>
}
`)
	_, err := Resolve("button.gsx", src, testResolved(t))
	if err == nil {
		t.Fatal("Resolve() = nil error, want error")
	}
	if !strings.Contains(err.Error(), "button.gsx:6:") {
		t.Errorf("error must carry a source position, got %q", err)
	}
}
```

Add `testResolved(t)` to the test file, building a `recipe.Resolved` with base `{"inline-flex","items-center"}`, variant `default`→`{"bg-primary"}`, `outline`→`{"border-border","bg-background"}`.

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/stylegen/ -run TestResolveDesugars -v`
Expected: FAIL — `Resolve` has the old signature.

- [ ] **Step 3: Write the implementation**

Change `inspectClassExpr` so that, when the parsed expression is a `*goast.CallExpr` whose `Fun` is a `*goast.SelectorExpr` — that is, a call of the form `<component>.<Method>(...)`:

1. Require `X` to be a bare `*goast.Ident`. Its name is the component name; reject with a positioned error if it does not match `resolved.Shape.Component`.
2. If `Sel` is `Role`, require zero arguments; record an edit replacing the span with `strconv.Quote(strings.Join(resolved.Base, " "))`.
3. Otherwise the dimension is `strings.ToLower(Sel.Name)`. Require exactly one argument and that it is a bare `*goast.Ident`; look the dimension up in `resolved.Shape` and reject an unknown one with a positioned error; record an edit replacing the span with the generated switch.

Switch generation — every declared value except the default gets a `case`, and the default's utilities go in the `default:` arm:

```go
func dimensionSwitch(resolved recipe.Resolved, dimension recipe.Dimension, argument string) string {
	var out strings.Builder
	fmt.Fprintf(&out, "switch %s {\n", argument)
	for _, value := range dimension.Values {
		if value == dimension.Default {
			continue
		}
		fmt.Fprintf(&out, "case %s:\n%s\n",
			strconv.Quote(value),
			strconv.Quote(strings.Join(resolved.Utilities(dimension.Name, value), " ")))
	}
	fmt.Fprintf(&out, "default:\n%s\n}",
		strconv.Quote(strings.Join(resolved.Utilities(dimension.Name, dimension.Default), " ")))
	return out.String()
}
```

Compute the edit span as established by the spike:

```go
start := r.fset.Position(part.ExprPos).Offset
end := r.fset.Position(part.End()).Offset
```

Keep `applyLiteralEdits`, `gen.Format`, the `gsxparser.ParseFile` reparse, and the residue checks exactly as they are, but change the residue check to also reject any surviving `.Role(`, `.Variant(`, or `.Size(` call.

**Staged retirement of the token path.** The legacy token-literal substitution
cannot be deleted here. Because Task 5 copies rather than moves,
`ui/button.gsx` is still in the old token form (`"gsxui-recipe-button-variant-outline"`,
`default: ""`), and `GenerateButton` resolves that file — so deleting
`resolveClassLiteral` would break the build and the committed-artifact tests,
violating the green-commit constraint.

Therefore in this task:

- Add the new `Resolve` (taking `recipe.Resolved`) and `HelperCalls` alongside
  the existing code.
- Rename the legacy entry point to an unexported `resolveTokens`, keeping
  `ResolveReport`, `resolveClassLiteral`, `declaredRecipeTokens` and the
  `used`/`declared` token-set machinery intact and still called by
  `generate.go`'s `GenerateButton` and by the existing token tests.
- Add a comment on `resolveTokens` stating it is the legacy path, retired by
  Task 7.

Task 7 deletes `resolveTokens`, `ResolveReport`, `resolveClassLiteral`,
`declaredRecipeTokens`, `GenerateButton`, and their tests in the same commit
that replaces `ui/button.gsx` with generated output — at which point nothing
authored is in token form and the capability is genuinely dead.

`HelperCalls` reuses the same traversal in a validation-only mode and returns the calls without editing — `GenerateAll` uses it to check the canonical against the shapes before touching any style.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/stylegen/ -v`
Expected: PASS. Update `testdata/resolve-input.gsx.txt` and `testdata/resolve-nova.golden.gsx.txt` to the helper-call form and its expansion.

- [ ] **Step 5: Commit**

```bash
git add internal/stylegen/resolve.go internal/stylegen/resolve_test.go internal/stylegen/testdata
git commit -m "feat: desugar recipe helper calls into generated switches"
```

---

### Task 7: Manifest-driven generation

**Files:**
- Modify: `internal/stylegen/generate.go`
- Modify: `cmd/stylegen/main.go`
- Test: `internal/stylegen/generate_test.go`

**Interfaces:**
- Consumes: `Resolve`, `HelperCalls` (Task 6), `recipe.Conform`, `recipe.CheckConflicts` (Tasks 3–4), `canonical.Shapes()` (Task 5)
- Produces: `const DefaultStyle = "nova"`, `func GenerateAll(root string, check bool) error`

Discovery replaces the hardcoded `buttonStyleSources` slice:

```text
for each registry/canonical/<component>.gsx
  for each registry/styles/<style>/<component>.css
    emit registry/generated/<style>/<component>.gsx   package ui
    emit site/stylepreview/<style>/<component>.gsx    package <style>
  copy registry/generated/<DefaultStyle>/<component>.gsx -> ui/<component>.gsx
```

- [ ] **Step 1: Write the test helpers**

Add these to `internal/stylegen/generate_test.go` first — later tasks use them too:

```go
// repoRoot walks up from the test's working directory to the module root.
func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("module root not found")
		}
		dir = parent
	}
}

// copyRepoFixture copies the inputs and artifacts GenerateAll touches into dst,
// so a test can corrupt one without mutating the working tree.
func copyRepoFixture(t *testing.T, dst string) {
	t.Helper()
	root := repoRoot(t)
	for _, dir := range []string{"registry", "ui", filepath.Join("site", "stylepreview")} {
		src := filepath.Join(root, dir)
		err := filepath.WalkDir(src, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			relative, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			target := filepath.Join(dst, relative)
			if entry.IsDir() {
				return os.MkdirAll(target, 0o755)
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			return os.WriteFile(target, content, 0o644)
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}
```

- [ ] **Step 2: Write the failing test**

```go
func TestGenerateAllIsIdempotent(t *testing.T) {
	root := repoRoot(t)
	if err := GenerateAll(root, false); err != nil {
		t.Fatalf("GenerateAll() error = %v", err)
	}
	if err := GenerateAll(root, true); err != nil {
		t.Fatalf("GenerateAll(check) after write = %v, want nil", err)
	}
}

func TestGenerateAllEmitsEveryStyleAndTheDefaultCopy(t *testing.T) {
	root := repoRoot(t)
	if err := GenerateAll(root, false); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{
		"registry/generated/nova/button.gsx",
		"registry/generated/maia/button.gsx",
		"site/stylepreview/nova/button.gsx",
		"site/stylepreview/maia/button.gsx",
		"ui/button.gsx",
	} {
		if _, err := os.Stat(filepath.Join(root, path)); err != nil {
			t.Errorf("missing generated artifact %s: %v", path, err)
		}
	}
	shipped, err := os.ReadFile(filepath.Join(root, "ui", "button.gsx"))
	if err != nil {
		t.Fatal(err)
	}
	fromDefault, err := os.ReadFile(filepath.Join(root, "registry", "generated", DefaultStyle, "button.gsx"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(shipped, fromDefault) {
		t.Error("ui/button.gsx must be a copy of the default style's generated output")
	}
}

func TestGeneratedSourcesAreFreeOfRecipeConstructs(t *testing.T) {
	root := repoRoot(t)
	if err := GenerateAll(root, false); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"ui/button.gsx", "registry/generated/maia/button.gsx"} {
		src, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(src, []byte(recipe.Prefix)) {
			t.Errorf("%s contains the recipe prefix", path)
		}
		for _, helper := range []string{".Role(", ".Variant(", ".Size("} {
			if bytes.Contains(src, []byte(helper)) {
				t.Errorf("%s contains a helper call %q", path, helper)
			}
		}
	}
}

func TestGenerateAllValidatesBeforeWriting(t *testing.T) {
	root := t.TempDir()
	copyRepoFixture(t, root) // helper: copies registry/, ui/, site/stylepreview/ into root
	styleFile := filepath.Join(root, "registry", "styles", "maia", "button.css")
	src, err := os.ReadFile(styleFile)
	if err != nil {
		t.Fatal(err)
	}
	broken := bytes.Replace(src, []byte(".gsxui-recipe-button-size-icon-lg"), []byte(".gsxui-recipe-button-size-icon-xl"), 1)
	if err := os.WriteFile(styleFile, broken, 0o644); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(filepath.Join(root, "ui", "button.gsx"))
	if err != nil {
		t.Fatal(err)
	}
	if err := GenerateAll(root, false); err == nil {
		t.Fatal("GenerateAll() = nil, want an error for the undeclared value")
	}
	after, err := os.ReadFile(filepath.Join(root, "ui", "button.gsx"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Error("a failing run must not mutate any artifact")
	}
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `go test ./internal/stylegen/ -run TestGenerateAll -v`
Expected: FAIL — `undefined: GenerateAll`.

- [ ] **Step 4: Write the implementation**

Rewrite `generate.go`:

```go
const DefaultStyle = "nova"

// GenerateAll resolves every canonical component against every authored style
// before it checks or mutates any artifact. A failing run writes nothing.
func GenerateAll(root string, check bool) error {
	outputs, err := resolveAll(root)
	if err != nil {
		return err
	}
	if check {
		return checkOutputs(root, outputs)
	}
	for _, output := range outputs {
		if err := writeGenerated(filepath.Join(root, output.relativePath), output.content); err != nil {
			return fmt.Errorf("write %s: %w", filepath.ToSlash(output.relativePath), err)
		}
	}
	return nil
}
```

`resolveAll` walks `registry/canonical/*.gsx`, looks each component's shape up in `canonical.Shapes()` (error if absent), validates the canonical's helper calls with `HelperCalls`, then for each `registry/styles/*/` directory: `recipe.ParseStyle` → `recipe.Conform` → `recipe.CheckConflicts(…, merge.Merge)` → `Resolve`. It accumulates three outputs per style (the `package ui` generated source, the `package <style>` preview via the existing `rewriteGSXPackage`, and — for `DefaultStyle` — the `ui/<component>.gsx` copy). Sort styles by name so output order is deterministic.

Delete the legacy token path that Task 6 staged for retirement: `resolveTokens`, `ResolveReport`, `resolveClassLiteral`, `declaredRecipeTokens`, `recipeSetDifference`, `GenerateButton`, and every test that exercises token-literal substitution (`TestResolveGolden`, the `UsedTokens` assertions, and the `TestGenerateButton*` artifact tests). After this task nothing authored is in token form, so the capability is dead code. The canonical is the only input, and `registry/canonical`'s own tests plus the new `GenerateAll` tests cover what those tests used to.

Rename the `generatedButtonSource` struct to `generatedSource` (same fields: `relativePath string`, `content []byte`). Keep `writeGeneratedButton` (rename to `writeGenerated`, and generalize its temp-file pattern from `.button.gsx-*` to `.gsx-*`), `checkButtonSources` (rename to `checkOutputs`), and the atomic temp-then-rename discipline unchanged.

In `cmd/stylegen/main.go`, replace `stylegen.GenerateButton(root, *check)` with `stylegen.GenerateAll(root, *check)`, and change `repositoryRoot`'s marker file from `ui/button.gsx` to `registry/canonical/button.gsx`.

- [ ] **Step 4b: Update the Makefile audit rules**

`make audit` is the first target of `make ci`, and several of its rules name
`ui/button.gsx` because that file used to be the recipe-token component. This
task makes it generated concrete Tailwind, which breaks the strictest rule.
Apply exactly these changes to the `audit:` target:

Delete this rule entirely — it allowed only `"group/button"` and
`gsxui-recipe-*` string literals in `ui/button.gsx`, which generated output
cannot satisfy:

```make
	@! rg -n -P '^[[:space:]]+"(?!(?:group/button|gsxui-recipe-[^"]*|)"[,]?[[:space:]]*$$)' ui/button.gsx
```

Replace it with the equivalent guarantee applied to the canonical instead —
the canonical must express recipe roles through helper calls, never literal
tokens, and must carry no concrete utilities:

```make
	@! rg -n 'gsxui-recipe-' registry/canonical -g '*.gsx'
```

Tighten the recipe-prefix rule, dropping the `ui/button.gsx` exemption. After
this task no file under `ui/` may contain a recipe token at all, which is
exactly success criterion 4:

```make
	@! rg -n 'gsxui-recipe-' ui -g '*.gsx'
```

Add `registry/canonical` to `audit-source-dirs` (line 3) so the `data-slot`
and `data-gsxui-slot` guards cover the canonical sources too.

Leave every other rule unchanged. In particular the `-g '!button.gsx'`
exemptions on the `group/`, `peer/` and `class=` rules stay: generated
`ui/button.gsx` still legitimately carries `group/button` and a `class={…}`
attribute.

Verify with `make audit` and confirm it exits 0.

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test ./internal/stylegen/ -v && go run ./cmd/stylegen --check`
Expected: PASS, and the check run exits 0.

- [ ] **Step 6: Commit**

```bash
git add internal/stylegen/generate.go internal/stylegen/generate_test.go cmd/stylegen/main.go
git commit -m "feat: drive style generation from a discovered manifest"
```

---

### Task 8: Contract JSON

**Files:**
- Create: `internal/recipe/contract.go`
- Modify: `internal/stylegen/generate.go`
- Test: `internal/recipe/contract_test.go`

**Interfaces:**
- Consumes: `Resolved` (Task 3)
- Produces: `type Contract`, `func BuildContract(components map[string]Shape, styles map[string]map[string]Resolved) Contract`, `func (Contract) MarshalIndent() ([]byte, error)`

Emitted to `registry/generated/recipes.json` by `GenerateAll` and drift-checked with every other artifact. The shape is emitted once and shared by all styles — strict conformance made visible in the artifact.

- [ ] **Step 1: Write the failing test**

```go
func TestBuildContractEmitsSharedShapeAndPerStyleUtilities(t *testing.T) {
	t.Parallel()
	shape := conformShape()
	resolved, err := Conform("nova/button.css", shape, mustParse(t, conformCSS))
	if err != nil {
		t.Fatal(err)
	}
	contract := BuildContract(
		map[string]Shape{"button": shape},
		map[string]map[string]Resolved{"nova": {"button": resolved}},
	)
	got, err := contract.MarshalIndent()
	if err != nil {
		t.Fatal(err)
	}
	var round map[string]any
	if err := json.Unmarshal(got, &round); err != nil {
		t.Fatalf("contract is not valid JSON: %v", err)
	}
	for _, want := range []string{
		`"version": 1`,
		`"components"`,
		`"default": "default"`,
		`"outline"`,
		`"styles"`,
		`"nova"`,
		`"border-border"`,
	} {
		if !strings.Contains(string(got), want) {
			t.Errorf("contract missing %q\nin: %s", want, got)
		}
	}
	// The shape must appear once, under components — never repeated per style.
	if strings.Count(string(got), `"values"`) != 1 {
		t.Errorf("shape must be emitted exactly once, got:\n%s", got)
	}
}

func TestContractIsDeterministic(t *testing.T) {
	t.Parallel()
	shape := conformShape()
	resolved, err := Conform("nova/button.css", shape, mustParse(t, conformCSS))
	if err != nil {
		t.Fatal(err)
	}
	build := func() []byte {
		contract := BuildContract(
			map[string]Shape{"button": shape},
			map[string]map[string]Resolved{"nova": {"button": resolved}, "maia": {"button": resolved}},
		)
		out, err := contract.MarshalIndent()
		if err != nil {
			t.Fatal(err)
		}
		return out
	}
	if !bytes.Equal(build(), build()) {
		t.Error("contract marshalling must be deterministic across map iteration orders")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/recipe/ -run TestBuildContract -v`
Expected: FAIL — `undefined: BuildContract`.

- [ ] **Step 3: Write the implementation**

```go
package recipe

import (
	"bytes"
	"encoding/json"
)

// ContractVersion is bumped whenever the emitted schema changes shape.
const ContractVersion = 1

// Contract is the serialized recipe model consumed by the theme editor. These
// types are marshalled, so their fields are exported.
type Contract struct {
	Version    int                                    `json:"version"`
	Components map[string]ContractComponent           `json:"components"`
	Styles     map[string]map[string]ContractUtilities `json:"styles"`
}

type ContractComponent struct {
	Base       bool                         `json:"base"`
	Dimensions map[string]ContractDimension `json:"dimensions"`
}

type ContractDimension struct {
	Default string   `json:"default"`
	Values  []string `json:"values"`
}

type ContractUtilities struct {
	Base       []string                       `json:"base,omitempty"`
	Dimensions map[string]map[string][]string `json:"dimensions"`
}

func BuildContract(components map[string]Shape, styles map[string]map[string]Resolved) Contract {
	contract := Contract{
		Version:    ContractVersion,
		Components: make(map[string]ContractComponent, len(components)),
		Styles:     make(map[string]map[string]ContractUtilities, len(styles)),
	}
	for name, shape := range components {
		dimensions := make(map[string]ContractDimension, len(shape.Dimensions))
		for _, dimension := range shape.Dimensions {
			dimensions[dimension.Name] = ContractDimension{
				Default: dimension.Default,
				Values:  slices.Clone(dimension.Values),
			}
		}
		contract.Components[name] = ContractComponent{Base: shape.Base, Dimensions: dimensions}
	}
	for style, resolvedComponents := range styles {
		entry := make(map[string]ContractUtilities, len(resolvedComponents))
		for name, resolved := range resolvedComponents {
			values := make(map[string]map[string][]string, len(resolved.Shape.Dimensions))
			for _, dimension := range resolved.Shape.Dimensions {
				byValue := make(map[string][]string, len(dimension.Values))
				for _, value := range dimension.Values {
					byValue[value] = resolved.Utilities(dimension.Name, value)
				}
				values[dimension.Name] = byValue
			}
			entry[name] = ContractUtilities{Base: slices.Clone(resolved.Base), Dimensions: values}
		}
		contract.Styles[style] = entry
	}
	return contract
}

// MarshalIndent renders the contract deterministically. encoding/json sorts map
// keys, so identical inputs always produce identical bytes.
func (c Contract) MarshalIndent() ([]byte, error) {
	out, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(bytes.TrimRight(out, "\n"), '\n'), nil
}
```

- [ ] **Step 4: Wire it into generation**

In `resolveAll`, accumulate the per-style `recipe.Resolved` values, then append one more output:

```go
contract := recipe.BuildContract(shapesByComponent, resolvedByStyle)
contractJSON, err := contract.MarshalIndent()
if err != nil {
	return nil, fmt.Errorf("marshal recipe contract: %w", err)
}
outputs = append(outputs, generatedSource{
	relativePath: filepath.Join("registry", "generated", "recipes.json"),
	content:      contractJSON,
})
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test ./internal/recipe/ ./internal/stylegen/ -v && go run ./cmd/stylegen && go run ./cmd/stylegen --check`
Expected: PASS; `registry/generated/recipes.json` exists and the check exits 0.

- [ ] **Step 6: Commit**

```bash
git add internal/recipe/contract.go internal/recipe/contract_test.go internal/stylegen/generate.go registry/generated/recipes.json
git commit -m "feat: emit the generated recipe contract"
```

---

### Task 9: Boundary test, downstream tests, and full verification

**Files:**
- Create: `internal/stylegen/architecture_test.go`
- Modify: `ui/button_test.go`
- Modify: `jstest/support/compiled-css-audit.test.ts`
- Modify: `Makefile` (if it has a stylegen target referencing `GenerateButton`)

**Interfaces:**
- Consumes: everything above
- Produces: no new API

- [ ] **Step 1: Write the boundary test**

```go
package stylegen

import (
	"go/build"
	"path/filepath"
	"strings"
	"testing"
)

// The canonical package is compiled so it type-checks and hosts
// style-independent behavior tests. It must never ship, so nothing outside
// stylegen and its own tests may import it.
func TestNothingOutsideStylegenImportsCanonical(t *testing.T) {
	const canonicalPath = "github.com/gsxhq/gsxui/registry/canonical"
	root := repoRoot(t)

	allowed := map[string]bool{
		"github.com/gsxhq/gsxui/internal/stylegen":     true,
		"github.com/gsxhq/gsxui/registry/canonical":    true,
	}

	out, err := exec.Command("go", "list", "./...").Output()
	if err != nil {
		t.Fatal(err)
	}
	packages := strings.Fields(string(out))
	for _, importPath := range packages {
		if allowed[importPath] {
			continue
		}
		pkg, err := build.Import(importPath, filepath.Join(root, "."), 0)
		if err != nil {
			continue
		}
		for _, imported := range append(pkg.Imports, pkg.TestImports...) {
			if imported == canonicalPath || strings.HasPrefix(imported, canonicalPath+"/") {
				t.Errorf("%s imports %s; the canonical package must never ship", importPath, canonicalPath)
			}
		}
	}
}
```

- [ ] **Step 2: Run it**

Run: `go test ./internal/stylegen/ -run TestNothingOutside -v`
Expected: PASS.

- [ ] **Step 3: Retarget the ui tests**

`ui/button_test.go` currently asserts recipe *token* classes via `canonicalButtonClass()`. `ui.Button` now renders concrete Nova utilities, so:

- Delete the tests that duplicate what `registry/canonical/button_test.go` now covers: `TestButtonDefault`'s structural assertions, the href/disabled/attrs cases. Behavior lives with the canonical.
- Keep `TestButtonPinned` and `TestButtonVariantAndSizeAxes` as **Nova pins**, replacing `canonicalButtonClass(variant, size)` with the concrete Nova utility strings read from `registry/generated/recipes.json` or inlined verbatim.
- Add a comment at the top of the file recording that this file pins the *generated Nova* output, and that style-independent behavior belongs in `registry/canonical/button_test.go`.

- [ ] **Step 4: Retarget the compiled-CSS audit**

`jstest/support/compiled-css-audit.test.ts` builds each recipe stylesheet with Tailwind. The recipes are unchanged, so this test keeps working as-is — but its second test, "compiled site Button fallback has the exact normal-site scope", asserts a `.gsxui-recipe-button` fallback selector in the site CSS that no longer exists now that `ui/button.gsx` carries concrete utilities. Delete that second test; the first (per-style recipe compilation) is the one that matters and stays.

Verify `web/site.css`'s `@source` globs still reach the generated sources — `@source "../ui/**/*.gsx"` covers `ui/button.gsx`, and `@source "../site/**/*.gsx"` covers the stylepreview packages. Add `@source "../registry/generated/**/*.gsx";` so styles a consumer might install are also scanned.

- [ ] **Step 4b: Stop the generation tests mutating the working tree**

Task 7's `TestGenerateAllIsIdempotent`, `TestGenerateAllEmitsEveryStyleAndTheDefaultCopy`
and `TestGeneratedSourcesAreFreeOfRecipeConstructs` call `GenerateAll(repoRoot(t), false)`,
so `go test` silently repairs stale artifacts in `ui/`, `registry/generated/`
and `site/stylepreview/`. A drift bug can therefore be masked simply by having
run the suite. This was mandated by the plan, not chosen by the implementer.

Point all three at an isolated copy instead, exactly as
`TestGenerateAllWritesDeterministicallyAndCheckNeverWrites` already does:

```go
root := t.TempDir()
copyRepoFixture(t, root)
```

`TestGenerateAllValidatesBeforeWriting` already uses this pattern — leave it.
After the change, `go test ./...` must leave `git status --short` empty.

- [ ] **Step 4c: Remove two tautological assertions**

`ui/input-group_test.go:112` and `:161` assert that `got` does NOT contain
`canonicalButtonClass(<other variant>, …)`. Since that helper returns a whole
`class="…"` attribute and an element has exactly one class attribute, those
negatives cannot be true once the positive assertion above them passes — they
read as coverage they do not provide. Delete both, and delete the dead
`tt.want != "default"` guard at `:110` (no case in that table has
`want == "default"`).

- [ ] **Step 4d0: Harden the contract's "shape emitted once" assertion**

`internal/recipe/contract_test.go`'s `strings.Count(string(got), `"values"`) != 1`
only works because the fixture happens to declare exactly one dimension. It
asserts "not duplicated" by coincidence, not "emitted once" in general. Derive
the expected count from the fixture instead:

```go
if got, want := strings.Count(string(out), `"values"`), len(shape.Dimensions); got != want {
	t.Errorf("shape emitted %d times, want %d:\n%s", got, want, out)
}
```

- [ ] **Step 4d: Replace the containsStyle loop**

`internal/stylegen/generate.go:211-218` — use `slices.Contains`. It is the only
`gopls check -severity=hint` hit in the changed files.

- [ ] **Step 5: Full verification**

Run each of these and confirm the stated expectation before claiming completion:

```bash
go build ./...                       # expect: no output
make audit                           # expect: exit 0
go run ./cmd/stylegen --check        # expect: exit 0, no drift
go test ./...                        # expect: all packages ok
npm test                             # expect: jstest suite passes
git status --short                   # expect: clean, no unstaged generated files
```

- [ ] **Step 6: Verify the success criteria from the spec**

Confirm each by inspection or command, and report the evidence:

1. `grep -c 'gsxui-recipe\|inline-flex' registry/canonical/button.gsx` → 0, and its class attributes contain only `"group/button"` plus `button.*` method calls
2. Shape declared once: `grep -rn 'recipe.Shape{' registry/ internal/` → one hit
3. Break a rule in `registry/styles/maia/button.css`, run `go run ./cmd/stylegen`, confirm the error names the exact missing `(dimension, value)`, then restore
4. `go run ./cmd/stylegen --check` → exit 0
5. `registry/canonical/button_test.go`'s `TestUnrecognizedVariantResolvesToDefault` passes, and the generated `default:` arm carries the default's utilities
6. `registry/generated/recipes.json` contains every variant, both styles, and shared shape
7. `go test ./registry/canonical/` → PASS
8. `go test ./internal/stylegen/ -run TestNothingOutside` → PASS

- [ ] **Step 7: Commit**

```bash
git add -A
git commit -m "test: enforce the canonical boundary and retarget downstream tests"
```

---

## Notes for the implementer

**Do not redesign the resolver's safety machinery.** The parse → edit → format → reparse → residue-check pipeline and the rejection rules (helper or token in a static class attribute, non-class attribute, interpolation, pipeline argument, switch tag, or nested markup) are load-bearing and already well tested. Task 6 changes *what* is substituted, not *how* substitution is verified.

**The `default:` arm is the point.** The old generated code emitted `default: ""`, so a misspelled variant produced an unstyled element with no error anywhere. Generating the declared default's utilities into that arm is the fix, and `TestResolveDefaultArmCarriesDeclaredDefault` plus `TestUnrecognizedVariantResolvesToDefault` are the tests that hold it in place. Do not weaken them.

**Order of validation matters.** `GenerateAll` must complete every parse, conformance check, conflict check, and resolve for every component and style *before* writing any file. `TestGenerateAllValidatesBeforeWriting` enforces this; a partial write on a failing run is a defect, not a detail.

**Two languages, deliberately.** Shapes are Go, recipes are CSS. `internal/recipe` bridges them and must not grow a dependency on gsx or the filesystem — that is what keeps it testable as pure data. The `merger` parameter in `CheckConflicts` exists for the same reason.

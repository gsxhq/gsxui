# Slot Axis (Stage 0) and Layer Gate (Stage 1) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Give the recipe model a slot axis so multi-slot components can be expressed, generate per-slot accessors without a bootstrap cycle, and build the layer-precedence gate that must exist before any bulk migration.

**Architecture:** `recipe.Shape` gains `Slots []Slot`, with `Base` and `Dimensions` moving onto `Slot`. Shape declarations move to a leaf package `registry/canonical/shapes` that compiles on its own, so `stylegen` can import them and generate per-slot accessor methods back into `registry/canonical` without the package having to already compile. Class decoding resolves slot names longest-match-first. Stage 1 then adds a `make audit` rule and a committed computed-style sweep harness, converting the CSS-layer precedence hazard from a matter of vigilance into a build failure.

**Tech Stack:** Go 1.24+, `github.com/gsxhq/gsx` (parser/ast/gen), `github.com/tdewolff/parse/v2/css`, `github.com/gsxhq/gsxui/merge`, `encoding/json`, Playwright + Chromium (sweep harness).

**Spec:** `docs/superpowers/specs/2026-07-29-slot-axis-and-component-migration-design.md`
**Base design (still binding):** `docs/superpowers/specs/2026-07-29-typed-recipe-model-design.md`

## Global Constraints

- Recipe class prefix is exactly `gsxui-recipe-`. It must never appear in any generated consumer artifact.
- Default style is `nova`, named in exactly one place (`stylegen.DefaultStyle`).
- Styles are strictly conformant: every style implements every `(slot, dimension, value)` in the shape.
- **Button's generated output must be byte-identical before and after Stage 0.** This is the proof the slot axis is backward compatible. Any diff to `ui/button.gsx` or `registry/generated/*/button.gsx` other than through a deliberate, justified change is a failure.
- **Intermediate commits need NOT be green.** This is one atomic cross-package refactor; the model change breaks every consumer at once. The FINAL state of the plan must be fully green. Reviewers: a broken build inside Tasks 1-7 is sanctioned by the plan, not a finding.
- `gofmt -l .` must print nothing (Makefile enforces `gofmt -l . | (! grep .)`).
- `make audit` must pass — it is the first target of `make ci`.
- `go run ./cmd/stylegen --check` must exit 0.
- Nothing outside `internal/stylegen` and its tests may import `registry/canonical`. The new `registry/canonical/shapes` package inherits the same rule, except that `internal/recipe` must NOT import it either (that would be a cycle).
- Use `{/* */}` for comments inside gsx markup, never `//`.
- Per the user's Go conventions: unexported unless serialisation or cross-package use requires otherwise.

## Known-Resolved Facts (do not re-derive)

- `ClassPart.Pos()` includes leading newline and indentation; `End()` excludes the trailing comma but DOES cover a `: cond` guard and pipeline stages. The rewrite span is `[ExprPos, End())`, used only after `validateHelperPart` rejects non-plain parts. A pending gsx fix (gsxhq/gsx#174) does not change this — the guard is on AST fields, independent of span semantics.
- `gsx.Attrs` is `[]Attr{Key,Value}`, not a map.
- On `data-*` names a plain Go bool already renders as presence; `gsx.Toggle` is redundant there and has been removed.
- Sidebar declares slots `menu`, `menu-action`, `menu-badge`, `menu-button`, `menu-button-tooltip`, `menu-button-tooltip-content` — prefix collisions are real, so longest-match decoding is mandatory.
- Per-slot dimensions already exist in `assets/css/styles/default.css`, e.g. `[data-gsxui-slot-badge][data-variant="destructive"]`.

## File Structure

**New:**
- `registry/canonical/shapes/shapes.go` — package doc + the `All()` registry
- `registry/canonical/shapes/button.go` — `var Button = recipe.Shape{…}` (moved)
- `internal/stylegen/accessors.go` — accessor code generation
- `internal/stylegen/accessors_test.go`
- `registry/canonical/button_recipe.gen.go` — GENERATED
- `jstest/support/computed-sweep.ts` — the sweep harness (Stage 1)
- `jstest/specs/layer-precedence.spec.ts` — the committed sweep spec (Stage 1)

**Modified:**
- `internal/recipe/shape.go` — `Slot` type, slot-aware encoding/decoding
- `internal/recipe/validate.go` — slot-aware `Conform`, `CheckConflicts`
- `internal/recipe/contract.go` — schema v2
- `internal/recipe/component.go` — generic slot access replacing fixed `Role`/`Variant`/`Size`
- `internal/stylegen/resolve.go` — matcher resolves method → `(slot, dimension)`
- `internal/stylegen/generate.go` — accessor generation wired into `GenerateAll`
- `registry/canonical/button_recipe.go` — becomes the `recipe.Component` binding only
- `registry/canonical/button.gsx` — `button.Role()` → `button.Root()`
- `Makefile` — the Stage 1 audit rule

**Deleted:**
- none

---

### Task 1: The Slot type and slot-aware Shape

**Files:**
- Modify: `internal/recipe/shape.go`
- Test: `internal/recipe/shape_test.go`

**Interfaces:**
- Consumes: existing `Dimension{Name, Default string; Values []string}`, `(Dimension) Has(value string) bool`
- Produces: `recipe.Slot{Name string; Base bool; Dimensions []Dimension}`, `Shape{Component string; Slots []Slot}`, `func (Shape) Slot(name string) (Slot, bool)`, `func (Slot) Dimension(name string) (Dimension, bool)`, `func (Shape) Validate() error`

The root slot has `Name: ""`. `Shape.Base` and `Shape.Dimensions` are removed — they move onto `Slot`.

- [ ] **Step 1: Write the failing test**

```go
func slotShape() Shape {
	return Shape{
		Component: "card",
		Slots: []Slot{
			{Name: "", Base: true},
			{Name: "header", Base: true, Dimensions: []Dimension{
				{Name: "variant", Default: "default", Values: []string{"default", "muted"}},
			}},
			{Name: "menu-button", Base: true},
		},
	}
}

func TestShapeValidateAcceptsSlots(t *testing.T) {
	t.Parallel()
	if err := slotShape().Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}
}

func TestShapeValidateRejectsSlotProblems(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		mutate  func(*Shape)
		wantErr string
	}{
		{"no slots", func(s *Shape) { s.Slots = nil }, "card: no slots declared"},
		{"duplicate slot", func(s *Shape) { s.Slots[2].Name = "header" }, `card: duplicate slot "header"`},
		{"slot with neither base nor dimensions", func(s *Shape) { s.Slots[2].Base = false }, `card: slot "menu-button" declares neither a base rule nor any dimension`},
		{"bad default in slot dimension", func(s *Shape) { s.Slots[1].Dimensions[0].Default = "loud" }, `card: slot "header" dimension "variant" default "loud" is not one of its values`},
		{"duplicate dimension in slot", func(s *Shape) {
			s.Slots[1].Dimensions = append(s.Slots[1].Dimensions, s.Slots[1].Dimensions[0])
		}, `card: slot "header" duplicate dimension "variant"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := slotShape()
			tt.mutate(&s)
			err := s.Validate()
			if err == nil {
				t.Fatalf("Validate() = nil, want %q", tt.wantErr)
			}
			if got := err.Error(); got != tt.wantErr {
				t.Errorf("Validate() = %q, want %q", got, tt.wantErr)
			}
		})
	}
}

func TestShapeSlotLookup(t *testing.T) {
	t.Parallel()
	s := slotShape()
	root, ok := s.Slot("")
	if !ok || !root.Base {
		t.Fatalf("Slot(\"\") = %+v, %v; want the root slot", root, ok)
	}
	header, ok := s.Slot("header")
	if !ok {
		t.Fatal("Slot(header) = false, want true")
	}
	if _, ok := header.Dimension("variant"); !ok {
		t.Error("header.Dimension(variant) = false, want true")
	}
	if _, ok := s.Slot("nope"); ok {
		t.Error("Slot(nope) = true, want false")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/recipe/ -run TestShapeValidateAcceptsSlots -v`
Expected: FAIL — `unknown field Slots in struct literal`.

- [ ] **Step 3: Write the implementation**

Replace `Shape`'s fields and add `Slot` in `internal/recipe/shape.go`:

```go
// Shape is a component's style interface: the slots it renders and, per slot,
// the dimensions that slot varies over. Every style must implement all of it.
type Shape struct {
	Component string
	Slots     []Slot
}

// Slot is one styled element of a component. Name is "" for the component's
// root element. Dimensions hang off the slot, not the component, because a
// style may vary one slot and not another.
type Slot struct {
	Name       string
	Base       bool
	Dimensions []Dimension
}

func (s Shape) Validate() error {
	if s.Component == "" {
		return fmt.Errorf("component name is empty")
	}
	if len(s.Slots) == 0 {
		return fmt.Errorf("%s: no slots declared", s.Component)
	}
	seen := make(map[string]struct{}, len(s.Slots))
	for _, slot := range s.Slots {
		if _, exists := seen[slot.Name]; exists {
			return fmt.Errorf("%s: duplicate slot %q", s.Component, slot.Name)
		}
		seen[slot.Name] = struct{}{}
		if err := slot.validate(s.Component); err != nil {
			return err
		}
	}
	return nil
}

func (s Slot) validate(component string) error {
	if !s.Base && len(s.Dimensions) == 0 {
		return fmt.Errorf("%s: slot %q declares neither a base rule nor any dimension",
			component, s.Name)
	}
	seen := make(map[string]struct{}, len(s.Dimensions))
	for _, dimension := range s.Dimensions {
		if dimension.Name == "" {
			return fmt.Errorf("%s: slot %q has an unnamed dimension", component, s.Name)
		}
		if _, exists := seen[dimension.Name]; exists {
			return fmt.Errorf("%s: slot %q duplicate dimension %q", component, s.Name, dimension.Name)
		}
		seen[dimension.Name] = struct{}{}
		if err := dimension.validate(component, s.Name); err != nil {
			return err
		}
	}
	return nil
}

// Slot returns the named slot; "" is the root.
func (s Shape) Slot(name string) (Slot, bool) {
	for _, slot := range s.Slots {
		if slot.Name == name {
			return slot, true
		}
	}
	return Slot{}, false
}

// Dimension returns the named dimension of this slot.
func (s Slot) Dimension(name string) (Dimension, bool) {
	for _, dimension := range s.Dimensions {
		if dimension.Name == name {
			return dimension, true
		}
	}
	return Dimension{}, false
}
```

Change `Dimension.validate` to take the slot name and prefix its messages:

```go
func (d Dimension) validate(component, slot string) error {
	if len(d.Values) == 0 {
		return fmt.Errorf("%s: slot %q dimension %q declares no values", component, slot, d.Name)
	}
	seen := make(map[string]struct{}, len(d.Values))
	for _, value := range d.Values {
		if value == "" {
			return fmt.Errorf("%s: slot %q dimension %q has an empty value", component, slot, d.Name)
		}
		if _, exists := seen[value]; exists {
			return fmt.Errorf("%s: slot %q dimension %q duplicates value %q", component, slot, d.Name, value)
		}
		seen[value] = struct{}{}
	}
	if !d.Has(d.Default) {
		return fmt.Errorf("%s: slot %q dimension %q default %q is not one of its values",
			component, slot, d.Name, d.Default)
	}
	return nil
}
```

Delete `Shape.Dimension` (the component-level lookup) — dimensions are per slot now. Existing tests referencing it will fail to compile; update them in Step 4.

- [ ] **Step 4: Update the pre-existing shape tests**

`internal/recipe/shape_test.go`'s `validShape()` and its dependants still use the flat form. Rewrite `validShape()` as a single-root-slot shape so it keeps testing what it tested:

```go
func validShape() Shape {
	return Shape{
		Component: "button",
		Slots: []Slot{{
			Name: "", Base: true,
			Dimensions: []Dimension{
				{Name: "variant", Default: "default", Values: []string{"default", "outline"}},
				{Name: "size", Default: "default", Values: []string{"default", "icon-lg"}},
			},
		}},
	}
}
```

Update every assertion whose expected error string changed to the new slot-qualified form. Do not weaken any assertion — the messages gain a `slot "…"` clause and nothing else.

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test ./internal/recipe/ -v`
Expected: PASS. `internal/stylegen` and `registry/canonical` will not compile yet — Tasks 5-7 fix them, and that is expected.

- [ ] **Step 6: Commit**

```bash
git add internal/recipe/shape.go internal/recipe/shape_test.go
git commit -m "feat: add a slot axis to the recipe shape"
```

---

### Task 2: Slot-aware class encoding and longest-match decoding

**Files:**
- Modify: `internal/recipe/shape.go`
- Test: `internal/recipe/shape_test.go`

**Interfaces:**
- Consumes: `Shape`, `Slot` (Task 1)
- Produces: `func (Shape) BaseClass(slot string) string`, `func (Shape) ValueClass(slot, dimension, value string) string`, `func (Shape) DecodeClass(class string) (slot, dimension, value string, kind ClassKind, err error)`

Longest-match slot resolution is the point of this task. Sidebar's `menu`, `menu-button` and `menu-button-tooltip-content` all coexist.

- [ ] **Step 1: Write the failing test**

```go
func sidebarShape() Shape {
	return Shape{
		Component: "sidebar",
		Slots: []Slot{
			{Name: "", Base: true},
			{Name: "menu", Base: true},
			{Name: "menu-button", Base: true, Dimensions: []Dimension{
				{Name: "size", Default: "default", Values: []string{"default", "lg"}},
			}},
			{Name: "menu-button-tooltip-content", Base: true},
		},
	}
}

func TestShapeSlotClassEncoding(t *testing.T) {
	t.Parallel()
	s := sidebarShape()
	for _, tt := range []struct{ got, want string }{
		{s.BaseClass(""), "gsxui-recipe-sidebar"},
		{s.BaseClass("menu-button"), "gsxui-recipe-sidebar-menu-button"},
		{s.ValueClass("menu-button", "size", "lg"), "gsxui-recipe-sidebar-menu-button-size-lg"},
		{s.ValueClass("", "size", "lg"), "gsxui-recipe-sidebar-size-lg"},
	} {
		if tt.got != tt.want {
			t.Errorf("got %q, want %q", tt.got, tt.want)
		}
	}
}

func TestShapeDecodeClassPrefersTheLongestSlot(t *testing.T) {
	t.Parallel()
	s := sidebarShape()
	tests := []struct {
		class                   string
		wantKind                ClassKind
		wantSlot, wantDim, wantVal string
	}{
		{"gsxui-recipe-sidebar", ClassBase, "", "", ""},
		{"gsxui-recipe-sidebar-menu", ClassBase, "menu", "", ""},
		// "menu-button" must win over "menu", or the remainder is nonsense.
		{"gsxui-recipe-sidebar-menu-button", ClassBase, "menu-button", "", ""},
		// The longest slot of all, three segments deep.
		{"gsxui-recipe-sidebar-menu-button-tooltip-content", ClassBase, "menu-button-tooltip-content", "", ""},
		// Slot plus dimension: longest slot first, then the dimension-value pair.
		{"gsxui-recipe-sidebar-menu-button-size-lg", ClassValue, "menu-button", "size", "lg"},
	}
	for _, tt := range tests {
		t.Run(tt.class, func(t *testing.T) {
			t.Parallel()
			slot, dim, val, kind, err := s.DecodeClass(tt.class)
			if err != nil {
				t.Fatalf("DecodeClass() error = %v", err)
			}
			if kind != tt.wantKind || slot != tt.wantSlot || dim != tt.wantDim || val != tt.wantVal {
				t.Errorf("DecodeClass() = (%q,%q,%q,%v), want (%q,%q,%q,%v)",
					slot, dim, val, kind, tt.wantSlot, tt.wantDim, tt.wantVal, tt.wantKind)
			}
		})
	}
}

func TestShapeDecodeClassRejectsSlotErrors(t *testing.T) {
	t.Parallel()
	s := sidebarShape()
	for _, class := range []string{
		"gsxui-recipe-sidebar-nosuchslot",
		"gsxui-recipe-sidebar-menu-button-size-xl", // undeclared value
		"gsxui-recipe-sidebar-menu-tone-quiet",     // undeclared dimension on that slot
		"gsxui-recipe-card",                        // wrong component
	} {
		t.Run(class, func(t *testing.T) {
			t.Parallel()
			if _, _, _, _, err := s.DecodeClass(class); err == nil {
				t.Errorf("DecodeClass(%q) = nil error, want error", class)
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/recipe/ -run TestShapeSlotClassEncoding -v`
Expected: FAIL — `too many arguments in call to s.BaseClass`.

- [ ] **Step 3: Write the implementation**

```go
func (s Shape) BaseClass(slot string) string {
	if slot == "" {
		return Prefix + s.Component
	}
	return Prefix + s.Component + "-" + slot
}

func (s Shape) ValueClass(slot, dimension, value string) string {
	return s.BaseClass(slot) + "-" + dimension + "-" + value
}

// DecodeClass resolves a recipe class against the shape. Slot names are matched
// LONGEST FIRST: sidebar declares "menu", "menu-button" and
// "menu-button-tooltip-content", so a shortest-match decode would assign
// "…-menu-button" to slot "menu" with a nonsense remainder.
func (s Shape) DecodeClass(class string) (slot, dimension, value string, kind ClassKind, err error) {
	if !strings.HasPrefix(class, Prefix) {
		return "", "", "", 0, fmt.Errorf("%q is not a recipe class", class)
	}
	if class == s.BaseClass("") {
		return "", "", "", ClassBase, nil
	}
	rest, ok := strings.CutPrefix(class, s.BaseClass("")+"-")
	if !ok {
		return "", "", "", 0, fmt.Errorf("recipe class %q does not belong to component %q", class, s.Component)
	}

	ordered := slices.Clone(s.Slots)
	sort.Slice(ordered, func(i, j int) bool { return len(ordered[i].Name) > len(ordered[j].Name) })

	for _, candidate := range ordered {
		remainder := rest
		if candidate.Name != "" {
			if remainder == candidate.Name {
				return candidate.Name, "", "", ClassBase, nil
			}
			trimmed, ok := strings.CutPrefix(remainder, candidate.Name+"-")
			if !ok {
				continue
			}
			remainder = trimmed
		}
		for _, dim := range candidate.Dimensions {
			suffix, ok := strings.CutPrefix(remainder, dim.Name+"-")
			if !ok {
				continue
			}
			if !dim.Has(suffix) {
				return "", "", "", 0, fmt.Errorf(
					"recipe class %q: slot %q dimension %q does not declare value %q",
					class, candidate.Name, dim.Name, suffix)
			}
			return candidate.Name, dim.Name, suffix, ClassValue, nil
		}
	}
	return "", "", "", 0, fmt.Errorf("recipe class %q names no declared slot or dimension of %q", class, s.Component)
}
```

Add `slices` and `sort` to the import block.

- [ ] **Step 4: Update the pre-existing encoding tests**

`TestShapeClassEncoding`, `TestShapeDecodeClass` and `TestShapeDecodeClassRejects` from the previous design call the old signatures. Update the call sites to pass `""` for the root slot and accept the extra return value. Keep every existing case — they still describe correct root-slot behavior.

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test ./internal/recipe/ -v`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/recipe/shape.go internal/recipe/shape_test.go
git commit -m "feat: encode and decode recipe classes with a slot axis"
```

---

### Task 3: Slot-aware conformance and conflict detection

**Files:**
- Modify: `internal/recipe/validate.go`
- Test: `internal/recipe/validate_test.go`

**Interfaces:**
- Consumes: `Shape`, `Slot`, `DecodeClass` (Tasks 1-2), `ParseStyle`/`Style`/`Rule`
- Produces: `Resolved{Shape Shape; Base map[string][]string; Values map[string]map[string]map[string][]string}`, `func (Resolved) BaseUtilities(slot string) []string`, `func (Resolved) Utilities(slot, dimension, value string) []string`, `Conform`, `CheckConflicts` (signatures unchanged)

`Base` is keyed by slot. `Values` is keyed `[slot][dimension][value]`.

- [ ] **Step 1: Write the failing test**

```go
func slotConformShape() Shape {
	return Shape{
		Component: "card",
		Slots: []Slot{
			{Name: "", Base: true},
			{Name: "header", Base: true, Dimensions: []Dimension{
				{Name: "variant", Default: "default", Values: []string{"default", "muted"}},
			}},
		},
	}
}

const slotConformCSS = `@layer components {
  .gsxui-recipe-card { @apply rounded-xl border; }
  .gsxui-recipe-card-header { @apply flex gap-1.5; }
  .gsxui-recipe-card-header-variant-default { @apply text-foreground; }
  .gsxui-recipe-card-header-variant-muted { @apply text-muted-foreground; }
}`

func TestConformAcceptsSlottedStyle(t *testing.T) {
	t.Parallel()
	resolved, err := Conform("nova/card.css", slotConformShape(), mustParse(t, slotConformCSS))
	if err != nil {
		t.Fatalf("Conform() error = %v", err)
	}
	if got, want := resolved.BaseUtilities(""), []string{"rounded-xl", "border"}; !reflect.DeepEqual(got, want) {
		t.Errorf("BaseUtilities(root) = %q, want %q", got, want)
	}
	if got, want := resolved.BaseUtilities("header"), []string{"flex", "gap-1.5"}; !reflect.DeepEqual(got, want) {
		t.Errorf("BaseUtilities(header) = %q, want %q", got, want)
	}
	if got, want := resolved.Utilities("header", "variant", "muted"), []string{"text-muted-foreground"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Utilities() = %q, want %q", got, want)
	}
}

func TestConformRejectsMissingSlotBase(t *testing.T) {
	t.Parallel()
	src := strings.Replace(slotConformCSS,
		"  .gsxui-recipe-card-header { @apply flex gap-1.5; }\n", "", 1)
	_, err := Conform("maia/card.css", slotConformShape(), mustParse(t, src))
	if err == nil {
		t.Fatal("Conform() = nil, want error")
	}
	want := `maia/card.css: slot "header" missing base rule .gsxui-recipe-card-header`
	if got := err.Error(); got != want {
		t.Errorf("Conform() = %q, want %q", got, want)
	}
}

func TestConformRejectsMissingSlotValue(t *testing.T) {
	t.Parallel()
	src := strings.Replace(slotConformCSS,
		"  .gsxui-recipe-card-header-variant-muted { @apply text-muted-foreground; }\n", "", 1)
	_, err := Conform("maia/card.css", slotConformShape(), mustParse(t, src))
	if err == nil {
		t.Fatal("Conform() = nil, want error")
	}
	want := `maia/card.css: slot "header" dimension "variant" missing value "muted"`
	if got := err.Error(); got != want {
		t.Errorf("Conform() = %q, want %q", got, want)
	}
}

func TestCheckConflictsNamesTheSlot(t *testing.T) {
	t.Parallel()
	src := strings.Replace(slotConformCSS,
		"  .gsxui-recipe-card-header { @apply flex gap-1.5; }",
		"  .gsxui-recipe-card-header { @apply rounded-lg rounded-md; }", 1)
	resolved, err := Conform("nova/card.css", slotConformShape(), mustParse(t, src))
	if err != nil {
		t.Fatal(err)
	}
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
	err = CheckConflicts("nova/card.css", resolved, stub)
	if err == nil {
		t.Fatal("CheckConflicts() = nil, want error")
	}
	want := "nova/card.css: .gsxui-recipe-card-header applies conflicting utilities: rounded-lg is superseded"
	if got := err.Error(); got != want {
		t.Errorf("CheckConflicts() = %q, want %q", got, want)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/recipe/ -run TestConformAcceptsSlottedStyle -v`
Expected: FAIL — `resolved.BaseUtilities undefined`.

- [ ] **Step 3: Write the implementation**

```go
// Resolved is a style proven to be a complete implementation of a shape.
type Resolved struct {
	Shape  Shape
	Base   map[string][]string                       // [slot]
	Values map[string]map[string]map[string][]string // [slot][dimension][value]
}

func (r Resolved) BaseUtilities(slot string) []string { return slices.Clone(r.Base[slot]) }

func (r Resolved) Utilities(slot, dimension, value string) []string {
	return slices.Clone(r.Values[slot][dimension][value])
}

func Conform(filename string, shape Shape, style Style) (Resolved, error) {
	if err := shape.Validate(); err != nil {
		return Resolved{}, err
	}

	resolved := Resolved{
		Shape:  shape,
		Base:   make(map[string][]string, len(shape.Slots)),
		Values: make(map[string]map[string]map[string][]string, len(shape.Slots)),
	}
	for _, slot := range shape.Slots {
		byDimension := make(map[string]map[string][]string, len(slot.Dimensions))
		for _, dimension := range slot.Dimensions {
			byDimension[dimension.Name] = make(map[string][]string, len(dimension.Values))
		}
		resolved.Values[slot.Name] = byDimension
	}

	// Style to shape: every rule must be declared.
	for _, rule := range style.Rules() {
		slot, dimension, value, kind, err := shape.DecodeClass(rule.Class)
		if err != nil {
			return Resolved{}, fmt.Errorf("%s: %w", filename, err)
		}
		if kind == ClassBase {
			declared, _ := shape.Slot(slot)
			if !declared.Base {
				return Resolved{}, fmt.Errorf("%s: slot %q declares no base rule, found %s",
					filename, slot, rule.Class)
			}
			resolved.Base[slot] = slices.Clone(rule.Utilities)
			continue
		}
		resolved.Values[slot][dimension][value] = slices.Clone(rule.Utilities)
	}

	// Shape to style: every declaration must be supplied.
	for _, slot := range shape.Slots {
		if slot.Base {
			if _, ok := resolved.Base[slot.Name]; !ok {
				return Resolved{}, fmt.Errorf("%s: slot %q missing base rule .%s",
					filename, slot.Name, shape.BaseClass(slot.Name))
			}
		}
		for _, dimension := range slot.Dimensions {
			for _, value := range dimension.Values {
				if _, ok := resolved.Values[slot.Name][dimension.Name][value]; !ok {
					return Resolved{}, fmt.Errorf("%s: slot %q dimension %q missing value %q",
						filename, slot.Name, dimension.Name, value)
				}
			}
		}
	}
	return resolved, nil
}
```

`CheckConflicts` keeps its existing `check(class string, utilities []string) error` closure and its superseded-before-repeat ordering verbatim. Only the traversal changes — walk `shape.Slots`, checking `resolved.Base[slot.Name]` under `shape.BaseClass(slot.Name)` and each value under `shape.ValueClass(slot.Name, dim.Name, value)`.

- [ ] **Step 4: Update the pre-existing validate tests**

`conformShape()`, `conformCSS` and their dependants use the flat form. Rewrite `conformShape()` as a single-root-slot shape and update expected messages to the slot-qualified form. Keep every case.

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test ./internal/recipe/ -v`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/recipe/validate.go internal/recipe/validate_test.go
git commit -m "feat: validate style conformance per slot"
```

---

### Task 4: Contract schema v2

**Files:**
- Modify: `internal/recipe/contract.go`
- Test: `internal/recipe/contract_test.go`

**Interfaces:**
- Consumes: `Shape`, `Slot`, `Resolved` (Tasks 1-3)
- Produces: `const ContractVersion = 2`, `ContractComponent{Slots map[string]ContractSlot}`, `ContractSlot{Base bool; Dimensions map[string]ContractDimension}`, `ContractUtilities{Slots map[string]ContractSlotUtilities}`, `ContractSlotUtilities{Base []string; Dimensions map[string]map[string][]string}`

The shape is still emitted once under `components`; `styles` still carries only utilities.

- [ ] **Step 1: Write the failing test**

```go
func TestBuildContractGroupsSlotsUnderComponent(t *testing.T) {
	t.Parallel()
	shape := slotConformShape()
	resolved, err := Conform("nova/card.css", shape, mustParse(t, slotConformCSS))
	if err != nil {
		t.Fatal(err)
	}
	out, err := BuildContract(
		map[string]Shape{"card": shape},
		map[string]map[string]Resolved{"nova": {"card": resolved}},
	).MarshalIndent()
	if err != nil {
		t.Fatal(err)
	}

	var parsed Contract
	if err := json.Unmarshal(out, &parsed); err != nil {
		t.Fatalf("contract is not valid JSON: %v", err)
	}
	if parsed.Version != 2 {
		t.Errorf("Version = %d, want 2", parsed.Version)
	}
	card, ok := parsed.Components["card"]
	if !ok {
		t.Fatal("components.card missing")
	}
	if _, ok := card.Slots["header"]; !ok {
		t.Error("components.card.slots.header missing — slots must group under their component")
	}
	if got := card.Slots["header"].Dimensions["variant"].Default; got != "default" {
		t.Errorf("header variant default = %q, want %q", got, "default")
	}
	// The shape must appear once, under components — never repeated per style.
	if got, want := strings.Count(string(out), `"values"`), 1; got != want {
		t.Errorf("shape emitted %d times, want %d:\n%s", got, want, out)
	}
	nova := parsed.Styles["nova"]["card"]
	if got, want := nova.Slots["header"].Dimensions["variant"]["muted"], []string{"text-muted-foreground"}; !reflect.DeepEqual(got, want) {
		t.Errorf("nova header muted = %q, want %q", got, want)
	}
}

func TestContractIsDeterministicWithSlots(t *testing.T) {
	t.Parallel()
	shape := slotConformShape()
	resolved, err := Conform("nova/card.css", shape, mustParse(t, slotConformCSS))
	if err != nil {
		t.Fatal(err)
	}
	build := func() []byte {
		out, err := BuildContract(
			map[string]Shape{"card": shape},
			map[string]map[string]Resolved{"nova": {"card": resolved}, "maia": {"card": resolved}},
		).MarshalIndent()
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

Run: `go test ./internal/recipe/ -run TestBuildContractGroupsSlots -v`
Expected: FAIL — `card.Slots undefined`.

- [ ] **Step 3: Write the implementation**

```go
// ContractVersion is bumped whenever the emitted schema changes shape.
// v2 introduced the slot axis.
const ContractVersion = 2

type ContractComponent struct {
	Slots map[string]ContractSlot `json:"slots"`
}

type ContractSlot struct {
	Base       bool                         `json:"base"`
	Dimensions map[string]ContractDimension `json:"dimensions"`
}

type ContractUtilities struct {
	Slots map[string]ContractSlotUtilities `json:"slots"`
}

type ContractSlotUtilities struct {
	Base       []string                       `json:"base,omitempty"`
	Dimensions map[string]map[string][]string `json:"dimensions"`
}
```

`Contract`, `ContractDimension` and `MarshalIndent` are unchanged. `BuildContract` walks slots for both halves. The root slot's key in the JSON maps is the empty string — that is deliberate and round-trips through `encoding/json` correctly.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/recipe/ -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/recipe/contract.go internal/recipe/contract_test.go
git commit -m "feat: group slots under their component in contract v2"
```

---

### Task 5: The shapes sub-package

**Files:**
- Create: `registry/canonical/shapes/shapes.go`
- Create: `registry/canonical/shapes/button.go`
- Create: `registry/canonical/shapes/shapes_test.go`
- Modify: `registry/canonical/recipe.go`, `registry/canonical/button_recipe.go`

**Interfaces:**
- Consumes: `recipe.Shape`, `recipe.Slot`, `recipe.Dimension`
- Produces: `shapes.Button` (exported `recipe.Shape`), `func shapes.All() map[string]recipe.Shape`

This package exists to break a bootstrap cycle: `stylegen` imports it to read shapes and generate accessors into `registry/canonical`, which cannot compile until those accessors exist. `shapes` has no dependents inside `registry/canonical`, so it always compiles.

- [ ] **Step 1: Write the failing test**

`registry/canonical/shapes/shapes_test.go`:

```go
package shapes

import "testing"

func TestAllShapesAreValid(t *testing.T) {
	t.Parallel()
	all := All()
	if len(all) == 0 {
		t.Fatal("All() is empty")
	}
	for name, shape := range all {
		if err := shape.Validate(); err != nil {
			t.Errorf("All()[%q].Validate() = %v", name, err)
		}
		if shape.Component != name {
			t.Errorf("All()[%q].Component = %q", name, shape.Component)
		}
	}
}

func TestButtonIsASingleRootSlot(t *testing.T) {
	t.Parallel()
	if got, want := len(Button.Slots), 1; got != want {
		t.Fatalf("Button.Slots = %d, want %d", got, want)
	}
	root := Button.Slots[0]
	if root.Name != "" {
		t.Errorf("Button root slot name = %q, want \"\"", root.Name)
	}
	if !root.Base {
		t.Error("Button root slot must declare a base rule")
	}
	variant, ok := root.Dimension("variant")
	if !ok {
		t.Fatal("Button root slot has no variant dimension")
	}
	want := []string{"default", "destructive", "outline", "secondary", "ghost", "link"}
	if len(variant.Values) != len(want) {
		t.Fatalf("variant values = %q, want %q", variant.Values, want)
	}
	for i, value := range want {
		if variant.Values[i] != value {
			t.Errorf("variant value %d = %q, want %q", i, variant.Values[i], value)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./registry/canonical/shapes/ -v`
Expected: FAIL — package does not exist.

- [ ] **Step 3: Write the implementation**

`registry/canonical/shapes/button.go`:

```go
package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Button is a single-element component: one root slot carrying both dimensions.
// Its dimensions mirror the public variant and size parameters of button.gsx.
var Button = recipe.Shape{
	Component: "button",
	Slots: []recipe.Slot{{
		Name: "",
		Base: true,
		Dimensions: []recipe.Dimension{
			{Name: "variant", Default: "default", Values: []string{
				"default", "destructive", "outline", "secondary", "ghost", "link"}},
			{Name: "size", Default: "default", Values: []string{
				"default", "xs", "sm", "lg", "icon", "icon-xs", "icon-sm", "icon-lg"}},
		},
	}},
}
```

`registry/canonical/shapes/shapes.go`:

```go
// Package shapes holds every component's style interface as pure data.
//
// It is a leaf package on purpose. internal/stylegen imports it to read shapes
// and to GENERATE the per-slot accessors in registry/canonical — which cannot
// compile until those accessors exist. Keeping shapes here means the generator
// always has something that compiles to read from.
//
// Nothing here may import registry/canonical, and internal/recipe must never
// import this package.
package shapes

import (
	"maps"

	"github.com/gsxhq/gsxui/internal/recipe"
)

var all = map[string]recipe.Shape{
	Button.Component: Button,
}

// All returns every declared component shape, keyed by component name.
func All() map[string]recipe.Shape {
	out := make(map[string]recipe.Shape, len(all))
	maps.Copy(out, all)
	return out
}
```

`registry/canonical/button_recipe.go` keeps only the binding:

```go
package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// button binds Button's shape to the accessor calls button.gsx authors. The
// variable name is the component name: stylegen resolves button.Root() by
// looking "button" up in shapes.All().
var button = recipe.Component{Shape: shapes.Button}
```

`registry/canonical/recipe.go` delegates:

```go
// Shapes returns every declared component shape, keyed by component name.
func Shapes() map[string]recipe.Shape { return shapes.All() }
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./registry/canonical/shapes/ -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add registry/canonical/shapes registry/canonical/recipe.go registry/canonical/button_recipe.go
git commit -m "feat: move recipe shapes to a leaf package"
```

---

### Task 6: Generated per-slot accessors

**Files:**
- Create: `internal/stylegen/accessors.go`
- Create: `internal/stylegen/accessors_test.go`
- Modify: `internal/recipe/component.go`
- Create: `registry/canonical/button_recipe.gen.go` (generated)
- Modify: `registry/canonical/button.gsx`

**Interfaces:**
- Consumes: `shapes.All()` (Task 5), `recipe.Shape`/`Slot`
- Produces: `func GenerateAccessors(shape recipe.Shape) ([]byte, error)`, and on `recipe.Component`: `func (Component) SlotClass(slot string) string`, `func (Component) SlotValueClass(slot, dimension, value string) string`

`recipe.Component`'s fixed `Role`/`Variant`/`Size` are replaced by two generic methods; the *typed* per-slot API is the generated code.

Accessor naming: kebab segments title-cased. Root slot base is `Root()`. A named slot's base is `<Slot>()`. A dimension is `<Slot><Dimension>(v)`, with the slot part omitted for the root.

- [ ] **Step 1: Write the failing test**

```go
func TestGenerateAccessorsForASingleRootSlot(t *testing.T) {
	t.Parallel()
	src, err := GenerateAccessors(shapes.Button)
	if err != nil {
		t.Fatalf("GenerateAccessors() error = %v", err)
	}
	got := string(src)
	for _, want := range []string{
		"// Code generated by stylegen; DO NOT EDIT.",
		"package canonical",
		"type buttonRecipe struct",
		"func (r buttonRecipe) Root() string",
		"func (r buttonRecipe) Variant(value string) string",
		"func (r buttonRecipe) Size(value string) string",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if _, err := goparser.ParseFile(token.NewFileSet(), "gen.go", src, 0); err != nil {
		t.Fatalf("generated source does not parse: %v", err)
	}
}

func TestGenerateAccessorsTitleCasesKebabSlots(t *testing.T) {
	t.Parallel()
	shape := recipe.Shape{
		Component: "sidebar",
		Slots: []recipe.Slot{
			{Name: "", Base: true},
			{Name: "menu-button", Base: true, Dimensions: []recipe.Dimension{
				{Name: "size", Default: "default", Values: []string{"default", "lg"}},
			}},
			{Name: "menu-button-tooltip-content", Base: true},
		},
	}
	src, err := GenerateAccessors(shape)
	if err != nil {
		t.Fatalf("GenerateAccessors() error = %v", err)
	}
	got := string(src)
	for _, want := range []string{
		"func (r sidebarRecipe) Root() string",
		"func (r sidebarRecipe) MenuButton() string",
		"func (r sidebarRecipe) MenuButtonSize(value string) string",
		"func (r sidebarRecipe) MenuButtonTooltipContent() string",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if _, err := goparser.ParseFile(token.NewFileSet(), "gen.go", src, 0); err != nil {
		t.Fatalf("generated source does not parse: %v", err)
	}
}

func TestGeneratedAccessorsResolveThroughTheShape(t *testing.T) {
	t.Parallel()
	// The generated bodies delegate to recipe.Component, so an unrecognized
	// value must still fall back to the dimension's declared default.
	c := recipe.Component{Shape: shapes.Button}
	if got, want := c.SlotClass(""), "gsxui-recipe-button"; got != want {
		t.Errorf("SlotClass() = %q, want %q", got, want)
	}
	if got, want := c.SlotValueClass("", "variant", "destructve"),
		"gsxui-recipe-button-variant-default"; got != want {
		t.Errorf("unrecognized value must resolve to the default: got %q, want %q", got, want)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/stylegen/ -run TestGenerateAccessors -v`
Expected: FAIL — `undefined: GenerateAccessors`.

- [ ] **Step 3: Rewrite recipe.Component**

```go
// Component binds a Shape to the accessor calls a canonical component authors.
// stylegen replaces every accessor call with concrete style source; the bodies
// exist so the canonical type-checks and its behavior tests can run.
//
// The typed, per-slot API is the GENERATED accessor type in registry/canonical.
// These two methods are what that generated code delegates to.
type Component struct {
	Shape Shape
}

// SlotClass is a slot's base recipe class.
func (c Component) SlotClass(slot string) string { return c.Shape.BaseClass(slot) }

// SlotValueClass is a slot dimension's recipe class. An empty or unrecognized
// value resolves to the dimension's declared default, matching the generated
// switch's default arm — so a behavior test written against the canonical
// asserts something true of every style.
func (c Component) SlotValueClass(slot, dimension, value string) string {
	declared, ok := c.Shape.Slot(slot)
	if !ok {
		panic(fmt.Sprintf("recipe: component %q declares no slot %q", c.Shape.Component, slot))
	}
	dim, ok := declared.Dimension(dimension)
	if !ok {
		panic(fmt.Sprintf("recipe: component %q slot %q declares no dimension %q",
			c.Shape.Component, slot, dimension))
	}
	if !dim.Has(value) {
		value = dim.Default
	}
	return c.Shape.ValueClass(slot, dimension, value)
}
```

- [ ] **Step 4: Write the generator**

`internal/stylegen/accessors.go`. Build the source with `text/template` or a `strings.Builder`, then run it through `format.Source` so output is gofmt-clean by construction:

```go
// accessorName title-cases kebab segments: "menu-button" -> "MenuButton".
func accessorName(slot string) string {
	if slot == "" {
		return "Root"
	}
	var out strings.Builder
	for _, segment := range strings.Split(slot, "-") {
		if segment == "" {
			continue
		}
		out.WriteString(strings.ToUpper(segment[:1]))
		out.WriteString(segment[1:])
	}
	return out.String()
}
```

`dimensionAccessorName` is the second half of the naming rule, and Task 7's
matcher imports it — define it here, once:

```go
// dimensionAccessorName names a slot dimension's accessor. The root slot
// contributes nothing, so button's variant accessor is Variant, not RootVariant.
func dimensionAccessorName(slot, dimension string) string {
	if slot == "" {
		return accessorName(dimension)
	}
	return accessorName(slot) + accessorName(dimension)
}
```

Note `accessorName("")` returns `"Root"`, so `dimensionAccessorName` must
special-case the root rather than concatenating.

The emitter, in full:

```go
// GenerateAccessors emits the typed per-slot accessor type for one component.
func GenerateAccessors(shape recipe.Shape) ([]byte, error) {
	if err := shape.Validate(); err != nil {
		return nil, err
	}
	typeName := shape.Component + "Recipe"

	var out strings.Builder
	out.WriteString("// Code generated by stylegen; DO NOT EDIT.\n\n")
	out.WriteString("package canonical\n\n")
	out.WriteString(`import "github.com/gsxhq/gsxui/internal/recipe"` + "\n\n")
	fmt.Fprintf(&out, "// %s is %s's typed slot accessor set.\n", typeName, shape.Component)
	fmt.Fprintf(&out, "type %s struct{ c recipe.Component }\n\n", typeName)

	slots := slices.Clone(shape.Slots)
	sort.Slice(slots, func(i, j int) bool { return slots[i].Name < slots[j].Name })
	for _, slot := range slots {
		if slot.Base {
			fmt.Fprintf(&out, "func (r %s) %s() string { return r.c.SlotClass(%s) }\n\n",
				typeName, accessorName(slot.Name), strconv.Quote(slot.Name))
		}
		dimensions := slices.Clone(slot.Dimensions)
		sort.Slice(dimensions, func(i, j int) bool { return dimensions[i].Name < dimensions[j].Name })
		for _, dimension := range dimensions {
			fmt.Fprintf(&out,
				"func (r %s) %s(value string) string { return r.c.SlotValueClass(%s, %s, value) }\n\n",
				typeName, dimensionAccessorName(slot.Name, dimension.Name),
				strconv.Quote(slot.Name), strconv.Quote(dimension.Name))
		}
	}

	formatted, err := format.Source([]byte(out.String()))
	if err != nil {
		return nil, fmt.Errorf("format generated accessors for %s: %w", shape.Component, err)
	}
	return formatted, nil
}
```

Imports: `fmt`, `go/format`, `slices`, `sort`, `strconv`, `strings`, plus
`internal/recipe`. Sorting slots and dimensions by name is what makes the output
deterministic, which the drift check requires.

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test ./internal/stylegen/ -run TestGenerate -v`
Expected: PASS.

- [ ] **Step 6: Switch Button's authored source to the generated accessors**

In `registry/canonical/button.gsx`, in BOTH class attributes (the `<a>` branch and the `<button>` branch), change `button.Role()` to `button.Root()`. Leave `button.Variant(variant)` and `button.Size(size)` as they are — the generated names match.

In `registry/canonical/button_recipe.go`, change the binding to the generated type:

```go
var button = buttonRecipe{c: recipe.Component{Shape: shapes.Button}}
```

- [ ] **Step 7: Commit**

```bash
git add internal/stylegen/accessors.go internal/stylegen/accessors_test.go \
        internal/recipe/component.go registry/canonical/button_recipe.go \
        registry/canonical/button.gsx registry/canonical/button_recipe.gen.go
git commit -m "feat: generate per-slot recipe accessors"
```

---

### Task 7: Slot-aware desugaring, generation, and Button byte-identity

**Files:**
- Modify: `internal/stylegen/resolve.go`
- Modify: `internal/stylegen/generate.go`
- Test: `internal/stylegen/resolve_test.go`, `internal/stylegen/generate_test.go`

**Interfaces:**
- Consumes: everything above
- Produces: `Call{Component, Slot, Dimension string}` (was `{Component, Dimension}`), unchanged `Resolve`/`HelperCalls`/`GenerateAll` signatures

The matcher no longer lowercases a method name to get a dimension. It resolves the method name against the shape's accessor names, which is the only way `MenuButtonSize` can split into slot `menu-button` and dimension `size`.

- [ ] **Step 1: Write the failing test**

```go
func TestResolveDesugarsSlotAccessors(t *testing.T) {
	t.Parallel()
	src := []byte(`package canonical

import "github.com/gsxhq/gsx"

component C(children gsx.Node) {
	<div class={ card.Root() }>
		<div class={ card.Header() }>{ children }</div>
	</div>
}
`)
	resolved := testSlotResolved(t) // card: root base, header base
	got, err := Resolve("card.gsx", src, resolved)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	for _, want := range []string{`"rounded-xl border"`, `"flex gap-1.5"`} {
		if !strings.Contains(string(got), want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	for _, forbidden := range []string{recipe.Prefix, "card.Root()", "card.Header()"} {
		if strings.Contains(string(got), forbidden) {
			t.Errorf("resolved source still contains %q:\n%s", forbidden, got)
		}
	}
}

func TestResolveSplitsSlotAndDimensionByAccessorName(t *testing.T) {
	t.Parallel()
	// MenuButtonSize must resolve to slot "menu-button" + dimension "size",
	// which a ToLower of the method name cannot do.
	calls, err := HelperCalls("sidebar.gsx", []byte(`package canonical

import "github.com/gsxhq/gsx"

component S(size string, children gsx.Node) {
	<div class={ sidebar.MenuButtonSize(size) }>{ children }</div>
}
`))
	if err != nil {
		t.Fatalf("HelperCalls() error = %v", err)
	}
	if len(calls) != 1 {
		t.Fatalf("HelperCalls() = %d calls, want 1", len(calls))
	}
	if got, want := calls[0].Slot, "menu-button"; got != want {
		t.Errorf("Slot = %q, want %q", got, want)
	}
	if got, want := calls[0].Dimension, "size"; got != want {
		t.Errorf("Dimension = %q, want %q", got, want)
	}
}
```

Add `testSlotResolved(t)` building a `recipe.Resolved` for the card shape from Task 3's fixture CSS via `ParseStyle` + `Conform`.

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/stylegen/ -run TestResolveDesugarsSlotAccessors -v`
Expected: FAIL — `calls[0].Slot undefined`.

- [ ] **Step 3: Write the implementation**

In `resolve.go`, replace the `strings.ToLower(selector.Sel.Name)` mapping with a shape-driven lookup. Build the accessor table from the shape once (the same `accessorName` function Task 6 defines — export it within the package so both files use one implementation, do NOT duplicate it):

```go
// accessorTarget resolves a generated accessor's method name back to the
// (slot, dimension) it was generated from. Name-splitting cannot do this:
// "MenuButtonSize" is slot "menu-button" + dimension "size", and nothing in
// the string says where the boundary is.
func accessorTarget(shape recipe.Shape, method string) (slot, dimension string, ok bool) {
	for _, s := range shape.Slots {
		if s.Base && accessorName(s.Name) == method {
			return s.Name, "", true
		}
		for _, d := range s.Dimensions {
			if dimensionAccessorName(s.Name, d.Name) == method {
				return s.Name, d.Name, true
			}
		}
	}
	return "", "", false
}
```

An unresolved method is a positioned error naming the component and the method.

Extend `Call` to `{Component, Slot, Dimension string}`. Update `checkHelperCalls` in `generate.go` to validate `(slot, dimension)` against the shape.

In `generate.go`, `resolveAll` gains one more generated artifact per component: `registry/canonical/<component>_recipe.gen.go` from `GenerateAccessors`. It flows through the same validation-before-mutation and `--check` drift machinery as every other artifact — do not special-case it.

Update the residue check's hardcoded method list. Rather than listing `.Root(`, `.Variant(`, `.Size(` and every future accessor, derive it: reject any surviving `<ident>.<Method>(` whose ident matches a known component name. That removes the three-uncoupled-edits problem noted as follow-up 8 of the base design.

- [ ] **Step 4: Verify Button byte-identity — THE CRITICAL CHECK**

```bash
git stash list  # ensure clean
cp ui/button.gsx /tmp/button-before.gsx
cp registry/generated/nova/button.gsx /tmp/nova-before.gsx
cp registry/generated/maia/button.gsx /tmp/maia-before.gsx
go run ./cmd/stylegen
diff /tmp/button-before.gsx ui/button.gsx && echo "ui/button.gsx IDENTICAL"
diff /tmp/nova-before.gsx registry/generated/nova/button.gsx && echo "nova IDENTICAL"
diff /tmp/maia-before.gsx registry/generated/maia/button.gsx && echo "maia IDENTICAL"
```

All three MUST be identical. If any differs, the slot axis is not backward compatible — STOP and report rather than accepting the diff. `registry/generated/recipes.json` WILL differ (schema v2); that is expected and is the only permitted change.

- [ ] **Step 5: Run the full suite**

Run: `go build ./... && go test ./... && go run ./cmd/stylegen --check && gofmt -l . && make audit`
Expected: all green, `--check` exit 0.

- [ ] **Step 6: Commit**

```bash
git add -A internal/stylegen registry ui
git commit -m "feat: desugar slot accessors and generate them per component"
```

---

### Task 8: The layer-precedence audit gate

**Files:**
- Modify: `Makefile`
- Create: `internal/stylegen/layercheck.go`
- Test: `internal/stylegen/layercheck_test.go`

**Interfaces:**
- Consumes: `shapes.All()`
- Produces: `func ComponentComposedMarkers(root string) ([]string, error)`, `func CheckLayerPrecedence(root string) error`

This is Stage 1's first half, and it must land before any bulk migration. The hazard: a rule in `@layer components` that overrides a migrated component's presentation silently stops winning, because the cascade orders layers before specificity. It bit three times during Button alone and twice produced a visible regression that a 313-test browser suite did not catch.

- [ ] **Step 1: Write the failing test**

```go
func TestCheckLayerPrecedenceAcceptsTheCurrentTree(t *testing.T) {
	if err := CheckLayerPrecedence(repoRoot(t)); err != nil {
		t.Fatalf("CheckLayerPrecedence() = %v, want nil", err)
	}
}

func TestCheckLayerPrecedenceRejectsAComponentsLayerOverride(t *testing.T) {
	root := t.TempDir()
	copyRepoFixture(t, root)
	path := filepath.Join(root, "assets", "css", "styles", "default.css")
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	// A components-layer rule setting a property Button emits from utilities.
	injected := bytes.Replace(src, []byte("@layer components {"),
		[]byte("@layer components {\n  [data-gsxui-slot-carousel-previous] { @apply rounded-full; }"), 1)
	if bytes.Equal(injected, src) {
		t.Fatal("fixture did not change — the @layer components anchor moved")
	}
	if err := os.WriteFile(path, injected, 0o644); err != nil {
		t.Fatal(err)
	}
	err = CheckLayerPrecedence(root)
	if err == nil {
		t.Fatal("CheckLayerPrecedence() = nil, want an error for the components-layer override")
	}
	if !strings.Contains(err.Error(), "carousel-previous") {
		t.Errorf("error must name the offending marker, got %q", err)
	}
	if !strings.Contains(err.Error(), "@layer utilities") {
		t.Errorf("error must say what to do, got %q", err)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/stylegen/ -run TestCheckLayerPrecedence -v`
Expected: FAIL — `undefined: CheckLayerPrecedence`.

- [ ] **Step 3: Write the implementation**

`ComponentComposedMarkers` finds every `data-gsxui-slot-*` marker on an element
that renders through a migrated component. Parse each `ui/*.gsx` with
`gsxparser.ParseFile` and walk it with `gsxast.Inspect` — do NOT regex the
source; the repo already rejects regex-based GSX rewriting:

```go
// ComponentComposedMarkers returns every data-gsxui-slot-* marker on an element
// that renders through a migrated component. Those markers are the ones whose
// presentation now comes from the utilities layer, so a components-layer rule
// against them is dead.
func ComponentComposedMarkers(root string) ([]string, error) {
	migrated := map[string]struct{}{}
	for name := range shapes.All() {
		migrated[name] = struct{}{} // "button" matches <ui.Button> and <Button>
	}

	paths, err := filepath.Glob(filepath.Join(root, "ui", "*.gsx"))
	if err != nil {
		return nil, err
	}
	sort.Strings(paths)

	var markers []string
	seen := map[string]struct{}{}
	for _, path := range paths {
		src, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		file, err := gsxparser.ParseFile(token.NewFileSet(), path, src, 0)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		gsxast.Inspect(file, func(node gsxast.Node) bool {
			element, ok := node.(*gsxast.Element)
			if !ok {
				return true
			}
			name := strings.TrimPrefix(element.Name, "ui.")
			if _, ok := migrated[strings.ToLower(name)]; !ok {
				return true
			}
			for _, attr := range element.Attrs {
				marker := attrName(attr)
				if !strings.HasPrefix(marker, "data-gsxui-slot-") {
					continue
				}
				if _, dup := seen[marker]; dup {
					continue
				}
				seen[marker] = struct{}{}
				markers = append(markers, marker)
			}
			return true
		})
	}
	sort.Strings(markers)
	return markers, nil
}
```

Add a small `attrName(gsxast.Attr) string` helper switching over the attribute
node types the package already handles (`*gsxast.StaticAttr`, `*gsxast.ExprAttr`,
`*gsxast.ClassAttr`, `*gsxast.EmbeddedAttr`) and returning their `Name`.

If `gsxast.Element`'s field names differ from `Name`/`Attrs` in the pinned gsx,
read `ast.go` and adapt — do not guess.

`CheckLayerPrecedence` parses `assets/css/styles/default.css` with the existing `recipe.ParseStyle`-adjacent CSS machinery, walks rules inside `@layer components`, and fails when a selector matches a composed marker AND declares any property in the migrated component's utility set (display, border-radius, border, height/size, padding-x, gap, font-size, svg sizing).

The error must name the marker, the property, and the fix:

```text
assets/css/styles/default.css:2643: [data-gsxui-slot-carousel-previous] sets
  border-radius in @layer components, but that element renders through
  ui.Button, whose own utilities win the layer ordering. Move this rule to
  @layer utilities and give it specificity >= (0,1,0) — see
  docs/superpowers/specs/2026-07-29-typed-recipe-model-design.md §9.
```

Report every violation, not just the first — a sweep is more useful than a bisect.

- [ ] **Step 4: Wire it into make audit**

Add to `cmd/stylegen/main.go` a `--check-layers` flag calling `CheckLayerPrecedence`, then add to the `audit:` target:

```make
	go run ./cmd/stylegen --check-layers
```

- [ ] **Step 5: Run tests and the gate**

Run: `go test ./internal/stylegen/ -v && make audit`
Expected: PASS, `make audit` exit 0.

- [ ] **Step 6: Commit**

```bash
git add internal/stylegen/layercheck.go internal/stylegen/layercheck_test.go \
        cmd/stylegen/main.go Makefile
git commit -m "feat: fail the build on components-layer overrides of compiled presentation"
```

---

### Task 9: The committed computed-style sweep harness

**Files:**
- Create: `jstest/support/computed-sweep.ts`
- Create: `jstest/specs/layer-precedence.spec.ts`
- Modify: `Makefile`

**Interfaces:**
- Consumes: the existing `jstest/support/manifest.ts` fixture list
- Produces: `sweepComputedStyles(page, fixtureURL): Promise<Record<string, Record<string,string>>>`, and a `make sweep-baseline` / `make sweep-compare` pair

Stage 1's second half. The Button work built this ad hoc and threw it away; every future migration needs it, and rebuilding it each time is how the Carousel regression survived.

- [ ] **Step 1: Write the harness**

`jstest/support/computed-sweep.ts`:

```ts
export const SWEPT_PROPERTIES = [
  "display", "position", "borderRadius", "borderWidth", "borderColor",
  "width", "height", "paddingLeft", "paddingRight", "paddingTop",
  "paddingBottom", "gap", "fontSize", "fontWeight", "lineHeight",
  "color", "backgroundColor", "opacity", "boxShadow", "textDecorationLine",
  "alignItems", "justifyContent",
] as const;

// sweepComputedStyles records the resting computed style of every marked
// element on a page. It is the only thing that reliably catches CSS-layer
// precedence regressions: they change rendering without failing any assertion.
export async function sweepComputedStyles(page, url: string) {
  await page.goto(url);
  return page.evaluate((props: readonly string[]) => {
    const out: Record<string, Record<string, string>> = {};
    const seen = new Map<string, number>();
    for (const el of document.querySelectorAll("*")) {
      const marker = [...el.attributes]
        .map((a) => a.name)
        .find((n) => n.startsWith("data-gsxui-slot-"));
      if (!marker) continue;
      const n = (seen.get(marker) ?? 0) + 1;
      seen.set(marker, n);
      const computed = getComputedStyle(el);
      const record: Record<string, string> = {};
      for (const p of props) record[p] = computed[p as never];
      out[`${marker}#${n}`] = record;
    }
    return out;
  }, SWEPT_PROPERTIES);
}
```

- [ ] **Step 2: Write the spec that uses it**

`jstest/specs/layer-precedence.spec.ts` sweeps every fixture in both colour schemes and writes the result to `jstest/.tmp/sweep-<scheme>.json`. It asserts nothing on its own — it is a baseline producer. Diffing against a pre-migration baseline is the assertion, and that is what `make sweep-compare` does.

Add a small number of REGRESSION PINS to the same file for the cases already known to have broken, so they cannot silently return:

```ts
test("carousel arrows stay circular", async ({ page }) => {
  await page.goto("/x/carousel/basic");
  const el = page.locator("[data-gsxui-slot-carousel-previous]").first();
  const radius = await el.evaluate((n) => parseFloat(getComputedStyle(n).borderRadius));
  const height = await el.evaluate((n) => n.getBoundingClientRect().height);
  expect(radius).toBeGreaterThanOrEqual(height / 2);
});
```

- [ ] **Step 3: Add the Makefile targets**

```make
sweep-baseline:
	SWEEP_OUT=jstest/.tmp/sweep-baseline npx playwright test --config jstest/playwright.config.ts jstest/specs/layer-precedence.spec.ts

sweep-compare:
	SWEEP_OUT=jstest/.tmp/sweep-current npx playwright test --config jstest/playwright.config.ts jstest/specs/layer-precedence.spec.ts
	node jstest/support/sweep-diff.mjs jstest/.tmp/sweep-baseline jstest/.tmp/sweep-current
```

`jstest/support/sweep-diff.mjs`:

```js
import { readFileSync, readdirSync } from "node:fs";
import { join } from "node:path";

const [baseDir, curDir] = process.argv.slice(2);
let differences = 0;

for (const file of readdirSync(baseDir).filter((f) => f.endsWith(".json"))) {
  const before = JSON.parse(readFileSync(join(baseDir, file), "utf8"));
  const after = JSON.parse(readFileSync(join(curDir, file), "utf8"));
  for (const [element, props] of Object.entries(before)) {
    if (!(element in after)) {
      console.log(`${file}  ${element}  DISAPPEARED`);
      differences++;
      continue;
    }
    for (const [prop, value] of Object.entries(props)) {
      const now = after[element][prop];
      if (now !== value) {
        console.log(`${file}  ${element}  ${prop}\n    before: ${value}\n    after:  ${now}`);
        differences++;
      }
    }
  }
  for (const element of Object.keys(after)) {
    if (!(element in before)) {
      console.log(`${file}  ${element}  APPEARED`);
      differences++;
    }
  }
}

if (differences > 0) {
  console.error(`\n${differences} computed-style difference(s) — each needs a fix or a justification.`);
  process.exit(1);
}
console.log("no computed-style differences");
```

- [ ] **Step 4: Prove the harness catches a real regression**

Temporarily revert the Carousel fix (move `rounded-full` back into `@layer components` with its `:where()`), run `make sweep-compare` against a baseline taken before the revert, and confirm it reports `data-gsxui-slot-carousel-previous#1 borderRadius`. Then restore the fix and confirm a clean compare. Paste both outputs into your report — a sweep harness that has never caught anything is not known to work.

- [ ] **Step 5: Run the full suite**

Run: `npx playwright test --config jstest/playwright.config.ts`
Expected: 1 failed (`sidebar-page.spec.ts:3`, pre-existing at merge base) or better, plus the new pins passing.

- [ ] **Step 6: Commit**

```bash
git add jstest/support/computed-sweep.ts jstest/support/sweep-diff.mjs \
        jstest/specs/layer-precedence.spec.ts Makefile
git commit -m "test: commit the computed-style sweep harness and regression pins"
```

---

### Task 10: Documentation and full verification

**Files:**
- Modify: `docs/superpowers/specs/2026-07-29-slot-axis-and-component-migration-design.md`
- Modify: `docs/superpowers/specs/2026-07-29-typed-recipe-model-design.md`
- Modify: `CHANGELOG.md`

- [ ] **Step 1: Record what implementation changed**

Earlier work in this repo repeatedly shipped specs describing designs that implementation had abandoned. Re-read both specs against the code and correct any statement that is no longer true — in particular the base design's §5 accessor table (`Role()` is now `Root()`, and `Variant`/`Size` are generated rather than fixed methods on `recipe.Component`), and its follow-up items 2 and 8, which this plan closes.

- [ ] **Step 2: Add the CHANGELOG entry**

Contract consumers need to know the schema moved:

```markdown
### Changed
- `registry/generated/recipes.json` is now schema version 2: a component's
  rules are grouped under named slots (`components.<c>.slots.<s>`), so
  multi-slot components can be expressed. Version 1 had no slot axis.
```

- [ ] **Step 3: Full verification**

Run each and confirm the stated expectation before claiming completion:

```bash
go build ./...                       # expect: no output
make audit                           # expect: exit 0, including the new layer gate
go run ./cmd/stylegen --check        # expect: exit 0, no drift
go test ./...                        # expect: all packages ok
gofmt -l .                           # expect: empty
git status --short                   # expect: clean
npx playwright test --config jstest/playwright.config.ts
```

- [ ] **Step 4: Verify the spec's success criteria**

Confirm each of §9's seven criteria with evidence:

1. `Shape` expresses a multi-slot component with per-slot dimensions — point at `sidebarShape()` in the tests.
2. Button emits byte-identical output — re-run Task 7 Step 4's diff and paste it.
3. `sidebar-menu-button-tooltip-content` decodes correctly against its shorter siblings — point at `TestShapeDecodeClassPrefersTheLongestSlot`.
4. Accessors generate with no bootstrap cycle; a typo'd slot is a compile error — demonstrate by introducing `button.Rooot()` and pasting the compile error, then reverting.
5. `recipes.json` v2 groups slots under components and still supports enumerate/diff/completeness with no CSS access — demonstrate with a short script as was done for v1.
6. `make audit` fails on a components-layer override — point at `TestCheckLayerPrecedenceRejectsAComponentsLayerOverride`.
7. The sweep harness is committed and proven to catch a regression — point at Task 9 Step 4's output.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "docs: reconcile specs with the slot axis and record contract v2"
```

---

## Notes for the implementer

**Button byte-identity is the whole proof.** Task 7 Step 4 is not a formality. The slot axis is a refactor of the model, and if Button's generated output moves by a single byte, something about the encoding changed that will silently change every component migrated later. Stop rather than accept a diff.

**Do not name-split accessor methods.** `MenuButtonSize` cannot be decomposed by string manipulation — nothing in it marks the slot/dimension boundary. Resolve method names against the shape's own generated accessor names. The same rule applies to class decoding, which is why slots are matched longest-first.

**The layer gate exists because tests do not see this bug class.** Carousel arrows lost their roundness and InputGroupButton lost a type-scale step, and both passed build, tests, drift check, `make audit` and a 313-test browser suite. If the gate feels redundant while writing it, that is the same intuition that let those two ship.

**Two-thirds of the catalogue has no dimensions.** Most components are pure base-per-slot. Do not add a dimension to a shape just to make it look like Button — an absent dimension is the correct model for most slots.

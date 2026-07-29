package recipe

import (
	"strings"
	"testing"
)

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
		{"no slots", func(s *Shape) { s.Slots = nil }, `button: no slots declared`},
		{"no base and no dimensions", func(s *Shape) {
			s.Slots[0].Base = false
			s.Slots[0].Dimensions = nil
		}, `button: slot "" declares neither a base rule nor any dimension`},
		{"empty dimension name", func(s *Shape) { s.Slots[0].Dimensions[0].Name = "" }, `button: slot "" has an unnamed dimension`},
		{"duplicate dimension", func(s *Shape) { s.Slots[0].Dimensions[1].Name = "variant" }, `button: slot "" duplicate dimension "variant"`},
		{"no values", func(s *Shape) { s.Slots[0].Dimensions[0].Values = nil }, `button: slot "" dimension "variant" declares no values`},
		{"duplicate value", func(s *Shape) {
			s.Slots[0].Dimensions[0].Values = []string{"outline", "outline"}
		}, `button: slot "" dimension "variant" duplicates value "outline"`},
		{"default not a value", func(s *Shape) {
			s.Slots[0].Dimensions[0].Default = "ghost"
		}, `button: slot "" dimension "variant" default "ghost" is not one of its values`},
		{"empty value", func(s *Shape) {
			s.Slots[0].Dimensions[0].Values = []string{"default", ""}
		}, `button: slot "" dimension "variant" has an empty value`},
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

func TestShapeDimensionLookup(t *testing.T) {
	t.Parallel()
	root, ok := validShape().Slot("")
	if !ok {
		t.Fatal("Slot(root) = false, want true")
	}
	d, ok := root.Dimension("size")
	if !ok {
		t.Fatal("Dimension(size) = false, want true")
	}
	if d.Default != "default" {
		t.Errorf("Default = %q, want %q", d.Default, "default")
	}
	if _, ok := root.Dimension("tone"); ok {
		t.Error("Dimension(tone) = true, want false")
	}
}

func TestShapeClassEncoding(t *testing.T) {
	t.Parallel()
	s := validShape()
	if got, want := s.BaseClass(""), "gsxui-recipe-button"; got != want {
		t.Errorf("BaseClass() = %q, want %q", got, want)
	}
	if got, want := s.ValueClass("", "size", "icon-lg"), "gsxui-recipe-button-size-icon-lg"; got != want {
		t.Errorf("ValueClass() = %q, want %q", got, want)
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
			slot, dimension, value, kind, err := s.DecodeClass(tt.class)
			if err != nil {
				t.Fatalf("DecodeClass() error = %v", err)
			}
			if slot != "" {
				t.Errorf("DecodeClass() slot = %q, want the root slot", slot)
			}
			if kind != tt.wantKind || dimension != tt.wantDimension || value != tt.wantValue {
				t.Errorf("DecodeClass() = (%q, %q, %v), want (%q, %q, %v)",
					dimension, value, kind, tt.wantDimension, tt.wantValue, tt.wantKind)
			}
		})
	}
}

func TestShapeDecodeClassPrefersTheLongestSlot(t *testing.T) {
	t.Parallel()
	s := sidebarShape()
	tests := []struct {
		class                      string
		wantKind                   ClassKind
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

func TestShapeDecodeClassRejects(t *testing.T) {
	t.Parallel()
	s := validShape()
	for _, class := range []string{
		"gsxui-recipe-card",                 // wrong component
		"gsxui-recipe-button-variant-plain", // undeclared value
		"gsxui-recipe-button-tone-quiet",    // undeclared dimension
		"inline-flex",                       // not a recipe class
		"gsxui-recipe-",                     // prefix only
	} {
		t.Run(class, func(t *testing.T) {
			t.Parallel()
			if _, _, _, _, err := s.DecodeClass(class); err == nil {
				t.Errorf("DecodeClass(%q) = nil error, want error", class)
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

// collidingShape is the confirmed non-injective encoding: slot "menu" varies
// over "button-size", slot "menu-button" varies over "size", and both reach
// gsxui-recipe-sidebar-menu-button-size-lg.
func collidingShape() Shape {
	return Shape{
		Component: "sidebar",
		Slots: []Slot{
			{Name: "menu", Dimensions: []Dimension{
				{Name: "button-size", Default: "md", Values: []string{"md", "lg", "xl"}},
			}},
			{Name: "menu-button", Dimensions: []Dimension{
				{Name: "size", Default: "default", Values: []string{"default", "lg"}},
			}},
		},
	}
}

// TestShapeValidateRejectsNonInjectiveClassEncoding pins finding 2. Without
// this the shape validates, one of the two axes becomes unreachable, and
// Conform then blames the wrong axis for a rule that is visibly present in the
// CSS file it is reading.
func TestShapeValidateRejectsNonInjectiveClassEncoding(t *testing.T) {
	t.Parallel()

	s := collidingShape()
	if a, b := s.ValueClass("menu", "button-size", "lg"), s.ValueClass("menu-button", "size", "lg"); a != b {
		t.Fatalf("fixture no longer collides: %q vs %q", a, b)
	}
	err := s.Validate()
	if err == nil {
		t.Fatal("Validate() = nil, want an error for the colliding class encoding")
	}
	for _, want := range []string{`slot "menu" dimension "button-size" value "lg"`, `slot "menu-button" dimension "size" value "lg"`} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("Validate() error must name both contributors, missing %q in %q", want, err)
		}
	}
}

// TestShapeValidateRejectsBaseVersusValueCollision is the same defect through
// the other door: a slot whose own name is another slot's value class.
func TestShapeValidateRejectsBaseVersusValueCollision(t *testing.T) {
	t.Parallel()

	s := Shape{
		Component: "sidebar",
		Slots: []Slot{
			{Name: "menu-button-size-lg", Base: true},
			{Name: "menu-button", Dimensions: []Dimension{
				{Name: "size", Default: "default", Values: []string{"default", "lg"}},
			}},
		},
	}
	if err := s.Validate(); err == nil {
		t.Fatal("Validate() = nil, want an error for the base-versus-value collision")
	}
}

// TestShapeDecodeClassSearchesEveryCandidate pins finding 2's second half.
// Longest-match must be a search, not a greedy first-match: a class that one
// candidate near-misses and another decodes exactly must decode.
func TestShapeDecodeClassSearchesEveryCandidate(t *testing.T) {
	t.Parallel()

	s := collidingShape()
	const class = "gsxui-recipe-sidebar-menu-button-size-xl"
	slot, dimension, value, kind, err := s.DecodeClass(class)
	if err != nil {
		t.Fatalf("DecodeClass(%q) error = %v — the longest candidate near-missed and the search stopped", class, err)
	}
	if slot != "menu" || dimension != "button-size" || value != "xl" || kind != ClassValue {
		t.Errorf("DecodeClass(%q) = (%q,%q,%q,%v), want (menu,button-size,xl,ClassValue)",
			class, slot, dimension, value, kind)
	}
}

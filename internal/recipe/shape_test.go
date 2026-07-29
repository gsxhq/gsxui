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
		"gsxui-recipe-card",                 // wrong component
		"gsxui-recipe-button-variant-plain", // undeclared value
		"gsxui-recipe-button-tone-quiet",    // undeclared dimension
		"inline-flex",                       // not a recipe class
		"gsxui-recipe-",                     // prefix only
	} {
		t.Run(class, func(t *testing.T) {
			t.Parallel()
			if _, _, _, err := s.DecodeClass(class); err == nil {
				t.Errorf("DecodeClass(%q) = nil error, want error", class)
			}
		})
	}
}

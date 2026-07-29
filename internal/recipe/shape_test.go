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

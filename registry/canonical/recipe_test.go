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

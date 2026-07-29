package canonical

import "testing"

// The shape-validity and Button-shape assertions live in
// registry/canonical/shapes, next to the data they describe. What remains here
// is the behavior of the GENERATED accessor set bound to that shape.

func TestRootClass(t *testing.T) {
	t.Parallel()
	if got, want := button.Root(), "gsxui-recipe-button"; got != want {
		t.Errorf("Root() = %q, want %q", got, want)
	}
}

func TestDimensionAccessorsResolveDeclaredValues(t *testing.T) {
	t.Parallel()
	if got, want := button.Variant("outline"), "gsxui-recipe-button-variant-outline"; got != want {
		t.Errorf("Variant() = %q, want %q", got, want)
	}
	if got, want := button.Size("icon-lg"), "gsxui-recipe-button-size-icon-lg"; got != want {
		t.Errorf("Size() = %q, want %q", got, want)
	}
}

func TestDimensionAccessorsFallBackToDefault(t *testing.T) {
	t.Parallel()
	// Empty and unrecognized values must both resolve to the declared default,
	// matching the generated switch's default arm.
	for _, value := range []string{"", "destructve"} {
		if got, want := button.Variant(value), "gsxui-recipe-button-variant-default"; got != want {
			t.Errorf("Variant(%q) = %q, want %q", value, got, want)
		}
	}
}

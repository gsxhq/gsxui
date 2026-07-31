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

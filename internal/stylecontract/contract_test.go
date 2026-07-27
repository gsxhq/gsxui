package stylecontract

import (
	"reflect"
	"testing"
)

func TestValidateRejectsInvalidContracts(t *testing.T) {
	tests := []struct {
		name       string
		components []Component
		want       string
	}{
		{
			name:       "empty component name",
			components: []Component{{}},
			want:       "component 0: empty component name",
		},
		{
			name: "duplicate component name",
			components: []Component{
				{Name: "Button"},
				{Name: "Button"},
			},
			want: "component \"Button\": duplicate component name",
		},
		{
			name: "empty slot name",
			components: []Component{{
				Name:  "Button",
				Slots: []Slot{{}},
			}},
			want: "component \"Button\": slot 0: empty slot name",
		},
		{
			name: "duplicate slot token globally",
			components: []Component{
				{Name: "Button", Slots: []Slot{{Name: "button"}}},
				{Name: "AlertDialogAction", Slots: []Slot{{Name: "button"}}},
			},
			want: "component \"AlertDialogAction\": slot \"button\": duplicate slot token (already declared by component \"Button\")",
		},
		{
			name: "empty axis attribute",
			components: []Component{{
				Name: "Button",
				Slots: []Slot{{
					Name: "button",
					Axes: []Axis{{}},
				}},
			}},
			want: "component \"Button\": slot \"button\": axis 0: empty attribute",
		},
		{
			name: "duplicate axis value",
			components: []Component{{
				Name: "Button",
				Slots: []Slot{{
					Name: "button",
					Axes: []Axis{{Attribute: "data-size", Values: []string{"sm", "lg", "sm"}}},
				}},
			}},
			want: "component \"Button\": slot \"button\": axis \"data-size\": duplicate value \"sm\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.components)
			if err == nil {
				t.Fatal("Validate() error = nil")
			}
			if err.Error() != tt.want {
				t.Fatalf("Validate() error = %q, want %q", err, tt.want)
			}
		})
	}
}

func TestValidateAcceptsPresenceOnlyAxis(t *testing.T) {
	err := Validate([]Component{{
		Name: "Button",
		Slots: []Slot{{
			Name: "button",
			Axes: []Axis{{Attribute: "data-disabled"}},
		}},
	}})
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestAllReturnsSortedCopy(t *testing.T) {
	originalPrimitives := primitiveContracts
	defer func() { primitiveContracts = originalPrimitives }()
	primitiveContracts = []Component{
		{Name: "Zebra"},
		{Name: "Alert"},
	}

	got := All()
	want := []Component{
		{Name: "Alert"},
		{Name: "Zebra"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("All() = %#v, want %#v", got, want)
	}

	got[0].Name = "Changed"
	if again := All(); again[0].Name != "Alert" {
		t.Fatalf("All() returned an aliased slice: %#v", again)
	}
}

func TestAllReturnsDeepSnapshot(t *testing.T) {
	originalPrimitives := primitiveContracts
	defer func() { primitiveContracts = originalPrimitives }()
	primitiveContracts = []Component{{
		Name: "Button",
		Slots: []Slot{{
			Name: "button",
			Axes: []Axis{{
				Attribute: "data-size",
				Values:    []string{"sm"},
			}},
		}},
	}}
	want := []Component{{
		Name: "Button",
		Slots: []Slot{{
			Name: "button",
			Axes: []Axis{{
				Attribute: "data-size",
				Values:    []string{"sm"},
			}},
		}},
	}}

	got := All()
	got[0].Name = "Changed"
	got[0].Slots[0].Name = "changed"
	got[0].Slots[0].Axes[0].Attribute = "data-changed"
	got[0].Slots[0].Axes[0].Values[0] = "changed"

	if again := All(); !reflect.DeepEqual(again, want) {
		t.Fatalf("All() returned an aliased nested value: got %#v, want %#v", again, want)
	}
}

func TestAllPreservesSliceNilness(t *testing.T) {
	originalPrimitives := primitiveContracts
	defer func() { primitiveContracts = originalPrimitives }()
	primitiveContracts = []Component{
		{Name: "NilSlots"},
		{Name: "EmptySlots", Slots: []Slot{}},
		{
			Name: "Axes",
			Slots: []Slot{
				{Name: "nil-axes"},
				{Name: "empty-axes", Axes: []Axis{}},
				{
					Name: "values",
					Axes: []Axis{
						{Attribute: "data-presence"},
						{Attribute: "data-empty", Values: []string{}},
					},
				},
			},
		},
	}

	components := make(map[string]Component)
	for _, component := range All() {
		components[component.Name] = component
	}
	if components["NilSlots"].Slots != nil {
		t.Fatalf("nil Slots = %#v, want nil", components["NilSlots"].Slots)
	}
	if components["EmptySlots"].Slots == nil {
		t.Fatal("empty Slots became nil")
	}

	slots := components["Axes"].Slots
	if slots[0].Axes != nil {
		t.Fatalf("nil Axes = %#v, want nil", slots[0].Axes)
	}
	if slots[1].Axes == nil {
		t.Fatal("empty Axes became nil")
	}
	if slots[2].Axes[0].Values != nil {
		t.Fatalf("presence-only Values = %#v, want nil", slots[2].Axes[0].Values)
	}
	if slots[2].Axes[1].Values == nil {
		t.Fatal("empty Values became nil")
	}
}

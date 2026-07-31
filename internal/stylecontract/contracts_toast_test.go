package stylecontract

import (
	"reflect"
	"testing"
)

func TestToastContracts(t *testing.T) {
	want := []Component{
		{
			Name:         "Toaster",
			RegistryName: "toaster",
			Slots: []Slot{
				{Name: "toaster", Axes: []Axis{{Attribute: "data-expanded", Values: []string{"false", "true"}, RuntimeValues: []string{"false", "true"}}}},
			},
		},
		{
			Name:         "Toast",
			RegistryName: "toast",
			Slots: []Slot{
				{Name: "toast", Axes: []Axis{
					{Attribute: "data-type", Values: []string{"default", "success", "info", "warning", "error", "loading"}},
					{Attribute: "data-state", Values: []string{"closed", "open"}, RuntimeValues: []string{"closed", "open"}},
					{Attribute: "data-visible", Values: []string{"false", "true"}, RuntimeValues: []string{"false", "true"}},
				}},
				{Name: "toast-icon"},
				{Name: "toast-content"},
				{Name: "toast-title"},
				{Name: "toast-description"},
				{Name: "toast-action"},
				{Name: "toast-cancel"},
				{Name: "toast-close"},
				{Name: "toast-close-icon"},
			},
		},
	}

	if !reflect.DeepEqual(toastContracts, want) {
		t.Fatalf("toastContracts = %#v\nwant %#v", toastContracts, want)
	}
}

// TestToastContractsUnionMatchesOriginalSonner pins the property that makes
// the toaster/toast split a safe refactor of the former single "sonner"
// contract: regrouping the slots into two registered components must not
// add, drop, rename or reorder a single slot or axis relative to what
// sonnerContracts used to declare. If this ever fails, the split has drifted
// into a redesign, not a regrouping.
func TestToastContractsUnionMatchesOriginalSonner(t *testing.T) {
	originalSonnerSlots := []Slot{
		{Name: "toaster", Axes: []Axis{{Attribute: "data-expanded", Values: []string{"false", "true"}, RuntimeValues: []string{"false", "true"}}}},
		{Name: "toast", Axes: []Axis{
			{Attribute: "data-type", Values: []string{"default", "success", "info", "warning", "error", "loading"}},
			{Attribute: "data-state", Values: []string{"closed", "open"}, RuntimeValues: []string{"closed", "open"}},
			{Attribute: "data-visible", Values: []string{"false", "true"}, RuntimeValues: []string{"false", "true"}},
		}},
		{Name: "toast-icon"},
		{Name: "toast-content"},
		{Name: "toast-title"},
		{Name: "toast-description"},
		{Name: "toast-action"},
		{Name: "toast-cancel"},
		{Name: "toast-close"},
		{Name: "toast-close-icon"},
	}

	var gotSlots []Slot
	for _, component := range toastContracts {
		gotSlots = append(gotSlots, component.Slots...)
	}

	if !reflect.DeepEqual(gotSlots, originalSonnerSlots) {
		t.Fatalf("union of toaster+toast slots = %#v\nwant (original sonner slots) %#v", gotSlots, originalSonnerSlots)
	}

	if len(toastContracts) != 2 {
		t.Fatalf("expected exactly 2 components after the split, got %d", len(toastContracts))
	}
	if toastContracts[0].RegistryName != "toaster" || toastContracts[1].RegistryName != "toast" {
		t.Fatalf("expected registry names [toaster toast], got [%s %s]", toastContracts[0].RegistryName, toastContracts[1].RegistryName)
	}
}

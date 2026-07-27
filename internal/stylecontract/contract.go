package stylecontract

import (
	"fmt"
	"sort"
)

type Axis struct {
	Attribute string
	Values    []string
}

type Slot struct {
	Name string
	Axes []Axis
}

type Component struct {
	Name  string
	Slots []Slot
}

func All() []Component {
	all := make([]Component, 0,
		len(primitiveContracts)+
			len(formContracts)+
			len(overlayContracts)+
			len(menuContracts)+
			len(compositeContracts)+
			len(sidebarContracts)+
			len(sonnerContracts),
	)
	appendComponents := func(components []Component) {
		for _, component := range components {
			all = append(all, cloneComponent(component))
		}
	}
	appendComponents(primitiveContracts)
	appendComponents(formContracts)
	appendComponents(overlayContracts)
	appendComponents(menuContracts)
	appendComponents(compositeContracts)
	appendComponents(sidebarContracts)
	appendComponents(sonnerContracts)
	sort.Slice(all, func(i, j int) bool {
		return all[i].Name < all[j].Name
	})
	return all
}

func cloneComponent(component Component) Component {
	clone := Component{Name: component.Name}
	if component.Slots == nil {
		return clone
	}
	clone.Slots = make([]Slot, len(component.Slots))
	for slotIndex, slot := range component.Slots {
		clone.Slots[slotIndex].Name = slot.Name
		if slot.Axes == nil {
			continue
		}
		clone.Slots[slotIndex].Axes = make([]Axis, len(slot.Axes))
		for axisIndex, axis := range slot.Axes {
			clone.Slots[slotIndex].Axes[axisIndex].Attribute = axis.Attribute
			if axis.Values != nil {
				clone.Slots[slotIndex].Axes[axisIndex].Values = make([]string, len(axis.Values))
				copy(clone.Slots[slotIndex].Axes[axisIndex].Values, axis.Values)
			}
		}
	}
	return clone
}

func Validate(components []Component) error {
	componentNames := make(map[string]struct{}, len(components))
	slotComponents := make(map[string]string)
	for componentIndex, component := range components {
		if component.Name == "" {
			return fmt.Errorf("component %d: empty component name", componentIndex)
		}
		if _, ok := componentNames[component.Name]; ok {
			return fmt.Errorf("component %q: duplicate component name", component.Name)
		}
		componentNames[component.Name] = struct{}{}

		for slotIndex, slot := range component.Slots {
			if slot.Name == "" {
				return fmt.Errorf("component %q: slot %d: empty slot name", component.Name, slotIndex)
			}
			if declaredBy, ok := slotComponents[slot.Name]; ok {
				return fmt.Errorf("component %q: slot %q: duplicate slot token (already declared by component %q)", component.Name, slot.Name, declaredBy)
			}
			slotComponents[slot.Name] = component.Name

			for axisIndex, axis := range slot.Axes {
				if axis.Attribute == "" {
					return fmt.Errorf("component %q: slot %q: axis %d: empty attribute", component.Name, slot.Name, axisIndex)
				}
				values := make(map[string]struct{}, len(axis.Values))
				for _, value := range axis.Values {
					if _, ok := values[value]; ok {
						return fmt.Errorf("component %q: slot %q: axis %q: duplicate value %q", component.Name, slot.Name, axis.Attribute, value)
					}
					values[value] = struct{}{}
				}
			}
		}
	}
	return nil
}

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
	all = append(all, primitiveContracts...)
	all = append(all, formContracts...)
	all = append(all, overlayContracts...)
	all = append(all, menuContracts...)
	all = append(all, compositeContracts...)
	all = append(all, sidebarContracts...)
	all = append(all, sonnerContracts...)
	sort.Slice(all, func(i, j int) bool {
		return all[i].Name < all[j].Name
	})
	return all
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

// Package recipe models component style recipes: the shape a component
// declares, and the utilities a style supplies for it.
package recipe

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

// Shape is a component's style interface: the slots it renders and, per slot,
// the dimensions that slot varies over. Every style must implement all of it.
type Shape struct {
	Component string
	Slots     []Slot
}

// Slot is one styled element of a component. Name is "" for the component's
// root element. Dimensions hang off the slot, not the component, because a
// style may vary one slot and not another.
type Slot struct {
	Name       string
	Base       bool
	Dimensions []Dimension
}

// Dimension is one axis of a Slot. Default must be a member of Values and is
// what an empty or unrecognized value resolves to.
type Dimension struct {
	Name    string
	Default string
	Values  []string
}

func (s Shape) Validate() error {
	if s.Component == "" {
		return fmt.Errorf("component name is empty")
	}
	if len(s.Slots) == 0 {
		return fmt.Errorf("%s: no slots declared", s.Component)
	}
	seen := make(map[string]struct{}, len(s.Slots))
	for _, slot := range s.Slots {
		if _, exists := seen[slot.Name]; exists {
			return fmt.Errorf("%s: duplicate slot %q", s.Component, slot.Name)
		}
		seen[slot.Name] = struct{}{}
		if err := slot.validate(s.Component); err != nil {
			return err
		}
	}
	return nil
}

func (s Slot) validate(component string) error {
	if !s.Base && len(s.Dimensions) == 0 {
		return fmt.Errorf("%s: slot %q declares neither a base rule nor any dimension",
			component, s.Name)
	}
	seen := make(map[string]struct{}, len(s.Dimensions))
	for _, dimension := range s.Dimensions {
		if dimension.Name == "" {
			return fmt.Errorf("%s: slot %q has an unnamed dimension", component, s.Name)
		}
		if _, exists := seen[dimension.Name]; exists {
			return fmt.Errorf("%s: slot %q duplicate dimension %q", component, s.Name, dimension.Name)
		}
		seen[dimension.Name] = struct{}{}
		if err := dimension.validate(component, s.Name); err != nil {
			return err
		}
	}
	return nil
}

func (d Dimension) validate(component, slot string) error {
	if len(d.Values) == 0 {
		return fmt.Errorf("%s: slot %q dimension %q declares no values", component, slot, d.Name)
	}
	seen := make(map[string]struct{}, len(d.Values))
	for _, value := range d.Values {
		if value == "" {
			return fmt.Errorf("%s: slot %q dimension %q has an empty value", component, slot, d.Name)
		}
		if _, exists := seen[value]; exists {
			return fmt.Errorf("%s: slot %q dimension %q duplicates value %q", component, slot, d.Name, value)
		}
		seen[value] = struct{}{}
	}
	if !d.Has(d.Default) {
		return fmt.Errorf("%s: slot %q dimension %q default %q is not one of its values",
			component, slot, d.Name, d.Default)
	}
	return nil
}

// Has reports whether value is declared by this dimension.
func (d Dimension) Has(value string) bool {
	for _, candidate := range d.Values {
		if candidate == value {
			return true
		}
	}
	return false
}

// Slot returns the named slot; "" is the root.
func (s Shape) Slot(name string) (Slot, bool) {
	for _, slot := range s.Slots {
		if slot.Name == name {
			return slot, true
		}
	}
	return Slot{}, false
}

// Dimension returns the named dimension of this slot.
func (s Slot) Dimension(name string) (Dimension, bool) {
	for _, dimension := range s.Dimensions {
		if dimension.Name == name {
			return dimension, true
		}
	}
	return Dimension{}, false
}

// Prefix namespaces every recipe class. It must never survive into a
// generated consumer artifact.
const Prefix = "gsxui-recipe-"

// ClassKind distinguishes a slot's base rule from a dimension value rule.
type ClassKind uint8

const (
	ClassBase ClassKind = iota
	ClassValue
)

func (s Shape) BaseClass(slot string) string {
	if slot == "" {
		return Prefix + s.Component
	}
	return Prefix + s.Component + "-" + slot
}

func (s Shape) ValueClass(slot, dimension, value string) string {
	return s.BaseClass(slot) + "-" + dimension + "-" + value
}

// DecodeClass resolves a recipe class against the shape. Slot names are matched
// LONGEST FIRST: sidebar declares "menu", "menu-button" and
// "menu-button-tooltip-content", so a shortest-match decode would assign
// "…-menu-button" to slot "menu" with a nonsense remainder.
func (s Shape) DecodeClass(class string) (slot, dimension, value string, kind ClassKind, err error) {
	if !strings.HasPrefix(class, Prefix) {
		return "", "", "", 0, fmt.Errorf("%q is not a recipe class", class)
	}
	if class == s.BaseClass("") {
		return "", "", "", ClassBase, nil
	}
	rest, ok := strings.CutPrefix(class, s.BaseClass("")+"-")
	if !ok {
		return "", "", "", 0, fmt.Errorf("recipe class %q does not belong to component %q", class, s.Component)
	}

	ordered := slices.Clone(s.Slots)
	sort.Slice(ordered, func(i, j int) bool { return len(ordered[i].Name) > len(ordered[j].Name) })

	for _, candidate := range ordered {
		remainder := rest
		if candidate.Name != "" {
			if remainder == candidate.Name {
				return candidate.Name, "", "", ClassBase, nil
			}
			trimmed, ok := strings.CutPrefix(remainder, candidate.Name+"-")
			if !ok {
				continue
			}
			remainder = trimmed
		}
		for _, dim := range candidate.Dimensions {
			suffix, ok := strings.CutPrefix(remainder, dim.Name+"-")
			if !ok {
				continue
			}
			if !dim.Has(suffix) {
				return "", "", "", 0, fmt.Errorf(
					"recipe class %q: slot %q dimension %q does not declare value %q",
					class, candidate.Name, dim.Name, suffix)
			}
			return candidate.Name, dim.Name, suffix, ClassValue, nil
		}
	}
	return "", "", "", 0, fmt.Errorf("recipe class %q names no declared slot or dimension of %q", class, s.Component)
}

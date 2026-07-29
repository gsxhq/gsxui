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
	return s.validateClassEncoding()
}

// validateClassEncoding requires the shape's class encoding to be injective.
// Slot, dimension and value are joined with "-" into one flat class name, so
// distinct axes can land on the same string:
//
//	slot "menu"        dimension "button-size" value "lg"
//	slot "menu-button" dimension "size"        value "lg"
//
// both encode as gsxui-recipe-<component>-menu-button-size-lg. Only one of them
// is reachable; the other's rules are dead, and DecodeClass reports the survivor,
// so Conform blames the wrong axis for a rule that is visibly present in the CSS.
// This also catches a base class colliding with a value class — a slot literally
// named "menu-button-size-lg".
func (s Shape) validateClassEncoding() error {
	owner := make(map[string]string)
	claim := func(class, by string) error {
		if previous, taken := owner[class]; taken {
			return fmt.Errorf(
				"%s: %s and %s both encode as class %q — one of them can never be selected; "+
					"rename a slot, dimension or value so the encoding stays unambiguous",
				s.Component, previous, by, class)
		}
		owner[class] = by
		return nil
	}
	for _, slot := range s.Slots {
		if err := claim(s.BaseClass(slot.Name), fmt.Sprintf("slot %q", slot.Name)); err != nil {
			return err
		}
		for _, dimension := range slot.Dimensions {
			for _, value := range dimension.Values {
				by := fmt.Sprintf("slot %q dimension %q value %q", slot.Name, dimension.Name, value)
				if err := claim(s.ValueClass(slot.Name, dimension.Name, value), by); err != nil {
					return err
				}
			}
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

	// Longest-first is a search, not a commitment. A candidate whose dimension
	// prefix matches but whose value is not declared is simply not the decode —
	// "menu-button-size-xl" prefix-matches slot "menu-button" dimension "size"
	// with an undeclared "xl", while slot "menu" dimension "button-size" decodes
	// it exactly. Returning at the first near-miss would reject that class
	// outright. Validate now rejects shapes where two candidates BOTH resolve,
	// so at most one ever survives this loop; the search is what makes the
	// remaining single answer the one that is found.
	var reasons []string
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
				reasons = append(reasons, fmt.Sprintf(
					"slot %q dimension %q does not declare value %q",
					candidate.Name, dim.Name, suffix))
				continue
			}
			return candidate.Name, dim.Name, suffix, ClassValue, nil
		}
	}
	if len(reasons) > 0 {
		return "", "", "", 0, fmt.Errorf("recipe class %q: %s", class, strings.Join(reasons, "; "))
	}
	return "", "", "", 0, fmt.Errorf("recipe class %q names no declared slot or dimension of %q", class, s.Component)
}

// Package recipe models component style recipes: the shape a component
// declares, and the utilities a style supplies for it.
package recipe

import (
	"fmt"
	"strings"
)

// Shape is a component's style interface: the dimensions it varies over and
// the values each dimension admits. Every style must implement all of it.
type Shape struct {
	Component  string
	Base       bool // component has a role() base rule
	Dimensions []Dimension
}

// Dimension is one axis of a Shape. Default must be a member of Values and is
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
	if len(s.Dimensions) == 0 {
		return fmt.Errorf("%s: no dimensions declared", s.Component)
	}
	seen := make(map[string]struct{}, len(s.Dimensions))
	for i, dimension := range s.Dimensions {
		if dimension.Name == "" {
			return fmt.Errorf("%s: dimension %d has no name", s.Component, i)
		}
		if _, exists := seen[dimension.Name]; exists {
			return fmt.Errorf("%s: duplicate dimension %q", s.Component, dimension.Name)
		}
		seen[dimension.Name] = struct{}{}
		if err := dimension.validate(s.Component); err != nil {
			return err
		}
	}
	return nil
}

func (d Dimension) validate(component string) error {
	if len(d.Values) == 0 {
		return fmt.Errorf("%s: dimension %q declares no values", component, d.Name)
	}
	seen := make(map[string]struct{}, len(d.Values))
	for _, value := range d.Values {
		if value == "" {
			return fmt.Errorf("%s: dimension %q has an empty value", component, d.Name)
		}
		if _, exists := seen[value]; exists {
			return fmt.Errorf("%s: dimension %q duplicates value %q", component, d.Name, value)
		}
		seen[value] = struct{}{}
	}
	if !d.Has(d.Default) {
		return fmt.Errorf("%s: dimension %q default %q is not one of its values",
			component, d.Name, d.Default)
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

// Dimension returns the named dimension.
func (s Shape) Dimension(name string) (Dimension, bool) {
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

// ClassKind distinguishes a component's base rule from a dimension value rule.
type ClassKind uint8

const (
	ClassBase ClassKind = iota
	ClassValue
)

func (s Shape) BaseClass() string { return Prefix + s.Component }

func (s Shape) ValueClass(dimension, value string) string {
	return Prefix + s.Component + "-" + dimension + "-" + value
}

// DecodeClass resolves a recipe class name against the shape. It matches
// declared dimensions and values rather than splitting on dashes, so a dashed
// value such as "icon-lg" is unambiguous.
func (s Shape) DecodeClass(class string) (dimension, value string, kind ClassKind, err error) {
	if !strings.HasPrefix(class, Prefix) {
		return "", "", 0, fmt.Errorf("%q is not a recipe class", class)
	}
	if class == s.BaseClass() {
		return "", "", ClassBase, nil
	}
	rest, ok := strings.CutPrefix(class, s.BaseClass()+"-")
	if !ok {
		return "", "", 0, fmt.Errorf("recipe class %q does not belong to component %q", class, s.Component)
	}
	for _, candidate := range s.Dimensions {
		suffix, ok := strings.CutPrefix(rest, candidate.Name+"-")
		if !ok {
			continue
		}
		if !candidate.Has(suffix) {
			return "", "", 0, fmt.Errorf("recipe class %q: dimension %q does not declare value %q",
				class, candidate.Name, suffix)
		}
		return candidate.Name, suffix, ClassValue, nil
	}
	return "", "", 0, fmt.Errorf("recipe class %q names no declared dimension of %q", class, s.Component)
}

// Package recipe models component style recipes: the shape a component
// declares, and the utilities a style supplies for it.
package recipe

import "fmt"

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

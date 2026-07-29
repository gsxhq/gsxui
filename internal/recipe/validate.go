package recipe

import (
	"fmt"
	"slices"
)

// Resolved is a style proven to be a complete implementation of a shape.
type Resolved struct {
	Shape  Shape
	Base   []string
	Values map[string]map[string][]string // [dimension][value]
}

// Utilities returns the utilities a resolved recipe supplies for one value.
func (r Resolved) Utilities(dimension, value string) []string {
	return slices.Clone(r.Values[dimension][value])
}

// Conform checks that style implements shape exactly: every declared
// (dimension, value) has a rule, and every rule maps to a declared one.
func Conform(filename string, shape Shape, style Style) (Resolved, error) {
	if err := shape.Validate(); err != nil {
		return Resolved{}, err
	}

	resolved := Resolved{Shape: shape, Values: make(map[string]map[string][]string, len(shape.Dimensions))}
	for _, dimension := range shape.Dimensions {
		resolved.Values[dimension.Name] = make(map[string][]string, len(dimension.Values))
	}

	// Style to shape: every rule must be declared.
	for _, rule := range style.Rules() {
		dimension, value, kind, err := shape.DecodeClass(rule.Class)
		if err != nil {
			return Resolved{}, fmt.Errorf("%s: %w", filename, err)
		}
		if kind == ClassBase {
			if !shape.Base {
				return Resolved{}, fmt.Errorf("%s: component %q declares no base rule, found %s",
					filename, shape.Component, rule.Class)
			}
			resolved.Base = slices.Clone(rule.Utilities)
			continue
		}
		resolved.Values[dimension][value] = slices.Clone(rule.Utilities)
	}

	// Shape to style: every declaration must be supplied.
	if shape.Base && resolved.Base == nil {
		return Resolved{}, fmt.Errorf("%s: missing base rule .%s", filename, shape.BaseClass())
	}
	for _, dimension := range shape.Dimensions {
		for _, value := range dimension.Values {
			if _, ok := resolved.Values[dimension.Name][value]; !ok {
				return Resolved{}, fmt.Errorf("%s: dimension %q missing value %q",
					filename, dimension.Name, value)
			}
		}
	}
	return resolved, nil
}

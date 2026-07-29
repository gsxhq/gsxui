package recipe

import "fmt"

// Component binds a Shape to the helper calls a canonical component authors.
// stylegen replaces every method call with concrete style source; the method
// bodies exist so the canonical type-checks and so its style-independent
// behavior tests can run.
type Component struct {
	Shape Shape
}

// Role is the component's base recipe class.
func (c Component) Role() string { return c.Shape.BaseClass() }

// Variant and Size are dimension helpers. An empty or unrecognized value
// resolves to the dimension's declared default, matching the generated
// switch's default arm — so a behavior test written against the canonical
// asserts something true of every style.
func (c Component) Variant(value string) string { return c.class("variant", value) }

func (c Component) Size(value string) string { return c.class("size", value) }

func (c Component) class(dimension, value string) string {
	declared, ok := c.Shape.Dimension(dimension)
	if !ok {
		panic(fmt.Sprintf("recipe: component %q declares no dimension %q", c.Shape.Component, dimension))
	}
	if !declared.Has(value) {
		value = declared.Default
	}
	return c.Shape.ValueClass(dimension, value)
}

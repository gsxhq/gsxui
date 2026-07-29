package recipe

import "fmt"

// Component binds a Shape to the helper calls a canonical component authors.
// stylegen replaces every method call with concrete style source; the method
// bodies exist so the canonical type-checks and so its style-independent
// behavior tests can run.
type Component struct {
	Shape Shape
}

// SlotClass is a slot's base recipe class. "" is the component's root element.
func (c Component) SlotClass(slot string) string {
	if _, ok := c.Shape.Slot(slot); !ok {
		panic(fmt.Sprintf("recipe: component %q declares no slot %q", c.Shape.Component, slot))
	}
	return c.Shape.BaseClass(slot)
}

// SlotValueClass resolves a dimension value on a slot. An empty or
// unrecognized value resolves to the dimension's declared default, matching the
// generated switch's default arm — so a behavior test written against the
// canonical asserts something true of every style.
func (c Component) SlotValueClass(slot, dimension, value string) string {
	declared, ok := c.Shape.Slot(slot)
	if !ok {
		panic(fmt.Sprintf("recipe: component %q declares no slot %q", c.Shape.Component, slot))
	}
	dim, ok := declared.Dimension(dimension)
	if !ok {
		panic(fmt.Sprintf("recipe: component %q slot %q declares no dimension %q",
			c.Shape.Component, slot, dimension))
	}
	if !dim.Has(value) {
		value = dim.Default
	}
	return c.Shape.ValueClass(slot, dimension, value)
}

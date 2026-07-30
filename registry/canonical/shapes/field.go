package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Field is the form-layout composite: a fieldset/legend/group frame around a
// row that arranges a label, a control, a description and an error. Almost all
// of its presentation is relational, so read
// registry/styles/nova/field.css and assets/css/styles/default/field.css
// together — the split between them is the interesting part, not this list.
//
// Only two of the contract's axes become dimensions:
//
//   - field's data-orientation (vertical|horizontal|responsive). A real
//     dimension: field.gsx stamps it from a Go param and each value has its own
//     rule set.
//   - field-legend's data-variant (legend|label). Same.
//
// The rest are runtime state nobody varies a recipe dimension on, exactly as
// the agreement check's doc comment anticipates: field's data-invalid and
// data-disabled arrive through attrs and only ever take the value "true" (a
// dimension needs a Default that is one of its values, which would make the
// invalid styling unconditional), and field-group's data-variant and
// field-separator-wrapper's data-content likewise. Each is carried as a
// runtime-state variant instead — see the recipe stylesheet.
var Field = recipe.Shape{
	Component: "field",
	Slots: []recipe.Slot{
		{Name: "set", Base: true},
		{Name: "legend", Base: true, Dimensions: []recipe.Dimension{
			{Name: "variant", Default: "legend", Values: []string{"legend", "label"}},
		}},
		{Name: "group", Base: true},
		{Name: "", Base: true, Dimensions: []recipe.Dimension{
			{Name: "orientation", Default: "vertical", Values: []string{"vertical", "horizontal", "responsive"}},
		}},
		{Name: "content", Base: true},
		{Name: "label", Base: true},
		{Name: "title", Base: true},
		{Name: "description", Base: true},
		{Name: "separator-wrapper", Base: true},
		{Name: "separator", Base: true},
		{Name: "separator-content", Base: true},
		{Name: "error", Base: true},
	},
}

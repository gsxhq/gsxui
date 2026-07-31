package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// InputGroup is the second hyphenated component on the axis (AspectRatio was
// the first). Registry and CSS identity stay kebab — "input-group",
// .gsxui-recipe-input-group — and internal/stylegen/identifier.go derives the
// camel Go identifiers (inputGroupRecipe, receiver inputGroup) from it. That
// mechanism's own test was written around a synthetic input-group shape long
// before this file existed; nothing needed adding to make it work.
//
// The contract's input-group-button slot is deliberately NOT declared. Its
// rules do not live in assets/css/styles/default/input-group.css at all —
// they are already promoted, unwrapped, into assets/css/styles/default.css's
// @layer utilities block, where they exist specifically to outrank Button's
// own migrated size ramp on the same element. Two reasons to leave them:
// migrating a component moves ITS file, not another component's fallout; and
// a "size" dimension here would resolve an unrecognised value (any Button size
// outside xs/sm/icon-xs/icon-sm, which InputGroupButton happily forwards) to
// the dimension's default and paint an xs button, where today no rule matches
// at all.
//
// input-group-control's aria-invalid axis is runtime state, not a dimension,
// the same treatment Input's and Textarea's own shapes give it.
var InputGroup = recipe.Shape{
	Component: "input-group",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "addon", Base: true, Dimensions: []recipe.Dimension{
			{Name: "align", Default: "inline-start", Values: []string{"inline-start", "inline-end", "block-start", "block-end"}},
		}},
		{Name: "text", Base: true},
		{Name: "control", Base: true},
	},
}

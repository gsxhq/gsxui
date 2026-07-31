package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Textarea is a single-slot form control, no dimensions. disabled and
// aria-invalid are runtime state (native pseudo-class / aria attribute), not
// shape dimensions — see internal/stylecontract's Textarea contract, which
// models them as Axes without Values (disabled) or with Values (aria-invalid),
// exactly like Checkbox's already-migrated aria-invalid handling.
var Textarea = recipe.Shape{
	Component: "textarea",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
	},
}

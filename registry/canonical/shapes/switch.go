package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Switch is a single-element component (a real
// <input type="checkbox" role="switch">): one root slot, no dimensions.
// The registry/CSS name is "switch" — a Go keyword. That does not mean the
// component gets renamed: componentIdentifier derives "switch_" (a trailing
// underscore, Go's own convention for exactly this collision), and every
// generated recipe class stays gsxui-recipe-switch, every marker stays
// data-gsxui-slot-switch. See internal/stylegen/identifier.go.
//
// Every relational rule in default.css is either a pseudo-class state
// (:focus-visible, :disabled, :checked), the .dark class, or the thumb's
// own ::before pseudo-element (and its checked/.dark combinations) — none
// chosen by the caller at render time, so all of it folds onto the one
// recipe rule as arbitrary/native Tailwind variants (`checked:`, `dark:`,
// `before:`, and stacked combinations like `dark:checked:before:`), the
// same mechanism Slider's pseudo-element variants and Kbd's `.dark`-baked
// arbitrary variant already proved — see registry/styles/nova/switch.css.
var Switch = recipe.Shape{
	Component: "switch",
	Slots:     []recipe.Slot{{Name: "", Base: true}},
}

package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Tabs is a plain compound component with no Radix client tree (each part
// stamps its own data-state/aria-selected explicitly — see ui/tabs.gsx's own
// doc comment). None of default.css's relational rules are author-facing
// dimensions: data-state="active" is JS/server-toggled RUNTIME state, not a
// caller-chosen variant, so it folds onto the trigger slot's one rule as the
// arbitrary attribute variant `data-[state=active]:`, the same treatment
// Alert's destructive-variant-painting-a-child rule and Card's :has() rules
// already proved — nothing here is a recipe.Dimension. has(>svg), the direct
// child >svg rule, :hover/:focus-visible/:disabled, and .dark (plain and
// combined with hover/data-state) all fold in the same way — see
// registry/styles/nova/tabs.css.
var Tabs = recipe.Shape{
	Component: "tabs",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "list", Base: true},
		{Name: "trigger", Base: true},
		{Name: "content", Base: true},
	},
}

package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Sidebar is the largest shape in the catalogue: one semantic tree rendered
// down two responsive branches (a Sheet-based mobile dialog and a fixed
// desktop column), plus a full menu family.
//
// NO DIMENSIONS, deliberately. The style contract declares data-variant /
// data-size / data-side / data-collapsible / data-state axes on several
// slots, but assets/css/styles/default/sidebar.css never gives every value of
// any of them a rule of its own: sidebar-menu-button's `default` variant and
// sidebar-menu-sub-button's `md` size carry no declarations at all (their
// presentation is the slot's base rule), and data-side/-collapsible/-state
// only ever appear as CONTEXT on an ancestor, never as a self-styling axis
// with one rule per value. recipe.Conform requires a rule for every declared
// value, so modelling those as dimensions would mean inventing empty arms.
// They are authored as runtime-state variants on the slot's own base rule
// instead — the same call the design doc records for aria-* axes, applied
// here for the second documented reason (a value with no rule of its own).
// The shape/contract agreement check is one-directional, so declaring fewer
// axes than the contract is legal and checked.
//
// Six markers sidebar.gsx stamps get no slot, because nothing styles them:
// sidebar-mobile-root (display mechanics only, assets/css/foundation.css),
// sidebar-mobile-title / -description / -inner (fully styled by the composed
// SheetTitle/SheetDescription, or by foundation), and
// sidebar-menu-button-tooltip / -tooltip-content (styled by the composed
// Tooltip, and gated by foundation). This mirrors Sheet's own shape, which
// omits its trigger and close for the same reason. Declaring them would need
// invented base rules; omitting them is safe because DecodeClass only ever
// sees classes this shape emits, and an undeclared compound would be rejected
// rather than silently decoded — see Shape.DecodeClass, whose longest-first
// search exists for exactly the menu / menu-button / menu-sub-button family
// declared below.
var Sidebar = recipe.Shape{
	Component: "sidebar",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "wrapper", Base: true},
		{Name: "mobile-content", Base: true},
		{Name: "mobile-header", Base: true},
		{Name: "desktop", Base: true},
		{Name: "gap", Base: true},
		{Name: "container", Base: true},
		{Name: "inner", Base: true},
		{Name: "trigger", Base: true},
		{Name: "trigger-label", Base: true},
		{Name: "rail", Base: true},
		{Name: "inset", Base: true},
		{Name: "input", Base: true},
		{Name: "header", Base: true},
		{Name: "footer", Base: true},
		{Name: "separator", Base: true},
		{Name: "content", Base: true},
		{Name: "group", Base: true},
		{Name: "group-label", Base: true},
		{Name: "group-action", Base: true},
		{Name: "group-content", Base: true},
		{Name: "menu", Base: true},
		{Name: "menu-item", Base: true},
		{Name: "menu-button", Base: true},
		{Name: "menu-action", Base: true},
		{Name: "menu-badge", Base: true},
		{Name: "menu-skeleton", Base: true},
		{Name: "menu-skeleton-icon", Base: true},
		{Name: "menu-skeleton-text", Base: true},
		{Name: "menu-sub", Base: true},
		{Name: "menu-sub-item", Base: true},
		{Name: "menu-sub-button", Base: true},
	},
}

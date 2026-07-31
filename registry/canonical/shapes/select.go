package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Select is the shadcn/ui custom listbox (native-popover machinery), the
// biggest inventory of this wave. "select" is a Go keyword: the registry/CSS
// identity stays "select" (gsxui-recipe-select, data-gsxui-slot-select), and
// componentIdentifier derives the Go receiver as "select_" — see
// internal/stylegen/identifier.go.
//
// Slots with no rule in default.css at all (select-bridge, select-group)
// carry no accessor — nothing to migrate. Two ancestor strays exist for
// "trigger" (assets/css/styles/default/button-group.css's
// `:where([data-gsxui-slot-button-group]) > :where(select-trigger):not([class*="w-"])`
// and assets/css/styles/default.css's own
// `:where([data-gsxui-slot-button-group]):has(> select[aria-hidden]:last-child) > :where(select-trigger):last-of-type`)
// — both styled SUBJECTs are select-trigger, but both rules react to
// ButtonGroup ancestor context, not to anything Select itself declares; they
// stay exactly where they are, the same "Calendar stray" shape Card's
// migration established.
//
// data-size is the one author-facing Dimension (default|sm), matching both
// default.css and internal/stylecontract's Select axis exactly. Every other
// data-*/aria-* attribute on trigger/content/item (data-state, data-side,
// data-placeholder, aria-invalid, data-disabled, aria-selected) is
// JS/server-toggled RUNTIME state, not a caller-chosen variant — all fold
// onto their slot's one rule as arbitrary/native Tailwind variants, matching
// the treatment Tabs' data-state="active" already proved. The
// SelectValue-inside-SelectTrigger rule and the checked-item-paints-its-own-
// indicator rule are the "one slot reacts to ANOTHER slot" shapes Alert's
// destructive-variant-paints-description and Card's :has() already proved,
// just as ancestor-context and same-element-conditional-child forms
// respectively. See registry/styles/nova/select.css.
var Select = recipe.Shape{
	Component: "select",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "trigger", Base: true, Dimensions: []recipe.Dimension{
			{Name: "size", Default: "default", Values: []string{"default", "sm"}},
		}},
		{Name: "value", Base: true},
		{Name: "content", Base: true},
		{Name: "label", Base: true},
		{Name: "item", Base: true},
		{Name: "item-text", Base: true},
		{Name: "item-indicator", Base: true},
		{Name: "separator", Base: true},
	},
}

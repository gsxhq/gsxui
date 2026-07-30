package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Combobox is the shadcn/ui Base UI filtered listbox (see
// registry/canonical/combobox.gsx's own doc comment for the port's
// MECHANISM/ADAPT/GAP record). No dimensions: neither
// assets/css/styles/default/combobox.css nor internal/stylecontract declares
// a variant/size axis for any Combobox part.
//
// Slots with no rule of their own in default.css carry no accessor and are
// absent here: combobox-input (the control's presentation is entirely
// InputGroupInput's), combobox-bridge (assets/css/foundation.css's sr-only
// group, not this style's) and combobox-group.
//
// Three of the slots sit on an element that ALSO carries another
// component's marker: combobox-input-group is InputGroup's own element,
// and combobox-trigger / combobox-clear are InputGroupButton's, i.e.
// Button's — which is already migrated. Their recipe classes are added
// through the composing call's own class attribute
// (registry/canonical/combobox.gsx),
// merged by merge.Merge alongside the composed component's, the same
// FieldLabel-composes-Label shape Wave 1 established. combobox-input-group's
// only utility is w-auto — upstream's own cn("w-auto", className) default,
// which Wave 4c had to compose into ComboboxInput's wrapper attrs by hand
// because CheckAuthoring forbids a class= attribute in an unmigrated
// component; it is an ordinary recipe utility again now.
//
// ONE rule from the old file is deliberately NOT here — the "Calendar
// stray" shape Card's migration named:
// `[data-gsxui-slot-combobox-content] > [data-gsxui-slot-input-group]`.
// Its styled SUBJECT is InputGroup, not a Combobox slot; its uncontested
// half (m-1 mb-0 shadow-none) moved verbatim to
// assets/css/styles/default/input-group.css rather than onto this recipe,
// because attaching it here would paint a foreign component from Combobox's
// own classes. Its contested half (h-8 border-input/30 bg-input/30) was
// already promoted to assets/css/styles/default.css's @layer utilities block
// by InputGroup's own migration and stays there.
//
// combobox-content's discrete-transition block came from the file this
// component shared with the menu family
// (assets/css/styles/default/menu-popover-transition.css) — Combobox was
// the last unmigrated selector in it, so that file is now deleted.
//
// RETAINED, deliberately not here (spec 10b): combobox-item's
// [data-disabled] opacity-50. data-disabled is caller-supplied through
// attrs (ComboboxItem declares no disabled param at all), so it stays a
// plain @layer components rule in assets/css/styles/default/combobox.css
// that a caller's own opacity-* can outrank. The pointer-events-none half
// of that rule IS migrated, matching Command and the menu family.
var Combobox = recipe.Shape{
	Component: "combobox",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "input-group", Base: true},
		{Name: "trigger", Base: true},
		{Name: "trigger-icon", Base: true},
		{Name: "clear", Base: true},
		{Name: "content", Base: true},
		{Name: "list", Base: true},
		{Name: "item", Base: true},
		{Name: "item-indicator", Base: true},
		{Name: "label", Base: true},
		{Name: "empty", Base: true},
		{Name: "separator", Base: true},
	},
}

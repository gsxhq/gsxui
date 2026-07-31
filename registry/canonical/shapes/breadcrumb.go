package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Breadcrumb declares only the slots default.css actually styles. The root
// <nav> is unstyled (same carve-out as Collapsible's root/content), so it is
// omitted. breadcrumb-ellipsis-label is fully styled by a shared sr-only rule
// in assets/css/foundation.css's @layer base (one selector list shared with
// Pagination, Dialog, Sheet, Command, Combobox, Carousel, Calendar's own
// "-label"/hidden-caption slots) — that rule is foundation MECHANICS
// legitimately shared across many unrelated components, not Breadcrumb's own
// rule to migrate, so this shape omits that slot too and no class is ever
// added to it.
//
// Separator and Ellipsis have no property of their own — default.css only
// reaches their child <svg> (`> svg { size-X }`), an unconditional descendant
// selector, not a conditional :has(). That translates as a same-slot
// `[&>svg]:size-X` arbitrary variant on the slot's own base rule (Ellipsis
// combines it with its own flex/size utilities in one rule); it is not the
// cross-slot "root paints a different slot" shape Alert/Resizable used, so it
// lives entirely on the owning slot itself.
var Breadcrumb = recipe.Shape{
	Component: "breadcrumb",
	Slots: []recipe.Slot{
		{Name: "list", Base: true},
		{Name: "item", Base: true},
		{Name: "link", Base: true},
		{Name: "page", Base: true},
		{Name: "separator", Base: true},
		{Name: "ellipsis", Base: true},
	},
}

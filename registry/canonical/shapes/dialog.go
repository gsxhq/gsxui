package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Dialog is a compound component. data-state (closed/open) on dialog-content
// is runtime state (native attribute variant), not a shape dimension — same
// treatment as aria-invalid elsewhere.
//
// trigger, close, close-label and footer-close are omitted (same carve-out as
// Breadcrumb's unstyled root/ellipsis-label, see that shape's own comment):
// DialogTrigger/DialogClose are plain unstyled buttons in default.css, close-
// label is fully styled by the shared sr-only rule in
// assets/css/foundation.css's @layer base (shared with Pagination, Sheet,
// Command, Combobox, Carousel, Calendar), and footer-close is a marker on an
// already-migrated <ui.Button> (DialogFooter composes it directly) which
// brings its own recipe utilities — none of the four have a rule of their own
// to migrate, so no accessor call is added and no class is ever added to
// them.
var Dialog = recipe.Shape{
	Component: "dialog",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "content", Base: true},
		{Name: "close-button", Base: true},
		{Name: "close-icon", Base: true},
		{Name: "header", Base: true},
		{Name: "footer", Base: true},
		{Name: "title", Base: true},
		{Name: "description", Base: true},
	},
}

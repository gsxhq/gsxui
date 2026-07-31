package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// Command is the shadcn/ui Command palette (cmdk's markup re-expressed as
// public slots — see registry/canonical/command.gsx's own doc comment). No
// dimensions: neither assets/css/styles/default/command.css nor
// internal/stylecontract declares a variant/size axis for any part.
//
// Slots with no rule of their own in default.css carry no accessor and are
// absent here: command-dialog (a bare <Dialog> wrapper),
// command-dialog-header (styled by assets/css/foundation.css's sr-only
// group, not by this style) and command-dialog-command.
//
// command-dialog-content IS a slot, but only for the one `overflow-hidden`
// declaration that survived Dialog's own migration. Its sibling p-0 /
// sm:max-w-lg pair stays promoted in assets/css/styles/default.css's
// @layer utilities block: CommandDialog stamps this marker onto the SAME
// <dialog> element that carries Dialog's own recipe class, and those two
// properties contest Dialog's p-4/sm:max-w-sm. overflow-hidden contests
// nothing Dialog sets, so it migrates. The `:where(dialog)` element
// qualifier is preserved as a `[dialog&]:` arbitrary variant.
//
// The command-dialog-content DESCENDANT rules (the palette's larger
// in-dialog density: h-12 rows, size-5 icons, px-2 groups, px-2/py-3 items,
// the group:not([hidden])~group pt-0 pair) all have a Command slot as their
// styled subject, so they fold onto that slot's own recipe as arbitrary
// ancestor variants — the same shape Select's SelectValue-inside-
// SelectTrigger rule already proved.
//
// RETAINED, deliberately not here (spec 10b): command-item's
// [data-disabled="true"] opacity-50. data-disabled is stamped by the CALLER
// through attrs (see CommandItem's doc comment), so it must stay a plain
// @layer components rule a caller's own opacity-* utility can outrank —
// the identical ruling assets/css/styles/default/menu.css already made.
// The pointer-events-none half of that same rule IS migrated, for the same
// reason menu.css gives: a caller must not be able to re-enable pointer
// interaction on a row they marked disabled.
var Command = recipe.Shape{
	Component: "command",
	Slots: []recipe.Slot{
		{Name: "", Base: true},
		{Name: "dialog-content", Base: true},
		{Name: "input-wrapper", Base: true},
		{Name: "input", Base: true},
		{Name: "list", Base: true},
		{Name: "empty", Base: true},
		{Name: "group", Base: true},
		{Name: "group-heading", Base: true},
		{Name: "separator", Base: true},
		{Name: "item", Base: true},
		{Name: "shortcut", Base: true},
	},
}

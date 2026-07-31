package shapes

import "github.com/gsxhq/gsxui/internal/recipe"

// AlertDialog is a compound component layered on top of Dialog: its root,
// trigger, action and cancel all delegate their presentation elsewhere
// (<ui.Dialog>, a plain unstyled button, and two already-migrated
// <ui.Button>s respectively), so none of them has a rule of its own to
// migrate and none gets an accessor — the same carve-out Dialog's own shape
// documents for its trigger/close/close-label/footer-close.
//
// content is the interesting slot. AlertDialogContent composes
// <ui.DialogContent> and its whole presentation is a NARROWING of Dialog's
// box (max-w-xs, sm:max-w-sm). Before this migration that narrowing was a
// same-element compound selector retained as plain CSS in
// assets/css/styles/default/alert-dialog.css; :where()-wrapped at (0,0,0) it
// could not beat Dialog's own (0,1,0) recipe class in the same layer, so it
// was dead. Carrying it as this slot's own utilities instead hands it to
// merge.Merge, which replaces Dialog's max-w-* outright — exactly what
// upstream does (shadcn alert-dialog.tsx puts AlertDialogContent's max-w on
// its own className and lets twMerge settle it against DialogContent's).
//
// data-state (closed/open) and the native open attribute are runtime state,
// not shape dimensions — same treatment as Dialog's own content slot.
var AlertDialog = recipe.Shape{
	Component: "alert-dialog",
	Slots: []recipe.Slot{
		{Name: "content", Base: true},
		{Name: "header", Base: true},
		{Name: "footer", Base: true},
		{Name: "title", Base: true},
		{Name: "description", Base: true},
	},
}

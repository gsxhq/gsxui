package ui

import "github.com/gsxhq/gsx"

// AlertDialog and its parts are the shadcn/ui AlertDialog
// (registry/new-york-v4/ui/alert-dialog.tsx), ported onto the exact same
// native <dialog> machinery Dialog already provides (ui/dialog.gsx,
// ui/dialog.js) rather than a second Radix-shaped implementation: an
// AlertDialog IS a Dialog that (a) cannot be light-dismissed by an outside
// click and (b) never renders the injected close X — everything else
// (top-layer stacking, Esc-to-close, trigger/content wiring by proximity,
// data-state, aria wiring) is identical, so this file composes ui.Dialog/
// ui.DialogContent instead of re-deriving them. See docs/jsx-parity.md
// `## alert-dialog` for the ledgered NOTE on which upstream revision these
// class strings were verified against (shadcn has since layered a `size`
// variant and an `AlertDialogMedia` part onto this component; both are out
// of scope here — see the ledger).
//
// data-gsxui-dialog-static is the one feature ui/dialog.js gained for this
// port: a content element carrying it is skipped by the backdrop-click
// light-dismiss path (Esc/cancel is untouched, exactly reproducing Radix's
// own AlertDialog behavior — outside clicks are ignored, Esc still closes).
// AlertDialogContent stamps it automatically; any ui.DialogContent can opt
// into the same behavior by adding the attribute directly.

// AlertDialog composes ui.Dialog directly — no state or markup of its own
// beyond overriding the inherited data-slot to match shadcn's
// "alert-dialog" (Dialog's own root renders data-slot="dialog"; the
// attrs-position override here is the same explicit-non-parameter-attribute
// mechanism ItemSeparator/FieldLabel use to override a composed component's
// own data-slot — see docs/jsx-parity.md `## item`/`## field`).
component AlertDialog(children gsx.Node, attrs gsx.Attrs) {
	<Dialog data-slot="alert-dialog" { attrs... }>{ children }</Dialog>
}

// AlertDialogTrigger renders its own <button>, same reasoning and the same
// button-in-button HTML trap as DialogTrigger (see ui/dialog.gsx's
// DialogTrigger doc comment and docs/jsx-parity.md `## dialog` FINDING):
// its children must be phrasing content, never a Button or other
// interactive element. For a styled trigger, skip the wrapper and put the
// data attribute on the Button itself: <ui.Button
// data-gsxui-dialog-trigger>Delete</ui.Button> — the same documented idiom,
// unaffected by which dialog flavor (Dialog or AlertDialog) sits behind it,
// since both key off the same data-gsxui-dialog-trigger contract.
component AlertDialogTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		data-slot="alert-dialog-trigger"
		data-gsxui-dialog-trigger
		type="button"
		aria-haspopup="dialog"
		aria-expanded="false"
		{ attrs... }
	>
		{ children }
	</button>
}

// AlertDialogContent composes ui.DialogContent with hideCloseButton always
// true (shadcn's AlertDialog never renders the injected X — there is no
// showCloseButton-equivalent prop on AlertDialogContent at all) plus three
// attrs layered on via the same override mechanism as AlertDialog's own
// data-slot: data-slot="alert-dialog-content" (shadcn's own slot name),
// role="alertdialog" (the one a11y difference from a plain Dialog — it
// tells assistive tech this dialog demands a response), and
// data-gsxui-dialog-static (opts this content out of dialog.js's
// backdrop-click light dismiss — see this file's header comment).
//
// No class attr is passed here (there is no third, merged class string).
// Diffed token-for-token against DialogContent's own base class, every
// non-variant utility alert-dialog.tsx's content class carries (bg-
// background, fixed, top-[50%], left-[50%], z-50, w-full, max-w-[calc(100%-
// 2rem)], translate-x/y-[-50%], gap-4, rounded-lg, border, p-6, shadow-lg,
// duration-200, sm:max-w-lg) and all six data-[state=…]:animate-in/out,
// fade-*, zoom-* tokens are already there — both dialogs share one
// centered-card recipe upstream (see the ledger NOTE on which revision).
// The one token that is NOT shared, a bare `grid`, is dropped rather than
// re-supplied: DialogContent's own `open:grid` exists specifically so the
// content stays display:none while the native <dialog> is closed (ui/
// dialog.gsx's own ADAPT); passing an unscoped `grid` alongside a
// `open:`-scoped one would not be resolved as a tailwind-merge conflict
// (variant scope is part of the conflict key — the same non-collision
// documented for accordion's rotate override, docs/jsx-parity.md
// `## accordion`) and, worse, the configured merger doesn't recognize the
// tw-animate-css tokens (`animate-in`/`fade-*`/`zoom-*`) as a conflict
// group at all (see `## animations`'s FINDING), so re-supplying any of
// them would literally duplicate them in the output rather than merge.
// Net effect: AlertDialogContent's rendered class is exactly DialogContent's
// own default, unchanged — pinned as such in alert-dialog_test.go, itself
// the parity claim (only role="alertdialog", data-gsxui-dialog-static, and
// hideCloseButton's injected-X/backdrop-dismiss suppression distinguish an
// alert dialog from a plain one).
component AlertDialogContent(children gsx.Node, attrs gsx.Attrs) {
	<DialogContent
		hideCloseButton={true}
		data-slot="alert-dialog-content"
		role="alertdialog"
		data-gsxui-dialog-static
		{ attrs... }
	>
		{ children }
	</DialogContent>
}

// AlertDialogHeader/Footer/Title/Description are NOT composed from Dialog's
// own Header/Footer/Title/Description — alert-dialog.tsx's class strings
// differ from dialog.tsx's (Title drops leading-none; Footer happens to
// coincide byte-for-byte with DialogFooter's own in the upstream revision
// this port targets, see the ledger NOTE) — each renders its own element
// with alert-dialog's own token-for-token classes.
component AlertDialogHeader(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="alert-dialog-header" class="flex flex-col gap-2 text-center sm:text-left" { attrs... }>{ children }</div>
}

component AlertDialogFooter(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="alert-dialog-footer" class="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end" { attrs... }>{ children }</div>
}

component AlertDialogTitle(children gsx.Node, attrs gsx.Attrs) {
	<h2 data-slot="alert-dialog-title" class="text-lg font-semibold" { attrs... }>{ children }</h2>
}

component AlertDialogDescription(children gsx.Node, attrs gsx.Attrs) {
	<p data-slot="alert-dialog-description" class="text-muted-foreground text-sm" { attrs... }>{ children }</p>
}

// AlertDialogAction is ui.Button (default variant/size, shadcn's own
// buttonVariants() default) plus data-gsxui-dialog-close — the same
// data-attribute idiom DialogFooter's own Close button and DialogClose use
// (docs/jsx-parity.md `## dialog` MECHANISM) standing in for shadcn's
// <AlertDialogAction asChild><Button/></AlertDialogAction> wrap.
// data-slot="alert-dialog-action" overrides Button's own "button" slot via
// the same attrs-position mechanism as AlertDialogContent's own override.
component AlertDialogAction(children gsx.Node, attrs gsx.Attrs) {
	<Button data-slot="alert-dialog-action" data-gsxui-dialog-close { attrs... }>{ children }</Button>
}

// AlertDialogCancel is ui.Button with variant="outline" (shadcn's own
// buttonVariants({variant: "outline"})) plus data-gsxui-dialog-close — same
// mechanism as AlertDialogAction.
component AlertDialogCancel(children gsx.Node, attrs gsx.Attrs) {
	<Button variant="outline" data-slot="alert-dialog-cancel" data-gsxui-dialog-close { attrs... }>{ children }</Button>
}

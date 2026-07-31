package canonical

import "github.com/gsxhq/gsx"

// AlertDialog composes Dialog's native top-layer and state machinery while
// opting out of backdrop light dismissal.
//
// AlertDialogContent's own class is a NARROWING of Dialog's box: it travels
// into <DialogContent>'s class attribute and merge.Merge settles the
// competing max-w-* utilities there, the same way upstream lets twMerge
// settle them. Every class attribute below is resolved to concrete utilities
// at generation time.
component AlertDialog(children gsx.Node, attrs gsx.Attrs) {
	<Dialog { attrs... } data-gsxui-slot-alert-dialog>{ children }</Dialog>
}

component AlertDialogTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		data-gsxui-dialog-trigger
		type="button"
		aria-haspopup="dialog"
		aria-expanded="false"
		{ attrs... }
		data-gsxui-slot-alert-dialog-trigger
	>
		{ children }
	</button>
}

component AlertDialogContent(children gsx.Node, attrs gsx.Attrs) {
	<DialogContent
		class={ alertDialog.Content() }
		hideCloseButton={true}
		role="alertdialog"
		data-gsxui-dialog-static
		{ attrs... }
		data-gsxui-slot-alert-dialog-content
	>
		{ children }
	</DialogContent>
}

component AlertDialogHeader(children gsx.Node, attrs gsx.Attrs) {
	<div class={ alertDialog.Header() } { attrs... } data-gsxui-slot-alert-dialog-header>{ children }</div>
}

component AlertDialogFooter(children gsx.Node, attrs gsx.Attrs) {
	<div class={ alertDialog.Footer() } { attrs... } data-gsxui-slot-alert-dialog-footer>{ children }</div>
}

component AlertDialogTitle(children gsx.Node, attrs gsx.Attrs) {
	<h2 class={ alertDialog.Title() } { attrs... } data-gsxui-slot-alert-dialog-title data-gsxui-slot-dialog-title>{ children }</h2>
}

component AlertDialogDescription(children gsx.Node, attrs gsx.Attrs) {
	<p class={ alertDialog.Description() } { attrs... } data-gsxui-slot-alert-dialog-description data-gsxui-slot-dialog-description>{ children }</p>
}

component AlertDialogAction(children gsx.Node, attrs gsx.Attrs) {
	<Button data-gsxui-dialog-close { attrs... } data-gsxui-slot-alert-dialog-action>{ children }</Button>
}

component AlertDialogCancel(children gsx.Node, attrs gsx.Attrs) {
	<Button
		variant="outline"
		data-gsxui-dialog-close
		{ attrs... }
		data-gsxui-slot-alert-dialog-cancel
	>
		{ children }
	</Button>
}

package ui

import "github.com/gsxhq/gsx"

// AlertDialog composes Dialog's native top-layer and state machinery while
// opting out of backdrop light dismissal.
component AlertDialog(children gsx.Node, attrs gsx.Attrs) {
	<Dialog { attrs... } data-gsxui-slot-alert-dialog>{ children }</Dialog>
}

component AlertDialogTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		data-gsxui-dialog-trigger
		type="button"
		aria-haspopup="dialog"
		aria-expanded="false"
		{ attrs... } data-gsxui-slot-alert-dialog-trigger
	>
		{ children }
	</button>
}

component AlertDialogContent(children gsx.Node, attrs gsx.Attrs) {
	<DialogContent
		hideCloseButton={true}
		role="alertdialog"
		data-gsxui-dialog-static
		{ attrs... } data-gsxui-slot-alert-dialog-content
	>
		{ children }
	</DialogContent>
}

component AlertDialogHeader(children gsx.Node, attrs gsx.Attrs) {
	<div { attrs... } data-gsxui-slot-alert-dialog-header>{ children }</div>
}

component AlertDialogFooter(children gsx.Node, attrs gsx.Attrs) {
	<div { attrs... } data-gsxui-slot-alert-dialog-footer>{ children }</div>
}

component AlertDialogTitle(children gsx.Node, attrs gsx.Attrs) {
	<h2 data-gsxui-dialog-title { attrs... } data-gsxui-slot-alert-dialog-title>{ children }</h2>
}

component AlertDialogDescription(children gsx.Node, attrs gsx.Attrs) {
	<p data-gsxui-dialog-description { attrs... } data-gsxui-slot-alert-dialog-description>{ children }</p>
}

component AlertDialogAction(children gsx.Node, attrs gsx.Attrs) {
	<Button data-gsxui-dialog-close { attrs... } data-gsxui-slot-alert-dialog-action>{ children }</Button>
}

component AlertDialogCancel(children gsx.Node, attrs gsx.Attrs) {
	<Button
		variant="outline"
		data-gsxui-dialog-close
		{ attrs... } data-gsxui-slot-alert-dialog-cancel
	>
		{ children }
	</Button>
}

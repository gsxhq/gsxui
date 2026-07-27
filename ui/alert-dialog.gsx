package ui

import "github.com/gsxhq/gsx"

// AlertDialog composes Dialog's native top-layer and state machinery while
// opting out of backdrop light dismissal.
component AlertDialog(children gsx.Node, attrs gsx.Attrs) {
	<Dialog { withSlot("alert-dialog", attrs)... }>{ children }</Dialog>
}

component AlertDialogTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		data-gsxui-dialog-trigger
		type="button"
		aria-haspopup="dialog"
		aria-expanded="false"
		{ withSlot("alert-dialog-trigger", attrs)... }
	>
		{ children }
	</button>
}

component AlertDialogContent(children gsx.Node, attrs gsx.Attrs) {
	<DialogContent
		hideCloseButton={true}
		role="alertdialog"
		data-gsxui-dialog-static
		{ withSlot("alert-dialog-content", attrs)... }
	>
		{ children }
	</DialogContent>
}

component AlertDialogHeader(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("alert-dialog-header", attrs)... }>{ children }</div>
}

component AlertDialogFooter(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("alert-dialog-footer", attrs)... }>{ children }</div>
}

component AlertDialogTitle(children gsx.Node, attrs gsx.Attrs) {
	<h2 data-gsxui-dialog-title { withSlot("alert-dialog-title", attrs)... }>{ children }</h2>
}

component AlertDialogDescription(children gsx.Node, attrs gsx.Attrs) {
	<p data-gsxui-dialog-description { withSlot("alert-dialog-description", attrs)... }>{ children }</p>
}

component AlertDialogAction(children gsx.Node, attrs gsx.Attrs) {
	<Button data-gsxui-dialog-close { withSlot("alert-dialog-action", attrs)... }>{ children }</Button>
}

component AlertDialogCancel(children gsx.Node, attrs gsx.Attrs) {
	<Button
		variant="outline"
		data-gsxui-dialog-close
		{ withSlot("alert-dialog-cancel", attrs)... }
	>
		{ children }
	</Button>
}

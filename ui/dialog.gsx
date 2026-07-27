package ui

import "github.com/gsxhq/gsx"

// Dialog uses the native <dialog> top layer. Trigger/content wiring is scoped
// by the dedicated root hook and implemented by ui/dialog.js.
component Dialog(children gsx.Node, attrs gsx.Attrs) {
	<div data-gsxui-dialog { attrs... } data-gsxui-slot-dialog>{ children }</div>
}

component DialogTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		data-gsxui-dialog-trigger
		type="button"
		aria-haspopup="dialog"
		aria-expanded="false"
		{ attrs... } data-gsxui-slot-dialog-trigger
	>
		{ children }
	</button>
}

component DialogContent(hideCloseButton bool, children gsx.Node, attrs gsx.Attrs) {
	<dialog
		data-gsxui-dialog-content
		data-state="closed"
		{ attrs... } data-gsxui-slot-dialog-content
	>
		{ children }
		{ if !hideCloseButton {
			<button
				type="button"
				data-gsxui-dialog-close
				data-gsxui-slot-dialog-close-button data-gsxui-slot-dialog-close
			>
				<svg
					aria-hidden="true"
					xmlns="http://www.w3.org/2000/svg"
					width="24"
					height="24"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
					data-gsxui-slot-dialog-close-icon
				>
					<path d="M18 6 6 18"/>
					<path d="m6 6 12 12"/>
				</svg>
				<span data-gsxui-slot-dialog-close-label>Close</span>
			</button>
		} }
	</dialog>
}

component DialogHeader(children gsx.Node, attrs gsx.Attrs) {
	<div { attrs... } data-gsxui-slot-dialog-header>{ children }</div>
}

component DialogFooter(showCloseButton bool, children gsx.Node, attrs gsx.Attrs) {
	<div { attrs... } data-gsxui-slot-dialog-footer>
		{ children }
		{ if showCloseButton {
			<Button
				variant="outline"
				data-gsxui-dialog-close
				data-gsxui-slot-dialog-footer-close
			>
				Close
			</Button>
		} }
	</div>
}

component DialogTitle(children gsx.Node, attrs gsx.Attrs) {
	<h2 data-gsxui-dialog-title { attrs... } data-gsxui-slot-dialog-title>{ children }</h2>
}

component DialogDescription(children gsx.Node, attrs gsx.Attrs) {
	<p data-gsxui-dialog-description { attrs... } data-gsxui-slot-dialog-description>{ children }</p>
}

component DialogClose(children gsx.Node, attrs gsx.Attrs) {
	<button data-gsxui-dialog-close type="button" { attrs... } data-gsxui-slot-dialog-close>{ children }</button>
}

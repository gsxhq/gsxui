package ui

import "github.com/gsxhq/gsx"

// Dialog uses the native <dialog> top layer. Trigger/content wiring is scoped
// by the dedicated root hook and implemented by ui/dialog.js.
component Dialog(children gsx.Node, attrs gsx.Attrs) {
	<div data-gsxui-dialog { withSlot("dialog", attrs)... }>{ children }</div>
}

component DialogTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		data-gsxui-dialog-trigger
		type="button"
		aria-haspopup="dialog"
		aria-expanded="false"
		{ withSlot("dialog-trigger", attrs)... }
	>
		{ children }
	</button>
}

component DialogContent(hideCloseButton bool, children gsx.Node, attrs gsx.Attrs) {
	<dialog
		data-gsxui-dialog-content
		data-state="closed"
		{ withSlot("dialog-content", attrs)... }
	>
		{ children }
		{ if !hideCloseButton {
			<button
				type="button"
				data-gsxui-dialog-close
				{ withSlot("dialog-close", withSlot("dialog-close-button", nil))... }
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
					{ withSlot("dialog-close-icon", nil)... }
				>
					<path d="M18 6 6 18"/>
					<path d="m6 6 12 12"/>
				</svg>
				<span { withSlot("dialog-close-label", nil)... }>Close</span>
			</button>
		} }
	</dialog>
}

component DialogHeader(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("dialog-header", attrs)... }>{ children }</div>
}

component DialogFooter(showCloseButton bool, children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("dialog-footer", attrs)... }>
		{ children }
		{ if showCloseButton {
			<Button
				variant="outline"
				data-gsxui-dialog-close
				{ withSlot("dialog-footer-close", nil)... }
			>
				Close
			</Button>
		} }
	</div>
}

component DialogTitle(children gsx.Node, attrs gsx.Attrs) {
	<h2 data-gsxui-dialog-title { withSlot("dialog-title", attrs)... }>{ children }</h2>
}

component DialogDescription(children gsx.Node, attrs gsx.Attrs) {
	<p data-gsxui-dialog-description { withSlot("dialog-description", attrs)... }>{ children }</p>
}

component DialogClose(children gsx.Node, attrs gsx.Attrs) {
	<button data-gsxui-dialog-close type="button" { withSlot("dialog-close", attrs)... }>{ children }</button>
}

package ui

import "github.com/gsxhq/gsx"

// Sheet composes Dialog's root and behavior with side-anchored presentation.
component Sheet(children gsx.Node, attrs gsx.Attrs) {
	<Dialog { withSlot("sheet", attrs)... }>{ children }</Dialog>
}

component SheetTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		data-gsxui-dialog-trigger
		type="button"
		aria-haspopup="dialog"
		aria-expanded="false"
		{ withSlot("sheet-trigger", attrs)... }
	>
		{ children }
	</button>
}

component SheetContent(side string, hideCloseButton bool, children gsx.Node, attrs gsx.Attrs) {
	<dialog
		data-gsxui-dialog-content
		data-state="closed"
		data-side={side |> default("right")}
		{ withSlot("dialog-content", withSlot("sheet-content", attrs))... }
	>
		{ children }
		{ if !hideCloseButton {
			<button
				type="button"
				data-gsxui-dialog-close
				{ withSlot("sheet-close", withSlot("sheet-close-button", nil))... }
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
					{ withSlot("sheet-close-icon", nil)... }
				>
					<path d="M18 6 6 18"/>
					<path d="m6 6 12 12"/>
				</svg>
				<span { withSlot("sheet-close-label", nil)... }>Close</span>
			</button>
		} }
	</dialog>
}

component SheetHeader(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("sheet-header", attrs)... }>{ children }</div>
}

component SheetFooter(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("sheet-footer", attrs)... }>{ children }</div>
}

component SheetTitle(children gsx.Node, attrs gsx.Attrs) {
	<h2 data-gsxui-dialog-title { withSlot("sheet-title", attrs)... }>{ children }</h2>
}

component SheetDescription(children gsx.Node, attrs gsx.Attrs) {
	<p data-gsxui-dialog-description { withSlot("sheet-description", attrs)... }>{ children }</p>
}

component SheetClose(children gsx.Node, attrs gsx.Attrs) {
	<button data-gsxui-dialog-close type="button" { withSlot("sheet-close", attrs)... }>{ children }</button>
}

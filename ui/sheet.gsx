package ui

import "github.com/gsxhq/gsx"

// Sheet composes Dialog's root and behavior with side-anchored presentation.
component Sheet(children gsx.Node, attrs gsx.Attrs) {
	<Dialog { attrs... } data-gsxui-slot-sheet>{ children }</Dialog>
}

component SheetTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		data-gsxui-dialog-trigger
		type="button"
		aria-haspopup="dialog"
		aria-expanded="false"
		{ attrs... } data-gsxui-slot-sheet-trigger
	>
		{ children }
	</button>
}

component SheetContent(side string, hideCloseButton bool, children gsx.Node, attrs gsx.Attrs) {
	<dialog
		data-gsxui-dialog-content
		data-state="closed"
		data-side={side |> default("right")}
		{ attrs... } data-gsxui-slot-sheet-content data-gsxui-slot-dialog-content
	>
		{ children }
		{ if !hideCloseButton {
			<button
				type="button"
				data-gsxui-dialog-close
				data-gsxui-slot-sheet-close-button data-gsxui-slot-sheet-close
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
					data-gsxui-slot-sheet-close-icon
				>
					<path d="M18 6 6 18"/>
					<path d="m6 6 12 12"/>
				</svg>
				<span data-gsxui-slot-sheet-close-label>Close</span>
			</button>
		} }
	</dialog>
}

component SheetHeader(children gsx.Node, attrs gsx.Attrs) {
	<div { attrs... } data-gsxui-slot-sheet-header>{ children }</div>
}

component SheetFooter(children gsx.Node, attrs gsx.Attrs) {
	<div { attrs... } data-gsxui-slot-sheet-footer>{ children }</div>
}

component SheetTitle(children gsx.Node, attrs gsx.Attrs) {
	<h2 data-gsxui-dialog-title { attrs... } data-gsxui-slot-sheet-title>{ children }</h2>
}

component SheetDescription(children gsx.Node, attrs gsx.Attrs) {
	<p data-gsxui-dialog-description { attrs... } data-gsxui-slot-sheet-description>{ children }</p>
}

component SheetClose(children gsx.Node, attrs gsx.Attrs) {
	<button data-gsxui-dialog-close type="button" { attrs... } data-gsxui-slot-sheet-close>{ children }</button>
}

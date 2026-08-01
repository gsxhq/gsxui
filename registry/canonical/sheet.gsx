package canonical

import "github.com/gsxhq/gsx"

// Sheet composes Dialog's root and behavior with side-anchored presentation.
//
// side is physical, matching shadcn: side="left"/"right" always anchors to
// that physical viewport edge (position, border and slide direction), in
// both dir="ltr" and dir="rtl" documents. Interior presentation — header and
// footer padding, title/description text alignment, the close button's
// logical end-* offset — is logical and follows dir normally.
//
// SheetContent deliberately does NOT reuse Dialog's content presentation —
// it must not inherit the plain-modal box — so it carries the shared
// <dialog> chrome as its own recipe utilities instead. Every class attribute
// below is resolved to concrete utilities at generation time.
component Sheet(children gsx.Node, attrs gsx.Attrs) {
	<Dialog { attrs... } data-gsxui-slot-sheet>{ children }</Dialog>
}

component SheetTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		type="button"
		aria-haspopup="dialog"
		aria-expanded="false"
		{ attrs... }
		data-gsxui-slot-sheet-trigger
	>
		{ children }
	</button>
}

component SheetContent(side string, hideCloseButton bool, children gsx.Node, attrs gsx.Attrs) {
	<dialog
		class={ sheet.Content(), sheet.ContentSide(side) }
		data-state="closed"
		data-side={side |> default("right")}
		{ attrs... }
		data-gsxui-slot-sheet-content
		data-gsxui-slot-dialog-content
	>
		{ children }
		{ if !hideCloseButton {
			<button
				type="button"
				class={ sheet.CloseButton() }
				data-gsxui-dialog-close
				data-gsxui-slot-sheet-close-button
				data-gsxui-slot-sheet-close
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
					class={ sheet.CloseIcon() }
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
	<div class={ sheet.Header() } { attrs... } data-gsxui-slot-sheet-header>{ children }</div>
}

component SheetFooter(children gsx.Node, attrs gsx.Attrs) {
	<div class={ sheet.Footer() } { attrs... } data-gsxui-slot-sheet-footer>{ children }</div>
}

component SheetTitle(children gsx.Node, attrs gsx.Attrs) {
	<h2 class={ sheet.Title() } { attrs... } data-gsxui-slot-sheet-title data-gsxui-slot-dialog-title>{ children }</h2>
}

component SheetDescription(children gsx.Node, attrs gsx.Attrs) {
	<p class={ sheet.Description() } { attrs... } data-gsxui-slot-sheet-description data-gsxui-slot-dialog-description>
		{ children }
	</p>
}

component SheetClose(children gsx.Node, attrs gsx.Attrs) {
	<button data-gsxui-dialog-close type="button" { attrs... } data-gsxui-slot-sheet-close>{ children }</button>
}

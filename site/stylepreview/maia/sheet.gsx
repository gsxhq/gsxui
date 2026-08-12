package maia

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
		class={
			"backdrop:bg-black/80 supports-backdrop-filter:backdrop:backdrop-blur-xs bg-popover text-popover-foreground fixed z-50 flex flex-col bg-clip-padding text-sm shadow-lg transition duration-200 ease-in-out backdrop:transition-none m-0 gap-4 max-h-none outline-none data-[state=closed]:animate-out data-[state=open]:animate-in backdrop:backdrop-blur-xs backdrop:duration-200 data-[state=open]:backdrop:animate-in data-[state=open]:backdrop:fade-in-0 data-[state=closed]:backdrop:animate-out data-[state=closed]:backdrop:fade-out-0 open:flex",
			switch side {
			case "bottom":
				"inset-x-0 bottom-0 h-auto border-t top-auto w-full max-w-none data-[state=closed]:slide-out-to-bottom data-[state=open]:slide-in-from-bottom"
			case "left":
				"inset-y-0 left-0 h-full w-3/4 border-r sm:max-w-sm right-auto data-[state=closed]:slide-out-to-left data-[state=open]:slide-in-from-left"
			case "top":
				"inset-x-0 top-0 h-auto border-b bottom-auto w-full max-w-none data-[state=closed]:slide-out-to-top data-[state=open]:slide-in-from-top"
			default:
				"inset-y-0 right-0 h-full w-3/4 border-l sm:max-w-sm left-auto data-[state=closed]:slide-out-to-right data-[state=open]:slide-in-from-right"
			}
		}
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
				class={ "absolute top-4 end-4" }
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
					class={ "size-4" }
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
	<div class={ "gap-1.5 p-6 flex flex-col" } { attrs... } data-gsxui-slot-sheet-header>{ children }</div>
}

component SheetFooter(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "gap-2 p-6 flex flex-col mt-auto" } { attrs... } data-gsxui-slot-sheet-footer>{ children }</div>
}

component SheetTitle(children gsx.Node, attrs gsx.Attrs) {
	<h2
		class={ "text-foreground text-base font-medium" }
		{ attrs... }
		data-gsxui-slot-sheet-title
		data-gsxui-slot-dialog-title
	>
		{ children }
	</h2>
}

component SheetDescription(children gsx.Node, attrs gsx.Attrs) {
	<p
		class={ "text-muted-foreground text-sm" }
		{ attrs... }
		data-gsxui-slot-sheet-description
		data-gsxui-slot-dialog-description
	>
		{ children }
	</p>
}

component SheetClose(children gsx.Node, attrs gsx.Attrs) {
	<button data-gsxui-dialog-close type="button" { attrs... } data-gsxui-slot-sheet-close>{ children }</button>
}

package maia

import "github.com/gsxhq/gsx"

// Sheet composes Dialog's root and behavior with side-anchored presentation.
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
		data-gsxui-dialog-trigger
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
			"m-0 flex-col gap-4 shadow-lg transition ease-in-out max-h-none bg-background text-foreground fixed z-50 text-sm duration-200 outline-none data-[state=closed]:animate-out data-[state=open]:animate-in backdrop:bg-[var(--overlay)] backdrop:backdrop-blur-xs backdrop:duration-200 data-[state=open]:backdrop:animate-in data-[state=open]:backdrop:fade-in-0 data-[state=closed]:backdrop:animate-out data-[state=closed]:backdrop:fade-out-0 open:flex",
			switch side {
			case "bottom":
				"inset-x-0 bottom-0 top-auto h-auto w-full max-w-none border-t data-[state=closed]:slide-out-to-bottom data-[state=open]:slide-in-from-bottom"
			case "left":
				"inset-y-0 left-0 right-auto h-full w-3/4 border-r sm:max-w-sm data-[state=closed]:slide-out-to-left data-[state=open]:slide-in-from-left"
			case "top":
				"inset-x-0 top-0 bottom-auto h-auto w-full max-w-none border-b data-[state=closed]:slide-out-to-top data-[state=open]:slide-in-from-top"
			default:
				"inset-y-0 right-0 left-auto h-full w-3/4 border-l sm:max-w-sm data-[state=closed]:slide-out-to-right data-[state=open]:slide-in-from-right"
			}
		}
		data-gsxui-dialog-content
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
				class={
					"absolute top-3 right-3 rounded-xs opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:ring-2 focus:ring-ring focus:ring-offset-2 focus:outline-hidden disabled:pointer-events-none"
				}
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
	<div class={ "flex flex-col gap-0.5 p-4" } { attrs... } data-gsxui-slot-sheet-header>{ children }</div>
}

component SheetFooter(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "mt-auto flex flex-col gap-2 p-4" } { attrs... } data-gsxui-slot-sheet-footer>{ children }</div>
}

component SheetTitle(children gsx.Node, attrs gsx.Attrs) {
	<h2
		class={ "text-base font-medium text-foreground" }
		data-gsxui-dialog-title
		{ attrs... }
		data-gsxui-slot-sheet-title
	>
		{ children }
	</h2>
}

component SheetDescription(children gsx.Node, attrs gsx.Attrs) {
	<p
		class={ "text-sm text-muted-foreground" }
		data-gsxui-dialog-description
		{ attrs... }
		data-gsxui-slot-sheet-description
	>
		{ children }
	</p>
}

component SheetClose(children gsx.Node, attrs gsx.Attrs) {
	<button data-gsxui-dialog-close type="button" { attrs... } data-gsxui-slot-sheet-close>{ children }</button>
}
